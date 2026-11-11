## ADDED Requirements

### Requirement: Table group ID on table group creation
The `tencentcloud_tcaplus_tablegroup` resource SHALL support an optional `table_group_id` parameter (TypeString, immutable after creation) that is passed to the `CreateTableGroup` API as `TableGroupId`. When not specified, the API defaults to auto-increment mode. Changes to `table_group_id` after creation SHALL be rejected with an error via the immutable args check in the Update function.

#### Scenario: Create table group with a user-specified table group id
- **WHEN** a user specifies `table_group_id = "101"` in the `tencentcloud_tcaplus_tablegroup` resource configuration
- **THEN** the provider SHALL pass `TableGroupId=101` in the `CreateTableGroup` API request

#### Scenario: Create table group without a table group id
- **WHEN** a user does NOT specify `table_group_id` in the `tencentcloud_tcaplus_tablegroup` resource configuration
- **THEN** the provider SHALL NOT set `TableGroupId` in the `CreateTableGroup` API request (API defaults to auto-increment)

#### Scenario: Read existing table group
- **WHEN** the provider reads an existing `tencentcloud_tcaplus_tablegroup` resource
- **THEN** `table_group_id` SHALL be refreshed from the `DescribeTableGroups` API response (`TableGroupInfo.TableGroupId`)

#### Scenario: Update table group id triggers immutable error
- **WHEN** a user changes `table_group_id` after creation
- **THEN** the provider SHALL return an error indicating `table_group_id` cannot be changed

### Requirement: CreateGroup service function passes optional table group id
The `CreateGroup` service function SHALL accept an optional `tableGroupId string` parameter and pass it to the `CreateTableGroup` API request as `TableGroupId` when the value is non-empty.

#### Scenario: CreateGroup passes table group id when provided
- **WHEN** `CreateGroup` is called with a non-empty `tableGroupId`
- **THEN** the service function SHALL set `request.TableGroupId` to the provided value before calling the `CreateTableGroup` API

#### Scenario: CreateGroup omits table group id when empty
- **WHEN** `CreateGroup` is called with an empty `tableGroupId`
- **THEN** the service function SHALL NOT set `request.TableGroupId` before calling the `CreateTableGroup` API (API auto-increments)

#### Scenario: CreateGroup returns the created table group id
- **WHEN** the `CreateTableGroup` API succeeds
- **THEN** the service function SHALL return the `TableGroupId` from the API response (`response.Response.TableGroupId`)

### Requirement: Immutable args check in Update function
The `tencentcloud_tcaplus_tablegroup` Update function SHALL use an `immutableArgs` array to validate that the `table_group_id` field is not changed, because the `ModifyTableGroupName` API does not support modifying the table group id.

#### Scenario: table_group_id in immutable args
- **WHEN** defining the immutable args for the Update function
- **THEN** the array SHALL contain `table_group_id`

#### Scenario: Immutable arg change returns error
- **WHEN** `table_group_id` changes in the Update function
- **THEN** the provider SHALL return a formatted error indicating the argument cannot be changed
