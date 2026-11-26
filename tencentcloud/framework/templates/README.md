# Framework reference code templates

This directory holds the **code convention templates** for the four
terraform-plugin-framework reference types (Resource / Data Source / Action /
Ephemeral Resource). They are learning assets for the AI code-generation
skills (`.codebuddy/skills/tf-fw-codegen*`), NOT business code.

## Why the `.go.tmpl` extension

Template files contain placeholder identifiers (fictional product `demo`)
and `// TEMPLATE:` decision notes, so they intentionally do **not** compile.
The `.tmpl` suffix keeps `go build ./...`, `go vet` and the repo lint scripts
from picking them up.

## Files

| File | resource_type | Purpose |
|---|---|---|
| `service_template.go.tmpl` | (shared by all) | Service layer skeleton. SDK calls NEVER go directly into lifecycle methods; they are wrapped in `service_tencentcloud_<product>.go` with retry + logging. |
| `general_template.go.tmpl` | `general` | Full CRUD managed resource: real cloud id or composite id, import support, refresh pattern. |
| `config_template.go.tmpl` | `config` | Configuration living on a parent resource: the locator input IS the state id (never a fixed constant), Create/Update share one apply path, no-op Delete, pre-write existence warning. |
| `attachment_template.go.tmpl` | `attachment` | Bind/unbind relation: bind APIs in create, unbind APIs in delete, composite `parent#child` id, no-op Update (force-recreate). |
| `datasource_template.go.tmpl` | `datasource` | Query data source: query inputs + computed list outputs. |
| `datasource_test_template.go.tmpl` | `datasource` | Unit test (gomonkey): tftypes raw config + ReadRequest stubs + service pagination test. |
| `action_template.go.tmpl` | `action` | One-shot operation: Invoke only, NO state — no `id`, no outputs, input attributes + `timeouts` block only. |
| `ephemeral_resource_template.go.tmpl` | `ephemeral_resource` | Ephemeral resource: Open (from `api.read`), sensitive values, optional Close. |
| `resource_test_template.go.tmpl` | `general` | Unit test (gomonkey): tfsdk plan/state helpers + full CRUD case matrix incl. async wait and timeout cases. |
| `config_test_template.go.tmpl` | `config` | Unit test: shared write path, locator-as-id, not-found sentinel, no-op Delete. |
| `attachment_test_template.go.tmpl` | `attachment` | Unit test: composite id, bind/unbind, parent-describe not-found, ImportState. |
| `action_test_template.go.tmpl` | `action` | Unit test: tftypes Config only (no plan/state), schema shape, async wait and timeout cases. |
| `ephemeral_test_template.go.tmpl` | `ephemeral_resource` | Unit test: tftypes Config + Result, sensitive output, every input reaches the request. |
| `sibling_md_template.md` | (shared by all) | Format of the mandatory sibling `<stem>.md` file consumed by `make doc`. |
| `verify_templates.py` | (tooling, not a template) | Gate: instantiates every `*.go.tmpl` in a throw-away package against a stub SDK and the REAL vendored framework — `go vet`, unused/colliding symbol check, `gofmt -s`, the sibling-md rules `make doc` enforces, the input validator self-test, then RUNS every test template and probe. Self-cleaning. |
| `validate_input.py` | (tooling, not a template) | Reference validator of the codegen JSON input (main skill step 1, checks 1-8): blocks illegal input BEFORE any file is written. `--selftest` runs the real inputs and ~30 illegal mutations; `--repo .` adds the repository-conflict checks. |
| `check_go_symbols.py` | (tooling, not a template) | Reports UNUSED private symbols and name COLLISIONS with sibling files of the same package. Run on every generated Go file. |
| `check_md_hcl.py` | (tooling, not a template) | Checks that the `Example Usage` HCL of a resource's sibling md uses the form each schema type requires (`*Attribute` => `name = { }` / `name = [ ]`, `*Block` => `name { }`), at every nesting depth. Run on every generated Go file before `make doc`; `--all` scans every framework resource. |

One implementation template and one test template per codegen
`resource_type` (the file names match the six JSON `resource_type` values),
plus the shared service layer and sibling md.

## Template conventions

- The fictional example uses product `demo`, SDK version `v20230101`, resource
  `widget`, type name `tencentcloud_demo_widget`. Substitute real values when
  generating code.
- `// TEMPLATE:` comments mark decision points that the generator (AI) must
  resolve; strip them from generated code.
- All templates are abstractions of reviewed, shipped implementations. When
  updating a template, keep it consistent with the patterns documented in
  `tencentcloud/framework/README.md` and the `tf-fw-codegen-common` skill.

## Maintenance rules

- **Every change to a `*.go.tmpl` file MUST be followed by**
  `python3 tencentcloud/framework/templates/verify_templates.py` and the run
  must end with `OK`. `gofmt -e` only proves syntax; wrong framework API
  usage (a `Computed` field in an action schema, `resp.State` on an action,
  a block placed in `Attributes`, a missing import, a redeclared named
  result) only fails at build time of a generated resource. A template that
  does not pass this gate must not be committed.
- When a repo-wide convention changes (e.g. the `fw` helper package, retry
  timeouts, registry wiring), update the matching template in the same change.
- Keep templates minimal: one representative field of each kind (required /
  optional+computed / computed / list-nested / sensitive), not every possible
  attribute shape.
