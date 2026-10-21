# fwm-rule-group-resource Specification

## ADDED Requirements

### Requirement: Resource MUST be registered as `tencentcloud_fwm_rule_group`

The provider SHALL register a new RESOURCE_KIND_GENERAL resource named `tencentcloud_fwm_rule_group`, whose Create / Read / Update / Delete callbacks invoke `CreateRuleGroup`, `DescribeRuleGroups`, `ModifyRuleGroup`, `DeleteRuleGroup` of `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611`. The resource SHALL live in the new service package `tencentcloud/services/fwm/` with file name `resource_tc_fwm_rule_group.go`.

#### Scenario: Resource registered in provider map

- **WHEN** the provider is loaded
- **THEN** `tencentcloud/provider.go` exposes the key `"tencentcloud_fwm_rule_group"` mapped to `fwm.ResourceTencentCloudFwmRuleGroup()`, and `tencentcloud/provider.md` lists `tencentcloud_fwm_rule_group` under a `Firewall Manager(FWM)` resource section.

#### Scenario: Client factory is available

- **WHEN** any CRUD callback needs the FWM client
- **THEN** `(*connectivity.TencentCloudClient).UseFwmV20250611Client()` returns a cached `*fwmv20250611.Client` created via `fwmv20250611.NewClient(me.Credential, me.Region, me.NewClientProfile(300))` with `WithHttpTransport(&LogRoundTripper{})`.

#### Scenario: Importer is configured

- **WHEN** a user runs `terraform import tencentcloud_fwm_rule_group.example fwmrg-wsrrhuvrib`
- **THEN** the importer SHALL pass the imported id straight into Read (`schema.ImportStatePassthrough`), and Read SHALL hydrate state from `DescribeRuleGroups` filtered by `GroupId`.

### Requirement: Schema MUST expose only writable cloud API fields plus flattened computed fields

The resource schema SHALL declare `group_name` (`TypeString`, Required), `product` (`TypeString`, Required, **ForceNew**, validated against `enterprise_sg` / `cfw_nat_acl` / `cfw_vpc_acl` / `cfw_edge_acl`), and the following Computed fields flattened from `RuleGroupResp` with NO `rule_groups` wrapper layer: `group_id` (String), `app_id` (String), `rule_count` (Int), `create_time` (String), `update_time` (String), `create_by` (String), `create_name` (String), `update_by` (String), `update_name` (String), `receive_account` (List of object), `exec_areas` (List of String), `direction` (Int).

`receive_account` SHALL contain `uin` (String), `nickname` (String), `receiver_type` (Int) and `members` (List of object); each `members` element SHALL contain `app_id`, `uin`, `nickname`, `member_id`, `node_name` (all String).

The query-side parameters of `DescribeRuleGroups` (`Product` as a filter, `Filters`, `Filters.Name`, `Filters.Values`, `Filters.OperatorType`, `Order`, `By`, `Offset`, `Limit`) SHALL NOT be exposed in the resource schema.

#### Scenario: Writable fields are required and typed correctly

- **WHEN** the schema is inspected
- **THEN** `group_name` is Required `TypeString` without ForceNew, and `product` is Required `TypeString` with `ForceNew: true` and a `validation.StringInSlice` validator over the 4 allowed product names.

#### Scenario: Invalid product is rejected at plan time

- **WHEN** a configuration sets `product = "all"` or any value outside the 4 allowed names
- **THEN** `terraform plan` fails schema validation before any cloud API call.

#### Scenario: Computed results are flattened

- **WHEN** the schema is inspected
- **THEN** there is no `rule_groups` (or equivalently named list-of-resource) attribute, and every `RuleGroupResp` scalar field appears as its own top-level Computed attribute.

#### Scenario: Query-only parameters are not exposed

- **WHEN** the schema is inspected
- **THEN** none of `filters`, `name`, `values`, `operator_type`, `order`, `by`, `offset`, `limit` exist as resource attributes.

### Requirement: Resource ID MUST be the cloud API `GroupId`

Create SHALL set `d.Id()` to the value of `response.Response.GroupId`; Read / Update / Delete SHALL read the identifier from `d.Id()` and use it directly as `GroupId` / `GroupIds[0]`. No composite id (no `tccommon.FILED_SP` splitting) is used because `GroupId` is globally unique.

#### Scenario: ID set from CreateRuleGroup response

- **GIVEN** a successful `CreateRuleGroup` returning `GroupId = "fwmrg-wsrrhuvrib"`
- **WHEN** Create completes
- **THEN** `d.Id()` equals `"fwmrg-wsrrhuvrib"` and the `group_id` attribute also equals `"fwmrg-wsrrhuvrib"`.

#### Scenario: Empty Create response id is fatal

- **WHEN** `CreateRuleGroup` returns `result == nil`, or `result.Response == nil`, or `result.Response.GroupId == nil`, or `*GroupId == ""`
- **THEN** Create SHALL return `resource.NonRetryableError` whose message mentions the missing `GroupId`, and `d.SetId` SHALL NOT be called with an empty value.

### Requirement: Read MUST locate the group via DescribeRuleGroups with a fixed GroupId filter

The service helper `(*FwmService).DescribeFwmRuleGroupById(ctx context.Context, groupId string) (*fwmv20250611.RuleGroupResp, error)` SHALL call `DescribeRuleGroups` with `Filters = [{Name: "GroupId", Values: [<groupId>], OperatorType: 9}]` and paginate via `Offset` / `Limit`, returning the entry whose `GroupId` exactly equals the requested id, or `(nil, nil)` when absent. The request SHALL NOT use `Product` as a filter.

#### Scenario: Read hydrates all fields

- **GIVEN** `DescribeRuleGroups` returns one `RuleGroupResp` with `GroupId == d.Id()`
- **WHEN** Read completes
- **THEN** `group_name`, `product`, `group_id`, `app_id`, `rule_count`, `create_time`, `update_time`, `create_by`, `create_name`, `update_by`, `update_name`, `receive_account`, `exec_areas` and `direction` are written to state, and each `d.Set` call is guarded by a nil check on the corresponding SDK field.

#### Scenario: Missing group clears state after logging

- **GIVEN** `DescribeRuleGroups` returns zero matching entries
- **WHEN** Read completes
- **THEN** Read logs `[CRUD] tencentcloud_fwm_rule_group id=<id>` **before** calling `d.SetId("")`, and returns `nil`.

#### Scenario: Product is normalized from numeric to name

- **GIVEN** the cloud API returns `Product = "2"` for the rule group
- **WHEN** Read completes
- **THEN** the `product` attribute in state equals `"cfw_nat_acl"` (mapping `1→enterprise_sg`, `2→cfw_nat_acl`, `3→cfw_vpc_acl`, `4→cfw_edge_acl`), and a value that is already a product name or is unmapped is stored unchanged.

### Requirement: Update MUST only modify `group_name` through ModifyRuleGroup

Update SHALL treat `group_name` as the only mutable argument. When `d.HasChange("group_name")` is true it SHALL call `ModifyRuleGroup` with `GroupId = d.Id()` and `GroupName` from state, then re-run Read. If a ForceNew-only argument such as `product` is detected as changed inside Update, Update SHALL return an error instead of silently ignoring it. Update SHALL NOT send `Product`.

#### Scenario: Renaming a rule group

- **GIVEN** an existing rule group with `group_name = "old-name"` and the configuration changed to `"new-name"`
- **WHEN** Update runs
- **THEN** exactly one `ModifyRuleGroup` call is issued with `GroupId` equal to `d.Id()` and `GroupName = "new-name"`, followed by a Read.

#### Scenario: Product change never reaches Update

- **WHEN** a user changes `product` in configuration
- **THEN** Terraform plans a destroy + create cycle because `product` is `ForceNew`, so Update is not invoked with a product change.

### Requirement: Delete MUST call DeleteRuleGroup with a single-element GroupIds list

Delete SHALL call `DeleteRuleGroup` with `GroupIds = []*string{&d.Id()}` and treat a non-error response as completion, without additional polling.

#### Scenario: Single group deletion

- **GIVEN** `d.Id() == "fwmrg-wsrrhuvrib"`
- **WHEN** Delete runs
- **THEN** the request contains exactly one element `"fwmrg-wsrrhuvrib"` in `GroupIds` and no error is returned on success.

### Requirement: Every cloud API call MUST be retry-wrapped with empty-response protection

All Create / Update / Delete calls SHALL be wrapped in `resource.Retry(tccommon.WriteRetryTimeout, ...)` returning `tccommon.RetryError(e)` on error, and the Read helper SHALL be wrapped in `resource.Retry(tccommon.ReadRetryTimeout, ...)`. Nil `result` or nil `Response` SHALL produce `resource.NonRetryableError`. No polling / `Timeouts` block is required because all four APIs are synchronous.

#### Scenario: Transient error is retried

- **WHEN** the SDK returns a retryable error from `CreateRuleGroup`, `ModifyRuleGroup` or `DeleteRuleGroup`
- **THEN** the call is retried within `tccommon.WriteRetryTimeout` rather than failing immediately.

#### Scenario: Nil response is not retried forever

- **WHEN** the SDK returns `nil` response or `nil` `Response` body
- **THEN** the callback returns `resource.NonRetryableError` so the retry loop stops.

### Requirement: Resource documentation and unit tests MUST be provided

The change SHALL add `tencentcloud/services/fwm/resource_tc_fwm_rule_group.md` with a one-line description mentioning FWM, an `Example Usage` HCL block and an `Import` section, without `Argument Reference` / `Attribute Reference` (generated by `make doc`). It SHALL add `tencentcloud/services/fwm/resource_tc_fwm_rule_group_test.go` using `gomonkey` mocks (not the terraform acceptance test harness) covering Create, Read (including product normalization and not-found), Update and Delete.

#### Scenario: Example file drives generated docs

- **WHEN** `make doc` runs in the finalize phase
- **THEN** `website/docs/r/fwm_rule_group.html.markdown` is generated from `resource_tc_fwm_rule_group.md` without manual edits to `website/`.

#### Scenario: Unit tests mock the cloud API

- **WHEN** the unit test suite for `tencentcloud_fwm_rule_group` runs
- **THEN** it patches `UseFwmV20250611Client` and the four `*WithContext` SDK methods with `gomonkey`, builds `schema.ResourceData` via `schema.TestResourceDataRaw`, and asserts request fields and resulting state without contacting the cloud.
