package main

import (
	"fmt"
	"os"

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
		fmt.Fprintln(os.Stderr, "failed to initialize:", err)
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
				return err
			}
			_, err = tea.NewProgram(model).Run()
			return err
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
		Short: "Create or switch branches",
		RunE: func(cmd *cobra.Command, args []string) error {
			if createName != "" {
				return ui.WithSpinner(
					fmt.Sprintf("Creating branch %s...", createName),
					func() error {
						return d.actions.CreateBranch(createName)
					},
				)
			}
			if switchName != "" {
				return ui.WithSpinner(
					fmt.Sprintf("Switching to %s...", switchName),
					func() error {
						return d.actions.SwitchBranch(switchName)
					},
				)
			}
			fmt.Fprintln(cmd.OutOrStdout(), ui.RenderHelp())
			return nil
		},
	}

	cmd.Flags().StringVarP(&createName, "create", "c", "", "Create new branch and push to remote")
	cmd.Flags().StringVarP(&switchName, "switch", "s", "", "Switch to existing branch")

	return cmd
}

func newMeCmd(d *deps) *cobra.Command {
	return &cobra.Command{
		Use:   "me",
		Short: "Show your GitHub profile",
		RunE: func(cmd *cobra.Command, args []string) error {
			profile, err := d.profileCtrl.ShowProfile("")
			if err != nil {
				return err
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
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			profile, err := d.profileCtrl.ShowProfile(args[0])
			if err != nil {
				return err
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

	if themeName, err := core.GetConfig("gli.theme"); err == nil {
		_ = ui.SetTheme(themeName)
	}

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
