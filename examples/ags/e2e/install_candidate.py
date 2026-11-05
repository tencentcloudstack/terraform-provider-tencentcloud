#!/usr/bin/env python3
"""Stage a local, explicitly unpublished Provider candidate in an isolated mirror."""

import argparse
import hashlib
import json
import os
import shutil
from pathlib import Path


def digest(path):
    hasher = hashlib.sha256()
    with path.open("rb") as stream:
        for chunk in iter(lambda: stream.read(1024 * 1024), b""):
            hasher.update(chunk)
    return hasher.hexdigest()


def install(binary, work_root):
    binary = binary.resolve(strict=True)
    if not os.access(binary, os.X_OK):
        raise ValueError("Provider binary is not executable")
    work_root.mkdir(parents=True, exist_ok=True)
    mirror = work_root / "provider-mirror"
    target = mirror / "registry.terraform.io/tencentcloudstack/tencentcloud/0.0.1-local/linux_amd64"
    target.mkdir(parents=True, exist_ok=True)
    executable = target / "terraform-provider-tencentcloud_v0.0.1-local"
    if executable.is_symlink():
        executable.unlink()
    elif executable.exists():
        raise ValueError("Refusing to overwrite a non-symlink Provider artifact")
    # Copy instead of symlinking a mutable build output: a later `make build`
    # must not change the candidate midway through a recorded run.
    shutil.copy2(binary, executable)
    config = work_root / "candidate.tfrc"
    config.write_text(
        "provider_installation {\n"
        "  filesystem_mirror {\n"
        f"    path = {json.dumps(str(mirror))}\n"
        '    include = ["registry.terraform.io/tencentcloudstack/tencentcloud"]\n'
        "  }\n"
        "  direct {\n"
        '    exclude = ["registry.terraform.io/tencentcloudstack/tencentcloud"]\n'
        "  }\n"
        "}\n"
    )
    manifest = {
        "source": "registry.terraform.io/tencentcloudstack/tencentcloud",
        "local_candidate_version": "0.0.1-local",
        "published": False,
        "source_binary_path": str(binary),
        "binary_path": str(executable),
        "binary_sha256": digest(executable),
        "terraform_cli_config": str(config),
    }
    (work_root / "candidate.json").write_text(json.dumps(manifest, indent=2) + "\n")
    return manifest


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--binary", type=Path, required=True)
    parser.add_argument("--work-root", type=Path, required=True)
    args = parser.parse_args()
    print(json.dumps(install(args.binary, args.work_root.resolve()), indent=2))
