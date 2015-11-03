package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	subsystemAugur = subsystem + "_augur"

	AugurProcessed = prometheus.NewCounter(
		prometheus.CounterOpts{
			Subsystem: subsystemAugur,
			Name:      "processed",
			Help:      "Number of Augur Insights processed.",
		},
	)
	AugurFailed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Subsystem: subsystemAugur,
			Name:      "failed",
			Help:      "Number of Augur Insights failed.",
		},
		[]string{"reason"},
	)
	AugurRequested = prometheus.NewCounter(
		prometheus.CounterOpts{
			Subsystem: subsystemAugur,
			Name:      "requested",
			Help:      "number of requests made to Augur API.",
		},
	)
	AugurRequestDur = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Subsystem: subsystemAugur,
			Name:      "request_duration_microseconds",
			Help:      "The Augur API request latency in microseconds.",
		},
	)
)

func init() {
	prometheus.MustRegister(AugurProcessed)
	prometheus.MustRegister(AugurFailed)
	prometheus.MustRegister(AugurRequested)
	prometheus.MustRegister(AugurRequestDur)
}
