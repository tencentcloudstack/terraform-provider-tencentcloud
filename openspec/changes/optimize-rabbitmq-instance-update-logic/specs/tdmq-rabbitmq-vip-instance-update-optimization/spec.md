# Spec: TDMQ RabbitMQ VIP Instance Update Optimization

This specification defines the requirements for optimizing the update logic of `tencentcloud_tdmq_rabbitmq_vip_instance` resource.

## ADDED Requirements

### Requirement: Dynamic Field Update Support
The resource SHALL support dynamic updates for fields that are supported by Tencent Cloud API, allowing users to modify instance properties without recreation.

#### Scenario: User updates node specification
- **GIVEN** an existing RabbitMQ VIP instance with `node_spec = "rabbit-vip-basic-1"`
- **WHEN** the user changes `node_spec` to `"rabbit-vip-profession-4c16g"` in Terraform configuration
- **THEN** `terraform plan` shows the node specification change
- **AND** `terraform apply` successfully updates the node specification via API
- **AND** the instance remains operational during the update
- **AND** the Terraform state reflects the updated node specification

#### Scenario: User updates node count
- **GIVEN** an existing RabbitMQ VIP instance with `node_num = 3`
- **WHEN** the user changes `node_num` to `5` in Terraform configuration
- **THEN** `terraform plan` shows the node count change
- **AND** `terraform apply` successfully updates the node count via API
- **AND** the instance scales to the new node count
- **AND** the Terraform state reflects the updated node count

#### Scenario: User updates storage size
- **GIVEN** an existing RabbitMQ VIP instance with `storage_size = 200`
- **WHEN** the user changes `storage_size` to `500` in Terraform configuration
- **THEN** `terraform plan` shows the storage size change
- **AND** `terraform apply` successfully updates the storage size via API
- **AND** the instance storage is expanded to the new size
- **AND** the Terraform state reflects the updated storage size

#### Scenario: User updates auto renew flag
- **GIVEN** an existing RabbitMQ VIP instance with `auto_renew_flag = false`
- **WHEN** the user changes `auto_renew_flag` to `true` in Terraform configuration
- **THEN** `terraform plan` shows the auto renew flag change
- **AND** `terraform apply` successfully updates the auto renew flag via API
- **AND** the instance is configured for automatic renewal
- **AND** the Terraform state reflects the updated auto renew flag

#### Scenario: User updates public network bandwidth
- **GIVEN** an existing RabbitMQ VIP instance with `band_width = 100` and public access enabled
- **WHEN** the user changes `band_width` to `200` in Terraform configuration
- **THEN** `terraform plan` shows the bandwidth change
- **AND** `terraform apply` successfully updates the bandwidth via API
- **AND** the instance public network bandwidth is updated
- **AND** the Terraform state reflects the updated bandwidth

### Requirement: Public Access Toggle Support
The resource SHALL support toggling public network access on/off dynamically.

#### Scenario: User enables public access on existing instance
- **GIVEN** an existing RabbitMQ VIP instance with `enable_public_access = false`
- **WHEN** the user changes `enable_public_access` to `true` in Terraform configuration
- **THEN** `terraform plan` shows the public access enable change
- **AND** `terraform apply` successfully enables public network access via API
- **AND** the instance becomes accessible from public network
- **AND** the Terraform state reflects the enabled status
- **AND** the `public_access_endpoint` computed field is populated

#### Scenario: User disables public access on existing instance
- **GIVEN** an existing RabbitMQ VIP instance with `enable_public_access = true`
- **WHEN** the user changes `enable_public_access` to `false` in Terraform configuration
- **THEN** `terraform plan` shows the public access disable change
- **AND** `terraform apply` successfully disables public network access via API
- **AND** the instance is no longer accessible from public network
- **AND** the Terraform state reflects the disabled status
- **AND** the `public_access_endpoint` computed field becomes empty

### Requirement: Multi-Field Update Support
The resource SHALL support updating multiple fields in a single Terraform apply operation.

#### Scenario: User updates multiple fields simultaneously
- **GIVEN** an existing RabbitMQ VIP instance
- **WHEN** the user changes multiple updatable fields (e.g., `node_spec`, `node_num`, `storage_size`) in Terraform configuration
- **THEN** `terraform plan` shows all field changes
- **AND** `terraform apply` successfully updates all fields via API calls
- **AND** all updates are applied in the correct order
- **AND** the Terraform state reflects all updated fields

#### Scenario: User updates both cluster name and tags
- **GIVEN** an existing RabbitMQ VIP instance
- **WHEN** the user changes both `cluster_name` and `resource_tags` in Terraform configuration
- **THEN** `terraform plan` shows both changes
- **AND** `terraform apply` successfully updates both fields
- **AND** the Terraform state reflects all updated fields

### Requirement: Immutable Field Constraints
The resource SHALL correctly enforce immutability for fields that cannot be updated after creation.

#### Scenario: User attempts to update zone IDs
- **GIVEN** an existing RabbitMQ VIP instance with `zone_ids = [1]`
- **WHEN** the user changes `zone_ids` to `[2, 3]` in Terraform configuration
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` fails with error message "argument `zone_ids` cannot be changed"
- **AND** the instance is not modified

#### Scenario: User attempts to update VPC ID
- **GIVEN** an existing RabbitMQ VIP instance with `vpc_id = "vpc-123"`
- **WHEN** the user changes `vpc_id` to `"vpc-456"` in Terraform configuration
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` fails with error message "argument `vpc_id` cannot be changed"
- **AND** the instance is not modified

#### Scenario: User attempts to update subnet ID
- **GIVEN** an existing RabbitMQ VIP instance with `subnet_id = "subnet-123"`
- **WHEN** the user changes `subnet_id` to `"subnet-456"` in Terraform configuration
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` fails with error message "argument `subnet_id` cannot be changed"
- **AND** the instance is not modified

#### Scenario: User attempts to update cluster version
- **GIVEN** an existing RabbitMQ VIP instance with `cluster_version = "3.8.30"`
- **WHEN** the user changes `cluster_version` to `"3.11.8"` in Terraform configuration
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` fails with error message "argument `cluster_version` cannot be changed"
- **AND** the instance is not modified

### Requirement: Update State Validation
The resource SHALL validate that updates have been successfully applied by reading the instance state after update.

#### Scenario: Validation confirms successful update
- **GIVEN** an existing RabbitMQ VIP instance
- **WHEN** the user updates `node_spec` and `terraform apply` completes
- **THEN** the provider calls the Read function to validate the update
- **AND** the API response shows the new `node_spec` value
- **AND** the Terraform state is updated with the confirmed value
- **AND** the apply operation reports success

#### Scenario: Validation detects update failure
- **GIVEN** an existing RabbitMQ VIP instance
- **WHEN** the user updates `node_num` and the API call appears successful but the actual value did not change
- **THEN** the provider detects the mismatch between requested and actual values
- **AND** an error is returned indicating the update may have failed
- **AND** the Terraform state is not modified to the inconsistent value

#### Scenario: Validation detects partially successful multi-field update
- **GIVEN** an existing RabbitMQ VIP instance
- **WHEN** the user updates both `node_spec` and `storage_size` and only `node_spec` is successfully applied
- **THEN** the provider detects that `storage_size` was not updated
- **AND** an error is returned indicating partial update failure
- **AND** the Terraform state reflects only the successful `node_spec` change

### Requirement: Update Error Handling and Retry
The resource SHALL implement robust error handling and retry logic for update operations.

#### Scenario: Retry on transient API failure
- **GIVEN** an existing RabbitMQ VIP instance
- **WHEN** the user updates `node_num` and the API call fails with a transient error (e.g., network timeout)
- **THEN** the provider automatically retries the API call according to retry policy
- **AND** if the retry succeeds, the update is applied
- **AND** the Terraform state reflects the updated value

#### Scenario: Clear error message for permanent API failure
- **GIVEN** an existing RabbitMQ VIP instance
- **WHEN** the user updates `storage_size` and the API returns a permanent error (e.g., insufficient quota)
- **THEN** the provider returns the error message to the user
- **AND** the error includes details about the quota limit
- **AND** the Terraform state remains unchanged
- **AND** subsequent apply operations can retry after quota issue is resolved

#### Scenario: Detailed debug logging for updates
- **GIVEN** an existing RabbitMQ VIP instance
- **WHEN** the user updates any field and TF_LOG is enabled
- **THEN** the provider logs the API request details including all parameters
- **AND** the provider logs the API response details
- **AND** the logs include the action name and instance ID
- **AND** the logs help diagnose update issues

### Requirement: Instance Status Validation Before Update
The resource SHALL validate that the instance is in an appropriate state before attempting updates.

#### Scenario: Update succeeds when instance is running
- **GIVEN** an existing RabbitMQ VIP instance in running state
- **WHEN** the user updates `node_spec`
- **THEN** the provider proceeds with the update
- **AND** the API call is made
- **AND** the update is applied successfully

#### Scenario: Update fails when instance is not running
- **GIVEN** an existing RabbitMQ VIP instance in creating or deleting state
- **WHEN** the user attempts to update `node_spec`
- **THEN** the provider checks the instance status before making API call
- **AND** an error is returned indicating the instance is not in a suitable state for updates
- **AND** no API call is made
- **AND** the Terraform state remains unchanged

### Requirement: Backward Compatibility
The resource SHALL maintain backward compatibility with existing resources and configurations.

#### Scenario: Existing resource without new updateable fields
- **GIVEN** an existing RabbitMQ VIP instance managed by Terraform before this optimization
- **WHEN** the provider is upgraded with the new update logic
- **THEN** `terraform plan` shows no changes for the existing resource
- **AND** the existing `cluster_name` and `resource_tags` updates continue to work
- **AND** the existing resource can still be updated without any configuration changes

#### Scenario: User adds new updateable field to existing resource
- **GIVEN** an existing RabbitMQ VIP instance without `node_num` specified in configuration
- **WHEN** the user adds `node_num` to the Terraform configuration
- **THEN** Terraform treats this as an update (not recreation)
- **AND** the `node_num` value is sent to the API
- **AND** no other instance properties are affected
- **AND** the Terraform state reflects the new `node_num` value

### Requirement: Idempotent Updates
The resource SHALL ensure that update operations are idempotent - applying the same configuration multiple times has the same result.

#### Scenario: Reapplying same configuration
- **GIVEN** an existing RabbitMQ VIP instance with current state matching configuration
- **WHEN** the user runs `terraform apply` without changing any values
- **THEN** `terraform plan` shows no changes
- **AND** no API calls are made
- **AND** the instance state remains unchanged

#### Scenario: Updating field to same value
- **GIVEN** an existing RabbitMQ VIP instance with `node_spec = "rabbit-vip-basic-1"`
- **WHEN** the user sets `node_spec = "rabbit-vip-basic-1"` explicitly in configuration
- **THEN** `terraform plan` shows no changes
- **AND** no API calls are made
- **AND** the instance state remains unchanged
