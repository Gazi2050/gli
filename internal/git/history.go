package git

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

type GitHistory struct {
	Core *GitCore
}

func NewGitHistory(core *GitCore) (*GitHistory, error) {
	if core == nil {
		return nil, errors.New("git core is required")
	}
	return &GitHistory{Core: core}, nil
}

func (g *GitHistory) ChangeCommitTime(scope, targetHash, dateTime string) error {
	scope, targetHash, err := normalizeScope(scope, targetHash)
	if err != nil {
		return err
	}

	dateTime = strings.TrimSpace(dateTime)
	if dateTime == "" {
		return errors.New("dateTime cannot be empty")
	}

	if scope == "single" {
		env := map[string]string{
			"GIT_AUTHOR_DATE":    dateTime,
			"GIT_COMMITTER_DATE": dateTime,
		}
		if _, stderr, err := g.runCommandWithEnv(env, "commit", "--amend", "--no-edit", "--date", dateTime); err != nil {
			if stderr != "" {
				return fmt.Errorf("failed to amend commit time: %s", stderr)
			}
			return fmt.Errorf("failed to amend commit time: %w", err)
		}
		return nil
	}

	filterScript := buildTimeFilterScript(scope, targetHash, dateTime)
	if _, stderr, err := g.Core.RunCommand("filter-branch", "-f", "--env-filter", filterScript, "--", "HEAD"); err != nil {
		if stderr != "" {
			return fmt.Errorf("failed to rewrite commit time history: %s", stderr)
		}
		return fmt.Errorf("failed to rewrite commit time history: %w", err)
	}

	return nil
}

func (g *GitHistory) ChangeCommitAuthor(scope, targetHash, name, email string) error {
	scope, targetHash, err := normalizeScope(scope, targetHash)
	if err != nil {
		return err
	}

	name = strings.TrimSpace(name)
	email = strings.TrimSpace(email)
	if name == "" {
		return errors.New("author name cannot be empty")
	}
	if email == "" {
		return errors.New("author email cannot be empty")
	}

	if scope == "single" {
		author := fmt.Sprintf("%s <%s>", name, email)
		if _, stderr, err := g.Core.RunCommand("commit", "--amend", "--no-edit", "--author", author); err != nil {
			if stderr != "" {
				return fmt.Errorf("failed to amend commit author: %s", stderr)
			}
			return fmt.Errorf("failed to amend commit author: %w", err)
		}
		return nil
	}

	filterScript := buildAuthorFilterScript(scope, targetHash, name, email)
	if _, stderr, err := g.Core.RunCommand("filter-branch", "-f", "--env-filter", filterScript, "--", "HEAD"); err != nil {
		if stderr != "" {
			return fmt.Errorf("failed to rewrite commit author history: %s", stderr)
		}
		return fmt.Errorf("failed to rewrite commit author history: %w", err)
	}

	return nil
}

func (g *GitHistory) ChangeCommitMessage(newMessage string) error {
	newMessage = strings.TrimSpace(newMessage)
	if newMessage == "" {
		return errors.New("commit message cannot be empty")
	}

	if _, stderr, err := g.Core.RunCommand("commit", "--amend", "-m", newMessage); err != nil {
		if stderr != "" {
			return fmt.Errorf("failed to amend commit message: %s", stderr)
		}
		return fmt.Errorf("failed to amend commit message: %w", err)
	}

	return nil
}

func normalizeScope(scope, targetHash string) (string, string, error) {
	scope = strings.TrimSpace(scope)
	targetHash = strings.TrimSpace(targetHash)

	switch scope {
	case "single", "all":
		return scope, targetHash, nil
	case "specific":
		if targetHash == "" {
			return "", "", errors.New("target hash is required when scope is specific")
		}
		return scope, targetHash, nil
	default:
		return "", "", fmt.Errorf("invalid scope %q: expected single, specific, or all", scope)
	}
}

func buildTimeFilterScript(scope, targetHash, dateTime string) string {
	date := singleQuote(dateTime)
	setScript := fmt.Sprintf("export GIT_AUTHOR_DATE=%s; export GIT_COMMITTER_DATE=%s", date, date)
	if scope == "specific" {
		return fmt.Sprintf("case \"$GIT_COMMIT\" in %s*) %s ;; esac", targetHash, setScript)
	}
	return setScript
}

func buildAuthorFilterScript(scope, targetHash, name, email string) string {
	nameQ := singleQuote(name)
	emailQ := singleQuote(email)
	setScript := fmt.Sprintf(
		"export GIT_AUTHOR_NAME=%s; export GIT_AUTHOR_EMAIL=%s; export GIT_COMMITTER_NAME=%s; export GIT_COMMITTER_EMAIL=%s",
		nameQ, emailQ, nameQ, emailQ,
	)
	if scope == "specific" {
		return fmt.Sprintf("case \"$GIT_COMMIT\" in %s*) %s ;; esac", targetHash, setScript)
	}
	return setScript
}

func singleQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "'\"'\"'") + "'"
}

func (g *GitHistory) runCommandWithEnv(env map[string]string, args ...string) (stdout string, stderr string, err error) {
	cmd := exec.Command("git", args...)
	if g.Core.WorkDir != "" {
		cmd.Dir = g.Core.WorkDir
	}

	cmd.Env = append(cmd.Environ(), mapToEnvList(env)...)

	var outBuf bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err = cmd.Run()
	return strings.TrimSpace(outBuf.String()), strings.TrimSpace(errBuf.String()), err
}

func mapToEnvList(env map[string]string) []string {
	if len(env) == 0 {
		return nil
	}

	out := make([]string, 0, len(env))
	for k, v := range env {
		out = append(out, fmt.Sprintf("%s=%s", k, v))
	}
	return out
}
