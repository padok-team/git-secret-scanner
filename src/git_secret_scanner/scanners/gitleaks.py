from __future__ import annotations
from typing import Self, Any

import json
from pathlib import Path
import subprocess

from git_secret_scanner import report

from .scanner import BaseScanner



GITLEAKS_TO_TRUFFLEHOG = {
    'slack-web-hook': 'SlackWebhook',
    'generic-api-key': 'GenericApiKey',
    'github-pat': 'Github',
    'jwt': 'JWT',
    'flickr-access-token': 'Flickr',
    'aws-access-token': 'AWS',
    'twilio-api-key': 'Twilio',
    'slack-access-token': 'Slack',
    'facebook': 'FacebookOAuth',
    'private-key': 'PrivateKey',
    'gcp-api-key': 'GoogleCloudApiKey',
    'mailgun-signing-key': 'MailgunSignKey',
    'mailgun-pub-key': 'MailgunPublicKey',
    'mailgun-private-api-token': 'MailgunApiToken',
}


class GitleaksReportItem:
    def __init__(self: Self,
        rule_id: str,
        file: str,
        start_line: int,
        secret: str,
    ) -> None:
        self.rule_id = rule_id
        self.file = file
        self.start_line = start_line
        self.secret = secret

    @staticmethod
    def from_json(json_dict: dict[str, Any]) -> GitleaksReportItem:
        return GitleaksReportItem(
            rule_id=json_dict['RuleID'],
            file=json_dict['File'],
            start_line=json_dict['StartLine'],
            secret=json_dict['Secret'],
        )


class GitleaksScanner(BaseScanner):
    def __map_rule(self: Self, rule: str) -> str:
        return GITLEAKS_TO_TRUFFLEHOG[rule] if rule in GITLEAKS_TO_TRUFFLEHOG else rule

    def scan(self: Self) -> None:
        report_path = Path(self.directory) / 'gitleaks.json'

        proc = subprocess.run([  # noqa: S603, S607
                'gitleaks', 'detect',
                    '--no-git',
                    '--source', self.directory,
                    '--report-format', 'json',
                    '--report-path', report_path,
                    '--no-banner',
                    '--no-color',
                    '--log-level', 'error',
                    '--exit-code', '0',
            ],
            stdout=subprocess.DEVNULL,
            stderr=subprocess.PIPE,
        )

        if proc.returncode != 0:
            error = RuntimeError(f'gitleaks scan failed for {self.repository}')
            error.add_note(proc.stderr.decode('utf-8'))
            raise error

        with Path(report_path).open('r') as report_file:
            raw_scan_results = report_file.read()

        # remove the report file to make sure it is not read by other scanners
        Path(report_path).unlink()

        if len(raw_scan_results) == 0:
            return

        scan_results: list[GitleaksReportItem] = json.loads(
            raw_scan_results,
            object_hook=GitleaksReportItem.from_json,
        )

        for item in scan_results:
            result = report.Secret(
                repository=self.repository,
                path=item.file.removeprefix(f'{self.directory}/'),
                kind=self.__map_rule(item.rule_id),
                line=item.start_line,
                valid=None,
                cleartext=item.secret,
            )
            self._results.add(result)
