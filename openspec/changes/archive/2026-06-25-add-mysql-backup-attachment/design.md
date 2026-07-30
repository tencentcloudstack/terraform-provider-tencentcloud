## Context

The TencentCloud Terraform Provider currently supports CDB MySQL instance management but lacks a dedicated attachment resource for managing manual backups. The CDB API provides three interfaces: `CreateBackup` (create manual backup), `DescribeBackups` (query backup list), and `DeleteBackup` (delete a backup). These enable a CRD lifecycle pattern suitable for an attachment resource.

Existing CDB resources under `tencentcloud/services/cdb/` include `resource_tc_mysql_backup_policy.go`, `resource_tc_mysql_backup_encryption_status.go`, and `resource_tc_mysql_backup_download_restriction.go`, but none manage the backup lifecycle itself. The attachment resource pattern is well-established in the codebase (e.g., `resource_tc_mysql_security_groups_attachment.go`, `resource_tc_mysql_cls_log_attachment.go`).

## Goals / Non-Goals

**Goals:**
- Implement a Terraform resource `tencentcloud_mysql_backup` that creates, reads, and deletes manual MySQL backups
- Use `CreateBackup` API for creation, `DescribeBackups` API for reading, `DeleteBackup` API for deletion
- Use composite ID format (`backup_id` + `tccommon.FILED_SP` + `instance_id`) for resource identity
- Support all CreateBackup parameters: instance_id, backup_method, backup_db_table_list, manual_backup_name, encryption_flag
- Follow existing CDB attachment resource patterns and conventions
- Provide import support via composite ID

**Non-Goals:**
- Modify or manage automatic backups (those are managed by backup policy)
- Support backup restoration (separate concern)
- Support backup download URL generation
- Support backup method or content updates (backup is immutable once created)

## Decisions

**Decision 1: RESOURCE_KIND_ATTACHMENT with CRD-only lifecycle**
- Rationale: Backups are immutable once created (no Update API exists). All schema fields use `ForceNew: true` so any change triggers destroy-and-recreate. This aligns with the attachment pattern where only Create/Read/Delete are implemented.
- Alternative considered: RESOURCE_KIND_GENERAL with Update as no-op. Rejected because it would mislead users into thinking updates are possible.

**Decision 2: Composite ID (`backup_id#instance_id`)**
- Rationale: `backup_id` alone is insufficient because different instances could theoretically have the same backup ID. Using both `backup_id` and `instance_id` ensures globally unique resource identity. This follows the pattern of `resource_tc_mysql_security_groups_attachment.go` which uses `security_group_id#instance_id`.
- Alternative considered: Using only `backup_id` as ID. Rejected due to potential collision across instances.

**Decision 3: Read via DescribeBackups with instance_id filter and backup_id match**
- Rationale: `DescribeBackups` is a list API that returns all backups for an instance. We query by `instance_id` and iterate to find the matching `backup_id`. The `Limit` parameter is set to the API maximum (1000) to minimize pagination. If the backup is not found (deleted externally), we call `d.SetId("")` to remove from state.
- Alternative considered: Direct `DescribeBackup` API. Rejected because no such single-backup query API exists.

**Decision 4: backup_db_table_list as TypeList of maps**
- Rationale: `BackupItem` in the SDK has `Db` (string) and `Table` (string). The Terraform schema uses `TypeList` with `TypeMap` elements containing "db" and "table" keys. This provides a natural HCL syntax for users.
- Alternative considered: TypeSet. Rejected because order may matter for the API.

**Decision 5: Non-async design (no polling after CreateBackup)**
- Rationale: `CreateBackup` is a synchronous API that returns `BackupId` immediately. The DescribeBackups API can be used to check backup status afterward, but no polling is needed in the Create path since the resource ID is available right away.
- Alternative considered: Polling until backup status is SUCCESS. Rejected because it would unnecessarily delay Terraform apply and the backup ID is valid immediately.

**Decision 6: Read timeout retry pattern**
- Rationale: Use `tccommon.ReadRetryTimeout` for the DescribeBackups call in the Read function, following the standard pattern used across the provider.

## Risks / Trade-offs

- **Risk**: `DescribeBackups` returns many backups for instances with long backup history, requiring pagination handling. → **Mitigation**: Set `Limit` to 1000 (API maximum per documentation) to minimize paging. For instances exceeding 1000 backups, implement pagination loop.

- **Risk**: Backup creation may fail if instance is in an invalid state. → **Mitigation**: Standard Terraform error propagation will surface API errors to the user.

- **Risk**: The `BackupId` returned by CreateBackup is `uint64` while `DescribeBackups.Items[].BackupId` is `int64`. → **Mitigation**: Convert types appropriately using `helper.Int64ToStr`/`helper.StrToInt64` helpers when comparing.

- **Risk**: backup_db_table_list parameter only valid when backup_method is "logical". → **Mitigation**: Document this constraint clearly. The API will reject invalid combinations and surface errors.