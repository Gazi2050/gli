import requests
import base64
from typing import Optional

class AIService:
    """
    Service for interacting with the AI commit generation backend.
    """
    
    def _get_api_url(self) -> str:
        """
        Retrieve the AI commit API endpoint.
        """
        encoded_url = "aHR0cHM6Ly9kaW55LWNsaS52ZXJjZWwuYXBwL2FwaS92Mi9jb21taXQ="
        return base64.b64decode(encoded_url).decode('utf-8')

    def generate_ai_commit(self, git_diff: str, username: str, repo_name: str, custom_instructions: str = "") -> Optional[str]:
        """
        Request an AI-generated commit message based on staged changes.
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
                        "feat": "🚀", "fix": "🐛", "docs": "📚", "style": "🎨",
                        "refactor": "♻️", "test": "✅", "chore": "🔧", "perf": "⚡"
                    },
                    "Tone": "casual",
                    "Length": "short",
                    "CustomInstructions": custom_instructions,
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
