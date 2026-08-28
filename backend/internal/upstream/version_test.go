package upstream

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIsValidStableVersion(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"1.94.2", true},
		{"1.94.2+20241002", true},
		{"1.94.2-insider", true},
		{"", false},
		{"1.94", false},
		{"1.94.2.3", false},
		{"abc", false},
		{" 1.94.2 ", true},
	}
	for _, c := range cases {
		if got := IsValidStableVersion(c.in); got != c.want {
			t.Errorf("IsValidStableVersion(%q)=%v want %v", c.in, got, c.want)
		}
	}
}

// TestFetchStableFindsVersion regression: previously FetchStable hit
// /api/version/stable/{platform}-{arch} which returns the latest version
// (not the queried one), so any non-latest version (e.g. 1.129.0) was
// wrongly reported as not found. The endpoint itself is now 404 upstream.
// FetchStable must look the requested version up in /api/releases/stable.
func TestFetchStableFindsVersion(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/releases/stable" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`["1.135.0","1.134.0","1.129.0","1.94.2"]`))
	}))
	defer srv.Close()

	orig := apiBase
	apiBase = srv.URL
	t.Cleanup(func() { apiBase = orig })

	rel, err := FetchStable(context.Background(), "1.129.0", "darwin", "arm64")
	if err != nil {
		t.Fatalf("FetchStable(1.129.0) err=%v, want nil", err)
	}
	if rel.Version != "1.129.0" {
		t.Errorf("Version=%q, want 1.129.0", rel.Version)
	}
	if rel.DownloadURL == "" {
		t.Errorf("DownloadURL empty")
	}
}

// TestFetchStableMissingVersion confirms unknown versions still surface as
// ErrVersionNotFound (the contract the handler relies on).
func TestFetchStableMissingVersion(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`["1.135.0","1.134.0"]`))
	}))
	defer srv.Close()

	orig := apiBase
	apiBase = srv.URL
	t.Cleanup(func() { apiBase = orig })

	if _, err := FetchStable(context.Background(), "9.9.9", "darwin", "arm64"); err != ErrVersionNotFound {
		t.Errorf("FetchStable(9.9.9) err=%v, want ErrVersionNotFound", err)
	}
}