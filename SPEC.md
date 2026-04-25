# lczyk's version format spec

## 1. scope

Defines:

1. **version string format** from four input parts.
2. **VERSION file format** read by build-time generators.
3. **build-time inputs** generator MUST collect.
4. **generated artifact contract** (four named values exposed to host program).

Not defined: SemVer rules for `<version>` field, git internals, distribution mechanics.

## 2. terminology

- MUST / SHOULD / MAY per RFC 2119.
- "Empty" = string length zero. Whitespace-only NOT empty here; trimming caller's job.

## 3. inputs

`FormatVersion` takes four ordered string params:

| name        | meaning                                                  |
|-------------|----------------------------------------------------------|
| `version`   | Project version (e.g. `1.2.3`). Opaque to this spec.     |
| `commitSHA` | Full git commit SHA, hex. May be empty.                  |
| `buildDate` | Build timestamp, ISO-8601 / RFC 3339 UTC. May be empty.  |
| `buildInfo` | Free-form build annotation (e.g. `dirty`). May be empty. |

All four strings. Formatter does no further validation.

## 4. output format

Grammar (ABNF-ish):

```
output     = version [ "+" commit7 ] [ SP "(" details ")" ]
details    = build-date [ ", " build-info ]
           / build-date
           / build-info
commit7    = 1*7(HEXDIG)            ; first min(7, len(commitSHA)) chars of commitSHA
SP         = %x20
```

Rules:

1. `version` field emitted verbatim, always.
2. If `commitSHA` non-empty, append `+` plus first `min(7, len(commitSHA))` chars of `commitSHA`. No truncation indicator.
3. Parenthesized `details` group emitted iff at least one of `buildDate`, `buildInfo` non-empty. Exactly one leading space before `(`. Group contains:
   - both fields → `(<buildDate>, <buildInfo>)` (comma-space separator);
   - only `buildDate` → `(<buildDate>)`;
   - only `buildInfo` → `(<buildInfo>)`.
4. No trailing whitespace.
5. No escaping on any field.

### 4.1 examples

| version | commitSHA       | buildDate              | buildInfo | output                                       |
|---------|-----------------|------------------------|-----------|----------------------------------------------|
| `1.2.3` | `""`            | `""`                   | `""`      | `1.2.3`                                      |
| `1.2.3` | `abc1234567890` | `""`                   | `""`      | `1.2.3+abc1234`                              |
| `1.2.3` | `abc`           | `""`                   | `""`      | `1.2.3+abc`                                  |
| `1.2.3` | `abcdefg`       | `""`                   | `""`      | `1.2.3+abcdefg`                              |
| `1.2.3` | `abc1234567890` | `2026-01-01T00:00:00Z` | `""`      | `1.2.3+abc1234 (2026-01-01T00:00:00Z)`       |
| `1.2.3` | `abc1234567890` | `""`                   | `dirty`   | `1.2.3+abc1234 (dirty)`                      |
| `1.2.3` | `abc1234567890` | `2026-01-01T00:00:00Z` | `dirty`   | `1.2.3+abc1234 (2026-01-01T00:00:00Z, dirty)`|

Normative cases; conforming impl MUST pass all rows.

## 5. VERSION file

Plain UTF-8 text at project root.

- Lines starting `#` = comments, MUST skip.
- Blank lines (after trim surrounding ASCII whitespace) MUST skip.
- First remaining line, trimmed, = version string.
- If no such line, generator MUST fail nonzero exit + diagnostic.
- Subsequent non-comment lines ignored (reserved future use; SHOULD NOT rely on).

### 5.1 project root discovery

Generator starts at cwd, walks upward, picks nearest ancestor (inclusive) with regular file `VERSION`. If none, generator MUST fail.

## 6. build-time inputs

Generator MUST collect:

| field       | source                                                                              |
|-------------|-------------------------------------------------------------------------------------|
| `version`   | `VERSION` file (§5).                                                                |
| `commitSHA` | `git rev-parse HEAD`, trimmed. Full SHA. Empty if unavailable.                       |
| `buildDate` | Current UTC time, formatted RFC 3339 / ISO-8601 (`YYYY-MM-DDTHH:MM:SSZ`). Empty if unavailable. |
| `buildInfo` | `"dirty"` iff `git status --porcelain` non-empty output; else `""`.    |

Impl MAY scope dirty check to subdir (e.g. `git status --porcelain -- .`) when ergonomic; whole-repo + scoped both permitted.

Git invoke failure (no git, not repo, etc.) implementation-defined: generator MAY error, or MAY emit empty strings. Go impl errors; Rust impl emits empty strings.

## 7. generated artifact

Generator emits source file in host language exposing four package-level values with these names and string types:

```
Version    string  // §6 version
CommitSHA  string  // §6 commitSHA (full SHA, NOT truncated)
BuildDate  string  // §6 buildDate
BuildInfo  string  // §6 buildInfo
```

Naming `PascalCase` for case-sensitive export langs (Go); idiomatic equivalents permitted elsewhere (e.g. env-var bridges in Rust: `VERSION_COMMIT_SHA`, `VERSION_BUILD_DATE`, `VERSION_BUILD_INFO`, `version` from build system package metadata).

Artifact MUST carry "do not edit" marker for the language. Artifact SHOULD be excluded from version control (`.gitignore`).

## 8. determinism

For fixed `(version, commitSHA, buildDate, buildInfo)` tuple, `FormatVersion` pure function: same inputs → same output bytes, every platform, every language. SHA truncation by code-unit count over input string; for hex SHAs unambiguous.

## 9. conformance

Impl conforms iff:

1. Exposes `FormatVersion`-equivalent function with §3 signature + §4 behavior.
2. Generator (if provided) follows §5–§7.
3. Passes every §4.1 row byte-for-byte.