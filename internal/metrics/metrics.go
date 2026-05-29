package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var CommitsTotal = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "gitstats_commits",
		Help: "Total number of commits received, labeled by repo, author, prefix, and conventional compliance.",
	},
	[]string{"repo", "author", "prefix", "conventional"},
)
