package git

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

var ErrNoStagedDiff = errors.New("no staged diff found")

type GitCore struct {
	WorkDir string
}

func NewGitCore(workDir string) *GitCore {
	return &GitCore{WorkDir: workDir}
}

func (g *GitCore) RunCommand(args ...string) (stdout string, stderr string, err error) {
	cmd := exec.Command("git", args...)
	if g.WorkDir != "" {
		cmd.Dir = g.WorkDir
	}

	var outBuf bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err = cmd.Run()
	return strings.TrimSpace(outBuf.String()), strings.TrimSpace(errBuf.String()), err
}

func (g *GitCore) GetConfig(key string) (string, error) {
	stdout, stderr, err := g.RunCommand("config", "--get", key)
	if err != nil {
		if stderr != "" {
			return "", fmt.Errorf("git config --get %s failed: %s", key, stderr)
		}
		return "", fmt.Errorf("git config --get %s failed: %w", key, err)
	}
	return stdout, nil
}

func (g *GitCore) GetGithubUsername() (string, error) {
	username, err := g.GetConfig("github.user")
	if err == nil && username != "" {
		return username, nil
	}

	username, err = g.GetConfig("user.name")
	if err == nil && username != "" {
		return username, nil
	}

	return "unknown-user", nil
}

func (g *GitCore) GetCurrentBranch() (string, error) {
	stdout, stderr, err := g.RunCommand("rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		if stderr != "" {
			return "", fmt.Errorf("failed to get current branch: %s", stderr)
		}
		return "", fmt.Errorf("failed to get current branch: %w", err)
	}
	return stdout, nil
}

func (g *GitCore) GetStagedDiff() (string, error) {
	stdout, stderr, err := g.RunCommand("diff", "--staged")
	if err != nil {
		if stderr != "" {
			return "", fmt.Errorf("failed to get staged diff: %s", stderr)
		}
		return "", fmt.Errorf("failed to get staged diff: %w", err)
	}

	if strings.TrimSpace(stdout) == "" {
		return "", ErrNoStagedDiff
	}

	return stdout, nil
}

func (g *GitCore) GetRepoName() (string, error) {
	remoteURL, err := g.GetConfig("remote.origin.url")
	if err != nil {
		return "", fmt.Errorf("failed to get remote origin url: %w", err)
	}

	normalized := strings.TrimSuffix(strings.TrimSpace(remoteURL), ".git")
	normalized = strings.TrimSuffix(normalized, "/")
	parts := strings.Split(normalized, "/")
	if len(parts) == 0 {
		return "", fmt.Errorf("unable to parse repository name from remote url: %s", remoteURL)
	}

	repoName := strings.TrimSpace(parts[len(parts)-1])
	if repoName == "" {
		return "", fmt.Errorf("unable to parse repository name from remote url: %s", remoteURL)
	}

	return repoName, nil
}
