package nssnmp

import (
	"fmt"
	"log"

	"github.com/pkg/errors"
	"github.com/soniah/gosnmp"
)

type oidGet struct {
	buf       interface{}
	gotten    bool
	valueType gosnmp.Asn1BER
}

func (ns *netscalerSnmpImpl) getOids(targets map[string]*oidGet) error {
	err := ns.snmp.Connect()
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Could not connect snmp - %s", ns.snmp.Target))
	}
	defer ns.snmp.Conn.Close()

	oids := []string{}
	for oid, _ := range targets {
		oids = append(oids, oid)
	}

	packets, err := ns.snmp.Get(oids)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Failed to get oids=%v from %s", oids, ns.snmp.Target))
	}

	pdumap := map[string]*gosnmp.SnmpPDU{}
	for i := range packets.Variables {
		pdumap[packets.Variables[i].Name] = &packets.Variables[i]
	}

	for _, oid := range oids {
		targets[oid].gotten = false
		// gosnmp appends dot prefix to PDU Name
		pdu, ok := pdumap["."+oid]

		if ok && pdu.Type == targets[oid].valueType {
			switch targets[oid].valueType {
			case gosnmp.Integer, gosnmp.Counter32, gosnmp.Gauge32, gosnmp.Counter64, gosnmp.Uinteger32:
				targetBuf := targets[oid].buf.(*int64)
				*targetBuf = gosnmp.ToBigInt(pdu.Value).Int64()
			case gosnmp.OctetString:
				targetBuf := targets[oid].buf.(*string)
				*targetBuf = string(pdu.Value.([]byte))
			default:
				return fmt.Errorf("Not implemented for oid type=%v: oid=%s, from %s", pdu.Type, oid, ns.snmp.Target)
			}
			targets[oid].gotten = true
		} else if ok && pdu.Type == gosnmp.NoSuchInstance {
			log.Printf("warning : No such oid=%s from %s", oid, ns.snmp.Target)
		} else if ok {
			log.Printf("warning : Unexpected SNMP type(%v) for oid=%s from %s", pdu.Type, oid, ns.snmp.Target)
		} else {
			log.Printf("warning : Unfortunately failed to get oid=%s from %s", oid, ns.snmp.Target)
		}
	}
	return nil
}
