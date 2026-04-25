def format_version(
    version: str, commit_sha: str, build_date: str, build_info: str
) -> str:
    result = version

    if commit_sha:
        result += "+" + commit_sha[: min(7, len(commit_sha))]

    if build_date and build_info:
        result += f" ({build_date}, {build_info})"
    elif build_date:
        result += f" ({build_date})"
    elif build_info:
        result += f" ({build_info})"

    return result
