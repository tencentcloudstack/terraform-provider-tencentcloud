#!/usr/bin/env python3
"""Validate examples against the actual local mux binary without cloud writes."""

import argparse
import hashlib
import json
import os
import shutil
import subprocess
from pathlib import Path

from install_candidate import digest, install


def run(argv, directory, env, evidence, expected=(0,)):
    result = subprocess.run(argv, cwd=directory, env=env, text=True, capture_output=True, timeout=120)
    evidence.write_text(result.stdout + result.stderr)
    if result.returncode not in expected:
        raise RuntimeError(f"{argv[1:]} returned {result.returncode}; inspect {evidence}")
    return result


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--binary", required=True, type=Path)
    parser.add_argument("--terraform", required=True, type=Path)
    parser.add_argument("--work-root", required=True, type=Path)
    args = parser.parse_args()
    os.umask(0o077)
    root = args.work_root.resolve()
    if (root / "examples").exists():
        raise ValueError("Use a fresh work root; existing state or lockfiles must not be overwritten")
    candidate = install(args.binary, root)
    source = Path(__file__).resolve().parents[1]
    examples = root / "examples"
    shutil.copytree(source, examples, ignore=shutil.ignore_patterns("e2e", ".terraform", "*.tfstate*", ".terraform.lock.hcl"))
    env = dict(os.environ)
    for key in list(env):
        if key.startswith("TF_CLI_ARGS") or key in ("TF_LOG", "TF_LOG_PATH"):
            del env[key]
    env.update(TF_CLI_CONFIG_FILE=candidate["terraform_cli_config"], TF_IN_AUTOMATION="1", TF_INPUT="0", CHECKPOINT_DISABLE="1")
    outcomes = []
    for name in ["code", "browser", "custom", "existing", "combined"]:
        directory = examples / name
        run([str(args.terraform), "fmt", "-check"], directory, env, root / f"{name}-fmt.log")
        run([str(args.terraform), "init", "-backend=false", "-input=false"], directory, env, root / f"{name}-init.log")
        result = run([str(args.terraform), "validate", "-json"], directory, env, root / f"{name}-validate.json")
        validation = json.loads(result.stdout)
        if not validation["valid"] or validation["error_count"]:
            raise RuntimeError(f"{name} validation failed")
        outcomes.append({"example": name, "valid": True, "error_count": 0})
    result = run([str(args.terraform), "providers", "schema", "-json"], examples / "code", env, root / "provider-schema-all.json")
    provider = json.loads(result.stdout)["provider_schemas"][candidate["source"]]
    selected = {
        "resource": provider["resource_schemas"]["tencentcloud_ags_sandbox_tool"],
        "data_source": provider["data_source_schemas"]["tencentcloud_ags_sandbox_tool"],
    }
    resource = selected["resource"]["block"]
    required_blocks = {"network_configuration", "custom_configuration", "storage_mounts", "log_configuration", "timeouts"}
    if not required_blocks.issubset(resource["block_types"]):
        raise RuntimeError("Mux schema omitted a required AGS block")
    selectors = selected["data_source"]["block"]["attributes"]
    if selectors["tool_id"].get("computed") or selectors["tool_name"].get("computed"):
        raise RuntimeError("Data Source selectors must not be populated as result fields")
    (root / "ags-schema.json").write_text(json.dumps(selected, indent=2) + "\n")
    # A fresh create plan is useful offline evidence, but never a cloud apply.
    result = run([str(args.terraform), "plan", "-refresh=false", "-input=false", "-detailed-exitcode", "-out=create.plan", "-var=tool_name=tf-ags-offline-plan"], examples / "code", env, root / "code-plan.log", expected=(2,))
    result = run([str(args.terraform), "show", "-json", "create.plan"], examples / "code", env, root / "code-plan.json")
    changes = [item for item in json.loads(result.stdout).get("resource_changes", []) if item["mode"] == "managed"]
    if len(changes) != 1 or changes[0]["change"]["actions"] != ["create"]:
        raise RuntimeError("Fresh code example did not plan exactly one Tool create")
    if digest(args.binary) != candidate["binary_sha256"]:
        raise RuntimeError("Provider binary changed during validation")
    report = {"candidate_sha256": candidate["binary_sha256"], "published": False, "examples": outcomes, "protocol5_schema": "PASS", "fresh_code_plan": "create-only", "cloud_applies": 0}
    (root / "validation-report.json").write_text(json.dumps(report, indent=2) + "\n")
    print(json.dumps(report, indent=2))


if __name__ == "__main__":
    main()
