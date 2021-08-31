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

	status   *statusExporter
	tuners   *tunersExporter
	programs *programsExporter
	services *servicesExporter
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

	var programsExporter *programsExporter
	if config.FetchPrograms {
		programsExporter = newProgramsExporter(ctx, client, logger)
	}

	var servicesExporter *servicesExporter
	if config.FetchServices {
		servicesExporter = newServicesExporter(ctx, client, logger)
	}

	return &Exporter{
		ctx:      ctx,
		logger:   logger,
		status:   statusExporter,
		tuners:   tunersExporter,
		programs: programsExporter,
		services: servicesExporter,
	}
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	if e.status != nil {
		e.status.Describe(ch)
	}
	if e.tuners != nil {
		e.tuners.Describe(ch)
	}
	if e.programs != nil {
		e.programs.Describe(ch)
	}
	if e.services != nil {
		e.services.Describe(ch)
	}
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	if e.status != nil {
		e.status.Collect(ch)
	}
	if e.tuners != nil {
		e.tuners.Collect(ch)
	}
	if e.programs != nil {
		e.programs.Collect(ch)
	}
	if e.services != nil {
		e.services.Collect(ch)
	}
}
