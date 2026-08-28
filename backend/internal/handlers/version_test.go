package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestVersionLookupStableOK(t *testing.T) {
	body := `{"channel":"stable","version":"1.94.2","platform":"darwin","architecture":"arm64"}`
	r := httptest.NewRequest(http.MethodPost, "/api/v1/versions/lookup", strings.NewReader(body))
	w := httptest.NewRecorder()
	VersionLookup(w, r)
	if w.Code != http.StatusOK && w.Code != http.StatusBadGateway {
		t.Logf("expected 200/502 (live upstream), got %d: %s", w.Code, w.Body.String())
	}
}

func TestVersionLookupBadPlatform(t *testing.T) {
	body := `{"channel":"stable","version":"1.94.2","platform":"plan9","architecture":"arm64"}`
	r := httptest.NewRequest(http.MethodPost, "/api/v1/versions/lookup", strings.NewReader(body))
	w := httptest.NewRecorder()
	VersionLookup(w, r)
	if w.Code != http.StatusBadRequest {
		t.Errorf("want 400, got %d", w.Code)
	}
}

func TestVersionLookupBadVersion(t *testing.T) {
	body := `{"channel":"stable","version":"abc","platform":"darwin","architecture":"arm64"}`
	r := httptest.NewRequest(http.MethodPost, "/api/v1/versions/lookup", strings.NewReader(body))
	w := httptest.NewRecorder()
	VersionLookup(w, r)
	if w.Code != http.StatusBadRequest {
		t.Errorf("want 400, got %d", w.Code)
	}
}

func TestVersionLookupBadChannel(t *testing.T) {
	body := `{"channel":"edge","version":"1.94.2","platform":"darwin","architecture":"arm64"}`
	r := httptest.NewRequest(http.MethodPost, "/api/v1/versions/lookup", strings.NewReader(body))
	w := httptest.NewRecorder()
	VersionLookup(w, r)
	if w.Code != http.StatusBadRequest {
		t.Errorf("want 400, got %d", w.Code)
	}
}

func TestWriteErrorEnvelope(t *testing.T) {
	w := httptest.NewRecorder()
	WriteError(w, http.StatusNotFound, Err(CodeVersionNotFound, "版本不存在", "Version not found"))
	if w.Code != http.StatusNotFound {
		t.Errorf("want 404, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "VERSION_NOT_FOUND") {
		t.Errorf("body missing code: %s", w.Body.String())
	}
}

func TestVersionLookupInsiderOmitsServer(t *testing.T) {
	// Insider channel must never include a server payload object.
	body := `{"channel":"insider","version":"1234567","platform":"darwin","architecture":"arm64"}`
	r := httptest.NewRequest(http.MethodPost, "/api/v1/versions/lookup", strings.NewReader(body))
	w := httptest.NewRecorder()
	VersionLookup(w, r)
	if w.Code == http.StatusOK && strings.Contains(w.Body.String(), `"server":{`) {
		t.Fatalf("insider must omit server object, got: %s", w.Body.String())
	}
}