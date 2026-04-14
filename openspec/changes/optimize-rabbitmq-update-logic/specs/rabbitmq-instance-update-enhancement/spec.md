## ADDED Requirements

### Requirement: Support remark parameter update
The system SHALL allow users to update the `remark` parameter of a RabbitMQ VIP instance through the ModifyRabbitMQVipInstance API. The `remark` parameter is an optional string field used to provide descriptive notes about the instance.

#### Scenario: Update remark parameter successfully
- **WHEN** user sets or changes the `remark` parameter in the Terraform configuration
- **THEN** the system SHALL call ModifyRabbitMQVipInstance API with the Remark field
- **THEN** the system SHALL successfully update the instance's remark in the cloud
- **THEN** the system SHALL read back the updated remark value and store it in the Terraform state

#### Scenario: Remove remark parameter
- **WHEN** user removes the `remark` parameter from the Terraform configuration
- **THEN** the system SHALL call ModifyRabbitMQVipInstance API with Remark set to empty string
- **THEN** the system SHALL successfully clear the instance's remark in the cloud
- **THEN** the system SHALL update the Terraform state with empty remark value

### Requirement: Support enable_deletion_protection parameter update
The system SHALL allow users to update the `enable_deletion_protection` parameter of a RabbitMQ VIP instance through the ModifyRabbitMQVipInstance API. The `enable_deletion_protection` parameter is an optional boolean field that controls whether deletion protection is enabled for the instance.

#### Scenario: Enable deletion protection
- **WHEN** user sets `enable_deletion_protection` to `true` in the Terraform configuration
- **THEN** the system SHALL call ModifyRabbitMQVipInstance API with EnableDeletionProtection set to true
- **THEN** the system SHALL successfully enable deletion protection for the instance
- **THEN** the system SHALL read back the enable_deletion_protection value as true in the Terraform state

#### Scenario: Disable deletion protection
- **WHEN** user sets `enable_deletion_protection` to `false` in the Terraform configuration
- **THEN** the system SHALL call ModifyRabbitMQVipInstance API with EnableDeletionProtection set to false
- **THEN** the system SHALL successfully disable deletion protection for the instance
- **THEN** the system SHALL read back the enable_deletion_protection value as false in the Terraform state

#### Scenario: Do not modify deletion protection when not specified
- **WHEN** the `enable_deletion_protection` parameter is not set or changed in the Terraform configuration
- **THEN** the system SHALL NOT include the EnableDeletionProtection field in the ModifyRabbitMQVipInstance API call
- **THEN** the system SHALL preserve the existing deletion protection setting in the cloud

### Requirement: Support enable_risk_warning parameter update
The system SHALL allow users to update the `enable_risk_warning` parameter of a RabbitMQ VIP instance through the ModifyRabbitMQVipInstance API. The `enable_risk_warning` parameter is an optional boolean field that controls whether cluster risk warning is enabled for the instance.

#### Scenario: Enable risk warning
- **WHEN** user sets `enable_risk_warning` to `true` in the Terraform configuration
- **THEN** the system SHALL call ModifyRabbitMQVipInstance API with EnableRiskWarning set to true
- **THEN** the system SHALL successfully enable cluster risk warning for the instance
- **THEN** the system SHALL read back the enable_risk_warning value as true in the Terraform state

#### Scenario: Disable risk warning
- **WHEN** user sets `enable_risk_warning` to `false` in the Terraform configuration
- **THEN** the system SHALL call ModifyRabbitMQVipInstance API with EnableRiskWarning set to false
- **THEN** the system SHALL successfully disable cluster risk warning for the instance
- **THEN** the system SHALL read back the enable_risk_warning value as false in the Terraform state

#### Scenario: Do not modify risk warning when not specified
- **WHEN** the `enable_risk_warning` parameter is not set or changed in the Terraform configuration
- **THEN** the system SHALL NOT include the EnableRiskWarning field in the ModifyRabbitMQVipInstance API call
- **THEN** the system SHALL preserve the existing risk warning setting in the cloud

### Requirement: Maintain backward compatibility with existing parameters
The system SHALL maintain full backward compatibility with existing update capabilities, including `cluster_name` and `resource_tags` parameters. The existing immutable parameters validation SHALL remain unchanged.

#### Scenario: Update cluster_name with new parameters
- **WHEN** user changes `cluster_name` and also sets `remark`, `enable_deletion_protection`, and `enable_risk_warning` in a single Terraform apply
- **THEN** the system SHALL call ModifyRabbitMQVipInstance API with all four parameters (ClusterName, Remark, EnableDeletionProtection, EnableRiskWarning)
- **THEN** the system SHALL successfully update all four parameters in the cloud
- **THEN** the system SHALL read back all updated values and update the Terraform state

#### Scenario: Attempt to update immutable parameter
- **WHEN** user attempts to modify an immutable parameter (e.g., `node_spec`, `storage_size`, `band_width`)
- **THEN** the system SHALL return a clear error message indicating that the parameter cannot be changed
- **THEN** the system SHALL NOT call any cloud API
- **THEN** the Terraform state SHALL remain unchanged

#### Scenario: Update only resource_tags with new parameters
- **WHEN** user changes `resource_tags` and also sets `remark` in a single Terraform apply
- **THEN** the system SHALL call ModifyRabbitMQVipInstance API with both Tags and Remark fields
- **THEN** the system SHALL successfully update both parameters in the cloud
- **THEN** the system SHALL read back all updated values and update the Terraform state

### Requirement: Read new parameters from cloud API
The system SHALL read the `remark`, `enable_deletion_protection`, and `enable_risk_warning` parameters from the DescribeRabbitMQVipInstance API response and store them in the Terraform state during the Read operation.

#### Scenario: Read instance with all new parameters
- **WHEN** the system performs a read operation on a RabbitMQ VIP instance
- **THEN** the system SHALL call DescribeRabbitMQVipInstance API to retrieve the instance details
- **THEN** the system SHALL extract the `remark` value from the API response and set it in the Terraform state
- **THEN** the system SHALL extract the `enable_deletion_protection` value from the API response and set it in the Terraform state
- **THEN** the system SHALL extract the `enable_risk_warning` value from the API response and set it in the Terraform state

#### Scenario: Read instance without optional parameters
- **WHEN** the system performs a read operation on a RabbitMQ VIP instance that does not have optional parameters set
- **THEN** the system SHALL handle missing `remark`, `enable_deletion_protection`, and `enable_risk_warning` values gracefully
- **THEN** the system SHALL set these parameters to their zero values (empty string for remark, false for boolean values) in the Terraform state
