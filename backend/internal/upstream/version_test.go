package upstream

import "testing"

func TestIsValidStableVersion(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"1.94.2", true},
		{"1.94.2+20241002", true},
		{"1.94.2-insider", true},
		{"", false},
		{"1.94", false},
		{"1.94.2.3", false},
		{"abc", false},
		{" 1.94.2 ", true},
	}
	for _, c := range cases {
		if got := IsValidStableVersion(c.in); got != c.want {
			t.Errorf("IsValidStableVersion(%q)=%v want %v", c.in, got, c.want)
		}
	}
}