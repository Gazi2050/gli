package ui

import (
	"fmt"
	"strings"

	"github.com/Gazi2050/gli/internal/api"
)

func RenderProfile(user *api.User, repos []api.Repo) string {
	st := GetStyles()

	if user == nil {
		return MessageBox(BoxWarning, "GitHub Profile", "No profile data available")
	}

	name := valueOr(user.Name, "GitHub User")
	login := valueOr(user.Login, "N/A")
	bio := valueOr(user.Bio, "No bio available.")

	header := st.Text.Bold(true).Render(name) + " " + st.Muted.Render("(@"+login+")")
	bioLine := st.Text.Render(wrapParagraphs(bio, boxWidth()-6))

	stats := strings.Join([]string{
		st.Accent.Render("Public Repos:") + " " + fmt.Sprintf("%d", user.PublicRepos),
		st.Accent.Render("Followers:") + " " + fmt.Sprintf("%d", user.Followers),
		st.Accent.Render("Following:") + " " + fmt.Sprintf("%d", user.Following),
	}, "  •  ")

	metaLines := []string{}
	if strings.TrimSpace(user.Location) != "" {
		metaLines = append(metaLines, st.Muted.Render("Location:")+" "+st.Text.Render(user.Location))
	}
	if strings.TrimSpace(user.TwitterUsername) != "" {
		metaLines = append(metaLines, st.Muted.Render("Twitter:")+" "+st.Text.Render("@"+user.TwitterUsername))
	}
	if strings.TrimSpace(user.Blog) != "" {
		metaLines = append(metaLines, st.Muted.Render("Site:")+" "+st.Text.Render(user.Blog))
	}
	joined := user.CreatedAt
	if len(joined) >= 10 {
		joined = joined[:10]
	}
	if strings.TrimSpace(joined) != "" {
		metaLines = append(metaLines, st.Muted.Render("Joined:")+" "+st.Text.Render(joined))
	}

	sections := []string{
		header,
		bioLine,
		"",
		st.Text.Render(stats),
	}

	if len(metaLines) > 0 {
		sections = append(sections, "", strings.Join(metaLines, "\n"))
	}

	return RenderBox(BoxPrimary, "GitHub Profile", strings.Join(sections, "\n"))
}

func valueOr(value, fallback string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return fallback
	}
	return value
}
