package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/yourorg/vscodemarket/internal/upstream"
)

// TestVersionLookupResponseUsesWhitelistedHosts guards SC-007:
// every downloadUrl in a successful response must come from the
// official-domain whitelist. The handler MUST refuse to surface any
// non-whitelisted URL.
func TestVersionLookupResponseUsesWhitelistedHosts(t *testing.T) {
	cases := []string{"1.94.2", "1.95.0", "1.100.0"}
	for _, v := range cases {
		body := `{"channel":"stable","version":"` + v + `","platform":"darwin","architecture":"arm64"}`
		r := httptest.NewRequest(http.MethodPost, "/api/v1/versions/lookup", strings.NewReader(body))
		w := httptest.NewRecorder()
		VersionLookup(w, r)
		// Either upstream succeeded (200) or upstream failed (502/404). In
		// the success path every emitted URL must validate.
		if w.Code != http.StatusOK {
			continue
		}
		out := w.Body.String()
		// Pull every "downloadUrl":"…" field and validate.
		for _, line := range strings.Split(out, ",") {
			if !strings.Contains(line, "downloadUrl") {
				continue
			}
			idx := strings.Index(line, "https://")
			if idx < 0 {
				continue
			}
			// Trim up to the closing quote.
			rest := line[idx:]
			if j := strings.IndexAny(rest, "\""); j > 0 {
				rest = rest[:j]
			}
			if err := upstream.AssertOfficial(rest); err != nil {
				t.Errorf("non-whitelisted URL leaked from handler: %q", rest)
			}
		}
	}
}