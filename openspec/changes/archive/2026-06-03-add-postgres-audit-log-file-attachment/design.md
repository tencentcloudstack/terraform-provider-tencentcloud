## Context

TencentCloud PostgreSQL provides audit log file management through three cloud APIs:
- `CreateAuditLogFile`: Creates an audit log file asynchronously (returns only RequestId, no FileName)
- `DescribeAuditLogFiles`: Queries audit log files with pagination and optional FileName filter
- `DeleteAuditLogFile`: Deletes an audit log file by FileName

The existing provider already has PostgreSQL resources under `tencentcloud/services/postgresql/`. The SDK package `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312` is already vendored with all three APIs available.

Key constraint: `CreateAuditLogFile` is asynchronous - it does not return a FileName. After creation, we must poll `DescribeAuditLogFiles` to discover the newly created file and wait for its status to become `success`.

## Goals / Non-Goals

**Goals:**
- Provide a Terraform resource `tencentcloud_postgres_audit_log_file` to manage audit log file lifecycle (create, read, delete)
- Support filter conditions (AuditLogFilter) during file creation
- Handle the asynchronous nature of file creation by polling until completion
- Follow existing provider patterns for error handling, retry, and composite IDs

**Non-Goals:**
- No Update operation (API does not support it)
- No data source for listing audit log files (out of scope for this change)
- No download functionality (DownloadUrl is exposed as computed attribute only)

## Decisions

### 1. Composite ID Format
**Decision**: Use `instance_id + tccommon.FILED_SP + file_name` as the resource ID.
**Rationale**: Both `instance_id` and `file_name` are required to uniquely identify and operate on an audit log file (Delete and Describe both require them). This follows the existing provider pattern for composite IDs.

### 2. Handling Asynchronous Creation
**Decision**: After calling `CreateAuditLogFile`, poll `DescribeAuditLogFiles` (without FileName filter, sorted by creation time) to find the newly created file, then wait for its status to become `success`.
**Rationale**: The Create API returns no FileName. We must list all files for the instance and identify the new one. Since we know the creation time window, we can find the most recently created file. Once identified, we set the resource ID and continue polling until status is `success`.
**Alternative considered**: Using a fixed sleep - rejected because file creation time is unpredictable.

### 3. All Input Fields as ForceNew
**Decision**: All user-provided input fields (`instance_id`, `start_time`, `end_time`, `product`, `filter`) are marked as `ForceNew: true`.
**Rationale**: There is no Update API. Any change to input parameters requires destroying and recreating the resource. This is the standard Terraform pattern for immutable resources.

### 4. Schema Timeouts
**Decision**: Declare `Create` and `Delete` timeouts in the schema. Create timeout is used for polling the async file creation.
**Rationale**: File creation is asynchronous and may take variable time. Timeouts allow users to control how long to wait.

### 5. Read Implementation
**Decision**: In the Read function, call `DescribeAuditLogFiles` with `FileName` filter to retrieve the specific file's metadata.
**Rationale**: The API supports filtering by FileName, which efficiently retrieves a single file's status without pagination.

## Risks / Trade-offs

- [Risk] CreateAuditLogFile returns no identifier → We must discover the file by listing. If multiple files are created simultaneously, there's a small window for ambiguity.
  → Mitigation: Poll immediately after creation and match by the most recent file that wasn't present before.

- [Risk] File creation may fail asynchronously (status becomes `failed`) → Mitigation: Check status during polling; if `failed`, return error with ErrMsg from the API response.

- [Trade-off] ForceNew on all fields means any parameter change destroys the file → This is acceptable because audit log files are inherently immutable once created.
