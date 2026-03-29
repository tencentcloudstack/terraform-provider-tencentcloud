## Context

当前 tencentcloud_teo_l7_acc_rule 资源通过 ImportZoneConfig API 进行更新操作。该 API 支持异步任务处理，但资源实现中未传入 TaskId 参数，这导致无法正确关联和追踪异步操作任务。TEO (TencentCloud EdgeOne) 服务需要通过 TaskId 来标识和管理异步操作任务，提升操作的可追踪性和可靠性。

资源文件位于 `tencentcloud/services/teo/resource_tencentcloud_teo_l7_acc_rule.go`，服务层逻辑位于 `tencentcloud/services/teo/service_tencentcloud_teo.go`。

## Goals / Non-Goals

**Goals:**
- 为 tencentcloud_teo_l7_acc_rule 资源添加 TaskId 可选参数
- 确保在调用 ImportZoneConfig API 时正确传递 TaskId 参数
- 保持完全向后兼容性，不影响现有配置
- 提供清晰的文档和示例说明 TaskId 的用途

**Non-Goals:**
- 不修改资源的其他行为或参数
- 不改变资源的 ID 生成逻辑
- 不引入新的外部依赖
- 不修改资源的创建和删除操作

## Decisions

### 1. Schema 设计
**Decision:** 在资源 schema 中添加 `task_id` 字段，类型为 `schema.TypeString`，标记为 `Optional` 和 `Computed`。

**Rationale:**
- 使用 `Optional` 保持向后兼容，现有配置无需修改
- 使用 `Computed` 允许 API 返回 TaskId 值并保存到 state
- 遵循 Terraform Provider 命名约定（snake_case）
- 使用 Go 基础类型 String，无需复杂验证逻辑

### 2. API 调用集成
**Decision:** 在 update 函数中，从 schema 读取 TaskId 值，并在调用 ImportZoneConfig API 时将该参数添加到请求中。

**Rationale:**
- 最小化代码变更，仅修改 update 操作相关逻辑
- 遵循现有 API 调用模式，使用 tencentcloud-sdk-go
- TaskId 为可选参数，API 支持不传该参数的情况

### 3. 文档更新
**Decision:** 更新 `website/docs/r/teo_l7_acc_rule.html.markdown` 文档，添加 TaskId 参数说明和使用示例。

**Rationale:**
- 提供清晰的用户指导
- 符合项目文档规范（所有资源必须有文档）
- 包含实际使用场景示例

## Risks / Trade-offs

### Risk 1: API 兼容性问题
[Risk] 如果 ImportZoneConfig API 对 TaskId 参数格式或值有特殊要求，可能导致更新操作失败。

**Mitigation:** 在实现前通过 API 文档确认 TaskId 的正确格式和约束条件，在测试阶段验证各种输入场景。

### Risk 2: State 不一致
[Risk] 如果 API 返回的 TaskId 与用户输入不一致，可能导致 state 漂移。

**Mitigation:** 将 TaskId 标记为 `Computed`，允许 API 返回值覆盖用户输入，保持 state 与实际 API 状态一致。

### Risk 3: 向后兼容性
[Risk] 添加新参数可能意外影响现有资源更新流程。

**Mitigation:** 确保参数为 `Optional`，不改变现有参数逻辑，并通过完整的回归测试验证。

## Migration Plan

### 部署步骤
1. 修改资源 schema 定义，添加 task_id 字段
2. 更新 update 函数，集成 TaskId 参数到 API 调用
3. 更新资源文档和使用示例
4. 添加单元测试和验收测试
5. 运行完整测试套件，确保无回归问题

### 回滚策略
如果发现问题，可以：
- 保持新参数为 Optional，立即回滚到旧版本不影响现有配置
- 用户可选择从配置中移除 task_id 参数

## Open Questions

无。需求明确，实现路径清晰。
