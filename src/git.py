import os, subprocess
from github import Github
from gitlab import Gitlab


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

class GitlabResource(GitResource):
    def get_repository_urls(self, visibility='all', no_archived=False) -> str:
        if visibility not in ['all', 'private', 'public']:
            raise AttributeError('Wrong visibility argument')

        # define visibility
        if visibility == 'all':
            visu = None
        else:
            visu = visibility

        # check GITLAB_TOKEN variable
        gitlab_token = os.environ.get('GITLAB_TOKEN')
        if not gitlab_token:
            raise AttributeError('Missing GITLAB_TOKEN env variable')
        
        # authenticate user and get the group to analyze
        gitlab = Gitlab('https://gitlab.com/', private_token=gitlab_token)
        group = gitlab.groups.get(self.organization)

        repository_urls: list[str] = []

        # get all projects of group
        for repo in group.projects.list(get_all=True,visibility=visu):
                repository_urls.append(repo.ssh_url_to_repo)

        # remove archived repositories if specified by user
        if no_archived:
            tmp: list[str] = []
            for repo in group.projects.list(get_all=True,archived=True,visibility=visu):
                    tmp.append(repo.ssh_url_to_repo)
            repository_urls = [x for x in repository_urls if x not in [y for y in tmp]]

        return repository_urls

