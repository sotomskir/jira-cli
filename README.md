# jira-cli
Atlassian's Jira REST API client written in Go for use with CI/CD tools to automate workflows.

## Install
### Linux
```bash
sudo curl -L "https://github.com/sotomskir/jira-cli/releases/download/0.1.0/jira-cli-$(uname -s)-$(uname -m)" -o /usr/local/bin/jira-cli
```
### Other platforms
You can download pre-build binary here: https://github.com/sotomskir/jira-cli/releases

## Usage
First you must login to Jira server. jira-cli will store your credentials in configuration file ~/.jira-cli.yaml
```bash
jira-cli login https://jira-server.example.com
```

```
Login to http://jira-server.example.com
Username: jira-user
Password:
Success, Logged in to http://jira-server.example.com
```
To login in non interactive mode use -u and -p flags. After successful login you can start using commands.
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
