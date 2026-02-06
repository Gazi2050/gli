from .core import GitCore
from .actions import GitActions
from .history import GitHistory
from .log import GitLog

class GitManager(GitCore, GitActions, GitHistory, GitLog):
    """
    Unified manager for all Git operations, composed of modular specialized classes.
    """
    pass
