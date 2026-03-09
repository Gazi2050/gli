import argparse
import sys
from utils.git import GitManager
from utils.api import GitHubAPI, AIService
from controllers import CommitController, ProfileController
from components.profile_view import ProfileView
from components.help_view import HelpView

class GLIApp:
    """
    Slim entry point for the gli CLI.
    
    Handles argument parsing and routes commands to specialized controllers.
    """
    
    def __init__(self):
        self.git = GitManager()
        self.github_api = GitHubAPI()
        self.ai_service = AIService()
        self.profile_view = ProfileView()
        self.help_view = HelpView()

        self.commit_ctrl = CommitController(self.git, self.ai_service)
        self.profile_ctrl = ProfileController(self.git, self.github_api, self.profile_view)

    def run(self):
        """
        Parse command-line arguments and dispatch to controllers.
        """
        VERSION = "[[STAMP]]"
        if VERSION.startswith("[[") and VERSION.endswith("]]"):
            VERSION = "dev-local"
            
        parser = argparse.ArgumentParser(description="gli - Modern Git Wrapper")
        
        parser.add_argument("-v", "--version", action="version", version=f"%(prog)s {VERSION}")
        parser.add_argument("-c", "--commit", nargs="?", const="prompt", help="Commit and push")
        parser.add_argument("-ac", "--ai-commit", action="store_true", help="AI-powered commit")
        parser.add_argument("-l", "--log", action="store_true", help="View git log")
        parser.add_argument("-rl", "--reflog", action="store_true", help="View git reflog")
        parser.add_argument("-rs", "--reset", choices=["soft", "hard"], help="Reset last commit")
        parser.add_argument("-s", "--switch", metavar="BRANCH", help="Branch management")
        parser.add_argument("-ct", "--changeTime", nargs="?", const="", metavar="DATE", help="History: Change time")
        parser.add_argument("-ca", "--changeAuthor", action="store_true", help="History: Change author")
        parser.add_argument("-cm", "--changeMessage", action="store_true", help="History: Change message")
        parser.add_argument("-nv", "--no-verify", action="store_true", help="Skip git hooks")
        
        parser.add_argument("-lb", "--local-branch", action="store_true", help="Branch: Create local only")
        parser.add_argument("-rb", "--remote-branch", action="store_true", help="Branch: Push to remote (default)")
        
        parser.add_argument("command", nargs="?", choices=["profile", "me"], help="Profile commands")
        parser.add_argument("username", nargs="?", help="Target username")

        args = parser.parse_args()

        if args.commit is not None:
            if args.commit == "prompt":
                self.commit_ctrl.handle_manual_commit(no_verify=args.no_verify)
            else:
                self.git.commit_and_push(args.commit, no_verify=args.no_verify)
        
        elif args.ai_commit:
            self.commit_ctrl.handle_ai_commit(no_verify=args.no_verify)
            
        elif args.log:
            self.git.show_log()
            
        elif args.reflog:
            self.git.show_reflog()
            
        elif args.reset:
            self.git.reset_commit(args.reset)
            
        elif args.switch:
            push_to_remote = True
            if args.local_branch:
                push_to_remote = False
            elif args.remote_branch:
                push_to_remote = True
                
            self.git.switch_branch(args.switch, push_to_remote=push_to_remote)
            
        elif args.changeTime is not None:
            self.git.change_commit_time(args.changeTime if args.changeTime != "" else None)
            
        elif args.changeAuthor:
            self.git.change_commit_author()
            
        elif args.changeMessage:
            self.git.change_commit_message()
            
        elif args.command == "profile":
            self.profile_ctrl.show_profile(args.username)
            
        elif args.command == "me":
            self.profile_ctrl.show_profile()
            
        else:
            self.help_view.render()

if __name__ == "__main__":
    try:
        GLIApp().run()
    except KeyboardInterrupt:
        print("\n\x1b[31mOperation cancelled.\x1b[0m")
        sys.exit(0)
