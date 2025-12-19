# myUplink.com exporter

[![license](https://img.shields.io/github/license/webdevops/myuplink-exporter.svg)](https://github.com/webdevops/myuplink-exporter/blob/master/LICENSE)
[![DockerHub](https://img.shields.io/badge/DockerHub-webdevops%2Fmyuplink--exporter-blue)](https://hub.docker.com/r/webdevops/myuplink-exporter/)
[![Quay.io](https://img.shields.io/badge/Quay.io-webdevops%2Fmyuplink--exporter-blue)](https://quay.io/repository/webdevops/myuplink-exporter)

A Prometheus exporter for myuplink.com device metrics

Usage
-----

```
Usage:
  myuplink-exporter [OPTIONS]

Application Options:
      --log.level=[trace|debug|info|warning|error] Log level (default: info) [$LOG_LEVEL]
      --log.format=[logfmt|json]                   Log format (default: logfmt) [$LOG_FORMAT]
      --log.source=[|short|file|full]              Show source for every log message (useful for debugging and bug reports) [$LOG_SOURCE]
      --log.color=[|auto|yes|no]                   Enable color for logs [$LOG_COLOR]
      --log.time                                   Show log time [$LOG_TIME]
      --myuplink.url=                              Url to myUplink API (default: https://api.myuplink.com) [$MYUPLINK_URL]
      --myuplink.auth.clientid=                    ClientID from myUplink [$MYUPLINK_AUTH_CLIENTID]
      --myuplink.auth.clientsecret=                ClientSecret from myUplink [$MYUPLINK_AUTH_CLIENTSECRET]
      --myuplink.device.allowed-connectionstates=  Allowed device connection states (default: Connected) [$MYUPLINK_DEVICE_ALLOWED_CONNECTIONSTATES]
      --myuplink.device.calc-total-parameters=     Calculate total metrics for these parameters (eg. energey log parameters) [$MYUPLINK_DEVICE_CALC_TOTAL_PARAMETRS]
      --server.bind=                               Server address (default: :8080) [$SERVER_BIND]
      --server.timeout.read=                       Server read timeout (default: 5s) [$SERVER_TIMEOUT_READ]
      --server.timeout.write=                      Server write timeout (default: 10s) [$SERVER_TIMEOUT_WRITE]

Help Options:
  -h, --help                                       Show this help message
```

## Auth (Client id and secret)

Register an application at [myUplinkDEVS](https://dev.myuplink.com/login) to get
client id and secret for your application.

## HTTP Endpoints

| Endpoint   | Description                             |
|------------|-----------------------------------------|
| `/metrics` | Default prometheus golang metrics       |
| `/probe`   | Probe devices metrics from myuplink.com |

### `/probe`

| GET parameter | Default                   | Required | Multiple | Description                                     |
|---------------|---------------------------|----------|----------|-------------------------------------------------|
| `cache`       |                           | no       | no       | Cache time in time.Duration (eg. `30s` or `5m`) |


## Metrics

| Metric                                 | Description                              |
|----------------------------------------|------------------------------------------|
| `myuplink_system_info`                 | System information                       |
| `myuplink_system_device_info`          | Device information (belongs to a system) |
| `myuplink_system_device_point`         | Device parameter metric point            |
