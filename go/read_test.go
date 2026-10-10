package version

import "testing"

func TestReadFallback(t *testing.T) {
	// In `go test` the binary's Main.Version is "(devel)", so fallback is used.
	info := Read("9.9.9-fallback")
	if info.Version != "9.9.9-fallback" {
		t.Errorf("Version = %q, want fallback %q", info.Version, "9.9.9-fallback")
	}
}

func TestInfoString(t *testing.T) {
	info := Info{Version: "1.0.0", CommitSHA: "abc1234567", BuildDate: "2026-01-01T00:00:00Z", BuildInfo: "dirty"}
	want := "1.0.0 @ abc1234 (2026-01-01T00:00:00Z, dirty)"
	if got := info.String(); got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}

func TestModuleVersion(t *testing.T) {
	tests := []struct {
		in   string
		want string
		ok   bool
	}{
		{"", "", false},
		{"(devel)", "", false},
		{"v0.0.0-20260916141054-30e9d853cd61", "", false},
		{"v0.0.0-20260916141054-30e9d853cd61+dirty", "", false},
		{"v1.2.4-0.20260916141054-30e9d853cd61", "", false},
		{"v1.2.3-rc.1.0.20260916141054-30e9d853cd61+dirty", "", false},
		{"v1.2.3", "1.2.3", true},
		{"v1.2.3+dirty", "1.2.3", true},
		{"v1.2.3-rc.1", "1.2.3-rc.1", true},
		{"v1.2.3+build.5", "1.2.3+build.5", true},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got, ok := moduleVersion(tt.in)
			if got != tt.want || ok != tt.ok {
				t.Errorf("moduleVersion(%q) = %q, %v; want %q, %v", tt.in, got, ok, tt.want, tt.ok)
			}
		})
	}
}
