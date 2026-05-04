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

<p align="center"><strong>$git wrapper for humans</strong></p>

gli is a lightweight, powerful CLI tool designed to make Git operations faster, smarter, and more intuitive. It combines essential Git commands with AI-powered enhancements and useful developer utilities.

## Installation

| Platform            | Install                                                                  | Uninstall                         |
| ------------------- | ------------------------------------------------------------------------ | --------------------------------- |
| Linux               | `sudo mv gli /usr/local/bin/gli && sudo chmod +x /usr/local/bin/gli`     | `sudo rm /usr/local/bin/gli`      |
| macOS               | `sudo mv gli-mac /usr/local/bin/gli && sudo chmod +x /usr/local/bin/gli` | `sudo rm /usr/local/bin/gli`      |
| Windows (CMD admin) | `move gli.exe C:\Windows\System32\`                                      | `del C:\Windows\System32\gli.exe` |

Grab the right binary from the [latest release](https://github.com/Gazi2050/gli/releases) and follow the platform-specific steps above.

## Features

| Feature | Description |
| --- | --- |
| AI Commit Assistant | Stages all changes, generates an AI commit message, and pushes through a single TUI flow. |
| Branch Management | Quickly create branches with upstream tracking or switch branches via concise flags. |
| Profile Lookup | Display your GitHub profile (`gli me`) or any user (`gli profile <user>`) without leaving the terminal. |

## Command Cheatsheet

| Command | What it does |
| --- | --- |
| `gli commit` | Stage every change, generate an AI commit message, and push (supports `--no-verify`). |
| `gli branch -c <name>` | Create a new branch, switch to it, and push with upstream tracking in one step. |
| `gli branch -s <name>` | Switch to an existing branch using a concise flag. |
| `gli me` | Show the GitHub profile detected from your local git config. |
| `gli profile <username>` | Display any GitHub user's profile directly in the terminal. |

<p align="center">Made with ❤️ for Git users</p>
