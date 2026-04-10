# Delta Spec: TDMQ RabbitMQ VIP Instance

This delta spec modifies the existing `tdmq-rabbitmq-vip-instance` specification to support enhanced update capabilities.

## MODIFIED Requirements

### Requirement: Public Access Field Immutability
The `enable_public_access` and `band_width` fields SHALL be mutable after instance creation, allowing users to update these fields through Terraform apply operations.

#### Scenario: User enables public access on existing instance
- **GIVEN** an existing RabbitMQ VIP instance with `enable_public_access = false`
- **WHEN** the user changes `enable_public_access` to `true` in the configuration
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` successfully enables public network access
- **AND** the `public_access_endpoint` computed field is updated with the new endpoint
- **AND** the Terraform state reflects the updated `enable_public_access` value

#### Scenario: User disables public access on existing instance
- **GIVEN** an existing RabbitMQ VIP instance with `enable_public_access = true`
- **WHEN** the user changes `enable_public_access` to `false` in the configuration
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` successfully disables public network access
- **AND** the `public_access_endpoint` field is cleared or set to empty string
- **AND** the Terraform state reflects the updated `enable_public_access` value

#### Scenario: User modifies bandwidth on existing instance
- **GIVEN** an existing RabbitMQ VIP instance with `band_width = 100` and public access enabled
- **WHEN** the user changes `band_width` to `200` in the configuration
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` successfully updates the bandwidth to 200 Mbps
- **AND** the Terraform state reflects the updated bandwidth value
- **AND** the update does not require instance recreation

## ADDED Requirements

### Requirement: Node Spec Update Support
The resource SHALL allow users to update `node_spec` (node specification) after instance creation.

#### Scenario: User upgrades node spec
- **GIVEN** an existing RabbitMQ VIP instance with `node_spec = "rabbit-vip-basic-1"` (4C8G)
- **WHEN** the user changes `node_spec` to `rabbit-vip-profession-4c16g` in the configuration
- **THEN** `terraform plan` shows the node spec change
- **AND** `terraform apply` successfully upgrades the node specification
- **AND** the Terraform state reflects the updated node spec
- **AND** the update does not require instance recreation

#### Scenario: API integration for node spec update
- **GIVEN** a user modifies `node_spec` in an existing instance configuration
- **WHEN** the Update operation detects `node_spec` changes via `d.HasChange()`
- **THEN** the request to `ModifyRabbitMQVipInstance` includes the updated `NodeSpec` parameter
- **AND** the API call succeeds and node spec is updated on the cloud resource

### Requirement: Node Number Update Support
The resource SHALL allow users to update `node_num` (node count) after instance creation.

#### Scenario: User scales up node count
- **GIVEN** an existing RabbitMQ VIP instance with `node_num = 3`
- **WHEN** the user changes `node_num` to `5` in the configuration
- **THEN** `terraform plan` shows the node count increase
- **AND** `terraform apply` successfully adds nodes to the instance
- **AND** the Terraform state reflects the updated node count
- **AND** the update does not require instance recreation

#### Scenario: User scales down node count
- **GIVEN** an existing RabbitMQ VIP instance with `node_num = 5`
- **WHEN** the user changes `node_num` to `3` in the configuration
- **THEN** `terraform plan` shows the node count decrease
- **AND** `terraform apply` successfully removes nodes from the instance
- **AND** the Terraform state reflects the updated node count
- **AND** the update does not require instance recreation

#### Scenario: API integration for node number update
- **GIVEN** a user modifies `node_num` in an existing instance configuration
- **WHEN** the Update operation detects `node_num` changes via `d.HasChange()`
- **THEN** the request to `ModifyRabbitMQVipInstance` includes the updated `NodeNum` parameter
- **AND** the API call succeeds and node count is updated on the cloud resource

#### Scenario: Wait for node number update completion
- **GIVEN** a user updates `node_num` from 3 to 5
- **WHEN** the update operation is triggered
- **THEN** the provider waits for the instance to reach stable status after node scaling
- **AND** the wait loop uses a timeout (default 10x ReadRetryTimeout)
- **AND** the operation only completes when instance status indicates successful scaling
- **AND** the Terraform state is refreshed to reflect the final node count

### Requirement: Storage Size Update Support
The resource SHALL allow users to update `storage_size` after instance creation.

#### Scenario: User increases storage size
- **GIVEN** an existing RabbitMQ VIP instance with `storage_size = 200`
- **WHEN** the user changes `storage_size` to `400` in the configuration
- **THEN** `terraform plan` shows the storage size increase
- **AND** `terraform apply` successfully expands the storage
- **AND** the Terraform state reflects the updated storage size
- **AND** the update does not require instance recreation

#### Scenario: API integration for storage size update
- **GIVEN** a user modifies `storage_size` in an existing instance configuration
- **WHEN** the Update operation detects `storage_size` changes via `d.HasChange()`
- **THEN** the request to `ModifyRabbitMQVipInstance` includes the updated `StorageSize` parameter
- **AND** the API call succeeds and storage size is updated on the cloud resource

#### Scenario: Wait for storage size update completion
- **GIVEN** a user updates `storage_size` from 200 to 400
- **WHEN** the update operation is triggered
- **THEN** the provider waits for the instance to reach stable status after storage expansion
- **AND** the wait loop uses a timeout (default 10x ReadRetryTimeout)
- **AND** the operation only completes when instance status indicates successful expansion
- **AND** the Terraform state is refreshed to reflect the final storage size

### Requirement: Auto Renew Flag Update Support
The resource SHALL allow users to update `auto_renew_flag` after instance creation.

#### Scenario: User enables auto renew
- **GIVEN** an existing RabbitMQ VIP instance with `auto_renew_flag = false`
- **WHEN** the user changes `auto_renew_flag` to `true` in the configuration
- **THEN** `terraform plan` shows the auto renew flag change
- **AND** `terraform apply` successfully enables automatic renewal
- **AND** the Terraform state reflects the updated flag value
- **AND** the update does not require instance recreation

#### Scenario: User disables auto renew
- **GIVEN** an existing RabbitMQ VIP instance with `auto_renew_flag = true`
- **WHEN** the user changes `auto_renew_flag` to `false` in the configuration
- **THEN** `terraform plan` shows the auto renew flag change
- **AND** `terraform apply` successfully disables automatic renewal
- **AND** the Terraform state reflects the updated flag value
- **AND** the update does not require instance recreation

#### Scenario: API integration for auto renew flag update
- **GIVEN** a user modifies `auto_renew_flag` in an existing instance configuration
- **WHEN** the Update operation detects `auto_renew_flag` changes via `d.HasChange()`
- **THEN** the request to `ModifyRabbitMQVipInstance` includes the updated `AutoRenewFlag` parameter
- **AND** the API call succeeds and auto renew flag is updated on the cloud resource

### Requirement: Multi-Parameter Update Support
The resource SHALL support updating multiple parameters in a single Terraform apply operation.

#### Scenario: User updates multiple parameters simultaneously
- **GIVEN** an existing RabbitMQ VIP instance with `node_spec = "rabbit-vip-basic-1"`, `node_num = 3`, `band_width = 100`
- **WHEN** the user changes `node_spec` to `rabbit-vip-profession-4c16g`, `node_num` to `5`, and `band_width` to `200` in one configuration update
- **THEN** `terraform plan` shows all three parameter changes
- **AND** `terraform apply` successfully updates all parameters in a single operation
- **AND** the Terraform state reflects all updated values
- **AND** the update does not require instance recreation
- **AND** the provider calls `ModifyRabbitMQVipInstance` once with all changed parameters

#### Scenario: API integration for multi-parameter update
- **GIVEN** a user modifies multiple parameters (`node_spec`, `node_num`, `storage_size`, `band_width`, `auto_renew_flag`, `enable_public_access`) in the configuration
- **WHEN** the Update operation detects changes via `d.HasChange()` for multiple fields
- **THEN** the request to `ModifyRabbitMQVipInstance` includes all changed parameters
- **AND** unchanged parameters are omitted from the API request
- **AND** the API call succeeds and all specified parameters are updated on the cloud resource
- **AND** the operation waits for the instance to reach stable status if any of the changed parameters require async operations

### Requirement: Immutable Parameters Preservation
The resource SHALL maintain immutability for parameters that cannot be changed after creation, as determined by Tencent Cloud API capabilities.

#### Scenario: User attempts to change immutable zone_ids
- **GIVEN** an existing RabbitMQ VIP instance with specific `zone_ids`
- **WHEN** the user attempts to change `zone_ids` in the configuration
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` fails with error message "argument `zone_ids` cannot be changed"
- **AND** the instance is not modified

#### Scenario: User attempts to change immutable vpc_id or subnet_id
- **GIVEN** an existing RabbitMQ VIP instance with specific `vpc_id` and `subnet_id`
- **WHEN** the user attempts to change `vpc_id` or `subnet_id` in the configuration
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` fails with error message "argument `vpc_id` cannot be changed" or "argument `subnet_id` cannot be changed"
- **AND** the instance is not modified

#### Scenario: User attempts to change immutable cluster_version
- **GIVEN** an existing RabbitMQ VIP instance with `cluster_version = "3.8.30"`
- **WHEN** the user attempts to change `cluster_version` to "3.11.8" in the configuration
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` fails with error message "argument `cluster_version` cannot be changed"
- **AND** the instance is not modified

### Requirement: Backward Compatibility for Enhanced Updates
The enhanced update functionality SHALL be backward compatible with existing Terraform configurations and state.

#### Scenario: Existing instance with old immutability constraints
- **GIVEN** a RabbitMQ VIP instance managed by Terraform before this enhancement
- **WHEN** the provider is upgraded to include enhanced update support
- **THEN** `terraform plan` shows no changes for resources that haven't been modified
- **AND** existing resources continue to function normally
- **AND** users can now update previously immutable fields without requiring state migration

#### Scenario: Existing instance state refresh
- **GIVEN** a RabbitMQ VIP instance with existing state before this enhancement
- **WHEN** `terraform refresh` is run with the upgraded provider
- **THEN** the state is successfully refreshed
- **AND** all mutable fields are correctly populated from the API
- **AND** no state migration is required
- **AND** no errors are raised

### Requirement: Error Handling for Enhanced Updates
The resource SHALL handle errors gracefully when updating enhanced parameters.

#### Scenario: API error during node spec update
- **GIVEN** a user attempts to update `node_spec` to an invalid specification
- **WHEN** the `ModifyRabbitMQVipInstance` API returns an invalid spec error
- **THEN** the error is logged and returned to the user with context
- **AND** Terraform state remains unchanged (previous node spec value)
- **AND** subsequent `terraform apply` attempts the update again after correction

#### Scenario: API error during node number update
- **GIVEN** a user attempts to update `node_num` to a value exceeding quota
- **WHEN** the `ModifyRabbitMQVipInstance` API returns a quota error
- **THEN** the error is logged and returned to the user with context
- **AND** Terraform state remains unchanged (previous node count)
- **AND** subsequent `terraform apply` attempts the update again after correction

#### Scenario: API error during storage size update
- **GIVEN** a user attempts to update `storage_size` to a value exceeding limits
- **WHEN** the `ModifyRabbitMQVipInstance` API returns a storage limit error
- **THEN** the error is logged and returned to the user with context
- **AND** Terraform state remains unchanged (previous storage size)
- **AND** subsequent `terraform apply` attempts the update again after correction

#### Scenario: Timeout during async update operations
- **GIVEN** a user updates `node_num` which requires async scaling
- **WHEN** the instance does not reach stable status within the timeout period
- **THEN** the operation returns a timeout error to the user
- **AND** the error message indicates that the update operation timed out
- **AND** the user can retry the operation with `terraform apply`
- **AND** documentation explains the timeout behavior and retry process

### Requirement: Documentation Updates for Enhanced Updates
The resource documentation SHALL be updated to reflect the enhanced update capabilities.

#### Scenario: Documentation for mutable parameters
- **GIVEN** the resource documentation file `tdmq_rabbitmq_vip_instance.html.markdown`
- **WHEN** reviewing the arguments reference section
- **THEN** `node_spec` is documented as mutable (can be changed after creation)
- **AND** `node_num` is documented as mutable with notes about scaling behavior
- **AND** `storage_size` is documented as mutable
- **AND** `auto_renew_flag` is documented as mutable
- **AND** `enable_public_access` is documented as mutable (changed from immutable)
- **AND** `band_width` is documented as mutable (changed from immutable)

#### Scenario: Documentation for immutable parameters
- **GIVEN** the resource documentation file
- **WHEN** reviewing the arguments reference section
- **THEN** `zone_ids`, `vpc_id`, `subnet_id`, `cluster_version`, `enable_create_default_ha_mirror_queue` are documented as immutable
- **AND** each immutable parameter includes a note that changes require recreation
- **AND** users are directed to use `terraform taint` or manual delete/create for immutable parameters

#### Scenario: Update operation examples in documentation
- **GIVEN** the resource documentation file
- **WHEN** reviewing the example usage section
- **THEN** an example shows updating `node_spec` and `node_num` in a single apply
- **AND** an example shows updating `enable_public_access` and `band_width`
- **AND** an example shows updating `auto_renew_flag`
- **AND** each example includes the terraform command to apply the changes

#### Scenario: Timeout and retry documentation
- **GIVEN** the resource documentation file
- **WHEN** reviewing the advanced usage or troubleshooting section
- **THEN** documentation explains that some updates (node_num, storage_size) are async operations
- **AND** default timeout values are documented
- **AND** retry instructions are provided if timeout occurs
- **AND** common error scenarios and solutions are explained
