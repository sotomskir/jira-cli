# jira-cli
Atlassian's Jira REST API client written in Go for use with CI/CD tools to automate workflows.

## Install
### Linux
```bash
sudo curl -L "https://github.com/sotomskir/jira-cli/releases/download/0.5.0/jira-cli-$(uname -s)-$(uname -m)" -o /usr/local/bin/jira-cli && sudo chmod +x /usr/local/bin/jira-cli
```
### Other platforms
You can download pre-build binary here: https://github.com/sotomskir/jira-cli/releases

## Usage
First you must login to Jira server. jira-cli will store your credentials in configuration file ~/.jira-cli.yaml
```bash
jira-cli login
```

```
JIRA server URL: https://someuser.atlassian.net
Username: admin
Password: 
Success, Logged in to: http://someuser.atlassian.net as: admin
```
To login in non interactive mode use `--server`, `--user` and `--password` flags. 
Alternately you can use jira-cli without login if you set environment variables: 
`JIRA_SERVER_URL`, `JIRA_USER` and `JIRA_PASSWORD`. After successful login you can start using commands.
To list available commands type:
```
jira-cli --help
```
```
jira-cli is a CLI for Atlassian Jira REST API.
It can be used with CI/CD pipelines to automate workflow.

Usage:
  jira-cli [command]

Available Commands:
  help        Help about any command
  issue       Manage Jira issues
  login       Login to Atlassian Jira server
  project     Manage Jira projects
  version     Manage Jira versions

Flags:
      --config string   config file (default is $HOME/.jira-cli.yaml)
  -h, --help            help for jira-cli
      --no-color        Disable ANSI color output
  -t, --toggle          Help message for toggle

Use "jira-cli [command] --help" for more information about a command.
```
Each command has it's own help:
```
jira-cli version --help
```
```
Manage Jira versions

Usage:
  jira-cli version [command]

Aliases:
  version, v

Available Commands:
  create      Create new version
  ls          List versions for project
  release     Set version status to Released

Flags:
  -h, --help   help for version

Global Flags:
      --config string   config file (default is $HOME/.jira-cli.yaml)
      --no-color        Disable ANSI color output

Use "jira-cli version [command] --help" for more information about a command.
```

## Bash completion
To load completion run
```bash
. <(jira-cli completion)
```

To configure your bash shell to load completions for each session add following line to your bashrc
 ~/.bashrc or ~/.profile
```bash
. <(jira-cli completion)
```

## Issue transition workflows
issue transition command require workflow definition in yaml file. 
Default filename is `workflow.yaml` and can be overridden by --workflow flag.
### workflow structure
```yaml
workflow:
  source status:
    target status: transition name
    default: default transition name
```
### example workflow definition
```yaml
workflow:
  code review:
    default: ready to test
  in test:
    done: done
    default: bug found
  to do:
    rejected: reject
    default: start progress
  in progress:
    default: code review
  done:
    default: reopen
  rejected:
    default: reopen
```
### corresponding Jira workflow
![Alt text](docs/workflow.png?raw=true "Example Jira workflow")

### Workflow from env variable
Alternatively workflow file content can be passed by `JIRA_WORKFLOW_CONTENT` environment variable.
```yaml
export JIRA_WORKFLOW_CONTENT=$(cat  <<- EOM
workflow:
  code review:
    default: ready to test
  in test:
    done: done
    default: bug found
  to do:
    rejected: reject
    default: start progress
  in progress:
    default: code review
  done:
    default: reopen
  rejected:
    default: reopen
EOM
)
```
