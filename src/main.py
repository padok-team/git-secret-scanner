from __future__ import annotations

import enum

import argparse
import threading
from concurrent import futures
import csv
import os
import shutil
import tempfile

from git import RepositoryVisibility, GitResource, GithubResource, GitlabResource
from scanners import TrufflehogScanner, GitleaksScanner
from secret import SecretReport


class ScanType(enum.StrEnum):
    Github = 'github'
    Gitlab = 'gitlab'


class GlobalArgs:
    scan_type: ScanType
    file: str
    visibility: RepositoryVisibility
    no_archived: bool
    repo_path: str
    no_clean_up: bool

class GithubArgs(GlobalArgs):
    org: str

class GitlabArgs(GlobalArgs):
    group: str


TEMP_DIR_NAME = 'github.padok.git-secret-scanner'


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
                    csv_writer.writerow([result.repository, result.path, result.kind, result.line, result.valid, result.cleartext, result.hash])
            lock.release()
    except Exception as error:
        # shutdown the whole pool if we catch an error
        pool.shutdown(wait=False, cancel_futures=True)
        # print the error we got
        print(error)


def run(args: GithubArgs | GitlabArgs, type: ScanType) -> None:
    org = ''
    git_resource: GitResource = None

    if type == ScanType.Github:
        org = args.org
        git_resource = GithubResource(org)
    elif type == ScanType.Gitlab:
        org = args.group
        git_resource = GitlabResource(org)
    else:
        raise AttributeError(f'Unknown scan type {type}')

    print(f'Retrieving {org} repositories...')

    repo_urls = git_resource.get_repository_urls(args.visibility, args.no_archived)

    # create tmp directory for cloned repositories
    repo_path = args.repo_path if args.repo_path else f'{tempfile.gettempdir()}/{TEMP_DIR_NAME}/{org}'
    if not os.path.exists(repo_path):
        os.makedirs(repo_path)

    # setup multithreading
    lock = threading.Lock()
    pool = futures.ThreadPoolExecutor(max_workers=5)

    print('Launching secret scan (this operation can take a while to complete)...')

    # setup the report file with column
    if not os.path.exists(args.file):
        with open(args.file, 'w') as report_file:
            csv_writer = csv.writer(report_file)
            csv_writer.writerow(['repository', 'path', 'kind', 'line', 'valid', 'cleartext', 'hash'])

    # submit tasks to the thread pool
    for url in repo_urls:
        pool.submit(repository_scan, url, repo_path, report_file, pool, lock)

    # wait for all tasks to complete
    pool.shutdown(wait=True)

    # delete cloned repositories when cleanup is enabled
    if not args.no_clean_up:
        print('Cleaning cloned repositories...')
        shutil.rmtree(repo_path)


def cli() -> None:
    def add_global_arguments(parser: argparse.ArgumentParser):
        parser.add_argument('-f', '--file',
            type=str,
            default='report.csv',
            help='path to the CSV report file to generate',
        )
        parser.add_argument('-v', '--visibility',
            type=str,
            default='all',
            choices=['all', 'private', 'public'],
            help='repositories visibility',
        )
        parser.add_argument('--no-archived',
            default=False,
            action='store_true',
            help='do not scan archived repositories',
        )
        parser.add_argument('--repo-path',
            type=str,
            help='folder path to store repositories',
        )
        parser.add_argument('--no-clean-up',
            default=False,
            action='store_true',
            help='do not clean repositories downloaded after the scan',
        )

    parser = argparse.ArgumentParser(description='Scan secrets in organization repositories')
    subparsers = parser.add_subparsers(required=True, dest='git_software')

    # GitHub
    github_parser = subparsers.add_parser(ScanType.Github, help='scan a GitHub organization')
    github_parser.add_argument('-o', '--org',
        type=str,
        required=True,
        help='organization to scan'
    )
    add_global_arguments(github_parser)

    # GitLab
    gitlab_parser = subparsers.add_parser(ScanType.Gitlab, help='scan a GitLab group')
    gitlab_parser.add_argument('-o', '--group',
        type=str,
        required=True,
        help='group to scan'
    )
    add_global_arguments(gitlab_parser)

    # parse arguments
    args = GlobalArgs()
    parser.parse_args(namespace=args)

    # run script
    run(args, args.scan_type)


if __name__ == '__main__':
    cli()
