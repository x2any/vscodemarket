package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestReleasesInvalidChannel(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/api/v1/releases?channel=edge", nil)
	w := httptest.NewRecorder()
	Releases(w, r)
	if w.Code != http.StatusBadRequest {
		t.Errorf("want 400, got %d", w.Code)
	}
}

func TestReleasesLiveOrBadGateway(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/api/v1/releases?channel=stable&page=1&pageSize=20", nil)
	w := httptest.NewRecorder()
	Releases(w, r)
	if w.Code != http.StatusOK && w.Code != http.StatusBadGateway {
		t.Errorf("want 200/502, got %d", w.Code)
	}
}