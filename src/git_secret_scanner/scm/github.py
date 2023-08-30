from github import Github as PyGithub

from .git import GitScm, GitProtocol, RepositoryVisibility


class Github(GitScm):
    def __init__(self,
        organization: str,
        visibility: RepositoryVisibility,
        include_archived: bool,
        server = 'github.com',
        protocol = GitProtocol.Https,
        token = '',
    ):
        super().__init__(
            organization,
            visibility,
            include_archived,
            server,
            protocol,
            token,
        )

    def list_repos(self) -> list[str]:
        base_url = f'https://api.{self.server}' if self.server == 'github.com' else f'https://{self.server}/api/v3'
        github = PyGithub(self._token, base_url=base_url)

        repos: list[str] = []

        for repo in github.get_organization(self.organization).get_repos(self.visibility):
            if self.include_archived or not repo.archived:
                repos.append(f'{self.organization}/{repo.name}')

        return repos
