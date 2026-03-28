# Spec Delta: TDMQ RabbitMQ VIP Instance

This document contains the delta specification for modifying the `tdmq-rabbitmq-vip-instance` capability to support dynamic updates.

## MODIFIED Requirements

### Requirement: Public Access Field Immutability
The `enable_public_access` and `band_width` fields SHALL support dynamic updates after instance creation through the ModifyRabbitMQVipInstance API.

#### Scenario: User enables public access on existing instance
- **GIVEN** an existing RabbitMQ VIP instance with `enable_public_access = false`
- **WHEN** the user changes `enable_public_access` to `true` in the configuration
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` calls `ModifyRabbitMQVipInstance` API with the new enable_public_access value
- **AND** the API successfully enables public network access
- **AND** the instance is updated without recreation
- **AND** the `public_access_endpoint` field is populated in the state

#### Scenario: User modifies bandwidth on existing instance
- **GIVEN** an existing RabbitMQ VIP instance with `band_width = 100` and `enable_public_access = true`
- **WHEN** the user changes `band_width` to `200` in the configuration
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` calls `ModifyRabbitMQVipInstance` API with the new band_width value
- **AND** the API successfully increases the bandwidth
- **AND** the instance is updated without recreation
- **AND** the Terraform state reflects the updated bandwidth value

### Requirement: Cluster Version Update Support
The `cluster_version` field SHALL support dynamic updates to allow cluster version upgrades.

#### Scenario: User upgrades cluster version
- **GIVEN** an existing RabbitMQ VIP instance with `cluster_version = "3.8.30"`
- **WHEN** the user changes `cluster_version` to "3.11.8" in the configuration
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` calls `ModifyRabbitMQVipInstance` API with the new cluster_version value
- **AND** the API successfully upgrades the cluster version
- **AND** the instance is updated without recreation
- **AND** the Terraform state reflects the updated cluster version

### Requirement: Mirror Queue Toggle Support
The `enable_create_default_ha_mirror_queue` field SHALL support dynamic updates to enable or disable mirrored queues.

#### Scenario: User enables mirror queue on existing instance
- **GIVEN** an existing RabbitMQ VIP instance with `enable_create_default_ha_mirror_queue = false`
- **WHEN** the user changes `enable_create_default_ha_mirror_queue` to `true` in the configuration
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` calls `ModifyRabbitMQVipInstance` API with the new enable_create_default_ha_mirror_queue value
- **AND** the API successfully enables mirrored queues
- **AND** the instance is updated without recreation
- **AND** the Terraform state reflects the updated setting

### Requirement: Node Specification Update Support
The `node_spec` field SHALL support dynamic updates to change instance node specifications.

#### Scenario: User updates node specification
- **GIVEN** an existing RabbitMQ VIP instance with `node_spec = "rabbit-vip-basic-1"`
- **WHEN** the user changes `node_spec` to "rabbit-vip-profession-4c16g" in the configuration
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` calls `ModifyRabbitMQVipInstance` API with the new node_spec value
- **AND** the API successfully updates the instance specification
- **AND** the instance is updated without recreation
- **AND** the Terraform state reflects the updated node_spec value

### Requirement: Node Count Update Support
The `node_num` field SHALL support dynamic updates to change the number of nodes in the instance.

#### Scenario: User increases node count
- **GIVEN** an existing RabbitMQ VIP instance with `node_num = 1`
- **WHEN** the user changes `node_num` to `3` in the configuration
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` calls `ModifyRabbitMQVipInstance` API with the new node_num value
- **AND** the API successfully adds nodes to the instance
- **AND** the instance is updated without recreation
- **AND** the Terraform state reflects the updated node count

#### Scenario: User decreases node count
- **GIVEN** an existing RabbitMQ VIP instance with `node_num = 3`
- **WHEN** the user changes `node_num` to `1` in the configuration
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` calls `ModifyRabbitMQVipInstance` API with the new node_num value
- **AND** the API successfully removes nodes from the instance
- **AND** the instance is updated without recreation
- **AND** the Terraform state reflects the updated node count

### Requirement: Storage Size Update Support
The `storage_size` field SHALL support dynamic updates to change the storage capacity of the instance.

#### Scenario: User increases storage size
- **GIVEN** an existing RabbitMQ VIP instance with `storage_size = 200`
- **WHEN** the user changes `storage_size` to `400` in the configuration
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` calls `ModifyRabbitMQVipInstance` API with the new storage_size value
- **AND** the API successfully increases the storage capacity
- **AND** the instance is updated without recreation
- **AND** the Terraform state reflects the updated storage size

## ADDED Requirements

### Requirement: Immutable Infrastructure Fields
The `zone_ids`, `vpc_id`, and `subnet_id` fields SHALL remain immutable and require instance recreation for changes.

#### Scenario: User attempts to change zone_ids
- **GIVEN** an existing RabbitMQ VIP instance with `zone_ids = [1]`
- **WHEN** the user changes `zone_ids` to `[1, 2, 3]` in the configuration
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` fails with error message "argument `zone_ids` cannot be changed"
- **AND** the instance is not modified
- **AND** the error message suggests recreating the resource

#### Scenario: User attempts to change vpc_id
- **GIVEN** an existing RabbitMQ VIP instance with `vpc_id = "vpc-12345"`
- **WHEN** the user changes `vpc_id` to "vpc-67890" in the configuration
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` fails with error message "argument `vpc_id` cannot be changed"
- **AND** the instance is not modified
- **AND** the error message suggests recreating the resource

#### Scenario: User attempts to change subnet_id
- **GIVEN** an existing RabbitMQ VIP instance with `subnet_id = "subnet-12345"`
- **WHEN** the user changes `subnet_id` to "subnet-67890" in the configuration
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` fails with error message "argument `subnet_id` cannot be changed"
- **AND** the instance is not modified
- **AND** the error message suggests recreating the resource

### Requirement: Immutable Billing Fields
The `pay_mode`, `time_span`, and `auto_renew_flag` fields SHALL remain immutable and require instance recreation for changes.

#### Scenario: User attempts to change pay_mode
- **GIVEN** an existing RabbitMQ VIP instance with `pay_mode = 1` (prepaid)
- **WHEN** the user changes `pay_mode` to `0` (postpaid) in the configuration
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` fails with error message "argument `pay_mode` cannot be changed"
- **AND** the instance is not modified
- **AND** the error message suggests recreating the resource

#### Scenario: User attempts to change time_span
- **GIVEN** an existing RabbitMQ VIP instance with `time_span = 1` (month)
- **WHEN** the user changes `time_span` to `12` (months) in the configuration
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` fails with error message "argument `time_span` cannot be changed"
- **AND** the instance is not modified
- **AND** the error message suggests recreating the resource

#### Scenario: User attempts to change auto_renew_flag
- **GIVEN** an existing RabbitMQ VIP instance with `auto_renew_flag = true`
- **WHEN** the user changes `auto_renew_flag` to `false` in the configuration
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` fails with error message "argument `auto_renew_flag` cannot be changed"
- **AND** the instance is not modified
- **AND** the error message suggests recreating the resource

### Requirement: Update Error Messages
When attempting to update immutable fields, the resource SHALL return clear error messages indicating which fields cannot be changed and suggest recreation.

#### Scenario: Clear error for immutable field change
- **GIVEN** an existing RabbitMQ VIP instance
- **WHEN** the user attempts to change an immutable field (e.g., `zone_ids`)
- **THEN** the update operation fails immediately with error: "argument `zone_ids` cannot be changed"
- **AND** the error message does not call any API
- **AND** the Terraform state remains unchanged
- **AND** subsequent operations can continue normally

### Requirement: Multiple Field Update with Mix of Mutable and Immutable
When a user changes multiple fields including both mutable and immutable fields, the update operation SHALL fail immediately with a clear error indicating the immutable field.

#### Scenario: User updates both mutable and immutable fields
- **GIVEN** an existing RabbitMQ VIP instance with `node_spec = "rabbit-vip-basic-1"` and `vpc_id = "vpc-12345"`
- **WHEN** the user changes both `node_spec` to "rabbit-vip-profession-4c16g" and `vpc_id` to "vpc-67890"
- **THEN** `terraform plan` shows both changes
- **AND** `terraform apply` fails immediately with error: "argument `vpc_id` cannot be changed"
- **AND** neither field is updated
- **AND** the instance is not modified
- **AND** the error message suggests recreating the resource
