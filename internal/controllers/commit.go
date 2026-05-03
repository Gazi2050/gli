package controllers

import (
	"fmt"

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

func (c *CommitController) RunAICommit(noVerify bool) (tea.Model, error) {
	if c.git == nil || c.ai == nil {
		return nil, fmt.Errorf("commit controller dependencies are not configured")
	}

	model := ui.NewCommitModel(c.git, c.ai, noVerify, "")
	return model, nil
}
