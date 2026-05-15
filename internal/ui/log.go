package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/mattn/go-runewidth"

	"github.com/Gazi2050/gli/internal/git"
)

type logModel struct {
	entries    []git.LogEntry
	title      string
	columnName string
	cursor     int
	pageSize   int
	quitting   bool
}

func NewLogModel(entries []git.LogEntry, title string, columnName string) logModel {
	return logModel{
		entries:    entries,
		title:      title,
		columnName: columnName,
		pageSize:   15,
	}
}

func (m logModel) Init() tea.Cmd { return nil }

func (m logModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "j", "down":
			if m.cursor < len(m.entries)-1 {
				m.cursor++
			}
		case "k", "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "pgdown":
			m.cursor = minInt(m.cursor+m.pageSize, len(m.entries)-1)
		case "pgup":
			m.cursor = maxInt(m.cursor-m.pageSize, 0)
		case "g", "home":
			m.cursor = 0
		case "G", "end":
			m.cursor = len(m.entries) - 1
		}
	}
	return m, nil
}

func (m logModel) View() string {
	if m.quitting {
		return ""
	}

	st := GetStyles()
	t := GetCurrentTheme()

	totalPages := (len(m.entries) + m.pageSize - 1) / m.pageSize
	currentPage := (m.cursor / m.pageSize) + 1

	start := (currentPage - 1) * m.pageSize
	end := minInt(start+m.pageSize, len(m.entries))
	pageEntries := m.entries[start:end]

	hashStyle := lipgloss.NewStyle().Foreground(t.Primary).Bold(true).Padding(0, 1)
	authorStyle := lipgloss.NewStyle().Foreground(t.Secondary).Padding(0, 1)
	commitStyle := lipgloss.NewStyle().Foreground(t.Text).Padding(0, 1)
	dateStyle := lipgloss.NewStyle().Foreground(t.Muted).Padding(0, 1)
	highlightStyle := lipgloss.NewStyle().Foreground(t.Accent).Bold(true).Padding(0, 1)
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(t.Text).Padding(0, 1)

	rows := [][]string{}
	for i, entry := range pageEntries {
		actualIdx := start + i
		commit := entry.Commit
		if runewidth.StringWidth(commit) > 50 {
			commit = runewidth.Truncate(commit, 47, "...")
		}
		author := entry.Author
		if runewidth.StringWidth(author) > 15 {
			author = runewidth.Truncate(author, 12, "...")
		}
		if actualIdx == m.cursor {
			rows = append(rows, []string{
				highlightStyle.Render("▸ " + entry.Hash),
				highlightStyle.Render(author),
				highlightStyle.Render(commit),
				highlightStyle.Render(entry.Date),
			})
		} else {
			rows = append(rows, []string{
				hashStyle.Render(entry.Hash),
				authorStyle.Render(author),
				commitStyle.Render(commit),
				dateStyle.Render(entry.Date),
			})
		}
	}

	tbl := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(t.BorderPrimary)).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == 0 {
				return headerStyle
			}
			return lipgloss.NewStyle().Padding(0, 1)
		}).
		Headers("Hash", "Author", m.columnName, "Date").
		Rows(rows...)

	pagination := st.Muted.Render(fmt.Sprintf("  Page %d/%d  •  %d entries  •  ↑/↓ navigate  •  q quit", currentPage, totalPages, len(m.entries)))

	title := st.Title.Render(m.title)

	return lipgloss.NewStyle().Padding(0, 1).Render(
		title + "\n\n" + tbl.Render() + "\n" + pagination,
	)
}

func RenderLog(entries []git.LogEntry, title string, columnName string) string {
	if len(entries) == 0 {
		st := GetStyles()
		return RenderCard(CardOptions{
			Variant: BoxPrimary,
			Title:   st.Title.Render(title),
			Content: st.Muted.Render("No entries found."),
		})
	}

	model := NewLogModel(entries, title, columnName)
	p := tea.NewProgram(model)
	finalModel, err := p.Run()
	if err != nil {
		return ErrorCard("Log Error", err.Error())
	}
	if lm, ok := finalModel.(logModel); ok && lm.quitting {
		return ""
	}
	return ""
}
