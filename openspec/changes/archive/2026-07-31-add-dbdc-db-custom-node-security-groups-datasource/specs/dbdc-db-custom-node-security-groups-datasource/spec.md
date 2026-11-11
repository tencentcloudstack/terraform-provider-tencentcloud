## ADDED Requirements

### Requirement: Data source reads DBCustom node security groups by node ID
The system SHALL provide a Terraform data source `tencentcloud_dbdc_db_custom_node_security_groups` that queries the security groups bound to a DBCustom node via the `DescribeDBCustomNodeSecurityGroups` cloud API.

#### Scenario: Query security groups for a valid node
- **WHEN** user configures the data source with a valid `node_id`
- **THEN** the system SHALL call `DescribeDBCustomNodeSecurityGroups` with the provided `node_id`
- **AND** the system SHALL return the `groups` list containing all security groups bound to the node

#### Scenario: Query security groups for a non-existent node
- **WHEN** user configures the data source with a non-existent `node_id`
- **THEN** the system SHALL return a `NonRetryableError` after retry exhaustion
- **AND** the system SHALL log the error with `[DATASOURCE]` prefix for troubleshooting

### Requirement: Security group fields are correctly mapped to Terraform schema
The system SHALL map all `SecurityGroup` fields from the cloud API response to the Terraform `groups` schema, including nested `PolicyRule` lists for inbound and outbound rules.

#### Scenario: Security group with all fields populated
- **WHEN** the API returns a security group with all fields populated (SecurityGroupId, SecurityGroupName, SecurityGroupRemark, ProjectId, CreateTime, Inbound, Outbound)
- **THEN** the Terraform state SHALL contain all corresponding fields in the `groups` list

#### Scenario: Security group with nil optional fields
- **WHEN** the API returns a security group with some optional fields set to nil
- **THEN** the Terraform state SHALL not include those nil fields
- **AND** the data source SHALL not crash on nil pointer dereference

### Requirement: Data source handles empty API responses gracefully
The system SHALL detect empty or nil API responses within the retry block and return a `NonRetryableError` to avoid data loss.

#### Scenario: API returns nil response
- **WHEN** the `DescribeDBCustomNodeSecurityGroups` API returns `result == nil` or `result.Response == nil`
- **THEN** the system SHALL return `resource.NonRetryableError` in the retry block
- **AND** the system SHALL log `[DATASOURCE] read empty, skip SetId` when retry is exhausted

### Requirement: Data source is registered in the provider
The system SHALL register the new data source `tencentcloud_dbdc_db_custom_node_security_groups` in `tencentcloud/provider.go` under the data sources section.

#### Scenario: Provider registration
- **WHEN** the Terraform provider is initialized
- **THEN** the data source `tencentcloud_dbdc_db_custom_node_security_groups` SHALL be available for use
- **AND** it SHALL be mapped to the `DataSourceTencentCloudDbdcDbCustomNodeSecurityGroups()` function