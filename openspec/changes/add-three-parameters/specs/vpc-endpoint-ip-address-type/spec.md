## ADDED Requirements

### Requirement: VPC endpoint resource supports IP address type configuration
The tencentcloud_vpc_end_point resource SHALL allow users to specify the IP address type for the VPC endpoint. The IpAddressType field SHALL be optional and accept a string value of either "Ipv4" or "Ipv6". The default value SHALL be "Ipv4".

#### Scenario: Create VPC endpoint with Ipv4 address type
- **WHEN** user creates a tencentcloud_vpc_end_point resource with IpAddressType set to "Ipv4"
- **THEN** the CreateVpcEndPoint API shall be called with IpAddressType="Ipv4"
- **AND** the VPC endpoint shall be configured to use IPv4

#### Scenario: Create VPC endpoint with Ipv6 address type
- **WHEN** user creates a tencentcloud_vpc_end_point resource with IpAddressType set to "Ipv6"
- **THEN** the CreateVpcEndPoint API shall be called with IpAddressType="Ipv6"
- **AND** the VPC endpoint shall be configured to use IPv6

#### Scenario: Create VPC endpoint without specifying address type
- **WHEN** user creates a tencentcloud_vpc_end_point resource without specifying IpAddressType
- **THEN** the CreateVpcEndPoint API shall be called without the IpAddressType parameter
- **AND** the VPC endpoint shall be created with the default IPv4 address type

#### Scenario: Read VPC endpoint with address type
- **WHEN** Terraform reads an existing tencentcloud_vpc_end_point resource
- **THEN** the DescribeVpcEndPoints API shall be called
- **AND** the returned IpAddressType value SHALL be stored in the Terraform state
- **AND** if the API does not return IpAddressType, the value SHALL default to "Ipv4"

#### Scenario: Update VPC endpoint address type
- **WHEN** user updates the IpAddressType field of an existing tencentcloud_vpc_end_point resource from "Ipv4" to "Ipv6"
- **THEN** the UpdateVpcEndPointAttribute API shall be called with the new IpAddressType value
- **AND** the VPC endpoint shall be configured to use the new IP address type

#### Scenario: Import VPC endpoint with address type
- **WHEN** user imports an existing VPC endpoint into Terraform state
- **THEN** the DescribeVpcEndPoints API shall be called
- **AND** the IpAddressType value SHALL be read and stored in the Terraform state