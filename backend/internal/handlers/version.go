package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/yourorg/vscodemarket/internal/upstream"
)

type lookupReq struct {
	Channel string `json:"channel"`
	Version string `json:"version"`
	// Platform and Architecture are optional. When omitted the response
	// returns the full platform × architecture matrix so the visitor can
	// pick the file matching their local machine without guessing.
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
	Channel string          `json:"channel"`
	Version string          `json:"version"`
	Commit  string          `json:"commit,omitempty"`
	Clients []clientPayload `json:"clients"`
	Servers []serverPayload `json:"servers,omitempty"`
}

// VersionLookup handles POST /api/v1/versions/lookup.
// Returns the full matrix when Platform/Architecture are omitted.
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
	if (req.Platform == "") != (req.Architecture == "") {
		WriteError(w, http.StatusBadRequest, Err(CodeInvalidRequest,
			"platform 与 architecture 必须同时给或同时省略", "platform and architecture must both be set or both omitted"))
		return
	}
	if req.Platform != "" && !isValidPlatformArch(req.Platform, req.Architecture) {
		WriteError(w, http.StatusBadRequest, Err(CodeInvalidPlatformArch,
			"平台/架构组合无效", "Invalid platform/architecture combination"))
		return
	}
	if !upstream.IsValidStableVersion(req.Version) && req.Channel == "stable" {
		WriteError(w, http.StatusBadRequest, Err(CodeInvalidRequest, "Stable 版本号格式非法", "Invalid stable version format"))
		return
	}

	platforms, archs := expandMatrix(req.Channel, req.Platform, req.Architecture)
	if len(platforms) == 0 {
		WriteError(w, http.StatusBadRequest, Err(CodeInvalidPlatformArch,
			"该频道无可用平台/架构", "no platform/architecture available"))
		return
	}

	// Single existence check before fanning out (FR-002).
	if _, err := upstream.FetchStable(r.Context(), req.Version, platforms[0], archs[0]); err != nil {
		writeUpstreamError(w, err)
		return
	}

	resp := lookupResp{Channel: req.Channel, Version: req.Version}

	clients := make([]clientPayload, 0, len(platforms))
	for i, p := range platforms {
		cr, err := upstream.FetchStable(r.Context(), req.Version, p, archs[i])
		if err != nil || cr == nil {
			continue
		}
		if upstream.AssertOfficial(cr.DownloadURL) != nil {
			continue
		}
		clients = append(clients, clientPayload{
			DownloadURL:  cr.DownloadURL,
			Platform:     cr.Platform,
			Architecture: cr.Architecture,
		})
	}
	resp.Clients = clients

	if req.Channel == "stable" {
		// vscode-server only ships for stable; insider path returns clients only.
		commit, _ := upstream.CommitHashForVersion(r.Context(), req.Version)
		resp.Commit = commit
		if commit != "" {
			servers := make([]serverPayload, 0, len(platforms))
			for i, p := range platforms {
				sr, err := upstream.FetchServer(r.Context(), commit, req.Version, p, archs[i])
				if err != nil || sr == nil {
					continue
				}
				if upstream.AssertOfficial(sr.DownloadURL) != nil {
					continue
				}
				servers = append(servers, serverPayload{
					DownloadURL:   sr.DownloadURL,
					CommitHash:    sr.CommitHash,
					ClientVersion: sr.ClientVersion,
					Platform:      sr.Platform,
					Architecture:  sr.Architecture,
				})
			}
			resp.Servers = servers
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

// expandMatrix returns the (platforms, archs) pairs to fan out. When the
// caller passed an explicit pair, that's the only row. Otherwise we return
// the full ADR-0006 matrix. Insider client builds ship darwin+linux+windows
// but use a single URL per platform (no per-arch URLs upstream), so for
// insider we de-duplicate by platform.
func expandMatrix(channel, platform, arch string) ([]string, []string) {
	if platform != "" {
		return []string{platform}, []string{arch}
	}
	if channel == "insider" {
		// Insider client URLs are platform-only (no arch suffix upstream);
		// we still emit one entry per supported arch so the UI looks uniform,
		// all pointing to the same URL.
		rows := []struct{ p, a string }{
			{"windows", "x86_64"}, {"windows", "arm64"},
			{"linux", "x86_64"}, {"linux", "arm64"}, {"linux", "armv7"},
			{"darwin", "x86_64"}, {"darwin", "arm64"},
		}
		ps, as := make([]string, len(rows)), make([]string, len(rows))
		for i, r := range rows {
			ps[i], as[i] = r.p, r.a
		}
		return ps, as
	}
	rows := []struct{ p, a string }{
		{"windows", "x86_64"}, {"windows", "arm64"},
		{"linux", "x86_64"}, {"linux", "arm64"}, {"linux", "armv7"},
		{"darwin", "x86_64"}, {"darwin", "arm64"},
	}
	ps, as := make([]string, len(rows)), make([]string, len(rows))
	for i, r := range rows {
		ps[i], as[i] = r.p, r.a
	}
	return ps, as
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