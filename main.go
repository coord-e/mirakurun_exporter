// Copyright 2021 coord_e
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  	 http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"net/http"
	"os"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/exporter-toolkit/web"
	webflag "github.com/prometheus/exporter-toolkit/web/kingpinflag"

	"github.com/coord-e/mirakurun_exporter/exporter"
	"github.com/coord-e/mirakurun_exporter/mirakurun"
)

var (
	// see Makefile
	BuildVersion   = "devel"
	BuildCommitSha = "unknown"
)

var (
	webConfig  = webflag.AddFlags(kingpin.CommandLine, ":9110")
	metricPath = kingpin.Flag("web.telemetry-path",
		"Path under which to expose metrics.").Default("/metrics").String()
	mirakurunURL = kingpin.Flag("exporter.mirakurun-url",
		"URL of the Mirakurun instance.").Required().String()
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

	client, err := mirakurun.NewClient(*mirakurunURL)
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

	server := &http.Server{
		ReadHeaderTimeout: 5 * time.Second,
	}
	if err := web.ListenAndServe(server, webConfig, logger); err != nil {
		level.Error(logger).Log("err", err)
		os.Exit(1)
	}
}
