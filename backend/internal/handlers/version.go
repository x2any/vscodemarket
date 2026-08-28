package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/yourorg/vscodemarket/internal/upstream"
)

type lookupReq struct {
	Channel      string `json:"channel"`
	Version      string `json:"version"`
	Platform     string `json:"platform"`
	Architecture string `json:"architecture"`
}

type clientPayload struct {
	DownloadURL  string `json:"downloadUrl"`
	Platform     string `json:"platform"`
	Architecture string `json:"architecture"`
}

type serverPayload struct {
	DownloadURL   string `json:"downloadUrl"`
	CommitHash    string `json:"commitHash"`
	ClientVersion string `json:"clientVersion"`
	Platform      string `json:"platform"`
	Architecture  string `json:"architecture"`
}

type lookupResp struct {
	Channel string         `json:"channel"`
	Version string         `json:"version"`
	Client  clientPayload  `json:"client"`
	Server  *serverPayload `json:"server,omitempty"`
}

// VersionLookup handles POST /api/v1/versions/lookup.
func VersionLookup(w http.ResponseWriter, r *http.Request) {
	var req lookupReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, Err(CodeInvalidRequest, "请求体非法", "Invalid request body"))
		return
	}
	if !isValidChannel(req.Channel) {
		WriteError(w, http.StatusBadRequest, Err(CodeInvalidRequest, "channel 必须是 stable 或 insider", "channel must be stable or insider"))
		return
	}
	if !isValidPlatformArch(req.Platform, req.Architecture) {
		WriteError(w, http.StatusBadRequest, Err(CodeInvalidPlatformArch,
			"平台/架构组合无效", "Invalid platform/architecture combination"))
		return
	}
	if !upstream.IsValidStableVersion(req.Version) && req.Channel == "stable" {
		WriteError(w, http.StatusBadRequest, Err(CodeInvalidRequest, "Stable 版本号格式非法", "Invalid stable version format"))
		return
	}

	client, err := upstream.FetchStable(r.Context(), req.Version, req.Platform, req.Architecture)
	if err != nil {
		writeUpstreamError(w, err)
		return
	}

	resp := lookupResp{
		Channel: req.Channel,
		Version: req.Version,
		Client: clientPayload{
			DownloadURL:  client.DownloadURL,
			Platform:     client.Platform,
			Architecture: client.Architecture,
		},
	}
	if req.Channel == "stable" {
		// vscode-server only ships for stable; insider path returns client only.
		commit, cerr := upstream.CommitHashForVersion(r.Context(), req.Version)
		if cerr == nil && commit != "" {
			srv, serr := upstream.FetchServer(r.Context(), commit, req.Version, req.Platform, req.Architecture)
			if serr == nil {
				if werr := upstream.AssertOfficial(srv.DownloadURL); werr == nil {
					resp.Server = &serverPayload{
						DownloadURL:   srv.DownloadURL,
						CommitHash:    srv.CommitHash,
						ClientVersion: srv.ClientVersion,
						Platform:      srv.Platform,
						Architecture:  srv.Architecture,
					}
				}
			}
		}
	}

	// Whitelist check on client URL too.
	if err := upstream.AssertOfficial(resp.Client.DownloadURL); err != nil {
		WriteError(w, http.StatusInternalServerError, Err(CodeNonOfficialURLBlocked,
			"检测到非官方客户端直链,已拒绝", "Non-official client URL blocked"))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

func writeUpstreamError(w http.ResponseWriter, err error) {
	if err == upstream.ErrVersionNotFound {
		WriteError(w, http.StatusNotFound, Err(CodeVersionNotFound, "版本不存在", "Version not found"))
		return
	}
	WriteError(w, http.StatusBadGateway, Err(CodeUpstreamFailure, "上游服务暂时不可用", "Upstream failure"))
}

func isValidChannel(c string) bool { return c == "stable" || c == "insider" }

// isValidPlatformArch implements the ADR-0006 matrix subset we expose.
func isValidPlatformArch(p, a string) bool {
	switch p {
	case "windows":
		return a == "x86_64" || a == "arm64"
	case "linux":
		return a == "x86_64" || a == "arm64" || a == "armv7"
	case "darwin":
		return a == "x86_64" || a == "arm64"
	}
	return false
}