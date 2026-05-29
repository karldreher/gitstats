// Package metrics defines the Prometheus metrics exported by tacoma.
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// CommitsTotal counts commits processed, partitioned by repo, author, type, and conventional-commit flag.
var CommitsTotal = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "gitstats_commits",
		Help: "Total number of commits received, labeled by repo, author, prefix, and conventional compliance.",
	},
	[]string{"repo", "author", "prefix", "conventional"},
)
