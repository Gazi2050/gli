from rich.console import Console
from rich.panel import Panel
from rich.table import Table
from rich.text import Text
from rich import box

class HelpView:
    """
    View component for rendering the custom 'gli' help screen.
    """
    
    ASCII_LOGO = r"""
 ██████╗ ██╗     ██╗
██╔════╝ ██║     ██║
██║  ███╗██║     ██║
██║   ██║██║     ██║
╚██████╔╝███████╗██║
 ╚═════╝ ╚══════╝╚═╝  
    """
    
    TAGLINE = "a git wrapper for make developer life easy"

    def __init__(self):
        """Initialize the help view."""
        self.console = Console()

    def render(self):
        """
        Render the logo, tagline, and command overview table to the console.
        """
        logo_text = Text(self.ASCII_LOGO, style="bold green")
        tagline_text = Text(self.TAGLINE, style="italic dim white")
        
        self.console.print(logo_text)
        self.console.print(tagline_text)
        self.console.print()

        table = Table(title="Available Commands", box=box.ROUNDED, border_style="dim", expand=True)
        table.add_column("Command", style="bold green", no_wrap=True)
        table.add_column("Flag", style="bold yellow", justify="center")
        table.add_column("Description", style="white")

        commands = [
            ("Commit & Push", "-c, --commit", "Stage all, commit with msg, and push"),
            ("AI Commit", "-ac, --ai-commit", "Generate AI message and push"),
            ("Log", "-l, --log", "View commit history graph"),
            ("Reflog", "-rl, --reflog", "View git reflog"),
            ("Reset", "-rs, --reset", "Reset last commit (soft/hard)"),
            ("Switch Branch", "-s, --switch", "Create, switch, and push new branch"),
            ("Change Time", "-ct, --changeTime", "Update commit timestamp(s)"),
            ("Change Author", "-ca, --changeAuthor", "Update commit author identity"),
            ("Change Message", "-cm, --changeMessage", "Update last commit message"),
            ("No Verify", "-nv, --no-verify", "Skip git hooks during commit"),
            ("My Profile", "me", "View your GitHub profile"),
            ("User Profile", "profile <user>", "View a specific GitHub profile"),
        ]

        for cmd, flag, desc in commands:
            table.add_row(cmd, flag, desc)

        self.console.print(table)
        self.console.print("\n[dim]Usage example: [bold green]gli -c 'feat: msg'[/] or [bold green]gli -ac[/][/]")
