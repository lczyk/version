# lczyk's version format spec

## 1. scope

defines:

1. **version string format** from four input parts.
2. **VERSION file format** read by build-time generators.
3. **build-time inputs** generator MUST collect.
4. **generated artifact contract** (four named values exposed to host program).

`<version>` field MUST be valid [SemVer 2.0.0](https://semver.org) string. formatter no validate; caller job. output = *display* string, not SemVer string — SemVer prefix trivially extractable (see §4.2).

not defined: SemVer grammar (see linked spec), git internals, distribution mechanics.

## 2. terminology

- MUST / SHOULD / MAY per RFC 2119.
- "empty" = string length zero. whitespace-only NOT empty; trim = caller job.

## 3. inputs

`FormatVersion` takes four ordered string params:

| name        | meaning                                                  |
|-------------|----------------------------------------------------------|
| `version`   | project version (e.g. `1.2.3`, `1.2.3-rc1`). MUST be SemVer 2.0.0. |
| `commitSHA` | full git commit SHA, hex. may be empty.                  |
| `buildDate` | build timestamp, ISO-8601 / RFC 3339 UTC. may be empty.  |
| `buildInfo` | free-form build annotation (e.g. `dirty`). may be empty. |

all four strings. formatter no further validation.

## 4. output format

grammar (ABNF-ish):

```
output     = semver [ SP "@" SP commit7 ] [ SP "(" details ")" ]
details    = build-date [ ", " build-info ]
           / build-date
           / build-info
semver     = <SemVer 2.0.0 version string>   ; opaque to this spec
commit7    = 1*7(HEXDIG)                     ; first min(7, len(commitSHA)) chars of commitSHA
SP         = %x20
```

rules:

1. `version` emit verbatim, always. first whitespace-delimited token = full SemVer string.
2. if `commitSHA` non-empty, append ` @ ` (space, at-sign, space) + first `min(7, len(commitSHA))` chars. no truncation indicator.
3. parenthesized `details` group emit iff at least one of `buildDate`, `buildInfo` non-empty. exactly one leading space before `(`. group:
   - both → `(<buildDate>, <buildInfo>)` (comma-space separator);
   - only `buildDate` → `(<buildDate>)`;
   - only `buildInfo` → `(<buildInfo>)`.
4. no trailing whitespace.
5. no escaping any field.

### 4.1 examples

| version     | commitSHA       | buildDate              | buildInfo | output                                          |
|-------------|-----------------|------------------------|-----------|-------------------------------------------------|
| `1.2.3`     | `""`            | `""`                   | `""`      | `1.2.3`                                         |
| `1.2.3-rc1` | `""`            | `""`                   | `""`      | `1.2.3-rc1`                                     |
| `1.2.3`     | `abc1234567890` | `""`                   | `""`      | `1.2.3 @ abc1234`                               |
| `1.2.3`     | `abc`           | `""`                   | `""`      | `1.2.3 @ abc`                                   |
| `1.2.3`     | `abcdefg`       | `""`                   | `""`      | `1.2.3 @ abcdefg`                               |
| `1.2.3`     | `""`            | `2026-01-01T00:00:00Z` | `""`      | `1.2.3 (2026-01-01T00:00:00Z)`                  |
| `1.2.3`     | `""`            | `""`                   | `dirty`   | `1.2.3 (dirty)`                                 |
| `1.2.3`     | `abc1234567890` | `2026-01-01T00:00:00Z` | `""`      | `1.2.3 @ abc1234 (2026-01-01T00:00:00Z)`        |
| `1.2.3`     | `abc1234567890` | `""`                   | `dirty`   | `1.2.3 @ abc1234 (dirty)`                       |
| `1.2.3`     | `abc1234567890` | `2026-01-01T00:00:00Z` | `dirty`   | `1.2.3 @ abc1234 (2026-01-01T00:00:00Z, dirty)` |

normative; conforming impl MUST pass all rows.

### 4.2 SemVer extraction

per rule 1, SemVer = leading whitespace-delimited token. no spaces (per SemVer 2.0.0 grammar), so field-split recovers:

```sh
foo --version | awk '{print $NF}' RS=' '          # if output is just the format
foo --version | awk '{print $2}'                  # if output is "<prog> <format>"
foo --version | sed -E 's/^([^ ]*).*/\1/'         # POSIX sed, no prog prefix
foo --version | sed -E 's/^[^ ]+ ([^ ]+).*/\1/'   # POSIX sed, with prog prefix
foo --version | cut -d' ' -f1                     # no prog prefix
foo --version | cut -d' ' -f2                     # with prog prefix
```

result = valid SemVer 2.0.0, pass to any SemVer-aware tool.

## 5. VERSION file

plain UTF-8 text at project root.

- lines start `#` = comments, MUST skip.
- blank lines (after trim ASCII whitespace) MUST skip.
- first remaining line, trimmed, = version string.
- MUST match [SemVer 2.0.0](https://semver.org). generator SHOULD validate via [official regex](https://semver.org/#is-there-a-suggested-regular-expression-regex-to-check-a-semver-string) (or equiv) and MUST fail nonzero exit + diagnostic on mismatch.
- if no such line, generator MUST fail nonzero exit + diagnostic.
- subsequent non-comment lines ignored (reserved; SHOULD NOT rely on).

### 5.1 project root discovery

generator starts at cwd, walks upward, picks nearest ancestor (inclusive) with regular file `VERSION`. if none, MUST fail.

## 6. input fields

regardless of wiring path (§7), the four §3 inputs have semantics:

| field       | meaning + canonical source                                                          |
|-------------|-------------------------------------------------------------------------------------|
| `version`   | project version per `VERSION` file (§5), or toolchain-reported module version where available. |
| `commitSHA` | full git commit SHA, hex. empty if unavailable.                                     |
| `buildDate` | build timestamp, RFC 3339 / ISO-8601 UTC (`YYYY-MM-DDTHH:MM:SSZ`). empty if unavailable. |
| `buildInfo` | `"dirty"` iff working tree had uncommitted changes at build time; else `""`.        |

empty-string semantics: when source unavailable (no git, VCS stamping disabled, etc.) field MUST be `""`, not synthesised. impls MUST NOT substitute placeholders (`unknown`, `n/a`, wall clock when VCS time absent, etc.).

## 7. wiring

impls MUST provide at least one mechanism for populating §6 fields into host program. two paths recognised; runtime sourcing preferred where toolchain supports.

### 7.1 runtime sourcing (preferred where supported)

impl SHOULD expose runtime reader deriving §6 fields from toolchain-baked metadata (e.g. Go `runtime/debug.ReadBuildInfo()` exposing `Main.Version` + `vcs.revision` / `vcs.time` / `vcs.modified`).

advantages over §7.2: works with package-manager fetch-and-build (e.g. `go run module@version`, `go install module@version`) where gitignored generated file absent from tagged tree; no generation step; values cannot go stale vs binary.

requirements:

1. reader MUST accept caller-provided `fallbackVersion` (or equiv), used as `version` when toolchain reports no meaningful module version (local untagged build, `(devel)`, etc.). caller typically passes trimmed `VERSION` file embedded into binary.
2. when toolchain VCS metadata unavailable (e.g. Go `-buildvcs=false`, non-VCS build), `commitSHA` / `buildDate` / `buildInfo` MUST be empty per §6, not synthesised.
3. reader output MUST be byte-identical to what generator (§7.2) would produce given same `(version, commitSHA, buildDate, buildInfo)` tuple via `FormatVersion`.

Go reference: `version.Read(fallbackVersion string) Info`.

### 7.2 generator (fallback / compile-time symbol path)

impl MAY expose generator collecting §6 fields at build time, emit host-language source declaring them as package-level constants/variables.

use when:

- toolchain provides no runtime build metadata (no §7.1 available);
- consumers need compile-time symbols (typo on `Version` ref = build error, not runtime empty string);
- builds run with VCS stamping disabled (e.g. Go `-buildvcs=false`, out-of-tree tarball builds) and `commitSHA` / `buildDate` must still populate from generation-time git state.

generator-specific source mapping (only when no toolchain metadata):

| field       | source                                                                              |
|-------------|-------------------------------------------------------------------------------------|
| `version`   | `VERSION` file (§5).                                                                |
| `commitSHA` | `git rev-parse HEAD`, trimmed. full SHA. empty if unavailable.                       |
| `buildDate` | current UTC time, RFC 3339 / ISO-8601. empty if unavailable.              |
| `buildInfo` | `"dirty"` iff `git status --porcelain` non-empty output; else `""`.                 |

impl MAY scope dirty check to subdir (e.g. `git status --porcelain -- .`); whole-repo + scoped both permitted.

git invoke failure (no git, not repo, etc.) impl-defined: generator MAY error, or MAY emit empty strings. Go impl errors; Rust impl emits empty strings.

emitted artifact exposes four package-level values, names + string types:

```
Version    string  // §6 version
CommitSHA  string  // §6 commitSHA (full SHA, NOT truncated)
BuildDate  string  // §6 buildDate
BuildInfo  string  // §6 buildInfo
```

`PascalCase` for case-sensitive export langs (Go); idiomatic equivalents permitted elsewhere (e.g. Rust env-var bridges: `VERSION_COMMIT_SHA`, `VERSION_BUILD_DATE`, `VERSION_BUILD_INFO`; `version` from build system package metadata).

artifact MUST carry "do not edit" marker for host language. SHOULD be excluded from VCS (`.gitignore`) — but for toolchains distributing via tagged source fetch (Go modules), gitignoring artifact breaks consumers; prefer §7.1.

### 7.3 coexistence

impls MAY provide both §7.1 and §7.2; presence of one MUST NOT preclude other. callers choose per-project.

## 8. determinism

for fixed `(version, commitSHA, buildDate, buildInfo)` tuple, `FormatVersion` pure: same inputs → same output bytes, every platform, every language. SHA truncation by code-unit count over input string; for hex SHAs unambiguous.

## 9. conformance

impl conforms iff:

1. exposes `FormatVersion`-equivalent function with §3 signature + §4 behavior.
2. provides at least one wiring path per §7; each provided path follows its subsection.
3. passes every §4.1 row byte-for-byte.