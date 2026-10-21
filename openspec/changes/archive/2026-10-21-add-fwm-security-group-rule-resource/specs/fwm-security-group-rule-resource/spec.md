## ADDED Requirements

### Requirement: Resource registration and schema

The provider SHALL register a new general resource `tencentcloud_fwm_security_group_rule` (function `fwm.ResourceTencentCloudFwmSecurityGroupRule()`), and the resource SHALL expose a schema that maps one enterprise security group rule of the `fwm/v20250611` API.

The schema MUST declare these input attributes (all `Optional` unless stated, all of them `Computed` so server-side defaults do not diff):

| Attribute | Type | Notes |
|---|---|---|
| `group_id` | String | `Optional` + `Computed` + `ForceNew`; existing rule group id |
| `group_name` | String | `Optional` + `ForceNew`; name used to create a new rule group |
| `ip_version` | String | `SecurityGroupRule.IpVersion` |
| `source_content` | String | `SecurityGroupRule.SourceContent` |
| `source_type` | String | `SecurityGroupRule.SourceType` |
| `dest_content` | String | `SecurityGroupRule.DestContent` |
| `dest_type` | String | `SecurityGroupRule.DestType` |
| `protocol` | String | `SecurityGroupRule.Protocol` |
| `port` | String | `SecurityGroupRule.Port` |
| `service_template_id` | String | `SecurityGroupRule.ServiceTemplateId` |
| `rule_action` | String | `SecurityGroupRule.RuleAction` |
| `description` | String | `SecurityGroupRule.Description` |
| `order_index` | Int | `SecurityGroupRule.OrderIndex` |
| `scope` | String | `SecurityGroupRule.Scope` |
| `protocol_port_type` | Int | `SecurityGroupRule.ProtocolPortType` |
| `belong_member_id` | String | `SecurityGroupRule.BelongMemberId` |

and these computed-only attributes: `rule_id`, `source_id`, `target_id`, `target_type`, `strategy`, `detail`, `region`, `sou_instance_name`, `sou_public_ip`, `sou_private_ip`, `sou_cidr`, `sou_parameter_name`, `instance_name`, `public_ip`, `private_ip`, `cidr`, `parameter_name`, `protocol_port_name`, `is_new`, `rule_uuid`, plus the `MaxItems: 1` computed blocks `dns_parse_count` (`valid_count`, `invalid_count`) and `belong_member` (`app_id`, `uin`, `nickname`, `member_id`, `node_name`).

The resource MUST NOT declare an attribute literally named `id`, and MUST NOT declare an `all_total_count` / `total_count` / `rule_ids` attribute (see design.md D5).

No `Timeouts` block MUST be declared, because every API used by this resource is synchronous.

#### Scenario: Resource is discoverable in the provider

- **WHEN** the provider's `ResourcesMap` is inspected
- **THEN** it contains the key `tencentcloud_fwm_security_group_rule` bound to `fwm.ResourceTencentCloudFwmSecurityGroupRule()`

#### Scenario: Only identity attributes force replacement

- **WHEN** the schema of `tencentcloud_fwm_security_group_rule` is inspected
- **THEN** exactly `group_id` and `group_name` are marked `ForceNew`, and every other input attribute is updatable in place

#### Scenario: Computed blocks are single-element lists

- **WHEN** a rule returned by the API contains `DnsParseCount` and `BelongMember`
- **THEN** `dns_parse_count` and `belong_member` each contain exactly one element whose keys mirror the SDK field names in snake_case

---

### Requirement: Rule group bootstrap on create

`Create` SHALL require exactly one of `group_id` / `group_name`, and SHALL create a rule group through `CreateRuleGroup` when the user supplies `group_name` instead of `group_id`.

- When only `group_name` is set, Create MUST call `CreateRuleGroup` with `GroupName = group_name` and `Product = "enterprise_sg"`, and MUST use the returned `GroupId` for the subsequent `CreateSecurityGroupRule` call.
- When `group_id` is set, Create MUST NOT call `CreateRuleGroup`.
- When neither or both are set, Create MUST fail with a `resource.NonRetryableError` whose message names `tencentcloud_fwm_security_group_rule` and explains that exactly one of `group_id` / `group_name` is required.
- If `CreateRuleGroup` returns a nil `Response` or a nil/empty `GroupId`, Create MUST fail with a `resource.NonRetryableError`.

#### Scenario: Group created from group_name

- **WHEN** a configuration sets `group_name = "tf-example"` and no `group_id`
- **THEN** `CreateRuleGroup` is invoked with `GroupName = "tf-example"` and `Product = "enterprise_sg"`, and the resulting `GroupId` is passed to `CreateSecurityGroupRule`

#### Scenario: Existing group reused

- **WHEN** a configuration sets `group_id = "fwmrg-xxxxxxxx"` and no `group_name`
- **THEN** `CreateRuleGroup` is never invoked and `CreateSecurityGroupRule` receives `GroupId = "fwmrg-xxxxxxxx"`

#### Scenario: Neither identity provided

- **WHEN** neither `group_id` nor `group_name` is set
- **THEN** Create fails with a non-retryable error and no API call is made

#### Scenario: Both identities provided

- **WHEN** both `group_id` and `group_name` are set
- **THEN** Create fails with a non-retryable error stating the two attributes are mutually exclusive

---

### Requirement: Create submits a single-element rule batch and persists the composite id

`Create` SHALL build one `SecurityGroupRule` from the schema, submit it as a one-element `Rules` slice to `CreateSecurityGroupRule`, and persist the returned rule id.

- The request MUST be built with `GroupId` plus `Rules: []*SecurityGroupRule{rule}`; schema fields MUST be copied to their SDK counterparts only when `d.GetOk(...)` reports a value.
- The SDK call MUST be wrapped in `resource.Retry(tccommon.WriteRetryTimeout, ...)` and errors MUST be wrapped with `tccommon.RetryError(e)`.
- If the response, `Response.RuleIds`, or `RuleIds[0]` is empty, Create MUST return a `resource.NonRetryableError` that includes the `logId`.
- On success Create MUST set `d.SetId(strings.Join([]string{groupId, ruleId}, tccommon.FILED_SP))` and set `rule_id`, then delegate to Read.

#### Scenario: Successful create

- **WHEN** `CreateSecurityGroupRule` returns `RuleIds = ["fwmrule-aaaaaaaa"]`
- **THEN** `d.Id()` becomes `<groupId>#fwmrule-aaaaaaaa`, `rule_id` is set to `fwmrule-aaaaaaaa`, and Read is executed

#### Scenario: Empty rule id returned

- **WHEN** `CreateSecurityGroupRule` returns an empty `RuleIds` slice
- **THEN** Create returns a non-retryable error and `d.Id()` is not set

#### Scenario: SDK error is retried then surfaced

- **WHEN** `CreateSecurityGroupRule` keeps failing
- **THEN** the failure is returned to Terraform after `WriteRetryTimeout` elapses and is logged with the `[CRITAL]` prefix

---

### Requirement: Read refreshes state from DescribeSecurityGroupRules by RuleId filter

`Read` SHALL parse the composite id, look the rule up through a new service helper `FwmService.DescribeFwmSecurityGroupRuleById`, and refresh state.

- The id MUST be split on `tccommon.FILED_SP`; anything other than two non-empty segments MUST produce `fmt.Errorf("id is broken,%s", d.Id())`.
- The helper MUST call `DescribeSecurityGroupRules` with `GroupId`, `Filters = [{Name: "RuleId", Values: [ruleId], OperatorType: 9}]`, `Offset = 0` and `Limit = 100`, inside `resource.Retry(tccommon.ReadRetryTimeout, ...)`.
- The helper MUST strict-equal the returned `RuleId` against the requested one and return `(nil, nil)` when nothing matches.
- When the helper returns `(nil, nil)`, Read MUST log `[CRUD] tencentcloud_fwm_security_group_rule id=%s` (keeping the current id visible in the log) and only then call `d.SetId("")`.
- Every `d.Set` MUST be preceded by a nil check on the corresponding response field.
- Read MUST NOT overwrite `source_type`, `dest_type`, `rule_action`, `description`, `source_content`, `dest_content`, `protocol_port_type` or `belong_member_id` (encoding drift, see design.md D6), and MUST overwrite `ip_version`, `protocol`, `port`, `service_template_id`, `scope`, `order_index` when present.

#### Scenario: Rule found

- **WHEN** the helper returns a `SecGroupRuleResp` for the requested rule id
- **THEN** `group_id`, `rule_id`, `protocol`, `port`, `scope`, `order_index`, and the computed attributes (including the `dns_parse_count` and `belong_member` blocks) are populated

#### Scenario: Rule deleted out of band

- **WHEN** the helper returns `(nil, nil)`
- **THEN** Read logs `[CRUD] tencentcloud_fwm_security_group_rule id=<id>` and sets the id to empty so Terraform plans a recreate

#### Scenario: Malformed id

- **WHEN** the state id does not contain exactly one `#` separator
- **THEN** Read returns an error containing `id is broken`

#### Scenario: Nil response fields are skipped

- **WHEN** the response omits `DnsParseCount` and `BelongMember`
- **THEN** Read does not call `d.Set` for `dns_parse_count` / `belong_member` and does not fail

---

### Requirement: Update modifies the rule through ModifySecurityGroupRule

`Update` SHALL send the changed rule body to `ModifySecurityGroupRule` and MUST always carry the rule identity.

- Update MUST short-circuit and skip the API call when none of `ip_version`, `source_content`, `source_type`, `dest_content`, `dest_type`, `protocol`, `port`, `service_template_id`, `rule_action`, `description`, `order_index`, `scope`, `protocol_port_type`, `belong_member_id` changed.
- The `ModifySecurityGroupRuleRequest` MUST always set `GroupId` (from the parsed id) and `Rule` with `RuleId` populated.
- The SDK call MUST be wrapped in `resource.Retry(tccommon.WriteRetryTimeout, ...)`; a nil `Response` MUST yield a `resource.NonRetryableError`.
- Update MUST end by delegating to Read.

#### Scenario: Body change is applied

- **WHEN** `description` is changed from `"a"` to `"b"`
- **THEN** `ModifySecurityGroupRule` is invoked with `Rule.RuleId` set to the rule id from state and `Rule.Description = "b"`

#### Scenario: No body change

- **WHEN** only computed attributes differ
- **THEN** `ModifySecurityGroupRule` is not invoked and the resource returns to Read

#### Scenario: Identity attributes never reach Update

- **WHEN** `group_id` or `group_name` changes
- **THEN** Terraform replaces the resource instead of calling Update

---

### Requirement: Delete removes the rule and any rule group created by this resource

`Delete` SHALL call `DeleteSecurityGroupRule` with a single-element `RuleIds`, and SHALL delete the rule group only when the group was created by this resource.

- The rule delete MUST send `GroupId` plus `RuleIds: []*string{ruleId}` inside `resource.Retry(tccommon.WriteRetryTimeout, ...)`.
- When `group_name` is present in the configuration (meaning this resource created the group), Delete MUST additionally call `DeleteRuleGroup` with `GroupIds: []*string{groupId}`.
- When only `group_id` is present, `DeleteRuleGroup` MUST NOT be called.
- The rule-group delete MUST be best-effort: failures are logged with a `[WARN]` prefix and do not fail the destroy, so a group already removed out of band does not block `terraform destroy`.

#### Scenario: Rule created in a resource-owned group

- **WHEN** the resource was created with `group_name`
- **THEN** Delete calls `DeleteSecurityGroupRule` and then `DeleteRuleGroup` for the same `groupId`

#### Scenario: Rule created in a pre-existing group

- **WHEN** the resource was created with `group_id`
- **THEN** Delete calls only `DeleteSecurityGroupRule` and leaves the rule group untouched

#### Scenario: Group delete fails

- **WHEN** `DeleteRuleGroup` returns an error
- **THEN** the error is logged with `[WARN]` and Delete still returns success after the rule is gone

---

### Requirement: Import support through the composite id

The resource SHALL implement `schema.ImportStatePassthrough`, and the documented import syntax MUST use the two-segment composite id.

#### Scenario: Import an existing rule

- **WHEN** a user runs `terraform import tencentcloud_fwm_security_group_rule.example fwmrg-xxxxxxxx#fwmrule-yyyyyyyy`
- **THEN** the following Read resolves the rule by the `RuleId` filter and populates state

#### Scenario: Imported resource does not own its group

- **WHEN** an imported resource is destroyed
- **THEN** only the rule is deleted, because `group_name` is empty in the imported state

---

### Requirement: Documentation and unit tests

The change SHALL ship resource documentation and gomonkey-based unit tests.

- `tencentcloud/services/fwm/resource_tc_fwm_security_group_rule.md` MUST contain a one-line description mentioning Firewall Manager (FWM), an example usage block, and an `Import` section with the composite id.
- `tencentcloud/provider.md` MUST list `tencentcloud_fwm_security_group_rule` under a `Firewall Manager(FWM)` Resources section so `make doc` regenerates `website/docs/r/fwm_security_group_rule.html.markdown`.
- `tencentcloud/services/fwm/resource_tc_fwm_security_group_rule_test.go` MUST use `gomonkey` patches over a mock `tccommon.ProviderMeta` instead of the Terraform acceptance suite, covering create (both identity modes), read (found / not found), update and delete (both ownership modes).

#### Scenario: Doc generation input is complete

- **WHEN** `make doc` runs after the change
- **THEN** `website/docs/r/fwm_security_group_rule.html.markdown` is produced from the resource schema and the `.md` example

#### Scenario: Unit tests cover the lifecycle without network access

- **WHEN** the `fwm` package tests are compiled
- **THEN** they reference only gomonkey patches and canned SDK responses, with no `TF_ACC` requirement
