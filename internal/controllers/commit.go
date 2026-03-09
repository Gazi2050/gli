package controllers

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/Gazi2050/gli/internal/ui"
)

type CommitController struct {
	git ui.GitManager
	ai  ui.AIService
}

func NewCommitController(git ui.GitManager, ai ui.AIService) (*CommitController, error) {
	if git == nil {
		return nil, fmt.Errorf("git manager is required")
	}
	if ai == nil {
		return nil, fmt.Errorf("ai service is required")
	}

	return &CommitController{git: git, ai: ai}, nil
}

func (c *CommitController) RunManualCommit(noVerify bool, message string) error {
	message = strings.TrimSpace(message)
	if message == "" {
		return fmt.Errorf("commit message cannot be empty")
	}

	return c.CommitAndPush(message, noVerify)
}

func (c *CommitController) RunAICommit(noVerify bool, parent tea.Model) (tea.Model, error) {
	_ = parent

	if c.git == nil || c.ai == nil {
		return nil, fmt.Errorf("commit controller dependencies are not configured")
	}

	model := ui.NewCommitModel(c.git, c.ai, noVerify, "")
	return model, nil
}

func (c *CommitController) CommitAndPush(message string, noVerify bool) error {
	if c.git == nil {
		return fmt.Errorf("git manager is not configured")
	}

	message = strings.TrimSpace(message)
	if message == "" {
		return fmt.Errorf("commit message cannot be empty")
	}

	if err := c.git.CommitAndPush(message, noVerify); err != nil {
		return fmt.Errorf("commit and push failed: %w", err)
	}

	return nil
}
