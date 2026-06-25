## Why

Currently, the TencentCloud Terraform Provider lacks a dedicated resource to manage manual backups for CDB (MySQL) instances. Users can only manage backups through the console or API directly. Adding `tencentcloud_mysql_backup` as an attachment resource enables users to create, query, and delete manual backups of MySQL instances declaratively through Terraform, providing consistent infrastructure-as-code management for database backup lifecycle.

## What Changes

- Add new RESOURCE_KIND_ATTACHMENT resource `tencentcloud_mysql_backup` to manage manual backups for CDB MySQL instances
- Support creating manual backups via `CreateBackup` API with configurable backup method, database/table selection, backup naming, and encryption settings
- Support reading backup status and details via `DescribeBackups` API
- Support deleting backups via `DeleteBackup` API
- New resource file: `resource_tc_mysql_backup_attachment.go` under `tencentcloud/services/cdb/`
- New documentation file for the resource

## Capabilities

### New Capabilities
- `mysql-backup-attachment`: Manage the lifecycle of manual MySQL backups (create, read, delete) through Terraform

### Modified Capabilities
<!-- None - this is a net new resource -->

## Impact

- Affected code: New file `tencentcloud/services/cdb/resource_tc_mysql_backup_attachment.go`
- Affected registration: `tencentcloud/provider.go` needs new resource registration entry
- Affected documentation: New file `tencentcloud/services/cdb/resource_tc_mysql_backup_attachment.md`
- Dependencies: Uses existing `tencentcloud-sdk-go/tencentcloud/cdb/v20170320` package (no new SDK dependency needed)
- No breaking changes to existing resources