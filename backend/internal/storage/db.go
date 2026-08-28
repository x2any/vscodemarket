package storage

import (
	"time"

	"gorm.io/gorm"
)

// BehaviorEvent is the persistent record for analytics.
// NOTE: We never store the raw IP — only CountryCode resolved via GeoIP.
type BehaviorEvent struct {
	ID               uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	EventType        string    `gorm:"size:16;index;not null" json:"eventType"`
	TargetType       string    `gorm:"size:16;index;not null" json:"targetType"`
	TargetIdentifier string    `gorm:"size:256;index;not null" json:"targetIdentifier"`
	Platform         string    `gorm:"size:16" json:"platform,omitempty"`
	Architecture     string    `gorm:"size:16" json:"architecture,omitempty"`
	Channel          string    `gorm:"size:16" json:"channel,omitempty"`
	CountryCode      string    `gorm:"size:8;index" json:"countryCode"`
	CreatedAt        time.Time `gorm:"index" json:"createdAt"`
}

func (BehaviorEvent) TableName() string { return "behavior_events" }

// Open initializes a SQLite-backed GORM DB and migrates.
func Open(path string) (*gorm.DB, error) {
	db, err := gorm.Open(sqliteDialector(path), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&BehaviorEvent{}); err != nil {
		return nil, err
	}
	return db, nil
}