package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var CommitsTotal = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "gitstats_commits_total",
		Help: "Total number of commits received, labeled by repo, author, type, and conventional compliance.",
	},
	[]string{"repo", "author", "commit_type", "conventional"},
)
