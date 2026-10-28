## ADDED Requirements

### Requirement: Resource registration
The provider SHALL expose a resource type named `tencentcloud_ga2_accelerate_area` that manages a single Tencent Cloud Global Accelerator V2 acceleration region per resource block. The resource MUST be registered in `tencentcloud/provider.go` under the `ga2` namespace, adjacent to the existing GA2 entries. The resource MUST also appear under the `Global Accelerator(GA2)` Resources section of `tencentcloud/provider.md` so that `make doc` generates a website page for it.

#### Scenario: Resource type is discoverable
- **WHEN** an operator runs `terraform plan` against a configuration that references `resource "tencentcloud_ga2_accelerate_area" "<name>"`
- **THEN** Terraform resolves the type without an "unknown resource" error and shows the planned create.

#### Scenario: Provider compiles
- **WHEN** the codebase is built with `go build ./tencentcloud/...`
- **THEN** the build succeeds with no compilation or vet errors related to the new resource.

#### Scenario: Website doc generated
- **WHEN** an operator runs `make doc`
- **THEN** `website/docs/r/ga2_accelerate_area.html.markdown` is generated from the resource Schema/Description and the example markdown file.

### Requirement: Schema mirrors `CreateAccelerateAreas`
The resource schema SHALL expose every input parameter accepted by the `CreateAccelerateAreas` API (the single `AcceleratorAreas` element flattened onto the resource) as a top-level attribute, with no renaming or merging:
- `global_accelerator_id` (string, required, **ForceNew**): the global accelerator instance ID.
- `accelerate_region` (string, required, **ForceNew**): the acceleration region; serves as the natural key used to resolve `AcceleratorAreaId` after Create.
- `bandwidth` (int, optional+computed): acceleration bandwidth.
- `isp_type` (string, optional+computed): ISP type, one of `BGP` / `三网` / `精品`, default `BGP`.
- `ip_version` (string, optional+computed): IP version; only `IPv4` supported, default `IPv4`.
- `ip_address` (set of string, optional+computed): bound IP addresses.

The resource SHALL additionally expose the following read-only attributes hydrated from the `DescribeAccelerateAreas` response item:
- `accelerator_area_id` (string, computed) — also stored as the second segment of `d.Id()`.
- `ip_address_info_set` (list of nested blocks, computed) — each block exposes `ip_address` (string, computed) and `isp_type` (string, computed).

#### Scenario: All required SDK input fields are present
- **WHEN** a developer inspects the resource schema
- **THEN** every field carried by the `CreateAccelerateAreas` payload (`GlobalAcceleratorId`, and the `AcceleratorAreas` element fields `AccelerateRegion`, `Bandwidth`, `IspType`, `IpVersion`, `IpAddress`) appears in the schema with semantically equivalent typing.

#### Scenario: No undocumented schema fields
- **WHEN** a developer inspects the resource schema
- **THEN** there are no input fields beyond those listed above; no derived flags or synthetic toggles are introduced. Fields that exist only in the describe response (`AcceleratorAreaId`, `IpAddressInfoSet`) are Computed-only.

### Requirement: Resource ID
The resource ID SHALL be the composite `<GlobalAcceleratorId><FILED_SP><AcceleratorAreaId>`, using the project-standard separator `tccommon.FILED_SP`. The resource SHALL support `terraform import` using the composite ID.

#### Scenario: Create sets the composite ID
- **WHEN** `CreateAccelerateAreas` succeeds, the polled task transitions to `SUCCESS`, and the `AcceleratorAreaId` is resolved by region
- **THEN** the resource calls `d.SetId(strings.Join([]string{gaId, areaId}, tccommon.FILED_SP))`.

#### Scenario: Import by composite ID
- **WHEN** an operator runs `terraform import tencentcloud_ga2_accelerate_area.x ga-xxxxxxxx#area-yyyyyyy`
- **THEN** the resource state is hydrated from `DescribeAccelerateAreas` using the parsed `(gaId, areaId)` pair.

#### Scenario: Malformed import ID
- **WHEN** the import ID does not contain exactly one `tccommon.FILED_SP` separator, or has empty components
- **THEN** the resource returns a descriptive error before any SDK call.

### Requirement: Async create with task polling and region-based ID resolution
On Create, the resource SHALL invoke `CreateAccelerateAreas` with `GlobalAcceleratorId` and a single-element `AcceleratorAreas` list built from the schema, capture the returned `TaskId`, and poll `DescribeTaskResult` via `Ga2Service.WaitForGa2TaskFinish(ctx, taskId, timeout)` until `Status == "SUCCESS"` or the user-supplied `Timeouts.Create` (default **5 minutes**) elapses. Because `CreateAccelerateAreas` does NOT return `AcceleratorAreaId`, after the task succeeds the resource SHALL call `Ga2Service.DescribeGa2AccelerateAreaByRegion(ctx, gaId, region)` and read the server-generated `AcceleratorAreaId` to construct the composite ID.

#### Scenario: Successful async create
- **WHEN** `CreateAccelerateAreas` succeeds, the polled task transitions to `SUCCESS` within the timeout, and the area is resolved by region
- **THEN** the resource sets the composite ID, invokes Read, and returns no error.

#### Scenario: Async create timeout
- **WHEN** the task does not reach `SUCCESS` before the configured `Timeouts.Create`
- **THEN** the resource returns an error containing the task ID and last observed status.

#### Scenario: Empty TaskId
- **WHEN** `CreateAccelerateAreas` returns a nil or empty `TaskId`
- **THEN** the resource returns an explicit error rather than dereferencing the nil pointer.

#### Scenario: Area not found after successful task
- **WHEN** the task transitions to `SUCCESS` but `DescribeGa2AccelerateAreaByRegion` returns `(nil, nil)` for `(gaId, region)`
- **THEN** the resource returns an explicit error indicating the acceleration region could not be resolved.

### Requirement: Read with retry and pagination
On Read, the resource SHALL call `Ga2Service.DescribeGa2AccelerateAreaById(ctx, gaId, areaId)`, which:
- Wraps the SDK call `DescribeAccelerateAreasWithContext` in `resource.Retry(tccommon.ReadRetryTimeout, ...)`.
- Sets `request.GlobalAcceleratorId` once outside the loop (the API has no per-area filter slot).
- Iterates pages with `Limit=100` (the documented maximum), constructing the request object once **outside** the loop and only mutating `Offset` / `Limit` per iteration.
- Strict-equals on `*item.AcceleratorAreaId == areaId` before returning.
- Returns `(nil, nil)` when no matching area exists.

When the helper returns `(nil, nil)`, the resource SHALL clear `d.SetId("")` and log a `[WARN]` line indicating the area may have been deleted out of band.

#### Scenario: Area present
- **WHEN** the helper finds a matching `AcceleratorAreas` item
- **THEN** the resource populates all schema fields (input + computed, including the nested `ip_address_info_set`) from the response, guarding every assignment against nil pointers.

#### Scenario: Area removed externally
- **WHEN** the helper returns `(nil, nil)` (no matching area)
- **THEN** the resource calls `d.SetId("")` and returns no error.

#### Scenario: Pagination request reuse
- **WHEN** the helper paginates through more than one page
- **THEN** a single `DescribeAccelerateAreasRequest` instance is reused across pages, with only `Offset` and `Limit` mutated.

### Requirement: Update path
The Update function SHALL:
- Build a single-element `ModifyAccelerateAreasRequest.AcceleratorAreas` list that always carries `AcceleratorAreaId` (from the parsed composite ID) and `AccelerateRegion` (from state), plus the current values of the mutable fields (`Bandwidth`, `IspType`, `IpVersion`, `IpAddress`).
- Always include `GlobalAcceleratorId` (mandatory by the API).
- Wrap the SDK call in `resource.Retry(tccommon.WriteRetryTimeout, ...)`.
- Wait for the returned `TaskId` via `WaitForGa2TaskFinish(ctx, taskId, d.Timeout(schema.TimeoutUpdate))`.
- Skip the Modify call entirely if no Modify-supported field (`bandwidth`, `isp_type`, `ip_version`, `ip_address`) changed.

#### Scenario: Mutable field change
- **WHEN** `bandwidth`, `isp_type`, `ip_version`, or `ip_address` changes
- **THEN** `ModifyAccelerateAreas` is called with the fully-populated single-element list, awaited via the task helper.

#### Scenario: No-op update
- **WHEN** `terraform apply` runs but no Modify-supported field has changed
- **THEN** the resource skips `ModifyAccelerateAreas` and immediately invokes Read.

### Requirement: Async delete
On Delete, the resource SHALL call `DeleteAccelerateAreas` with `GlobalAcceleratorId` and `AcceleratorAreaIds=[areaId]` (parsed from the composite ID), capture the returned `TaskId`, and poll the task to completion via `WaitForGa2TaskFinish(ctx, taskId, d.Timeout(schema.TimeoutDelete))` (default **5 minutes**).

#### Scenario: Successful async delete
- **WHEN** the delete task transitions to `SUCCESS`
- **THEN** the resource returns no error and Terraform marks the resource as destroyed.

#### Scenario: Empty TaskId on delete
- **WHEN** `DeleteAccelerateAreas` returns a nil response or nil `TaskId`
- **THEN** the resource returns an explicit error rather than dereferencing the nil pointer.

### Requirement: Retry coverage
Every SDK call (`CreateAccelerateAreasWithContext`, `DescribeAccelerateAreasWithContext`, `ModifyAccelerateAreasWithContext`, `DeleteAccelerateAreasWithContext`) SHALL be invoked from inside a `resource.Retry` block. The retry budget is `tccommon.WriteRetryTimeout` for write operations and `tccommon.ReadRetryTimeout` for read operations.

#### Scenario: Transient SDK error
- **WHEN** any SDK call returns a transient TencentCloud SDK error
- **THEN** the call is retried via `tccommon.RetryError(e)` until it succeeds or the retry budget is exhausted.

#### Scenario: Nil response defense
- **WHEN** an SDK call returns a nil `result` / `result.Response` (or a nil critical sub-field such as `TaskId`)
- **THEN** the wrapper returns `resource.NonRetryableError` with a descriptive message rather than dereferencing the nil pointer.

### Requirement: Logging conventions
The resource SHALL emit:
- `defer tccommon.LogElapsed("resource.tencentcloud_ga2_accelerate_area.<op>")()` at the top of every CRUD function.
- `defer tccommon.InconsistentCheck(d, meta)()` at the top of every CRUD function.
- A `[DEBUG]` line per SDK invocation containing the request action, request body, and response body (matching the existing GA2 log format).
- A `[CRITAL]%s ... failed, reason:%+v` line on every retry-block failure.
- A `[WARN]` line when the resource is detected as deleted out of band during Read.

#### Scenario: Standard log lines emitted
- **WHEN** any CRUD operation runs
- **THEN** the operation's elapsed time is logged via `tccommon.LogElapsed` and inconsistency is checked via `tccommon.InconsistentCheck`.

### Requirement: Documentation and tests
The change SHALL include:
- A markdown document at `tencentcloud/services/ga2/resource_tc_ga2_accelerate_area.md` containing a self-contained `resource "tencentcloud_ga2_accelerate_area" "..." { ... }` example and a `terraform import` example using the composite ID. Filename pattern follows `resource_tc_config_compliance_pack.md`.
- An acceptance-test file at `tencentcloud/services/ga2/resource_tc_ga2_accelerate_area_test.go` exercising at minimum: create, basic update (e.g. `bandwidth`), import, and destroy. Filename pattern follows `resource_tc_config_compliance_pack_test.go`.

#### Scenario: Documentation present
- **WHEN** the change is merged
- **THEN** the markdown documentation file exists and contains both an HCL example and an `import` example using the composite ID.

#### Scenario: Test file present
- **WHEN** the change is merged
- **THEN** the `_test.go` file declares `TestAccTencentCloudGa2AccelerateAreaResource_basic` (or equivalent) using `resource.Test` with at least two `Steps` (create + update) plus an `ImportState` verification step.

### Requirement: SDK constraint
The implementation SHALL NOT modify any file under `vendor/github.com/tencentcloud/tencentcloud-sdk-go/`. If a required API is missing from the vendored SDK, the implementer MUST halt and request an SDK upgrade rather than authoring or patching SDK source.

#### Scenario: Vendored SDK is sufficient
- **WHEN** the implementation begins
- **THEN** the four AccelerateAreas APIs (`CreateAccelerateAreas`, `DescribeAccelerateAreas`, `ModifyAccelerateAreas`, `DeleteAccelerateAreas`) plus `DescribeTaskResult` are confirmed present under `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ga2/v20250115/` before any code is written.
