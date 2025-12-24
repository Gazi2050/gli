from rich.console import Console
from rich.panel import Panel
from rich.table import Table
from rich import box

class ProfileView:
    """
    View component for rendering the simplified GitHub profile dashboard.
    """
    
    def __init__(self):
        """Initialize the view with a Rich console."""
        self.console = Console()

    def render(self, user_data):
        """
        Render the user data into a high-density terminal dashboard.

        Args:
            user_data (dict): The public profile data fetched from the GitHub API.
        """
        header = Panel(
            f"[bold cyan]{user_data.get('name', 'N/A')}[/] (@{user_data.get('login', 'N/A')})\n"
            f"[italic white]{user_data.get('bio', 'No bio available.')}[/]",
            title="GitHub Profile", border_style="cyan", box=box.ROUNDED
        )

        stats_table = Table(show_header=False, box=box.SIMPLE, padding=(0, 2))
        stats_table.add_row("[bold green]Followers[/]", str(user_data.get('followers', 0)))
        stats_table.add_row("[bold green]Following[/]", str(user_data.get('following', 0)))
        stats_table.add_row("[bold green]Repositories[/]", str(user_data.get('public_repos', 0)))
        stats_table.add_row("[bold green]Gists[/]", str(user_data.get('public_gists', 0)))

        account_table = Table(show_header=False, box=box.SIMPLE, padding=(0, 2))
        account_table.add_row("[bold yellow]Location[/]", user_data.get('location', 'Unknown'))
        account_table.add_row("[bold yellow]Company[/]", user_data.get('company', 'None'))
        account_table.add_row("[bold yellow]Blog[/]", user_data.get('blog', 'None'))
        account_table.add_row("[bold yellow]Created[/]", user_data.get('created_at', 'N/A')[:10])

        main_table = Table.grid(expand=True)
        main_table.add_column(ratio=1)
        main_table.add_column(ratio=1)
        main_table.add_row(
            Panel(stats_table, title="Statistics", border_style="green", box=box.ROUNDED),
            Panel(account_table, title="Account Details", border_style="yellow", box=box.ROUNDED)
        )

        self.console.print(header)
        self.console.print(main_table)
