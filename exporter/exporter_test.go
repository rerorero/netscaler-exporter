package exporter

import (
	"reflect"
	"sync"
	"testing"

	"github.com/rerorero/netscaler-exporter/exporter/netscaler/nssnmp"

	"github.com/rerorero/netscaler-exporter/exporter/netscaler/nshttp"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/rerorero/netscaler-exporter/exporter/metric"
	"github.com/rerorero/netscaler-exporter/exporter/netscaler"
)

// Mocked netscaler client
type mockedNetscaler struct {
	getStats func() (*netscaler.NetscalerStats, []error)
	getHost  string
}

func (ns *mockedNetscaler) GetStats() (*netscaler.NetscalerStats, []error) { return ns.getStats() }
func (ns *mockedNetscaler) GetHost() string                                { return ns.getHost }

// Mocked prometheus Collector
type mockedCollector struct {
	describe func(chan<- *prometheus.Desc)
	collect  func(chan<- prometheus.Metric)
}

func (c *mockedCollector) Describe(ch chan<- *prometheus.Desc) { c.describe(ch) }
func (c *mockedCollector) Collect(ch chan<- prometheus.Metric) { c.collect(ch) }

// Mocked prometheus Metric
type mockedPromMetric struct{ name string }

func (m *mockedPromMetric) Desc() *prometheus.Desc     { return nil }
func (m *mockedPromMetric) Write(mm *dto.Metric) error { return nil }

// Mocked NsMetric
type mockedNsMetric struct {
	collector prometheus.Collector
	update    func(*netscaler.NetscalerStats, prometheus.Labels) bool
	reset     func()
}

func (nsm *mockedNsMetric) GetCollector() prometheus.Collector { return nsm.collector }
func (nsm *mockedNsMetric) Update(stats *netscaler.NetscalerStats, labels prometheus.Labels) bool {
	return nsm.update(stats, labels)
}
func (nsm *mockedNsMetric) Reset() { nsm.reset() }

func TestExporterCollect(t *testing.T) {
	mockStatsMap := map[string]*netscaler.NetscalerStats{}

	mock1Stats := &netscaler.NetscalerStats{
		Http: nshttp.NetscalerHttpStats{
			VServers: map[string]*nshttp.HttpVServerStats{
				"mock1_1": nil,
				"mock1_2": nil,
			},
		},
		Snmp: nssnmp.NetscalerSnmpStats{
			VServers: map[string]*nssnmp.SnmpVserverStats{
				"mock1_1": nil,
				"mock1_2": nil,
			},
		},
	}
	mock1 := &mockedNetscaler{
		getStats: func() (*netscaler.NetscalerStats, []error) {
			return mock1Stats, nil
		},
		getHost: "mock1",
	}
	mockStatsMap[mock1.getHost] = mock1Stats

	mock2Stats := &netscaler.NetscalerStats{
		Http: nshttp.NetscalerHttpStats{
			VServers: map[string]*nshttp.HttpVServerStats{
				"mock2_1": nil,
				"mock2_2": nil,
				"mock2_3": nil,
				"mock2_4": nil,
			},
		},
	}
	mock2 := &mockedNetscaler{
		getStats: func() (*netscaler.NetscalerStats, []error) {
			return mock2Stats, nil
		},
		getHost: "mock2",
	}
	mockStatsMap[mock2.getHost] = mock2Stats

	// mocked metrics
	updateCalled := map[string]int{}
	updateCalledMutex := sync.Mutex{}
	updateCall := func(name string) {
		updateCalledMutex.Lock()
		count := updateCalled[name]
		updateCalled[name] = count + 1
		updateCalledMutex.Unlock()
	}

	updateGlobalMock := func(name string) func(*netscaler.NetscalerStats, prometheus.Labels) bool {
		return func(s *netscaler.NetscalerStats, l prometheus.Labels) bool {
			if len(l) != 1 {
				t.Fatal("invalid labels", l)
			}
			host := l[LabelNsHost]
			if !reflect.DeepEqual(mockStatsMap[host], s) {
				t.Fatal("Unexpected stats: ", mockStatsMap[host], s)
			}
			updateCall(name)
			return true
		}
	}
	testGlobalMetrics := []metric.NsMetric{
		&mockedNsMetric{
			collector: &mockedCollector{
				describe: func(ch chan<- *prometheus.Desc) { ch <- &prometheus.Desc{} },
				collect:  func(ch chan<- prometheus.Metric) { ch <- &mockedPromMetric{name: "global_metrics1"} },
			},
			update: updateGlobalMock("global_metrics1"),
			reset:  func() {},
		},
		&mockedNsMetric{
			collector: &mockedCollector{
				describe: func(ch chan<- *prometheus.Desc) { ch <- &prometheus.Desc{} },
				collect:  func(ch chan<- prometheus.Metric) { ch <- &mockedPromMetric{name: "global_metrics2"} },
			},
			update: updateGlobalMock("global_metrics2"),
			reset:  func() {},
		},
		&mockedNsMetric{
			collector: &mockedCollector{
				describe: func(ch chan<- *prometheus.Desc) { ch <- &prometheus.Desc{} },
				collect:  func(ch chan<- prometheus.Metric) { ch <- &mockedPromMetric{name: "global_metrics3"} },
			},
			update: updateGlobalMock("global_metrics3"),
			reset:  func() {},
		},
	}

	updateVserverMock := func(name string) func(s *netscaler.NetscalerStats, l prometheus.Labels) bool {
		return func(s *netscaler.NetscalerStats, l prometheus.Labels) bool {
			if len(l) != 2 {
				t.Fatal("invalid labels", l)
			}
			_, ok := l[LabelVServer]
			if !ok {
				t.Fatal("invalid labels", l)
			}
			host := l[LabelNsHost]
			if !reflect.DeepEqual(mockStatsMap[host], s) {
				t.Fatal("Unexpected stats: ", mockStatsMap[host], s)
			}
			updateCall(name)
			return true
		}
	}
	testVserverMetrics := []metric.NsMetric{
		&mockedNsMetric{
			collector: &mockedCollector{
				describe: func(ch chan<- *prometheus.Desc) { ch <- &prometheus.Desc{} },
				collect:  func(ch chan<- prometheus.Metric) { ch <- &mockedPromMetric{name: "vserver_metrics1"} },
			},
			update: updateVserverMock("vserver_metrics1"),
			reset:  func() {},
		},
		&mockedNsMetric{
			collector: &mockedCollector{
				describe: func(ch chan<- *prometheus.Desc) { ch <- &prometheus.Desc{} },
				collect:  func(ch chan<- prometheus.Metric) { ch <- &mockedPromMetric{name: "vserver_metrics2"} },
			},
			update: updateVserverMock("vserver_metrics2"),
			reset:  func() {},
		},
	}

	// create sut
	sut := &exporterImpl{
		netscalers:     []netscaler.Netscaler{mock1, mock2},
		globalMetrics:  testGlobalMetrics,
		vserverMetrics: testVserverMetrics,
	}

	// create mocked channel
	ch := make(chan prometheus.Metric)
	defer close(ch)

	// start collect
	go sut.Collect(ch)

	// wait
	actualCollectCount := map[string]int{}
	expectedMetricNum := (len(sut.globalMetrics) + len(sut.vserverMetrics)) * len(sut.netscalers)
	for i := 0; i < expectedMetricNum; i++ {
		m := <-ch
		current := actualCollectCount[m.(*mockedPromMetric).name]
		actualCollectCount[m.(*mockedPromMetric).name] = current + 1
	}

	// verify
	expectedCollect := map[string]int{
		"global_metrics1":  len(sut.netscalers),
		"global_metrics2":  len(sut.netscalers),
		"global_metrics3":  len(sut.netscalers),
		"vserver_metrics1": len(sut.netscalers),
		"vserver_metrics2": len(sut.netscalers),
	}
	if !reflect.DeepEqual(expectedCollect, actualCollectCount) {
		t.Fatal("Unexpected collected count", expectedCollect, actualCollectCount)
	}

	vserverMetricsNum := 0
	for _, stat := range mockStatsMap {
		vserverMetricsNum += len(stat.Http.VServers)
	}
	expectedUpdates := map[string]int{
		"global_metrics1":  len(sut.netscalers),
		"global_metrics2":  len(sut.netscalers),
		"global_metrics3":  len(sut.netscalers),
		"vserver_metrics1": vserverMetricsNum,
		"vserver_metrics2": vserverMetricsNum,
	}
	if !reflect.DeepEqual(expectedUpdates, updateCalled) {
		t.Fatal("Unexpected update count", expectedUpdates, updateCalled)
	}
}
