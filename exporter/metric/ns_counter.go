package metric

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rerorero/netscaler-vpx-exporter/exporter/netscaler"
)

type NsCounter interface {
	NsMetric
}

type GetCounterArg struct {
	*CollectArg
	Counter prometheus.Counter
}

type nsCounter struct {
	*nsMetricBase
	metric *prometheus.CounterVec
	get    func(GetCounterArg) float64
}

func NewNsCounter(
	name string,
	help string,
	labels []string,
	exists func(CollectArg) bool,
	get func(GetCounterArg) float64,
) NsCounter {
	vec := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: Namespace,
			Name:      name,
			Help:      help,
		},
		labels,
	)
	return &nsCounter{
		metric:       vec,
		nsMetricBase: &nsMetricBase{exists},
		get:          get,
	}
}

func (metric *nsCounter) GetCollector() prometheus.Collector {
	return metric.metric
}

func (nsc *nsCounter) Update(stats *netscaler.NetscalerStats, labels prometheus.Labels) bool {
	collected := false
	nsc.metric.Reset()
	if stats != nil && nsc.exists(stats, labels) {
		counter := nsc.metric.With(labels)
		v := nsc.get(GetCounterArg{&CollectArg{stats, labels}, counter})
		// Add() panics if value is 0.
		if v > 0 {
			counter.Add(v)
			collected = true
		}
	}
	return collected
}
