## Why

Tencent Cloud Global Accelerator V2 (`ga2`) currently exposes Terraform resources for `GlobalAccelerator`, `Listener`, `EndpointGroup`, and `ForwardingRule`. The `ForwardingPolicy` layer sits between the HTTP/HTTPS listener and the forwarding rules — it manages the L7 host-based routing domain. Without a dedicated Terraform resource for `ForwardingPolicy`, users cannot manage forwarding policies (host domains) through IaC, forcing them to use the console for this critical piece of the `Accelerator → Listener → ForwardingPolicy → ForwardingRule` chain. Adding `tencentcloud_ga2_forwarding_policy` completes the full GA2 resource hierarchy in Terraform.

## What Changes

- Add a new resource `tencentcloud_ga2_forwarding_policy` backed by the `ga2` v20250115 SDK.
- Map the Create/Read/Modify/Delete API operations to standard Terraform CRUD:
  - `CreateForwardingPolicy` → create a forwarding policy with `Host` domain under a `(global_accelerator_id, listener_id)` pair. Returns `ForwardingPolicyId` + async `TaskId`.
  - `DescribeForwardingPolicy` → read the policy by listing all policies under the parent listener and matching `ForwardingPolicyId`. Returns `ForwardingPolicySet` items with `Host`, `DefaultHostFlag`.
  - `ModifyForwardingPolicy` → update the `Host` domain. Returns async `TaskId`.
  - `DeleteForwardingPolicy` → delete the policy. Returns async `TaskId`.
- All Create/Modify/Delete operations are asynchronous; each returns a `TaskId` that must be polled via `DescribeTaskResult` until `Status == "SUCCESS"`.
- Resource ID: 3-segment composite `<global_accelerator_id>#<listener_id>#<forwarding_policy_id>` using `tccommon.FILED_SP`, matching the pattern established by other GA2 resources.
- `global_accelerator_id` and `listener_id` are `ForceNew` — `ModifyForwardingPolicy` uses them only as identifiers and cannot relocate a policy.
- `host` is the only updatable field via `ModifyForwardingPolicy`.

## Capabilities

### New Capabilities
- `ga2-forwarding-policy-resource`: Full lifecycle management (create / read / update / delete / import) of a Tencent Cloud Global Accelerator V2 layer-7 forwarding policy, including async task polling and schema parity with the GA2 SDK.

### Modified Capabilities
<!-- None: this change only introduces a new resource; it does not alter requirement-level behavior of any existing capability. -->

## Impact

- **New code**:
  - `tencentcloud/services/ga2/resource_tc_ga2_forwarding_policy.go` (CRUD + schema + build/flatten helpers)
  - `tencentcloud/services/ga2/resource_tc_ga2_forwarding_policy.md` (resource doc + import syntax)
  - `tencentcloud/services/ga2/resource_tc_ga2_forwarding_policy_test.go` (acceptance test)
- **Modified code**:
  - `tencentcloud/services/ga2/service_tencentcloud_ga2.go`: add `DescribeGa2ForwardingPolicyById` helper
  - `tencentcloud/provider.go`: register `tencentcloud_ga2_forwarding_policy` resource
  - `tencentcloud/provider.md`: append `tencentcloud_ga2_forwarding_policy` to the GA2 Resources section
- **APIs consumed**: `CreateForwardingPolicy`, `DescribeForwardingPolicy`, `ModifyForwardingPolicy`, `DeleteForwardingPolicy`, `DescribeTaskResult` (all already vendored in `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ga2/v20250115/`).
- **No breaking change**: purely additive; no existing schema or state is modified.