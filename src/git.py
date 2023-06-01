import os, subprocess
from github import Github


class GitResource():
    def __init__(self, organization: str):
        self.organization = organization

    @staticmethod
    def clone(url: str, directory: str) -> None:
        try:
            subprocess.call(['git', 'clone', url, directory], stderr=subprocess.DEVNULL)
        except subprocess.CalledProcessError:
            raise Exception('Failed to clone repository')

    def get_repository_urls(self, visibility='all', no_archived=False) -> str:
        raise NotImplementedError('"get_repository_urls" method not implemented')


class GithubResource(GitResource):
    def get_repository_urls(self, visibility='all', no_archived=False) -> str:
        if visibility not in ['all', 'private', 'public']:
            raise AttributeError('Wrong visibility argument')

        # check GITHUB_TOKEN variable
        github_token = os.environ.get('GITHUB_TOKEN')
        if not github_token:
            raise AttributeError('Missing GITHUB_TOKEN env variable')

        github = Github(github_token)

        repository_urls: list[str] = []
        for repo in github.get_organization(self.organization).get_repos(visibility):
            if not (no_archived and repo.archived):
                repository_urls.append(repo.ssh_url)

        return repository_urls
