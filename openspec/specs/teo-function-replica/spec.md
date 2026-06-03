## ADDED Requirements

### Requirement: Create TEO function replica
The system SHALL allow users to create an edge function replica by specifying zone_id, function_id, replica_name, content, and optionally remark. Upon successful creation, the resource ID SHALL be set to the composite of zone_id, function_id, and replica_name joined by `#` separator.

#### Scenario: Successful creation with all parameters
- **WHEN** user provides zone_id, function_id, replica_name, content, and remark in the Terraform configuration
- **THEN** the system calls CreateFunctionReplica API with all parameters and sets the resource ID to `zone_id#function_id#replica_name`

#### Scenario: Successful creation without optional remark
- **WHEN** user provides zone_id, function_id, replica_name, and content without remark
- **THEN** the system calls CreateFunctionReplica API without remark and sets the resource ID to `zone_id#function_id#replica_name`

#### Scenario: API returns error during creation
- **WHEN** the CreateFunctionReplica API returns an error
- **THEN** the system SHALL retry the request using tccommon.RetryError and return the error if retries are exhausted

### Requirement: Read TEO function replica
The system SHALL read the current state of an edge function replica by calling DescribeFunctionReplicas API with zone_id, function_id, and a filter on replica-name, then set the content and remark fields from the matching replica in the response.

#### Scenario: Successful read of existing replica
- **WHEN** the resource exists and DescribeFunctionReplicas returns a matching replica
- **THEN** the system SHALL set content and remark from the FunctionReplica response fields

#### Scenario: Replica not found during read
- **WHEN** DescribeFunctionReplicas returns no matching replica for the given replica_name
- **THEN** the system SHALL remove the resource from state (d.SetId(""))

#### Scenario: Read uses filter to locate specific replica
- **WHEN** reading the resource state
- **THEN** the system SHALL use Filters with Name=replica-name and Values=[replica_name] to query the specific replica, and set Limit to 200 (maximum allowed value)

### Requirement: Update TEO function replica
The system SHALL allow users to update the content and/or remark of an existing edge function replica by calling ModifyFunctionReplica API.

#### Scenario: Update content only
- **WHEN** user changes the content field in Terraform configuration
- **THEN** the system calls ModifyFunctionReplica API with the new content value

#### Scenario: Update remark only
- **WHEN** user changes the remark field in Terraform configuration
- **THEN** the system calls ModifyFunctionReplica API with the new remark value

#### Scenario: Update both content and remark
- **WHEN** user changes both content and remark fields
- **THEN** the system calls ModifyFunctionReplica API with both new values in a single request

### Requirement: Delete TEO function replica
The system SHALL delete an edge function replica by calling DeleteFunctionReplica API with zone_id, function_id, and the replica_name wrapped in a single-element list for ReplicaNames.

#### Scenario: Successful deletion
- **WHEN** the resource is destroyed
- **THEN** the system calls DeleteFunctionReplica API with ReplicaNames containing the single replica_name

#### Scenario: API returns error during deletion
- **WHEN** the DeleteFunctionReplica API returns an error
- **THEN** the system SHALL retry the request using tccommon.RetryError and return the error if retries are exhausted

### Requirement: Import TEO function replica
The system SHALL support importing an existing edge function replica using the composite ID format `zone_id#function_id#replica_name`.

#### Scenario: Successful import with composite ID
- **WHEN** user runs terraform import with ID in format `zone_id#function_id#replica_name`
- **THEN** the system SHALL parse the composite ID, set zone_id, function_id, and replica_name, then call Read to populate remaining fields

#### Scenario: Invalid import ID format
- **WHEN** user provides an import ID that does not contain exactly 3 parts separated by `#`
- **THEN** the system SHALL return an error indicating the expected format

### Requirement: ForceNew on identity fields
The system SHALL force resource recreation when zone_id, function_id, or replica_name are changed, as the API does not support modifying these identity fields.

#### Scenario: Change replica_name triggers recreation
- **WHEN** user changes the replica_name in Terraform configuration
- **THEN** Terraform SHALL plan to destroy the old resource and create a new one

#### Scenario: Change zone_id triggers recreation
- **WHEN** user changes the zone_id in Terraform configuration
- **THEN** Terraform SHALL plan to destroy the old resource and create a new one

### Requirement: Provider registration
The system SHALL register the `tencentcloud_teo_function_replica` resource in provider.go and document it in provider.md.

#### Scenario: Resource is available after provider registration
- **WHEN** the provider is initialized
- **THEN** the resource `tencentcloud_teo_function_replica` SHALL be available for use in Terraform configurations
