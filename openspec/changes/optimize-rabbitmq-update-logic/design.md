## Context

当前 `tencentcloud_tdmq_rabbitmq_vip_instance` 资源的 update 实现过于严格，将许多可以通过 ModifyRabbitMQVipInstance API 修改的参数标记为不可变。通过分析腾讯云 TDMQ API v20200217，我们发现 API 实际上支持更多参数的更新，包括 remark（备注）、enable_deletion_protection（删除保护）和 enable_risk_warning（风险提示）。当前实现将这些参数添加到不可变列表中，导致用户无法通过 Terraform 更新这些参数。

## Goals / Non-Goals

**Goals:**
- 在 resource_tc_tdmq_rabbitmq_vip_instance schema 中添加 remark、enable_deletion_protection、enable_risk_warning 三个新参数
- 在 Create、Read、Update 方法中完整支持这些新参数
- 在 Update 方法中移除对这些参数的不可变限制，实现它们的更新逻辑
- 确保所有更改保持向后兼容性
- 添加相应的单元测试覆盖新增功能
- 更新文档以反映新参数

**Non-Goals:**
- 修改 API 接口或其他资源的实现
- 引入新的外部依赖
- 改变现有的参数行为（除了添加对它们的更新支持）
- 修改数据源实现

## Decisions

### 1. Schema 设计决策
选择使用 `Optional: true` 和 `Computed: true` 来定义新参数，因为：
- **Optional**: 用户可以选择是否指定这些参数，不指定时使用 API 默认值
- **Computed**: 这些参数的值可能由 API 设置（如默认值），需要从读取结果中获取
- 这与现有字段（如 `cluster_version`）的设计模式保持一致

### 2. Update 逻辑实现决策
选择在 Update 方法中检测这些参数的变化并调用 API，而不是在不可变列表中阻止它们，因为：
- API 明确支持这些参数的修改（ModifyRabbitMQVipInstanceRequest）
- 提供更好的用户体验，无需销毁重建资源
- 符合 Terraform 资源的最佳实践（可变参数应该可更新）

### 3. API 调用策略
选择将所有变更的参数在一次 API 调用中发送，而不是分别调用多次，因为：
- ModifyRabbitMQVipInstance API 接受多个参数的同时更新
- 减少 API 调用次数，提高性能
- 简化错误处理逻辑（一次成功或失败）

### 4. 向后兼容性保证策略
选择不修改现有字段的行为，只添加新参数，因为：
- 确保现有 Terraform 配置不会受影响
- 不会破坏现有的 state 文件
- 用户可以渐进地采用新功能

## Risks / Trade-offs

### Risk 1: API 参数未来可能变为不可变
**Mitigation:**
- 在代码中添加注释说明当前 API 支持更新这些参数
- 监控腾讯云 API 变更日志
- 如果 API 变更导致参数变为不可变，通过错误消息告知用户

### Risk 2: 新参数可能影响现有资源的行为
**Mitigation:**
- 所有新参数都标记为 Optional，不强制要求
- 读取时使用 Computed 属性从 API 获取实际值
- 确保在未设置时不改变现有行为

### Risk 3: API 返回 nil 值可能导致 panic
**Mitigation:**
- 在 Read 方法中添加 nil 检查
- 使用条件判断避免解引用 nil 指针
- 参考现有代码中的 nil 处理模式（如 resource_tags 的处理）

### Trade-off: 代码复杂度 vs 功能完整性
**Decision:**
- 选择保持适度的代码复杂度，实现完整的功能支持
- 不添加过度抽象，保持代码清晰易读
- 参考现有代码模式（如 cluster_name 和 resource_tags 的更新逻辑）

## Migration Plan

此变更不需要特殊的迁移步骤，因为：

1. **现有资源**: 当用户升级 provider 版本后，执行 `terraform plan` 时，新的可选参数不会触发任何变更
2. **首次使用**: 用户可以在现有资源中添加新参数，provider 会正确处理为更新操作
3. **无状态迁移**: state 文件不需要迁移，新参数会在下次 refresh 时从 API 读取
4. **无破坏性变更**: 所有新参数都是 Optional 的，不会影响现有配置

**部署策略:**
1. 代码审查确保向后兼容性
2. 在测试环境验证新功能
3. 发布新版本到生产环境
4. 更新文档和示例

**回滚策略:**
- 如果发现严重问题，可以快速回退到之前的版本
- 由于是纯功能增强，回退不会影响现有资源
- 不会破坏用户的 Terraform 配置

## Open Questions

1. **问**: 是否需要添加对这些新参数的验证逻辑？
   - **答**: 不需要，API 层面已经包含验证逻辑，Terraform 的类型检查也会捕获基本错误

2. **问**: 是否需要添加对这些参数的 Timeout 支持？
   - **答**: 不需要，ModifyRabbitMQVipInstance API 的响应时间是可预测的，现有重试机制已足够

3. **问**: 是否需要为这些新参数添加 diff suppress 功能？
   - **答**: 不需要，这些参数应该正确反映在 diff 中，帮助用户理解变更
