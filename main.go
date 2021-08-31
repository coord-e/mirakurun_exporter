package main

import (
	"net/http"
	"os"

	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/exporter-toolkit/web"
	webflag "github.com/prometheus/exporter-toolkit/web/kingpinflag"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/coord-e/mirakurun_exporter/exporter"
	"github.com/coord-e/mirakurun_exporter/mirakurun"
)

var (
	// see Makefile
	BuildVersion   = "devel"
	BuildCommitSha = "unknown"
)

var (
	webConfig     = webflag.AddFlags(kingpin.CommandLine)
	listenAddress = kingpin.Flag("web.listen-address", "The address to listen on for HTTP requests.").Default(":9110").String()
	metricPath    = kingpin.Flag("web.telemetry-path",
		"Path under which to expose metrics.").Default("/metrics").String()
	mirakurunPath = kingpin.Flag("exporter.mirakurun-path",
		"Path to the Mirakurun instance.").Required().String()
	fetchStatus = kingpin.Flag("exporter.status",
		"Whether to export metrics from /api/status.").Default("true").Bool()
	fetchTuners = kingpin.Flag("exporter.tuners",
		"Whether to export metrics from /api/tuners.").Default("true").Bool()
	fetchPrograms = kingpin.Flag("exporter.programs",
		"Whether to export metrics from /api/programs.").Default("true").Bool()
	fetchServices = kingpin.Flag("exporter.services",
		"Whether to export metrics from /api/services.").Default("true").Bool()
)

func main() {
	promlogConfig := &promlog.Config{}
	flag.AddFlags(kingpin.CommandLine, promlogConfig)
	kingpin.Version(BuildVersion)
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	logger := promlog.New(promlogConfig)

	level.Info(logger).Log("msg", "Starting mirakurun_exporter", "version", BuildVersion, "commit", BuildCommitSha)

	client, err := mirakurun.NewClient(*mirakurunPath)
	if err != nil {
		level.Error(logger).Log("msg", "failed to create Mirakurun client", "err", err)
		os.Exit(1)
	}

	config := exporter.Config{
		FetchStatus:   *fetchStatus,
		FetchTuners:   *fetchTuners,
		FetchPrograms: *fetchPrograms,
		FetchServices: *fetchServices,
	}
	var handler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		registry := prometheus.NewRegistry()
		exporter := exporter.New(r.Context(), client, config, logger)
		registry.MustRegister(exporter)
		h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
		h.ServeHTTP(w, r)
	}
	http.Handle(*metricPath, promhttp.InstrumentMetricHandler(prometheus.DefaultRegisterer, handler))

	level.Info(logger).Log("msg", "Listening on", "address", *listenAddress)
	server := &http.Server{Addr: *listenAddress}
	if err := web.ListenAndServe(server, *webConfig, logger); err != nil {
		level.Error(logger).Log("err", err)
		os.Exit(1)
	}
}
