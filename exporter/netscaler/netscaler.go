package netscaler

import (
	"github.com/rerorero/netscaler-exporter/exporter/conf"
	"github.com/rerorero/netscaler-exporter/exporter/netscaler/nshttp"
	"github.com/rerorero/netscaler-exporter/exporter/netscaler/nssnmp"
)

type Netscaler interface {
	GetStats() (*NetscalerStats, []error)
	GetHost() string
}

type NetscalerStats struct {
	Http nshttp.NetscalerHttpStats
	Snmp nssnmp.NetscalerSnmpStats
}

type netscalerImpl struct {
	Http     nshttp.NetscalerHttp
	Snmp     nssnmp.NetscalerSnmp
	hostname string
}

func NewNetscalerClient(node conf.NetscalerNode) (Netscaler, error) {
	ns := &netscalerImpl{
		hostname: node.Host,
	}

	if node.EnableHttpStat {
		httpClient, err := nshttp.NewNetscalerHttpClient(node.Host, node.HTTPPort, node.Username, node.Password, node.TimeoutSec)
		if err != nil {
			return nil, err
		}
		ns.Http = httpClient
	}

	if node.EnableSnmpStat {
		ns.Snmp = nssnmp.NewNetscalerSnmp(node.Host, node.SNMPPort, node.SNMPCommunity, node.TimeoutSec)
	}

	return ns, nil
}

func (ns *netscalerImpl) GetStats() (*NetscalerStats, []error) {
	stats := &NetscalerStats{}
	errors := []error{}

	if ns.Http != nil {
		s, err := ns.Http.GetStats()
		if s != nil {
			stats.Http = *s
		}
		if err != nil {
			errors = append(errors, err)
		}
	}

	if ns.Snmp != nil {
		s, err := ns.Snmp.GetStats()
		if s != nil {
			stats.Snmp = *s
		}
		if err != nil {
			errors = append(errors, err)
		}
	}

	return stats, errors
}

func (ns *netscalerImpl) GetHost() string {
	return ns.hostname
}
