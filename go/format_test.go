package version

import "testing"

func TestFormatVersion(t *testing.T) {
	tests := []struct {
		name      string
		version   string
		commitSHA string
		buildDate string
		buildInfo string
		want      string
	}{
		{
			name:    "version only",
			version: "1.0.0",
			want:    "1.0.0",
		},
		{
			name:      "version with commit",
			version:   "1.0.0",
			commitSHA: "abc1234567890",
			want:      "1.0.0 @ abc1234",
		},
		{
			name:      "version with commit and date",
			version:   "1.0.0",
			commitSHA: "abc1234567890",
			buildDate: "2026-01-01T00:00:00Z",
			want:      "1.0.0 @ abc1234 (2026-01-01T00:00:00Z)",
		},
		{
			name:      "version with commit, date, and dirty",
			version:   "1.0.0",
			commitSHA: "abc1234567890",
			buildDate: "2026-01-01T00:00:00Z",
			buildInfo: "dirty",
			want:      "1.0.0 @ abc1234 (2026-01-01T00:00:00Z, dirty)",
		},
		{
			name:      "dirty without date",
			version:   "1.0.0",
			commitSHA: "abc1234567890",
			buildInfo: "dirty",
			want:      "1.0.0 @ abc1234 (dirty)",
		},
		{
			name:      "version with short commit SHA",
			version:   "1.0.0",
			commitSHA: "abc",
			want:      "1.0.0 @ abc",
		},
		{
			name:      "version with 7-char commit SHA",
			version:   "1.0.0",
			commitSHA: "abcdefg",
			want:      "1.0.0 @ abcdefg",
		},
		{
			name:      "date without commit",
			version:   "1.0.0",
			buildDate: "2026-01-01T00:00:00Z",
			want:      "1.0.0 (2026-01-01T00:00:00Z)",
		},
		{
			name:      "dirty without commit or date",
			version:   "1.0.0",
			buildInfo: "dirty",
			want:      "1.0.0 (dirty)",
		},
		{
			name:    "pre-release semver passes through",
			version: "1.2.3-rc1",
			want:    "1.2.3-rc1",
		},
		{
			name:      "8-char commit truncates to 7",
			version:   "1.0.0",
			commitSHA: "abcdefgh",
			want:      "1.0.0 @ abcdefg",
		},
		{
			name:      "pre-release with commit",
			version:   "1.0.0-alpha.1",
			commitSHA: "abc1234567890",
			want:      "1.0.0-alpha.1 @ abc1234",
		},
		{
			name:      "pre-release with full metadata",
			version:   "1.0.0-rc.1",
			commitSHA: "abc1234567890",
			buildDate: "2026-01-01T00:00:00Z",
			buildInfo: "dirty",
			want:      "1.0.0-rc.1 @ abc1234 (2026-01-01T00:00:00Z, dirty)",
		},
		{
			name:      "buildinfo with multiple words",
			version:   "1.0.0",
			commitSHA: "abc1234567890",
			buildInfo: "dirty, uncommitted",
			want:      "1.0.0 @ abc1234 (dirty, uncommitted)",
		},
		{
			name:      "1-char commit",
			version:   "1.0.0",
			commitSHA: "a",
			want:      "1.0.0 @ a",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatVersion(tt.version, tt.commitSHA, tt.buildDate, tt.buildInfo)
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}
