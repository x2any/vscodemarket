package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/yourorg/vscodemarket/internal/storage"
)

type trendingDeps struct {
	Repo *storage.EventRepo
}

// TrendingDepsFrom wires the repo for the trending handler.
func TrendingDepsFrom(r *storage.EventRepo) trendingDeps {
	return trendingDeps{Repo: r}
}

type trendingResp struct {
	TargetType string                    `json:"targetType"`
	Window     string                    `json:"window"`
	Items      []storage.TrendingRow     `json:"items"`
}

// Trending handles GET /api/v1/trending?targetType=&window=
func Trending(deps trendingDeps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		tt := q.Get("targetType")
		window := q.Get("window")
		if !validTargetType(tt) {
			WriteError(w, http.StatusBadRequest, Err(CodeInvalidRequest,
				"targetType 非法", "invalid targetType"))
			return
		}
		rows, err := deps.Repo.Trending(tt, window, 10)
		if err != nil {
			// Treat invalid window / DB hiccup as empty list per FR-016.
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			_ = json.NewEncoder(w).Encode(trendingResp{
				TargetType: tt, Window: window, Items: []storage.TrendingRow{},
			})
			return
		}
		if rows == nil {
			rows = []storage.TrendingRow{}
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(trendingResp{
			TargetType: tt, Window: window, Items: rows,
		})
	}
}