# CAM Policy Resource Specification

## ADDED Requirements

### Requirement: Tags Management
The `tencentcloud_cam_policy` resource SHALL support tags management to enable resource organization and access control.

#### Scenario: Create policy with tags
- **WHEN** user creates a CAM policy with tags specified
- **THEN** the tags are applied to the policy during creation via CreatePolicy API

#### Scenario: Read policy tags
- **WHEN** user imports or refreshes a CAM policy
- **THEN** the current tags are retrieved from TencentCloud and stored in state

#### Scenario: Update policy tags
- **WHEN** user modifies tags on an existing CAM policy
- **THEN** the tags are updated using ModifyTags API without recreating the policy

#### Scenario: Remove all tags
- **WHEN** user removes all tags from a CAM policy configuration
- **THEN** all tags are deleted from the policy resource

### Requirement: Schema Definition
The tags field SHALL be defined as an optional map of strings in the resource schema.

#### Scenario: Tags field specification
- **WHEN** defining the resource schema
- **THEN** tags field uses TypeMap, is Optional, and includes description "Instance tag"

### Requirement: API Integration
The resource SHALL integrate with TencentCloud APIs for tags operations.

#### Scenario: Create API integration
- **WHEN** creating a policy with tags
- **THEN** tags are converted to SDK Tag array format and included in CreatePolicy request

#### Scenario: Read API integration
- **WHEN** reading policy state
- **THEN** tags are retrieved using TagService.DescribeResourceTags with resource name format `qcs::cam:{region}:uin/:policy/{policyId}`

#### Scenario: Update API integration
- **WHEN** updating policy tags
- **THEN** TagService.ModifyTags is called with calculated replace and delete tag sets

### Requirement: Backward Compatibility
Tags support SHALL maintain full backward compatibility with existing configurations.

#### Scenario: Existing policies without tags
- **WHEN** managing policies created before tags support
- **THEN** resource operations work normally with tags field empty or null
