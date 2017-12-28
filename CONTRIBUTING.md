# Test with mocked SNMP server
```
pip install snmpsim
snmpsimd.py --data-dir=./misc --agent-udpv4-endpoint=127.0.0.1:9161

# Started mock server

```
Now we can run SNMP tests.

```
go test -v ./...
```

# Test without mock
In short mode, SNMP tests are ignored.
```
go test -v -short ./...
```
