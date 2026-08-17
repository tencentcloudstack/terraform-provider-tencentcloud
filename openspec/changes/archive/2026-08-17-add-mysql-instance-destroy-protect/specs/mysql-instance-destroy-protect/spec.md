## ADDED Requirements

### Requirement: destroy_protect parameter on mysql instance creation
The `tencentcloud_mysql_instance` resource SHALL support an optional `destroy_protect` parameter (TypeString, values `on`/`off`, Computed) that is passed to the `CreateDBInstance` and `CreateDBInstanceHour` APIs as `DestroyProtect`. When the user does not specify `destroy_protect`, the provider SHALL NOT set `DestroyProtect` in the create API request (API uses default).

#### Scenario: Create mysql instance with destroy protection enabled
- **WHEN** a user specifies `destroy_protect = "on"` in the `tencentcloud_mysql_instance` resource configuration
- **THEN** the provider SHALL pass `DestroyProtect=on` in the `CreateDBInstance` or `CreateDBInstanceHour` API request

#### Scenario: Create mysql instance with destroy protection disabled
- **WHEN** a user specifies `destroy_protect = "off"` in the `tencentcloud_mysql_instance` resource configuration
- **THEN** the provider SHALL pass `DestroyProtect=off` in the `CreateDBInstance` or `CreateDBInstanceHour` API request

#### Scenario: Create mysql instance without destroy_protect
- **WHEN** a user does NOT specify `destroy_protect` in the `tencentcloud_mysql_instance` resource configuration
- **THEN** the provider SHALL NOT set `DestroyProtect` in the create API request

### Requirement: Read destroy_protect from DescribeDBInstances response
The provider SHALL read `DestroyProtect` from the `DescribeDBInstances` API response (`InstanceInfo.DestroyProtect`) during the Read operation and set it in the Terraform state.

#### Scenario: Read destroy_protect from existing instance
- **WHEN** the provider reads an existing `tencentcloud_mysql_instance` resource
- **THEN** `destroy_protect` SHALL be refreshed from the `DescribeDBInstances` API response (`InstanceInfo.DestroyProtect`)
- **AND** the provider SHALL perform a nil check on the `DestroyProtect` pointer before setting the state value

#### Scenario: Read handles nil DestroyProtect gracefully
- **WHEN** the `DescribeDBInstances` API response returns `nil` for `InstanceInfo.DestroyProtect`
- **THEN** the provider SHALL skip setting `destroy_protect` in state (no panic, no error)

### Requirement: destroy_protect schema definition
The `destroy_protect` schema field SHALL be defined in the `TencentMsyqlBasicInfo()` function with the following properties: TypeString, Optional, Computed.

#### Scenario: Schema field properties
- **WHEN** the schema is defined for `destroy_protect`
- **THEN** the field SHALL have Type set to `schema.TypeString`
- **AND** the field SHALL have Optional set to `true`
- **AND** the field SHALL have Computed set to `true`
- **AND** the field SHALL have a Description explaining the valid values `on` and `off`
