## Why

TencentCloud Global Accelerator 2.0 (GA2) supports Layer-7 forwarding rules that allow users to route traffic based on conditions (e.g., host, path) and perform actions (e.g., forward to endpoint groups). Currently, there is no Terraform resource to manage these forwarding rules, requiring users to configure them manually through the console or API.

## What Changes

- Add a new Terraform resource `tencentcloud_ga2_forwarding_rule` that provides full CRUD lifecycle management for GA2 Layer-7 forwarding rules.
- The resource supports:
  - Creating forwarding rules with conditions, actions, origin headers, origin SNI, and origin host settings.
  - Reading forwarding rule state via the DescribeForwardingRule API with pagination.
  - Modifying forwarding rule conditions, actions, and origin settings.
  - Deleting forwarding rules.
- All Create/Modify/Delete operations are asynchronous (return TaskId), so the resource will poll DescribeForwardingRule after each mutation to confirm the operation has taken effect.
- The resource uses a composite ID consisting of `global_accelerator_id`, `listener_id`, `forwarding_policy_id`, and `forwarding_rule_id` joined by `tccommon.FILED_SP`.

## Capabilities

### New Capabilities
- `ga2-forwarding-rule-crud`: Full CRUD lifecycle management for GA2 Layer-7 forwarding rules, including async operation polling and composite ID handling.

### Modified Capabilities

## Impact

- New files:
  - `tencentcloud/services/ga2/resource_tc_ga2_forwarding_rule.go` - Resource implementation
  - `tencentcloud/services/ga2/resource_tc_ga2_forwarding_rule_test.go` - Unit tests with gomonkey mocks
  - `tencentcloud/services/ga2/resource_tc_ga2_forwarding_rule.md` - Example usage documentation
- Modified files:
  - `tencentcloud/services/ga2/service_tencentcloud_ga2.go` - Add service-layer helper functions for forwarding rule CRUD
  - `tencentcloud/provider.go` - Register the new resource
  - `tencentcloud/provider.md` - Add resource entry to provider documentation
- Dependencies: Uses existing `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ga2/v20250115` package (already vendored)
