package sweeper

import (
	"context"
	"log"
	"time"
)

// Repo is the storage contract the sweeper needs.
// Concrete impl lives in internal/storage/event_repo.go.
type Repo interface {
	SweepOlderThan(cutoff time.Time) (int64, error)
}

// Run executes an initial sweep and then ticks every interval until ctx is cancelled.
// Per Constitution Principle IV, this replaces a cron daemon.
func Run(ctx context.Context, repo Repo, interval time.Duration, retention time.Duration) {
	if interval <= 0 {
		interval = 24 * time.Hour
	}
	if retention <= 0 {
		retention = 90 * 24 * time.Hour
	}
	sweepOnce(repo, retention)
	t := time.NewTicker(interval)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			sweepOnce(repo, retention)
		}
	}
}

func sweepOnce(repo Repo, retention time.Duration) {
	cutoff := time.Now().UTC().Add(-retention)
	n, err := repo.SweepOlderThan(cutoff)
	if err != nil {
		log.Printf("sweeper: error: %v", err)
		return
	}
	if n > 0 {
		log.Printf("sweeper: removed %d events older than %s", n, cutoff.Format(time.RFC3339))
	}
}