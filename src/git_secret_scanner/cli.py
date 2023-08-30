from typing import Annotated

import os
import shutil
import typer

from . import console, scm, scan


REQUIREMENTS=('git', 'trufflehog', 'gitleaks')


pretty_debug = os.environ.get('PRETTY_DEBUG') in ['1', 'True']

cli = typer.Typer(pretty_exceptions_enable=pretty_debug)


visibility_option = typer.Option('--visibility', '-v',
    help='Repositories visibility.',
)
no_archived_option = typer.Option('--no-archived',
    help='Do not scan archived repositories.',
)
report_path_option = typer.Option('--report-path', '-r',
    metavar='<path>',
    help='Path to the CSV report file to generate.',
)
clone_path_option = typer.Option('--clone-path', '-c',
    metavar='<path>',
    help='Folder path to store repositories.',
)
no_clean_up_option = typer.Option('--no-clean-up',
    help='Do not clean up repositories downloaded after the scan.',
)
ssh_clone_option = typer.Option('--ssh-clone',
    help='Use SSH to clone repositories instead of HTTPS.',
)
baseline_path_option = typer.Option('--baseline-path', '-b',
    metavar='<path>',
    help='Path to the CSV report to use as baseline for the scan.',
)
fingerprints_ignore_path_option = typer.Option('--fingerprints-ignore-path', '-i',
    metavar='<path>',
    help='Path to file with newline separated fingerprints (SHA-256) of secrets to ignore during the scan.',  # noqa: E501
)


@cli.callback()
def check_requirements(_: typer.Context) -> None:
    for tool in REQUIREMENTS:
        if shutil.which(tool) is None:
            console.exit_with_error(f'Required tool missing: {tool} was not found.')


@cli.command(help="Scan secrets in a GitHub organization's repositories")
def github(
    org: Annotated[str, typer.Option('-o', '--org',
        metavar='<organization>',
        help='Organization to scan.',
    )],
    visibility: Annotated[
        scm.RepositoryVisibility,
        visibility_option,
    ] = scm.RepositoryVisibility.All,
    no_archived: Annotated[bool, no_archived_option] = False,
    report_path: Annotated[str, report_path_option] = 'report.csv',
    clone_path: Annotated[str, clone_path_option] = '',
    no_clean_up: Annotated[bool, no_clean_up_option] = False,
    ssh_clone: Annotated[bool, ssh_clone_option] = False,
    fingerprints_ignore_path: Annotated[str, fingerprints_ignore_path_option] = '',
    baseline_path: Annotated[str, baseline_path_option] = '',
) -> None:
    # look for the requirement GITHUB_TOKEN environment variable
    token = os.environ.get('GITHUB_TOKEN')
    if not token:
        console.exit_with_error('Missing environment variable: GITHUB_TOKEN is not defined.')
        return

    git_scm = scm.Github(
        organization=org,
        visibility=visibility,
        include_archived=(not no_archived),
        protocol=(scm.GitProtocol.Ssh if ssh_clone else scm.GitProtocol.Https),
        token=token,
    )

    context = scan.Context(
        report_path=report_path,
        clone_path=clone_path,
        no_clean_up=no_clean_up,
        fingerprints_ignore_path=fingerprints_ignore_path,
        baseline_path=baseline_path,
        git_scm=git_scm,
    )

    scan.run(context)


@cli.command(help="Scan secrets in a GitLab group's repositories")
def gitlab(
    group: Annotated[str, typer.Option('-o', '--group',
        metavar='<group>',
        help='Group to scan.',
    )],
    visibility: Annotated[
        scm.RepositoryVisibility,
        visibility_option,
    ] = scm.RepositoryVisibility.All,
    no_archived: Annotated[bool, no_archived_option] = False,
    report_path: Annotated[str, report_path_option] = 'report.csv',
    clone_path: Annotated[str, clone_path_option] = '',
    no_clean_up: Annotated[bool, no_clean_up_option] = False,
    ssh_clone: Annotated[bool, ssh_clone_option] = False,
    fingerprints_ignore_path: Annotated[str, fingerprints_ignore_path_option] = '',
    baseline_path: Annotated[str, baseline_path_option] = '',
) -> None:
    # look for the requirement GITLAB_TOKEN environment variable
    token = os.environ.get('GITLAB_TOKEN')
    if not token:
        console.exit_with_error('Missing environment variable: GITLAB_TOKEN is not defined.')
        return

    git_scm = scm.Gitlab(
        organization=group,
        visibility=visibility,
        include_archived=(not no_archived),
        protocol=(scm.GitProtocol.Ssh if ssh_clone else scm.GitProtocol.Https),
        token=token,
    )

    context = scan.Context(
        report_path=report_path,
        clone_path=clone_path,
        no_clean_up=no_clean_up,
        fingerprints_ignore_path=fingerprints_ignore_path,
        baseline_path=baseline_path,
        git_scm=git_scm,
    )

    scan.run(context)


if __name__ == '__main__':
    cli()
