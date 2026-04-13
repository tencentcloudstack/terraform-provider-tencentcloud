## ADDED Requirements

### Requirement: Resource supports reading and updating remark field
The system SHALL allow users to read and update the `remark` field for tencentcloud_tdmq_rabbitmq_vip_instance resource.

#### Scenario: Read remark field from cloud API
- **WHEN** user queries the tencentcloud_tdmq_rabbitmq_vip_instance resource
- **THEN** system reads the Remark field from DescribeRabbitMQVipInstance API response
- **AND** system sets the remark value in the resource state

#### Scenario: Update remark field via Modify API
- **WHEN** user updates the remark field in Terraform configuration
- **THEN** system calls ModifyRabbitMQVipInstance API with the new Remark value
- **AND** system reads the updated value back from cloud API

#### Scenario: Read remark from instance list API as fallback
- **WHEN** system reads from DescribeRabbitMQVipInstances API response
- **THEN** system reads the Remark field from RabbitMQVipInstance response
- **AND** system sets the remark value in the resource state

### Requirement: Resource supports reading and updating enable_deletion_protection field
The system SHALL allow users to read and update the `enable_deletion_protection` field for tencentcloud_tdmq_rabbitmq_vip_instance resource.

#### Scenario: Read enable_deletion_protection field from cloud API
- **WHEN** user queries the tencentcloud_tdmq_rabbitmq_vip_instance resource
- **THEN** system reads the EnableDeletionProtection field from DescribeRabbitMQVipInstance API response
- **AND** system sets the enable_deletion_protection value in the resource state
- **AND** system converts the boolean value to match Terraform type

#### Scenario: Update enable_deletion_protection field via Modify API
- **WHEN** user updates the enable_deletion_protection field in Terraform configuration
- **THEN** system calls ModifyRabbitMQVipInstance API with the new EnableDeletionProtection value
- **AND** system reads the updated value back from cloud API

#### Scenario: Read enable_deletion_protection from instance list API
- **WHEN** system reads from DescribeRabbitMQVipInstances API response
- **THEN** system reads the EnableDeletionProtection field from RabbitMQVipInstance response
- **AND** system sets the enable_deletion_protection value in the resource state

### Requirement: Resource supports reading and updating enable_risk_warning field
The system SHALL allow users to read and update the `enable_risk_warning` field for tencentcloud_tdmq_rabbitmq_vip_instance resource.

#### Scenario: Read enable_risk_warning field from instance detail API
- **WHEN** user queries the tencentcloud_tdmq_rabbitmq_vip_instance resource
- **THEN** system reads the EnableRiskWarning field from DescribeRabbitMQVipInstance API response (RabbitMQClusterInfo)
- **AND** system sets the enable_risk_warning value in the resource state
- **AND** system converts the boolean value to match Terraform type

#### Scenario: Update enable_risk_warning field via Modify API
- **WHEN** user updates the enable_risk_warning field in Terraform configuration
- **THEN** system calls ModifyRabbitMQVipInstance API with the new EnableRiskWarning value
- **AND** system reads the updated value back from cloud API

#### Scenario: enable_risk_warning not available in instance list API
- **WHEN** system reads from DescribeRabbitMQVipInstances API response
- **THEN** system does not read EnableRiskWarning field (not available in this API)
- **AND** system relies on DescribeRabbitMQVipInstance API for this field

### Requirement: New fields are optional and computed
The system SHALL define all new fields as Optional and Computed to maintain backward compatibility.

#### Scenario: Existing resources work without new fields
- **WHEN** user has existing tencentcloud_tdmq_rabbitmq_vip_instance resources without new fields
- **THEN** system does not require these fields to be present
- **AND** system reads these fields as computed values from cloud API

#### Scenario: New resources can omit new fields
- **WHEN** user creates a new tencentcloud_tdmq_rabbitmq_vip_instance resource
- **THEN** user can omit remark, enable_deletion_protection, and enable_risk_warning fields
- **AND** system uses default values from cloud API

### Requirement: Resource tags update logic handles empty tags correctly
The system SHALL correctly handle resource_tags update by setting RemoveAllTags flag when all tags are removed.

#### Scenario: Update resource_tags with new tags
- **WHEN** user updates resource_tags with non-empty tag list
- **THEN** system sets the Tags parameter in ModifyRabbitMQVipInstance API request
- **AND** system does not set RemoveAllTags parameter

#### Scenario: Remove all resource_tags
- **WHEN** user updates resource_tags with empty tag list
- **THEN** system sets RemoveAllTags parameter to true in ModifyRabbitMQVipInstance API request
- **AND** system does not set Tags parameter

#### Scenario: resource_tags unchanged
- **WHEN** user does not modify resource_tags field
- **THEN** system does not call ModifyRabbitMQVipInstance API for tags
- **AND** system preserves existing tags
