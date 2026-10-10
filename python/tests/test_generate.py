import os
import subprocess
from pathlib import Path

import pytest

from version.generate import main


def _git(cwd: Path, *args: str, env: dict[str, str] | None = None) -> None:
    subprocess.run(
        [
            "git",
            "-c",
            "user.name=t",
            "-c",
            "user.email=t@example.invalid",
            "-c",
            "commit.gpgsign=false",
            *args,
        ],
        cwd=cwd,
        check=True,
        capture_output=True,
        env={**os.environ, **(env or {})},
    )


def test_stable_across_regeneration(
    tmp_path: Path, monkeypatch: pytest.MonkeyPatch
) -> None:
    (tmp_path / "VERSION").write_text("1.2.3\n")
    # a tracked placeholder, the setup that used to report "dirty" forever
    out = tmp_path / "_version.py"
    out.write_text("version = ''\n")
    _git(tmp_path, "init", "-q")
    _git(tmp_path, "add", "VERSION", "_version.py")
    _git(
        tmp_path,
        "commit",
        "-q",
        "-m",
        "chore: init",
        env={"GIT_COMMITTER_DATE": "2026-01-02T03:04:05+02:00"},
    )
    monkeypatch.chdir(tmp_path)

    assert main(["--out", "_version.py"]) == 0
    first = out.read_text()
    assert "version = '1.2.3'" in first
    assert "build_date = '2026-01-02T01:04:05Z'" in first
    assert "build_info = ''" in first

    os.utime(out, (1, 1))
    assert main(["--out", "_version.py"]) == 0
    assert out.read_text() == first
    assert out.stat().st_mtime == 1

    (tmp_path / "VERSION").write_text("1.2.3\n# edited\n")
    assert main(["--out", "_version.py"]) == 0
    assert "build_info = 'dirty'" in out.read_text()
