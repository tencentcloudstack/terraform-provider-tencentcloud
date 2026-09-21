## Why

Tencent Cloud EdgeOne (`teo`) 的日志分析模块支持创建日志下载任务，用户可下载站点指定时间范围的访问日志或托管规则日志。当前 Terraform Provider 未覆盖该能力，用户只能在控制台手动创建下载任务，无法纳入声明式运维流程。腾讯云 TEO SDK 已提供 `CreateLogAnalysisDownloadTask` 与 `DescribeLogAnalysisDownloadTasks` 接口，可以通过新增资源 `tencentcloud_teo_log_analysis_download_task` 将日志下载任务的生命周期纳入 Terraform 管理。

## What Changes

- 新增资源 `tencentcloud_teo_log_analysis_download_task`，映射 `CreateLogAnalysisDownloadTask` 入参为 schema 字段：`zone_id`（必填，ForceNew）、`area`（必填，ForceNew）、`start_time`（必填，ForceNew）、`end_time`（必填，ForceNew）、`log_type`（可选，ForceNew）、`condition`（可选，ForceNew）、`format`（可选，ForceNew）、`sort`（可选，ForceNew）。
- Create 调用 `CreateLogAnalysisDownloadTask`，返回 `TaskId` 作为资源 ID；由于云端无 Delete/Modify 接口，资源为 CR-only：Delete 为 no-op（仅清除 state），Update 通过将全部顶层字段加入 `immutableArgs` 数组返回 error，任何字段变更均要求重建资源。
- Read 通过新增 service 方法 `DescribeTeoLogAnalysisDownloadTaskById` 封装 `DescribeLogAnalysisDownloadTasks`，使用 `Filters`（`task-id`）过滤并按 `Limit=100`（API 注释最大值）分页查询，回填所有字段及计算字段 `task_id`、`status`、`create_time`、`url`、`expire_time`。
- 资源 ID 为裸 `TaskId`（无复合 ID）；支持 `terraform import`。
- 在 `tencentcloud/provider.go` 与 `tencentcloud/provider.md` 中注册该资源。
- 编写资源文档 `resource_tc_teo_log_analysis_download_task.md`（含 Example Usage 与 Import 说明）。
- 编写单元测试 `resource_tc_teo_log_analysis_download_task_test.go`，使用 gomonkey mock 云 API 进行业务逻辑测试。

## Capabilities

### New Capabilities
- `teo-log-analysis-download-task-resource`: 管理 EdgeOne 日志分析下载任务的生命周期（创建 / 查询 / 导入），CR-only 资源，字段变更触发重建，Delete 为 no-op 清除本地 state。

### Modified Capabilities
<!-- 无现有 capability 变更，本次仅新增资源 -->

## Impact

- **新增代码**:
  - `tencentcloud/services/teo/resource_tc_teo_log_analysis_download_task.go`（schema + Create/Read/Update/Delete + 辅助函数，单文件）
  - `tencentcloud/services/teo/resource_tc_teo_log_analysis_download_task.md`（资源文档 + Import 示例）
  - `tencentcloud/services/teo/resource_tc_teo_log_analysis_download_task_test.go`（gomonkey mock 单元测试）
- **修改代码**:
  - `tencentcloud/services/teo/service_tencentcloud_teo.go`：新增 `DescribeTeoLogAnalysisDownloadTaskById` 方法
  - `tencentcloud/provider.go`：在 `ResourcesMap` 中注册 `tencentcloud_teo_log_analysis_download_task`
  - `tencentcloud/provider.md`：新增资源描述条目
- **APIs consumed**: `CreateLogAnalysisDownloadTask`、`DescribeLogAnalysisDownloadTasks`（均已 vendored 在 `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901/`）
- **No breaking change**: 纯新增，不修改任何已有资源 schema 或 state。
- **No SDK upgrade required**: 所需 API 均已在 vendored SDK 中。
