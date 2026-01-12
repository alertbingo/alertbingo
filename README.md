## alert.bingo

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
2026/01/12 19:45:52 Required flags "dashboard, site, service, name, token" not set
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
