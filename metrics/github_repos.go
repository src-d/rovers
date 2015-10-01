package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	subsystemGitHubRepos = subsystem + "_github_repos"

	GitHubReposProcessed = prometheus.NewCounter(
		prometheus.CounterOpts{
			Subsystem: subsystemGitHubRepos,
			Name:      "processed",
			Help:      "Number of GitHub repositories processed.",
		},
	)
	GitHubReposFailed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Subsystem: subsystemGitHubRepos,
			Name:      "failed",
			Help:      "Number of GitHub repositories failed.",
		},
		[]string{"reason"},
	)
	GitHubReposRequested = prometheus.NewCounter(
		prometheus.CounterOpts{
			Subsystem: subsystemGitHubRepos,
			Name:      "requested",
			Help:      "number of requests made to GitHub API.",
		},
	)
	GitHubReposRequestDur = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Subsystem: subsystemGitHubRepos,
			Name:      "request_duration_microseconds",
			Help:      "The GitHub API request latency in microseconds.",
		},
	)
)

func init() {
	prometheus.MustRegister(GitHubReposProcessed)
	prometheus.MustRegister(GitHubReposFailed)
	prometheus.MustRegister(GitHubReposRequested)
	prometheus.MustRegister(GitHubReposRequestDur)
}
