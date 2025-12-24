# gli
```text
  ________  ___       ___     
 /  _____/ /  /      /__/     
/  /  ___ /  /       __       
\  \__\  /  /____   /  /      
 \______/ /______/ /__/       
```
> **a git wrapper to make developer life easy**

A lightweight Git wrapper simplifying complex workflows into intuitive terminal commands.

### **📥 Installation**

Download the binary for your OS from the [Latest Release](https://github.com/Gazi2050/gli/releases).

#### **🐧 Linux / 🍎 macOS**
Move it to your system folder to use `gli` from any directory:
#### **🔄 Updating**
To update to the latest version, just download the new binary and run the move command again to overwrite the old one:
```bash
sudo mv gli-Linux /usr/local/bin/gli
```
Check your current version anytime with `gli --version`.

#### **🪟 Windows**
1.  Download `gli-Windows.exe`.
2.  Optional: Rename it to `gli.exe`.
3.  **Pro Tip:** Add the folder to your **System PATH** to use it in any terminal (CMD, PowerShell).

### **🚀 Usage**
```bash
gli
```

### **💎 Features Overview**

| Feature | Commands | Quick Link |
| :--- | :--- | :--- |
| **Commit** | `-c`, `--commit` | [View Usage](#-commit--push) |
| **History** | `-l`, `--log` | [View Usage](#-history-log) |
| **Recovery** | `-rl`, `--reflog` | [View Usage](#-recovery-reflog) |
| **Reset** | `-rs`, `--reset` | [View Usage](#-reset-commit) |
| **Branch** | `-s`, `--switch` | [View Usage](#-branch-management) |
| **Profile** | `me`, `profile` | [View Usage](#-github-profile) |
| **Time Warp**| `-ct`, `--changeTime`| [View Usage](#-time-warp) |
| **Identity** | `-ca`, `--changeAuthor`| [View Usage](#-identity-swap) |
| **Rewrite**  | `-cm`, `--changeMessage`| [View Usage](#-message-rewrite) |

---

### **📖 Detailed Usage**

#### **📦 Commit & Push**
Automatically stages all changes, creates a commit with your message, and pushes to your current branch.
```bash
gli -c "feat: implementation"
gli --commit "feat: implementation"
```

#### **📜 History (Log)**
Displays a high-density, color-coded table of your recent commit history including hashes, timestamps, and authors.
```bash
gli -l
gli --log
```

#### **🔄 Recovery (Reflog)**
Views the Git reference log with numeric indices, making it easy to identify points for history recovery.
```bash
gli -rl
gli --reflog
```

#### **⏪ Reset Commit**
Resets the current branch head to the previous commit. 
- `soft`: Keeps your changes staged in the index.
- `hard`: Completely discards all changes.
```bash
gli -rs soft
gli --reset hard
```

#### **🌿 Branch Management**
Creates a new local branch, switches to it immediately, and sets up remote tracking on origin.
```bash
gli -s "new-feature"
gli --switch "new-feature"
```

#### **👤 GitHub Profile**
Displays a high-density dashboard of a GitHub user's public profile stats and metadata.
```bash
gli me
gli profile gazi2050
```

#### **⏳ Time Warp**
Interactively updates the timestamp of your commits. Supports modifying the **Last**, **Specific**, or **All** commits in history.
```bash
# Shorthand
gli -ct

# Full-hand
gli --changeTime
```

#### **🆔 Identity Swap**
Interactively update the Author name and Email for your commits. Target the **Last**, a **Specific Hash**, or the **Entire History**.
```bash
gli -ca
gli --changeAuthor
```

#### **✍️ Message Rewrite**
Quickly update the commit message of your most recent commit with a beautiful interactive prompt.
```bash
gli -cm
gli --changeMessage
```

<p align="center">Made with ❤️ for Git users</p>