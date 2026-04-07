# Delta Spec: TDMQ RabbitMQ VIP Instance - Update Logic Optimization

## MODIFIED Requirements

### Requirement: Public Access Field Immutability
The `enable_public_access` and `band_width` fields SHALL remain immutable after instance creation to prevent service disruption.

**Reason**: These fields require infrastructure-level changes that cannot be safely performed on running instances without service interruption or recreation.

**Migration**: Users who need to modify public access settings must destroy and recreate the instance using `terraform taint` or manual deletion and re-creation workflow.

### Requirement: Backward Compatibility with Public Access
The public access fields SHALL remain backward compatible with existing resources.

**Reason**: Maintaining existing behavior ensures no breaking changes for existing users.

**Migration**: No migration required; existing instances continue to work as before.

## ADDED Requirements

### Requirement: Enhanced Update Operation Support
The `tencentcloud_tdmq_rabbitmq_vip_instance` resource SHALL support updating additional fields through the `ModifyRabbitMQVipInstance` API.

#### Scenario: User updates cluster name
- **GIVEN** an existing RabbitMQ VIP instance with `cluster_name = "old-name"`
- **WHEN** the user changes `cluster_name` to `"new-name"` in the configuration
- **THEN** `terraform plan` shows the cluster name change
- **AND** `terraform apply` successfully updates the cluster name
- **AND** the Terraform state reflects the updated cluster name

#### Scenario: User updates resource tags
- **GIVEN** an existing RabbitMQ VIP instance with resource_tags
- **WHEN** the user modifies, adds, or removes resource_tags blocks
- **THEN** `terraform plan` shows the tag changes
- **AND** `terraform apply` successfully updates the tags on the instance
- **AND** the Terraform state reflects the updated tags
- **AND** tags are updated using full replacement (not incremental) via `Tags` array
- **AND** when all tags are removed, `RemoveAllTags` flag is set to true

#### Scenario: User sets remark
- **GIVEN** an existing RabbitMQ VIP instance without remark
- **WHEN** the user adds `remark = "This is a production instance"` to the configuration
- **THEN** `terraform plan` shows the remark addition
- **AND** `terraform apply` successfully sets the remark
- **AND** the Terraform state reflects the updated remark
- **AND** the remark field is Optional (can be added later)

#### Scenario: User modifies existing remark
- **GIVEN** an existing instance with `remark = "Old remark"`
- **WHEN** the user changes `remark` to `"New remark"` in the configuration
- **THEN** `terraform plan` shows the remark change
- **AND** `terraform apply` successfully updates the remark
- **AND** the Terraform state reflects the updated remark

#### Scenario: User removes remark
- **GIVEN** an existing instance with `remark = "This is a remark"`
- **WHEN** the user removes the `remark` field from the configuration
- **THEN** `terraform plan` shows the remark removal
- **AND** `terraform apply` successfully removes the remark by sending empty string or omitting the field
- **AND** the Terraform state no longer contains the remark field

#### Scenario: User enables deletion protection
- **GIVEN** an existing RabbitMQ VIP instance with deletion protection disabled
- **WHEN** the user sets `enable_deletion_protection = true` in the configuration
- **THEN** `terraform plan` shows the deletion protection change
- **AND** `terraform apply` successfully enables deletion protection
- **AND** the instance cannot be deleted without explicit user action
- **AND** the Terraform state reflects the updated setting

#### Scenario: User disables deletion protection
- **GIVEN** an existing instance with `enable_deletion_protection = true`
- **WHEN** the user sets `enable_deletion_protection = false` in the configuration
- **THEN** `terraform plan` shows the deletion protection change
- **AND** `terraform apply` successfully disables deletion protection
- **AND** the instance can be deleted normally
- **AND** the Terraform state reflects the updated setting

#### Scenario: User enables risk warning
- **GIVEN** an existing RabbitMQ VIP instance with risk warning disabled
- **WHEN** the user sets `enable_risk_warning = true` in the configuration
- **THEN** `terraform plan` shows the risk warning change
- **AND** `terraform apply` successfully enables risk warning
- **AND** the Terraform state reflects the updated setting

#### Scenario: User disables risk warning
- **GIVEN** an existing instance with `enable_risk_warning = true`
- **WHEN** the user sets `enable_risk_warning = false` in the configuration
- **THEN** `terraform plan` shows the risk warning change
- **AND** `terraform apply` successfully disables risk warning
- **AND** the Terraform state reflects the updated setting

### Requirement: Immutable Fields List
The resource SHALL correctly identify which fields cannot be modified after instance creation.

#### Scenario: User attempts to modify immutable fields
- **GIVEN** an existing RabbitMQ VIP instance
- **WHEN** the user attempts to change any of the following immutable fields:
  - `zone_ids` - availability zones
  - `vpc_id` - VPC ID
  - `subnet_id` - subnet ID
  - `node_spec` - node specifications
  - `node_num` - number of nodes
  - `storage_size` - storage size
  - `enable_create_default_ha_mirror_queue` - mirrored queue setting
  - `auto_renew_flag` - automatic renewal flag
  - `time_span` - purchase duration
  - `pay_mode` - payment method
  - `cluster_version` - cluster version
  - `band_width` - public network bandwidth
  - `enable_public_access` - public network access
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` fails with error message "argument `<field_name>` cannot be changed"
- **AND** the instance is not modified
- **AND** the error message clearly indicates which field caused the failure

#### Scenario: Update with no immutable field changes
- **GIVEN** an existing RabbitMQ VIP instance
- **WHEN** the user changes only mutable fields (cluster_name, resource_tags, remark, enable_deletion_protection, enable_risk_warning)
- **THEN** `terraform plan` shows only the mutable field changes
- **AND** `terraform apply` successfully updates the instance
- **AND** no errors are raised for immutable field checks

### Requirement: Update Operation with State Consistency
The Update operation SHALL ensure state consistency by waiting for the instance to reach a stable state after modifications.

#### Scenario: Successful update with state waiting
- **GIVEN** an existing RabbitMQ VIP instance
- **WHEN** the user modifies a mutable field (e.g., cluster_name)
- **AND** the Update operation calls `ModifyRabbitMQVipInstance` API successfully
- **THEN** the Update operation waits for the instance to reach a stable state
- **AND** the Read operation is called to refresh the state after update
- **AND** the Terraform state reflects the updated values
- **AND** the operation completes without timeout (within default timeout period)

#### Scenario: Update API timeout handling
- **GIVEN** an existing RabbitMQ VIP instance
- **WHEN** the `ModifyRabbitMQVipInstance` API call times out or retries are exhausted
- **THEN** the error is logged with clear context
- **AND** the error is returned to the user with retry information
- **AND** the Terraform state remains unchanged (previous values)
- **AND** the user can retry the update operation

#### Scenario: Update API failure handling
- **GIVEN** an existing RabbitMQ VIP instance
- **WHEN** the `ModifyRabbitMQVipInstance` API returns an error
- **THEN** the error is propagated to the user with API error details
- **AND** the Terraform state remains unchanged (previous values)
- **AND** subsequent `terraform apply` attempts the update again
- **AND** no partial state is written

### Requirement: Schema Definition for New Update Fields
The resource schema SHALL define the new update fields with appropriate attributes.

#### Scenario: remark schema properties
- **GIVEN** the resource schema definition
- **WHEN** examining the `remark` field
- **THEN** the field type is `schema.TypeString`
- **AND** the field is marked as `Optional: true`
- **AND** the field has a clear description for documentation

#### Scenario: enable_deletion_protection schema properties
- **GIVEN** the resource schema definition
- **WHEN** examining the `enable_deletion_protection` field
- **THEN** the field type is `schema.TypeBool`
- **AND** the field is marked as `Optional: true` and `Computed: true`
- **AND** the field has a clear description for documentation

#### Scenario: enable_risk_warning schema properties
- **GIVEN** the resource schema definition
- **WHEN** examining the `enable_risk_warning` field
- **THEN** the field type is `schema.TypeBool`
- **AND** the field is marked as `Optional: true` and `Computed: true`
- **AND** the field has a clear description for documentation

### Requirement: API Integration for New Update Fields
The resource SHALL correctly map new update fields between Terraform and Tencent Cloud APIs.

#### Scenario: Update API integration for remark
- **GIVEN** a user modifies the `remark` field
- **WHEN** the Update operation calls `ModifyRabbitMQVipInstance`
- **THEN** the request includes `Remark: helper.String(remark_value)`
- **AND** the API updates the instance remark
- **AND** when remark is removed, the request includes `Remark: helper.String("")` or omits the field

#### Scenario: Update API integration for deletion protection
- **GIVEN** a user modifies the `enable_deletion_protection` field
- **WHEN** the Update operation calls `ModifyRabbitMQVipInstance`
- **THEN** the request includes `EnableDeletionProtection: helper.Bool(value)`
- **AND** the API updates the deletion protection setting
- **AND** the field is only included in the request when it has changed

#### Scenario: Update API integration for risk warning
- **GIVEN** a user modifies the `enable_risk_warning` field
- **WHEN** the Update operation calls `ModifyRabbitMQVipInstance`
- **THEN** the request includes `EnableRiskWarning: helper.Bool(value)`
- **AND** the API updates the risk warning setting
- **AND** the field is only included in the request when it has changed

### Requirement: Read Operation for New Update Fields
The resource SHALL correctly read the new update fields from Tencent Cloud API responses.

#### Scenario: Read remark from API response
- **GIVEN** a RabbitMQ VIP instance exists with a remark
- **WHEN** the Read operation calls `DescribeRabbitMQVipInstance`
- **THEN** the remark value is extracted from the API response
- **AND** the value is set in Terraform state as `remark` field
- **AND** nil or empty values are handled gracefully

#### Scenario: Read deletion protection from API response
- **GIVEN** a RabbitMQ VIP instance exists with deletion protection enabled
- **WHEN** the Read operation calls `DescribeRabbitMQVipInstance`
- **THEN** the deletion protection status is extracted from the API response
- **AND** the boolean value is set in Terraform state as `enable_deletion_protection` field
- **AND** nil values default to false

#### Scenario: Read risk warning from API response
- **GIVEN** a RabbitMQ VIP instance exists with risk warning enabled
- **WHEN** the Read operation calls `DescribeRabbitMQVipInstance`
- **THEN** the risk warning status is extracted from the API response
- **AND** the boolean value is set in Terraform state as `enable_risk_warning` field
- **AND** nil values default to false

### Requirement: Error Handling for New Update Fields
The resource SHALL handle errors related to new update fields gracefully.

#### Scenario: API error during remark update
- **GIVEN** a user updates the `remark` field with invalid content
- **WHEN** the `ModifyRabbitMQVipInstance` API returns an error
- **THEN** the error is logged and returned to the user
- **AND** the Terraform state remains unchanged (previous remark value)
- **AND** subsequent `terraform apply` attempts the update again

#### Scenario: API error during deletion protection update
- **GIVEN** a user updates the `enable_deletion_protection` field
- **WHEN** the `ModifyRabbitMQVipInstance` API returns an error
- **THEN** the error is logged and returned to the user
- **AND** the Terraform state remains unchanged (previous value)
- **AND** subsequent `terraform apply` attempts the update again

#### Scenario: API error during risk warning update
- **GIVEN** a user updates the `enable_risk_warning` field
- **WHEN** the `ModifyRabbitMQVipInstance` API returns an error
- **THEN** the error is logged and returned to the user
- **AND** the Terraform state remains unchanged (previous value)
- **AND** subsequent `terraform apply` attempts the update again

### Requirement: Backward Compatibility for New Update Fields
The new update fields SHALL be backward compatible with existing resources.

#### Scenario: Existing instance without new update fields
- **GIVEN** a RabbitMQ VIP instance managed by Terraform before this feature
- **WHEN** the provider is upgraded to include new update fields support
- **THEN** `terraform plan` shows no changes for resources without these fields
- **AND** existing resources continue to function normally
- **AND** state refresh correctly populates the new fields from API response if present

#### Scenario: First-time addition of new update fields to existing resource
- **GIVEN** an existing instance without the new update fields in configuration
- **WHEN** the user adds any of the new fields (remark, enable_deletion_protection, enable_risk_warning)
- **THEN** Terraform treats this as an update (not recreation)
- **AND** only the `ModifyRabbitMQVipInstance` API is called
- **AND** no other instance properties are affected

### Requirement: Test Coverage for New Update Fields
All new update field operations SHALL have corresponding test coverage.

#### Scenario: Unit test for remark update
- **GIVEN** a test suite for the RabbitMQ VIP instance resource
- **WHEN** running the test suite
- **THEN** tests exist for creating, reading, updating, and deleting the remark field
- **AND** all tests pass successfully

#### Scenario: Unit test for deletion protection update
- **GIVEN** a test suite for the RabbitMQ VIP instance resource
- **WHEN** running the test suite
- **THEN** tests exist for enabling and disabling deletion protection
- **AND** all tests pass successfully

#### Scenario: Unit test for risk warning update
- **GIVEN** a test suite for the RabbitMQ VIP instance resource
- **WHEN** running the test suite
- **THEN** tests exist for enabling and disabling risk warning
- **AND** all tests pass successfully
