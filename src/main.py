from __future__ import annotations

import argparse, threading
import csv, os, shutil
from concurrent import futures

from git import GithubResource, GitResource
from scanners import TrufflehogScanner, GitleaksScanner
from secret import SecretReport


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
        if os.path.exists(destination):
            print('Repository '+repository+' already scanned !')
        else:
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
             # generate CSV report
            with open(report_file, 'a') as file:
                csv_writer = csv.writer(file)
                for result in results:
                    csv_writer.writerow([result.repository, result.path, result.kind, result.line, result.valid, result.cleartext, result.hash])
            lock.release()
    except Exception as error:
        # shutdown the whole pool if we catch an error
        pool.shutdown(wait=False, cancel_futures=True)
        # print the error we got
        print(error)


def run(args: argparse.Namespace) -> None:
    if args.file == None:
        report_file: str = './reports/'+args.org+'.csv'
    else:
        report_file: str = args.file
    
    if args.repo_path == None:
        repo_path:   str = '/tmp/'+args.org
    else:
        repo_path:   str = args.repo_path
    
    organization:    str = args.org
    visibility:      str = args.visibility
    no_archived:    bool = args.no_archived
    clean_up:       bool = args.clean_up

    print(f'Retrieving {organization} repositories...')

    github = GithubResource(organization)
    repo_urls = github.get_repository_urls(visibility, no_archived)

    # create tmp directory for cloned repositories
    if not os.path.exists(repo_path):
        os.makedirs(repo_path)
    
    # setup multithreading
    lock = threading.Lock()
    pool = futures.ThreadPoolExecutor(max_workers=5)

    print('Launching secret scan (this operation can take a while to complete)...')
    
    if not os.path.exists(report_file):
        with open(report_file, 'w') as file:
            csv_writer = csv.writer(file)
            csv_writer.writerow(['repository', 'path', 'kind', 'line', 'valid', 'cleartext', 'hash'])
    
    # submit tasks to the thread pool
    for url in repo_urls:
        pool.submit(repository_scan, url, repo_path, report_file, pool, lock)

    # wait for all tasks to complete
    pool.shutdown(wait=True)

    # delete dir
    if clean_up == True:
        print('Cleaning folder with cloned repositories...')
        shutil.rmtree(repo_path)

    print('Report available !')
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
    parser.add_argument('--repo-path',
        type=str,
        help='Folder path to store repositories',
    )
    parser.add_argument('--clean-up',
        default=False,
        action='store_true',
        help='Clean up folder created after scan',
    )

    # parse arguments
    args = parser.parse_args()
    # run script
    run(args)


if __name__ == '__main__':
    main()
