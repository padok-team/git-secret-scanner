from typing_extensions import Annotated

import os
import shutil
import typer

from git_secret_scanner.console import exit_with_error
from git_secret_scanner.git import RepositoryVisibility, GithubResource, GitlabResource
from git_secret_scanner.scan import ScanContext, ScanType, run_scan


REQUIREMENTS=('git', 'trufflehog', 'gitleaks')


pretty_debug = True if os.environ.get('PRETTY_DEBUG') in ['1', 'True'] else False

cli = typer.Typer(pretty_exceptions_enable=pretty_debug)


visibility_option = typer.Option('--visibility', '-v',
    help='Repositories visibility',
)
no_archived_option = typer.Option('--no-archived',
    help='Do not scan archived repositories',
)
file_option = typer.Option('--file', '-f',
    metavar='<file>',
    help='Path to the CSV report file to generate',
)
repo_path_option = typer.Option('--repo-path', '-p',
    metavar='<path>',
    help='Folder path to store repositories',
)
no_clean_up_option = typer.Option('--no-clean-up',
    help='Do not clean up repositories downloaded after the scan',
)


@cli.callback()
def check_requirements(ctx: typer.Context):
    for tool in REQUIREMENTS:
        if shutil.which(tool) is None:
            exit_with_error(f'Required tool missing: {tool} is not installed.')


@cli.command(help='Scan secrets in a GitHub organization\'s repositories')
def github(
    org: Annotated[str, typer.Option('-o', '--org',
        metavar='<organization>',
        help='Organization to scan',
    )],
    visibility: Annotated[RepositoryVisibility, visibility_option] = RepositoryVisibility.All,
    no_archived: Annotated[bool, no_archived_option] = False,
    file: Annotated[str, file_option] = 'report.csv',
    repo_path: Annotated[str, repo_path_option] = '',
    no_clean_up: Annotated[bool, no_clean_up_option] = False,
):
    context = ScanContext()
    context.scan_type = ScanType.Github
    context.file = file
    context.repo_path = repo_path
    context.no_clean_up = no_clean_up

    # look for the requirement GITHUB_TOKEN environment variable
    token = os.environ.get('GITHUB_TOKEN')
    if not token:
        exit_with_error('Missing environment variable: GITHUB_TOKEN is not defined.')
        return

    git_resource = GithubResource(org, visibility, not no_archived, token=token)

    run_scan(context, git_resource)


@cli.command(help='Scan secrets in a GitLab group\'s repositories')
def gitlab(
    group: Annotated[str, typer.Option('-o', '--group',
        metavar='<group>',
        help='Group to scan',
    )],
    visibility: Annotated[RepositoryVisibility, visibility_option] = RepositoryVisibility.All,
    no_archived: Annotated[bool, no_archived_option] = False,
    file: Annotated[str, file_option] = 'report.csv',
    repo_path: Annotated[str, repo_path_option] = '',
    no_clean_up: Annotated[bool, no_clean_up_option] = False,
):
    context = ScanContext()
    context.scan_type = ScanType.Gitlab
    context.file = file
    context.repo_path = repo_path
    context.no_clean_up = no_clean_up

    # look for the requirement GITLAB_TOKEN environment variable
    token = os.environ.get('GITLAB_TOKEN')
    if not token:
        exit_with_error('Missing environment variable: GITLAB_TOKEN is not defined.')
        return

    git_resource = GitlabResource(group, visibility, not no_archived, token=token)

    run_scan(context, git_resource)


if __name__ == '__main__':
    cli()
