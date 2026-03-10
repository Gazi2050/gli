package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/Gazi2050/gli/internal/api"
	"github.com/Gazi2050/gli/internal/controllers"
	gitpkg "github.com/Gazi2050/gli/internal/git"
	"github.com/Gazi2050/gli/internal/ui"
)

var version = "dev-local"

type gitManagerAdapter struct {
	core    *gitpkg.GitCore
	actions *gitpkg.GitActions
}

func (g *gitManagerAdapter) RunCommand(args ...string) (string, string, error) {
	return g.core.RunCommand(args...)
}

func (g *gitManagerAdapter) GetStagedDiff() (string, error) {
	return g.core.GetStagedDiff()
}

func (g *gitManagerAdapter) GetGithubUsername() (string, error) {
	return g.core.GetGithubUsername()
}

func (g *gitManagerAdapter) GetRepoName() (string, error) {
	return g.core.GetRepoName()
}

func (g *gitManagerAdapter) CommitAndPush(message string, noVerify bool) error {
	return g.actions.CommitAndPush(message, noVerify)
}

type appDeps struct {
	core        *gitpkg.GitCore
	actions     *gitpkg.GitActions
	history     *gitpkg.GitHistory
	commitCtrl  *controllers.CommitController
	profileCtrl *controllers.ProfileController
}

func main() {
	rootCmd, err := newRootCmd()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to initialize app:", err)
		os.Exit(1)
	}

	rootCmd.SetArgs(normalizeLegacyArgs(os.Args[1:]))

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newRootCmd() (*cobra.Command, error) {
	deps, err := buildDeps()
	if err != nil {
		return nil, err
	}

	var (
		showVersion   bool
		commitMessage string
		aiCommit      bool
		showLog       bool
		showReflog    bool
		resetMode     string
		switchBranch  string
		localBranch   bool
		remoteBranch  bool
		changeTime    string
		changeAuthor  bool
		changeMessage bool
		noVerify      bool
	)

	cmd := &cobra.Command{
		Use:           "gli",
		Short:         "gli - Modern Git Wrapper",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if showVersion {
				fmt.Fprintln(cmd.OutOrStdout(), version)
				return nil
			}

			if len(args) > 0 {
				switch args[0] {
				case "version":
					fmt.Fprintln(cmd.OutOrStdout(), version)
					return nil
				case "me":
					if len(args) > 1 {
						return fmt.Errorf("me does not take arguments")
					}
					profile, err := deps.profileCtrl.ShowProfile("")
					if err != nil {
						return err
					}
					fmt.Fprintln(cmd.OutOrStdout(), profile)
					return nil
				case "profile":
					username := ""
					if len(args) > 2 {
						return fmt.Errorf("usage: gli profile [username]")
					}
					if len(args) == 2 {
						username = args[1]
					}
					profile, err := deps.profileCtrl.ShowProfile(username)
					if err != nil {
						return err
					}
					fmt.Fprintln(cmd.OutOrStdout(), profile)
					return nil
				default:
					if cmd.Flags().Changed("commit") && strings.TrimSpace(commitMessage) == "prompt" {
						commitMessage = strings.Join(args, " ")
					} else {
						return fmt.Errorf("unknown command %q", args[0])
					}
				}
			}

			if cmd.Flags().Changed("commit") {
				if commitMessage == "prompt" {
					message, err := promptLine("Enter commit message: ")
					if err != nil {
						return err
					}
					commitMessage = message
				}
				return deps.commitCtrl.RunManualCommit(noVerify, commitMessage)
			}

			if aiCommit {
				model, err := deps.commitCtrl.RunAICommit(noVerify, nil)
				if err != nil {
					return err
				}
				_, err = tea.NewProgram(model).Run()
				return err
			}

			if showLog {
				entries, err := deps.core.ShowLog(20)
				if err != nil {
					return err
				}
				printLog(entries)
				return nil
			}

			if showReflog {
				entries, err := deps.core.ShowReflog(20)
				if err != nil {
					return err
				}
				printReflog(entries)
				return nil
			}

			if resetMode != "" {
				return deps.actions.ResetCommit(resetMode)
			}

			if switchBranch != "" {
				pushToRemote := true
				if localBranch {
					pushToRemote = false
				} else if remoteBranch {
					pushToRemote = true
				}
				return deps.actions.SwitchBranch(switchBranch, pushToRemote)
			}

			if cmd.Flags().Changed("changeTime") {
				if strings.TrimSpace(changeTime) != "" && changeTime != "prompt" {
					return deps.history.ChangeCommitTime("single", "", strings.TrimSpace(changeTime))
				}
				model := ui.NewHistoryModel(deps.history, ui.HistoryOperationTime)
				_, err := tea.NewProgram(model).Run()
				return err
			}

			if changeAuthor {
				model := ui.NewHistoryModel(deps.history, ui.HistoryOperationAuthor)
				_, err := tea.NewProgram(model).Run()
				return err
			}

			if changeMessage {
				model := ui.NewHistoryModel(deps.history, ui.HistoryOperationMessage)
				_, err := tea.NewProgram(model).Run()
				return err
			}

			fmt.Fprintln(cmd.OutOrStdout(), ui.RenderHelp())
			return nil
		},
	}

	cmd.Flags().BoolVarP(&showVersion, "version", "v", false, "Show version")
	cmd.Flags().StringVarP(&commitMessage, "commit", "c", "", "Commit and push")
	cmd.Flags().Lookup("commit").NoOptDefVal = "prompt"
	cmd.Flags().BoolVar(&aiCommit, "ai-commit", false, "AI-powered commit")
	cmd.Flags().BoolVarP(&showLog, "log", "l", false, "View git log")
	cmd.Flags().BoolVar(&showReflog, "reflog", false, "View git reflog")
	cmd.Flags().StringVar(&resetMode, "reset", "", "Reset last commit: soft or hard")
	cmd.Flags().StringVarP(&switchBranch, "switch", "s", "", "Create/switch branch")
	cmd.Flags().BoolVar(&localBranch, "local-branch", false, "Branch: create local only")
	cmd.Flags().BoolVar(&remoteBranch, "remote-branch", false, "Branch: push to remote")
	cmd.Flags().StringVar(&changeTime, "changeTime", "", "History: change commit time")
	cmd.Flags().Lookup("changeTime").NoOptDefVal = "prompt"
	cmd.Flags().BoolVar(&changeAuthor, "changeAuthor", false, "History: change commit author")
	cmd.Flags().BoolVar(&changeMessage, "changeMessage", false, "History: change commit message")
	cmd.Flags().BoolVar(&noVerify, "no-verify", false, "Skip git hooks")

	return cmd, nil
}

func buildDeps() (*appDeps, error) {
	core := gitpkg.NewGitCore("")

	// Optional: theme selection via `git config gli.theme <name>`.
	if themeName, err := core.GetConfig("gli.theme"); err == nil {
		_ = ui.SetTheme(themeName)
	}

	actions, err := gitpkg.NewGitActions(core)
	if err != nil {
		return nil, err
	}

	history, err := gitpkg.NewGitHistory(core)
	if err != nil {
		return nil, err
	}

	aiService, err := api.NewAIService(nil)
	if err != nil {
		return nil, err
	}

	githubClient := api.NewGitHubClient(nil, "")
	gitAdapter := &gitManagerAdapter{core: core, actions: actions}

	commitCtrl, err := controllers.NewCommitController(gitAdapter, aiService)
	if err != nil {
		return nil, err
	}

	profileCtrl, err := controllers.NewProfileController(core, githubClient, ui.RenderProfile)
	if err != nil {
		return nil, err
	}

	return &appDeps{
		core:        core,
		actions:     actions,
		history:     history,
		commitCtrl:  commitCtrl,
		profileCtrl: profileCtrl,
	}, nil
}

func promptLine(prompt string) (string, error) {
	fmt.Fprint(os.Stdout, prompt)
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return "", err
	}
	return strings.TrimSpace(line), nil
}

func printLog(entries []gitpkg.LogEntry) {
	for _, entry := range entries {
		fmt.Printf("%s  %s  %s  %s\n", entry.Hash, entry.Time, entry.Author, entry.Message)
	}
}

func printReflog(entries []gitpkg.ReflogEntry) {
	for _, entry := range entries {
		fmt.Printf("%s  %s  %s  %s\n", entry.Ref, entry.Hash, entry.Time, entry.Message)
	}
}

func normalizeLegacyArgs(args []string) []string {
	flagMap := map[string]string{
		"-ac": "--ai-commit",
		"-rl": "--reflog",
		"-rs": "--reset",
		"-lb": "--local-branch",
		"-rb": "--remote-branch",
		"-ct": "--changeTime",
		"-ca": "--changeAuthor",
		"-cm": "--changeMessage",
		"-nv": "--no-verify",
	}

	out := make([]string, 0, len(args))
	for _, arg := range args {
		if long, ok := flagMap[arg]; ok {
			out = append(out, long)
			continue
		}

		replaced := false
		for short, long := range flagMap {
			prefix := short + "="
			if strings.HasPrefix(arg, prefix) {
				out = append(out, long+arg[len(short):])
				replaced = true
				break
			}
		}
		if replaced {
			continue
		}

		out = append(out, arg)
	}

	return out
}
