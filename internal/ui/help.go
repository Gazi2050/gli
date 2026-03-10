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

type helpGroup struct {
	Title string
	Rows  []helpRow
}

var helpGroups = []helpGroup{
	{
		Title: "Commit",
		Rows: []helpRow{
			{Command: "Commit & Push", Flag: "-c, --commit", Description: "Stage all, commit with msg, and push"},
			{Command: "AI Commit", Flag: "-ac, --ai-commit", Description: "Generate AI message and push"},
			{Command: "No Verify", Flag: "-nv, --no-verify", Description: "Skip git hooks during commit"},
		},
	},
	{
		Title: "History",
		Rows: []helpRow{
			{Command: "Change Time", Flag: "-ct, --changeTime", Description: "Update commit timestamp(s)"},
			{Command: "Change Author", Flag: "-ca, --changeAuthor", Description: "Update commit author identity"},
			{Command: "Change Message", Flag: "-cm, --changeMessage", Description: "Update last commit message"},
		},
	},
	{
		Title: "Branches & Log",
		Rows: []helpRow{
			{Command: "Log", Flag: "-l, --log", Description: "View commit history graph"},
			{Command: "Reflog", Flag: "-rl, --reflog", Description: "View git reflog"},
			{Command: "Reset", Flag: "-rs, --reset", Description: "Reset last commit (soft/hard)"},
			{Command: "Switch Branch", Flag: "-s, --switch\n-lb, --local-branch\n-rb, --remote-branch", Description: "Create branch\n[-lb] Local only\n[-rb] Push remote"},
		},
	},
	{
		Title: "Profile",
		Rows: []helpRow{
			{Command: "My Profile", Flag: "me", Description: "View your GitHub profile"},
			{Command: "User Profile", Flag: "profile <user>", Description: "View a specific GitHub profile"},
		},
	},
}

func RenderHelp() string {
	st := GetStyles()
	logoStyle := lipgloss.NewStyle().Foreground(GetCurrentTheme().Primary).Bold(true)
	taglineStyle := st.Muted.Italic(true)
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(GetCurrentTheme().Primary)
	flagHeaderStyle := lipgloss.NewStyle().Bold(true).Foreground(GetCurrentTheme().Warning)
	descHeaderStyle := lipgloss.NewStyle().Bold(true).Foreground(GetCurrentTheme().Text)
	commandCellStyle := lipgloss.NewStyle().Foreground(GetCurrentTheme().Primary).Bold(true)
	flagCellStyle := lipgloss.NewStyle().Foreground(GetCurrentTheme().Warning)
	descriptionCellStyle := lipgloss.NewStyle().Foreground(GetCurrentTheme().Text)
	groupStyle := st.Muted.Bold(true)

	commandWidth, flagWidth, descriptionWidth := columnWidths(helpGroups)

	header := strings.Join([]string{
		headerStyle.Render(padRight("Command", commandWidth)),
		flagHeaderStyle.Render(padRight("Flag", flagWidth)),
		descHeaderStyle.Render(padRight("Description", descriptionWidth)),
	}, "  ")

	lines := []string{header}
	for gi, group := range helpGroups {
		if gi > 0 {
			lines = append(lines, "")
		}
		lines = append(lines, groupStyle.Render(group.Title))

		for _, row := range group.Rows {
			commandLines := wrapCellLines(row.Command, commandWidth)
			flagLines := wrapCellLines(row.Flag, flagWidth)
			descriptionLines := wrapCellLines(row.Description, descriptionWidth)

			maxLines := max3(len(commandLines), len(flagLines), len(descriptionLines))
			for i := 0; i < maxLines; i++ {
				commandValue := lineAt(commandLines, i)
				flagValue := lineAt(flagLines, i)
				descriptionValue := lineAt(descriptionLines, i)

				line := strings.Join([]string{
					commandCellStyle.Render(padRight(commandValue, commandWidth)),
					flagCellStyle.Render(padRight(flagValue, flagWidth)),
					descriptionCellStyle.Render(padRight(descriptionValue, descriptionWidth)),
				}, "  ")
				lines = append(lines, line)
			}
		}
	}

	table := strings.Join(lines, "\n")
	usage := st.Muted.Render("Usage example: ") + st.Accent.Render("gli -c 'feat: msg'") + st.Muted.Render(" or ") + st.Accent.Render("gli -ac")

	headerBlock := strings.Join([]string{
		logoStyle.Render(asciiLogo),
		taglineStyle.Render(tagline),
	}, "\n")

	return strings.Join([]string{
		headerBlock,
		RenderBox(BoxPrimary, "Available Commands", strings.Join([]string{table, "", usage}, "\n")),
	}, "\n\n")
}

func columnWidths(groups []helpGroup) (commandWidth, flagWidth, descriptionWidth int) {
	// RenderBox uses width = boxWidth(), with border+padding roughly consuming 6 columns.
	inner := boxWidth() - 6
	if inner < 50 {
		inner = 50
	}

	// 3 columns + 2 separators ("  ") => 4 extra spaces.
	usable := inner - 4
	if usable < 30 {
		usable = 30
	}

	// Start from sensible minima so narrow terminals still look like a table.
	minCmd, minFlag, minDesc := 12, 18, 20
	maxCmd, maxFlag := 18, 26

	longestCmd, longestFlag := len("Command"), len("Flag")
	for _, g := range groups {
		for _, row := range g.Rows {
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
	}

	commandWidth = clamp(longestCmd, minCmd, maxCmd)
	flagWidth = clamp(longestFlag, minFlag, maxFlag)
	descriptionWidth = usable - commandWidth - flagWidth

	// Ensure description has enough space; shrink other cols if needed.
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
