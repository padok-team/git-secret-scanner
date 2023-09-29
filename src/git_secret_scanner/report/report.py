from __future__ import annotations
from typing import Self, Any
from types import TracebackType

import enum

import csv
import hashlib
from pathlib import Path

from .secret_kind import SecretKind


class ReportColumn(enum.StrEnum):
    Repository = 'repository'
    Path = 'path'
    Kind = 'kind'
    Line = 'line'
    Valid = 'valid'
    Cleartext = 'cleartext'
    Fingerprint = 'fingerprint'


def read_report(path: str) -> set[ReportSecret]:
    if Path(path).exists():
        with Path(path).open('r') as file:
            try:
                if ','.join(list(ReportColumn)) not in file.readline():
                    msg = f'file {path} is not a valid report file: wrong header'
                    raise ValueError(msg)
                reader = csv.DictReader(file, fieldnames=list(ReportColumn))
                return {ReportSecret.from_dict(secret) for secret in reader}
            except csv.Error as error:
                msg = f'file {path} is not a valid report file'
                raise ValueError(msg) from error
    msg = f'report file {path} does not exist'
    raise FileExistsError(msg)


class ReportWriter:
    def __init__(self: Self, path: str, force_recreate: bool = False) -> None:
        self.__file = Path(path).open('w+' if force_recreate else 'a+', newline='')  # noqa: SIM115
        self.writer = csv.DictWriter(self.__file, fieldnames=list(ReportColumn))
        self.__write_header()

    def __write_header(self: Self) -> None:
        if self.__file.tell() == 0:
            self.writer.writeheader()

    def add_secret(self: Self, secret: ReportSecret) -> None:
        self.writer.writerow(secret.to_dict())

    def __enter__(self: Self) -> Self:
        return self

    def __exit__(self: Self,
        exc_type: type[BaseException] | None,
        exc_val: BaseException | None,
        exc_tb: TracebackType | None,
    ) -> bool:
        self.__file.close()
        return False


class ReportSecret:
    def __init__(self: Self,
        repository: str,
        path: str,
        kind: SecretKind,
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
        self.fingerprint = fingerprint

        if fingerprint is None and cleartext is not None:
            self.fingerprint = hashlib.sha256(
                bytes(cleartext, 'utf-8')
                    # make sure that fingerprints are equal even when the strings are not
                    # encoded in the exact same way
                    .decode('unicode-escape')
                    .strip()
                    # gitleaks and trufflhog do not keep as many "-" in the cleartext of PrivateKeys
                    # we strip them to make sure we end up with the same fingerprints
                    .strip('-')
                    .encode('utf-8'),
            ).hexdigest()
        elif fingerprint is None and cleartext is None:
            msg = 'SecretReport cannot have both "None" cleartext and fingerprint'
            raise AttributeError(msg)

    def to_dict(self: Self) -> dict[ReportColumn, Any]:
        return {
            ReportColumn.Repository: self.repository,
            ReportColumn.Path: self.path,
            ReportColumn.Kind: self.kind,
            ReportColumn.Line: self.line,
            ReportColumn.Valid: self.valid,
            ReportColumn.Cleartext: self.cleartext,
            ReportColumn.Fingerprint: self.fingerprint,
        }

    def __hash__(self: Self) -> int:
        return hash((self.repository, self.path, self.fingerprint))

    def __eq__(self: Self, other: ReportSecret) -> bool:
        if not isinstance(other, type(self)):
            raise NotImplementedError
        return self.repository == other.repository \
            and self.path == other.path \
            and (self.kind == other.kind
                or self.kind == SecretKind.Generic
                or other.kind == SecretKind.Generic)  \
            and (self.line == other.line
                or self.line is None
                or other.line is None) \
            and (self.valid == other.valid
                or self.valid is None
                or other.valid is None) \
            and self.fingerprint == other.fingerprint

    def __str__(self: Self) -> str:
        return ('SecretReport('
            f'repository={self.repository},'
            f'path={self.path},'
            f'kind={self.kind},'
            f'line={self.line},'
            f'valid={self.valid},'
            f'cleartext={self.cleartext},'
            f'fingerprint={self.fingerprint})')

    def __repr__(self: Self) -> str:
        return ('SecretReport('
            f'repository={self.repository},'
            f'path={self.path},'
            f'kind={self.kind},'
            f'line={self.line},'
            f'valid={self.valid},'
            f'cleartext={self.cleartext},'
            f'fingerprint={self.fingerprint})')

    @classmethod
    def merge(cls: type[ReportSecret], first: ReportSecret, second: ReportSecret) -> ReportSecret:
        if first == second:
            return cls(
                repository=first.repository,
                path=first.path,
                kind=(first.kind if first.kind != SecretKind.Generic else second.kind),
                line=(first.line if first.line is not None else second.line),
                valid=(first.valid if first.valid is not None else second.valid),
                cleartext=first.cleartext,
                fingerprint=first.fingerprint,
            )
        msg = 'non equal secrets cannot be merged'
        raise AttributeError(msg)

    @classmethod
    def from_dict(cls: type[ReportSecret], row: dict[ReportColumn, str]) -> ReportSecret:
        return cls(
            repository=row[ReportColumn.Repository],
            path=row[ReportColumn.Path],
            kind=SecretKind[row[ReportColumn.Kind]],
            line=(int(row[ReportColumn.Line])
                if row[ReportColumn.Line] != ''
                else None),
            valid=(bool(row[ReportColumn.Valid])
                if row[ReportColumn.Valid] != ''
                else None),
            cleartext=(row[ReportColumn.Cleartext]
                if ReportColumn.Cleartext in row
                else None),
            fingerprint=row[ReportColumn.Fingerprint],
        )
