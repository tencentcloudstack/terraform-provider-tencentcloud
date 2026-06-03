## 1. Service Layer

- [x] 1.1 Add audit service helper methods to `tencentcloud/services/postgresql/service_tencentcloud_postgresql.go`: OpenAuditService, DescribeAuditInstanceList (filtered by InstanceId), ModifyAuditService, CloseAuditService. Each method wraps the corresponding SDK call with retry logic using `tccommon.ReadRetryTimeout` and `tccommon.RetryError()`.

## 2. Resource Implementation

- [x] 2.1 Create `tencentcloud/services/postgresql/resource_tc_postgres_audit_service.go` with schema definition including: instance_id (Required, ForceNew, String), log_expire_day (Required, Int), hot_log_expire_day (Required, Int), audit_type (Required, String), product (Optional, String, default "postgres"), and computed attributes (audit_status, cold_log_expire_day, hot_log_size, cold_log_size, create_time). Implement CRUD functions: Create calls OpenAuditService, Read calls DescribeAuditInstanceList with InstanceId filter, Update calls ModifyAuditService, Delete calls CloseAuditService. Support import by instance_id.

## 3. Provider Registration

- [x] 3.1 Register `tencentcloud_postgres_audit_service` resource in `tencentcloud/provider.go` and `tencentcloud/provider.md`.

## 4. Documentation

- [x] 4.1 Create `tencentcloud/services/postgresql/resource_tc_postgres_audit_service.md` with resource description, example usage, and import section.

## 5. Unit Tests

- [x] 5.1 Create `tencentcloud/services/postgresql/resource_tc_postgres_audit_service_test.go` with unit tests using gomonkey to mock cloud API calls. Test Create, Read, Update, and Delete operations. Run tests with `go test -gcflags=all=-l`.
