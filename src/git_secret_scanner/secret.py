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
        cleartext: str | None,
        fingerprint: str | None = None,
    ) -> None:
        self.repository = repository
        self.path = path
        self.kind = kind
        self.line = line
        self.valid = valid
        self.cleartext = cleartext
        self.fingerprint = None
        if cleartext is not None:
            self.fingerprint = fingerprint or hashlib.sha256(cleartext.encode('utf-8')).hexdigest()

    def __hash__(self) -> int:
        return hash((self.repository, self.path, self.fingerprint))

    def __eq__(self, other) -> bool:
        if not isinstance(other, type(self)):
            raise NotImplementedError
        return self.repository == other.repository \
            and self.path == other.path \
            and (self.kind == other.kind
                or self.kind == 'GenericApiKey'
                or other.kind == 'GenericApiKey')  \
            and (self.line == other.line
                or self.line is None
                or other.line is None) \
            and (self.valid == other.valid
                or self.valid is None
                or other.valid is None) \
            and self.fingerprint == other.fingerprint

    def __str__(self) -> str:
        return ('SecretReport('
            f'repository={self.repository},'
            f'path={self.path},'
            f'kind={self.kind},'
            f'line={self.line},'
            f'valid={self.valid},'
            f'cleartext={self.cleartext},'
            f'fingerprint={self.fingerprint})')

    def __repr__(self) -> str:
        return ('SecretReport('
            f'repository={self.repository},'
            f'path={self.path},'
            f'kind={self.kind},'
            f'line={self.line},'
            f'valid={self.valid},'
            f'cleartext={self.cleartext},'
            f'fingerprint={self.fingerprint})')

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
        raise AttributeError('non equal secrets cannot be merged')
