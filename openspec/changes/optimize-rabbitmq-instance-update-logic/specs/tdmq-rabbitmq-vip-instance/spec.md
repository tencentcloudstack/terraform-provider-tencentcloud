# Delta Spec: TDMQ RabbitMQ VIP Instance Update Behavior

This delta specification updates the update behavior requirements for the `tencentcloud_tdmq_rabbitmq_vip_instance` Terraform resource.

## ADDED Requirements

### Requirement: Updateable Parameters List
The resource SHALL support updating a specific set of parameters through the `ModifyRabbitMQVipInstance` API.

#### Scenario: List of updateable parameters
- **GIVEN** the resource definition
- **WHEN** examining the update function
- **THEN** the following parameters SHALL be updateable:
  - `cluster_name` - Instance cluster name
  - `remark` - Instance remark information
  - `enable_deletion_protection` - Deletion protection flag
  - `enable_risk_warning` - Cluster risk warning flag
  - `resource_tags` - Instance resource tags

#### Scenario: Update operations only use updateable parameters
- **GIVEN** a user modifies a configuration
- **WHEN** only updateable parameters are changed
- **THEN** the `terraform apply` executes the Update function
- **AND** only the modified updateable parameters are sent to the API

### Requirement: Immutable Parameters List
The resource SHALL reject attempts to update parameters that are not supported by the Tencent Cloud API.

#### Scenario: List of immutable parameters
- **GIVEN** the resource definition
- **WHEN** examining the update function
- **THEN** the following parameters SHALL be immutable (cannot be changed after creation):
  - `zone_ids` - Availability zones
  - `vpc_id` - VPC ID
  - `subnet_id` - Subnet ID
  - `node_spec` - Node specification
  - `node_num` - Number of nodes
  - `storage_size` - Storage size
  - `enable_create_default_ha_mirror_queue` - Mirror queue flag
  - `auto_renew_flag` - Auto renewal flag
  - `time_span` - Purchase duration
  - `pay_mode` - Payment mode
  - `cluster_version` - Cluster version
  - `band_width` - Public network bandwidth
  - `enable_public_access` - Public access flag

#### Scenario: Attempting to update immutable parameters fails
- **GIVEN** an existing RabbitMQ VIP instance
- **WHEN** a user attempts to modify an immutable parameter (e.g., `node_spec`)
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` fails with error message "argument `<parameter_name>` cannot be changed"
- **AND** the instance is not modified

#### Scenario: Clear error messages for immutable parameters
- **GIVEN** a user attempts to update an immutable parameter
- **WHEN** the Update function validates the parameters
- **THEN** the error message explicitly states which parameter cannot be changed
- **AND** the error message is user-friendly and actionable

### Requirement: Update Function Behavior
The resource Update function SHALL implement a clear and maintainable structure for handling parameter updates.

#### Scenario: Parameter change detection
- **GIVEN** the Update function is called
- **WHEN** the function checks for parameter changes
- **THEN** it uses `d.HasChange()` to detect which parameters have changed
- **AND** it only includes changed parameters in the API request

#### Scenario: Update request construction
- **GIVEN** one or more parameters have changed
- **WHEN** constructing the `ModifyRabbitMQVipInstance` request
- **THEN** the request includes only the changed parameters
- **AND** unchanged parameters are not included in the request
- **AND** the `InstanceId` is always included in the request

#### Scenario: Update operation executes when parameters change
- **GIVEN** at least one updateable parameter has changed
- **WHEN** the Update function executes
- **THEN** it calls the `ModifyRabbitMQVipInstance` API with the modified parameters
- **AND** the API call is wrapped in retry logic for transient failures
- **AND** on success, it calls the Read function to refresh the state

#### Scenario: No update when no parameters change
- **GIVEN** the Update function is called
- **WHEN** no updateable parameters have changed
- **THEN** the function returns without calling the API
- **AND** the Read function is still called to refresh state

### Requirement: Error Handling in Update Function
The resource Update function SHALL handle errors appropriately and provide meaningful feedback.

#### Scenario: API failure during update
- **GIVEN** an update operation is in progress
- **WHEN** the `ModifyRabbitMQVipInstance` API returns an error
- **THEN** the error is logged with appropriate context
- **AND** the error is returned to the user
- **AND** the Terraform state remains unchanged
- **AND** subsequent `terraform apply` attempts can retry the operation

#### Scenario: Transient error retry
- **GIVEN** an update operation encounters a transient error (e.g., network issue)
- **WHEN** the retry logic is triggered
- **THEN** the operation is retried according to the configured timeout
- **AND** if retries are exhausted, the final error is returned to the user

#### Scenario: Validation error handling
- **GIVEN** a user provides invalid input for an updateable parameter
- **WHEN** the API validates the input and returns an error
- **THEN** the validation error is propagated to the user
- **AND** the error message includes details about which parameter failed validation
- **AND** the update operation is aborted

### Requirement: State Consistency After Update
The resource SHALL ensure state consistency after successful update operations.

#### Scenario: State reflects updated values
- **GIVEN** an update operation succeeds
- **WHEN** the Read function is called to refresh state
- **THEN** all updateable parameters in state match the actual values in the cloud
- **AND** immutable parameters remain unchanged in state
- **AND** no drift exists between Terraform state and cloud resource

#### Scenario: Computed fields are updated
- **GIVEN** an update operation modifies a parameter
- **WHEN** the Read function refreshes state
- **THEN** any computed fields that may have changed are updated
- **AND** the state accurately reflects the current resource state

### Requirement: Concurrent Update Handling
The resource SHALL handle concurrent update attempts gracefully.

#### Scenario: Optimistic locking behavior
- **GIVEN** two Terraform runs attempt to update the same instance simultaneously
- **WHEN** the second update is attempted after the first completes
- **THEN** the state from the first update is used as the base
- **AND** the second update reconciles any conflicts
- **AND** the final state reflects the intended changes from both operations

### Requirement: Partial Update Support
The resource SHALL support updating a subset of updateable parameters without requiring all parameters to be specified.

#### Scenario: Update only remark field
- **GIVEN** an existing instance with multiple updateable parameters configured
- **WHEN** a user only changes the `remark` field
- **THEN** only the `remark` parameter is sent to the API
- **AND** other updateable parameters remain unchanged in the cloud

#### Scenario: Update multiple fields independently
- **GIVEN** an existing instance
- **WHEN** a user changes multiple updateable parameters in one apply
- **THEN** all changed parameters are sent in a single API call
- **AND** unchanged parameters are not included in the API request
- **AND** the update is atomic (all changes succeed or all fail together)
