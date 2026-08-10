## Context

当前 `tencentcloud_tag_attachment` 资源（位于 `tencentcloud/services/tag/resource_tc_tag_attachment.go`）的 CRUD 实现如下：

- **Create**: 调用 `AddResourceTag` API，参数为 `TagKey`/`TagValue`/`Resource`，成功后设置复合 ID `tagKey#tagValue#resource`（使用 `tccommon.FILED_SP` 分隔符）。
- **Read**: 从 ID 中拆分出三段，调用 service 层 `DescribeTagTagAttachmentById`（封装 `GetResources` API）查询资源标签关联，并 set 各字段。
- **Delete**: 从 ID 中拆分出 `tagKey` 和 `resource`，调用 service 层 `DeleteTagTagAttachmentById`（封装 `DeleteResourceTag` API）。
- **Update**: 当前**不存在** Update 函数。所有字段（`tag_key`、`tag_value`、`resource`）均设置了 `ForceNew: true`，任何修改都会触发删除重建。

**约束：**
- vendor SDK（`tencentcloud-sdk-go/tencentcloud/tag/v20180813`）中已存在 `UpdateResourceTagValue` 接口，请求参数为 `TagKey`、`TagValue`（修改后的新值）、`Resource`（资源六段式），用于"修改资源已关联的标签值（标签键不变）"。
- 该 API 是同步操作，无需轮询状态。
- 复合 ID 中包含 `tagValue`，update 成功后必须用新值重建 ID，否则后续 Read/Delete 会用旧 tag_value 查询导致状态不一致。

**关键文件：**
- `tencentcloud/services/tag/resource_tc_tag_attachment.go` - 资源定义（Schema、CRUD）
- `tencentcloud/services/tag/service_tencentcloud_tag.go` - service 层 API 封装

## Goals / Non-Goals

**Goals:**
- 将 `tag_value` 字段的 `ForceNew` 改为 `false`，支持原地更新。
- 新增 `resourceTencentCloudTagAttachmentUpdate` 函数，检测 `tag_value` 变更时调用 `UpdateResourceTagValue` API。
- 在 service 层封装 `UpdateResourceTagValue` 调用，复用项目标准的 `resource.Retry` + `tccommon.WriteRetryTimeout` 重试模式。
- update 成功后同步更新复合 ID（`tagKey#新tagValue#resource`）。
- 在 Schema 中注册 `Update` 回调。
- 保持 `tag_key` 和 `resource` 字段的 `ForceNew: true` 不变（云 API `UpdateResourceTagValue` 仅支持改值不支持改键）。
- 保持向后兼容性。

**Non-Goals:**
- 不修改 `tag_key` 或 `resource` 字段的可更新性（云 API 不支持修改标签键或换绑资源，且 `UpdateResourceTagValue` 仅改值）。
- 不引入 Timeouts 块（`UpdateResourceTagValue` 为同步操作，无需异步等待）。
- 不修改 Create / Read / Delete 逻辑。
- 不修改资源 ID 格式（仍为 `tagKey#tagValue#resource` 三段式）。

## Decisions

### Decision 1: 仅 tag_value 可更新，tag_key 与 resource 保持 ForceNew

**决策：** 仅将 `tag_value` 的 `ForceNew` 设为 `false`，`tag_key` 和 `resource` 维持 `ForceNew: true`。

**理由：**
- 云 API `UpdateResourceTagValue` 的语义是"修改资源已关联的标签值（标签键不变）"，仅支持改值，不支持改键，也不支持更换关联的资源。
- 若允许 `tag_key` 可更新，需先删除旧标签键再新增新标签键，属于不同 API（`DeleteResourceTag` + `AddResourceTag`），且语义与"更新"不符，应保持 ForceNew 触发重建。
- `resource` 字段代表被打标签的云资源六段式，更换资源本质上是新的标签关联关系，应保持 ForceNew。

**替代方案：**
- 全部去掉 ForceNew 并在 update 中处理多种变更 → 超出本次需求范围，且 tag_key/resource 的"更新"需组合多个 API，复杂度高且语义不清，暂不实现。

### Decision 2: Update 函数通过 d.HasChange("tag_value") 检测变更

**决策：** 在 `resourceTencentCloudTagAttachmentUpdate` 中使用 `d.HasChange("tag_value")` 检测变更，仅在 tag_value 变化时调用 API。

**理由：**
- 当前只有 `tag_value` 一个可更新字段，使用 `d.HasChange` 精确检测即可，无需引入 `mutableArgs` 数组遍历（参考 igtm_strategy 的做法适用于多字段场景）。
- `tag_key` 和 `resource` 仍为 ForceNew，不会进入 Update 函数，无需额外处理。

**实现要点：**
- 从 `d.Id()` 拆分出 `tagKey`、旧 `tagValue`、`resource` 三段。
- 当 `d.HasChange("tag_value")` 时，取新值 `d.Get("tag_value")`，构造 `UpdateResourceTagValueRequest`（`TagKey`、`TagValue`=新值、`Resource`），调用 service 层方法。
- API 调用成功后，用新 tag_value 重建 ID：`d.SetId(tagKey + FILED_SP + newTagValue + FILED_SP + resource)`。
- 最后调用 `resourceTencentCloudTagAttachmentRead(d, meta)` 同步状态。

### Decision 3: service 层封装与重试策略

**决策：** 在 `TagService` 中新增 `UpdateTagAttachmentTagValue` 方法，使用 `resource.Retry(tccommon.WriteRetryTimeout, ...)` + `tccommon.RetryError(e)` 标准重试模式。

**理由：**
- 与项目中其他 service 方法（如 `DeleteTagResourceById`、`ModifyTags`）保持一致。
- `UpdateResourceTagValue` 为同步操作，无需轮询状态，retry 块内仅执行 API 调用。
- 按 openspec 使用须知，设置 id 等成功操作放到 retry 块外。

**替代方案：**
- 直接在 resource 函数中调用 client API → 违反项目分层规范，service 层应封装 API 调用。

### Decision 4: ID 更新时机

**决策：** 在 `UpdateResourceTagValue` API 调用成功（retry 返回 nil）后，于 retry 块外执行 `d.SetId(...)` 重建 ID，再调用 Read。

**理由：**
- 复合 ID `tagKey#tagValue#resource` 中包含 tag_value，更新后若不重建 ID，后续 Read 会用旧 tag_value 查询，可能查不到匹配的标签关联。
- 按代码生成要求，设置 id 等成功操作放到 retry 块外、retry 错误处理后。
- Read 函数会从新 ID 拆分并用新 tag_value 查询，确保状态正确同步。

### Decision 5: Update 回调注册

**决策：** 在 `ResourceTencentCloudTagAttachment()` 返回的 `schema.Resource` 中添加 `Update: resourceTencentCloudTagAttachmentUpdate`。

**理由：**
- 当前资源未注册 Update 回调，Terraform SDK 在没有 Update 函数时，任何非 ForceNew 字段变更都会报错。必须显式注册 Update 函数才能让 `tag_value`（ForceNew=false）的变更生效。

## Risks / Trade-offs

### Risk 1: tag_value 更新后 ID 未同步导致状态不一致
**风险：** 若 Update 成功但 ID 未更新，后续 Read 使用旧 tag_value 查询，会认为资源已删除并清空 state。
**缓解措施：** 在 API 成功后立即用新 tag_value 重建 ID（`d.SetId`），随后调用 Read 用新 ID 同步状态。

### Risk 2: 行为变更影响存量用户
**风险：** 去掉 `tag_value` 的 ForceNew 后，存量用户修改 tag_value 不再触发重建而是原地更新。
**缓解措施：**
- 这是正向变化，符合用户预期（原地更新标签值比删除重建更合理）。
- 云 API `UpdateResourceTagValue` 直接完成修改，不涉及资源重建。
- 在变更日志中说明此行为变更。

### Risk 3: 并发修改标签值冲突
**风险：** 多个 Terraform run 或控制台同时修改同一资源的标签值。
**缓解措施：**
- Terraform state 机制本身不处理并发，依赖用户避免并发操作。
- Read 函数会从云端同步最新状态，下次 plan 会发现差异。
- API 层面的并发冲突由云服务端处理，返回错误时通过 retry 机制重试。

## Migration Plan

**对现有资源的影响：**
1. **已创建的资源**：
   - 不修改 tag_value → 无影响。
   - 修改 tag_value → 从"触发删除重建"变为"原地更新标签值"。

2. **用户迁移步骤**：
   - 无需特殊迁移步骤。
   - 用户首次修改 tag_value 时会调用 `UpdateResourceTagValue` 而非删除重建。

3. **回滚计划**：
   - 如发现严重问题，可将 `tag_value` 的 `ForceNew` 重新设为 `true` 并移除 Update 函数。
   - 已更新标签值的资源不受影响（云端标签值已是新值）。

## Open Questions

1. **Q**: `UpdateResourceTagValue` 是否需要处理资源不存在的错误？
   **A**: 若资源或标签键不存在，API 会返回 `ResourceNotFound.AttachedTagKeyNotFound` 错误，该错误为非重试错误，会直接返回给用户。Read 函数中已有对资源不存在的处理（`d.SetId("")`），无需在 Update 中特殊处理。

2. **Q**: 是否需要在 update 中处理 tag_key 同时变更的情况？
   **A**: 不需要。`tag_key` 仍为 `ForceNew: true`，变更 tag_key 会触发重建，不会进入 Update 函数。
