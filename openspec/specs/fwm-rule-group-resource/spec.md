# fwm-rule-group-resource Specification

## Purpose
TBD - created by archiving change add-fwm-rule-group-resource. Update Purpose after archive.
## Requirements
### Requirement: Resource manages FWM rule group lifecycle

The provider SHALL expose a resource `tencentcloud_fwm_rule_group` that manages a Tencent Cloud Firewall Manager rule group via `CreateRuleGroup`, `DescribeRuleGroups`, `ModifyRuleGroup` and `DeleteRuleGroup` in the `fwm/v20250611` SDK namespace. All four APIs are synchronous; no async task polling is required.

#### Scenario: Create provisions a rule group

- **WHEN** a user defines the resource with `group_name` and `product`
- **THEN** the provider SHALL call `CreateRuleGroup` with `GroupName` and `Product`, require a non-nil non-empty `Response.GroupId`, and set the resource ID to that bare `GroupId` value (`fwmrg-` prefix)

#### Scenario: Create rejects an empty group id

- **WHEN** `CreateRuleGroup` returns `Response == nil`, `Response.GroupId == nil`, or an empty string id
- **THEN** the provider SHALL return `resource.NonRetryableError` and SHALL NOT write a resource ID

#### Scenario: Read reflects current cloud state

- **WHEN** the provider reads the resource
- **THEN** it SHALL call `DescribeRuleGroups` with a `GroupId` filter, and map every non-nil field of the matching `RuleGroupResp` into Terraform state

#### Scenario: Read clears state when the group is gone

- **WHEN** `DescribeRuleGroups` returns no record whose `GroupId` exactly equals `d.Id()`
- **THEN** the provider SHALL log `[CRUD] tencentcloud_fwm_rule_group id=<id>` and then call `d.SetId("")`, returning nil

#### Scenario: Update renames the group

- **WHEN** `group_name` changes
- **THEN** the provider SHALL call `ModifyRuleGroup` with `GroupId = d.Id()` and the new `GroupName`, then re-read the resource

#### Scenario: Update is a no-op without changes

- **WHEN** no mutable attribute has changed
- **THEN** the provider SHALL skip the `ModifyRuleGroup` call and only perform a read

#### Scenario: Delete removes the group

- **WHEN** the resource is destroyed
- **THEN** the provider SHALL call `DeleteRuleGroup` with `GroupIds` containing exactly one element — `d.Id()`

### Requirement: Resource is importable by rule group id

The resource ID SHALL be the bare rule group ID, and the resource SHALL support `terraform import` with that ID.

#### Scenario: Import with an existing id

- **WHEN** a user runs `terraform import tencentcloud_fwm_rule_group.example fwmrg-wsrrhuvrib`
- **THEN** the provider SHALL accept the ID through `schema.ImportStatePassthrough` and populate all fields from `DescribeRuleGroups` on the subsequent read

#### Scenario: Id is used directly in every API call

- **WHEN** Read, Update or Delete is invoked
- **THEN** the provider SHALL use `d.Id()` verbatim as `GroupId` (or as the single element of `GroupIds`) without splitting or re-joining it on any separator

### Requirement: Schema mirrors the mutable API parameters

The resource schema SHALL contain `group_name` and `product` as the only user-configurable attributes, and SHALL expose all describe-only fields as computed attributes, flattened to the top level without a `rule_groups` wrapper layer.

#### Scenario: Required attributes enforced

- **WHEN** a user omits `group_name` or `product`
- **THEN** Terraform SHALL return a validation error before any API call

#### Scenario: Product is validated and immutable

- **WHEN** `product` is set to a value outside `enterprise_sg`, `cfw_nat_acl`, `cfw_vpc_acl`, `cfw_edge_acl`
- **THEN** Terraform SHALL reject the configuration at plan time
- **AND WHEN** a user changes `product` on an existing resource
- **THEN** Terraform SHALL plan a destroy-and-recreate, because `ModifyRuleGroup` cannot change the product

#### Scenario: Computed attributes mirror RuleGroupResp

- **WHEN** the resource is read
- **THEN** `group_id`, `app_id`, `rule_count`, `create_time`, `update_time`, `create_by`, `create_name`, `update_by`, `update_name`, `receive_account`, `exec_areas` and `direction` SHALL be set from the corresponding `RuleGroupResp` fields, each guarded by a nil check

#### Scenario: Nested receive_account mirrors SDK structs

- **WHEN** `RuleGroupResp.ReceiveAccount` is non-empty
- **THEN** each element SHALL be exposed with `uin`, `nickname`, `receiver_type`, and a nested `members` list whose elements carry `app_id`, `uin`, `nickname`, `member_id` and `node_name`

#### Scenario: The rule_groups list is not wrapped

- **WHEN** the schema is defined
- **THEN** it SHALL NOT contain a `rule_groups` (or `rule_groups_set` / `rule_groups_list`) attribute; the fields of the single matched `RuleGroupResp` SHALL live at the top level

### Requirement: Product numeric code is normalized on read

`DescribeRuleGroups` returns `Product` as a numeric code while `CreateRuleGroup` accepts a string name. Read SHALL translate the code back to the configured name so that no spurious diff is produced.

#### Scenario: Numeric code is mapped to its product name

- **WHEN** the describe response returns `Product` equal to `1`, `2`, `3` or `4`
- **THEN** the provider SHALL set the `product` attribute to `enterprise_sg`, `cfw_nat_acl`, `cfw_vpc_acl` or `cfw_edge_acl` respectively

#### Scenario: An already-named product passes through

- **WHEN** the describe response returns `Product` as a product name rather than a numeric code
- **THEN** the provider SHALL store that value unchanged

### Requirement: Single-group lookup uses a GroupId filter with exact client-side match

Because the SDK exposes no get-by-id API, Read SHALL resolve a single rule group through `DescribeRuleGroups`.

#### Scenario: Filter shape

- **WHEN** the service helper queries one rule group
- **THEN** it SHALL send `Filters: [{Name: "GroupId", Values: [<id>], OperatorType: 9}]` using the SDK type `CommonFilter`

#### Scenario: Exact match guards the fuzzy operator

- **WHEN** the fuzzy filter returns multiple `RuleGroups`
- **THEN** the helper SHALL return only the element whose `GroupId` exactly equals the requested id, and SHALL return `(nil, nil)` if none matches

#### Scenario: Pagination uses the maximum page size

- **WHEN** the helper paginates `DescribeRuleGroups`
- **THEN** it SHALL pass `Limit = 100` and increment `Offset` by `100` until the returned page is shorter than the limit

### Requirement: Every SDK invocation is wrapped with project retry semantics

Every SDK call SHALL be invoked inside `resource.Retry(...)` using the project timeout constants.

#### Scenario: Read calls use ReadRetryTimeout

- **WHEN** `DescribeRuleGroups` is invoked from Read or from the service helper
- **THEN** the call SHALL be wrapped in `resource.Retry(tccommon.ReadRetryTimeout, ...)` and use `tccommon.RetryError` for retryable errors

#### Scenario: Mutating calls use WriteRetryTimeout

- **WHEN** `CreateRuleGroup`, `ModifyRuleGroup` or `DeleteRuleGroup` is invoked
- **THEN** the call SHALL be wrapped in `resource.Retry(tccommon.WriteRetryTimeout, ...)` and use `tccommon.RetryError` for retryable errors

#### Scenario: Nil SDK response is non-retryable

- **WHEN** any SDK call returns `result == nil` or `result.Response == nil`
- **THEN** the call SHALL return `resource.NonRetryableError` with a message naming `tencentcloud_fwm_rule_group`

#### Scenario: No nested retries

- **WHEN** a resource CRUD function calls a service helper that already retries
- **THEN** the resource SHALL NOT add a second `resource.Retry` around that helper call

### Requirement: Create waits for read-after-create visibility

`CreateRuleGroup` is synchronous but list queries are eventually consistent, so Create SHALL confirm visibility before delegating to Read.

#### Scenario: Visibility confirmed before read

- **WHEN** `CreateRuleGroup` returns a `GroupId` and `d.SetId` has been called
- **THEN** the provider SHALL poll `DescribeRuleGroups` (via the service helper) inside `resource.Retry(tccommon.ReadRetryTimeout, ...)` returning `resource.RetryableError` while the group is not yet visible, and SHALL only then invoke Read

#### Scenario: SetId happens before any polling

- **WHEN** the create response carries a valid `GroupId`
- **THEN** `d.SetId(<group_id>)` SHALL be executed before the visibility poll, so a failed poll still leaves an id in state that `terraform destroy` can clean up

