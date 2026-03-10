package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

type BoxVariant int

const (
	BoxPrimary BoxVariant = iota
	BoxSuccess
	BoxWarning
	BoxError
)

type Styles struct {
	Title   lipgloss.Style
	Muted   lipgloss.Style
	Text    lipgloss.Style
	Hint    lipgloss.Style
	Success lipgloss.Style
	Warning lipgloss.Style
	Error   lipgloss.Style
	Accent  lipgloss.Style
}

type spinnerDoneMsg struct{ err error }

func GetStyles() Styles {
	t := GetCurrentTheme()
	return Styles{
		Title:   lipgloss.NewStyle().Bold(true).Foreground(t.Primary),
		Muted:   lipgloss.NewStyle().Foreground(t.Muted),
		Text:    lipgloss.NewStyle().Foreground(t.Text),
		Hint:    lipgloss.NewStyle().Foreground(t.Muted),
		Success: lipgloss.NewStyle().Bold(true).Foreground(t.Success),
		Warning: lipgloss.NewStyle().Bold(true).Foreground(t.Warning),
		Error:   lipgloss.NewStyle().Bold(true).Foreground(t.Error),
		Accent:  lipgloss.NewStyle().Bold(true).Foreground(t.Primary),
	}
}

func TerminalWidth() int {
	fd := int(os.Stdout.Fd())
	w, _, err := term.GetSize(fd)
	if err != nil || w <= 0 {
		return 80
	}
	return w
}

func BoxWidthForTerminalWidth(terminalWidth int) int {
	w := terminalWidth
	// Leave some breathing room to avoid wrapping against terminal edges.
	w = w - 4
	if w < 50 {
		w = 50
	}
	if w > 100 {
		w = 100
	}
	return w
}

func boxWidth() int {
	return BoxWidthForTerminalWidth(TerminalWidth())
}

func BoxInnerWidth() int {
	// RenderBox uses border+padding that is roughly 6 columns total.
	return boxWidth() - 6
}

func RenderBox(variant BoxVariant, title string, content string) string {
	t := GetCurrentTheme()
	st := GetStyles()

	border := t.BorderPrimary
	switch variant {
	case BoxSuccess:
		border = t.BorderSuccess
	case BoxWarning:
		border = t.BorderWarning
	case BoxError:
		border = t.BorderError
	}

	w := boxWidth()
	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(border).
		Padding(1, 2).
		Width(w)

	if strings.TrimSpace(title) != "" {
		content = st.Title.Render(title) + "\n\n" + content
	}

	panel := box.Render(content)
	return lipgloss.PlaceHorizontal(TerminalWidth(), lipgloss.Center, panel)
}

func RenderInnerBox(variant BoxVariant, title string, content string) string {
	t := GetCurrentTheme()
	st := GetStyles()

	border := t.BorderPrimary
	switch variant {
	case BoxSuccess:
		border = t.BorderSuccess
	case BoxWarning:
		border = t.BorderWarning
	case BoxError:
		border = t.BorderError
	}

	w := BoxInnerWidth()
	if w < 30 {
		w = 30
	}

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(border).
		Padding(1, 2).
		Width(w)

	if strings.TrimSpace(title) != "" {
		content = st.Title.Render(title) + "\n\n" + content
	}

	return box.Render(content)
}

func MessageBox(variant BoxVariant, title, body string) string {
	st := GetStyles()
	body = wrapParagraphs(body, boxWidth()-6) // approximate: border+padding
	switch variant {
	case BoxSuccess:
		body = st.Success.Render(body)
	case BoxWarning:
		body = st.Warning.Render(body)
	case BoxError:
		body = st.Error.Render(body)
	default:
		body = st.Text.Render(body)
	}
	return RenderBox(variant, title, body)
}

func RenderTitle(text string) string {
	st := GetStyles()
	return st.Title.Render(text)
}

func WithSpinner(message string, fn func() error) error {
	if fn == nil {
		return fmt.Errorf("fn cannot be nil")
	}

	// Non-TTY fallback: just run the function.
	if !term.IsTerminal(int(os.Stdout.Fd())) {
		if strings.TrimSpace(message) != "" {
			fmt.Fprintln(os.Stdout, message)
		}
		return fn()
	}

	m := spinnerModel{
		message: message,
		spinner: spinner.New(spinner.WithSpinner(spinner.Dot)),
	}

	p := tea.NewProgram(m)
	resultCh := make(chan error, 1)
	go func() {
		err := fn()
		resultCh <- err
		p.Send(spinnerDoneMsg{err: err})
	}()

	finalModel, err := p.Run()
	if err != nil {
		// If Bubble Tea can't attach to TTY, still return the underlying result.
		return <-resultCh
	}

	if sm, ok := finalModel.(spinnerModel); ok {
		return sm.doneErr
	}
	return nil
}

type spinnerModel struct {
	message string
	spinner spinner.Model
	doneErr error
}

func (m spinnerModel) Init() tea.Cmd { return m.spinner.Tick }

func (m spinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case spinnerDoneMsg:
		m.doneErr = msg.err
		return m, tea.Quit
	default:
		return m, nil
	}
}

func (m spinnerModel) View() string {
	st := GetStyles()
	return st.Text.Render(m.spinner.View() + " " + strings.TrimSpace(m.message))
}

func wrapParagraphs(s string, width int) string {
	lines := strings.Split(s, "\n")
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimRight(line, " ")
		if strings.TrimSpace(line) == "" {
			out = append(out, "")
			continue
		}
		out = append(out, wrapLine(line, width)...)
	}
	return strings.Join(out, "\n")
}

func wrapLine(line string, width int) []string {
	if width <= 10 {
		return []string{line}
	}

	words := strings.Fields(line)
	if len(words) == 0 {
		return []string{""}
	}

	var out []string
	cur := words[0]
	for _, w := range words[1:] {
		if lipgloss.Width(cur)+1+lipgloss.Width(w) <= width {
			cur = cur + " " + w
			continue
		}
		out = append(out, cur)
		cur = w
	}
	out = append(out, cur)
	return out
}
