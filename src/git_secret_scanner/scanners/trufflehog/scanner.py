from __future__ import annotations
from typing import Self, Any

import json
import subprocess

from git_secret_scanner.report import ReportSecret, SecretKind
from git_secret_scanner.scanners.base import BaseScanner

from .mapping import TRUFFLEHOG_DETECTOR_TO_SECRET_KIND


TRUFFLEHOG_IGNORE_TAG = 'trufflehog:ignore'


class TrufflehogReportItem:
    def __init__(self: Self,
        file: str,
        line: int | None,
        detector_name: str,
        verified: bool | None,
        raw: str,
    ) -> None:
        self.file = file
        self.line = line
        self.detector_name = detector_name
        self.verified = verified
        self.raw = raw

    @classmethod
    def from_json(cls: type[TrufflehogReportItem], json_dict: dict[str, Any]) -> TrufflehogReportItem | dict[str, Any]:
        if 'SourceMetadata' in json_dict:
            return cls(
                file=json_dict['SourceMetadata']['Data']['Git']['file'],
                line=(json_dict['SourceMetadata']['Data']['Git']['line']
                    if 'line' in json_dict['SourceMetadata']['Data']['Git']
                    else None),
                detector_name=json_dict['DetectorName'],
                verified=(json_dict['Verified']
                    # validity checks are not relevent on PrivateKeys
                    if len(str(json_dict['Verified'])) > 0 and json_dict['DetectorName'] != 'PrivateKey'
                    else None),
                raw=json_dict['Raw'],
            )
        return json_dict


class TrufflehogScanner(BaseScanner):
    def __map_detector(self: Self, detector: str) -> SecretKind:
        return (TRUFFLEHOG_DETECTOR_TO_SECRET_KIND[detector]
            if detector in TRUFFLEHOG_DETECTOR_TO_SECRET_KIND
            else SecretKind.Generic)

    def scan(self: Self) -> None:
        proc = subprocess.run([  # noqa: S603, S607
                # truffle filesystem is no longer used as it does not compute line numbers
                # in the right way
                'trufflehog', 'git',
                    '--no-update',
                    # this works since we shallow clone repositories with depth = 1
                    '--max-depth', '1',
                    '--json',
                    f'file://{self.directory}',
            ],
            capture_output=True,
        )

        if proc.returncode != 0:
            error = RuntimeError(f'trufflehog scan failed for {self.repository}')
            error.add_note(proc.stderr.decode('utf-8'))
            raise error

        raw_report = proc.stdout.decode('utf-8')

        if len(raw_report) == 0:
            return

        for raw_item in raw_report.split('\n')[:-1]:
            item: TrufflehogReportItem = json.loads(
                raw_item,
                object_hook=TrufflehogReportItem.from_json,
            )
            result = ReportSecret(
                repository=self.repository,
                path=item.file.removeprefix(f'{self.directory}/'),
                kind=self.__map_detector(item.detector_name),
                line=item.line,
                valid=item.verified,
                cleartext=item.raw,
            )
            self._results.add(result)
