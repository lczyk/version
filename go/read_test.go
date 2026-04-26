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
