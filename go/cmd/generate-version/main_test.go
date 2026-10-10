package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

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

// gitT runs git in dir for test fixtures, failing the test on error.
func gitT(t *testing.T, dir string, env []string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-c", "user.name=t", "-c", "user.email=t@example.invalid", "-c", "commit.gpgsign=false"}, args...)...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), env...)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
}

func chdirT(t *testing.T, dir string) {
	t.Helper()
	orig, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(orig) })
}

func TestRunStableAcrossRegeneration(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "VERSION"), []byte("1.2.3\n"), 0644); err != nil {
		t.Fatal(err)
	}
	// A tracked placeholder, the setup that used to report "dirty" forever.
	out := filepath.Join(dir, "version.go")
	if err := os.WriteFile(out, []byte("package version\n"), 0644); err != nil {
		t.Fatal(err)
	}
	gitT(t, dir, nil, "init", "-q")
	gitT(t, dir, nil, "add", "VERSION", "version.go")
	gitT(t, dir, []string{"GIT_COMMITTER_DATE=2026-01-02T03:04:05+02:00"}, "commit", "-q", "-m", "chore: init")
	chdirT(t, dir)

	if err := run([]string{"-out", "version.go"}); err != nil {
		t.Fatal(err)
	}
	first, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{`Version   = "1.2.3"`, `BuildDate = "2026-01-02T01:04:05Z"`, `BuildInfo = ""`} {
		if !strings.Contains(string(first), want) {
			t.Errorf("generated file missing %s:\n%s", want, first)
		}
	}

	old := time.Unix(1, 0)
	if err := os.Chtimes(out, old, old); err != nil {
		t.Fatal(err)
	}
	if err := run([]string{"-out", "version.go"}); err != nil {
		t.Fatal(err)
	}
	second, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	if string(second) != string(first) {
		t.Errorf("regenerated file changed:\n%s\nvs\n%s", first, second)
	}
	if fi, err := os.Stat(out); err != nil || !fi.ModTime().Equal(old) {
		t.Errorf("unchanged file was rewritten")
	}

	if err := os.WriteFile(filepath.Join(dir, "VERSION"), []byte("1.2.3\n# edited\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := run([]string{"-out", "version.go"}); err != nil {
		t.Fatal(err)
	}
	third, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(third), `BuildInfo = "dirty"`) {
		t.Errorf("edit outside the output file not reported dirty:\n%s", third)
	}
}
