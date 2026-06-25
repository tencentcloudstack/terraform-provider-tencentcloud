## ADDED Requirements

### Requirement: Resource registration
The provider SHALL expose a resource type named `tencentcloud_ga2_global_accelerator` that manages a single Tencent Cloud Global Accelerator V2 instance per resource block. The resource MUST be registered in `tencentcloud/provider.go` under the `ga2` namespace.

#### Scenario: Resource type is discoverable
- **WHEN** an operator runs `terraform plan` against a configuration that references `resource "tencentcloud_ga2_global_accelerator" "<name>"`
- **THEN** Terraform resolves the type without an "unknown resource" error and shows the planned create.

#### Scenario: Provider compiles
- **WHEN** the codebase is built with `go build ./tencentcloud/...`
- **THEN** the build succeeds with no compilation or vet errors related to the new resource.

### Requirement: Schema mirrors `CreateGlobalAccelerator`
The resource schema SHALL expose every input parameter accepted by the `CreateGlobalAccelerator` API as a top-level attribute, with no renaming or merging:
- `name` (string, optional+computed): instance name, ≤60 bytes.
- `instance_charge_type` (string, optional+computed, **ForceNew**): billing mode (`POSTPAID` is the only currently supported value).
- `description` (string, optional+computed): instance description, ≤100 bytes.
- `cross_border_type` (string, optional+computed): cross-border tier (`HighQuality` / `Unicom`).
- `cross_border_promise_flag` (bool, optional+computed): cross-border-service commitment flag.
- `tags` (map[string]string, optional): resource tags.

The resource SHALL additionally expose the following read-only attributes hydrated from `DescribeGlobalAccelerators` response:
- `global_accelerator_id` (string, computed) — also assigned as the resource ID.
- `state` (string, computed)
- `status` (string, computed)
- `cname` (string, computed)
- `ddos_id` (string, computed)
- `create_time` (string, computed)
- `listener_counts` (int, computed)
- `accelerator_area_counts` (int, computed)

#### Scenario: Required SDK fields are present
- **WHEN** a developer inspects the resource schema
- **THEN** every field declared in `ga2v20250115.CreateGlobalAcceleratorRequestParams` (Name, InstanceChargeType, Description, CrossBorderType, CrossBorderPromiseFlag, Tags) appears in the schema with semantically equivalent typing.

#### Scenario: No undocumented schema fields
- **WHEN** a developer inspects the resource schema
- **THEN** there are no fields beyond those listed above; no derived flags or synthetic toggles are introduced.

### Requirement: Resource ID
The resource ID SHALL be the bare `GlobalAcceleratorId` returned by `CreateGlobalAccelerator`. The resource SHALL support `terraform import` using just that ID.

#### Scenario: Create sets the ID
- **WHEN** `CreateGlobalAccelerator` returns a non-nil `GlobalAcceleratorId`
- **THEN** the resource calls `d.SetId(<gaId>)` after the async task completes successfully.

#### Scenario: Import by ID
- **WHEN** an operator runs `terraform import tencentcloud_ga2_global_accelerator.x ga-xxxxxxxx`
- **THEN** the resource state is hydrated from `DescribeGlobalAccelerators` using the imported ID, with no manual composite-ID parsing required.

### Requirement: Async create with task polling
On Create, the resource SHALL invoke `CreateGlobalAccelerator`, capture the returned `TaskId` and `GlobalAcceleratorId`, and poll `DescribeTaskResult` via the existing `Ga2Service.WaitForGa2TaskFinish(ctx, taskId, timeout)` helper until `Status == "SUCCESS"` or the user-supplied `Timeouts.Create` (default **5 minutes**) elapses.

#### Scenario: Successful async create
- **WHEN** `CreateGlobalAccelerator` succeeds and the polled task transitions to `SUCCESS` within the timeout
- **THEN** the resource sets the ID, invokes Read, and returns no error.

#### Scenario: Async create timeout
- **WHEN** the task does not reach `SUCCESS` before the configured `Timeouts.Create`
- **THEN** the resource returns an error containing the task ID and last observed status.

#### Scenario: Empty TaskId
- **WHEN** `CreateGlobalAccelerator` returns a nil `TaskId`
- **THEN** the resource returns an explicit error rather than silently calling the polling helper.

### Requirement: Read with retry and pagination
On Read, the resource SHALL call `Ga2Service.DescribeGa2GlobalAcceleratorById(ctx, gaId)`, which:
- Wraps the SDK call `DescribeGlobalAcceleratorsWithContext` in `resource.Retry(tccommon.ReadRetryTimeout, ...)`.
- Filters by `Filters=[{Name:"global-accelerator-id", Values:[gaId]}]`.
- Iterates pages with `Limit=100` (the documented maximum), constructing the request object once **outside** the loop and only mutating `Offset` / `Limit` per iteration.
- Strict-equals `*item.GlobalAcceleratorId == gaId` before returning.
- Returns `(nil, nil)` when the instance is absent.

When the helper returns `(nil, nil)`, the resource SHALL clear `d.SetId("")` and log a `[WARN]` line indicating the instance may have been deleted out of band.

#### Scenario: Resource present
- **WHEN** the helper finds a matching `GlobalAcceleratorSet`
- **THEN** the resource populates all schema fields (input + computed + tags) from the response.

#### Scenario: Resource removed externally
- **WHEN** the helper returns `(nil, nil)` (no matching instance)
- **THEN** the resource calls `d.SetId("")` and returns no error.

#### Scenario: Pagination request reuse
- **WHEN** the helper paginates through more than one page
- **THEN** a single `DescribeGlobalAcceleratorsRequest` instance is reused across pages, with only `Offset` and `Limit` mutated.

### Requirement: Update path
The Update function SHALL:
- Call `ModifyGlobalAccelerator` (wrapped in `resource.Retry(tccommon.WriteRetryTimeout, ...)`) **only when** at least one of `name` / `description` / `cross_border_type` / `cross_border_promise_flag` has changed; the SDK request MUST include only the changed fields plus the required `GlobalAcceleratorId`.
- Wait for the returned `TaskId` via `WaitForGa2TaskFinish(ctx, taskId, d.Timeout(schema.TimeoutUpdate))`.
- When `tags` has changed, reconcile drift via `svctag.NewTagService(...).ModifyTags(ctx, "qcs::ga2:<region>:uin/:globalAccelerator/<gaId>", replaceTags, deleteTags)` after (or independently of) the Modify call.

#### Scenario: Tags-only change
- **WHEN** only the `tags` map changes
- **THEN** the resource skips `ModifyGlobalAccelerator` entirely and reconciles tags via the tag service.

#### Scenario: Mixed change
- **WHEN** both `name` and `tags` change
- **THEN** `ModifyGlobalAccelerator` is called with `Name` populated, awaited via the task helper, and tags are reconciled afterwards.

### Requirement: Async delete
On Delete, the resource SHALL call `DeleteGlobalAccelerator`, capture the returned `TaskId`, and poll the task to completion via `WaitForGa2TaskFinish(ctx, taskId, d.Timeout(schema.TimeoutDelete))` (default **5 minutes**).

#### Scenario: Successful async delete
- **WHEN** the delete task transitions to `SUCCESS`
- **THEN** the resource returns no error and Terraform marks the resource as destroyed.

### Requirement: Retry coverage
Every SDK call (`CreateGlobalAcceleratorWithContext`, `DescribeGlobalAcceleratorsWithContext`, `ModifyGlobalAcceleratorWithContext`, `DeleteGlobalAcceleratorWithContext`) SHALL be invoked from inside a `resource.Retry` block. The retry budget is `tccommon.WriteRetryTimeout` for write operations and `tccommon.ReadRetryTimeout` for read operations.

#### Scenario: Transient SDK error
- **WHEN** any of the four SDK calls returns a transient TencentCloud SDK error
- **THEN** the call is retried via `tccommon.RetryError(e)` until it succeeds or the retry budget is exhausted.

#### Scenario: Nil response defense
- **WHEN** an SDK call returns a nil `Response` (or a nil critical sub-field such as `TaskId`)
- **THEN** the wrapper returns `resource.NonRetryableError` with a descriptive message rather than dereferencing the nil pointer.

### Requirement: Logging conventions
The resource SHALL emit:
- `defer tccommon.LogElapsed("resource.tencentcloud_ga2_global_accelerator.<op>")()` at the top of every CRUD function.
- `defer tccommon.InconsistentCheck(d, meta)()` at the top of every CRUD function.
- A `[DEBUG]` line per SDK invocation containing the request action, request body, and response body (matching the existing `tencentcloud_ga2_endpoint_group` log format).
- A `[CRITAL]%s ... failed, reason:%+v` line on every retry-block failure.
- A `[WARN]` line when the resource is detected as deleted out of band during Read.

#### Scenario: Standard log lines emitted
- **WHEN** any CRUD operation runs
- **THEN** the operation's elapsed time is logged via `tccommon.LogElapsed` and inconsistency is checked via `tccommon.InconsistentCheck`.

### Requirement: Documentation and tests
The change SHALL include:
- A markdown document at `tencentcloud/services/ga2/resource_tc_ga2_global_accelerator.md` containing a self-contained `terraform { ... } resource "tencentcloud_ga2_global_accelerator" "..." { ... }` example and a `terraform import` example. Filename pattern follows `resource_tc_config_compliance_pack.md`.
- An acceptance-test file at `tencentcloud/services/ga2/resource_tc_ga2_global_accelerator_test.go` exercising at minimum: create, basic update (e.g. `description`), import, and destroy. Filename pattern follows `resource_tc_config_compliance_pack_test.go`.

#### Scenario: Documentation present
- **WHEN** the change is merged
- **THEN** the markdown documentation file exists and contains both an HCL example and an `import` example using the bare `GlobalAcceleratorId`.

#### Scenario: Test file present
- **WHEN** the change is merged
- **THEN** the `_test.go` file declares a `TestAccTencentCloudGa2GlobalAcceleratorResource_basic` (or equivalent) test using `resource.Test` with a `CheckDestroy` and at least two `Steps` (create + update + import).

### Requirement: SDK constraint
The implementation SHALL NOT modify any file under `vendor/github.com/tencentcloud/tencentcloud-sdk-go/`. If a required API is missing from the vendored SDK, the implementer MUST halt and request an SDK upgrade rather than authoring or patching SDK source.

#### Scenario: Vendored SDK is sufficient
- **WHEN** the implementation begins
- **THEN** the four GlobalAccelerator APIs plus `DescribeTaskResult` are confirmed present under `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ga2/v20250115/` before any code is written.
