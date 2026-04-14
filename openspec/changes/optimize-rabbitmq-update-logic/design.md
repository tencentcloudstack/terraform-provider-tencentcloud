## Context

当前 Terraform Provider 的 RabbitMQ VIP 实例资源 (`tencentcloud_tdmq_rabbitmq_vip_instance`) 支持通过 ModifyRabbitMQVipInstance API 更新 `cluster_name` 和 `resource_tags` 参数。然而，腾讯云的 ModifyRabbitMQVipInstance API 还支持其他参数的更新，包括 `remark`（备注）、`enable_deletion_protection`（是否开启删除保护）和 `enable_risk_warning`（是否开启集群风险提示）。这些参数在当前的 Provider 实现中缺失，导致用户无法通过 Terraform 管理这些重要的实例配置。

当前的实现位于 `tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance.go`，包含完整的 CRUD 操作。资源通过 tencentcloud-sdk-go 调用腾讯云 TDMQ 服务的 API。

### Current State

- **Supported updates**: cluster_name, resource_tags
- **Unsupported updates** (but available in API): remark, enable_deletion_protection, enable_risk_warning
- **Immutable parameters**: zone_ids, vpc_id, subnet_id, node_spec, node_num, storage_size, enable_create_default_ha_mirror_queue, auto_renew_flag, time_span, pay_mode, cluster_version, band_width, enable_public_access
- **API capabilities**: ModifyRabbitMQVipInstance 支持 InstanceId, ClusterName, Remark, EnableDeletionProtection, RemoveAllTags, Tags, EnableRiskWarning 参数

### Constraints

- 必须保持向后兼容性，不能破坏现有 Terraform 配置和 state
- 只能新增 Optional 字段到 schema，不能修改现有字段的属性
- 必须遵循 Terraform Plugin SDK v2 的最佳实践
- 需要确保与云 API 的兼容性，只能调用 API 支持的参数

## Goals / Non-Goals

**Goals:**

1. 新增 `remark` 参数到 schema，支持实例备注的更新和读取
2. 新增 `enable_deletion_protection` 参数到 schema，支持删除保护的更新和读取
3. 新增 `enable_risk_warning` 参数到 schema，支持集群风险提示的更新和读取
4. 在 Update 函数中实现对这三个参数的更新逻辑
5. 在 Read 函数中实现对这三个参数的读取逻辑
6. 保持现有功能完全向后兼容，确保现有配置继续正常工作

**Non-Goals:**

1. 不修改现有不可变参数的验证逻辑
2. 不添加任何云 API 不支持的参数
3. 不修改资源的基本 CRUD 流程或架构
4. 不引入新的外部依赖或数据模型变更
5. 不修改 Timeouts 或其他高级配置

## Decisions

### Schema Design

**Decision**: 将三个新参数都标记为 `Optional` 和 `Computed`

**Rationale**:
- **Optional**: 这些参数在创建实例时可以不设置，后续可以单独更新
- **Computed**: 读取实例时，如果这些参数在云端已设置（通过控制台或其他工具），需要能够读取并显示在 Terraform state 中
- 符合 Terraform Provider 的最佳实践，允许用户灵活管理这些参数

**Alternatives considered**:
- *Option A*: Optional only - 不推荐，因为无法读取云端的现有值
- *Option B*: Required - 不推荐，因为用户可能不想设置这些参数
- *Option C*: Only in Create - 不推荐，因为用户需要在创建后也能更新这些参数

### API Parameter Mapping

**Decision**: 直接使用 tencentcloud-sdk-go 中的字段名称和类型

**Rationale**:
- 保持与 SDK 的一致性，降低维护成本
- 简化代码，减少类型转换的逻辑
- SDK 已经提供了正确的类型定义和序列化逻辑

**Mapping**:
- `remark` → `request.Remark` (string)
- `enable_deletion_protection` → `request.EnableDeletionProtection` (bool)
- `enable_risk_warning` → `request.EnableRiskWarning` (bool)

### Update Logic Implementation

**Decision**: 只有当参数发生变化时才包含在 API 请求中

**Rationale**:
- 避免不必要的 API 调用
- 防止意外覆盖云端的现有值
- 提高更新效率
- 符合 Terraform Provider 的惯用模式

**Implementation**:
```go
if d.HasChange("remark") {
    if v, ok := d.GetOk("remark"); ok {
        request.Remark = helper.String(v.(string))
        needUpdate = true
    } else {
        request.Remark = helper.String("")
        needUpdate = true
    }
}
```

### Read Logic Implementation

**Decision**: 从 DescribeRabbitMQVipInstance API 响应中读取新参数

**Rationale**:
- 确保 Terraform state 与云端状态保持一致
- 支持从云端读取这些参数的初始值
- 允许用户通过其他工具（如控制台）修改这些参数后，通过 Terraform 读取回来

**Implementation**:
需要先检查 DescribeRabbitMQVipInstance API 响应的结构，确定新参数在响应中的字段名称和位置。

### Error Handling

**Decision**: 保持与现有代码相同的错误处理模式

**Rationale**:
- 保持代码一致性
- 复用现有的错误处理逻辑
- 确保错误信息对用户友好

**Pattern**:
- 使用 `defer tccommon.LogElapsed()` 记录执行时间
- 使用 `defer tccommon.InconsistentCheck()` 检查状态一致性
- 使用 `helper.Retry()` 进行重试
- 使用日志记录调试信息

### Testing Strategy

**Decision**: 使用 mock API 进行单元测试，不调用真实云 API

**Rationale**:
- 提高测试速度和稳定性
- 避免产生实际的云资源费用
- 可以模拟各种 API 响应场景
- 符合项目的测试约定

**Test coverage**:
- 更新 remark 参数
- 更新 enable_deletion_protection 参数
- 更新 enable_risk_warning 参数
- 同时更新多个参数
- 读取包含新参数的实例
- 读取不包含新参数的实例

## Risks / Trade-offs

### Risk 1: API 响应字段不明确

**Risk**: DescribeRabbitMQVipInstance API 响应中可能不包含新参数的完整信息，或者字段名称与 Modify API 不同。

**Mitigation**:
- 先查看 vendor 目录下的 API 定义文件，确认响应结构
- 如果响应中不包含这些参数，可能需要使用其他 API 或字段
- 在实现过程中添加日志和错误处理，便于调试

### Risk 2: 向后兼容性问题

**Risk**: 新增 Computed 参数可能导致 Terraform plan 显示意外的 diff，如果用户没有在配置中显式设置这些参数。

**Mitigation**:
- 在文档中明确说明新参数的作用
- 确保只有在云端存在值时才设置到 state
- 测试各种场景，包括创建、更新、导入等

### Risk 3: API 限制或行为变更

**Risk**: 云 API 可能对某些参数有限制（如最大长度、特定值范围），或者在未来版本中改变行为。

**Mitigation**:
- 在代码中添加参数验证逻辑
- 监控云 API 的变更日志
- 在文档中记录已知的 API 限制

### Trade-off: Complexity vs. Flexibility

**Trade-off**: 增加更多可选参数会增加 schema 的复杂性，但提供了更大的灵活性。

**Decision**: 优先考虑灵活性，因为：
- 用户可以根据需要选择是否使用这些参数
- 不影响不使用这些参数的用户
- 符合基础设施即代码的理念，提供完整的资源管理能力

## Migration Plan

由于这是一个纯新增功能的变更（只新增 Optional/Computed 参数），不需要迁移计划。现有配置可以继续正常工作，新参数是可选的。

如果用户已经在云端设置了这些参数（通过控制台或其他工具），Terraform 在下次 Read 操作时会自动读取这些值到 state 中，用户可以选择是否在 Terraform 配置中管理这些参数。

## Open Questions

1. **API 响应字段位置**: DescribeRabbitMQVipInstance API 响应中，新参数是否在顶层结构中，还是在嵌套的结构（如 ClusterInfo）中？需要在实现时确认。

2. **参数默认值**: 如果云端未设置这些参数，API 响应会返回什么值（nil、空字符串、false）？需要在实现时处理这些情况。

3. **API 限制**: 新参数是否有任何 API 层面的限制（如 remark 的最大长度、enable_deletion_protection 的使用条件）？需要在实现前查阅云 API 文档。
