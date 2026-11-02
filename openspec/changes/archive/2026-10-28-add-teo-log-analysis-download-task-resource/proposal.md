## Why

Tencent Cloud EdgeOne（`teo`）日志分析模块支持创建日志下载任务，用户可下载站点指定时间范围的访问日志或托管规则日志。当前 Terraform Provider 未覆盖该能力，用户只能在控制台手动创建下载任务，无法纳入声明式运维流程。腾讯云 TEO SDK 已提供 `CreateLogAnalysisDownloadTask` 与 `DescribeLogAnalysisDownloadTasks` 接口。

该能力**不适合用 resource 表达**：日志下载任务没有 Delete / Modify 接口，任务在云端自动过期（创建后保留 3 天），也没有任何需要 Terraform 持续维护的生命周期状态——它本质上是一次性操作（创建任务 → 等待打包 → 拿到下载链接）。因此本次以 **terraform-plugin-framework Action**（`action` 块）的形式提供 `tencentcloud_teo_log_analysis_download_task`，避免用 CR-only resource 假装管理一个不可变更、不可删除的对象。

## What Changes

- 新增 framework action `tencentcloud_teo_log_analysis_download_task`（`action_tc_teo_log_analysis_download_task.go`），映射 `CreateLogAnalysisDownloadTask` 入参为 action config 字段：`zone_id`（必填）、`area`（必填）、`start_time`（必填）、`end_time`（必填）、`log_type`（可选）、`condition`（可选）、`format`（可选）、`sort`（可选）。
- `Invoke` 调用 `CreateLogAnalysisDownloadTask` 拿到 `TaskId`，随后轮询 `DescribeLogAnalysisDownloadTasks`（`Filters` 使用 `task-id`，`Limit=100`）直到任务脱离 `loading` 状态，并通过 `SendProgress` 向 Terraform 输出任务 ID 与最终下载链接 `Url`。
- 新增 service 方法：
  - `CreateTeoLogAnalysisDownloadTask(ctx, request)`：封装 `CreateLogAnalysisDownloadTaskWithContext`，`WriteRetryTimeout` 重试，返回 `TaskId`。
  - `DescribeTeoLogAnalysisDownloadTaskById(ctx, zoneId, area, taskId)`：封装 `DescribeLogAnalysisDownloadTasksWithContext`，`ReadRetryTimeout` 重试 + 分页（`Limit=100`，API 注释最大值）。
  - `WaitTeoLogAnalysisDownloadTaskById(ctx, timeout, zoneId, area, taskId)`：轮询直到 `Status != loading`。
- 在 `tencentcloud/framework/registry.go` 的 `actionFactories` 中注册 `teo.NewTeoLogAnalysisDownloadTask`；在 `tencentcloud/provider.md` 的 TEO `Action` 小节登记该 action。
- 编写 action 文档 `action_tc_teo_log_analysis_download_task.md`（含 Example Usage）。
- 编写单元测试 `action_tc_teo_log_analysis_download_task_test.go`，使用 gomonkey mock 云 API 进行业务逻辑测试。
- 不修改 vendored SDK：所需 API 在 `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo v1.3.175` 中已存在。

## Capabilities

### New Capabilities
- `teo-log-analysis-download-task-action`: 通过 Terraform action 一次性创建 EdgeOne 日志分析下载任务，等待打包完成并输出下载链接。

### Modified Capabilities
<!-- 无现有 capability 变更，本次仅新增 action -->

## Impact

- **新增代码**:
  - `tencentcloud/services/teo/action_tc_teo_log_analysis_download_task.go`（action schema + Invoke + 参数校验辅助函数）
  - `tencentcloud/services/teo/action_tc_teo_log_analysis_download_task.md`（action 文档 + Example Usage）
  - `tencentcloud/services/teo/action_tc_teo_log_analysis_download_task_test.go`（gomonkey mock 单元测试）
- **修改代码**:
  - `tencentcloud/services/teo/service_tencentcloud_teo.go`：新增 `CreateTeoLogAnalysisDownloadTask`、`DescribeTeoLogAnalysisDownloadTaskById`、`WaitTeoLogAnalysisDownloadTaskById`
  - `tencentcloud/framework/registry.go`：`actionFactories` 注册 `teo.NewTeoLogAnalysisDownloadTask`
  - `tencentcloud/provider.md`：TEO `Action` 小节新增 `tencentcloud_teo_log_analysis_download_task`
- **APIs consumed**: `CreateLogAnalysisDownloadTask`、`DescribeLogAnalysisDownloadTasks`（均已 vendored 在 `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901/`）
- **No breaking change**: 纯新增，不修改任何已有 resource / data source / action。
- **No SDK upgrade required**: 所需 API 均已在 vendored SDK `teo v1.3.175` 中，`go.mod` / `vendor/` 无改动。
