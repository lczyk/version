package main

import "testing"

func TestSemverRegex(t *testing.T) {
	tests := []struct {
		version string
		valid   bool
	}{
		// valid
		{"1.0.0", true},
		{"0.0.0", true},
		{"1.2.3", true},
		{"10.20.30", true},
		{"1.0.0-alpha", true},
		{"1.0.0-alpha.1", true},
		{"1.0.0-0.3.7", true},
		{"1.0.0-x.7.z.92", true},
		{"1.0.0+20130313144700", true},
		{"1.0.0-beta+exp.sha.5114f85", true},
		{"1.0.0+21AF26D3-117B344092BD", true},
		{"2.0.0-rc.1+build.123", true},
		// invalid
		{"1", false},
		{"1.2", false},
		{"v1.0.0", false},
		{"1.0.0.0", false},
		{"01.0.0", false},
		{"1.02.0", false},
		{"1.0.00", false},
		{"1.0.0-", false},
		{"1.0.0-+", false},
		{"1.0.0-+build", false},
		{"1.0.0-alpha..1", false},
		{"1.0.0-alpha.+build", false},
		{"", false},
		{"latest", false},
		{"main", false},
	}

	for _, tt := range tests {
		t.Run(tt.version, func(t *testing.T) {
			got := semverRE.MatchString(tt.version)
			if got != tt.valid {
				t.Errorf("semver %q: got valid=%v, want %v", tt.version, got, tt.valid)
			}
		})
	}
}
