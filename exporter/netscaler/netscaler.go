package netscaler

import (
	"github.com/rerorero/netscaler-vpx-exporter/exporter/netscaler/nshttp"
	"github.com/rerorero/netscaler-vpx-exporter/exporter/netscaler/nssnmp"
)

type Netscaler interface {
	GetStats() error
}

type netscalerImpl struct {
	Http nshttp.NetscalerHttp
	Snmp nssnmp.NetscalerSnmp
}

func NewNetscalerClient(
	host string,
	httpPort int,
	username string,
	password string,
	enableHttpStat bool,
	enableSmtpStat bool,
	timeoutSec int,
) (Netscaler, error) {
	ns := &netscalerImpl{}

	if enableHttpStat {
		httpClient, err := nshttp.NewNetscalerHttpClient(host, httpPort, username, password, timeoutSec)
		if err != nil {
			return nil, err
		}
		ns.Http = httpClient
	}

	return ns, nil
}

func (ns *netscalerImpl) GetStats() error {
	// TODO
	return nil
}
