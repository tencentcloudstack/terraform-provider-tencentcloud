## ADDED Requirements

### Requirement: Resource registration
The provider SHALL expose a resource type named `tencentcloud_teo_log_analysis_download_task` that manages a single Tencent Cloud EdgeOne log analysis download task per resource block. The resource MUST be registered in `tencentcloud/provider.go` under the `teo` namespace.

#### Scenario: Resource type is discoverable
- **WHEN** an operator runs `terraform plan` against a configuration that references `resource "tencentcloud_teo_log_analysis_download_task" "<name>"`
- **THEN** Terraform resolves the type without an "unknown resource" error and shows the planned create.

#### Scenario: Provider compiles
- **WHEN** the codebase is built with `go build ./tencentcloud/...`
- **THEN** the build succeeds with no compilation errors related to the new resource.

### Requirement: Schema mirrors `CreateLogAnalysisDownloadTask`
The resource schema SHALL expose every input parameter accepted by the `CreateLogAnalysisDownloadTask` API as a top-level attribute, with no renaming or merging:
- `zone_id` (string, required, ForceNew): 站点 ID。
- `area` (string, required, ForceNew): 数据归属地区，取值 `mainland` / `overseas`。
- `start_time` (string, required, ForceNew): 开始时间，示例 `2020-04-29T00:00:00Z`。
- `end_time` (string, required, ForceNew): 结束时间，示例 `2020-04-30T00:00:00Z`。
- `log_type` (string, optional, ForceNew): 日志类型，取值 `l7-access-logs` / `web-attack`。
- `condition` (string, optional, ForceNew): 日志匹配条件，最大长度 12KB。
- `format` (string, optional, ForceNew): 文件格式，取值 `csv`。
- `sort` (string, optional, ForceNew): 原始日志时间排序，取值 `asc` / `desc`。

The resource SHALL additionally expose the following read-only computed attributes hydrated from `DescribeLogAnalysisDownloadTasks` response (`LogAnalysisDownloadTask`):
- `task_id` (string, computed) — also stored as part of the resource ID.
- `status` (string, computed): 任务状态，取值 `loading` / `failed` / `completed`。
- `create_time` (string, computed): 任务创建时间。
- `url` (string, computed): 下载地址，仅当 `status = completed` 时有返回值。
- `expire_time` (string, computed): 下载任务过期时间。

#### Scenario: Required SDK fields are present
- **WHEN** a developer inspects the resource schema
- **THEN** every field declared in `teo.v20220901.CreateLogAnalysisDownloadTaskRequestParams` (ZoneId, Area, StartTime, EndTime, LogType, Condition, Format, Sort) appears in the schema with semantically equivalent typing.

#### Scenario: No undocumented schema fields
- **WHEN** a developer inspects the resource schema
- **THEN** there are no input fields beyond those listed above; no derived flags or synthetic toggles are introduced.

### Requirement: Composite resource ID
The resource ID SHALL be the composite `ZoneId#Area#TaskId` (using `tccommon.FILED_SP` as the separator), because `DescribeLogAnalysisDownloadTasks` requires `ZoneId` and `Area` as mandatory inputs to locate the task. The resource SHALL support `terraform import` using the composite ID.

#### Scenario: Create sets the ID
- **WHEN** `CreateLogAnalysisDownloadTask` returns a non-nil, non-empty `TaskId`
- **THEN** the resource calls `d.SetId(zoneId + tccommon.FILED_SP + area + tccommon.FILED_SP + taskId)` after the SDK call succeeds.

#### Scenario: Import by composite ID
- **WHEN** an operator runs `terraform import tencentcloud_teo_log_analysis_download_task.example zone-xxx#mainland#task-yyy`
- **THEN** the resource state is hydrated from `DescribeLogAnalysisDownloadTasks` using the parsed ZoneId, Area, and TaskId, with no error.

#### Scenario: Broken ID rejected
- **WHEN** the resource ID does not contain exactly 3 segments when split by the field separator
- **THEN** the Read/Update/Delete function returns an error indicating the ID is broken.

### Requirement: Create with retry and nil defense
On Create, the resource SHALL invoke `CreateLogAnalysisDownloadTask`, wrapped in `resource.Retry(tccommon.WriteRetryTimeout, ...)`. Inside the retry block:
- The SDK call `CreateLogAnalysisDownloadTaskWithContext` is invoked; transient errors are retried via `tccommon.RetryError(e)`.
- If `result == nil` or `result.Response == nil`, the wrapper returns `resource.NonRetryableError` with a descriptive message.
- A `[DEBUG]` log line containing the request action, request body, and response body is emitted on success.

After the retry block, the resource SHALL:
- Print the current `logId` and `d.Id()` for troubleshooting.
- Check that `response.Response.TaskId` is neither nil nor empty string; if empty, return a `NonRetryableError` (using `fmt.Errorf`).
- Set the composite resource ID.
- Invoke Read to hydrate state.

#### Scenario: Successful create
- **WHEN** `CreateLogAnalysisDownloadTask` succeeds and returns a non-empty `TaskId`
- **THEN** the resource sets the composite ID, invokes Read, and returns no error.

#### Scenario: Empty TaskId rejected
- **WHEN** `CreateLogAnalysisDownloadTask` returns a nil or empty `TaskId`
- **THEN** the resource returns an explicit error rather than setting an empty ID.

#### Scenario: Nil response defense
- **WHEN** the SDK call returns a nil `Response`
- **THEN** the wrapper returns `resource.NonRetryableError` with a descriptive message rather than dereferencing the nil pointer.

### Requirement: Read with retry, pagination, and empty-response handling
On Read, the resource SHALL call `TeoService.DescribeTeoLogAnalysisDownloadTaskById(ctx, zoneId, area, taskId)`, which:
- Wraps the SDK call `DescribeLogAnalysisDownloadTasksWithContext` in `resource.Retry(tccommon.ReadRetryTimeout, ...)`.
- Builds the request and `Filters` slice (with `Name: "task-id"`, `Values: [taskId]`) outside the pagination loop; only `Offset` / `Limit` mutate per iteration.
- Iterates pages with `Limit=100` (the documented maximum) until all results are consumed.
- Strict-equals `*task.TaskId == taskId` before returning the matched task.
- Returns `(nil, nil)` when the task is absent.

When the helper returns `(nil, nil)`, the resource SHALL:
- Print `log.Printf("[CRUD] teo_log_analysis_download_task id=%s", d.Id())` to preserve context.
- Call `d.SetId("")` and return nil.

When the helper returns a task, the resource SHALL call `_ = d.Set(...)` for every input field and every computed field, each guarded by a nil check on the corresponding response field before calling `d.Set`.

#### Scenario: Resource present
- **WHEN** the helper finds a matching `LogAnalysisDownloadTask`
- **THEN** the resource populates all schema fields (input + computed) from the response.

#### Scenario: Resource removed externally
- **WHEN** the helper returns `(nil, nil)` (no matching task)
- **THEN** the resource prints a `[CRUD]` log line with the ID, calls `d.SetId("")`, and returns no error.

#### Scenario: Pagination request reuse
- **WHEN** the helper paginates through more than one page
- **THEN** a single `DescribeLogAnalysisDownloadTasksRequest` instance is reused across pages, with only `Offset` and `Limit` mutated.

#### Scenario: Nil field guard
- **WHEN** a response field (e.g. `Url`, `ExpireTime`) is nil
- **THEN** the resource skips calling `d.Set` for that field rather than setting an empty/zero value.

### Requirement: Update is immutable (CR-only resource)
The Update function SHALL define an `immutableArgs` array containing all top-level input schema fields (`zone_id`, `area`, `start_time`, `end_time`, `log_type`, `condition`, `format`, `sort`). When any of these fields has changed (`d.HasChange`), the function SHALL return an error instructing the operator to recreate the resource. No SDK Modify call SHALL be invoked because the cloud API provides no update interface.

#### Scenario: Field change rejected
- **WHEN** any input field changes between plan steps
- **THEN** the Update function returns an error indicating the field is immutable and the resource must be recreated.

#### Scenario: No change passes through
- **WHEN** no input field has changed
- **THEN** the Update function calls Read and returns no error.

### Requirement: Delete is a no-op
On Delete, the resource SHALL NOT invoke any cloud API (no `DeleteLogAnalysisDownloadTask` exists). The function SHALL only emit the standard `defer tccommon.LogElapsed(...)` and `defer tccommon.InconsistentCheck(...)` calls and return nil, clearing the resource from Terraform local state.

#### Scenario: Successful no-op delete
- **WHEN** an operator runs `terraform destroy`
- **THEN** the Delete function returns nil without calling any SDK API, and Terraform removes the resource from state.

### Requirement: Retry coverage and logging conventions
Every SDK call (`CreateLogAnalysisDownloadTaskWithContext`, `DescribeLogAnalysisDownloadTasksWithContext`) SHALL be invoked from inside a `resource.Retry` block. The retry budget is `tccommon.WriteRetryTimeout` for Create and `tccommon.ReadRetryTimeout` for Read.

The resource SHALL emit:
- `defer tccommon.LogElapsed("resource.tencentcloud_teo_log_analysis_download_task.<op>")()` at the top of every CRUD function.
- `defer tccommon.InconsistentCheck(d, meta)()` at the top of every CRUD function.
- A `[DEBUG]` line per SDK invocation containing the request action, request body, and response body.
- A `[CRITAL]%s ... failed, reason:%+v` line on every retry-block failure.
- A `[CRUD]` line when the resource is detected as deleted out of band during Read.

#### Scenario: Transient SDK error
- **WHEN** an SDK call returns a transient TencentCloud SDK error
- **THEN** the call is retried via `tccommon.RetryError(e)` until it succeeds or the retry budget is exhausted.

#### Scenario: Standard log lines emitted
- **WHEN** any CRUD operation runs
- **THEN** the operation's elapsed time is logged via `tccommon.LogElapsed` and inconsistency is checked via `tccommon.InconsistentCheck`.

### Requirement: Documentation and tests
The change SHALL include:
- A markdown document at `tencentcloud/services/teo/resource_tc_teo_log_analysis_download_task.md` containing a self-contained `resource "tencentcloud_teo_log_analysis_download_task" "..." { ... }` example and a `terraform import` example using the composite `ZoneId#Area#TaskId` ID. The one-line description SHALL mention the cloud product EdgeOne (TEO). The document SHALL NOT include `Argument Reference` or `Attribute Reference` sections (auto-generated).
- A unit-test file at `tencentcloud/services/teo/resource_tc_teo_log_analysis_download_task_test.go` using gomonkey to mock the cloud API (`CreateLogAnalysisDownloadTaskWithContext`, `DescribeLogAnalysisDownloadTasksWithContext`) and covering at minimum: successful Create + Read, Update-rejects-immutable-field, and no-op Delete. Terraform test suites SHALL NOT be used (gomonkey mock only).

#### Scenario: Documentation present
- **WHEN** the change is merged
- **THEN** the markdown documentation file exists, mentions EdgeOne (TEO) in the one-line description, contains an HCL example and an `import` example using the composite ID, and does not contain `Argument Reference` or `Attribute Reference` sections.

#### Scenario: Test file present
- **WHEN** the change is merged
- **THEN** the `_test.go` file declares unit tests using gomonkey mocks (not Terraform acceptance test suites) covering Create, Read, Update-reject, and Delete no-op scenarios.

### Requirement: SDK constraint
The implementation SHALL NOT modify any file under `vendor/github.com/tencentcloud/tencentcloud-sdk-go/`. If a required API is missing from the vendored SDK, the implementer MUST halt and request an SDK upgrade rather than authoring or patching SDK source.

#### Scenario: Vendored SDK is sufficient
- **WHEN** the implementation begins
- **THEN** `CreateLogAnalysisDownloadTask` and `DescribeLogAnalysisDownloadTasks` are confirmed present under `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901/` before any code is written.
