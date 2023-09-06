from typing import Self

from gitlab import Gitlab as PythonGitlab

from .git import GitScm, GitProtocol, RepositoryVisibility


class Gitlab(GitScm):
    def __init__(self: Self,
        organization: str,
        visibility: RepositoryVisibility,
        include_archived: bool,
        server: str = 'gitlab.com',
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

    def list_repos(self: Self) -> list[str]:
        gitlab = PythonGitlab(url=f'https://{self.server}', private_token=self._token)
        group = gitlab.groups.get(self.organization)

        projects = group.projects.list(
            visibility=self.visibility.gitlab_conv(),
            archived=(None if self.include_archived else False),
            include_subgroups=True,
            iterator=True,
        )

        return [repo.path_with_namespace for repo in projects]
