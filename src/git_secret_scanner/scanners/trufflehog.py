from __future__ import annotations
from typing import Any

import json
import subprocess

import git_secret_scanner.report as report

from .scanner import BaseScanner


class TrufflehogReportItem():
    def __init__(self,
        file: str,
        line: int | None,
        detector_name: str,
        verified: bool | None,
        raw: str,
    ):
        self.file = file
        self.line = line
        self.detector_name = detector_name
        self.verified = verified
        self.raw = raw

    @staticmethod
    def from_json(json_dict: dict[str, Any]) -> TrufflehogReportItem:
        return TrufflehogReportItem(
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


class TrufflehogScanner(BaseScanner):
    def scan(self) -> None:
        proc = subprocess.run([
                # truffle filesystem is no longer used as it does not compute line numbers
                # in the right way
                'trufflehog', 'git',
                    '--no-update',
                    # this works since we shallow clone repositories with depth = 1
                    '--max-depth', '1',
                    '--json',
                    f'file://{self.directory}',
            ],
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
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
            result = report.Secret(
                repository=self.repository,
                path=item.file.removeprefix(f'{self.directory}/'),
                kind=item.detector_name,
                line=item.line,
                valid=item.verified,
                cleartext=item.raw,
            )
            self._results.add(result)
