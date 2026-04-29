## ADDED Requirements

### Requirement: Resource CRUD lifecycle management
The system SHALL provide a Terraform resource `tencentcloud_teo_load_balancer` that supports full CRUD lifecycle: Create, Read, Update, and Delete operations for TEO Load Balancer instances.

#### Scenario: Create a new load balancer instance
- **WHEN** user applies a Terraform configuration with `tencentcloud_teo_load_balancer` resource specifying `zone_id`, `name`, `type`, `origin_groups`, and optionally `health_checker`, `steering_policy`, `failover_policy`
- **THEN** the system SHALL call `CreateLoadBalancer` API with the specified parameters, set the resource ID to `zone_id:instance_id` (using FILED_SP separator), and populate all computed attributes from the API response

#### Scenario: Read an existing load balancer instance
- **WHEN** Terraform performs a refresh/read operation on an existing `tencentcloud_teo_load_balancer` resource
- **THEN** the system SHALL call `DescribeLoadBalancerList` API with ZoneId and InstanceId filter (Limit=1) to retrieve the current state, and update all schema fields from the response. If the instance is not found, the resource SHALL be marked as removed from state.

#### Scenario: Update an existing load balancer instance
- **WHEN** user modifies updatable fields (`name`, `origin_groups`, `health_checker`, `steering_policy`, `failover_policy`) in the Terraform configuration
- **THEN** the system SHALL call `ModifyLoadBalancer` API with ZoneId, InstanceId, and the changed parameters to update the instance

#### Scenario: Delete a load balancer instance
- **WHEN** user destroys a `tencentcloud_teo_load_balancer` resource
- **THEN** the system SHALL call `DeleteLoadBalancer` API with ZoneId and InstanceId to delete the instance

#### Scenario: Import an existing load balancer instance
- **WHEN** user runs `terraform import tencentcloud_teo_load_balancer.example zone_id:instance_id`
- **THEN** the system SHALL parse the compound ID, read the current state via `DescribeLoadBalancerList`, and populate the Terraform state

### Requirement: Resource ID uses compound key
The system SHALL use `ZoneId` + `InstanceId` as a compound resource ID, separated by `tccommon.FILED_SP`.

#### Scenario: Resource ID format after creation
- **WHEN** a load balancer instance is successfully created with ZoneId="zone-abc" and InstanceId="lb-123"
- **THEN** the Terraform resource ID SHALL be "zone-abc:lb-123" (using FILED_SP separator)

#### Scenario: Resource ID parsing in Read/Update/Delete
- **WHEN** the system needs to perform Read, Update, or Delete operations
- **THEN** it SHALL split the resource ID by FILED_SP to extract ZoneId and InstanceId, and use them as API request parameters

### Requirement: Schema field definitions
The system SHALL define the following schema fields for `tencentcloud_teo_load_balancer`:

#### Scenario: Required fields
- **WHEN** user creates a load balancer resource
- **THEN** the following fields SHALL be required: `zone_id` (ForceNew), `name`, `type` (ForceNew), `origin_groups`

#### Scenario: Optional fields
- **WHEN** user creates a load balancer resource
- **THEN** the following fields SHALL be optional: `health_checker`, `steering_policy`, `failover_policy`

#### Scenario: Computed fields
- **WHEN** Terraform reads the resource state
- **THEN** the following fields SHALL be computed (read-only): `instance_id`, `status`, `origin_group_health_status`, `l4_used_list`, `l7_used_list`, `references`

### Requirement: OriginGroups nested block
The system SHALL support `origin_groups` as a TypeList of nested blocks, each containing `priority` (string, required) and `origin_group_id` (string, required).

#### Scenario: Set origin groups with priorities
- **WHEN** user specifies `origin_groups` with multiple entries each having `priority` and `origin_group_id`
- **THEN** the system SHALL map each entry to an `OriginGroupInLoadBalancer` API struct with Priority and OriginGroupId fields

### Requirement: HealthChecker nested block
The system SHALL support `health_checker` as a TypeList (MaxItems: 1) nested block with the following fields: `type` (required), `port` (optional), `interval` (optional), `timeout` (optional), `health_threshold` (optional), `critical_threshold` (optional), `path` (optional), `method` (optional), `expected_codes` (optional, list of strings), `headers` (optional, list of nested blocks with `key` and `value`), `follow_redirect` (optional), `send_context` (optional), `recv_context` (optional).

#### Scenario: Configure HTTP health checker
- **WHEN** user specifies `health_checker` with type="HTTP", port=80, interval=30, path="/health"
- **THEN** the system SHALL create a `HealthChecker` API struct with the corresponding fields

#### Scenario: No health checker specified
- **WHEN** user does not specify `health_checker`
- **THEN** the system SHALL not send HealthChecker in the Create/Modify request, and the cloud API will default to NoCheck (no health check)

### Requirement: Type field is immutable
The system SHALL set `type` field as ForceNew, since ModifyLoadBalancer API does not support changing the instance type.

#### Scenario: Changing type triggers recreation
- **WHEN** user changes the `type` field from "HTTP" to "GENERAL"
- **THEN** Terraform SHALL destroy the existing resource and create a new one

### Requirement: ZoneId field is immutable
The system SHALL set `zone_id` field as ForceNew, since a load balancer instance cannot be moved between zones.

#### Scenario: Changing zone_id triggers recreation
- **WHEN** user changes the `zone_id` field
- **THEN** Terraform SHALL destroy the existing resource and create a new one

### Requirement: API call retry handling
The system SHALL use `tccommon.ReadRetryTimeout` for read operations and `tccommon.WriteRetryTimeout` for write operations, wrapping errors with `tccommon.RetryError()`.

#### Scenario: Create API call with retry
- **WHEN** the `CreateLoadBalancer` API call fails with a retryable error
- **THEN** the system SHALL retry up to `WriteRetryTimeout` and wrap the error with `tccommon.RetryError()`

#### Scenario: Read API call with retry
- **WHEN** the `DescribeLoadBalancerList` API call fails with a retryable error
- **THEN** the system SHALL retry up to `ReadRetryTimeout` and wrap the error with `tccommon.RetryError()`

### Requirement: Create response validation
The system SHALL validate the CreateLoadBalancer API response, checking that the response is not nil and InstanceId is not empty. If validation fails, return `tccommon.NonRetryableError`.

#### Scenario: Create returns empty InstanceId
- **WHEN** `CreateLoadBalancer` API response has empty InstanceId
- **THEN** the system SHALL return a NonRetryableError

### Requirement: Resource registration in provider
The system SHALL register `tencentcloud_teo_load_balancer` resource in `provider.go` and document it in `provider.md`.

#### Scenario: Resource is available in provider
- **WHEN** user writes a Terraform configuration using `tencentcloud_teo_load_balancer` resource type
- **THEN** Terraform SHALL recognize the resource type and validate its schema

### Requirement: Resource documentation
The system SHALL generate a `.md` documentation file for the resource following the gendoc/README.md format, including one-line description, example usage, and import section.

#### Scenario: Documentation file exists
- **WHEN** the resource is implemented
- **THEN** a documentation file `resource_tc_teo_load_balancer.md` SHALL exist in the teo service directory with proper format
