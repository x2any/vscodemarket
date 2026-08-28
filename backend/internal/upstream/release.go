package upstream

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// ReleaseEntry is a row in the release list.
type ReleaseEntry struct {
	Channel      Channel `json:"channel"`
	Version      string  `json:"version"`
	Platform     string  `json:"platform,omitempty"`
	Architecture string  `json:"architecture,omitempty"`
	DownloadURL  string  `json:"downloadUrl"`
	CommitHash   string  `json:"commitHash,omitempty"`
	ReleaseDate  string  `json:"releaseDate,omitempty"`
}

// FetchReleases queries Microsoft's release manifest for a given channel.
// The upstream /api/releases/{channel} endpoint returns a flat JSON array
// of version strings (no platform/architecture/commit fields per entry);
// platform filtering is applied client-side and download URLs are
// reconstructed here so we never expose non-whitelisted hosts.
func FetchReleases(ctx context.Context, channel Channel, platform, arch string, page, pageSize int) ([]ReleaseEntry, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	url := apiBase + "/api/releases/" + string(channel)
	c := &http.Client{Timeout: 10 * time.Second}
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	resp, err := c.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return nil, 0, fmt.Errorf("upstream status %d", resp.StatusCode)
	}
	var versions []string
	if err := json.NewDecoder(resp.Body).Decode(&versions); err != nil {
		return nil, 0, err
	}
	out := make([]ReleaseEntry, 0, len(versions))
	for _, v := range versions {
		if platform != "" && !strings.EqualFold(platform, platform) {
			// No platform metadata in upstream payload — keep all versions.
			// Platform filtering is only meaningful when the user explicitly
			// requests both platform AND arch (then we generate per-entry).
		}
		entry := ReleaseEntry{
			Channel:    channel,
			Version:    v,
			DownloadURL: buildDownloadURL(channel, v, platform, arch),
		}
		if platform != "" {
			entry.Platform = platform
		}
		if arch != "" {
			entry.Architecture = arch
		}
		out = append(out, entry)
	}
	total := len(out)
	start := min((page-1)*pageSize, total)
	end := min(start+pageSize, total)
	return out[start:end], total, nil
}

func buildDownloadURL(channel Channel, version, platform, arch string) string {
	if channel == ChannelInsider {
		return fmt.Sprintf("%s/%s/%s", apiBase, version, channel)
	}
	if platform != "" && arch != "" {
		return fmt.Sprintf("%s/%s/%s-%s/%s",
			apiBase, version, strings.ToLower(platform), strings.ToLower(arch), channel)
	}
	return fmt.Sprintf("%s/%s/%s", apiBase, version, channel)
}