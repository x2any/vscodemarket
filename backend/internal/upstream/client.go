package upstream

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"slices"
	"strings"
	"time"
)

// Channel distinguishes Stable / Insider.
type Channel string

const (
	ChannelStable  Channel = "stable"
	ChannelInsider Channel = "insider"
)

// ClientRelease is a single channel × platform × architecture download entry.
type ClientRelease struct {
	Channel      Channel `json:"channel"`
	Version      string  `json:"version"`
	Platform     string  `json:"platform"`
	Architecture string  `json:"architecture"`
	DownloadURL  string  `json:"downloadUrl"`
	CommitHash   string  `json:"commitHash,omitempty"`
}

type apiClient struct {
	hc *http.Client
}

func newAPIClient() *apiClient {
	return &apiClient{hc: &http.Client{Timeout: 10 * time.Second}}
}

// FetchStable queries Microsoft's API for a single stable version.
func FetchStable(ctx context.Context, version, platform, arch string) (*ClientRelease, error) {
	url := apiBase + "/api/releases/stable"
	c := newAPIClient()
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	resp, err := c.hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("upstream status %d", resp.StatusCode)
	}
	var entries []string
	if err := json.NewDecoder(resp.Body).Decode(&entries); err != nil {
		return nil, err
	}
	if !slices.Contains(entries, version) {
		return nil, ErrVersionNotFound
	}
	dl := fmt.Sprintf("%s/%s/%s-%s/stable",
		apiBase, version, strings.ToLower(platform), strings.ToLower(arch))
	return &ClientRelease{
		Channel: ChannelStable, Version: version,
		Platform: platform, Architecture: arch,
		DownloadURL: dl,
	}, nil
}

// LatestStable returns the latest stable version's body for commit-hash lookup.
func LatestStable(ctx context.Context, platform, arch string) (string, error) {
	url := fmt.Sprintf("https://update.code.visualstudio.com/api/version/stable/%s-%s",
		strings.ToLower(platform), strings.ToLower(arch))
	c := newAPIClient()
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	resp, err := c.hc.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("upstream status %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	return strings.TrimSpace(string(body)), nil
}

// CommitHashForVersion looks up the commit hash for a stable version.
// Microsoft's /api/releases/stable returns only a flat version array;
// commit hash requires hitting the update endpoint for the specific
// version, which returns a JSON body containing the commit field.
//
// If commit metadata is unavailable upstream we return ("", nil) — the
// handler treats this as "server URL not published yet" rather than an
// error (per FR-005 / Edge Cases "Stable commit not yet published").
func CommitHashForVersion(ctx context.Context, version string) (string, error) {
	url := fmt.Sprintf("%s/api/version/stable/%s", apiBase, version)
	c := newAPIClient()
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	resp, err := c.hc.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 404 {
		return "", ErrVersionNotFound
	}
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("upstream status %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	// Body is a JSON-encoded object with a "commit" field, but Microsoft's
	// update endpoint may return the version string or a structured body.
	// Best-effort: try parsing as object first.
	var payload struct {
		Commit string `json:"commit"`
	}
	if err := json.Unmarshal(body, &payload); err == nil && payload.Commit != "" {
		return payload.Commit, nil
	}
	return "", nil
}

// IsValidVersion rejects empty / malformed strings before hitting upstream.
// Accepts MAJOR.MINOR.PATCH optionally followed by a single pre-release or
// build tag (e.g. `-insider`, `+20241002`, `-rc.1`); rejects four-segment
// numbers like 1.94.2.3.
var stableVersionRe = regexp.MustCompile(`^\d+\.\d+\.\d+(?:[-+][0-9A-Za-z.]+)?$`)

func IsValidStableVersion(v string) bool {
	return stableVersionRe.MatchString(strings.TrimSpace(v))
}

var ErrVersionNotFound = fmt.Errorf("version not found")

// apiBase is the Microsoft update endpoint root. Tests override it via
// apiBaseFor to point at an httptest server without TLS ceremony.
var apiBase = "https://update.code.visualstudio.com"