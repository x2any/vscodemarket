package upstream

import (
	"net/url"
	"strings"
)

// Allowed download hosts per Constitution Principle II.
// SC-007: every downloadUrl in a response must match this set.
var allowedHosts = map[string]struct{}{
	"update.code.visualstudio.com": {},
	"marketplace.visualstudio.com":  {},
	"vscode.download.prss.microsoft.com": {},
	"aka.ms": {},
}

// AssertOfficial rejects any URL whose host is not in the whitelist.
// Returned error message is safe for logs and 5xx responses.
func AssertOfficial(raw string) error {
	u, err := url.Parse(raw)
	if err != nil {
		return ErrNonOfficial
	}
	host := strings.ToLower(u.Host)
	if _, ok := allowedHosts[host]; !ok {
		return ErrNonOfficial
	}
	return nil
}

// ErrNonOfficial signals a non-whitelisted host.
var ErrNonOfficial = nonOfficialError{}

type nonOfficialError struct{}

func (nonOfficialError) Error() string { return "non-official download host blocked" }
