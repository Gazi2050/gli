package controllers

import (
	"fmt"
	"strings"

	"github.com/Gazi2050/gli/internal/api"
	"github.com/Gazi2050/gli/internal/ui"
)

type GitProfileIdentity interface {
	GetGithubUsername() (string, error)
}

type GitHubProfileAPI interface {
	GetUser(username string) (*api.User, error)
	GetUserRepos(username string, limit int) ([]api.Repo, error)
}

type ProfileController struct {
	git       GitProfileIdentity
	github    GitHubProfileAPI
	repoLimit int
}

func NewProfileController(git GitProfileIdentity, github GitHubProfileAPI) (*ProfileController, error) {
	if git == nil {
		return nil, fmt.Errorf("git identity source is required")
	}
	if github == nil {
		return nil, fmt.Errorf("github client is required")
	}

	return &ProfileController{
		git:       git,
		github:    github,
		repoLimit: 5,
	}, nil
}

func (c *ProfileController) SetRepoLimit(limit int) {
	if limit > 0 {
		c.repoLimit = limit
	}
}

func (c *ProfileController) ShowProfile(username string) (string, error) {
	if c.git == nil || c.github == nil {
		return "", fmt.Errorf("profile controller dependencies are not configured")
	}

	resolvedUsername := strings.TrimSpace(username)
	if resolvedUsername == "" {
		gitUsername, err := c.git.GetGithubUsername()
		if err != nil {
			return "", fmt.Errorf("failed to resolve username from git config: %w", err)
		}
		resolvedUsername = strings.TrimSpace(gitUsername)
	}

	if resolvedUsername == "" {
		return "", fmt.Errorf("username could not be detected or provided")
	}

	user, err := c.github.GetUser(resolvedUsername)
	if err != nil {
		return "", fmt.Errorf("failed to fetch GitHub user %q: %w", resolvedUsername, err)
	}

	repos, err := c.github.GetUserRepos(resolvedUsername, c.repoLimit)
	if err != nil {
		return "", fmt.Errorf("failed to fetch repositories for %q: %w", resolvedUsername, err)
	}

	return ui.RenderProfile(user, repos), nil
}
