<p align="center">
<pre> 
                                       ██████╗ ██╗     ██╗
                                      ██╔════╝ ██║     ██║
                                      ██║  ███╗██║     ██║
                                      ██║   ██║██║     ██║
                                      ╚██████╔╝███████╗██║
                                       ╚═════╝ ╚══════╝╚═╝
</pre>
</p>

<p align="center"><strong>$ git wrapper for humans</strong></p>

gli is a CLI tool that makes Git simpler. Generate commit messages with AI, manage branches, and look up GitHub profiles — right from your terminal.

## Table of Contents

- [Installation](#installation)
- [Features](#features)

## Installation

Download the appropriate binary from the [latest release](https://github.com/Gazi2050/gli/releases), extract it, then follow the commands below for your platform.

---

### Linux

**Install (amd64)**

```bash
sudo mv gli_v0.0.8_linux_amd64 /usr/local/bin/gli && sudo chmod +x /usr/local/bin/gli
```

**Verify** — `gli version`

**Install (arm64)**

```bash
sudo mv gli_v0.0.8_linux_arm64 /usr/local/bin/gli && sudo chmod +x /usr/local/bin/gli
```

**Verify** — `gli version`

**Uninstall**

```bash
sudo rm /usr/local/bin/gli
```

---

### macOS

**Install (amd64)**

```bash
sudo mv gli_v0.0.8_darwin_amd64 /usr/local/bin/gli && sudo chmod +x /usr/local/bin/gli
```

**Install (arm64)**

```bash
sudo mv gli_v0.0.8_darwin_arm64 /usr/local/bin/gli && sudo chmod +x /usr/local/bin/gli
```

**Verify** — `gli version`

**Uninstall**

```bash
sudo rm /usr/local/bin/gli
```

---

### Windows (CMD as Admin)

**Install (amd64)**

```cmd
move gli_v0.0.8_windows_amd64.exe C:\Windows\System32\gli.exe
```

**Install (arm64)**

```cmd
move gli_v0.0.8_windows_arm64.exe C:\Windows\System32\gli.exe
```

**Verify** — `gli version`

**Uninstall**

```cmd
del C:\Windows\System32\gli.exe
```

## Features

### AI Commit Assistant

Stages all changes, generates an AI commit message, and pushes through a single TUI workflow.

> **Commands:**
>
> - `gli commit` — Stage changes, generate AI commit message, and push
> - `gli commit --no-verify` — Skip git hooks

### Branch Management

Quickly create branches with upstream tracking or switch branches using concise flags.

> **Commands:**
>
> - `gli branch -c <name>` — Create branch and push with upstream
> - `gli branch -s <name>` — Switch to existing branch

### Profile Lookup

Display your GitHub profile or any user's profile directly in the terminal.

> **Commands:**
>
> - `gli me` — Show your GitHub profile
> - `gli profile <username>` — Show any GitHub user profile

## Command Cheatsheet

| Command                  | What it does                                                                                       |
| ------------------------ | -------------------------------------------------------------------------------------------------- |
| `gli commit`             | Stage every change, generate an AI commit message, and push (use `--no-verify` to skip git hooks). |
| `gli branch -c <name>`   | Create a new branch, switch to it, and push with upstream tracking in one step.                    |
| `gli branch -s <name>`   | Switch to an existing branch using a concise flag.                                                 |
| `gli me`                 | Show the GitHub profile detected from your local git config.                                       |
| `gli profile <username>` | Display any GitHub user's profile directly in the terminal.                                        |

<p align="center">Made with ❤️ for Git users</p>
