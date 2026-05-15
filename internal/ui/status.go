package ui

import (
	"fmt"
	"strings"

	"github.com/Gazi2050/gli/internal/git"
)

func RenderStatus(data *git.StatusData) string {
	st := GetStyles()

	labelStyle := st.Muted.Width(10)

	branch := labelStyle.Render("Branch") + "  " + st.Accent.Bold(true).Render(data.Branch)
	remote := labelStyle.Render("Remote") + "  " + st.Text.Render(data.Remote)

	var lines []string
	lines = append(lines, branch, remote)

	if len(data.ModifiedFiles) > 0 {
		lines = append(lines, "")
		lines = append(lines, st.Warning.Render(fmt.Sprintf("modified (%d)", len(data.ModifiedFiles))))
		for _, f := range data.ModifiedFiles {
			stats := st.Success.Render(fmt.Sprintf("+%d", f.Additions)) + " " + st.Error.Render(fmt.Sprintf("-%d", f.Deletions))
			lines = append(lines, "  "+st.Warning.Render("~ "+f.Path)+"  "+stats)
		}
	}

	if len(data.StagedFiles) > 0 {
		lines = append(lines, "")
		lines = append(lines, st.Warning.Render(fmt.Sprintf("modified staged (%d)", len(data.StagedFiles))))
		for _, f := range data.StagedFiles {
			stats := st.Success.Render(fmt.Sprintf("+%d", f.Additions)) + " " + st.Error.Render(fmt.Sprintf("-%d", f.Deletions))
			lines = append(lines, "  "+st.Warning.Render("~ "+f.Path)+"  "+stats)
		}
	}

	if len(data.NewFiles) > 0 {
		lines = append(lines, "")
		lines = append(lines, st.Success.Render(fmt.Sprintf("new (%d)", len(data.NewFiles))))
		for _, f := range data.NewFiles {
			lines = append(lines, "  "+st.Success.Render("+ "+f))
		}
	}

	if len(data.DeletedFiles) > 0 {
		lines = append(lines, "")
		lines = append(lines, st.Error.Render(fmt.Sprintf("deleted (%d)", len(data.DeletedFiles))))
		for _, f := range data.DeletedFiles {
			lines = append(lines, "  "+st.Error.Render("- "+f))
		}
	}

	content := strings.Join(lines, "\n")

	return RenderCard(CardOptions{
		Variant: BoxPrimary,
		Title:   st.Title.Render("gli status"),
		Content: content,
	})
}
