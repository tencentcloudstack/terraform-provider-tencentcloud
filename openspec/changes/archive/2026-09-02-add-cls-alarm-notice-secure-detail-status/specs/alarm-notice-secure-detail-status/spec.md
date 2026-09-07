## ADDED Requirements

### Requirement: Support secure_detail_status parameter

The `tencentcloud_cls_alarm_notice` resource SHALL support a `secure_detail_status` parameter (type int, optional) to control the alarm detail security authentication redirect switch, mapping to the CLS API `SecureDetailStatus` field with enum values 1 (off, default) and 2 (on).

#### Scenario: Create alarm notice with secure_detail_status

- **WHEN** user configures `secure_detail_status` with value 2 in the resource configuration
- **THEN** the resource creates an alarm notice via `CreateAlarmNoticeRequest` with `SecureDetailStatus` set to 2
- **AND** the alarm notice is created with the security authentication redirect switch enabled

#### Scenario: Create alarm notice without secure_detail_status

- **WHEN** user does not configure `secure_detail_status`
- **THEN** the resource creates an alarm notice without explicitly setting `SecureDetailStatus`
- **AND** the cloud API applies the default value (1, off)

#### Scenario: Read alarm notice with secure_detail_status

- **WHEN** the resource reads an existing alarm notice via `DescribeAlarmNotices`
- **THEN** the resource reads the `SecureDetailStatus` field from the `AlarmNotice` response
- **AND** populates `secure_detail_status` in Terraform state with the returned value

#### Scenario: Read alarm notice where SecureDetailStatus is nil

- **WHEN** the API response returns `SecureDetailStatus` as nil
- **THEN** the resource skips setting `secure_detail_status` in state
- **AND** no error is raised

#### Scenario: Update alarm notice secure_detail_status

- **WHEN** user changes `secure_detail_status` from 1 to 2
- **THEN** the resource detects the change via `d.HasChange("secure_detail_status")`
- **AND** calls `ModifyAlarmNotice` with `SecureDetailStatus` set to 2
- **AND** the alarm notice security authentication redirect switch is updated

### Requirement: secure_detail_status schema definition

The `secure_detail_status` field SHALL be defined as `schema.TypeInt` with `Optional: true` and `Computed: true` in the resource schema, consistent with existing numeric switch fields like `deliver_status` and `alarm_shield_status`.

#### Scenario: Field is optional and computed

- **WHEN** user creates a configuration without `secure_detail_status`
- **THEN** Terraform plan succeeds without requiring the field
- **AND** after apply, the field is populated in state with the API-returned default value
