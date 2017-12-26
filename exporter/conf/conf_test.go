package conf

import (
	"reflect"
	"testing"
)

func TestParseConf(t *testing.T) {
	s := `
netscaler:
  static_targets:
    - host: 192.168.10.10
      http_port: 8080
      username: foo
      password: bar
      snmp_port: 9161
      snmp_community: public
      enable_http: true
      enable_snmp: false
      timeout: 100
    - host: 192.168.10.20
      username: aaa
      password: bbb
      enable_http: no
      enable_snmp: yes
`
	actual, err := NewConfFromYaml([]byte(s))
	if err != nil {
		t.Fatal(err.Error())
	}

	expected := &Conf{
		Netscaler: NetscalerConf{
			StaticTargets: []NetscalerNode{
				NetscalerNode{
					Host:           "192.168.10.10",
					HTTPPort:       8080,
					Username:       "foo",
					Password:       "bar",
					SNMPPort:       9161,
					SNMPCommunity:  "public",
					EnableHttpStat: true,
					EnableSnmpStat: false,
					TimeoutSec:     100,
				},
				NetscalerNode{
					Host:           "192.168.10.20",
					HTTPPort:       80,
					Username:       "aaa",
					Password:       "bbb",
					SNMPPort:       161,
					EnableHttpStat: false,
					EnableSnmpStat: true,
					TimeoutSec:     15,
				},
			},
		},
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Unexpected parsing result:\n%v\n%v", expected, actual)
	}
}

func TestParseEmptyNetscalerConf(t *testing.T) {
	s := `netscaler: {}`
	actual, err := NewConfFromYaml([]byte(s))
	if err != nil {
		t.Fatal(err)
	}

	expected := &Conf{
		Netscaler: NetscalerConf{
			StaticTargets: nil,
		},
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Unexpected parsing result:\n%v\n%v", expected, actual)
	}
}
