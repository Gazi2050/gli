package git

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
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

type FileStat struct {
	Path        string
	Additions   int
	Deletions   int
}

type BranchEntry struct {
	Name    string
	Hash    string
	Current bool
}

func (g *GitCore) GetBranches() ([]BranchEntry, error) {
	out, _, err := g.RunCommand("branch", "-v", "--format=%(refname:short)|%(objectname:short)")
	if err != nil {
		return nil, fmt.Errorf("failed to get branches: %w", err)
	}

	currentBranch, _ := g.GetCurrentBranch()

	var branches []BranchEntry
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "|", 2)
		if len(parts) != 2 {
			continue
		}
		name := strings.TrimPrefix(parts[0], "* ")
		branches = append(branches, BranchEntry{
			Name:    name,
			Hash:    parts[1],
			Current: name == currentBranch,
		})
	}
	return branches, nil
}

type LogEntry struct {
	Hash    string
	Author  string
	Commit  string
	Date    string
}

func (g *GitCore) GetLogData(limit int) ([]LogEntry, error) {
	if limit <= 0 {
		limit = 50
	}
	out, _, err := g.RunCommand("log", "--pretty=format:%h|%an|%s|%cr", "-n", strconv.Itoa(limit))
	if err != nil {
		return nil, fmt.Errorf("failed to get log: %w", err)
	}

	var entries []LogEntry
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "|", 4)
		if len(parts) != 4 {
			continue
		}
		entries = append(entries, LogEntry{
			Hash:    parts[0],
			Author:  parts[1],
			Commit:  parts[2],
			Date:    parts[3],
		})
	}
	return entries, nil
}

func (g *GitCore) GetReflogData(limit int) ([]LogEntry, error) {
	if limit <= 0 {
		limit = 50
	}
	out, _, err := g.RunCommand("reflog", "--pretty=format:%h|%an|%gs|%cr", "-n", strconv.Itoa(limit))
	if err != nil {
		return nil, fmt.Errorf("failed to get reflog: %w", err)
	}

	var entries []LogEntry
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "|", 4)
		if len(parts) != 4 {
			continue
		}
		entries = append(entries, LogEntry{
			Hash:    parts[0],
			Author:  parts[1],
			Commit:  parts[2],
			Date:    parts[3],
		})
	}
	return entries, nil
}

type StatusData struct {
	Branch         string
	Remote         string
	Ahead          int
	Behind         int
	ModifiedFiles  []FileStat
	StagedFiles    []FileStat
	NewFiles       []string
	DeletedFiles   []string
	LastCommit     string
	LastDate       string
}

func (g *GitCore) GetStatusData() (*StatusData, error) {
	branch, err := g.GetCurrentBranch()
	if err != nil {
		return nil, err
	}

	remote, _ := g.GetConfig("remote.origin.url")
	remote = strings.TrimSuffix(remote, ".git")
	if parts := strings.Split(remote, ":"); len(parts) == 2 {
		remote = parts[1]
	}
	if parts := strings.Split(remote, "/"); len(parts) >= 2 {
		remote = parts[len(parts)-2] + "/" + parts[len(parts)-1]
	}

	ahead, behind := 0, 0
	if countStr, _, err := g.RunCommand("rev-list", "--left-right", "--count", "@{u}...HEAD"); err == nil {
		parts := strings.Fields(countStr)
		if len(parts) == 2 {
			behind, _ = strconv.Atoi(parts[0])
			ahead, _ = strconv.Atoi(parts[1])
		}
	}

	var modifiedFiles, stagedFiles []FileStat
	var newFiles, deletedFiles []string

	if out, _, err := g.RunCommand("diff", "--numstat"); err == nil && out != "" {
		for _, line := range strings.Split(out, "\n") {
			parts := strings.Fields(line)
			if len(parts) < 3 {
				continue
			}
			adds, _ := strconv.Atoi(parts[0])
			dels, _ := strconv.Atoi(parts[1])
			modifiedFiles = append(modifiedFiles, FileStat{Path: parts[2], Additions: adds, Deletions: dels})
		}
	}

	if out, _, err := g.RunCommand("diff", "--cached", "--numstat"); err == nil && out != "" {
		for _, line := range strings.Split(out, "\n") {
			parts := strings.Fields(line)
			if len(parts) < 3 {
				continue
			}
			adds, _ := strconv.Atoi(parts[0])
			dels, _ := strconv.Atoi(parts[1])
			if adds == 0 && dels == 0 {
				stagedFiles = append(stagedFiles, FileStat{Path: parts[2]})
			} else {
				stagedFiles = append(stagedFiles, FileStat{Path: parts[2], Additions: adds, Deletions: dels})
			}
		}
	}

	if out, _, err := g.RunCommand("diff", "--cached", "--name-only", "--diff-filter=A"); err == nil && out != "" {
		for _, f := range strings.Split(out, "\n") {
			if f = strings.TrimSpace(f); f != "" {
				newFiles = append(newFiles, f)
			}
		}
	}

	if out, _, err := g.RunCommand("ls-files", "--others", "--exclude-standard"); err == nil && out != "" {
		for _, f := range strings.Split(out, "\n") {
			if f = strings.TrimSpace(f); f != "" {
				newFiles = append(newFiles, f)
			}
		}
	}

	if out, _, err := g.RunCommand("diff", "--name-only", "--diff-filter=D"); err == nil && out != "" {
		for _, f := range strings.Split(out, "\n") {
			if f = strings.TrimSpace(f); f != "" {
				deletedFiles = append(deletedFiles, f)
			}
		}
	}

	if out, _, err := g.RunCommand("diff", "--cached", "--name-only", "--diff-filter=D"); err == nil && out != "" {
		for _, f := range strings.Split(out, "\n") {
			found := false
			for _, existing := range deletedFiles {
				if existing == f {
					found = true
					break
				}
			}
			if !found {
				deletedFiles = append(deletedFiles, f)
			}
		}
	}

	lastCommit, _, _ := g.RunCommand("log", "-1", "--format=%s")
	lastDate, _, _ := g.RunCommand("log", "-1", "--format=%cr")

	return &StatusData{
		Branch:         branch,
		Remote:         remote,
		Ahead:          ahead,
		Behind:         behind,
		ModifiedFiles:  modifiedFiles,
		StagedFiles:    stagedFiles,
		NewFiles:       newFiles,
		DeletedFiles:   deletedFiles,
		LastCommit:     lastCommit,
		LastDate:       lastDate,
	}, nil
}
