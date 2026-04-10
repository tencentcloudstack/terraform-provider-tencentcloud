# Spec: TDMQ RabbitMQ Instance Update Enhancement

This specification defines the enhanced update capabilities for RabbitMQ VIP instance resources, adding support for additional parameters that can be updated through Terraform.

## ADDED Requirements

### Requirement: Remark Field Support
The `tencentcloud_tdmq_rabbitmq_vip_instance` resource SHALL support a `remark` field for managing instance remark information.

#### Scenario: User creates instance with remark
- **GIVEN** a user defines a RabbitMQ VIP instance configuration
- **WHEN** the user includes a `remark` field with a text value in the configuration
- **THEN** the instance is created with the specified remark
- **AND** the remark is visible in Terraform state after creation
- **AND** the remark is applied to the cloud resource

#### Scenario: User creates instance without remark
- **GIVEN** a user defines a RabbitMQ VIP instance configuration
- **WHEN** the user does not include the `remark` field
- **THEN** the instance is created successfully without a remark
- **AND** the `remark` field in state reflects the API default value (if any)

#### Scenario: User updates remark on existing instance
- **GIVEN** an existing RabbitMQ VIP instance with `remark = "initial value"`
- **WHEN** the user changes `remark` to "updated value" in the configuration
- **THEN** `terraform plan` shows the remark change
- **AND** `terraform apply` successfully updates the remark
- **AND** the Terraform state reflects the updated remark

#### Scenario: User reads instance remark into state
- **GIVEN** a RabbitMQ VIP instance exists with a remark in Tencent Cloud
- **WHEN** Terraform performs a refresh operation
- **THEN** the `remark` field in state is populated with the current remark from the cloud

### Requirement: Deletion Protection Field Support
The `tencentcloud_tdmq_rabbitmq_vip_instance` resource SHALL support an `enable_deletion_protection` field for controlling whether deletion protection is enabled.

#### Scenario: User creates instance with deletion protection enabled
- **GIVEN** a user defines a RabbitMQ VIP instance configuration
- **WHEN** the user sets `enable_deletion_protection = true` in the configuration
- **THEN** the instance is created with deletion protection enabled
- **AND** the `enable_deletion_protection` field shows `true` in Terraform state
- **AND** the deletion protection is applied to the cloud resource

#### Scenario: User creates instance with deletion protection disabled (default)
- **GIVEN** a user defines a RabbitMQ VIP instance configuration
- **WHEN** the user sets `enable_deletion_protection = false` or omits the field
- **THEN** the instance is created without deletion protection (default behavior)
- **AND** the `enable_deletion_protection` field in state reflects the actual value from API

#### Scenario: User enables deletion protection on existing instance
- **GIVEN** an existing RabbitMQ VIP instance with `enable_deletion_protection = false`
- **WHEN** the user changes `enable_deletion_protection` to `true` in the configuration
- **THEN** `terraform plan` shows the deletion protection change
- **AND** `terraform apply` successfully enables deletion protection
- **AND** the Terraform state reflects the updated value

#### Scenario: User disables deletion protection on existing instance
- **GIVEN** an existing RabbitMQ VIP instance with `enable_deletion_protection = true`
- **WHEN** the user changes `enable_deletion_protection` to `false` in the configuration
- **THEN** `terraform plan` shows the deletion protection change
- **AND** `terraform apply` successfully disables deletion protection
- **AND** the Terraform state reflects the updated value

### Requirement: Risk Warning Field Support
The `tencentcloud_tdmq_rabbitmq_vip_instance` resource SHALL support an `enable_risk_warning` field for controlling whether cluster risk warning is enabled.

#### Scenario: User creates instance with risk warning enabled
- **GIVEN** a user defines a RabbitMQ VIP instance configuration
- **WHEN** the user sets `enable_risk_warning = true` in the configuration
- **THEN** the instance is created with risk warning enabled
- **AND** the `enable_risk_warning` field shows `true` in Terraform state
- **AND** the risk warning setting is applied to the cloud resource

#### Scenario: User creates instance with risk warning disabled (default)
- **GIVEN** a user defines a RabbitMQ VIP instance configuration
- **WHEN** the user sets `enable_risk_warning = false` or omits the field
- **THEN** the instance is created without risk warning enabled (default behavior)
- **AND** the `enable_risk_warning` field in state reflects the actual value from API

#### Scenario: User enables risk warning on existing instance
- **GIVEN** an existing RabbitMQ VIP instance with `enable_risk_warning = false`
- **WHEN** the user changes `enable_risk_warning` to `true` in the configuration
- **THEN** `terraform plan` shows the risk warning change
- **AND** `terraform apply` successfully enables risk warning
- **AND** the Terraform state reflects the updated value

#### Scenario: User disables risk warning on existing instance
- **GIVEN** an existing RabbitMQ VIP instance with `enable_risk_warning = true`
- **WHEN** the user changes `enable_risk_warning` to `false` in the configuration
- **THEN** `terraform plan` shows the risk warning change
- **AND** `terraform apply` successfully disables risk warning
- **AND** the Terraform state reflects the updated value

### Requirement: Update API Integration
The resource SHALL correctly map the new updateable parameters between Terraform and Tencent Cloud APIs.

#### Scenario: Remark API integration for update
- **GIVEN** a user modifies the `remark` field in an existing instance configuration
- **WHEN** the Update operation detects `remark` changes via `d.HasChange()`
- **THEN** the request to `ModifyRabbitMQVipInstance` includes `Remark` field with the new value
- **AND** the API call succeeds and the remark is updated on the cloud resource

#### Scenario: Deletion protection API integration for update
- **GIVEN** a user modifies the `enable_deletion_protection` field in an existing instance configuration
- **WHEN** the Update operation detects `enable_deletion_protection` changes via `d.HasChange()`
- **THEN** the request to `ModifyRabbitMQVipInstance` includes `EnableDeletionProtection` field with the new boolean value
- **AND** the API call succeeds and the deletion protection setting is updated on the cloud resource

#### Scenario: Risk warning API integration for update
- **GIVEN** a user modifies the `enable_risk_warning` field in an existing instance configuration
- **WHEN** the Update operation detects `enable_risk_warning` changes via `d.HasChange()`
- **THEN** the request to `ModifyRabbitMQVipInstance` includes `EnableRiskWarning` field with the new boolean value
- **AND** the API call succeeds and the risk warning setting is updated on the cloud resource

#### Scenario: Multiple fields updated in single operation
- **GIVEN** a user modifies multiple updateable fields (e.g., `remark`, `enable_deletion_protection`) in one configuration change
- **WHEN** the Update operation detects changes in multiple fields
- **THEN** a single `ModifyRabbitMQVipInstance` API call includes all changed fields
- **AND** the API call succeeds and all fields are updated atomically

### Requirement: Read API Integration
The resource SHALL correctly read the new updateable parameters from Tencent Cloud API responses.

#### Scenario: Read remark from API response
- **GIVEN** a RabbitMQ VIP instance exists with a remark in Tencent Cloud
- **WHEN** the Read operation calls `DescribeRabbitMQVipInstances`
- **THEN** the remark value is extracted from `response.ClusterInfo.Remark`
- **AND** the value is set in Terraform state as `remark` field
- **AND** nil values are handled gracefully (not set in state)

#### Scenario: Read deletion protection from API response
- **GIVEN** a RabbitMQ VIP instance exists in Tencent Cloud
- **WHEN** the Read operation calls `DescribeRabbitMQVipInstances`
- **THEN** the deletion protection value is extracted from `response.ClusterInfo.EnableDeletionProtection`
- **AND** the boolean value is set in Terraform state as `enable_deletion_protection`
- **AND** nil values are handled gracefully (not set in state)

#### Scenario: Read risk warning from API response
- **GIVEN** a RabbitMQ VIP instance exists in Tencent Cloud
- **WHEN** the Read operation calls `DescribeRabbitMQVipInstances`
- **THEN** the risk warning value is extracted from `response.ClusterInfo.EnableRiskWarning`
- **AND** the boolean value is set in Terraform state as `enable_risk_warning`
- **AND** nil values are handled gracefully (not set in state)

### Requirement: Schema Definition for New Fields
The resource schema SHALL define the new updateable fields with appropriate attributes.

#### Scenario: remark schema properties
- **GIVEN** the resource schema definition
- **WHEN** examining the `remark` field
- **THEN** the field type is `schema.TypeString`
- **AND** the field is marked as `Optional: true`
- **AND** the field is marked as `Computed: true`
- **AND** the field has a clear description for documentation

#### Scenario: enable_deletion_protection schema properties
- **GIVEN** the resource schema definition
- **WHEN** examining the `enable_deletion_protection` field
- **THEN** the field type is `schema.TypeBool`
- **AND** the field is marked as `Optional: true`
- **AND** the field is marked as `Computed: true`
- **AND** the field has a clear description for documentation

#### Scenario: enable_risk_warning schema properties
- **GIVEN** the resource schema definition
- **WHEN** examining the `enable_risk_warning` field
- **THEN** the field type is `schema.TypeBool`
- **AND** the field is marked as `Optional: true`
- **AND** the field is marked as `Computed: true`
- **AND** the field has a clear description for documentation

### Requirement: Backward Compatibility
The new updateable fields SHALL be backward compatible with existing resources.

#### Scenario: Existing instance without new fields
- **GIVEN** a RabbitMQ VIP instance managed by Terraform before this feature
- **WHEN** the provider is upgraded to include the new updateable fields
- **THEN** `terraform plan` shows no changes for resources without these fields
- **AND** existing resources continue to function normally

#### Scenario: First-time addition of new fields to existing resource
- **GIVEN** an existing instance without `remark`, `enable_deletion_protection`, or `enable_risk_warning` in configuration
- **WHEN** the user adds one or more of these fields for the first time
- **THEN** Terraform treats this as an update (not recreation)
- **AND** only the `ModifyRabbitMQVipInstance` API is called
- **AND** no other instance properties are affected

#### Scenario: State refresh populates new fields from API
- **GIVEN** an existing instance without the new fields in Terraform state
- **WHEN** the user runs `terraform refresh`
- **THEN** the new fields are populated from the API response
- **AND** the state is updated without requiring user intervention

### Requirement: Error Handling
The resource SHALL handle errors related to the new updateable fields gracefully.

#### Scenario: API error during remark update
- **GIVEN** a user updates the `remark` field on an existing instance
- **WHEN** the `ModifyRabbitMQVipInstance` API fails
- **THEN** the error is logged and returned to the user
- **AND** Terraform state remains unchanged (previous value)
- **AND** subsequent `terraform apply` attempts the update again

#### Scenario: Nil value handling during read
- **GIVEN** API response contains nil `Remark`, `EnableDeletionProtection`, or `EnableRiskWarning`
- **WHEN** the Read operation processes the response
- **THEN** nil values are handled without panicking
- **AND** fields are not set in state (or set to default values)
- **AND** no error is returned for nil values

#### Scenario: Invalid value handling
- **GIVEN** a user provides an invalid value for `remark` (e.g., exceeds length limit)
- **WHEN** the API validation rejects the value
- **THEN** the error is propagated to the user with context
- **AND** the update operation fails
- **AND** no partial state is written

### Requirement: Immutable Fields Validation
The resource SHALL maintain the existing immutable fields validation logic and ensure new updateable fields are not incorrectly marked as immutable.

#### Scenario: Immutable fields remain unchanged
- **GIVEN** the resource update function
- **WHEN** a user attempts to modify an immutable field (e.g., `zone_ids`, `node_spec`)
- **THEN** the update fails with error "argument `<field_name>` cannot be changed"
- **AND** the immutable fields list includes all previously immutable fields

#### Scenario: New updateable fields are not marked as immutable
- **GIVEN** the resource update function
- **WHEN** a user modifies `remark`, `enable_deletion_protection`, or `enable_risk_warning`
- **THEN** the update succeeds without triggering immutable field validation error
- **AND** the new fields are not included in the immutable fields list

### Requirement: Code Formatting
All code changes SHALL be formatted using `go fmt`.

#### Scenario: File formatting after modification
- **GIVEN** code changes are made to the resource file
- **WHEN** all changes are complete
- **THEN** `go fmt` is executed on the file
- **AND** all Go code adheres to standard formatting rules
- **AND** no formatting warnings or errors exist
