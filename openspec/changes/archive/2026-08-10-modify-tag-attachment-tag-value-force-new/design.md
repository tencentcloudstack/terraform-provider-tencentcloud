## Context

当前 `tencentcloud_tag_attachment` 资源只实现了 Create/Read/Delete 三个操作，`tag_value` 字段为 `ForceNew: true`。用户修改 `tag_value` 时，Terraform 会先删除旧标签关联再创建新标签关联，导致标签短暂脱离资源，影响基于标签的权限策略、计费分账和自动化运维。

腾讯云 tag 服务（v20180813）提供了 `UpdateResourceTagValue` 接口，其描述为"本接口用于修改资源已关联的标签值（标签键不变）"，入参为 `TagKey`、`TagValue`（修改后的值）、`Resource`，与当前资源 schema 的三个字段完全对应。该接口已存在于 vendor 目录中，无需升级 SDK。

**关键文件：**
- `tencentcloud/services/tag/resource_tc_tag_attachment.go` — 资源定义（Schema + CRUD）
- `tencentcloud/services/tag/service_tencentcloud_tag.go` — 服务层，已有 `DescribeTagTagAttachmentById` 和 `DeleteTagTagAttachmentById`
- `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tag/v20180813/models.go:2697` — `UpdateResourceTagValueRequest` 定义

**资源 ID 格式：** `tagKey + FILED_SP + tagValue + FILED_SP + resourceId`，修改 `tag_value` 后 ID 会变化，需重新 `SetId`。

## Goals / Non-Goals

**Goals:**
- 将 `tag_value` 字段从 `ForceNew: true` 改为 `ForceNew: false`，支持原地更新
- 新增 `Update` 函数，当 `tag_value` 发生变化时调用 `UpdateResourceTagValue` API
- 在服务层封装 `UpdateResourceTagValue` 调用
- 保持向后兼容，不影响不修改 `tag_value` 的现有配置
- 补充单元测试覆盖 update 逻辑

**Non-Goals:**
- 不支持修改 `tag_key`（`UpdateResourceTagValue` API 要求标签键不变，修改 `tag_key` 仍需重建）
- 不支持修改 `resource`（资源六段式变更等价于换资源，仍为 ForceNew）
- 不引入 Timeouts 块（`UpdateResourceTagValue` 是同步操作，无需轮询等待）

## Decisions

### Decision 1: Update 函数处理 tag_value 变更
**决策：** 新增 `resourceTencentCloudTagAttachmentUpdate` 函数，在其中检测 `d.HasChange("tag_value")`，调用 `UpdateResourceTagValue` API 完成原地修改。

**理由：**
- `UpdateResourceTagValue` API 的语义是"标签键不变，仅修改标签值"，与 `tag_value` 可更新需求完全契合
- API 入参 `TagKey`、`TagValue`、`Resource` 与资源 schema 字段一一对应，无需额外转换

**替代方案：**
- 在 Update 中删除旧标签关联再创建新关联 → 仍会导致标签短暂脱离资源，违背需求初衷，且浪费两次 API 调用

### Decision 2: 更新后重新设置资源 ID
**决策：** 调用 `UpdateResourceTagValue` 成功后，使用新的 `tag_value` 重新拼接并 `d.SetId(tagKey + FILED_SP + newTagValue + FILED_SP + resourceId)`，随后调用 `Read` 同步状态。

**理由：**
- 资源 ID 包含 `tag_value`，更新后旧 ID 不再匹配云端实际标签值
- 重新 `SetId` 保证后续 Read/Delete 操作能基于新值定位资源

**实现：**
```go
// retry 块外、错误处理后
d.SetId(tagKey + tccommon.FILED_SP + newTagValue + tccommon.FILED_SP + resourceId)
return resourceTencentCloudTagAttachmentRead(d, meta)
```

### Decision 3: tag_key 和 resource 变更仍触发重建
**决策：** `tag_key` 和 `resource` 字段保持 `ForceNew: true` 不变。

**理由：**
- `UpdateResourceTagValue` API 明确要求标签键不变（"标签键不变"）
- `resource` 变更意味着绑定到不同资源，语义上属于新关联，应重建

### Decision 4: 服务层方法封装
**决策：** 在 `service_tencentcloud_tag.go` 新增 `UpdateTagTagAttachment(ctx, tagKey, newTagValue, resource string)` 方法，封装 `UpdateResourceTagValue` 请求构造、调用和错误处理。

**理由：**
- 与现有 `DeleteTagTagAttachmentById` 方法风格保持一致
- 集中管理 API 调用与日志，便于维护

### Decision 5: API 调用使用 WriteRetryTimeout 重试
**决策：** Update 中的 API 调用使用 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 包装，失败时使用 `tccommon.RetryError()` 包装错误。

**理由：**
- 符合项目规范：调用云 API 接口时以 `tccommon.WriteRetryTimeout` 作为超时时间添加 retry 处理
- 与 Create 函数中的 retry 逻辑保持一致

### Decision 6: 单元测试使用 gomonkey mock
**决策：** 在 `resource_tc_tag_attachment_test.go` 中使用 gomonkey 对云 API 进行 mock，仅测试业务逻辑，不使用 terraform 测试套件。

**理由：**
- 本次为修改现有资源（新增 Update），但 Update 为全新函数，按规范使用 mock 方式补充单元测试
- 避免依赖真实云资源和 `TF_ACC` 环境变量

## Risks / Trade-offs

### Risk 1: tag_value 更新失败导致 state 不一致
**风险：** `UpdateResourceTagValue` API 调用失败，state 中已记录新值但云端仍为旧值。
**缓解措施：** API 调用失败时直接返回 error，不执行 `SetId`，Terraform 会保持旧 state，用户可重试。

### Risk 2: 并发修改 tag_value
**风险：** 多个 Terraform 进程同时修改同一资源的 tag_value。
**缓解措施：** 依赖云 API 的并发控制和错误返回，Provider 层不做额外加锁（与项目其他资源一致）。

### Risk 3: 行为变更影响存量用户
**风险：** 存量用户修改 `tag_value` 时，从"重建"变为"原地更新"。
**缓解措施：** 该变更属于正向增强，原地更新比重建更安全，用户预期一致。
