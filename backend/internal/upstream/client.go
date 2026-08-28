package upstream

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
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
	url := fmt.Sprintf("https://update.code.visualstudio.com/api/version/stable/%s-%s",
		strings.ToLower(platform), strings.ToLower(arch))
	c := newAPIClient()
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	resp, err := c.hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 404 {
		return nil, ErrVersionNotFound
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("upstream status %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	// Microsoft returns the version string when found, empty body otherwise.
	if strings.TrimSpace(string(body)) == "" {
		return nil, ErrVersionNotFound
	}
	if !strings.EqualFold(strings.TrimSpace(string(body)), version) {
		return nil, ErrVersionNotFound
	}
	dl := fmt.Sprintf("https://update.code.visualstudio.com/%s/%s-%s/stable",
		version, strings.ToLower(platform), strings.ToLower(arch))
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

// CommitHashForVersion fetches the commit hash associated with a stable build.
// Microsoft's release.json endpoint exposes the commit field.
func CommitHashForVersion(ctx context.Context, version string) (string, error) {
	url := "https://update.code.visualstudio.com/api/releases/stable"
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
	var entries []releaseJSONEntry
	if err := json.NewDecoder(resp.Body).Decode(&entries); err != nil {
		return "", err
	}
	for _, e := range entries {
		if e.Version == version && e.Commit != "" {
			return e.Commit, nil
		}
	}
	return "", ErrVersionNotFound
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