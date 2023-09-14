from typing import Self

from git_secret_scanner.report import ReportSecret


class BaseScanner:
    def __init__(self: Self, directory: str, repository: str) -> None:
        self.directory = directory
        self.repository = repository
        self._results: set[ReportSecret] = set()

    def get_results(self: Self) -> set[ReportSecret]:
        return self._results

    def scan(self: Self) -> None:
        msg = '"scan" method not implemented'
        raise NotImplementedError(msg)
