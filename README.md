# Offbeat

Offbeat automatically add worklogs into your Tempo timesheet.

## Usage

```shell
offbeat add [key1] [keyN...] [flags]

Flags:
  -e, --end string     Date to finish adding worklogs. Format: YYYY-MM-DD
  -h, --help           help for add
  -s, --start string   Date to start adding worklogs from. Format: YYYY-MM-DD
  -w, --weekends       Include weekends.
```

## Examples

Automatically add timesheets for ABC-123 and ABC-456 between 01-01-2023 and 07-01-2023:
```shell
offbeat add ABC-123 ABC-456 -s 2021-01-01 -e 2021-01-07
``` 

## Configuration

Create a `offbeat.yaml` file in the `$HOME/.config/offbeat/` directory containing:

```yaml
tempo:
  userId: <USER_ID> # Tempo userId to add worklogs for
  apiKey: <API_KEY> # Tempo API key with view and manage worklog scopes
atlassian:
  instance: <INSTANCE> # Base URL of your Atlassian cloud instance
  email: <EMAIL> # Email address used to log into your Atlassian instance
  apiKey: <API_KEY> # Atlassian API key
```

## Build

You must have Go 1.20+ installed. Clone the repository and run: 

```shell
go build 
```