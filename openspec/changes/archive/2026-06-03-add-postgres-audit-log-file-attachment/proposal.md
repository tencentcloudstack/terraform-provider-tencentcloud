## Why

TencentCloud PostgreSQL supports audit log file management through cloud APIs (CreateAuditLogFile, DescribeAuditLogFiles, DeleteAuditLogFile). Currently, there is no Terraform resource to manage the lifecycle of audit log files. Users need a Terraform resource to create, query, and delete audit log files for their PostgreSQL instances in an infrastructure-as-code workflow.

## What Changes

- Add a new Terraform resource `tencentcloud_postgres_audit_log_file` (RESOURCE_KIND_ATTACHMENT) that manages the creation and deletion of PostgreSQL audit log files.
- The resource supports:
  - Creating an audit log file with time range and optional filter conditions
  - Reading the audit log file status and metadata via DescribeAuditLogFiles
  - Deleting the audit log file via DeleteAuditLogFile
- The resource uses a composite ID of `instance_id` + `file_name` (separated by `tccommon.FILED_SP`)
- Since CreateAuditLogFile is asynchronous (returns no FileName), the Create function will poll DescribeAuditLogFiles until the file status becomes `success`
- Only CRD operations are supported (no Update API), so all input fields are ForceNew

## Capabilities

### New Capabilities
- `postgres-audit-log-file`: Terraform resource to create, read, and delete PostgreSQL audit log files with filter support

### Modified Capabilities

## Impact

- New file: `tencentcloud/services/postgresql/resource_tc_postgres_audit_log_file_attachment.go`
- New file: `tencentcloud/services/postgresql/resource_tc_postgres_audit_log_file_attachment_test.go`
- New file: `tencentcloud/services/postgresql/resource_tc_postgres_audit_log_file.md`
- Modified: `tencentcloud/provider.go` (register new resource)
- Modified: `tencentcloud/provider.md` (add resource entry)
- Service layer: may extend existing `tencentcloud/services/postgresql/service_tencentcloud_postgresql.go`
- Dependencies: uses existing `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312` package (already vendored)
