class ProfileController:
    """
    Controller for fetching and rendering GitHub profiles.
    """

    def __init__(self, git_manager, github_api, profile_view):
        self.git = git_manager
        self.api = github_api
        self.view = profile_view

    def show_profile(self, username: str = None):
        """
        Retrieve and render a GitHub profile dashboard.
        """
        username = username or self.git.get_github_username()
        if not username:
            self.git.console.print("[bold red]✗ Error:[/] Username could not be detected or provided.")
            return

        user_data = self.api.fetch_user_data(username)
        if "error" in user_data:
            self.git.console.print(f"[bold red]✗ Error:[/] {user_data['error']}")
            return

        self.view.render(user_data)
