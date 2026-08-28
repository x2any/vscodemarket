package upstream

import (
	"context"
	"fmt"
	"strings"
)

// ServerRelease is a vscode-server build bound to a client commit.
type ServerRelease struct {
	CommitHash    string `json:"commitHash"`
	ClientVersion string `json:"clientVersion"`
	Platform      string `json:"platform"`
	Architecture  string `json:"architecture"`
	DownloadURL   string `json:"downloadUrl"`
}

// FetchServer builds the official vscode-server URL for a stable client commit.
func FetchServer(ctx context.Context, commit, clientVersion, platform, arch string) (*ServerRelease, error) {
	if commit == "" {
		return nil, ErrVersionNotFound
	}
	dl := fmt.Sprintf("https://update.code.visualstudio.com/commit:%s/server-%s-%s/stable",
		commit, strings.ToLower(platform), strings.ToLower(arch))
	return &ServerRelease{
		CommitHash:    commit,
		ClientVersion: clientVersion,
		Platform:      platform,
		Architecture:  arch,
		DownloadURL:   dl,
	}, nil
}