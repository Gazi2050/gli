package git

import (
	"errors"
	"fmt"
	"strings"
)

type GitActions struct {
	Core *GitCore
}

func NewGitActions(core *GitCore) (*GitActions, error) {
	if core == nil {
		return nil, errors.New("git core is required")
	}
	return &GitActions{Core: core}, nil
}

func (g *GitActions) AddAll() error {
	if _, stderr, err := g.Core.RunCommand("add", "."); err != nil {
		if stderr != "" {
			return fmt.Errorf("git add failed: %s", stderr)
		}
		return fmt.Errorf("git add failed: %w", err)
	}
	return nil
}

func (g *GitActions) CommitAndPush(message string, noVerify bool) error {
	if strings.TrimSpace(message) == "" {
		return errors.New("commit message cannot be empty")
	}

	commitArgs := []string{"commit", "-m", message}
	if noVerify {
		commitArgs = append(commitArgs, "--no-verify")
	}
	if _, stderr, err := g.Core.RunCommand(commitArgs...); err != nil {
		if stderr != "" {
			return fmt.Errorf("git commit failed: %s", stderr)
		}
		return fmt.Errorf("git commit failed: %w", err)
	}

	return g.Push(noVerify)
}

func (g *GitActions) Push(noVerify bool) error {
	hasUpstream, err := g.hasUpstream()
	if err != nil {
		return err
	}

	if hasUpstream {
		pushArgs := []string{"push"}
		if noVerify {
			pushArgs = append(pushArgs, "--no-verify")
		}
		if _, stderr, err := g.Core.RunCommand(pushArgs...); err != nil {
			if stderr != "" {
				return fmt.Errorf("git push failed: %s", stderr)
			}
			return fmt.Errorf("git push failed: %w", err)
		}
		return nil
	}

	branch, err := g.Core.GetCurrentBranch()
	if err != nil {
		return fmt.Errorf("failed to detect branch for upstream push: %w", err)
	}

	pushArgs := []string{"push"}
	if noVerify {
		pushArgs = append(pushArgs, "--no-verify")
	}
	pushArgs = append(pushArgs, "--set-upstream", "origin", branch)
	if _, stderr, err := g.Core.RunCommand(pushArgs...); err != nil {
		if stderr != "" {
			return fmt.Errorf("git push --set-upstream failed: %s", stderr)
		}
		return fmt.Errorf("git push --set-upstream failed: %w", err)
	}

	return nil
}

func (g *GitActions) CreateBranch(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("branch name cannot be empty")
	}
	if _, stderr, err := g.Core.RunCommand("switch", "-c", name); err != nil {
		if stderr != "" {
			return fmt.Errorf("git switch -c %s failed: %s", name, stderr)
		}
		return fmt.Errorf("git switch -c %s failed: %w", name, err)
	}
	if _, stderr, err := g.Core.RunCommand("push", "-u", "origin", name); err != nil {
		if stderr != "" {
			return fmt.Errorf("git push -u origin %s failed: %s", name, stderr)
		}
		return fmt.Errorf("git push -u origin %s failed: %w", name, err)
	}
	return nil
}

func (g *GitActions) SwitchBranch(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("branch name cannot be empty")
	}
	if _, stderr, err := g.Core.RunCommand("switch", name); err != nil {
		if stderr != "" {
			return fmt.Errorf("git switch %s failed: %s", name, stderr)
		}
		return fmt.Errorf("git switch %s failed: %w", name, err)
	}
	return nil
}

func (g *GitActions) HasUpstream() (bool, error) {
	return g.hasUpstream()
}

func (g *GitActions) hasUpstream() (bool, error) {
	_, stderr, err := g.Core.RunCommand("rev-parse", "--symbolic-full-name", "@{u}")
	if err == nil {
		return true, nil
	}

	low := strings.ToLower(stderr)
	if strings.Contains(low, "no upstream configured") ||
		strings.Contains(low, "does not have any upstream branch") {
		return false, nil
	}

	if stderr != "" {
		return false, fmt.Errorf("failed checking upstream branch: %s", stderr)
	}
	return false, fmt.Errorf("failed checking upstream branch: %w", err)
}
