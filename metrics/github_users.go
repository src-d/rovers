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
	GitHubUsersTotalDur = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Subsystem: subsystemGitHubUsers,
			Name:      "total_duration_microseconds",
			Help:      "The time spent by the command to process all GitHub users in microseconds.",
		},
	)
)

func init() {
	prometheus.MustRegister(GitHubUsersProcessed)
	prometheus.MustRegister(GitHubUsersFailed)
	prometheus.MustRegister(GitHubUsersTotalDur)
}
