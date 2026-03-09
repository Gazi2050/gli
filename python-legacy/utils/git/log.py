import subprocess
from rich.table import Table
from rich import box

class GitLog:
    """
    Visualization tools for Git history: log and reflog.
    """

    def show_log(self, count: int = 10):
        """
        Retrieve and render the commit history in a formatted Rich table.
        """
        try:
            result = subprocess.run(
                ["git", "log", f"-n", str(count), "--pretty=format:%h|%ad|%an|%s", "--date=format:%Y-%m-%d %H:%M"],
                capture_output=True, text=True, check=True
            )
            
            table = Table(title="Git History", box=box.ROUNDED, border_style="green")
            table.add_column("Hash", style="dim green", no_wrap=True)
            table.add_column("Date & Time", style="green")
            table.add_column("Author", style="yellow")
            table.add_column("Message", style="white")

            for line in result.stdout.strip().split("\n"):
                if "|" in line:
                    table.add_row(*line.split("|"))
            
            self.console.print(table)
        except Exception:
            self.console.print("[bold red]Error:[/] Could not fetch log.")

    def show_reflog(self, count: int = 10):
        """
        Retrieve and render the Git reflog in a formatted Rich table.
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
