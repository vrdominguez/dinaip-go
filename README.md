# dinaIP GO

This is an updated version of the dinaIP client written in Go, developed as an opportunity to explore and experiment with Golang.

## Configuration

To configure the dinaIP GO client, you need to create a YAML file with the following format:

```yaml
username: userDinahostingAPI
password: passwordDinahostingAPI

logs:
  path: /path/to/file.log
  level: DEBUG

zones:
  mydomain.tld:
    - mysubdomain1
    - mysubdomain2
  mydomain2.tld:
    - subdomain
```

## Example system.d configuration file

```systemd
[Unit]
Description=dinaIP GO
Wants=network-online.target
After=network-online.target

[Service]
User=systemUserForDinaIP
Group=systemGroupForDinaIP
ExecStart=/path/to/dinapi-go -c /path/to/config.yaml


[Install]
WantedBy=multi-user.target
```

Make sure to replace userDinahostingAPI and passwordDinahostingAPI with your actual username and password for the Dinahosting API. Additionally, specify the log file path and the desired logging level. Lastly, provide the domains and their respective subdomains that you want to manage.

## Use it

Before running the dinaIP GO client, ensure that you have a valid configuration file with the correct values, including a log file path where the program can write logs.

### Runing from source

1. Clone this repository.
2. Open a terminal and navigate to the directory where the program is located.
3. Run the code with the command `go run ./main.go -c /path/to/config.yaml`, replacing /path/to/config.yaml with the path to your configuration file.

### Running in the Command Line

1. Download the compiled program from the [lastest release](https://github.com/vrdominguez/dinaip-go/releases/latest) on the [project's GitHub page](https://github.com/vrdominguez/dinaip-go/releases) or compile it yourself.
2. Open a terminal and navigate to the directory where the program is located.
3. Execute the command `./dinapi-go -c /path/to/config.yaml`, replacing /path/to/config.yaml with the path to your configuration file.

### Running as a Service (with systemd)

1. Download the compiled program from the [lastest release](https://github.com/vrdominguez/dinaip-go/releases/latest) on the [project's GitHub page](https://github.com/vrdominguez/dinaip-go/releases) or compile it yourself.
2. Move the binary to `/opt/dinaip-go/` or other path of your choice.
3. Create a yaml configuration for the program. (I placed mine at `/etc/dinaIP.yaml`)
4. If you don't want to run the service as root (you may not run it as root), create an user and a group for it.
5. Create your system.d file from example on this README.
6. Save the unit file with a .service extension (e.g., `dinaip-go.service`).
7. Move the unit file to the appropriate location for systemd unit files (e.g., `/etc/systemd/system/`).
8. Reload the systemd daemon to load the new unit file: `sudo systemctl daemon-reload`.
9. Enable the service to start on boot: `sudo systemctl enable dinaip-go`.
10. Start the service: `sudo systemctl start dinaip-go`.
11. Verify that the service is running: `sudo systemctl status dinaip-go`.

Please note that you may need to adjust the unit file and paths according to your specific setup and requirements.

Also, you may want to configure log rotation with your prefered log rotating system. Here is an example with logrotate, for logs placed at `/var/log/dinaip.log`:

```logrotate
/var/log/dinaip.log {
    daily
    rotate 30
    copytruncate
    missingok
    notifempty
    create 644 systemUserForDinaIP systemGroupForDinaIP
}
```

In this example:

- ```daily``` indicates that log rotation should occur daily.
- ```rotate 30``` specifies that only 30 rotated log files should be kept.
- ```copytruncate``` copies the current log file and truncates the original file, allowing ongoing logging without disrupting the program that writes to the log.
- ```missingok``` ignores an error if the log file does not exist.
- ```notifempty``` does not rotate the log file if it is empty.
- ```create 600 systemUserForDinaIP systemGroupForDinaIP``` creates a new log file with the specified permissions and ownership if it does not exist, using the provided systemUserForDinaIP user and systemGroupForDinaIP group.

To use this configuration, save it as /etc/logrotate.d/dinaip or any other appropriate name. Logrotate will automatically rotate the /var/log/dinaip.log file based on the specified settings.
