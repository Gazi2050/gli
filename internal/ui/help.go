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
  ╚═════╝ ╚══════╝╚═╝`

const tagline = "a git wrapper for modern developers"

type helpRow struct {
	Command string
	Flag    string
	Desc    string
}

var helpRows = []helpRow{
	{Command: "AI Commit", Flag: "gli commit", Desc: "Stage all, generate AI message, commit & push"},
	{Command: "Create Branch", Flag: "gli branch -c <name>", Desc: "Create branch + push to remote"},
	{Command: "Switch Branch", Flag: "gli branch -s <name>", Desc: "Switch to existing branch"},
	{Command: "My Profile", Flag: "gli me", Desc: "Show your GitHub profile"},
	{Command: "User Profile", Flag: "gli profile <user>", Desc: "Show any GitHub profile"},
}

func RenderHelp() string {
	st := GetStyles()
	t := GetCurrentTheme()

	logoStyle := lipgloss.NewStyle().Foreground(t.Primary).Bold(true)
	taglineStyle := st.Muted.Italic(true)

	commandStyle := lipgloss.NewStyle().Foreground(t.Primary).Bold(true)
	flagStyle := lipgloss.NewStyle().Foreground(t.Warning)
	descStyle := lipgloss.NewStyle().Foreground(t.Text)
	borderStyle := st.Muted
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(t.Text)

	tl, tr, bl, br, h, v, tsep, bsep, l, r, c := roundedBorderRunes()

	colW := []int{18, 26, 42}
	padCol := func(s string, w int) string {
		return fmt.Sprintf("%-*s", w, s)
	}

	top := borderStyle.Render(tl) + borderStyle.Render(strings.Repeat(h, colW[0])) +
		borderStyle.Render(tsep) + borderStyle.Render(strings.Repeat(h, colW[1])) +
		borderStyle.Render(tsep) + borderStyle.Render(strings.Repeat(h, colW[2])) +
		borderStyle.Render(tr)

	header := borderStyle.Render(v) +
		headerStyle.Render(" "+padCol("Command", colW[0])+" ") + borderStyle.Render(v) +
		headerStyle.Render(" "+padCol("Usage", colW[1])+" ") + borderStyle.Render(v) +
		headerStyle.Render(" "+padCol("Description", colW[2])+" ") + borderStyle.Render(v)

	sep := borderStyle.Render(l) + borderStyle.Render(strings.Repeat(h, colW[0])) +
		borderStyle.Render(c) + borderStyle.Render(strings.Repeat(h, colW[1])) +
		borderStyle.Render(c) + borderStyle.Render(strings.Repeat(h, colW[2])) +
		borderStyle.Render(r)

	rows := []string{}
	for _, row := range helpRows {
		line := borderStyle.Render(v) +
			commandStyle.Render(" "+padCol(row.Command, colW[0])+" ") +
			borderStyle.Render(v) +
			flagStyle.Render(" "+padCol(row.Flag, colW[1])+" ") +
			borderStyle.Render(v) +
			descStyle.Render(" "+padCol(row.Desc, colW[2])+" ") +
			borderStyle.Render(v)
		rows = append(rows, line)
	}

	bottom := borderStyle.Render(bl) + borderStyle.Render(strings.Repeat(h, colW[0])) +
		borderStyle.Render(bsep) + borderStyle.Render(strings.Repeat(h, colW[1])) +
		borderStyle.Render(bsep) + borderStyle.Render(strings.Repeat(h, colW[2])) +
		borderStyle.Render(br)

	table := strings.Join(append([]string{top, header, sep}, append(rows, bottom)...), "\n")

	usage := st.Muted.Render("Example: ") + st.Accent.Render("gli commit") + st.Muted.Render(" or ") + st.Accent.Render("gli branch -c feat/auth")

	return strings.Join([]string{
		logoStyle.Render(asciiLogo),
		taglineStyle.Render(tagline),
		table,
		usage,
	}, "\n\n")
}
