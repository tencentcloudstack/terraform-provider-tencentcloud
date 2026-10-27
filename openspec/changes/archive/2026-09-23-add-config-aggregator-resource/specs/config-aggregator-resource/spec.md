## ADDED Requirements

### Requirement: Manage Config Aggregator lifecycle
The provider SHALL manage a Tencent Cloud Config Aggregator (账号组) through a `tencentcloud_config_aggregator` resource supporting create, read, update, and delete operations against the config v20220802 APIs.

#### Scenario: Create a custom aggregator
- **WHEN** a user applies a `tencentcloud_config_aggregator` resource with `name`, `description`, `type=CUSTOM`, `owner_uin`, and an `aggregator_accounts` list
- **THEN** the provider SHALL call `CreateAggregator` and set the resource id to the composite `account_group_id#owner_uin`

#### Scenario: Read an existing aggregator
- **WHEN** the provider reads the resource
- **THEN** the provider SHALL call `DescribeAggregator` with `AccountGroupId` and `OwnerUin` parsed from the composite id and populate `name`, `description`, `type`, `aggregator_accounts`, `aggregator_status`, and `account_group_id`

#### Scenario: Update name/description/members
- **WHEN** `name`, `description`, or `aggregator_accounts` changes
- **THEN** the provider SHALL call `UpdateAggregator` with `AccountGroupId` and `OwnerUin` from the id plus the changed fields

#### Scenario: Delete an aggregator
- **WHEN** the resource is destroyed
- **THEN** the provider SHALL call `DeleteAggregators` with `AccountGroupId` and `OwnerUin` from the id

### Requirement: Aggregator schema fields
The resource SHALL expose the following fields: `name` (required string), `description` (required string), `type` (required string, ForceNew, enum `RD`/`CUSTOM`), `owner_uin` (required string, ForceNew), `aggregator_accounts` (optional list of objects with `member_uin` int and `member_name` string), `account_group_id` (computed string), and `aggregator_status` (computed int).

#### Scenario: Nested aggregator_accounts round-trip
- **WHEN** `aggregator_accounts` is set with one or more `{member_uin, member_name}` entries
- **THEN** on read the provider SHALL flatten the returned `AggregatorAccounts` list back into `aggregator_accounts` preserving each member's `member_uin` and `member_name`

### Requirement: Aggregator resource id and import
The resource id SHALL be the composite of `account_group_id` and `owner_uin` joined by `tccommon.FILED_SP`. The resource SHALL support `terraform import` using that composite id.

#### Scenario: Import by composite id
- **WHEN** a user runs `terraform import tencentcloud_config_aggregator.foo "ca-xxxxxxxx#100012345678"`
- **THEN** the provider SHALL parse the id into `account_group_id=ca-xxxxxxxx` and `owner_uin=100012345678` and read the resource

### Requirement: Missing resource handling on read
The provider SHALL, when `DescribeAggregator` returns an empty response for an existing state id, log the id via `log.Printf("[CRUD] tencentcloud_config_aggregator id=%s", d.Id())` and then `d.SetId("")` so Terraform removes it from state.

#### Scenario: Aggregator already deleted out-of-band
- **WHEN** `DescribeAggregator` returns nil/empty during read
- **THEN** the provider SHALL log the id and clear the state id instead of erroring