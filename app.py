import argparse
from utils.api import GitHubAPI
from utils.git_utils import GitManager
from components.profile_view import ProfileView
from components.help_view import HelpView

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
        parser = argparse.ArgumentParser(description="gli - Modern Git Wrapper")
        
        parser.add_argument("-c", "--commit", help="Commit and push: gli -c 'msg'")
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

        if args.commit:
            self.git.commit_and_push(args.commit)
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

if __name__ == "__main__":
    GLIApp().run()
