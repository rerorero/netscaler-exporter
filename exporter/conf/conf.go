package conf

import (
	"fmt"
	"io/ioutil"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type NetscalerNode struct {
	Host           string `yaml:"host"`
	HTTPPort       int    `yaml:"http_port,omitempty"`
	Username       string `yaml:"username,omitempty"`
	Password       string `yaml:"password,omitempty"`
	SNMPPort       int    `yaml:"snmp_port,omitempty"`
	SNMPCommunity  string `yaml:"snmp_community,omitempty"`
	EnableHttpStat bool   `yaml:"enable_http"`
	EnableSnmpStat bool   `yaml:"enable_snmp"`
	TimeoutSec     int    `yaml:"timeout,omitempty"`
}

type NetscalerConf struct {
	StaticTargets []NetscalerNode `yaml:"static_targets,omitempty"`
}

type Conf struct {
	Netscaler NetscalerConf `yaml:"netscaler"`
	BindPort  int           `yaml:"bind_port"`
}

func NewConfFromFile(confPath string) (*Conf, error) {
	buf, err := ioutil.ReadFile(confPath)
	if err != nil {
		return nil, fmt.Errorf("Fatal: Could not read config - %s \n", err)
	}

	return NewConfFromYaml(buf)
}

func NewConfFromYaml(yamlbytes []byte) (*Conf, error) {
	conf := &Conf{}
	err := yaml.Unmarshal(yamlbytes, conf)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Fatal: Could not parse config(%s)", string(yamlbytes)))
	}

	for i, _ := range conf.Netscaler.StaticTargets {
		if conf.Netscaler.StaticTargets[i].HTTPPort == 0 {
			conf.Netscaler.StaticTargets[i].HTTPPort = 80
		}
		if conf.Netscaler.StaticTargets[i].SNMPPort == 0 {
			conf.Netscaler.StaticTargets[i].SNMPPort = 161
		}
		if conf.Netscaler.StaticTargets[i].TimeoutSec == 0 {
			conf.Netscaler.StaticTargets[i].TimeoutSec = 15
		}
	}

	return conf, nil
}
