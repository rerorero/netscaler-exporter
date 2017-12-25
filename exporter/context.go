package exporter

import (
	"github.com/pkg/errors"
	"github.com/rerorero/netscaler-vpx-exporter/exporter/conf"
	"github.com/rerorero/netscaler-vpx-exporter/exporter/netscaler"
)

type Context struct {
	Conf       *conf.Conf
	Netscalers []netscaler.Netscaler
}

func NewContext(config *conf.Conf) (*Context, error) {
	nsary := []netscaler.Netscaler{}
	for i, nsconf := range config.Netscaler.Targets {
		ns, err := netscaler.NewNetscalerClient(
			nsconf.Host,
			nsconf.HTTPPort,
			nsconf.Username,
			nsconf.Password)

		if err != nil {
			return nil, errors.Wrap(err, "error : Failed to instantiate Netscaler client")
		}
		nsary[i] = ns
	}

	return &Context{
		Conf:       config,
		Netscalers: nsary,
	}, nil
}
