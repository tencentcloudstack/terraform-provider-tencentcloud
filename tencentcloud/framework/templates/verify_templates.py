#!/usr/bin/env python3
"""Compile-check every framework template against the REAL terraform-plugin-framework.

Why this exists: `gofmt -e` only checks syntax. Wrong framework API usage (a
Computed field on an action schema, resp.State on an action, a block placed
in Attributes, a missing import, a redeclared named result ...) compiles
nowhere and only shows up when a generated resource is built. This script
instantiates each template in its own throw-away Go package, backed by a
stub SDK, and runs `go vet` + `go test` on them with the repository's
vendored framework.

Usage (from the repository root):

    python3 tencentcloud/framework/templates/verify_templates.py

Exit code 0 = every template compiles and every probe passes. The throw-away
directory is removed in a `finally` block, so a failed run leaves nothing
behind. Run it after EVERY change to a *.go.tmpl file.
"""
import pathlib
import re
import shutil
import subprocess
import sys

sys.dont_write_bytecode = True  # never leave a __pycache__ next to the templates
sys.path.insert(0, str(pathlib.Path(__file__).resolve().parent))
import check_go_symbols  # noqa: E402

REPO = pathlib.Path(__file__).resolve().parents[3]
TPL = REPO / "tencentcloud/framework/templates"
BASE = REPO / "tencentcloud/framework/zz_tmplcheck"
MOD = "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/zz_tmplcheck"

SDK_ORIG = "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/demo/v20230101"
SDK_NEW = MOD + "/sdk/v20230101"
SVDEMO_ORIG = "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/demo"

STRUCTS = {
    "Widget": [("WidgetId", "*string"), ("Name", "*string"), ("ZoneId", "*string"),
               ("Remark", "*string"), ("Status", "*string"), ("CreateOn", "*string"),
               ("Tags", "[]*string")],
    "WidgetSettings": [("RetentionDays", "*int64"), ("Enabled", "*bool"), ("NotifyEmail", "*string")],
    "WidgetBinding": [("GadgetId", "*string"), ("WidgetId", "*string"), ("BindMode", "*string"),
                      ("BoundAt", "*string"), ("BoundStatus", "*string")],
    "Gadget": [("GadgetId", "*string"), ("WidgetBindings", "[]*WidgetBinding")],
}
# api -> (request fields, response fields)
APIS = {
    "CreateWidget": ([("Name", "*string"), ("ZoneId", "*string"), ("Remark", "*string"), ("Tags", "[]*string")],
                     [("WidgetId", "*string"), ("TaskId", "*string")]),
    "DescribeWidgets": ([("Name", "*string"), ("WidgetIds", "[]*string"), ("Offset", "*int64"), ("Limit", "*int64")],
                        [("WidgetSet", "[]*Widget"), ("TotalCount", "*int64")]),
    "ModifyWidget": ([("WidgetId", "*string"), ("Name", "*string"), ("Remark", "*string")],
                     [("TaskId", "*string")]),
    "DeleteWidget": ([("WidgetId", "*string")], []),
    "GetAsyncStatus": ([("TaskId", "*string")], [("Status", "*string")]),
    "UpdateWidgetConfig": ([("InstanceId", "*string"), ("Settings", "*WidgetSettings")], []),
    "DescribeWidgetConfig": ([("InstanceId", "*string")],
                             [("Settings", "*WidgetSettings"), ("UpdateTime", "*string"), ("Modifier", "*string")]),
    "AttachWidget": ([("GadgetId", "*string"), ("WidgetId", "*string"), ("BindMode", "*string")], []),
    "DetachWidget": ([("GadgetId", "*string"), ("WidgetId", "*string")], []),
    "DescribeGadgets": ([("GadgetIds", "[]*string")], [("GadgetSet", "[]*Gadget")]),
    "RestartWidget": ([("WidgetId", "*string")], [("TaskId", "*string")]),
    "DescribeTempCredential": ([("WidgetId", "*string"), ("Duration", "*int64")],
                               [("Token", "*string"), ("ExpiredAt", "*string")]),
}
# the SDK client methods the service template itself invokes
CLIENT_METHODS = [
    "CreateWidget", "DescribeWidgets", "ModifyWidget", "DeleteWidget", "GetAsyncStatus",
    "UpdateWidgetConfig", "DescribeWidgetConfig", "AttachWidget", "DetachWidget",
    "DescribeGadgets", "RestartWidget", "DescribeTempCredential",
]

# Methods the service template lists in its shape map but does not repeat in
# code. The harness generates them in EXACTLY the same retry-closure shape so
# the type templates and their test templates run against a realistic service
# layer (they are NOT shipped templates; the shapes come from service_template).
_EXTRA_HEAD = """package demo

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	demov20230101 "__SDK__"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	fwHelper "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

var _ = fwHelper.ErrNotFound
"""

_RETRY_HEAD = r"""
func (me *DemoService) @API@(ctx context.Context, request *demov20230101.@API@Request) @SIG@ {
	logId := tccommon.GetLogId(ctx)
	errRet := resource.Retry(tccommon.@TIMEOUT@, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := demoSdkClient().@API@(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("@API@ failed, Response is nil."))
		}

@BODY@
		return nil
	})

	if errRet != nil {
		log.Printf("[CRITAL]%s @API@ failed, reason:%+v", logId, errRet)

		return @FAIL@
	}

	return @OK@
}
"""


def _method(api, sig, timeout, body, fail, ok):
    return (_RETRY_HEAD.replace("@API@", api).replace("@SIG@", sig).replace("@TIMEOUT@", timeout)
            .replace("@BODY@", body).replace("@FAIL@", fail).replace("@OK@", ok))


EXTRA_SERVICE_METHODS = (
    _EXTRA_HEAD
    # plain write shape: error only
    + _method("UpdateWidgetConfig", "error", "WriteRetryTimeout", "", "errRet", "nil")
    + _method("AttachWidget", "error", "WriteRetryTimeout", "", "errRet", "nil")
    + _method("DetachWidget", "error", "WriteRetryTimeout", "", "errRet", "nil")
    # async write shape: returns the task id
    + _method("RestartWidget", "(taskId string, err error)", "WriteRetryTimeout",
              "\t\tif result.Response.TaskId != nil {\n\t\t\ttaskId = *result.Response.TaskId\n\t\t}\n",
              '"", errRet', "taskId, nil")
    # list describe shape
    + _method("DescribeGadgets", "(gadgets []*demov20230101.Gadget, err error)", "ReadRetryTimeout",
              "\t\tgadgets = result.Response.GadgetSet\n", "nil, errRet", "gadgets, nil")
    # singleton describe shape: a response without the object is ErrNotFound
    + _method("DescribeWidgetConfig", "(params *demov20230101.DescribeWidgetConfigResponseParams, err error)",
              "ReadRetryTimeout",
              "\t\tif result.Response.Settings == nil {\n\t\t\treturn resource.NonRetryableError(fwHelper.ErrNotFound)\n\t\t}\n\n\t\tparams = result.Response\n",
              "nil, errRet", "params, nil")
    + _method("DescribeTempCredential", "(token string, expiredAt string, err error)", "ReadRetryTimeout",
              "\t\tif result.Response.Token == nil {\n\t\t\treturn resource.NonRetryableError(fwHelper.ErrNotFound)\n\t\t}\n\n"
              "\t\ttoken = *result.Response.Token\n\t\tif result.Response.ExpiredAt != nil {\n\t\t\texpiredAt = *result.Response.ExpiredAt\n\t\t}\n",
              '"", "", errRet', "token, expiredAt, nil")
).replace("__SDK__", SDK_NEW)

# probe tests: behaviours a compiler cannot prove
PROBES = {
    "config": ('probe_test.go', '''package demo

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// A cloud answer without a settings object must not panic.
func TestProbeFlattenNilSettings(t *testing.T) {
	var diags diag.Diagnostics
	got := flattenWidgetSettings(context.Background(), nil, types.ObjectNull(DemoWidgetConfigSettingsAttrTypes()), &diags)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	if got.IsNull() {
		t.Fatal("flatten with nil settings must still produce an object")
	}
}
'''),
    "general": ('probe_test.go', '''package demo

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// Import (fw.WithImportByID) followed by a Read-side State.Get must decode
// the model, including the timeouts block, from an id-only state.
func TestProbeImportThenDecode(t *testing.T) {
	ctx := context.Background()
	r := &DemoWidgetResource{}

	var sresp resource.SchemaResponse
	r.Schema(ctx, resource.SchemaRequest{}, &sresp)

	resp := resource.ImportStateResponse{State: tfsdk.State{
		Schema: sresp.Schema,
		Raw:    tftypes.NewValue(sresp.Schema.Type().TerraformType(ctx), nil),
	}}
	r.ImportState(ctx, resource.ImportStateRequest{ID: "widget-1"}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("ImportState failed: %v", resp.Diagnostics)
	}

	var m DemoWidgetResourceModel
	if diags := resp.State.Get(ctx, &m); diags.HasError() {
		t.Fatalf("State.Get after import failed: %v", diags)
	}
	if m.Id.ValueString() != "widget-1" {
		t.Fatalf("id = %q, want widget-1", m.Id.ValueString())
	}
}
'''),
    "action": ('probe_test.go', '''package demo

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/action"
)

// An action schema must build without diagnostics and expose NO id/state
// attribute and a `timeouts` block.
func TestProbeActionSchema(t *testing.T) {
	var resp action.SchemaResponse
	(&DemoRestartWidget{}).Schema(context.Background(), action.SchemaRequest{}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("schema diagnostics: %v", resp.Diagnostics)
	}
	if _, ok := resp.Schema.Attributes["id"]; ok {
		t.Fatal("an action schema must not declare an id attribute")
	}
	if _, ok := resp.Schema.Blocks["timeouts"]; !ok {
		t.Fatal("the timeouts block must live in Blocks")
	}
}
'''),
}

ACCESSOR_PATCH_RX = re.compile(r'^[ \t]*patches\.ApplyMethodReturn\(client, "UseDemoV20230101Client", demoClient\)\n', re.M)

# dir -> (type template, [(test template, file name)])
UNITS = {
    "general": ("general_template.go.tmpl", [("resource_test_template.go.tmpl", "resource_test.go")]),
    "config": ("config_template.go.tmpl", [("config_test_template.go.tmpl", "config_test.go")]),
    "attachment": ("attachment_template.go.tmpl", [("attachment_test_template.go.tmpl", "attachment_test.go")]),
    "datasource": ("datasource_template.go.tmpl", [("datasource_test_template.go.tmpl", "datasource_test.go")]),
    "action": ("action_template.go.tmpl", [("action_test_template.go.tmpl", "action_test.go")]),
    "ephemeral": ("ephemeral_resource_template.go.tmpl", [("ephemeral_test_template.go.tmpl", "ephemeral_test.go")]),
}


def build_sdk():
    sdk = ['package v20230101', '', 'import (\n\t"encoding/json"\n\t"fmt"\n)', '', 'type Client struct{}', '']
    for name, fields in STRUCTS.items():
        sdk.append(f"type {name} struct {{")
        sdk += [f"\t{f} {t}" for f, t in fields]
        sdk.append("}\n")
    for api, (req, resp) in APIS.items():
        sdk.append(f"type {api}Request struct {{")
        sdk += [f"\t{f} {t}" for f, t in req]
        sdk.append("}")
        sdk.append(f"func New{api}Request() *{api}Request {{ return &{api}Request{{}} }}")
        sdk.append(f'func (r *{api}Request) GetAction() string {{ return "{api}" }}')
        sdk.append(f"func (r *{api}Request) ToJsonString() string {{ b, _ := json.Marshal(r); return string(b) }}")
        sdk.append(f"type {api}ResponseParams struct {{")
        sdk += [f"\t{f} {t}" for f, t in resp]
        sdk.append("\tRequestId *string")
        sdk.append("}")
        sdk.append(f"type {api}Response struct {{ Response *{api}ResponseParams }}")
        sdk.append(f"func (r *{api}Response) ToJsonString() string {{ b, _ := json.Marshal(r); return string(b) }}")
        if api in CLIENT_METHODS:
            # Deliberately NOT a tiny body: gomonkey overwrites the first bytes
            # of a patched function with a jump (16 bytes on arm64); a stub
            # shorter than that corrupts its neighbour. An unpatched call fails
            # loudly instead of silently returning (nil, nil).
            sdk.append(
                f"func (c *Client) {api}(request *{api}Request) (*{api}Response, error) {{\n"
                f"\tif request == nil {{\n"
                f"\t\treturn nil, fmt.Errorf(\"stub client: %s called with a nil request\", \"{api}\")\n"
                f"\t}}\n"
                f"\treturn nil, fmt.Errorf(\"stub client: %s is not patched by the test\", request.GetAction())\n"
                f"}}"
            )
        sdk.append("")
    (BASE / "sdk/v20230101").mkdir(parents=True)
    (BASE / "sdk/v20230101/sdk.go").write_text("\n".join(sdk))


def build_units():
    service_src = (TPL / "service_template.go.tmpl").read_text()
    # the real accessor lives on connectivity.TencentCloudClient; the harness
    # swaps it for a plain constructor of the stub client.
    service_src = service_src.replace("me.client.UseDemoV20230101Client()", "demoSdkClient()")
    service_src = service_src.replace(SDK_ORIG, SDK_NEW)
    accessor = ('package demo\n\nimport demov20230101 "%s"\n\n'
                'func demoSdkClient() *demov20230101.Client { return &demov20230101.Client{} }\n') % SDK_NEW

    for d, (main, tests) in UNITS.items():
        out = BASE / d
        out.mkdir(parents=True)
        (out / "service.go").write_text(service_src)
        (out / "accessor.go").write_text(accessor)
        (out / "extra_service.go").write_text(EXTRA_SERVICE_METHODS)
        (out / "main.go").write_text((TPL / main).read_text().replace(SDK_ORIG, SDK_NEW))
        for t, fname in tests:
            body = (TPL / t).read_text().replace(SDK_ORIG, SDK_NEW).replace(SVDEMO_ORIG, MOD + "/" + d)
            # The fictional demo product has no versioned accessor on the REAL
            # connectivity.TencentCloudClient, so the harness service builds
            # its stub SDK client directly (see build_units) and the accessor
            # patch line is dropped; every other line runs as generated.
            body = ACCESSOR_PATCH_RX.sub("", body)
            (out / fname).write_text(body)
        if d in PROBES:
            fname, body = PROBES[d]
            (out / fname).write_text(body)


def check_symbols():
    """UNUSED / COLLIDE gate (see check_go_symbols.py) over the template-derived
    files of every package: main.go and the test templates."""
    problems = []
    for d, (_, tests) in UNITS.items():
        files = [BASE / d / "main.go"] + [BASE / d / fname for _, fname in tests]
        problems += [f"{d}/{line}" for line in check_go_symbols.check(files)]
    print(f"--- unused / colliding Go symbols: {len(problems)}")
    for line in problems:
        print("   ", line)
    return len(problems)


def check_gofmt():
    """Every *.go.tmpl must already be `gofmt -s` clean — the same flavour
    `make ready` (fmt-faster) applies — checked on the template text itself,
    not on the harness copies."""
    fmt_dir = BASE / "_gofmt"
    fmt_dir.mkdir(parents=True, exist_ok=True)
    originals = {}
    for tmpl in sorted(TPL.glob("*.go.tmpl")):
        copy = fmt_dir / (tmpl.name + ".go")
        copy.write_text(tmpl.read_text())
        originals[copy.name] = tmpl.name
    proc = subprocess.run(["gofmt", "-s", "-l", str(fmt_dir)], cwd=REPO, capture_output=True, text=True)
    bad = [originals.get(pathlib.Path(line).name, line) for line in proc.stdout.split() if line]
    print(f"--- gofmt -s (*.go.tmpl): {len(bad)} unformatted")
    for line in bad:
        print("   ", line)
    return len(bad) + (1 if proc.returncode != 0 else 0)


def check_sibling_md():
    """The sibling md template must satisfy the structural rules gendoc
    (`make doc`) enforces on every framework resource: after the template
    comments are stripped the text is a non-empty one-line description, the
    standalone line `Example Usage` (REQUIRED, preceded by a description),
    and an optional standalone line `Import`; no markdown heading anywhere
    before Example Usage."""
    import re

    raw = (TPL / "sibling_md_template.md").read_text()
    text = re.sub(r"<!--.*?-->", "", raw, flags=re.S).strip()
    problems = []
    if "\nExample Usage\n" not in text:
        problems.append("missing the standalone `Example Usage` line (gendoc: example usage missing)")
    else:
        head = text[: text.index("\nExample Usage\n")]
        if not head.strip():
            problems.append("description before `Example Usage` is empty (gendoc: description empty)")
        if re.search(r"^#", head, re.M):
            problems.append("a markdown heading precedes `Example Usage`")
        if "_" in head.strip().splitlines()[0].replace("tencentcloud_", ""):
            problems.append("the one-line description contains an underscore")
    if text.count("\nImport\n") > 1:
        problems.append("more than one standalone `Import` line")
    print(f"--- sibling md template vs gendoc rules: {len(problems)} problem(s)")
    for line in problems:
        print("   ", line)
    return len(problems)


def check_md_hcl_gate():
    """Example Usage HCL must use the form each schema type REQUIRES (attribute
    `= {` / `= [` vs block `name {`; conventions section 6). Runs the checker's
    own mutation self-test, then every framework resource that has a sibling md,
    and finally the HCL examples inside the templates / skill docs themselves."""
    import check_md_hcl

    check_md_hcl.selftest()
    problems = []
    pairs = check_md_hcl.framework_pairs()
    for gof, md in pairs:
        problems += check_md_hcl.check_pair(gof, md)
    problems += check_md_hcl.check_doc_examples()
    print(f"--- md HCL form vs schema type ({len(pairs)} framework md + doc examples): {len(problems)} problem(s)")
    for line in problems:
        print("   ", line)
    return len(problems)


def check_validator():
    """The codegen-input validator must accept the real inputs and reject every
    illegal mutation (its own self-test)."""
    proc = subprocess.run([sys.executable, str(TPL / "validate_input.py"), "--selftest"],
                          cwd=REPO, capture_output=True, text=True)
    out = proc.stdout + proc.stderr
    oks = len([l for l in out.splitlines() if l.startswith("[ok]")])
    print(f"--- validate_input.py self-test: exit {proc.returncode}, {oks} cases ok")
    if proc.returncode != 0:
        print(out.strip())
    return proc.returncode


def check_helper_docs():
    """The helper table in conventions section 3.2 must match the helper
    package: every exported function is listed, and every `fwHelper.X` the
    skills / templates / README mention exists. A helper that is not in the
    table is invisible to the generator (it hand-rolls the conversion
    instead); a documented name that does not exist is a compile error in the
    generated code."""
    import re

    helper_dir = REPO / "tencentcloud/framework/helper"
    exported = set()
    funcs = set()
    for f in helper_dir.glob("*.go"):
        if f.name.endswith("_test.go"):
            continue
        src = f.read_text()
        funcs |= set(re.findall(r"^func ([A-Z]\w*)", src, re.M))
        exported |= set(re.findall(r"^func ([A-Z]\w*)", src, re.M))
        exported |= set(re.findall(r"^(?:var|const|type) ([A-Z]\w*)", src, re.M))
        for blk in re.findall(r"^(?:var|const) \(\n(.*?)^\)", src, re.M | re.S):
            exported |= set(re.findall(r"^\t([A-Z]\w*)", blk, re.M))

    skills = REPO / ".codebuddy/skills"
    conv_path = skills / "tf-fw-codegen-common/references/conventions.md"
    if not conv_path.exists():
        print("--- helper docs vs helper package: skipped (skills directory not present)")
        return 0
    conv = conv_path.read_text()

    problems = []
    docs = [p for p in skills.rglob("*.md") if "tf-fw-codegen" in str(p)] + [TPL / "README.md"]
    for p in docs + sorted(TPL.glob("*.tmpl")):
        for m in re.finditer(r"fwHelper\.([A-Za-z0-9_*]+)", p.read_text()):
            name = m.group(1)
            if "*" not in name and name not in exported:
                problems.append(f"{p.name}: fwHelper.{name} does not exist")
    for name in sorted(funcs):
        if name not in conv:
            problems.append(f"helper function `{name}` is not documented in conventions section 3.2")
    print(f"--- helper docs vs helper package: {len(problems)} problem(s)")
    for line in problems:
        print("   ", line)
    return len(problems)


def run(args, label):
    proc = subprocess.run(args, cwd=REPO, capture_output=True, text=True)
    out = (proc.stdout + proc.stderr).replace(MOD + "/", "").replace("tencentcloud/framework/zz_tmplcheck/", "")
    passed = len(re.findall(r"^\s*--- PASS:", out, re.M))
    summary = f", {passed} test cases passed" if "-v" in args else ""
    print(f"--- {label}: exit {proc.returncode}{summary}")
    if proc.returncode != 0:
        print(out.strip())
    elif "-v" in args and passed == 0:
        print("    no test case ran: a gate that executes nothing proves nothing")
        return 1
    return proc.returncode


def main():
    if BASE.exists():
        shutil.rmtree(BASE)
    failed = 0
    try:
        build_sdk()
        build_units()
        pkg = "./tencentcloud/framework/zz_tmplcheck/..."
        # -gcflags=-e lifts the 10-error-per-package cap: ALL errors at once.
        failed += run(["go", "vet", "-gcflags=-e", pkg], "go vet (every template package, incl. test templates)") != 0
        failed += check_symbols() != 0
        failed += check_gofmt() != 0
        failed += check_sibling_md() != 0
        failed += check_validator() != 0
        failed += check_md_hcl_gate() != 0
        failed += check_helper_docs() != 0
        # Runs the two TEST templates (gomonkey needs -l: no inlining) AND the
        # behavioural probes — what the compiler cannot prove: nil safety,
        # import decoding, action schema shape, the wait/timeouts paths,
        # not-found handling.
        failed += run(["go", "test", "-v", "-gcflags=all=-l", "-count=1", pkg],
                      "go test (test templates + behavioural probes)") != 0
    finally:
        shutil.rmtree(BASE, ignore_errors=True)

    if failed:
        print("\nFAILED: at least one template does not compile or a probe failed.")
        return 1
    print("\nOK: all 6 type templates + service + 2 test templates compile against the real "
          "framework, no unused private symbol exists, and every test case and probe passes.")
    return 0


if __name__ == "__main__":
    sys.exit(main())
