package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

const asciiLogo = `
  ██████╗ ██╗     ██╗
 ██╔════╝ ██║     ██║
 ██║  ███╗██║     ██║
 ██║   ██║██║     ██║
 ╚██████╔╝███████╗██║
  ╚═════╝ ╚══════╝╚═╝`

const tagline = "a git wrapper for minimal developers"

type helpRow struct {
	Command string
	Flag    string
	Desc    string
}

var helpRows = []helpRow{
	{Command: "AI Commit", Flag: "gli commit", Desc: "AI commit & push"},
	{Command: "Status", Flag: "gli status", Desc: "Show branch status"},
	{Command: "Branch", Flag: "gli branch", Desc: "List local branches"},
	{Command: "Create Branch", Flag: "gli branch -c <name>", Desc: "Create & push branch"},
	{Command: "Switch Branch", Flag: "gli branch -s <name>", Desc: "Switch branch"},
	{Command: "Log", Flag: "gli log", Desc: "Show commit history"},
	{Command: "Reflog", Flag: "gli reflog", Desc: "Show reference log"},
	{Command: "My Profile", Flag: "gli me", Desc: "Show your profile"},
	{Command: "User Profile", Flag: "gli profile <user>", Desc: "Show any profile"},
}

func RenderHelp() string {
	st := GetStyles()
	t := GetCurrentTheme()

	logoStyle := lipgloss.NewStyle().Foreground(t.Primary).Bold(true)
	taglineStyle := st.Muted.Italic(true)

	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(t.Text).Padding(0, 1)
	commandStyle := lipgloss.NewStyle().Foreground(t.Primary).Bold(true).Padding(0, 1)
	flagStyle := lipgloss.NewStyle().Foreground(t.Warning).Padding(0, 1)
	descStyle := lipgloss.NewStyle().Foreground(t.Text).Padding(0, 1)

	rows := [][]string{}
	for _, r := range helpRows {
		rows = append(rows, []string{r.Command, r.Flag, r.Desc})
	}

	tbl := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(t.BorderPrimary)).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == 0:
				return headerStyle
			case col == 0:
				return commandStyle
			case col == 1:
				return flagStyle
			default:
				return descStyle
			}
		}).
		Headers("Command", "Usage", "Description").
		Rows(rows...)

	usage := st.Muted.Render("Example: ") + st.Accent.Render("gli commit") + st.Muted.Render(" or ") + st.Accent.Render("gli branch -c feat/auth")

	return strings.Join([]string{
		logoStyle.Render(asciiLogo),
		taglineStyle.Render(tagline),
		tbl.Render(),
		usage,
	}, "\n\n")
}
