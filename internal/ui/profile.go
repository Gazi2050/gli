package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/Gazi2050/gli/internal/api"
)

func RenderProfile(user *api.User, repos []api.Repo) string {
	st := GetStyles()
	theme := GetCurrentTheme()

	if user == nil {
		return WarningCard("GitHub Profile", "No profile data available")
	}

	name := valueOr(user.Name, "GitHub User")
	login := valueOr(user.Login, "N/A")
	bio := valueOr(user.Bio, "")

	var lines []string

	lines = append(lines, st.Text.Bold(true).Render(name))
	lines = append(lines, st.Muted.Render("@"+login))

	if bio != "" {
		lines = append(lines, bio)
	}

	lines = append(lines, "")

	label := func(s string) string { return st.Muted.Render("◈ "+s+":") }

	lines = append(lines, label("repos")+" "+st.Accent.Render(fmt.Sprint(user.PublicRepos)))
	lines = append(lines, label("followers")+" "+st.Text.Render(fmt.Sprint(user.Followers)))
	lines = append(lines, label("following")+" "+st.Text.Render(fmt.Sprint(user.Following)))

	if strings.TrimSpace(user.Location) != "" {
		lines = append(lines, label("location")+" "+st.Text.Render(user.Location))
	}
	if strings.TrimSpace(user.TwitterUsername) != "" {
		lines = append(lines, label("twitter")+" "+st.Text.Render("@"+user.TwitterUsername))
	}
	if strings.TrimSpace(user.Blog) != "" {
		lines = append(lines, label("site")+" "+st.Text.Render(user.Blog))
	}

	joined := user.CreatedAt
	if len(joined) >= 10 {
		joined = joined[:10]
	}
	if t, err := time.Parse("2006-01-02", joined); err == nil {
		joined = t.Format("Jan 2, 2006")
	}
	if strings.TrimSpace(joined) != "" {
		lines = append(lines, label("joined")+" "+st.Text.Render(joined))
	}

	content := strings.Join(lines, "\n")

	return RenderCard(CardOptions{
		Variant: BoxPrimary,
		Title:   st.Title.Render(theme.IconSuccess) + " " + st.Title.Render("GitHub Profile"),
		Content: content,
	})
}

func valueOr(value, fallback string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return fallback
	}
	return value
}
