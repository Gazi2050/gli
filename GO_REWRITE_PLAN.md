## gli – Go + Bubble Tea Rewrite Plan

### 1. setup – Go module & dependencies

- **Decide module name**
  - Example: `github.com/Gazi2050/gli` or `github.com/Gazi2050/gli-go`.
- **Initialize Go module**
  - Run: `go mod init <module-path>`.
- **Add dependencies**
  - `github.com/spf13/cobra` for CLI routing.
  - `github.com/charmbracelet/bubbletea` for TUI state/update loop.
  - `github.com/charmbracelet/bubbles` for text input, spinner, table.
  - `github.com/charmbracelet/lipgloss` for styling.
  - Run `go get` for each.
- **Create base folders**
  - `cmd/`, `internal/git/`, `internal/api/`, `internal/ui/`, `internal/controllers/`.

---

### 2. git-core – `internal/git/core.go`

- **Design `GitCore` struct**
  - Fields: `WorkDir string` (optional), maybe a logger if needed.
- **Implement command runner**
  - `RunCommand(args ...string) (stdout string, stderr string, err error)` using `exec.Command("git", args...)`.
  - Set `cmd.Dir` to `WorkDir` if non-empty.
- **Implement helpers**
  - `GetConfig(key string) (string, error)` → `git config --get key`.
  - `GetGithubUsername() (string, error)`:
    - Try `github.user`, fallback `user.name`, fallback `"unknown-user"`.
  - `GetCurrentBranch() (string, error)` → `git rev-parse --abbrev-ref HEAD`.
  - `GetStagedDiff() (string, error)` → `git diff --staged`.
    - Return `""` or special error if no diff.
  - `GetRepoName() (string, error)`:
    - `git config --get remote.origin.url`, split by `/`, strip `.git`.

---

### 3. git-actions – `internal/git/actions.go`

- **Design `GitActions`**
  - Struct embeds or has reference to `GitCore`.
  - Methods should return rich errors (so UI can display).
- **`CommitAndPush(message string, noVerify bool) error`**
  - Call `git add .` (or configurable path).
  - Build `git commit -m <message>` + `--no-verify` if `noVerify`.
  - Detect upstream:
    - Run `git rev-parse --symbolic-full-name @{u}` to check if branch has upstream.
    - If yes: `git push`.
    - If no: `git push --set-upstream origin <branch>`.
- **`ResetCommit(mode string) error`**
  - Validate `mode` ∈ `{soft, hard}`.
  - Run `git reset --<mode> HEAD~1`.
- **`SwitchBranch(name string, pushToRemote bool) error`**
  - Run `git checkout -b <name>`.
  - If `pushToRemote`:
    - `git push -u origin <name>`.
  - Return errors for UI to show.

---

### 4. git-history – `internal/git/history.go`

- **Design history API to match Python semantics**
  - Expose non-interactive methods:
    - `ChangeCommitTime(scope, targetHash, dateTime string) error`.
    - `ChangeCommitAuthor(scope, targetHash, name, email string) error`.
    - `ChangeCommitMessage(newMessage string) error`.
  - Let UI handle user prompts; history layer just runs git commands.
- **Implement scope handling**
  - `scope` string: `"single"`, `"specific"`, `"all"`.
- **Time changes**
  - For `"single"`:
    - Use `git commit --amend --no-edit --date "<dateTime>"` and env `GIT_AUTHOR_DATE`, `GIT_COMMITTER_DATE` if needed.
  - For `"all"` / `"specific"`:
    - Build `git filter-branch -f --env-filter "<script>" -- HEAD`.
    - `script` sets `GIT_AUTHOR_DATE` and `GIT_COMMITTER_DATE`.
    - For `"specific"`, wrap in `case "$GIT_COMMIT" in <hash>* ) ... ;; esac`.
- **Author changes**
  - Similar to time:
    - For `"single"`: `git commit --amend --no-edit --author="<name> <email>"`.
    - For `"all"/"specific"`: `filter-branch` env for name/email.
- **Message changes**
  - Only last commit:
    - `git commit --amend -m "<newMessage>"`.

---

### 5. git-log – `internal/git/log.go`

- **Log commands**
  - `ShowLog(count int) ([]LogEntry, error)`:
    - Run: `git log -n <count> --pretty=format:%h|%ad|%an|%s --date=format:%Y-%m-%d %H:%M`.
    - Parse into struct `{Hash, Time, Author, Message}`.
  - `ShowReflog(count int) ([]ReflogEntry, error)`:
    - Run: `git reflog -n <count> --pretty=format:%h|%ad|%gs --date=format:%Y-%m-%d %H:%M`.
    - Wrap with index to mimic `HEAD@{i}`.
- **UI responsibility**
  - Return slices to UI; UI uses Lipgloss/bubbles table to render.

---

### 6. api-ai – `internal/api/ai.go`

- **Design `AIService` struct**
  - Fields: `Client *http.Client`, `Endpoint string`.
  - Private method to decode base64 URL like Python version.
- **Define request/response structs**
  - `AIRequest` with fields `GitDiff`, `Version`, `Name`, `RepoName`, `System`, `Config`.
  - `AIConfig` with nested `Commit` struct, mirroring Python JSON.
  - `AIResponse` with nested `Data.CommitMessage`.
- **Implement `GenerateCommitMessage`**
  - Signature: `GenerateCommitMessage(diff, username, repoName, customInstructions string) (string, error)`.
  - Populate config: same `Theme`, `Conventional`, `Tone`, `Length`, etc.
  - POST JSON with timeout, handle non-200 error codes.
  - Parse JSON, return commit message or error.

---

### 7. api-github – `internal/api/github.go`

- **Design `GitHubClient`**
  - Fields: `Client *http.Client`, `BaseURL string` (default `https://api.github.com/`).
- **Data models**
  - `User` struct: `Login`, `Name`, `Bio`, `Location`, `PublicRepos`, `Followers`, `Following`, `CreatedAt`, `TwitterUsername`, `Blog`, etc.
  - `Repo` struct: `Name`, `Description`, `StargazersCount`, `Language`, `UpdatedAt`, etc.
- **Methods**
  - `GetUser(username string) (*User, error)`:
    - GET `BaseURL + "users/" + username"`.
  - `GetUserRepos(username string, limit int) ([]Repo, error)`:
    - GET `BaseURL + "users/" + username + "/repos?sort=updated&per_page=<limit>"`.

---

### 8. ui-help – `internal/ui/help.go`

- **ASCII logo + tagline**
  - Port current `HelpView.ASCII_LOGO` and tagline text.
- **Commands table**
  - Recreate Python help table using Lipgloss:
    - Columns: Command, Flag, Description.
    - Same rows as current `HelpView` (commit, AI commit, log, reflog, reset, switch, history, profile).
- **API**
  - Provide a simple `RenderHelp() string` returning styled string.
  - Or `HelpModel` that implements a Bubble Tea model if you want it interactive.

---

### 9. ui-commit – `internal/ui/commit.go`

- **Model design**
  - Struct fields:
    - `state` (enum of Idle/Staging/LoadingAI/ShowingProposal/Editing/Committing/Done).
    - `aiMessage string`.
    - `textInput textinput.Model`.
    - `spinner spinner.Model`.
    - `gitDiff string`.
    - `errorMsg string` (for display).
  - Dependencies passed in:
    - `GitManager` interface.
    - `AIService` interface.
- **State machine**
  - **Idle**: waiting for command to start.
  - **Staging**: call `git add .` via cmd.
  - **LoadingAI**: send Bubble Tea `Cmd` to call `GenerateCommitMessage`.
  - **ShowingProposal**: render AI message and keybindings:
    - `1` → commit and push with AI message.
    - `2` → go back to `LoadingAI` (regenerate).
    - `3` → enter `Editing` with text input pre-filled.
    - `4` / `q` → cancel and exit.
  - **Editing**:
    - Use `textinput.Model` for editing commit message.
    - On `enter`, commit with edited message.
  - **Committing**:
    - Call `CommitAndPush`.
  - **Done**:
    - Show success / return to parent model.
- **Messages / commands**
  - Define custom `tea.Msg` types:
    - `aiGeneratedMsg{Message string, Err error}`.
    - `commitFinishedMsg{Err error}`.
  - Use `tea.Cmd` wrappers that run goroutines to call git/AI methods.

---

### 10. ui-history – `internal/ui/history.go`

- **Model design**
  - Fields:
    - `state` (ScopeSelect, HashInput, TimeInput / AuthorInput / MessageInput, Confirm, Running, Done).
    - `scope` string.
    - `targetHash` string.
    - `inputs` (maybe a series of textinputs).
    - `errorMsg` for validation.
  - Reuse `GitHistory` methods from `internal/git/history.go`.
- **Flow**
  - **ScopeSelect**:
    - Show: `[1] Last commit`, `[2] Specific hash`, `[3] All commits`.
    - Map to `"single" | "specific" | "all"`.
  - **HashInput**:
    - Only if `"specific"`.
  - **FieldInput**:
    - For time: two inputs (date + time) or a single datetime string.
    - For author: name + email.
    - For message: just a message input.
  - **Confirm**:
    - Show summary and `[y] Yes / [n] No`.
  - **Running**:
    - Spinner + calling appropriate `Change*` method via `tea.Cmd`.
  - **Done**:
    - Show result and exit back to main model.

---

### 11. ui-profile – `internal/ui/profile.go`

- **Design**
  - Simple renderer (does not need to be interactive unless you want).
  - Take a `User` struct and maybe a slice of `Repo`s.
- **Layout**
  - Top line: name + `(@login)`.
  - Bio line.
  - Stats line: Public repos, Followers, Following.
  - Meta lines: Location, Twitter, Blog, Joined.
- **Implementation**
  - Use Lipgloss styles:
    - Panel-style border (rounded).
    - Color palette similar to current Rich output.
  - Provide `RenderProfile(user *User, repos []Repo) string`.

---

### 12. controllers – `internal/controllers/commit.go`, `internal/controllers/profile.go`

- **Commit controller**
  - Expose high-level methods:
    - `RunManualCommit(noVerify bool, message string) error`.
    - `RunAICommit(noVerify bool, parent tea.Model) (tea.Model, error)` or similar.
  - For non-TUI (pure CLI path), allow:
    - `CommitAndPush(message string)` without opening Bubble Tea.
- **Profile controller**
  - `ShowProfile(username string) (string, error)`:
    - Resolve username (param or from git config).
    - Call `GetUser` and `GetUserRepos`.
    - Pass to `ui/profile.go` renderer.
    - Return rendered string for printing.

---

### 13. cmd-main – `cmd/main.go`

- **Set up Cobra root command**
  - Name: `gli`.
  - `PersistentFlags`: maybe global options later.
- **Flags mirroring Python CLI**
  - `--commit`, `--ai-commit`, `--log`, `--reflog`, `--reset`, `--switch`, `--local-branch`, `--remote-branch`, `--changeTime`, `--changeAuthor`, `--changeMessage`, `--no-verify`.
  - Positional commands: `me`, `profile <username>`.
- **Version**
  - Global `var version = "dev-local"`.
  - Overridable by `-ldflags "-X main.version=vX.Y.Z"`.
  - Hook up `--version` or `gli version`.
- **Dispatch logic**
  - After parsing flags:
    - If `--commit` set → manual commit path.
    - Else if `--ai-commit` → run AI Bubble Tea model.
    - Else if `--log` / `--reflog` / `--reset` / `--switch` / history flags → call appropriate Git methods and/or open history UI model.
    - Else if `me`/`profile` → call profile controller and print.
    - Else → show help UI (`ui-help`).

---

### 14. ci-cd – `.github/workflows/build.yml`

- **Version job**
  - Keep current versioning logic (compute next tag).
  - Expose `new_tag` output.
- **Build job**
  - Single `ubuntu-latest` job with steps:
    - Checkout.
    - Set up Go (`actions/setup-go`).
    - Build binaries:
      - `GOOS=linux GOARCH=amd64 go build -ldflags="-X main.version=$TAG" -o dist/gli ./cmd`.
      - `GOOS=darwin GOARCH=amd64 go build -ldflags="-X main.version=$TAG" -o dist/gli-mac ./cmd`.
      - `GOOS=windows GOARCH=amd64 go build -ldflags="-X main.version=$TAG" -o dist/gli.exe ./cmd`.
    - Upload all `dist/gli*` to the release with the new tag.


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

