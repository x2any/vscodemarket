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
//     x86_64 / WOW64 / Win64 → x86_64; otherwise FALLBACK → x86_64.
func infer(ua string) (string, string, string) {
	low := strings.ToLower(ua)
	platform := "linux"
	switch {
	case strings.Contains(low, "mac os x") || strings.Contains(low, "macintosh"):
		platform = "darwin"
	case strings.Contains(low, "windows"):
		platform = "windows"
	}

	arch := "x86_64"
	confidence := "FALLBACK"
	switch {
	case strings.Contains(low, "arm64") || strings.Contains(low, "aarch64"):
		arch = "arm64"
		confidence = "HIGH"
	case strings.Contains(low, "armv7"):
		arch = "armv7"
		confidence = "HIGH"
	case strings.Contains(low, "wow64") || strings.Contains(low, "win64") || strings.Contains(low, "x86_64") || strings.Contains(low, "x64"):
		arch = "x86_64"
		confidence = "HIGH"
	}
	return platform, arch, confidence
}
