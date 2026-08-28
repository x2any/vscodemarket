package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/yourorg/vscodemarket/internal/upstream"
)

type releasesResp struct {
	Results  []upstream.ReleaseEntry `json:"results"`
	Page     int                      `json:"page"`
	PageSize int                      `json:"pageSize"`
	Total    int                      `json:"total"`
}

// Releases handles GET /api/v1/releases
// Accepts channel + page + pageSize only. Platform / architecture are
// intentionally absent here — version detail lives at /versions/lookup.
func Releases(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	channel := upstream.Channel(strings.ToLower(q.Get("channel")))
	if channel == "" {
		channel = upstream.ChannelStable
	}
	if channel != upstream.ChannelStable && channel != upstream.ChannelInsider {
		WriteError(w, http.StatusBadRequest, Err(CodeInvalidRequest, "channel 非法", "invalid channel"))
		return
	}
	page := atoiDefault(q.Get("page"), 1)
	pageSize := atoiDefault(q.Get("pageSize"), 30)
	results, total, err := upstream.FetchReleases(r.Context(), channel, "", "", page, pageSize)
	if err != nil {
		WriteError(w, http.StatusBadGateway, Err(CodeUpstreamFailure, "版本清单上游失败", "releases upstream failure"))
		return
	}
	if results == nil {
		results = []upstream.ReleaseEntry{}
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(releasesResp{
		Results:  results,
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	})
}