# Spec: TDMQ RabbitMQ Instance Update

This specification defines the requirements for optimizing the RabbitMQ VIP instance update logic to support more field modifications.

## ADDED Requirements

### Requirement: Node Spec Field Update Support
The `tencentcloud_tdmq_rabbitmq_vip_instance` resource SHALL support updating the `node_spec` field after instance creation.

#### Scenario: User upgrades node spec on existing instance
- **GIVEN** an existing RabbitMQ VIP instance with `node_spec = "rabbit-vip-basic-1"`
- **WHEN** the user changes `node_spec` to `"rabbit-vip-profession-4c16g"` in the configuration
- **THEN** `terraform plan` shows the node spec change
- **AND** `terraform apply` successfully updates the instance specification
- **AND** the Terraform state reflects the updated node_spec value
- **AND** the update operation calls the appropriate Tencent Cloud API

#### Scenario: User validates node spec values
- **GIVEN** a user provides a `node_spec` value in the configuration
- **WHEN** Terraform validates the configuration
- **THEN** the value must be one of the supported specifications: rabbit-vip-basic-5, rabbit-vip-profession-2c8g, rabbit-vip-basic-1, rabbit-vip-profession-4c16g, rabbit-vip-basic-2, rabbit-vip-profession-8c32g, rabbit-vip-basic-4, rabbit-vip-profession-16c64g
- **AND** invalid values are rejected with a clear error message

### Requirement: Node Count Field Update Support
The `tencentcloud_tdmq_rabbitmq_vip_instance` resource SHALL support updating the `node_num` field after instance creation.

#### Scenario: User increases node count on existing instance
- **GIVEN** an existing RabbitMQ VIP instance with `node_num = 3`
- **WHEN** the user changes `node_num` to `5` in the configuration
- **THEN** `terraform plan` shows the node count change
- **AND** `terraform apply` successfully updates the instance node count
- **AND** the Terraform state reflects the updated node_num value
- **AND** the update operation calls the appropriate Tencent Cloud API

#### Scenario: User decreases node count on existing instance
- **GIVEN** an existing RabbitMQ VIP instance with `node_num = 5`
- **WHEN** the user changes `node_num` to `3` in the configuration
- **THEN** `terraform plan` shows the node count change
- **AND** `terraform apply` successfully updates the instance node count
- **AND** the Terraform state reflects the updated node_num value
- **AND** the update operation calls the appropriate Tencent Cloud API

### Requirement: Storage Size Field Update Support
The `tencentcloud_tdmq_rabbitmq_vip_instance` resource SHALL support updating the `storage_size` field after instance creation.

#### Scenario: User increases storage size on existing instance
- **GIVEN** an existing RabbitMQ VIP instance with `storage_size = 200`
- **WHEN** the user changes `storage_size` to `400` in the configuration
- **THEN** `terraform plan` shows the storage size change
- **AND** `terraform apply` successfully updates the instance storage size
- **AND** the Terraform state reflects the updated storage_size value
- **AND** the update operation calls the appropriate Tencent Cloud API

#### Scenario: User validates storage size constraints
- **GIVEN** a user provides a `storage_size` value in the configuration
- **WHEN** Terraform validates the configuration
- **THEN** the value must be a positive integer
- **AND** the value must meet Tencent Cloud's minimum storage requirement
- **AND** the value must not exceed the maximum storage limit for the instance type
- **AND** invalid values are rejected with a clear error message

### Requirement: Bandwidth Field Update Support
The `tencentcloud_tdmq_rabbitmq_vip_instance` resource SHALL support updating the `band_width` field after instance creation.

#### Scenario: User increases bandwidth on existing instance
- **GIVEN** an existing RabbitMQ VIP instance with `band_width = 100`
- **WHEN** the user changes `band_width` to `200` in the configuration
- **THEN** `terraform plan` shows the bandwidth change
- **AND** `terraform apply` successfully updates the instance bandwidth
- **AND** the Terraform state reflects the updated band_width value
- **AND** the update operation calls the appropriate Tencent Cloud API

#### Scenario: User decreases bandwidth on existing instance
- **GIVEN** an existing RabbitMQ VIP instance with `band_width = 200`
- **WHEN** the user changes `band_width` to `100` in the configuration
- **THEN** `terraform plan` shows the bandwidth change
- **AND** `terraform apply` successfully updates the instance bandwidth
- **AND** the Terraform state reflects the updated band_width value
- **AND** the update operation calls the appropriate Tencent Cloud API

### Requirement: Public Access Toggle Field Update Support
The `tencentcloud_tdmq_rabbitmq_vip_instance` resource SHALL support updating the `enable_public_access` field after instance creation.

#### Scenario: User enables public access on existing instance
- **GIVEN** an existing RabbitMQ VIP instance with `enable_public_access = false`
- **WHEN** the user changes `enable_public_access` to `true` in the configuration
- **THEN** `terraform plan` shows the public access change
- **AND** `terraform apply` successfully enables public network access
- **AND** the Terraform state reflects the updated enable_public_access value
- **AND** the update operation calls the appropriate Tencent Cloud API

#### Scenario: User disables public access on existing instance
- **GIVEN** an existing RabbitMQ VIP instance with `enable_public_access = true`
- **WHEN** the user changes `enable_public_access` to `false` in the configuration
- **THEN** `terraform plan` shows the public access change
- **AND** `terraform apply` successfully disables public network access
- **AND** the Terraform state reflects the updated enable_public_access value
- **AND** the update operation calls the appropriate Tencent Cloud API

### Requirement: Enhanced Error Handling for Update Operations
The resource SHALL provide clear error messages when update operations fail.

#### Scenario: API returns error for unsupported field modification
- **GIVEN** a user attempts to modify a field that is not supported by the Tencent Cloud API
- **WHEN** the update operation calls the API
- **THEN** the error is caught and logged
- **AND** a clear error message is returned to the user indicating which field cannot be modified
- **AND** the Terraform state remains unchanged

#### Scenario: API returns error for invalid field value
- **GIVEN** a user provides an invalid value for an updatable field
- **WHEN** the update operation calls the API
- **THEN** the error is caught and logged
- **AND** a clear error message is returned to the user indicating the validation error
- **AND** the Terraform state remains unchanged

### Requirement: Backward Compatibility for Update Operations
The update logic SHALL remain backward compatible with existing Terraform configurations.

#### Scenario: Existing instance without update-capable fields
- **GIVEN** a RabbitMQ VIP instance managed by Terraform before this feature
- **WHEN** the provider is upgraded to include update support for more fields
- **THEN** `terraform plan` shows no changes for resources without field modifications
- **AND** existing resources continue to function normally
- **AND** the provider does not force any unintended updates

#### Scenario: Existing instance with previously immutable fields
- **GIVEN** a RabbitMQ VIP instance with a field that was previously immutable
- **WHEN** the user does not modify that field in the configuration
- **THEN** `terraform plan` shows no changes
- **AND** the instance is not modified
- **AND** the provider maintains the current state

### Requirement: API Integration for Update Operations
The resource SHALL correctly map update operations to Tencent Cloud APIs.

#### Scenario: Update node spec via ModifyRabbitMQVipInstance API
- **GIVEN** a user modifies the `node_spec` field in the configuration
- **WHEN** the Update operation detects the change via `d.HasChange("node_spec")`
- **THEN** the request to `ModifyRabbitMQVipInstance` includes the updated node_spec value
- **AND** the API call succeeds and the instance specification is updated

#### Scenario: Update multiple fields in single operation
- **GIVEN** a user modifies multiple updatable fields in the configuration (e.g., node_spec, node_num, storage_size)
- **WHEN** the Update operation detects changes via `d.HasChange()` for each field
- **THEN** the request to the API includes all modified fields
- **AND** the API call succeeds and all specified fields are updated
- **AND** the Terraform state reflects all updated values

### Requirement: Update Operation Idempotency
The update operation SHALL be idempotent, meaning repeated applications with the same configuration should not cause unintended changes.

#### Scenario: User applies same configuration multiple times
- **GIVEN** a RabbitMQ VIP instance with specific field values
- **WHEN** the user runs `terraform apply` without modifying the configuration
- **THEN** `terraform plan` shows no changes
- **AND** `terraform apply` does not call any update APIs
- **AND** the instance remains unchanged

#### Scenario: User applies configuration with same field values
- **GIVEN** a RabbitMQ VIP instance with `node_spec = "rabbit-vip-basic-1"`
- **WHEN** the user applies a configuration that sets `node_spec = "rabbit-vip-basic-1"` (same value)
- **THEN** `terraform plan` shows no changes for node_spec
- **AND** no update API is called for node_spec
- **AND** the instance remains unchanged
