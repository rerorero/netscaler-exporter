package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rerorero/netscaler-exporter/exporter"
	"github.com/rerorero/netscaler-exporter/exporter/conf"
)

var (
	confPath = flag.String("conf.file", "/etc/vpx-exporter/config.yaml", "Path to configuration file.")
)

func main() {
	flag.Parse()

	conf, err := conf.NewConfFromFile(*confPath)
	if err != nil {
		log.Fatal(err.Error())
	}

	exporter, err := exporter.NewExporter(conf)
	if err != nil {
		log.Fatal(err.Error())
	}

	prometheus.MustRegister(exporter)

	http.Handle("/metrics", promhttp.Handler())

	port := fmt.Sprintf(":%d", conf.BindPort)
	err = http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal(err.Error())
	}
}
