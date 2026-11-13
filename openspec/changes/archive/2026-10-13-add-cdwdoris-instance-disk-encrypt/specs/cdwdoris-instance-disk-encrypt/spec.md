## ADDED Requirements

### Requirement: Disk encryption on cdwdoris instance creation
The `tencentcloud_cdwdoris_instance` resource SHALL support an optional `disk_encrypt` parameter (TypeInt, immutable after creation) that is passed to the `CreateInstanceNew` API request as `DiskEncrypt`. Changes to `disk_encrypt` after creation SHALL be rejected with an error via the immutable args check in the Update function.

#### Scenario: Create cdwdoris instance with disk_encrypt
- **WHEN** a user specifies `disk_encrypt = 1` in the `tencentcloud_cdwdoris_instance` resource configuration
- **THEN** the provider SHALL pass `DiskEncrypt=1` in the `CreateInstanceNew` API request

#### Scenario: Create cdwdoris instance without disk_encrypt
- **WHEN** a user does NOT specify `disk_encrypt` in the `tencentcloud_cdwdoris_instance` resource configuration
- **THEN** the provider SHALL NOT set `DiskEncrypt` in the `CreateInstanceNew` API request

#### Scenario: Update disk_encrypt triggers immutable error
- **WHEN** a user changes `disk_encrypt` after creation
- **THEN** the provider SHALL return an error: `argument 'disk_encrypt' cannot be changed`

### Requirement: Refresh encryption status on cdwdoris instance read
The `tencentcloud_cdwdoris_instance` resource SHALL support a Computed `encrypt` parameter (TypeInt) that is refreshed from the `DescribeInstance` API response. The value is read from `response.InstanceInfo.MasterSummary.Encrypt`, falling back to `response.InstanceInfo.CoreSummary.Encrypt` when `MasterSummary` or its `Encrypt` is nil.

#### Scenario: Read existing cdwdoris instance and populate encrypt
- **WHEN** the provider reads an existing `tencentcloud_cdwdoris_instance` resource and `MasterSummary.Encrypt` is not nil
- **THEN** `encrypt` SHALL be set to the value of `MasterSummary.Encrypt`

#### Scenario: Read existing cdwdoris instance with encrypt only on CoreSummary
- **WHEN** the provider reads an existing `tencentcloud_cdwdoris_instance` resource and `MasterSummary` is nil or `MasterSummary.Encrypt` is nil but `CoreSummary.Encrypt` is not nil
- **THEN** `encrypt` SHALL be set to the value of `CoreSummary.Encrypt`

#### Scenario: Read existing cdwdoris instance with no encrypt returned
- **WHEN** the provider reads an existing `tencentcloud_cdwdoris_instance` resource and both `MasterSummary.Encrypt` and `CoreSummary.Encrypt` are nil
- **THEN** the provider SHALL NOT set the `encrypt` field

### Requirement: Immutable args check includes disk_encrypt
The `tencentcloud_cdwdoris_instance` Update function SHALL include `disk_encrypt` in the `immutableArgs` array so that any change to it is rejected.

#### Scenario: Immutable args array definition
- **WHEN** defining the immutable args in the Update function
- **THEN** the array SHALL contain `disk_encrypt` alongside the existing immutable arguments

#### Scenario: disk_encrypt change returns error
- **WHEN** `disk_encrypt` changes in the Update function
- **THEN** the provider SHALL return a formatted error `argument 'disk_encrypt' cannot be changed`