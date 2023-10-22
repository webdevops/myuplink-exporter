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
      --log.debug                   debug mode [$LOG_DEBUG]
      --log.devel                   development mode [$LOG_DEVEL]
      --log.json                    Switch log output to json format [$LOG_JSON]
      --myuplink.url=               Url to myUplink API (default: https://api.myuplink.com) [$MYUPLINK_URL]
      --myuplink.auth.clientid=     ClientID from myUplink [$MYUPLINK_AUTH_CLIENTID]
      --myuplink.auth.clientsecret= ClientSecret from myUplink [$MYUPLINK_AUTH_CLIENTSECRET]
      --server.bind=                Server address (default: :8080) [$SERVER_BIND]
      --server.timeout.read=        Server read timeout (default: 5s) [$SERVER_TIMEOUT_READ]
      --server.timeout.write=       Server write timeout (default: 10s) [$SERVER_TIMEOUT_WRITE]

Help Options:
  -h, --help                        Show this help message
```

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
