## gli – Go + Bubble Tea Rewrite To‑Dos

1. **setup**  
   - Initialize Go module (`go mod init` and create `go.mod`)  
   - Add dependencies: `cobra`, `bubbletea`, `bubbles`, `lipgloss`

2. **git-core** (`internal/git/core.go`)  
   - Implement `RunCommand(args ...string)` using `os/exec`  
   - Implement `GetConfig(key)`, `GetCurrentBranch()`, `GetStagedDiff()`, `GetRepoName()`, `GetGithubUsername()`

3. **git-actions** (`internal/git/actions.go`)  
   - Implement `CommitAndPush(msg, noVerify)`  
   - Implement `ResetCommit(mode)`  
   - Implement `SwitchBranch(name, pushToRemote)`

4. **git-history** (`internal/git/history.go`)  
   - Implement `ChangeCommitTime()` using `git filter-branch` or `git commit --amend`  
   - Implement `ChangeCommitAuthor()` using `git filter-branch` or `git commit --amend`  
   - Implement `ChangeCommitMessage()` using `git commit --amend`

5. **git-log** (`internal/git/log.go`)  
   - Implement `ShowLog(count)` and `ShowReflog(count)`  
   - Render output as Lipgloss-styled tables

6. **api-ai** (`internal/api/ai.go`)  
   - Implement `GenerateCommitMessage` HTTP POST  
   - Mirror the JSON payload of the existing Python `AIService`

7. **api-github** (`internal/api/github.go`)  
   - Implement `GetUser(username)` (GitHub user profile)  
   - Implement `GetUserRepos(username)` (recent repos)

8. **ui-help** (`internal/ui/help.go`)  
   - Render ASCII logo and tagline in Lipgloss  
   - Render commands table similar to current Rich-based `HelpView`

9. **ui-commit** (`internal/ui/commit.go`)  
   - Bubble Tea model and view for AI commit workflow  
   - States: `Idle → Staging → LoadingAI → ShowingProposal → Editing → Committing → Done`

10. **ui-history** (`internal/ui/history.go`)  
    - Bubble Tea model for scope selection and multi-step inputs  
    - States: `ScopeSelect → HashInput (if specific) → FieldInput → Confirm → Done`

11. **ui-profile** (`internal/ui/profile.go`)  
    - Lipgloss-based GitHub profile card view  
    - Show name, login, bio, stats, meta info (location, blog, etc.)

12. **controllers** (`internal/controllers/commit.go`, `internal/controllers/profile.go`)  
    - Wire Git layer + API layer + UI components  
    - Mirror orchestration in current Python controllers

13. **cmd-main** (`cmd/main.go`)  
    - Implement Cobra root command  
    - Add flags/subcommands mirroring existing Python CLI

14. **ci-cd** (`.github/workflows/build.yml`)  
    - Replace PyInstaller matrix with Go cross-compilation  
    - Build Linux, macOS, and Windows binaries in a single job

