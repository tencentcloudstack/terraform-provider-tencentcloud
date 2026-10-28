## Context

腾讯云 EdgeOne（TEO）日志分析模块提供日志下载任务能力：用户指定站点（ZoneId）、数据归属地区（Area）、时间范围（StartTime/EndTime）等条件后，云端异步生成可下载的日志文件并返回下载链接。

vendored SDK（`teo v1.3.175`）已暴露以下 API：
- `CreateLogAnalysisDownloadTaskWithContext` → 返回 `{ TaskId }`（同步返回任务 ID）
- `DescribeLogAnalysisDownloadTasksWithContext` → 返回 `{ TotalCount, Tasks []*LogAnalysisDownloadTask }`，支持 `Filters`（`task-id`）、`Limit`（最大 100）、`Offset`

API 关键约束：
- **无 Delete 接口**：SDK 中不存在 `DeleteLogAnalysisDownloadTask`，下载任务在创建成功后保留 3 天并自动过期。
- **无 Modify 接口**：SDK 中不存在 `ModifyLogAnalysisDownloadTask`，任务参数不可变更。
- `CreateLogAnalysisDownloadTaskRequest` 入参：`ZoneId`、`Area`、`StartTime`、`EndTime`、`LogType`、`Condition`、`Format`、`Sort`，全部为扁平 string 字段。
- `DescribeLogAnalysisDownloadTasksRequest` 入参：`ZoneId`、`Area`、`LogType`、`Filters`、`Limit`、`Offset`。
- `LogAnalysisDownloadTask` 出参包含：`TaskId`、`ZoneId`、`Area`、`StartTime`、`EndTime`、`LogType`、`Condition`、`Format`、`Sort`、`Status`、`CreateTime`、`Url`、`ExpireTime`；`Status` 取值 `loading` / `failed` / `completed`，`Url` 仅在 `Status = completed` 时返回。

provider 已有 framework provider（与 SDKv2 provider mux 共存）与 `tencentcloud/framework/fw` 中的 `ActionWithConfigure`，TEO 也已有 framework action 先例 `tencentcloud_teo_confirm_origin_acl_update`；服务侧 `TeoService` 可直接复用。

## Goals / Non-Goals

**Goals:**
- 以 terraform-plugin-framework Action 形式提供 `tencentcloud_teo_log_analysis_download_task`，通过 `action` 块一次性创建日志分析下载任务。
- action config 字段与 `CreateLogAnalysisDownloadTaskRequest` 入参一一对应，无合并/重命名、不引入任何云 API 不支持的开关。
- 创建成功后轮询任务状态，把 `TaskId` 与最终下载链接 `Url` 通过 `SendProgress` 反馈给用户；任务 `failed` 时返回 error diagnostic。
- 代码风格参考已有 framework action（`action_tc_teo_confirm_origin_acl_update.go`、`action_tc_bdrc_run_copy_pair_tasks.go`）：`fw.ActionWithConfigure`、`Metadata`/`Schema`/`Invoke`、业务逻辑下沉到 service。
- 单元测试使用 gomonkey mock 云 API，覆盖 Invoke 成功 / 任务失败 / API 错误 / 参数校验 / 客户端未配置。

**Non-Goals:**
- 不新增 resource 或 data source：任务不可删除、不可修改、无 Terraform 生命周期可管理，resource 语义不成立（`terraform destroy` 无法删除任务，`Update` 只能报错要求重建）。
- 不下发日志文件内容；action 只产出云端下载链接（`url`）。
- 不提供后台异步执行/定时触发：action 同步等待打包完成，超时控制由 `tccommon.ReadRetryTimeout` 决定。
- 不新增任何 extension 兜底逻辑（云 API 无删除能力）。

## Decisions

### D1. 使用 framework Action 而非 SDKv2 Resource

**选择**：在 `tencentcloud/services/teo/` 下新增 `action_tc_teo_log_analysis_download_task.go`，实现 `action.Action` + `action.ActionWithConfigure`，并在 `tencentcloud/framework/registry.go` 的 `actionFactories` 中注册。

**备选**：SDKv2 resource（`ResourcesMap` 注册，Delete no-op、Update 全字段 `immutableArgs` 报错）。

**理由**：
- 云 API 只有 Create / Describe，没有 Delete / Modify，resource 只能靠“Delete no-op + 全字段 ForceNew”模拟生命周期，用户 `terraform destroy` 后云端仍有任务、`state` 还会因为任务 3 天后过期而被清空，语义混乱。
- Terraform Action 的定位正是“一次性操作”，与 `CreateLogAnalysisDownloadTask` 完全匹配；provider 已有两个 framework action 先例，`ActionWithConfigure` 已能拿到共享客户端。
- framework action 与 SDKv2 resource 通过 mux 共存，注册方式集中在 `framework/registry.go`，无需改 `tencentcloud/provider.go` 的 `ResourcesMap`。

### D2. config 字段与 Create 入参一一对应

**选择**：`zone_id`、`area`、`start_time`、`end_time` 为 Required；`log_type`、`condition`、`format`、`sort` 为 Optional（为空则不传给云 API，由云端按默认值 `l7-access-logs` / `csv` / `desc` 处理）。

**理由**：不引入派生字段（如 `wait` 开关），保证 schema 与云 API 能力严格一致；Optional 字段只有在用户显式配置时才出现在请求体中。

### D3. 枚举值在 Invoke 中校验

**选择**：`area`（`mainland` / `overseas`）、`log_type`（`l7-access-logs` / `web-attack`）、`format`（`csv`）、`sort`（`asc` / `desc`）在 Invoke 开头做白名单校验，非法值直接 `AddError` 返回，不发起云 API 调用。

**理由**：provider 未 vendored `terraform-plugin-framework-validators`，无法在 schema 层使用 `stringvalidator.OneOf`；在 Invoke 前置校验等价且能给出更明确的错误信息。必填字段同时校验 null / unknown / 空串。

### D4. 业务逻辑下沉到 TeoService

**选择**：新增三个 service 方法：
- `CreateTeoLogAnalysisDownloadTask(ctx, request)`：`resource.Retry(tccommon.WriteRetryTimeout, ...)` 包裹 `CreateLogAnalysisDownloadTaskWithContext`，retry 块内检查 `result == nil || result.Response == nil` 返回 `NonRetryableError`；retry 块外检查 `TaskId` 非空（nil / `""` 均报错）。返回 `TaskId`。
- `DescribeTeoLogAnalysisDownloadTaskById(ctx, zoneId, area, taskId)`：请求与 `Filters`（`Name: "task-id"`）在分页循环外构建，仅循环内变更 `Offset`/`Limit`（`Limit=100`，API 注释最大值）；命中 `*task.TaskId == taskId` 立即返回，未命中且结果已取完返回 `(nil, nil)`。
- `WaitTeoLogAnalysisDownloadTaskById(ctx, timeout, zoneId, area, taskId)`：普通轮询循环（间隔 3s，受 `ctx` 与 deadline 双重约束），直到 `Status != loading`。

**理由**：
- 与 provider 现有 service 分层一致（action 只做 config 解析、校验与 diagnostic 组装）。
- 轮询不使用嵌套 `resource.Retry`：`Describe` 内部已有 retry，`Wait` 只用普通循环 + `time.After`，避免“retry 中再 retry”导致超时时间叠加。
- 轮询而不是单次查询：`Url` 只在 `Status = completed` 时返回，单次查询几乎必然拿到 `loading`，无法给用户任何可用信息。

### D5. Invoke 的失败语义

**选择**：
- `CreateLogAnalysisDownloadTask` 失败 → 返回 error diagnostic（任务未创建）。
- 等待阶段超时 / 查询失败 → **不报错**，只 `SendProgress` 提示“任务已创建但最终状态未确认”，并提示用 `DescribeLogAnalysisDownloadTasks` 查询下载链接。
- `Status = failed` → 返回 error diagnostic（任务已创建但打包失败）。

**理由**：任务已创建后让 action 失败会误导用户重复 `apply` 并创建重复任务；而 `failed` 是业务层面的真实失败，需要显式告警。

### D6. 结果输出方式

**选择**：仅通过 `resp.SendProgress(action.InvokeProgressEvent{...})` 输出三条进度：开始创建、创建成功（含 `TaskId`）、最终结果（含 `Url`）。所有 `SendProgress` 调用前判断非 nil（单测中手工构造 `InvokeResponse` 时该函数为空）。

**理由**：framework action 的 `InvokeResponse` 不承载 state/output，`SendProgress` 是向用户回传信息的唯一通道；`Url` 无法写入 state，因此在进度消息中完整给出。

### D7. TEO `TaskId` 无复合 ID / 无 import

**选择**：action 不落 state，因此不涉及资源 ID、`import`、`d.SetId` 等逻辑；`TaskId` 仅作为日志与进度消息中的定位信息。

**理由**：action 无 state 概念，`DescribeLogAnalysisDownloadTasks` 所需的 `ZoneId`/`Area` 直接来自 config。

## Risks / Trade-offs

- **[Risk]** action 同步等待打包完成，最坏情况阻塞 `terraform apply` 约 `tccommon.ReadRetryTimeout`（默认 3 分钟，可通过 `TENCENTCLOUD_READ_RETRY_TIMEOUT` 调整）→ **Mitigation**：超时只输出进度提示、不使 action 失败，用户可改小超时；任务本身已在云端创建成功。
- **[Risk]** 下载链接仅在进度消息中出现，不会写入 state，无法用于其他资源的引用 → **Mitigation**：这是 action 语义的固有限制；文档中说明如需在 Terraform 中引用下载链接，应改用数据源或 API 查询。
- **[Risk]** 任务 3 天后自动过期，过期后链接不可用 → **Mitigation**：文档说明任务保留期，用户需在有效期内下载文件。
- **[Trade-off]** 未提供 `condition` 的 JSON 语法校验 → **Trade-off**：`condition` 为最大 12KB 的自由格式匹配表达式，provider 无法离线校验，交由云 API 报错。

## Migration Plan

纯新增，无 state 迁移：
1. 合并 action 实现 + service 方法 + registry 注册 + provider.md 登记。
2. 发布后用户通过 `action "tencentcloud_teo_log_analysis_download_task" "example" { config { ... } }` + `terraform apply -invoke=...` 使用。
3. 回滚：删除新增文件与 `registry.go`/`provider.md` 中的登记行；无 state 变更需回退。

## Open Questions

- 无。SDK 已提供所需全部 API（`teo v1.3.175`），实现路径完全确定。
