import requests

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
