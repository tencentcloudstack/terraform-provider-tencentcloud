# fwm-nat-acl-rule-resource Specification

## Purpose
TBD - created by archiving change add-fwm-nat-acl-rule-resource. Update Purpose after archive.
## Requirements
### Requirement: Resource MUST be registered as `tencentcloud_fwm_nat_acl_rule`

The provider SHALL register a new general-type resource named `tencentcloud_fwm_nat_acl_rule` that manages the full CRUD lifecycle of a single NAT edge ACL rule in the FWM product, using the `tencentcloud-sdk-go/tencentcloud/fwm/v20250611` SDK. The fwm service is new to the provider, so the connectivity client binding `UseFwmV20250611Client()` SHALL also be added.

#### Scenario: Resource registered in provider map

- **WHEN** the provider is loaded
- **THEN** `provider.go` exposes the resource via key `"tencentcloud_fwm_nat_acl_rule"` mapped to `fwm.ResourceTencentCloudFwmNatAclRule()`, and the import `services/fwm` is present.

#### Scenario: fwm connectivity client binding exists

- **WHEN** any fwm CRUD handler needs an API client
- **THEN** `connectivity/client.go` provides `UseFwmV20250611Client()` returning a cached `*fwmv20250611.Client` (mirror of `UseCfwV20190904Client`).

#### Scenario: Resource appears in gendoc index

- **WHEN** `tencentcloud/provider.md` is scanned by `make doc`
- **THEN** the FWM Resource section MUST include `tencentcloud_fwm_nat_acl_rule` so that `website/docs/r/fwm_nat_acl_rule.html.markdown` is generated.

### Requirement: Schema MUST mirror the NAT ACL rule writable fields plus read-only display fields

The resource schema SHALL declare the required rule-group / direction / rule-attribute keys with semantics matching the SDK `NatAclRule` input struct, and the computed keys from the `NatAclRuleResp` output struct.

Required (ForceNew marked):

| HCL key | SDK field | Type | ForceNew |
|---|---|---|---|
| `group_id` | `GroupId` | TypeString | Yes |
| `direction` | `Direction` | TypeInt | Yes |
| `source_content` | `SourceContent` | TypeString | No |
| `source_type` | `SourceType` | TypeString | Yes |
| `target_content` | `TargetContent` | TypeString | No |
| `target_type` | `TargetType` | TypeString | Yes |
| `protocol` | `Protocol` | TypeString | No |
| `rule_action` | `RuleAction` | TypeString | No |
| `order_index` | `OrderIndex` | TypeInt | No |
| `scope` | `Scope` | TypeString | No |

Optional:

| HCL key | SDK field | Type |
|---|---|---|
| `port` | `Port` | TypeString |
| `description` | `Description` | TypeString |
| `param_template_id` | `ParamTemplateId` | TypeString |
| `belong_member_id` | `BelongMemberId` | TypeString |

Computed (read-only): `rule_id`, `sequence`, `scope_desc`, `create_time`, `update_time`, `source_name`, `target_name`, `param_template_name`, `fw_ins_id`, `country_name`, `city_name`, `country_code`, `city_code`, `country_key`, `city_key`, `is_region`, `cloud_code`, `is_cloud`, `invalid`, `app_id`, `uin`, `nickname`, `member_id`, `node_name`, `dns_parse_cnt`, `hit_count`, `last_hit_time`.

#### Scenario: Required fields enforce on plan

- **WHEN** the user writes a config that omits `group_id`, `direction`, `source_content`, `source_type`, `target_content`, `target_type`, `protocol`, `rule_action`, `order_index`, or `scope`
- **THEN** `terraform plan` SHALL fail validation pointing at the missing required attribute.

#### Scenario: Changing direction / source_type / target_type forces replacement

- **GIVEN** state has `direction = 0`, `source_type = "ip"`, `target_type = "ip"`
- **WHEN** the user changes any of them in HCL
- **THEN** Terraform's plan reports a destroy + create cycle (ForceNew).

#### Scenario: bool/int optional fields use GetOkExists

- **WHEN** the code reads `direction` or `order_index`
- **THEN** it uses `d.GetOkExists(...)` so that `0` is not treated as unset.

### Requirement: Field name mismatch between input and output MUST be explicitly mapped

The input struct `NatAclRule` and output struct `NatAclRuleResp` use different names for two equivalent fields. The implementation SHALL map them explicitly to avoid permanent plan diff after import:

1. Input `NatAclRule.OrderIndex` (int64) ↔ Output `NatAclRuleResp.Sequence` (int64). On read, `d.Set("order_index", *rule.Sequence)`; on write, `rule.OrderIndex` from `d.Get("order_index")`.
2. Input `NatAclRule.BelongMemberId` (string) ↔ Output `NatAclRuleResp.BelongMember` (`MemberInfo`). On read, if `BelongMember != nil && BelongMember.MemberId != nil` then `d.Set("belong_member_id", *BelongMember.MemberId)`.

#### Scenario: order_index round-trips via Sequence

- **GIVEN** the user set `order_index = 5`
- **WHEN** Create then Read runs
- **THEN** the state `order_index` equals `5` (read from `Sequence`), producing no diff.

#### Scenario: belong_member_id round-trips via BelongMember.MemberId

- **GIVEN** the rule has a non-nil `BelongMember.MemberId = "mem-123"`
- **WHEN** Read runs
- **THEN** state `belong_member_id = "mem-123"`.

### Requirement: Resource ID MUST be composite GroupId + RuleId

After Create, the Terraform resource ID SHALL be `strings.Join([]string{groupId, ruleId}, tccommon.FILED_SP)` (`FILED_SP = "#"`). Read, Update, and Delete SHALL split `d.Id()` by `tccommon.FILED_SP` into `groupId` (index 0) and `ruleId` (index 1), erroring if the split length is not 2.

#### Scenario: ID encodes both GroupId and RuleId

- **GIVEN** Create succeeded with `GroupId = "fwmrg-abc"` and returned `RuleIds[0] = "fwmrule-xyz"`
- **WHEN** Create completes
- **THEN** `d.Id()` returns `"fwmrg-abc#fwmrule-xyz"`.

#### Scenario: Import supports GroupId:RuleId

- **GIVEN** the user runs `terraform import tencentcloud_fwm_nat_acl_rule.my fwmrg-abc:fwmrule-xyz`
- **WHEN** import runs
- **THEN** the resource ID is normalized to `"fwmrg-abc#fwmrule-xyz"` and Read populates state.

### Requirement: Create MUST pass a single-element Rules array and bind RuleIds[0]

The Create callback SHALL build a `CreateNatAclRuleRequest` with `Rules` containing exactly one `NatAclRule` element (single-instance resource), call `CreateNatAclRuleWithContext` inside `resource.Retry(tccommon.WriteRetryTimeout, ...)`, guard `response == nil || response.Response == nil || len(response.Response.RuleIds) == 0` as `NonRetryableError`, then bind `RuleId = *response.Response.RuleIds[0]` and `d.SetId(...)` BEFORE returning. Create MUST end by calling Read. All four interfaces are synchronous — NO `waitReady`/async polling SHALL be generated.

#### Scenario: Create issues a single-element Rules request

- **WHEN** the user runs `terraform apply` with a valid resource config
- **THEN** exactly one `CreateNatAclRule` request is issued with `Rules` of length 1, and `d.Id()` is set to `GroupId#RuleId` from `RuleIds[0]`.

#### Scenario: Empty RuleIds returns NonRetryableError

- **GIVEN** the API returns `RuleIds` empty
- **WHEN** the retry callback runs
- **THEN** the callback returns `resource.NonRetryableError(...)` and does NOT set an empty ID.

### Requirement: Read MUST query DescribeNatAclRules with Direction and RuleId filter

The Read callback SHALL call a service-layer `DescribeFwmNatAclRuleById(ctx, groupId, ruleId, direction)` that issues `DescribeNatAclRules` with `GroupId`, `Direction` (required, 0/1), `Limit=100`, `Offset=0`, and `Filters=[{Name:"RuleId", Values:[ruleId], OperatorType:1}]` (1=EQUAL). It SHALL take `Rules[0]`. If the API returns a `ResourceNotFound` error OR an empty `Rules` list, the handler SHALL `log.Printf("[CRUD] fwm_nat_acl_rule id=%s", d.Id())` then `d.SetId("")` and return nil (externally deleted). All `d.Set` calls SHALL be nil-guarded.

#### Scenario: Externally deleted rule clears state

- **GIVEN** the rule no longer exists server-side
- **WHEN** Read runs
- **THEN** the handler logs the id, calls `d.SetId("")`, and returns nil (no error).

#### Scenario: ScopeDesc and display names are not plan-compared

- **WHEN** Read populates state
- **THEN** `scope_desc`, `source_name`, `target_name`, `param_template_name` are written as Computed fields only; they do not appear as user-writable arguments and never force a diff.

### Requirement: Update MUST send a full-overwrite ModifyNatAclRule

The Update callback SHALL build a `ModifyNatAclRuleRequest` with `GroupId` and a single `Rule *NatAclRule` containing `RuleId` (from the composite ID) plus ALL writable fields from current state (full overwrite semantics; the server reads the old record then writes the complete new value). It SHALL call `ModifyNatAclRuleWithContext` inside `resource.Retry(tccommon.WriteRetryTimeout, ...)`, guard nil response, and end by calling Read. `GroupId` MUST NOT be changed.

#### Scenario: Editing scope triggers one overwrite Modify

- **GIVEN** state has `scope = "ALL"`
- **WHEN** the user changes `scope` to `"ap-guangzhou"`
- **THEN** Update issues exactly one `ModifyNatAclRule` request with `Rule.RuleId`, `Rule.Scope = "ap-guangzhou"`, and all other writable fields present.

#### Scenario: Transient SDK error is retried

- **GIVEN** the first invocation of `ModifyNatAclRule` returns a retriable error
- **WHEN** the retry callback re-runs
- **THEN** the second attempt's response is observed; the original error is surfaced only after the retry budget is exhausted.

### Requirement: Delete MUST be idempotent on ResourceNotFound

The Delete callback SHALL build a `DeleteNatAclRuleRequest` with `GroupId` and `RuleIds` containing exactly one element (the ruleId from the composite ID), call `DeleteNatAclRuleWithContext` inside `resource.Retry(tccommon.WriteRetryTimeout, ...)`. If the error contains `ResourceNotFound`, the handler SHALL treat it as success (return nil).

#### Scenario: Deleting a non-existent rule is a no-op

- **GIVEN** the rule was already deleted server-side
- **WHEN** the user runs `terraform destroy`
- **THEN** the `DeleteNatAclRule` call returns `ResourceNotFound`, and Delete returns nil without error.

### Requirement: Documentation and unit test MUST follow project conventions

- The HCL example markdown SHALL live at `tencentcloud/services/fwm/resource_tc_fwm_nat_acl_rule.md` containing a one-line description mentioning FWM, an Example Usage (with `jsonencode` if JSON string fields involved), and an Import section documenting `GroupId:RuleId`.
- A unit test SHALL live at `tencentcloud/services/fwm/resource_tc_fwm_nat_acl_rule_test.go` using gomonkey to mock the cloud APIs (NOT the terraform test suite), covering Create, Read, Update, and Delete (including the `ResourceNotFound` idempotent case).
- Running `make doc` SHALL regenerate `website/docs/r/fwm_nat_acl_rule.html.markdown`.

#### Scenario: Test uses gomonkey mocks

- **WHEN** the test file is opened
- **THEN** the package is `fwm_test`, tests mock `UseFwmV20250611Client` and the `*WithContext` SDK methods, and no terraform AccTest are used.

