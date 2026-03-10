package ui

import (
	"errors"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type GitManager interface {
	RunCommand(args ...string) (stdout string, stderr string, err error)
	GetStagedDiff() (string, error)
	GetGithubUsername() (string, error)
	GetRepoName() (string, error)
	CommitAndPush(message string, noVerify bool) error
}

type AIService interface {
	GenerateCommitMessage(diff, username, repoName, customInstructions string) (string, error)
}

type commitState int

const (
	commitStateIdle commitState = iota
	commitStateStaging
	commitStateLoadingAI
	commitStateShowingProposal
	commitStateEditing
	commitStateCommitting
	commitStateDone
)

type CommitModel struct {
	state       commitState
	aiMessage   string
	textInput   textinput.Model
	spinner     spinner.Model
	gitDiff     string
	errorMsg    string
	successMsg  string
	cancelled   bool
	noVerify    bool
	customNotes string
	username    string
	repoName    string
	menuIndex   int

	git GitManager
	ai  AIService
}

type stagedMsg struct {
	Diff     string
	Username string
	RepoName string
	Err      error
}

type aiGeneratedMsg struct {
	Message string
	Err     error
}

type commitFinishedMsg struct {
	Err error
}

func NewCommitModel(git GitManager, ai AIService, noVerify bool, customInstructions string) CommitModel {
	ti := textinput.New()
	ti.Placeholder = "Write commit message"
	ti.Prompt = "> "
	ti.CharLimit = 300
	ti.Width = min(80, max(30, BoxInnerWidth()-4))

	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(GetCurrentTheme().Primary)

	return CommitModel{
		state:       commitStateStaging,
		textInput:   ti,
		spinner:     sp,
		noVerify:    noVerify,
		customNotes: customInstructions,
		git:         git,
		ai:          ai,
	}
}

func (m CommitModel) Init() tea.Cmd {
	if m.git == nil || m.ai == nil {
		return tea.Quit
	}

	return tea.Batch(m.spinner.Tick, m.stageCmd())
}

func (m CommitModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Keep the input usable across narrow terminals.
		boxW := BoxWidthForTerminalWidth(msg.Width)
		inner := boxW - 6
		m.textInput.Width = min(80, max(30, inner-4))
		return m, nil
	case tea.KeyMsg:
		return m.handleKey(msg)
	case spinner.TickMsg:
		if m.state == commitStateLoadingAI || m.state == commitStateCommitting || m.state == commitStateStaging {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
		return m, nil
	case stagedMsg:
		if msg.Err != nil {
			m.state = commitStateDone
			m.errorMsg = msg.Err.Error()
			return m, nil
		}

		m.gitDiff = msg.Diff
		m.username = msg.Username
		m.repoName = msg.RepoName
		m.state = commitStateLoadingAI
		m.errorMsg = ""
		return m, tea.Batch(m.spinner.Tick, m.generateAICmd())
	case aiGeneratedMsg:
		if msg.Err != nil {
			m.state = commitStateDone
			m.errorMsg = msg.Err.Error()
			return m, nil
		}
		m.aiMessage = strings.TrimSpace(msg.Message)
		m.state = commitStateShowingProposal
		m.menuIndex = 0
		m.errorMsg = ""
		return m, nil
	case commitFinishedMsg:
		m.state = commitStateDone
		if msg.Err != nil {
			m.errorMsg = msg.Err.Error()
			return m, tea.Quit
		}
		m.successMsg = "commit and push completed"
		return m, tea.Quit
	}

	if m.state == commitStateEditing {
		var cmd tea.Cmd
		m.textInput, cmd = m.textInput.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m CommitModel) View() string {
	st := GetStyles()

	var content string
	switch m.state {
	case commitStateIdle:
		content = st.Muted.Render("Waiting to start AI commit workflow...")
	case commitStateStaging:
		content = fmt.Sprintf("%s %s", m.spinner.View(), st.Text.Render("Staging files..."))
	case commitStateLoadingAI:
		content = fmt.Sprintf("%s %s", m.spinner.View(), st.Text.Render("Generating AI commit message..."))
	case commitStateShowingProposal:
		proposal := RenderInnerBox(BoxPrimary, "AI Commit Proposal", wrapParagraphs(m.aiMessage, BoxInnerWidth()-6))
		menu := strings.Join([]string{
			renderCommitMenuItem(m.menuIndex, 0, "Commit and push"),
			renderCommitMenuItem(m.menuIndex, 1, "Regenerate message"),
			renderCommitMenuItem(m.menuIndex, 2, "Edit message"),
			renderCommitMenuItem(m.menuIndex, 3, "Cancel (or q)"),
		}, "\n")
		content = strings.Join([]string{
			proposal,
			"",
			menu,
			"",
			st.Hint.Render("Use ↑/↓ and Enter. Press q to cancel."),
		}, "\n")
	case commitStateEditing:
		content = strings.Join([]string{
			RenderTitle("Edit Commit Message"),
			"",
			m.textInput.View(),
			"",
			st.Error.Render(m.errorMsg),
			"",
			st.Hint.Render("Enter to commit • Esc to go back • Ctrl+C to quit"),
		}, "\n")
	case commitStateCommitting:
		content = fmt.Sprintf("%s %s", m.spinner.View(), st.Text.Render("Committing and pushing changes..."))
	case commitStateDone:
		switch {
		case m.errorMsg != "":
			return MessageBox(BoxError, "Commit Failed", m.errorMsg)
		case m.cancelled:
			return MessageBox(BoxWarning, "Cancelled", "Workflow cancelled.")
		default:
			return MessageBox(BoxSuccess, "Success", m.successMsg)
		}
	}

	return RenderBox(BoxPrimary, "AI Commit", content)
}

func (m CommitModel) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	}

	switch m.state {
	case commitStateShowingProposal:
		switch msg.String() {
		case "up":
			if m.menuIndex > 0 {
				m.menuIndex--
			}
			return m, nil
		case "down":
			if m.menuIndex < 3 {
				m.menuIndex++
			}
			return m, nil
		case "enter":
			return m.handleCommitMenuSelect(m.menuIndex)
		case "q":
			m.state = commitStateDone
			m.cancelled = true
			return m, tea.Quit
		}
	case commitStateEditing:
		switch msg.String() {
		case "esc":
			m.state = commitStateShowingProposal
			m.textInput.Blur()
			m.errorMsg = ""
			return m, nil
		case "enter":
			message := strings.TrimSpace(m.textInput.Value())
			if message == "" {
				m.errorMsg = "commit message cannot be empty"
				return m, nil
			}
			m.aiMessage = message
			m.textInput.Blur()
			m.state = commitStateCommitting
			m.errorMsg = ""
			return m, tea.Batch(m.spinner.Tick, m.commitCmd(message))
		default:
			var cmd tea.Cmd
			m.textInput, cmd = m.textInput.Update(msg)
			return m, cmd
		}
	case commitStateDone:
		return m, tea.Quit
	}

	return m, nil
}

func renderCommitMenuItem(selected, idx int, label string) string {
	st := GetStyles()
	prefix := "  "
	if selected == idx {
		prefix = st.Accent.Render("> ")
	}
	return prefix + st.Hint.Render(label)
}

func (m CommitModel) handleCommitMenuSelect(idx int) (tea.Model, tea.Cmd) {
	switch idx {
	case 0:
		m.state = commitStateCommitting
		m.errorMsg = ""
		return m, tea.Batch(m.spinner.Tick, m.commitCmd(m.aiMessage))
	case 1:
		m.state = commitStateLoadingAI
		m.errorMsg = ""
		return m, tea.Batch(m.spinner.Tick, m.generateAICmd())
	case 2:
		m.state = commitStateEditing
		m.textInput.SetValue(m.aiMessage)
		m.textInput.Focus()
		m.errorMsg = ""
		return m, textinput.Blink
	case 3:
		m.state = commitStateDone
		m.cancelled = true
		return m, tea.Quit
	default:
		return m, nil
	}
}

func (m CommitModel) stageCmd() tea.Cmd {
	return func() tea.Msg {
		if m.git == nil {
			return stagedMsg{Err: errors.New("git service is not configured")}
		}

		if _, stderr, err := m.git.RunCommand("add", "."); err != nil {
			if stderr != "" {
				return stagedMsg{Err: fmt.Errorf("git add failed: %s", stderr)}
			}
			return stagedMsg{Err: fmt.Errorf("git add failed: %w", err)}
		}

		diff, err := m.git.GetStagedDiff()
		if err != nil {
			return stagedMsg{Err: fmt.Errorf("failed to get staged diff: %w", err)}
		}

		username, err := m.git.GetGithubUsername()
		if err != nil {
			return stagedMsg{Err: fmt.Errorf("failed to get git username: %w", err)}
		}
		repoName, err := m.git.GetRepoName()
		if err != nil {
			return stagedMsg{Err: fmt.Errorf("failed to get repository name: %w", err)}
		}

		return stagedMsg{
			Diff:     diff,
			Username: username,
			RepoName: repoName,
		}
	}
}

func (m CommitModel) generateAICmd() tea.Cmd {
	return func() tea.Msg {
		if m.ai == nil {
			return aiGeneratedMsg{Err: errors.New("ai service is not configured")}
		}

		message, err := m.ai.GenerateCommitMessage(m.gitDiff, m.username, m.repoName, m.customNotes)
		if err != nil {
			return aiGeneratedMsg{Err: fmt.Errorf("failed to generate AI commit message: %w", err)}
		}

		return aiGeneratedMsg{Message: message}
	}
}

func (m CommitModel) commitCmd(message string) tea.Cmd {
	return func() tea.Msg {
		if m.git == nil {
			return commitFinishedMsg{Err: errors.New("git service is not configured")}
		}

		err := m.git.CommitAndPush(message, m.noVerify)
		if err != nil {
			return commitFinishedMsg{Err: fmt.Errorf("commit and push failed: %w", err)}
		}

		return commitFinishedMsg{}
	}
}
