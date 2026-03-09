package git

import (
	"fmt"
	"strconv"
	"strings"
)

type LogEntry struct {
	Hash    string
	Time    string
	Author  string
	Message string
}

type ReflogEntry struct {
	Ref     string
	Hash    string
	Time    string
	Message string
}

func (g *GitCore) ShowLog(count int) ([]LogEntry, error) {
	count = normalizeCount(count)

	stdout, stderr, err := g.RunCommand(
		"log",
		"-n",
		strconv.Itoa(count),
		"--pretty=format:%h|%ad|%an|%s",
		"--date=format:%Y-%m-%d %H:%M",
	)
	if err != nil {
		if stderr != "" {
			return nil, fmt.Errorf("git log failed: %s", stderr)
		}
		return nil, fmt.Errorf("git log failed: %w", err)
	}

	if strings.TrimSpace(stdout) == "" {
		return []LogEntry{}, nil
	}

	lines := strings.Split(stdout, "\n")
	entries := make([]LogEntry, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, "|", 4)
		if len(parts) != 4 {
			return nil, fmt.Errorf("failed to parse git log line: %q", line)
		}

		entries = append(entries, LogEntry{
			Hash:    strings.TrimSpace(parts[0]),
			Time:    strings.TrimSpace(parts[1]),
			Author:  strings.TrimSpace(parts[2]),
			Message: strings.TrimSpace(parts[3]),
		})
	}

	return entries, nil
}

func (g *GitCore) ShowReflog(count int) ([]ReflogEntry, error) {
	count = normalizeCount(count)

	stdout, stderr, err := g.RunCommand(
		"reflog",
		"-n",
		strconv.Itoa(count),
		"--pretty=format:%h|%ad|%gs",
		"--date=format:%Y-%m-%d %H:%M",
	)
	if err != nil {
		if stderr != "" {
			return nil, fmt.Errorf("git reflog failed: %s", stderr)
		}
		return nil, fmt.Errorf("git reflog failed: %w", err)
	}

	if strings.TrimSpace(stdout) == "" {
		return []ReflogEntry{}, nil
	}

	lines := strings.Split(stdout, "\n")
	entries := make([]ReflogEntry, 0, len(lines))
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, "|", 3)
		if len(parts) != 3 {
			return nil, fmt.Errorf("failed to parse git reflog line: %q", line)
		}

		entries = append(entries, ReflogEntry{
			Ref:     fmt.Sprintf("HEAD@{%d}", i),
			Hash:    strings.TrimSpace(parts[0]),
			Time:    strings.TrimSpace(parts[1]),
			Message: strings.TrimSpace(parts[2]),
		})
	}

	return entries, nil
}

func normalizeCount(count int) int {
	if count <= 0 {
		return 10
	}
	return count
}
