package exporter

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/coord-e/mirakurun_exporter/mirakurun"
)

const namespace = "mirakurun"

type Config struct {
	FetchStatus   bool
	FetchTuners   bool
	FetchPrograms bool
	FetchServices bool
}

type Exporter struct {
	ctx    context.Context
	logger log.Logger

	status *statusExporter
	tuners *tunersExporter
	// programs *programsExporter
	// services *servicesExporter
}

// Verify if Exporter implements prometheus.Collector
var _ prometheus.Collector = (*Exporter)(nil)

func New(ctx context.Context, client *mirakurun.Client, config Config, logger log.Logger) *Exporter {
	var statusExporter *statusExporter
	if config.FetchStatus {
		statusExporter = newStatusExporter(ctx, client, logger)
	}

	var tunersExporter *tunersExporter
	if config.FetchTuners {
		tunersExporter = newTunersExporter(ctx, client, logger)
	}

	// var programsExporter *programsExporter
	// if config.FetchPrograms {
	// 	programsExporter := newProgramsExporter(ctx, client, logger)
	// }

	// var servicesExporter *servicesExporter
	// if config.FetchServices {
	// 	servicesExporter := newServicesExporter(ctx, client, logger)
	// }

	return &Exporter{
		ctx:    ctx,
		logger: logger,
		status: statusExporter,
		tuners: tunersExporter,
		// programs: programsExporter,
		// services: servicesExporter,
	}
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	if e.status != nil {
		e.status.Describe(ch)
	}
	if e.tuners != nil {
		e.tuners.Describe(ch)
	}
	// if e.programs != nil {
	// 	e.programs.Describe(ch)
	// }
	// if e.services != nil {
	// 	e.services.Describe(ch)
	// }
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	if e.status != nil {
		e.status.Collect(ch)
	}
	if e.tuners != nil {
		e.tuners.Collect(ch)
	}
	// if e.programs != nil {
	// 	e.programs.Collect(ch)
	// }
	// if e.services != nil {
	// 	e.services.Collect(ch)
	// }
}
