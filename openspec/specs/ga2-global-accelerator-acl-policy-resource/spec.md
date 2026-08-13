# ga2-global-accelerator-acl-policy-resource Specification

## Purpose
TBD - created by archiving change add-ga2-global-accelerator-acl-policy-resource. Update Purpose after archive.
## Requirements
### Requirement: Resource registration
The provider SHALL expose a resource type named `tencentcloud_ga2_global_accelerator_acl_policy` that manages a single Tencent Cloud Global Accelerator V2 (GA2) access-control policy per resource block. The resource MUST be registered in `tencentcloud/provider.go` under the `ga2` namespace, adjacent to the existing `ga2` resource entries.

#### Scenario: Resource type is discoverable
- **WHEN** an operator runs `terraform plan` against a configuration that references `resource "tencentcloud_ga2_global_accelerator_acl_policy" "<name>"`
- **THEN** Terraform resolves the type without an "unknown resource" error and shows the planned create.

#### Scenario: Provider compiles
- **WHEN** the codebase is built with `go build ./tencentcloud/...`
- **THEN** the build succeeds with no compilation errors related to the new resource.

### Requirement: Schema mirrors the cloud APIs
The resource schema SHALL expose every input parameter accepted by the `CreateGlobalAcceleratorAclPolicy` / `ModifyGlobalAcceleratorAclPolicy` APIs as top-level attributes, with no renaming or merging:
- `global_accelerator_id` (string, Required, **ForceNew**): the parent global accelerator instance ID; required by all CRUD calls and immutable.
- `default_action` (string, Required, **ForceNew**): default traffic action. Enumerated values: `ACCEPT` (allow all traffic by default), `DROP` (deny all traffic by default). Accepted only by `CreateGlobalAcceleratorAclPolicy`; `ModifyGlobalAcceleratorAclPolicy` has no slot for it, so it is ForceNew.
- `status` (string, Optional, Computed): ACL policy state. Enumerated values: `OPEN` (enabled), `CLOSE` (disabled). Updatable in place via `ModifyGlobalAcceleratorAclPolicy`.

The resource SHALL additionally expose the following read-only (Computed) attributes hydrated from the cloud APIs:
- `global_accelerator_acl_policy_id` (string, Computed): the ACL policy ID returned by `CreateGlobalAcceleratorAclPolicy` and re-read from `DescribeGlobalAcceleratorAclPolicies`.
- `task_id` (string, Computed): the async task ID returned by the most recent write call (Create / Modify / Delete).

The resource SHALL declare a `Timeouts` block with `Create` / `Update` / `Delete` defaults of 5 minutes.

#### Scenario: Required fields are enforced
- **WHEN** a user creates a `tencentcloud_ga2_global_accelerator_acl_policy` resource without `global_accelerator_id`
- **THEN** Terraform SHALL report a validation error indicating the field is required.

#### Scenario: Required fields are enforced for default_action
- **WHEN** a user creates the resource without `default_action`
- **THEN** Terraform SHALL report a validation error indicating the field is required.

#### Scenario: No undocumented schema fields
- **WHEN** a developer inspects the resource schema
- **THEN** the only fields present are `global_accelerator_id`, `default_action`, `status`, `global_accelerator_acl_policy_id`, `task_id`, and the standard `Timeouts` block; no derived flags or synthetic toggles are introduced.

### Requirement: Composite resource ID
The resource ID SHALL be the composite `<global_accelerator_id>#<global_accelerator_acl_policy_id>` using `tccommon.FILED_SP` ("#") as the separator, because every Read/Update/Delete cloud API call requires both `GlobalAcceleratorId` and `GlobalAcceleratorAclPolicyId` and the Describe API has no per-policy filter. Read/Update/Delete SHALL split `d.Id()` on `tccommon.FILED_SP` to recover both IDs. The resource SHALL support `terraform import` using this composite ID.

#### Scenario: Create sets the composite ID
- **WHEN** `CreateGlobalAcceleratorAclPolicy` returns a non-nil `GlobalAcceleratorAclPolicyId` and the async task reaches `SUCCESS`
- **THEN** the resource calls `d.SetId(<global_accelerator_id> + tccommon.FILED_SP + <global_accelerator_acl_policy_id>)`.

#### Scenario: ID split recovers both components
- **WHEN** Read/Update/Delete need the cloud API IDs
- **THEN** the resource splits `d.Id()` on `tccommon.FILED_SP` into exactly two parts and uses them as `GlobalAcceleratorId` and `GlobalAcceleratorAclPolicyId` respectively; if the split does not yield exactly two parts, the resource returns a descriptive error.

#### Scenario: Import by composite ID
- **WHEN** an operator runs `terraform import tencentcloud_ga2_global_accelerator_acl_policy.x <gaId>#<policyId>`
- **THEN** the resource state is hydrated from `DescribeGlobalAcceleratorAclPolicies` using the split IDs, with no manual multi-arg parsing required.

### Requirement: Async create with task polling
On Create, the resource SHALL invoke `CreateGlobalAcceleratorAclPolicyWithContext` with `GlobalAcceleratorId` and `DefaultAction` from the schema, capture the returned `TaskId` and `GlobalAcceleratorAclPolicyId`, and poll `DescribeTaskResult` via the existing `Ga2Service.WaitForGa2TaskFinish(ctx, taskId, timeout)` helper until `Status == "SUCCESS"` or the user-supplied `Timeouts.Create` (default 5 minutes) elapses. `d.SetId(...)` SHALL happen only after the task succeeds, outside the retry block.

#### Scenario: Successful async create
- **WHEN** `CreateGlobalAcceleratorAclPolicy` succeeds and the polled task transitions to `SUCCESS` within the timeout
- **THEN** the resource sets the composite ID, invokes Read, and returns no error.

#### Scenario: Async create timeout
- **WHEN** the task does not reach `SUCCESS` before the configured `Timeouts.Create`
- **THEN** the resource returns an error containing the task ID and last observed status.

#### Scenario: Empty TaskId or PolicyId on create
- **WHEN** `CreateGlobalAcceleratorAclPolicy` returns a nil `Response`, a nil `TaskId`, or a nil/empty `GlobalAcceleratorAclPolicyId`
- **THEN** the resource returns a `NonRetryableError` with a descriptive message rather than silently proceeding to poll or set state.

### Requirement: Read with retry, pagination, and client-side match
On Read, the resource SHALL split the composite ID, then call `Ga2Service.DescribeGa2GlobalAcceleratorAclPolicyById(ctx, gaId, policyId)`, which:
- Wraps the SDK call `DescribeGlobalAcceleratorAclPoliciesWithContext` in `resource.Retry(tccommon.ReadRetryTimeout, ...)`.
- Sets `GlobalAcceleratorId` on the request; iterates pages with `Limit="200"` (the documented maximum for the `*string` Limit field) and `Offset` (`*uint64`) starting at 0, constructing the request object once **outside** the loop and only mutating `Offset` / `Limit` per iteration.
- Strict-equals `*item.GlobalAcceleratorAclPolicyId == policyId` before returning the matched `GlobalAcceleratorAclPolicies`.
- Returns `(nil, nil)` when the policy is absent.

When the helper returns `(nil, nil)`, the resource SHALL **first** log `log.Printf("[CRUD] ga2_global_accelerator_acl_policy id=%s", d.Id())` to preserve the id, **then** call `d.SetId("")` and return nil. If the resource is `IsNewResource()` at that point, the resource SHALL return an explicit error instead of clearing the ID silently.

When the helper returns a matched policy, the resource SHALL call `_ = d.Set(...)` for `global_accelerator_id`, `default_action`, `status`, `global_accelerator_acl_policy_id`, and `task_id` (only when the relevant response field is non-nil).

#### Scenario: Resource present
- **WHEN** the helper finds a matching `GlobalAcceleratorAclPolicies`
- **THEN** the resource populates all schema fields from the response, guarding each `d.Set` on the corresponding response field being non-nil.

#### Scenario: Resource removed externally
- **WHEN** the helper returns `(nil, nil)` (no matching policy)
- **THEN** the resource logs the id via `[CRUD]`, calls `d.SetId("")`, and returns no error.

#### Scenario: Pagination request reuse
- **WHEN** the helper paginates through more than one page
- **THEN** a single `DescribeGlobalAcceleratorAclPoliciesRequest` instance is reused across pages, with only `Offset` and `Limit` mutated.

### Requirement: Update path — status only
The Update function SHALL call `ModifyGlobalAcceleratorAclPolicyWithContext` (wrapped in `resource.Retry(tccommon.WriteRetryTimeout, ...)`) when `status` has changed. The SDK request MUST include `GlobalAcceleratorId`, `GlobalAcceleratorAclPolicyId` (both recovered from the composite ID), and `Status`. After the SDK call succeeds, the resource SHALL wait for the returned `TaskId` via `WaitForGa2TaskFinish(ctx, taskId, d.Timeout(schema.TimeoutUpdate))`.

`default_action` and `global_accelerator_id` are ForceNew and SHALL NOT be sent on the Modify call; Terraform handles them by destroy+recreate.

#### Scenario: Status change triggers Modify
- **WHEN** only `status` changes (e.g. `OPEN` → `CLOSE`)
- **THEN** the resource calls `ModifyGlobalAcceleratorAclPolicy` with `Status` populated, awaits the task, and re-reads.

#### Scenario: No updatable field changed
- **WHEN** no updatable field (`status`) has changed
- **THEN** the resource skips the Modify call and returns directly (or after Read).

### Requirement: Async delete
On Delete, the resource SHALL call `DeleteGlobalAcceleratorAclPolicyWithContext` (wrapped in `resource.Retry(tccommon.WriteRetryTimeout, ...)`) with `GlobalAcceleratorId` and `GlobalAcceleratorAclPolicyId` recovered from the composite ID, capture the returned `TaskId`, and poll the task to completion via `WaitForGa2TaskFinish(ctx, taskId, d.Timeout(schema.TimeoutDelete))` (default 5 minutes).

#### Scenario: Successful async delete
- **WHEN** the delete task transitions to `SUCCESS`
- **THEN** the resource returns no error and Terraform marks the resource as destroyed.

#### Scenario: Delete of already-absent policy
- **WHEN** the delete SDK call returns a `ResourceNotFound`-style error
- **THEN** the resource treats the policy as already deleted and returns no error.

### Requirement: Retry coverage and nil defenses
Every SDK call (`CreateGlobalAcceleratorAclPolicyWithContext`, `DescribeGlobalAcceleratorAclPoliciesWithContext`, `ModifyGlobalAcceleratorAclPolicyWithContext`, `DeleteGlobalAcceleratorAclPolicyWithContext`) SHALL be invoked from inside a `resource.Retry` block. The retry budget is `tccommon.WriteRetryTimeout` for write operations and `tccommon.ReadRetryTimeout` for read operations. Errors from SDK calls SHALL be routed through `tccommon.RetryError(e)`. `d.SetId(...)` and other success-path state mutations SHALL occur outside and after the retry block, never inside the retry closure.

#### Scenario: Transient SDK error
- **WHEN** any of the four SDK calls returns a transient TencentCloud SDK error
- **THEN** the call is retried via `tccommon.RetryError(e)` until it succeeds or the retry budget is exhausted.

#### Scenario: Nil response defense on write
- **WHEN** a write SDK call returns a nil `Response` (or a nil critical sub-field such as `TaskId` on Create, or `TaskId` on Modify/Delete)
- **THEN** the wrapper returns `resource.NonRetryableError` with a descriptive message rather than dereferencing the nil pointer.

#### Scenario: State mutation outside retry
- **WHEN** a Create or Modify or Delete retry block succeeds
- **THEN** no `d.SetId(...)`, `d.Set(...)`, or id-assignment statement executes inside the retry closure; all such statements run after the retry block returns and (for writes) after `WaitForGa2TaskFinish` returns.

### Requirement: Logging conventions
The resource SHALL emit:
- `defer tccommon.LogElapsed("resource.tencentcloud_ga2_global_accelerator_acl_policy.<op>")()` at the top of every CRUD function (op = create / read / update / delete).
- `defer tccommon.InconsistentCheck(d, meta)()` at the top of every CRUD function.
- A `[DEBUG]` line per SDK invocation containing the request action, request body, and response body (matching the existing `ga2` resource log format).
- A `[CRITAL]%s ... failed, reason:%+v` line on every retry-block failure.
- A `[CRUD] ga2_global_accelerator_acl_policy id=%s` line before clearing the ID in Read when the policy is absent.

The resource SHALL use the literal resource name `ga2_global_accelerator_acl_policy` (lowercase snake_case) in all log and error messages, never vague phrases like "该资源" or "当前资源".

#### Scenario: Standard log lines emitted
- **WHEN** any CRUD operation runs
- **THEN** the operation's elapsed time is logged via `tccommon.LogElapsed` and inconsistency is checked via `tccommon.InconsistentCheck`, and the lowercase resource name appears in all error/log strings.

### Requirement: Documentation and tests
The change SHALL include:
- A markdown document at `tencentcloud/services/ga2/resource_tc_ga2_global_accelerator_acl_policy.md` containing: a one-line description that mentions the GA2 (Global Accelerator V2) product; a self-contained `terraform { ... } resource "tencentcloud_ga2_global_accelerator_acl_policy" "..." { ... }` Example Usage block; and an `Import` section documenting the composite-ID import syntax (`terraform import tencentcloud_ga2_global_accelerator_acl_policy.x <global_accelerator_id>#<global_accelerator_acl_policy_id>`). The markdown SHALL NOT contain `Argument Reference` or `Attribute Reference` sections (those are auto-generated by `make doc`). Filename pattern follows `resource_tc_ga2_global_accelerator.md`.
- A test file at `tencentcloud/services/ga2/resource_tc_ga2_global_accelerator_acl_policy_test.go` using `gomonkey` mocks against the cloud-API client methods (NOT the Terraform acceptance suite), exercising at minimum: successful Create (with task-wait + ID set), successful Read (present and absent), successful Update (status change), and successful Delete. Tests SHALL be runnable via `go test -gcflags=all=-l ./tencentcloud/services/ga2/...`.

#### Scenario: Documentation present
- **WHEN** the change is merged
- **THEN** the markdown documentation file exists and contains a one-line GA2-mentioning description, an HCL Example Usage, and an `Import` section using the composite ID.

#### Scenario: Test file present and runnable
- **WHEN** the change is merged
- **THEN** the `_test.go` file declares gomonkey-based unit tests and `go test -gcflags=all=-l ./tencentcloud/services/ga2/...` succeeds without requiring `TENCENTCLOUD_SECRET_ID` / `TENCENTCLOUD_SECRET_KEY`.

### Requirement: SDK constraint
The implementation SHALL NOT modify any file under `vendor/github.com/tencentcloud/tencentcloud-sdk-go/`. If a required API is missing from the vendored SDK, the implementer MUST halt and request an SDK upgrade rather than authoring or patching SDK source. No `_extension.go` file SHALL be generated for this resource, and no comment header SHALL be added to the top of the resource `.go` file.

#### Scenario: Vendored SDK is sufficient
- **WHEN** the implementation begins
- **THEN** the four AclPolicy APIs plus `DescribeTaskResult` are confirmed present under `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ga2/v20250115/` before any resource code is written.

#### Scenario: No vendor modifications
- **WHEN** the change is merged
- **THEN** `git status` shows no modifications under `vendor/github.com/tencentcloud/tencentcloud-sdk-go/` attributable to this change.
