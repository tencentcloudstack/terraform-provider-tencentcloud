# fwm-vpc-acl-rule-resource Specification

## Purpose
TBD - created by archiving change add-fwm-vpc-acl-rule-resource. Update Purpose after archive.
## Requirements
### Requirement: Resource MUST be registered as `tencentcloud_fwm_vpc_acl_rule`

The provider SHALL register a new CRUD-type resource (RESOURCE_KIND_GENERAL) named `tencentcloud_fwm_vpc_acl_rule` whose Create/Read/Update/Delete callbacks invoke the FWM VPC ACL rule APIs (`CreateVpcAclRule`, `DescribeVpcAclRules`, `ModifyVpcAclRule`, `DeleteVpcAclRule`) of `tencentcloud-sdk-go/tencentcloud/fwm/v20250611`. The `connectivity.TencentCloudClient` SHALL expose a `UseFwmV20250611Client()` accessor so the fwm service layer can obtain the SDK client.

#### Scenario: Resource registered in provider map

- **WHEN** the provider is loaded
- **THEN** `provider.go` exposes the resource via key `"tencentcloud_fwm_vpc_acl_rule"` mapped to `fwm.ResourceTencentCloudFwmVpcAclRule()`, alongside other `tencentcloud_fwm_*` resource entries.

#### Scenario: FWM client accessor available

- **WHEN** the fwm service layer calls `meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseFwmV20250611Client()`
- **THEN** a `*fwmv20250611.Client` is returned (lazily initialized on first call).

#### Scenario: Importer is configured

- **WHEN** a user runs `terraform import tencentcloud_fwm_vpc_acl_rule.example GroupId#RuleId`
- **THEN** the importer SHALL pass the supplied `GroupId#RuleId` straight into the Read callback (using `schema.ImportStatePassthrough`), and Read SHALL hydrate state from `DescribeVpcAclRules` filtered by that GroupId + RuleId.

### Requirement: Schema MUST mirror the VpcAclRule API input with flat top-level fields

The resource schema SHALL declare exactly these top-level argument keys, with semantics matching the SDK fields. Rule fields are flattened to the top level (no enclosing `rules` list layer), mirroring `tencentcloud_cfw_vpc_policy`.

| HCL key | SDK input field | SDK output field | Type | Required | ForceNew | Computed |
|---|---|---|---|---|---|---|
| `group_id` | `request.GroupId` | — | TypeString | Yes | Yes | No |
| `edge_id` | `Rules.EdgeId` | `Rules.EdgeId` | TypeString | Yes | Yes | No |
| `fw_group_id` | `Rules.FwGroupId` | `Rules.FwGroupId` | TypeString | Yes | Yes | No |
| `ip_version` | — (input has no field) | `Rules.IpVersion` | TypeString | Yes | Yes | No |
| `source_type` | `Rules.SourceType` | `Rules.SourceType` | TypeString | Yes | Yes | No |
| `source_content` | `Rules.SourceContent` | `Rules.SourceContent` | TypeString | Yes | Yes | No |
| `dest_type` | `Rules.DestType` | `Rules.DestType` | TypeString | Yes | Yes | No |
| `dest_content` | `Rules.DestContent` | `Rules.DestContent` | TypeString | Yes | Yes | No |
| `protocol` | `Rules.Protocol` | `Rules.Protocol` | TypeString | Yes | No | No |
| `rule_action` | `Rules.RuleAction` | `Rules.RuleAction` | TypeString | Yes | No | No |
| `order_index` | `Rules.OrderIndex` | `Rules.Sequence` | TypeInt | Yes | No | No |
| `description` | `Rules.Description` | `Rules.Description` | TypeString | Yes | No | No |
| `port` | `Rules.Port` | `Rules.Port` | TypeString | No | No | No |
| `param_template_id` | `Rules.ParamTemplateId` | `Rules.ParamTemplateId` | TypeString | No | No | No |
| `belong_member_id` | `Rules.BelongMemberId` | `Rules.BelongMember.MemberId` | TypeString | No | No | No |
| `rule_id` | `Rules.RuleId` (modify) | `Rules.RuleId` | TypeString | No | — | Yes |
| `sequence` | — | `Rules.Sequence` | TypeInt | — | — | Yes |
| `source_name` | — | `Rules.SourceName` | TypeString | — | — | Yes |
| `dest_name` | — | `Rules.DestName` | TypeString | — | — | Yes |
| `param_template_name` | — | `Rules.ParamTemplateName` | TypeString | — | — | Yes |
| `fw_group_name` | — | `Rules.FwGroupName` | TypeString | — | — | Yes |
| `invalid` | — | `Rules.Invalid` | TypeInt | — | — | Yes |
| `app_id` | — | `Rules.BelongMember.AppId` | TypeString | — | — | Yes |
| `uin` | — | `Rules.BelongMember.Uin` | TypeString | — | — | Yes |
| `nickname` | — | `Rules.BelongMember.Nickname` | TypeString | — | — | Yes |
| `member_id` | — | `Rules.BelongMember.MemberId` | TypeString | — | — | Yes |
| `node_name` | — | `Rules.BelongMember.NodeName` | TypeString | — | — | Yes |
| `create_time` | — | `Rules.CreateTime` | TypeString | — | — | Yes |
| `update_time` | — | `Rules.UpdateTime` | TypeString | — | — | Yes |
| `dns_parse_cnt` | — | `Rules.DnsParseCnt` | TypeInt | — | — | Yes |
| `hit_count` | — | `Rules.HitCount` | TypeInt | — | — | Yes |
| `last_hit_time` | — | `Rules.LastHitTime` | TypeString | — | — | Yes |
| `all_total_count` | — | `Response.AllTotalCount` | TypeInt | — | — | Yes |

The forceNew set is: `group_id`, `edge_id`, `fw_group_id`, `ip_version`, `source_type`, `source_content`, `dest_type`, `dest_content`. Changes to any other writable field (`protocol`, `rule_action`, `order_index`, `description`, `port`, `param_template_id`, `belong_member_id`) are in-place updates.

`ip_version` is Required + ForceNew. The Create input struct `VpcAclRule` has NO `IpVersion` field (server infers from rule group), so Create MUST NOT send it; it exists only as a schema field for distinguishing ipv4/ipv6 and is populated from the Describe response `Rules.IpVersion` during Read.

#### Scenario: ForceNew field change triggers replacement

- **GIVEN** state has `source_content = "10.0.0.0/24"`
- **WHEN** the user changes `source_content` to `"10.0.1.0/24"`
- **THEN** Terraform's plan reports a destroy + create cycle (resource is replaced because `source_content` is ForceNew). Update is NOT invoked.

#### Scenario: Writable field change triggers in-place update

- **GIVEN** state has `description = "old"`
- **WHEN** the user changes `description` to `"new"`
- **THEN** Update is invoked with a single `ModifyVpcAclRule` request carrying `RuleId`, all writable fields fully populated, and `GroupId`; no replacement occurs.

### Requirement: Resource ID MUST be the composite GroupId#RuleId

After Create, the resource's Terraform ID SHALL be set to `strings.Join([]string{groupId, ruleId}, tccommon.FILED_SP)` (i.e. `GroupId#RuleId`). Read/Update/Delete SHALL split `d.Id()` by `tccommon.FILED_SP` and require exactly 2 parts.

#### Scenario: ID is set from CreateVpcAclRule response

- **GIVEN** a successful `CreateVpcAclRule` returns `RuleIds = ["r-abc123"]` and the user passed `group_id = "g-001"`
- **WHEN** Create completes
- **THEN** `d.Id()` returns `"g-001#r-abc123"`.

#### Scenario: Create response missing RuleIds is fatal

- **WHEN** `CreateVpcAclRule` returns a non-nil response but `RuleIds` is nil or empty
- **THEN** the Create callback SHALL return a non-retryable error mentioning `RuleIds is empty`.

### Requirement: Create MUST call CreateVpcAclRule with a single-element Rules array

The Create callback SHALL build a `CreateVpcAclRuleRequest` with:
- `request.GroupId = group_id`
- `request.Rules = []{*VpcAclRule{...}}` (exactly 1 element, assembled from the flat top-level schema fields)

The call SHALL be wrapped in `resource.Retry(tccommon.WriteRetryTimeout, ...)` with `tccommon.RetryError(e)`. After the retry succeeds, the callback SHALL validate `Response.RuleIds` is non-empty, take `RuleIds[0]`, then set `d.SetId(GroupId#RuleId)`, then invoke Read. The `d.SetId` MUST happen after the retry block but before any further logic. Create MUST NOT send `IpVersion` (the input struct has no such field).

#### Scenario: New required field appears in plan

- **WHEN** the user writes a config omitting a required field such as `group_id` or `source_content`
- **THEN** `terraform plan` SHALL fail validation pointing at the missing required attribute.

### Requirement: Read MUST locate the resource via DescribeVpcAclRules + RuleId filter

The service-layer helper `(*FwmService).DescribeFwmVpcAclRuleById(ctx, groupId, ruleId)` SHALL build a `DescribeVpcAclRulesRequest` with:
- `GroupId = groupId`
- `Filters = [{Name: "RuleId", Values: [ruleId], OperatorType: 1}]` (1 = EQUAL exact match)
- `Limit = 0` (use default 100)
- `Offset = 0`

…and invoke `UseFwmV20250611Client().DescribeVpcAclRules` inside `resource.Retry(tccommon.ReadRetryTimeout, ...)`.

If `Response.Rules` is empty, the helper SHALL return `(nil, nil, nil)`. The Read callback SHALL detect that case, log `[CRUD] fwm_vpc_acl_rule id=<d.Id()>` to preserve context, then call `d.SetId("")` and return nil.

Read field mapping:
- `order_index` ← `Rules[0].Sequence` (NOT OrderIndex; the response has no OrderIndex)
- `belong_member_id` ← `Rules[0].BelongMember.MemberId` (BelongMember is a `MemberInfo` object; only set when BelongMember != nil and MemberId != nil)
- `app_id`, `uin`, `nickname`, `member_id`, `node_name` ← flattened from `Rules[0].BelongMember.*` (only when BelongMember != nil and the inner field != nil)
- `all_total_count` ← `Response.AllTotalCount`

#### Scenario: Resource exists

- **GIVEN** the API holds a rule with `RuleId="r-abc"`
- **WHEN** Read runs
- **THEN** all schema attributes are populated from the API response with nil-safe pointer dereferencing.

#### Scenario: Resource removed out-of-band

- **GIVEN** the user runs `terraform refresh` after the rule was deleted externally
- **WHEN** Read calls `DescribeVpcAclRules` and `len(Response.Rules) == 0`
- **THEN** Read SHALL log `[CRUD] fwm_vpc_acl_rule id=<id>`, call `d.SetId("")`, and return `nil` (no error), so the next plan proposes a re-create.

### Requirement: Update MUST call ModifyVpcAclRule as a full overwrite

ModifyVpcAclRule is a full overwrite (not incremental) and reuses the same field validation as Create. The Update callback SHALL detect changes against the mutable field set `["protocol", "rule_action", "order_index", "description", "port", "param_template_id", "belong_member_id"]`. When any has changed, Update SHALL build a single `ModifyVpcAclRuleRequest` with:
- `request.GroupId = group_id`
- `request.Rule = &VpcAclRule{...}` (single object, NOT array)
- The `VpcAclRule` SHALL include `RuleId = ruleId` and ALL writable fields populated from current `d.Get(...)` (full overwrite)

The call SHALL be wrapped in `resource.Retry(tccommon.WriteRetryTimeout, ...)` with `tccommon.RetryError(e)`. Update MUST end by re-invoking Read. GroupId is NOT modifiable (it is part of the ID, never sent as a changed value).

#### Scenario: Editing only description

- **GIVEN** state has `description = "old"`
- **WHEN** the user changes `description` to `"new"`
- **THEN** Update issues exactly one `ModifyVpcAclRule` request with `RuleId`, `GroupId`, and the full set of writable field values (including unchanged ones), ending with `description = "new"`.

### Requirement: Delete MUST call DeleteVpcAclRule and be idempotent on ResourceNotFound

The Delete callback SHALL build a `DeleteVpcAclRuleRequest` with `GroupId = group_id` and `RuleIds = [ruleId]` (1 element), and call `UseFwmV20250611Client().DeleteVpcAclRule` inside `resource.Retry(tccommon.WriteRetryTimeout, ...)`. If the API returns ResourceNotFound (the rule was already deleted), Delete SHALL treat it as success and return nil (idempotent destroy).

#### Scenario: Idempotent delete on a missing resource

- **GIVEN** the rule is already gone in the cloud
- **WHEN** Delete is called
- **THEN** the ResourceNotFound error is swallowed and Delete returns nil, completing the destroy without error.

### Requirement: Every API call MUST be wrapped in resource.Retry

Every invocation of `UseFwmV20250611Client().CreateVpcAclRule`, `DescribeVpcAclRules`, `ModifyVpcAclRule`, `DeleteVpcAclRule` SHALL be wrapped in `resource.Retry(...)` with the appropriate timeout (`ReadRetryTimeout` for Describe, `WriteRetryTimeout` for Create/Modify/Delete) and SHALL forward errors via `tccommon.RetryError(e)` so transient errors are retried. The service layer SHALL call `ratelimit.Check(request.GetAction())` inside the retry block. Setting `d.SetId` and other success-side operations SHALL happen OUTSIDE the retry block (after retry error handling).

#### Scenario: Transient SDK error is retried

- **GIVEN** the first invocation returns a retriable error
- **WHEN** the retry callback re-runs
- **THEN** the second attempt's response is observed; the original error surfaces only after the retry budget is exhausted.

### Requirement: All response field reads MUST be nil-safe

Every `_ = d.Set("<key>", respData.<Field>)` access SHALL be guarded by `if respData.<Field> != nil` before dereference. Nested struct fields (e.g. `respData.BelongMember.MemberId`) SHALL check the outer pointer (BelongMember) then the inner field before reading.

#### Scenario: Optional string field absent in response

- **GIVEN** the API returns `Port == nil`
- **WHEN** Read runs
- **THEN** `d.Set("port", ...)` is NOT called for that field; state retains whatever the previous value was (or stays unset for fresh imports).

### Requirement: Documentation and unit test MUST follow project conventions

- The HCL example markdown SHALL live at `tencentcloud/services/fwm/resource_tc_fwm_vpc_acl_rule.md`, containing a one-line description (with FWM product name), an Example Usage block (using the testValue config), and an Import section explaining the `GroupId#RuleId` composite ID.
- A unit test SHALL live at `tencentcloud/services/fwm/resource_tc_fwm_vpc_acl_rule_test.go` using gomonkey mocking (NOT terraform acceptance test suite), covering Create/Read/Update/Delete/Schema with mock SDK responses.
- Running `make doc` SHALL regenerate `website/docs/r/fwm_vpc_acl_rule.html.markdown` (handled in finalize stage).

#### Scenario: Generated website doc lists the resource

- **WHEN** `make doc` runs
- **THEN** `website/docs/r/fwm_vpc_acl_rule.html.markdown` exists and lists every schema attribute defined in the Schema requirement above.

#### Scenario: Unit test structure

- **WHEN** the test file is opened
- **THEN** the package is `fwm_test`, it imports `gomonkey/v2` and `testify/assert`, mocks `UseFwmV20250611Client` via `ApplyMethodReturn`, mocks the four CRUD SDK methods via `ApplyMethodFunc`, and asserts on `d.Id()` and key state attributes.

