# version (Rust)

Format a version string from its parts. Conforms to [SPEC.md](../SPEC.md).

Single crate, no workspace. Used as both a runtime dep (formatter + macro) and a build-script dep (env-var emitter).

## format_version

```rust
let s = version::format_version("1.0.0", "abc1234567890", "", "dirty");
// "1.0.0 @ abc1234 (dirty)"
```

Output format and edge cases: see [SPEC §4](../SPEC.md).

## version!() macro

Compile-time string from `CARGO_PKG_VERSION` plus env vars set by [`emit`] in `build.rs`.

```rust
println!("mytool {}", version::version!());
// "mytool 0.7.0 @ 5f2fc35 (2026-04-25T10:01:45Z, dirty)"
```

## emit (build-script helper)

`Cargo.toml`:

```toml
[dependencies]
version = { git = "https://github.com/lczyk/version", branch = "main", subdir = "rust" }

[build-dependencies]
version = { git = "https://github.com/lczyk/version", branch = "main", subdir = "rust" }
```

`build.rs`:

```rust
fn main() {
    version::emit();
}
```

`emit()` sets these env vars consumed by `version!()`:

| env var               | source                                        |
|-----------------------|-----------------------------------------------|
| `VERSION_COMMIT_SHA`  | `git rev-parse HEAD` (full SHA, or `""`)      |
| `VERSION_BUILD_DATE`  | `date -u +%Y-%m-%dT%H:%M:%SZ` (or `""`)       |
| `VERSION_BUILD_INFO`  | `"dirty"` if `git status --porcelain` non-empty, else `""` |

Per [SPEC §6](../SPEC.md), git failure emits empty strings (no error).

## Development

```
make test     # cargo test
make lint     # cargo clippy -- -D warnings + cargo fmt --check
make format   # cargo fmt
```
