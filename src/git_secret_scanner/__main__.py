import sys

from .cli import cli


if __name__ == '__main__':
    sys.exit(cli(prog_name='git-secret-scanner'))
