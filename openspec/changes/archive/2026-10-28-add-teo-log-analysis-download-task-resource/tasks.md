## 1. Service 层实现

- [x] 1.1 在 `tencentcloud/services/teo/service_tencentcloud_teo.go` 中新增 `CreateTeoLogAnalysisDownloadTask(ctx context.Context, request *teo.CreateLogAnalysisDownloadTaskRequest) (taskId string, errRet error)`：SDK 调用 `CreateLogAnalysisDownloadTaskWithContext` 包装在 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 中，错误经 `tccommon.RetryError(e)` 分类；`result == nil || result.Response == nil` 返回 `NonRetryableError`；成功打印 `[DEBUG]`；retry 块外校验 `Response.TaskId` 非 nil 且非空
- [x] 1.2 在 `tencentcloud/services/teo/service_tencentcloud_teo.go` 中新增 `DescribeTeoLogAnalysisDownloadTaskById(ctx, zoneId, area, taskId) (*teov20220901.LogAnalysisDownloadTask, error)`：在分页循环外构建 `DescribeLogAnalysisDownloadTasksRequest` 与 `Filters`（`Name: "task-id"`, `Values: [taskId]`），仅循环内变更 `Offset`/`Limit`；`Limit=100`（API 注释最大值）；严格相等 `*task.TaskId == taskId` 命中后立即返回；结果取完仍未命中返回 `(nil, nil)`；每次分页 SDK 调用包装在 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 中并防御性检查 `result == nil || result.Response == nil`
- [x] 1.3 新增 `WaitTeoLogAnalysisDownloadTaskById(ctx, timeout, zoneId, area, taskId)`：普通轮询循环（3s 间隔 + deadline + `ctx.Done()`），复用 `DescribeTeoLogAnalysisDownloadTaskById`（其内部已有 retry），不在 retry 块内再套 retry，直到 `Status != loading` 返回任务
- [x] 1.4 确认 `connectivity.UseTeoV20220901Client()` 可用，SDK 导入路径为 `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901`；确认所需 API 已在 vendored `teo v1.3.175` 中存在，`go.mod`/`go.sum`/`vendor/` 不做改动

## 2. Action 实现

- [x] 2.1 创建 `tencentcloud/services/teo/action_tc_teo_log_analysis_download_task.go`（单文件），实现 `NewTeoLogAnalysisDownloadTask() action.Action`、`TeoLogAnalysisDownloadTask`（内嵌 `fw.ActionWithConfigure`）与 `TeoLogAnalysisDownloadTaskModel`
- [x] 2.2 `Metadata` 返回 `tencentcloud_teo_log_analysis_download_task`；`Schema` 声明 `zone_id`/`area`/`start_time`/`end_time`（Required）与 `log_type`/`condition`/`format`/`sort`（Optional），字段与 `CreateLogAnalysisDownloadTaskRequestParams` 一一对应，不引入派生开关与 computed 字段
- [x] 2.3 `Invoke` 解析 config，校验必填字段（null/unknown/空串均报 `Missing <name>`）、枚举字段（`area`、`log_type`、`format`、`sort` 白名单）与 `a.Client() == nil`（`Provider not configured`）；校验失败直接返回，不发起云 API 调用
- [x] 2.4 `Invoke` 构建 `CreateLogAnalysisDownloadTaskRequest`（Optional 字段为空时不赋值，交由云端默认值），调用 `service.CreateTeoLogAnalysisDownloadTask`；失败时 `AddError` 且 summary 包含 `creating TEO log analysis download task`
- [x] 2.5 创建成功后通过 `SendProgress` 输出任务 ID，调用 `service.WaitTeoLogAnalysisDownloadTaskById(ctx, tccommon.ReadRetryTimeout, ...)`；等待失败只输出进度提示（不报错），`Status = failed` 输出 error diagnostic，`completed` 时输出含 `Url` 的进度消息
- [x] 2.6 封装 `teoLogAnalysisDownloadTaskSendProgress`，在 `resp.SendProgress == nil` 时直接返回，保证单测中手工构造的 `InvokeResponse` 不 panic

## 3. Action 注册

- [x] 3.1 在 `tencentcloud/framework/registry.go` 的 `actionFactories` 中新增 `teo.NewTeoLogAnalysisDownloadTask`（紧跟 `teo.NewTeoConfirmOriginAclUpdate`，按产品分组）
- [x] 3.2 在 `tencentcloud/provider.md` 的 TEO `Action` 小节新增 `tencentcloud_teo_log_analysis_download_task`；`tencentcloud/provider.go` 的 `ResourcesMap` 不注册该名称（无 Delete/Modify 接口，不用 resource 表达）

## 4. 文档示例

- [x] 4.1 创建 `tencentcloud/services/teo/action_tc_teo_log_analysis_download_task.md`：一句话描述提及 EdgeOne（TEO），说明任务保留 3 天且不可删除/修改；包含 `~> NOTE` Terraform 1.14+ 提示；`Example Usage` 使用 `action "tencentcloud_teo_log_analysis_download_task" "example" { config { ... } }`；不要添加 `Argument Reference` 和 `Attribute Reference` 部分

## 5. 单元测试

- [x] 5.1 创建 `tencentcloud/services/teo/action_tc_teo_log_analysis_download_task_test.go`（`package teo_test`），使用 gomonkey mock `UseTeoV20220901Client`、`CreateLogAnalysisDownloadTaskWithContext`、`DescribeLogAnalysisDownloadTasksWithContext`，不使用 terraform 测试套件
- [x] 5.2 `TestTeoLogAnalysisDownloadTask_MetadataAndSchema`：校验 type name 与 required/optional 属性集合
- [x] 5.3 `TestTeoLogAnalysisDownloadTask_Invoke_Success`：mock Create 返回 TaskId、Describe 返回 `completed` + `Url`，校验请求字段（含 `Filters` 的 `task-id`）、三条进度消息与无 error diagnostic
- [x] 5.4 `TestTeoLogAnalysisDownloadTask_Invoke_TaskFailed`：Describe 返回 `failed`，校验 error diagnostic 含 TaskId 与 failed
- [x] 5.5 `TestTeoLogAnalysisDownloadTask_Invoke_CreateAPIError`：Create 返回错误，校验 error diagnostic 且未调用 Describe
- [x] 5.6 `TestTeoLogAnalysisDownloadTask_Invoke_MissingInput`：空 `zone_id`、非法 `area`、非法 `sort` 三种场景均报错且不调用云 API
- [x] 5.7 `TestTeoLogAnalysisDownloadTask_Invoke_ClientNotConfigured`：未注入 client 时报 `Provider not configured`

## 6. 验证

- [x] 6.1 确认代码可编译（构建由其他流程执行，本任务仅确认代码生成正确）
- [x] 6.2 确认所有函数返回的 error 已检查；必定不出错的函数用 `_ = func()` 处理
- [x] 6.3 确认 Create 调用的入参在 `CreateLogAnalysisDownloadTask` 中存在，Read 调用的入参在 `DescribeLogAnalysisDownloadTasks` 中存在（参数一致性检查）
- [x] 6.4 确认未修改 `vendor/`、`go.mod`、`go.sum`；未手写 `website/` 下文件（由 `make doc` 在收尾阶段生成 `website/docs/actions/teo_log_analysis_download_task.html.markdown`）
