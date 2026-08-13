## 1. Resource Implementation

- [x] 1.1 Create resource file `tencentcloud/services/postgresql/resource_tc_postgres_audit_log_file_attachment.go` with schema definition (Required+ForceNew: instance_id, start_time, end_time, product; Optional+ForceNew: filter block; Computed: file_name, status, file_size, create_time, download_url, err_msg, progress, finish_time; Timeouts: Create, Delete)
- [x] 1.2 Implement Create function: call CreateAuditLogFile API with retry, then poll DescribeAuditLogFiles until file status is `success`, set composite ID (instance_id#file_name)
- [x] 1.3 Implement Read function: parse composite ID, call DescribeAuditLogFiles with FileName filter, populate computed attributes (check nil before setting), remove from state if not found
- [x] 1.4 Implement Delete function: parse composite ID, call DeleteAuditLogFile API with retry

## 2. Provider Registration

- [x] 2.1 Register `tencentcloud_postgres_audit_log_file` resource in `tencentcloud/provider.go`
- [x] 2.2 Add resource entry in `tencentcloud/provider.md`

## 3. Documentation

- [x] 3.1 Create example documentation file `tencentcloud/services/postgresql/resource_tc_postgres_audit_log_file.md` with Example Usage and Import sections

## 4. Unit Tests

- [x] 4.1 Create test file `tencentcloud/services/postgresql/resource_tc_postgres_audit_log_file_attachment_test.go` with gomonkey-based unit tests for Create, Read, and Delete functions
- [x] 4.2 Run unit tests with `go test -gcflags=all=-l` to verify they pass
