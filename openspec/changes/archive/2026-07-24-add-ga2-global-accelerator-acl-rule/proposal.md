## Why

Tencent Cloud Global Accelerator V2 (`ga2`) exposes ACL rules as a sub-resource of ACL policies on global accelerator instances. Currently, the provider has no way to manage these ACL rules through Terraform, forcing users to configure them manually in the console. Adding `tencentcloud_ga2_global_accelerator_acl_rule` closes this gap, enabling full Terraform-native management of ACL rules for GA2 instances.

## What Changes

- Add a new resource `tencentcloud_ga2_global_accelerator_acl_rule` backed by the `ga2` v20250115 SDK.
- Map all `CreateGlobalAcceleratorAclRule` / `ModifyGlobalAcceleratorAclRule` request parameters to schema fields: `global_accelerator_id`, `global_accelerator_acl_policy_id`, `global_accelerator_acl_rule_id`, `protocol`, `port`, `source_cidr_block`, `policy`, `description`.
- Implement async-aware CRUD: `CreateGlobalAcceleratorAclRule`, `ModifyGlobalAcceleratorAclRule`, and `DeleteGlobalAcceleratorAclRule` all return a `TaskId` that must be polled via `DescribeTaskResult` until `Status == "SUCCESS"`. Reuse the existing `Ga2Service.WaitForGa2TaskFinish(ctx, taskId, timeout)` helper.
- Add a new service helper `Ga2Service.DescribeGa2GlobalAcceleratorAclRuleById(ctx, policyId, ruleId string) (*GlobalAcceleratorAclRuleSet, error)` that wraps `DescribeGlobalAcceleratorAclRules` with pagination (Limit=200, the documented maximum) and matches by `GlobalAcceleratorAclRuleId`.
- Resource ID is a composite of `GlobalAcceleratorId#GlobalAcceleratorAclPolicyId#GlobalAcceleratorAclRuleId` using `tccommon.FILED_SP` as separator, since the describe API is keyed by policy rather than the rule ID alone.
- Wire the new resource into `tencentcloud/provider.go` under the `ga2` namespace.
- Author resource markdown documentation `resource_tc_ga2_global_accelerator_acl_rule.md`.
- Author acceptance-test scaffolding `resource_tc_ga2_global_accelerator_acl_rule_test.go`.
- All SDK calls are wrapped with `resource.Retry` (write paths use `tccommon.WriteRetryTimeout`, read paths use `tccommon.ReadRetryTimeout`).

## Capabilities

### New Capabilities
- `ga2-global-accelerator-acl-rule-resource`: Lifecycle management (create / read / update / delete / import) of a Tencent Cloud Global Accelerator V2 ACL rule, including async task polling and full schema parity with the `CreateGlobalAcceleratorAclRule` / `ModifyGlobalAcceleratorAclRule` APIs.

### Modified Capabilities
<!-- None: this change only introduces a new resource; it does not alter requirement-level behavior of any existing capability. -->

## Impact

- **New code**:
  - `tencentcloud/services/ga2/resource_tc_ga2_global_accelerator_acl_rule.go` (CRUD + schema + build/flatten helpers).
  - `tencentcloud/services/ga2/resource_tc_ga2_global_accelerator_acl_rule.md` (resource doc + import syntax).
  - `tencentcloud/services/ga2/resource_tc_ga2_global_accelerator_acl_rule_test.go` (acceptance test skeleton).
- **Modified code**:
  - `tencentcloud/services/ga2/service_tencentcloud_ga2.go`: add `DescribeGa2GlobalAcceleratorAclRuleById`.
  - `tencentcloud/provider.go`: register `tencentcloud_ga2_global_accelerator_acl_rule` in `ResourcesMap`.
  - `tencentcloud/provider.md`: register resource in documentation.
- **APIs consumed**: `CreateGlobalAcceleratorAclRule`, `DescribeGlobalAcceleratorAclRules`, `ModifyGlobalAcceleratorAclRule`, `DeleteGlobalAcceleratorAclRule`, `DescribeTaskResult` (already vendored in `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ga2/v20250115/`).
- **No breaking change**: purely additive; no existing schema or state is modified.
- **No SDK upgrade required**: all required APIs are already present in the vendored SDK.