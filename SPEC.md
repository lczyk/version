# lczyk's version format spec

## 1. scope

defines:

1. **version string format** from four input parts.
2. **VERSION file format** read by build-time generators.
3. **build-time inputs** generator MUST collect.
4. **generated artifact contract** (four named values exposed to host program).

`<version>` field MUST be valid [SemVer 2.0.0](https://semver.org) string. formatter no validate; caller job. output = *display* string, not SemVer string — SemVer prefix designed trivially extractable (see §4.2).

not defined: SemVer grammar (see linked spec), git internals, distribution mechanics.

## 2. terminology

- MUST / SHOULD / MAY per RFC 2119.
- "empty" = string length zero. whitespace-only NOT empty here; trim = caller job.

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

1. `version` field emit verbatim, always. first whitespace-delimited token of output = full SemVer string.
2. if `commitSHA` non-empty, append ` @ ` (space, at-sign, space) plus first `min(7, len(commitSHA))` chars of `commitSHA`. no truncation indicator.
3. parenthesized `details` group emit iff at least one of `buildDate`, `buildInfo` non-empty. exactly one leading space before `(`. group contains:
   - both fields → `(<buildDate>, <buildInfo>)` (comma-space separator);
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

normative cases; conforming impl MUST pass all rows.

### 4.2 SemVer extraction

by rule 1, SemVer string = leading whitespace-delimited token of output. contains no spaces (per SemVer 2.0.0 grammar), so trivial field-split recovers:

```sh
foo --version | awk '{print $NF}' RS=' '          # if output is just the format
foo --version | awk '{print $2}'                  # if output is "<prog> <format>"
foo --version | sed -E 's/^([^ ]*).*/\1/'         # POSIX sed, no prog prefix
foo --version | sed -E 's/^[^ ]+ ([^ ]+).*/\1/'   # POSIX sed, with prog prefix
foo --version | cut -d' ' -f1                     # no prog prefix
foo --version | cut -d' ' -f2                     # with prog prefix
```

result = valid SemVer 2.0.0 string, pass to any SemVer-aware tool.

## 5. VERSION file

plain UTF-8 text at project root.

- lines start `#` = comments, MUST skip.
- blank lines (after trim surrounding ASCII whitespace) MUST skip.
- first remaining line, trimmed, = version string.
- version string MUST match [SemVer 2.0.0](https://semver.org). generator SHOULD validate via [official regex](https://semver.org/#is-there-a-suggested-regular-expression-regex-to-check-a-semver-string) (or equivalent) and MUST fail nonzero exit + diagnostic on mismatch.
- if no such line, generator MUST fail nonzero exit + diagnostic.
- subsequent non-comment lines ignored (reserved future use; SHOULD NOT rely on).

### 5.1 project root discovery

generator starts at cwd, walks upward, picks nearest ancestor (inclusive) with regular file `VERSION`. if none, generator MUST fail.

## 6. build-time inputs

generator MUST collect:

| field       | source                                                                              |
|-------------|-------------------------------------------------------------------------------------|
| `version`   | `VERSION` file (§5).                                                                |
| `commitSHA` | `git rev-parse HEAD`, trimmed. full SHA. empty if unavailable.                       |
| `buildDate` | current UTC time, formatted RFC 3339 / ISO-8601 (`YYYY-MM-DDTHH:MM:SSZ`). empty if unavailable. |
| `buildInfo` | `"dirty"` iff `git status --porcelain` non-empty output; else `""`.    |

impl MAY scope dirty check to subdir (e.g. `git status --porcelain -- .`) when ergonomic; whole-repo + scoped both permitted.

git invoke failure (no git, not repo, etc.) impl-defined: generator MAY error, or MAY emit empty strings. Go impl errors; Rust impl emits empty strings.

## 7. generated artifact

generator emits source file in host language exposing four package-level values with these names and string types:

```
Version    string  // §6 version
CommitSHA  string  // §6 commitSHA (full SHA, NOT truncated)
BuildDate  string  // §6 buildDate
BuildInfo  string  // §6 buildInfo
```

naming `PascalCase` for case-sensitive export langs (Go); idiomatic equivalents permitted elsewhere (e.g. env-var bridges in Rust: `VERSION_COMMIT_SHA`, `VERSION_BUILD_DATE`, `VERSION_BUILD_INFO`, `version` from build system package metadata).

artifact MUST carry "do not edit" marker for language. artifact SHOULD be excluded from version control (`.gitignore`).

## 8. determinism

for fixed `(version, commitSHA, buildDate, buildInfo)` tuple, `FormatVersion` pure function: same inputs → same output bytes, every platform, every language. SHA truncation by code-unit count over input string; for hex SHAs unambiguous.

## 9. conformance

impl conforms iff:

1. exposes `FormatVersion`-equivalent function with §3 signature + §4 behavior.
2. generator (if provided) follows §5–§7.
3. passes every §4.1 row byte-for-byte.
