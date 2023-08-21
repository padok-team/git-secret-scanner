from concurrent import futures
import csv
import os
import shutil
import tempfile

from git_secret_scanner.console import exit_with_error, ProgressSpinner, ProgressBar
from git_secret_scanner.git import GitResource
from git_secret_scanner.scanners import TrufflehogScanner, GitleaksScanner
from git_secret_scanner.secret import SecretReport


TEMP_DIR_NAME = 'github.padok.git-secret-scanner'


class ScanContext:
    def __init__(
        self,
        report_path: str,
        clone_path: str,
        no_clean_up: bool,
        git_resource: GitResource,
    ):
        self.report_path = report_path
        self.clone_path = clone_path
        self.no_clean_up = no_clean_up
        self.git_resource = git_resource


def repository_scan(
    repo: str,
    report_path: str,
    clone_path: str,
    git_resource: GitResource,
):
    destination = f'{clone_path}/{repo}'

    # check if repository has already been scanned
    if os.path.exists(destination):
        print('Repository '+repo+' already scanned !')
    else:
        # clone repositories and run the tools
        git_resource.clone_repo(repo, destination=destination)

        trufflehog = TrufflehogScanner(destination, repo)
        trufflehog.scan()
        trufflehog_results = trufflehog.get_results()
        
        gitleaks = GitleaksScanner(destination, repo)
        gitleaks.scan()
        gitleaks_results = gitleaks.get_results()

        results = list(set(trufflehog_results))
        for gr in gitleaks_results:
            if gr in results:
                tr = results.pop(results.index(gr))
                results.append(SecretReport.merge(gr, tr))
            else:
                results.append(gr)

        # generate CSV report
        with open(report_path, 'a') as report_file:
            csv_writer = csv.writer(report_file)
            for result in results:
                csv_writer.writerow([
                    result.repository,
                    result.path,
                    result.kind,
                    result.line,
                    result.valid,
                    result.cleartext,
                    result.hash,
                ])


def run_scan(context: ScanContext) -> None:
    git_resource = context.git_resource

    repos = []

    try:
        with ProgressSpinner(f'Listing {git_resource.organization} repositories...') as progress:
            repos = git_resource.list_repos()
    except Exception as error:
        exit_with_error('Failed to list repositories', error)

    # create tmp directory for cloned repositories
    clone_path = context.clone_path
    if clone_path == '':
        clone_path = f'{tempfile.gettempdir()}/{TEMP_DIR_NAME}'

    if not os.path.exists(clone_path):
        os.makedirs(clone_path)

    # setup the report file with column
    if not os.path.exists(context.report_path):
        with open(context.report_path, 'w') as report_file:
            csv_writer = csv.writer(report_file)
            csv_writer.writerow([
                'repository',
                'path',
                'kind',
                'line',
                'valid',
                'cleartext',
                'hash',
            ])

    try:
        with ProgressBar('Scanning repositories...', len(repos)) as progress:
            # submit tasks to the thread pool
            with futures.ThreadPoolExecutor(max_workers=5) as executor:
                scan_futures = [
                    executor.submit(
                        repository_scan,
                        repo,
                        context.report_path,
                        clone_path,
                        git_resource,
                    ) for repo in repos
                ]

                # iterate over completed futures that are yielded
                for future in futures.as_completed(scan_futures):
                    try:
                        # check that the future did not raise an exception
                        future.result()
                        # update progress
                        progress.update(1)
                    except Exception as error:
                        # if the future is canceled, it is intended and not an error
                        if not isinstance(error, futures.CancelledError):
                            # cancel remaning futures on error
                            executor.shutdown(wait=True, cancel_futures=True)
                            raise error
    except Exception as error:
        exit_with_error('Scan failed', error)

    # delete cloned repositories when cleanup is enabled
    if not context.no_clean_up:
        try:
            with ProgressSpinner('Cleaning up cloned repositories...') as progress:
                    shutil.rmtree(clone_path)
        except Exception as error:
            exit_with_error('Failed to perform cleanup', error)
