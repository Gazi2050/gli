import argparse
import sys
from utils.api import fetch_github_user, fetch_avatar
from utils.git_utils import get_github_username, commit_and_push
from components.ascii_converter import image_to_ascii
from components.profile_view import display_profile

def handle_profile(username):
    """Handle the profile command."""
    if not username:
        print("Error: No username provided.")
        return

    user_data = fetch_github_user(username)
    if "error" in user_data:
        print(f"Error: {user_data['error']}")
        return

    avatar_url = user_data.get("avatar_url")
    ascii_art = ""
    if avatar_url:
        avatar_bytes = fetch_avatar(avatar_url)
        if avatar_bytes:
            ascii_art = image_to_ascii(avatar_bytes, width=30)

    display_profile(user_data, ascii_art)

def main():
    parser = argparse.ArgumentParser(description="gli - A lightweight Git wrapper")
    subparsers = parser.add_subparsers(dest="command")

    # Profile command
    profile_parser = subparsers.add_parser("profile", help="Display a GitHub user profile")
    profile_parser.add_argument("username", nargs="?", help="GitHub username")

    # Me command
    me_parser = subparsers.add_parser("me", help="Display your own GitHub profile")

    # Commit command
    commit_parser = subparsers.add_parser("commit", help="Commit and push changes")
    commit_parser.add_argument("-c", "--commit", help="Commit message")
    
    # Aliases and direct flags
    args = sys.argv[1:]
    if not args:
        parser.print_help()
        return

    # Handle direct flags like -c or --commit before subparser
    if args[0] in ["-c", "--commit"]:
        message = None
        if len(args) > 1:
            message = args[1]
        
        if not message:
            print("Error: Commit message required. Usage: gli -c 'message'")
            return
        
        commit_and_push(message)
        return
    
    elif args[0] == "me":
        username = get_github_username()
        if not username:
            print("Error: Could not detect GitHub username from git config.")
            return
        handle_profile(username)
        return
    
    parsed_args = parser.parse_args()

    if parsed_args.command == "profile":
        handle_profile(parsed_args.username)
    elif parsed_args.command == "commit":
        if parsed_args.commit:
            commit_and_push(parsed_args.commit)
        else:
            print("Error: Commit message required.")
    else:
        parser.print_help()

if __name__ == "__main__":
    main()
