package ui

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
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
	scopeIndex int
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

type historyDoneTimeoutMsg struct{}

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
			return m, tea.Tick(300*time.Millisecond, func(time.Time) tea.Msg { return historyDoneTimeoutMsg{} })
		}
		m.successMsg = "history updated successfully"
		return m, tea.Tick(300*time.Millisecond, func(time.Time) tea.Msg { return historyDoneTimeoutMsg{} })
	case historyDoneTimeoutMsg:
		return m, tea.Quit
	}

	return m, nil
}

func (m HistoryModel) View() string {
	st := GetStyles()

	title := ""
	var content string
	switch m.state {
	case historyStateScopeSelect:
		title = historyTitle(m.operation)
		content = strings.Join([]string{
			st.Muted.Render("Scope"),
			"",
			renderScopeOption(m.scopeIndex, 0, "Last commit"),
			renderScopeOption(m.scopeIndex, 1, "Specific hash"),
			renderScopeOption(m.scopeIndex, 2, "All commits"),
			"",
			st.Warning.Render("Warning: rewriting history can be destructive, especially for all commits."),
			"",
			st.Hint.Render("Use arrows or j/k and Enter. Shortcuts: 1/2/3. Press q to cancel."),
		}, "\n")
	case historyStateHashInput:
		title = historyTitle(m.operation)
		content = strings.Join([]string{
			st.Text.Render("Enter target commit hash:"),
			m.inputs[0].View(),
			"",
			st.Error.Render(m.errorMsg),
			st.Hint.Render("Enter to continue • Esc to go back • Ctrl+C to quit"),
		}, "\n")
	case historyStateTimeInput:
		title = "Change Commit Time"
		content = strings.Join([]string{
			st.Text.Render("Date (YYYY-MM-DD):"),
			m.inputs[0].View(),
			"",
			st.Text.Render("Time (HH:MM or HH:MM:SS):"),
			m.inputs[1].View(),
			"",
			st.Error.Render(m.errorMsg),
			st.Hint.Render("Tab switch • Enter next/confirm • Esc back • Ctrl+C quit"),
		}, "\n")
	case historyStateAuthorInput:
		title = "Change Commit Author"
		content = strings.Join([]string{
			st.Text.Render("Author name:"),
			m.inputs[0].View(),
			"",
			st.Text.Render("Author email:"),
			m.inputs[1].View(),
			"",
			st.Error.Render(m.errorMsg),
			st.Hint.Render("Tab switch • Enter next/confirm • Esc back • Ctrl+C quit"),
		}, "\n")
	case historyStateMessageInput:
		title = "Change Commit Message"
		content = strings.Join([]string{
			st.Text.Render("New commit message:"),
			m.inputs[0].View(),
			"",
			st.Error.Render(m.errorMsg),
			st.Hint.Render("Enter to continue • Esc to cancel • Ctrl+C to quit"),
		}, "\n")
	case historyStateConfirm:
		title = "Confirm History Rewrite"
		content = strings.Join([]string{
			st.Text.Render(confirmSummary(m)),
			"",
			st.Hint.Render("Proceed? [y] Yes / [n] No (q to cancel)"),
		}, "\n")
	case historyStateRunning:
		title = historyTitle(m.operation)
		content = fmt.Sprintf("%s %s", m.spinner.View(), st.Text.Render("Applying history changes..."))
	case historyStateDone:
		switch {
		case m.errorMsg != "":
			return MessageBox(BoxError, "History Update Failed", m.errorMsg) + "\n"
		case m.cancelled:
			return MessageBox(BoxWarning, "Cancelled", "History update cancelled.") + "\n"
		default:
			return MessageBox(BoxSuccess, "Success", m.successMsg) + "\n"
		}
	}

	return RenderBox(BoxPrimary, title, content)
}

func (m HistoryModel) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if msg.String() == "ctrl+c" {
		return m, tea.Quit
	}

	switch m.state {
	case historyStateScopeSelect:
		switch msg.String() {
		case "up", "k":
			if m.scopeIndex > 0 {
				m.scopeIndex--
			}
			return m, nil
		case "down", "j":
			if m.scopeIndex < 2 {
				m.scopeIndex++
			}
			return m, nil
		case "enter":
			switch m.scopeIndex {
			case 0:
				m.scope = "single"
				return m.toFieldInput(), nil
			case 1:
				m.scope = "specific"
				m.state = historyStateHashInput
				m.inputs = []textinput.Model{newInput("Commit hash", 64)}
				m.inputs[0].Focus()
				m.active = 0
				m.errorMsg = ""
				return m, textinput.Blink
			case 2:
				m.scope = "all"
				return m.toFieldInput(), nil
			}
			return m, nil
		case "1":
			m.scope = "single"
			m.scopeIndex = 0
			return m.toFieldInput(), nil
		case "2":
			m.scope = "specific"
			m.scopeIndex = 1
			m.state = historyStateHashInput
			m.inputs = []textinput.Model{newInput("Commit hash", 64)}
			m.inputs[0].Focus()
			m.active = 0
			m.errorMsg = ""
			return m, textinput.Blink
		case "3":
			m.scope = "all"
			m.scopeIndex = 2
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
		return m, nil
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

func renderScopeOption(selectedIdx, idx int, label string) string {
	st := GetStyles()
	prefix := "  "
	if selectedIdx == idx {
		prefix = st.Accent.Render("> ")
	}
	return prefix + st.Text.Render(label)
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
