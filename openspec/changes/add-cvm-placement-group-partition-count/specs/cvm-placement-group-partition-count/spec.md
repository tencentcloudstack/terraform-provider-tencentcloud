## ADDED Requirements

### Requirement: Strategy on placement group creation
The `tencentcloud_placement_group` resource SHALL support an optional `strategy` parameter (TypeString, values SPREAD/PARTITION, Computed, immutable after creation) that is passed to the `CreateDisasterRecoverGroup` API as `Strategy`. Default strategy is `SPREAD`. Changes to `strategy` after creation SHALL be rejected with an error via the immutable args check in the Update function.

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
The `tencentcloud_placement_group` resource SHALL support an optional `partition_count` parameter (TypeInt, range 2-30, immutable after creation) that is passed to the `CreateDisasterRecoverGroup` API as `PartitionCount`. This parameter is only valid when `strategy` is `PARTITION`. Changes to `partition_count` after creation SHALL be rejected with an error via the immutable args check in the Update function.

#### Scenario: Create placement group with strategy=PARTITION and partition_count
- **WHEN** a user specifies `strategy = "PARTITION"` and `partition_count = 5` in the `tencentcloud_placement_group` resource configuration
- **THEN** the provider SHALL pass `Strategy=PARTITION` and `PartitionCount=5` in the `CreateDisasterRecoverGroup` API request

#### Scenario: Create placement group with partition_count but without PARTITION strategy
- **WHEN** a user specifies `partition_count = 5` but does NOT set `strategy = "PARTITION"` in the configuration
- **THEN** the provider SHALL return an error indicating `partition_count` is only valid when `strategy` is set to `PARTITION`

#### Scenario: Create placement group without partition_count
- **WHEN** a user does NOT specify `partition_count` in the `tencentcloud_placement_group` resource configuration
- **THEN** the provider SHALL NOT set `PartitionCount` in the `CreateDisasterRecoverGroup` API request

#### Scenario: Read existing placement group (after SDK update)
- **WHEN** the provider reads an existing `tencentcloud_placement_group` resource
- **THEN** `strategy` and `partition_count` SHALL be refreshed from the `DescribeDisasterRecoverGroups` API response (`DisasterRecoverGroup.Strategy` and `DisasterRecoverGroup.PartitionCount`)

#### Scenario: Update strategy triggers immutable error
- **WHEN** a user changes `strategy` after creation
- **THEN** the provider SHALL return an error: `argument 'strategy' cannot be changed`

#### Scenario: Update partition_count triggers immutable error
- **WHEN** a user changes `partition_count` after creation
- **THEN** the provider SHALL return an error: `argument 'partition_count' cannot be changed`

#### Scenario: Validation of partition_count range
- **WHEN** a user specifies `partition_count = 1` or `partition_count = 31`
- **THEN** the provider SHALL return a validation error indicating the value must be between 2 and 30

### Requirement: Immutable args check in Update function
The `tencentcloud_placement_group` Update function SHALL use an `immutableArgs` array pattern to validate that immutable fields are not changed.

#### Scenario: Immutable args array definition
- **WHEN** defining the immutable args
- **THEN** the array SHALL contain `["type", "strategy", "affinity", "partition_count"]`
- **AND** the array SHALL be defined as a package-level variable `CVM_PLACEMENT_GROUP_IMMUTABLE_ARGS` in `extension_cvm.go` for extensibility

#### Scenario: Any immutable arg change returns error
- **WHEN** any field in `immutableArgs` changes in the Update function
- **THEN** the provider SHALL return a formatted error `argument '%s' cannot be changed`

### Requirement: CreatePlacementGroup returns full response
The `CreatePlacementGroup` service function SHALL return the complete `*cvm.CreateDisasterRecoverGroupResponse` instead of individual fields, for better extensibility.

#### Scenario: Caller extracts fields from response
- **WHEN** `CreatePlacementGroup` succeeds
- **THEN** the caller SHALL extract `DisasterRecoverGroupId` from `response.Response.DisasterRecoverGroupId`
- **AND** no individual field values SHALL be pre-extracted by the service function