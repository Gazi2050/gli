# gli
```text
 ██████╗ ██╗     ██╗
██╔════╝ ██║     ██║
██║  ███╗██║     ██║
██║   ██║██║     ██║
╚██████╔╝███████╗██║
 ╚═════╝ ╚══════╝╚═╝    
```
> **A modern Git wrapper that simplifies developer workflows.**

**gli** is a lightweight, powerful CLI tool designed to make Git operations faster, smarter, and more intuitive. It combines essential Git commands with AI-powered enhancements and useful developer utilities.

---

## 📥 Installation

Download the standalone binary for your OS from the [Latest Release](https://github.com/Gazi2050/gli/releases).

| Platform | One-Line Install Command |
| :--- | :--- |
| **🐧 Linux** | `sudo mv gli /usr/local/bin/gli && sudo chmod +x /usr/local/bin/gli` |
| **🍎 macOS** | `sudo mv gli-mac /usr/local/bin/gli && sudo chmod +x /usr/local/bin/gli` |
| **🪟 Windows** | `move gli.exe C:\Windows\System32\` (CMD Admin) |

### 🔄 Updating
Simply download the new binary and overwrite the existing one.
```bash
gli --version
```

---

## 🚀 Features at a Glance

| Command | Flag | Description | Usage |
| :--- | :--- | :--- | :--- |
| **Commit & Push** | `-c`, `--commit` | Stage all, commit with msg, and push | [View](#manual-commit) |
| **AI Commit** | `-ac`, `--ai-commit` | Generate AI message and push | [View](#ai-powered-commit) |
| **Log** | `-l`, `--log` | View commit history graph | [View](#git-log) |
| **Reflog** | `-rl`, `--reflog` | View git reflog | [View](#reflog-recovery) |
| **Reset** | `-rs`, `--reset` | Reset last commit (soft/hard) | [View](#reset) |
| **Switch Branch** | `-s`, `--switch`<br>`-lb`, `--local-branch`<br>`-rb`, `--remote-branch` | Create branch<br>[-lb] Local only<br>[-rb] Push remote | [View](#2-branch-management) |
| **Change Time** | `-ct`, `--changeTime` | Update commit timestamp(s) | [View](#history-modification-the-fix-it-tools) |
| **Change Author** | `-ca`, `--changeAuthor` | Update commit author identity | [View](#history-modification-the-fix-it-tools) |
| **Change Message** | `-cm`, `--changeMessage` | Update last commit message | [View](#history-modification-the-fix-it-tools) |
| **No Verify** | `-nv`, `--no-verify` | Skip git hooks during commit | [View](#manual-commit) |
| **My Profile** | `me` | View your GitHub profile | [View](#5-github-profile) |
| **User Profile** | `profile <user>` | View a specific GitHub profile | [View](#5-github-profile) |

---

## 📖 Command Reference

### 1. Commit & Push
Streamline your daily workflow. Stages all changes (`git add .`), commits, and pushes to the current branch.

#### Manual Commit
Pop a prompt to enter your message, or provide it directly.
```bash
# Interactive prompt
gli -c

# Direct message
gli -c "feat: implementation completed"

# Skip verification (pre-commit hooks)
gli -c "wip: save work" -nv
```

#### AI-Powered Commit
Let AI analyze your `git diff` and suggest a Conventional Commit message.
```bash
gli -ac
```
**Options:**
- `[1] Commit & Push`: Accept the AI suggestion.
- `[2] Regenerate`: Ask for a different suggestion.
- `[3] Edit`: Modify the AI suggestion manually.

---

### 2. Branch Management
Create and switch branches instantly. `gli` handles upstream tracking automatically.

```bash
# Create branch & push to origin (Default)
gli -s "feature/new-ui"

# Create LOCAL branch only (No push)
gli -s "experiment/test" -lb

# Create & push matching REMOTE branch explicitly
gli -s "fix/bug-123" -rb
```

---

### 3. History & Logging

#### Git Log
A cleaner, color-coded alternative to `git log`.
```bash
gli -l
```

#### Reflog (Recovery)
View the reference log to recover lost commits or reset states.
```bash
gli -rl
```

---

### 4. Resets & Corrections
Fix mistakes without memorizing complex Git flags.

#### Reset
Reset your branch to the previous commit.
```bash
# Soft Reset (Keep changes staged)
gli -rs soft

# Hard Reset (Discard all changes)
gli -rs hard
```

#### History Modification (The "Fix It" Tools)
Interactive tools to modify the **most recent commit**.
- **Change Message**: `gli -cm`
- **Change Author**: `gli -ca`
- **Change Date**: `gli -ct`

---

### 5. GitHub Profile
View a stylish dashboard of GitHub statistics for yourself or others.
```bash
# Your profile
gli me

# Another user's profile
gli profile torvalds
```

---

<p align="center">Made with ❤️ for Git users</p>