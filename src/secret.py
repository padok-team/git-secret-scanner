from __future__ import annotations

import hashlib


class SecretReport():
    def __init__(
        self,
        repository: str,
        path: str,
        kind: str,
        line: int | None,
        valid: bool | None,
        cleartext: str,
    ) -> None:
        self.repository = repository
        self.path = path
        self.kind = kind
        self.line = line
        self.valid = valid
        self.cleartext = cleartext
        self.hash = hashlib.sha256(cleartext.encode('utf-8')).hexdigest()

    def __hash__(self) -> int:
        return hash((self.repository, self.path, self.hash))

    def __eq__(self, other) -> bool:
        if not isinstance(other, type(self)):
            raise NotImplementedError
        return self.repository == other.repository \
            and self.path == other.path \
            and self.hash == other.hash

    def __str__(self) -> str:
        return f'SecretReport{{ \
            repository={self.repository}, \
            path={self.path}, \
            kind={self.kind}, \
            line={self.line}, \
            valid={self.valid}, \
            cleartext{self.cleartext}, \
            hash={self.hash}}}'

    def __repr__(self) -> str:
        return f'SecretReport{{ \
            repository={self.repository}, \
            path={self.path}, \
            kind={self.kind}, \
            line={self.line}, \
            valid={self.valid}, \
            cleartext{self.cleartext}, \
            hash={self.hash}}}'

    @staticmethod
    def merge(first: SecretReport, second: SecretReport) -> SecretReport:
        if first == second:
            return SecretReport(
                repository=first.repository,
                path=first.path,
                kind=(first.kind if first.kind != 'GenericApiKey' else second.kind),
                line=(first.line if first.line else second.line),
                valid=(first.valid if first.valid is not None else second.valid),
                cleartext=first.cleartext,
            )
        raise AttributeError('Non equal secrets cannot be merged')
