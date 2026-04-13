# Spec: TDMQ RabbitMQ VIP Instance - Delta

This delta specification adds new field support to the `tencentcloud_tdmq_rabbitmq_vip_instance` Terraform resource for enhanced instance management capabilities.

## ADDED Requirements

### Requirement: Remark Field Support
The `tencentcloud_tdmq_rabbitmq_vip_instance` resource SHALL support a `remark` field for managing instance remarks/notes.

#### Scenario: User creates instance with remark
- **GIVEN** a user defines a RabbitMQ VIP instance configuration
- **WHEN** user includes `remark = "production cluster"` in configuration
- **THEN** instance is created with the specified remark
- **AND** remark is visible in Terraform state after creation
- **AND** remark is applied to the cloud resource

#### Scenario: User creates instance without remark
- **GIVEN** a user defines a RabbitMQ VIP instance configuration
- **WHEN** user does not include `remark` field
- **THEN** instance is created successfully without remark
- **AND** `remark` field in state reflects API response (may be empty or default)

#### Scenario: User updates instance remark
- **GIVEN** an existing RabbitMQ VIP instance with `remark = "initial note"`
- **WHEN** user changes `remark` to `"updated note"` in configuration
- **THEN** `terraform plan` shows remark change
- **AND** `terraform apply` successfully updates remark on instance
- **AND** Terraform state reflects updated remark

### Requirement: Remark API Integration
The resource SHALL correctly map remark between Terraform and Tencent Cloud APIs.

#### Scenario: Create API integration
- **GIVEN** a user creates an instance with `remark = "test instance"`
- **WHEN** Create operation calls `CreateRabbitMQVipInstance`
- **THEN** request includes `Remark: helper.String("test instance")`
- **AND** API response confirms successful creation

#### Scenario: Read API integration
- **GIVEN** a RabbitMQ VIP instance exists with remark in Tencent Cloud
- **WHEN** Read operation calls `DescribeRabbitMQVipInstances`
- **THEN** remark is extracted from `response.ClusterInfo.Remark`
- **AND** result is set in Terraform state as `remark` field

#### Scenario: Update API integration
- **GIVEN** a user modifies remark in an existing instance configuration
- **WHEN** Update operation detects `remark` changes via `d.HasChange()`
- **THEN** request to `ModifyRabbitMQVipInstance` includes `Remark` field
- **AND** API call succeeds and remark is updated on the cloud resource

### Requirement: Enable Deletion Protection Field Support
The `tencentcloud_tdmq_rabbitmq_vip_instance` resource SHALL support an `enable_deletion_protection` field for managing instance deletion protection.

#### Scenario: User creates instance with deletion protection enabled
- **GIVEN** a user defines a RabbitMQ VIP instance configuration
- **WHEN** user sets `enable_deletion_protection = true` in configuration
- **THEN** instance is created with deletion protection enabled
- **AND** field is visible in Terraform state after creation
- **AND** deletion protection is applied to the cloud resource

#### Scenario: User creates instance with deletion protection disabled (default)
- **GIVEN** a user defines a RabbitMQ VIP instance configuration
- **WHEN** user does not include `enable_deletion_protection` field or sets it to `false`
- **THEN** instance is created with deletion protection disabled
- **AND** `enable_deletion_protection` field in state reflects actual API value

#### Scenario: User updates deletion protection status
- **GIVEN** an existing RabbitMQ VIP instance with `enable_deletion_protection = false`
- **WHEN** user changes `enable_deletion_protection` to `true` in configuration
- **THEN** `terraform plan` shows deletion protection change
- **AND** `terraform apply` successfully updates deletion protection on instance
- **AND** Terraform state reflects updated deletion protection status

### Requirement: Enable Deletion Protection API Integration
The resource SHALL correctly map enable_deletion_protection between Terraform and Tencent Cloud APIs.

#### Scenario: Create API integration
- **GIVEN** a user creates an instance with `enable_deletion_protection = true`
- **WHEN** Create operation calls `CreateRabbitMQVipInstance`
- **THEN** request includes `EnableDeletionProtection: helper.Bool(true)`
- **AND** API response confirms successful creation

#### Scenario: Read API integration
- **GIVEN** a RabbitMQ VIP instance exists with deletion protection in Tencent Cloud
- **WHEN** Read operation calls `DescribeRabbitMQVipInstances`
- **THEN** enable_deletion_protection is extracted from `response.ClusterInfo.EnableDeletionProtection`
- **AND** boolean value is set in Terraform state as `enable_deletion_protection` field
- **AND** nil values are handled gracefully (not set in state)

#### Scenario: Update API integration
- **GIVEN** a user modifies enable_deletion_protection in an existing instance configuration
- **WHEN** Update operation detects `enable_deletion_protection` changes via `d.HasChange()`
- **THEN** request to `ModifyRabbitMQVipInstance` includes `EnableDeletionProtection` field
- **AND** API call succeeds and deletion protection is updated on the cloud resource

### Requirement: Enable Risk Warning Field Support
The `tencentcloud_tdmq_rabbitmq_vip_instance` resource SHALL support an `enable_risk_warning` field for managing cluster risk warning settings.

#### Scenario: User creates instance with risk warning enabled
- **GIVEN** a user defines a RabbitMQ VIP instance configuration
- **WHEN** user sets `enable_risk_warning = true` in configuration
- **THEN** instance is created with risk warning enabled
- **AND** field is visible in Terraform state after creation
- **AND** risk warning is applied to the cloud resource

#### Scenario: User creates instance with risk warning disabled (default)
- **GIVEN** a user defines a RabbitMQ VIP instance configuration
- **WHEN** user does not include `enable_risk_warning` field or sets it to `false`
- **THEN** instance is created with risk warning disabled (API default)
- **AND** `enable_risk_warning` field in state reflects actual API value

#### Scenario: User updates risk warning status
- **GIVEN** an existing RabbitMQ VIP instance with `enable_risk_warning = false`
- **WHEN** user changes `enable_risk_warning` to `true` in configuration
- **THEN** `terraform plan` shows risk warning change
- **AND** `terraform apply` successfully updates risk warning on instance
- **AND** Terraform state reflects updated risk warning status

### Requirement: Enable Risk Warning API Integration
The resource SHALL correctly map enable_risk_warning between Terraform and Tencent Cloud APIs.

#### Scenario: Create API integration
- **GIVEN** a user creates an instance with `enable_risk_warning = true`
- **WHEN** Create operation calls `CreateRabbitMQVipInstance`
- **THEN** request includes `EnableRiskWarning: helper.Bool(true)`
- **AND** API response confirms successful creation

#### Scenario: Read API integration
- **GIVEN** a RabbitMQ VIP instance exists with risk warning setting in Tencent Cloud
- **WHEN** Read operation calls `DescribeRabbitMQVipInstances`
- **THEN** enable_risk_warning is extracted from `response.ClusterInfo.EnableRiskWarning`
- **AND** boolean value is set in Terraform state as `enable_risk_warning` field
- **AND** nil values are handled gracefully (not set in state)

#### Scenario: Update API integration
- **GIVEN** a user modifies enable_risk_warning in an existing instance configuration
- **WHEN** Update operation detects `enable_risk_warning` changes via `d.HasChange()`
- **THEN** request to `ModifyRabbitMQVipInstance` includes `EnableRiskWarning` field
- **AND** API call succeeds and risk warning is updated on the cloud resource

### Requirement: Schema Definition for New Fields
The resource schema SHALL define new fields with appropriate attributes.

#### Scenario: remark schema properties
- **GIVEN** the resource schema definition
- **WHEN** examining the `remark` field
- **THEN** field type is `schema.TypeString`
- **AND** field is marked as `Optional: true` and `Computed: true`
- **AND** field has a clear description for documentation

#### Scenario: enable_deletion_protection schema properties
- **GIVEN** the resource schema definition
- **WHEN** examining the `enable_deletion_protection` field
- **THEN** field type is `schema.TypeBool`
- **AND** field is marked as `Optional: true` and `Computed: true`
- **AND** field has a clear description for documentation

#### Scenario: enable_risk_warning schema properties
- **GIVEN** the resource schema definition
- **WHEN** examining the `enable_risk_warning` field
- **THEN** field type is `schema.TypeBool`
- **AND** field is marked as `Optional: true` and `Computed: true`
- **AND** field has a clear description for documentation

### Requirement: Update Logic for New Fields
The resource update operation SHALL support updating the new fields.

#### Scenario: Update supports remark changes
- **GIVEN** an existing RabbitMQ VIP instance
- **WHEN** user changes `remark` field in configuration
- **THEN** update operation detects the change via `d.HasChange("remark")`
- **AND** `ModifyRabbitMQVipInstance` API is called with new remark value
- **AND** Terraform state is refreshed after successful update

#### Scenario: Update supports enable_deletion_protection changes
- **GIVEN** an existing RabbitMQ VIP instance
- **WHEN** user changes `enable_deletion_protection` field in configuration
- **THEN** update operation detects the change via `d.HasChange("enable_deletion_protection")`
- **AND** `ModifyRabbitMQVipInstance` API is called with new deletion protection value
- **AND** Terraform state is refreshed after successful update

#### Scenario: Update supports enable_risk_warning changes
- **GIVEN** an existing RabbitMQ VIP instance
- **WHEN** user changes `enable_risk_warning` field in configuration
- **THEN** update operation detects the change via `d.HasChange("enable_risk_warning")`
- **AND** `ModifyRabbitMQVipInstance` API is called with new risk warning value
- **AND** Terraform state is refreshed after successful update

#### Scenario: Update handles multiple field changes simultaneously
- **GIVEN** an existing RabbitMQ VIP instance
- **WHEN** user changes multiple new fields (remark, enable_deletion_protection, enable_risk_warning) in configuration
- **THEN** update operation detects all changes
- **AND** `ModifyRabbitMQVipInstance` API is called once with all modified fields
- **AND** all fields are updated in a single API call

### Requirement: Backward Compatibility
The new fields SHALL be backward compatible with existing resources.

#### Scenario: Existing instance without new fields
- **GIVEN** a RabbitMQ VIP instance managed by Terraform before this feature
- **WHEN** provider is upgraded to include new field support
- **THEN** `terraform plan` shows no changes for resources without these fields
- **AND** existing resources continue to function normally
- **AND** state refresh correctly populates new fields from API response

#### Scenario: First-time field addition to existing resource
- **GIVEN** an existing instance without new fields in configuration
- **WHEN** user adds `remark`, `enable_deletion_protection`, or `enable_risk_warning` field for the first time
- **THEN** Terraform treats this as an update (not recreation)
- **AND** only the `ModifyRabbitMQVipInstance` API is called
- **AND** no other instance properties are affected

### Requirement: Error Handling
The resource SHALL handle errors related to new fields gracefully.

#### Scenario: API error during field update
- **GIVEN** a user updates a new field on an existing instance
- **WHEN** the `ModifyRabbitMQVipInstance` API fails
- **THEN** error is logged and returned to the user
- **AND** Terraform state remains unchanged (previous field values)
- **AND** subsequent `terraform apply` attempts to update again

#### Scenario: Nil value handling during read
- **GIVEN** API response contains nil values for new fields
- **WHEN** Read operation processes the response
- **THEN** nil values are handled without panicking
- **AND** fields are not set in state (or set to default values)
- **AND** no error is returned

### Requirement: Documentation Completeness
The resource documentation SHALL clearly describe the new fields.

#### Scenario: Field documentation
- **GIVEN** the resource documentation file `tdmq_rabbitmq_vip_instance.html.markdown`
- **WHEN** reviewing the arguments reference section
- **THEN** `remark` is documented with type and description
- **AND** `enable_deletion_protection` is documented with type and description
- **AND** `enable_risk_warning` is documented with type and description
- **AND** all fields are marked as optional and updateable

#### Scenario: Usage example
- **GIVEN** the resource documentation file
- **WHEN** reviewing the example usage section
- **THEN** examples demonstrate usage of new fields
- **AND** examples show creating instances with various combinations of new fields
- **AND** examples show updating new fields on existing instances

### Requirement: Code Formatting
All code changes SHALL be formatted using `go fmt`.

#### Scenario: File formatting after modification
- **GIVEN** code changes are made to the resource file
- **WHEN** all changes are complete
- **THEN** `go fmt` is executed on the file
- **AND** all Go code adheres to standard formatting rules
- **AND** no formatting warnings or errors exist
