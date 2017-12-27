package exporter

import (
	m "github.com/rerorero/netscaler-vpx-exporter/exporter/metric"
	"github.com/rerorero/netscaler-vpx-exporter/exporter/netscaler/nshttp"
	"github.com/rerorero/netscaler-vpx-exporter/exporter/netscaler/nssnmp"
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
			"http_enabled",
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
			"snmp_enabled",
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
			"http_busy_errors",
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
			"vserver_surge_count",
			"Number of requests in the surge queue.",
			vserverMetricsLabels,
			func(arg m.CollectArg) bool { return httpVsrvOf(&arg) != nil },
			func(arg m.GetGaugeArg) float64 { return httpVsrvOf(arg.CollectArg).VserverSurgeCount },
		),
		m.NewNsGauge(
			"vserver_health",
			"The percentage of UP services bound to this vserver.",
			vserverMetricsLabels,
			func(arg m.CollectArg) bool { return httpVsrvOf(&arg) != nil },
			func(arg m.GetGaugeArg) float64 { return httpVsrvOf(arg.CollectArg).Health },
		),
	}

	/*
			Name                  string  `json:"name"`
			VserverSrugeCount     float64 `json:"vsvrsurgecount,string"`
		EstablishedConn       float64 `json:"establishedconn,string"`
		InactiveServices      float64 `json:"inactsvcs,string"`
			Health                float64 `json:"vslbhealth,string"`
		PrimaryIpAddress      string  `json:"primaryipaddress"`
		PrimaryPort           int     `json:"primaryport"`
		Type                  string  `json:"type"`
		State                 string  `json:"state"`
		ActiveServices        float64 `json:"actsvcs,string"`
		TotalHits             float64 `json:"tothits,string"`
		HitsRate              float64 `json:"hitsrate"`
		TotalRequests         float64 `json:"totalrequests,string"`
		RequestsRate          float64 `json:"requestsrate"`
		TotalResponses        float64 `json:"totalresponses,string"`
		ResponsesRate         float64 `json:"responsesrate"`
		TotalRequestBytes     float64 `json:"totalrequestbytes,string"`
		RequestBytesRate      float64 `json:"requestbytesrate"`
		TotalResponseBytes    float64 `json:"totalresponsebytes,string"`
		ResponseBytesRate     float64 `json:"responsebytesrate"`
		TotalPackateReceived  float64 `json:"totalpktsrecvd,string"`
		PackageReceivedRate   float64 `json:"pktsrecvdrate"`
		TotalPackageSent      float64 `json:"totalpktssent,string"`
		PackageSentRate       float64 `json:"pktssentrate"`
		SurgeCount            float64 `json:"surgecount,string"`
		ServiceSurgeCount     float64 `json:"svcsurgecount,string"`
		InvlidRequestResponse float64 `json:"invalidrequestresponse,string"`
	*/
)

func httpVsrvOf(arg *m.CollectArg) *nshttp.HttpVServerStats {
	return arg.Stats.Http.VServers[arg.Labels[LabelVServer]]
}

func snmpVsrvOf(arg *m.CollectArg) *nssnmp.SnmpVserverStats {
	return arg.Stats.Snmp.VServers[arg.Labels[LabelVServer]]
}
