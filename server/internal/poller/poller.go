package poller

import (
	"context"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/karldreher/gitstats/server/internal/github"
	"github.com/karldreher/gitstats/server/internal/metrics"
	"github.com/karldreher/gitstats/server/internal/persistence"
)

var pollReady atomic.Bool

var conventionalRE = regexp.MustCompile(
	`^(feat|fix|docs|style|refactor|perf|test|chore)(\(.+\))?: .{1,50}`,
)
var commitTypeRE = regexp.MustCompile(
	`^(feat|fix|docs|style|refactor|perf|test|chore)`,
)

// IsReady reports whether at least one poll has completed successfully.
func IsReady() bool { return pollReady.Load() }

// Run starts the polling loop. It runs an immediate first poll, then ticks
// every POLL_INTERVAL_MINUTES (default 15). Blocks until ctx is cancelled.
func Run(ctx context.Context, client *github.Client, store persistence.StateStore) {
	interval := pollInterval()
	doPoll(client, store)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			doPoll(client, store)
		case <-ctx.Done():
			return
		}
	}
}

func doPoll(client *github.Client, store persistence.StateStore) {
	since := resolveSince(client.Mode(), store)
	log.Printf("polling since %s", since.Format(time.RFC3339))

	commits, err := client.FetchCommits(since)
	if err != nil {
		log.Println("poll error:", err)
		return
	}
	log.Printf("fetched %d commits", len(commits))

	for _, c := range commits {
		commitType, conventional := parseConventional(c.Message)
		metrics.CommitsTotal.WithLabelValues(c.Repo, c.Author, commitType, conventional).Inc()
		if store != nil {
			_ = store.Increment(persistence.LabelKey(c.Repo, c.Author, commitType, conventional))
		}
	}

	if store != nil {
		_ = store.Set(persistence.KeyLastPolledAt, float64(time.Now().Unix()))
		if err := store.Flush(); err != nil {
			log.Println("persistence flush error:", err)
		}
	}

	pollReady.Store(true)
}

func resolveSince(mode github.Mode, store persistence.StateStore) time.Time {
	if store != nil {
		if state, err := store.Load(); err == nil {
			if ts, ok := state[persistence.KeyLastPolledAt]; ok && ts > 0 {
				return time.Unix(int64(ts), 0)
			}
		}
	}
	if mode == github.PersonalMode {
		return time.Now().AddDate(0, 0, -30)
	}
	return time.Now()
}

func parseConventional(message string) (commitType, conventional string) {
	subject := strings.SplitN(message, "\n", 2)[0]
	if conventionalRE.MatchString(subject) {
		if m := commitTypeRE.FindString(subject); m != "" {
			return m, "true"
		}
	}
	return "undefined", "false"
}

func pollInterval() time.Duration {
	if v := os.Getenv("POLL_INTERVAL_MINUTES"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return time.Duration(n) * time.Minute
		}
	}
	return 15 * time.Minute
}
