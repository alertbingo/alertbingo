## [alert.bingo](https://alert.bingo) CLI

Download the latest CLI: https://github.com/alertbingo/alertbingo/releases

Place the downloaded binary at /usr/local/bin/alertbingo (or another location in your $PATH).

The hoststats command posts the following checks for the host:

* CPU usage
* Disk Inodes Used %
* Disk Space Used %
* Memory usage
* Uptime


### Configuration Example

Create a shell script (e.g. send_checks.sh) to configure and run the CLI:
```bash
#!/bin/bash

export ALERTBINGO_TOKEN="xxx"
export ALERTBINGO_SITE="staging"
export ALERTBINGO_DASHBOARD="MyDashboard"
export ALERTBINGO_INACTIVE_ESCALATE="10m"

/usr/local/bin/alertbingo hoststats --service "$(hostname -s)"
```

Make the script executable:

```bash
chmod +x send_checks.sh
```

Scheduling with Cron - To send checks every minute, add a cron job:

```
* * * * * /path/to/send_checks.sh
```



## Command summary

### check 

```
NAME:
   alertbingo check - Send a check to Alert Bingo

OPTIONS:
   --dashboard string, -d string    Dashboard name [$ALERTBINGO_DASHBOARD]
   --site string, -s string         Site identifier (e.g., myapp-prod) [$ALERTBINGO_SITE]
   --service string                 Service name (e.g., postgres) [$ALERTBINGO_SERVICE]
   --name string, -n string         Check name (e.g., postgres-rds-space-free) [$ALERTBINGO_NAME]
   --alert-level string, -l string  Alert level: ok, warn, or alert (default: "ok") [$ALERTBINGO_ALERT_LEVEL]
   --message string, -m string      Optional long-form status message [$ALERTBINGO_MESSAGE]
   --value string, -v string        Short-form status value [$ALERTBINGO_VALUE]
   --inactive-expire string         Optional duration string for inactive expiry (e.g., 48h or 30m) [$ALERTBINGO_INACTIVE_EXPIRE]
   --inactive-escalate string       Optional duration string for inactive escalation (e.g., 1h or 30m) [$ALERTBINGO_INACTIVE_ESCALATE]
   --highlighted string             Optional highlighted status [$ALERTBINGO_HIGHLIGHTED]
   --token string, -t string        API Bearer token [$ALERTBINGO_TOKEN]
   --api-url string                 API URL (default: "https://app.alert.bingo/api/v1/checks") [$ALERTBINGO_API_URL]
   --help, -h                       show help
```


### hoststats
```
NAME:
   alertbingo hoststats - Send host statistics checks (memory, uptime, CPU) to Alert Bingo

OPTIONS:
   --dashboard string, -d string  Dashboard name [$ALERTBINGO_DASHBOARD]
   --site string, -s string       Site identifier (e.g., myapp-prod) [$ALERTBINGO_SITE]
   --service string               Service name (e.g., host) [$ALERTBINGO_SERVICE]
   --message string, -m string    Optional long-form status message [$ALERTBINGO_MESSAGE]
   --inactive-expire string       Optional duration string for inactive expiry (e.g., 48h or 30m) [$ALERTBINGO_INACTIVE_EXPIRE]
   --inactive-escalate string     Optional duration string for inactive escalation (e.g., 1h or 30m) [$ALERTBINGO_INACTIVE_ESCALATE]
   --highlighted string           Optional highlighted status [$ALERTBINGO_HIGHLIGHTED]
   --token string, -t string      API Bearer token [$ALERTBINGO_TOKEN]
   --api-url string               API URL (default: "https://app.alert.bingo/api/v1/checks") [$ALERTBINGO_API_URL]
   --help, -h                     show help
```

### certcheck

Check SSL/TLS certificate expiry for one or more URLs. Warns when certificates expire within 14 days, and alerts when expired.

```
NAME:
   alertbingo certcheck - Check SSL/TLS certificate expiry for one or more URLs

USAGE:
   alertbingo certcheck [options] <url> [url...]

OPTIONS:
   --dashboard string, -d string  Dashboard name [$ALERTBINGO_DASHBOARD]
   --site string, -s string       Site identifier (e.g., myapp-prod) [$ALERTBINGO_SITE]
   --name string, -n string       Check name (e.g., ssl) [$ALERTBINGO_NAME]
   --message string, -m string    Optional long-form status message [$ALERTBINGO_MESSAGE]
   --inactive-expire string       Optional duration string for inactive expiry (e.g., 48h or 30m) [$ALERTBINGO_INACTIVE_EXPIRE]
   --inactive-escalate string     Optional duration string for inactive escalation (e.g., 1h or 30m) [$ALERTBINGO_INACTIVE_ESCALATE]
   --highlighted string           Optional highlighted status [$ALERTBINGO_HIGHLIGHTED]
   --token string, -t string      API Bearer token [$ALERTBINGO_TOKEN]
   --api-url string               API URL (default: "https://app.alert.bingo/api/v1/checks") [$ALERTBINGO_API_URL]
   --timeout duration             Timeout for TLS connection (default: 10s) [$ALERTBINGO_TIMEOUT]
   --help, -h                     show help
```

Note: The `--name` flag sets the check Name, while the Service field is automatically set to the URL/hostname being checked.

Example:
```bash
alertbingo certcheck --dashboard MyDashboard --site prod --name ssl \
  https://example.com https://api.example.com
```

### urlcheck

Check URL availability, HTTP status code, and optionally verify response body content.

```
NAME:
   alertbingo urlcheck - Check URL availability, status code, and optionally body content

USAGE:
   alertbingo urlcheck [options] <url> [expected_code] [expected_body]

OPTIONS:
   --dashboard string, -d string  Dashboard name [$ALERTBINGO_DASHBOARD]
   --site string, -s string       Site identifier (e.g., myapp-prod) [$ALERTBINGO_SITE]
   --name string, -n string       Check name (e.g., http) [$ALERTBINGO_NAME]
   --message string, -m string    Optional long-form status message [$ALERTBINGO_MESSAGE]
   --inactive-expire string       Optional duration string for inactive expiry (e.g., 48h or 30m) [$ALERTBINGO_INACTIVE_EXPIRE]
   --inactive-escalate string     Optional duration string for inactive escalation (e.g., 1h or 30m) [$ALERTBINGO_INACTIVE_ESCALATE]
   --highlighted string           Optional highlighted status [$ALERTBINGO_HIGHLIGHTED]
   --token string, -t string      API Bearer token [$ALERTBINGO_TOKEN]
   --api-url string               API URL (default: "https://app.alert.bingo/api/v1/checks") [$ALERTBINGO_API_URL]
   --timeout duration             Timeout for HTTP request (default: 10s) [$ALERTBINGO_TIMEOUT]
   --help, -h                     show help
```

Note: The `--name` flag sets the check Name, while the Service field is automatically set to the URL being checked.

Arguments:
- `url` - The URL to check (required)
- `expected_code` - Expected HTTP status code (optional, defaults to any 2xx)
- `expected_body` - String that must be present in the response body (optional)

Examples:
```bash
# Check URL returns 2xx
alertbingo urlcheck --dashboard MyDashboard --site prod --name http \
  https://example.com/health

# Check URL returns specific status code
alertbingo urlcheck --dashboard MyDashboard --site prod --name http \
  https://example.com/health 200

# Check URL returns 200 and body contains "OK"
alertbingo urlcheck --dashboard MyDashboard --site prod --name http \
  https://example.com/health 200 "OK"
```
