from rich.panel import Panel
from rich import box
from typing import Optional

class GitHistory:
    """
    Advanced Git history manipulation: time, author, and message rewriting.
    """
    # This class expects run_command, console, and _get_modification_scope to be available.

    def _get_modification_scope(self, title: str) -> str:
        """
        Prompt the user to select the breadth of a history modification.
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

    def change_commit_time(self, date_str: Optional[str] = None) -> bool:
        """
        Update the timestamp for one or more commits.
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

    def change_commit_author(self) -> bool:
        """
        Interactively update the author identity for one or more commits.
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

    def change_commit_message(self) -> bool:
        """
        Interactively update the commit message of the most recent commit.
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
