# Mirakurun Exporter

[![CI](https://github.com/coord-e/mirakurun_exporter/actions/workflows/ci.yml/badge.svg)](https://github.com/coord-e/mirakurun_exporter/actions/workflows/ci.yml)

Prometheus exporter for [Mirakurun](https://github.com/Chinachu/Mirakurun) metrics.
Pre-built binaries are available at [the releases](https://github.com/coord-e/mirakurun_exporter/releases).
Container images are available at [the packages](https://github.com/coord-e?tab=packages&repo_name=mirakurun_exporter).

## Usage

```
$ mirakurun_exporter -h
usage: mirakurun_exporter --exporter.mirakurun-url=EXPORTER.MIRAKURUN-URL [<flags>]

Flags:
  -h, --help                Show context-sensitive help (also try --help-long and --help-man).
      --web.systemd-socket  Use systemd socket activation listeners instead of port listeners (Linux
                            only).
      --web.listen-address=:9110 ...
                            Addresses on which to expose metrics and web interface. Repeatable for
                            multiple addresses.
      --web.config.file=""  [EXPERIMENTAL] Path to configuration file that can enable TLS or
                            authentication.
      --web.telemetry-path="/metrics"
                            Path under which to expose metrics.
      --exporter.mirakurun-url=EXPORTER.MIRAKURUN-URL
                            URL of the Mirakurun instance.
      --exporter.status     Whether to export metrics from /api/status.
      --exporter.tuners     Whether to export metrics from /api/tuners.
      --exporter.programs   Whether to export metrics from /api/programs.
      --exporter.services   Whether to export metrics from /api/services.
      --log.level=info      Only log messages with the given severity or above. One of: [debug, info, warn, error]
      --log.format=logfmt   Output format of log messages. One of: [logfmt, json]
      --version             Show application version.
```

### Example

To run against a Mirakurun instance running at `localhost:40772`:

```
$ mirakurun_exporter --exporter.mirakurun-url=http://localhost:40772/
```

## Build

```
$ make build
```
