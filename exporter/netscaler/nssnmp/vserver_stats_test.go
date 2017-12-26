package nssnmp

import "testing"
import "reflect"

func TestGetSnmpVserverStatsSucceed(t *testing.T) {
	if testing.Short() {
		t.Skip("Skip SNMP tests in short mode becase of lack of mocked SNMP server.")
	}
	snmp := newTestSnmp()
	actual, err := snmp.GetSnmpVserverStats()
	if err != nil {
		t.Fatal(err)
	}
	expected := map[string]*SnmpVserverStats{}
	expected["lbvserver1"] = &SnmpVserverStats{
		Name:                 "lbvserver1",
		LbVserverAverageTTFB: 100,
	}
	expected["lbvserver2"] = &SnmpVserverStats{
		Name:                 "lbvserver2",
		LbVserverAverageTTFB: 200,
	}
	expected["lbvserver3"] = &SnmpVserverStats{
		Name:                 "lbvserver3",
		LbVserverAverageTTFB: 300,
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("lbvservers snmp stats are not matched %v - %v", expected, actual)
	}
}
