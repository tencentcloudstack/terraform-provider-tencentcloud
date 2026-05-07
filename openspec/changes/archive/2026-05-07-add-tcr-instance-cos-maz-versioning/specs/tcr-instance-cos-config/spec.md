## ADDED Requirements

### Requirement: TCR instance SHALL support COS multi-AZ configuration

The `tencentcloud_tcr_instance` resource SHALL provide an `enable_cos_maz` field that allows users to enable or disable COS bucket multi-AZ (multi-availability zone) feature during instance creation.

#### Scenario: User enables COS multi-AZ during creation
- **WHEN** user sets `enable_cos_maz = true` in the Terraform configuration
- **THEN** the provider SHALL pass `EnableCosMAZ = true` to the CreateInstance API
- **THEN** the created TCR instance SHALL have COS multi-AZ enabled

#### Scenario: User disables COS multi-AZ during creation
- **WHEN** user sets `enable_cos_maz = false` in the Terraform configuration
- **THEN** the provider SHALL pass `EnableCosMAZ = false` to the CreateInstance API
- **THEN** the created TCR instance SHALL have COS multi-AZ disabled

#### Scenario: User does not specify COS multi-AZ setting
- **WHEN** user omits the `enable_cos_maz` field from the Terraform configuration
- **THEN** the provider SHALL NOT pass the `EnableCosMAZ` parameter to the CreateInstance API
- **THEN** the API SHALL use its default value (false)
- **THEN** the provider SHALL read the actual value from the API response and store it in state

#### Scenario: Provider reads COS multi-AZ setting
- **WHEN** the provider performs a read operation on an existing TCR instance
- **THEN** the provider SHALL retrieve the `EnableCosMAZ` field from the DescribeInstances API response
- **THEN** the provider SHALL set the `enable_cos_maz` field in the Terraform state with the retrieved value

### Requirement: TCR instance SHALL support COS versioning configuration

The `tencentcloud_tcr_instance` resource SHALL provide an `enable_cos_versioning` field that allows users to enable or disable COS bucket versioning feature during instance creation.

#### Scenario: User enables COS versioning during creation
- **WHEN** user sets `enable_cos_versioning = true` in the Terraform configuration
- **THEN** the provider SHALL pass `EnableCosVersioning = true` to the CreateInstance API
- **THEN** the created TCR instance SHALL have COS versioning enabled

#### Scenario: User disables COS versioning during creation
- **WHEN** user sets `enable_cos_versioning = false` in the Terraform configuration
- **THEN** the provider SHALL pass `EnableCosVersioning = false` to the CreateInstance API
- **THEN** the created TCR instance SHALL have COS versioning disabled

#### Scenario: User does not specify COS versioning setting
- **WHEN** user omits the `enable_cos_versioning` field from the Terraform configuration
- **THEN** the provider SHALL NOT pass the `EnableCosVersioning` parameter to the CreateInstance API
- **THEN** the API SHALL use its default value (false)
- **THEN** the provider SHALL read the actual value from the API response and store it in state

#### Scenario: Provider reads COS versioning setting
- **WHEN** the provider performs a read operation on an existing TCR instance
- **THEN** the provider SHALL retrieve the `EnableCosVersioning` field from the DescribeInstances API response
- **THEN** the provider SHALL set the `enable_cos_versioning` field in the Terraform state with the retrieved value

### Requirement: Schema fields SHALL be backward compatible

The new fields SHALL NOT break existing Terraform configurations or state files.

#### Scenario: Existing resources without new fields
- **WHEN** an existing TCR instance is refreshed or planned after provider upgrade
- **THEN** the provider SHALL NOT mark the resource as needing replacement
- **THEN** the provider SHALL add the new fields to the state with values retrieved from the API
- **THEN** no configuration changes SHALL be required from the user

#### Scenario: New fields in schema
- **WHEN** the new fields are added to the resource schema
- **THEN** both fields SHALL be marked as `Optional` to allow user specification
- **THEN** both fields SHALL be marked as `Computed` to allow API-provided default values
- **THEN** both fields SHALL use `TypeBool` to match the API data type
