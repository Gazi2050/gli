from rich.console import Console
from rich.panel import Panel
from rich.columns import Columns
from rich import box

class ProfileView:
    """
    View component for rendering a minimalist and stylish GitHub profile card.
    """
    
    def __init__(self):
        """Initialize the view with a Rich console."""
        self.console = Console()

    def render(self, user_data, repos_data=None):
        name = user_data.get('name', 'GitHub User')
        login = user_data.get('login', 'N/A')
        bio = user_data.get('bio') or "[italic white]No bio available.[/]"
        
        stats = [
            f"[bold cyan]Public Repos:[/] {user_data.get('public_repos', 0)}",
            f"[bold magenta]Followers:[/] {user_data.get('followers', 0)}",
            f"[bold yellow]Following:[/] {user_data.get('following', 0)}"
        ]
        stats_row = "  •  ".join(stats)

        location = user_data.get('location')
        twitter = user_data.get('twitter_username')
        blog = user_data.get('blog')
        created = user_data.get('created_at', '')[:10]

        meta_info = []
        if location: meta_info.append(f"[bold cyan]📍  Location:[/] {location}")
        if twitter: meta_info.append(f"[bold blue]🐦  Twitter:[/] @{twitter}")
        if blog: meta_info.append(f"[bold green]🔗  Site:[/] {blog}")
        meta_info.append(f"[bold white]🗓️  Joined:[/] {created}")

        content = [
            f"[bold white font_size=20]{name}[/] [dim](@{login})[/]",
            f"{bio}\n",
            f"{stats_row}\n",
            "\n\n".join(meta_info)
        ]

        self.console.print(
            Panel(
                "\n".join(content),
                border_style="cyan",
                box=box.ROUNDED,
                title="[bold cyan]GitHub Profile[/]",
                padding=(1, 2)
            )
        )
