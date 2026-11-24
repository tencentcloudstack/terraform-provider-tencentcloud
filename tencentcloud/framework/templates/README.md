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
| `config_template.go.tmpl` | `config` | Singleton account-wide configuration: fixed constant id, Create/Update share one apply path, no-op Delete, pre-write existence warning. |
| `attachment_template.go.tmpl` | `attachment` | Bind/unbind relation: bind APIs in create, unbind APIs in delete, composite `parent#child` id, no-op Update (force-recreate). |
| `datasource_template.go.tmpl` | `datasource` | Query data source: query inputs + computed list outputs. |
| `datasource_test_template.go.tmpl` | `datasource` | Unit test (gomonkey): tftypes raw config + ReadRequest stubs + service pagination test. |
| `action_template.go.tmpl` | `action` | One-shot operation: Invoke, stateless. |
| `ephemeral_resource_template.go.tmpl` | `ephemeral_resource` | Ephemeral resource: Open (from `api.read`), sensitive values, optional Close. |
| `resource_test_template.go.tmpl` | `general` (+ other resource subtypes) | Unit test (gomonkey): tfsdk plan/state helpers + full CRUD case matrix. |
| `sibling_md_template.md` | (shared by all) | Format of the mandatory sibling `<stem>.md` file consumed by `make doc`. |

One template per codegen `resource_type` — the file names match the six
JSON `resource_type` values exactly, plus two shared files (service layer
and sibling md) used by every type.

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

- When a repo-wide convention changes (e.g. the `fw` helper package, retry
  timeouts, registry wiring), update the matching template in the same change.
- Keep templates minimal: one representative field of each kind (required /
  optional+computed / computed / list-nested / sensitive), not every possible
  attribute shape.
