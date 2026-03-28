# Spec: RabbitMQ Instance Dynamic Update

This specification defines requirements for RabbitMQ VIP instance core configuration parameters dynamic update capability.

## ADDED Requirements

### Requirement: Node Specification Update Support
The `tencentcloud_tdmq_rabbitmq_vip_instance` resource SHALL support updating node specification parameter through appropriate API calls or instance recreation strategies.

#### Scenario: User updates node specification to higher spec
- **GIVEN** an existing RabbitMQ VIP instance with `node_spec = "rabbit-vip-basic-1"` (4C8G)
- **WHEN** the user changes `node_spec` to `rabbit-vip-profession-4c16g"` (4C16G) in the configuration
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` successfully updates the node specification
- **AND** the instance is modified through appropriate API calls or recreation strategy
- **AND** the Terraform state reflects the updated node specification

#### Scenario: User updates node specification to lower spec
- **GIVEN** an existing RabbitMQ VIP instance with `node_spec = "rabbit-vip-profession-4c16g"` (4C16G)
- **WHEN** the user changes `node_spec` to `rabbit-vip-basic-1"` (4C8G) in the configuration
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` successfully updates the node specification
- **AND** the instance is modified through appropriate API calls or recreation strategy
- **AND** the Terraform state reflects the updated node specification

### Requirement: Node Count Update Support
The `tencentcloud_tdmq_rabbitmq_vip_instance` resource SHALL support updating node count parameter through appropriate API calls or instance recreation strategies.

#### Scenario: User increases node count
- **GIVEN** an existing RabbitMQ VIP instance with `node_num = 1`
- **WHEN** the user changes `node_num` to `3` in the configuration
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` successfully updates the node count
- **AND** the instance is modified through appropriate API calls or recreation strategy
- **AND** the Terraform state reflects the updated node count

#### Scenario: User decreases node count
- **GIVEN** an existing RabbitMQ VIP instance with `node_num = 3`
- **WHEN** the user changes `node_num` to `1` in the configuration
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` successfully updates the node count
- **AND** the instance is modified through appropriate API calls or recreation strategy
- **AND** the Terraform state reflects the updated node count

### Requirement: Storage Size Update Support
The `tencentcloud_tdmq_rabbitmq_vip_instance` resource SHALL support updating storage size parameter through appropriate API calls or instance recreation strategies.

#### Scenario: User increases storage size
- **GIVEN** an existing RabbitMQ VIP instance with `storage_size = 200`
- **WHEN** the user changes `storage_size` to `400` in the configuration
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` successfully updates the storage size
- **AND** the instance is modified through appropriate API calls or recreation strategy
- **AND** the Terraform state reflects the updated storage size

#### Scenario: User decreases storage size (if supported by API)
- **GIVEN** an existing RabbitMQ VIP instance with `storage_size = 400`
- **WHEN** the user changes `storage_size` to `200` in the configuration
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` successfully updates the storage size (if API supports reduction)
- **OR** returns clear error if storage reduction is not supported by the API
- **AND** the Terraform state reflects the updated storage size or remains unchanged

### Requirement: Bandwidth Update Support
The `tencentcloud_tdmq_rabbitmq_vip_instance` resource SHALL support updating bandwidth parameter when public access is enabled.

#### Scenario: User increases bandwidth
- **GIVEN** an existing RabbitMQ VIP instance with `enable_public_access = true` and `band_width = 100`
- **WHEN** the user changes `band_width` to `200` in the configuration
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` successfully updates the bandwidth
- **AND** the instance is modified through appropriate API calls
- **AND** the Terraform state reflects the updated bandwidth value

#### Scenario: User decreases bandwidth
- **GIVEN** an existing RabbitMQ VIP instance with `enable_public_access = true` and `band_width = 200`
- **WHEN** the user changes `band_width` to `100` in the configuration
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` successfully updates the bandwidth
- **AND** the instance is modified through appropriate API calls
- **AND** the Terraform state reflects the updated bandwidth value

### Requirement: Public Access Toggle Support
The `tencentcloud_tdmq_rabbitmq_vip_instance` resource SHALL support toggling public access on/off for existing instances.

#### Scenario: User enables public access on existing instance
- **GIVEN** an existing RabbitMQ VIP instance with `enable_public_access = false`
- **WHEN** the user changes `enable_public_access` to `true` in the configuration
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` successfully enables public access on the instance
- **AND** the instance is modified through appropriate API calls
- **AND** the Terraform state reflects the updated configuration
- **AND** `public_access_endpoint` computed field shows the public access endpoint

#### Scenario: User disables public access on existing instance
- **GIVEN** an existing RabbitMQ VIP instance with `enable_public_access = true`
- **WHEN** the user changes `enable_public_access` to `false` in the configuration
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` successfully disables public access on the instance
- **AND** the instance is modified through appropriate API calls
- **AND** the Terraform state reflects the updated configuration
- **AND** `public_access_endpoint` computed field is cleared or shows empty value

### Requirement: Multiple Parameters Update Support
The `tencentcloud_tdmq_rabbitmq_vip_instance` resource SHALL support updating multiple core configuration parameters simultaneously in a single apply operation.

#### Scenario: User updates multiple parameters at once
- **GIVEN** an existing RabbitMQ VIP instance
- **WHEN** the user changes multiple core parameters (`node_spec`, `node_num`, `storage_size`) simultaneously in the configuration
- **THEN** `terraform plan` detects all changes
- **AND** `terraform apply` successfully updates all parameters in appropriate order
- **AND** the instance is modified through appropriate API calls or recreation strategies
- **AND** the Terraform state reflects all updated parameters
- **AND** no partial updates occur - all changes are applied atomically or rolled back together

### Requirement: Update Operation Error Handling
The `tencentcloud_tdmq_rabbitmq_vip_instance` resource SHALL provide comprehensive error handling for update operations to ensure stability and reliability.

#### Scenario: API returns retryable error during update
- **GIVEN** a user is updating an instance configuration
- **WHEN** the API returns a retryable error (e.g., rate limit, temporary service unavailability)
- **THEN** the provider automatically retries the operation using `helper.Retry()`
- **AND** the retry logic respects configured timeout limits
- **AND** the update eventually succeeds or fails with a clear error message
- **AND** the user is informed about retry attempts in logs

#### Scenario: API returns non-retryable error during update
- **GIVEN** a user is updating an instance configuration
- **WHEN** the API returns a non-retryable error (e.g., invalid parameter, permission denied)
- **THEN** the provider immediately returns the error to the user
- **AND** the error message includes context about what went wrong
- **AND** the Terraform state remains unchanged (no partial updates)
- **AND** the user can retry with corrected configuration

#### Scenario: Update operation times out
- **GIVEN** a user is updating an instance configuration that takes a long time
- **WHEN** the update operation exceeds the configured timeout
- **THEN** the provider returns a timeout error to the user
- **AND** the error message includes the timeout duration and suggests increasing timeout if needed
- **AND** the Terraform state remains unchanged (no partial updates)
- **AND** the user can retry with increased timeout configuration

### Requirement: Update Operation Logging
The `tencentcloud_tdmq_rabbitmq_vip_instance` resource SHALL provide detailed logging for update operations to facilitate troubleshooting and audit.

#### Scenario: Update operation logs start
- **GIVEN** a user initiates an update operation
- **WHEN** the update function is called
- **THEN** the provider logs the start of the update operation with timestamp
- **AND** the log includes the instance ID
- **AND** the log includes which parameters are being changed
- **AND** the log uses `tccommon.LogElapsed()` to track duration

#### Scenario: Update operation logs API calls
- **GIVEN** the provider is making API calls to update an instance
- **WHEN** each API call is executed
- **THEN** the provider logs the API request details (excluding sensitive information)
- **AND** the provider logs the API response details
- **AND** the logs include request and response JSON for debugging
- **AND** the logs include API action name

#### Scenario: Update operation logs successful completion
- **GIVEN** an update operation has been initiated
- **WHEN** the update completes successfully
- **THEN** the provider logs the successful completion with timestamp
- **AND** the log includes the duration of the update operation
- **AND** the log includes the final state of the updated parameters
- **AND** the log indicates that state refresh is being performed

#### Scenario: Update operation logs failure
- **GIVEN** an update operation has been initiated
- **WHEN** the update fails with an error
- **THEN** the provider logs the error details
- **AND** the log includes the error message and stack trace
- **AND** the log includes the state of parameters before the failure
- **AND** the provider returns the error to the user for resolution

### Requirement: Update State Consistency
The `tencentcloud_tdmq_rabbitmq_vip_instance` resource SHALL ensure Terraform state remains consistent with the actual cloud resource state during and after update operations.

#### Scenario: State consistency after successful update
- **GIVEN** a user successfully updates an instance configuration
- **WHEN** the update operation completes
- **THEN** the provider immediately calls the Read function to refresh state
- **AND** the Terraform state reflects the actual state of the cloud resource
- **AND** `tccommon.InconsistentCheck()` validates state consistency
- **AND** no drift exists between state and cloud resource

#### Scenario: State consistency after failed update
- **GIVEN** a user attempts to update an instance but the update fails
- **WHEN** the update operation fails
- **THEN** the Terraform state remains unchanged (no partial updates)
- **AND** the state still reflects the pre-update configuration
- **AND** the user can retry the update or revert to previous configuration
- **AND** `tccommon.InconsistentCheck()` ensures no state corruption
- **AND** the provider returns the error without corrupting state

### Requirement: Backward Compatibility
All dynamic update capabilities SHALL be backward compatible with existing resources and configurations.

#### Scenario: Existing resource before dynamic update enhancements
- **GIVEN** a RabbitMQ VIP instance managed by Terraform before this feature
- **WHEN** the provider is upgraded to include dynamic update capabilities
- **THEN** `terraform plan` shows no unexpected changes
- **AND** existing resources continue to function normally
- **AND** the user can still use all existing update operations
- **AND** state refresh works correctly without forcing updates

#### Scenario: Gradual adoption of new update capabilities
- **GIVEN** a user has an existing RabbitMQ VIP instance
- **WHEN** the user decides to use newly supported update operations (e.g., updating bandwidth)
- **THEN** the user can simply modify the configuration and run `terraform apply`
- **AND** no manual state modification or resource recreation is required
- **AND** the update operation succeeds if the API supports it
- **AND** the user is not forced to upgrade all instances at once
