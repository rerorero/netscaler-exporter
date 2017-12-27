# netscaler-vpx-exporter

[![Build Status](https://travis-ci.org/rerorero/netscaler-vpx-exporter.svg?branch=test)](https://travis-ci.org/rerorero/netscaler-vpx-exporter)

Simple Prometheus exoprter server that collects Citrix Netscaler VPX load balancer stats.

# To run it
```
TODO

```
Testing on localhost.
```
# Run mocked SNMP server on localhost
pip install snmpsim
snmpsimd.py --data-dir=./snmpTesting --agent-udpv4-endpoint=127.0.0.1:9161

# Start exporter
go run main.go --conf.file=./snmpTesting/conf.yml

# Get metrics
curl localhost:8080/metrics
```

# Configuration format
```
# Port to bind the exporter
bind_port: 8080

# Node information of Netscaler
netscaler:
  static_targets:
    - host: 192.168.10.10
      http_port: 8080       # REST API port
      username: foo         # to authorize REST API
      password: bar         # to authorize REST API
      snmp_port: 9161
      snmp_community: public
      enable_http: yes      # If set to no, metrics which retrieved from REST API are not exported.
      enable_snmp: yes      # If set to no, metrics which retrieved from SNMP are not exported.
      timeout: 100
    - host: 192.168.10.20   # You can configure multiple hosts
      snmp_port: 9161
      snmp_community: public
      enable_http: no
      enable_snmp: yes
```

TODO

# Exported Metrics
TODO

# Using Docker
TODO
