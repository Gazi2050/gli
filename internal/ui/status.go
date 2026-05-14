package ui

import (
	"fmt"

	"github.com/Gazi2050/gli/internal/git"
)

func RenderStatus(data *git.StatusData) string {
	st := GetStyles()
	theme := GetCurrentTheme()

	labelStyle := st.Muted.Width(10)
	valueStyle := st.Text

	branch := labelStyle.Render("Branch") + "  " + st.Accent.Bold(true).Render(data.Branch)
	remote := labelStyle.Render("Remote") + "  " + valueStyle.Render(data.Remote)

	syncLine := labelStyle.Render("Sync")
	if data.Ahead > 0 {
		syncLine += "  " + st.Success.Render(fmt.Sprintf("↑%d", data.Ahead))
	}
	if data.Behind > 0 {
		if data.Ahead > 0 {
			syncLine += " "
		}
		syncLine += "  " + st.Error.Render(fmt.Sprintf("↓%d", data.Behind))
	}
	if data.Ahead == 0 && data.Behind == 0 {
		syncLine += "  " + st.Success.Render(theme.IconSuccess + " up to date")
	}

	staged := st.Success.Render(fmt.Sprintf("%d staged", data.Staged))
	unstaged := st.Warning.Render(fmt.Sprintf("%d modified", data.Unstaged))
	untracked := st.Error.Render(fmt.Sprintf("%d untracked", data.Untracked))
	files := labelStyle.Render("Files") + "  " + staged + "  •  " + unstaged + "  •  " + untracked

	lastCommit := labelStyle.Render("Last") + "  " + valueStyle.Render(data.LastCommit) + "  " + st.Muted.Render(data.LastDate)

	content := branch + "\n" + remote + "\n" + syncLine + "\n" + files + "\n" + lastCommit

	return RenderCard(CardOptions{
		Variant: BoxPrimary,
		Title:   st.Title.Render(theme.IconBranch) + " " + st.Title.Render("gli status"),
		Content: content,
	})
}
