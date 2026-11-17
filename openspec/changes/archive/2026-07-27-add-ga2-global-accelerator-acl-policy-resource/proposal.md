## Why

Tencent Cloud Global Accelerator V2 (`ga2`) already ships Terraform resources for the accelerator instance (`tencentcloud_ga2_global_accelerator`) and its child objects (listener, endpoint group, forwarding policy/rule, accelerate area), but there is **no** Terraform coverage for the **访问控制策略 (AclPolicy)** that gates which source traffic is allowed/dropped on an accelerator instance. Today operators must create, enable/disable and delete these ACL policies through the console or raw API calls, which breaks the fully-Terraform-native workflow goal of the `ga2` namespace and makes security posture un-reviewable in `terraform plan`. This change closes that gap by adding the `tencentcloud_ga2_global_accelerator_acl_policy` resource.

## What Changes

- Add a new resource `tencentcloud_ga2_global_accelerator_acl_policy` backed by the `ga2` v20250115 SDK, implementing full CRUD (RESOURCE_KIND_GENERAL) for a global accelerator ACL policy.
- Schema fields mirror the cloud APIs:
  - `global_accelerator_id` (string, Required, **ForceNew**) — the parent global accelerator instance ID; required by every CRUD call and immutable.
  - `default_action` (string, Required, **ForceNew**) — default traffic action (`ACCEPT` / `DROP`). Only accepted by `CreateGlobalAcceleratorAclPolicy`; `ModifyGlobalAcceleratorAclPolicy` has no slot for it, so it is modeled ForceNew to avoid a silently-un-updateable field.
  - `status` (string, Optional, Computed) — ACL policy state (`OPEN` / `CLOSE`); updatable in place via `ModifyGlobalAcceleratorAclPolicy`.
  - `global_accelerator_acl_policy_id` (string, Computed) — the ACL policy ID returned by Create.
  - `task_id` (string, Computed) — the async task ID surfaced from the most recent write call (Create / Modify / Delete), for operator observability.
- Implement async-aware CRUD: `CreateGlobalAcceleratorAclPolicy`, `ModifyGlobalAcceleratorAclPolicy` and `DeleteGlobalAcceleratorAclPolicy` all return a `TaskId` that must be polled to `Status == "SUCCESS"` via the existing `Ga2Service.WaitForGa2TaskFinish(ctx, taskId, timeout)` helper (which polls `DescribeTaskResult`).
- Add a new service helper `Ga2Service.DescribeGa2GlobalAcceleratorAclPolicyById(ctx, gaId, policyId string) (*GlobalAcceleratorAclPolicies, error)` that wraps `DescribeGlobalAcceleratorAclPolicies`. Because that API only accepts `GlobalAcceleratorId` (+ Offset/Limit) and has no per-policy filter, the helper paginates with `Limit="200"` (the documented maximum) and strictly matches `*item.GlobalAcceleratorAclPolicyId == policyId` client-side; returns `(nil, nil)` when absent.
- Use a **composite resource ID** `<global_accelerator_id>#<global_accelerator_acl_policy_id>` (separator `tccommon.FILED_SP`, i.e. `#`) because every Read/Update/Delete call needs both IDs. Read/Update/Delete split the ID back into the two components. Import therefore requires the composite ID (documented in the resource markdown).
- Wire the new resource into `tencentcloud/provider.go` under the `ga2` namespace, adjacent to the existing `ga2` entries.
- Author resource markdown documentation `resource_tc_ga2_global_accelerator_acl_policy.md` (one-line description mentioning the GA2 product, Example Usage HCL snippet, and an `Import` section that documents the composite-ID import syntax). Do NOT hand-edit `website/docs/` — it will be regenerated via `make doc`.
- Author unit tests `resource_tc_ga2_global_accelerator_acl_policy_test.go` using `gomonkey` mocks (per the project rule for newly-added resources: no Terraform test suite, only business-logic unit tests runnable via `go test -gcflags=all=-l`).
- All SDK calls are wrapped with `resource.Retry` (write paths use `tccommon.WriteRetryTimeout`, read paths use `tccommon.ReadRetryTimeout`); id assignment / state-setting happens **outside** the retry block; nil-response defenses return `NonRetryableError`.

## Capabilities

### New Capabilities
- `ga2-global-accelerator-acl-policy-resource`: Lifecycle management (create / read / update / delete / import) of a single Tencent Cloud Global Accelerator V2 access-control policy (AclPolicy) scoped to one global accelerator instance, including async task polling, composite-ID import, schema parity with the four cloud APIs, and unit-tested business logic.

### Modified Capabilities
<!-- None: this change only introduces a new resource; it does not alter requirement-level behavior of any existing capability. -->

## Impact

- **New code**:
  - `tencentcloud/services/ga2/resource_tc_ga2_global_accelerator_acl_policy.go` (schema + CRUD + composite-ID split helpers, single file, mirroring `tencentcloud_ga2_forwarding_policy` / `tencentcloud_igtm_strategy` style).
  - `tencentcloud/services/ga2/resource_tc_ga2_global_accelerator_acl_policy.md` (resource doc + composite-ID import syntax).
  - `tencentcloud/services/ga2/resource_tc_ga2_global_accelerator_acl_policy_test.go` (gomonkey-based unit tests).
- **Modified code**:
  - `tencentcloud/services/ga2/service_tencentcloud_ga2.go`: add `DescribeGa2GlobalAcceleratorAclPolicyById`.
  - `tencentcloud/provider.go`: register `tencentcloud_ga2_global_accelerator_acl_policy` in `ResourcesMap`.
- **APIs consumed**: `CreateGlobalAcceleratorAclPolicy`, `DescribeGlobalAcceleratorAclPolicies`, `ModifyGlobalAcceleratorAclPolicy`, `DeleteGlobalAcceleratorAclPolicy`, `DescribeTaskResult` (all already vendored in `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ga2/v20250115/`).
- **No breaking change**: purely additive; no existing schema or state is modified.
- **No SDK upgrade required**: all required APIs are already present in the vendored SDK; no `vendor/` file is modified.
