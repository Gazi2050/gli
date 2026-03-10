package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const asciiLogo = `
 ██████╗ ██╗     ██╗
██╔════╝ ██║     ██║
██║  ███╗██║     ██║
██║   ██║██║     ██║
╚██████╔╝███████╗██║
 ╚═════╝ ╚══════╝╚═╝
`

const tagline = "a git wrapper for make developer life easy"

type helpRow struct {
	Command     string
	Flag        string
	Description string
}

var helpRows = []helpRow{
	{Command: "Commit & Push", Flag: "-c, --commit", Description: "Stage all, commit with msg, and push"},
	{Command: "AI Commit", Flag: "-ac, --ai-commit", Description: "Generate AI message and push"},
	{Command: "Log", Flag: "-l, --log", Description: "View commit history graph"},
	{Command: "Reflog", Flag: "-rl, --reflog", Description: "View git reflog"},
	{Command: "Reset", Flag: "-rs, --reset", Description: "Reset last commit (soft/hard)"},
	{Command: "Switch Branch", Flag: "-s, --switch\n-lb, --local-branch\n-rb, --remote-branch", Description: "Create branch\n[-lb] Local only\n[-rb] Push remote"},
	{Command: "", Flag: "", Description: ""},
	{Command: "Change Time", Flag: "-ct, --changeTime", Description: "Update commit timestamp(s)"},
	{Command: "Change Author", Flag: "-ca, --changeAuthor", Description: "Update commit author identity"},
	{Command: "Change Message", Flag: "-cm, --changeMessage", Description: "Update last commit message"},
	{Command: "", Flag: "", Description: ""},
	{Command: "No Verify", Flag: "-nv, --no-verify", Description: "Skip git hooks during commit"},
	{Command: "My Profile", Flag: "me", Description: "View your GitHub profile"},
	{Command: "User Profile", Flag: "profile <user>", Description: "View a specific GitHub profile"},
}

func RenderHelp() string {
	st := GetStyles()
	logoStyle := lipgloss.NewStyle().Foreground(GetCurrentTheme().Primary).Bold(true)
	taglineStyle := st.Muted.Italic(true)
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(GetCurrentTheme().Text)
	commandCellStyle := lipgloss.NewStyle().Foreground(GetCurrentTheme().Primary).Bold(true)
	flagCellStyle := lipgloss.NewStyle().Foreground(GetCurrentTheme().Warning)
	descriptionCellStyle := lipgloss.NewStyle().Foreground(GetCurrentTheme().Text)
	borderStyle := st.Muted

	commandWidth, flagWidth, descriptionWidth := columnWidths(helpRows)
	table := renderTable(
		[]string{"Command", "Flag", "Description"},
		helpRows,
		[]int{commandWidth, flagWidth, descriptionWidth},
		headerStyle,
		commandCellStyle,
		flagCellStyle,
		descriptionCellStyle,
		borderStyle,
	)
	usage := st.Muted.Render("Usage example: ") + st.Accent.Render("gli -c 'feat: msg'") + st.Muted.Render(" or ") + st.Accent.Render("gli -ac")

	headerBlock := strings.Join([]string{
		logoStyle.Render(asciiLogo),
		taglineStyle.Render(tagline),
	}, "\n")

	return strings.Join([]string{
		headerBlock,
		table,
		"",
		usage,
	}, "\n\n")
}

func columnWidths(rows []helpRow) (commandWidth, flagWidth, descriptionWidth int) {
	tableWidth := BoxWidthForTerminalWidth(TerminalWidth())
	usable := tableWidth - 4 - (2 * 3) // borders/separators + padding per column
	if usable < 40 {
		usable = 40
	}

	minCmd, minFlag, minDesc := 14, 18, 20
	maxCmd, maxFlag := 22, 30

	longestCmd, longestFlag := len("Command"), len("Flag")
	for _, row := range rows {
		if row.Command == "" && row.Flag == "" && row.Description == "" {
			continue
		}
		for _, line := range strings.Split(row.Command, "\n") {
			if lipgloss.Width(line) > longestCmd {
				longestCmd = lipgloss.Width(line)
			}
		}
		for _, line := range strings.Split(row.Flag, "\n") {
			if lipgloss.Width(line) > longestFlag {
				longestFlag = lipgloss.Width(line)
			}
		}
	}

	commandWidth = clamp(longestCmd, minCmd, maxCmd)
	flagWidth = clamp(longestFlag, minFlag, maxFlag)
	descriptionWidth = usable - commandWidth - flagWidth

	if descriptionWidth < minDesc {
		need := minDesc - descriptionWidth

		shrinkCmd := min(need, max(0, commandWidth-minCmd))
		commandWidth -= shrinkCmd
		need -= shrinkCmd

		shrinkFlag := min(need, max(0, flagWidth-minFlag))
		flagWidth -= shrinkFlag
		need -= shrinkFlag

		descriptionWidth = usable - commandWidth - flagWidth
		if descriptionWidth < 10 {
			descriptionWidth = 10
		}
	}

	return commandWidth, flagWidth, descriptionWidth
}

func padRight(s string, width int) string {
	return fmt.Sprintf("%-*s", width, s)
}

func lineAt(lines []string, idx int) string {
	if idx < 0 || idx >= len(lines) {
		return ""
	}
	return lines[idx]
}

func max3(a, b, c int) int {
	max := a
	if b > max {
		max = b
	}
	if c > max {
		max = c
	}
	return max
}

func wrapCellLines(s string, width int) []string {
	parts := strings.Split(s, "\n")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimRight(p, " ")
		if strings.TrimSpace(p) == "" {
			out = append(out, "")
			continue
		}
		out = append(out, wrapLine(p, width)...)
	}
	return out
}

func renderTable(headers []string, rows []helpRow, widths []int, headerStyle, cmdStyle, flagStyle, descStyle, borderStyle lipgloss.Style) string {
	tl, tr, bl, br, h, v, t, b, l, r, c := roundedBorderRunes()
	colW := make([]int, len(widths))
	for i, w := range widths {
		if w < 4 {
			w = 4
		}
		colW[i] = w + 2
	}

	top := borderStyle.Render(tl) + borderStyle.Render(strings.Repeat(h, colW[0])) +
		borderStyle.Render(t) + borderStyle.Render(strings.Repeat(h, colW[1])) +
		borderStyle.Render(t) + borderStyle.Render(strings.Repeat(h, colW[2])) +
		borderStyle.Render(tr)

	header := borderStyle.Render(v) +
		headerStyle.Render(" "+padRight(headers[0], colW[0]-2)+" ") + borderStyle.Render(v) +
		headerStyle.Render(" "+padRight(headers[1], colW[1]-2)+" ") + borderStyle.Render(v) +
		headerStyle.Render(" "+padRight(headers[2], colW[2]-2)+" ") + borderStyle.Render(v)

	sep := borderStyle.Render(l) + borderStyle.Render(strings.Repeat(h, colW[0])) +
		borderStyle.Render(c) + borderStyle.Render(strings.Repeat(h, colW[1])) +
		borderStyle.Render(c) + borderStyle.Render(strings.Repeat(h, colW[2])) +
		borderStyle.Render(r)

	lines := []string{top, header, sep}
	for _, row := range rows {
		if row.Command == "" && row.Flag == "" && row.Description == "" {
			blank := borderStyle.Render(v) + padRight("", colW[0]) +
				borderStyle.Render(v) + padRight("", colW[1]) +
				borderStyle.Render(v) + padRight("", colW[2]) +
				borderStyle.Render(v)
			lines = append(lines, blank)
			continue
		}

		commandLines := wrapCellLines(row.Command, widths[0])
		flagLines := wrapCellLines(row.Flag, widths[1])
		descLines := wrapCellLines(row.Description, widths[2])
		maxLines := max3(len(commandLines), len(flagLines), len(descLines))

		for i := 0; i < maxLines; i++ {
			commandValue := truncateToWidth(lineAt(commandLines, i), colW[0]-2)
			flagValue := truncateToWidth(lineAt(flagLines, i), colW[1]-2)
			descValue := truncateToWidth(lineAt(descLines, i), colW[2]-2)

			line := borderStyle.Render(v) +
				cmdStyle.Render(" "+padRight(commandValue, colW[0]-2)+" ") +
				borderStyle.Render(v) +
				flagStyle.Render(" "+padRight(flagValue, colW[1]-2)+" ") +
				borderStyle.Render(v) +
				descStyle.Render(" "+padRight(descValue, colW[2]-2)+" ") +
				borderStyle.Render(v)
			lines = append(lines, line)
		}
	}

	bottom := borderStyle.Render(bl) + borderStyle.Render(strings.Repeat(h, colW[0])) +
		borderStyle.Render(b) + borderStyle.Render(strings.Repeat(h, colW[1])) +
		borderStyle.Render(b) + borderStyle.Render(strings.Repeat(h, colW[2])) +
		borderStyle.Render(br)

	lines = append(lines, bottom)
	return strings.Join(lines, "\n")
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
