package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/yourorg/vscodemarket/internal/geoip"
	"github.com/yourorg/vscodemarket/internal/handlers"
	"github.com/yourorg/vscodemarket/internal/storage"
	"github.com/yourorg/vscodemarket/internal/sweeper"
)

func main() {
	dbPath := envOr("DB_PATH", "/data/vscodemarket.db")
	geoipPath := envOr("GEOIP_PATH", "/data/GeoLite2-Country.mmdb")

	db, err := storage.Open(dbPath)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	repo := storage.NewEventRepo(db)

	resolver, err := geoip.New(geoipPath)
	if err != nil {
		log.Printf("geoip: failed to load %s, continuing with UNKNOWN-only mode: %v", geoipPath, err)
	}
	defer resolver.Close()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	go sweeper.Run(ctx, repo, 24*time.Hour, 90*24*time.Hour)

	r := chi.NewRouter()
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
		})
		r.Post("/versions/lookup", handlers.VersionLookup)
		r.Post("/ua/infer", handlers.UAInfer)
		r.Get("/extensions/search", handlers.ExtensionSearch)
		r.Get("/extensions/{pub}/{name}/versions", handlers.ExtensionVersions)
		r.Get("/extensions/{pub}/{name}/versions/{ver}", handlers.ExtensionVersion)
		r.Get("/releases", handlers.Releases)
		r.Post("/events", handlers.RecordEvent(handlers.DepsFrom(repo, resolver)))
		r.Get("/trending", handlers.Trending(handlers.TrendingDepsFrom(repo)))
	})

	addr := ":8081"
	log.Printf("vscodemarket backend listening on %s (db=%s, geoip=%s)", addr, dbPath, geoipPath)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}

func envOr(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}