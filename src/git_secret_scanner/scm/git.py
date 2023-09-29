from __future__ import annotations
from typing import Self

import enum

from pathlib import Path
import subprocess
import shutil

from git_secret_scanner import console, constants


@enum.unique
class RepositoryVisibility(enum.StrEnum):
    All = 'all'
    Private = 'private'
    Public = 'public'

    def gitlab_conv(self: Self) -> RepositoryVisibility | None:
        return None if self == RepositoryVisibility.All else self


@enum.unique
class GitProtocol(enum.StrEnum):
    Https = 'https'
    Ssh = 'ssh'


class GitScm:
    def __init__(self: Self,
        organization: str,
        visibility: RepositoryVisibility,
        include_archived: bool,
        server: str | None,
        protocol: GitProtocol = GitProtocol.Https,
        token: str = '',
    ) -> None:
        self.organization = organization
        self.visibility = visibility
        self.include_archived = include_archived
        self.server = server
        self.protocol = protocol
        self.__token = token

    def clone_repo(self: Self,
        repo: str,
        destination: str,
        shallow_clone: bool = False,
        no_git: bool = False,
    ) -> None:
        if not Path(destination).exists():
            if self.protocol == GitProtocol.Https:
                clone_url = f'https://x-access-token:{self.__token}@{self.server}/{repo}'
            else:
                clone_url = f'git@{self.server}:{repo}'

            shallow_args = []
            if shallow_clone:
                shallow_args = ['--depth', '1']

            proc = subprocess.run([  # noqa: S603, S607
                    'git', 'clone', '--quiet', *shallow_args, clone_url, destination,
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
        elif constants.DEFAULT_CLONE_PATH not in destination:
            console.warn(f'directory {destination} already exists, it will be scanned without cloning.')

    def list_repos(self: Self) -> set[str]:
        msg = '"list_repos()" method not implemented'
        raise NotImplementedError(msg)
