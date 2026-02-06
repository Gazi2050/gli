import requests
from typing import Dict, List, Optional

class GitHubAPI:
    """
    Client for interacting with the GitHub REST API.
    """
    
    BASE_URL = "https://api.github.com/users/"

    def fetch_user_data(self, username: str) -> Dict[str, str]:
        """
        Retrieve public profile information for a specified GitHub user.
        """
        try:
            response = requests.get(f"{self.BASE_URL}{username}")
            response.raise_for_status()
            return response.json()
        except requests.exceptions.RequestException as e:
            return {"error": f"Failed to fetch data for '{username}': {str(e)}"}

    def fetch_user_repos(self, username: str) -> List[Dict]:
        """
        Retrieve public repositories for a specified GitHub user.
        """
        try:
            params = {"sort": "updated", "per_page": 5}
            response = requests.get(f"{self.BASE_URL}{username}/repos", params=params)
            response.raise_for_status()
            return response.json()
        except requests.exceptions.RequestException:
            return []
