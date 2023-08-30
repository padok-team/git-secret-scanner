from gitlab import Gitlab as PythonGitlab

from .git import GitScm, GitProtocol, RepositoryVisibility


class Gitlab(GitScm):
    def __init__(self,
        organization: str,
        visibility: RepositoryVisibility,
        include_archived: bool,
        server = 'gitlab.com',
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
        gitlab = PythonGitlab(url=f'https://{self.server}', private_token=self._token)
        group = gitlab.groups.get(self.organization)

        repos: list[str] = []

        projects = group.projects.list(
            visibility=self.visibility.gitlab_conv(),
            archived=(None if self.include_archived else False),
            include_subgroups=True,
            iterator=True,
        )
        for repo in projects:
            repos.append(repo.path_with_namespace)

        return repos
