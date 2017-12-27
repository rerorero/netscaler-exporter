package nssnmp

import "testing"

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
