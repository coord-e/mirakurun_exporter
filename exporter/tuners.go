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

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/coord-e/mirakurun_exporter/mirakurun"
)

type tunersExporter struct {
	ctx    context.Context
	client *mirakurun.Client
	logger log.Logger

	availableTunerDevices *prometheus.Desc
	faultTunerDevices     *prometheus.Desc
	remoteTunerDevices    *prometheus.Desc
	grTunerDevices        *prometheus.Desc
	bsTunerDevices        *prometheus.Desc
	csTunerDevices        *prometheus.Desc
	skyTunerDevices       *prometheus.Desc
	tunerDevices          *prometheus.Desc
	users                 *prometheus.Desc
	streamDrops           *prometheus.Desc
	streamPackets         *prometheus.Desc
}

// Verify if tunersExporter implements prometheus.Collector
var _ prometheus.Collector = (*tunersExporter)(nil)

func newTunersExporter(ctx context.Context, client *mirakurun.Client, logger log.Logger) *tunersExporter {
	const subsystem = "tuners"

	return &tunersExporter{
		ctx:    ctx,
		client: client,
		logger: logger,

		availableTunerDevices: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "available_tuner_devices"),
			"Number of available tuner devices in Mirakurun.",
			[]string{"state"}, nil),
		faultTunerDevices: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "fault_tuner_devices"),
			"Number of fault tuner devices in Mirakurun.",
			nil, nil),
		remoteTunerDevices: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "remote_tuner_devices"),
			"Number of remote tuner devices in Mirakurun.",
			nil, nil),
		grTunerDevices: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "GR_tuner_devices"),
			"Number of GR tuner devices in Mirakurun.",
			nil, nil),
		bsTunerDevices: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "BS_tuner_devices"),
			"Number of BS tuner devices in Mirakurun.",
			nil, nil),
		csTunerDevices: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "CS_tuner_devices"),
			"Number of CS tuner devices in Mirakurun.",
			nil, nil),
		skyTunerDevices: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "SKY_tuner_devices"),
			"Number of SKY tuner devices in Mirakurun.",
			nil, nil),
		tunerDevices: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "tuner_devices"),
			"Number of all tuner devices in Mirakurun.",
			nil, nil),
		users: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "users"),
			"Number of tuner users in Mirakurun labeled by tuner device name.",
			[]string{"tuner_device"}, nil),
		streamDrops: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "stream_drops_total"),
			"Total number of drops in a TS stream of Mirakurun labeled by tuner device name.",
			[]string{"tuner_device"}, nil),
		streamPackets: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "stream_packets_total"),
			"Total number of packets in a TS stream of Mirakurun labeled by tuner device name.",
			[]string{"tuner_device"}, nil),
	}
}

func (e *tunersExporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.availableTunerDevices
	ch <- e.faultTunerDevices
	ch <- e.remoteTunerDevices
	ch <- e.grTunerDevices
	ch <- e.bsTunerDevices
	ch <- e.csTunerDevices
	ch <- e.skyTunerDevices
	ch <- e.tunerDevices
	ch <- e.users
	ch <- e.streamDrops
	ch <- e.streamPackets
}

func (e *tunersExporter) Collect(ch chan<- prometheus.Metric) {
	tuners, err := e.client.GetTuners(e.ctx)
	if err != nil {
		level.Error(e.logger).Log("msg", "failed to fetch Mirakurun tuners", "err", err)
		return
	}

	var availableFree, availableUsed, fault, remote, gr, bs, cs, sky int
	users := map[string]int{}
	drops := map[string]int64{}
	packets := map[string]int64{}
	for _, tuner := range *tuners {
		if tuner.IsFree {
			availableFree++
		}
		if tuner.IsUsing {
			availableUsed++
		}
		if tuner.IsFault {
			fault++
		}
		if tuner.IsRemote {
			remote++
		}
		for _, ty := range tuner.Types {
			switch ty {
			case "GR":
				gr++
			case "BS":
				bs++
			case "CS":
				cs++
			case "SKY":
				sky++
			default:
				level.Warn(e.logger).Log("msg", "unknown channel type", "type", ty)
			}
		}
		users[tuner.Name] = 0
		drops[tuner.Name] = 0
		packets[tuner.Name] = 0
		for _, user := range tuner.Users {
			users[tuner.Name]++

			if user.StreamInfo == nil {
				continue
			}
			for _, info := range *user.StreamInfo {
				drops[tuner.Name] += info.Drop
				packets[tuner.Name] += info.Packet
			}
		}
	}

	ch <- prometheus.MustNewConstMetric(e.availableTunerDevices, prometheus.GaugeValue, float64(availableFree), "free")
	ch <- prometheus.MustNewConstMetric(e.availableTunerDevices, prometheus.GaugeValue, float64(availableUsed), "used")
	ch <- prometheus.MustNewConstMetric(e.faultTunerDevices, prometheus.GaugeValue, float64(fault))
	ch <- prometheus.MustNewConstMetric(e.remoteTunerDevices, prometheus.GaugeValue, float64(remote))
	ch <- prometheus.MustNewConstMetric(e.grTunerDevices, prometheus.GaugeValue, float64(gr))
	ch <- prometheus.MustNewConstMetric(e.bsTunerDevices, prometheus.GaugeValue, float64(bs))
	ch <- prometheus.MustNewConstMetric(e.csTunerDevices, prometheus.GaugeValue, float64(cs))
	ch <- prometheus.MustNewConstMetric(e.skyTunerDevices, prometheus.GaugeValue, float64(sky))
	ch <- prometheus.MustNewConstMetric(e.tunerDevices, prometheus.GaugeValue, float64(len(*tuners)))
	for tunerDevice, count := range users {
		ch <- prometheus.MustNewConstMetric(e.users, prometheus.GaugeValue, float64(count), tunerDevice)
	}
	for tunerDevice, count := range drops {
		ch <- prometheus.MustNewConstMetric(e.streamDrops, prometheus.CounterValue, float64(count), tunerDevice)
	}
	for tunerDevice, count := range packets {
		ch <- prometheus.MustNewConstMetric(e.streamPackets, prometheus.CounterValue, float64(count), tunerDevice)
	}
}
