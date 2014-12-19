# Third Law

A small monitoring tool meant to stay out of your way. Specifically, it doesn't have to own a process to be able to work with it.

## Example

```json
{
  "frequency": 2000,
  "outputs": {
    "failure": [{"type": "stderr"}]
  },
  "checks": [{
    "name": "redis",
    "type": "shell",
    "command": "redis-cli",
    "dir": "/opt/redis/",
    "arguments": ["ping"],
    "out": "PONG\n",
    "recover": ["restart redis"]
  }],
  "actions": {
    "restart redis": {
      "type": "shell",
      "command": "/etc/init.d/redis",
      "arguments": ["start"]
    }
  }
}
```

The `include` directive can be used to specify a directory which contains additional configuration files. These files are limited to the `checks`, `check` and `actions` fields. A check in one file can reference an action in a different file.

# Configuration Options

## frequency
The time, in milliseconds, to run the checks. Checks are run in series, so the real "sleep" time is going to be the specified frequency + however long it takes to do all the checks and actions.

Defaults to 10 seconds (10000 milliseconds).

## include
Specifies a directory to load additional configuration files from. All files within the directory are loaded. These child files are limited to defining `checks`, `check` and `actions`

## outputs
Defines the outputs to send the results of an iteration to. Outputs can be sent for the case where all checks pass (`success`) or when at least one check fails (`failure`):

```json
{
  "outputs": {
    "success": [{"type": "stdout"}],
    "failure": [
      {"type": "stderr"},
      {
        "type": "file",
        "path": "/var/log/opt/thirdlaw.log",
        "truncate": true
      }
    ]
  }
}
```
### stdout
Writes the results to stdout

### stderr
Writes the results to stderr

### file
Writes the results to a specified file. The file output accepts the following configuration values:

- `path`: the filepath to write results to (defaults to failures.log)
- `truncate`: whether or not to truncate the file before each write (defaults to false)

## checks and check
`checks` and `check` define the code to execute on each iteration. The two fields are only different in that `checks` is an array of `check`.

All `checks` accept a `recover` option which is an array of `action` names to run in case of failure.

### http
Makes an HTTP request. Any error or a response with a status code of 300 or more will result in a failure. The http check accepts the following configuration values:

- `address`: the full address (scheme, host, port, path) to make the request to (defaults to http://127.0.0.1/)
- `timeout`: the timeout, in milliseconds, to wait before getting a response (defaults to 5000).

### shell
Invokes the shell and runs the specified command. Any error running the command, include an exit code not equal to 0, will result in a failure. It's also possible to specify the expected output. The shell check accepts the following configuration values:

- `command`: REQUIRED command to run (no default)
- `arguments`: array of arguments to pass to the command (defaults to none)
- `dir`: the worker directory to use (defaults to thirdlaw's working directory)
- `out`: expected stdout text (defaults to ignoring any output)

## actions
Actions are invoked when a check fails. Actions are run as invoked in a blocking manner. As a general rule, you'll want your actions to invoke scripts with launch whatever process has failed as a daemon. All actions accept the following configurations:

- `retries`: how many times to retry the action should it fail (defaults to 0)
- `delay`: how long to wait, in milliseconds, between failed retries (detauls to 1000)

Remember, if you're frequence is set to 10000 and you have a single check with a single action that has a retry of 5 and a delay of 1000, the worst case check frequency is 15 seconds (10000 + 1000 * 5) as everything happens synchronously.

### shell
Invokes the shell and runs the specified command. The shell action accepts the following configuration values:

- `command`: REQUIRED command to run (no default)
- `arguments`: array of arguments to pass to the command (defaults to none)
- `dir`: the worker directory to use (defaults to thirdlaw's working directory)
