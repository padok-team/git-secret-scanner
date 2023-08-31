import sys

from git_secret_scanner.cli import cli


if __name__ == '__main__':
    sys.exit(cli(prog_name='git-secret-scanner'))
