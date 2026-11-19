## ADDED Requirements

### Requirement: Manage aggregate delivery configuration
The Terraform provider SHALL provide a resource `tencentcloud_config_update_aggregate_config_deliver` that manages the delivery settings (投递设置) of a Tencent Cloud Config account group (账号组). The resource SHALL support Read and Update operations against the `DescribeAggregateConfigDeliver` and `UpdateAggregateConfigDeliver` cloud APIs.

#### Scenario: Create a new aggregate delivery config
- **WHEN** a user applies a `tencentcloud_config_update_aggregate_config_deliver` resource with a valid `account_group_id` and `status`
- **THEN** the provider SHALL call `UpdateAggregateConfigDeliver` to set the delivery configuration, set the resource id to the `account_group_id`, and refresh state via Read

#### Scenario: Read the aggregate delivery config
- **WHEN** the provider reads an existing `tencentcloud_config_update_aggregate_config_deliver` resource
- **THEN** it SHALL call `DescribeAggregateConfigDeliver` with the `account_group_id` (from `d.Id()`) and populate all schema fields that the API returns as non-nil

#### Scenario: Update the aggregate delivery config
- **WHEN** a user changes a mutable field (e.g. `status`, `deliver_name`, `target_arn`) on an existing resource
- **THEN** the provider SHALL call `UpdateAggregateConfigDeliver` with the updated fields and the `account_group_id`, then refresh state via Read

#### Scenario: Delete the resource
- **WHEN** a user destroys a `tencentcloud_config_update_aggregate_config_deliver` resource
- **THEN** the provider SHALL remove it from state without calling any cloud delete API (no delete API exists), and SHALL NOT alter the cloud delivery configuration

### Requirement: Schema fields match cloud API
The resource schema SHALL expose `account_group_id` and `status` as required, `deliver_name`, `target_arn`, `deliver_prefix`, `deliver_type`, `deliver_uin`, `deliver_content_type` as optional, and `create_time` as computed. Field types SHALL correspond to the cloud API types (`account_group_id` string, `status`/`deliver_content_type` uint64→int, `deliver_uin` int64→int, others string).

#### Scenario: account_group_id is the resource key and ForceNew
- **WHEN** a user changes `account_group_id` on an existing resource
- **THEN** the provider SHALL treat it as ForceNew (destroy + recreate) because `account_group_id` identifies a different account group's delivery config

#### Scenario: create_time is computed
- **WHEN** the resource is read
- **THEN** `create_time` SHALL be populated only from the `DescribeAggregateConfigDeliver` response and MUST NOT be settable by the user

### Requirement: Import support
The resource SHALL support `terraform import` using the `account_group_id` as the import id.

#### Scenario: Import an existing aggregate delivery config
- **WHEN** a user runs `terraform import tencentcloud_config_update_aggregate_config_deliver.example <account_group_id>`
- **THEN** the provider SHALL set the id to the account group id and perform a Read to populate state

### Requirement: Nil-safety in Read
The Read handler SHALL check that each response field is non-nil before calling `d.Set`, and if the response or its `Response` field is nil it SHALL log the resource id and clear the state.

#### Scenario: API returns empty response
- **WHEN** `DescribeAggregateConfigDeliver` returns a nil response or nil `Response`
- **THEN** the Read handler SHALL log `[CRUD]` with the resource id and call `d.SetId("")` so the resource is removed from state
