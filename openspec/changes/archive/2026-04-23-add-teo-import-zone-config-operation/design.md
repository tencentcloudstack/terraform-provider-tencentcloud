## Context

TEO（TencentCloud EdgeOne）是腾讯云的边缘安全加速平台，支持站点配置管理。当前 Terraform Provider 中已有多个 TEO 操作资源（如 `tencentcloud_teo_check_cname_status_operation`、`tencentcloud_teo_identify_zone_operation`、`tencentcloud_teo_create_cls_index_operation` 等），均遵循 RESOURCE_KIND_OPERATION 模式。

本次新增 `tencentcloud_teo_import_zone_config` 操作资源，调用 `ImportZoneConfig` API 导入站点配置。该 API 为异步接口，返回 `TaskId`，需要通过 `DescribeZoneConfigImportResult` 轮询任务状态直到完成（Status 为 success 或 failure）。

参考实现：
- `resource_tc_teo_check_cname_status_operation.go`：TEO 操作资源标准模式
- `resource_tc_teo_create_cls_index_operation.go`：简化版 TEO 操作资源
- `resource_tc_redis_backup_operation.go`：异步操作轮询模式

## Goals / Non-Goals

**Goals:**
- 新增 `tencentcloud_teo_import_zone_config` 操作资源，支持通过 Terraform 导入 TEO 站点配置
- 正确处理异步接口：调用 ImportZoneConfig 后轮询 DescribeZoneConfigImportResult 直到任务完成
- 遵循现有 TEO 操作资源的代码风格和架构模式
- 提供完整的单元测试（使用 gomonkey mock）

**Non-Goals:**
- 不支持 Update 操作（一次性操作资源）
- 不支持资源导入（Import）
- 不修改已有 TEO 资源的行为

## Decisions

### 1. 资源 ID 格式
使用 `zone_id` 和 `task_id` 的联合 ID（以 `tccommon.FILED_SP` 分隔）作为资源 ID。

**理由**：zone_id 和 task_id 共同唯一标识一次导入操作，使用联合 ID 可以确保每次导入操作都有唯一的资源标识。DescribeZoneConfigImportResult 也需要 zone_id 和 task_id 进行查询。

### 2. 异步轮询策略
调用 ImportZoneConfig 获取 TaskId 后，使用 `resource.Retry` 轮询 DescribeZoneConfigImportResult，超时时间为 `6 * tccommon.ReadRetryTimeout`。

**理由**：参考 `resource_tc_redis_backup_operation.go` 的异步轮询模式。轮询条件：
- Status 为 "doing" → 继续轮询
- Status 为 "success" → 轮询成功
- Status 为 "failure" → 返回错误（包含 Message 信息）

### 3. Schema 设计
- `zone_id`（Required, ForceNew）：站点 ID
- `content`（Required, ForceNew）：待导入的配置内容（JSON 格式）
- `task_id`（Computed）：异步任务 ID
- `status`（Computed）：导入任务状态
- `message`（Computed）：状态提示信息
- `import_time`（Computed）：导入开始时间
- `finish_time`（Computed）：导入完成时间

**理由**：zone_id 和 content 是 ImportZoneConfig API 的必填入参，设为 Required 且 ForceNew（一次性操作无 Update）。task_id 为 API 返回值，设为 Computed。DescribeZoneConfigImportResult 的出参设为 Computed，在 Create 中轮询成功后设置。

### 4. API 客户端使用
使用 `UseTeoV20220901Client()` 方法获取 TEO 客户端，与 `resource_tc_teo_create_cls_index_operation.go` 保持一致。

**理由**：现有 TEO 操作资源均使用 `UseTeoV20220901Client()`，保持一致性。

## Risks / Trade-offs

- [异步轮询超时] → 使用 6 倍 ReadRetryTimeout 作为轮询上限，对于大型配置导入应足够
- [配置内容可能很大] → content 字段使用 schema.TypeString 存储，与 API 的字符串参数一致，Terraform 本身对字符串长度无硬限制
- [任务失败处理] → 轮询到 failure 状态时返回错误信息，包含 API 返回的 Message 字段
