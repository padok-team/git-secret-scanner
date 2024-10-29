# git-secret-scanner

This tool aims to find secrets and credentials in git repositories owned by Organizations or Groups using the libraries [TruffleHog](https://trufflesecurity.com/) &amp; [Gitleaks](https://gitleaks.io/).

> **Warning**
> 
> This tool is only designed for Linux and MacOS.
> The current version only supports Gitlab and GitHub.

## Why this tool?

Trufflehog and Gitleaks are already designed to find secrets in git repositories. So you may wonder *"what is the purpose of a tool combining both scanners?"* 

These two tools have both their own strenghts and weaknesses:
- TruffleHog is very effective at classifying different secrets, but cannot find them all. It relies on detectors that can easily detect specific types of secrets, but not general secrets or general API keys.
- Gitleaks is able to find many more secrets, but is not as good as Trufflehog at classification. It contains fewer detectors and relies on string entropy to detect potential secrets that are not found by its detectors.

We designed this tool to combine the strenghts of both previous tools in order to find as many secrets as possible and to have an efficient classification of these secrets.

## Requirements

`git-secret-scanner` requires the following tools to work:
- [git](https://git-scm.com/book/fr/v2/D%C3%A9marrage-rapide-Installation-de-Git)
- [TruffleHog](https://github.com/trufflesecurity/trufflehog) (>= 3.82.13)
- [Gitleaks](https://github.com/gitleaks/gitleaks) (>= 8.21.1)

You can easily check that all requirements are met with the commands below:

```shell
git --version
trufflehog --version
gitleaks version
```

## Installation

### Using `homebrew`

The simplest way to install `git-secret-scanner` is with `homebrew`.

```shell
brew tap padok-team/tap
brew install git-secret-scanner
```

### With binary

Download the binary for your platform and OS on the [realeases page](https://github.com/zricethezav/gitleaks/releases).

### From source

1. Clone the repository

```shell
git clone https://github.com/padok-team/git-secret-scanner.git
cd git-secret-scanner
```

2. Build the binary

```shell
make build
```

## Usage

To get detailed usage information about how to use this tool, run 

```shell
git-secret-scanner --help
```

### Simple

Add a personal access token ([GitHub](https://docs.github.com/en/enterprise-server@3.4/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token) / [Gitlab](https://docs.gitlab.com/ee/user/profile/personal_access_tokens.html)) for your git SaaS in your environment variables.

```shell
# GitHub
export GITHUB_TOKEN="<token>"
# Gitlab
export GITLAB_TOKEN="<token>"
```

> GitHub tokens require the `repo` scope, Gitlab tokens require both `read_api` and `read_repository` scopes.

```shell
# With GITHUB_TOKEN set
git-secret-scanner github -o "<org>"
# With GITLAB_TOKEN set
git-secret-scanner gitlab -g "<group>"
```

### Ignore secrets

You can instruct `git-secret-scanner` to ignore some specific secrets in its results. This is useful to ignore false positives or to ignore secrets that have already been dealt with.

#### Ignore secrets with comments

`git-secret-scanner` understands Gitleaks and Trufflehog annotations to ignore secrets (`gitleaks:allow` and `trufflehog:ignore`). You can add a comment with one of these annotations on the line that has the secret to have `git-secret-scanner` ignore it.

#### Ignore secrets with fingerprints

To ignore specific fingerprints, create a file with a list of all secret fingerprints to ignore during the scan. A fingerprint is computed in the following way:

```
<repo_name>:<commit_sha>:<file>:<line>
```

Then run `git-secret-scanner` with the `-i` flag:

```shell
git-secret-scanner github -o "<org>" -i "<path_to_fingerprints_ignore_file>"
git-secret-scanner gitlab -g "<group>" -i "<path_to_fingerprints_ignore_file>"
```

### Baseline

`git-secret-scanner` supports using a previous report as a baseline for a scan. All previous secrets found in the baseline are ignored in the final report. This is useful to detect added secrets between two scans.

```shell
git-secret-scanner github -o "<org>" -b "<path_to_previous_report_csv>"
git-secret-scanner gitlab -g "<group>" -b "<path_to_previous_report_csv>"
```

## Questions?

Open an issue to contact us or to give us suggestions. We are open to collaboration.

## License

[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
