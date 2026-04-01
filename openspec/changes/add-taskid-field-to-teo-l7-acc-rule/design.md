## Context

本次变更在已有的 `tencentcloud_teo_l7_acc_rule` 资源中新增 `task_id` 字段。该字段由 `ImportZoneConfig` API 在 Read 操作中返回，表示配置导入任务的任务 ID。用户可以通过该 ID 使用 `DescribeZoneConfigImportResult` 接口查询最近 7 天内的导入任务执行结果。

当前实现：
- 资源文件：`tencentcloud/services/teo/resource_tencentcloud_teo_l7_acc_rule.go`
- 测试文件：`tencentcloud/services/teo/resource_tencentcloud_teo_l7_acc_rule_test.go`
- 文档文件：`website/docs/r/teo_l7_acc_rule.md`

约束条件：
- 必须保持向后兼容，不能破坏现有 TF 配置和 state
- `task_id` 字段应为 Computed 属性，不由用户设置
- 不能修改已有字段的 schema 定义

## Goals / Non-Goals

**Goals:**
- 在 `tencentcloud_teo_l7_acc_rule` 资源的 Schema 中新增 `task_id` 字段
- 在 Read 函数中从 `ImportZoneConfig` API 响应中读取 `TaskId` 并填充到 resource state
- 更新单元测试以验证 `task_id` 字段正确读取
- 更新资源文档，添加 `task_id` 字段的说明
- 确保完全向后兼容

**Non-Goals:**
- 不修改 Create、Update、Delete 操作（`task_id` 仅在 Read 时返回）
- 不修改其他资源的实现
- 不引入新的外部依赖
- 不改变资源的行为逻辑

## Decisions

### Schema 定义
在 `resourceTencentcloudTeoL7AccRule()` 函数中添加 `task_id` 字段：
- 类型：`schema.TypeString`
- 属性：`Computed: true`（仅由 API 返回，不由用户设置）
- 描述：明确说明该字段用于追踪配置导入任务

### Read 操作实现
在 Read 函数中，从 `ImportZoneConfig` API 响应的 `TaskId` 字段读取任务 ID，并使用 `d.Set("task_id", response.TaskId)` 设置到 state。

### 测试策略
在单元测试中添加测试用例，验证当 API 返回 `TaskId` 时，该字段正确填充到 resource state。同时确保不包含 `TaskId` 的响应也不会导致错误。

### 文档更新
在资源文档的 **Arguments Reference** 和 **Attributes Reference** 部分添加 `task_id` 字段的说明。

## Risks / Trade-offs

| Risk | Mitigation |
|------|------------|
| API 返回格式变更 | 使用 API SDK 的类型安全访问，确保字段读取的健壮性 |
| 测试环境不返回 TaskId | 单元测试中添加可选的 TaskId 测试用例，验证有和无 TaskId 两种情况 |
| 向后兼容性问题 | 确保新增字段为 Computed 属性，不影响现有配置 |

## Migration Plan

无需迁移计划。此变更仅添加 Computed 字段，完全向后兼容。

## Open Questions

无。这是一个简单的字段添加变更，所有实现细节都已明确。
