# Spec Delta: TDMQ RabbitMQ VIP Instance Update Logic Enhancement

This spec delta modifies the existing TDMQ RabbitMQ VIP Instance specification to support additional parameters that can be updated after instance creation.

## MODIFIED Requirements

### Requirement: Public Access Field Immutability
The `enable_public_access` and `band_width` fields SHALL be mutable after instance creation, allowing users to modify public network access configuration dynamically.

#### Scenario: User enables public access on existing instance
- **GIVEN** an existing RabbitMQ VIP instance with `enable_public_access = false`
- **WHEN** user changes `enable_public_access` to `true` in the configuration
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` successfully enables public network access
- **AND** public access endpoint becomes available
- **AND** Terraform state reflects the updated `enable_public_access` value

#### Scenario: User disables public access on existing instance
- **GIVEN** an existing RabbitMQ VIP instance with `enable_public_access = true`
- **WHEN** user changes `enable_public_access` to `false` in the configuration
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` successfully disables public network access
- **AND** public access endpoint becomes unavailable
- **AND** Terraform state reflects the updated `enable_public_access` value

#### Scenario: User modifies bandwidth on existing instance
- **GIVEN** an existing RabbitMQ VIP instance with `band_width = 100`
- **WHEN** user changes `band_width` to `200` in the configuration
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` successfully updates the public network bandwidth
- **AND** the update operation waits for the change to complete
- **AND** Terraform state reflects the updated `band_width` value

#### Scenario: User modifies both enable_public_access and band_width simultaneously
- **GIVEN** an existing RabbitMQ VIP instance with public access configuration
- **WHEN** user changes both `enable_public_access` and `band_width` values
- **THEN** `terraform plan` detects both changes
- **AND** `terraform apply` sends both parameters in a single ModifyRabbitMQVipInstance API call
- **AND** both parameters are updated successfully
- **AND** Terraform state reflects both updated values

## ADDED Requirements

### Requirement: Auto Renew Flag Update Support
The resource SHALL allow users to update the `auto_renew_flag` parameter through Terraform apply operations.

#### Scenario: User enables auto renew on existing instance
- **GIVEN** an existing RabbitMQ VIP instance with `auto_renew_flag = false`
- **WHEN** user changes `auto_renew_flag` to `true` in the configuration
- **THEN** `terraform plan` shows the auto renew flag change
- **AND** `terraform apply` successfully enables automatic renewal
- **AND** Terraform state reflects the updated `auto_renew_flag` value

#### Scenario: User disables auto renew on existing instance
- **GIVEN** an existing RabbitMQ VIP instance with `auto_renew_flag = true`
- **WHEN** user changes `auto_renew_flag` to `false` in the configuration
- **THEN** `terraform plan` shows the auto renew flag change
- **AND** `terraform apply` successfully disables automatic renewal
- **AND** Terraform state reflects the updated `auto_renew_flag` value

#### Scenario: User modifies auto renew flag with other parameters
- **GIVEN** an existing RabbitMQ VIP instance
- **WHEN** user modifies `auto_renew_flag` along with other mutable parameters (cluster_name, resource_tags, enable_public_access, band_width)
- **THEN** `terraform plan` detects all changes
- **AND** `terraform apply` sends all modified parameters in a single ModifyRabbitMQVipInstance API call
- **AND** all parameters are updated successfully
- **AND** Terraform state reflects all updated values

### Requirement: Public Access Update API Integration
The resource SHALL correctly map public access field updates to the `ModifyRabbitMQVipInstance` API.

#### Scenario: Update enable_public_access via Modify API
- **GIVEN** a user modifies `enable_public_access` in an existing instance configuration
- **WHEN** Update operation detects the change via `d.HasChange("enable_public_access")`
- **THEN** request to `ModifyRabbitMQVipInstance` includes `EnablePublicAccess` parameter with the new boolean value
- **AND** API call succeeds and the public access status is updated on the cloud resource
- **AND** operation waits for the change to complete before returning

#### Scenario: Update band_width via Modify API
- **GIVEN** a user modifies `band_width` in an existing instance configuration
- **WHEN** Update operation detects the change via `d.HasChange("band_width")`
- **THEN** request to `ModifyRabbitMQVipInstance` includes `Bandwidth` parameter with the new integer value
- **AND** API call succeeds and the public network bandwidth is updated
- **AND** operation waits for the change to complete before returning

### Requirement: Auto Renew Flag Update API Integration
The resource SHALL correctly map the auto renew flag update to the `ModifyRabbitMQVipInstance` API.

#### Scenario: Update auto_renew_flag via Modify API
- **GIVEN** a user modifies `auto_renew_flag` in an existing instance configuration
- **WHEN** Update operation detects the change via `d.HasChange("auto_renew_flag")`
- **THEN** request to `ModifyRabbitMQVipInstance` includes `AutoRenewFlag` parameter with the new boolean value
- **AND** API call succeeds and the auto renew status is updated on the cloud resource
- **AND** operation completes immediately without waiting (no async operation required)

### Requirement: Async Operation Waiting for Public Access Updates
The resource SHALL support asynchronous operation waiting for public access parameter updates that require time to complete.

#### Scenario: Wait for bandwidth update completion
- **GIVEN** a user modifies `band_width` on an existing instance
- **WHEN** Update operation calls `ModifyRabbitMQVipInstance` API
- **THEN** after successful API call, the operation polls the instance status using `DescribeRabbitMQVipInstances` API
- **AND** polling continues until `PublicNetworkTps` in the response matches the new bandwidth value
- **AND** polling timeout is set to `tccommon.WriteRetryTimeout * 2`
- **AND** upon success, the operation returns and triggers a Read operation to refresh state
- **AND** if timeout is reached, an error is returned to the user

#### Scenario: Wait for public access toggle completion
- **GIVEN** a user modifies `enable_public_access` on an existing instance
- **WHEN** Update operation calls `ModifyRabbitMQVipInstance` API
- **THEN** after successful API call, the operation polls the instance status using `DescribeRabbitMQVipInstances` API
- **AND** polling continues until `PublicDataStreamStatus` in the response matches the new expected state ("ON" for true, "OFF" for false)
- **AND** polling timeout is set to `tccommon.WriteRetryTimeout * 2`
- **AND** upon success, the operation returns and triggers a Read operation to refresh state
- **AND** if timeout is reached, an error is returned to the user

#### Scenario: No waiting for auto_renew_flag update
- **GIVEN** a user modifies `auto_renew_flag` on an existing instance
- **WHEN** Update operation calls `ModifyRabbitMQVipInstance` API
- **THEN** after successful API call, the operation does not wait for async completion
- **AND** operation immediately triggers a Read operation to refresh state
- **AND** no polling is performed for auto_renew_flag updates

### Requirement: Update Operation Error Handling and Logging
The resource SHALL provide comprehensive error handling and logging for update operations.

#### Scenario: API error during public access update
- **GIVEN** a user modifies `enable_public_access` or `band_width` on an existing instance
- **WHEN** `ModifyRabbitMQVipInstance` API returns an error
- **THEN** error is logged with detailed context including parameter name and value
- **AND** error message is returned to the user with clear explanation
- **AND** Terraform state remains unchanged (previous values are preserved)
- **AND** subsequent `terraform apply` attempts the update again

#### Scenario: API error during auto renew flag update
- **GIVEN** a user modifies `auto_renew_flag` on an existing instance
- **WHEN** `ModifyRabbitMQVipInstance` API returns an error
- **THEN** error is logged with detailed context including the new flag value
- **AND** error message is returned to the user with clear explanation
- **AND** Terraform state remains unchanged (previous flag value is preserved)
- **AND** subsequent `terraform apply` attempts the update again

#### Scenario: Timeout during async operation waiting
- **GIVEN** a user modifies `enable_public_access` or `band_width` on an existing instance
- **WHEN** API call succeeds but polling timeout is reached while waiting for completion
- **THEN** a timeout error is logged with context about the parameter being modified
- **AND** timeout error is returned to the user with clear message
- **AND** user can retry the operation or check the console for current status

### Requirement: Immutable Parameters Validation
The resource SHALL continue to enforce immutability for parameters that cannot be changed.

#### Scenario: User attempts to modify immutable parameters
- **GIVEN** an existing RabbitMQ VIP instance
- **WHEN** user changes parameters in the `immutableArgs` list (zone_ids, vpc_id, subnet_id, node_spec, node_num, storage_size, enable_create_default_ha_mirror_queue, time_span, pay_mode, cluster_version)
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` fails with error message "argument `<parameter_name>` cannot be changed"
- **AND** no API call is made to modify the instance
- **AND** error message clearly indicates which parameter cannot be changed

#### Scenario: User attempts to modify auto_renew_flag, band_width, or enable_public_access
- **GIVEN** an existing RabbitMQ VIP instance
- **WHEN** user changes `auto_renew_flag`, `band_width`, or `enable_public_access`
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` successfully processes the update
- **AND** no immutability error is raised for these parameters

### Requirement: Backward Compatibility for Update Logic Changes
The update logic changes SHALL be backward compatible with existing resources and configurations.

#### Scenario: Existing instance with public access fields
- **GIVEN** a RabbitMQ VIP instance managed by Terraform before this enhancement
- **WHEN** provider is upgraded to support public access updates
- **THEN** `terraform plan` shows no changes for instances that don't modify `enable_public_access` or `band_width`
- **AND** existing resources continue to function normally
- **AND** state refresh correctly reads the current values of these fields

#### Scenario: Existing instance with auto_renew_flag field
- **GIVEN** a RabbitMQ VIP instance managed by Terraform before this enhancement
- **WHEN** provider is upgraded to support `auto_renew_flag` updates
- **THEN** `terraform plan` shows no changes for instances that don't modify `auto_renew_flag`
- **AND** existing resources continue to function normally
- **AND** state refresh correctly reads the current value of this field

#### Scenario: First-time modification of newly mutable parameters
- **GIVEN** an existing instance with immutable parameters in the previous version
- **WHEN** user modifies `enable_public_access`, `band_width`, or `auto_renew_flag` after provider upgrade
- **THEN** Terraform treats this as an update (not recreation)
- **AND** only the `ModifyRabbitMQVipInstance` API is called
- **AND** no other instance properties are affected
- **AND** operation completes successfully with the new update logic
