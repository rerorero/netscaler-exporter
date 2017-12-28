package exporter

import (
	"log"
	"sync"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/pkg/errors"
	"github.com/rerorero/netscaler-exporter/exporter/conf"
	"github.com/rerorero/netscaler-exporter/exporter/netscaler"
)

type Exporter interface {
	prometheus.Collector
}

type exporterImpl struct {
	conf       *conf.Conf
	netscalers []netscaler.Netscaler
}

func NewExporter(config *conf.Conf) (Exporter, error) {
	nsary := []netscaler.Netscaler{}
	for _, nsconf := range config.Netscaler.StaticTargets {
		ns, err := netscaler.NewNetscalerClient(nsconf)
		if err != nil {
			return nil, errors.Wrap(err, "error : Failed to instantiate Netscaler client")
		}
		nsary = append(nsary, ns)
	}

	return &exporterImpl{
		conf:       config,
		netscalers: nsary,
	}, nil
}

func (e *exporterImpl) Describe(ch chan<- *prometheus.Desc) {
	// global
	for _, metric := range globalMetrics {
		metric.GetCollector().Describe(ch)
	}

	// vserver
	for _, metric := range vserverMetrics {
		metric.GetCollector().Describe(ch)
	}
}

func (e *exporterImpl) Collect(ch chan<- prometheus.Metric) {
	wg := &sync.WaitGroup{}
	for _, ns := range e.netscalers {
		wg.Add(1)
		go doCollect(ns, ch, wg)
	}
	wg.Wait()
}

func doCollect(ns netscaler.Netscaler, ch chan<- prometheus.Metric, wg *sync.WaitGroup) {
	stats, errors := ns.GetStats()
	for _, err := range errors {
		log.Println("warn : Failed to get stats from ", ns.GetHost(), err.Error())
	}

	// global metrics
	for _, metric := range globalMetrics {
		labels := prometheus.Labels{LabelNsHost: ns.GetHost()}
		metric.Reset()
		metric.Update(stats, labels)
		metric.GetCollector().Collect(ch)
	}

	// vserver metrics
	vservers := collectVservers(stats)
	for _, metric := range vserverMetrics {
		metric.Reset()
		if stats != nil {
			for _, vserver := range vservers {
				labels := prometheus.Labels{
					LabelNsHost:  ns.GetHost(),
					LabelVServer: vserver,
				}
				metric.Update(stats, labels)
			}
		} else {
			metric.Update(nil, nil)
		}
		metric.GetCollector().Collect(ch)
	}

	wg.Done()
}

func collectVservers(stats *netscaler.NetscalerStats) []string {
	vservers := []string{}
	vserversSet := map[string]struct{}{}
	if stats != nil {
		for s, _ := range stats.Http.VServers {
			vserversSet[s] = struct{}{}
		}
		for s, _ := range stats.Snmp.VServers {
			vserversSet[s] = struct{}{}
		}
	}
	for s, _ := range vserversSet {
		vservers = append(vservers, s)
	}
	return vservers
}
