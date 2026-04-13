## Context

当前 `tencentcloud_tdmq_rabbitmq_vip_instance` 资源仅支持有限的可更新字段（cluster_name 和 resource_tags），而腾讯云 TDMQ RabbitMQ API 的 ModifyRabbitMQVipInstanceRequest 实际上支持更多可更新参数：

- **Remark**: 实例备注说明
- **EnableDeletionProtection**: 是否开启删除保护
- **EnableRiskWarning**: 是否开启集群风险提示

这些字段在 DescribeRabbitMQVipInstance API 响应的 RabbitMQClusterInfo 结构中可用，但当前代码没有利用这些功能。用户无法通过 Terraform 管理这些重要的实例属性，特别是在生产环境中删除保护和风险提示是关键的安全特性。

## Goals / Non-Goals

**Goals:**

1. 在 Schema 中新增 remark、enable_deletion_protection 和 enable_risk_warning 字段
2. 在 Update 函数中支持这些字段的更新，调用 ModifyRabbitMQVipInstance API
3. 在 Read 函数中从云 API 响应中读取这些字段的值并设置到 Terraform state
4. 保持向后兼容性，不破坏现有配置和 state
5. 为新增字段添加单元测试

**Non-Goals:**

1. 不修改已有的字段定义和行为
2. 不涉及其他 RabbitMQ 资源的变更
3. 不改变实例创建流程（Create 函数）
4. 不提供从外部系统迁移这些字段的功能

## Decisions

**1. Schema 字段设计**

- `remark`: 类型为 `TypeString`，Optional，不设置 Computed
  - 理由：备注是一个简单的文本字段，用户可以自由设置和修改
  - 验证：不需要特殊验证，符合云 API 的要求（3-64 个字符，仅数字、字母、"-"和"_"）

- `enable_deletion_protection`: 类型为 `TypeBool`，Optional，Computed
  - 理由：删除保护是实例的安全特性，用户可以在创建时设置，也可以在运行时修改。设置为 Computed 是因为云 API 可能提供默认值
  - 默认值：不设置默认值，依赖云 API 的默认行为

- `enable_risk_warning`: 类型为 `TypeBool`，Optional，Computed
  - 理由：风险提示是集群的安全特性，用户可以控制。设置为 Computed 是因为云 API 可能提供默认值
  - 默认值：不设置默认值，依赖云 API 的默认行为

**2. Create 函数处理**

- 在 Create 函数中，对于这三个新增字段，如果用户提供了值，则传递给 CreateRabbitMQVipInstance API
- 检查 CreateRabbitMQVipInstanceRequest 是否支持这些字段
  - 如果不支持，这些字段只能在 Update 时设置
  - 如果支持，则在 Create 时也设置这些字段

**3. Update 函数实现**

- 保持现有的不可变参数列表（immutableArgs），因为这些参数确实不支持通过 ModifyRabbitMQVipInstance API 更新
- 对于新增的三个字段：
  - 检测字段是否发生变化（d.HasChange）
  - 如果发生变化，构建 ModifyRabbitMQVipInstanceRequest 并调用云 API
  - 支持部分字段更新（不需要一次性更新所有字段）

**4. Read 函数实现**

- 从 DescribeRabbitMQVipInstance API 响应的 ClusterInfo 中读取这些字段：
  - ClusterInfo.Remark → `remark`
  - ClusterInfo.EnableDeletionProtection → `enable_deletion_protection`
  - ClusterInfo.EnableRiskWarning → `enable_risk_warning`
- 处理 nil 值：如果字段为 nil，则不设置或设置为默认值

**5. 向后兼容性**

- 新增字段均为 Optional，不会影响现有配置
- 不修改现有字段的定义和行为
- 不改变现有 API 调用逻辑

## Risks / Trade-offs

**Risk 1: Create API 可能不支持新字段**

- **风险**：CreateRabbitMQVipInstanceRequest 可能不支持 remark、enable_deletion_protection 和 enable_risk_warning 字段
- **缓解措施**：
  1. 检查 vendor 目录下的 models.go，确认 CreateRabbitMQVipInstanceRequest 的字段定义
  2. 如果不支持，则在 Create 函数中不传递这些字段
  3. 在提案文档中明确说明这些字段只能通过 Update 设置

**Risk 2: 云 API 响应中字段可能为 nil**

- **风险**：某些旧实例或特定配置下，云 API 响应中的这些字段可能为 nil
- **缓解措施**：
  1. 在 Read 函数中进行 nil 检查
  2. 仅在字段不为 nil 时才设置到 state
  3. 添加单元测试覆盖 nil 值场景

**Risk 3: 更新操作可能失败或超时**

- **风险**：ModifyRabbitMQVipInstance API 可能因为实例状态或其他原因失败
- **缓解措施**：
  1. 使用现有的 resource.Retry 机制进行重试
  2. 返回清晰的错误信息，帮助用户排查问题
  3. 遵循现有的错误处理模式（defer tccommon.LogElapsed(), defer tccommon.InconsistentCheck()）

**Trade-off: 字段验证的复杂度**

- **权衡**：remark 字段在云 API 中有长度和字符限制，但 Terraform provider 通常不在 schema 层面进行这类验证
- **决策**：不在 Schema 中添加 ValidateFunc，依赖云 API 的验证
- **理由**：保持一致性，现有字段大多没有复杂的 ValidateFunc；云 API 的验证更准确

**Trade-off: Computed 字段的初始值**

- **权衡**：enable_deletion_protection 和 enable_risk_warning 设置为 Optional + Computed，意味着用户不指定时会有默认值
- **决策**：让云 API 决定默认值，provider 不预设默认值
- **理由**：避免与云 API 的默认行为不一致；用户可以通过 Read 操作查看实际值
