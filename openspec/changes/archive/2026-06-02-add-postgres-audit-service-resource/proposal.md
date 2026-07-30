## Why

TencentCloud PostgreSQL supports database audit services that allow users to enable, configure, and disable audit logging for their instances. Currently, there is no Terraform resource to manage this audit service lifecycle. Users need a declarative way to enable and manage PostgreSQL audit services through Terraform.

## What Changes

- Add a new Terraform resource `tencentcloud_postgres_audit_service` that manages the full lifecycle (CRUD) of PostgreSQL database audit services.
- The resource uses `OpenAuditService` to create (enable), `DescribeAuditInstanceList` to read, `ModifyAuditService` to update, and `CloseAuditService` to delete (disable) the audit service.
- Register the new resource in `tencentcloud/provider.go` and `tencentcloud/provider.md`.
- Add documentation for the new resource.

## Capabilities

### New Capabilities
- `postgres-audit-service`: Terraform resource to manage PostgreSQL database audit service lifecycle including enabling, reading status, modifying configuration, and disabling audit.

### Modified Capabilities

## Impact

- New file: `tencentcloud/services/postgresql/resource_tc_postgres_audit_service.go`
- New file: `tencentcloud/services/postgresql/resource_tc_postgres_audit_service_test.go`
- New file: `tencentcloud/services/postgresql/resource_tc_postgres_audit_service.md`
- Modified: `tencentcloud/provider.go` (register new resource)
- Modified: `tencentcloud/provider.md` (add resource entry)
- Dependencies: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312` (already vendored)
