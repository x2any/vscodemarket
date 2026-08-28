package upstream

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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
func FetchReleases(ctx context.Context, channel Channel, platform, arch string, page, pageSize int) ([]ReleaseEntry, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	url := fmt.Sprintf("https://update.code.visualstudio.com/api/releases/%s", channel)
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
	body, _ := io.ReadAll(resp.Body)
	var entries []releaseJSONEntry
	if err := json.Unmarshal(body, &entries); err != nil {
		return nil, 0, err
	}
	// Stable entries carry platform; Insider manifests do not.
	out := make([]ReleaseEntry, 0, len(entries))
	for _, e := range entries {
		entry := ReleaseEntry{
			Channel:     channel,
			Version:     e.Version,
			Platform:    e.Platform,
			Architecture: e.Architecture,
			CommitHash:  e.Commit,
			ReleaseDate: e.Date,
		}
		if platform != "" && entry.Platform != "" && !strings.EqualFold(entry.Platform, platform) {
			continue
		}
		if arch != "" && entry.Architecture != "" && !strings.EqualFold(entry.Architecture, arch) {
			continue
		}
		if platform != "" && arch != "" {
			entry.DownloadURL = buildDownloadURL(channel, e.Version, e.Platform, e.Architecture)
		} else {
			entry.DownloadURL = fmt.Sprintf("https://update.code.visualstudio.com/%s/%s",
				e.Version, channel)
		}
		out = append(out, entry)
	}
	total := len(out)
	start := min((page-1)*pageSize, total)
	end := min(start+pageSize, total)
	return out[start:end], total, nil
}

// releaseJSONEntry mirrors the official manifest payload.
// Upstream uses different fields for stable vs insider; we accept both.
type releaseJSONEntry struct {
	Version      string `json:"version"`
	Commit       string `json:"commit"`
	Platform     string `json:"platform"`
	Architecture string `json:"architecture"`
	Date         string `json:"date"`
}

func buildDownloadURL(channel Channel, version, platform, arch string) string {
	if channel == ChannelInsider {
		return fmt.Sprintf("https://update.code.visualstudio.com/%s/%s",
			version, channel)
	}
	return fmt.Sprintf("https://update.code.visualstudio.com/%s/%s-%s/%s",
		version, strings.ToLower(platform), strings.ToLower(arch), channel)
}