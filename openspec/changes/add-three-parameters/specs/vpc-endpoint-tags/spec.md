## ADDED Requirements

### Requirement: VPC endpoint resource supports tag management
The tencentcloud_vpc_end_point resource SHALL allow users to specify a list of tags to attach to the VPC endpoint. Each tag SHALL consist of a Key (required) and Value (optional) field. The Tags field SHALL be optional and accept a list of tag objects.

#### Scenario: Create VPC endpoint with single tag
- **WHEN** user creates a tencentcloud_vpc_end_point resource with Tags set to a list containing one tag with Key="env" and Value="prod"
- **THEN** the CreateVpcEndPoint API shall be called with the Tags parameter
- **AND** the VPC endpoint shall be tagged with the specified tag

#### Scenario: Create VPC endpoint with multiple tags
- **WHEN** user creates a tencentcloud_vpc_end_point resource with Tags set to a list containing multiple tags
- **THEN** the CreateVpcEndPoint API shall be called with all specified tags
- **AND** the VPC endpoint shall be tagged with all the specified tags

#### Scenario: Create VPC endpoint without tags
- **WHEN** user creates a tencentcloud_vpc_end_point resource without specifying Tags
- **THEN** the CreateVpcEndPoint API shall be called without the Tags parameter
- **AND** the VPC endpoint shall be created without any tags

#### Scenario: Create VPC endpoint with tag without Value
- **WHEN** user creates a tencentcloud_vpc_end_point resource with Tags set to a list containing a tag with Key="category" and Value omitted
- **THEN** the CreateVpcEndPoint API shall be called with the tag Key and empty Value
- **AND** the VPC endpoint shall be tagged with a tag that has only the Key field

#### Scenario: Read VPC endpoint with tags
- **WHEN** Terraform reads an existing tencentcloud_vpc_end_point resource that has tags
- **THEN** the DescribeVpcEndPoints API shall be called
- **AND** the returned Tags list SHALL be stored in the Terraform state

#### Scenario: Update VPC endpoint tags
- **WHEN** user updates the Tags field of an existing tencentcloud_vpc_end_point resource
- **THEN** the UpdateVpcEndPointAttribute API shall be called with the new Tags value
- **AND** the VPC endpoint shall be tagged with the new list of tags

#### Scenario: Import VPC endpoint with tags
- **WHEN** user imports an existing VPC endpoint that has tags into Terraform state
- **THEN** the DescribeVpcEndPoints API shall be called
- **AND** the Tags list SHALL be read and stored in the Terraform state