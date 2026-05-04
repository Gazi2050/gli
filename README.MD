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

<p align="center"><strong>`$ -git wrapper for humans`</strong></p>

gli is a lightweight, powerful CLI tool designed to make Git operations faster, smarter, and more intuitive. It combines essential Git commands with AI-powered enhancements and useful developer utilities.

## Table of Contents

- [Installation](#installation)
- [Features](#features)

## Installation

Grab the appropriate binary from the [latest release](https://github.com/Gazi2050/gli/releases) and follow the commands below for your platform.

### Linux

**Install**

```bash
sudo mv gli /usr/local/bin/gli && sudo chmod +x /usr/local/bin/gli
```

**Uninstall**

```bash
sudo rm /usr/local/bin/gli
```

### macOS

**Install**

```bash
sudo mv gli-mac /usr/local/bin/gli && sudo chmod +x /usr/local/bin/gli
```

**Uninstall**

```bash
sudo rm /usr/local/bin/gli
```

### Windows (CMD as Admin)

**Install**

```cmd
move gli.exe C:\Windows\System32\
```

**Uninstall**

```cmd
del C:\Windows\System32\gli.exe
```

## Features

### AI Commit Assistant

Stages all changes, generates an AI commit message, and pushes through a single TUI workflow.

> **Command:** `gli commit`  
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

| Command                  | What it does                                                                                               |
| ------------------------ | ------------------------------------------------------------------------------------------------------- |
| `gli commit`             | Stage every change, generate an AI commit message, and push (use `--no-verify` to skip git hooks).   |
| `gli branch -c <name>`   | Create a new branch, switch to it, and push with upstream tracking in one step.       |
| `gli branch -s <name>`   | Switch to an existing branch using a concise flag.                                    |
| `gli me`                 | Show the GitHub profile detected from your local git config.                          |
| `gli profile <username>` | Display any GitHub user's profile directly in the terminal.                           |

<p align="center">Made with ❤️ for Git users</p>
