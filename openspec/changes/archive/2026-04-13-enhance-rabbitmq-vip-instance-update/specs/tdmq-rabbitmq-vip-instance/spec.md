## ADDED Requirements

### Requirement: Remark Field Support
The `tencentcloud_tdmq_rabbitmq_vip_instance` resource SHALL support a `remark` field for managing instance remarks.

#### Scenario: User creates instance with remark
- **GIVEN** a user defines a RabbitMQ VIP instance configuration
- **WHEN** the user includes `remark = "Production instance for team A"` in the configuration
- **THEN** the instance is created with the specified remark
- **AND** the remark is visible in Terraform state after creation
- **AND** the remark is applied to the cloud resource

#### Scenario: User updates instance remark
- **GIVEN** an existing RabbitMQ VIP instance with remark in Terraform configuration
- **WHEN** the user changes the remark value to "Updated remark for team A"
- **THEN** `terraform plan` shows the remark change
- **AND** `terraform apply` successfully updates the remark on the instance
- **AND** the Terraform state reflects the updated remark

#### Scenario: User reads instance remark into state
- **GIVEN** a RabbitMQ VIP instance exists with a remark in Tencent Cloud
- **WHEN** Terraform performs a refresh operation
- **THEN** the `remark` field in state is populated with the current remark from the cloud
- **AND** nil remark values are handled gracefully

#### Scenario: User removes remark from configuration
- **GIVEN** an existing instance with `remark` field in configuration
- **WHEN** the user removes the `remark` field from the configuration
- **THEN** `terraform plan` shows no changes to remark
- **AND** the remark on the cloud resource remains unchanged

### Requirement: Remark API Integration
The resource SHALL correctly map the remark field between Terraform and Tencent Cloud APIs.

#### Scenario: Create API integration for remark
- **GIVEN** a user creates an instance with `remark = "Test remark"`
- **WHEN** the Create operation calls `CreateRabbitMQVipInstance`
- **THEN** the request includes remark in the parameters
- **AND** the API response confirms successful instance creation

#### Scenario: Read API integration for remark
- **GIVEN** a RabbitMQ VIP instance exists with a remark in Tencent Cloud
- **WHEN** the Read operation calls `DescribeRabbitMQVipInstances`
- **THEN** the remark is extracted from the API response
- **AND** the value is set in Terraform state as `remark` field
- **AND** nil values are handled gracefully (not set in state)

#### Scenario: Update API integration for remark
- **GIVEN** a user modifies the remark in an existing instance configuration
- **WHEN** the Update operation detects `remark` changes via `d.HasChange()`
- **THEN** the request to `ModifyRabbitMQVipInstance` includes the new remark value
- **AND** the API call succeeds and the remark is updated on the cloud resource

### Requirement: Enable Deletion Protection Field Support
The `tencentcloud_tdmq_rabbitmq_vip_instance` resource SHALL support an `enable_deletion_protection` field for controlling instance deletion protection.

#### Scenario: User creates instance with deletion protection enabled
- **GIVEN** a user defines a RabbitMQ VIP instance configuration
- **WHEN** the user includes `enable_deletion_protection = true` in the configuration
- **THEN** the instance is created with deletion protection enabled
- **AND** the field is visible in Terraform state after creation
- **AND** the deletion protection is applied to the cloud resource

#### Scenario: User updates deletion protection
- **GIVEN** an existing RabbitMQ VIP instance with `enable_deletion_protection = false`
- **WHEN** the user changes `enable_deletion_protection` to `true`
- **THEN** `terraform plan` shows the deletion protection change
- **AND** `terraform apply` successfully updates the deletion protection on the instance
- **AND** the Terraform state reflects the updated value

#### Scenario: User reads deletion protection into state
- **GIVEN** a RabbitMQ VIP instance exists with deletion protection configuration in Tencent Cloud
- **WHEN** Terraform performs a refresh operation
- **THEN** the `enable_deletion_protection` field in state is populated with the current value from the cloud
- **AND** nil values are handled gracefully

### Requirement: Enable Deletion Protection API Integration
The resource SHALL correctly map the enable_deletion_protection field between Terraform and Tencent Cloud APIs.

#### Scenario: Create API integration for deletion protection
- **GIVEN** a user creates an instance with `enable_deletion_protection = true`
- **WHEN** the Create operation calls `CreateRabbitMQVipInstance`
- **THEN** the request includes the deletion protection flag in the parameters
- **AND** the API response confirms successful instance creation

#### Scenario: Read API integration for deletion protection
- **GIVEN** a RabbitMQ VIP instance exists with deletion protection configuration in Tencent Cloud
- **WHEN** the Read operation calls `DescribeRabbitMQVipInstances`
- **THEN** the deletion protection status is extracted from the API response
- **AND** the value is set in Terraform state as `enable_deletion_protection` field
- **AND** nil values are handled gracefully

#### Scenario: Update API integration for deletion protection
- **GIVEN** a user modifies the deletion protection in an existing instance configuration
- **WHEN** the Update operation detects `enable_deletion_protection` changes via `d.HasChange()`
- **THEN** the request to `ModifyRabbitMQVipInstance` includes the new deletion protection value
- **AND** the API call succeeds and the deletion protection is updated on the cloud resource

### Requirement: Enable Risk Warning Field Support
The `tencentcloud_tdmq_rabbitmq_vip_instance` resource SHALL support an `enable_risk_warning` field for controlling cluster risk warning.

#### Scenario: User creates instance with risk warning enabled
- **GIVEN** a user defines a RabbitMQ VIP instance configuration
- **WHEN** the user includes `enable_risk_warning = true` in the configuration
- **THEN** the instance is created with cluster risk warning enabled
- **AND** the field is visible in Terraform state after creation
- **AND** the risk warning configuration is applied to the cloud resource

#### Scenario: User updates risk warning
- **GIVEN** an existing RabbitMQ VIP instance with `enable_risk_warning = false`
- **WHEN** the user changes `enable_risk_warning` to `true`
- **THEN** `terraform plan` shows the risk warning change
- **AND** `terraform apply` successfully updates the risk warning on the instance
- **AND** the Terraform state reflects the updated value

#### Scenario: User reads risk warning into state
- **GIVEN** a RabbitMQ VIP instance exists with risk warning configuration in Tencent Cloud
- **WHEN** Terraform performs a refresh operation
- **THEN** the `enable_risk_warning` field in state is populated with the current value from the cloud
- **AND** nil values are handled gracefully

### Requirement: Enable Risk Warning API Integration
The resource SHALL correctly map the enable_risk_warning field between Terraform and Tencent Cloud APIs.

#### Scenario: Create API integration for risk warning
- **GIVEN** a user creates an instance with `enable_risk_warning = true`
- **WHEN** the Create operation calls `CreateRabbitMQVipInstance`
- **THEN** the request includes the risk warning flag in the parameters
- **AND** the API response confirms successful instance creation

#### Scenario: Read API integration for risk warning
- **GIVEN** a RabbitMQ VIP instance exists with risk warning configuration in Tencent Cloud
- **WHEN** the Read operation calls `DescribeRabbitMQVipInstances`
- **THEN** the risk warning status is extracted from the API response
- **AND** the value is set in Terraform state as `enable_risk_warning` field
- **AND** nil values are handled gracefully

#### Scenario: Update API integration for risk warning
- **GIVEN** a user modifies the risk warning in an existing instance configuration
- **WHEN** the Update operation detects `enable_risk_warning` changes via `d.HasChange()`
- **THEN** the request to `ModifyRabbitMQVipInstance` includes the new risk warning value
- **AND** the API call succeeds and the risk warning is updated on the cloud resource

### Requirement: Schema Definition for New Fields
The resource schema SHALL define the new fields with appropriate attributes.

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
- **AND** the field is marked as `Optional: true`
- **AND** the field has a clear description for documentation

#### Scenario: enable_risk_warning schema properties
- **GIVEN** the resource schema definition
- **WHEN** examining the `enable_risk_warning` field
- **THEN** the field type is `schema.TypeBool`
- **AND** the field is marked as `Optional: true`
- **AND** the field has a clear description for documentation

### Requirement: Backward Compatibility with New Fields
The new fields SHALL be backward compatible with existing resources.

#### Scenario: Existing instance without new fields
- **GIVEN** a RabbitMQ VIP instance managed by Terraform before this feature
- **WHEN** the provider is upgraded to include the new fields
- **THEN** `terraform plan` shows no changes for resources without these fields
- **AND** existing resources continue to function normally
- **AND** state refresh correctly populates the new fields from API response if they exist in the cloud

#### Scenario: First-time addition of new fields to existing resource
- **GIVEN** an existing instance without the new fields in configuration
- **WHEN** the user adds any of the new fields (remark, enable_deletion_protection, enable_risk_warning) for the first time
- **THEN** Terraform treats this as an update (not recreation)
- **AND** only the `ModifyRabbitMQVipInstance` API is called
- **AND** no other instance properties are affected

### Requirement: Error Handling for New Fields
The resource SHALL handle errors related to the new fields gracefully.

#### Scenario: API error during update of new fields
- **GIVEN** a user updates any of the new fields on an existing instance
- **WHEN** the `ModifyRabbitMQVipInstance` API fails
- **THEN** the error is logged and returned to the user
- **AND** Terraform state remains unchanged (previous values)
- **AND** subsequent `terraform apply` attempts the update again

#### Scenario: Nil value handling during read
- **GIVEN** API response contains nil values for any of the new fields
- **WHEN** the Read operation processes the response
- **THEN** nil values are handled without panicking
- **AND** fields are not set in state (or set to default values)
- **AND** no error is returned
