package exporter

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/coord-e/mirakurun_exporter/mirakurun"
)

type statusExporter struct {
	ctx    context.Context
	client *mirakurun.Client
	logger log.Logger

	residentMemory   *prometheus.Desc
	totalMemory      *prometheus.Desc
	usedMemory       *prometheus.Desc
	programsDBEvents *prometheus.Desc
	rpcConnections   *prometheus.Desc
	streams          *prometheus.Desc
	errors           *prometheus.Desc
	timerError1      *prometheus.Desc
	timerError5      *prometheus.Desc
	timerError15     *prometheus.Desc
	info             *prometheus.Desc
}

// Verify if statusExporter implements prometheus.Collector
var _ prometheus.Collector = (*statusExporter)(nil)

func newStatusExporter(ctx context.Context, client *mirakurun.Client, logger log.Logger) *statusExporter {
	const subsystem = "status"

	return &statusExporter{
		ctx:    ctx,
		client: client,
		logger: logger,

		residentMemory: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "resident_memory_bytes"),
			"Amount of space occupied in the main memory device for the Mirakurun process in bytes.",
			nil, nil),
		totalMemory: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "total_memory_bytes"),
			"Total heap size of the Mirakurun process in bytes.",
			nil,
			nil),
		usedMemory: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "used_memory_bytes"),
			"Used heap size of the Mirakurun process in bytes.",
			nil,
			nil),
		programsDBEvents: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "programs_db_events"),
			"Number of EPG programs stored in Mirakurun.",
			nil,
			nil),
		rpcConnections: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "rpc_connections"),
			"Number of JSON-RPC connections to Mirakurun.",
			nil,
			nil),
		streams: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "streams"),
			"Number of streams in Mirakurun.",
			[]string{"stream"},
			nil),
		errors: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "errors_total"),
			"Total number of errors in Mirakurun.",
			[]string{"error"},
			nil),
		timerError1: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "timer_error1_seconds"),
			"1m average difference from clock time in Mirakurun to the real time in seconds.",
			nil,
			nil),
		timerError5: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "timer_error5_seconds"),
			"5m average difference from clock time in Mirakurun to the real time in seconds.",
			nil,
			nil),
		timerError15: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "timer_error15_seconds"),
			"15m average difference from clock time in Mirakurun to the real time in seconds.",
			nil,
			nil),
		info: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "info"),
			"A metric with a constant '1' value labeled by metadata of Mirakurun.",
			[]string{"nodeversion", "version", "arch"},
			nil),
	}
}

func (e *statusExporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.residentMemory
	ch <- e.totalMemory
	ch <- e.usedMemory
	ch <- e.programsDBEvents
	ch <- e.rpcConnections
	ch <- e.streams
	ch <- e.errors
	ch <- e.timerError1
	ch <- e.timerError5
	ch <- e.timerError15
	ch <- e.info
}

func (e *statusExporter) Collect(ch chan<- prometheus.Metric) {
	status, err := e.client.GetStatus(e.ctx)
	if err != nil {
		level.Error(e.logger).Log("msg", "failed to fetch Mirakurun status", "err", err)
		return
	}

	ch <- prometheus.MustNewConstMetric(e.residentMemory, prometheus.GaugeValue, float64(status.Process.MemoryUsage.RSS))
	ch <- prometheus.MustNewConstMetric(e.totalMemory, prometheus.GaugeValue, float64(status.Process.MemoryUsage.HeapTotal))
	ch <- prometheus.MustNewConstMetric(e.usedMemory, prometheus.GaugeValue, float64(status.Process.MemoryUsage.HeapUsed))
	ch <- prometheus.MustNewConstMetric(e.programsDBEvents, prometheus.GaugeValue, float64(status.EPG.StoredEvents))
	if status.RPCCount != nil {
		ch <- prometheus.MustNewConstMetric(e.rpcConnections, prometheus.GaugeValue, float64(*status.RPCCount))
	}
	ch <- prometheus.MustNewConstMetric(e.streams, prometheus.GaugeValue, float64(status.StreamCount.TunerDevice), "tuner_device")
	ch <- prometheus.MustNewConstMetric(e.streams, prometheus.GaugeValue, float64(status.StreamCount.TSFilter), "ts_filter")
	ch <- prometheus.MustNewConstMetric(e.streams, prometheus.GaugeValue, float64(status.StreamCount.Decoder), "decoder")
	ch <- prometheus.MustNewConstMetric(e.errors, prometheus.CounterValue, float64(status.ErrorCount.UncaughtException), "uncaught_exception")
	ch <- prometheus.MustNewConstMetric(e.errors, prometheus.CounterValue, float64(status.ErrorCount.UnhandledRejection), "unhandled_rejection")
	ch <- prometheus.MustNewConstMetric(e.errors, prometheus.CounterValue, float64(status.ErrorCount.BufferOverflow), "buffer_overflow")
	ch <- prometheus.MustNewConstMetric(e.errors, prometheus.CounterValue, float64(status.ErrorCount.TunerDeviceRespawn), "tuner_device_respawn")
	ch <- prometheus.MustNewConstMetric(e.errors, prometheus.CounterValue, float64(status.ErrorCount.DecoderRespawn), "decoder_respawn")
	ch <- prometheus.MustNewConstMetric(e.timerError1, prometheus.GaugeValue, status.TimerAccuracy.M1.Avg/1000000)
	ch <- prometheus.MustNewConstMetric(e.timerError5, prometheus.GaugeValue, status.TimerAccuracy.M5.Avg/1000000)
	ch <- prometheus.MustNewConstMetric(e.timerError15, prometheus.GaugeValue, status.TimerAccuracy.M15.Avg/1000000)
	ch <- prometheus.MustNewConstMetric(e.info, prometheus.UntypedValue, 1.0, status.Process.Versions["node"], status.Version, status.Process.Arch)
}
