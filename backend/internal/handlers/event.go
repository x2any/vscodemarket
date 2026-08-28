package handlers

import (
	"encoding/json"
	"net"
	"net/http"
	"strings"

	"github.com/yourorg/vscodemarket/internal/geoip"
	"github.com/yourorg/vscodemarket/internal/storage"
)

type eventReq struct {
	EventType        string `json:"eventType"`
	TargetType       string `json:"targetType"`
	TargetIdentifier string `json:"targetIdentifier"`
	Platform         string `json:"platform"`
	Architecture     string `json:"architecture"`
	Channel          string `json:"channel"`
}

type eventDeps struct {
	Repo     *storage.EventRepo
	Resolver *geoip.Resolver
}

// DepsFrom wires concrete adapters for the handler.
func DepsFrom(r *storage.EventRepo, g *geoip.Resolver) eventDeps {
	return eventDeps{Repo: r, Resolver: g}
}

// RecordEvent handles POST /api/v1/events.
// Failures are logged but never propagated — the caller always receives 202
// (per FR-012 / SC-005:埋点失败不能阻塞主请求).
func RecordEvent(deps eventDeps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req eventReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			// Still 202 — bad payload shouldn't block UX.
			w.WriteHeader(http.StatusAccepted)
			return
		}
		if !validEventType(req.EventType) || !validTargetType(req.TargetType) || strings.TrimSpace(req.TargetIdentifier) == "" {
			w.WriteHeader(http.StatusAccepted)
			return
		}
		ip := clientIP(r)
		ev := &storage.BehaviorEvent{
			EventType:        req.EventType,
			TargetType:       req.TargetType,
			TargetIdentifier: req.TargetIdentifier,
			Platform:         req.Platform,
			Architecture:     req.Architecture,
			Channel:          req.Channel,
			CountryCode:      deps.Resolver.Lookup(ip),
		}
		if err := deps.Repo.Insert(ev); err != nil {
			// Log only — main path continues.
			w.Header().Set("X-Event-Status", "degraded")
		}
		w.WriteHeader(http.StatusAccepted)
		_ = json.NewEncoder(w).Encode(map[string]bool{"accepted": true})
	}
}

func validEventType(s string) bool {
	return s == "SEARCH" || s == "DOWNLOAD"
}
func validTargetType(s string) bool {
	return s == "CLIENT" || s == "SERVER" || s == "EXTENSION"
}

func clientIP(r *http.Request) string {
	if v := r.Header.Get("X-Forwarded-For"); v != "" {
		if i := strings.IndexByte(v, ','); i > 0 {
			return strings.TrimSpace(v[:i])
		}
		return strings.TrimSpace(v)
	}
	if v := r.Header.Get("X-Real-IP"); v != "" {
		return strings.TrimSpace(v)
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}