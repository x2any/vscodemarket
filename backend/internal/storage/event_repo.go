package storage

import (
	"time"

	"gorm.io/gorm"
)

type EventRepo struct{ db *gorm.DB }

func NewEventRepo(db *gorm.DB) *EventRepo { return &EventRepo{db: db} }

func (r *EventRepo) Insert(e *BehaviorEvent) error {
	if e.CreatedAt.IsZero() {
		e.CreatedAt = time.Now().UTC()
	}
	if e.CountryCode == "" {
		e.CountryCode = "UNKNOWN"
	}
	return r.db.Create(e).Error
}

// Trending aggregates top-N targetIdentifier counts within a window.
// Empty windows return an empty slice (HTTP 200 with empty items).
func (r *EventRepo) Trending(targetType, window string, limit int) ([]TrendingRow, error) {
	if limit <= 0 || limit > 10 {
		limit = 10
	}
	since, err := windowSince(window)
	if err != nil {
		return nil, err
	}
	var rows []TrendingRow
	q := r.db.Model(&BehaviorEvent{}).
		Select("target_identifier as target_identifier, count(*) as cnt").
		Where("target_type = ? AND created_at >= ?", targetType, since).
		Group("target_identifier").
		Order("cnt DESC").
		Limit(limit)
	if err := q.Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

type TrendingRow struct {
	TargetIdentifier string `json:"targetIdentifier"`
	Cnt              int64  `json:"count"`
}

// SweepOlderThan deletes events whose CreatedAt is strictly before cutoff.
func (r *EventRepo) SweepOlderThan(cutoff time.Time) (int64, error) {
	res := r.db.Where("created_at < ?", cutoff).Delete(&BehaviorEvent{})
	return res.RowsAffected, res.Error
}

func windowSince(window string) (time.Time, error) {
	now := time.Now().UTC()
	switch window {
	case "24h":
		return now.Add(-24 * time.Hour), nil
	case "7d":
		return now.Add(-7 * 24 * time.Hour), nil
	case "30d":
		return now.Add(-30 * 24 * time.Hour), nil
	}
	return now, ErrInvalidWindow
}

var ErrInvalidWindow = errInvalidWindow{}

type errInvalidWindow struct{}

func (errInvalidWindow) Error() string { return "invalid time window" }