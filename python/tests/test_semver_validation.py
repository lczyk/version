import pytest
import tempfile
from pathlib import Path

from version.generate import read_version


def test_valid_semver_versions() -> None:
    """test valid SemVer 2.0.0 strings are accepted."""
    valid = [
        "1.0.0",
        "0.0.0",
        "1.2.3",
        "10.20.30",
        "1.0.0-alpha",
        "1.0.0-alpha.1",
        "1.0.0-0.3.7",
        "1.0.0-x.7.z.92",
        "1.0.0+20130313144700",
        "1.0.0-beta+exp.sha.5114f85",
        "1.0.0+21AF26D3-117B344092BD",
        "2.0.0-rc.1+build.123",
    ]
    for v in valid:
        with tempfile.NamedTemporaryFile(mode="w", suffix=".txt", delete=False) as f:
            f.write(v)
            f.flush()
            try:
                result = read_version(Path(f.name))
                assert result == v, f"expected {v}, got {result}"
            finally:
                Path(f.name).unlink()


def test_invalid_semver_versions() -> None:
    """test invalid SemVer 2.0.0 strings are rejected."""
    invalid = [
        "1",
        "1.2",
        "v1.0.0",
        "1.0.0.0",
        "01.0.0",
        "1.02.0",
        "1.0.00",
        "1.0.0-",
        "1.0.0-+",
        "1.0.0-+build",
        "1.0.0-alpha..1",
        "1.0.0-alpha.+build",
        "latest",
        "main",
    ]
    for v in invalid:
        with tempfile.NamedTemporaryFile(mode="w", suffix=".txt", delete=False) as f:
            f.write(v)
            f.flush()
            try:
                with pytest.raises(
                    ValueError, match="not a valid SemVer 2.0.0 string"
                ):
                    read_version(Path(f.name))
            finally:
                Path(f.name).unlink()


def test_comments_and_blanks_skipped() -> None:
    """test comments and blank lines are skipped, first valid line is used."""
    with tempfile.NamedTemporaryFile(mode="w", suffix=".txt", delete=False) as f:
        f.write("# comment\n\n1.2.3\ninvalid-line\n")
        f.flush()
        try:
            result = read_version(Path(f.name))
            assert result == "1.2.3"
        finally:
            Path(f.name).unlink()


def test_whitespace_trimmed() -> None:
    """test whitespace is trimmed from version line."""
    with tempfile.NamedTemporaryFile(mode="w", suffix=".txt", delete=False) as f:
        f.write("  1.2.3  \n")
        f.flush()
        try:
            result = read_version(Path(f.name))
            assert result == "1.2.3"
        finally:
            Path(f.name).unlink()
