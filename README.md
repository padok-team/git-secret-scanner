# git-secret-scanner

This tool aims to find secrets and credentials in git repositories using the libraries TrufflehogScanner &amp; GitleaksScanner

> :warning: Only works for linux based distribution

## Why this tool ?

During many pentesting missions, we used the following tools : 
- [TrufflehogScanner](https://github.com/trufflesecurity/trufflehog), a shell tool used to get credentials from git repositories
- [GitleaksScanner](https://github.com/gitleaks/gitleaks), a tool with the same purpose as the previous one

So you may ask *"why do we need another tool ? We already have two !"* 

Well .. during the pentests we did, we found that these two tools have strenght and weaknesses :
- TrufflehogScanner is very good to categorize the different secrets but doesn't find all of them
- Gitleaks finds a lot of secrets but sort the biggest majority of them in the same category (often as api key)

So we created this tool which combines the strenght of both previous tools to get all the secrets well categorized.

## Requirements

Here are the required tools to do the installation:
- [python3](https://www.python.org/downloads/)
- [pip](https://pip.pypa.io/en/stable/installation/)
- [git](https://git-scm.com/book/fr/v2/D%C3%A9marrage-rapide-Installation-de-Git)

Install these tools also :

- [TruffehogScanner](https://github.com/trufflesecurity/trufflehog)
- [GitleaksScanner](https://github.com/gitleaks/gitleaks)

To test if the tool are correctly installed run the following commands:

```bash
# Linux or MacOS
$ python3 --version
$ pip --version
$ git --version
```

> You can also try to run the tools TrufflehogScanner (running `trufflehog` command) et GitleaksScanner (running `gitleaks` command).

## Installation

Here are the installation steps to use the tool:

1. Clone the repository

```bash
$ git clone https://github.com/padok-team/git-secret-scanner.git # using https
# or
$ git clone git@github.com:padok-team/git-secret-scanner.git # using ssh
```

2. Install the requirements to run the tool

```bash
# Linux or MacOS
$ cd git-secret-scanner
$ pip install -r requirements.txt
```

3. Create a Github access key

This is for the proper functionning of the tool. You can follow this very clear [tutorial](https://docs.github.com/en/enterprise-server@3.4/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token) provided by Github.

4. Add your Github access key to your environment variables

```bash
# Linux or MacOS
$ GITHUB_TOKEN=<your_token_value>
$ export GITHUB_TOKEN
```

And it's done, you can see the **Usage** section to run the tool.

## Usage

To get information about how to use this tool you can first run `python3 main.py -h`.

### Example

Analyze the repositories of the organization *my-org* and write the output in the file *output.csv*: 

```bash
$ cd src/ # go in folder with the python files
$ python3 main.py --org my-org -f output.csv
```

## Releases

The releases must meet these standarts to be approve :
- [ ] each release is bound to the repository
- [ ] each release includes a changelog
- [ ] each release is versioned using [semantic versioning](https://semver.org/)
- [ ] artefacts are published and served by Github

## Questions ?

Open an issue to contact us or to give us suggestions. We are open to collaboration!

## License

[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)