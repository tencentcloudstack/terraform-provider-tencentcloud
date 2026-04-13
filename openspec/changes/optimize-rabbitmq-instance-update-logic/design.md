## Context

当前的 RabbitMQ VIP 实例资源 (`tencentcloud_tdmq_rabbitmq_vip_instance`) 的 update 逻辑存在功能限制。虽然腾讯云的 `ModifyRabbitMQVipInstance` API 支持更新多个实例属性，但当前实现仅支持更新 `cluster_name` 和 `resource_tags` 两个参数。

从 vendor 目录下的 API 定义可以看到，`ModifyRabbitMQVipInstanceRequest` 支持以下可更新参数：
- `ClusterName`（已支持）
- `Tags` 和 `RemoveAllTags`（已支持）
- `Remark`（备注信息）- 未支持
- `EnableDeletionProtection`（删除保护）- 未支持
- `EnableRiskWarning`（风险提示）- 未支持

此外，`DescribeRabbitMQVipInstanceResponseParams` 中的 `RabbitMQClusterInfo` 结构也包含了这些字段的数据，说明这些属性是可查询和可管理的。

当前的限制导致用户无法通过 Terraform 完全管理 RabbitMQ 实例的可配置属性，需要在控制台或其他工具中手动配置，降低了自动化程度。

## Goals / Non-Goals

**Goals:**
1. 新增 `remark` 参数到资源 schema，支持实例备注信息的读取和更新
2. 新增 `enable_deletion_protection` 参数到资源 schema，支持删除保护开关的读取和更新
3. 新增 `enable_risk_warning` 参数到资源 schema，支持风险提示开关的读取和更新
4. 更新 Update 函数，将这些新参数正确地传递给 `ModifyRabbitMQVipInstance` API
5. 更新 Read 函数，从 API 响应中读取这些新参数的值并设置到资源状态
6. 确保所有变更保持向后兼容，不影响现有用户配置

**Non-Goals:**
- 不修改任何已有的参数定义或行为
- 不引入新的 API 调用或依赖
- 不改变现有的错误处理逻辑
- 不涉及性能优化或架构重构

## Decisions

### 1. Schema 字段定义

**Decision:** 新增三个 Optional 字段到资源 schema 中。

**Rationale:**
- 所有三个参数在 `ModifyRabbitMQVipInstance` API 中都是可选的（`omitempty`），因此在 Terraform schema 中也应该定义为 Optional
- 这些参数在创建时可能不需要设置（比如 remark 可以为空），但后续可以通过 update 修改
- 使用 Computed 标记允许 API 返回默认值，同时用户可以显式设置

**Alternatives considered:**
- *Computed-only:* 如果设置为 Computed 而非 Optional，用户将无法显式设置这些值，只能读取 API 返回的值。这不符合需求，因为用户需要能够更新这些参数。

### 2. Update 逻辑实现

**Decision:** 在现有的 `resourceTencentCloudTdmqRabbitmqVipInstanceUpdate` 函数中，为每个新参数添加类似 `cluster_name` 的更新逻辑：
- 检测参数是否发生变化（`d.HasChange`）
- 将新值设置到 API request 对象中
- 标记需要调用 update API（`needUpdate = true`）

**Rationale:**
- 保持与现有代码风格一致，降低维护成本
- 使用 `HasChange` 避免不必要的 API 调用，提高效率
- 三个新参数的更新逻辑简单且独立，可以直接集成到现有函数中

**Alternatives considered:**
- *独立的 update 函数:* 为不同类型的参数创建独立的 update 函数。虽然可能更模块化，但这些参数的更新逻辑非常简单，拆分反而增加复杂性。

### 3. Read 逻辑实现

**Decision:** 在 `resourceTencentCloudTdmqRabbitmqVipInstanceRead` 函数中，从 `rabbitmqVipInstance.ClusterInfo` 对象中读取新参数的值：
- `Remark` → `remark`
- `EnableDeletionProtection` → `enable_deletion_protection`
- `EnableRiskWarning` → `enable_risk_warning`

**Rationale:**
- API 响应的 `RabbitMQClusterInfo` 结构中已经包含这些字段
- 直接从响应中读取并设置到资源状态，确保 Terraform state 与云资源状态一致
- 使用 `d.Set()` 方法将值设置到资源状态，忽略可能的错误（按照现有代码模式）

**Alternatives considered:**
- *错误处理增强:* 为每个 `d.Set()` 调用添加错误检查并返回。虽然更严谨，但现有代码模式并不检查 `d.Set()` 的返回值，保持一致即可。

### 4. 不可变参数的处理

**Decision:** 保持现有的 `immutableArgs` 列表不变，不将新参数添加到其中。

**Rationale:**
- 新增的三个参数在 API 中都是可修改的
- 用户需要在实例生命周期中能够调整这些设置
- 这符合这些参数的业务语义（备注、删除保护、风险提示都是可以随时开启/关闭或修改的）

**Alternatives considered:**
- *部分参数不可变:* 考虑将 `enable_deletion_protection` 设为不可变，以防止意外删除。但这限制了用户的灵活性，云 API 本身支持修改，因此不应在 Terraform provider 中添加额外的限制。

### 5. 文档和测试更新

**Decision:** 更新测试文件但不在此阶段生成 website 文档。

**Rationale:**
- 测试文件需要添加测试用例来验证新参数的读写功能
- website 文档通过 `make doc` 命令自动生成，不属于此设计的范围
- 根据 OpenSpec 工作流程，文档生成将在收尾阶段（`tfpacer-finalize` skill）执行

**Alternatives considered:**
- *手动更新文档:* 虽然可以直接修改 `website/docs/r/trabbit.html.md`，但这违反了禁止事项（"禁止直接新增/修改 `website/` 目录下的任何文件，只能在收尾阶段通过 `make doc` 命令来生成"）。

## Risks / Trade-offs

### Risk 1: API 参数变更
[Risk] 腾讯云 API 可能在未来移除或修改 `Remark`、`EnableDeletionProtection` 或 `EnableRiskWarning` 参数
→ **Mitigation:** 使用 vendor 目录下固定的 API SDK 版本，确保短期内 API 兼容性；持续关注 API 变更通知

### Risk 2: 向后兼容性
[Risk] 新增的 Optional 参数可能在某些旧版本的实例上没有值或 API 行为不一致
→ **Mitigation:** 使用 Computed 标记允许 API 返回默认值；确保所有新参数都是 Optional，不强制用户设置

### Risk 3: 测试覆盖不足
[Risk] 测试用例可能无法覆盖所有参数组合和边界情况
→ **Mitigation:** 在设计 tasks.md 时确保包含全面的测试场景；使用 mock 云 API 的方式进行单元测试，不依赖真实的云环境

### Trade-off 1: 代码重复
[Trade-off] Update 函数中的参数检查逻辑存在一定重复，但没有引入抽象层
→ **Rationale:** 这些逻辑简单且独立，过度抽象反而增加复杂性。保持代码简洁，遵循现有项目模式。

### Trade-off 2: 错误处理严格性
[Trade-off] 没有对 `d.Set()` 的返回值进行错误检查
→ **Rationale:** 现有代码模式普遍忽略 `d.Set()` 的错误，因为 schema 字段已定义，不太可能出现运行时错误。保持一致即可。
