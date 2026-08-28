package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestExtensionSearchMissingQuery(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/api/v1/extensions/search", nil)
	w := httptest.NewRecorder()
	ExtensionSearch(w, r)
	if w.Code != http.StatusBadRequest {
		t.Errorf("want 400, got %d", w.Code)
	}
}

func TestExtensionSearchLiveOrBadGateway(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/api/v1/extensions/search?q=python", nil)
	w := httptest.NewRecorder()
	ExtensionSearch(w, r)
	if w.Code != http.StatusOK && w.Code != http.StatusBadGateway {
		t.Errorf("want 200/502, got %d: %s", w.Code, w.Body.String())
	}
	if w.Code == http.StatusOK {
		var body searchResp
		if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
			t.Errorf("decode: %v", err)
		}
	}
}

func TestExtensionVersionRequiresFields(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/api/v1/extensions//x/versions/1", nil)
	w := httptest.NewRecorder()
	ExtensionVersion(w, r)
	if w.Code != http.StatusBadRequest {
		t.Errorf("want 400, got %d", w.Code)
	}
}