package ui

import (
	"fmt"
	"strings"

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
	bio := valueOr(user.Bio, "No bio available.")

	header := st.Text.Bold(true).Render(name) + " " + st.Muted.Render("(@"+login+")")

	stats := strings.Join([]string{
		st.Accent.Render(fmt.Sprintf("%d repos", user.PublicRepos)),
		st.Text.Render(fmt.Sprintf("%d followers", user.Followers)),
		st.Text.Render(fmt.Sprintf("%d following", user.Following)),
	}, "  •  ")

	meta := []string{}
	if strings.TrimSpace(user.Location) != "" {
		meta = append(meta, st.Muted.Render("Location:")+" "+st.Text.Render(user.Location))
	}
	if strings.TrimSpace(user.TwitterUsername) != "" {
		meta = append(meta, st.Muted.Render("Twitter:")+" "+st.Text.Render("@"+user.TwitterUsername))
	}
	if strings.TrimSpace(user.Blog) != "" {
		meta = append(meta, st.Muted.Render("Site:")+" "+st.Text.Render(user.Blog))
	}
	joined := user.CreatedAt
	if len(joined) >= 10 {
		joined = joined[:10]
	}
	if strings.TrimSpace(joined) != "" {
		meta = append(meta, st.Muted.Render("Joined:")+" "+st.Text.Render(joined))
	}

	sections := []string{header, bio, stats}
	if len(meta) > 0 {
		sections = append(sections, strings.Join(meta, "  •  "))
	}
	content := strings.Join(sections, "\n")

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
