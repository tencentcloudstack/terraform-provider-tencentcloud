## Context

TencentCloud MongoDB provides audit log file management through three cloud APIs:
- `CreateAuditLogFile`: Creates an audit log file with time range, sorting, and filtering options. Returns a `FileName`.
- `DescribeAuditLogFiles`: Queries audit log files by `InstanceId` and optional `FileName`. Returns file details including status, size, download URL.
- `DeleteAuditLogFile`: Deletes an audit log file by `InstanceId` and `FileName`.

There is no Update API. The resource is immutable after creation. The SDK package `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mongodb/v20190725` is already vendored.

## Goals / Non-Goals

**Goals:**
- Provide a `tencentcloud_mongodb_audit_log_file` resource that manages the full CRD lifecycle of MongoDB audit log files.
- Support all CreateAuditLogFile parameters: instance_id, start_time, end_time, order, order_by, and filter (with nested sub-fields).
- Expose computed attributes from DescribeAuditLogFiles response (items with file details).
- Use composite ID (`instance_id` + `file_name`) with `tccommon.FILED_SP` separator.
- Follow existing provider patterns (retry with `tccommon.ReadRetryTimeout`, error wrapping with `tccommon.RetryError`).

**Non-Goals:**
- No Update operation (API does not support it). Immutable fields trigger ForceNew or return error.
- No data source for audit log files in this change.
- No async polling after Create (CreateAuditLogFile is synchronous, returns FileName immediately).

## Decisions

1. **Composite ID format**: Use `instance_id + tccommon.FILED_SP + file_name` as the resource ID. This allows proper Read and Delete operations which require both fields.
   - Alternative: Use only `file_name` as ID. Rejected because Delete and Describe both require `InstanceId`.

2. **No Update, use immutableArgs pattern**: Since there is no Update API, the resource's Update function will check if any mutable field changed and return an error. `instance_id` is marked ForceNew. All other input fields (start_time, end_time, order, order_by, filter) are added to `immutableArgs` in the Update function.
   - Alternative: Mark all fields ForceNew. Rejected because the project convention for CRD-only resources uses the immutableArgs pattern for non-ID fields.

3. **Filter as TypeList with MaxItems 1**: The `Filter` parameter in CreateAuditLogFile is a single struct (`*AuditLogFilter`). In Terraform schema, this maps to a `TypeList` with `MaxItems: 1` containing the nested fields (host, user, exec_time, affect_rows, atype, result, param).

4. **Read uses DescribeAuditLogFiles with FileName filter**: The Describe API supports filtering by `FileName`. We pass both `InstanceId` and `FileName` to get the specific file. If no items returned, mark resource as removed from state.

5. **Items as Computed TypeList**: The `items` attribute stores the full audit log file details returned by DescribeAuditLogFiles, including file_name, create_time, status, file_size, download_url, err_msg, progress_rate.

## Risks / Trade-offs

- [Risk] CreateAuditLogFile may return an empty FileName if the request fails silently → Mitigation: Check if response or FileName is nil/empty after Create call, return `NonRetryableError` if so.
- [Risk] Audit log file creation is asynchronous (file status goes from "creating" to "success"/"failed") → Mitigation: The resource will be created successfully once the API returns a FileName. The `items` computed attribute will reflect the current status on subsequent reads. Users can check the `status` field in `items`.
- [Trade-off] No Update support means any parameter change requires resource destruction and recreation for ForceNew fields, or returns an error for immutable fields. This is acceptable given the API constraints.
