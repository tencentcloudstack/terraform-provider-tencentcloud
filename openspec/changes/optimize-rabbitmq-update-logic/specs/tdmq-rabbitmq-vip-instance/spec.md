# Delta Spec: TDMQ RabbitMQ VIP Instance Update Logic Optimization

This is a delta specification for the `tencentcloud_tdmq_rabbitmq_vip_instance` Terraform resource.

## ADDED Requirements

### Requirement: Remark Field Support
The `tencentcloud_tdmq_rabbitmq_vip_instance` resource SHALL support a `remark` field for managing instance remarks.

#### Scenario: User creates instance with remark
- **GIVEN** a user defines a RabbitMQ VIP instance configuration
- **WHEN** the user includes `remark = "Production RabbitMQ instance"` in the configuration
- **THEN** the instance is created with the specified remark
- **AND** the remark is visible in Terraform state after creation
- **AND** the remark is applied to the cloud resource

#### Scenario: User creates instance without remark
- **GIVEN** a user defines a RabbitMQ VIP instance configuration
- **WHEN** the user does not include the `remark` field
- **THEN** the instance is created successfully without a remark
- **AND** the `remark` field in state is not set or shows the API default value

#### Scenario: User reads instance remark into state
- **GIVEN** a RabbitMQ VIP instance exists with a remark in Tencent Cloud
- **WHEN** Terraform performs a refresh operation
- **THEN** the `remark` field in state is populated with the current remark from the cloud
- **AND** nil values are handled gracefully

### Requirement: Remark Update Support
The resource SHALL allow users to update the remark through Terraform apply operations.

#### Scenario: User adds remark to existing instance
- **GIVEN** an existing RabbitMQ VIP instance without remark in Terraform configuration
- **WHEN** the user adds `remark = "Updated remark"` to the configuration
- **THEN** `terraform plan` shows the remark addition
- **AND** `terraform apply` successfully updates the remark on the instance
- **AND** the Terraform state reflects the updated remark

#### Scenario: User modifies existing remark
- **GIVEN** an existing instance with `remark = "Old remark"`
- **WHEN** the user changes the remark to `remark = "New remark"`
- **THEN** `terraform plan` shows the remark change
- **AND** `terraform apply` successfully updates the remark
- **AND** the Terraform state reflects the updated remark

#### Scenario: User removes remark
- **GIVEN** an existing instance with a remark
- **WHEN** the user removes the `remark` field from the configuration
- **THEN** `terraform plan` may show the remark removal (depending on API behavior)
- **AND** the remark is either removed or reset to API default value
- **AND** the Terraform state reflects the change

### Requirement: Remark API Integration
The resource SHALL correctly map remark between Terraform and Tencent Cloud APIs.

#### Scenario: Create API integration
- **GIVEN** a user creates an instance with `remark = "Test remark"`
- **WHEN** the Create operation calls `CreateRabbitMQVipInstance`
- **THEN** the request may include remark if supported (depends on API version)
- **AND** the API response confirms successful instance creation

#### Scenario: Read API integration
- **GIVEN** a RabbitMQ VIP instance exists with a remark in Tencent Cloud
- **WHEN** the Read operation calls `DescribeRabbitMQVipInstance`
- **THEN** the remark is extracted from `response.ClusterInfo.Remark`
- **AND** the result is set in Terraform state as `remark` field
- **AND** nil values are handled gracefully

#### Scenario: Update API integration
- **GIVEN** a user modifies the remark in an existing instance configuration
- **WHEN** the Update operation detects `remark` changes via `d.HasChange()`
- **THEN** the request to `ModifyRabbitMQVipInstance` includes `Remark: helper.String("new value")`
- **AND** the API call succeeds and the remark is updated on the cloud resource

### Requirement: Enable Deletion Protection Field Support
The `tencentcloud_tdmq_rabbitmq_vip_instance` resource SHALL support an `enable_deletion_protection` field for managing deletion protection.

#### Scenario: User creates instance with deletion protection enabled
- **GIVEN** a user defines a RabbitMQ VIP instance configuration
- **WHEN** the user includes `enable_deletion_protection = true` in the configuration
- **THEN** the instance is created with deletion protection enabled
- **AND** the field is visible in Terraform state after creation

#### Scenario: User creates instance without deletion protection (default)
- **GIVEN** a user defines a RabbitMQ VIP instance configuration
- **WHEN** the user does not specify `enable_deletion_protection` field
- **THEN** the instance is created with deletion protection disabled (default)
- **AND** the `enable_deletion_protection` field in state shows `false` or API default

#### Scenario: User reads deletion protection status into state
- **GIVEN** a RabbitMQ VIP instance exists in Tencent Cloud
- **WHEN** Terraform performs a refresh operation
- **THEN** the `enable_deletion_protection` field in state is populated from the cloud
- **AND** nil values are handled gracefully

### Requirement: Enable Deletion Protection Update Support
The resource SHALL allow users to update deletion protection through Terraform apply operations.

#### Scenario: User enables deletion protection on existing instance
- **GIVEN** an existing RabbitMQ VIP instance with `enable_deletion_protection = false`
- **WHEN** the user changes `enable_deletion_protection` to `true` in the configuration
- **THEN** `terraform plan` shows the change
- **AND** `terraform apply` successfully enables deletion protection
- **AND** the Terraform state reflects the updated value

#### Scenario: User disables deletion protection on existing instance
- **GIVEN** an existing instance with `enable_deletion_protection = true`
- **WHEN** the user changes `enable_deletion_protection` to `false` in the configuration
- **THEN** `terraform plan` shows the change
- **AND** `terraform apply` successfully disables deletion protection
- **AND** the Terraform state reflects the updated value

### Requirement: Enable Deletion Protection API Integration
The resource SHALL correctly map enable_deletion_protection between Terraform and Tencent Cloud APIs.

#### Scenario: Read API integration
- **GIVEN** a RabbitMQ VIP instance exists in Tencent Cloud
- **WHEN** the Read operation calls `DescribeRabbitMQVipInstances`
- **THEN** the deletion protection status is extracted from `response.EnableDeletionProtection`
- **AND** the boolean value is set in Terraform state as `enable_deletion_protection`
- **AND** nil values are handled gracefully

#### Scenario: Update API integration
- **GIVEN** a user modifies deletion protection in an existing instance configuration
- **WHEN** the Update operation detects `enable_deletion_protection` changes via `d.HasChange()`
- **THEN** the request to `ModifyRabbitMQVipInstance` includes `EnableDeletionProtection: helper.Bool(true/false)`
- **AND** the API call succeeds and deletion protection is updated

### Requirement: Enable Risk Warning Field Support
The `tencentcloud_tdmq_rabbitmq_vip_instance` resource SHALL support an `enable_risk_warning` field for managing cluster risk warning.

#### Scenario: User creates instance with risk warning enabled
- **GIVEN** a user defines a RabbitMQ VIP instance configuration
- **WHEN** the user includes `enable_risk_warning = true` in the configuration
- **THEN** the instance is created with risk warning enabled
- **AND** the field is visible in Terraform state after creation

#### Scenario: User creates instance without risk warning (default)
- **GIVEN** a user defines a RabbitMQ VIP instance configuration
- **WHEN** the user does not specify `enable_risk_warning` field
- **THEN** the instance is created with risk warning disabled (default)
- **AND** the `enable_risk_warning` field in state shows `false` or API default

#### Scenario: User reads risk warning status into state
- **GIVEN** a RabbitMQ VIP instance exists in Tencent Cloud
- **WHEN** Terraform performs a refresh operation
- **THEN** the `enable_risk_warning` field in state is populated from the cloud
- **AND** nil values are handled gracefully

### Requirement: Enable Risk Warning Update Support
The resource SHALL allow users to update risk warning through Terraform apply operations.

#### Scenario: User enables risk warning on existing instance
- **GIVEN** an existing RabbitMQ VIP instance with `enable_risk_warning = false`
- **WHEN** the user changes `enable_risk_warning` to `true` in the configuration
- **THEN** `terraform plan` shows the change
- **AND** `terraform apply` successfully enables risk warning
- **AND** the Terraform state reflects the updated value

#### Scenario: User disables risk warning on existing instance
- **GIVEN** an existing instance with `enable_risk_warning = true`
- **WHEN** the user changes `enable_risk_warning` to `false` in the configuration
- **THEN** `terraform plan` shows the change
- **AND** `terraform apply` successfully disables risk warning
- **AND** the Terraform state reflects the updated value

### Requirement: Enable Risk Warning API Integration
The resource SHALL correctly map enable_risk_warning between Terraform and Tencent Cloud APIs.

#### Scenario: Read API integration
- **GIVEN** a RabbitMQ VIP instance exists in Tencent Cloud
- **WHEN** the Read operation calls `DescribeRabbitMQVipInstances`
- **THEN** the risk warning status is extracted from `response.EnableRiskWarning`
- **AND** the boolean value is set in Terraform state as `enable_risk_warning`
- **AND** nil values are handled gracefully

#### Scenario: Update API integration
- **GIVEN** a user modifies risk warning in an existing instance configuration
- **WHEN** the Update operation detects `enable_risk_warning` changes via `d.HasChange()`
- **THEN** the request to `ModifyRabbitMQVipInstance` includes `EnableRiskWarning: helper.Bool(true/false)`
- **AND** the API call succeeds and risk warning is updated

### Requirement: Schema Definition for New Fields
The resource schema SHALL define the new fields with appropriate attributes.

#### Scenario: remark schema properties
- **GIVEN** the resource schema definition
- **WHEN** examining the `remark` field
- **THEN** the field type is `schema.TypeString`
- **AND** the field is marked as `Optional: true` and `Computed: true`
- **AND** the field has description: "Remarks for the RabbitMQ instance."

#### Scenario: enable_deletion_protection schema properties
- **GIVEN** the resource schema definition
- **WHEN** examining the `enable_deletion_protection` field
- **THEN** the field type is `schema.TypeBool`
- **AND** the field is marked as `Optional: true` and `Computed: true`
- **AND** the field has description: "Whether to enable deletion protection. Default is false."

#### Scenario: enable_risk_warning schema properties
- **GIVEN** the resource schema definition
- **WHEN** examining the `enable_risk_warning` field
- **THEN** the field type is `schema.TypeBool`
- **AND** the field is marked as `Optional: true` and `Computed: true`
- **AND** the field has description: "Whether to enable cluster risk warning."

### Requirement: Backward Compatibility for New Fields
The new fields SHALL be backward compatible with existing resources.

#### Scenario: Existing instance without new fields
- **GIVEN** a RabbitMQ VIP instance managed by Terraform before this feature
- **WHEN** the provider is upgraded to include new field support
- **THEN** `terraform plan` shows no changes for resources without these fields
- **AND** existing resources continue to function normally
- **AND** state refresh correctly populates the new fields from API response

#### Scenario: First-time field addition to existing resource
- **GIVEN** an existing instance without new fields in configuration
- **WHEN** the user adds any new field for the first time
- **THEN** Terraform treats this as an update (not recreation)
- **AND** only the `ModifyRabbitMQVipInstance` API is called
- **AND** no other instance properties are affected

### Requirement: Error Handling for New Fields
The resource SHALL handle errors related to new fields gracefully.

#### Scenario: API error during field update
- **GIVEN** a user updates a new field on an existing instance
- **WHEN** the `ModifyRabbitMQVipInstance` API fails
- **THEN** the error is logged and returned to the user
- **AND** Terraform state remains unchanged
- **AND** subsequent `terraform apply` attempts the update again

#### Scenario: Nil value handling during read
- **GIVEN** API response contains nil values for new fields
- **WHEN** the Read operation processes the response
- **THEN** nil values are handled without panicking
- **AND** fields are not set in state (or set to default values)
- **AND** no error is returned
