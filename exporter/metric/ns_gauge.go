package metric

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rerorero/netscaler-vpx-exporter/exporter/netscaler"
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
			Namespace: Namespace,
			Name:      name,
			Help:      help,
		},
		labels,
	)
	return &nsGauge{
		metric:       vec,
		nsMetricBase: &nsMetricBase{exists},
		get:          get,
	}
}

func (metric *nsGauge) GetCollector() prometheus.Collector {
	return metric.metric
}

func (nsg *nsGauge) Update(stats *netscaler.NetscalerStats, labels prometheus.Labels) bool {
	collected := false
	nsg.metric.Reset()
	if stats != nil && nsg.exists(stats, labels) {
		gauge := nsg.metric.With(labels)
		v := nsg.get(GetGaugeArg{&CollectArg{stats, labels}, gauge})
		gauge.Set(v)
	}
	return collected
}
