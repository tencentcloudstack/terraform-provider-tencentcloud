# 优化 RabbitMQ 实例 update 逻辑 - 规范定义

## ADDED Requirements

### Requirement: RabbitMQ 实例应支持通过 API 更新节点配置

#### Scenario: 用户需要扩容 RabbitMQ 实例的节点数量

当用户通过 Terraform 修改 `tencentcloud_tdmq_rabbitmq_vip_instance` 资源的 `node_num` 字段时：

**前置条件**：
- 已存在一个 RabbitMQ VIP 实例
- 实例状态为运行中（Success）
- 目标节点数量符合腾讯云限制规则

**操作步骤**：
1. 用户在 Terraform 配置中将 `node_num` 从 3 修改为 5
2. 执行 `terraform apply`
3. Provider 调用 `ModifyRabbitMQVipInstanceSpec` API
4. API 返回成功，实例状态变为 "Running"（更新中）
5. Provider 轮询实例状态，等待状态变为 "Success"

**预期结果**：
- Terraform 显示 "Modifying instance node count from 3 to 5"
- 调用 API `ModifyRabbitMQVipInstanceSpec`，参数包含 `InstanceId` 和 `NodeNum=5`
- 等待实例状态变为 "Success" 后完成更新
- 不需要删除并重新创建资源
- 资源状态更新后，Read 函数返回新的 `node_num` 值

**验证点**：
- API 调用参数正确
- 实例状态检查逻辑正确
- Update 成功后状态正确读取
- 不会因 `node_num` 变更而触发资源重建

---

### Requirement: RabbitMQ 实例应支持通过 API 更新存储规格

#### Scenario: 用户需要扩容 RabbitMQ 实例的存储空间

当用户通过 Terraform 修改 `tencentcloud_tdmq_rabbitmq_vip_instance` 资源的 `storage_size` 字段时：

**前置条件**：
- 已存在一个 RabbitMQ VIP 实例
- 实例状态为运行中（Success）
- 目标存储大小符合腾讯云限制规则

**操作步骤**：
1. 用户在 Terraform 配置中将 `storage_size` 从 200 修改为 500
2. 执行 `terraform apply`
3. Provider 调用 `ModifyRabbitMQVipInstanceSpec` API
4. API 返回成功，实例状态变为 "Running"（更新中）
5. Provider 轮询实例状态，等待状态变为 "Success"

**预期结果**：
- Terraform 显示 "Modifying instance storage size from 200GB to 500GB"
- 调用 API `ModifyRabbitMQVipInstanceSpec`，参数包含 `InstanceId` 和 `StorageSize=500`
- 等待实例状态变为 "Success" 后完成更新
- 不需要删除并重新创建资源
- 资源状态更新后，Read 函数返回新的 `storage_size` 值

**验证点**：
- API 调用参数正确
- 实例状态检查逻辑正确
- Update 成功后状态正确读取
- 不会因 `storage_size` 变更而触发资源重建

---

### Requirement: RabbitMQ 实例应支持通过 API 更新节点规格

#### Scenario: 用户需要升级 RabbitMQ 实例的节点规格

当用户通过 Terraform 修改 `tencentcloud_tdmq_rabbitmq_vip_instance` 资源的 `node_spec` 字段时：

**前置条件**：
- 已存在一个 RabbitMQ VIP 实例
- 实例状态为运行中（Success）
- 目标规格在腾讯云可选规格列表中

**操作步骤**：
1. 用户在 Terraform 配置中将 `node_spec` 从 "rabbit-vip-basic-1" 修改为 "rabbit-vip-profession-4c16g"
2. 执行 `terraform apply`
3. Provider 调用 `ModifyRabbitMQVipInstanceSpec` API
4. API 返回成功，实例状态变为 "Running"（更新中）
5. Provider 轮询实例状态，等待状态变为 "Success"

**预期结果**：
- Terraform 显示 "Modifying instance node spec from rabbit-vip-basic-1 to rabbit-vip-profession-4c16g"
- 调用 API `ModifyRabbitMQVipInstanceSpec`，参数包含 `InstanceId` 和 `NodeSpec=rabbit-vip-profession-4c16g`
- 等待实例状态变为 "Success" 后完成更新
- 不需要删除并重新创建资源
- 资源状态更新后，Read 函数返回新的 `node_spec` 值

**验证点**：
- API 调用参数正确
- 实例状态检查逻辑正确
- Update 成功后状态正确读取
- 不会因 `node_spec` 变更而触发资源重建

---

### Requirement: RabbitMQ 实例应支持通过 API 更新公网访问配置

#### Scenario: 用户需要开启或关闭 RabbitMQ 实例的公网访问

当用户通过 Terraform 修改 `tencentcloud_tdmq_rabbitmq_vip_instance` 资源的 `enable_public_access` 字段时：

**前置条件**：
- 已存在一个 RabbitMQ VIP 实例
- 实例状态为运行中（Success）

**操作步骤**：
1. 用户在 Terraform 配置中将 `enable_public_access` 从 false 修改为 true
2. 执行 `terraform apply`
3. Provider 调用 `ModifyRabbitMQVipInstancePublicAccess` API
4. API 返回成功

**预期结果**：
- Terraform 显示 "Enabling public access for instance"
- 调用 API `ModifyRabbitMQVipInstancePublicAccess`，参数包含 `InstanceId` 和 `EnablePublicAccess=true`
- 更新完成后，Read 函数返回新的公网访问配置
- 不需要删除并重新创建资源

**验证点**：
- API 调用参数正确
- Update 成功后状态正确读取
- 不会因 `enable_public_access` 变更而触发资源重建

---

### Requirement: RabbitMQ 实例应支持通过 API 更新公网带宽

#### Scenario: 用户需要调整 RabbitMQ 实例的公网带宽

当用户通过 Terraform 修改 `tencentcloud_tdmq_rabbitmq_vip_instance` 资源的 `band_width` 字段时：

**前置条件**：
- 已存在一个 RabbitMQ VIP 实例
- 实例状态为运行中（Success）
- 目标带宽符合腾讯云限制规则
- 公网访问已启用

**操作步骤**：
1. 用户在 Terraform 配置中将 `band_width` 从 3 修改为 10
2. 执行 `terraform apply`
3. Provider 调用 `ModifyRabbitMQVipInstancePublicAccess` API
4. API 返回成功

**预期结果**：
- Terraform 显示 "Modifying public bandwidth from 3Mbps to 10Mbps"
- 调用 API `ModifyRabbitMQVipInstancePublicAccess`，参数包含 `InstanceId` 和 `Bandwidth=10`
- 更新完成后，Read 函数返回新的带宽配置
- 不需要删除并重新创建资源

**验证点**：
- API 调用参数正确
- Update 成功后状态正确读取
- 不会因 `band_width` 变更而触发资源重建

---

### Requirement: RabbitMQ 实例 Update 函数应支持多字段同时更新

#### Scenario: 用户同时修改多个可更新字段

当用户在同一个 `terraform apply` 操作中同时修改多个可更新字段时：

**前置条件**：
- 已存在一个 RabbitMQ VIP 实例
- 实例状态为运行中（Success）

**操作步骤**：
1. 用户在 Terraform 配置中同时修改：
   - `cluster_name`: "old-cluster" → "new-cluster"
   - `node_num`: 3 → 5
   - `enable_public_access`: false → true
2. 执行 `terraform apply`

**预期结果**：
- 按顺序调用三个不同的 API：
  1. `ModifyRabbitMQVipInstance` - 修改集群名称
  2. `ModifyRabbitMQVipInstanceSpec` - 修改节点数量（需要等待）
  3. `ModifyRabbitMQVipInstancePublicAccess` - 修改公网访问
- 如果任一 API 调用失败，整个 Update 操作失败并返回错误
- 所有更新成功后，Read 函数返回所有字段的新值
- 不需要删除并重新创建资源

**验证点**：
- API 调用顺序正确
- 错误处理逻辑正确（失败时停止后续操作）
- 所有字段最终状态正确
- 不会因多字段变更而触发资源重建

---

### Requirement: RabbitMQ 实例应拒绝修改不可变字段

#### Scenario: 用户尝试修改不可变字段

当用户尝试修改标记为不可变的字段时：

**前置条件**：
- 已存在一个 RabbitMQ VIP 实例

**操作步骤**：
1. 用户在 Terraform 配置中尝试修改 `vpc_id` 字段
2. 执行 `terraform apply`

**预期结果**：
- Terraform 显示错误信息："argument `vpc_id` cannot be changed"
- 不调用任何 API
- 不修改资源状态
- 提示用户某些字段在创建后不可修改

**验证点**：
- 不可变字段列表正确
- 错误信息清晰明确
- 不会调用 API

---

### Requirement: RabbitMQ 实例规格变更后应等待实例状态恢复

#### Scenario: 规格变更后等待实例状态变为 Success

当调用 `ModifyRabbitMQVipInstanceSpec` API 修改规格后：

**前置条件**：
- 已存在一个 RabbitMQ VIP 实例
- 执行了规格变更操作

**操作步骤**：
1. 调用 `ModifyRabbitMQVipInstanceSpec` API 成功
2. 实例状态变为 "Running"
3. Provider 轮询 `DescribeRabbitMQVipInstances` API 查询状态

**预期结果**：
- 使用 `resource.Retry` 进行重试等待
- 重试超时时间设置为 `ReadRetryTimeout * 10`
- 每次重试查询实例状态
- 当状态为 "Success" 时返回成功
- 如果超时仍未变为 "Success"，返回错误
- 如果状态为其他非法状态，立即返回错误

**验证点**：
- 重试逻辑正确
- 超时时间设置合理
- 状态判断正确
- 错误处理正确

---

### Requirement: RabbitMQ 实例 Update 函数应正确处理 API 错误

#### Scenario: API 调用失败时的错误处理

当 Update 操作中的任一 API 调用失败时：

**前置条件**：
- 已存在一个 RabbitMQ VIP 实例

**操作步骤**：
1. 用户修改配置并执行 `terraform apply`
2. 某个 API 调用失败（如配额不足、参数错误等）

**预期结果**：
- 使用 `resource.Retry` 处理临时性错误
- 如果重试后仍失败，返回错误给用户
- 错误信息包含详细的失败原因
- 不执行后续的 Update 操作
- 资源状态保持不变（或部分更新已提交的操作）

**验证点**：
- 错误重试逻辑正确
- 错误信息清晰
- 不会执行后续操作
- 状态一致性

---

### Requirement: RabbitMQ 实例 Update 函数应在完成后刷新状态

#### Scenario: Update 成功后调用 Read 函数刷新状态

当所有 Update 操作成功完成后：

**前置条件**：
- 已存在一个 RabbitMQ VIP 实例
- 执行了 Update 操作

**操作步骤**：
1. 所有 API 调用成功
2. Update 函数结束前调用 `resourceTencentCloudTdmqRabbitmqVipInstanceRead(d, meta)`

**预期结果**：
- Read 函数重新查询实例状态
- 更新 Terraform state 中的所有字段
- 用户看到的 `terraform show` 输出反映最新的配置
- 后续的 `terraform plan` 不会显示这些字段的变更

**验证点**：
- Update 成功后必定调用 Read
- Read 函数正确查询最新状态
- State 正确更新
- 不会产生虚假的变更计划
