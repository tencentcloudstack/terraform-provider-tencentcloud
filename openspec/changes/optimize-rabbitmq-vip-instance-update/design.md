## Context

当前 `tencentcloud_tdmq_rabbitmq_vip_instance` 资源在 Update 操作中存在两个关键问题：
1. 云 API 支持修改的字段（Remark、EnableDeletionProtection、EnableRiskWarning）在 Terraform Provider 中未被实现
2. 虽然云 API 的 `ModifyRabbitMQVipInstance` 接口支持这些字段的修改，但当前的 Schema 定义和 Update 逻辑未能充分利用这些能力

通过分析 `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217/` 目录下的云 API 定义：
- `ModifyRabbitMQVipInstanceRequest` 接口支持修改：ClusterName、Remark、EnableDeletionProtection、RemoveAllTags、Tags、EnableRiskWarning
- `CreateRabbitMQVipInstanceRequest` 接口不支持 Remark、EnableDeletionProtection、EnableRiskWarning 参数
- `DescribeRabbitMQVipInstance` 和 `DescribeRabbitMQVipInstances` 接口响应包含 Remark、EnableDeletionProtection、EnableRiskWarning 字段

## Goals / Non-Goals

**Goals:**
1. 在 Schema 中添加 remark、enable_deletion_protection、enable_risk_warning 字段的定义
2. 在 Read 函数中从云 API 响应中读取这些字段并设置到 state
3. 在 Update 函数中支持通过 ModifyRabbitMQVipInstance 接口更新这些字段
4. 完善 resource_tags 的更新逻辑，正确处理 RemoveAllTags 标志

**Non-Goals:**
1. 不修改现有字段的行为或定义（保持向后兼容）
2. 不添加 CreateRabbitMQVipInstanceRequest 接口不支持的字段
3. 不修改不可变参数的逻辑（zone_ids、vpc_id、subnet_id、node_spec、node_num、storage_size、enable_create_default_ha_mirror_queue、auto_renew_flag、time_span、pay_mode、cluster_version、band_width、enable_public_access）

## Decisions

### 1. Schema 字段定义
- 所有新增字段均设置为 `Optional: true, Computed: true`，确保向后兼容
- 字段命名遵循 Terraform Provider 的命名约定（snake_case）
- 不设置默认值，让用户明确指定

### 2. Read 函数实现
- 从 `DescribeTdmqRabbitmqVipInstanceById` 的 `RabbitMQClusterInfo` 响应中读取 Remark、EnableDeletionProtection、EnableRiskWarning
- 从 `DescribeTdmqRabbitmqVipInstanceByFilter` 的 `RabbitMQVipInstance` 响应中读取 Remark、EnableDeletionProtection（注意：此接口响应不包含 EnableRiskWarning）
- 优先使用 `DescribeTdmqRabbitmqVipInstanceById` 获取完整信息

### 3. Update 函数实现
- 在 `immutableArgs` 列表中保留现有的不可变参数，不添加新字段
- 当检测到 remark、enable_deletion_protection、enable_risk_warning 字段变化时，调用 ModifyRabbitMQVipInstance 接口
- 优化 resource_tags 的更新逻辑：
  - 当 resource_tags 为非空时，设置 Tags 参数
  - 当 resource_tags 为空时，设置 RemoveAllTags 参数为 true

### 4. Cloud API 映射
- Schema 字段 `remark` → API 字段 `Remark`
- Schema 字段 `enable_deletion_protection` → API 字段 `EnableDeletionProtection`
- Schema 字段 `enable_risk_warning` → API 字段 `EnableRiskWarning`

## Risks / Trade-offs

### Risk 1: EnableRiskWarning 字段在 DescribeRabbitMQVipInstances 中不可用
**影响**: 如果用户只查询实例列表而不查询单个实例详情，可能无法获取到 EnableRiskWarning 字段的值

**缓解措施**:
- 在 Read 函数中优先使用 `DescribeTdmqRabbitmqVipInstanceById` 获取完整信息
- 文档中说明该字段的读取行为

### Risk 2: 新增字段在 Create 阶段无法设置
**影响**: 用户无法在创建实例时设置这些字段，必须在创建后通过 update 操作设置

**缓解措施**:
- 在文档中明确说明这些字段只能在创建后更新
- 这是云 API 的限制，无法在 Provider 层面解决

### Risk 3: 向后兼容性
**影响**: 新增字段可能导致某些场景下的行为变化

**缓解措施**:
- 所有字段均设置为 `Optional: true, Computed: true`
- 不修改现有字段的行为
- 更新时只检查并修改变化的新字段
