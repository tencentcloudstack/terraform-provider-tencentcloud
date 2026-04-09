## Context

当前 Terraform Provider for TencentCloud 缺少 TEO (TencentCloud EdgeOne) 边缘 KV 查询资源。TEO SDK 已经提供了 EdgeKVGet API，可以直接在 Provider 中调用该接口来实现查询功能。Provider 遵循 Terraform Plugin SDK v2 规范，代码组织在 `tencentcloud/services/teo/` 目录下。

## Goals / Non-Goals

**Goals:**
- 实现完整的 Terraform Resource `tencentcloud_teo_edge_k_v_get`，支持 CRUD 操作
- 定义符合 Terraform 规范的 Resource Schema
- 通过 TEO SDK 的 EdgeKVGet API 实现查询功能
- 提供完整的单元测试和验收测试
- 生成对应的文档和示例

**Non-Goals:**
- 不修改现有 TEO 资源的行为
- 不涉及数据迁移或状态迁移
- 不引入新的外部依赖

## Decisions

### Resource Schema 设计

采用 Terraform Plugin SDK v2 的 `schema.Resource` 定义 Schema：
- 使用 `TypeString` 定义 ZoneId 和 Namespace 参数
- 使用 `TypeList` + `Elem: &schema.Schema{Type: TypeString}` 定义 Keys 参数
- 使用 `TypeList` + `Elem: &schema.Resource{Schema: {...}}` 定义 Data 参数，包含 Key、Value 和 Expiration 字段

**理由**：符合 Terraform Provider 规范，TypeList 支持数组类型，嵌套的 Resource 支持复杂结构。

### CRUD 操作实现

- **Create**: 调用 EdgeKVGet API 查询数据，将结果设置到 State
- **Read**: 再次调用 EdgeKVGet API，更新 State 中的数据
- **Update**: 由于是查询操作，Update 直接调用 Read 即可
- **Delete**: 删除 State 中的数据（从 Terraform 状态中移除）

**理由**：对于查询型资源，Create/Read 是核心操作，Update/Delete 遵循 Terraform 惯例。

### 错误处理和重试

- 使用 `defer tccommon.LogElapsed(ctx, "resource.tencentcloud_teo_edge_k_v_get.create")` 记录耗时
- 使用 `defer tccommon.InconsistentCheck(d, meta)` 检查最终一致性
- 使用 `helper.Retry()` 处理可能的 API 重试

**理由**：遵循 Provider 项目的通用模式，确保稳定性和可观测性。

### 资源 ID 设计

使用复合 ID 格式：`zoneId#namespace#keysHash`
- zoneId: 站点 ID
- namespace: 命名空间
- keysHash: Keys 列表的哈希值（用于唯一标识不同的键组合）

**理由**：复合 ID 能够唯一标识资源，哈希处理避免 ID 过长。

## Risks / Trade-offs

### 风险
- [API 限流] → 如果 Keys 列表过大，可能触发 API 限流。通过限制 Keys 数组长度（上限 20）来降低风险。
- [键不存在] → 当查询的键不存在时，Value 字段返回空字符串。需要在文档中明确说明此行为。

### 权衡
- [查询性能 vs 数据完整性] → 一次查询最多支持 20 个键，如果需要查询更多键，需要多次调用。这是 API 的限制，需要在文档中说明。
- [ID 可读性 vs 唯一性] → 使用 keysHash 可能降低 ID 的可读性，但确保了唯一性。可以在文档中说明 ID 的格式。

## Open Questions

无
