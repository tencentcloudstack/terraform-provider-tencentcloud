## Why

Tencent Cloud Firewall Manager (FWM, 云防火墙规则组管理) exposes a complete synchronous CRUD surface for firewall rule groups — `CreateRuleGroup`, `DescribeRuleGroups`, `ModifyRuleGroup`, `DeleteRuleGroup` — in the `fwm/v20250611` SDK namespace. A rule group (`fwmrg-` prefix) is the container that four firewall products share (enterprise security group, VPC border, NAT border, internet border), and it is the prerequisite object for every firewall rule. The provider currently exposes no FWM resource at all, so users cannot manage rule groups as infrastructure-as-code.

The SDK package `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm v1.3.196` (namespace `v20250611`) has already been added to `go.mod`, `tools.go` and materialized under `vendor/`, which unblocks this change.

## What Changes

- Add new general (CRUD) resource `tencentcloud_fwm_rule_group` in a new `fwm` service package.
- Resource schema exposes `group_name` and `product` as user-configurable fields (`group_name` updatable, `product` immutable because `ModifyRuleGroup` cannot change it), plus read-only attributes returned by `DescribeRuleGroups`: `group_id`, `app_id`, `rule_count`, `create_time`, `update_time`, `create_by`, `create_name`, `update_by`, `update_name`, `receive_account` (nested list with `members`), `exec_areas`, `direction`.
- Resource ID is the bare rule group ID (`fwmrg-xxxx`) returned by `CreateRuleGroup`; the resource is importable with that ID.
- `product` accepts exactly four product codes — `enterprise_sg`, `cfw_nat_acl`, `cfw_vpc_acl`, `cfw_edge_acl` — validated at plan time.
- Because `DescribeRuleGroups` returns `Product` as a numeric code (`1|2|3|4`), Read normalizes the code back to the product string so no spurious diff is produced against the user's configuration.
- Add service-layer helpers in `service_tencentcloud_fwm.go`: `DescribeFwmRuleGroupById` (paginated `DescribeRuleGroups` with a `GroupId` filter) and `DescribeFwmRuleGroupsByFilter`.
- Register a new client entry `UseFwmV20250611Client` in `tencentcloud/connectivity/client.go`.
- Register the resource in `tencentcloud/provider.go` and `tencentcloud/provider.md`.
- Add `resource_tc_fwm_rule_group.md` example documentation and a gomonkey-based unit test file (new resources are unit-tested with mocked cloud APIs, not acceptance tests).

## Capabilities

### New Capabilities
- `fwm-rule-group-resource`: Full synchronous CRUD lifecycle, plus import, for a FWM rule group bound to a single firewall product.

### Modified Capabilities

## Impact

- New files:
  - `tencentcloud/services/fwm/resource_tc_fwm_rule_group.go`
  - `tencentcloud/services/fwm/resource_tc_fwm_rule_group.md`
  - `tencentcloud/services/fwm/resource_tc_fwm_rule_group_test.go`
  - `tencentcloud/services/fwm/service_tencentcloud_fwm.go`
- Modified files:
  - `tencentcloud/connectivity/client.go` — add `fwm` import and `UseFwmV20250611Client`
  - `tencentcloud/provider.go` — register `tencentcloud_fwm_rule_group`
  - `tencentcloud/provider.md` — document `tencentcloud_fwm_rule_group`
- SDK dependency: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm v1.3.196` (already added to `go.mod` / `tools.go` and vendored).
- No existing resource schema or state is touched, so the change is fully backward compatible.
