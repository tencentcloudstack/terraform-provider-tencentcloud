## Context

Terraform Provider for TencentCloud 已支持多种 CLS 资源（alarm、topic、cos_shipper 等），但尚未覆盖 RemoteWrite 投递任务。CLS RemoteWrite 任务允许用户将日志通过 Prometheus RemoteWrite 协议投递到外部目标服务，是一个通用的日志投递能力。

云 API 已在 `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016` 中提供完整的 CRUD 接口：
- `CreateRemoteWriteTask`：创建任务，返回 `TaskId`
- `DescribeRemoteWriteTasks`：按 Filters 过滤查询任务列表（支持按 taskId 过滤）
- `ModifyRemoteWriteTask`：修改任务配置
- `DeleteRemoteWriteTask`：删除任务（需 TaskId + TopicId）

当前资源类型为 RESOURCE_KIND_GENERAL，需参考 `tencentcloud_igtm_strategy` 资源的代码风格实现完整 CRUD 生命周期。

## Goals / Non-Goals

**Goals:**
- 实现 `tencentcloud_cls_remote_write_task` 资源的完整 CRUD 生命周期管理
- 支持 Create、Read、Update、Delete 四个操作，分别映射到对应云 API
- 使用 `task_id` + `topic_id` 组合作为复合 ID（`task_id#topic_id`），因为 Modify 和 Delete 接口均要求同时传入这两个字段
- 支持 `auth_info` 嵌套块（含 username、password、token）的读写
- 支持资源 Import（使用复合 ID）
- 在 provider.go 中注册资源
- 使用 gomonkey mock 云 API 进行单元测试

**Non-Goals:**
- 不实现数据源（datasource）变体
- 不修改已有 CLS 资源的 schema 或行为
- 不处理 RemoteWrite 任务的异步状态轮询（云 API 为同步接口，Create/Modify/Delete 直接返回结果）

## Decisions

### 1. 复合 ID 格式：`task_id#topic_id`

**决策**: 使用 `task_id` + `topic_id` 以 `tccommon.FIELD_SP`（`#`）分隔符组合作为资源 ID。

**理由**:
- `DeleteRemoteWriteTask` 和 `ModifyRemoteWriteTask` 接口均要求同时传入 `TaskId` 和 `TopicId`
- Read 操作通过 `DescribeRemoteWriteTasks` 按 `taskId` 过滤获取详情，但 Delete/Modify 需要 topic_id
- 复合 ID 确保 Delete 和 Modify 操作能从 `d.Id()` 解析出两个必要字段
- 符合项目约定（使用 `tccommon.FILED_SP` 作为分隔符）

**备选方案**: 仅用 `task_id` 作为 ID，将 `topic_id` 设为 ForceNew 字段。但这会导致 topic_id 变更时重建资源，而 Modify 接口本身支持修改 topic_id，因此不合适。

### 2. Read 操作使用 DescribeRemoteWriteTasks + Filters

**决策**: Read 操作通过 `DescribeRemoteWriteTasks` 接口，设置 `Filters` 中 `Key=taskId`、`Values=[task_id]` 进行过滤查询。

**理由**:
- 云 API 未提供 `DescribeRemoteWriteTask`（单数）接口，只有 `DescribeRemoteWriteTasks`（列表接口）
- 按 taskId 过滤可精确定位单个任务
- 从 `Response.Infos` 列表的第一项提取任务详情

### 3. auth_info 作为嵌套块（schema.TypeList, MaxItems: 1）

**决策**: `auth_info` 使用 `schema.TypeList` + `MaxItems: 1` + `Elem: schema.Resource` 结构，内含 `username`、`password`、`token` 三个可选字符串字段。

**理由**:
- 云 API 中 `AuthInfo` 是一个 `RemoteWriteAuthInfo` 结构体（含 Username、Password、Token）
- 使用嵌套块比三个独立的顶层字段更贴合 API 结构，且语义更清晰
- 参考 CLS 其他资源（如 cos_shipper 的 content 块）的嵌套块模式

### 4. enable 字段在 Update 时可修改

**决策**: `enable` 字段设为可选非 ForceNew 字段，在 Update 操作中通过 `ModifyRemoteWriteTask` 的 `Enable` 参数修改。

**理由**:
- `ModifyRemoteWriteTask` 接口支持 `Enable` 参数（0 关闭 / 1 开启）
- 用户可能需要在不重建资源的情况下启停任务

### 5. 只读字段：status、create_time、update_time、logset_id

**决策**: 这些字段仅在 Read 操作中从 API 响应设置到 state，不在 schema 中标记为用户可配置（Computed: true）。

**理由**:
- `DescribeRemoteWriteTasks` 返回这些字段，但 `CreateRemoteWriteTask` 和 `ModifyRemoteWriteTask` 入参不包含它们
- `status`：任务运行状态（1 运行中 / 2 暂停 / 3 失败），由云服务维护
- `create_time`、`update_time`：由云服务维护
- `logset_id`：日志集 ID，由云服务关联返回

### 6. auth_type 和 net_type 使用 uint64 → int 类型映射

**决策**: 云 API 中 `AuthType` 和 `NetType` 为 `uint64` 类型，Terraform schema 中使用 `schema.TypeInt`。

**理由**:
- Terraform Plugin SDK 中整数类型统一使用 `schema.TypeInt`
- `uint64` 与 `int` 在实际使用范围内一致（0/1/2 等小整数）
- 转换时使用 `helper.IntUint64()` 或直接类型转换

### 7. 单元测试使用 gomonkey mock 云 API

**决策**: 单元测试不使用 Terraform 测试套件（TF_ACC），而是使用 gomonkey 对云 API 调用进行 mock，仅测试业务逻辑。

**理由**:
- 遵循项目规范：新增 Terraform 资源的测试使用 mock 方法，不依赖实际云环境
- 避免 TF_ACC 测试需要真实凭证和网络访问
- 专注于验证 CRUD 逻辑的正确性

## Risks / Trade-offs

- **[风险] DescribeRemoteWriteTasks 返回列表而非单条** → 通过 taskId 过滤精确匹配，并在结果为空时按规范先打印日志再 `d.SetId("")`，避免静默清空 state
- **[风险] AuthInfo 可能为 nil** → 在 Read 操作中对 `AuthInfo` 指针做 nil 检查，仅非 nil 时才设置 `auth_info` 块
- **[风险] Create 返回 TaskId 可能为空** → 在 Create 操作中检查返回值是否为空，若为空则返回 `NonRetryableError`，避免写入空 ID 导致后续状态混乱
- **[风险] 字段在 Create 和 Modify 接口中可用性不一致** → 已核对：Create 入参不含 `Enable`，Modify 入参含 `Enable`；Create 必填 `TopicId/Name/Target/RemoteWriteURL/AuthType/NetType`，Modify 中这些字段为可选。Schema 中按 Create 必填字段标记 Required，其余 Optional
- **[权衡] auth_info 嵌套块 vs 平铺字段** → 选择嵌套块以匹配 API 结构，但增加了 schema 复杂度。可接受，因为符合项目现有模式
