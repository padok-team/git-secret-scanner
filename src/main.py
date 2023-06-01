from __future__ import annotations

import argparse, threading
import tempfile, csv
from concurrent import futures

from git import GithubResource, GitResource
from scanners import TrufflehogScanner, GitleaksScanner
from secret import SecretReport


def repository_scan(
    url: str,
    destination_folder: str,
    scan_results: list[SecretReport],
    pool: futures.ThreadPoolExecutor,
    lock: threading.Lock,
):
    repository = url.split('/')[-1].removesuffix('.git')
    destination = f'{destination_folder}/{repository}'

    try:
        GithubResource.clone(url, destination)

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
        scan_results += results
        lock.release()
    except Exception as error:
        # shutdown the whole pool if we catch an error
        pool.shutdown(wait=False, cancel_futures=True)
        # print the error we got
        print(error)


def run(args: argparse.Namespace) -> None:
    organization: str = args.org
    report_file: str = args.file
    visibility: str = args.visibility
    no_archived: bool = args.no_archived

    print(f'Retrieving {organization} repositories...')

    github = GithubResource(organization)
    repo_urls = github.get_repository_urls(visibility, no_archived)

    scan_results: list[SecretReport] = []

    # create tmp directory for cloned repositories
    tmp_directory = tempfile.TemporaryDirectory()

    # setup multithreading
    lock = threading.Lock()
    pool = futures.ThreadPoolExecutor(max_workers=5)

    print('Launching secret scan (this operation can take a while to complete)...')

    # submit tasks to the thread pool
    for url in repo_urls:
        pool.submit(repository_scan, url, tmp_directory.name, scan_results, pool, lock)

    # wait for all tasks to complete
    pool.shutdown(wait=True)

    # delete tmp dir
    tmp_directory.cleanup()

    print('Generating report...')

    # generate CSV report
    with open(report_file, 'w') as file:
        csv_writer = csv.writer(file)
        csv_writer.writerow(['repository', 'path', 'kind', 'line', 'valid', 'cleartext', 'hash'])
        for secret in scan_results:
            csv_writer.writerow([secret.repository, secret.path, secret.kind, secret.line, secret.valid, secret.cleartext, secret.hash])

    # graceful shutdown
    exit(0)


def main() -> None:
    description = 'Scan secrets in organization repositories'
    parser = argparse.ArgumentParser(description=description)

    parser.add_argument('--org',
        type=str,
        help='Organization to scan',
        required=True
    )
    parser.add_argument('-f', '--file',
        type=str,
        default='secret_scanner_report.csv',
        help='Path to the CSV report file to generate',
    )
    parser.add_argument('-v', '--visibility',
        type=str,
        default='all',
        choices=['all', 'private', 'public'],
        help='Repositories visibility',
    )
    parser.add_argument('--no-archived',
        default=False,
        action='store_true',
        help='Do not scan archived repositories',
    )

    # parse arguments
    args = parser.parse_args()

    # run script
    run(args)


if __name__ == '__main__':
    main()
