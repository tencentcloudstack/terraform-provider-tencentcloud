## ADDED Requirements

### Requirement: Support SecurityGroupId parameter in VPC Endpoint resource
The tencentcloud_vpc_end_point resource SHALL support SecurityGroupId as an optional string parameter for configuring security groups.

#### Scenario: Create VPC Endpoint with SecurityGroupId
- **WHEN** user specifies SecurityGroupId parameter in tencentcloud_vpc_end_point resource
- **THEN** the provider SHALL pass SecurityGroupId to CreateVpcEndPoint API
- **AND** the provider SHALL store SecurityGroupId in Terraform state

#### Scenario: Read VPC Endpoint SecurityGroupId
- **WHEN** the provider reads an existing VPC Endpoint resource
- **THEN** the provider SHALL retrieve SecurityGroupId from DescribeVpcEndPoints API
- **AND** the provider SHALL update Terraform state with the retrieved SecurityGroupId

#### Scenario: Update VPC Endpoint SecurityGroupId
- **WHEN** user updates SecurityGroupId parameter in tencentcloud_vpc_end_point resource
- **THEN** the provider SHALL update SecurityGroupId via ModifyVpcEndPointAttribute API or recreate resource
- **AND** the provider SHALL update Terraform state with the new SecurityGroupId

### Requirement: Support Tags parameter in VPC Endpoint resource
The tencentcloud_vpc_end_point resource SHALL support Tags as an optional list parameter for resource tagging and management.

#### Scenario: Create VPC Endpoint with Tags
- **WHEN** user specifies Tags parameter in tencentcloud_vpc_end_point resource
- **THEN** the provider SHALL pass Tags list to CreateVpcEndPoint API
- **AND** the provider SHALL validate that each Tag entry has required Key field
- **AND** the provider SHALL store Tags in Terraform state

#### Scenario: Read VPC Endpoint Tags
- **WHEN** the provider reads an existing VPC Endpoint resource
- **THEN** the provider SHALL retrieve Tags from DescribeVpcEndPoints API
- **AND** the provider SHALL update Terraform state with the retrieved Tags

#### Scenario: Update VPC Endpoint Tags
- **WHEN** user updates Tags parameter in tencentcloud_vpc_end_point resource
- **THEN** the provider SHALL update Tags via ModifyVpcEndPointAttribute API or related tag management API
- **AND** the provider SHALL update Terraform state with the new Tags

#### Scenario: Validate Tag structure
- **WHEN** user provides Tags parameter
- **THEN** the provider SHALL validate that each Tag has Key field (required)
- **AND** the provider SHALL accept optional Value field
- **AND** the provider SHALL return validation error if Key is missing

### Requirement: Support IpAddressType parameter in VPC Endpoint resource
The tencentcloud_vpc_end_point resource SHALL support IpAddressType as an optional string parameter for IP address type configuration, supporting Ipv4 and Ipv6 with default value Ipv4.

#### Scenario: Create VPC Endpoint with IpAddressType
- **WHEN** user specifies IpAddressType parameter in tencentcloud_vpc_end_point resource
- **THEN** the provider SHALL pass IpAddressType to CreateVpcEndPoint API
- **AND** the provider SHALL validate that IpAddressType is either "Ipv4" or "Ipv6"
- **AND** the provider SHALL store IpAddressType in Terraform state

#### Scenario: Create VPC Endpoint without IpAddressType (default behavior)
- **WHEN** user creates tencentcloud_vpc_end_point resource without specifying IpAddressType
- **THEN** the provider SHALL use default value "Ipv4"
- **AND** the provider SHALL pass "Ipv4" to CreateVpcEndPoint API
- **AND** the provider SHALL store "Ipv4" in Terraform state

#### Scenario: Read VPC Endpoint IpAddressType
- **WHEN** the provider reads an existing VPC Endpoint resource
- **THEN** the provider SHALL retrieve IpAddressType from DescribeVpcEndPoints API
- **AND** the provider SHALL update Terraform state with the retrieved IpAddressType

#### Scenario: Validate IpAddressType values
- **WHEN** user provides IpAddressType parameter
- **THEN** the provider SHALL validate that value is "Ipv4" or "Ipv6"
- **AND** the provider SHALL return validation error for invalid values
