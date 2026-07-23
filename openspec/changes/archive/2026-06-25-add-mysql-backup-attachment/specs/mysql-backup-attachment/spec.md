## ADDED Requirements

### Requirement: Create manual MySQL backup
The system SHALL allow users to create a manual backup for a CDB MySQL instance by calling the `CreateBackup` API with the specified parameters.

#### Scenario: Create a full instance backup with physical method
- **WHEN** user provides `instance_id`, `backup_method` set to "physical", and optional `manual_backup_name` and `encryption_flag`
- **THEN** the system SHALL call `CreateBackup` API, obtain a `BackupId`, set the resource ID as `backup_id#instance_id`, and read back the backup details

#### Scenario: Create a logical backup with database/table selection
- **WHEN** user provides `instance_id`, `backup_method` set to "logical", and `backup_db_table_list` with database and table entries
- **THEN** the system SHALL call `CreateBackup` API with the `BackupDBTableList` parameter populated from the user's configuration

#### Scenario: Create snapshot backup for basic edition instance
- **WHEN** user provides `instance_id`, `backup_method` set to "snapshot"
- **THEN** the system SHALL call `CreateBackup` API and create a snapshot backup

### Requirement: Read MySQL backup details
The system SHALL allow users to query the status and details of a created backup via the `DescribeBackups` API.

#### Scenario: Read existing backup
- **WHEN** user runs `terraform refresh` or the Read function is called
- **THEN** the system SHALL call `DescribeBackups` with the `instance_id`, iterate the returned items to find the matching `backup_id`, and populate the resource state with backup details

#### Scenario: Backup deleted externally
- **WHEN** the `DescribeBackups` API returns no item matching the `backup_id`
- **THEN** the system SHALL call `d.SetId("")` to remove the resource from state and log a warning

### Requirement: Delete MySQL backup
The system SHALL allow users to delete a manual backup by calling the `DeleteBackup` API.

#### Scenario: Delete an existing backup
- **WHEN** user runs `terraform destroy` or removes the resource from configuration
- **THEN** the system SHALL call `DeleteBackup` API with `instance_id` and `backup_id` parsed from the resource ID

### Requirement: Import existing MySQL backup
The system SHALL support importing existing backups into Terraform state using a composite ID.

#### Scenario: Import backup by composite ID
- **WHEN** user runs `terraform import tencentcloud_mysql_backup.foo backup_id#instance_id`
- **THEN** the system SHALL parse the composite ID, call `DescribeBackups` to verify the backup exists, and populate the resource state

### Requirement: Immutable backup parameters
The system SHALL treat all configurable parameters as ForceNew, so any parameter change triggers destroy-and-recreate of the backup resource.

#### Scenario: Change backup method
- **WHEN** user modifies `backup_method` from "physical" to "logical"
- **THEN** the system SHALL destroy the existing backup (call `DeleteBackup`) and create a new backup (call `CreateBackup`) with the new parameters
