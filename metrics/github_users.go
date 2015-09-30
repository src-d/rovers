package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	subsystemGitHubUsers = subsystem + "_github_users"

	GitHubUsersProcessed = prometheus.NewCounter(
		prometheus.CounterOpts{
			Subsystem: subsystemGitHubUsers,
			Name:      "processed",
			Help:      "Number of GitHub users processed.",
		},
	)
	GitHubUsersFailed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Subsystem: subsystemGitHubUsers,
			Name:      "failed",
			Help:      "Number of GitHub users failed.",
		},
		[]string{"reason"},
	)
	GitHubUsersRequested = prometheus.NewCounter(
		prometheus.CounterOpts{
			Subsystem: subsystemGitHubUsers,
			Name:      "requested",
			Help:      "number of requests made to GitHub Users API.",
		},
	)
	GitHubUsersRequestDur = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Subsystem: subsystemGitHubUsers,
			Name:      "request_duration_microseconds",
			Help:      "The GitHub API request latency in microseconds.",
		},
	)
)

func init() {
	prometheus.MustRegister(GitHubUsersProcessed)
	prometheus.MustRegister(GitHubUsersFailed)
	prometheus.MustRegister(GitHubUsersRequested)
	prometheus.MustRegister(GitHubUsersRequestDur)
}
