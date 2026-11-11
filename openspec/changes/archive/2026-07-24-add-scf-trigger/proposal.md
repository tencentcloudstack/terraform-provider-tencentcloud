## Why

The Terraform Provider for TencentCloud currently only offers `tencentcloud_scf_trigger_config`, which manages an existing SCF trigger through `UpdateTrigger` and has an empty Delete operation — it does not model the full trigger lifecycle. Users cannot create, read, update, and destroy SCF function triggers in a declarative, idempotent way through Terraform.

By adding a new `tencentcloud_scf_trigger` resource that maps to the native SCF trigger CRUD APIs (`CreateTrigger`, `ListTriggers`, `UpdateTrigger`, `DeleteTrigger`), users can fully manage the SCF (Serverless Cloud Function) trigger lifecycle as part of their infrastructure-as-code workflows.

## What Changes

- Add a new Terraform resource `tencentcloud_scf_trigger` of kind RESOURCE_KIND_GENERAL under `tencentcloud/services/scf/`.
- Implement full CRUD:
  - **Create** → calls SCF `CreateTrigger` API.
  - **Read** → calls SCF `ListTriggers` API, filtered by trigger name, to fetch the `TriggerInfo`.
  - **Update** → calls SCF `UpdateTrigger` API for mutable fields.
  - **Delete** → calls SCF `DeleteTrigger` API.
- Add the resource schema with parameters: `function_name`, `trigger_name`, `type`, `trigger_desc`, `namespace`, `qualifier`, `enable`, `custom_argument`, `description`, and computed fields from `TriggerInfo` (e.g. `available_status`, `add_time`, `mod_time`).
- Use a composite id (`function_name` + `namespace` + `trigger_name`) joined by `tccommon.FILED_SP`.
- Register the new resource in `tencentcloud/provider.go` and `tencentcloud/provider.md`.
- Add a service-layer helper `DescribeScfTriggerById` in `service_tencentcloud_scf.go`.
- Add the resource documentation `resource_tc_scf_trigger.md`.
- Add unit tests with gomonkey mocks (`resource_tc_scf_trigger_test.go`).

## Capabilities

### New Capabilities
- `scf-trigger-resource`: Full lifecycle management (create, read, update, delete) of a TencentCloud SCF (Serverless Cloud Function) function trigger via the `tencentcloud_scf_trigger` Terraform resource.

### Modified Capabilities
<!-- None. The existing `tencentcloud_scf_trigger_config` resource is not modified. -->

## Impact

- **New files**:
  - `tencentcloud/services/scf/resource_tc_scf_trigger.go` (resource schema + CRUD)
  - `tencentcloud/services/scf/resource_tc_scf_trigger_test.go` (gomonkey-based unit tests)
  - `tencentcloud/services/scf/resource_tc_scf_trigger.md` (user docs)
- **Modified files**:
  - `tencentcloud/services/scf/service_tencentcloud_scf.go` (add `DescribeScfTriggerById` helper)
  - `tencentcloud/provider.go` (register `tencentcloud_scf_trigger`)
  - `tencentcloud/provider.md` (register `tencentcloud_scf_trigger`)
- **APIs used** (package `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/scf/v20180416`):
  - `CreateTrigger`, `ListTriggers`, `UpdateTrigger`, `DeleteTrigger`
- **Dependencies**: Uses already-vendored `scf/v20180416` SDK; no new external dependencies.
- **Backward compatibility**: Fully additive; no existing resource schema or behavior is changed.
