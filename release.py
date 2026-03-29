#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
GitHub Release 自动发布工具

功能：
  1. 自动从 version.py 读取版本号
  2. 将指定 skill 目录打成 zip 包（自动排除 __pycache__、*.log、*.zip、.git 等）
  3. 通过 GitHub REST API 创建 Release 并上传 zip 附件
  4. 支持从 changelog.md 自动提取当前版本的 Release Notes

使用方式：
  # 基本用法（需要设置 GITHUB_TOKEN 环境变量）
  python release.py --skill dicom-doctor

  # 指定版本号（覆盖 version.py 中的版本）
  python release.py --skill dicom-doctor --version 2.10.0

  # 仅打包不发布（dry-run 模式）
  python release.py --skill dicom-doctor --dry-run

  # 指定仓库（默认从 git remote 自动检测）
  python release.py --skill dicom-doctor --repo eyjian/ai-skills

环境变量：
  GITHUB_TOKEN  - GitHub Personal Access Token（需要 repo 权限）
                  创建地址: https://github.com/settings/tokens

作者：DICOM Doctor 项目组
"""

import argparse
import json
import os
import re
import subprocess
import sys
import tempfile
import zipfile
from datetime import datetime
from pathlib import Path
from urllib.error import HTTPError
from urllib.request import Request, urlopen


# ============================================================
# 常量
# ============================================================

# 打 zip 包时排除的文件/目录模式
EXCLUDE_PATTERNS = [
    "__pycache__",
    "*.pyc",
    "*.pyo",
    "*.log",
    "*.zip",
    ".git",
    ".gitignore",
    ".DS_Store",
    "Thumbs.db",
    "*.egg-info",
    "*.swp",
    "*.swo",
    "*~",
    "input/",       # 测试用的 DICOM 输入文件（通常很大）
    "output/",      # 输出目录
    "docs/",        # 技术文章等文档（非 skill 运行所需）
    "dicom_output_*",
]

GITHUB_API_BASE = "https://api.github.com"
GITHUB_UPLOAD_BASE = "https://uploads.github.com"


# ============================================================
# 工具函数
# ============================================================

def _log(msg: str, level: str = "INFO"):
    """简单日志输出。"""
    colors = {
        "INFO": "\033[36m",    # 青色
        "OK": "\033[32m",      # 绿色
        "WARN": "\033[33m",    # 黄色
        "ERROR": "\033[31m",   # 红色
    }
    reset = "\033[0m"
    color = colors.get(level, "")
    print(f"{color}[{level}]{reset} {msg}")


def _read_version(skill_dir: Path) -> str:
    """从 version.py 读取版本号。"""
    version_file = skill_dir / "scripts" / "version.py"
    if not version_file.exists():
        raise FileNotFoundError(f"找不到版本文件: {version_file}")

    content = version_file.read_text(encoding="utf-8")
    match = re.search(r'__version__\s*=\s*["\']([^"\']+)["\']', content)
    if not match:
        raise ValueError(f"无法从 {version_file} 中解析版本号")

    return match.group(1)


def _extract_changelog(skill_dir: Path, version: str) -> str:
    """从 changelog.md 中提取指定版本的变更说明。"""
    changelog_file = skill_dir / "references" / "changelog.md"
    if not changelog_file.exists():
        _log(f"未找到 changelog 文件: {changelog_file}，将使用默认 Release Notes", "WARN")
        return ""

    content = changelog_file.read_text(encoding="utf-8")

    # 匹配 ## vX.Y.Z 开头的章节，提取到下一个 ## 之前
    # 支持 ## v2.10.0 — 标题 这种格式
    pattern = rf"(## v{re.escape(version)}\b.*?)(?=\n## v\d|\Z)"
    match = re.search(pattern, content, re.DOTALL)

    if match:
        notes = match.group(1).strip()
        _log(f"从 changelog.md 提取到 v{version} 的变更说明（{len(notes)} 字符）")
        return notes

    _log(f"changelog.md 中未找到 v{version} 的条目，将使用默认 Release Notes", "WARN")
    return ""


def _detect_repo_from_git() -> str:
    """从 git remote 自动检测 owner/repo。"""
    try:
        result = subprocess.run(
            ["git", "remote", "get-url", "origin"],
            capture_output=True, text=True, check=True,
        )
        url = result.stdout.strip()

        # 支持 SSH 和 HTTPS 格式
        # git@github.com:eyjian/ai-skills.git
        # https://github.com/eyjian/ai-skills.git
        patterns = [
            r"github\.com[:/]([^/]+/[^/]+?)(?:\.git)?$",
        ]
        for p in patterns:
            m = re.search(p, url)
            if m:
                return m.group(1)

        raise ValueError(f"无法从 remote URL 解析仓库信息: {url}")
    except (subprocess.CalledProcessError, FileNotFoundError):
        raise RuntimeError(
            "无法从 git remote 检测仓库信息。请使用 --repo owner/repo 手动指定。"
        )


def _should_exclude(rel_path: str) -> bool:
    """判断文件是否应被排除。"""
    import fnmatch
    parts = Path(rel_path).parts

    for pattern in EXCLUDE_PATTERNS:
        # 检查路径中的每个部分
        for part in parts:
            if fnmatch.fnmatch(part, pattern):
                return True
        # 检查完整相对路径
        if fnmatch.fnmatch(rel_path, pattern):
            return True
        if fnmatch.fnmatch(Path(rel_path).name, pattern):
            return True

    return False


def _create_zip(skill_dir: Path, output_path: Path, skill_name: str) -> int:
    """
    将 skill 目录打成 zip 包。

    Returns:
        打包的文件数量
    """
    file_count = 0
    total_size = 0

    with zipfile.ZipFile(output_path, "w", zipfile.ZIP_DEFLATED, compresslevel=9) as zf:
        for root, dirs, files in os.walk(skill_dir):
            # 跳过隐藏目录和排除目录
            dirs[:] = [
                d for d in dirs
                if not d.startswith(".") and not _should_exclude(d)
            ]

            for fname in files:
                fpath = Path(root) / fname
                rel_path = fpath.relative_to(skill_dir.parent)

                if _should_exclude(str(rel_path)):
                    continue

                zf.write(fpath, str(rel_path))
                file_count += 1
                total_size += fpath.stat().st_size

    zip_size = output_path.stat().st_size
    _log(
        f"打包完成: {file_count} 个文件, "
        f"原始大小 {total_size / 1024 / 1024:.1f}MB → "
        f"压缩后 {zip_size / 1024 / 1024:.1f}MB "
        f"(压缩率 {(1 - zip_size / max(total_size, 1)) * 100:.0f}%)"
    )
    return file_count


def _github_api(
    method: str,
    url: str,
    token: str,
    data: dict = None,
    content_type: str = "application/json",
    binary_data: bytes = None,
) -> dict:
    """调用 GitHub REST API。"""
    headers = {
        "Authorization": f"token {token}",
        "Accept": "application/vnd.github.v3+json",
    }

    if binary_data is not None:
        body = binary_data
        headers["Content-Type"] = content_type
    elif data is not None:
        body = json.dumps(data).encode("utf-8")
        headers["Content-Type"] = "application/json"
    else:
        body = None

    req = Request(url, data=body, headers=headers, method=method)

    try:
        with urlopen(req) as resp:
            resp_body = resp.read().decode("utf-8")
            if resp_body:
                return json.loads(resp_body)
            return {}
    except HTTPError as e:
        error_body = e.read().decode("utf-8") if e.fp else ""
        _log(f"GitHub API 错误 ({e.code}): {error_body}", "ERROR")
        raise


def _create_release(
    repo: str,
    token: str,
    tag: str,
    name: str,
    body: str,
    draft: bool = False,
    prerelease: bool = False,
) -> dict:
    """创建 GitHub Release。"""
    url = f"{GITHUB_API_BASE}/repos/{repo}/releases"
    data = {
        "tag_name": tag,
        "target_commitish": "main",
        "name": name,
        "body": body,
        "draft": draft,
        "prerelease": prerelease,
    }
    return _github_api("POST", url, token, data=data)


def _upload_asset(
    upload_url: str,
    token: str,
    file_path: Path,
    content_type: str = "application/zip",
) -> dict:
    """上传 Release 附件。"""
    # upload_url 格式: https://uploads.github.com/repos/.../releases/.../assets{?name,label}
    # 需要替换掉模板部分
    upload_url = re.sub(r"\{[^}]*\}", "", upload_url)
    upload_url = f"{upload_url}?name={file_path.name}"

    with open(file_path, "rb") as f:
        binary_data = f.read()

    _log(f"上传附件: {file_path.name} ({len(binary_data) / 1024 / 1024:.1f}MB)...")
    return _github_api(
        "POST", upload_url, token,
        binary_data=binary_data,
        content_type=content_type,
    )


def _check_release_exists(repo: str, token: str, tag: str) -> bool:
    """检查指定 tag 的 Release 是否已存在。"""
    url = f"{GITHUB_API_BASE}/repos/{repo}/releases/tags/{tag}"
    try:
        _github_api("GET", url, token)
        return True
    except HTTPError as e:
        if e.code == 404:
            return False
        raise


# ============================================================
# 主流程
# ============================================================

def main():
    parser = argparse.ArgumentParser(
        description="GitHub Release 自动发布工具 —— 打 zip 包并上传到 GitHub Releases",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
示例:
  # 发布 dicom-doctor（自动读取版本号和 changelog）
  python release.py --skill dicom-doctor

  # 仅打包，不上传（用于本地测试）
  python release.py --skill dicom-doctor --dry-run

  # 发布为草稿（不公开，可在 GitHub 页面确认后再发布）
  python release.py --skill dicom-doctor --draft

  # 指定版本号
  python release.py --skill dicom-doctor --version 2.10.0

环境变量:
  GITHUB_TOKEN  GitHub Personal Access Token（需要 repo 权限）
        """,
    )
    parser.add_argument(
        "--skill", required=True,
        help="Skill 名称（即目录名，如 dicom-doctor）",
    )
    parser.add_argument(
        "--version", default=None,
        help="版本号（默认从 scripts/version.py 自动读取）",
    )
    parser.add_argument(
        "--repo", default=None,
        help="GitHub 仓库（owner/repo 格式，默认从 git remote 自动检测）",
    )
    parser.add_argument(
        "--dry-run", action="store_true",
        help="仅打包，不上传到 GitHub（用于本地测试）",
    )
    parser.add_argument(
        "--draft", action="store_true",
        help="创建为草稿 Release（不公开，可在 GitHub 页面确认后再发布）",
    )
    parser.add_argument(
        "--prerelease", action="store_true",
        help="标记为预发布版本",
    )
    parser.add_argument(
        "--output-dir", default=None,
        help="zip 包输出目录（默认为临时目录；--dry-run 时默认为当前目录）",
    )
    parser.add_argument(
        "--notes", default=None,
        help="自定义 Release Notes（默认从 changelog.md 自动提取）",
    )
    parser.add_argument(
        "--force", action="store_true",
        help="如果同版本 Release 已存在，强制删除后重新创建",
    )

    args = parser.parse_args()

    # ---- 1. 定位 skill 目录 ----
    script_dir = Path(__file__).resolve().parent
    skill_dir = script_dir / args.skill

    if not skill_dir.is_dir():
        _log(f"Skill 目录不存在: {skill_dir}", "ERROR")
        sys.exit(1)

    _log(f"Skill 目录: {skill_dir}")

    # ---- 2. 读取版本号 ----
    if args.version:
        version = args.version
        _log(f"使用指定版本号: v{version}")
    else:
        version = _read_version(skill_dir)
        _log(f"从 version.py 读取版本号: v{version}")

    tag = f"v{version}"
    skill_name = args.skill

    # ---- 3. 提取 Release Notes ----
    if args.notes:
        release_notes = args.notes
    else:
        release_notes = _extract_changelog(skill_dir, version)

    if not release_notes:
        release_notes = f"## {skill_name} {tag}\n\n发布于 {datetime.now().strftime('%Y-%m-%d %H:%M')}"

    # ---- 4. 打 zip 包 ----
    zip_filename = f"{skill_name}-{tag}.zip"

    if args.dry_run and not args.output_dir:
        output_dir = Path(".")
    elif args.output_dir:
        output_dir = Path(args.output_dir)
        output_dir.mkdir(parents=True, exist_ok=True)
    else:
        output_dir = Path(tempfile.mkdtemp(prefix=f"{skill_name}-release-"))

    zip_path = output_dir / zip_filename

    _log(f"开始打包: {skill_dir.name}/ → {zip_path}")
    file_count = _create_zip(skill_dir, zip_path, skill_name)

    if file_count == 0:
        _log("打包结果为空，请检查 skill 目录", "ERROR")
        sys.exit(1)

    # ---- 5. dry-run 模式：到此为止 ----
    if args.dry_run:
        _log(f"[dry-run] zip 包已生成: {zip_path.resolve()}", "OK")
        _log(f"[dry-run] Tag: {tag}")
        _log(f"[dry-run] Release Notes 预览（前 200 字符）:")
        print(release_notes[:200])
        print("...")
        return

    # ---- 6. 获取 GitHub Token ----
    token = os.environ.get("GITHUB_TOKEN")
    if not token:
        _log(
            "未设置 GITHUB_TOKEN 环境变量。\n"
            "请前往 https://github.com/settings/tokens 创建 Personal Access Token（需要 repo 权限），\n"
            "然后设置环境变量: export GITHUB_TOKEN=ghp_xxxxxxxxxxxx",
            "ERROR",
        )
        sys.exit(1)

    # ---- 7. 检测仓库 ----
    if args.repo:
        repo = args.repo
    else:
        repo = _detect_repo_from_git()

    _log(f"目标仓库: {repo}")

    # ---- 8. 检查是否已存在同版本 Release ----
    if _check_release_exists(repo, token, tag):
        if args.force:
            _log(f"Release {tag} 已存在，--force 模式：删除旧 Release...", "WARN")
            # 获取旧 Release 的 ID 并删除
            url = f"{GITHUB_API_BASE}/repos/{repo}/releases/tags/{tag}"
            old_release = _github_api("GET", url, token)
            delete_url = f"{GITHUB_API_BASE}/repos/{repo}/releases/{old_release['id']}"
            _github_api("DELETE", delete_url, token)
            _log(f"已删除旧 Release: {tag}")
        else:
            _log(
                f"Release {tag} 已存在！使用 --force 强制覆盖，或修改版本号后重试。",
                "ERROR",
            )
            sys.exit(1)

    # ---- 9. 创建 Release ----
    release_title = f"{skill_name} {tag}"
    _log(f"创建 Release: {release_title}...")

    release = _create_release(
        repo=repo,
        token=token,
        tag=tag,
        name=release_title,
        body=release_notes,
        draft=args.draft,
        prerelease=args.prerelease,
    )

    release_id = release["id"]
    release_url = release["html_url"]
    upload_url = release["upload_url"]

    _log(f"Release 创建成功 (ID: {release_id})")

    # ---- 10. 上传 zip 附件 ----
    asset = _upload_asset(upload_url, token, zip_path)

    download_url = asset.get("browser_download_url", "N/A")
    asset_size = asset.get("size", 0)

    # ---- 11. 完成 ----
    _log("=" * 60, "OK")
    _log(f"🎉 发布成功！", "OK")
    _log(f"   Release 页面: {release_url}", "OK")
    _log(f"   下载地址:     {download_url}", "OK")
    _log(f"   附件大小:     {asset_size / 1024 / 1024:.1f}MB", "OK")
    if args.draft:
        _log(f"   ⚠️  当前为草稿状态，请前往 GitHub 页面确认发布", "WARN")
    _log("=" * 60, "OK")

    # 清理临时文件（非 dry-run 且非指定输出目录时）
    if not args.output_dir:
        zip_path.unlink(missing_ok=True)


if __name__ == "__main__":
    main()
