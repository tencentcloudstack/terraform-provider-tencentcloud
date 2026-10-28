## Why

TencentCloud MongoDB supports audit log file management through cloud APIs (CreateAuditLogFile, DescribeAuditLogFiles, DeleteAuditLogFile), but there is currently no Terraform resource to manage these audit log files. Users need a Terraform resource to create, query, and delete MongoDB audit log files as part of their infrastructure-as-code workflow.

## What Changes

- Add a new Terraform resource `tencentcloud_mongodb_audit_log_file` that manages the lifecycle of MongoDB audit log files.
- The resource supports Create (CreateAuditLogFile), Read (DescribeAuditLogFiles), and Delete (DeleteAuditLogFile) operations.
- No Update API is available; changes to input parameters will force resource recreation.
- The resource uses a composite ID of `instance_id` and `file_name` joined by `tccommon.FILED_SP`.

## Capabilities

### New Capabilities
- `mongodb-audit-log-file`: Terraform resource for creating, reading, and deleting MongoDB audit log files with filtering and sorting options.

### Modified Capabilities

## Impact

- New resource file: `tencentcloud/services/mongodb/resource_tc_mongodb_audit_log_file.go`
- New test file: `tencentcloud/services/mongodb/resource_tc_mongodb_audit_log_file_test.go`
- New documentation: `tencentcloud/services/mongodb/resource_tc_mongodb_audit_log_file.md`
- Modified: `tencentcloud/provider.go` (register new resource)
- Modified: `tencentcloud/provider.md` (add resource entry)
- Dependencies: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mongodb/v20190725` (already vendored)
