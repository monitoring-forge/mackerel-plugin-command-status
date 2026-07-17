# mackerel-plugin-command-status
Mackerel custom plugin for exec command and get exit status as metrics

## Usage

```
Usage:
  mackerel-plugin-command-status mackerel-plugin-command-status [OPTIONS] -- command args1 args2 ...

Application Options:
      --timeout= Timeout to wait for command finished (default: 30s)
  -n, --name=    Metrics name
  -q, --quiet    Suppress error output of sub command
  -v, --version  Show version

Help Options:
  -h, --help     Show this help message
```

Sample

```
$ mackerel-plugin-command-status --name update-cache -- /path/to/cmd-fetch-cache
command-status.time-taken.update-cache  0.008186        1606958816
command-status.exit-code.update-cache  0       1606958816
```

```
$ mackerel-plugin-command-status -n test -- false 
2025/04/15 16:17:07 Command false exit with err: exit status 1
command-status.time-taken.test  0.001459        1744701427
command-status.exit-code.test   1       1744701427

$ mackerel-plugin-command-status --quiet --name test -- false
command-status.time-taken.test  0.001316        1744701374
command-status.exit-code.test   1       1744701374
```

`--quiet` (`-q`) suppresses error output of sub command.

```
$ mackerel-plugin-command-status -n sleep --timeout 3s -- sleep 30
2020/12/03 10:29:11 Command sleep timeout. killed
command-status.time-taken.sleep 3.016507        1614308642
command-status.exit-code.sleep  137     1614308642
```

Configure for mackerel-agent to use mackerel-agent as crontab and get their exit status code as metrics.

```
[plugin.metrics.update-cache]
command = "/path/to/mackerel-plugin-command-status --name update-cache --timeout 10s -- /path/to/cmd-fetch-cache"
```

and set a monitor with `warning(critical)  > 0` and `maxCheckAttempts` appropriate to your SLO.

## Install

Please download release page or `mkr plugin install monitoring-forge/mackerel-plugin-command-status`.
