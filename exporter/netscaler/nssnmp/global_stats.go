package nssnmp

import (
	"github.com/soniah/gosnmp"
)

type SnmpGlobalStats struct {
	ServerBusyErrors int64
}

func (ns *netscalerSnmpImpl) getSnmpGlobalStats() (*SnmpGlobalStats, error) {
	stats := SnmpGlobalStats{}

	info := map[string]*oidGet{
		oidHttpErrorServerBusy: &oidGet{buf: &stats.ServerBusyErrors, valueType: gosnmp.Counter64},
	}

	err := ns.getOids(info)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}
