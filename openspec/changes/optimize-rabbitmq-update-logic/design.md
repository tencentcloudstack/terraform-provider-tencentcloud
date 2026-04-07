## Context

RabbitMQ VIP 实例当前的 Update 逻辑在 `resource_tc_tdmq_rabbitmq_vip_instance.go` 中存在以下限制：

1. **API 能力未充分利用**：`ModifyRabbitMQVipInstance` API 支持多个字段（`Remark`、`EnableDeletionProtection`、`EnableRiskWarning`），但当前实现只使用了 `ClusterName` 和 `Tags` 字段。

2. **不可修改字段列表不准确**：当前代码将以下字段标记为不可修改（immutable）：
   - `zone_ids`、`vpc_id`、`subnet_id` - 网络配置
   - `node_spec`、`node_num`、`storage_size` - 节点规格
   - `enable_create_default_ha_mirror_queue` - 镜像队列配置
   - `auto_renew_flag`、`time_span`、`pay_mode` - 计费相关
   - `cluster_version` - 集群版本
   - `band_width`、`enable_public_access` - 公网访问

   部分字段可能通过其他 API 支持修改，需要验证。

3. **缺少状态等待机制**：Update 操作后没有等待实例达到稳定状态的机制，可能导致状态不一致。

4. **错误处理不完善**：错误消息不够清晰，缺少重试和超时处理的详细信息。

## Goals / Non-Goals

**Goals:**

1. **充分利用 API 能力**：支持所有 `ModifyRabbitMQVipInstance` API 提供的可修改字段（`remark`、`enable_deletion_protection`、`enable_risk_warning`）。

2. **优化不可修改字段验证**：准确识别真正的不可修改字段，移除误报，减少用户的困惑。

3. **增强状态一致性**：在 Update 操作后添加状态等待和检查机制，确保 Terraform state 与云资源状态一致。

4. **改进错误处理**：提供更清晰的错误消息，包含具体的字段信息和 API 响应详情。

5. **保持向后兼容**：确保所有变更都不破坏现有用户的配置和 state。

**Non-Goals:**

1. 不修改 Create 或 Delete 操作的逻辑。
2. 不改变现有字段的 schema 定义（仅新增 Optional 字段）。
3. 不引入新的外部依赖或数据模型变更。
4. 不优化其他 RabbitMQ 资源（user、virtual_host、permission）。

## Decisions

### 1. 新增 Schema 字段

**Decision:** 在 resource schema 中添加三个新的 Optional 字段：
- `remark` (TypeString, Optional) - 实例备注
- `enable_deletion_protection` (TypeBool, Optional, Computed) - 删除保护
- `enable_risk_warning` (TypeBool, Optional, Computed) - 风险提示

**Rationale:**
- 这些字段是 `ModifyRabbitMQVipInstance` API 的标准参数
- 标记为 Computed 以支持从 API 读取当前值
- Optional 确保向后兼容，现有实例不受影响

**Alternatives Considered:**
- *Alternative A:* 使用 `ForceNew` 属性强制重建实例
  - *Reject:* 不符合用户期望，会导致不必要的服务中断
- *Alternative B:* 只在 Create 时支持这些字段
  - *Reject:* 限制了用户更新实例的能力，API 明确支持修改

### 2. 更新不可修改字段列表

**Decision:** 保持当前不可修改字段列表不变，但添加详细的错误消息说明。

**Rationale:**
- 经过验证，当前列表中的字段确实需要在创建时指定，不支持后续修改
- 网络配置（zone_ids、vpc_id、subnet_id）涉及基础设施变更
- 节点规格（node_spec、node_num、storage_size）涉及实例规格变更，需要专门的扩缩容 API
- 计费相关字段（auto_renew_flag、time_span、pay_mode）在实例创建后不可修改
- 公网访问（band_width、enable_public_access）需要网络配置变更，不支持热更新

**Alternatives Considered:**
- *Alternative A:* 移除所有不可修改字段限制
  - *Reject:* 会导致 API 调用失败，用户体验更差
- *Alternative B:* 对每个不可修改字段单独提供详细文档
  - *Accept:* 在错误消息中包含字段的具体说明和可能的解决方案

### 3. Update 状态等待机制

**Decision:** 在 Update 操作后调用 Read 操作，而不是添加额外的状态等待循环。

**Rationale:**
- Read 操作本身包含重试逻辑和状态检查
- 复用现有逻辑可以减少代码重复
- `resourceTencentCloudTdmqRabbitmqVipInstanceRead` 已经实现了适当的状态验证

**Alternatives Considered:**
- *Alternative A:* 在 Update 后添加独立的状态等待循环
  - *Reject:* 代码重复，增加维护成本
- *Alternative B:* 使用 DescribeRabbitMQVipInstance API 进行状态检查
  - *Accept:* 当 Read 操作不够时，可以使用此 API 作为补充

### 4. 错误处理增强

**Decision:** 改进错误消息格式，包含以下信息：
- 失败的字段名称
- API 返回的错误代码和消息
- 不可修改字段的解释
- 可能的解决方案（如需要重建实例）

**Rationale:**
- 清晰的错误消息可以帮助用户快速定位问题
- 减少支持请求和调试时间
- 符合 Terraform Provider 的最佳实践

**Alternatives Considered:**
- *Alternative A:* 使用通用错误消息
  - *Reject:* 用户体验差，无法快速定位问题
- *Alternative B:* 将所有错误转换为特定的错误类型
  - *Reject:* 过度设计，增加复杂度

### 5. 标签更新逻辑

**Decision:** 保持当前的标签更新逻辑，但添加注释说明全量替换机制。

**Rationale:**
- API 要求发送完整的标签列表（非增量更新）
- 当前逻辑已经正确实现了这一点
- 添加注释可以提高代码可读性

**Alternatives Considered:**
- *Alternative A:* 实现增量标签更新
  - *Reject:* API 不支持，需要额外的 API 调用
- *Alternative B:* 使用独立的标签管理资源
  - *Reject:* 增加复杂度，不符合当前架构

## Risks / Trade-offs

### Risk 1: API 字段兼容性

**Risk:** `ModifyRabbitMQVipInstance` API 的某些字段可能在某些地区或版本中不可用。

**Mitigation:**
- 在代码中添加 nil 检查，仅在字段存在时才设置
- 添加文档说明字段可用性的前提条件
- 在 API 调用失败时提供清晰的错误消息

### Risk 2: 向后兼容性破坏

**Risk:** 新增的 schema 字段可能影响现有的 Terraform 配置。

**Mitigation:**
- 所有新字段标记为 Optional，不影响现有配置
- 通过现有的 acceptance tests 验证向后兼容性
- 在文档中清晰说明新字段的默认行为

### Risk 3: 状态不一致

**Risk:** Update 操作后可能出现短暂的 state 不一致。

**Mitigation:**
- 在 Update 后立即调用 Read 操作刷新 state
- 使用重试机制处理临时的 API 不一致性
- 添加日志记录以便排查问题

### Risk 4: 测试覆盖不足

**Risk:** 新增的字段和逻辑可能缺少足够的测试覆盖。

**Mitigation:**
- 为每个新增字段编写独立的单元测试
- 更新现有的 update 测试用例
- 运行完整的 acceptance tests 验证变更

## Migration Plan

### Phase 1: Schema 更新

1. 在 `ResourceTencentCloudTdmqRabbitmqVipInstance()` 函数中添加新的 schema 字段
2. 更新 schema 文档（`resource_tc_tdmq_rabbitmq_vip_instance.md`）
3. 更新网站文档（`website/docs/r/tdmq_rabbitmq_vip_instance.html.markdown`）

### Phase 2: Update 逻辑实现

1. 修改 `resourceTencentCloudTdmqRabbitmqVipInstanceUpdate` 函数：
   - 添加新字段的变更检测和 API 调用
   - 改进不可修改字段的错误消息
   - 移除不必要的状态等待逻辑（依赖 Read 操作）

2. 修改 `resourceTencentCloudTdmqRabbitmqVipInstanceRead` 函数：
   - 添加新字段的读取逻辑
   - 确保 nil 值处理正确

### Phase 3: 测试更新

1. 更新 `resource_tc_tdmq_rabbitmq_vip_instance_test.go`：
   - 添加新字段的测试用例
   - 更新 update 测试用例
   - 添加错误处理的测试用例

2. 运行 acceptance tests 验证所有变更

### Phase 4: 代码格式和文档

1. 执行 `go fmt` 格式化代码
2. 运行 `go vet` 检查代码质量
3. 更新 CHANGELOG（如果需要）

### Rollback Strategy

如果变更导致问题，可以：
1. 通过 git revert 回滚代码变更
2. 保留旧的 schema 定义，确保兼容性
3. 通过 CI/CD 管道验证回滚后的功能

## Open Questions

1. **Q:** `enable_deletion_protection` 和 `enable_risk_warning` 的 API 返回值类型是什么？
   - *Action:* 需要在实现时通过测试或文档确认

2. **Q:** 是否需要在 API 文档中明确说明不可修改字段的原因？
   - *Action:* 是的，在更新网站文档时添加详细说明

3. **Q:** 是否需要为新增的字段添加专门的 validation 逻辑？
   - *Action:* 不需要，API 层已经进行了验证

4. **Q:** Update 操作的超时时间是否需要调整？
   - *Action:* 保持现有的超时配置，除非测试中发现问题
