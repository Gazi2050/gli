import subprocess
import requests
import os
from rich.console import Console
from rich.table import Table
from rich.panel import Panel
from rich import box

class GitManager:
    """
    Controller for Git operations and repository metadata retrieval.
    
    This class encapsulates subprocess calls to the Git CLI, providing a
    visually rich terminal interface and automated environment overrides.
    """
    
    def __init__(self):
        """Initialize the GitManager with a Rich console instance."""
        self.console = Console()

    def get_config(self, key):
        """
        Retrieve a value from the local or global Git configuration.

        Args:
            key (str): The configuration key to lookup (e.g., 'user.name').

        Returns:
            str: The configuration value if found, otherwise None.
        """
        try:
            return subprocess.check_output(["git", "config", "--get", key]).decode("utf-8").strip()
        except subprocess.CalledProcessError:
            return None

    def get_github_username(self):
        """
        Attempt to detect the GitHub username from Git configuration.

        Returns:
            str: The detected username, defaulting to 'user.name' if 'github.user' is missing.
        """
        return self.get_config("github.user") or self.get_config("user.name")

    def run_command(self, args, env=None):
        """
        Execute a Git command with optional environment variable overrides.

        Args:
            args (list): List of command arguments (excluding the 'git' prefix).
            env (dict, optional): Environment variables to inject into the process.

        Returns:
            bool: True if the command executed successfully, False otherwise.
        """
        current_env = os.environ.copy()
        current_env["FILTER_BRANCH_SQUELCH_WARNING"] = "1"
        if env:
            current_env.update(env)
        
        try:
            subprocess.run(["git"] + args, capture_output=True, text=True, check=True, env=current_env)
            return True
        except subprocess.CalledProcessError as e:
            self.console.print(f"[bold red]✗ Error:[/] {e.stderr.strip()}")
            return False

    def commit_and_push(self, message, path="."):
        """
        Stage changes, commit, and push to origin using current system time.

        Args:
            message (str): The commit message.
            path (str, optional): The path to stage. Defaults to current directory.

        Returns:
            bool: Success status of the entire workflow.
        """
        with self.console.status("[bold green]Working on your commit...[/]"):
            if not self.run_command(["add", path]): return False
            if not self.run_command(["commit", "-m", message]): return False
            if not self.run_command(["push"]): return False
        
        self.console.print(Panel(
            f"Message: [bold white]{message}[/]\nStatus: [bold green]Pushed to Remote[/]",
            title="Commit & Push", border_style="green", box=box.ROUNDED
        ))
        return True

    def _get_modification_scope(self, title):
        """
        Prompt the user to select the breadth of a history modification.

        Args:
            title (str): The title to display in the selection panel.

        Returns:
            str: One of 'single', 'specific', or 'all'.
        """
        self.console.print(Panel(
            f"[bold cyan]{title}[/]\n\n"
            f"[1] Last Commit Only\n"
            f"[2] Specific Commit (by Hash)\n"
            f"[3] All Commits in Branch", 
            box=box.ROUNDED
        ))
        choice = self.console.input("[bold white]Select scope (1/2/3): [/]")
        if choice == "2": return "specific"
        if choice == "3": return "all"
        return "single"

    def change_commit_time(self, date_str=None):
        """
        Update the timestamp for one or more commits via environment filters.

        Args:
            date_str (str, optional): The target date string. If None, prompts interactively.

        Returns:
            bool: Success status of the modification.
        """
        scope = self._get_modification_scope("Time Warp Controller")
        
        target_hash = ""
        if scope == "specific":
            target_hash = self.console.input("[bold white]Enter commit hash: [/]").strip()

        if not date_str:
            date = self.console.input("[bold white]Select commit date (YYYY-MM-DD): [/]")
            time = self.console.input("[bold white]Select commit time (HH:MM, 24h): [/]")
            date_str = f"{date} {time}:00"

        if scope in ["all", "specific"]:
            filter_script = f"export GIT_AUTHOR_DATE='{date_str}'; export GIT_COMMITTER_DATE='{date_str}'"
            if scope == "specific":
                filter_script = f"case \"$GIT_COMMIT\" in {target_hash}*) {filter_script} ;; esac"
            
            cmd = ["filter-branch", "-f", "--env-filter", filter_script, "--", "HEAD"]
            status_msg = "[bold yellow]Rewriting history...[/]"
        else:
            env = {"GIT_AUTHOR_DATE": date_str, "GIT_COMMITTER_DATE": date_str}
            cmd = ["commit", "--amend", "--no-edit", "--date", date_str]
            status_msg = f"[bold yellow]Updating last commit to {date_str}...[/]"

        with self.console.status(status_msg):
            success = self.run_command(cmd, env=None if scope in ["all", "specific"] else env)
        
        if success:
            detail = f"Target: [bold]{target_hash if scope == 'specific' else scope.capitalize()}[/]"
            self.console.print(Panel(
                f"{detail}\nNew Timestamp: [bold yellow]{date_str}[/]",
                title="Time Warp", border_style="yellow", box=box.ROUNDED
            ))
        return success

    def change_commit_author(self):
        """
        Interactively update the author identity for one or more commits.

        Returns:
            bool: Success status of the modification.
        """
        scope = self._get_modification_scope("Identity Swapper")
        
        target_hash = ""
        if scope == "specific":
            target_hash = self.console.input("[bold white]Enter commit hash: [/]").strip()

        name = self.console.input("[bold white]Enter author name: [/]")
        email = self.console.input("[bold white]Enter author email: [/]")
        author_str = f"{name} <{email}>"

        if scope in ["all", "specific"]:
            filter_script = f"export GIT_AUTHOR_NAME='{name}'; export GIT_AUTHOR_EMAIL='{email}'; export GIT_COMMITTER_NAME='{name}'; export GIT_COMMITTER_EMAIL='{email}'"
            if scope == "specific":
                filter_script = f"case \"$GIT_COMMIT\" in {target_hash}*) {filter_script} ;; esac"
                
            cmd = ["filter-branch", "-f", "--env-filter", filter_script, "--", "HEAD"]
            status_msg = "[bold yellow]Rewriting author in history...[/]"
        else:
            cmd = ["commit", "--amend", "--no-edit", f"--author={author_str}"]
            status_msg = f"[bold yellow]Changing last commit identity...[/]"

        with self.console.status(status_msg):
            success = self.run_command(cmd)
        
        if success:
            detail = f"Target: [bold]{target_hash if scope == 'specific' else scope.capitalize()}[/]"
            self.console.print(Panel(
                f"{detail}\nNew Author: [bold cyan]{author_str}[/]",
                title="Identity Swapped", border_style="cyan", box=box.ROUNDED
            ))
        return success

    def change_commit_message(self):
        """
        Interactively update the commit message of the most recent commit.

        Returns:
            bool: Success status of the modification.
        """
        self.console.print("[bold cyan]Message Rewriter[/]")
        new_msg = self.console.input("[bold white]Enter new commit message: [/]")

        with self.console.status("[bold green]Updating message...[/]"):
            success = self.run_command(["commit", "--amend", "-m", new_msg])
        
        if success:
            self.console.print(Panel(
                f"New Message: [bold white]{new_msg}[/]",
                title="Message Updated", border_style="green", box=box.ROUNDED
            ))
        return success

    def show_log(self, count=10):
        """
        Retrieve and render the commit history in a formatted Rich table.

        Args:
            count (int, optional): Number of commits to display. Defaults to 10.
        """
        try:
            result = subprocess.run(
                ["git", "log", f"-n", str(count), "--pretty=format:%h|%ad|%an|%s", "--date=format:%Y-%m-%d %H:%M"],
                capture_output=True, text=True, check=True
            )
            
            table = Table(title="Git History", box=box.ROUNDED, border_style="cyan")
            table.add_column("Hash", style="dim cyan", no_wrap=True)
            table.add_column("Date & Time", style="green")
            table.add_column("Author", style="yellow")
            table.add_column("Message", style="white")

            for line in result.stdout.strip().split("\n"):
                if "|" in line:
                    table.add_row(*line.split("|"))
            
            self.console.print(table)
        except Exception:
            self.console.print("[bold red]Error:[/] Could not fetch log.")

    def show_reflog(self, count=10):
        """
        Retrieve and render the Git reflog in a formatted Rich table.

        Args:
            count (int, optional): Number of reflog entries to display. Defaults to 10.
        """
        try:
            result = subprocess.run(
                ["git", "reflog", f"-n", str(count), "--pretty=format:%h|%ad|%gs", "--date=format:%Y-%m-%d %H:%M"],
                capture_output=True, text=True, check=True
            )
            
            table = Table(title="Reflog (Recovery)", box=box.ROUNDED, border_style="magenta")
            table.add_column("Index", style="dim white", justify="right")
            table.add_column("Hash", style="dim magenta", no_wrap=True)
            table.add_column("Time", style="green")
            table.add_column("Operation", style="white")

            for i, line in enumerate(result.stdout.strip().split("\n")):
                if "|" in line:
                    table.add_row(f"HEAD@{{{i}}}", *line.split("|"))
            
            self.console.print(table)
        except Exception:
            self.console.print("[bold red]Error:[/] Could not fetch reflog.")

    def reset_commit(self, mode="soft"):
        """
        Reset the current branch head to the previous commit.

        Args:
            mode (str, optional): The reset mode ('soft' or 'hard'). Defaults to 'soft'.

        Returns:
            bool: Success status of the reset operation.
        """
        if mode not in ["soft", "hard"]:
            self.console.print("[bold red]Error:[/] Invalid mode.")
            return False
            
        with self.console.status(f"[bold red]Resetting: {mode}...[/]"):
            success = self.run_command(["reset", f"--{mode}", "HEAD~1"])
        
        if success:
            self.console.print(Panel(
                f"Reset to [bold]HEAD~1[/] using [bold red]{mode}[/] mode.",
                title="Git Reset", border_style="red", box=box.ROUNDED
            ))
        return success

    def switch_branch(self, name):
        """
        Create a new branch and switch to it immediately.

        Args:
            name (str): The name of the new branch.

        Returns:
            bool: Success status of the branch operation.
        """
        with self.console.status(f"[bold blue]Switching to {name}...[/]"):
            if not self.run_command(["checkout", "-b", name]):
                return False
            success = self.run_command(["push", "-u", "origin", name])
            
        if success:
            self.console.print(Panel(
                f"Active Branch: [bold cyan]{name}[/]\nTracking origin.",
                title="Branch Switch", border_style="blue", box=box.ROUNDED
            ))
        return success
