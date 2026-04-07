# Spec Delta: TDMQ RabbitMQ VIP Instance Update Logic Optimization

This delta spec modifies the existing `tdmq-rabbitmq-vip-instance` spec to support additional updateable fields and add async state waiting.

## MODIFIED Requirements

### Requirement: Public Access Field Immutability
The `enable_public_access` and `band_width` fields SHALL be updateable after instance creation through the ModifyRabbitMQVipInstance API.

#### Scenario: User enables public access on existing instance
- **GIVEN** an existing RabbitMQ VIP instance with `enable_public_access = false`
- **WHEN** user changes `enable_public_access` to `true` in configuration
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` calls `ModifyRabbitMQVipInstance` API with updated `EnablePublicAccess` flag
- **AND** API successfully updates the instance configuration
- **AND** instance is not recreated
- **AND** Terraform state reflects the updated `enable_public_access` value

#### Scenario: User disables public access on existing instance
- **GIVEN** an existing RabbitMQ VIP instance with `enable_public_access = true`
- **WHEN** user changes `enable_public_access` to `false` in configuration
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` calls `ModifyRabbitMQVipInstance` API to disable public access
- **AND** API successfully updates the instance configuration
- **AND** instance is not recreated
- **AND** `public_access_endpoint` computed field becomes empty or reflects the change

#### Scenario: User modifies bandwidth on existing instance
- **GIVEN** an existing RabbitMQ VIP instance with `band_width = 100`
- **WHEN** user changes `band_width` to `200` in configuration
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` calls `ModifyRabbitMQVipInstance` API with updated `Bandwidth` parameter
- **AND** API successfully updates the bandwidth
- **AND** instance is not recreated
- **AND** Terraform state reflects the updated `band_width` value

#### Scenario: User updates both public access and bandwidth simultaneously
- **GIVEN** an existing RabbitMQ VIP instance with `enable_public_access = false` and `band_width = 100`
- **WHEN** user changes both `enable_public_access = true` and `band_width = 200` in configuration
- **THEN** `terraform plan` detects both changes
- **AND** `terraform apply` calls `ModifyRabbitMQVipInstance` API once with both parameters
- **AND** API successfully updates both fields in a single operation
- **AND** instance is not recreated
- **AND** Terraform state reflects both updated values

### Requirement: Auto Renew Flag Update Support
The resource SHALL support updating the `auto_renew_flag` field after instance creation.

#### Scenario: User enables auto renew on existing instance
- **GIVEN** an existing RabbitMQ VIP instance with `auto_renew_flag = false`
- **WHEN** user changes `auto_renew_flag` to `true` in configuration
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` calls `ModifyRabbitMQVipInstance` API with updated `AutoRenewFlag` parameter
- **AND** API successfully enables auto renew for the instance
- **AND** instance is not recreated
- **AND** Terraform state reflects the updated `auto_renew_flag` value

#### Scenario: User disables auto renew on existing instance
- **GIVEN** an existing RabbitMQ VIP instance with `auto_renew_flag = true`
- **WHEN** user changes `auto_renew_flag` to `false` in configuration
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` calls `ModifyRabbitMQVipInstance` API to disable auto renew
- **AND** API successfully disables auto renew
- **AND** instance is not recreated
- **AND** Terraform state reflects the updated `auto_renew_flag` value

### Requirement: Update API Integration for Additional Fields
The resource SHALL correctly pass additional updateable fields to `ModifyRabbitMQVipInstance` API.

#### Scenario: Send enable_public_access to Update API
- **GIVEN** a user updates an instance with `enable_public_access = true`
- **WHEN** Update operation calls `ModifyRabbitMQVipInstance`
- **THEN** request includes `EnablePublicAccess: helper.Bool(true)`
- **AND** API accepts the parameter and updates the instance

#### Scenario: Send band_width to Update API
- **GIVEN** a user updates an instance with `band_width = 200`
- **WHEN** Update operation calls `ModifyRabbitMQVipInstance`
- **THEN** request includes `Bandwidth: helper.IntUint64(200)`
- **AND** API accepts the parameter and updates the bandwidth

#### Scenario: Send auto_renew_flag to Update API
- **GIVEN** a user updates an instance with `auto_renew_flag = true`
- **WHEN** Update operation calls `ModifyRabbitMQVipInstance`
- **THEN** request includes `AutoRenewFlag: helper.Bool(true)`
- **AND** API accepts the parameter and updates auto renew settings

#### Scenario: Update multiple fields in single API call
- **GIVEN** a user updates an instance with changes to `enable_public_access`, `band_width`, `cluster_name`, and `resource_tags`
- **WHEN** Update operation detects multiple changes via `d.HasChange()`
- **THEN** request to `ModifyRabbitMQVipInstance` includes all changed parameters
- **AND** API call succeeds and updates all fields in a single operation
- **AND** only one API call is made (not multiple calls for each field)

### Requirement: Update State Wait
The resource SHALL wait for instance status to stabilize after update operations.

#### Scenario: Successful update with status waiting
- **GIVEN** an existing RabbitMQ VIP instance
- **WHEN** Update operation calls `ModifyRabbitMQVipInstance` API
- **AND** API returns success
- **AND** resource enters "Updating" status
- **THEN** Update operation polls instance status using `DescribeTdmqRabbitmqVipInstanceByFilter`
- **AND** operation waits until status changes to "Running" or "Success"
- **AND** retry mechanism is used with timeout (`tccommon.ReadRetryTimeout*10`)
- **AND** operation returns successfully after status stabilizes

#### Scenario: Update timeout handling
- **GIVEN** an existing RabbitMQ VIP instance
- **WHEN** Update operation calls `ModifyRabbitMQVipInstance` API
- **AND** API returns success but instance status remains in "Updating" state beyond timeout period
- **THEN** operation returns a timeout error to user
- **AND** error message indicates the instance is still updating
- **AND** user can retry `terraform apply` after manual intervention

#### Scenario: Update with status polling errors
- **GIVEN** an existing RabbitMQ VIP instance
- **WHEN** Update operation polls instance status during wait period
- **AND** intermediate API calls return retryable errors (e.g., rate limiting)
- **THEN** operation uses `tccommon.RetryError` to handle retryable errors
- **AND** polling continues with exponential backoff
- **AND** operation completes successfully once status stabilizes

#### Scenario: Update with invalid status transitions
- **GIVEN** an existing RabbitMQ VIP instance
- **WHEN** Update operation polls instance status and detects an unexpected state (e.g., "Failed", "Rollback")
- **THEN** operation returns a non-retryable error
- **AND** error message clearly indicates the illegal status
- **AND** update operation fails without waiting for timeout

#### Scenario: Update with no status change (idempotent)
- **GIVEN** an existing RabbitMQ VIP instance
- **WHEN** Update operation is called but API detects no actual changes
- **AND** instance status is already "Running" or "Success"
- **THEN** operation may skip waiting or completes immediately
- **AND** operation returns successfully without errors

### Requirement: Reduced Immutable Parameters
The resource SHALL only prohibit updating truly immutable parameters that require instance recreation.

#### Scenario: Zone IDs remain immutable
- **GIVEN** an existing RabbitMQ VIP instance with `zone_ids = [1]`
- **WHEN** user attempts to change `zone_ids` to `[1, 2, 3]` in configuration
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` fails with error message "argument `zone_ids` cannot be changed after instance creation"
- **AND** error message includes guidance: "Please recreate the instance if you need to modify this parameter"
- **AND** instance is not modified or recreated automatically

#### Scenario: VPC ID remains immutable
- **GIVEN** an existing RabbitMQ VIP instance with `vpc_id = "vpc-abc123"`
- **WHEN** user attempts to change `vpc_id` to another VPC ID
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` fails with clear error message
- **AND** instance is not modified

#### Scenario: Subnet ID remains immutable
- **GIVEN** an existing RabbitMQ VIP instance with `subnet_id = "subnet-xyz789"`
- **WHEN** user attempts to change `subnet_id`
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` fails with clear error message
- **AND** instance is not modified

#### Scenario: Node spec remains immutable
- **GIVEN** an existing RabbitMQ VIP instance with `node_spec = "rabbit-vip-basic-1"`
- **WHEN** user attempts to change `node_spec`
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` fails with clear error message
- **AND** instance is not modified

#### Scenario: Node count remains immutable
- **GIVEN** an existing RabbitMQ VIP instance with `node_num = 1`
- **WHEN** user attempts to change `node_num` to `3`
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` fails with clear error message
- **AND** instance is not modified

#### Scenario: Storage size remains immutable
- **GIVEN** an existing RabbitMQ VIP instance with `storage_size = 200`
- **WHEN** user attempts to change `storage_size` to `500`
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` fails with clear error message
- **AND** instance is not modified

#### Scenario: HA mirror queue flag remains immutable
- **GIVEN** an existing RabbitMQ VIP instance with `enable_create_default_ha_mirror_queue = false`
- **WHEN** user attempts to change `enable_create_default_ha_mirror_queue`
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` fails with clear error message
- **AND** instance is not modified

#### Scenario: Time span remains immutable
- **GIVEN** an existing RabbitMQ VIP instance with `time_span = 1`
- **WHEN** user attempts to change `time_span` to `12`
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` fails with clear error message
- **AND** instance is not modified

#### Scenario: Pay mode remains immutable
- **GIVEN** an existing RabbitMQ VIP instance with `pay_mode = 1` (prepaid)
- **WHEN** user attempts to change `pay_mode` to `0` (postpaid)
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` fails with clear error message
- **AND** instance is not modified

#### Scenario: Cluster version remains immutable
- **GIVEN** an existing RabbitMQ VIP instance with `cluster_version = "3.8.30"`
- **WHEN** user attempts to change `cluster_version` to "3.11.8"
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` fails with clear error message
- **AND** error message indicates version upgrades require special procedures
- **AND** instance is not modified

### Requirement: Enhanced Error Messages for Immutable Parameters
The resource SHALL provide clear, actionable error messages when users attempt to modify immutable parameters.

#### Scenario: Error message includes parameter name
- **GIVEN** an existing RabbitMQ VIP instance
- **WHEN** user attempts to modify an immutable parameter (e.g., `zone_ids`)
- **THEN** error message explicitly states: "argument `zone_ids` cannot be changed after instance creation"
- **AND** parameter name is visible in the error
- **AND** user can quickly identify which parameter caused the error

#### Scenario: Error message includes guidance
- **GIVEN** an existing RabbitMQ VIP instance
- **WHEN** user attempts to modify an immutable parameter
- **THEN** error message includes: "Please recreate the instance if you need to modify this parameter"
- **AND** guidance is actionable and clear
- **AND** user understands how to proceed if modification is required

#### Scenario: Error message format is consistent
- **GIVEN** an existing RabbitMQ VIP instance
- **WHEN** user attempts to modify any of the 9 immutable parameters
- **THEN** all error messages follow the same format
- **AND** structure is: "argument `{parameter_name}` cannot be changed after instance creation. Please recreate the instance if you need to modify this parameter"
- **AND** user experience is consistent across different parameters

### Requirement: Backward Compatibility for Update Logic Changes
The update logic changes SHALL be fully backward compatible with existing configurations.

#### Scenario: Existing configuration continues to work
- **GIVEN** an existing RabbitMQ VIP instance managed by Terraform
- **WHEN** provider is upgraded with optimized update logic
- **THEN** `terraform plan` shows no changes for configurations that haven't changed
- **AND** existing updateable fields (`cluster_name`, `resource_tags`) continue to work
- **AND** no recreation is triggered for existing resources

#### Scenario: Gradual adoption of new update capabilities
- **GIVEN** an existing RabbitMQ VIP instance
- **WHEN** user has not previously updated `enable_public_access`, `band_width`, or `auto_renew_flag`
- **THEN** provider reads current values from API into state
- **AND** user can now modify these fields without recreating the instance
- **AND** migration is seamless and requires no manual intervention

#### Scenario: State refresh works correctly
- **GIVEN** an existing RabbitMQ VIP instance with any configuration
- **WHEN** `terraform refresh` is run
- **THEN** state correctly reflects all updateable fields from API
- **AND** `enable_public_access`, `band_width`, and `auto_renew_flag` are correctly populated
- **AND** immutable fields are correctly populated but not editable
- **AND** no conflicts or errors occur during refresh

### Requirement: Update Operation Idempotency
The update operation SHALL be idempotent and handle retry scenarios correctly.

#### Scenario: Re-apply with same configuration
- **GIVEN** an existing RabbitMQ VIP instance with configuration A
- **WHEN** user runs `terraform apply` multiple times without changing configuration
- **THEN** first apply successfully updates the instance
- **AND** subsequent applies detect no changes (`terraform plan` shows no changes)
- **AND** no API calls are made for operations with no changes
- **AND** instance remains in stable state

#### Scenario: Apply after failed update
- **GIVEN** an existing RabbitMQ VIP instance
- **WHEN** update operation fails due to transient error (e.g., network issue)
- **THEN** Terraform state remains unchanged (previous values)
- **AND** user can retry `terraform apply` immediately
- **AND** retry successfully completes the update
- **AND** no partial state corruption occurs

#### Scenario: Concurrent update attempts
- **GIVEN** an existing RabbitMQ VIP instance
- **WHEN** multiple `terraform apply` operations run concurrently (e.g., from different machines)
- **THEN** last successful operation wins
- **AND** state converges to final configuration
- **AND** no race condition errors occur
