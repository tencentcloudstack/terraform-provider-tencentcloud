## Context

腾讯云 CLS 服务已提供 Splunk 投递相关的云 API（`CreateSplunkDeliver`、`DescribeSplunkDelivers`、`ModifySplunkDeliver`、`DeleteSplunkDeliver`），且对应的 Go SDK 已存在于 vendor 中。当前 Terraform Provider 尚未封装这些 API，需要新增 `tencentcloud_cls_splunk_deliver` 资源。

该资源属于 RESOURCE_KIND_GENERAL 类型，需要实现完整的 CRUD + Import 能力。

## Goals / Non-Goals

**Goals:**
- 提供 `tencentcloud_cls_splunk_deliver` 资源，支持 Splunk 投递任务的创建、读取、更新、删除和导入
- Schema 设计覆盖所有云 API 参数，包括嵌套复杂类型（`net_info`、`metadata_info`、`external_role`）
- 使用 `task_id` 和 `topic_id` 联合 ID（`task_id#topic_id`）作为 Terraform 资源 ID
- 遵循现有 CLS 资源（如 `tencentcloud_cls_cos_shipper`）的代码风格和模式

**Non-Goals:**
- 不修改已有 CLS 资源或数据源
- 不新增数据源（datasource），仅创建资源

## Decisions

### 1. 资源 ID 设计：使用联合 ID `task_id#topic_id`

**决策**: 资源 ID 使用 `task_id` + `tccommon.FILED_SP` + `topic_id` 格式。

**理由**: `DeleteSplunkDeliver` 和 `ModifySplunkDeliver` 接口都需要同时传入 `TaskId` 和 `TopicId`，且 `DescribeSplunkDelivers` 的查询需要 `TopicId`。将两者联合存储可以确保 Read/Update/Delete 操作都能获取到必要的参数。

**备选方案**: 仅使用 `task_id` 作为 ID。但这样在 Read/Delete 时无法获取 `topic_id`，需要额外做 API 查询，增加复杂度。

### 2. Schema 嵌套类型设计

**决策**: 对于 `NetInfo`、`MetadataInfo`、`ExternalRole` 三个复杂类型，使用 `TypeList` + `MaxItems: 1` 的嵌套 schema 结构。

**理由**: 与现有 CLS 资源（如 `cos_shipper` 的 `compress`、`content` 块）保持一致。`TypeList` + `MaxItems: 1` 是 Terraform Plugin SDK v2 中表示单个嵌套对象的惯用模式。

**`net_info` 子字段**:
- `host` (Required, TypeString): 网络地址
- `port` (Required, TypeInt): 端口
- `token` (Required, TypeString): 认证 token
- `net_type` (Required, TypeInt): 网络类型，1=内网，2=外网
- `vpc_id` (Optional, TypeString): 所属网络（内网时必填）
- `virtual_gateway_type` (Optional, TypeInt): 网络服务类型（内网时必填）
- `is_ssl` (Optional, TypeBool): 是否使用 SSL

**`metadata_info` 子字段**:
- `format` (Required, TypeString): 数据格式，rawlog/json
- `meta_fields` (Required, TypeSet of TypeString): 投递字段
- `enable_tag` (Required, TypeBool): 是否投递 TAG 字段
- `tag_json_tiled` (Optional, TypeBool): JSON 是否平铺

**`external_role` 子字段**:
- `role_arn` (Required, TypeString): 跨账户投递角色 RoleArn
- `external_id` (Required, TypeString): 跨账户投递角色名称

### 3. Enable 字段处理

**决策**: `enable` 字段设为 Optional + Computed。

**理由**: `CreateSplunkDeliver` 接口不支持 `Enable` 参数，但 `ModifySplunkDeliver` 支持，且 `DescribeSplunkDelivers` 的响应中包含 `Enable`。创建时由云 API 默认开启，用户可在后续通过 Update 修改。

### 4. Read 操作：通过 DescribeSplunkDelivers 列表接口查询

**决策**: Read 操作使用 `DescribeSplunkDelivers` 接口，通过 `Filters` 参数按 `taskId` 过滤，查询单个投递任务。

**理由**: 该接口是当前唯一可用的查询接口，支持按 `taskId` 过滤。Read 时需从 ID 中解析出 `task_id` 和 `topic_id`，设置 `TopicId` 和 `Filters` 参数。

### 5. 异步操作处理

**决策**: 不在 Schema 中声明 Timeouts 块，Create/Update/Delete 操作使用 retry 机制处理临时失败。

**理由**: 根据云 API 文档，这些接口为同步接口，不涉及异步等待。仅在网络波动或临时错误时使用 `tccommon.RetryError` 和 `tccommon.ReadRetryTimeout` 进行重试。

### 6. Import 支持

**决策**: 支持通过联合 ID `task_id#topic_id` 导入已有资源。

**理由**: 用户可能已有通过控制台或其他方式创建的 Splunk 投递任务，需要纳入 Terraform 管理。

## Risks / Trade-offs

- **[风险] DescribeSplunkDelivers 是列表接口，无单独查询接口**: 需要通过 Filters 精确过滤，若过滤结果为空或返回多条，需要做异常处理。→ 缓解：Read 时若 `Infos` 为空，记录日志后 `d.SetId("")`；若返回多条，取第一条并记录警告。
- **[风险] NetInfo 中 Token 是敏感信息**: Terraform state 中会明文存储。→ 缓解：与现有 CLS 资源行为一致，由用户自行管理 state 安全。
- **[风险] 复杂嵌套类型可能导致 plan diff 问题**: `TypeList` 嵌套结构在 Terraform 中可能出现不必要的 diff。→ 缓解：在 Read 方法中设置所有字段时进行 nil 检查，避免 service 端返回 null 与 terraform 配置默认值不一致导致的 diff。

## Open Questions

<!-- 无待解决问题 -->