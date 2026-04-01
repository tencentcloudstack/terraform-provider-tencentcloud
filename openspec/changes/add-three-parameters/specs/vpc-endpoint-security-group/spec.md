## ADDED Requirements

### Requirement: VPC endpoint resource supports security group binding
The tencentcloud_vpc_end_point resource SHALL allow users to specify a security group ID to bind to the VPC endpoint. The SecurityGroupId field SHALL be optional and accept a string value representing a valid security group ID.

#### Scenario: Create VPC endpoint with security group
- **WHEN** user creates a tencentcloud_vpc_end_point resource with SecurityGroupId set to a valid security group ID
- **THEN** the CreateVpcEndPoint API shall be called with the SecurityGroupId parameter
- **AND** the VPC endpoint shall be bound to the specified security group

#### Scenario: Create VPC endpoint without security group
- **WHEN** user creates a tencentcloud_vpc_end_point resource without specifying SecurityGroupId
- **THEN** the CreateVpcEndPoint API shall be called without the SecurityGroupId parameter
- **AND** the VPC endpoint shall be created without a security group binding

#### Scenario: Read VPC endpoint with security group
- **WHEN** Terraform reads an existing tencentcloud_vpc_end_point resource that has a security group bound
- **THEN** the DescribeVpcEndPoints API shall be called
- **AND** the returned SecurityGroupId value SHALL be stored in the Terraform state

#### Scenario: Update VPC endpoint security group
- **WHEN** user updates the SecurityGroupId field of an existing tencentcloud_vpc_end_point resource
- **THEN** the UpdateVpcEndPointAttribute API shall be called with the new SecurityGroupId value
- **AND** the VPC endpoint shall be bound to the new security group

#### Scenario: Import VPC endpoint with security group
- **WHEN** user imports an existing VPC endpoint that has a security group bound into Terraform state
- **THEN** the DescribeVpcEndPoints API shall be called
- **AND** the SecurityGroupId value SHALL be read and stored in the Terraform state