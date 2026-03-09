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
	{Command: "Change Time", Flag: "-ct, --changeTime", Description: "Update commit timestamp(s)"},
	{Command: "Change Author", Flag: "-ca, --changeAuthor", Description: "Update commit author identity"},
	{Command: "Change Message", Flag: "-cm, --changeMessage", Description: "Update last commit message"},
	{Command: "No Verify", Flag: "-nv, --no-verify", Description: "Skip git hooks during commit"},
	{Command: "My Profile", Flag: "me", Description: "View your GitHub profile"},
	{Command: "User Profile", Flag: "profile <user>", Description: "View a specific GitHub profile"},
}

func RenderHelp() string {
	logoStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Bold(true)
	taglineStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("248")).Italic(true)
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("252"))
	headerStyle := lipgloss.NewStyle().Bold(true)
	commandHeaderStyle := headerStyle.Foreground(lipgloss.Color("42"))
	flagHeaderStyle := headerStyle.Foreground(lipgloss.Color("220"))
	descriptionHeaderStyle := headerStyle.Foreground(lipgloss.Color("252"))
	commandCellStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Bold(true)
	flagCellStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("220"))
	descriptionCellStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	borderStyle := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("240")).Padding(0, 1)
	usageStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	usageCmdStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Bold(true)

	commandWidth, flagWidth, descriptionWidth := columnWidths(helpRows)

	header := strings.Join([]string{
		commandHeaderStyle.Render(padRight("Command", commandWidth)),
		flagHeaderStyle.Render(padRight("Flag", flagWidth)),
		descriptionHeaderStyle.Render(padRight("Description", descriptionWidth)),
	}, "  ")

	lines := []string{titleStyle.Render("Available Commands"), "", header}
	for _, row := range helpRows {
		commandLines := strings.Split(row.Command, "\n")
		flagLines := strings.Split(row.Flag, "\n")
		descriptionLines := strings.Split(row.Description, "\n")

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

	table := borderStyle.Render(strings.Join(lines, "\n"))
	usage := usageStyle.Render("Usage example: ") + usageCmdStyle.Render("gli -c 'feat: msg'") + usageStyle.Render(" or ") + usageCmdStyle.Render("gli -ac")

	sections := []string{
		logoStyle.Render(asciiLogo),
		taglineStyle.Render(tagline),
		table,
		usage,
	}

	return strings.Join(sections, "\n\n")
}

func columnWidths(rows []helpRow) (commandWidth, flagWidth, descriptionWidth int) {
	commandWidth = len("Command")
	flagWidth = len("Flag")
	descriptionWidth = len("Description")

	for _, row := range rows {
		for _, line := range strings.Split(row.Command, "\n") {
			if len(line) > commandWidth {
				commandWidth = len(line)
			}
		}
		for _, line := range strings.Split(row.Flag, "\n") {
			if len(line) > flagWidth {
				flagWidth = len(line)
			}
		}
		for _, line := range strings.Split(row.Description, "\n") {
			if len(line) > descriptionWidth {
				descriptionWidth = len(line)
			}
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
