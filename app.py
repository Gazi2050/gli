import argparse
from utils.api import GitHubAPI
from utils.git_utils import GitManager
from components.profile_view import ProfileView
from components.help_view import HelpView
from rich.panel import Panel
from rich import box
from rich.prompt import Prompt

class GLIApp:
    """
    Main application controller for the gli CLI.
    
    Coordinates interaction between Git utilities, GitHub API, and terminal UI components.
    """
    
    def __init__(self):
        """Initialize core application services."""
        self.git = GitManager()
        self.api = GitHubAPI()
        self.view = ProfileView()
        self.help = HelpView()

    def show_profile(self, username=None):
        """
        Retrieve and render a GitHub profile dashboard.

        Args:
            username (str, optional): The GitHub username. Defaults to local config detection.
        """
        username = username or self.git.get_github_username()
        if not username:
            print("Error: Username could not be detected or provided.")
            return

        user_data = self.api.fetch_user_data(username)
        if "error" in user_data:
            print(user_data["error"])
            return

        self.view.render(user_data)

    def run(self):
        """
        Execute the application based on command-line arguments.
        
        Parses flags and subcommands to dispatch tasks to appropriate service methods.
        """
        VERSION = "[[STAMP]]"
        if VERSION.startswith("[[") and VERSION.endswith("]]"):
            VERSION = "dev-local"
            
        parser = argparse.ArgumentParser(description="gli - Modern Git Wrapper")
        
        parser.add_argument("-v", "--version", action="version", version=f"%(prog)s {VERSION}")
        parser.add_argument("-c", "--commit", nargs="?", const="prompt", help="Commit and push: gli -c 'msg'")
        parser.add_argument("-ac", "--ai-commit", action="store_true", help="Generate message using AI and commit")
        parser.add_argument("-l", "--log", action="store_true", help="View git log")
        parser.add_argument("-rl", "--reflog", action="store_true", help="View git reflog")
        parser.add_argument("-rs", "--reset", choices=["soft", "hard"], help="Reset last commit")
        parser.add_argument("-s", "--switch", metavar="BRANCH", help="Create and switch branch")
        parser.add_argument("-ct", "--changeTime", nargs="?", const="", metavar="DATE", help="Change commit time")
        parser.add_argument("-ca", "--changeAuthor", action="store_true", help="Change commit author")
        parser.add_argument("-cm", "--changeMessage", action="store_true", help="Change commit message")
        
        parser.add_argument("command", nargs="?", choices=["profile", "me"], help="Command to run")
        parser.add_argument("username", nargs="?", help="GitHub username for profile")

        args = parser.parse_args()

        if args.commit is not None:
            if args.commit == "prompt":
                self.handle_manual_commit()
            else:
                self.git.commit_and_push(args.commit)
        elif args.ai_commit:
            self.handle_ai_commit()
        elif args.log:
            self.git.show_log()
        elif args.reflog:
            self.git.show_reflog()
        elif args.reset:
            self.git.reset_commit(args.reset)
        elif args.switch:
            self.git.switch_branch(args.switch)
        elif args.changeTime is not None:
            self.git.change_commit_time(args.changeTime if args.changeTime != "" else None)
        elif args.changeAuthor:
            self.git.change_commit_author()
        elif args.changeMessage:
            self.git.change_commit_message()
        elif args.command == "profile":
            self.show_profile(args.username)
        elif args.command == "me":
            self.show_profile()
        else:
            self.help.render()

    def handle_manual_commit(self):
        """
        Prompt the user for a commit message, then stage and push.
        """
        self.git.console.print("[bold cyan]Manual Commit[/]")
        message = self.git.console.input("[bold white]Enter commit message: [/]").strip()
        
        if not message:
            self.git.console.print("[bold red]✗ Error:[/] Commit message cannot be empty.")
            return
            
        self.git.commit_and_push(message)

    def handle_ai_commit(self):
        """
        Stage all changes, generate a commit message using AI, and push.
        """
        # Stage all changes first
        if not self.git.run_command(["add", "."]):
            return

        diff = self.git.get_staged_diff()
        if not diff:
            self.git.console.print("[bold yellow]⚠ Info:[/] No changes detected in the repository.")
            return

        username = self.git.get_github_username() or "unknown-user"
        repo_name = self.git.get_repo_name()
        
        custom_instructions = ""

        while True:
            with self.git.console.status("[bold cyan]Analyzing changes with AI...[/]"):
                message = self.api.generate_ai_commit(diff, username, repo_name, custom_instructions)
                
            if not message:
                self.git.console.print("[bold red]✗ Error:[/] Failed to generate message from AI.")
                return

            # Show the generated message
            self.git.console.print()
            self.git.console.print(Panel(
                f"[bold white]{message}[/]",
                title="AI PROPOSAL", border_style="cyan", box=box.ROUNDED
            ))

            # Prompt for action
            self.git.console.print(f"\n[bold cyan][1][/] [white]Commit & Push[/]")
            self.git.console.print(f"[bold yellow][2][/] [white]Regenerate[/]")
            self.git.console.print(f"[bold magenta][4][/] [white]Refine with prompt[/]")
            self.git.console.print(f"[bold blue][5][/] [white]Edit message manually[/]")
            self.git.console.print(f"[bold red][3][/] [white]Cancel[/]")
            
            choice = self.git.console.input("\n[bold white]Select action (1/2/3/4/5): [/]").strip()

            if choice == "1":
                self.git.commit_and_push(message)
                break
            elif choice == "2":
                custom_instructions = "" # Reset instructions for a clean regeneration
                continue 
            elif choice == "4":
                refinement = self.git.console.input("\n[bold magenta]Enter refinement instructions (e.g. 'be more technical'): [/]").strip()
                if refinement:
                    if custom_instructions:
                        custom_instructions += f", {refinement}"
                    else:
                        custom_instructions = refinement
                continue # Loop back with instructions
            elif choice == "5":
                edited_message = Prompt.ask(f"\n[bold blue]Edit message[/]", default=message).strip()
                if edited_message:
                    self.git.commit_and_push(edited_message)
                    break
            else:
                self.git.console.print("[bold yellow]Aborted.[/]")
                break

if __name__ == "__main__":
    GLIApp().run()
