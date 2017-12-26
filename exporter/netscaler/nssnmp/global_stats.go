package nssnmp

import (
	"fmt"
	"log"

	"github.com/soniah/gosnmp"

	"github.com/pkg/errors"
)

type SnmpGlobalStats struct {
	ServerBusyErrors int64
}

func (ns *netscalerSnmpImpl) GetSnmpGlobalStats() (*SnmpGlobalStats, error) {
	err := ns.snmp.Connect()
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Could not connect snmp - %s", ns.snmp.Target))
	}
	defer ns.snmp.Conn.Close()

	serverBusy, err := getServerBusy(ns.snmp)
	if err != nil {
		return nil, err
	}

	return &SnmpGlobalStats{
		ServerBusyErrors: serverBusy,
	}, nil
}

func getServerBusy(snmp *gosnmp.GoSNMP) (int64, error) {
	packets, err := snmp.Get([]string{oidHttpErrorServerBusy})
	if err != nil {
		return 0, errors.Wrap(err, fmt.Sprintf("Failed to get httpErrorServerBusy from %s", snmp.Target))
	}

	if len(packets.Variables) == 0 {
		log.Printf("warning : httpErrorServerBusy is not found in %s", snmp.Target)
		return 0, nil
	} else if len(packets.Variables) > 1 {
		log.Printf("warning : Unexpected httpErrorServerBusy schema in %s: %v", snmp.Target, packets.Variables)
		return 0, nil
	} else if packets.Variables[0].Type != gosnmp.Counter64 {
		log.Printf("warning : Unexpected type for httpErrorServerBusy in %s: %v", snmp.Target, packets.Variables[0].Type)
		return 0, nil
	}

	pdu := packets.Variables[0]
	bi := gosnmp.ToBigInt(pdu.Value)
	return bi.Int64(), nil
}
