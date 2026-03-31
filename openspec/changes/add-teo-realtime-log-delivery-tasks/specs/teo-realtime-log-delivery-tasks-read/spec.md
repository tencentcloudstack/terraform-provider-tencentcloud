## ADDED Requirements

### Requirement: RealtimeLogDeliveryTasks parameter SHALL be available in the resource
The tencentcloud_teo_realtime_log_delivery resource SHALL provide a `realtime_log_delivery_tasks` parameter that is computed and optional. This parameter SHALL contain a list of realtime log delivery tasks with their detailed information retrieved from DescribeRealtimeLogDeliveryTasks API.

#### Scenario: Read resource returns realtime_log_delivery_tasks parameter
- **WHEN** a user performs a read operation on a tencentcloud_teo_realtime_log_delivery resource
- **THEN** the resource SHALL return a `realtime_log_delivery_tasks` parameter containing the list of realtime log delivery tasks
- **AND** each task in the list SHALL contain the detailed information including TaskId, ZoneId, TaskName, TaskType, Status, and other relevant fields returned by DescribeRealtimeLogDeliveryTasks API

#### Scenario: realtime_log_delivery_tasks parameter is optional
- **WHEN** a user creates or updates a tencentcloud_teo_realtime_log_delivery resource without specifying the `realtime_log_delivery_tasks` parameter
- **THEN** the create or update operation SHALL succeed
- **AND** the resource SHALL still compute and return the `realtime_log_delivery_tasks` parameter after the operation completes

### Requirement: RealtimeLogDeliveryTasks parameter SHALL reflect the current state from API
The `realtime_log_delivery_tasks` parameter SHALL always reflect the current state of the realtime log delivery tasks as returned by DescribeRealtimeLogDeliveryTasks API during a read operation.

#### Scenario: realtime_log_delivery_tasks reflects API response
- **WHEN** DescribeRealtimeLogDeliveryTasks API returns a list of realtime log delivery tasks
- **THEN** the `realtime_log_delivery_tasks` parameter SHALL be populated with exactly the same data structure and values as returned by the API
- **AND** all task details from the API response SHALL be included in the parameter

#### Scenario: realtime_log_delivery_tasks handles empty response
- **WHEN** DescribeRealtimeLogDeliveryTasks API returns no tasks or an empty list
- **THEN** the `realtime_log_delivery_tasks` parameter SHALL be an empty list
- **AND** the read operation SHALL not fail

### Requirement: RealtimeLogDeliveryTasks parameter SHALL be backward compatible
The addition of the `realtime_log_delivery_tasks` parameter SHALL not break any existing Terraform configurations or states.

#### Scenario: Existing configurations remain valid
- **WHEN** a user has an existing Terraform configuration for tencentcloud_teo_realtime_log_delivery resource that does not reference `realtime_log_delivery_tasks`
- **THEN** the configuration SHALL continue to work without any modifications
- **AND** terraform apply SHALL not propose any changes due to the new parameter

#### Scenario: Existing states remain valid
- **WHEN** a user has an existing Terraform state for tencentcloud_teo_realtime_log_delivery resource created before the addition of `realtime_log_delivery_tasks` parameter
- **THEN** the state SHALL remain valid
- **AND** terraform refresh SHALL successfully populate the new `realtime_log_delivery_tasks` parameter without requiring state migration

### Requirement: RealtimeLogDeliveryTasks parameter SHALL use correct data types
The `realtime_log_delivery_tasks` parameter SHALL use appropriate Terraform schema types to represent the data structure returned by DescribeRealtimeLogDeliveryTasks API.

#### Scenario: realtime_log_delivery_tasks uses TypeList
- **WHEN** the resource schema defines the `realtime_log_delivery_tasks` parameter
- **THEN** the type SHALL be `schema.TypeList`
- **AND** it SHALL be marked as both `Computed` and `Optional`

#### Scenario: realtime_log_delivery_tasks elements use nested schema
- **WHEN** the `realtime_log_delivery_tasks` parameter contains task elements
- **THEN** each element SHALL be a nested resource with appropriate schema fields
- **AND** each field SHALL use the correct Terraform type (String, Int, Bool, List, or Map) to match the API response type
