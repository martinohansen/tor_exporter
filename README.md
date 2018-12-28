# tor_exporter
Small Prometheus exporter for the tor daemon.

## Usage

`tor_exporter (--tor.control-socket /var/run/tor/control | --tor.control-port localhost:9051)`

The default port for serving the metrics is `9105`. This can be changed with flag: `--web.listen-address`.

### Docker

A docker image is also provided.

```shell
docker run -d -p 9105:9105 -v /var/lib/tor/control_auth_cookie:/var/lib/tor/control_auth_cookie:ro,z \
martinohansen/tor_exporter --tor.control-port example:9051
```

## Installation

The typical way of installing in Go should work.

```shell
go install
```

## Metrics exposed

* `tor_connection_circuits` - the number of open circuits
* `tor_connection_streams` - the number of open streams
* `tor_connection_orconns` - the number of open ORConns
* `tor_traffic_read_bytes`, `tor_traffic_written_bytes` - cumulative inbound and outbound traffic

## To-do

* Let's cheat and assume the CircuitID is sequential. Instant counter!