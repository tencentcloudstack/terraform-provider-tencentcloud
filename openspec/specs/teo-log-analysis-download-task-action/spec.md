# teo-log-analysis-download-task-action Specification

## Purpose
TBD - created by archiving change add-teo-log-analysis-download-task-resource. Update Purpose after archive.
## Requirements
### Requirement: Action registration
The provider SHALL expose a terraform-plugin-framework action named `tencentcloud_teo_log_analysis_download_task` that creates a single Tencent Cloud EdgeOne log analysis download task per invocation. The action MUST be registered in `tencentcloud/framework/registry.go` by appending `teo.NewTeoLogAnalysisDownloadTask` to `actionFactories`, and MUST be listed under the TEO `Action` section of `tencentcloud/provider.md`.

#### Scenario: Action type is discoverable
- **WHEN** an operator runs `terraform apply -invoke=action.tencentcloud_teo_log_analysis_download_task.example` against a configuration that declares the action
- **THEN** Terraform resolves the type without an "unknown action" error and invokes the action.

#### Scenario: No SDKv2 resource is registered
- **WHEN** a developer inspects `tencentcloud/provider.go`
- **THEN** `tencentcloud_teo_log_analysis_download_task` does not appear in `ResourcesMap`, because the cloud API has no delete or update interface and the capability is modelled as an action.

#### Scenario: Provider compiles
- **WHEN** the codebase is built with `go build ./tencentcloud/...`
- **THEN** the build succeeds with no compilation errors related to the new action.

### Requirement: Config schema mirrors `CreateLogAnalysisDownloadTask`
The action config schema SHALL expose every input parameter accepted by the `CreateLogAnalysisDownloadTask` API as a top-level attribute, with no renaming, merging or derived toggles:
- `zone_id` (string, required): 站点 ID。
- `area` (string, required): 数据归属地区，取值 `mainland` / `overseas`。
- `start_time` (string, required): 开始时间，示例 `2020-04-29T00:00:00Z`。
- `end_time` (string, required): 结束时间，示例 `2020-04-30T00:00:00Z`，与开始时间跨度最大 31 天。
- `log_type` (string, optional): 日志类型，取值 `l7-access-logs` / `web-attack`，云端默认 `l7-access-logs`。
- `condition` (string, optional): 日志匹配条件，最大长度 12KB。
- `format` (string, optional): 文件格式，取值 `csv`，云端默认 `csv`。
- `sort` (string, optional): 原始日志时间排序，取值 `asc` / `desc`，云端默认 `desc`。

The action SHALL NOT expose any attribute that is not accepted by `CreateLogAnalysisDownloadTaskRequestParams`, and SHALL NOT expose computed attributes because framework actions have no state.

#### Scenario: Required SDK fields are present
- **WHEN** a developer inspects the action schema
- **THEN** every field declared in `teo.v20220901.CreateLogAnalysisDownloadTaskRequestParams` (ZoneId, Area, StartTime, EndTime, LogType, Condition, Format, Sort) appears in the schema with semantically equivalent typing.

#### Scenario: Optional fields are omitted when unset
- **WHEN** an optional attribute (`log_type`, `condition`, `format`, `sort`) is null or an empty string
- **THEN** the corresponding request field is left nil so the cloud API applies its default value.

### Requirement: Config validation before API call
`Invoke` SHALL validate the config before any cloud API call:
- `zone_id`, `area`, `start_time`, `end_time` MUST be non-null, known and non-empty; otherwise an error diagnostic naming the missing attribute (`Missing <name>`) is added and no API call is made.
- `area` MUST be one of `mainland`, `overseas`.
- `log_type`, when set, MUST be one of `l7-access-logs`, `web-attack`.
- `format`, when set, MUST be `csv`.
- `sort`, when set, MUST be one of `asc`, `desc`.
- If `a.Client()` is nil, an error diagnostic with summary `Provider not configured` is added and no API call is made.

#### Scenario: Missing required attribute
- **WHEN** `zone_id` is empty
- **THEN** `Invoke` adds an error diagnostic containing `Missing zone_id` and does not call `CreateLogAnalysisDownloadTask`.

#### Scenario: Invalid enum value
- **WHEN** `area` is `invalid-area` or `sort` is `unknown`
- **THEN** `Invoke` adds an error diagnostic containing `Invalid area` / `Invalid sort` and does not call `CreateLogAnalysisDownloadTask`.

#### Scenario: Provider not configured
- **WHEN** the action is invoked without a configured provider client
- **THEN** `Invoke` adds an error diagnostic containing `Provider not configured`.

### Requirement: Create via service with retry and nil defense
The action SHALL call `TeoService.CreateTeoLogAnalysisDownloadTask(ctx, request)`, which:
- Wraps `CreateLogAnalysisDownloadTaskWithContext` in `resource.Retry(tccommon.WriteRetryTimeout, ...)`.
- Classifies SDK errors through `tccommon.RetryError(e)`.
- Returns `resource.NonRetryableError` when `result == nil` or `result.Response == nil`.
- Emits a `[DEBUG]` log line with the request action, request body and response body on success.
- After the retry block, returns an error when `Response.TaskId` is nil or empty.

When the service returns an error, `Invoke` SHALL add an error diagnostic whose summary contains `creating TEO log analysis download task` and whose detail contains the SDK error text, and SHALL NOT query the task.

#### Scenario: Successful create
- **WHEN** `CreateLogAnalysisDownloadTask` returns a non-empty `TaskId`
- **THEN** the action proceeds to wait for the task result and the create request carries exactly the configured fields.

#### Scenario: Create API error
- **WHEN** `CreateLogAnalysisDownloadTask` returns an error
- **THEN** `Invoke` returns an error diagnostic containing the SDK error and does not call `DescribeLogAnalysisDownloadTasks`.

#### Scenario: Empty TaskId rejected
- **WHEN** `CreateLogAnalysisDownloadTask` returns a nil or empty `TaskId`
- **THEN** the service returns an explicit error rather than an empty task identifier.

### Requirement: Wait for the task result and report the download URL
After a successful create, `Invoke` SHALL call `TeoService.WaitTeoLogAnalysisDownloadTaskById(ctx, tccommon.ReadRetryTimeout, zoneId, area, taskId)`, which polls `DescribeTeoLogAnalysisDownloadTaskById` every 3 seconds until `Status != loading`, the deadline expires, or `ctx` is cancelled. It SHALL NOT wrap the describe call in a second `resource.Retry` block, because the describe helper already retries internally.

`DescribeTeoLogAnalysisDownloadTaskById` SHALL:
- Wrap `DescribeLogAnalysisDownloadTasksWithContext` in `resource.Retry(tccommon.ReadRetryTimeout, ...)` and return `NonRetryableError` on a nil response.
- Build the request and the `Filters` slice (`Name: "task-id"`, `Values: [taskId]`) outside the pagination loop; only `Offset` and `Limit` mutate per iteration.
- Use `Limit=100` (the documented maximum) and stop when the task is matched or all results are consumed.
- Return `(nil, nil)` when the task is absent.

`Invoke` SHALL report the outcome through `resp.SendProgress`:
- On wait success with a non-empty `Url`: a message containing the `TaskId` and the download URL.
- On wait success without a `Url`: a message stating that the task is completed but the URL was not returned.
- On wait error (timeout or API failure): a message stating that the task was created but its final status could not be confirmed, referencing the `DescribeLogAnalysisDownloadTasks` API.

`SendProgress` MUST be called only when `resp.SendProgress` is non-nil.

#### Scenario: Task completes successfully
- **WHEN** the described task reaches `Status = completed` with a `Url`
- **THEN** `Invoke` emits progress messages for creation start, creation success (with `TaskId`), and completion (with `Url`), and returns no diagnostics.

#### Scenario: Task fails to pack
- **WHEN** the described task reaches `Status = failed`
- **THEN** `Invoke` adds an error diagnostic whose summary contains the `TaskId` and `failed`.

#### Scenario: Wait did not finish in time
- **WHEN** the task is still `loading` when `tccommon.ReadRetryTimeout` expires
- **THEN** `Invoke` emits a progress message that the task was created but the final status could not be confirmed, and returns no error diagnostic so the created task is not re-created by a retry.

#### Scenario: Progress channel unavailable
- **WHEN** `resp.SendProgress` is nil (for example in a manually constructed `InvokeResponse`)
- **THEN** the action does not panic and continues to completion.

### Requirement: No state, no import, no delete
The action SHALL NOT persist any Terraform state, SHALL NOT declare an import path, and SHALL NOT call any delete or modify API (the cloud API provides neither).

#### Scenario: Destroy does not touch the task
- **WHEN** an operator runs `terraform destroy` after invoking the action
- **THEN** no cloud API is called for the download task, because the action has no state to destroy.

### Requirement: Documentation and tests
The change SHALL include:
- A markdown document at `tencentcloud/services/teo/action_tc_teo_log_analysis_download_task.md` containing a one-line description that mentions EdgeOne (TEO), an `Example Usage` section with an `action "tencentcloud_teo_log_analysis_download_task" "..." { config { ... } }` block, and a note that actions require Terraform 1.14+. The document SHALL NOT include `Argument Reference` or `Attribute Reference` sections (auto-generated).
- A unit-test file at `tencentcloud/services/teo/action_tc_teo_log_analysis_download_task_test.go` using gomonkey to mock the cloud API (`CreateLogAnalysisDownloadTaskWithContext`, `DescribeLogAnalysisDownloadTasksWithContext`) and covering at minimum: successful Invoke with progress messages, task-failed diagnostic, create API error, config validation (missing/invalid values) and nil-client guard. Terraform test suites SHALL NOT be used.

#### Scenario: Documentation present
- **WHEN** the change is merged
- **THEN** the markdown documentation file exists next to the Go source, mentions EdgeOne (TEO) in the one-line description, contains an HCL action example, and does not contain `Argument Reference` or `Attribute Reference` sections.

#### Scenario: Test file present
- **WHEN** the change is merged
- **THEN** the `_test.go` file declares unit tests using gomonkey mocks (not Terraform acceptance test suites) covering Invoke success, task failure, API error, validation, and the nil-client guard.

### Requirement: SDK constraint
The implementation SHALL NOT modify `go.mod`, `go.sum` or any file under `vendor/github.com/tencentcloud/tencentcloud-sdk-go/`. The required APIs (`CreateLogAnalysisDownloadTask`, `DescribeLogAnalysisDownloadTasks`) already exist in the vendored `teo v1.3.175` SDK.

#### Scenario: Vendored SDK is sufficient
- **WHEN** the implementation begins
- **THEN** `CreateLogAnalysisDownloadTask` and `DescribeLogAnalysisDownloadTasks` are confirmed present under `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901/` and no dependency change is made.

