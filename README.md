# version

[![lint_and_test](https://github.com/lczyk/version/actions/workflows/lint_and_test.yml/badge.svg)](https://github.com/lczyk/version/actions/workflows/lint_and_test.yml)

Tiny build-time generated version-string format + per-language helpers to wire it into a binary.

Output:

```
<version>+<commit7> (<date>, <buildinfo>)
```

Example: `0.7.0+5f2fc35 (2026-04-25T10:01:45Z, dirty)`

Trailing parts drop when empty. Full grammar, conformance rules, and edge cases in [SPEC.md](SPEC.md).

## Implementations

| language | dir            | helper                                       |
|----------|----------------|----------------------------------------------|
| Go       | [go/](go/)         | `FormatVersion` + `generate-version` codegen |
| Rust     | [rust/](rust/)     | `format_version` + `version!()` macro + `emit()` build-script helper |
| Python   | [python/](python/) | `format_version` + `generate-version` codegen |

All produce byte-identical output for the same inputs (per [SPEC §8](SPEC.md)).
