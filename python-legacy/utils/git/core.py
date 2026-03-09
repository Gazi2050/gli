import subprocess
import os
from rich.console import Console
from typing import Optional, List, Dict

class GitCore:
    """
    Foundation for Git operations and configuration retrieval.
    """
    
    def __init__(self, console: Optional[Console] = None):
        """Initialize with a Rich console instance."""
        self.console = console or Console()

    def get_config(self, key: str) -> Optional[str]:
        """
        Retrieve a value from the local or global Git configuration.
        """
        try:
            return subprocess.check_output(["git", "config", "--get", key]).decode("utf-8").strip()
        except subprocess.CalledProcessError:
            return None

    def get_github_username(self) -> str:
        """
        Attempt to detect the GitHub username from Git configuration.
        """
        return self.get_config("github.user") or self.get_config("user.name") or "unknown-user"

    def get_repo_name(self) -> str:
        """
        Extract the repository name from the remote origin URL.
        """
        url = self.get_config("remote.origin.url")
        if not url:
            return "unknown-repo"
        
        repo_name = url.split("/")[-1]
        if repo_name.endswith(".git"):
            repo_name = repo_name[:-4]
        return repo_name

    def get_current_branch(self) -> Optional[str]:
        """
        Retrieve the name of the currently active branch.
        """
        try:
            return subprocess.check_output(["git", "rev-parse", "--abbrev-ref", "HEAD"]).decode("utf-8").strip()
        except subprocess.CalledProcessError:
            return None

    def get_staged_diff(self) -> Optional[str]:
        """
        Retrieve the diff of staged changes.
        """
        try:
            diff = subprocess.check_output(["git", "diff", "--staged"]).decode("utf-8").strip()
            return diff if diff else None
        except subprocess.CalledProcessError:
            return None

    def run_command(self, args: List[str], env: Optional[Dict[str, str]] = None) -> bool:
        """
        Execute a Git command with optional environment variable overrides.
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
