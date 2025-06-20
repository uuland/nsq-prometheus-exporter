package main

import (
	"flag"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"

	"caozhipan/nsq-prometheus-exporter/controllers"
)

var (
	nsqLookupdAddress = flag.String("nsq.lookupd.address", "127.0.0.1:4161", "nsqllookupd address list with comma")
	clientMetrics     = flag.Bool("client.metrics", false, "export client metrics, default false")
	nodeMetrics       = flag.Bool("node.metrics", false, "export node metrics, default false")
)

func main() {
	flag.Parse()

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		for {
			controllers.SyncNodeList(*nsqLookupdAddress)
			<-ticker.C
		}
	}()

	controllers.Collector.ScrapeClient = *clientMetrics
	controllers.Collector.ScrapeNode = *nodeMetrics

	prometheus.MustRegister(controllers.Collector)

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":9527", nil))

}
