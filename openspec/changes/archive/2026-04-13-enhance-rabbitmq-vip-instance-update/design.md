## Context

当前 Terraform Provider 的 RabbitMQ VIP 实例资源 (`tencentcloud_tdmq_rabbitmq_vip_instance`) 通过 `ModifyRabbitMQVipInstance` API 支持 update 操作。根据云 API 的最新版本（vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217/models.go），该接口支持更新以下字段：
- `ClusterName`（集群名称）- 已支持
- `Remark`（备注）- 未支持
- `EnableDeletionProtection`（删除保护）- 未支持
- `EnableRiskWarning`（集群风险提示）- 未支持
- `Tags`（标签）- 已支持
- `RemoveAllTags`（删除所有标签）- 已支持

当前实现仅支持 `cluster_name` 和 `resource_tags` 两个字段的更新，其他字段未在 update 逻辑中实现。这限制了用户通过 Terraform 完全管理 RabbitMQ 实例配置的能力。

## Goals / Non-Goals

**Goals:**
- 扩展 RabbitMQ VIP 实例的 update 能力，支持所有云 API 允许更新的字段
- 保持完全向后兼容，不破坏现有 Terraform 配置和 state
- 确保新增字段的正确性和可靠性
- 提供完整的测试覆盖

**Non-Goals:**
- 修改资源的 create/delete 逻辑
- 改变现有字段的行为或类型
- 修改 API 调用方式或错误处理模式
- 添加新的自定义 API 或封装层

## Decisions

### 1. Schema 字段设计

新增三个 Optional 字段到 Resource Schema：
- `remark`: Type: String, Optional, 描述：实例备注信息
- `enable_deletion_protection`: Type: Bool, Optional, 描述：是否开启删除保护
- `enable_risk_warning`: Type: Bool, Optional, 描述：是否开启集群风险提示

**Rationale:** 这些字段直接对应云 API 的参数，类型和语义保持一致。所有字段都是 Optional，确保向后兼容性。

### 2. Update 函数实现策略

在 `resourceTencentCloudTdmqRabbitmqVipInstanceUpdate` 函数中：
- 保持现有的 immutableArgs 列表（不修改，因为这些字段确实不支持通过 API 更新）
- 为三个新字段添加 update 逻辑，使用 `d.HasChange()` 检测变更
- 使用 `helper.String()` 和 `helper.Bool()` 等辅助函数进行类型转换
- 设置 `needUpdate` 标志，只有存在变更时才调用 API

**Rationale:** 遵循现有代码模式，确保一致性和可维护性。只在实际有变更时调用 API，避免不必要的 API 调用。

### 3. Read 函数更新

在 `resourceTencentCloudTdmqRabbitmqVipInstanceRead` 函数中：
- 添加三个新字段的读取逻辑，从 `DescribeTdmqRabbitmqVipInstanceById` 或 `DescribeTdmqRabbitmqVipInstanceByFilter` 返回的数据中获取对应字段
- 使用 `_ = d.Set()` 设置字段值，忽略错误（遵循现有模式）

**Rationale:** 确保 read 操作能够正确反映实例的最新状态，包括新增的字段。

### 4. 测试策略

更新 `resource_tc_tdmq_rabbitmq_vip_instance_test.go`：
- 添加测试用例验证新字段的更新功能
- 测试用例应覆盖正常更新、批量更新、以及不更新时的行为
- 使用 mock API 方式测试逻辑，避免调用真实云 API

**Rationale:** 确保新增功能的正确性，防止回归问题。根据约束，测试文件必须使用 mock 云API的方式，不要在单测中调用真实的云API。

## Risks / Trade-offs

### Risk: 新增字段可能在不同云 API 版本中有不同的行为或限制
**Mitigation:** 使用 vendor 目录中的云 API 定义（v20200217），确保与 provider 依赖的 SDK 版本一致。在文档中注明这些字段的 API 版本要求。

### Risk: 用户可能误以为某些字段可以更新，但实际上云 API 不支持
**Mitigation:** 保持现有 immutableArgs 列表不变，只添加明确支持更新的字段。在错误信息中明确指出哪些字段不可更改。

### Risk: 新增字段可能在某些实例配置下无效或被忽略
**Mitigation:** 在文档中说明各字段的适用场景和限制条件。用户可以通过 read 操作验证字段值是否正确设置。

## Migration Plan

1. **代码变更顺序：**
   - 更新 Schema，新增三个 Optional 字段
   - 更新 Read 函数，添加新字段的读取逻辑
   - 更新 Update 函数，添加新字段的更新逻辑
   - 更新测试文件，添加测试用例

2. **向后兼容性验证：**
   - 新字段均为 Optional，现有配置无需修改
   - 运行现有的 acceptance tests，确保不破坏现有功能

3. **文档更新：**
   - 在变更实现后，通过 `make doc` 生成更新的文档（由 tfpacer-finalize skill 执行）

## Open Questions

无。所有技术决策已经明确，可以直接实施。