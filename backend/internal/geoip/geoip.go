package geoip

import (
	"net"
	"os"
	"sync"

	"github.com/oschwald/geoip2-golang"
)

// Resolver looks up the ISO 3166-1 alpha-2 country code for an IP.
// When no mmdb is loaded (or lookup fails), it returns "UNKNOWN"
// per Constitution Principle I and NFR3.
type Resolver struct {
	mu   sync.RWMutex
	db   *geoip2.Reader
	path string
}

// New loads mmdb from path if the file exists; missing file is NOT an error.
func New(path string) (*Resolver, error) {
	r := &Resolver{path: path}
	if _, err := os.Stat(path); err == nil {
		db, err := geoip2.Open(path)
		if err != nil {
			return r, nil // LoadOrNil semantics — keep working without db.
		}
		r.db = db
	}
	return r, nil
}

// Close releases the underlying mmdb handle.
func (r *Resolver) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.db != nil {
		return r.db.Close()
	}
	return nil
}

// Lookup returns the ISO alpha-2 country code, or "UNKNOWN".
func (r *Resolver) Lookup(ipStr string) string {
	r.mu.RLock()
	db := r.db
	r.mu.RUnlock()
	if db == nil || ipStr == "" {
		return "UNKNOWN"
	}
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return "UNKNOWN"
	}
	record, err := db.Country(ip)
	if err != nil || record == nil || record.Country.IsoCode == "" {
		return "UNKNOWN"
	}
	return record.Country.IsoCode
}
