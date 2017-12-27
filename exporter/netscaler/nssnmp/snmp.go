package nssnmp

import (
	"time"

	"github.com/soniah/gosnmp"
)

type NetscalerSnmp interface {
	GetStats() (*NetscalerSnmpStats, error)
}

type NetscalerSnmpStats struct {
	VServers   map[string]*SnmpVserverStats
	Global     *SnmpGlobalStats
	SnmpHealth bool
}

type netscalerSnmpImpl struct {
	snmp *gosnmp.GoSNMP
}

const (
	// ref. https://github.com/citrix/netscaler-snmp-oid-reference/blob/master/docs/index.md
	oidHttpErrorServerBusy   = "1.3.6.1.4.1.5951.4.1.1.48.61"
	oidVserverTable          = "1.3.6.1.4.1.5951.4.1.3.1"
	oidVserverNames          = "1.3.6.1.4.1.5951.4.1.3.1.1.1"
	oidLbVserverAverageTTFBs = "1.3.6.1.4.1.5951.4.1.3.6.1.5"
)

func NewNetscalerSnmp(
	host string,
	port int,
	community string,
	timeoutSec int,
) NetscalerSnmp {
	snmp := &gosnmp.GoSNMP{
		Target:    host,
		Port:      uint16(port),
		Community: community,
		Version:   gosnmp.Version2c,
		Timeout:   time.Duration(timeoutSec) * time.Second,
		// Logger:    log.New(os.Stdout, "debug :", log.Flags()),
	}

	return &netscalerSnmpImpl{
		snmp: snmp,
	}
}

func (ns *netscalerSnmpImpl) GetStats() (*NetscalerSnmpStats, error) {
	global, err := ns.getSnmpGlobalStats()
	if err != nil {
		return nil, err
	}

	vservers, err := ns.getSnmpVserverStats()
	if err != nil {
		return nil, err
	}

	return &NetscalerSnmpStats{
		Global:     global,
		VServers:   vservers,
		SnmpHealth: true,
	}, nil
}
