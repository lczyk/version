# version

[![lint_and_test](https://github.com/lczyk/version/actions/workflows/lint_and_test.yml/badge.svg)](https://github.com/lczyk/version/actions/workflows/lint_and_test.yml)

Tiny build-time generated version-string format + per-language helpers to wire it into a binary.

Output:

```
<semver> @ <commit7> (<date>, <buildinfo>)
```

Example: `0.7.0 @ 5f2fc35 (2026-04-25T10:01:45Z, dirty)`

Trailing parts drop when empty. The `<semver>` prefix is a strict [SemVer 2.0.0](https://semver.org) string and is the first whitespace-delimited token of the output, so it is trivial to extract:

```sh
foo --version | awk '{print $2}'                  # "<prog> <semver> @ ..." → <semver>
foo --version | cut -d' ' -f2                     # same
foo --version | sed -E 's/^[^ ]+ ([^ ]+).*/\1/'   # same, POSIX sed
```

Full grammar, conformance rules, edge cases, and more extraction snippets in [SPEC.md](SPEC.md) (see [§4.2](SPEC.md#42-semver-extraction)).

## Implementations

| language | dir            | helper                                       |
|----------|----------------|----------------------------------------------|
| Go       | [go/](go/)         | `FormatVersion` + `generate-version` codegen |
| Rust     | [rust/](rust/)     | `format_version` + `version!()` macro + `emit()` build-script helper |
| Python   | [python/](python/) | `format_version` + `generate-version` codegen |

All produce byte-identical output for the same inputs (per [SPEC §8](SPEC.md)).
