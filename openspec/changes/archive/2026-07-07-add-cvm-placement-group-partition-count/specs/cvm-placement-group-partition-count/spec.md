## ADDED Requirements

### Requirement: Partition count on placement group creation
The `tencentcloud_placement_group` resource SHALL support an optional `partition_count` parameter (TypeInt, range 2-30, ForceNew) that is passed to the `CreateDisasterRecoverGroup` API as `PartitionCount`. The value SHALL be captured from the Create API response and persisted in Terraform state.

#### Scenario: Create placement group with partition_count
- **WHEN** a user specifies `partition_count = 5` in the `tencentcloud_placement_group` resource configuration
- **THEN** the provider SHALL pass `PartitionCount=5` in the `CreateDisasterRecoverGroup` API request
- **AND** the provider SHALL persist `partition_count = 5` in Terraform state after creation

#### Scenario: Create placement group without partition_count
- **WHEN** a user does NOT specify `partition_count` in the `tencentcloud_placement_group` resource configuration
- **THEN** the provider SHALL NOT set `PartitionCount` in the `CreateDisasterRecoverGroup` API request
- **AND** `partition_count` SHALL NOT appear in Terraform state

#### Scenario: Read existing placement group
- **WHEN** the provider reads an existing `tencentcloud_placement_group` resource
- **THEN** `partition_count` SHALL NOT be refreshed from the Describe API (as the API does not return this field)
- **AND** `partition_count` value (if set) SHALL be retained from the previous state (ForceNew behavior)

#### Scenario: Update partition_count triggers recreation
- **WHEN** a user changes `partition_count` from `5` to `10` in the configuration
- **THEN** the provider SHALL destroy the existing placement group and create a new one with `partition_count = 10`

#### Scenario: Validation of partition_count range
- **WHEN** a user specifies `partition_count = 1` or `partition_count = 31`
- **THEN** the provider SHALL return a validation error indicating the value must be between 2 and 30