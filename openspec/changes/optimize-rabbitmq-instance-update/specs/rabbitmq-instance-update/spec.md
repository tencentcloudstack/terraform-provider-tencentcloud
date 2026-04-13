## ADDED Requirements

### Requirement: Resource supports remark field CRUD
The system SHALL support the `remark` field in the `tencentcloud_tdmq_rabbitmq_vip_instance` resource.

#### Scenario: Create instance with remark
- **WHEN** user creates a RabbitMQ VIP instance with `remark` field set to "Production environment"
- **THEN** the instance is created and the remark is set to "Production environment"
- **AND** the remark value is persisted in the Terraform state

#### Scenario: Read instance returns remark
- **WHEN** user reads an existing RabbitMQ VIP instance
- **THEN** the `remark` field is returned from the cloud API
- **AND** the remark value is set in the Terraform state

#### Scenario: Update instance remark
- **WHEN** user updates the `remark` field of an existing RabbitMQ VIP instance from "Development" to "Staging"
- **THEN** the ModifyRabbitMQVipInstance API is called with the new remark value
- **AND** the remark is successfully updated on the instance
- **AND** the updated remark is reflected in the Terraform state

### Requirement: Resource supports enable_deletion_protection field CRUD
The system SHALL support the `enable_deletion_protection` field in the `tencentcloud_tdmq_rabbitmq_vip_instance` resource.

#### Scenario: Create instance with deletion protection enabled
- **WHEN** user creates a RabbitMQ VIP instance with `enable_deletion_protection` set to `true`
- **THEN** the instance is created with deletion protection enabled
- **AND** the enable_deletion_protection value is persisted in the Terraform state

#### Scenario: Read instance returns deletion protection status
- **WHEN** user reads an existing RabbitMQ VIP instance
- **THEN** the `enable_deletion_protection` field is returned from the cloud API
- **AND** the enable_deletion_protection value is set in the Terraform state

#### Scenario: Enable deletion protection on existing instance
- **WHEN** user updates `enable_deletion_protection` from `false` to `true` on an existing instance
- **THEN** the ModifyRabbitMQVipInstance API is called with EnableDeletionProtection set to true
- **AND** deletion protection is enabled on the instance
- **AND** the updated enable_deletion_protection is reflected in the Terraform state

#### Scenario: Disable deletion protection on existing instance
- **WHEN** user updates `enable_deletion_protection` from `true` to `false` on an existing instance
- **THEN** the ModifyRabbitMQVipInstance API is called with EnableDeletionProtection set to false
- **AND** deletion protection is disabled on the instance
- **AND** the updated enable_deletion_protection is reflected in the Terraform state

### Requirement: Resource supports enable_risk_warning field CRUD
The system SHALL support the `enable_risk_warning` field in the `tencentcloud_tdmq_rabbitmq_vip_instance` resource.

#### Scenario: Create instance with risk warning enabled
- **WHEN** user creates a RabbitMQ VIP instance with `enable_risk_warning` set to `true`
- **THEN** the instance is created with risk warning enabled
- **AND** the enable_risk_warning value is persisted in the Terraform state

#### Scenario: Read instance returns risk warning status
- **WHEN** user reads an existing RabbitMQ VIP instance
- **THEN** the `enable_risk_warning` field is returned from the cloud API
- **AND** the enable_risk_warning value is set in the Terraform state

#### Scenario: Enable risk warning on existing instance
- **WHEN** user updates `enable_risk_warning` from `false` to `true` on an existing instance
- **THEN** the ModifyRabbitMQVipInstance API is called with EnableRiskWarning set to true
- **AND** risk warning is enabled on the instance
- **AND** the updated enable_risk_warning is reflected in the Terraform state

#### Scenario: Disable risk warning on existing instance
- **WHEN** user updates `enable_risk_warning` from `true` to `false` on an existing instance
- **THEN** the ModifyRabbitMQVipInstance API is called with EnableRiskWarning set to false
- **AND** risk warning is disabled on the instance
- **AND** the updated enable_risk_warning is reflected in the Terraform state

### Requirement: Update operation supports partial field updates
The system SHALL support updating any combination of the new fields (remark, enable_deletion_protection, enable_risk_warning) independently.

#### Scenario: Update only remark
- **WHEN** user updates only the `remark` field without changing enable_deletion_protection or enable_risk_warning
- **THEN** the ModifyRabbitMQVipInstance API is called with only the remark field
- **AND** other fields remain unchanged on the instance

#### Scenario: Update only enable_deletion_protection
- **WHEN** user updates only the `enable_deletion_protection` field without changing remark or enable_risk_warning
- **THEN** the ModifyRabbitMQVipInstance API is called with only the EnableDeletionProtection field
- **AND** other fields remain unchanged on the instance

#### Scenario: Update only enable_risk_warning
- **WHEN** user updates only the `enable_risk_warning` field without changing remark or enable_deletion_protection
- **THEN** the ModifyRabbitMQVipInstance API is called with only the EnableRiskWarning field
- **AND** other fields remain unchanged on the instance

#### Scenario: Update multiple fields simultaneously
- **WHEN** user updates multiple fields (e.g., remark and enable_deletion_protection) in a single update
- **THEN** the ModifyRabbitMQVipInstance API is called with all changed fields
- **AND** all specified fields are updated on the instance

### Requirement: Backward compatibility is maintained
The system SHALL maintain backward compatibility with existing Terraform configurations.

#### Scenario: Existing configuration without new fields
- **WHEN** user applies an existing configuration that does not include remark, enable_deletion_protection, or enable_risk_warning fields
- **THEN** the resource is created/updated successfully
- **AND** the new fields use cloud API default values
- **AND** existing behavior is not affected

#### Scenario: Import existing instance without new fields
- **WHEN** user imports an existing RabbitMQ VIP instance
- **THEN** the import operation succeeds
- **AND** the new fields are populated with values from the cloud API
- **AND** the imported state is consistent with the cloud instance
