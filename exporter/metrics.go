package exporter

import (
	m "github.com/rerorero/netscaler-exporter/exporter/metric"
	"github.com/rerorero/netscaler-exporter/exporter/netscaler/nshttp"
	"github.com/rerorero/netscaler-exporter/exporter/netscaler/nssnmp"
)

const (
	LabelNsHost  = "ns_host"
	LabelVServer = "vserver"
)

var (
	// global metrics
	globalMetricsLabels = []string{LabelNsHost}
	globalMetrics       = []m.NsMetric{
		m.NewNsGauge(
			"ns_http_enabled",
			"Current state of the Netscaler Rest API port.",
			globalMetricsLabels,
			func(arg m.CollectArg) bool { return true },
			func(arg m.GetGaugeArg) float64 {
				if arg.Stats.Http.HttpHealth {
					return 1
				} else {
					return 0
				}
			},
		),
		m.NewNsGauge(
			"ns_snmp_enabled",
			"Current state of the Netscaler SNMP port.",
			globalMetricsLabels,
			func(arg m.CollectArg) bool { return true },
			func(arg m.GetGaugeArg) float64 {
				if arg.Stats.Snmp.SnmpHealth {
					return 1
				} else {
					return 0
				}
			},
		),
		m.NewNsCounter(
			"ns_http_busy_errors",
			"Total number of HTTP error responses received.",
			globalMetricsLabels,
			func(arg m.CollectArg) bool { return arg.Stats.Snmp.Global != nil },
			func(arg m.GetCounterArg) float64 { return float64(arg.Stats.Snmp.Global.ServerBusyErrors) },
		),
	}

	// Vitrual server metrics
	vserverMetricsLabels = []string{LabelNsHost, LabelVServer}
	vserverMetrics       = []m.NsMetric{
		m.NewNsGauge(
			"ns_vserver_ttfb",
			"Average TTFB(Time to First Byte) between the NetScaler appliance and the server.",
			vserverMetricsLabels,
			func(arg m.CollectArg) bool { return snmpVsrvOf(&arg) != nil },
			func(arg m.GetGaugeArg) float64 { return float64(snmpVsrvOf(arg.CollectArg).LbVserverAverageTTFB) },
		),
		m.NewNsGauge(
			"ns_vserver_established_conn",
			"Number of client connections in ESTABLISHED state.",
			vserverMetricsLabels,
			func(arg m.CollectArg) bool { return httpVsrvOf(&arg) != nil },
			func(arg m.GetGaugeArg) float64 { return httpVsrvOf(arg.CollectArg).EstablishedConn },
		),
		m.NewNsGauge(
			"ns_vserver_inactive_services",
			"Number of inactive services.",
			vserverMetricsLabels,
			func(arg m.CollectArg) bool { return httpVsrvOf(&arg) != nil },
			func(arg m.GetGaugeArg) float64 { return httpVsrvOf(arg.CollectArg).InactiveServices },
		),
		m.NewNsGauge(
			"ns_vserver_surge_count",
			"Number of requests in the surge queue.",
			vserverMetricsLabels,
			func(arg m.CollectArg) bool { return httpVsrvOf(&arg) != nil },
			func(arg m.GetGaugeArg) float64 { return httpVsrvOf(arg.CollectArg).VserverSurgeCount },
		),
		m.NewNsGauge(
			"ns_vserver_health",
			"The percentage of UP services bound to this vserver.",
			vserverMetricsLabels,
			func(arg m.CollectArg) bool { return httpVsrvOf(&arg) != nil },
			func(arg m.GetGaugeArg) float64 { return httpVsrvOf(arg.CollectArg).Health },
		),
		m.NewNsGauge(
			"ns_vserver_type",
			"Protocol associated with the vserver.HTTP=1,TCP=2,SSL=3,UDP=4,OTHER=0",
			vserverMetricsLabels,
			func(arg m.CollectArg) bool { return httpVsrvOf(&arg) != nil },
			func(arg m.GetGaugeArg) float64 {
				switch httpVsrvOf(arg.CollectArg).Type {
				case "HTTP":
					return 1
				case "TCP":
					return 2
				case "SSL":
					return 3
				case "UDP":
					return 4
				default:
					return 0
				}
			},
		),
		m.NewNsGauge(
			"ns_vserver_state",
			"Current state of the server.1=UP, 0=DOWN or other state.",
			vserverMetricsLabels,
			func(arg m.CollectArg) bool { return httpVsrvOf(&arg) != nil },
			func(arg m.GetGaugeArg) float64 {
				if httpVsrvOf(arg.CollectArg).State == "UP" {
					return 1
				} else {
					return 0
				}
			},
		),
		m.NewNsGauge(
			"ns_vserver_active_services",
			"Number of active services.",
			vserverMetricsLabels,
			func(arg m.CollectArg) bool { return httpVsrvOf(&arg) != nil },
			func(arg m.GetGaugeArg) float64 { return httpVsrvOf(arg.CollectArg).ActiveServices },
		),
		m.NewNsCounter(
			"ns_vserver_total_hits",
			"Total vserver hits",
			vserverMetricsLabels,
			func(arg m.CollectArg) bool { return httpVsrvOf(&arg) != nil },
			func(arg m.GetCounterArg) float64 { return httpVsrvOf(arg.CollectArg).TotalHits },
		),
		m.NewNsGauge(
			"ns_vserver_hits_rate",
			"Rate of vserver hits.",
			vserverMetricsLabels,
			func(arg m.CollectArg) bool { return httpVsrvOf(&arg) != nil },
			func(arg m.GetGaugeArg) float64 { return httpVsrvOf(arg.CollectArg).HitsRate },
		),
		m.NewNsCounter(
			"ns_vserver_total_requests",
			"Total number of requests received on this vserver.This applies to HTTP/SSL.",
			vserverMetricsLabels,
			func(arg m.CollectArg) bool { return httpVsrvOf(&arg) != nil },
			func(arg m.GetCounterArg) float64 { return httpVsrvOf(arg.CollectArg).TotalRequests },
		),
		m.NewNsGauge(
			"ns_vserver_requests_rate",
			"Rate of requests received on this vserver.This applies to HTTP/SSL.",
			vserverMetricsLabels,
			func(arg m.CollectArg) bool { return httpVsrvOf(&arg) != nil },
			func(arg m.GetGaugeArg) float64 { return httpVsrvOf(arg.CollectArg).RequestsRate },
		),
		m.NewNsCounter(
			"ns_vserver_total_responses",
			"Total number of responses received on this vserver.This applies to HTTP/SSL.",
			vserverMetricsLabels,
			func(arg m.CollectArg) bool { return httpVsrvOf(&arg) != nil },
			func(arg m.GetCounterArg) float64 { return httpVsrvOf(arg.CollectArg).TotalResponses },
		),
		m.NewNsGauge(
			"ns_vserver_responses_rate",
			"Rate of responses received on this vserver.This applies to HTTP/SSL.",
			vserverMetricsLabels,
			func(arg m.CollectArg) bool { return httpVsrvOf(&arg) != nil },
			func(arg m.GetGaugeArg) float64 { return httpVsrvOf(arg.CollectArg).ResponsesRate },
		),
		m.NewNsCounter(
			"ns_vserver_total_req_bytes",
			"Total number of request bytes received on this vserver.",
			vserverMetricsLabels,
			func(arg m.CollectArg) bool { return httpVsrvOf(&arg) != nil },
			func(arg m.GetCounterArg) float64 { return httpVsrvOf(arg.CollectArg).TotalRequestBytes },
		),
		m.NewNsGauge(
			"ns_vserver_req_bytes_rate",
			"Rate of requests received on this vserver.",
			vserverMetricsLabels,
			func(arg m.CollectArg) bool { return httpVsrvOf(&arg) != nil },
			func(arg m.GetGaugeArg) float64 { return httpVsrvOf(arg.CollectArg).RequestBytesRate },
		),
		m.NewNsCounter(
			"ns_vserver_total_resp_bytes",
			"Total number of response bytes received on this vserver.",
			vserverMetricsLabels,
			func(arg m.CollectArg) bool { return httpVsrvOf(&arg) != nil },
			func(arg m.GetCounterArg) float64 { return httpVsrvOf(arg.CollectArg).TotalResponseBytes },
		),
		m.NewNsGauge(
			"ns_vserver_resp_bytes_rate",
			"Rate of responses received on this vserver.",
			vserverMetricsLabels,
			func(arg m.CollectArg) bool { return httpVsrvOf(&arg) != nil },
			func(arg m.GetGaugeArg) float64 { return httpVsrvOf(arg.CollectArg).ResponseBytesRate },
		),
	}
)

func httpVsrvOf(arg *m.CollectArg) *nshttp.HttpVServerStats {
	return arg.Stats.Http.VServers[arg.Labels[LabelVServer]]
}

func snmpVsrvOf(arg *m.CollectArg) *nssnmp.SnmpVserverStats {
	return arg.Stats.Snmp.VServers[arg.Labels[LabelVServer]]
}
