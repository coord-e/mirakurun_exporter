package exporter

import (
	"context"
	"strconv"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/coord-e/mirakurun_exporter/mirakurun"
)

type servicesExporter struct {
	ctx    context.Context
	client *mirakurun.Client
	logger log.Logger

	grServices *prometheus.Desc
	services   *prometheus.Desc
}

// Verify if servicesExporter implements prometheus.Collector
var _ prometheus.Collector = (*servicesExporter)(nil)

func newServicesExporter(ctx context.Context, client *mirakurun.Client, logger log.Logger) *servicesExporter {
	const subsystem = "services"

	return &servicesExporter{
		ctx:    ctx,
		client: client,
		logger: logger,

		grServices: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "GR_services"),
			"Number of GR services available in Mirakurun.",
			[]string{"channel"}, nil),
		services: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "services"),
			"Number of all services available in Mirakurun.",
			[]string{"network_id"}, nil),
	}
}

func (e *servicesExporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.grServices
	ch <- e.services
}

func (e *servicesExporter) Collect(ch chan<- prometheus.Metric) {
	services, err := e.client.GetServices(e.ctx)
	if err != nil {
		level.Error(e.logger).Log("msg", "failed to fetch Mirakurun services", "err", err)
		return
	}

	grCounts := map[string]int{}
	counts := map[int]int{}
	for _, service := range *services {
		if service.Channel != nil && service.Channel.Type == "GR" {
			grCounts[service.Channel.Channel]++
		}
		counts[service.NetworkID]++
	}

	for channel, count := range grCounts {
		ch <- prometheus.MustNewConstMetric(e.grServices, prometheus.GaugeValue, float64(count), channel)
	}
	for networkID, count := range counts {
		ch <- prometheus.MustNewConstMetric(e.services, prometheus.GaugeValue, float64(count), strconv.Itoa(networkID))
	}
}
