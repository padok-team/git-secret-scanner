import enum

import threading
from concurrent import futures
import csv
import os
import shutil
import tempfile

from git_secret_scanner.console import print, exit_with_error, ProgressSpinner, ProgressBar
from git_secret_scanner.git import GitResource
from git_secret_scanner.scanners import TrufflehogScanner, GitleaksScanner
from git_secret_scanner.secret import SecretReport


TEMP_DIR_NAME = 'github.padok.git-secret-scanner'


class ScanType(enum.StrEnum):
    Github = 'github'
    Gitlab = 'gitlab'


class ScanContext:
    scan_type: ScanType
    file: str
    repo_path: str
    no_clean_up: bool


# lock for multithreading
lock = threading.Lock()


def repository_scan(url: str, folder: str, report_file: str):
    repository = url.split('/')[-1].removesuffix('.git')
    destination = f'{folder}/{repository}'

    # check if repository has already been scanned
    if os.path.exists(destination):
        print('Repository '+repository+' already scanned !')
    else:
        # clone repositories and run the tools
        GitResource.clone(url, destination)
        trufflehog = TrufflehogScanner(destination, repository)
        trufflehog.scan()
        trufflehog_results = trufflehog.get_results()
        
        gitleaks = GitleaksScanner(destination, repository)
        gitleaks.scan()
        gitleaks_results = gitleaks.get_results()
        results = list(set(trufflehog_results))
        for gr in gitleaks_results:
            if gr in results:
                tr = results.pop(results.index(gr))
                results.append(SecretReport.merge(gr, tr))
            else:
                results.append(gr)

        # take the lock while manipulating the CSV file
        # TODO: this should normally be performed by the underlying OS, this should be removed
        with lock:
            # generate CSV report
            with open(report_file, 'a') as file:
                csv_writer = csv.writer(file)
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


def run_scan(context: ScanContext, git_resource: GitResource) -> None:
    with ProgressSpinner(f'Listing {git_resource.organization} repositories...') as progress:
        try:
            # retrieving org repositories
            repo_urls = git_resource.get_repository_urls()
        except Exception as error:
            progress.error()
            exit_with_error(error)

    # create tmp directory for cloned repositories
    repo_path = context.repo_path
    if repo_path == '':
        repo_path = f'{tempfile.gettempdir()}/{TEMP_DIR_NAME}/{git_resource.organization}'

    if not os.path.exists(repo_path):
        os.makedirs(repo_path)

    # setup the report file with column
    if not os.path.exists(context.file):
        with open(context.file, 'w') as report_file:
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

    with ProgressBar('Scanning repositories...', len(repo_urls)) as progress:
        # submit tasks to the thread pool
        with futures.ThreadPoolExecutor(max_workers=5) as executor:
            scan_futures = [
                executor.submit(repository_scan, url, repo_path, context.file) for url in repo_urls
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
                        executor.shutdown(wait=False, cancel_futures=True)
                        progress.error()
                        exit_with_error('Scan failed', error)

    # delete cloned repositories when cleanup is enabled
    if not context.no_clean_up:
        with ProgressSpinner(f'Listing {git_resource.organization} repositories...') as progress:
            try:
                shutil.rmtree(repo_path)
            except Exception as error:
                progress.error()
                exit_with_error(error)
