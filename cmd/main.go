package main

import (
	"fmt"
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

func main() {
	rootCmd, err := newRootCmd()
	if err != nil {
		fmt.Fprintln(os.Stderr, ui.ErrorCard("Initialization Error", err.Error()))
		os.Exit(1)
	}

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

	rootCmd := &cobra.Command{
		Use:           "gli",
		Short:         "gli - Modern Git Wrapper",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintln(cmd.OutOrStdout(), ui.RenderHelp())
			return nil
		},
	}

	rootCmd.AddCommand(
		newCommitCmd(deps),
		newBranchCmd(deps),
		newStatusCmd(deps),
		newLogCmd(deps),
		newReflogCmd(deps),
		newMeCmd(deps),
		newProfileCmd(deps),
		newVersionCmd(),
	)

	return rootCmd, nil
}

func newCommitCmd(d *deps) *cobra.Command {
	var noVerify bool

	cmd := &cobra.Command{
		Use:   "commit",
		Short: "Stage all, generate AI commit message, commit and push",
		RunE: func(cmd *cobra.Command, args []string) error {
			model, err := d.commitCtrl.RunAICommit(noVerify)
			if err != nil {
				return fmt.Errorf("%s", ui.ErrorCard("Commit Error", friendlyGitError(err)))
			}
			if _, err = tea.NewProgram(model).Run(); err != nil {
				return fmt.Errorf("%s", ui.ErrorCard("Commit Error", friendlyError(err)))
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&noVerify, "no-verify", false, "Skip git hooks")

	return cmd
}

func newBranchCmd(d *deps) *cobra.Command {
	var createName string
	var switchName string

	cmd := &cobra.Command{
		Use:   "branch",
		Short: "List, create, or switch branches",
		RunE: func(cmd *cobra.Command, args []string) error {
			out := cmd.OutOrStdout()
			if createName != "" {
				err := ui.WithSpinner(
					fmt.Sprintf("Creating branch %s...", createName),
					func() error {
						return d.actions.CreateBranch(createName)
					},
				)
				if err != nil {
					return fmt.Errorf("%s", ui.ErrorCard("Branch Error", friendlyGitError(err)))
				}
				fmt.Fprintln(out, ui.SuccessCard("Done", fmt.Sprintf("Branch %s created and pushed to remote.", createName)))
				return nil
			}
			if switchName != "" {
				err := ui.WithSpinner(
					fmt.Sprintf("Switching to %s...", switchName),
					func() error {
						return d.actions.SwitchBranch(switchName)
					},
				)
				if err != nil {
					return fmt.Errorf("%s", ui.ErrorCard("Branch Error", friendlyGitError(err)))
				}
				fmt.Fprintln(out, ui.SuccessCard("Done", fmt.Sprintf("Switched to branch %s.", switchName)))
				return nil
			}
			branches, err := d.core.GetBranches()
			if err != nil {
				return fmt.Errorf("%s", ui.ErrorCard("Branch Error", friendlyGitError(err)))
			}
			output := ui.RenderBranches(branches)
			if output != "" {
				fmt.Fprintln(out, output)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&createName, "create", "c", "", "Create new branch and push to remote")
	cmd.Flags().StringVarP(&switchName, "switch", "s", "", "Switch to existing branch")

	return cmd
}

func newStatusCmd(d *deps) *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show current branch status",
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := d.core.GetStatusData()
			if err != nil {
				return fmt.Errorf("%s", ui.ErrorCard("Status Error", "Not a git repository. Run this inside a git project."))
			}
			fmt.Fprintln(cmd.OutOrStdout(), ui.RenderStatus(data))
			return nil
		},
	}
}

func newLogCmd(d *deps) *cobra.Command {
	var limit int

	cmd := &cobra.Command{
		Use:   "log",
		Short: "Show commit history",
		RunE: func(cmd *cobra.Command, args []string) error {
			entries, err := d.core.GetLogData(limit)
			if err != nil {
				return fmt.Errorf("%s", ui.ErrorCard("Log Error", friendlyGitError(err)))
			}
			output := ui.RenderLog(entries, "gli log", "Commit")
			if output != "" {
				fmt.Fprintln(cmd.OutOrStdout(), output)
			}
			return nil
		},
	}

	cmd.Flags().IntVarP(&limit, "limit", "n", 50, "Number of commits to show")

	return cmd
}

func newReflogCmd(d *deps) *cobra.Command {
	var limit int

	cmd := &cobra.Command{
		Use:   "reflog",
		Short: "Show reference log",
		RunE: func(cmd *cobra.Command, args []string) error {
			entries, err := d.core.GetReflogData(limit)
			if err != nil {
				return fmt.Errorf("%s", ui.ErrorCard("Reflog Error", friendlyGitError(err)))
			}
			output := ui.RenderLog(entries, "gli reflog", "Action")
			if output != "" {
				fmt.Fprintln(cmd.OutOrStdout(), output)
			}
			return nil
		},
	}

	cmd.Flags().IntVarP(&limit, "limit", "n", 50, "Number of entries to show")

	return cmd
}

func newMeCmd(d *deps) *cobra.Command {
	return &cobra.Command{
		Use:   "me",
		Short: "Show your GitHub profile",
		RunE: func(cmd *cobra.Command, args []string) error {
			profile, err := d.profileCtrl.ShowProfile("")
			if err != nil {
				return fmt.Errorf("%s", ui.ErrorCard("Profile Error", friendlyProfileError(err)))
			}
			fmt.Fprintln(cmd.OutOrStdout(), profile)
			return nil
		},
	}
}

func newProfileCmd(d *deps) *cobra.Command {
	return &cobra.Command{
		Use:   "profile <username>",
		Short: "Show a GitHub user's profile",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("%s", ui.ErrorCard("Missing Username", "Usage: gli profile <username>\nExample: gli profile torvalds"))
			}
			if len(args) > 1 {
				return fmt.Errorf("%s", ui.ErrorCard("Too Many Arguments", "Usage: gli profile <username>\nExample: gli profile torvalds"))
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			username := strings.TrimSpace(args[0])
			if username == "" {
				return fmt.Errorf("%s", ui.ErrorCard("Invalid Username", "Username cannot be empty."))
			}
			profile, err := d.profileCtrl.ShowProfile(username)
			if err != nil {
				return fmt.Errorf("%s", ui.ErrorCard("Profile Error", friendlyProfileError(err)))
			}
			fmt.Fprintln(cmd.OutOrStdout(), profile)
			return nil
		},
	}
}

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show gli version",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintln(cmd.OutOrStdout(), version)
			return nil
		},
	}
}

type deps struct {
	core        *gitpkg.GitCore
	actions     *gitpkg.GitActions
	commitCtrl  *controllers.CommitController
	profileCtrl *controllers.ProfileController
}

func buildDeps() (*deps, error) {
	core := gitpkg.NewGitCore("")

	actions, err := gitpkg.NewGitActions(core)
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

	return &deps{
		core:        core,
		actions:     actions,
		commitCtrl:  commitCtrl,
		profileCtrl: profileCtrl,
	}, nil
}

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

func (g *gitManagerAdapter) AddAll() error {
	return g.actions.AddAll()
}

func (g *gitManagerAdapter) CommitAndPush(message string, noVerify bool) error {
	return g.actions.CommitAndPush(message, noVerify)
}

func (g *gitManagerAdapter) HasUpstream() (bool, error) {
	return g.actions.HasUpstream()
}

func friendlyGitError(err error) string {
	msg := strings.ToLower(err.Error())
	switch {
	case strings.Contains(msg, "not a git repository"):
		return "Not a git repository. Run this inside a git project."
	case strings.Contains(msg, "no staged diff") || strings.Contains(msg, "no staged changes"):
		return "No changes to commit. Make some changes first."
	case strings.Contains(msg, "already exists") || strings.Contains(msg, "already exists and"):
		return "This branch already exists."
	case strings.Contains(msg, "did not match") || strings.Contains(msg, "not found"):
		return "Branch not found. Check the name and try again."
	case strings.Contains(msg, "network") || strings.Contains(msg, "timeout") || strings.Contains(msg, "connection refused"):
		return "Network error. Check your internet connection."
	case strings.Contains(msg, "permission denied"):
		return "Permission denied. Check your access rights."
	case strings.Contains(msg, "push") && strings.Contains(msg, "failed"):
		return "Push failed. Check your remote and permissions."
	default:
		return err.Error()
	}
}

func friendlyProfileError(err error) string {
	msg := strings.ToLower(err.Error())
	switch {
	case strings.Contains(msg, "status 404"):
		return "User not found on GitHub."
	case strings.Contains(msg, "status 403"):
		return "GitHub API rate limit exceeded. Try again later."
	case strings.Contains(msg, "network") || strings.Contains(msg, "timeout") || strings.Contains(msg, "connection refused"):
		return "Network error. Check your internet connection."
	case strings.Contains(msg, "username could not be detected"):
		return "Could not detect your GitHub username. Set it with:\n  git config --global github.user <username>"
	default:
		return err.Error()
	}
}

func friendlyError(err error) string {
	msg := err.Error()
	if strings.Contains(strings.ToLower(msg), "failed to initialize") {
		return "Something went wrong. Please try again."
	}
	return msg
}
