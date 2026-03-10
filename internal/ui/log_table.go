package ui

import (
	"fmt"
	"strings"

	"github.com/Gazi2050/gli/internal/git"
	"github.com/charmbracelet/lipgloss"
)

func RenderLogTable(title string, entries []git.LogEntry) string {
	st := GetStyles()
	headers := []string{"Hash", "Date & Time", "Author", "Message"}
	rows := make([][]string, 0, len(entries))
	for _, e := range entries {
		rows = append(rows, []string{e.Hash, e.Time, e.Author, e.Message})
	}

	widths := computeColumnWidths(headers, rows, []int{8, 16, 10, 30})
	borderStyle := lipgloss.NewStyle().Foreground(GetCurrentTheme().Success)
	table := renderGridTable(headers, rows, widths, []lipgloss.Style{
		st.Success,
		st.Success,
		st.Warning,
		st.Text,
	}, []string{"left", "left", "left", "left"}, borderStyle, title)

	return table
}

func RenderReflogTable(title string, entries []git.ReflogEntry) string {
	st := GetStyles()
	headers := []string{"Index", "Hash", "Time", "Operation"}
	rows := make([][]string, 0, len(entries))
	for _, e := range entries {
		rows = append(rows, []string{e.Ref, e.Hash, e.Time, e.Message})
	}

	widths := computeColumnWidths(headers, rows, []int{8, 7, 16, 30})
	borderStyle := lipgloss.NewStyle().Foreground(GetCurrentTheme().Secondary)
	table := renderGridTable(headers, rows, widths, []lipgloss.Style{
		st.Muted,
		st.Secondary,
		st.Success,
		st.Text,
	}, []string{"right", "left", "left", "left"}, borderStyle, title)

	return table
}

func computeColumnWidths(headers []string, rows [][]string, mins []int) []int {
	tableWidth := BoxWidthForTerminalWidth(TerminalWidth())
	usable := tableWidth - 4 - (2 * len(headers))
	if usable < 20 {
		usable = 20
	}
	if usable < 40 {
		usable = 40
	}

	widths := make([]int, len(headers))
	for i := range headers {
		widths[i] = lipgloss.Width(headers[i])
	}

	for _, row := range rows {
		for i, cell := range row {
			for _, line := range strings.Split(cell, "\n") {
				if lipgloss.Width(line) > widths[i] {
					widths[i] = lipgloss.Width(line)
				}
			}
		}
	}

	for i := range widths {
		if i < len(mins) && widths[i] < mins[i] {
			widths[i] = mins[i]
		}
	}

	sum := 0
	for _, w := range widths {
		sum += w
	}
	remaining := usable - sum
	if remaining < 0 {
		// shrink wider columns first
		for remaining < 0 {
			maxIdx := -1
			maxW := 0
			for i, w := range widths {
				minW := 8
				if i < len(mins) {
					minW = mins[i]
				}
				if w > minW && w > maxW {
					maxW = w
					maxIdx = i
				}
			}
			if maxIdx == -1 {
				break
			}
			widths[maxIdx]--
			remaining++
		}
	} else if remaining > 0 {
		// give extra space to last column (message/operation)
		widths[len(widths)-1] += remaining
	}

	return widths
}

func renderGridTable(headers []string, rows [][]string, widths []int, styles []lipgloss.Style, align []string, borderStyle lipgloss.Style, title string) string {
	border := borderStyle
	tl, tr, bl, br, h, v, t, b, l, r, c := roundedBorderRunes()
	colW := make([]int, len(widths))
	for i, w := range widths {
		if w < 4 {
			w = 4
		}
		colW[i] = w + 2 // add padding (left/right)
	}

	top := border.Render(tl) + border.Render(strings.Repeat(h, colW[0])) +
		border.Render(t) + border.Render(strings.Repeat(h, colW[1])) +
		border.Render(t) + border.Render(strings.Repeat(h, colW[2])) +
		border.Render(t) + border.Render(strings.Repeat(h, colW[3])) +
		border.Render(tr)

	header := border.Render(v) +
		GetStyles().Text.Bold(true).Render(" "+padRight(headers[0], colW[0]-2)+" ") + border.Render(v) +
		GetStyles().Text.Bold(true).Render(" "+padRight(headers[1], colW[1]-2)+" ") + border.Render(v) +
		GetStyles().Text.Bold(true).Render(" "+padRight(headers[2], colW[2]-2)+" ") + border.Render(v) +
		GetStyles().Text.Bold(true).Render(" "+padRight(headers[3], colW[3]-2)+" ") + border.Render(v)

	sep := border.Render(l) + border.Render(strings.Repeat(h, colW[0])) +
		border.Render(c) + border.Render(strings.Repeat(h, colW[1])) +
		border.Render(c) + border.Render(strings.Repeat(h, colW[2])) +
		border.Render(c) + border.Render(strings.Repeat(h, colW[3])) +
		border.Render(r)

	lines := []string{top, header, sep}
	for _, row := range rows {
		cells := make([][]string, len(row))
		maxLines := 1
		for i, cell := range row {
			cells[i] = wrapCellLines(cell, widths[i])
			if len(cells[i]) > maxLines {
				maxLines = len(cells[i])
			}
		}

		for i := 0; i < maxLines; i++ {
			line := border.Render(v)
			for c := range row {
				val := truncateToWidth(lineAt(cells[c], i), colW[c]-2)
				style := styles[c]
				if c < len(align) && align[c] == "right" {
					line += style.Render(" "+padLeft(val, colW[c]-2)+" ") + border.Render(v)
				} else {
					line += style.Render(" "+padRight(val, colW[c]-2)+" ") + border.Render(v)
				}
			}
			lines = append(lines, line)
		}
	}

	bottom := border.Render(bl) + border.Render(strings.Repeat(h, colW[0])) +
		border.Render(b) + border.Render(strings.Repeat(h, colW[1])) +
		border.Render(b) + border.Render(strings.Repeat(h, colW[2])) +
		border.Render(b) + border.Render(strings.Repeat(h, colW[3])) +
		border.Render(br)

	lines = append(lines, bottom)
	table := strings.Join(lines, "\n")

	if strings.TrimSpace(title) != "" {
		return title + "\n" + table
	}
	return table
}

func roundedBorderRunes() (tl, tr, bl, br, h, v, t, b, l, r, c string) {
	return "╭", "╮", "╰", "╯", "─", "│", "┬", "┴", "├", "┤", "┼"
}

func padLeft(s string, width int) string {
	return fmt.Sprintf("%*s", width, s)
}

// removed centering helpers
