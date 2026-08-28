package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/yourorg/vscodemarket/internal/geoip"
	"github.com/yourorg/vscodemarket/internal/storage"
)

func TestRecordEventDegradesOnBadJSON(t *testing.T) {
	// With real deps but a bad body the handler must still return 202.
	dir := t.TempDir()
	db, err := storage.Open(dir + "/t.db")
	if err != nil {
		t.Fatal(err)
	}
	r, _ := geoip.New("/nope")
	h := RecordEvent(DepsFrom(storage.NewEventRepo(db), r))
	req := httptest.NewRequest(http.MethodPost, "/api/v1/events", strings.NewReader("not-json"))
	w := httptest.NewRecorder()
	h(w, req)
	if w.Code != http.StatusAccepted {
		t.Errorf("want 202 even on bad JSON, got %d", w.Code)
	}
}

func TestRecordEventValidates(t *testing.T) {
	dir := t.TempDir()
	db, _ := storage.Open(dir + "/t.db")
	res, _ := geoip.New("/nope")
	h := RecordEvent(DepsFrom(storage.NewEventRepo(db), res))
	req := httptest.NewRequest(http.MethodPost, "/api/v1/events",
		strings.NewReader(`{"eventType":"FOO","targetType":"CLIENT","targetIdentifier":"1.0.0"}`))
	w := httptest.NewRecorder()
	h(w, req)
	if w.Code != http.StatusAccepted {
		t.Errorf("want 202 even on invalid enum, got %d", w.Code)
	}
}