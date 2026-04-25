import pytest

from version import format_version


@pytest.mark.parametrize(
    ("version", "commit_sha", "build_date", "build_info", "expected"),
    [
        ("1.0.0", "", "", "", "1.0.0"),
        ("1.0.0", "abc1234567890", "", "", "1.0.0 @ abc1234"),
        (
            "1.0.0",
            "abc1234567890",
            "2026-01-01T00:00:00Z",
            "",
            "1.0.0 @ abc1234 (2026-01-01T00:00:00Z)",
        ),
        (
            "1.0.0",
            "abc1234567890",
            "2026-01-01T00:00:00Z",
            "dirty",
            "1.0.0 @ abc1234 (2026-01-01T00:00:00Z, dirty)",
        ),
        ("1.0.0", "abc1234567890", "", "dirty", "1.0.0 @ abc1234 (dirty)"),
        ("1.0.0", "abc", "", "", "1.0.0 @ abc"),
        ("1.0.0", "abcdefg", "", "", "1.0.0 @ abcdefg"),
        ("1.0.0", "", "2026-01-01T00:00:00Z", "", "1.0.0 (2026-01-01T00:00:00Z)"),
        ("1.0.0", "", "", "dirty", "1.0.0 (dirty)"),
        ("1.2.3-rc1", "", "", "", "1.2.3-rc1"),
        ("1.0.0", "abcdefgh", "", "", "1.0.0 @ abcdefg"),
        ("1.0.0-alpha.1", "abc1234567890", "", "", "1.0.0-alpha.1 @ abc1234"),
        ("1.0.0-rc.1", "abc1234567890", "2026-01-01T00:00:00Z", "dirty", "1.0.0-rc.1 @ abc1234 (2026-01-01T00:00:00Z, dirty)"),
        ("1.0.0", "abc1234567890", "", "dirty, uncommitted", "1.0.0 @ abc1234 (dirty, uncommitted)"),
        ("1.0.0", "a", "", "", "1.0.0 @ a"),
    ],
)
def test_format_version(
    version: str, commit_sha: str, build_date: str, build_info: str, expected: str
) -> None:
    assert format_version(version, commit_sha, build_date, build_info) == expected
