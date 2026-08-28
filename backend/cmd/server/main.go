package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/yourorg/vscodemarket/internal/handlers"
)

func main() {
	r := chi.NewRouter()

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
		})
		r.Post("/versions/lookup", handlers.VersionLookup)
		r.Post("/ua/infer", handlers.UAInfer)
	})

	addr := ":8081"
	log.Printf("vscodemarket backend listening on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}