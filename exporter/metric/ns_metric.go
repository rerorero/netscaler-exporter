package metric

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rerorero/netscaler-vpx-exporter/exporter/netscaler"
)

type NsMetric interface {
	// returns prometheus collector instance.
	GetCollector() prometheus.Collector

	// Update the metric by Netscaler stats.
	// It returns where the metric collected or not.
	Update(*netscaler.NetscalerStats, prometheus.Labels) bool

	// Clear persistent value.
	// Since the collector sets a fixed value, it calls Reset() in Collect() every time.
	Reset()
}

type CollectArg struct {
	Stats  *netscaler.NetscalerStats
	Labels prometheus.Labels
}

type nsMetricBase struct {
	existsf func(CollectArg) bool
}

func (metric *nsMetricBase) exists(stats *netscaler.NetscalerStats, labels prometheus.Labels) bool {
	return metric.existsf(CollectArg{stats, labels})
}
