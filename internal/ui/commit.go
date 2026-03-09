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
	ti.Width = 80

	sp := spinner.New()
	sp.Spinner = spinner.Dot

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
		m.errorMsg = ""
		return m, nil
	case commitFinishedMsg:
		m.state = commitStateDone
		if msg.Err != nil {
			m.errorMsg = msg.Err.Error()
			return m, nil
		}
		m.successMsg = "commit and push completed"
		return m, nil
	}

	if m.state == commitStateEditing {
		var cmd tea.Cmd
		m.textInput, cmd = m.textInput.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m CommitModel) View() string {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("42"))
	errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	successStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Bold(true)
	hintStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	boxStyle := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("240")).Padding(1, 2)

	var content string
	switch m.state {
	case commitStateIdle:
		content = "Waiting to start AI commit workflow..."
	case commitStateStaging:
		content = fmt.Sprintf("%s Staging files and loading git context...", m.spinner.View())
	case commitStateLoadingAI:
		content = fmt.Sprintf("%s Generating AI commit message...", m.spinner.View())
	case commitStateShowingProposal:
		content = strings.Join([]string{
			titleStyle.Render("AI Commit Proposal"),
			"",
			m.aiMessage,
			"",
			hintStyle.Render("[1] Commit and push"),
			hintStyle.Render("[2] Regenerate message"),
			hintStyle.Render("[3] Edit message"),
			hintStyle.Render("[4] Cancel (or q)"),
		}, "\n")
	case commitStateEditing:
		content = strings.Join([]string{
			titleStyle.Render("Edit Commit Message"),
			"",
			m.textInput.View(),
			"",
			errorStyle.Render(m.errorMsg),
			"",
			hintStyle.Render("Enter to commit, Esc to go back"),
		}, "\n")
	case commitStateCommitting:
		content = fmt.Sprintf("%s Committing and pushing changes...", m.spinner.View())
	case commitStateDone:
		switch {
		case m.errorMsg != "":
			content = errorStyle.Render("Error: "+m.errorMsg) + "\n" + hintStyle.Render("Press q to exit")
		case m.cancelled:
			content = hintStyle.Render("Workflow cancelled. Press q to exit")
		default:
			content = successStyle.Render(m.successMsg) + "\n" + hintStyle.Render("Press q to exit")
		}
	}

	return boxStyle.Render(content)
}

func (m CommitModel) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	}

	switch m.state {
	case commitStateShowingProposal:
		switch msg.String() {
		case "1":
			m.state = commitStateCommitting
			m.errorMsg = ""
			return m, tea.Batch(m.spinner.Tick, m.commitCmd(m.aiMessage))
		case "2":
			m.state = commitStateLoadingAI
			m.errorMsg = ""
			return m, tea.Batch(m.spinner.Tick, m.generateAICmd())
		case "3":
			m.state = commitStateEditing
			m.textInput.SetValue(m.aiMessage)
			m.textInput.Focus()
			m.errorMsg = ""
			return m, textinput.Blink
		case "4", "q":
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
		if msg.String() == "q" || msg.String() == "enter" {
			return m, tea.Quit
		}
	}

	return m, nil
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
