# dbdc-db-custom-disaster-recover-group-resource Specification

## Purpose
TBD - created by archiving change add-dbdc-db-custom-disaster-recover-group. Update Purpose after archive.
## Requirements
### Requirement: Create placement group
The system SHALL allow users to create a DB Custom placement group via
`tencentcloud_dbdc_db_custom_disaster_recover_group`, calling
`CreateDBCustomDisasterRecoverGroup` with `name` (required), and optional
`type`, `strategy`, `affinity`, `tags`, and `client_token`. After a successful
create call, the system SHALL poll `DescribeDBCustomDisasterRecoverGroups`
until the group `Status` is `Available` before returning, and SHALL set the
resource id to the returned `DisasterRecoverGroupId`.

#### Scenario: Successful create with required name only
- **WHEN** a user applies a configuration specifying only `name`
- **THEN** the provider calls `CreateDBCustomDisasterRecoverGroup` with `name`,
  polls until `Status == "Available"`, sets the resource id to the returned
  `DisasterRecoverGroupId`, and reads the computed fields

#### Scenario: Create fails when group id is empty
- **WHEN** `CreateDBCustomDisasterRecoverGroup` returns a nil/empty
  `DisasterRecoverGroupId`
- **THEN** the provider returns a `NonRetryableError` and does not write an
  empty id to state

#### Scenario: Create fails on CreateFailed status
- **WHEN** the polled group `Status` becomes `CreateFailed`
- **THEN** the provider returns a non-retryable error describing the failure

### Requirement: Read placement group
The system SHALL read a placement group by calling
`DescribeDBCustomDisasterRecoverGroups` filtered by `DisasterRecoverGroupIds`
set to the resource id. It SHALL populate all schema fields from the returned
`DisasterRecoverGroup`, guarding every pointer/slice for nil. If the group is
not found, the system SHALL log the id and clear the resource id from state.

#### Scenario: Group exists
- **WHEN** the group with the resource id exists
- **THEN** the provider sets `name`, `type`, `strategy`, `affinity`, `tags`,
  `status`, `node_quota_total`, `current_num`, `created_time`, `node_ids`, and
  `disaster_recover_group_id` from the API response (skipping nil fields)

#### Scenario: Group not found
- **WHEN** `DescribeDBCustomDisasterRecoverGroups` returns no matching group
- **THEN** the provider logs `[CRUD] disaster_recover_group id=<id>` and calls
  `d.SetId("")` so Terraform removes it from state

### Requirement: Update placement group attributes
The system SHALL support in-place update of `name` and `affinity` by calling
`ModifyDBCustomDisasterRecoverGroupAttribute` with `DisasterRecoverGroupId`,
`Name`, and `Affinity`. The system SHALL reject changes to immutable fields
(`type`, `strategy`, `tags`, `client_token`) by checking an `immutableArgs`
list and returning an error when a change is detected.

#### Scenario: Update name and affinity
- **WHEN** a user changes `name` and/or `affinity`
- **THEN** the provider calls `ModifyDBCustomDisasterRecoverGroupAttribute`
  with the new values and re-reads the resource

#### Scenario: Reject change to immutable field
- **WHEN** a user changes `type`, `strategy`, `tags`, or `client_token`
- **THEN** the provider returns an error indicating the field is immutable and
  requires recreating the resource

### Requirement: Delete placement group
The system SHALL delete a placement group by calling
`DeleteDBCustomDisasterRecoverGroups` with `DisasterRecoverGroupIds` set to the
resource id. Because delete is asynchronous, the system SHALL poll
`DescribeDBCustomTaskStatus` using the returned `TaskId` until the status is
`Succeeded`, reusing the shared `waitDBCustomTaskSucceeded` helper, and SHALL
honor the configured delete timeout.

#### Scenario: Successful delete
- **WHEN** a user destroys the resource
- **THEN** the provider calls `DeleteDBCustomDisasterRecoverGroups`, polls the
  returned `TaskId` until `Succeeded`, and returns nil

#### Scenario: Delete task fails
- **WHEN** the polled task `Status` is `Failed`
- **THEN** the provider returns an error describing the task failure

### Requirement: Import placement group
The system SHALL support importing a placement group by its
`DisasterRecoverGroupId` via `ImportStatePassthrough`, after which the Read
flow populates the state.

#### Scenario: Import by id
- **WHEN** a user runs `terraform import
  tencentcloud_dbdc_db_custom_disaster_recover_group.example <group-id>`
- **THEN** the provider sets the resource id to `<group-id>` and runs Read to
  populate the remaining fields

### Requirement: Retry and nil-safety conventions
Every cloud API call SHALL be wrapped in `resource.Retry` using
`tccommon.ReadRetryTimeout` for reads and `tccommon.WriteRetryTimeout` for
writes, with errors wrapped via `tccommon.RetryError`. All response and pointer
fields SHALL be nil-checked before access, and setting id/state SHALL occur
outside the retry block.

#### Scenario: Transient API failure retried
- **WHEN** a cloud API call fails with a retryable error
- **THEN** the provider retries within the configured timeout using
  `tccommon.RetryError` semantics

#### Scenario: Non-retryable error surfaces immediately
- **WHEN** a cloud API call fails with a non-retryable error
- **THEN** the provider returns the error without further retries

