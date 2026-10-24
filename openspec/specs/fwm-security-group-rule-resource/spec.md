# fwm-security-group-rule-resource Specification

## Purpose
TBD - created by archiving change add-fwm-security-group-rule-resource. Update Purpose after archive.
## Requirements
### Requirement: Resource MUST be registered as `tencentcloud_fwm_security_group_rule`

The provider SHALL register a new resource named `tencentcloud_fwm_security_group_rule` whose Create / Read / Update / Delete callbacks invoke the `CreateSecurityGroupRule` / `DescribeSecurityGroupRules` / `ModifySecurityGroupRule` / `DeleteSecurityGroupRule` APIs of `tencentcloud-sdk-go/tencentcloud/fwm/v20250611` respectively. All four APIs are synchronous (no FlowId/TaskId/JobId, no intermediate state); the resource MUST NOT implement waitReady/async polling.

The provider SHALL also register the FWM service client: `tencentcloud/connectivity/client.go` MUST add the `fwmv20250611` import, the `fwmv20250611Conn *fwmv20250611.Client` field and the `UseFwmV20250611Client()` accessor (mirroring `UseIgtmV20231024Client`).

#### Scenario: Resource registered in provider map

- **WHEN** the provider is loaded
- **THEN** `tencentcloud/provider.go` exposes the resource via key `"tencentcloud_fwm_security_group_rule"` mapped to `fwm.ResourceTencentCloudFwmSecurityGroupRule()`.

#### Scenario: FWM client accessor exists

- **WHEN** a CRUD callback obtains the FWM client
- **THEN** it calls `meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseFwmV20250611Client()` which lazily initializes and returns a `*fwmv20250611.Client`.

#### Scenario: Resource appears in gendoc index

- **WHEN** `tencentcloud/provider.md` is scanned by `make doc`
- **THEN** the Firewall(FWM) Resource section includes `tencentcloud_fwm_security_group_rule` so that `website/docs/r/fwm_security_group_rule.html.markdown` is generated.

### Requirement: Schema MUST be flat with required/forceNew constraints per FWM API semantics

The resource schema SHALL declare top-level argument keys (flat, no extra nesting layer wrapping the rule fields as a list/object). Writable input keys and their constraints:

| HCL key | SDK request field (Create `Rules[0]` / Modify `Rule`) | Type | Required | ForceNew |
|---|---|---|---|---|
| `group_id` | `request.GroupId` | TypeString | Yes | No |
| `ip_version` | `IpVersion` | TypeString | Yes | **Yes** |
| `order_index` | `OrderIndex` | TypeInt | Yes | No |
| `source_type` | `SourceType` | TypeString | Yes | **Yes** |
| `source_content` | `SourceContent` | TypeString | Yes | No |
| `dest_type` | `DestType` | TypeString | Yes | **Yes** |
| `dest_content` | `DestContent` | TypeString | Yes | No |
| `protocol` | `Protocol` | TypeString | Yes | No |
| `port` | `Port` | TypeString | Yes | No |
| `rule_action` | `RuleAction` | TypeString | Yes | No |
| `scope` | `Scope` | TypeString | No | **Yes** |
| `service_template_id` | `ServiceTemplateId` | TypeString | No | No |
| `description` | `Description` | TypeString | No | No |
| `protocol_port_type` | `ProtocolPortType` | TypeInt | No | No |
| `belong_member_id` | `BelongMemberId` | TypeString | No | No |

Read-only Computed keys (server display fields, MUST NOT be mapped into comparable writable fields): `rule_id`, `source_id`, `target_id`, `strategy`, `detail`, `region`, `protocol_port_name`, `sou_instance_name`, `sou_public_ip`, `sou_private_ip`, `sou_cidr`, `sou_parameter_name`, `instance_name`, `public_ip`, `private_ip`, `cidr`, `parameter_name`, `dns_parse_count` (nested `valid_count`/`invalid_count`), `is_new`, `id`, `all_total_count`.

#### Scenario: Required fields enforce on plan

- **WHEN** the user writes a config omitting `group_id`, `ip_version`, `order_index`, `source_type`, `source_content`, `dest_type`, `dest_content`, `protocol`, `port` or `rule_action`
- **THEN** `terraform plan` SHALL fail validation pointing at the missing required attribute.

#### Scenario: Changing forceNew fields forces replacement

- **GIVEN** an existing rule in state
- **WHEN** the user changes `ip_version`, `source_type`, `dest_type` or `scope` in HCL
- **THEN** Terraform's plan reports a destroy + create cycle (resource is replaced).

#### Scenario: Flat schema without list wrapper

- **WHEN** a code reviewer inspects the schema
- **THEN** the writable rule fields are top-level keys; there is NO `rules`/`rule` nesting layer, and Create builds exactly one element in the SDK `Rules` array.

### Requirement: Create MUST call CreateSecurityGroupRule with a single-element Rules array and take RuleIds[0]

The Create callback SHALL build a `CreateSecurityGroupRuleRequest` with `GroupId` and `Rules` containing exactly one `*SecurityGroupRule` populated from the schema, call `CreateSecurityGroupRuleWithContext(ctx, request)` inside `resource.Retry(tccommon.WriteRetryTimeout, ...)`, and validate the response: if `result == nil`, `result.Response == nil`, `len(result.Response.RuleIds) == 0`, `RuleIds[0] == nil` or `*RuleIds[0] == ""`, the callback MUST return `resource.NonRetryableError` after logging the current `logId` (the id must NOT be written to state in this failure case). On success the resource ID SHALL be `strings.Join([]string{groupId, ruleIds[0]}, tccommon.FILED_SP)` (FILED_SP = `#`). Create MUST NOT implement async polling; on success it ends by calling Read.

#### Scenario: Create produces one API call and a compound id

- **GIVEN** HCL with all required fields (`group_id`, `ip_version = "ipv4"`, `order_index = 1`, `source_type = "net"`, `source_content = "1.1.1.1"`, `dest_type = "net"`, `dest_content = "2.2.2.2"`, `protocol = "TCP"`, `port = "80"`, `rule_action = "accept"`, `scope = "SG"`)
- **WHEN** the user runs `terraform apply`
- **THEN** exactly one `CreateSecurityGroupRule` request is issued whose `Rules` has length 1, and `d.Id()` returns `<GroupId>#<RuleIds[0]>`.

#### Scenario: Empty RuleIds is a non-retryable failure

- **GIVEN** the API responds with `Response.RuleIds` empty or nil
- **WHEN** the retry callback runs
- **THEN** it returns `resource.NonRetryableError` with a clear message and does NOT `d.SetId` with an empty rule id.

#### Scenario: Transient SDK error is retried

- **GIVEN** the first invocation of `CreateSecurityGroupRule` returns a retriable error
- **WHEN** the retry callback re-runs
- **THEN** the original error is surfaced only after `tccommon.WriteRetryTimeout` retry budget is exhausted.

### Requirement: Read MUST query by RuleId filter and convert int enums back to input strings

The Read callback SHALL split `d.Id()` by `tccommon.FILED_SP` into `groupId`/`ruleId`, then call (via the FWM service layer) `DescribeSecurityGroupRulesWithContext` with `GroupId`, `Filters = [{Name: "RuleId", Values: [ruleId], OperatorType: 1}]` (1 = EQUAL precise match), wrapped in `resource.Retry(tccommon.ReadRetryTimeout, ...)`. Inside the retry, a nil `result`/`result.Response` MUST return `resource.NonRetryableError`. Outside the retry, if the returned `Rules` list is empty, the callback MUST log `log.Printf("[CRUD] fwm_security_group_rule id=%s", d.Id())` to preserve the id context and then `d.SetId("")` and return nil (rule deleted externally; removal from state, not an error).

For the first element of `Rules`, Read SHALL set state fields with reverse-mapping of response int enums to input strings so that plan does not permanently diff:

- `SourceType`/`TargetType` (int) → `source_type`/`dest_type` (string): net=0, vpc=1, subnet=2, cvm=3, clb=4, eni=5, cdb=6, template=7, tag=8, region=9, resourcegroup=100. If the int is not in the map, keep the decimal string form of the int and log a warning (no panic).
- `Strategy` (int) → `rule_action` (string): accept=2, drop=1, log=0.
- `Detail` → `description`; `SourceId` → `source_content` (and also to read-only `source_id`); `TargetId` → `dest_content` (and also to read-only `target_id`); `BelongMember.MemberId` → `belong_member_id` only when `BelongMember` is non-nil and `MemberId` is non-empty.
- Every `d.Set(...)` MUST be guarded by a nil check on the corresponding response field.
- `Id` (int64 auto-increment primary key with overflow risk) MUST be Computed only; state and import use `RuleId` (uuid) exclusively.

#### Scenario: Read filters by RuleId with EQUAL operator

- **WHEN** Read executes
- **THEN** the `DescribeSecurityGroupRules` request contains exactly one filter `{Name: "RuleId", Values: [<ruleId>], OperatorType: 1}` plus `GroupId`, and the first element of `Rules` is used to populate state.

#### Scenario: Externally deleted rule is removed from state without error

- **GIVEN** state contains a rule whose RuleId no longer exists on the server
- **WHEN** Read runs and `Response.Rules` is empty
- **THEN** the callback logs `[CRUD] fwm_security_group_rule id=<id>` first, then sets `d.SetId("")` and returns nil (no error).

#### Scenario: int enums are reverse-mapped to strings in state

- **GIVEN** the server returns a rule with `SourceType = 0`, `TargetType = 0`, `Strategy = 2`
- **WHEN** Read populates state
- **THEN** `source_type = "net"`, `dest_type = "net"`, `rule_action = "accept"`, so a subsequent `terraform plan` reports no diff.

#### Scenario: Unknown int enum does not panic

- **GIVEN** the server returns `SourceType = 55` (value outside the map)
- **WHEN** Read populates state
- **THEN** `source_type` is set to the string `"55"`, a warning is logged, and no panic occurs.

### Requirement: Update MUST call ModifySecurityGroupRule with full-coverage semantics

The Update callback SHALL detect changes among the mutable args (`order_index`, `source_content`, `dest_content`, `protocol`, `port`, `service_template_id`, `rule_action`, `description`, `protocol_port_type`, `belong_member_id`) via `d.HasChange`. When any change is detected, it SHALL build a `ModifySecurityGroupRuleRequest` with `GroupId` (split from `d.Id()`) and `Rule` (a single `*SecurityGroupRule` containing `RuleId` plus ALL writable fields re-populated from the current config, because the API semantics are full coverage — not incremental). The call `ModifySecurityGroupRuleWithContext(ctx, request)` MUST be wrapped in `resource.Retry(tccommon.WriteRetryTimeout, ...)` with nil-response guards. `GroupId` MUST NOT be modifiable (it is not in the mutable args). `ip_version`, `source_type`, `dest_type` and `scope` are ForceNew so they never reach Update. The unused-by-server `ModifyType` field MUST NOT be relied upon. Update ends by calling Read.

#### Scenario: Editing a mutable field triggers a single full Modify call

- **GIVEN** state has `port = "80"`, `description = "old"`
- **WHEN** the user changes `port` to `"443"`
- **THEN** Update issues exactly one `ModifySecurityGroupRule` request whose `Rule` contains `RuleId` and ALL writable fields (`order_index`, `source_content`, `dest_content`, `protocol`, `port = "443"`, `rule_action`, `description = "old"`, and optional fields when set).

#### Scenario: No mutable change produces no Modify call

- **GIVEN** only a ForceNew field changed (handled by replace) or nothing changed
- **WHEN** Update runs with no `d.HasChange` hit among mutable args
- **THEN** no `ModifySecurityGroupRule` request is issued and Read is invoked.

### Requirement: Delete MUST call DeleteSecurityGroupRule with one RuleId and treat ResourceNotFound as success

The Delete callback SHALL split `d.Id()` into `groupId`/`ruleId`, build a `DeleteSecurityGroupRuleRequest` with `GroupId` and `RuleIds` containing exactly one element (the single rule id), and call `DeleteSecurityGroupRuleWithContext(ctx, request)` wrapped in `resource.Retry(tccommon.WriteRetryTimeout, ...)`. When the API returns `ResourceNotFound` (rule already deleted), the callback SHALL treat it as success (idempotent destroy) and return nil; all other errors go through `tccommon.RetryError(e)`.

#### Scenario: Delete issues one batch call with a single element

- **GIVEN** state has the compound id `<GroupId>#<RuleId>`
- **WHEN** the user runs `terraform destroy`
- **THEN** exactly one `DeleteSecurityGroupRule` request is issued with `RuleIds = [<RuleId>]` and length 1.

#### Scenario: ResourceNotFound during destroy is success

- **GIVEN** the rule was already deleted out-of-band
- **WHEN** Delete runs and receives a `ResourceNotFound` error
- **THEN** the callback returns nil (destroy succeeds without retry storm).

### Requirement: Optional int fields MUST use GetOkExists to avoid zero-value drift

For optional numeric fields (`order_index` is required; `protocol_port_type` is optional int), the code SHALL read values with `d.GetOkExists(...)` so that a legitimate `0`/`false` value is not treated as "unset" (avoiding state drift). Optional string fields may use `d.GetOk(...)`. All error returns MUST be checked; a function that cannot fail is called as `_ = fn(...)`.

#### Scenario: protocol_port_type zero is preserved

- **GIVEN** HCL sets `protocol_port_type = 0`
- **WHEN** Create/Update builds the request
- **THEN** the SDK field `ProtocolPortType` is populated with 0 (not omitted), so plan does not drift.

### Requirement: Import MUST support the compound id `GroupId#RuleId`

The resource SHALL declare `Importer: &schema.ResourceImporter{State: schema.ImportStatePassthrough}`. The resource markdown Import section MUST document that the import id is the compound `GroupId#RuleId` (with `#` being `tccommon.FILED_SP`), e.g. `terraform import tencentcloud_fwm_security_group_rule.example <group-id>#<rule-id>`.

#### Scenario: Import passthrough works with compound id

- **GIVEN** an existing rule on the server
- **WHEN** the user runs `terraform import tencentcloud_fwm_security_group_rule.example <GroupId>#<RuleId>`
- **THEN** `d.Id()` equals the provided string, Read splits it into `GroupId`/`RuleId` and populates state from the server.

### Requirement: Documentation and unit test MUST follow project conventions

- The resource markdown SHALL live at `tencentcloud/services/fwm/resource_tc_fwm_security_group_rule.md` with: a one-line description mentioning the cloud product name (FWM), an `Example Usage` HCL block (values: `ip_version = "ipv4"`, `order_index = 1`, `source_type = "net"`, `source_content = "1.1.1.1"`, `dest_type = "net"`, `dest_content = "2.2.2.2"`, `protocol = "TCP"`, `port = "80"`, `rule_action = "accept"`, `scope = "SG"`), and an `Import` section documenting the compound `GroupId#RuleId` id. It MUST NOT contain `Argument Reference` / `Attribute Reference` sections (tool-generated).
- The unit test SHALL live at `tencentcloud/services/fwm/resource_tc_fwm_security_group_rule_test.go`, package `fwm_test`, and use gomonkey mocks (NOT the Terraform acceptance test suite) to mock `CreateSecurityGroupRuleWithContext` / `DescribeSecurityGroupRulesWithContext` / `ModifySecurityGroupRuleWithContext` / `DeleteSecurityGroupRuleWithContext`, covering Create / Read / Update / Delete business logic with the testValue set above.
- Logging/error messages MUST use the resource name `fwm_security_group_rule` (not vague wording like "the resource").
- No `_extension.go` file SHALL be generated.

#### Scenario: Unit test mocks the four FWM APIs

- **WHEN** the test file is opened
- **THEN** the package is `fwm_test`, gomonkey `ApplyMethod` (or equivalent) replaces the four `fwmv20250611.Client` context methods, and each CRUD callback's business logic (request build, id join/split, enum reverse-mapping, empty-response handling) is asserted without real API calls.

#### Scenario: Markdown contains import with compound id

- **WHEN** the resource markdown is opened
- **THEN** the Import section shows a `terraform import tencentcloud_fwm_security_group_rule.example <group-id>#<rule-id>` example and explains that both `GroupId` and `RuleId` are required in the compound id.

