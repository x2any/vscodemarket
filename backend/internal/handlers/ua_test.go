package handlers

import "testing"

func TestInfer(t *testing.T) {
	cases := []struct {
		ua      string
		wantP   string
		wantA   string
		wantC   string
	}{
		{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36", "darwin", "x86_64", "HIGH"},
		{"Mozilla/5.0 (Macintosh; ARM Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Safari/605.1.15", "darwin", "arm64", "HIGH"},
		{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36", "windows", "x86_64", "HIGH"},
		{"Mozilla/5.0 (Windows NT 10.0; ARM64) AppleWebKit/537.36", "windows", "arm64", "HIGH"},
		{"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36", "linux", "x86_64", "HIGH"},
		{"Mozilla/5.0 (X11; Linux aarch64) AppleWebKit/537.36", "linux", "arm64", "HIGH"},
		{"Mozilla/5.0 (X11; Linux armv7l) AppleWebKit/537.36", "linux", "armv7", "HIGH"},
		{"curl/8.4.0", "linux", "x86_64", "FALLBACK"},
	}
	for _, c := range cases {
		p, a, conf := infer(c.ua)
		if p != c.wantP || a != c.wantA || conf != c.wantC {
			t.Errorf("infer(%q)=(%s,%s,%s) want (%s,%s,%s)", c.ua, p, a, conf, c.wantP, c.wantA, c.wantC)
		}
	}
}
