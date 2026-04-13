# Spec: TDMQ RabbitMQ VIP Instance

Delta spec for optimizing RabbitMQ instance update logic.

## MODIFIED Requirements

### Requirement: Public Access Field Immutability
The `enable_public_access` and `band_width` fields SHALL be immutable after instance creation.

#### Scenario: User attempts to enable public access on existing instance
- **GIVEN** an existing RabbitMQ VIP instance with `enable_public_access = false`
- **WHEN** user changes `enable_public_access` to `true` in configuration
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` fails with error message "argument `enable_public_access` cannot be changed"
- **AND** instance is not modified

#### Scenario: User attempts to modify bandwidth on existing instance
- **GIVEN** an existing RabbitMQ VIP instance with `band_width = 100`
- **WHEN** user changes `band_width` to `200` in configuration
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` fails with error message "argument `band_width` cannot be changed"
- **AND** instance is not modified

#### Scenario: Recreating instance with different public access settings
- **GIVEN** an existing instance with public access configuration
- **WHEN** user changes immutable fields (`enable_public_access` or `band_width`)
- **THEN** Terraform requires manual resource destruction and recreation
- **AND** users must use `terraform taint` or manual delete/create workflow

### Requirement: Other Immutable Fields
The `zone_ids`, `vpc_id`, `subnet_id`, `node_spec`, `node_num`, `storage_size`, `enable_create_default_ha_mirror_queue`, `auto_renew_flag`, `time_span`, `pay_mode`, and `cluster_version` fields SHALL be immutable after instance creation.

#### Scenario: User attempts to modify immutable fields
- **GIVEN** an existing RabbitMQ VIP instance
- **WHEN** user changes any of `zone_ids`, `vpc_id`, `subnet_id`, `node_spec`, `node_num`, `storage_size`, `enable_create_default_ha_mirror_queue`, `auto_renew_flag`, `time_span`, `pay_mode`, or `cluster_version` in configuration
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` fails with error message "argument `<field_name>` cannot be changed"
- **AND** instance is not modified

## ADDED Requirements

### Requirement: Remark Field Support
The `tencentcloud_tdmq_rabbitmq_vip_instance` resource SHALL support a `remark` field for managing instance remarks.

#### Scenario: User creates instance with remark
- **GIVEN** a user defines a RabbitMQ VIP instance configuration
- **WHEN** user sets `remark = "Example remark for RabbitMQ instance"` in configuration
- **THEN** instance is created with the specified remark
- **AND** remark is visible in Terraform state after creation
- **AND** remark is applied to the cloud resource

#### Scenario: User creates instance without remark
- **GIVEN** a user defines a RabbitMQ VIP instance configuration
- **WHEN** user does not include `remark` field
- **THEN** instance is created successfully without a remark
- **AND** `remark` field in state is empty or reflects API default value

#### Scenario: User reads instance remark into state
- **GIVEN** a RabbitMQ VIP instance exists with a remark in Tencent Cloud
- **WHEN** Terraform performs a refresh operation
- **THEN** `remark` field in state is populated with the current remark from cloud
- **AND** nil remark values are handled gracefully

#### Scenario: User updates instance remark
- **GIVEN** an existing RabbitMQ VIP instance with `remark = "Old remark"`
- **WHEN** user changes `remark` to `"New remark"` in configuration
- **THEN** `terraform plan` shows the remark change
- **AND** `terraform apply` successfully updates the remark on the instance
- **AND** Terraform state reflects the updated remark

### Requirement: Enable Deletion Protection Field Support
The `tencentcloud_tdmq_rabbitmq_vip_instance` resource SHALL support an `enable_deletion_protection` field for managing deletion protection status.

#### Scenario: User creates instance with deletion protection enabled
- **GIVEN** a user defines a RabbitMQ VIP instance configuration
- **WHEN** user sets `enable_deletion_protection = true` in configuration
- **THEN** instance is created with deletion protection enabled
- **AND** `enable_deletion_protection` field shows `true` in Terraform state after creation
- **AND** deletion protection is applied to the cloud resource

#### Scenario: User creates instance without deletion protection (default)
- **GIVEN** a user defines a RabbitMQ VIP instance configuration
- **WHEN** user does not include `enable_deletion_protection` field
- **THEN** instance is created with deletion protection disabled (API default: false)
- **AND** `enable_deletion_protection` field in state shows `false`

#### Scenario: User reads deletion protection status into state
- **GIVEN** a RabbitMQ VIP instance exists in Tencent Cloud
- **WHEN** Terraform performs a refresh operation
- **THEN** `enable_deletion_protection` field in state is populated with the current status from cloud
- **AND** nil values are handled gracefully

#### Scenario: User enables deletion protection on existing instance
- **GIVEN** an existing RabbitMQ VIP instance with `enable_deletion_protection = false`
- **WHEN** user changes `enable_deletion_protection` to `true` in configuration
- **THEN** `terraform plan` shows the deletion protection change
- **AND** `terraform apply` successfully enables deletion protection on the instance
- **AND** Terraform state reflects the updated status

#### Scenario: User disables deletion protection on existing instance
- **GIVEN** an existing RabbitMQ VIP instance with `enable_deletion_protection = true`
- **WHEN** user changes `enable_deletion_protection` to `false` in configuration
- **THEN** `terraform plan` shows the deletion protection change
- **AND** `terraform apply` successfully disables deletion protection on the instance
- **AND** Terraform state reflects the updated status

### Requirement: Enable Risk Warning Field Support
The `tencentcloud_tdmq_rabbitmq_vip_instance` resource SHALL support an `enable_risk_warning` field for managing cluster risk warning status.

#### Scenario: User creates instance with risk warning enabled
- **GIVEN** a user defines a RabbitMQ VIP instance configuration
- **WHEN** user sets `enable_risk_warning = true` in configuration
- **THEN** instance is created with risk warning enabled
- **AND** `enable_risk_warning` field shows `true` in Terraform state after creation
- **AND** risk warning is applied to the cloud resource

#### Scenario: User creates instance without risk warning (default)
- **GIVEN** a user defines a RabbitMQ VIP instance configuration
- **WHEN** user does not include `enable_risk_warning` field
- **THEN** instance is created with risk warning disabled (API default: false)
- **AND** `enable_risk_warning` field in state shows `false`

#### Scenario: User updates risk warning on existing instance
- **GIVEN** an existing RabbitMQ VIP instance with `enable_risk_warning = false`
- **WHEN** user changes `enable_risk_warning` to `true` in configuration
- **THEN** `terraform plan` shows the risk warning change
- **AND** `terraform apply` successfully enables risk warning on the instance
- **AND** Terraform state reflects the updated status

#### Scenario: User disables risk warning on existing instance
- **GIVEN** an existing RabbitMQ VIP instance with `enable_risk_warning = true`
- **WHEN** user changes `enable_risk_warning` to `false` in configuration
- **THEN** `terraform plan` shows the risk warning change
- **AND** `terraform apply` successfully disables risk warning on the instance
- **AND** Terraform state reflects the updated status

### Requirement: Create API Integration for New Fields
The resource SHALL correctly pass new fields to `CreateRabbitMQVipInstance` API.

#### Scenario: Send remark to Create API
- **GIVEN** a user creates an instance with `remark = "Example remark"`
- **WHEN** Create operation calls `CreateRabbitMQVipInstance`
- **THEN** request includes `Remark: helper.String("Example remark")`
- **AND** API accepts the parameter

#### Scenario: Send enable_deletion_protection to Create API
- **GIVEN** a user creates an instance with `enable_deletion_protection = true`
- **WHEN** Create operation calls `CreateRabbitMQVipInstance`
- **THEN** request includes `EnableDeletionProtection: helper.Bool(true)`
- **AND** API enables deletion protection

#### Scenario: Omit new fields when not specified
- **GIVEN** a user creates an instance without `remark`, `enable_deletion_protection`, or `enable_risk_warning`
- **WHEN** Create operation calls `CreateRabbitMQVipInstance`
- **THEN** request does not include these fields
- **AND** API applies default values

### Requirement: Read API Integration for New Fields
The resource SHALL correctly read new fields from `DescribeRabbitMQVipInstances` API.

#### Scenario: Read remark from API response
- **GIVEN** a RabbitMQ VIP instance exists with a remark
- **WHEN** Read operation calls `DescribeRabbitMQVipInstances`
- **THEN** remark value is extracted from `response.RabbitMQVipInstance.Remark`
- **AND** value is set in Terraform state as `remark` field
- **AND** nil values are handled gracefully (not set in state)

#### Scenario: Read enable_deletion_protection from API response
- **GIVEN** a RabbitMQ VIP instance exists with deletion protection enabled
- **WHEN** Read operation calls `DescribeRabbitMQVipInstances`
- **THEN** deletion protection status is extracted from `response.RabbitMQVipInstance.EnableDeletionProtection`
- **AND** boolean value is set in Terraform state as `enable_deletion_protection`
- **AND** nil values are handled gracefully (not set in state)

#### Scenario: Handle missing enable_risk_warning in API response
- **GIVEN** a RabbitMQ VIP instance exists in Tencent Cloud
- **WHEN** Read operation calls `DescribeRabbitMQVipInstances`
- **THEN** if `EnableRiskWarning` field does not exist in API response, `enable_risk_warning` field is not set in state
- **AND** no error is raised during state refresh

### Requirement: Update API Integration for New Fields
The resource SHALL correctly pass new fields to `ModifyRabbitMQVipInstance` API.

#### Scenario: Update remark via Modify API
- **GIVEN** an existing RabbitMQ VIP instance
- **WHEN** user changes `remark` in configuration
- **THEN** Update operation calls `ModifyRabbitMQVipInstance` with `Remark: helper.String(new_value)`
- **AND** API successfully updates the remark
- **AND** no other instance properties are affected

#### Scenario: Update enable_deletion_protection via Modify API
- **GIVEN** an existing RabbitMQ VIP instance
- **WHEN** user changes `enable_deletion_protection` in configuration
- **THEN** Update operation calls `ModifyRabbitMQVipInstance` with `EnableDeletionProtection: helper.Bool(new_value)`
- **AND** API successfully updates the deletion protection status
- **AND** no other instance properties are affected

#### Scenario: Update enable_risk_warning via Modify API
- **GIVEN** an existing RabbitMQ VIP instance
- **WHEN** user changes `enable_risk_warning` in configuration
- **THEN** Update operation calls `ModifyRabbitMQVipInstance` with `EnableRiskWarning: helper.Bool(new_value)`
- **AND** API successfully updates the risk warning status
- **AND** no other instance properties are affected

### Requirement: Schema Definition for New Fields
The resource schema SHALL define new fields with appropriate attributes.

#### Scenario: remark schema properties
- **GIVEN** resource schema definition
- **WHEN** examining the `remark` field
- **THEN** field type is `schema.TypeString`
- **AND** field is marked as `Optional: true` and `Computed: true`
- **AND** field has description: "Remark of the RabbitMQ instance."

#### Scenario: enable_deletion_protection schema properties
- **GIVEN** resource schema definition
- **WHEN** examining the `enable_deletion_protection` field
- **THEN** field type is `schema.TypeBool`
- **AND** field is marked as `Optional: true` and `Computed: true`
- **AND** field has description: "Whether to enable deletion protection. Default is false."

#### Scenario: enable_risk_warning schema properties
- **GIVEN** resource schema definition
- **WHEN** examining the `enable_risk_warning` field
- **THEN** field type is `schema.TypeBool`
- **AND** field is marked as `Optional: true` and `Computed: true`
- **AND** field has description: "Whether to enable cluster risk warning. Default is false."

### Requirement: Documentation Completeness for New Fields
The resource documentation SHALL clearly describe new fields.

#### Scenario: Field documentation
- **GIVEN** resource documentation file `tdmq_rabbitmq_vip_instance.html.markdown`
- **WHEN** reviewing the arguments reference section
- **THEN** `remark` is documented with type and description
- **AND** `enable_deletion_protection` is documented with type, default value, and description
- **AND** `enable_risk_warning` is documented with type, default value, and description
- **AND** all three fields are marked as mutable (can be changed after creation)

#### Scenario: Usage example with new fields
- **GIVEN** resource documentation file
- **WHEN** reviewing the example usage section
- **THEN** an example shows creating an instance with `remark` field
- **AND** example demonstrates the use of `enable_deletion_protection` field
- **AND** example demonstrates updating these fields

### Requirement: Backward Compatibility with New Fields
The new fields SHALL be backward compatible with existing resources.

#### Scenario: Existing instance without new fields
- **GIVEN** a RabbitMQ VIP instance managed by Terraform before this feature
- **WHEN** provider is upgraded to include new field support
- **THEN** `terraform plan` shows no changes for resources without these fields
- **AND** existing resources continue to function normally
- **AND** state refresh correctly populates new fields from API response

#### Scenario: First-time new field addition to existing resource
- **GIVEN** an existing instance without `remark`, `enable_deletion_protection`, or `enable_risk_warning` in configuration
- **WHEN** user adds any of these fields for the first time
- **THEN** Terraform treats this as an update (not recreation)
- **AND** only `ModifyRabbitMQVipInstance` API is called
- **AND** no other instance properties are affected

### Requirement: Error Handling for New Fields
The resource SHALL handle new field errors gracefully.

#### Scenario: API error during creation with new fields
- **GIVEN** a user creates an instance with invalid `remark` (e.g., too long)
- **WHEN** `CreateRabbitMQVipInstance` API returns a validation error
- **THEN** error is propagated to user with context
- **AND** instance creation is not completed
- **AND** no partial state is written

#### Scenario: API error during update with new fields
- **GIVEN** a user updates `remark` on an existing instance
- **WHEN** `ModifyRabbitMQVipInstance` API fails
- **THEN** error is logged and returned to user
- **AND** Terraform state remains unchanged (previous values)
- **AND** subsequent `terraform apply` attempts the update again

#### Scenario: Nil value handling during read for new fields
- **GIVEN** API response contains nil values for new fields
- **WHEN** Read operation processes the response
- **THEN** nil values are handled without panicking
- **AND** fields are not set in state (or set to default values)
- **AND** no error is returned for nil values
