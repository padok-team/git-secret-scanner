from typing import Self

from concurrent import futures
import csv
from pathlib import Path
import shutil
import sys

from git_secret_scanner import console
from git_secret_scanner.report import ReportColumn, ReportSecret, SecretKind
from git_secret_scanner.scanners import GitleaksScanner, TrufflehogScanner
from git_secret_scanner.scm import GitScm


# reports may contain very large secrets
# we increase the field size limit not to run into issues because of that
csv.field_size_limit(sys.maxsize)


class Scan:
    def __init__(self: Self,
        report_path: str,
        clone_path: str,
        no_clean_up: bool,
        fingerprints_ignore_path: str | None,
        baseline_path: str | None,
        git_scm: GitScm,
    ) -> None:
        self.report_path = report_path
        self.clone_path = clone_path
        self.no_clean_up = no_clean_up
        self.baseline_path = baseline_path
        self.fingerprints_ignore_path = fingerprints_ignore_path
        self.git_scm = git_scm

    def __repository_scan(self: Self, repo: str) -> set[ReportSecret]:
        destination = f'{self.clone_path}/{repo}'

        # check if repository has already been scanned
        if Path(destination).exists():
            console.print('Repository '+repo+' already scanned !')
            return set()

        # clone repositories and run the tools
        self.git_scm.clone_repo(
            repo=repo,
            destination=destination,
            shallow_clone=True,
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
            t_secret = (trufflehog_results & {secret}).pop()
            g_secret = (gitleaks_results & {secret}).pop()
            results.add(ReportSecret.merge(t_secret, g_secret))

        return results


    def run(self: Self) -> None:
        # retrieve fingerprints to ignore from the scan
        ignored_fingerprints: set[str] = set()
        if self.fingerprints_ignore_path is not None:
            try:
                with Path(self.fingerprints_ignore_path).open('r') as fingerprints_ignore_file:
                    ignored_fingerprints = {
                        fingerprint.rstrip() for fingerprint in fingerprints_ignore_file
                    }
            except FileNotFoundError as error:
                msg = f"fingerprints ignore file not found: '{self.fingerprints_ignore_path}'"
                raise FileNotFoundError(msg) from error

        # retrieve the baseline
        baseline: set[ReportSecret] = set()
        if self.baseline_path is not None:
            try:
                with Path(self.baseline_path).open('r') as baseline_file:
                    csv_reader = csv.DictReader(baseline_file)
                    for secret in csv_reader:
                        baseline.add(
                            ReportSecret(
                                repository=secret[ReportColumn.Repository],
                                path=secret[ReportColumn.Path],
                                kind=SecretKind[secret[ReportColumn.Kind].lower().capitalize()],
                                line=(int(secret[ReportColumn.Line])
                                    if secret[ReportColumn.Line] != ''
                                    else None),
                                valid=(bool(secret[ReportColumn.Valid])
                                    if secret[ReportColumn.Valid] != ''
                                    else None),
                                cleartext=(secret[ReportColumn.Cleartext]
                                    if ReportColumn.Cleartext in secret
                                    else None),
                                fingerprint=secret['fingerprint'],
                            ),
                        )
            except FileNotFoundError as error:
                msg = f"baseline file not found: '{self.baseline_path}'"
                raise FileNotFoundError(msg) from error

        repos = []

        with console.ProgressSpinner(f'Listing {self.git_scm.organization} repositories...') as progress:
            repos = self.git_scm.list_repos()

        # create clone path if missing
        if not Path(self.clone_path).exists():
            Path(self.clone_path).mkdir(parents=True)

        # setup the report file with columns
        if not Path(self.report_path).exists():
            with Path(self.report_path).open('w', newline='') as report_file:
                csv_writer = csv.writer(report_file)
                csv_writer.writerow(list(ReportColumn))

        with console.ProgressBar('Scanning repositories...', len(repos)) as progress: # noqa: SIM117
            # submit tasks to the thread pool
            with futures.ThreadPoolExecutor(max_workers=5) as executor:
                scan_futures = {
                    executor.submit(self.__repository_scan, repo): repo for repo in repos
                }

                # iterate over completed futures that are yielded
                for future in futures.as_completed(scan_futures):
                    try:
                        # check that the future did not raise an exception
                        # and retrieve the results from the scan
                        results = future.result()

                        # append the results to the report
                        with Path(self.report_path).open('a', newline='') as report_file:
                            csv_writer = csv.writer(report_file)
                            # only add secrets that are not already in the baseline
                            for result in results - baseline:
                                # only add secret in report if its fingerprint is not ignored
                                if result.fingerprint not in ignored_fingerprints:
                                    csv_writer.writerow(result.to_row())

                        # update progress
                        progress.update(1)
                    except Exception as error:  # noqa: PERF203
                        # if the future is canceled, it is intended and not an error
                        if not isinstance(error, futures.CancelledError):
                            # cancel remaning futures on error
                            executor.shutdown(wait=True, cancel_futures=True)
                            msg = f'repository scan failed for {scan_futures[future]}'
                            raise RuntimeError(msg) from error  # noqa: TRY004

        # delete cloned repositories when cleanup is not disabled
        if not self.no_clean_up:
            with console.ProgressSpinner('Cleaning up cloned repositories...') as progress:
                    shutil.rmtree(self.clone_path)
