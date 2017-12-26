package nssnmp

import (
	"fmt"
	"log"
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
	err := ns.snmp.Connect()
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Could not connect snmp - %s", ns.snmp.Target))
	}
	defer ns.snmp.Conn.Close()

	serverMap, err := getVserverIndex(ns.snmp)
	if err != nil {
		return nil, err
	}
	inversed := inverseMap(serverMap)

	// get vserver Time to First Byte
	bttfs, err := getVserverTTFBs(ns.snmp, inversed)
	if err != nil {
		return nil, err
	}

	stats := map[string]*SnmpVserverStats{}
	for name, _ := range serverMap {
		stat := &SnmpVserverStats{Name: name}

		bttf, ok := bttfs[name]
		if ok {
			stat.LbVserverAverageTTFB = bttf
		}

		stats[stat.Name] = stat
	}

	return stats, nil
}

func getVserverIndex(snmp *gosnmp.GoSNMP) (map[string]int, error) {
	pdus, err := snmp.BulkWalkAll(oidVserverNames)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Failed to walk vservs in %s", snmp.Target))
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

func getVserverTTFBs(snmp *gosnmp.GoSNMP, severMap map[int]string) (map[string]int64, error) {
	bttfs := map[string]int64{}
	oids := []string{}
	for index, _ := range severMap {
		oids = append(oids, fmt.Sprintf("%s.%d", oidLbVserverAverageTTFBs, index))
	}

	packets, err := snmp.Get(oids)
	if err != nil {
		log.Printf("warning : Could not failed to retrive vserverTTFB from %s: %s", snmp.Target, err.Error())
	} else {
		for _, pdu := range packets.Variables {
			if pdu.Type != gosnmp.Gauge32 {
				log.Printf("warning : Unexpected type for vserverTTFB in %s: %v", snmp.Target, pdu.Type)
			} else {
				// retrieve vserver's index number
				index, err := tailOid(pdu.Name)
				if err != nil {
					return nil, err
				}
				name, ok := severMap[index]
				if ok {
					bi := gosnmp.ToBigInt(pdu.Value)
					bttfs[name] = bi.Int64()
				} else {
					log.Printf("warning : vserverTTFB is not found in %s: i=%d", snmp.Target, index)
				}
			}
		}
	}
	return bttfs, nil
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

func inverseMap(m map[string]int) map[int]string {
	inversed := map[int]string{}
	for key, value := range m {
		inversed[value] = key
	}
	return inversed
}
