package main

import (
	"flag"
	"fmt"

	"github.com/rerorero/netscaler-vpx-exporter/exporter"
	"github.com/rerorero/netscaler-vpx-exporter/exporter/conf"
)

var (
	confPath = flag.String("conf.file", "/etc/vpx-exporter/config.yaml", "Path to configuration file.")
)

func main() {
	flag.Parse()

	conf, err := conf.NewConfFrom(*confPath)
	if err != nil {
		panic(err)
	}

	context, _ := exporter.NewContext(conf)
	fmt.Println(context)
	context.Netscalers[0].Authorize()
}
