//! Format a version string from its parts.
//!
//! Output: `<semver> @ <commit7> (<date>, <buildinfo>)`. Trailing parts omitted when empty.
//! `<semver>` is the caller's version field, which MUST be a valid SemVer 2.0.0 string
//! (<https://semver.org>). It is emitted verbatim, with no separator characters following it,
//! so it can be extracted with `awk '{print $1}'` or `cut -d' ' -f1`.
//!
//! # Runtime use
//!
//! ```
//! let s = version::format_version("1.0.0", "abc1234567890", "", "dirty");
//! assert_eq!(s, "1.0.0 @ abc1234 (dirty)");
//! ```
//!
//! Or via macro, which reads `CARGO_PKG_VERSION` plus env vars set by [`emit`]:
//!
//! ```ignore
//! println!("mytool {}", version::version!());
//! ```
//!
//! # Build-script use
//!
//! Add this crate as a `build-dependency` and call [`emit`] from `build.rs`:
//!
//! ```no_run
//! // build.rs
//! version::emit();
//! ```

use std::process::Command;

pub fn format_version(
    version: &str,
    commit_sha: &str,
    build_date: &str,
    build_info: &str,
) -> String {
    let mut result = String::from(version);

    if !commit_sha.is_empty() {
        let n = commit_sha.len().min(7);
        result.push_str(" @ ");
        result.push_str(&commit_sha[..n]);
    }

    match (!build_date.is_empty(), !build_info.is_empty()) {
        (true, true) => {
            result.push_str(" (");
            result.push_str(build_date);
            result.push_str(", ");
            result.push_str(build_info);
            result.push(')');
        }
        (true, false) => {
            result.push_str(" (");
            result.push_str(build_date);
            result.push(')');
        }
        (false, true) => {
            result.push_str(" (");
            result.push_str(build_info);
            result.push(')');
        }
        (false, false) => {}
    }

    result
}

/// Compile-time version string from `CARGO_PKG_VERSION` plus env vars set by [`emit`]:
/// `VERSION_COMMIT_SHA`, `VERSION_BUILD_DATE`, `VERSION_BUILD_INFO`.
#[macro_export]
macro_rules! version {
    () => {
        $crate::format_version(
            env!("CARGO_PKG_VERSION"),
            option_env!("VERSION_COMMIT_SHA").unwrap_or(""),
            option_env!("VERSION_BUILD_DATE").unwrap_or(""),
            option_env!("VERSION_BUILD_INFO").unwrap_or(""),
        )
    };
}

fn cmd_output(program: &str, args: &[&str], dir: Option<&str>) -> Option<String> {
    let mut c = Command::new(program);
    c.args(args);
    if let Some(d) = dir {
        c.current_dir(d);
    }
    let out = c.output().ok()?;
    if !out.status.success() {
        return None;
    }
    Some(String::from_utf8_lossy(&out.stdout).trim().to_string())
}

/// Emit the env vars consumed by [`version!`]. Call from `build.rs`.
///
/// Sets:
/// - `VERSION_COMMIT_SHA` — full git SHA, or `""` if unavailable
/// - `VERSION_BUILD_DATE` — UTC ISO-8601 from `date -u`, or `""` if unavailable
/// - `VERSION_BUILD_INFO` — `"dirty"` if working tree dirty, else `""`
pub fn emit() {
    let manifest_dir = std::env::var("CARGO_MANIFEST_DIR").unwrap_or_else(|_| ".".into());

    let sha = cmd_output("git", &["rev-parse", "HEAD"], Some(&manifest_dir))
        .filter(|s| !s.is_empty())
        .unwrap_or_default();

    let info = match cmd_output(
        "git",
        &["status", "--porcelain", "--", "."],
        Some(&manifest_dir),
    ) {
        Some(s) if !s.is_empty() => "dirty".to_string(),
        _ => String::new(),
    };

    let date = cmd_output("date", &["-u", "+%Y-%m-%dT%H:%M:%SZ"], None)
        .filter(|s| !s.is_empty())
        .unwrap_or_default();

    println!("cargo:rustc-env=VERSION_COMMIT_SHA={sha}");
    println!("cargo:rustc-env=VERSION_BUILD_DATE={date}");
    println!("cargo:rustc-env=VERSION_BUILD_INFO={info}");

    println!("cargo:rerun-if-changed=build.rs");
    println!("cargo:rerun-if-changed=.git/HEAD");
    println!("cargo:rerun-if-changed=.git/index");
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn version_only() {
        assert_eq!(format_version("1.0.0", "", "", ""), "1.0.0");
    }

    #[test]
    fn version_with_commit() {
        assert_eq!(
            format_version("1.0.0", "abc1234567890", "", ""),
            "1.0.0 @ abc1234"
        );
    }

    #[test]
    fn version_with_commit_and_date() {
        assert_eq!(
            format_version("1.0.0", "abc1234567890", "2026-01-01T00:00:00Z", ""),
            "1.0.0 @ abc1234 (2026-01-01T00:00:00Z)"
        );
    }

    #[test]
    fn version_with_commit_date_dirty() {
        assert_eq!(
            format_version("1.0.0", "abc1234567890", "2026-01-01T00:00:00Z", "dirty"),
            "1.0.0 @ abc1234 (2026-01-01T00:00:00Z, dirty)"
        );
    }

    #[test]
    fn dirty_without_date() {
        assert_eq!(
            format_version("1.0.0", "abc1234567890", "", "dirty"),
            "1.0.0 @ abc1234 (dirty)"
        );
    }

    #[test]
    fn short_commit_sha() {
        assert_eq!(format_version("1.0.0", "abc", "", ""), "1.0.0 @ abc");
    }

    #[test]
    fn seven_char_commit_sha() {
        assert_eq!(
            format_version("1.0.0", "abcdefg", "", ""),
            "1.0.0 @ abcdefg"
        );
    }

    #[test]
    fn date_without_commit() {
        assert_eq!(
            format_version("1.0.0", "", "2026-01-01T00:00:00Z", ""),
            "1.0.0 (2026-01-01T00:00:00Z)"
        );
    }

    #[test]
    fn dirty_without_commit_or_date() {
        assert_eq!(format_version("1.0.0", "", "", "dirty"), "1.0.0 (dirty)");
    }

    #[test]
    fn prerelease_semver_passes_through() {
        assert_eq!(format_version("1.2.3-rc1", "", "", ""), "1.2.3-rc1");
    }

    #[test]
    fn eight_char_commit_truncates_to_seven() {
        assert_eq!(
            format_version("1.0.0", "abcdefgh", "", ""),
            "1.0.0 @ abcdefg"
        );
    }

    #[test]
    fn prerelease_with_commit() {
        assert_eq!(
            format_version("1.0.0-alpha.1", "abc1234567890", "", ""),
            "1.0.0-alpha.1 @ abc1234"
        );
    }

    #[test]
    fn prerelease_with_full_metadata() {
        assert_eq!(
            format_version("1.0.0-rc.1", "abc1234567890", "2026-01-01T00:00:00Z", "dirty"),
            "1.0.0-rc.1 @ abc1234 (2026-01-01T00:00:00Z, dirty)"
        );
    }

    #[test]
    fn buildinfo_with_multiple_words() {
        assert_eq!(
            format_version("1.0.0", "abc1234567890", "", "dirty, uncommitted"),
            "1.0.0 @ abc1234 (dirty, uncommitted)"
        );
    }

    #[test]
    fn one_char_commit() {
        assert_eq!(format_version("1.0.0", "a", "", ""), "1.0.0 @ a");
    }
}
