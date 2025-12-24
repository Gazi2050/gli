from rich.console import Console
from rich.layout import Layout
from rich.panel import Panel
from rich.table import Table
from rich import box
from rich.text import Text

def display_profile(user_data, ascii_art):
    """Display GitHub profile in a modern, neofetch-style layout."""
    console = Console()
    
    if "error" in user_data:
        console.print(f"[bold red]Error:[/] {user_data['error']}")
        return

    # Create the info table
    info_table = Table(show_header=False, box=None, padding=(0, 1))
    info_table.add_column("Key", style="bold cyan")
    info_table.add_column("Value", style="white")
    
    info_table.add_row("Name", user_data.get("name") or "N/A")
    info_table.add_row("Username", f"@{user_data.get('login')}")
    info_table.add_row("Bio", user_data.get("bio") or "No bio provided.")
    info_table.add_row("Location", user_data.get("location") or "Unknown")
    info_table.add_row("Company", user_data.get("company") or "N/A")
    info_table.add_row("Blog", user_data.get("blog") or "N/A")
    info_table.add_row("Followers", str(user_data.get("followers")))
    info_table.add_row("Following", str(user_data.get("following")))
    info_table.add_row("Public Repos", str(user_data.get("public_repos")))
    
    # Layout with ASCII and Info
    layout = Layout()
    layout.split_row(
        Layout(Panel(Text(ascii_art, style="green"), box=box.ROUNDED, title="Avatar", border_style="dim"), ratio=1),
        Layout(Panel(info_table, box=box.ROUNDED, title=f"GitHub Profile: {user_data.get('login')}", border_style="cyan"), ratio=2)
    )
    
    console.print(layout)
