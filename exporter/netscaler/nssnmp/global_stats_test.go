package nssnmp

import "testing"

func newTestSnmp() NetscalerSnmp {
	return NewNetscalerSnmp("127.0.0.1", 9161, "public", 30)
}

func TestGetSnmpGlobalStatsSucceed(t *testing.T) {
	if testing.Short() {
		t.Skip("Skip SNMP tests in short mode becase of lack of mocked SNMP server.")
	}
	snmp := newTestSnmp()
	actual, err := snmp.GetSnmpGlobalStats()
	if err != nil {
		t.Fatal(err)
	}

	if 12345678910 != actual.ServerBusyErrors {
		t.Fatalf("ServerBusyErrors is wrong: %d", actual.ServerBusyErrors)
	}
}
