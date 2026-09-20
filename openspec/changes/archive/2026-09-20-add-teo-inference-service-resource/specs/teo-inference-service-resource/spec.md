## ADDED Requirements

### Requirement: Resource registration
The provider SHALL expose a resource type named `tencentcloud_teo_inference_service` that manages a single Tencent Cloud EdgeOne inference service per resource block. The resource MUST be registered in `tencentcloud/provider.go` under the `teo` namespace.

#### Scenario: Resource type is discoverable
- **WHEN** an operator runs `terraform plan` against a configuration that references `resource "tencentcloud_teo_inference_service" "<name>"`
- **THEN** Terraform resolves the type without an "unknown resource" error and shows the planned create.

#### Scenario: Provider compiles
- **WHEN** the codebase is built with `go build ./tencentcloud/...`
- **THEN** the build succeeds with no compilation errors related to the new resource.

### Requirement: Schema mirrors `CreateInferenceService`
The resource schema SHALL expose every input parameter accepted by the `CreateInferenceService` API as a top-level or nested attribute, with no renaming or merging:
- `zone_id` (string, Required, **ForceNew**): site ID.
- `name` (string, Required, **ForceNew**): inference service name; not accepted by `ModifyInferenceService`, hence ForceNew.
- `listen_port` (int, Required): model service listen port, 1-65535.
- `containers` (list, Required): container config list; each element contains `image_type` (string, Required), `tcr_repository_config` (list, Optional) with `tcr_type`/`image` (Required) and `registry_id`/`region_name` (Optional), `startup_command` (string, Optional), `environment_variables` (list, Optional) with `key` (Required) and `value` (Optional).
- `resource_config` (list, Required): resource config; contains `scaling_mode` (string, Required), `hardware_spec` (string, Optional), `hardware_spec_id` (string, Optional), `hardware_config` (list, Optional) with `gpu_num`/`cpu_num`/`mem_size`/`disk_size` (Optional), `auto_scaling_config` (list, Optional) with `min_instance_count` (int, Required) and `scaling_policies` (list, Optional), `manual_instance_config` (list, Optional) with `fixed_instance_count` (int, Required), `concurrency` (int, Optional).
- `affinity_config` (list, Optional): affinity config; contains `switch` (string, Required), `affinity_mode` (string, Optional), `session_id_affinity_config` (list, Optional) with `source`/`header_name` (Optional). Note: `DescribeInferenceServices` response does not include `AffinityConfig`, so this field is write-only and not read back.
- `request_paths` (list of string, Optional): request path list.
- `description` (string, Optional): description.

The resource SHALL additionally expose the following read-only Computed attributes hydrated from `DescribeInferenceServices`:
- `service_id` (string, Computed) — also part of the composite resource ID.
- `status` (string, Computed)
- `scaling_status` (string, Computed)
- `current_instance_count` (int, Computed)
- `inference_url` (string, Computed)
- `create_time` (string, Computed)
- `update_time` (string, Computed)

The `scaling_policies` list elements SHALL contain `policy_name` (string, Required), `policy_type` (string, Required), and `scheduled_scaling_policy` (list, Optional). The `scheduled_scaling_policy` element SHALL contain `scheduled_actions` (list, Required) with `cron_expression` (string, Required) and `min_instance_count` (int, Required), plus `effective_range` (list, Required) with `effective_type` (string, Required) and `start_date`/`end_date` (string, Optional), and `time_zone` (string, Optional).

#### Scenario: Required SDK fields are present
- **WHEN** a developer inspects the resource schema
- **THEN** every field declared in `CreateInferenceServiceRequestParams` (ZoneId, Name, ListenPort, Containers, ResourceConfig, AffinityConfig, RequestPaths, Description) appears in the schema with semantically equivalent typing.

#### Scenario: Computed fields exposed
- **WHEN** a developer inspects the resource schema
- **THEN** the read-only fields from `InferenceService` (Status, ScalingStatus, CurrentInstanceCount, InferenceURL, CreateTime, UpdateTime, ServiceId) appear as Computed attributes.

#### Scenario: No undocumented schema fields
- **WHEN** a developer inspects the resource schema
- **THEN** there are no fields beyond those sourced from the SDK or required by the composite-ID/Computed conventions; no derived flags or synthetic toggles are introduced.

### Requirement: Composite resource ID
The resource ID SHALL be the composite `<zoneId><tccommon.FILED_SP><serviceId>` returned after `CreateInferenceService` succeeds. The resource SHALL support `terraform import` using this composite ID, and the Read/Update/Delete functions SHALL split `d.Id()` to recover `zoneId` and `serviceId`.

#### Scenario: Create sets the composite ID
- **WHEN** `CreateInferenceService` returns a non-empty `ServiceId`
- **THEN** the resource calls `d.SetId(strings.Join([]string{zoneId, serviceId}, tccommon.FILED_SP))`.

#### Scenario: Import by composite ID
- **WHEN** an operator runs `terraform import tencentcloud_teo_inference_service.x zone-123#is-456`
- **THEN** the resource splits the ID into `zoneId=zone-123` and `serviceId=is-456` and hydrates state from `DescribeInferenceServices`.

#### Scenario: Broken ID is rejected
- **WHEN** `d.Id()` does not split into exactly 2 segments
- **THEN** the Read/Update/Delete function returns an `id is broken` error.

### Requirement: Create with retry and nil defense
On Create, the resource SHALL build a `CreateInferenceServiceRequest` from the schema (translating nested lists into the corresponding SDK structs), invoke `CreateInferenceServiceWithContext` inside `resource.Retry(tccommon.WriteRetryTimeout, ...)`, and defend against nil `result` / `result.Response` by returning `resource.NonRetryableError`. After the retry block, the resource SHALL check that `response.Response.ServiceId` is non-nil and non-empty (printing `logId` and `d.Id()` first for diagnostics); if empty, return a `NonRetryableError`. Finally set the composite ID and call Read.

#### Scenario: Successful create
- **WHEN** `CreateInferenceService` returns a non-nil `ServiceId`
- **THEN** the resource sets the composite ID, invokes Read, and returns no error.

#### Scenario: Empty ServiceId rejected
- **WHEN** `CreateInferenceService` returns a nil or empty `ServiceId`
- **THEN** the resource returns an explicit error and does NOT call `d.SetId`.

#### Scenario: Nil response defense
- **WHEN** the SDK call returns a nil `Response`
- **THEN** the retry wrapper returns `resource.NonRetryableError` with a descriptive message.

### Requirement: Read with retry and pagination
On Read, the resource SHALL split `d.Id()` into `zoneId` and `serviceId`, then call `TeoService.DescribeTeoInferenceServiceById(ctx, zoneId, serviceId)`, which:
- Wraps the SDK call `DescribeInferenceServicesWithContext` in `resource.Retry(tccommon.ReadRetryTimeout, ...)`.
- Filters by `Filters=[{Name:"service-id", Values:[serviceId]}]` and `ZoneId`.
- Iterates pages with `Limit=200` (the documented maximum), constructing the request object once **outside** the loop and only mutating `Offset` per iteration.
- Strict-equals `*item.ServiceId == serviceId` before returning.
- Returns `(nil, nil)` when the instance is absent.

When the helper returns `(nil, nil)`, the resource SHALL first print `log.Printf("[CRUD] inference_service id=%s", d.Id())` to preserve context, then call `d.SetId("")` and return no error. When present, the resource SHALL populate every schema field (input + computed) with nil checks before each `_ = d.Set(...)`.

#### Scenario: Resource present
- **WHEN** the helper finds a matching `InferenceService`
- **THEN** the resource populates all schema fields (input + computed) from the response.

#### Scenario: Resource removed externally
- **WHEN** the helper returns `(nil, nil)`
- **THEN** the resource logs the id, calls `d.SetId("")`, and returns no error.

#### Scenario: Pagination request reuse
- **WHEN** the helper paginates through more than one page
- **THEN** a single `DescribeInferenceServicesRequest` instance is reused across pages, with only `Offset` and `Limit` mutated.

### Requirement: Update path with ForModify structs
The Update function SHALL, when any updatable field changes (`listen_port`, `containers`, `resource_config`, `affinity_config`, `request_paths`, `description`), build a `ModifyInferenceServiceRequest` from the schema and call `ModifyInferenceServiceWithContext` inside `resource.Retry(tccommon.WriteRetryTimeout, ...)`.

The Update function SHALL:
- Use the `InferenceContainerConfigForModify` and `InferenceResourceConfigForModify` structs (not the Create variants), reflecting the SDK's Modify request types.
- NOT include `hardware_spec`, `hardware_spec_id`, or `gpu_num` in the Modify request body, because `InferenceResourceConfigForModify` does not accept `HardwareSpec`/`HardwareSpecId` and `InferenceHardwareConfigForModify` does not contain `GPUNum`.
- Include `zone_id` and `service_id` (parsed from `d.Id()`) as required inputs.
- After success, call Read to refresh state.

#### Scenario: Updatable field change
- **WHEN** `description` or other updatable field changes
- **THEN** `ModifyInferenceService` is called with the changed fields plus `ZoneId` and `ServiceId`.

#### Scenario: ForceNew field change triggers recreate
- **WHEN** `zone_id` or `name` changes
- **THEN** Terraform destroys and recreates the resource rather than calling Modify.

### Requirement: Delete via OperateInferenceService
On Delete, the resource SHALL call `OperateInferenceServiceWithContext` inside `resource.Retry(tccommon.WriteRetryTimeout, ...)` with `ZoneId`, `ServiceId` (parsed from `d.Id()`), and `Operation="Delete"`. The operation is synchronous (response contains only `RequestId`); no task polling is required.

#### Scenario: Successful delete
- **WHEN** `OperateInferenceService` with `Operation=Delete` succeeds
- **THEN** the resource returns no error and Terraform marks the resource as destroyed.

#### Scenario: Already deleted returns success
- **WHEN** a subsequent delete attempt fails with `ResourceNotFound.InferenceService`
- **THEN** the resource treats it as already deleted and returns no error (via `tccommon.RetryError` classification).

### Requirement: Retry coverage
Every SDK call (`CreateInferenceServiceWithContext`, `DescribeInferenceServicesWithContext`, `ModifyInferenceServiceWithContext`, `OperateInferenceServiceWithContext`) SHALL be invoked from inside a `resource.Retry` block. The retry budget is `tccommon.WriteRetryTimeout` for write operations and `tccommon.ReadRetryTimeout` for read operations. Retryable errors SHALL be wrapped via `tccommon.RetryError(e)`; the retry block SHALL only contain the API call, with ID-setting and other success actions placed outside the block.

#### Scenario: Transient SDK error retried
- **WHEN** any SDK call returns a transient TencentCloud SDK error
- **THEN** the call is retried via `tccommon.RetryError(e)` until it succeeds or the retry budget is exhausted.

#### Scenario: Retry block contains only API call
- **WHEN** a developer inspects the Create/Update/Delete retry block
- **THEN** it contains only the SDK invocation and nil-defense; `d.SetId` and other success-side actions are outside the retry block.

### Requirement: Logging conventions
The resource SHALL emit:
- `defer tccommon.LogElapsed("resource.tencentcloud_teo_inference_service.<op>")()` at the top of every CRUD function.
- `defer tccommon.InconsistentCheck(d, meta)()` at the top of every CRUD function.
- A `[DEBUG]` line per SDK invocation containing the request action, request body, and response body (matching `tencentcloud_igtm_strategy` log format).
- A `[CRITAL]%s ... failed, reason:%+v` line on every retry-block failure, using the resource name `inference_service` (not "该资源").
- A `[CRUD] inference_service id=%s` line before `d.SetId("")` when the resource is detected as deleted out of band during Read.

#### Scenario: Standard log lines emitted
- **WHEN** any CRUD operation runs
- **THEN** elapsed time is logged via `tccommon.LogElapsed` and inconsistency is checked via `tccommon.InconsistentCheck`, and the resource name `inference_service` is used in error/log messages.

### Requirement: Documentation and tests
The change SHALL include:
- A markdown document at `tencentcloud/services/teo/resource_tc_teo_inference_service.md` containing a self-contained `terraform { ... } resource "tencentcloud_teo_inference_service" "..." { ... }` example and a `terraform import` example showing the composite ID (`zoneId#serviceId`). The one-line description must mention TEO. Do NOT include `Argument Reference` / `Attribute Reference` sections (auto-generated). Do NOT hand-edit `website/` files (generated by `make doc`).
- A unit-test file at `tencentcloud/services/teo/resource_tc_teo_inference_service_test.go` using **gomonkey** to mock the cloud API (NOT the terraform test suite), covering Create/Read/Update/Delete business logic paths.

#### Scenario: Documentation present
- **WHEN** the change is merged
- **THEN** the markdown documentation file exists and contains both an HCL example and an `import` example using the composite ID.

#### Scenario: Unit test uses gomonkey
- **WHEN** a developer inspects the test file
- **THEN** it uses `gomonkey` to mock SDK methods and does NOT use `resource.Test`/`TF_ACC` acceptance-test scaffolding.

### Requirement: SDK constraint
The implementation SHALL NOT modify any file under `vendor/github.com/tencentcloud/tencentcloud-sdk-go/`. The four required APIs (`CreateInferenceService`, `DescribeInferenceServices`, `ModifyInferenceService`, `OperateInferenceService`) are confirmed present in the vendored SDK before any code is written.

#### Scenario: Vendored SDK is sufficient
- **WHEN** the implementation begins
- **THEN** the four InferenceService APIs are confirmed present under `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901/` (client.go + models.go) before any code is written.