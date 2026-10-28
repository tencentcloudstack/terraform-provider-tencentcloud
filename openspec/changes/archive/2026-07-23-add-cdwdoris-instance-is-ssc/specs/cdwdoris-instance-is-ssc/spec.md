## ADDED Requirements

### Requirement: Storage-compute separation parameter on instance creation
The `tencentcloud_cdwdoris_instance` resource SHALL support an optional `is_ssc` parameter of type `TypeBool` that controls whether the CDW Doris instance uses storage-compute separation (存算分离) architecture. The parameter SHALL be passed to the `CreateInstanceNew` API as `IsSSC`.

#### Scenario: Create instance with storage-compute separation enabled
- **WHEN** a user sets `is_ssc = true` in the `tencentcloud_cdwdoris_instance` resource configuration
- **THEN** the Terraform provider SHALL pass `IsSSC = true` in the `CreateInstanceNew` API request

#### Scenario: Create instance without storage-compute separation
- **WHEN** a user does not set `is_ssc` in the `tencentcloud_cdwdoris_instance` resource configuration
- **THEN** the Terraform provider SHALL NOT pass `IsSSC` in the `CreateInstanceNew` API request (default behavior preserved)

#### Scenario: Parameter is immutable after creation
- **WHEN** a user attempts to change the `is_ssc` value on an existing `tencentcloud_cdwdoris_instance` resource
- **THEN** the Terraform provider SHALL return an error indicating that the parameter cannot be changed, forcing a destroy-and-recreate of the resource