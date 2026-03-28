## Context

当前 RabbitMQ VIP 实例资源 `tencentcloud_tdmq_rabbitmq_vip_instance` 的 update 逻辑存在严重限制。通过分析代码发现，`ModifyRabbitMQVipInstance` API 仅支持修改元数据信息（集群名称、标签、备注、删除保护），而大量核心配置参数（节点规格、节点数量、存储大小、带宽、公网访问等）被标记为不可修改。

这种限制导致用户在实际使用中无法通过 Terraform 灵活地调整实例配置，特别是在业务增长需要扩展资源或进行配置优化的场景下，降低了资源管理的效率和灵活性。

## Goals / Non-Goals

**Goals:**
- 移除不必要的参数限制，使核心配置参数支持 update 操作
- 为无法通过 `ModifyRabbitMQVipInstance` API 修改的参数，提供合理的替代方案
- 保持向后兼容性，不破坏现有用户配置
- 改进错误处理和重试逻辑，提升 update 操作的稳定性
- 增强日志记录，便于问题排查和审计

**Non-Goals:**
- 不修改现有资源的 schema（只新增 Optional 字段）
- 不破坏 Terraform state 的兼容性
- 不改变已存在资源的默认行为（除非有明确的 bug）

## Decisions

### 决策 1：分阶段实现参数更新支持

**理由：**
- 部分参数（如 `node_spec`、`node_num`、`storage_size`）可能涉及实例重建，需要与用户沟通
- 部分参数（如 `band_width`、`enable_public_access`）可能有独立的 API 或需要特殊处理
- 分阶段实现可以降低风险，便于逐步验证和回滚

**替代方案考虑：**
- 方案 A：一次性放开所有参数限制 - 风险过高，可能导致大量用户实例被重建
- 方案 B：仅支持可以通过 `ModifyRabbitMQVipInstance` API 修改的参数 - 过于保守，无法解决核心问题
- 方案 C（采用）：分阶段实现，先探索和验证可行的参数更新方式，再逐步扩展

### 决策 2：采用 ForceNew 机制处理需要重建的参数

**理由：**
- 部分参数（如 `zone_ids`、`vpc_id`、`subnet_id`、`cluster_version`）确实需要重建实例
- 使用 ForceNew 可以向用户明确传达重建的影响
- Terraform 的 ForceNew 机制是处理不可变参数的标准方式

**实现方式：**
- 将需要重建的参数的 `ForceNew` 属性设置为 true
- 在 schema 注释中明确说明该参数的变更会导致实例重建
- 在 update 函数中，检查这些参数的变化并返回明确的错误信息

### 决策 3：探索和利用其他 API 实现参数更新

**理由：**
- `ModifyRabbitMQVipInstance` API 的限制可能不是绝对的
- 可能有其他隐含的 API 或组合操作可以实现参数更新
- 需要与 TDMQ 团队沟通，了解 API 的发展方向

**具体策略：**
- 搜索和测试其他 TDMQ API，如 `ModifyCluster`、`ModifyNetwork` 等
- 参考其他云服务的实现方式，如 `ModifyRocketMQInstanceSpec`
- 如果确认不存在支持更新的 API，则使用 ForceNew 机制

### 决策 4：改进错误处理和重试逻辑

**理由：**
- 当前 update 函数的错误处理较为简单
- 异步操作（如实例配置变更）可能需要更长的重试时间
- 需要更详细的错误信息帮助用户理解问题

**实现方式：**
- 使用 `helper.Retry()` 进行最终一致性重试
- 为不同的操作设置不同的重试超时时间
- 在日志中记录详细的请求和响应信息
- 提供友好的错误信息，指导用户如何解决问题

### 决策 5：增强日志记录和审计

**理由：**
- 当前 update 操作的日志记录较为简单
- 需要更好的可观测性来监控 update 操作
- 便于问题排查和审计

**实现方式：**
- 使用 `tccommon.LogElapsed()` 记录每个操作的时间
- 记录 update 操作的详细参数
- 记录 API 调用的请求和响应
- 在日志中标识哪些参数发生了变化

## Risks / Trade-offs

### 风险 1：API 限制导致部分参数无法更新
- **风险**：TDMQ API 可能确实不支持某些参数的更新，只能通过重建实例实现
- **缓解措施**：使用 ForceNew 机制，明确告知用户重建的影响；在文档中说明 API 的限制

### 风险 2：向后兼容性问题
- **风险**：修改 update 逻辑可能破坏现有用户配置
- **缓解措施**：不修改 schema，只放宽参数限制；新增 ForceNew 属性而不是移除参数；充分测试现有配置

### 风险 3：状态一致性问题
- **风险**：异步操作可能导致 Terraform state 与实际状态不一致
- **缓解措施**：使用 `tccommon.InconsistentCheck()` 进行状态一致性检查；在 update 后立即调用 Read 函数刷新状态

### 风险 4：测试覆盖不足
- **风险**：新的 update 逻辑可能缺乏充分的测试
- **缓解措施**：编写全面的测试用例，覆盖各种 update 场景；使用 TENCENTCLOUD_SECRET_ID/KEY 环境变量运行验收测试

### 风险 5：性能影响
- **风险**：update 操作可能需要更长的时间，影响 Terraform apply 的性能
- **缓解措施**：优化重试逻辑，减少不必要的 API 调用；提供清晰的进度提示；支持 Timeout 配置

## Migration Plan

1. **分析阶段**：深入分析 TDMQ API，确定哪些参数可以通过 API 更新，哪些需要重建
2. **设计阶段**：根据分析结果，设计新的 update 逻辑，确定哪些参数支持 ForceNew
3. **实现阶段**：分阶段实现，先实现可以 API 更新的参数，再实现需要重建的参数
4. **测试阶段**：编写全面的测试用例，运行验收测试，确保向后兼容性
5. **文档阶段**：更新 provider 文档和示例，说明新的 update 能力和限制
6. **发布阶段**：发布新版本，监控用户反馈，及时修复问题

## Open Questions

1. **API 能力确认**：TDMQ API 是否支持通过其他方式修改 `node_spec`、`node_num`、`storage_size`、`band_width`、`enable_public_access` 等参数？
   - 需要与 TDMQ 团队沟通，确认 API 的发展方向
   - 需要测试其他可能的 API，如 `ModifyCluster`、`ModifyNetwork` 等

2. **重建策略确认**：对于需要重建的参数，是否应该提供自动重建能力，还是要求用户手动重建？
   - 需要与用户沟通，了解他们的需求和期望
   - 需要考虑 Terraform 的最佳实践

3. **回滚策略**：如果 update 操作失败，如何回滚到之前的状态？
   - 需要设计合理的回滚策略
   - 需要考虑数据的备份和恢复
