from __future__ import annotations

import enum

import subprocess
import shutil
from github import Github
from gitlab import Gitlab


class RepositoryVisibility(enum.StrEnum):
    All = 'all'
    Private = 'private'
    Public = 'public'


class GitProtocol(enum.StrEnum):
    Https = 'https'
    Ssh = 'ssh'


class GitResource():
    def __init__(self,
        organization: str,
        visibility: RepositoryVisibility,
        include_archived: bool,
        server: str,
        protocol = GitProtocol.Https,
        token = '',
    ):
        self.organization = organization
        self.visibility = visibility
        self.include_archived = include_archived
        self.server = server
        self.protocol = protocol
        self._token = token

    def clone_repo(self,
        repo: str,
        destination,
        shallow_clone = False,
        no_git = False,
    ) -> None:
        if self.protocol == GitProtocol.Https:
            clone_url = f'https://x-access-token:{self._token}@{self.server}/{repo}'
        else:
            clone_url = f'git@{self.server}:{repo}'

        shallow_args = []
        if shallow_clone:
            shallow_args = ['--depth', '1']

        proc = subprocess.run([
                'git', 'clone', '--quiet', *shallow_args, clone_url, destination
            ],
            stdout=subprocess.DEVNULL,
            stderr=subprocess.PIPE,
        )

        if no_git:
            shutil.rmtree(f'{destination}/.git')

        if proc.returncode != 0:
            error = RuntimeError(f'failed to clone repository {repo}')
            error.add_note(proc.stderr.decode('utf-8'))
            raise error

    def list_repos(self) -> list[str]:
        raise NotImplementedError('"list_repos()" method not implemented')


class GithubResource(GitResource):
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
        github = Github(self._token, base_url=base_url)

        repos: list[str] = []

        for repo in github.get_organization(self.organization).get_repos(self.visibility):
            if self.include_archived or not repo.archived:
                repos.append(f'{self.organization}/{repo.name}')

        return repos


class GitlabResource(GitResource):
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
        visibility = None if self.visibility == RepositoryVisibility.All else self.visibility
        archived = None if self.include_archived else False

        # authenticate user and get the group to analyze
        gitlab = Gitlab(url=f'https://{self.server}', private_token=self._token)
        group = gitlab.groups.get(self.organization)

        repos: list[str] = []

        projects = group.projects.list(
            visibility=visibility,
            archived=archived,
            include_subgroups=True,
            iterator=True,
        )
        for repo in projects:
            repos.append(repo.path_with_namespace)

        return repos
