package storage

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func openTestDB(t *testing.T) (*EventRepo, func()) {
	t.Helper()
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")
	db, err := Open(dbPath)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	return NewEventRepo(db), func() { _ = os.Remove(dbPath) }
}

func TestInsertAndTrendingAndSweep(t *testing.T) {
	repo, cleanup := openTestDB(t)
	defer cleanup()

	now := time.Now().UTC()
	// Two recent SEARCH events for ext1, one for ext2, one old event.
	must := func(e error) { t.Helper(); if e != nil { t.Fatal(e) } }
	must(repo.Insert(&BehaviorEvent{EventType: "SEARCH", TargetType: "EXTENSION", TargetIdentifier: "ext1", CountryCode: "CN", CreatedAt: now}))
	must(repo.Insert(&BehaviorEvent{EventType: "SEARCH", TargetType: "EXTENSION", TargetIdentifier: "ext1", CountryCode: "CN", CreatedAt: now.Add(-time.Minute)}))
	must(repo.Insert(&BehaviorEvent{EventType: "SEARCH", TargetType: "EXTENSION", TargetIdentifier: "ext2", CountryCode: "US", CreatedAt: now}))
	must(repo.Insert(&BehaviorEvent{EventType: "SEARCH", TargetType: "EXTENSION", TargetIdentifier: "ext-old", CountryCode: "CN", CreatedAt: now.AddDate(0, 0, -100)}))

	rows, err := repo.Trending("EXTENSION", "24h", 10)
	if err != nil {
		t.Fatalf("trending: %v", err)
	}
	if len(rows) != 2 {
		t.Fatalf("want 2 rows, got %d", len(rows))
	}
	if rows[0].TargetIdentifier != "ext1" || rows[0].Cnt != 2 {
		t.Errorf("want ext1/2, got %s/%d", rows[0].TargetIdentifier, rows[0].Cnt)
	}

	// Sweep cuts anything older than 90 days.
	n, err := repo.SweepOlderThan(now.AddDate(0, 0, -90))
	if err != nil {
		t.Fatalf("sweep: %v", err)
	}
	if n != 1 {
		t.Errorf("want 1 swept, got %d", n)
	}
}

func TestTrendingInvalidWindow(t *testing.T) {
	repo, cleanup := openTestDB(t)
	defer cleanup()
	if _, err := repo.Trending("EXTENSION", "1y", 10); err == nil {
		t.Error("want error for invalid window")
	}
}