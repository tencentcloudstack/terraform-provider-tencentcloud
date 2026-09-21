## Context

腾讯云 EdgeOne（TEO）日志分析模块提供日志下载任务能力：用户指定站点（ZoneId）、数据归属地区（Area）、时间范围（StartTime/EndTime）等条件后，云端异步生成可下载的日志文件并返回下载链接。

vendored SDK `v20220901` 已暴露以下 API：
- `CreateLogAnalysisDownloadTaskWithContext` → 返回 `{ TaskId }`（同步创建，TaskId 为任务唯一标识）
- `DescribeLogAnalysisDownloadTasksWithContext` → 返回 `{ TotalCount, Tasks []*LogAnalysisDownloadTask }`，按 `Filters`（`task-id`）过滤，支持分页（`Limit` 最大 100、`Offset`）

API 关键约束：
- **无 Delete 接口**：SDK 中不存在 `DeleteLogAnalysisDownloadTask`，下载任务在云端会按时过期保留（任务创建成功后保留 3 天）。因此资源 Delete 为 no-op，仅清除本地 state。
- **无 Modify 接口**：SDK 中不存在 `ModifyLogAnalysisDownloadTask`。因此资源 Update 中任意字段变更均不可调用云 API，只能通过 TF 层重建资源。
- `CreateLogAnalysisDownloadTaskRequest` 入参：`ZoneId`、`Area`、`StartTime`、`EndTime`、`LogType`、`Condition`、`Format`、`Sort`，全部为扁平 string 字段。
- `DescribeLogAnalysisDownloadTasksRequest` 入参：`ZoneId`、`Area`、`LogType`、`Filters`、`Limit`、`Offset`。注意：查询时 `Area` 是必填项，必须与创建时一致才能命中。
- `LogAnalysisDownloadTask` 出参包含：`TaskId`、`ZoneId`、`Area`、`StartTime`、`EndTime`、`LogType`、`Condition`、`Format`、`Sort`、`Status`、`CreateTime`、`Url`、`ExpireTime`。

provider 已有 TEO 服务的连通性绑定 `UseTeoV20220901Client` 及 service helper 文件 `service_tencentcloud_teo.go`，可直接复用。

## Goals / Non-Goals

**Goals:**
- 提供 `tencentcloud_teo_log_analysis_download_task` 资源，通过 Terraform 创建 TEO 日志分析下载任务并查询回填状态。
- schema 字段与 `CreateLogAnalysisDownloadTaskRequest` 入参一一对应，无合并/重命名。
- 资源 ID 为裸 `TaskId`，支持 `terraform import`。
- CR-only 资源的正确处理：Delete no-op、Update 字段变更拦截（按规则将全部顶层字段加入 `immutableArgs`，变更返回 error 要求重建）。
- 代码风格参考 `tencentcloud_igtm_strategy`（单文件、retry 包装、防御性 nil 检查）。
- 单元测试使用 gomonkey mock 云 API，覆盖 Create / Read / Update / Delete 业务逻辑。

**Non-Goals:**
- 不为下载任务实现数据源（data source），本次仅新增资源。
- 不下载日志文件内容；下载任务只产出云端下载链接（`url`），实际文件下载由用户通过链接完成。
- 不轮询任务状态至 `completed`：任务为异步生成下载链接，`CreateLogAnalysisDownloadTask` 接口本身为同步返回 TaskId 并非异步任务模式（无 TaskId 轮询语义），资源 Create 后直接 Read 回填即可；下载链接在 `Status == completed` 时由 Read 自然刷新。
- 不新增 `DeleteLogAnalysisDownloadTask` 的 extension 兜底逻辑（云 API 未提供该能力）。

## Decisions

### D1. 资源 ID 为裸 TaskId（无复合 ID）

**选择**：`d.SetId(*response.Response.TaskId)`，仅使用 `TaskId` 作为资源 ID。

**备选**：使用 `ZoneId` + `TaskId` 组合的复合 ID（如 `zoneId#taskId`）。

**理由**：
- `TaskId` 全局唯一，`DescribeLogAnalysisDownloadTasks` 可仅凭 `Filters=[{Name:"task-id", Values:[taskId]}]` 命中，无需额外携带 ZoneId。
- 但 Read 接口需要 `ZoneId` 与 `Area` 作为必填参数才能查询。因此资源 ID 中需携带 `ZoneId` 与 `Area`，采用复合 ID：`ZoneId#Area#TaskId`（使用 `tccommon.FILED_SP` 分隔），在 Read/Update/Delete 中拆分出各字段。
- 复合 ID 同时满足 `terraform import` 的需求，import 时需提供 `ZoneId#Area#TaskId` 形式的联合 ID。

（修正）最终选择：**复合 ID = `ZoneId#Area#TaskId`**，理由是 `DescribeLogAnalysisDownloadTasks` 的 `ZoneId` 与 `Area` 均为必填入参，无法仅凭 TaskId 完成查询；import 也必须提供这些信息。

### D2. 字段全部 ForceNew

**选择**：所有 schema 字段（除 `task_id` 等 computed 字段）均设置 `ForceNew: true`。

**理由**：
- 云端无 Modify 接口，任何字段变更都无法通过 API 更新，只能重建。
- 设置 ForceNew 可让 Terraform 在检测到字段变更时自动走重建路径（先 Delete 再 Create），无需在 Update 函数中手动处理。
- 按照规则要求，CR-only 资源仍会在 Update 方法中将全部顶层字段加入 `immutableArgs` 数组：若发现 update 方法传入的参数在 `immutableArgs` 中，则返回 error，作为双保险（防止 ForceNew 未命中场景）。

### D3. Delete 为 no-op

**选择**：Delete 函数仅 `defer tccommon.LogElapsed(...)` 与 `defer tccommon.InconsistentCheck(...)`，直接 `return nil`。

**备选**：Delete 返回 error 提示"云端不支持删除"。

**理由**：
- 云端无删除接口，且下载任务会自动过期（保留 3 天）。让 Terraform Delete 成功清除本地 state 是最符合用户预期的行为（`terraform destroy` 不应因云端无删除接口而失败）。
- 与 provider 中已有的 no-op Delete 模式一致（参考 `resource_tc_teo_purge_task_operation.go`）。

### D4. Update 拦截：immutableArgs 返回 error

**选择**：Update 函数中定义 `immutableArgs` 为所有顶层入参字段，遍历检查 `d.HasChange(v)`，命中则返回 error，提示该字段不可变更需重建资源。

**理由**：
- 满足规则"若一个资源只有CRD接口，则只将Id()字段设置成ForceNew，并在资源update方法中将其余顶层字段加入immutableArgs数组"。
- 实际由于 D2 已对全部字段设置 ForceNew，Update 分支一般不会被触发；此为双保险。

### D5. service 方法 DescribeTeoLogAnalysisDownloadTaskById

**选择**：在 `service_tencentcloud_teo.go` 新增方法，封装 `DescribeLogAnalysisDownloadTasks`：
- 入参：`ctx, zoneId, area, taskId`
- 在 for 循环外构建 `request` 与 `Filters`，仅循环中变更 `Offset`/`Limit`。
- `Limit=100`（API 注释最大值），分页直到累计 `len(tasks) >= TotalCount`。
- 找到 `*task.TaskId == taskId` 的元素后返回；未找到返回 `(nil, nil)`。
- 每次分页 SDK 调用用 `tccommon.ReadRetryTimeout` 的 `resource.Retry` 包装。

**理由**：
- 遵循现有 teo service 方法的分页模式与命名规范。

### D6. retry 与错误处理

**选择**：
- Create 调用 `CreateLogAnalysisDownloadTaskWithContext` 包在 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 中。
- Read（service 方法中）调用 `DescribeLogAnalysisDownloadTasksWithContext` 包在 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 中。
- retry 块内检查 `result == nil || result.Response == nil`，返回 `resource.NonRetryableError`。
- Create 成功后检查 `response.Response.TaskId` 是否为空（nil 或 `""`），为空则返回 `NonRetryableError`；检查前打印 `logId` 与 `d.Id()` 便于排障。
- retry 块外（Create）：设置 `d.SetId(...)`，再调用 Read。

**理由**：遵循规则中 Create/Read 的检查要求与 retry 不嵌套要求。

### D7. Read 空响应处理

**选择**：在资源 Read 函数中，若 service 方法返回 `(nil, nil)`，先打印 `log.Printf("[CRUD] teo_log_analysis_download_task id=%s", d.Id())` 保留现场，再 `d.SetId("")`。

**理由**：遵循规则"若云API返回了空，先打印 log 保留现场，再 d.SetId("")"。

### D8. schema 字段不含无关 wrapper

**选择**：schema 为扁平结构，`zone_id`、`area`、`start_time`、`end_time`、`log_type`、`condition`、`format`、`sort` 为入参字段；`task_id`、`status`、`create_time`、`url`、`expire_time` 为 computed 字段。

**理由**：本资源为单实例资源（非列表型），入参与出参字段直接平铺到顶层，无需 wrapper 嵌套层。

## Risks / Trade-offs

- **[Risk]** 下载任务在云端自动过期（保留 3 天），过期后 Read 返回 null 触发 `d.SetId("")`，Terraform 会计划重建资源 → **Mitigation**：这是云端语义，资源文档中说明任务自动过期行为；用户如需长期保留下载结果应在 `url` 可用时即下载文件。
- **[Risk]** ForceNew 导致任何字段变更都触发重建（会创建一个新的下载任务，旧任务保留至过期）→ **Mitigation**：与云 API 能力一致（无 Modify 接口），无更优方案；文档说明此限制。
- **[Risk]** Delete no-op 不会真正删除云端任务 → **Mitigation**：任务自动过期，且无删除接口；no-op 是 provider 中已确立的模式。
- **[Trade-off]** 复合 ID `ZoneId#Area#TaskId` 使 import 需提供三个值而非单个 TaskId → **Trade-off**：`DescribeLogAnalysisDownloadTasks` 要求 ZoneId/Area 必填，复合 ID 是唯一可行方案；文档 import 示例中明确说明联合 ID 格式。

## Migration Plan

纯新增，无 state 迁移：
1. 合并新资源 + service 方法 + provider 注册。
2. 发布后用户通过添加 `resource "tencentcloud_teo_log_analysis_download_task" "example" { ... }` 配置使用。
3. 回滚：删除新增文件与 provider.go 注册行；无 state 变更需回退。

## Open Questions

- 无。SDK 已提供所需全部 API，实现路径完全确定。