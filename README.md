# PagerDuty on call report generator

Generate a report for the oncall rotation using PagerDuty API

## Installation

Using homebrew

```bash
brew tap rogersole/tap
brew install pd-report
```

## Usage

```bash
Generate on-call rotation reports automatically
from your PagerDuty account.

Usage:
  pd-report [command]

Available Commands:
  help        Help about any command
  report      generates the report(s) for the given schedule(s) id(s)
  schedules   list schedules on PagerDuty
  services    list services on PagerDuty
  teams       list teams on PagerDuty
  users       list users on PagerDuty

Flags:
      --config string   configuration file (default is $HOME/.pd-report-config.yml)
  -h, --help            help for pd-report

Use "pd-report [command] --help" for more information about a command.
```

And `report` specific flags:

```bash
Usage:
  pd-report report [flags]

Flags:
  -h, --help                   help for report
  -o, --output-format string   pdf, console (default "console")
  -s, --schedules strings      schedule ids to report (comma-separated with no spaces), or 'all' (default [all])

Global Flags:
      --config string   configuration file (default is ~/.pd-report-config.yml)
```

## Configuration

The configuration must be a `.yml` file (specified by the `--config` flag) with the following content:

```yml
# PagerDuty auth token
pdAuthToken: 12345

# Rotation starting time
rotationStartHour: 08:00:00

# Currency to be shown when calculating the daily rates
currency: £

# Prices assigned to each day type
# Three types must be specified: weekday, weekend and bankholiday
rotationPrices:
  - type: weekday
    price: 1
  - type: weekend
    price: 2
  - type: bankholiday
    price: 2

# List of users to be considered for the rotation
# Each one should be specifying a calendar for the bank holidays
# and the ID defined in PagerDuty
rotationUsers:
  - name: "User 1"
    holidaysCalendar: uk
    userId: P11A11B
  - name: "User 2"
    holidaysCalendar: uk
    userId: P22A22B
  - name: "Roger Solé Navarro"
    holidaysCalendar: sp_premia
    userId: P33A33B

# List of schedule IDs that can be ignored when generating the report
schedulesToIgnore:
  - SCHED_1
  - SCHED_2
```

> The default configuration file is `~/pd-report-config.yml`.
> To specify the path and the filename, the flag `--config` can be used on commands execution.


## Known limitations

- `list-people` command: the pagination is not implemented (yet)
- `report` command: no way to specify the output folder/filename for the pdf report
- `calendars`:
  - there are only 2018 and 2019 uk bank holiday calendars defined
  - there is no possibility to load external calendars (yet)


## Roadmap

- Add support for calendars loaded from outside the application
- Generate a summary for all the users involved in the rotations reported at the end
