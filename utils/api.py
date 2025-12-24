import requests

def fetch_github_user(username):
    """Fetch public profile data for a GitHub user."""
    url = f"https://api.github.com/users/{username}"
    try:
        response = requests.get(url)
        if response.status_code == 200:
            return response.json()
        elif response.status_code == 404:
            return {"error": f"User '{username}' not found."}
        else:
            return {"error": f"Failed to fetch data: {response.status_code}"}
    except Exception as e:
        return {"error": str(e)}

def fetch_avatar(url):
    """Fetch avatar image bytes."""
    try:
        response = requests.get(url)
        if response.status_code == 200:
            return response.content
        return None
    except:
        return None
