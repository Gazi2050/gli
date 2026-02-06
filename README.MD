# gli
```text
 ██████╗ ██╗     ██╗
██╔════╝ ██║     ██║
██║  ███╗██║     ██║
██║   ██║██║     ██║
╚██████╔╝███████╗██║
 ╚═════╝ ╚══════╝╚═╝    
```
> **a git wrapper to make developer life easy**

A lightweight Git wrapper simplifying complex workflows into intuitive terminal commands with built-in AI intelligence.

### **📥 Installation**

Download the standalone binary for your OS from the [Latest Release](https://github.com/Gazi2050/gli/releases).

| Platform | One-Line Install Command | 🗑️ Uninstall |
| :--- | :--- | :--- |
| **🐧 Linux** | `sudo mv gli /usr/local/bin/gli && sudo chmod +x /usr/local/bin/gli` | `sudo rm /usr/local/bin/gli` |
| **🍎 macOS** | `sudo mv gli-mac /usr/local/bin/gli && sudo chmod +x /usr/local/bin/gli` | `sudo rm /usr/local/bin/gli` |
| **🪟 Windows** | `move gli.exe C:\Windows\System32\` (CMD Admin) | `del C:\Windows\System32\gli.exe` |

#### **🔄 Updating**
To update, just download the new binary and run the install command again to overwrite the old one.
```bash
gli --version  # Verify your version
```

### **🚀 AI-Powered Workflow**

`gli` now integrates directly with a powerful AI backend to analyze your changes and suggest the perfect commit messages.

- **Smart Analysis**: Analyzes your staged diff to suggest Conventional Commit messages.
- **Interactive Review**: Confirm, regenerate, or edit the AI proposal directly on the same line.
- **Zero Configuration**: Works out of the box with built-in API integration.

### **💎 Features Overview**

| Feature | Commands | Quick Link |
| :--- | :--- | :--- |
| **Manual Commit** | `-c`, `--commit` | [Usage](#-manual-commit--push) |
| **AI Commit** | `-ac`, `--ai-commit` | [Usage](#-ai-commit--push) |
| **History** | `-l`, `--log` | [Usage](#-history-log) |
| **Recovery** | `-rl`, `--reflog` | [Usage](#-recovery-reflog) |
| **Reset** | `-rs`, `--reset` | [Usage](#-reset-commit) |
| **Branch** | `-s`, `--switch` | [Usage](#-branch-management) |
| **Time Warp**| `-ct`, `--changeTime`| [Usage](#-time-warp) |
| **Identity Swap** | `-ca`, `--changeAuthor`| [Usage](#-identity-swap) |
| **Message Rewrite**  | `-cm`, `--changeMessage`| [Usage](#-message-rewrite) |
| **No Verify** | `-nv`, `--no-verify` | [Usage](#-manual-commit--push) |
| **Profile Dashboard** | `me`, `profile` | [Usage](#-github-profile) |

---

### **📖 Detailed Usage**

#### **📦 Manual Commit & Push**
Stages all changes and commits with your message. 
- If you omit the message, `gli` starts an **interactive prompt**.
- **Skipping Hooks**: Use `-nv` to skip tests/linting during commit.
```bash
# With message
gli -c "feat: added new feature"
gli -c "fix: fast fix" -nv

# Interactive prompt
gli -c
gli -c -nv
```

#### **🤖 AI Commit & Push**
Stages all changes and generates a suggested commit message using AI. You can then review the proposal.
- **Option 1**: Commit & Push (Standard flow)
- **Option 2**: Regenerate (Ask AI for a new suggestion)
- **Option 3**: Edit manual (Modify the message directly in-line)
- **Skipping Hooks**: Add `-nv` or `--no-verify` to bypass pre-commit hooks.
```bash
gli -ac
gli -ac -nv
```

#### **📜 History (Log)**
Displays a clean, color-coded table of recent commit history.
```bash
gli -l
```

#### **🔄 Recovery (Reflog)**
Views the Git reference log to find and recover lost commits.
```bash
gli -rl
```

#### **⏪ Reset Commit**
Resets the branch to the previous commit (`soft` or `hard`).
```bash
gli -rs soft
```

#### **🌿 Branch Management**
Creates a new branch, switches to it, and pushes to remote.
```bash
gli -s "feature-branch"
```

#### **⏳ Time Warp**
Interactively update commit timestamps for history correction.
```bash
gli -ct
```

#### **🆔 Identity Swap**
Interactively update author name/email for commits.
```bash
gli -ca
```

#### **✍️ Message Rewrite**
Update the most recent commit message with an interactive prompt.
```bash
gli -cm
```

#### **👤 GitHub Profile**
View a stylish dashboard of public GitHub profile statistics.
```bash
gli me
gli profile gazi2050
```

<p align="center">Made with ❤️ for Git users</p>