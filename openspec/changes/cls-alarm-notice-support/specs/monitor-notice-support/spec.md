## ADDED Requirements

### Requirement: Support MonitorNotice parameter configuration

The `tencentcloud_cls_alarm` resource SHALL support a `monitor_notice` configuration block to associate alarms with the observable platform notification templates via the CLS CreateAlarm API's `MonitorNotice` parameter.

#### Scenario: Create alarm with monitor_notice

- **WHEN** user configures `monitor_notice` block with valid `notices` entries
- **THEN** the resource creates an alarm using `CreateAlarmRequest.MonitorNotice` field
- **AND** the alarm is successfully associated with the specified observable platform notification templates

#### Scenario: Read alarm with monitor_notice

- **WHEN** an alarm was created or updated with `monitor_notice` configuration
- **THEN** the resource reads the `MonitorNotice` field from the API response
- **AND** populates the `monitor_notice` block in Terraform state with the returned values

#### Scenario: Update alarm from alarm_notice_ids to monitor_notice

- **WHEN** user changes alarm configuration from `alarm_notice_ids` to `monitor_notice`
- **THEN** the resource updates the alarm via ModifyAlarm API
- **AND** the alarm switches from CLS notification channels to observable platform notification templates

### Requirement: MonitorNotice structure support

The `monitor_notice` block SHALL contain a `notices` list, where each entry includes `notice_id` (required), `content_tmpl_id` (optional), and `alarm_levels` (required) fields matching the SDK's `MonitorNoticeRule` structure.

#### Scenario: Configure notice with all fields

- **WHEN** user specifies `notice_id`, `content_tmpl_id`, and `alarm_levels` in a `notices` entry
- **THEN** all three fields are mapped to `MonitorNoticeRule` and sent to the API
- **AND** the alarm uses the specified content template for the notification

#### Scenario: Configure notice without content_tmpl_id

- **WHEN** user omits `content_tmpl_id` in a `notices` entry
- **THEN** the `ContentTmplId` field in `MonitorNoticeRule` is not set or set to nil
- **AND** the API uses the default content template for that notice

#### Scenario: Configure multiple notices with different alarm levels

- **WHEN** user configures multiple `notices` entries with different `alarm_levels`
- **THEN** each notice rule is correctly mapped to the `MonitorNotice.Notices` array
- **AND** each notice triggers only for its specified alarm levels

### Requirement: alarm_notice_ids becomes optional

The `alarm_notice_ids` field SHALL be changed from `Required` to `Optional` to align with the CLS API specification where it is optional and mutually exclusive with `MonitorNotice`.

#### Scenario: Existing configuration with alarm_notice_ids continues to work

- **WHEN** an existing Terraform configuration uses `alarm_notice_ids` (without `monitor_notice`)
- **THEN** the configuration remains valid after upgrading the provider
- **AND** the alarm continues to use CLS notification channels

#### Scenario: Create new alarm with alarm_notice_ids only

- **WHEN** user creates a new alarm with only `alarm_notice_ids` configured
- **THEN** the alarm is created successfully using `CreateAlarmRequest.AlarmNoticeIds`
- **AND** no `MonitorNotice` field is sent to the API

### Requirement: Mutual exclusivity between alarm_notice_ids and monitor_notice

The resource SHALL enforce mutual exclusivity between `alarm_notice_ids` and `monitor_notice` using Terraform schema's `ExactlyOneOf` constraint, ensuring users configure exactly one notification method.

#### Scenario: Configuration with both fields is rejected

- **WHEN** user attempts to configure both `alarm_notice_ids` and `monitor_notice`
- **THEN** Terraform plan phase fails with a validation error
- **AND** the error message indicates the two fields are mutually exclusive

#### Scenario: Configuration with neither field is rejected

- **WHEN** user attempts to create an alarm without either `alarm_notice_ids` or `monitor_notice`
- **THEN** Terraform plan phase fails with a validation error
- **AND** the error message indicates at least one notification method must be configured

#### Scenario: Switching between notification methods

- **WHEN** user updates configuration to remove `alarm_notice_ids` and add `monitor_notice`
- **THEN** Terraform detects a configuration change and updates the alarm
- **AND** the alarm switches notification methods successfully

### Requirement: Documentation update

The resource documentation SHALL clearly describe the `monitor_notice` block structure, field meanings, mutual exclusivity with `alarm_notice_ids`, and provide usage examples for both notification methods.

#### Scenario: User reads documentation for monitor_notice

- **WHEN** user consults the resource documentation
- **THEN** the documentation includes a `monitor_notice` section with field descriptions
- **AND** includes an example HCL configuration using `monitor_notice`
- **AND** explains that `monitor_notice` and `alarm_notice_ids` cannot be used together
