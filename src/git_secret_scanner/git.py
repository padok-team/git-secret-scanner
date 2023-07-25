from __future__ import annotations

import enum

import os
import subprocess
from github import Github
from gitlab import Gitlab


class RepositoryVisibility(enum.StrEnum):
    All = 'all'
    Private = 'private'
    Public = 'public'


class GitResource():
    def __init__(self, organization: str, visibility: RepositoryVisibility, include_archived: bool):
        self.organization = organization
        self.visibility = visibility
        self.include_archived = include_archived

    @staticmethod
    def clone(url: str, directory: str) -> None:
        try:
            subprocess.call(['git', 'clone', url, directory], stderr=subprocess.DEVNULL)
        except subprocess.CalledProcessError:
            raise Exception('Failed to clone repository')

    def get_repository_urls(self) -> list[str]:
        raise NotImplementedError('"get_repository_urls" method not implemented')


class GithubResource(GitResource):
    def get_repository_urls(self) -> list[str]:
        # check GITHUB_TOKEN variable
        github_token = os.environ.get('GITHUB_TOKEN')
        if not github_token:
            raise AttributeError('Missing GITHUB_TOKEN env variable')

        github = Github(github_token)

        repository_urls: list[str] = []
        for repo in github.get_organization(self.organization).get_repos(self.visibility):
            if self.include_archived or not repo.archived:
                repository_urls.append(repo.ssh_url)

        return repository_urls


class GitlabResource(GitResource):
    def get_repository_urls(self) -> list[str]:
        visibility = None if self.visibility == RepositoryVisibility.All else self.visibility

        # check GITLAB_TOKEN variable
        gitlab_token = os.environ.get('GITLAB_TOKEN')
        if not gitlab_token:
            raise AttributeError('Missing GITLAB_TOKEN env variable')

        # authenticate user and get the group to analyze
        gitlab = Gitlab('https://gitlab.com/', private_token=gitlab_token)
        group = gitlab.groups.get(self.organization)

        repository_urls: list[str] = []

        # get all projects of group
        for repo in group.projects.list(get_all=True, visibility=visibility):
            repository_urls.append(repo.ssh_url_to_repo)

        # remove archived repositories if specified by user
        if not self.include_archived:
            tmp: list[str] = []
            for repo in group.projects.list(get_all=True, archived=True, visibility=visibility):
                    tmp.append(repo.ssh_url_to_repo)
            repository_urls = [x for x in repository_urls if x not in [y for y in tmp]]

        return repository_urls
