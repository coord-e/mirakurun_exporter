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
	"strconv"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/coord-e/mirakurun_exporter/mirakurun"
)

type programsExporter struct {
	ctx    context.Context
	client *mirakurun.Client
	logger log.Logger

	programs *prometheus.Desc
}

// Verify if programsExporter implements prometheus.Collector
var _ prometheus.Collector = (*programsExporter)(nil)

func newProgramsExporter(ctx context.Context, client *mirakurun.Client, logger log.Logger) *programsExporter {
	const subsystem = "programs"

	return &programsExporter{
		ctx:    ctx,
		client: client,
		logger: logger,

		programs: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "stored_programs"),
			"Number of programs stored in Mirakurun.",
			[]string{"service_id"}, nil),
	}
}

func (e *programsExporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.programs
}

func (e *programsExporter) Collect(ch chan<- prometheus.Metric) {
	programs, err := e.client.GetPrograms(e.ctx)
	if err != nil {
		level.Error(e.logger).Log("msg", "failed to fetch Mirakurun programs", "err", err)
		return
	}

	counts := map[int]int{}
	for _, program := range *programs {
		counts[program.ServiceID]++
	}

	for serviceID, count := range counts {
		ch <- prometheus.MustNewConstMetric(e.programs, prometheus.GaugeValue, float64(count), strconv.Itoa(serviceID))
	}
}
