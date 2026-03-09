package ui

import (
	"fmt"
	"strings"

	"github.com/Gazi2050/gli/internal/api"
	"github.com/charmbracelet/lipgloss"
)

func RenderProfile(user *api.User, repos []api.Repo) string {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("42"))
	nameStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("252"))
	loginStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	bioStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	labelGreen := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("42"))
	labelMagenta := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	labelYellow := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("220"))
	metaLabelStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("111"))
	metaValueStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	repoNameStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("111"))
	repoDescStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("250"))
	panelStyle := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("42")).Padding(1, 2)

	if user == nil {
		return panelStyle.Render("No profile data available")
	}

	name := valueOr(user.Name, "GitHub User")
	login := valueOr(user.Login, "N/A")
	bio := valueOr(user.Bio, "No bio available.")

	header := nameStyle.Render(name) + " " + loginStyle.Render("(@"+login+")")
	bioLine := bioStyle.Render(bio)

	stats := strings.Join([]string{
		labelGreen.Render("Public Repos:") + " " + fmt.Sprintf("%d", user.PublicRepos),
		labelMagenta.Render("Followers:") + " " + fmt.Sprintf("%d", user.Followers),
		labelYellow.Render("Following:") + " " + fmt.Sprintf("%d", user.Following),
	}, "  •  ")

	metaLines := []string{}
	if strings.TrimSpace(user.Location) != "" {
		metaLines = append(metaLines, metaLabelStyle.Render("Location:")+" "+metaValueStyle.Render(user.Location))
	}
	if strings.TrimSpace(user.TwitterUsername) != "" {
		metaLines = append(metaLines, metaLabelStyle.Render("Twitter:")+" "+metaValueStyle.Render("@"+user.TwitterUsername))
	}
	if strings.TrimSpace(user.Blog) != "" {
		metaLines = append(metaLines, metaLabelStyle.Render("Blog:")+" "+metaValueStyle.Render(user.Blog))
	}
	joined := user.CreatedAt
	if len(joined) >= 10 {
		joined = joined[:10]
	}
	if strings.TrimSpace(joined) != "" {
		metaLines = append(metaLines, metaLabelStyle.Render("Joined:")+" "+metaValueStyle.Render(joined))
	}

	sections := []string{
		titleStyle.Render("GitHub Profile"),
		"",
		header,
		bioLine,
		"",
		stats,
	}

	if len(metaLines) > 0 {
		sections = append(sections, "", strings.Join(metaLines, "\n"))
	}

	if len(repos) > 0 {
		sections = append(sections, "", titleStyle.Render("Recent Repositories"))
		for _, repo := range repos {
			repoLine := "- " + repoNameStyle.Render(valueOr(repo.Name, "unnamed"))
			if strings.TrimSpace(repo.Language) != "" {
				repoLine += " " + metaValueStyle.Render("("+repo.Language+")")
			}
			repoLine += " " + metaValueStyle.Render(fmt.Sprintf("★ %d", repo.StargazersCount))
			sections = append(sections, repoLine)
			if strings.TrimSpace(repo.Description) != "" {
				sections = append(sections, "  "+repoDescStyle.Render(repo.Description))
			}
		}
	}

	return panelStyle.Render(strings.Join(sections, "\n"))
}

func valueOr(value, fallback string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return fallback
	}
	return value
}
