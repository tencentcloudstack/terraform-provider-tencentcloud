## ADDED Requirements

### Requirement: Update instance remark
The system SHALL allow users to update the `remark` parameter of a RabbitMQ VIP instance. The `remark` parameter is an optional string field that provides instance description information.

#### Scenario: Update remark on existing instance
- **WHEN** user updates the `remark` parameter in the Terraform configuration
- **THEN** the system calls `ModifyRabbitMQVipInstance` API with the new `Remark` value
- **AND** the system reads the updated `Remark` value from the API response after update

#### Scenario: Update instance without changing remark
- **WHEN** user updates other parameters of the instance without changing `remark`
- **THEN** the system does not include the `Remark` field in the API request
- **AND** the system preserves the existing `remark` value in the resource state

#### Scenario: Read instance with remark set
- **WHEN** system reads a RabbitMQ VIP instance that has a remark
- **THEN** the system retrieves the `Remark` value from the `RabbitMQClusterInfo` object in the API response
- **AND** the system sets the `remark` parameter in the resource state with the retrieved value

#### Scenario: Read instance without remark
- **WHEN** system reads a RabbitMQ VIP instance that does not have a remark
- **THEN** the system retrieves a nil or empty `Remark` value from the `RabbitMQClusterInfo` object
- **AND** the system sets the `remark` parameter in the resource state to an empty string or nil

### Requirement: Update deletion protection flag
The system SHALL allow users to update the `enable_deletion_protection` parameter of a RabbitMQ VIP instance. The `enable_deletion_protection` parameter is an optional boolean field that controls whether deletion protection is enabled for the instance.

#### Scenario: Enable deletion protection
- **WHEN** user sets `enable_deletion_protection` to `true` in the Terraform configuration
- **THEN** the system calls `ModifyRabbitMQVipInstance` API with `EnableDeletionProtection` set to `true`
- **AND** the system reads the updated `EnableDeletionProtection` value from the API response after update

#### Scenario: Disable deletion protection
- **WHEN** user sets `enable_deletion_protection` to `false` in the Terraform configuration
- **THEN** the system calls `ModifyRabbitMQVipInstance` API with `EnableDeletionProtection` set to `false`
- **AND** the system reads the updated `EnableDeletionProtection` value from the API response after update

#### Scenario: Update instance without changing deletion protection
- **WHEN** user updates other parameters of the instance without changing `enable_deletion_protection`
- **THEN** the system does not include the `EnableDeletionProtection` field in the API request
- **AND** the system preserves the existing `enable_deletion_protection` value in the resource state

#### Scenario: Read instance with deletion protection
- **WHEN** system reads a RabbitMQ VIP instance that has deletion protection enabled
- **THEN** the system retrieves the `EnableDeletionProtection` value from the `RabbitMQClusterInfo` object in the API response
- **AND** the system sets the `enable_deletion_protection` parameter in the resource state to `true`

#### Scenario: Read instance without deletion protection
- **WHEN** system reads a RabbitMQ VIP instance that does not have deletion protection enabled
- **THEN** the system retrieves the `EnableDeletionProtection` value from the `RabbitMQClusterInfo` object
- **AND** the system sets the `enable_deletion_protection` parameter in the resource state to `false`

### Requirement: Update risk warning flag
The system SHALL allow users to update the `enable_risk_warning` parameter of a RabbitMQ VIP instance. The `enable_risk_warning` parameter is an optional boolean field that controls whether cluster risk warning is enabled for the instance.

#### Scenario: Enable risk warning
- **WHEN** user sets `enable_risk_warning` to `true` in the Terraform configuration
- **THEN** the system calls `ModifyRabbitMQVipInstance` API with `EnableRiskWarning` set to `true`
- **AND** the system reads the updated `EnableRiskWarning` value from the API response after update

#### Scenario: Disable risk warning
- **WHEN** user sets `enable_risk_warning` to `false` in the Terraform configuration
- **THEN** the system calls `ModifyRabbitMQVipInstance` API with `EnableRiskWarning` set to `false`
- **AND** the system reads the updated `EnableRiskWarning` value from the API response after update

#### Scenario: Update instance without changing risk warning
- **WHEN** user updates other parameters of the instance without changing `enable_risk_warning`
- **THEN** the system does not include the `EnableRiskWarning` field in the API request
- **AND** the system preserves the existing `enable_risk_warning` value in the resource state

#### Scenario: Read instance with risk warning enabled
- **WHEN** system reads a RabbitMQ VIP instance that has risk warning enabled
- **THEN** the system retrieves the `EnableRiskWarning` value from the `RabbitMQClusterInfo` object in the API response
- **AND** the system sets the `enable_risk_warning` parameter in the resource state to `true`

#### Scenario: Read instance without risk warning
- **WHEN** system reads a RabbitMQ VIP instance that does not have risk warning enabled
- **THEN** the system retrieves the `EnableRiskWarning` value from the `RabbitMQClusterInfo` object
- **AND** the system sets the `enable_risk_warning` parameter in the resource state to `false`

### Requirement: Backward compatibility
The system SHALL ensure that existing Terraform configurations without the new parameters continue to work without any modifications or errors.

#### Scenario: Existing configuration without new parameters
- **WHEN** user applies an existing Terraform configuration that does not include `remark`, `enable_deletion_protection`, or `enable_risk_warning` parameters
- **THEN** the system successfully updates the instance with the existing parameters
- **AND** the system does not require the user to add the new parameters
- **AND** the system reads the values for the new parameters from the API response and sets them in the resource state

#### Scenario: Import existing instance
- **WHEN** user imports an existing RabbitMQ VIP instance into Terraform state
- **THEN** the system reads all parameters including `remark`, `enable_deletion_protection`, and `enable_risk_warning` from the API
- **AND** the system sets all parameters in the resource state with the values from the API
- **AND** the imported state matches the actual cloud resource configuration
