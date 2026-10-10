package version

import (
	"regexp"
	"runtime/debug"
	"strings"
)

// pseudoVersionRE matches Go pseudo-versions (v0.0.0-20060102150405-abcdef123456
// and the vX.Y.(Z+1)-0. / vX.Y.Z-pre.0. forms), as in golang.org/x/mod/module.
var pseudoVersionRE = regexp.MustCompile(`^v[0-9]+\.(0\.0-|\d+\.\d+-([^+]*\.)?0\.)\d{14}-[A-Za-z0-9]+(\+[0-9A-Za-z-]+(\.[0-9A-Za-z-]+)*)?$`)

// Info bundles build-time facts about a binary.
type Info struct {
	Version, CommitSHA, BuildDate, BuildInfo string
}

// String formats the Info via FormatVersion.
func (i Info) String() string {
	return FormatVersion(i.Version, i.CommitSHA, i.BuildDate, i.BuildInfo)
}

// Read derives Info from runtime/debug.ReadBuildInfo() — the metadata that the
// Go toolchain stamps into binaries from the VCS state at build time.
//
// fallbackVersion is used as Version when the runtime info reports no
// meaningful module version: Main.Version empty, "(devel)", or a
// pseudo-version (Go 1.24+ stamps local builds with one, e.g.
// v0.0.0-20260101000000-abcdef123456+dirty). A real tag, from
// `go install module@vX.Y.Z` or a local build at a tagged commit, is used with
// its leading "v" and any "+dirty" suffix dropped, so Version stays SemVer and
// dirtiness is reported once, in BuildInfo.
//
// CommitSHA, BuildDate, BuildInfo come from the vcs.* settings the toolchain
// records (vcs.revision, vcs.time, vcs.modified). BuildDate is therefore the
// commit time. They will be empty for builds where VCS info is stripped (e.g.
// `-buildvcs=false`).
func Read(fallbackVersion string) Info {
	out := Info{Version: fallbackVersion}
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return out
	}
	if v, ok := moduleVersion(bi.Main.Version); ok {
		out.Version = v
	}
	for _, s := range bi.Settings {
		switch s.Key {
		case "vcs.revision":
			out.CommitSHA = s.Value
		case "vcs.time":
			out.BuildDate = s.Value
		case "vcs.modified":
			if s.Value == "true" {
				out.BuildInfo = "dirty"
			}
		}
	}
	return out
}

// moduleVersion turns a toolchain-reported module version into a SemVer
// string, reporting false when it carries no release information.
func moduleVersion(v string) (string, bool) {
	if v == "" || v == "(devel)" || pseudoVersionRE.MatchString(v) {
		return "", false
	}
	v = strings.TrimSuffix(v, "+dirty")
	return strings.TrimPrefix(v, "v"), true
}
