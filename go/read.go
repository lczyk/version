package version

import "runtime/debug"

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
// meaningful module version (i.e. when Main.Version is empty or "(devel)",
// which happens for `go run` from a local tree, or builds outside a tagged
// commit). When the binary is built via `go install module@vX.Y.Z` or
// `go run module@vX.Y.Z`, Main.Version is the module version and is used.
//
// CommitSHA, BuildDate, BuildInfo come from the vcs.* settings the toolchain
// records (vcs.revision, vcs.time, vcs.modified). They will be empty for
// builds where VCS info is stripped (e.g. `-buildvcs=false`).
func Read(fallbackVersion string) Info {
	out := Info{Version: fallbackVersion}
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return out
	}
	if v := bi.Main.Version; v != "" && v != "(devel)" {
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
