package nssnmp

import "testing"
import "github.com/soniah/gosnmp"

func newTestSnmp() *netscalerSnmpImpl {
	return NewNetscalerSnmp("127.0.0.1", 9161, "public", 30).(*netscalerSnmpImpl)
}

func TestGetOidSuccess(t *testing.T) {
	if testing.Short() {
		t.Skip("Skip SNMP tests in short mode becase of lack of mocked SNMP server.")
	}
	snmp := newTestSnmp()

	busy := int64(-1)
	vserverName := ""
	unknown := int64(-1)
	info := map[string]*oidGet{
		oidHttpErrorServerBusy:  &oidGet{buf: &busy, valueType: gosnmp.Counter64},
		oidVserverNames + ".1":  &oidGet{buf: &vserverName, valueType: gosnmp.OctetString},
		"1.1.1.1.1.1.1.1.1.1.1": &oidGet{buf: &unknown, valueType: gosnmp.Counter32},
	}

	err := snmp.getOids(info)
	if err != nil {
		t.Fatal(err)
	}

	if !info[oidHttpErrorServerBusy].gotten || busy != 12345678910 {
		t.Fatalf("failed - %s", oidHttpErrorServerBusy)
	}
	if !info[oidVserverNames+".1"].gotten || vserverName != "lbvserver1" {
		t.Fatalf("failed - %s", oidVserverNames+".1")
	}
	if info["1.1.1.1.1.1.1.1.1.1.1"].gotten || unknown != -1 {
		t.Fatalf("failed - unknown")
	}
}
