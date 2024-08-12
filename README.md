# hikup

hikup is a Go program that automatically updates and recreates Docker containers when new images are available.

## Features

- Automatically update running Docker containers to use the latest image version
- Support for configuration file to include or exclude specific containers
- Dynamic configuration reloading via SIGHUP
- Logging to syslog for easy integration with system log management

## Installation

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/hikup.git
   ```

2. Build the program:
   ```
   cd hikup
   go build
   ```

3. (Optional) Install the program system-wide:
   ```
   sudo cp hikup /usr/local/bin/
   ```

## Usage

hikup can be run with the following options:

- `-a`: Recreate all running containers
- `-c <path>`: Specify a path to a configuration file

These options are mutually exclusive.

### Examples

1. Update all running containers:
   ```
   hikup -a
   ```

2. Use a configuration file:
   ```
   hikup -c /etc/hikup.conf
   ```

3. Run as a system service (after setting up the systemd unit file):
   ```
   sudo systemctl start hikup
   ```

## Configuration File

The configuration file can be in JSON or YAML format. It supports the following options:

- `include_containers`: List of container names to include for updates
- `exclude_containers`: List of container names to exclude from updates

Using `"*"` in the `include_containers` list will update all containers except those in the `exclude_containers` list.

### Example Configuration (YAML)

```yaml
include_containers:
  - "*"
exclude_containers:
  - database
  - cache
```

This configuration will update all containers except "database" and "cache".

## Logging

hikup logs to syslog. You can view the logs using journalctl or by checking your system's syslog files.

To view logs with journalctl:

```
journalctl -u hikup.service
```

## Reloading Configuration

To reload the configuration without restarting the service, send a SIGHUP signal:

```
kill -SIGHUP $(pgrep hikup)
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

[MIT License](LICENSE)
