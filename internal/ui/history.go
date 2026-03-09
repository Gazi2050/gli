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

type HistoryService interface {
	ChangeCommitTime(scope, targetHash, dateTime string) error
	ChangeCommitAuthor(scope, targetHash, name, email string) error
	ChangeCommitMessage(newMessage string) error
}

type HistoryOperation string

const (
	HistoryOperationTime    HistoryOperation = "time"
	HistoryOperationAuthor  HistoryOperation = "author"
	HistoryOperationMessage HistoryOperation = "message"
)

type historyState int

const (
	historyStateScopeSelect historyState = iota
	historyStateHashInput
	historyStateTimeInput
	historyStateAuthorInput
	historyStateMessageInput
	historyStateConfirm
	historyStateRunning
	historyStateDone
)

type HistoryModel struct {
	state      historyState
	operation  HistoryOperation
	scope      string
	targetHash string
	inputs     []textinput.Model
	active     int
	errorMsg   string
	successMsg string
	cancelled  bool

	spinner spinner.Model
	history HistoryService
}

type historyFinishedMsg struct {
	Err error
}

func NewHistoryModel(history HistoryService, operation HistoryOperation) HistoryModel {
	sp := spinner.New()
	sp.Spinner = spinner.Dot

	model := HistoryModel{
		history:   history,
		operation: operation,
		spinner:   sp,
	}

	switch operation {
	case HistoryOperationTime:
		model.state = historyStateScopeSelect
	case HistoryOperationAuthor:
		model.state = historyStateScopeSelect
	case HistoryOperationMessage:
		model.scope = "single"
		model.state = historyStateMessageInput
		model.inputs = []textinput.Model{newInput("New commit message", 120)}
		model.inputs[0].Focus()
	default:
		model.state = historyStateDone
		model.errorMsg = "invalid history operation"
	}

	return model
}

func (m HistoryModel) Init() tea.Cmd {
	if m.history == nil {
		return tea.Quit
	}
	if m.state == historyStateDone {
		return tea.Quit
	}
	return nil
}

func (m HistoryModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKey(msg)
	case spinner.TickMsg:
		if m.state == historyStateRunning {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
	case historyFinishedMsg:
		m.state = historyStateDone
		if msg.Err != nil {
			m.errorMsg = msg.Err.Error()
			return m, nil
		}
		m.successMsg = "history updated successfully"
		return m, nil
	}

	return m, nil
}

func (m HistoryModel) View() string {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("42"))
	errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	successStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Bold(true)
	hintStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	boxStyle := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("240")).Padding(1, 2)

	var content string
	switch m.state {
	case historyStateScopeSelect:
		content = strings.Join([]string{
			titleStyle.Render(historyTitle(m.operation)),
			"",
			"[1] Last commit",
			"[2] Specific hash",
			"[3] All commits",
			"",
			hintStyle.Render("Choose 1/2/3, or q to cancel"),
		}, "\n")
	case historyStateHashInput:
		content = strings.Join([]string{
			titleStyle.Render(historyTitle(m.operation)),
			"",
			"Enter target commit hash:",
			m.inputs[0].View(),
			"",
			errorStyle.Render(m.errorMsg),
			hintStyle.Render("Enter to continue, Esc to go back"),
		}, "\n")
	case historyStateTimeInput:
		content = strings.Join([]string{
			titleStyle.Render("Change Commit Time"),
			"",
			"Date (YYYY-MM-DD):",
			m.inputs[0].View(),
			"",
			"Time (HH:MM or HH:MM:SS):",
			m.inputs[1].View(),
			"",
			errorStyle.Render(m.errorMsg),
			hintStyle.Render("Enter next field, Tab switch, Esc back"),
		}, "\n")
	case historyStateAuthorInput:
		content = strings.Join([]string{
			titleStyle.Render("Change Commit Author"),
			"",
			"Author name:",
			m.inputs[0].View(),
			"",
			"Author email:",
			m.inputs[1].View(),
			"",
			errorStyle.Render(m.errorMsg),
			hintStyle.Render("Enter next field, Tab switch, Esc back"),
		}, "\n")
	case historyStateMessageInput:
		content = strings.Join([]string{
			titleStyle.Render("Change Commit Message"),
			"",
			"New commit message:",
			m.inputs[0].View(),
			"",
			errorStyle.Render(m.errorMsg),
			hintStyle.Render("Enter to continue, Esc to cancel"),
		}, "\n")
	case historyStateConfirm:
		content = strings.Join([]string{
			titleStyle.Render("Confirm History Rewrite"),
			"",
			confirmSummary(m),
			"",
			hintStyle.Render("Proceed? [y] Yes / [n] No (q to cancel)"),
		}, "\n")
	case historyStateRunning:
		content = fmt.Sprintf("%s Applying history changes...", m.spinner.View())
	case historyStateDone:
		switch {
		case m.errorMsg != "":
			content = errorStyle.Render("Error: "+m.errorMsg) + "\n" + hintStyle.Render("Press q to exit")
		case m.cancelled:
			content = hintStyle.Render("History update cancelled. Press q to exit")
		default:
			content = successStyle.Render(m.successMsg) + "\n" + hintStyle.Render("Press q to exit")
		}
	}

	return boxStyle.Render(content)
}

func (m HistoryModel) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if msg.String() == "ctrl+c" {
		return m, tea.Quit
	}

	switch m.state {
	case historyStateScopeSelect:
		switch msg.String() {
		case "1":
			m.scope = "single"
			return m.toFieldInput(), nil
		case "2":
			m.scope = "specific"
			m.state = historyStateHashInput
			m.inputs = []textinput.Model{newInput("Commit hash", 64)}
			m.inputs[0].Focus()
			m.active = 0
			m.errorMsg = ""
			return m, textinput.Blink
		case "3":
			m.scope = "all"
			return m.toFieldInput(), nil
		case "q":
			m.state = historyStateDone
			m.cancelled = true
			return m, tea.Quit
		}
	case historyStateHashInput:
		switch msg.String() {
		case "esc":
			m.state = historyStateScopeSelect
			m.inputs = nil
			m.active = 0
			m.errorMsg = ""
			return m, nil
		case "enter":
			hash := strings.TrimSpace(m.inputs[0].Value())
			if hash == "" {
				m.errorMsg = "commit hash is required"
				return m, nil
			}
			m.targetHash = hash
			m.errorMsg = ""
			return m.toFieldInput(), nil
		default:
			var cmd tea.Cmd
			m.inputs[0], cmd = m.inputs[0].Update(msg)
			return m, cmd
		}
	case historyStateTimeInput, historyStateAuthorInput, historyStateMessageInput:
		switch msg.String() {
		case "tab":
			if len(m.inputs) > 1 {
				m.focusInput((m.active + 1) % len(m.inputs))
				return m, textinput.Blink
			}
			return m, nil
		case "esc":
			if m.operation == HistoryOperationMessage {
				m.state = historyStateDone
				m.cancelled = true
				return m, tea.Quit
			}
			if m.scope == "specific" {
				m.state = historyStateHashInput
				m.inputs = []textinput.Model{newInput("Commit hash", 64)}
				m.inputs[0].SetValue(m.targetHash)
				m.inputs[0].Focus()
				m.active = 0
				m.errorMsg = ""
				return m, textinput.Blink
			}
			m.state = historyStateScopeSelect
			m.inputs = nil
			m.active = 0
			m.errorMsg = ""
			return m, nil
		case "enter":
			if m.active < len(m.inputs)-1 {
				m.focusInput(m.active + 1)
				return m, textinput.Blink
			}
			if err := m.validateInputs(); err != nil {
				m.errorMsg = err.Error()
				return m, nil
			}
			m.errorMsg = ""
			m.state = historyStateConfirm
			return m, nil
		default:
			if len(m.inputs) == 0 {
				return m, nil
			}
			var cmd tea.Cmd
			m.inputs[m.active], cmd = m.inputs[m.active].Update(msg)
			return m, cmd
		}
	case historyStateConfirm:
		switch msg.String() {
		case "y":
			m.state = historyStateRunning
			m.errorMsg = ""
			return m, tea.Batch(m.spinner.Tick, m.runHistoryCmd())
		case "n":
			return m.toFieldInput(), nil
		case "q":
			m.state = historyStateDone
			m.cancelled = true
			return m, tea.Quit
		}
	case historyStateDone:
		if msg.String() == "q" || msg.String() == "enter" {
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m HistoryModel) toFieldInput() HistoryModel {
	m.inputs = nil
	m.active = 0
	m.errorMsg = ""

	switch m.operation {
	case HistoryOperationTime:
		m.state = historyStateTimeInput
		m.inputs = []textinput.Model{
			newInput("YYYY-MM-DD", 10),
			newInput("HH:MM", 8),
		}
	case HistoryOperationAuthor:
		m.state = historyStateAuthorInput
		m.inputs = []textinput.Model{
			newInput("Author name", 80),
			newInput("Author email", 120),
		}
	case HistoryOperationMessage:
		m.state = historyStateMessageInput
		m.inputs = []textinput.Model{newInput("New commit message", 120)}
	default:
		m.state = historyStateDone
		m.errorMsg = "invalid history operation"
		return m
	}

	if len(m.inputs) > 0 {
		m.focusInput(0)
	}
	return m
}

func (m *HistoryModel) focusInput(idx int) {
	if idx < 0 || idx >= len(m.inputs) {
		return
	}
	m.active = idx
	for i := range m.inputs {
		if i == idx {
			m.inputs[i].Focus()
		} else {
			m.inputs[i].Blur()
		}
	}
}

func (m HistoryModel) validateInputs() error {
	switch m.operation {
	case HistoryOperationTime:
		date := strings.TrimSpace(m.inputs[0].Value())
		timeValue := strings.TrimSpace(m.inputs[1].Value())
		if date == "" {
			return errors.New("date is required")
		}
		if timeValue == "" {
			return errors.New("time is required")
		}
	case HistoryOperationAuthor:
		name := strings.TrimSpace(m.inputs[0].Value())
		email := strings.TrimSpace(m.inputs[1].Value())
		if name == "" {
			return errors.New("author name is required")
		}
		if email == "" {
			return errors.New("author email is required")
		}
	case HistoryOperationMessage:
		message := strings.TrimSpace(m.inputs[0].Value())
		if message == "" {
			return errors.New("new commit message is required")
		}
	default:
		return errors.New("invalid history operation")
	}

	return nil
}

func (m HistoryModel) runHistoryCmd() tea.Cmd {
	return func() tea.Msg {
		if m.history == nil {
			return historyFinishedMsg{Err: errors.New("history service is not configured")}
		}

		var err error
		switch m.operation {
		case HistoryOperationTime:
			date := strings.TrimSpace(m.inputs[0].Value())
			timeValue := strings.TrimSpace(m.inputs[1].Value())
			if strings.Count(timeValue, ":") == 1 {
				timeValue += ":00"
			}
			dt := date + " " + timeValue
			err = m.history.ChangeCommitTime(m.scope, m.targetHash, dt)
		case HistoryOperationAuthor:
			name := strings.TrimSpace(m.inputs[0].Value())
			email := strings.TrimSpace(m.inputs[1].Value())
			err = m.history.ChangeCommitAuthor(m.scope, m.targetHash, name, email)
		case HistoryOperationMessage:
			message := strings.TrimSpace(m.inputs[0].Value())
			err = m.history.ChangeCommitMessage(message)
		default:
			err = errors.New("invalid history operation")
		}

		if err != nil {
			return historyFinishedMsg{Err: fmt.Errorf("history update failed: %w", err)}
		}
		return historyFinishedMsg{}
	}
}

func historyTitle(op HistoryOperation) string {
	switch op {
	case HistoryOperationTime:
		return "Change Commit Time"
	case HistoryOperationAuthor:
		return "Change Commit Author"
	case HistoryOperationMessage:
		return "Change Commit Message"
	default:
		return "History"
	}
}

func confirmSummary(m HistoryModel) string {
	lines := []string{
		fmt.Sprintf("Operation: %s", string(m.operation)),
		fmt.Sprintf("Scope: %s", m.scope),
	}
	if m.scope == "specific" {
		lines = append(lines, fmt.Sprintf("Target hash: %s", m.targetHash))
	}

	switch m.operation {
	case HistoryOperationTime:
		date := strings.TrimSpace(m.inputs[0].Value())
		timeValue := strings.TrimSpace(m.inputs[1].Value())
		if strings.Count(timeValue, ":") == 1 {
			timeValue += ":00"
		}
		lines = append(lines, fmt.Sprintf("New datetime: %s %s", date, timeValue))
	case HistoryOperationAuthor:
		name := strings.TrimSpace(m.inputs[0].Value())
		email := strings.TrimSpace(m.inputs[1].Value())
		lines = append(lines, fmt.Sprintf("New author: %s <%s>", name, email))
	case HistoryOperationMessage:
		message := strings.TrimSpace(m.inputs[0].Value())
		lines = append(lines, fmt.Sprintf("New message: %s", message))
	}

	return strings.Join(lines, "\n")
}

func newInput(placeholder string, charLimit int) textinput.Model {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.Prompt = "> "
	ti.CharLimit = charLimit
	ti.Width = 80
	return ti
}
