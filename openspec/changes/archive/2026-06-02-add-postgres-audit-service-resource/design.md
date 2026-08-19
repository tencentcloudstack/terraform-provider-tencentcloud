## Context

TencentCloud PostgreSQL provides database audit service APIs that allow enabling, modifying, querying, and disabling audit logging on PostgreSQL instances. The cloud APIs are:
- `OpenAuditService`: Enable audit service on an instance
- `DescribeAuditInstanceList`: Query audit instance list (used for reading audit status)
- `ModifyAuditService`: Modify audit service configuration
- `CloseAuditService`: Disable audit service on an instance

All APIs are in the `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312` package which is already vendored. The existing PostgreSQL service layer is at `tencentcloud/services/postgresql/service_tencentcloud_postgresql.go`.

## Goals / Non-Goals

**Goals:**
- Provide a `tencentcloud_postgres_audit_service` resource that manages the full lifecycle of PostgreSQL audit service
- Support enabling audit with configurable log retention days, hot log retention days, and audit type
- Support modifying audit configuration (log_expire_day, hot_log_expire_day, audit_type)
- Support disabling audit service (resource deletion)
- Support importing existing audit service by instance_id
- Follow existing provider patterns (retry logic, error handling, service layer)

**Non-Goals:**
- Managing audit log content or querying audit logs
- Managing audit rules/policies (separate concern)
- Supporting the `DescribeAuditInstanceList` response fields (cold_log_expire_day, hot_log_size, cold_log_size, etc.) as writable attributes - these are computed/read-only

## Decisions

1. **Resource ID**: Use `instance_id` as the resource ID since the audit service is a per-instance configuration. No composite ID needed.

2. **CRUD Mapping**:
   - Create → `OpenAuditService`
   - Read → `DescribeAuditInstanceList` (filter by InstanceId, AuditSwitch=1)
   - Update → `ModifyAuditService`
   - Delete → `CloseAuditService`

3. **Schema Design**:
   - `instance_id` (Required, ForceNew, String): The PostgreSQL instance ID
   - `log_expire_day` (Required, Int): Log retention days (7/30/90/180/365/1095/1825)
   - `hot_log_expire_day` (Required, Int): Hot log retention days (7/30/90/180/365/1095/1825)
   - `audit_type` (Required, String): Audit type ("complex" or "simple")
   - `product` (Optional, String, default "postgres"): Product name, fixed to "postgres"
   - Computed attributes from Read: `audit_status`, `cold_log_expire_day`, `hot_log_size`, `cold_log_size`, `create_time`

4. **Read Implementation**: Use `DescribeAuditInstanceList` with `Product="postgres"`, `AuditSwitch=1`, and `Filters` with `InstanceId` filter to find the specific instance's audit info. If not found or AuditStatus is "OFF", treat as resource deleted.

5. **Retry Logic**: All API calls wrapped with `tccommon.ReadRetryTimeout` and `resource.Retry` pattern using `tccommon.RetryError()`.

6. **Service Layer**: Add helper methods to the existing PostgreSQL service file for the audit API calls.

## Risks / Trade-offs

- [Risk] `DescribeAuditInstanceList` is a list API, not a direct get-by-ID → Mitigation: Filter by InstanceId and iterate results to find the matching instance.
- [Risk] Audit service state may not be immediately consistent after Open/Modify/Close → Mitigation: Use retry with timeout on Read operations.
- [Risk] `product` field is always "postgres" but required by API → Mitigation: Make it Optional with default value "postgres" to reduce user burden while maintaining API compatibility.
