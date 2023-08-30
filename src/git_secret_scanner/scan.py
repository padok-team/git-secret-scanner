from concurrent import futures
import csv
import os
import shutil
import sys
import tempfile

from . import console, scm, scanners, report


# directory name of the temporary directory used by the tool to clone repositories
TEMP_DIR_NAME = 'github.padok.git-secret-scanner'


# reports may contain very large secrets
# we increase the field size limit not to run into issues because of that
csv.field_size_limit(sys.maxsize)


class Context:
    def __init__(
        self,
        report_path: str,
        clone_path: str,
        no_clean_up: bool,
        fingerprints_ignore_path: str,
        baseline_path: str,
        git_scm: scm.GitScm,
    ):
        self.report_path = report_path
        self.clone_path = clone_path
        self.no_clean_up = no_clean_up
        self.baseline_path = baseline_path
        self.fingerprints_ignore_path = fingerprints_ignore_path
        self.git_scm = git_scm


def repository_scan(
    repo: str,
    clone_path: str,
    git_scm: scm.GitScm,
) -> set[report.Secret]:
    destination = f'{clone_path}/{repo}'

    # check if repository has already been scanned
    if os.path.exists(destination):
        print('Repository '+repo+' already scanned !')
        return set()
    else:
        # clone repositories and run the tools
        git_scm.clone_repo(
            repo=repo,
            destination=destination,
            shallow_clone=True,
        )

        trufflehog = scanners.TrufflehogScanner(destination, repo)
        trufflehog.scan()
        trufflehog_results = trufflehog.get_results()

        gitleaks = scanners.GitleaksScanner(destination, repo)
        gitleaks.scan()
        gitleaks_results = gitleaks.get_results()

        trufflehog_set = set(trufflehog_results)
        gitleaks_set = set(gitleaks_results)

        results = trufflehog_set ^ gitleaks_set
        intersect = trufflehog_set & gitleaks_set
    
        for secret in intersect:
            t_secret = (trufflehog_results & {secret}).pop()
            g_secret = (gitleaks_results & {secret}).pop()
            results.add(report.Secret.merge(t_secret, g_secret))

        return results
 

def run(context: Context) -> None:
    # retrieve fingerprints to ignore from the scan
    ignored_fingerprints: set[str] = set()
    if len(context.fingerprints_ignore_path) > 0:
        try:
            with open(context.fingerprints_ignore_path, 'r') as fingerprints_ignore_file:
                ignored_fingerprints = {
                    fingerprint.rstrip() for fingerprint in fingerprints_ignore_file
                }
        except FileNotFoundError as error:
            console.exit_with_error('Failed to open fingerprints ingore file', error)

    # retrieve the baseline
    baseline: set[report.Secret] = set()
    if len(context.baseline_path) > 0:
        try:
            with open(context.baseline_path, 'r') as baseline_file:
                csv_reader = csv.DictReader(baseline_file)
                for secret in csv_reader:
                    baseline.add(
                        report.Secret(
                            repository=secret[report.Column.Repository],
                            path=secret[report.Column.Path],
                            kind=secret[report.Column.Kind],
                            line=(int(secret[report.Column.Line])
                                if len(secret[report.Column.Line]) > 0
                                else None),
                            valid=(bool(secret[report.Column.Valid])
                                if len(secret[report.Column.Valid]) > 0
                                else None),
                            cleartext=(secret[report.Column.Cleartext]
                                if report.Column.Cleartext in secret
                                else None),
                            fingerprint=secret['fingerprint'],
                        )
                    )
        except FileNotFoundError as error:
            console.exit_with_error('Failed to open baseline file', error)

    repos = []

    try:
        with console.ProgressSpinner(
            f'Listing {context.git_scm.organization} repositories...'
        ) as progress:
            repos = context.git_scm.list_repos()
    except Exception as error:
        console.exit_with_error('Failed to list repositories', error)

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
            csv_writer.writerow(list(report.Column))

    try:
        with console.ProgressBar('Scanning repositories...', len(repos)) as progress:
            # submit tasks to the thread pool
            with futures.ThreadPoolExecutor(max_workers=5) as executor:
                scan_futures = [
                    executor.submit(
                        repository_scan,
                        repo,
                        clone_path,
                        context.git_scm,
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
                                        csv_writer.writerow(result.to_row())

                        # update progress
                        progress.update(1)
                    except Exception as error:
                        # if the future is canceled, it is intended and not an error
                        if not isinstance(error, futures.CancelledError):
                            # cancel remaning futures on error
                            executor.shutdown(wait=True, cancel_futures=True)
                            raise error
    except Exception as error:
        console.exit_with_error('Scan failed', error)

    # delete cloned repositories when cleanup is not disabled
    if not context.no_clean_up:
        try:
            with console.ProgressSpinner('Cleaning up cloned repositories...') as progress:
                    shutil.rmtree(clone_path)
        except Exception as error:
            console.exit_with_error('Failed to perform cleanup', error)
