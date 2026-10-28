## Why

Tencent Cloud Global Accelerator V2 (`ga2`) currently exposes resources for the parent `GlobalAccelerator` instance, the `Listener`, and the `EndpointGroup`. The L7 listener model also has a **forwarding rule** layer (`ForwardingRule`, attached to a `ForwardingPolicy` under an HTTP/HTTPS listener) that lets users route HTTP(S) traffic by host/path/header conditions to specific origin behaviors. Without a Terraform resource for this object, users must hand-create forwarding rules in the console after Terraform provisions the listener. Adding `tencentcloud_ga2_forwarding_rule` finally makes the entire `Accelerator → Listener → ForwardingPolicy → ForwardingRule` chain Terraform-native.

## What Changes

- Add a new resource `tencentcloud_ga2_forwarding_rule` backed by the `ga2` v20250115 SDK.
- Map every `CreateForwardingRule` request parameter to a schema field, in particular:
  - identifying inputs: `global_accelerator_id`, `listener_id`, `forwarding_policy_id`
  - rule body: `rule_conditions` (set of `{rule_condition_type, rule_condition_value[]}`), `rule_actions` (set of `{rule_action_type, rule_action_value}`)
  - origin behavior: `origin_headers` (set of `{key, value}`), `enable_origin_sni`, `origin_sni`, `origin_host`
- Implement async-aware CRUD: `CreateForwardingRule`, `ModifyForwardingRule`, `DeleteForwardingRule` each return a `TaskId` that must be polled via `DescribeTaskResult` until `Status == "SUCCESS"`. Reuse the existing `Ga2Service.WaitForGa2TaskFinish(ctx, taskId, timeout)` helper.
- Add a new service helper `Ga2Service.DescribeGa2ForwardingRuleById(ctx, gaId, listenerId, policyId, ruleId) (*ga2v20250115.ForwardingRuleSet, error)` that wraps `DescribeForwardingRule` (the SDK call lists all rules under a `(gaId, listenerId, policyId)`) with `Limit=100` (documented maximum) and pagination, then strict-equals `*item.ForwardingRuleId == ruleId` before returning.
- Surface read-only computed fields from the describe response: `forwarding_rule_id` (also stored as the last segment of `d.Id()`).
- Wire the new resource into `tencentcloud/provider.go` under the `ga2` namespace block (adjacent to the three existing GA2 resources).
- Append `tencentcloud_ga2_forwarding_rule` to the `Global Accelerator(GA2)` Resources section of `tencentcloud/provider.md` so `make doc` picks it up.
- Author resource markdown documentation `resource_tc_ga2_forwarding_rule.md` (example HCL snippet + `terraform import` syntax).
- Author acceptance-test scaffolding `resource_tc_ga2_forwarding_rule_test.go`.
- Resource ID is the **4-segment composite** `<GlobalAcceleratorId>#<ListenerId>#<ForwardingPolicyId>#<ForwardingRuleId>` using `tccommon.FILED_SP`, since `Modify*`/`Delete*` and the `DescribeForwardingRule` lookup all require the four IDs together.
- All SDK calls are wrapped with `resource.Retry` (write paths use `tccommon.WriteRetryTimeout`, read paths use `tccommon.ReadRetryTimeout`).
- `Timeouts` block defaults to **5 minutes** for Create/Update/Delete, matching the other GA2 resources.
- `global_accelerator_id`, `listener_id`, `forwarding_policy_id` are **ForceNew** because `ModifyForwardingRule` does not allow moving a rule across (accelerator, listener, policy) boundaries.

## Capabilities

### New Capabilities
- `ga2-forwarding-rule-resource`: Lifecycle management (create / read / update / delete / import) of a Tencent Cloud Global Accelerator V2 layer-7 forwarding rule, including async task polling and full schema parity with `CreateForwardingRule`.

### Modified Capabilities
<!-- None: this change only introduces a new resource; it does not alter requirement-level behavior of any existing capability. -->

## Impact

- **New code**:
  - `tencentcloud/services/ga2/resource_tc_ga2_forwarding_rule.go` (CRUD + schema + build/flatten helpers, single file, mirroring `tencentcloud_igtm_monitor` style).
  - `tencentcloud/services/ga2/resource_tc_ga2_forwarding_rule.md` (resource doc + import syntax, mirroring `resource_tc_config_compliance_pack.md` filename convention).
  - `tencentcloud/services/ga2/resource_tc_ga2_forwarding_rule_test.go` (acceptance test skeleton, mirroring `resource_tc_config_compliance_pack_test.go` filename convention).
- **Modified code**:
  - `tencentcloud/services/ga2/service_tencentcloud_ga2.go`: add `DescribeGa2ForwardingRuleById`.
  - `tencentcloud/provider.go`: register `tencentcloud_ga2_forwarding_rule` adjacent to the three existing GA2 resources.
  - `tencentcloud/provider.md`: append `tencentcloud_ga2_forwarding_rule` to the existing `Global Accelerator(GA2)` Resources section.
- **APIs consumed**: `CreateForwardingRule`, `DescribeForwardingRule`, `ModifyForwardingRule`, `DeleteForwardingRule`, `DescribeTaskResult` (all already vendored in `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ga2/v20250115/`).
- **No breaking change**: purely additive; no existing schema or state is modified.
- **No SDK upgrade required**: all required APIs and structs (`RuleCondition`, `RuleAction`, `OriginHeader`, `ForwardingRuleSet`) are already present in the vendored SDK.
