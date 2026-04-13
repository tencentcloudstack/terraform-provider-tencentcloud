## Context

当前 RabbitMQ VIP 实例资源（`tencentcloud_tdmq_rabbitmq_vip_instance`）的 update 逻辑过于严格，限制了多个实际上云 API 支持更新的字段。根据腾讯云 TDMQ 云 API 文档，`ModifyRabbitMQVipInstance` 接口支持更新 `ClusterName`（已支持）、`Remark`、`EnableDeletionProtection`、`Tags`（已支持）、`EnableRiskWarning` 等字段。

然而，当前的实现将 `remark`、`enable_deletion_protection`、`enable_risk_warning` 等字段标记为不可变（在 `immutableArgs` 列表中），导致用户无法通过 Terraform 更新这些字段。这些字段在 Create API 中已经存在（如 `EnableDeletionProtection`），在 Read API 中也能获取（如 `Remark`），但 Update API 中却没有正确使用。

## Goals / Non-Goals

**Goals:**
- 新增 `remark` 参数到 schema，支持更新实例备注信息
- 新增 `enable_deletion_protection` 参数到 schema，支持更新删除保护状态
- 新增 `enable_risk_warning` 参数到 schema，支持更新集群风险提示状态
- 更新 Update 函数以支持这些字段的更新
- 更新 Read 函数以正确读取这些字段
- 保持向后兼容，不破坏现有 Terraform 配置和 state

**Non-Goals:**
- 修改现有参数的行为（如 `cluster_name`、`resource_tags`）
- 修改 Create 或 Delete 逻辑
- 新增或修改其他 Terraform 资源

## Decisions

### 1. Schema 设计

**决策：** 在现有 schema 中新增三个 Optional 参数
- `remark`（字符串类型，Computed）：实例备注信息
- `enable_deletion_protection`（布尔类型，Computed）：是否开启删除保护
- `enable_risk_warning`（布尔类型，Computed）：是否开启集群风险提示

**理由：**
- 所有参数均为 Optional，不破坏现有配置
- 设置 Computed 标志，允许从 Read API 中获取默认值
- 与云 API 的 `ModifyRabbitMQVipInstanceRequest` 参数一一对应

### 2. Update 逻辑实现

**决策：** 修改 `resourceTencentCloudTdmqRabbitmqVipInstanceUpdate` 函数
- 将 `remark`、`enable_deletion_protection`、`enable_risk_warning` 从 `immutableArgs` 列表中移除
- 使用 `d.HasChange()` 检测这些参数的变化
- 将变化后的值设置到 `request` 对象中

**理由：**
- 使用 `HasChange()` 是 Terraform Provider SDK v2 的标准做法
- 只在参数变化时才调用云 API，减少不必要的 API 调用
- 与现有的 `cluster_name` 和 `resource_tags` 更新逻辑保持一致

### 3. Read 逻辑实现

**决策：** 修改 `resourceTencentCloudTdmqRabbitmqVipInstanceRead` 函数
- 从 `DescribeTdmqRabbitmqVipInstanceById` 返回的对象中读取 `Remark`、`EnableDeletionProtection` 字段
- 设置到 schema 的相应字段中

**理由：**
- 保持 Terraform state 与云端实际状态同步
- 注意 `EnableRiskWarning` 字段在 Describe API 中可能不存在，需要处理 nil 情况

### 4. Create 逻辑实现

**决策：** 修改 `resourceTencentCloudTdmqRabbitmqVipInstanceCreate` 函数
- 添加 `enable_deletion_protection` 参数到 Create API 请求中（云 API 已支持）

**理由：**
- 允许用户在创建实例时就设置删除保护
- 与云 API 的 `CreateRabbitMQVipInstanceRequest` 保持一致

## Risks / Trade-offs

### Risk 1: `EnableRiskWarning` 字段在 Describe API 中可能不存在
- **风险：** 云 API 的 `DescribeRabbitMQVipInstance` 返回的 `RabbitMQVipInstance` 结构体中可能不包含 `EnableRiskWarning` 字段
- **缓解措施：** 在 Read 函数中添加 nil 检查，如果字段不存在则不设置到 schema 中

### Risk 2: 向后兼容性
- **风险：** 新增参数可能导致旧的 Terraform state 在 `terraform apply` 时出现偏差
- **缓解措施：** 所有新参数均为 Optional 和 Computed，Terraform 会自动从云端读取并更新 state，不会破坏现有配置

### Risk 3: 参数类型和命名一致性
- **风险：** 参数命名与云 API 不一致导致混淆
- **缓解措施：** 严格遵循现有命名规范（snake_case），与云 API 的 PascalCase 保持语义一致

## Migration Plan

### 部署步骤
1. 修改 `resource_tc_tdmq_rabbitmq_vip_instance.go` 文件
2. 运行 `make doc` 生成文档
3. 运行单元测试验证功能
4. 运行验收测试（`TF_ACC=1`）验证完整流程

### 回滚策略
- 如果出现问题，可以通过 git revert 快速回滚代码
- 由于所有新参数均为 Optional，不会影响现有用户

### Open Questions
- 无

## Open Questions

无
