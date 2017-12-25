package conf

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type NetscalerNode struct {
	Host     string `yaml:host`
	HTTPPort int    `yaml:http_port`
	Username string `yaml:username`
	Password string `yaml:password`
}

type NetscalerConf struct {
	Targets []NetscalerNode `yaml:"targets,omitempty"`
}

type Conf struct {
	Netscaler NetscalerConf `yaml:netscaler`
}

func NewConfFrom(confPath string) (*Conf, error) {
	buf, err := ioutil.ReadFile(confPath)
	if err != nil {
		return nil, fmt.Errorf("Fatal: Could not read config - %s \n", err)
	}

	conf := &Conf{}
	err2 := yaml.Unmarshal(buf, conf)
	if err2 != nil {
		return nil, fmt.Errorf("Fatal: Could not parse config(%s) - %s \n", confPath, err)
	}

	return conf, nil
}
