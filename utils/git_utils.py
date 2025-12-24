import subprocess

def get_git_config(key):
    """Retrieve a git configuration value."""
    try:
        return subprocess.check_output(["git", "config", "--get", key]).decode("utf-8").strip()
    except:
        return None

def get_github_username():
    """Attempt to detect the GitHub username from git config."""
    # Try github.user first, then user.name or user.email as fallbacks
    username = get_git_config("github.user")
    if not username:
        username = get_git_config("user.name")
    return username

def run_git_command(args, title=None):
    """Run a git command and return success status and output."""
    from rich.console import Console
    console = Console()
    
    if title:
        console.print(f"[bold blue]→ {title}...[/]")
    
    try:
        result = subprocess.run(["git"] + args, capture_output=True, text=True, check=True)
        return True, result.stdout.strip()
    except subprocess.CalledProcessError as e:
        console.print(f"[bold red]✗ Error:[/] {e.stderr.strip()}")
        return False, e.stderr.strip()

def commit_and_push(message):
    """Stage all changes, commit, and push."""
    from rich.console import Console
    console = Console()
    
    # 1. git add .
    success, _ = run_git_command(["add", "."], "Staging all changes")
    if not success: return False
    
    # 2. git commit -m "..."
    success, _ = run_git_command(["commit", "-m", message], f"Committing with message: {message}")
    if not success: return False
    
    # 3. git push
    success, _ = run_git_command(["push"], "Pushing to remote")
    if not success: return False
    
    console.print("[bold green]✔ Changes committed and pushed successfully![/]")
    return True
