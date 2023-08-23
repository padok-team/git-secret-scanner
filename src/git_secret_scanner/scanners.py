from __future__ import annotations

from typing import TypedDict, Any

import os
import subprocess
import json

from git_secret_scanner.secret import SecretReport


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


class Scanner():
    def __init__(self, directory: str, repository: str) -> None:
        self.directory = directory
        self.repository = repository
        self._results: list[SecretReport] = []

    def get_results(self) -> list[SecretReport]:
        return self._results

    def scan(self) -> None:
        raise NotImplementedError('"scan" method not implemented')


TrufflehogReportItem = TypedDict('TrufflehogReportItem', {
    'SourceMetadata': Any,
    'DetectorName': str,
    'Verified': bool,
    'Raw': str,
})

class TrufflehogScanner(Scanner):
    def scan(self) -> None:
        proc = subprocess.run([
                'trufflehog', 'filesystem',
                    '--no-update',
                    '--json',
                    self.directory,
            ],
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
        )

        if proc.returncode != 0:
            error = RuntimeError(f'trufflehog scan failed for {self.repository}')
            error.add_note(proc.stderr.decode('utf-8'))
            raise error

        report = proc.stdout.decode('utf-8')

        if len(report) == 0:
            return

        for raw_secret in report.split('\n')[:-1]:
            secret: TrufflehogReportItem = json.loads(raw_secret)
            result = SecretReport(
                repository=self.repository,
                path=secret['SourceMetadata']['Data']['Filesystem']['file'].removeprefix(f'{self.directory}/'),
                kind=secret['DetectorName'],
                line=None,
                valid=(secret['Verified'] if len(str(secret['Verified'])) > 0 else None),
                cleartext=secret['Raw'],
            )
            self._results.append(result)


GitleaksReportItem = TypedDict('GitleaksReportItem', {
    'RuleID': str,
    'File': str,
    'StartLine': int,
    'Secret': str,
})

GitleaksReport = list[GitleaksReportItem]


class GitleaksScanner(Scanner):
    def __map_rule(self, rule: str) -> str:
        return GITLEAKS_TO_TRUFFLEHOG[rule] if rule in GITLEAKS_TO_TRUFFLEHOG else rule

    def scan(self) -> None:
        report_path = os.path.join(self.directory, 'gitleaks.json')

        proc = subprocess.run([
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

        with open(report_path, 'r') as report:
            raw_scan_results = report.read()

        # remove the report file to make sure it is not read by other scanners
        os.remove(report_path)

        if len(raw_scan_results) == 0:
            return

        scan_results: GitleaksReport = json.loads(raw_scan_results)

        for secret in scan_results:
            result = SecretReport(
                repository=self.repository,
                path=secret['File'].removeprefix(f'{self.directory}/'),
                kind=self.__map_rule(secret['RuleID']),
                line=secret['StartLine'],
                valid=None,
                cleartext=secret['Secret'],
            )
            self._results.append(result)
