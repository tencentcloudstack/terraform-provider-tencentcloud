#!/usr/bin/env python3
"""Check that the `Example Usage` HCL of a framework resource's sibling md
uses the HCL form its schema attribute types REQUIRE.

Why this exists: in terraform-plugin-framework the HCL syntax of a field is
decided by its Go schema type, and writing the wrong form is a hard Terraform
error ("Unsupported block type" / "Unsupported argument"):

    schema type                                          HCL form
    ---------------------------------------------------  ----------------------
    SingleNestedAttribute, ObjectAttribute,
    MapAttribute, MapNestedAttribute                     name = { ... }
    ListNestedAttribute, SetNestedAttribute,
    ListAttribute, SetAttribute                          name = [ ... ]
    SingleNestedBlock, ListNestedBlock, SetNestedBlock,
    anything registered under `Blocks:` (e.g. timeouts)  name { ... }
    String/Int64/Int32/Float64/Float32/Number/Bool       name = <scalar>

Usage (repository root):

    python3 tencentcloud/framework/templates/check_md_hcl.py <file.go> [...]
    python3 tencentcloud/framework/templates/check_md_hcl.py --all     # every framework resource with a sibling md

For each <file.go> the sibling `<stem>.md` is checked. Exit code 0 = no
problem. Run it on every generated resource BEFORE `make doc` (conventions
section 7.D step 2); `verify_templates.py` runs `--all` plus the self-test.
"""
import pathlib
import re
import sys

sys.dont_write_bytecode = True

REPO = pathlib.Path(__file__).resolve().parents[3]

KIND_BY_TYPE = {
    "SingleNestedAttribute": "obj", "ObjectAttribute": "obj",
    "MapAttribute": "obj", "MapNestedAttribute": "obj",
    "ListNestedAttribute": "list", "SetNestedAttribute": "list",
    "ListAttribute": "list", "SetAttribute": "list",
    "SingleNestedBlock": "block", "ListNestedBlock": "block", "SetNestedBlock": "block",
    "StringAttribute": "scalar", "Int64Attribute": "scalar", "Int32Attribute": "scalar",
    "Float64Attribute": "scalar", "Float32Attribute": "scalar", "NumberAttribute": "scalar",
    "BoolAttribute": "scalar",
}
# Map-valued kinds: the keys inside the braces are user data, not schema names.
FREE_KEY_TYPES = {"MapAttribute", "MapNestedAttribute"}
FORM_HINT = {
    "obj": "name = { ... }",
    "list": "name = [ ... ]",
    "block": "name { ... }",
    "scalar": "name = <value>",
}

DECL_RE = re.compile(r'"(\w+)"\s*:\s*(?:\w+\.)?(\w+(?:Attribute|Block))\{')
# `"timeouts": fooTimeoutsBlock(ctx)` / `resourcetimeouts.Block(ctx, ...)` helpers.
HELPER_BLOCK_RE = re.compile(r'"(\w+)"\s*:\s*(?:\w+\.)?\w*Block\w*\(')
SKIP_TOP = {"resource", "data", "action", "ephemeral", "provider", "locals", "variable",
            "output", "terraform", "module", "list", "import", "moved", "removed"}


def schema_kinds(go_src):
    """name -> set of kinds ({obj,list,block,scalar}); name -> True when ALL its
    declarations are map-valued (so the keys inside are free)."""
    kinds, free = {}, {}
    for name, typ in DECL_RE.findall(go_src):
        kind = KIND_BY_TYPE.get(typ)
        if kind is None:
            continue
        kinds.setdefault(name, set()).add(kind)
        free[name] = free.get(name, True) and typ in FREE_KEY_TYPES
    for name in HELPER_BLOCK_RE.findall(go_src):
        kinds.setdefault(name, set()).add("block")
        free[name] = False
    return kinds, free


def strip_line(line):
    """Drop comments and string contents so brackets inside strings are ignored."""
    out, i, in_str = [], 0, False
    while i < len(line):
        c = line[i]
        if in_str:
            if c == "\\":
                i += 2
                continue
            if c == '"':
                in_str = False
                out.append('"')
            i += 1
            continue
        if c == '"':
            in_str = True
            out.append('"')
        elif c == "#" or line[i:i + 2] == "//":
            break
        else:
            out.append(c)
        i += 1
    return "".join(out)


LITERAL_RE = re.compile(r'^(?:"|-?\d|true\b|false\b)')
KEY_RE = re.compile(r'^\s*("?[\w-]+"?)\s*(=)?\s*(.*)$')


def check_hcl(hcl, kinds, free):
    """Return a list of human-readable problems for one HCL snippet."""
    problems = []
    stack = ["top"]  # frame kinds: top / schema / free / list
    heredoc = None
    for no, raw in enumerate(hcl.splitlines(), 1):
        if heredoc:
            if raw.strip() == heredoc:
                heredoc = None
            continue
        m = re.search(r"<<-?\s*(\w+)\s*$", raw)
        line = strip_line(raw)
        if m:
            heredoc = m.group(1)
        text = line.strip()
        if not text:
            continue

        opener_kind = "schema"  # kind of the frame pushed by the first opener
        km = KEY_RE.match(line)
        keyed = bool(km) and text[0] not in "{}[]"
        if keyed:
            key = km.group(1).strip('"')
            has_eq = bool(km.group(2))
            rest = km.group(3).lstrip()
            frame = stack[-1]
            if frame == "schema" and key in kinds:
                if rest.startswith("{"):
                    form = "obj" if has_eq else "block"
                elif rest.startswith("["):
                    form = "list" if has_eq else None
                elif has_eq and LITERAL_RE.match(rest):
                    # Only a positively-identified scalar literal is judged;
                    # `null`, references (var.x / local.x) and function calls
                    # are valid for every schema type.
                    form = "scalar"
                else:
                    form = None
                if form is not None and form not in kinds[key]:
                    want = " or ".join(FORM_HINT[k] for k in sorted(kinds[key]))
                    got = {"obj": f"{key} = {{ ... }}", "list": f"{key} = [ ... ]",
                           "block": f"{key} {{ ... }}", "scalar": f"{key} = <value>"}[form]
                    problems.append(f"line {no}: `{key}` is written as `{got}` but its schema type requires `{want.replace('name', key)}`")
            if frame in ("top",) and key in SKIP_TOP:
                opener_kind = "schema"
            elif frame == "schema" and key in kinds and free.get(key) and has_eq and rest.startswith("{"):
                opener_kind = "free"

        # Walk the brackets of the line to maintain the frame stack.
        first = True
        for ch in line:
            if ch in "{[":
                if ch == "[":
                    kind = "list"
                elif first and keyed:
                    kind = opener_kind
                elif stack[-1] == "free":
                    kind = "schema"
                else:
                    kind = "schema"
                stack.append(kind)
                first = False
            elif ch in "}]":
                if len(stack) > 1:
                    stack.pop()
        # `key = value` inside a free (map) frame that opens an object holds schema-less data
        # only at the key level; its nested object is a normal object value.
    if len(stack) != 1:
        problems.append("unbalanced brackets in the HCL example")
    return problems


HCL_FENCE_RE = re.compile(r"```hcl\n(.*?)```", re.S)


def check_pair(go_path, md_path=None):
    go_path = pathlib.Path(go_path)
    md_path = pathlib.Path(md_path) if md_path else go_path.with_suffix(".md")
    if not md_path.exists():
        return [f"{md_path}: sibling md missing"]
    kinds, free = schema_kinds(go_path.read_text())
    out = []
    for block in HCL_FENCE_RE.findall(md_path.read_text()):
        for p in check_hcl(block, kinds, free):
            out.append(f"{md_path}: {p}")
    return out


def framework_pairs():
    pairs = []
    for gof in sorted((REPO / "tencentcloud/services").rglob("*.go")):
        if gof.name.endswith("_test.go"):
            continue
        md = gof.with_suffix(".md")
        if md.exists() and "terraform-plugin-framework/" in gof.read_text():
            pairs.append((gof, md))
    return pairs


def check_doc_examples():
    """The templates and skill docs themselves must never teach a wrong form.
    Names and kinds come from the demo schemas in the *.go.tmpl files; every
    doc line that STARTS with `name = {` / `name = [` / `name {` is judged
    against them (table rows and inline code spans never start a line)."""
    tdir = REPO / "tencentcloud/framework/templates"
    kinds, free = {}, {}
    for tmpl in sorted(tdir.glob("*.go.tmpl")):
        k, f = schema_kinds(tmpl.read_text())
        for name, ks in k.items():
            kinds.setdefault(name, set()).update(ks)
            free[name] = free.get(name, True) and f[name]
    docs = sorted(tdir.glob("*.md")) + sorted((REPO / ".codebuddy/skills").glob("tf-fw-codegen*/**/*.md"))
    problems = []
    line_re = re.compile(r"^\s*(\w+)\s*(=)?\s*([{\[])")
    for doc in docs:
        in_fence = False
        for no, raw in enumerate(doc.read_text().splitlines(), 1):
            if raw.lstrip().startswith("```"):
                in_fence = not in_fence
                continue
            m = line_re.match(strip_line(raw))
            if not m or m.group(1) not in kinds:
                continue
            name, has_eq, opener = m.group(1), bool(m.group(2)), m.group(3)
            form = ("obj" if has_eq else "block") if opener == "{" else ("list" if has_eq else None)
            if form is not None and form not in kinds[name]:
                rel = doc.relative_to(REPO)
                problems.append(f"{rel}: line {no}: `{name}` is shown in the wrong HCL form for its schema type")
    # the shipped sibling md example must pass the same rules as a real resource
    gen = (tdir / "general_template.go.tmpl").read_text()
    k, f = schema_kinds(gen)
    for block in HCL_FENCE_RE.findall((tdir / "sibling_md_template.md").read_text()):
        for p in check_hcl(block, k, f):
            problems.append(f"sibling_md_template.md: {p}")
    return problems


def selftest():
    """The detector must catch every wrong form and accept every right one."""
    src = '''
        "enc": schema.SingleNestedAttribute{
        "rules": schema.ListNestedAttribute{
        "labels": schema.MapAttribute{
        "tags": schema.ListNestedBlock{
        "spec": schema.SingleNestedBlock{
        "name": schema.StringAttribute{
        "zones": schema.SetAttribute{
        "envs": schema.MapNestedAttribute{
        "timeouts": fooTimeoutsBlock(ctx),
    '''
    kinds, free = schema_kinds(src)
    good = '''resource "x" "y" {
  name = "a"
  enc = {
    name = "n"
  }
  rules = [
    {
      name = "r1"
    },
    { name = "r2" },
  ]
  labels = {
    spec = "free key, not a block"
  }
  tags {
    name = "t1"
  }
  tags {
    name = "t2"
  }
  spec {
    name = "s"
  }
  zones = ["a", "b"]
  envs = {
    prod = {
      name = "p"
    }
  }
  timeouts {
    create = "15m"
  }
  enc = null
  rules = var.rules
  zones = null
  command = <<EOT
enc {
EOT
}'''
    assert check_hcl(good, kinds, free) == [], check_hcl(good, kinds, free)
    bad = {
        "enc written as block": ('resource "x" "y" {\n  enc {\n    name = "n"\n  }\n}', "enc"),
        "rules written as block": ('resource "x" "y" {\n  rules {\n    name = "r"\n  }\n}', "rules"),
        "rules written as object": ('resource "x" "y" {\n  rules = {\n    name = "r"\n  }\n}', "rules"),
        "tags written as list attr": ('resource "x" "y" {\n  tags = [{\n    name = "r"\n  }]\n}', "tags"),
        "spec written as object attr": ('resource "x" "y" {\n  spec = {\n    name = "r"\n  }\n}', "spec"),
        "enc written as string literal": ('resource "x" "y" {\n  enc = "x"\n}', "enc"),
        "timeouts written as attr": ('resource "x" "y" {\n  timeouts = {\n    create = "1m"\n  }\n}', "timeouts"),
        "zones written as block": ('resource "x" "y" {\n  zones {\n  }\n}', "zones"),
        "labels written as block": ('resource "x" "y" {\n  labels {\n  }\n}', "labels"),
        "scalar written as block": ('resource "x" "y" {\n  name {\n  }\n}', "name"),
        "nested wrong form": ('resource "x" "y" {\n  enc = {\n    spec = {\n    }\n  }\n}', "spec"),
        "nested inside list": ('resource "x" "y" {\n  rules = [\n    {\n      enc {\n      }\n    },\n  ]\n}', "enc"),
        "inside map value": ('resource "x" "y" {\n  envs = {\n    prod = {\n      enc {\n      }\n    }\n  }\n}', "enc"),
    }
    for label, (hcl, name) in bad.items():
        got = check_hcl(hcl, kinds, free)
        assert got and f"`{name}`" in got[0], f"selftest: detector missed `{label}`: {got}"
    assert check_hcl('resource "x" "y" {\n  enc = {\n', kinds, free), "selftest: unbalanced not caught"


def main(argv):
    selftest()
    if argv == ["--all"]:
        pairs = framework_pairs()
        problems = []
        for gof, md in pairs:
            problems += check_pair(gof, md)
        problems += check_doc_examples()
        print(f"--- md HCL form vs schema type ({len(pairs)} framework md scanned): {len(problems)} problem(s)")
    elif argv:
        problems = []
        for g in argv:
            problems += check_pair(g)
        print(f"--- md HCL form vs schema type ({len(argv)} file(s)): {len(problems)} problem(s)")
    else:
        print(__doc__)
        return 2
    for line in problems:
        print("   ", line)
    return 1 if problems else 0


if __name__ == "__main__":
    sys.exit(main(sys.argv[1:]))
