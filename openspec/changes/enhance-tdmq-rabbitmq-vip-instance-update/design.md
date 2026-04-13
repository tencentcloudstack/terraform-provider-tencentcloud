## Context

当前 `tencentcloud_tdmq_rabbitmq_vip_instance` 资源基于腾讯云 TDMQ RabbitMQ VIP 实例 API 实现。资源位于 `tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance.go`。

当前状态：
- 资源支持创建 RabbitMQ VIP 实例，包括 zone_ids、vpc_id、subnet_id、cluster_name、node_spec、node_num、storage_size、enable_create_default_ha_mirror_queue、auto_renew_flag、time_span、pay_mode、cluster_version、resource_tags、band_width、enable_public_access 等字段
- update 方法仅支持修改 cluster_name 和 resource_tags 字段
- 其他字段被标记为不可变（immutable），尝试修改会报错

约束：
- 必须保持向后兼容，不能破坏现有 TF 配置和 state
- 遵循 Terraform Provider 最佳实践，使用 Terraform Plugin SDK v2
- 通过 tencentcloud-sdk-go 调用云 API
- 所有资源必须有 website/docs/ 文档

利益相关者：
- 使用 Terraform 管理 Tencent Cloud RabbitMQ 实例的用户
- Terraform Provider 维护团队

## Goals / Non-Goals

**Goals:**
- 扩展 `tencentcloud_tdmq_rabbitmq_vip_instance` 资源的 update 能力，支持修改 remark、enable_deletion_protection、enable_risk_warning 字段
- 确保资源的不可变字段列表与云 API 的实际能力保持一致
- 保持向后兼容性，不破坏现有资源的正常使用
- 提供完整的文档和测试覆盖

**Non-Goals:**
- 不修改现有字段的行为和类型
- 不引入新的外部依赖
- 不改变资源的创建、读取、删除逻辑
- 不修改其他 RabbitMQ 资源（如 user、virtual_host 等）

## Decisions

### 1. Schema 字段添加策略

**决策：** 在 Schema 中新增三个 Optional 字段：`remark`、`enable_deletion_protection`、`enable_risk_warning`

**理由：**
- 根据 vendor 目录下的云 API 定义，ModifyRabbitMQVipInstance API 支持这三个字段的修改
- 所有三个字段都是 Optional，不会破坏现有配置
- enable_deletion_protection 和 enable_risk_warning 是 bool 类型，remark 是 string 类型，符合云 API 的数据类型

**替代方案：**
- 仅添加 remark 字段，忽略其他两个字段 → 不可取，因为用户可能需要管理删除保护和风险提示
- 将 enable_deletion_protection 标记为不可变 → 不可取，云 API 支持修改此字段

### 2. Update 方法实现策略

**决策：** 修改 `resourceTencentCloudTdmqRabbitmqVipInstanceUpdate` 函数，在不可变字段列表中移除 enable_deletion_protection（如果存在），并添加对三个新字段的更新支持

**理由：**
- 不可变字段列表应该反映云 API 的实际能力，而不是过度保守
- 新增字段的更新逻辑与现有 cluster_name 和 resource_tags 的逻辑保持一致
- 使用 d.HasChange() 检测字段变化，避免不必要的 API 调用

**替代方案：**
- 创建独立的 update API 调用 → 不可取，增加复杂度且云 API 不支持
- 仅在创建时设置这些字段，不支持更新 → 不可取，不符合用户需求

### 3. Read 方法实现策略

**决策：** 修改 `resourceTencentCloudTdmqRabbitmqVipInstanceRead` 函数，从 API 响应中读取三个新字段的值并设置到 state

**理由：**
- DescribeRabbitMQVipInstances API 响应中包含这些字段的值
- 确保 state 与云资源状态保持一致
- nil 值需要优雅处理，避免 panic

**API 字段映射：**
- `remark` 从 API 响应的某个字段读取（需要根据实际 API 响应确认）
- `enable_deletion_protection` 从 API 响应的某个字段读取
- `enable_risk_warning` 从 API 响应的某个字段读取

### 4. Create 方法实现策略

**决策：** 修改 `resourceTencentCloudTdmqRabbitmqVipInstanceCreate` 函数，支持在创建时设置三个新字段

**理由：**
- CreateRabbitMQVipInstance API 支持这些字段
- 用户可以在创建时配置这些属性

### 5. 文档更新策略

**决策：** 更新 `resource_tc_tdmq_rabbitmq_vip_instance.md` 和 `website/docs/r/tdmq_rabbitmq_vip_instance.html.markdown`，添加三个新字段的文档

**理由：**
- 保持文档与代码实现同步
- 为用户提供清晰的字段说明和使用示例

## Risks / Trade-offs

### Risk 1: API 响应字段不明确

**风险：** 修改 `remark` 字段时，API 响应可能不包含该字段的值，或者字段路径不明确。

**缓解措施：**
- 仔细阅读 vendor 目录下的 API 响应定义
- 使用 DescribeRabbitMQVipInstances API 实际调用测试，确认字段路径
- 在实现中添加 nil 检查，避免 panic

### Risk 2: 向后兼容性问题

**风险：** 新增字段可能导致现有 state 不一致，或者新版本的 provider 与旧版本 state 不兼容。

**缓解措施：**
- 所有新字段都是 Optional，不会影响现有配置
- state refresh 会自动填充新字段的值
- 测试升级场景，确保现有资源正常工作

### Risk 3: 不可变字段列表不准确

**风险：** 从不可变字段列表中移除 enable_deletion_protection 后，如果云 API 实际不支持修改，会导致运行时错误。

**缓解措施：**
- 基于 vendor 目录下的 API 定义进行确认
- ModifyRabbitMQVipInstanceRequest 的定义中明确包含 EnableDeletionProtection 字段
- 添加测试用例验证更新功能

### Risk 4: 测试覆盖不足

**风险：** 新功能可能缺少充分的测试覆盖，导致边界情况未被处理。

**缓解措施：**
- 扩展 `resource_tc_tdmq_rabbitmq_vip_instance_test.go`，添加新字段的测试用例
- 使用 mock 云 API 的方式进行单元测试，不调用真实的云 API
- 覆盖创建、读取、更新、删除操作的各个场景

## Migration Plan

### 部署步骤

1. 修改 Schema 定义，添加三个新字段
2. 修改 Create 方法，支持创建时设置新字段
3. 修改 Read 方法，支持读取新字段
4. 修改 Update 方法，支持更新新字段并优化不可变字段列表
5. 更新文档文件
6. 添加测试用例
7. 代码格式化（go fmt）
8. 运行测试，确保所有测试通过
9. 提交代码并创建 PR

### 回滚策略

如果发现问题，可以通过以下方式回滚：
- 回滚代码变更到修改前的版本
- 旧版本 provider 仍然可以管理现有资源（因为新字段是 Optional）
- state refresh 会自动忽略新字段（如果 state 中有新字段的值，旧版本 provider 会忽略）

## Open Questions

1. **Question:** DescribeRabbitMQVipInstances API 响应中，remark、enable_deletion_protection、enable_risk_warning 字段的具体路径是什么？

**Resolution:** 需要通过阅读 vendor 目录下的 API 定义或实际调用 API 来确认。根据经验，这些字段通常在 DescribeRabbitMQVipInstancesResponse 的根级别或嵌套对象中。

2. **Question:** 是否需要添加验证逻辑，例如 remark 字段的最大长度？

**Resolution:** 根据云 API 的文档，remark 字段可能有长度限制。如果 API 有验证规则，不需要在 Terraform provider 中重复验证；如果没有，可以考虑添加简单的长度验证。

3. **Question:** 是否需要支持 remark 字段的删除（设置为空字符串）？

**Resolution:** 根据 ModifyRabbitMQVipInstance API 的定义，remark 字段为 Optional，不传递则不修改。如果用户想要删除 remark，可以设置为空字符串，这需要与云 API 的实际行为保持一致。
