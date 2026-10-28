## Requirements

### Requirement: Strategy on placement group creation
The `tencentcloud_placement_group` resource SHALL support an optional `strategy` parameter (TypeString, values SPREAD/PARTITION, ForceNew, Computed) that is passed to the `CreateDisasterRecoverGroup` API as `Strategy`. Default strategy is `SPREAD`.

#### Scenario: Create placement group with PARTITION strategy
- **WHEN** a user specifies `strategy = "PARTITION"` in the `tencentcloud_placement_group` resource configuration
- **THEN** the provider SHALL pass `Strategy=PARTITION` in the `CreateDisasterRecoverGroup` API request

#### Scenario: Create placement group without explicit strategy
- **WHEN** a user does NOT specify `strategy` in the `tencentcloud_placement_group` resource configuration
- **THEN** the provider SHALL NOT set `Strategy` in the `CreateDisasterRecoverGroup` API request (API defaults to SPREAD)

#### Scenario: Validation of strategy values
- **WHEN** a user specifies `strategy = "INVALID"` 
- **THEN** the provider SHALL return a validation error indicating valid values are `SPREAD` and `PARTITION`

### Requirement: Partition count on placement group creation
The `tencentcloud_placement_group` resource SHALL support an optional `partition_count` parameter (TypeInt, range 2-30, ForceNew) that is passed to the `CreateDisasterRecoverGroup` API as `PartitionCount`. The value SHALL be captured from the Create API response and persisted in Terraform state. This parameter is only valid when `strategy` is `PARTITION`.

#### Scenario: Create placement group with strategy=PARTITION and partition_count
- **WHEN** a user specifies `strategy = "PARTITION"` and `partition_count = 5` in the `tencentcloud_placement_group` resource configuration
- **THEN** the provider SHALL pass `Strategy=PARTITION` and `PartitionCount=5` in the `CreateDisasterRecoverGroup` API request
- **AND** the provider SHALL persist `partition_count = 5` in Terraform state after creation

#### Scenario: Create placement group with partition_count but without PARTITION strategy
- **WHEN** a user specifies `partition_count = 5` but does NOT set `strategy = "PARTITION"` in the configuration
- **THEN** the provider SHALL return an error indicating `partition_count` is only valid when `strategy` is set to `PARTITION`

#### Scenario: Create placement group without partition_count
- **WHEN** a user does NOT specify `partition_count` in the `tencentcloud_placement_group` resource configuration
- **THEN** the provider SHALL NOT set `PartitionCount` in the `CreateDisasterRecoverGroup` API request
- **AND** `partition_count` SHALL NOT appear in Terraform state

#### Scenario: Read existing placement group
- **WHEN** the provider reads an existing `tencentcloud_placement_group` resource
- **THEN** `partition_count` and `strategy` SHALL NOT be refreshed from the Describe API (as the API does not return these fields)
- **AND** `partition_count` and `strategy` values (if set) SHALL be retained from the previous state (ForceNew behavior)

#### Scenario: Update partition_count triggers recreation
- **WHEN** a user changes `partition_count` from `5` to `10` in the configuration
- **THEN** the provider SHALL destroy the existing placement group and create a new one with `partition_count = 10`

#### Scenario: Validation of partition_count range
- **WHEN** a user specifies `partition_count = 1` or `partition_count = 31`
- **THEN** the provider SHALL return a validation error indicating the value must be between 2 and 30