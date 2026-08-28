package storage

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// sqliteDialector keeps the GORM driver dependency in one place so callers
// don't import gorm.io/driver/sqlite directly.
func sqliteDialector(path string) gorm.Dialector {
	return sqlite.Open(path)
}