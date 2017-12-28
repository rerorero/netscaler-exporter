package metric

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rerorero/netscaler-exporter/exporter/netscaler"
)

type NsGauge interface {
	NsMetric
}

type GetGaugeArg struct {
	*CollectArg
	Gauge prometheus.Gauge
}

type nsGauge struct {
	*nsMetricBase
	metric *prometheus.GaugeVec
	get    func(GetGaugeArg) float64
}

func NewNsGauge(
	name string,
	help string,
	labels []string,
	exists func(CollectArg) bool,
	get func(GetGaugeArg) float64,
) NsGauge {
	vec := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: name,
			Help: help,
		},
		labels,
	)
	return &nsGauge{
		metric:       vec,
		nsMetricBase: &nsMetricBase{exists},
		get:          get,
	}
}

func (nsg *nsGauge) GetCollector() prometheus.Collector {
	return nsg.metric
}

func (nsg *nsGauge) Update(stats *netscaler.NetscalerStats, labels prometheus.Labels) bool {
	collected := false
	if stats != nil && nsg.exists(stats, labels) {
		gauge := nsg.metric.With(labels)
		v := nsg.get(GetGaugeArg{&CollectArg{stats, labels}, gauge})
		gauge.Set(v)
		collected = true
	}
	return collected
}

func (nsg *nsGauge) Reset() {
	nsg.metric.Reset()
}
