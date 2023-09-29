from typing import Self

from concurrent import futures
import csv
from pathlib import Path
import shutil
import sys

from git_secret_scanner import console
from git_secret_scanner.report import read_report, ReportSecret, ReportWriter
from git_secret_scanner.scanners import GitleaksScanner, TrufflehogScanner, is_ignored
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
        max_concurrency: int,
        git_scm: GitScm,
    ) -> None:
        self.report_path = report_path
        self.clone_path = clone_path
        self.no_clean_up = no_clean_up
        self.baseline_path = baseline_path
        self.fingerprints_ignore_path = fingerprints_ignore_path
        self.max_concurrency = max_concurrency
        self.git_scm = git_scm

    def __repository_scan(self: Self, repo: str) -> set[ReportSecret]:
        destination = f'{self.clone_path}/{repo}'

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

        # now that we have all our secrets, remove secrets ignored by scanners
        ignored_secrets: set[ReportSecret] = set()
        for secret in results:
            path, line = secret.path, secret.line
            # if no line were reported, skip (this should never happen in theory...)
            if line:
                with Path(destination, path).open('r') as file:
                    for idx, content in enumerate(file):
                        if idx == line - 1 and is_ignored(content):
                            ignored_secrets.add(secret)

        return results - ignored_secrets

    def __load_ignored_fingerprints(self: Self) -> set[str]:
        if self.fingerprints_ignore_path is not None:
            try:
                with Path(self.fingerprints_ignore_path).open('r') as fingerprints_ignore_file:
                    return {fingerprint.rstrip() for fingerprint in fingerprints_ignore_file}
            except FileNotFoundError as error:
                msg = f"fingerprints ignore file not found: '{self.fingerprints_ignore_path}'"
                raise FileNotFoundError(msg) from error
        return set()

    def __load_baseline(self: Self) -> set[ReportSecret]:
        return set() if self.baseline_path is None else read_report(self.baseline_path)

    def run(self: Self) -> None:
        # retrieve fingerprints to ignore from the scan
        ignored_fingerprints = self.__load_ignored_fingerprints()
        # retrieve the baseline
        baseline = self.__load_baseline()

        # retrieve the list of all repositories in the organization
        repos: set[str] = set()
        with console.ProgressSpinner(f'Listing {self.git_scm.organization} repositories...') as progress:
            repos = self.git_scm.list_repos()

        # create clone path if missing
        if not Path(self.clone_path).exists():
            Path(self.clone_path).mkdir(parents=True)

        scanned_repos: set[str] = set()
        if Path(self.report_path).exists() and Path(self.report_path).stat().st_size > 0:
            scanned_repos = {secret.repository for secret in read_report(self.report_path)}

        # if the report exist and is not empty, ask the user what to do with the current report
        force_recreate = False
        if len(scanned_repos) > 0:
            msg = (f'\nA report "{self.report_path}" already exists. Do you want to override the current report?\n'
                '   [red]yes[/red]: the current report will be overriden\n'
                '   [blue]no[/blue]: only repositories missing from the report will be scanned\n')
            console.print(msg)
            force_recreate = console.confirm('Choice')
            console.print('')

        # remove already scanned repos from the scan when force recreate is False
        if not force_recreate:
            repos = repos - scanned_repos

        with console.ProgressBar('Scanning repositories...', len(repos)) as progress: # noqa: SIM117
            # submit tasks to the thread pool
            with ReportWriter(self.report_path, force_recreate=force_recreate) as report_writer:
                with futures.ThreadPoolExecutor(max_workers=self.max_concurrency) as executor:
                    scan_futures = {
                        executor.submit(self.__repository_scan, repo): repo for repo in repos
                    }

                    # iterate over completed futures that are yielded
                    for future in futures.as_completed(scan_futures):
                        try:
                            # check that the future did not raise an exception
                            # and retrieve the results from the scan
                            results = future.result()

                            # only add secrets that are not already in the baseline
                            for result in results - baseline:
                                # only add secret in report if its fingerprint is not ignored
                                if result.fingerprint not in ignored_fingerprints:
                                    # add the secret to the report
                                    report_writer.add_secret(result)

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
