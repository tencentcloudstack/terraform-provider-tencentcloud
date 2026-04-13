# Delta Spec: TDMQ RabbitMQ VIP Instance

This is a delta spec that adds requirements to the existing `tdmq-rabbitmq-vip-instance` spec.

## ADDED Requirements

### Requirement: Remark Field Support
The `tencentcloud_tdmq_rabbitmq_vip_instance` resource SHALL support a `remark` field for managing instance remarks.

#### Scenario: User creates instance with remark
- **GIVEN** a user defines a RabbitMQ VIP instance configuration
- **WHEN** user includes `remark = "production instance"` in the configuration
- **THEN** the instance is created with the specified remark
- **AND** the remark field is visible in Terraform state after creation
- **AND** the remark is applied to the cloud resource

#### Scenario: User creates instance without remark
- **GIVEN** a user defines a RabbitMQ VIP instance configuration
- **WHEN** user does not include `remark` field
- **THEN** the instance is created successfully without a remark
- **AND** the `remark` field is not present or is empty in Terraform state

#### Scenario: User updates remark on existing instance
- **GIVEN** an existing RabbitMQ VIP instance with `remark = "old remark"`
- **WHEN** user changes `remark` to `"new remark"` in the configuration
- **THEN** `terraform plan` shows the remark change
- **AND** `terraform apply` successfully updates the remark on the instance
- **AND** Terraform state reflects the updated remark

#### Scenario: User reads remark from cloud
- **GIVEN** a RabbitMQ VIP instance exists with a remark in Tencent Cloud
- **WHEN** Terraform performs a refresh operation
- **THEN** the `remark` field in state is populated with the current remark from the cloud
- **AND** nil values are handled gracefully

### Requirement: Deletion Protection Field Support
The `tencentcloud_tdmq_rabbitmq_vip_instance` resource SHALL support an `enable_deletion_protection` field for managing instance deletion protection.

#### Scenario: User creates instance with deletion protection enabled
- **GIVEN** a user defines a RabbitMQ VIP instance configuration
- **WHEN** user includes `enable_deletion_protection = true` in the configuration
- **THEN** the instance is created with deletion protection enabled
- **AND** the `enable_deletion_protection` field is visible in Terraform state after creation

#### Scenario: User creates instance with deletion protection disabled (default)
- **GIVEN** a user defines a RabbitMQ VIP instance configuration
- **WHEN** user does not include `enable_deletion_protection` field or sets it to `false`
- **THEN** the instance is created without deletion protection
- **AND** the `enable_deletion_protection` field in state shows `false`

#### Scenario: User updates deletion protection on existing instance
- **GIVEN** an existing RabbitMQ VIP instance with `enable_deletion_protection = false`
- **WHEN** user changes `enable_deletion_protection` to `true` in the configuration
- **THEN** `terraform plan` shows the deletion protection change
- **AND** `terraform apply` successfully enables deletion protection on the instance
- **AND** Terraform state reflects the updated value

#### Scenario: User disables deletion protection on existing instance
- **GIVEN** an existing RabbitMQ VIP instance with `enable_deletion_protection = true`
- **WHEN** user changes `enable_deletion_protection` to `false` in the configuration
- **THEN** `terraform plan` shows the deletion protection change
- **AND** `terraform apply` successfully disables deletion protection on the instance
- **AND** Terraform state reflects the updated value

### Requirement: Risk Warning Field Support
The `tencentcloud_tdmq_rabbitmq_vip_instance` resource SHALL support an `enable_risk_warning` field for managing cluster risk warning.

#### Scenario: User creates instance with risk warning enabled
- **GIVEN** a user defines a RabbitMQ VIP instance configuration
- **WHEN** user includes `enable_risk_warning = true` in the configuration
- **THEN** the instance is created with cluster risk warning enabled
- **AND** the `enable_risk_warning` field is visible in Terraform state after creation

#### Scenario: User creates instance with risk warning disabled (default)
- **GIVEN** a user defines a RabbitMQ VIP instance configuration
- **WHEN** user does not include `enable_risk_warning` field or sets it to `false`
- **THEN** the instance is created without cluster risk warning
- **AND** the `enable_risk_warning` field in state shows `false`

#### Scenario: User updates risk warning on existing instance
- **GIVEN** an existing RabbitMQ VIP instance with `enable_risk_warning = false`
- **WHEN** user changes `enable_risk_warning` to `true` in the configuration
- **THEN** `terraform plan` shows the risk warning change
- **AND** `terraform apply` successfully enables risk warning on the instance
- **AND** Terraform state reflects the updated value

#### Scenario: User disables risk warning on existing instance
- **GIVEN** an existing RabbitMQ VIP instance with `enable_risk_warning = true`
- **WHEN** user changes `enable_risk_warning` to `false` in the configuration
- **THEN** `terraform plan` shows the risk warning change
- **AND** `terraform apply` successfully disables risk warning on the instance
- **AND** Terraform state reflects the updated value

### Requirement: Update API Integration for New Fields
The resource SHALL correctly map new fields between Terraform and Tencent Cloud APIs.

#### Scenario: Create API integration for remark
- **GIVEN** a user creates an instance with `remark = "production"`
- **WHEN** Create operation calls `CreateRabbitMQVipInstance`
- **THEN** the request includes `Remark: helper.String("production")`
- **AND** the API response confirms successful creation

#### Scenario: Create API integration for deletion protection
- **GIVEN** a user creates an instance with `enable_deletion_protection = true`
- **WHEN** Create operation calls `CreateRabbitMQVipInstance`
- **THEN** the request includes `EnableDeletionProtection: helper.Bool(true)`
- **AND** the API response confirms successful creation

#### Scenario: Update API integration for remark
- **GIVEN** a user modifies the remark in an existing instance configuration
- **WHEN** Update operation detects `remark` changes via `d.HasChange()`
- **THEN** the request to `ModifyRabbitMQVipInstance` includes `Remark: helper.String("new remark")`
- **AND** the API call succeeds and the remark is updated on the cloud resource

#### Scenario: Update API integration for deletion protection
- **GIVEN** a user modifies the deletion protection in an existing instance configuration
- **WHEN** Update operation detects `enable_deletion_protection` changes via `d.HasChange()`
- **THEN** the request to `ModifyRabbitMQVipInstance` includes `EnableDeletionProtection: helper.Bool(true)`
- **AND** the API call succeeds and the deletion protection is updated on the cloud resource

#### Scenario: Update API integration for risk warning
- **GIVEN** a user modifies the risk warning in an existing instance configuration
- **WHEN** Update operation detects `enable_risk_warning` changes via `d.HasChange()`
- **THEN** the request to `ModifyRabbitMQVipInstance` includes `EnableRiskWarning: helper.Bool(true)`
- **AND** the API call succeeds and the risk warning is updated on the cloud resource

### Requirement: Schema Definition for New Fields
The resource schema SHALL define new fields with appropriate attributes.

#### Scenario: remark schema properties
- **GIVEN** the resource schema definition
- **WHEN** examining the `remark` field
- **THEN** field type is `schema.TypeString`
- **AND** field is marked as `Optional: true`
- **AND** field has description: "Instance remark."

#### Scenario: enable_deletion_protection schema properties
- **GIVEN** the resource schema definition
- **WHEN** examining the `enable_deletion_protection` field
- **THEN** field type is `schema.TypeBool`
- **AND** field is marked as `Optional: true`
- **AND** field has description: "Whether to enable deletion protection. Default is false."

#### Scenario: enable_risk_warning schema properties
- **GIVEN** the resource schema definition
- **WHEN** examining the `enable_risk_warning` field
- **THEN** field type is `schema.TypeBool`
- **AND** field is marked as `Optional: true`
- **AND** field has description: "Whether to enable cluster risk warning. Default is false."

### Requirement: Documentation Completeness for New Fields
The resource documentation SHALL clearly describe new fields.

#### Scenario: Field documentation
- **GIVEN** the resource documentation file `tdmq_rabbitmq_vip_instance.html.markdown`
- **WHEN** reviewing the arguments reference section
- **THEN** `remark` is documented with type, default value (none), and description
- **AND** `enable_deletion_protection` is documented with type, default value, and description
- **AND** `enable_risk_warning` is documented with type, default value, and description

#### Scenario: Usage example
- **GIVEN** the resource documentation file
- **WHEN** reviewing the example usage section
- **THEN** an example shows creating an instance with new fields
- **AND** the example demonstrates the relationship between these fields

### Requirement: Backward Compatibility for New Fields
The new fields SHALL be backward compatible with existing resources.

#### Scenario: Existing instance without new fields
- **GIVEN** a RabbitMQ VIP instance managed by Terraform before this feature
- **WHEN** the provider is upgraded to include new fields support
- **THEN** `terraform plan` shows no changes for resources without these fields
- **AND** existing resources continue to function normally
- **AND** state refresh correctly populates the new fields from API response with default values

#### Scenario: First-time addition of new fields to existing resource
- **GIVEN** an existing instance without `remark`, `enable_deletion_protection`, or `enable_risk_warning` in configuration
- **WHEN** user adds these fields for the first time
- **THEN** Terraform treats this as an update (not recreation)
- **AND** only the `ModifyRabbitMQVipInstance` API is called for updatable fields
- **AND** no other instance properties are affected

### Requirement: Error Handling for New Fields
The resource SHALL handle errors for new fields gracefully.

#### Scenario: API error during creation with new fields
- **GIVEN** a user creates an instance with invalid remark value
- **WHEN** the `CreateRabbitMQVipInstance` API returns an error
- **THEN** error is propagated to the user with context
- **AND** the instance creation is not completed
- **AND** no partial state is written

#### Scenario: API error during update with new fields
- **GIVEN** a user updates new fields on an existing instance
- **WHEN** the `ModifyRabbitMQVipInstance` API fails
- **THEN** error is logged and returned to the user
- **AND** Terraform state remains unchanged (previous values)
- **AND** subsequent `terraform apply` attempts to update again

#### Scenario: Nil value handling during read
- **GIVEN** API response contains nil values for new fields
- **WHEN** the Read operation processes the response
- **THEN** nil values are handled without panicking
- **AND** fields are not set in state (or set to default values)
- **AND** no error is returned
