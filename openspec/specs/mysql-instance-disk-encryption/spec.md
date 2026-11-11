# mysql-instance-disk-encryption Specification

## Purpose
TBD - created by archiving change add-mysql-instance-disk-encryption. Update Purpose after archive.
## Requirements
### Requirement: Disk encryption parameter on mysql instance creation
The `tencentcloud_mysql_instance` resource SHALL support an optional `disk_encryption` schema parameter (TypeString, Optional, ForceNew, Computed) that controls whether disk encryption is enabled when creating a CDB instance. When the value is `on`, disk encryption SHALL be enabled; otherwise disk encryption SHALL not be enabled. The parameter SHALL be `ForceNew` because no Modify/Update CDB API accepts `DiskEncryption`.

#### Scenario: Create month-paid instance with disk encryption enabled
- **WHEN** a user sets `disk_encryption = "on"` in the `tencentcloud_mysql_instance` resource configuration and selects month-paid billing
- **THEN** the provider SHALL set `request.DiskEncryption = "on"` in the `CreateDBInstance` API request

#### Scenario: Create hour-paid instance with disk encryption enabled
- **WHEN** a user sets `disk_encryption = "on"` in the `tencentcloud_mysql_instance` resource configuration and selects hourly-paid billing
- **THEN** the provider SHALL set `request.DiskEncryption = "on"` in the `CreateDBInstanceHour` API request

#### Scenario: Create instance without explicit disk encryption
- **WHEN** a user does NOT specify `disk_encryption` in the `tencentcloud_mysql_instance` resource configuration
- **THEN** the provider SHALL NOT set `DiskEncryption` in either the `CreateDBInstance` or `CreateDBInstanceHour` API request

#### Scenario: Change disk encryption after creation is rejected
- **WHEN** a user changes `disk_encryption` on an existing `tencentcloud_mysql_instance` resource
- **THEN** the provider SHALL force replacement of the resource (ForceNew) because the value cannot be modified in place

### Requirement: Read disk encryption from Describe API
The `tencentcloud_mysql_instance` resource Read function SHALL refresh the `disk_encryption` state field from the `DescribeDBInstances` API response (`InstanceInfo.DiskEncryption`). The provider SHALL guard against a nil `DiskEncryption` field before setting it into state.

#### Scenario: Refresh disk encryption from Describe response
- **WHEN** the provider reads an existing `tencentcloud_mysql_instance` resource and the `DescribeDBInstances` response returns a non-nil `DiskEncryption` field
- **THEN** the provider SHALL set `disk_encryption` in Terraform state to the returned value

#### Scenario: Describe returns nil disk encryption
- **WHEN** the provider reads an existing `tencentcloud_mysql_instance` resource and the `DescribeDBInstances` response `DiskEncryption` field is nil
- **THEN** the provider SHALL skip setting `disk_encryption` into state (nil guard)

### Requirement: Consistent create wiring across billing models and instance roles
The provider SHALL set `DiskEncryption` in the shared `mysqlAllInstanceRoleSet` helper so that the parameter is applied consistently to both `CreateDBInstanceRequest` and `CreateDBInstanceHourRequest` for all instance roles (master, dr, ro).

#### Scenario: Shared helper sets field on both request types
- **WHEN** `disk_encryption` is set in the resource configuration
- **THEN** the `mysqlAllInstanceRoleSet` helper SHALL set `DiskEncryption` on the `CreateDBInstanceRequest` when billing is month-paid
- **AND** SHALL set `DiskEncryption` on the `CreateDBInstanceHourRequest` when billing is hourly-paid

