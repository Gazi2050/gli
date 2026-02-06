import requests
import base64

class GitHubAPI:
    """
    Client for interacting with the GitHub REST API.
    
    Provides methods to fetch public user metadata and handle rate limits or API errors.
    """
    
    BASE_URL = "https://api.github.com/users/"

    def fetch_user_data(self, username):
        """
        Retrieve public profile information for a specified GitHub user.

        Args:
            username (str): The GitHub username to query.

        Returns:
            dict: Parsed JSON response from GitHub, or a dict containing an 'error' key.
        """
        try:
            response = requests.get(f"{self.BASE_URL}{username}")
            response.raise_for_status()
            return response.json()
        except requests.exceptions.RequestException as e:
            return {"error": f"Failed to fetch data for '{username}': {str(e)}"}

    def fetch_user_repos(self, username):
        """
        Retrieve public repositories for a specified GitHub user, sorted by last update.

        Args:
            username (str): The GitHub username to query.

        Returns:
            list: A list of repository dictionaries, or an empty list if an error occurs.
        """
        try:
            params = {"sort": "updated", "per_page": 5}
            response = requests.get(f"{self.BASE_URL}{username}/repos", params=params)
            response.raise_for_status()
            return response.json()
        except requests.exceptions.RequestException:
            return []

    def _get_api_url(self):
        """
        Decode and return the Base64 encoded AI commit API endpoint.
        
        Returns:
            str: The decoded API URL.
        """
        # Encrypted URL: https://diny-cli.vercel.app/api/v2/commit
        encoded_url = "aHR0cHM6Ly9kaW55LWNsaS52ZXJjZWwuYXBwL2FwaS92Mi9jb21taXQ="
        return base64.b64decode(encoded_url).decode('utf-8')

    def generate_ai_commit(self, git_diff, username, repo_name):
        """
        Send a git diff to the AI backend to generate a meaningful commit message.

        Args:
            git_diff (str): The staged changes (git diff --staged).
            username (str): The local git username for identification.
            repo_name (str): The repository name for identification.

        Returns:
            str: The generated commit message, or None if the request fails.
        """
        url = self._get_api_url()
        payload = {
            "gitDiff": git_diff,
            "version": "v1.0.0",
            "name": username,
            "repoName": repo_name,
            "system": "linux",
            "config": {
                "Theme": "catppuccin",
                "Commit": {
                    "Conventional": True,
                    "ConventionalFormat": ["feat", "fix", "docs", "chore", "style", "refactor", "test", "perf"],
                    "Emoji": False,
                    "EmojiMap": {
                        "feat": "🚀",
                        "fix": "🐛",
                        "docs": "📚",
                        "style": "🎨",
                        "refactor": "♻️",
                        "test": "✅",
                        "chore": "🔧",
                        "perf": "⚡"
                    },
                    "Tone": "casual",
                    "Length": "short",
                    "CustomInstructions": "",
                    "HashAfterCommit": False
                }
            }
        }

        try:
            response = requests.post(url, json=payload, timeout=30)
            response.raise_for_status()
            data = response.json()
            
            if "data" in data and "commitMessage" in data["data"]:
                return data["data"]["commitMessage"]
            return None
        except requests.exceptions.RequestException:
            return None
