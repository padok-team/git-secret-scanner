from typing import Annotated, Optional, Callable, cast

from functools import update_wrapper
from inspect import signature
import os
import shutil
import typer

from git_secret_scanner import constants
from git_secret_scanner.scan import Scan
from git_secret_scanner.scm import RepositoryVisibility, GitProtocol, GitScm, Github, Gitlab


# required tools that must be found in PATH
REQUIREMENTS=('git', 'trufflehog', 'gitleaks')


pretty_debug = os.environ.get('PRETTY_DEBUG') in ['1', 'True']

cli = typer.Typer(pretty_exceptions_enable=pretty_debug)


org_option = typer.Option('--org', '-o',
    metavar='<organization>',
    show_default=False,
    help='Organization to scan.',
)
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
    show_default=False,
    help='Folder path to store repositories.',
)
no_clean_up_option = typer.Option('--no-clean-up',
    help='Do not clean up repositories downloaded after the scan.',
)
ssh_clone_option = typer.Option('--ssh-clone',
    help='Use SSH to clone repositories instead of HTTPS.',
)
server_option = typer.Option('--server',
    metavar='<hostname>',
    help='The hostname of the private server (eg. github.com)',
)
baseline_path_option = typer.Option('--baseline-path', '-b',
    metavar='<path>',
    show_default=False,
    help='Path to the CSV report to use as baseline for the scan.',
)
fingerprints_ignore_path_option = typer.Option('--fingerprints-ignore-path', '-i',
    metavar='<path>',
    show_default=False,
    help='Path to file with newline separated fingerprints (SHA-256) of secrets to ignore during the scan.',
)
max_concurrency_option = typer.Option('--max-concurrency',
    metavar='<number>',
    help='Maximum number of concurrent workers.',
)


def scm_command(scm_cls: type[GitScm], token_var: str) -> Callable[[Callable], Callable]:
    def inner(fn: Callable) -> Callable:
        def wrapper(
            org: Annotated[str, org_option],
            visibility: Annotated[RepositoryVisibility, visibility_option] = RepositoryVisibility.All,
            no_archived: Annotated[bool, no_archived_option] = False,
            report_path: Annotated[str, report_path_option] = 'report.csv',
            clone_path: Annotated[Optional[str], clone_path_option] = None,
            no_clean_up: Annotated[bool, no_clean_up_option] = False,
            server: Annotated[Optional[str], server_option] = None,
            ssh_clone: Annotated[bool, ssh_clone_option] = False,
            fingerprints_ignore_path: Annotated[Optional[str], fingerprints_ignore_path_option] = None,
            baseline_path: Annotated[Optional[str], baseline_path_option] = None,
            max_concurrency: Annotated[int, max_concurrency_option] = 5,
        ) -> None:
            # look for the required 'token_var' environment variable
            token = os.environ.get(token_var)
            if not token:
                msg = f'Missing environment variable: {token_var} is not defined'
                raise RuntimeError(msg)

            git_scm = scm_cls(
                organization=org,
                visibility=visibility,
                include_archived=(not no_archived),
                protocol=(GitProtocol.Ssh if ssh_clone else GitProtocol.Https),
                server=server,
                token=token,
            )

            scan = Scan(
                report_path=report_path,
                clone_path=(constants.DEFAULT_CLONE_PATH if clone_path is None else clone_path),
                no_clean_up=no_clean_up,
                fingerprints_ignore_path=fingerprints_ignore_path,
                baseline_path=baseline_path,
                max_concurrency=max_concurrency,
                git_scm=git_scm,
            )

            scan.run()

        fn_sign = signature(fn)
        wrapper_sign = signature(wrapper)

        params = dict(wrapper_sign.parameters)

        for param in fn_sign.parameters:
            params[param] = fn_sign.parameters[param]

        updated_wrapper = cast(Callable, update_wrapper(wrapper=wrapper, wrapped=fn))
        updated_wrapper.__signature__ = wrapper_sign.replace(parameters=list(params.values()))

        return updated_wrapper
    return inner


@cli.callback()
def check_requirements(_: typer.Context) -> None:
    for tool in REQUIREMENTS:
        if shutil.which(tool) is None:
            msg = f'required tool missing: {tool}'
            raise FileNotFoundError(msg)


@cli.command(help="Scan secrets in a GitHub organization's repositories")
@scm_command(Github, token_var='GITHUB_TOKEN')  # noqa: S106
def github(
    server: Annotated[str, typer.Option('--server',  # noqa: ARG001
        metavar='<hostname>',
        show_default=False,
        help='Server name of the GitHub Enterprise private server.',
    )] = 'github.com',
) -> None:
    pass


@cli.command(help="Scan secrets in a GitLab group's repositories")
@scm_command(Gitlab, token_var='GITLAB_TOKEN')  # noqa: S106
def gitlab(
    org: Annotated[str, typer.Option('-o', '--group',  # noqa: ARG001
        metavar='<group>',
        show_default=False,
        help='Group to scan.',
    )],
    server: Annotated[str, typer.Option('--server',  # noqa: ARG001
        metavar='<hostname>',
        show_default=False,
        help='Server name of the self-hosted GitLab instance.',
    )] = 'gitlab.com',
) -> None:
    pass


if __name__ == '__main__':
    cli()
