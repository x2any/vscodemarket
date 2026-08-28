package geoip

import "testing"

// TestUnknownWhenNoDatabase asserts the LoadOrNil semantics: with no mmdb
// the resolver must always return "UNKNOWN" rather than erroring.
func TestUnknownWhenNoDatabase(t *testing.T) {
	r, err := New("/nonexistent/GeoLite2-Country.mmdb")
	if err != nil {
		t.Fatalf("constructor should not error on missing file: %v", err)
	}
	if got := r.Lookup("8.8.8.8"); got != "UNKNOWN" {
		t.Errorf("want UNKNOWN, got %q", got)
	}
	if got := r.Lookup(""); got != "UNKNOWN" {
		t.Errorf("empty IP should return UNKNOWN, got %q", got)
	}
	if got := r.Lookup("not-an-ip"); got != "UNKNOWN" {
		t.Errorf("invalid IP should return UNKNOWN, got %q", got)
	}
	_ = r.Close()
}
