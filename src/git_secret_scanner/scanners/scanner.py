import git_secret_scanner.report as report


class BaseScanner():
    def __init__(self, directory: str, repository: str) -> None:
        self.directory = directory
        self.repository = repository
        self._results: set[report.Secret] = set()

    def get_results(self) -> set[report.Secret]:
        return self._results

    def scan(self) -> None:
        raise NotImplementedError('"scan" method not implemented')
