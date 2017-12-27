package nssnmp

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/soniah/gosnmp"
)

type SnmpVserverStats struct {
	Name                 string
	LbVserverAverageTTFB int64
}

func (ns *netscalerSnmpImpl) GetSnmpVserverStats() (map[string]*SnmpVserverStats, error) {
	serverMap, err := ns.getVserverIndex()
	if err != nil {
		return nil, err
	}

	stats := map[string]*SnmpVserverStats{}
	oidTargets := map[string]*oidGet{}

	for name, serverIndex := range serverMap {
		stat := &SnmpVserverStats{Name: name}

		// stats oid mapping
		oidTargets[fmt.Sprintf("%s.%d", oidLbVserverAverageTTFBs, serverIndex)] = &oidGet{buf: &stat.LbVserverAverageTTFB, valueType: gosnmp.Gauge32}

		stats[name] = stat
	}

	// request SNMP
	err = ns.getOids(oidTargets)
	if err != nil {
		return nil, err
	}

	return stats, nil
}

func (ns *netscalerSnmpImpl) getVserverIndex() (map[string]int, error) {
	err := ns.snmp.Connect()
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Could not connect snmp - %s", ns.snmp.Target))
	}
	defer ns.snmp.Conn.Close()

	pdus, err := ns.snmp.BulkWalkAll(oidVserverNames)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Failed to walk vservs in %s", ns.snmp.Target))
	}

	m := map[string]int{}
	for _, pdu := range pdus {
		if pdu.Type != gosnmp.OctetString {
			return nil, fmt.Errorf("Unexpected snmp type of vservername - %v", pdu)
		}

		// retrieve vserver's index number
		index, err := tailOid(pdu.Name)
		if err != nil {
			return nil, err
		}

		name := string(pdu.Value.([]byte))
		m[name] = index
	}

	return m, nil
}

// retrieve tail of oid as an index number
func tailOid(oid string) (int, error) {
	parts := strings.Split(oid, ".")
	index, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		return -1, fmt.Errorf("Unexpected oid name - %s, %v", oid, parts)
	}
	return index, nil
}
