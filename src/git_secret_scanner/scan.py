from concurrent import futures
import csv
import os
import shutil
import sys
import tempfile

from git_secret_scanner.console import exit_with_error, ProgressSpinner, ProgressBar
from git_secret_scanner.git import GitResource
from git_secret_scanner.scanners import TrufflehogScanner, GitleaksScanner
from git_secret_scanner.secret import SecretReport


TEMP_DIR_NAME = 'github.padok.git-secret-scanner'


# reports may contain very large secrets
# we increase the field size limit not to run into issues because of that
csv.field_size_limit(sys.maxsize)


class ScanContext:
    def __init__(
        self,
        report_path: str,
        clone_path: str,
        no_clean_up: bool,
        fingerprints_ignore_path: str,
        baseline_path: str,
        git_resource: GitResource,
    ):
        self.report_path = report_path
        self.clone_path = clone_path
        self.no_clean_up = no_clean_up
        self.baseline_path = baseline_path
        self.fingerprints_ignore_path = fingerprints_ignore_path
        self.git_resource = git_resource


def repository_scan(
    repo: str,
    clone_path: str,
    git_resource: GitResource,
) -> set[SecretReport]:
    destination = f'{clone_path}/{repo}'

    # check if repository has already been scanned
    if os.path.exists(destination):
        print('Repository '+repo+' already scanned !')
        return set()
    else:
        # clone repositories and run the tools
        git_resource.clone_repo(
            repo=repo,
            destination=destination,
            shallow_clone=True,
            no_git=True
        )

        trufflehog = TrufflehogScanner(destination, repo)
        trufflehog.scan()
        trufflehog_results = trufflehog.get_results()

        gitleaks = GitleaksScanner(destination, repo)
        gitleaks.scan()
        gitleaks_results = gitleaks.get_results()

        trufflehog_set = set(trufflehog_results)
        gitleaks_set = set(gitleaks_results)

        results = trufflehog_set ^ gitleaks_set
        intersect = trufflehog_set & gitleaks_set
    
        for secret in intersect:
            t_secret = trufflehog_results[trufflehog_results.index(secret)]
            g_secret = gitleaks_results[gitleaks_results.index(secret)]
            results.add(SecretReport.merge(t_secret, g_secret))

        return results
 

def run_scan(context: ScanContext) -> None:
    # retrieve fingerprints to ignore from the scan
    ignored_fingerprints: set[str] = set()
    if len(context.fingerprints_ignore_path) > 0:
        try:
            with open(context.fingerprints_ignore_path, 'r') as fingerprints_ignore_file:
                ignored_fingerprints = {
                    fingerprint.rstrip() for fingerprint in fingerprints_ignore_file
                }
        except FileNotFoundError as error:
            exit_with_error('Failed to open fingerprints ingore file', error)

    # retrieve the baseline
    baseline: set[SecretReport] = set()
    if len(context.baseline_path) > 0:
        try:
            with open(context.baseline_path, 'r') as baseline_file:
                csv_reader = csv.DictReader(baseline_file)
                for secret in csv_reader:
                    baseline.add(
                        SecretReport(
                            repository=secret['repository'],
                            path=secret['path'],
                            kind=secret['kind'],
                            line=(int(secret['line']) if len(secret['line']) > 0 else None),
                            valid=(bool(secret['valid']) if len(secret['valid']) > 0 else None),
                            cleartext=(secret['cleartext'] if 'cleartext' in secret else None),
                            fingerprint=secret['fingerprint'],
                        )
                    )
        except FileNotFoundError as error:
            exit_with_error('Failed to open baseline file', error)

    git_resource, repos = context.git_resource, []

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

    # setup the report file with columns
    if not os.path.exists(context.report_path):
        with open(context.report_path, 'w', newline='') as report_file:
            csv_writer = csv.writer(report_file)
            csv_writer.writerow([
                'repository',
                'path',
                'kind',
                'line',
                'valid',
                'cleartext',
                'fingerprint',
            ])

    try:
        with ProgressBar('Scanning repositories...', len(repos)) as progress:
            # submit tasks to the thread pool
            with futures.ThreadPoolExecutor(max_workers=5) as executor:
                scan_futures = [
                    executor.submit(
                        repository_scan,
                        repo,
                        clone_path,
                        git_resource,
                    ) for repo in repos
                ]

                # iterate over completed futures that are yielded
                for future in futures.as_completed(scan_futures):
                    try:
                        # check that the future did not raise an exception
                        # and retrieve the results from the scan
                        results = future.result()

                        # append the results to the report
                        with open(context.report_path, 'a', newline='') as report_file:
                            csv_writer = csv.writer(report_file)
                            # only add secrets that are not already in the baseline
                            for result in results - baseline:
                                # only add secret in report if its fingerprint is not ignored
                                if result.fingerprint not in ignored_fingerprints:
                                    if result not in baseline:
                                        csv_writer.writerow([
                                            result.repository,
                                            result.path,
                                            result.kind,
                                            result.line,
                                            result.valid,
                                            result.cleartext,
                                            result.fingerprint,
                                        ])

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

    # delete cloned repositories when cleanup is not disabled
    if not context.no_clean_up:
        try:
            with ProgressSpinner('Cleaning up cloned repositories...') as progress:
                    shutil.rmtree(clone_path)
        except Exception as error:
            exit_with_error('Failed to perform cleanup', error)
