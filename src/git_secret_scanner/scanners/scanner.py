from __future__ import annotations
from typing import Self

from git_secret_scanner import report


class BaseScanner:
    def __init__(self: Self, directory: str, repository: str) -> None:
        self.directory = directory
        self.repository = repository
        self._results: set[report.Secret] = set()

    def get_results(self: Self) -> set[report.Secret]:
        return self._results

    def scan(self: Self) -> None:
        msg = '"scan" method not implemented'
        raise NotImplementedError(msg)
