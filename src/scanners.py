from __future__ import annotations

import os, subprocess
import json

from secret import SecretReport


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


class TrufflehogScanner(Scanner):
    def scan(self) -> None:
        try:
            scan_results = subprocess.check_output([
                'trufflehog', 'filesystem', '--no-update', '--json', self.directory,
            ], stderr=subprocess.DEVNULL)
        except subprocess.CalledProcessError:
            raise Exception('Failed to run trufflehog scan')
        
        if len(scan_results) == 0:
            return []
        for secret in scan_results.decode('utf-8').split('\n')[:-1]:
            secret = json.loads(secret)
            result = SecretReport(
                repository=self.repository,
                path=secret['SourceMetadata']['Data']['Filesystem']['file'].removeprefix(f'{self.directory}/'),
                kind=secret['DetectorName'],
                line=None,
                valid=(secret['Verified'] or None),
                cleartext=secret['Raw'],
            )
            self._results.append(result)


class GitleaksScanner(Scanner):
    def __map_rule(self, rule: str) -> str:
        return GITLEAKS_TO_TRUFFLEHOG[rule] if rule in GITLEAKS_TO_TRUFFLEHOG else rule

    def scan(self) -> None:
        report_path = os.path.join(self.directory, 'gitleaks.json')

        try:
            subprocess.call([
                'gitleaks', 'detect', '--no-git', '-s', self.directory, '-f', 'json', '--exit-code', '0', '-r', report_path,
            ], stderr=subprocess.DEVNULL)
            scan_results = subprocess.check_output(['cat', report_path], stderr=subprocess.DEVNULL)
        except subprocess.CalledProcessError:
            raise Exception('Failed to run gitleaks scan')

        if len(scan_results) == 0:
            return []

        scan_results = json.loads(scan_results)

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
