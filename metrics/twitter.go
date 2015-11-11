package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	subsystemTwitter = subsystem + "_twitter"

	TwitterProcessed = prometheus.NewGauge(
		prometheus.CounterOpts{
			Subsystem: subsystemTwitter,
			Name:      "processed",
			Help:      "Number of Twitter profiles processed.",
		},
	)
	TwitterFailed = prometheus.NewGaugeVec(
		prometheus.CounterOpts{
			Subsystem: subsystemTwitter,
			Name:      "failed",
			Help:      "Number of Twitter profiles failed.",
		},
		[]string{"reason"},
	)
	TwitterRequested = prometheus.NewGauge(
		prometheus.CounterOpts{
			Subsystem: subsystemTwitter,
			Name:      "requested",
			Help:      "Number of requests made to Twitter web.",
		},
	)
	TwitterRequestDur = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Subsystem: subsystemTwitter,
			Name:      "request_duration_microseconds",
			Help:      "The Twitter web request latency in microseconds.",
		},
	)
)

func init() {
	prometheus.MustRegister(TwitterProcessed)
	prometheus.MustRegister(TwitterFailed)
	prometheus.MustRegister(TwitterRequested)
	prometheus.MustRegister(TwitterRequestDur)
}
