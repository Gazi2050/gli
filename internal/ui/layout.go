package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
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
	Title     lipgloss.Style
	Muted     lipgloss.Style
	Text      lipgloss.Style
	Hint      lipgloss.Style
	Success   lipgloss.Style
	Warning   lipgloss.Style
	Error     lipgloss.Style
	Accent    lipgloss.Style
	Secondary lipgloss.Style
	Dim       lipgloss.Style
}

type spinnerDoneMsg struct{ err error }

func GetStyles() Styles {
	t := GetCurrentTheme()
	return Styles{
		Title:     lipgloss.NewStyle().Bold(true).Foreground(t.Primary),
		Muted:     lipgloss.NewStyle().Foreground(t.Muted),
		Text:      lipgloss.NewStyle().Foreground(t.Text),
		Hint:      lipgloss.NewStyle().Foreground(t.Muted),
		Success:   lipgloss.NewStyle().Bold(true).Foreground(t.Success),
		Warning:   lipgloss.NewStyle().Bold(true).Foreground(t.Warning),
		Error:     lipgloss.NewStyle().Bold(true).Foreground(t.Error),
		Accent:    lipgloss.NewStyle().Bold(true).Foreground(t.Accent),
		Secondary: lipgloss.NewStyle().Foreground(t.Secondary),
		Dim:       lipgloss.NewStyle().Foreground(t.Dim),
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
	if terminalWidth <= 0 {
		return 80
	}
	w := terminalWidth - 2
	if w < 20 {
		w = terminalWidth
	}
	if w > 100 {
		w = 100
	}
	if w > terminalWidth {
		w = terminalWidth
	}
	if w < 10 {
		w = 10
	}
	return w
}

func boxWidth() int {
	return BoxWidthForTerminalWidth(TerminalWidth())
}

func BoxInnerWidth() int {
	return boxWidth() - 6
}

func borderColorForVariant(v BoxVariant) lipgloss.Color {
	t := GetCurrentTheme()
	switch v {
	case BoxSuccess:
		return t.BorderSuccess
	case BoxWarning:
		return t.BorderWarning
	case BoxError:
		return t.BorderError
	default:
		return t.BorderPrimary
	}
}

func RenderBox(variant BoxVariant, title string, content string) string {
	st := GetStyles()

	border := borderColorForVariant(variant)

	w := boxWidth()
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
	body = wrapParagraphs(body, boxWidth()-6)
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

type CardOptions struct {
	Variant   BoxVariant
	Title     string
	Content   string
	FullWidth bool
	MaxWidth  int
}

func RenderCard(opts CardOptions) string {
	st := GetStyles()

	border := borderColorForVariant(opts.Variant)
	maxW := opts.MaxWidth
	if maxW <= 0 {
		maxW = 80
	}

	cardStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(border).
		Padding(0, 1)

	if opts.FullWidth {
		cardStyle = cardStyle.Width(boxWidth())
	} else {
		contentWidth := measureContentWidth(opts.Content, opts.Title)
		w := minInt(maxW, maxInt(contentWidth+6, 30))
		tw := TerminalWidth() - 4
		if w > tw {
			w = tw
		}
		cardStyle = cardStyle.Width(w)
	}

	content := opts.Content
	if strings.TrimSpace(opts.Title) != "" {
		content = st.Title.Render(opts.Title) + "\n" + content
	}

	return cardStyle.Render(content)
}

func SuccessCard(title, body string) string {
	t := GetCurrentTheme()
	st := GetStyles()
	icon := st.Success.Render(t.IconSuccess)
	fullTitle := icon + " " + st.Success.Render(title)
	return RenderCard(CardOptions{
		Variant: BoxSuccess,
		Title:   fullTitle,
		Content: st.Text.Render(body),
	})
}

func ErrorCard(title, body string) string {
	t := GetCurrentTheme()
	st := GetStyles()
	icon := st.Error.Render(t.IconError)
	fullTitle := icon + " " + st.Error.Render(title)
	return RenderCard(CardOptions{
		Variant: BoxError,
		Title:   fullTitle,
		Content: st.Text.Render(body),
	})
}

func WarningCard(title, body string) string {
	t := GetCurrentTheme()
	st := GetStyles()
	icon := st.Warning.Render(t.IconWarning)
	fullTitle := icon + " " + st.Warning.Render(title)
	return RenderCard(CardOptions{
		Variant: BoxWarning,
		Title:   fullTitle,
		Content: st.Text.Render(body),
	})
}

func measureContentWidth(content string, title string) int {
	maxW := 0
	if title != "" {
		for _, line := range strings.Split(runewidth.Truncate(title, 200, ""), "\n") {
			w := runewidth.StringWidth(line)
			if w > maxW {
				maxW = w
			}
		}
	}
	for _, line := range strings.Split(content, "\n") {
		w := runewidth.StringWidth(runewidth.Truncate(stripAnsi(line), 200, ""))
		if w > maxW {
			maxW = w
		}
	}
	return maxW
}

func stripAnsi(s string) string {
	var out strings.Builder
	i := 0
	for i < len(s) {
		if s[i] == '\x1b' {
			i++
			if i < len(s) && s[i] == '[' {
				i++
				for i < len(s) {
					if (s[i] >= '0' && s[i] <= '9') || s[i] == ';' || s[i] == '?' {
						i++
					} else if s[i] >= 'a' && s[i] <= 'z' || s[i] >= 'A' && s[i] <= 'Z' {
						i++
						break
					} else {
						break
					}
				}
			}
			continue
		}
		out.WriteByte(s[i])
		i++
	}
	return out.String()
}

func RenderTitle(text string) string {
	st := GetStyles()
	return st.Title.Render(text)
}

func WithSpinner(message string, fn func() error) error {
	if fn == nil {
		return fmt.Errorf("fn cannot be nil")
	}

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
	cur := ""
	for _, w := range words {
		if runewidth.StringWidth(w) > width {
			if cur != "" {
				out = append(out, cur)
				cur = ""
			}
			out = append(out, w)
			continue
		}
		if cur == "" {
			cur = w
			continue
		}
		if runewidth.StringWidth(cur)+1+runewidth.StringWidth(w) <= width {
			cur = cur + " " + w
			continue
		}
		out = append(out, cur)
		cur = w
	}
	if cur != "" {
		out = append(out, cur)
	}
	return out
}

func roundedBorderRunes() (tl, tr, bl, br, h, v, t, b, l, r, c string) {
	return "╭", "╮", "╰", "╯", "─", "│", "┬", "┴", "├", "┤", "┼"
}

func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
