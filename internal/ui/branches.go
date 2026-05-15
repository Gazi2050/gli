package ui

import (
	"strings"

	"github.com/Gazi2050/gli/internal/git"
)

func RenderBranches(branches []git.BranchEntry) string {
	if len(branches) == 0 {
		st := GetStyles()
		return RenderCard(CardOptions{
			Variant: BoxPrimary,
			Title:   st.Title.Render("gli branch"),
			Content: st.Muted.Render("No branches found."),
		})
	}

	st := GetStyles()

	var lines []string
	for _, b := range branches {
		if b.Current {
			lines = append(lines, st.Accent.Bold(true).Render("▸ "+b.Name))
		} else {
			lines = append(lines, st.Text.Render("  "+b.Name))
		}
	}

	title := st.Title.Render("gli branch")

	return RenderCard(CardOptions{
		Variant: BoxPrimary,
		Title:   title,
		Content: strings.Join(lines, "\n"),
	})
}
