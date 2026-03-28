# Spec: RabbitMQ Instance Dynamic Update

This specification defines the requirements for dynamically updating RabbitMQ VIP Instance specifications.

## ADDED Requirements

### Requirement: Node Specification Update
The `tencentcloud_tdmq_rabbitmq_vip_instance` resource SHALL allow updating the `node_spec` field to change instance node specifications (e.g., rabbit-vip-basic-1, rabbit-vip-profession-4c16g).

#### Scenario: User updates node specification
- **GIVEN** an existing RabbitMQ VIP instance with `node_spec = "rabbit-vip-basic-1"`
- **WHEN** the user changes `node_spec` to `rabbit-vip-profession-4c16g` in the Terraform configuration
- **THEN** `terraform plan` shows the node_spec change as an update operation
- **AND** `terraform apply` calls the `ModifyRabbitMQVipInstance` API with the new node_spec value
- **AND** the API successfully updates the instance specification
- **AND** `terraform apply` completes successfully
- **AND** subsequent `terraform show` or read operations reflect the new node_spec value

### Requirement: Node Count Update
The resource SHALL allow updating the `node_num` field to change the number of nodes in the RabbitMQ instance.

#### Scenario: User increases node count
- **GIVEN** an existing RabbitMQ VIP instance with `node_num = 1`
- **WHEN** the user changes `node_num` to `3` in the Terraform configuration
- **THEN** `terraform plan` shows the node_num change as an update operation
- **AND** `terraform apply` calls the `ModifyRabbitMQVipInstance` API with the new node_num value
- **AND** the API successfully adds nodes to the instance
- **AND** `terraform apply` completes successfully
- **AND** the instance now has 3 nodes

#### Scenario: User decreases node count
- **GIVEN** an existing RabbitMQ VIP instance with `node_num = 3`
- **WHEN** the user changes `node_num` to `1` in the Terraform configuration
- **THEN** `terraform plan` shows the node_num change as an update operation
- **AND** `terraform apply` calls the `ModifyRabbitMQVipInstance` API with the new node_num value
- **AND** the API successfully removes nodes from the instance
- **AND** `terraform apply` completes successfully
- **AND** the instance now has 1 node

### Requirement: Storage Size Update
The resource SHALL allow updating the `storage_size` field to change the storage capacity of the RabbitMQ instance.

#### Scenario: User increases storage size
- **GIVEN** an existing RabbitMQ VIP instance with `storage_size = 200`
- **WHEN** the user changes `storage_size` to `400` in the Terraform configuration
- **THEN** `terraform plan` shows the storage_size change as an update operation
- **AND** `terraform apply` calls the `ModifyRabbitMQVipInstance` API with the new storage_size value
- **AND** the API successfully increases the storage capacity
- **AND** `terraform apply` completes successfully
- **AND** the instance now has 400G storage

#### Scenario: User decreases storage size
- **GIVEN** an existing RabbitMQ VIP instance with `storage_size = 400`
- **WHEN** the user changes `storage_size` to `200` in the Terraform configuration
- **THEN** `terraform plan` shows the storage_size change as an update operation
- **AND** `terraform apply` calls the `ModifyRabbitMQVipInstance` API with the new storage_size value
- **AND** the API successfully decreases the storage capacity (if supported by Tencent Cloud)
- **AND** `terraform apply` completes successfully
- **AND** the instance now has 200G storage

### Requirement: Public Network Bandwidth Update
The resource SHALL allow updating the `band_width` field to change the public network bandwidth of the RabbitMQ instance.

#### Scenario: User increases bandwidth
- **GIVEN** an existing RabbitMQ VIP instance with `band_width = 100` and `enable_public_access = true`
- **WHEN** the user changes `band_width` to `200` in the Terraform configuration
- **THEN** `terraform plan` shows the band_width change as an update operation
- **AND** `terraform apply` calls the `ModifyRabbitMQVipInstance` API with the new band_width value
- **AND** the API successfully increases the bandwidth
- **AND** `terraform apply` completes successfully
- **AND** the instance now has 200 Mbps public bandwidth

#### Scenario: User decreases bandwidth
- **GIVEN** an existing RabbitMQ VIP instance with `band_width = 200` and `enable_public_access = true`
- **WHEN** the user changes `band_width` to `100` in the Terraform configuration
- **THEN** `terraform plan` shows the band_width change as an update operation
- **AND** `terraform apply` calls the `ModifyRabbitMQVipInstance` API with the new band_width value
- **AND** the API successfully decreases the bandwidth
- **AND** `terraform apply` completes successfully
- **AND** the instance now has 100 Mbps public bandwidth

### Requirement: Public Access Toggle
The resource SHALL allow updating the `enable_public_access` field to enable or disable public network access for the RabbitMQ instance.

#### Scenario: User enables public access
- **GIVEN** an existing RabbitMQ VIP instance with `enable_public_access = false`
- **WHEN** the user changes `enable_public_access` to `true` in the Terraform configuration
- **THEN** `terraform plan` shows the enable_public_access change as an update operation
- **AND** `terraform apply` calls the `ModifyRabbitMQVipInstance` API with `EnablePublicAccess = true`
- **AND** the API successfully enables public network access
- **AND** `terraform apply` completes successfully
- **AND** the instance now has public access enabled
- **AND** the `public_access_endpoint` field is populated in the state

#### Scenario: User disables public access
- **GIVEN** an existing RabbitMQ VIP instance with `enable_public_access = true`
- **WHEN** the user changes `enable_public_access` to `false` in the Terraform configuration
- **THEN** `terraform plan` shows the enable_public_access change as an update operation
- **AND** `terraform apply` calls the `ModifyRabbitMQVipInstance` API with `EnablePublicAccess = false`
- **AND** the API successfully disables public network access
- **AND** `terraform apply` completes successfully
- **AND** the instance now has public access disabled
- **AND** the `public_access_endpoint` field is empty or null in the state

### Requirement: Cluster Version Update
The resource SHALL allow updating the `cluster_version` field to upgrade the RabbitMQ cluster version (e.g., from 3.8.30 to 3.11.8).

#### Scenario: User upgrades cluster version
- **GIVEN** an existing RabbitMQ VIP instance with `cluster_version = "3.8.30"`
- **WHEN** the user changes `cluster_version` to `3.11.8` in the Terraform configuration
- **THEN** `terraform plan` shows the cluster_version change as an update operation
- **AND** `terraform apply` calls the `ModifyRabbitMQVipInstance` API with the new cluster_version value
- **AND** the API successfully upgrades the cluster version
- **AND** `terraform apply` completes successfully
- **AND** the instance now runs cluster version 3.11.8

#### Scenario: API does not support cluster version update
- **GIVEN** an existing RabbitMQ VIP instance with `cluster_version = "3.8.30"`
- **WHEN** the user changes `cluster_version` to `3.11.8` in the Terraform configuration
- **AND** the `ModifyRabbitMQVipInstance` API returns an error indicating cluster_version cannot be modified
- **THEN** `terraform plan` shows the cluster_version change
- **AND** `terraform apply` returns the API error to the user
- **AND** the error message indicates that cluster_version is immutable
- **AND** the instance remains unchanged

### Requirement: Mirror Queue Toggle
The resource SHALL allow updating the `enable_create_default_ha_mirror_queue` field to enable or disable mirrored queues.

#### Scenario: User enables mirror queue
- **GIVEN** an existing RabbitMQ VIP instance with `enable_create_default_ha_mirror_queue = false`
- **WHEN** the user changes `enable_create_default_ha_mirror_queue` to `true` in the Terraform configuration
- **THEN** `terraform plan` shows the enable_create_default_ha_mirror_queue change as an update operation
- **AND** `terraform apply` calls the `ModifyRabbitMQVipInstance` API with `EnableCreateDefaultHaMirrorQueue = true`
- **AND** the API successfully enables mirrored queues
- **AND** `terraform apply` completes successfully
- **AND** the instance now has mirrored queues enabled

#### Scenario: User disables mirror queue
- **GIVEN** an existing RabbitMQ VIP instance with `enable_create_default_ha_mirror_queue = true`
- **WHEN** the user changes `enable_create_default_ha_mirror_queue` to `false` in the Terraform configuration
- **THEN** `terraform plan` shows the enable_create_default_ha_mirror_queue change as an update operation
- **AND** `terraform apply` calls the `ModifyRabbitMQVipInstance` API with `EnableCreateDefaultHaMirrorQueue = false`
- **AND** the API successfully disables mirrored queues
- **AND** `terraform apply` completes successfully
- **AND** the instance now has mirrored queues disabled

### Requirement: Multiple Field Update
The resource SHALL allow updating multiple mutable fields in a single Terraform apply operation.

#### Scenario: User updates node_spec and storage_size together
- **GIVEN** an existing RabbitMQ VIP instance with `node_spec = "rabbit-vip-basic-1"` and `storage_size = 200`
- **WHEN** the user changes `node_spec` to `rabbit-vip-profession-4c16g` and `storage_size` to `400` in the Terraform configuration
- **THEN** `terraform plan` shows both field changes as update operations
- **AND** `terraform apply` calls the `ModifyRabbitMQVipInstance` API once with both updated fields
- **AND** the API successfully updates both node_spec and storage_size
- **AND** `terraform apply` completes successfully
- **AND** both fields reflect the new values in the state

### Requirement: Immutable Field Protection
The resource SHALL prevent updates to immutable fields (zone_ids, vpc_id, subnet_id, pay_mode, time_span, auto_renew_flag) and return clear error messages.

#### Scenario: User attempts to change vpc_id
- **GIVEN** an existing RabbitMQ VIP instance with `vpc_id = "vpc-12345"`
- **WHEN** the user changes `vpc_id` to `vpc-67890` in the Terraform configuration
- **THEN** `terraform plan` shows the vpc_id change
- **AND** `terraform apply` returns an error: "argument `vpc_id` cannot be changed"
- **AND** the instance is not modified
- **AND** the Terraform state remains unchanged

#### Scenario: User attempts to change zone_ids
- **GIVEN** an existing RabbitMQ VIP instance with `zone_ids = [1]`
- **WHEN** the user changes `zone_ids` to `[1, 2, 3]` in the Terraform configuration
- **THEN** `terraform plan` shows the zone_ids change
- **AND** `terraform apply` returns an error: "argument `zone_ids` cannot be changed"
- **AND** the instance is not modified
- **AND** the Terraform state remains unchanged

### Requirement: Update API Integration
The resource SHALL correctly call the `ModifyRabbitMQVipInstance` API with the appropriate parameters when updating mutable fields.

#### Scenario: API call for node_spec update
- **GIVEN** a user changes `node_spec` from `rabbit-vip-basic-1` to `rabbit-vip-profession-4c16g`
- **WHEN** the Update operation calls `ModifyRabbitMQVipInstance`
- **THEN** the request includes `InstanceId = <instance_id>`
- **AND** the request includes `NodeSpec = "rabbit-vip-profession-4c16g"`
- **AND** the API response confirms successful update
- **AND** no other fields are modified in the request

#### Scenario: API call for node_num update
- **GIVEN** a user changes `node_num` from `1` to `3`
- **WHEN** the Update operation calls `ModifyRabbitMQVipInstance`
- **THEN** the request includes `InstanceId = <instance_id>`
- **AND** the request includes `NodeNum = 3`
- **AND** the API response confirms successful update
- **AND** no other fields are modified in the request

### Requirement: Update Error Handling
The resource SHALL handle API errors gracefully during update operations and propagate meaningful error messages to users.

#### Scenario: API returns validation error for invalid node_spec
- **GIVEN** a user changes `node_spec` to an invalid value `invalid-spec`
- **WHEN** the Update operation calls `ModifyRabbitMQVipInstance`
- **AND** the API returns a validation error
- **THEN** the error is propagated to the user with context
- **AND** the Terraform state remains unchanged (previous node_spec value)
- **AND** subsequent `terraform apply` attempts the update again

#### Scenario: API returns quota exceeded error for storage_size
- **GIVEN** a user changes `storage_size` to a value exceeding account quota
- **WHEN** the Update operation calls `ModifyRabbitMQVipInstance`
- **AND** the API returns a quota exceeded error
- **THEN** the error is propagated to the user with context indicating quota limit
- **AND** the Terraform state remains unchanged
- **AND** the instance remains with the original storage_size

### Requirement: Update Retry Logic
The resource SHALL implement retry logic for update operations to handle transient API failures.

#### Scenario: Transient network error during update
- **GIVEN** a user changes `node_num` from `1` to `3`
- **WHEN** the Update operation calls `ModifyRabbitMQVipInstance`
- **AND** the API call fails with a transient network error
- **THEN** the operation is retried according to the retry policy
- **AND** on retry, the API call succeeds
- **AND** `terraform apply` completes successfully
- **AND** the instance now has 3 nodes

### Requirement: Update State Consistency
The resource SHALL ensure that Terraform state is updated to reflect the actual resource state after a successful update operation.

#### Scenario: State reflects updated node_spec
- **GIVEN** a user changes `node_spec` from `rabbit-vip-basic-1` to `rabbit-vip-profession-4c16g`
- **WHEN** `terraform apply` completes successfully
- **THEN** the Terraform state shows `node_spec = "rabbit-vip-profession-4c16g"`
- **AND** subsequent `terraform show` displays the updated value
- **AND** no drift occurs between state and actual resource

#### Scenario: State remains unchanged on failed update
- **GIVEN** a user changes `storage_size` from `200` to `1000`
- **WHEN** the Update operation fails due to API error
- **THEN** the Terraform state remains with `storage_size = 200`
- **AND** no partial state is written
- **AND** the instance remains with the original storage_size

### Requirement: Backward Compatibility
The update behavior SHALL be backward compatible with existing Terraform configurations and states.

#### Scenario: Existing resource without update operations
- **GIVEN** an existing RabbitMQ VIP instance managed by Terraform before this feature
- **WHEN** the user runs `terraform plan` with the updated provider
- **AND** the user does not change any fields in the configuration
- **THEN** `terraform plan` shows no changes
- **AND** the resource continues to function normally
- **AND** the state remains consistent

#### Scenario: Existing resource with new update capability
- **GIVEN** an existing RabbitMQ VIP instance with `node_spec = "rabbit-vip-basic-1"`
- **WHEN** the user upgrades the provider to include dynamic update support
- **AND** the user changes `node_spec` to `rabbit-vip-profession-4c16g`
- **THEN** `terraform plan` shows the update operation
- **AND** `terraform apply` successfully updates the node_spec
- **AND** no state migration or manual intervention is required
