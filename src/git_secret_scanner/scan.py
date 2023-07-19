from typing import Callable, Any
import enum

import threading
from concurrent import futures
import csv
import os
import shutil
import tempfile
from rich.progress import Progress, SpinnerColumn, TextColumn

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


def repository_scan(
    url: str,
    folder: str,
    report_file: str,
    pool: futures.ThreadPoolExecutor,
    lock: threading.Lock,
):
    repository = url.split('/')[-1].removesuffix('.git')
    destination = f'{folder}/{repository}'

    try:
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

            lock.acquire()
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
            lock.release()
    except Exception as error:
        # shutdown the whole pool if we catch an error
        pool.shutdown(wait=False, cancel_futures=True)
        # print the error we got
        print(error)


def task_with_progress_spiner(description: str, task: Callable) -> Any:
    with Progress(SpinnerColumn(), TextColumn('[progress.description]{task.description}')) as progress:  # noqa: E501
        progress.add_task(description=description, total=None)
        result = task()
    return result


def run_scan(context: ScanContext, git_resource: GitResource) -> None:
    # retrieving org repositories
    repo_urls: list[str] = task_with_progress_spiner(
        f'Listing {git_resource.organization} repositories...',
        git_resource.get_repository_urls,
    )

    # create tmp directory for cloned repositories
    repo_path = context.repo_path
    if repo_path == '':
        repo_path = f'{tempfile.gettempdir()}/{TEMP_DIR_NAME}/{git_resource.organization}'

    if not os.path.exists(repo_path):
        os.makedirs(repo_path)

    # setup multithreading
    lock = threading.Lock()
    pool = futures.ThreadPoolExecutor(max_workers=5)

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

    with Progress() as progress:
        task = progress.add_task("Scanning repositories...", total=len(repo_urls))

        # submit tasks to the thread pool
        for url in repo_urls:
            future = pool.submit(repository_scan, url, repo_path, context.file, pool, lock)
            future.add_done_callback(lambda _: progress.update(task, advance=1))

        # wait for all tasks to complete
        pool.shutdown(wait=True)

    # delete cloned repositories when cleanup is enabled
    if not context.no_clean_up:
        task_with_progress_spiner(
            'Cleaning cloned repositories...',
            lambda: shutil.rmtree(repo_path),
        )
