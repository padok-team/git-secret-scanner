from .gitleaks import GitleaksScanner, GITLEAKS_IGNORE_TAG
from .trufflehog import TrufflehogScanner, TRUFFLEHOG_IGNORE_TAG


def is_ignored(line: str) -> bool:
    return GITLEAKS_IGNORE_TAG in line or TRUFFLEHOG_IGNORE_TAG in line


__all__ = ['GitleaksScanner', 'TrufflehogScanner', 'is_ignored']
