# Spec: RabbitMQ VIP Instance Update Enhancement

This specification defines the enhancement to the RabbitMQ VIP instance update capability, allowing modification of configuration parameters that were previously immutable.

## ADDED Requirements

### Requirement: User can modify deletion protection
The system SHALL allow users to modify the `enable_deletion_protection` parameter of an existing RabbitMQ VIP instance through Terraform.

#### Scenario: Enable deletion protection
- **WHEN** user sets `enable_deletion_protection` to `true` in the Terraform configuration
- **THEN** the system SHALL call ModifyRabbitMQVipInstance API with `EnableDeletionProtection=true`
- **AND** the deletion protection SHALL be enabled on the instance

#### Scenario: Disable deletion protection
- **WHEN** user sets `enable_deletion_protection` to `false` in the Terraform configuration
- **THEN** the system SHALL call ModifyRabbitMQVipInstance API with `EnableDeletionProtection=false`
- **AND** the deletion protection SHALL be disabled on the instance

### Requirement: User can modify instance remark
The system SHALL allow users to modify the `remark` (description) parameter of an existing RabbitMQ VIP instance through Terraform.

#### Scenario: Set instance remark
- **WHEN** user sets `remark` to a value in the Terraform configuration
- **THEN** the system SHALL call ModifyRabbitMQVipInstance API with the `Remark` parameter
- **AND** the instance remark SHALL be updated to the specified value

#### Scenario: Clear instance remark
- **WHEN** user removes the `remark` parameter or sets it to empty
- **THEN** the system SHALL call ModifyRabbitMQVipInstance API with `Remark` set to empty string or null
- **AND** the instance remark SHALL be cleared

### Requirement: User can modify risk warning setting
The system SHALL allow users to modify the `enable_risk_warning` parameter of an existing RabbitMQ VIP instance through Terraform.

#### Scenario: Enable risk warning
- **WHEN** user sets `enable_risk_warning` to `true` in the Terraform configuration
- **THEN** the system SHALL call ModifyRabbitMQVipInstance API with `EnableRiskWarning=true`
- **AND** the risk warning SHALL be enabled on the instance

#### Scenario: Disable risk warning
- **WHEN** user sets `enable_risk_warning` to `false` in the Terraform configuration
- **THEN** the system SHALL call ModifyRabbitMQVipInstance API with `EnableRiskWarning=false`
- **AND** the risk warning SHALL be disabled on the instance

### Requirement: Error handling for unsupported modifications
The system SHALL provide clear error messages when users attempt to modify parameters that are not supported by the ModifyRabbitMQVipInstance API.

#### Scenario: Modify immutable infrastructure parameter
- **WHEN** user attempts to modify `zone_ids`, `vpc_id`, `subnet_id`, `node_spec`, `node_num`, `storage_size`, `cluster_version`, `auto_renew_flag`, `time_span`, or `pay_mode`
- **THEN** the system SHALL return an error message indicating the parameter cannot be changed
- **AND** the error message SHALL be clear about which parameter is not modifiable

### Requirement: Partial updates supported
The system SHALL support updating a subset of the supported parameters without requiring all parameters to be specified.

#### Scenario: Update only remark
- **WHEN** user modifies only the `remark` parameter
- **THEN** the system SHALL call ModifyRabbitMQVipInstance API with only `Remark` parameter
- **AND** other parameters SHALL remain unchanged

#### Scenario: Update multiple parameters
- **WHEN** user modifies multiple supported parameters (e.g., `remark` and `enable_deletion_protection`)
- **THEN** the system SHALL call ModifyRabbitMQVipInstance API with all changed parameters
- **AND** the modification SHALL be applied in a single API call

### Requirement: Backward compatibility maintained
The system SHALL maintain backward compatibility with existing Terraform configurations and state.

#### Scenario: Existing configuration without new parameters
- **WHEN** user has an existing RabbitMQ VIP instance resource without `remark` or `enable_risk_warning` fields
- **THEN** the existing configuration SHALL continue to work without modification
- **AND** the Read function SHALL populate these fields with values from the API response

#### Scenario: Existing state compatibility
- **WHEN** user upgrades to the new provider version
- **THEN** the existing state SHALL remain compatible
- **AND** no state migration SHALL be required
