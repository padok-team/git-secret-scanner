from typing import Self

from github import Github as PyGithub

from .git import GitScm, GitProtocol, RepositoryVisibility


class Github(GitScm):
    def __init__(self: Self,
        organization: str,
        visibility: RepositoryVisibility,
        include_archived: bool,
        server: str = 'github.com',
        protocol: GitProtocol = GitProtocol.Https,
        token: str = '',
    ) -> None:
        super().__init__(
            organization,
            visibility,
            include_archived,
            server,
            protocol,
            token,
        )
        base_url = f'https://api.{server}' if server == 'github.com' else f'https://{server}/api/v3'
        self._github = PyGithub(token, base_url=base_url)

    def list_repos(self: Self) -> set[str]:
        repos: set[str] = set()

        for repo in self._github.get_organization(self.organization).get_repos(self.visibility):
            if self.include_archived or not repo.archived:
                repos.add(f'{self.organization}/{repo.name}')

        return repos
