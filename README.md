# git-secret-scanner

This tool aims to find secrets and credentials in git repositories owned by Organizations or Groups using the libraries [TruffleHog](https://trufflesecurity.com/) &amp; [Gitleaks](https://gitleaks.io/).

> **Warning**
> 
> This tool is only designed for Linux and MacOS.
> The current version only supports GitLab and GitHub.

## Why this tool?

Trufflehog and Gitleaks are already designed to find secrets in git repositories. So you may wonder *"what is the purpose of a tool combining both scanners?"* 

These two tools have both their own strenghts and weaknesses:
- TruffleHog is very effective at classifying different secrets, but cannot find them all. It relies on detectors that can easily detect specific types of secrets, but not general secrets or general API keys.
- Gitleaks is able to find many more secrets, but is not as good as Trufflehog at classification. It contains fewer detectors and relies on string entropy to detect potential secrets that are not found by its detectors.

We designed this tool to combine the strenghts of both previous tools in order to find as many secrets as possible and to have an efficient classification of these secrets.

## Requirements

`git-secret-scanner` requires the following tools to work:
- [Python 3](https://www.python.org/downloads/) (>= 3.11)
- [pip](https://pip.pypa.io/en/stable/installation/)
- [git](https://git-scm.com/book/fr/v2/D%C3%A9marrage-rapide-Installation-de-Git)
- [TruffleHog](https://github.com/trufflesecurity/trufflehog) (>= 3.0)
- [Gitleaks](https://github.com/gitleaks/gitleaks) (>= 8.0)

You can easily check that all requirements are met with the commands below:

```bash
$ python --version
$ pip --version
$ git --version
$ trufflehog --version
$ gitleaks version
```

## Installation

### Using `pip`

The simplest way to install `git-secret-scanner` is with `pip`.

```bash
$ pip install git-secret-scanner
```

Then export your personal access token for ([GitHub](https://docs.github.com/en/enterprise-server@3.4/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token) or [GitLab](https://docs.gitlab.com/ee/user/profile/personal_access_tokens.html)):

```bash
# GitHub
$ export GITHUB_TOKEN="<token>"
# GitLab
$ export GITLAB_TOKEN="<token>"
```

### From source

1. Clone the repository

```bash
$ git clone https://github.com/padok-team/git-secret-scanner.git # using https
# or
$ git clone git@github.com:padok-team/git-secret-scanner.git # using ssh
$ cd git-secret-scanner
```

2. Install the Python requirements to run the tool

```bash
$ pip install -r requirements.txt
```

3. Add your personal access token ([GitHub](https://docs.github.com/en/enterprise-server@3.4/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token) / [GitLab](https://docs.gitlab.com/ee/user/profile/personal_access_tokens.html)) for your git SaaS in your environment variables:

```bash
# GitHub
$ export GITHUB_TOKEN="<token>"
# GitLab
$ export GITLAB_TOKEN="<token>"
```

> GitHub tokens require the `repo` scope, GitLab tokens require both `read_api` and `read_repository` scopes.

## Usage

To get detailed usage information about how to use this tool, run 

```bash
$ git-secret-scanner --help
```

### Examples

#### GitHub

Scan the repositories of the organization *my-org* and write the output in the file *output.csv*: 

```bash
$ git-secret-scanner github -o <my-org>
```

#### GitLab

Scan the repositories of the group *my-group* and write the output in the file *output.csv*: 

```bash
$ git-secret-scanner gitlab -o <my-org>
```

## Questions?

Open an issue to contact us or to give us suggestions. We are open to collaboration!

## License

[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
