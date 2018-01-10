# netscaler-exporter

[![Build Status](https://travis-ci.org/rerorero/netscaler-exporter.svg?branch=master)](https://travis-ci.org/rerorero/netscaler-exporter)

Simple Prometheus exoprter server that collects Citrix Netscaler stats.

# To run it
```
go get github.com/rerorero/netscaler-exporter
echo 'bind_port: 8080' > ./config.yml
netscaler-exporter --conf.file=./config.yml

curl localhost:8080/metrics
```

Testing on localhost.
```
# Run mocked SNMP server on localhost
pip install snmpsim
snmpsimd.py --data-dir=./misc --agent-udpv4-endpoint=127.0.0.1:9161

# Start exporter
go run main.go --conf.file=./misc/snmpconf.yml

# Get metrics
curl localhost:8080/metrics
```

# Configuration file
```
# Port to bind the exporter
bind_port: 8080

# Netscaler hosts
netscaler:
  static_targets:
    - host: 192.168.10.10
      http_port: 8080       # REST API port
      username: foo         # to authorize REST API
      password: bar         # to authorize REST API
      snmp_port: 9161
      snmp_community: public
      enable_http: yes      # If set to no, metrics that retrieved from REST API are not exported.
      enable_snmp: yes      # If set to no, metrics that retrieved from SNMP are not exported.
      timeout: 100
    - host: 192.168.10.20   # You can configure multiple hosts
      snmp_port: 9161
      snmp_community: public
      enable_http: no
      enable_snmp: yes
```

# Exported Metrics
See [metrics.go](exporter/metrics.go)

# Using Docker
```
echo 'bind_port: 8080' > /var/nsx/config.yml
docker run -p 8080:8080 -v /var/nsx/config.yml:/etc/nsx/nsxconf.yml rerorero/netscaler-exporter:latest

```
