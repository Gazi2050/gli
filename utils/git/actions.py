import subprocess
from rich.panel import Panel
from rich import box
from typing import Optional

class GitActions:
    """
    Primary Git workflows: commit, push, branch management, and reset.
    """
    # This class expects run_command, get_current_branch, and console to be available 
    # (via GitCore mixin).

    def commit_and_push(self, message: str, path: str = ".", no_verify: bool = False) -> bool:
        """
        Stage changes, commit, and push. Automatically handles upstream tracking.
        """
        with self.console.status("[bold green]Working on your commit...[/]"):
            if not self.run_command(["add", path]): return False
            
            commit_cmd = ["commit", "-m", message]
            if no_verify:
                commit_cmd.append("--no-verify")
                
            if not self.run_command(commit_cmd): return False
            
            branch = self.get_current_branch()
            has_upstream = False
            if branch:
                try:
                    subprocess.check_output(["git", "rev-parse", "--symbolic-full-name", "@{u}"], stderr=subprocess.DEVNULL)
                    has_upstream = True
                except subprocess.CalledProcessError:
                    has_upstream = False

            if has_upstream:
                push_success = self.run_command(["push"])
            else:
                push_success = self.run_command(["push", "--set-upstream", "origin", branch])

            if not push_success: return False
        
        self.console.print(Panel(
            f"Message: [bold white]{message}[/]\nStatus: [bold green]Pushed to Remote[/]",
            title="Commit & Push", border_style="green", box=box.ROUNDED
        ))
        return True

    def reset_commit(self, mode: str = "soft") -> bool:
        """
        Reset the current branch head to the previous commit.
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

    def switch_branch(self, name: str) -> bool:
        """
        Create a new branch and switch to it immediately.
        """
        with self.console.status(f"[bold blue]Switching to {name}...[/]"):
            if not self.run_command(["checkout", "-b", name]):
                return False
            success = self.run_command(["push", "-u", "origin", name])
            
        if success:
            self.console.print(Panel(
                f"Active Branch: [bold green]{name}[/]\nTracking origin.",
                title="Branch Switch", border_style="blue", box=box.ROUNDED
            ))
        return success
