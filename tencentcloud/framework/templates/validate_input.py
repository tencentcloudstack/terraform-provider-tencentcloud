#!/usr/bin/env python3
"""Reference validator for the codegen JSON input of the tf-fw-codegen skill.

It implements, mechanically, step 1 of the main skill (checks 1-8) so that an
illegal input is blocked BEFORE any file is written. The skill text and this
script describe the same rules; when they disagree, fix the script and the
skill together.

Usage (repository root):

    python3 tencentcloud/framework/templates/validate_input.py INPUT.json [--repo .]
    python3 tencentcloud/framework/templates/validate_input.py --selftest

Exit code 0 = legal, 1 = at least one problem (ALL problems are printed),
2 = usage error. `--repo` additionally runs the repository-conflict checks
(check 8); without it only the in-memory checks (1-7) run.
"""
import json
import pathlib
import re
import sys

sys.dont_write_bytecode = True

TYPES = ("general", "config", "attachment", "datasource", "action", "ephemeral_resource")
TOP_KEYS = {"resource_type", "service", "resource_name", "resource_id", "shipyard_authorization",
            "x_shipyard_ns", "env", "x_virtual_env", "api", "options"}
VERBS = ("create", "read", "update", "delete")
ENTRY_KEYS = {"name", "capi_uid", "async_info", "purpose"}
ASYNC_KEYS = {"name", "capi_uid", "async_info"}

# resource_type -> (min create, min read, min update, min delete); None = must be empty
MATRIX = {
    "general": (1, 1, 0, 0),
    "config": (None, 1, 1, None),
    "attachment": (1, 1, None, 1),
    "datasource": (None, 1, None, None),
    "action": (1, None, None, None),
    "ephemeral_resource": (None, 1, None, 0),
}

FILE_PREFIX = {"general": "resource", "config": "resource", "attachment": "resource",
               "datasource": "data_source", "action": "action", "ephemeral_resource": "ephemeral"}
FACTORY_SUFFIX = {"general": "Resource", "config": "Resource", "attachment": "Resource",
                  "datasource": "DataSource", "action": "", "ephemeral_resource": "EphemeralResource"}

RX_SERVICE = re.compile(r"^[a-z][a-z0-9]*$")
RX_NAME = re.compile(r"^tencentcloud_[a-z0-9]+(_[a-z0-9]+)*$")
RX_API = re.compile(r"^[A-Z][A-Za-z0-9]*$")
RX_UID = re.compile(r"^iacpcapi-[A-Za-z0-9]+$")


def strip_line_comments(text):
    """Remove `//` comments that are not inside a JSON string."""
    out = []
    for line in text.splitlines():
        in_str = False
        esc = False
        cut = None
        for i, ch in enumerate(line):
            if esc:
                esc = False
            elif ch == "\\" and in_str:
                esc = True
            elif ch == '"':
                in_str = not in_str
            elif ch == "/" and not in_str and line[i:i + 2] == "//":
                cut = i
                break
        out.append(line if cut is None else line[:cut])
    return "\n".join(out)


def camel(stem):
    return "".join(p[:1].upper() + p[1:] for p in stem.split("_"))


def target_names(data):
    """Return (stem, files, factory) derived from the input (conventions section 1)."""
    stem = data["resource_name"][len("tencentcloud_"):]
    prefix = FILE_PREFIX[data["resource_type"]]
    files = [f"{prefix}_tc_{stem}.go", f"{prefix}_tc_{stem}.md", f"{prefix}_tc_{stem}_test.go"]
    factory = "New" + camel(stem) + FACTORY_SUFFIX[data["resource_type"]]
    return stem, files, factory


def _entries(api):
    for verb in VERBS:
        arr = api.get(verb, [])
        if isinstance(arr, list):
            for i, e in enumerate(arr):
                if isinstance(e, dict):
                    yield verb, i, e


def validate(data):
    """In-memory checks 2-7. Returns the list of problems."""
    errs = []
    if not isinstance(data, dict):
        return ["top level must be a JSON object"]

    # 2. unknown keys
    for k in sorted(set(data) - TOP_KEYS):
        errs.append(f"unknown top-level key `{k}`")
    api = data.get("api")
    if not isinstance(api, dict):
        errs.append("`api` must be an object")
        api = {}
    for k in sorted(set(api) - set(VERBS)):
        errs.append(f"unknown key `api.{k}` (allowed: create/read/update/delete)")
    for verb in VERBS:
        arr = api.get(verb, [])
        if not isinstance(arr, list):
            errs.append(f"`api.{verb}` must be an array")
            continue
        for i, e in enumerate(arr):
            loc = f"api.{verb}[{i}]"
            if not isinstance(e, dict):
                errs.append(f"`{loc}` must be an object")
                continue
            for k in sorted(set(e) - ENTRY_KEYS):
                errs.append(f"unknown key `{loc}.{k}`")

    # 3. top-level fields
    rt = data.get("resource_type")
    if rt not in TYPES:
        errs.append(f"`resource_type` must be one of {', '.join(TYPES)} (got {rt!r})")
    svc = data.get("service")
    if not isinstance(svc, str) or not RX_SERVICE.match(svc):
        errs.append(f"`service` must match {RX_SERVICE.pattern} (got {svc!r})")
    rn = data.get("resource_name")
    if not isinstance(rn, str) or not RX_NAME.match(rn):
        errs.append(f"`resource_name` must match {RX_NAME.pattern} (got {rn!r})")
    rid = data.get("resource_id")
    if not isinstance(rid, str) or not rid.strip():
        errs.append("`resource_id` must be a non-empty string")
    auth = data.get("shipyard_authorization")
    if not isinstance(auth, str) or not auth.startswith("token ") or not auth[6:].strip():
        errs.append("`shipyard_authorization` must start with `token ` followed by the token")
    ns = data.get("x_shipyard_ns")
    if not isinstance(ns, str) or not ns.strip():
        errs.append("`x_shipyard_ns` must be non-empty")
    env = data.get("env")
    if env not in ("prod", "test", "dev"):
        errs.append(f"`env` must be prod/test/dev (got {env!r})")
    if env == "dev":
        vx = data.get("x_virtual_env")
        if not isinstance(vx, str) or not vx.strip():
            errs.append("`x_virtual_env` is required (e.g. dev.<user>) when env is dev")

    # 4. resource_id shape
    if isinstance(rid, str) and rid.strip():
        if rt == "datasource" and rid != "fwHelper.BuildToken()":
            errs.append("datasource `resource_id` must be exactly `fwHelper.BuildToken()`")
        if rt in ("action", "ephemeral_resource") and rid != "-":
            errs.append(f"{rt} `resource_id` must be exactly `-`")
        if rt == "attachment":
            segs = rid.split("#")
            if len(segs) < 2 or any(not s.strip() for s in segs):
                errs.append("attachment `resource_id` must be `#`-joined with at least two non-empty segments")
        if rt in ("general", "config"):
            if rid in ("-", "fwHelper.BuildToken()"):
                errs.append(f"{rt} `resource_id` must name a field, not a placeholder literal")
            if rt == "config" and not all(re.match(r"^[A-Za-z][A-Za-z0-9]*$", s) for s in rid.split("#")):
                errs.append("config `resource_id` must be a locator field name (or `#`-joined names)")
            if rt == "general" and not re.match(r"^[A-Za-z][A-Za-z0-9]*(#[A-Za-z][A-Za-z0-9]*)*$", rid):
                errs.append("general `resource_id` must be an SDK field name (or `#`-joined names)")

    # 5. entries
    uid_names = {}
    for verb, i, e in _entries(api):
        loc = f"api.{verb}[{i}]"
        name, uid = e.get("name"), e.get("capi_uid")
        if not isinstance(name, str) or not RX_API.match(name):
            errs.append(f"`{loc}.name` must match {RX_API.pattern} (got {name!r})")
        if not isinstance(uid, str) or not RX_UID.match(uid):
            errs.append(f"`{loc}.capi_uid` must match {RX_UID.pattern} (got {uid!r})")
        if isinstance(name, str) and isinstance(uid, str):
            uid_names.setdefault(("uid", uid), set()).add(name)
            uid_names.setdefault(("name", name), set()).add(uid)
        if "async_info" not in e:
            errs.append(f"`{loc}.async_info` is required (`{{}}` for a synchronous API)")
            continue
        ai = e["async_info"]
        if not isinstance(ai, dict):
            errs.append(f"`{loc}.async_info` must be an object")
            continue
        if ai:
            for k in sorted(set(ai) - ASYNC_KEYS):
                errs.append(f"unknown key `{loc}.async_info.{k}`")
            an, au = ai.get("name"), ai.get("capi_uid")
            if not isinstance(an, str) or not RX_API.match(an):
                errs.append(f"`{loc}.async_info.name` must match {RX_API.pattern} (got {an!r})")
            if not isinstance(au, str) or not RX_UID.match(au):
                errs.append(f"`{loc}.async_info.capi_uid` must match {RX_UID.pattern} (got {au!r})")
            if isinstance(an, str) and isinstance(au, str):
                uid_names.setdefault(("uid", au), set()).add(an)
                uid_names.setdefault(("name", an), set()).add(au)
            if ai.get("async_info", {}) != {}:
                errs.append(f"`{loc}.async_info.async_info` must be `{{}}` (chained polling is not supported)")

    # 6. consistency
    for (kind, key), others in sorted(uid_names.items()):
        if len(others) > 1:
            what = "names" if kind == "uid" else "capi_uids"
            errs.append(f"{kind} `{key}` maps to several {what}: {sorted(others)}")
    for verb in VERBS:
        arr = api.get(verb, [])
        if isinstance(arr, list):
            seen = set()
            for e in arr:
                n = e.get("name") if isinstance(e, dict) else None
                if isinstance(n, str):
                    if n in seen:
                        errs.append(f"`api.{verb}` lists `{n}` twice")
                    seen.add(n)

    # 7. key matrix
    if rt in MATRIX:
        for verb, need in zip(VERBS, MATRIX[rt]):
            arr = api.get(verb, [])
            n = len(arr) if isinstance(arr, list) else 0
            if need is None and n > 0:
                errs.append(f"{rt}: `api.{verb}` must be empty (has {n})")
            elif need is not None and n < need:
                errs.append(f"{rt}: `api.{verb}` needs at least {need} entr{'y' if need == 1 else 'ies'} (has {n})")
    return errs


def check_repo(data, repo):
    """Check 8: repository conflicts (code is the truth, never provider.md)."""
    errs = []
    repo = pathlib.Path(repo)
    name = data["resource_name"]
    for rel in ("tencentcloud/framework/registry.go", "tencentcloud/provider.go"):
        f = repo / rel
        if not f.exists():
            errs.append(f"{rel} not found under {repo}: pass the provider repository root with --repo")
            continue
        if re.search(r'"%s"\s*:' % re.escape(name), f.read_text()):
            errs.append(f"resource_name `{name}` is already registered in {rel}")
    _, files, factory = target_names(data)
    pkg = repo / "tencentcloud/services" / data["service"]
    if not pkg.is_dir():
        print(f"note: tencentcloud/services/{data['service']}/ does not exist: a NEW product (client accessor is created by the sub-skill)")
        return errs
    for fn in files:
        if (pkg / fn).exists():
            errs.append(f"target file already exists: tencentcloud/services/{data['service']}/{fn}")
    rx = re.compile(r"^func %s\(" % re.escape(factory), re.M)
    for go in sorted(pkg.glob("*.go")):
        if rx.search(go.read_text()):
            errs.append(f"factory `{factory}` is already declared in {go.relative_to(repo)}")
            break
    return errs


REAL_DATASOURCE = """{
  "resource_type": "datasource", // comment
  "service": "tdmysql",
  "resource_name": "tencentcloud_tdmysql_database_objects",
  "resource_id": "fwHelper.BuildToken()",
  "shipyard_authorization": "token 6b978f90caa24dfd87c82d12ecfe4e56",
  "x_shipyard_ns": "asyproj-qSP83cgZ48",
  "env": "prod",
  "x_virtual_env": "dev.seanyyan",
  "api": {"read": [{"name": "DescribeDatabaseObjects", "capi_uid": "iacpcapi-DNF3kysLlE", "async_info": {}}]}
}"""

_FLOW = {"name": "DescribeFlow", "capi_uid": "iacpcapi-mrlBaaFrIM", "async_info": {}}
REAL_GENERAL = {
    "resource_type": "general", "service": "tdmysql", "resource_name": "tencentcloud_tdmysql_db_instance_v2",
    "resource_id": "instanceId", "shipyard_authorization": "token 6b978f90caa24dfd87c82d12ecfe4e56",
    "x_shipyard_ns": "asyproj-qSP83cgZ48", "env": "prod", "x_virtual_env": "",
    "api": {
        "create": [
            {"name": "CreateDBInstances", "capi_uid": "iacpcapi-iRLoBGMPuV", "async_info": dict(_FLOW)},
            {"name": "ModifyAutoRenewFlag", "capi_uid": "iacpcapi-l8XAZphkvI", "async_info": {}},
        ],
        "read": [{"name": "DescribeDBInstanceDetail", "capi_uid": "iacpcapi-bXmM4kbjaO", "async_info": {}}],
        "update": [
            {"name": "ModifyDBInstance", "capi_uid": "iacpcapi-ED4AgQO6gn", "async_info": dict(_FLOW)},
            {"name": "ModifyAutoRenewFlag", "capi_uid": "iacpcapi-l8XAZphkvI", "async_info": {}},
        ],
        "delete": [
            {"name": "IsolateDBInstance", "capi_uid": "iacpcapi-vapWaOcpYh",
             "async_info": {"name": "DescribeDBInstanceDetail", "capi_uid": "iacpcapi-bXmM4kbjaO", "async_info": {}}},
            {"name": "DestroyInstances", "capi_uid": "iacpcapi-80qFg3z34t",
             "async_info": {"name": "DescribeDBInstanceDetail", "capi_uid": "iacpcapi-bXmM4kbjaO", "async_info": {}}},
        ],
    },
}


def selftest():
    import copy

    failures = []

    def expect(label, data, legal, needle=None):
        errs = validate(data)
        ok = (not errs) if legal else bool(errs) and (needle is None or any(needle in e for e in errs))
        if not ok:
            failures.append(f"{label}: legal={legal} errs={errs}")
        print(f"[{'ok' if ok else 'FAIL'}] {label}")

    expect("real datasource JSON (x_virtual_env set while env=prod)", json.loads(strip_line_comments(REAL_DATASOURCE)), True)
    expect("real general JSON (cross-verb reuse, nested async_info {})", REAL_GENERAL, True)

    def mutated(fn):
        d = copy.deepcopy(REAL_GENERAL)
        fn(d)
        return d

    expect("typo top-level key", mutated(lambda d: d.update(resouce_id="x")), False, "unknown top-level key")
    expect("typo verb key", mutated(lambda d: d["api"].update(creat=[])), False, "unknown key `api.creat`")
    expect("bad resource_type", mutated(lambda d: d.update(resource_type="resource")), False, "resource_type")
    expect("uppercase resource_name", mutated(lambda d: d.update(resource_name="tencentcloud_TdMysql_x")), False, "resource_name")
    expect("resource_name without prefix", mutated(lambda d: d.update(resource_name="tdmysql_db_instance")), False, "resource_name")
    expect("service with a dash", mutated(lambda d: d.update(service="td-mysql")), False, "service")
    expect("token without prefix", mutated(lambda d: d.update(shipyard_authorization="abc")), False, "shipyard_authorization")
    expect("bad env", mutated(lambda d: d.update(env="staging")), False, "env")
    expect("dev without x_virtual_env", mutated(lambda d: d.update(env="dev", x_virtual_env="")), False, "x_virtual_env")
    expect("lowercase api name", mutated(lambda d: d["api"]["read"][0].update(name="describeDBInstanceDetail")), False, "name")
    expect("bad capi_uid", mutated(lambda d: d["api"]["read"][0].update(capi_uid="bXmM4kbjaO")), False, "capi_uid")
    expect("missing async_info", mutated(lambda d: d["api"]["read"][0].pop("async_info")), False, "async_info")
    expect("half-filled async_info", mutated(lambda d: d["api"]["create"][0].update(async_info={"name": "DescribeFlow"})), False, "async_info.capi_uid")
    expect("chained async_info", mutated(lambda d: d["api"]["create"][0]["async_info"].update(async_info={"name": "X"})), False, "chained polling")
    expect("same capi_uid, two names", mutated(lambda d: d["api"]["read"][0].update(capi_uid="iacpcapi-iRLoBGMPuV")), False, "maps to several")
    expect("same name twice in one verb", mutated(lambda d: d["api"]["create"].append(copy.deepcopy(d["api"]["create"][1]))), False, "twice")
    expect("general without read", mutated(lambda d: d["api"].update(read=[])), False, "needs at least 1")
    expect("general with null async_info", mutated(lambda d: d["api"]["read"][0].update(async_info=None)), False, "async_info")

    cfg = {"resource_type": "config", "service": "tdmysql", "resource_name": "tencentcloud_tdmysql_x_config",
           "resource_id": "InstanceId", "shipyard_authorization": "token abc", "x_shipyard_ns": "ns", "env": "prod",
           "api": {"read": [{"name": "DescribeX", "capi_uid": "iacpcapi-aaa", "async_info": {}}],
                   "update": [{"name": "ModifyX", "capi_uid": "iacpcapi-bbb", "async_info": {}}]}}
    expect("legal config (x_virtual_env omitted)", cfg, True)
    expect("config with create", {**cfg, "api": {**cfg["api"], "create": cfg["api"]["update"]}}, False, "must be empty")
    expect("config fixed-constant id", {**cfg, "resource_id": "-"}, False, "placeholder")

    att = {**cfg, "resource_type": "attachment", "resource_name": "tencentcloud_tdmysql_x_attachment",
           "resource_id": "InstanceId#SecurityGroupId",
           "api": {"create": cfg["api"]["update"], "read": cfg["api"]["read"], "delete": cfg["api"]["update"]}}
    expect("legal attachment (same API as create and delete)", att, True)
    expect("attachment single segment id", {**att, "resource_id": "InstanceId"}, False, "at least two")
    expect("attachment with update", {**att, "api": {**att["api"], "update": cfg["api"]["update"]}}, False, "must be empty")

    act = {**cfg, "resource_type": "action", "resource_name": "tencentcloud_tdmysql_do_x", "resource_id": "-",
           "api": {"create": cfg["api"]["update"]}}
    expect("legal action", act, True)
    expect("action with fixed BuildToken id", {**act, "resource_id": "fwHelper.BuildToken()"}, False, "exactly `-`")
    expect("action with read", {**act, "api": {**act["api"], "read": cfg["api"]["read"]}}, False, "must be empty")

    eph = {**cfg, "resource_type": "ephemeral_resource", "resource_name": "tencentcloud_tdmysql_x_detail",
           "resource_id": "-", "api": {"read": cfg["api"]["read"]}}
    expect("legal ephemeral", eph, True)
    expect("legal ephemeral with delete (Close)", {**eph, "api": {**eph["api"], "delete": cfg["api"]["update"]}}, True)
    expect("datasource with literal -", {**eph, "resource_type": "datasource"}, False, "BuildToken")

    sink = target_names(REAL_GENERAL)
    if sink != ("tdmysql_db_instance_v2",
                ["resource_tc_tdmysql_db_instance_v2.go", "resource_tc_tdmysql_db_instance_v2.md",
                 "resource_tc_tdmysql_db_instance_v2_test.go"], "NewTdmysqlDbInstanceV2Resource"):
        failures.append(f"target_names: {sink}")
    for rtype, want in (("datasource", "NewTdmysqlDbInstanceV2DataSource"), ("action", "NewTdmysqlDbInstanceV2"),
                        ("ephemeral_resource", "NewTdmysqlDbInstanceV2EphemeralResource")):
        got = target_names({**REAL_GENERAL, "resource_type": rtype})[2]
        if got != want:
            failures.append(f"factory for {rtype}: {got} != {want}")
    print(f"[{'ok' if not failures else 'FAIL'}] derived file stem / factory names")

    print(f"\nvalidate_input selftest: {'OK' if not failures else 'FAILED'}")
    for f in failures:
        print("  ", f)
    return 1 if failures else 0


def main(argv):
    if argv == ["--selftest"]:
        return selftest()
    repo = None
    if "--repo" in argv:
        i = argv.index("--repo")
        if i + 1 >= len(argv):
            print(__doc__)
            return 2
        repo = argv[i + 1]
        argv = argv[:i] + argv[i + 2:]
    if len(argv) != 1:
        print(__doc__)
        return 2
    try:
        data = json.loads(strip_line_comments(pathlib.Path(argv[0]).read_text()))
    except (OSError, ValueError) as e:
        print(f"cannot parse {argv[0]}: {e}")
        return 1
    errs = validate(data)
    if not errs and repo:
        errs += check_repo(data, repo)
    for e in errs:
        print("ERROR:", e)
    print(f"validate_input: {len(errs)} problem(s)")
    return 1 if errs else 0


if __name__ == "__main__":
    sys.exit(main(sys.argv[1:]))
