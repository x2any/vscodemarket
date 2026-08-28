package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
)

type inferReq struct {
	UserAgent string `json:"userAgent"`
}

type inferResp struct {
	Platform     string `json:"platform"`
	Architecture string `json:"architecture"`
	Confidence   string `json:"confidence"` // HIGH | FALLBACK
}

// UAInfer handles POST /api/v1/ua/infer.
func UAInfer(w http.ResponseWriter, r *http.Request) {
	var req inferReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.UserAgent == "" {
		WriteError(w, http.StatusBadRequest, Err(CodeInvalidRequest, "缺少 userAgent", "missing userAgent"))
		return
	}
	p, a, conf := infer(req.UserAgent)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(inferResp{Platform: p, Architecture: a, Confidence: conf})
}

// infer returns (platform, architecture, confidence).
// Rules (kept intentionally tiny; expand only when observed false cases appear):
//   - OS: Mac → darwin; Windows → windows; else Linux.
//   - Arch: explicit arm64 / aarch64 / armv7 → matches;
//     x86_64 / WOW64 / Win64 → x86_64; otherwise default to x86_64.
//   - Confidence: HIGH if platform OR arch is explicitly matched;
//     FALLBACK only when the UA carries neither an OS nor arch hint
//     (e.g. "curl/8.4.0").
func infer(ua string) (string, string, string) {
	low := strings.ToLower(ua)
	platform := ""
	switch {
	case strings.Contains(low, "mac os x") || strings.Contains(low, "macintosh"):
		platform = "darwin"
	case strings.Contains(low, "windows"):
		platform = "windows"
	case strings.Contains(low, "linux") || strings.Contains(low, "x11"):
		platform = "linux"
	}

	arch := ""
	switch {
	case strings.Contains(low, "aarch64") || strings.Contains(low, "arm64") || strings.Contains(low, " arm "):
		// Apple Safari on Apple Silicon reports "Macintosh; ARM Mac OS X"
		// with a bare ARM token — treat as arm64.
		arch = "arm64"
	case strings.Contains(low, "armv7"):
		arch = "armv7"
	case strings.Contains(low, "wow64") || strings.Contains(low, "win64") || strings.Contains(low, "x86_64") || strings.Contains(low, "x64"):
		arch = "x86_64"
	}
	if platform == "" && arch == "" {
		return "linux", "x86_64", "FALLBACK"
	}
	if platform == "" {
		platform = "linux"
	}
	if arch == "" {
		arch = "x86_64"
	}
	return platform, arch, "HIGH"
}
