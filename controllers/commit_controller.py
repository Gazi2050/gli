import readline
from rich.panel import Panel
from rich import box
from typing import Optional

class CommitController:
    """
    Controller for managing manual and AI-powered commit workflows.
    """

    def __init__(self, git_manager, ai_service):
        self.git = git_manager
        self.ai = ai_service

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
        Orchestrate the AI-powered commit workflow.
        """
        if not self.git.run_command(["add", "."]):
            return

        diff = self.git.get_staged_diff()
        if not diff:
            self.git.console.print("[bold yellow]⚠ Info:[/] No changes detected in the repository.")
            return

        username = self.git.get_github_username() or "unknown-user"
        repo_name = self.git.get_repo_name()

        while True:
            with self.git.console.status("[bold cyan]Analyzing changes with AI...[/]"):
                message = self.ai.generate_ai_commit(diff, username, repo_name)
                
            if not message:
                self.git.console.print("[bold red]✗ Error:[/] Failed to generate message from AI.")
                return

            self.git.console.print(f"\n[bold cyan]AI Proposal:[/] [bold white]{message}[/]")

            self.git.console.print(f"\n[bold cyan][1][/] [white]Commit & Push[/]")
            self.git.console.print(f"[bold yellow][2][/] [white]Regenerate[/]")
            self.git.console.print(f"[bold blue][3][/] [white]Edit message manually[/]")
            self.git.console.print(f"[bold red][4][/] [white]Cancel[/]")
            
            choice = self.git.console.input("\n[bold white]Select action (1/2/3/4): [/]").strip()

            if choice == "1":
                self.git.commit_and_push(message)
                break
            elif choice == "2":
                continue 
            elif choice == "3":
                def hook():
                    readline.insert_text(message)
                    readline.redisplay()
                
                readline.set_pre_input_hook(hook)
                try:
                    # \x01 and \x02 are mandatory for readline to correctly handle ANSI wide chars
                    prompt = "\x01\033[1;34m\x02Edit message:\x01\033[0m\x02 "
                    edited_message = input(prompt).strip()
                except (EOFError, KeyboardInterrupt):
                    edited_message = None
                    print()
                finally:
                    readline.set_pre_input_hook(None)

                if edited_message:
                    self.git.commit_and_push(edited_message)
                    break
                else:
                    self.git.console.print("[bold yellow]⚠ Info:[/] Message was empty or cancelled. Returning to proposal.")
                    continue
            else:
                self.git.console.print("[bold yellow]Aborted.[/]")
                break
