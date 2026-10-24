## Why

Tencent Cloud Firewall Manager (`fwm`) exposes an "enterprise security group" (企业安全组) rule model through the `fwm/v20250611` OpenAPI, but the Terraform provider currently has no resource for it — only the legacy `cfw` product (`tencentcloud_cfw_edge_policy` etc.) is covered. Users who manage enterprise security group rules must create them in the console, which breaks the "everything as code" workflow and makes rule drift auditing impossible. This change adds `tencentcloud_fwm_security_group_rule` so a single enterprise security group rule (including its owning rule group) can be created, read, updated, deleted and imported by Terraform.

## What Changes

- Add a new resource `tencentcloud_fwm_security_group_rule` backed by the vendored `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611` SDK.
- Add a new service package `tencentcloud/services/fwm/` with:
  - `service_tencentcloud_fwm.go` — `FwmService` wrapper with `DescribeFwmSecurityGroupRuleById(ctx, groupId, ruleId)`, which calls `DescribeSecurityGroupRules` with `Filters = [{Name: "RuleId", Values: [ruleId], OperatorType: 9}]` plus `GroupId`, then strict-equals the returned `RuleId` before returning.
  - `resource_tc_fwm_security_group_rule.go` — schema + Create/Read/Update/Delete.
  - `resource_tc_fwm_security_group_rule.md` — resource example + `terraform import` syntax.
  - `resource_tc_fwm_security_group_rule_test.go` — gomonkey-based unit tests (no terraform acceptance suite for a brand-new resource).
- Register the API client: add `fwmV20250611Conn *fwmv20250611.Client` cache field and `UseFwmV20250611Client()` to `tencentcloud/connectivity/client.go`.
- Model **one rule per resource instance**, mirroring the API capabilities:
  - `CreateSecurityGroupRule` accepts a batch (`Rules []*SecurityGroupRule`) and returns `RuleIds []*string`; the resource submits a one-element `Rules` slice and persists `RuleIds[0]` as the rule identity.
  - `DescribeSecurityGroupRules` can pinpoint exactly one rule via the `RuleId` filter.
  - `ModifySecurityGroupRule` mutates exactly one rule (`Rule` field, must carry `RuleId`).
  - `DeleteSecurityGroupRule` deletes by `RuleIds`; for a single rule it is called with one element.
- Rule-group bootstrap: `CreateSecurityGroupRule` requires a `GroupId`. When the user does not pass an existing `group_id`, Create first calls `CreateRuleGroup` with `GroupName = group_name` and `Product = "enterprise_sg"`, then uses the returned `GroupId` for the rule creation. `group_id` is therefore `Optional + Computed` and `ForceNew`; `group_name` is `Optional` and `ForceNew` (ignored when `group_id` is given).
- Resource ID is the **2-segment composite** `<GroupId>#<RuleId>` using `tccommon.FILED_SP`, because Read/Modify/Delete all need both IDs.
- All four entry points are **synchronous** (no `TaskId`/`FlowId` returned), so no `Timeouts` block and no polling loop is introduced; every SDK call is still wrapped in `resource.Retry` with `tccommon.WriteRetryTimeout` (Create/Update/Delete) or `tccommon.ReadRetryTimeout` (Read).
- Read surfaces the describe-only attributes as `Computed` fields (`rule_id`, `source_id`, `target_id`, `strategy`, `detail`, `region`, `protocol_port_name`, `parameter_name`, `sou_instance_name`, `sou_public_ip`, `sou_private_ip`, `sou_cidr`, `sou_parameter_name`, `instance_name`, `public_ip`, `private_ip`, `cidr`, `is_new`, `rule_uuid`, `dns_parse_count`, `belong_member`) while **not** overwriting input fields whose API encoding differs (see design.md D6).
- Wire the resource into `tencentcloud/provider.go` and append it to a new `Firewall Manager(FWM)` Resources section of `tencentcloud/provider.md` so `make doc` picks it up.

## Capabilities

### New Capabilities
- `fwm-security-group-rule-resource`: Lifecycle management (create / read / update / delete / import) of one Tencent Cloud Firewall Manager enterprise security group rule, including on-demand rule-group bootstrap through `CreateRuleGroup` and precise single-rule reads via the `RuleId` filter.

### Modified Capabilities
<!-- None: this change is purely additive; no existing capability's requirements change. -->

## Impact

- **New code**:
  - `tencentcloud/services/fwm/service_tencentcloud_fwm.go`
  - `tencentcloud/services/fwm/resource_tc_fwm_security_group_rule.go`
  - `tencentcloud/services/fwm/resource_tc_fwm_security_group_rule.md`
  - `tencentcloud/services/fwm/resource_tc_fwm_security_group_rule_test.go`
- **Modified code**:
  - `tencentcloud/connectivity/client.go`: import `fwmv20250611` and add the cached client accessor `UseFwmV20250611Client()`.
  - `tencentcloud/provider.go`: register `"tencentcloud_fwm_security_group_rule": fwm.ResourceTencentCloudFwmSecurityGroupRule()`.
  - `tencentcloud/provider.md`: add a `Firewall Manager(FWM)` section listing the new resource.
- **APIs consumed** (all already vendored in `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611/`):
  `CreateRuleGroup`, `CreateSecurityGroupRule`, `DescribeSecurityGroupRules`, `ModifySecurityGroupRule`, `DeleteSecurityGroupRule`.
- **No breaking change**: purely additive; no existing schema, state or API client is modified.
- **No SDK upgrade required**: the vendored module `fwm v1.3.196` already contains every struct used (`SecurityGroupRule`, `SecGroupRuleResp`, `CommonFilter`, `SgDnsParseCount`, `MemberInfo`).
