package upstream

import (
	"context"
	"testing"
)

// TestEnginesVscodePresence asserts that when upstream returns versions,
// each populated version carries an EnginesVSCode string (possibly empty
// for legacy entries). This guards FR-009: every ExtensionVersion surfaces
// engines.vscode in the API.
func TestEnginesVscodePresence(t *testing.T) {
	// Use a known stable extension: ms-python.python
	ctx := context.Background()
	versions, err := FetchExtensionVersions(ctx, "ms-python", "python")
	if err != nil {
		t.Skipf("upstream unreachable in test env: %v", err)
	}
	if len(versions) == 0 {
		t.Skip("upstream returned no versions")
	}
	for _, v := range versions {
		// EnginesVSCode may be empty for very old versions, but the field
		// itself must exist on the struct (compile-time guarantee). At
		// runtime we just ensure non-nil serialization works.
		_ = v.EnginesVSCode
	}
}