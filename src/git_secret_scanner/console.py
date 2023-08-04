from __future__ import annotations

import rich


stdout = rich.console.Console()
stderr = rich.console.Console(stderr=True)


def print(message: str):
    stdout.print(message)


def print_error(message: str):
    stderr.print(f'[red]{message}[/red]')


def print_error_and_fail(message: str):
    print_error(message)
    exit(1)
