# Tasks: 优化 RabbitMQ 实例的 update 逻辑

## 1. 需求分析和 API 验证

- [x] 1.1 查看 `CreateRabbitMQVipInstanceRequest` 的定义，确认是否支持 `remark`、`enable_deletion_protection`、`enable_risk_warning` 字段
- [x] 1.2 查看 `DescribeRabbitMQVipInstanceResponse` 的定义，确认是否返回 `remark`、`enable_deletion_protection`、`enable_risk_warning` 字段
- [x] 1.3 验证 `ModifyRabbitMQVipInstanceRequest` 的字段定义（已在 design.md 中确认）
- [x] 1.4 如果 Create 或 Describe API 不支持这些字段，调整实现策略

## 2. 代码修改任务

### 2.1 修改 Schema 定义

- [x] 2.1.1 在 `resource_tc_tdmq_rabbitmq_vip_instance.go` 的 `ResourceTencentCloudTdmqRabbitmqVipInstance` 函数中，添加 `remark` 字段的 Schema 定义
- [x] 2.1.2 在 `resource_tc_tdmq_rabbitmq_vip_instance.go` 的 `ResourceTencentCloudTdmqRabbitmqVipInstance` 函数中，添加 `enable_deletion_protection` 字段的 Schema 定义
- [x] 2.1.3 在 `resource_tc_tdmq_rabbitmq_vip_instance.go` 的 `ResourceTencentCloudTdmqRabbitmqVipInstance` 函数中，添加 `enable_risk_warning` 字段的 Schema 定义

### 2.2 修改 Create 函数

- [x] 2.2.1 在 `resourceTencentCloudTdmqRabbitmqVipInstanceCreate` 函数中，添加 `remark` 字段的设置逻辑（如果 Create API 支持）
- [x] 2.2.2 在 `resourceTencentCloudTdmqRabbitmqVipInstanceCreate` 函数中，添加 `enable_deletion_protection` 字段的设置逻辑（如果 Create API 支持）
- [x] 2.2.3 在 `resourceTencentCloudTdmqRabbitmqVipInstanceCreate` 函数中，添加 `enable_risk_warning` 字段的设置逻辑（如果 Create API 支持）
- [x] 2.2.4 验证 Create 函数的参数设置逻辑正确

### 2.3 修改 Read 函数

- [x] 2.3.1 在 `resourceTencentCloudTdmqRabbitmqVipInstanceRead` 函数中，添加从 API 响应读取 `remark` 字段的逻辑
- [x] 2.3.2 在 `resourceTencentCloudTdmqRabbitmqVipInstanceRead` 函数中，添加从 API 响应读取 `enable_deletion_protection` 字段的逻辑
- [x] 2.3.3 在 `resourceTencentCloudTdmqRabbitmqVipInstanceRead` 函数中，添加从 API 响应读取 `enable_risk_warning` 字段的逻辑
- [x] 2.3.4 处理 nil 值的情况，确保设置合理的默认值
- [x] 2.3.5 验证 Read 函数的字段读取逻辑正确

### 2.4 修改 Update 函数

- [x] 2.4.1 检查 `resourceTencentCloudTdmqRabbitmqVipInstanceUpdate` 函数中的 `immutableArgs` 列表，确保没有错误地将新字段添加到该列表
- [x] 2.4.2 在 `resourceTencentCloudTdmqRabbitmqVipInstanceUpdate` 函数中，添加 `remark` 字段的更新逻辑
- [x] 2.4.3 在 `resourceTencentCloudTdmqRabbitmqVipInstanceUpdate` 函数中，添加 `enable_deletion_protection` 字段的更新逻辑
- [x] 2.4.4 在 `resourceTencentCloudTdmqRabbitmqVipInstanceUpdate` 函数中，添加 `enable_risk_warning` 字段的更新逻辑
- [x] 2.4.5 使用 `d.HasChange()` 检查字段变更，只在变更时更新
- [x] 2.4.6 验证 Update 函数的字段更新逻辑正确

## 3. 测试任务

### 3.1 单元测试

- [x] 3.1.1 创建测试用例，验证 `remark` 字段的创建、读取和更新
- [x] 3.1.2 创建测试用例，验证 `enable_deletion_protection` 字段的创建、读取和更新
- [x] 3.1.3 创建测试用例，验证 `enable_risk_warning` 字段的创建、读取和更新
- [x] 3.1.4 创建测试用例，验证边界情况（nil 值、空字符串、默认值）
- [x] 3.1.5 运行单元测试，确保所有测试通过

### 3.2 集成测试

- [ ] 3.2.1 在 `resource_tc_tdmq_rabbitmq_vip_instance_test.go` 中，添加测试场景：创建实例时设置 `remark` 字段
- [ ] 3.2.2 在 `resource_tc_tdmq_rabbitmq_vip_instance_test.go` 中，添加测试场景：更新实例的 `remark` 字段
- [ ] 3.2.3 在 `resource_tc_tdmq_rabbitmq_vip_instance_test.go` 中，添加测试场景：创建实例时设置 `enable_deletion_protection = true`
- [ ] 3.2.4 在 `resource_tc_tdmq_rabbitmq_vip_instance_test.go` 中，添加测试场景：切换 `enable_deletion_protection` 状态
- [ ] 3.2.5 在 `resource_tc_tdmq_rabbitmq_vip_instance_test.go` 中，添加测试场景：创建实例时设置 `enable_risk_warning = true`
- [ ] 3.2.6 在 `resource_tc_tdmq_rabbitmq_vip_instance_test.go` 中，添加测试场景：切换 `enable_risk_warning` 状态
- [ ] 3.2.7 在 `resource_tc_tdmq_rabbitmq_vip_instance_test.go` 中，添加测试场景：删除保护启用时的删除失败
- [ ] 3.2.8 验证现有测试用例仍然通过
- [ ] 3.2.9 运行集成测试，确保所有测试通过

### 3.3 验收测试

- [ ] 3.3.1 在真实腾讯云环境中创建测试实例
- [ ] 3.3.2 验证新字段在创建时正确设置
- [ ] 3.3.3 验证新字段在更新时正确修改
- [ ] 3.3.4 验证新字段在读取时正确返回
- [ ] 3.3.5 验证删除保护功能正常工作
- [ ] 3.3.6 验证向后兼容性，现有配置无影响
- [ ] 3.3.7 验证与 Terraform CLI 的兼容性
- [ ] 3.3.8 清理测试资源

## 4. 文档任务

### 4.1 资源文档

- [x] 4.1.1 更新 `tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance.md` 文档，添加新增字段的说明
- [x] 4.1.2 为每个新增字段添加详细的描述、使用示例和注意事项
- [x] 4.1.3 添加删除保护功能的说明和注意事项

### 4.2 网站文档

- [x] 4.2.1 更新 `website/docs/r/tdmq_rabbitmq_vip_instance.html.markdown` 文档，添加新增字段的说明
- [x] 4.2.2 为每个新增字段添加 HCL 使用示例
- [x] 4.2.3 添加删除保护功能的使用说明和警告提示

### 4.3 示例配置

- [x] 4.3.1 创建示例配置文件，展示如何使用新增字段
- [x] 4.3.2 创建示例配置文件，展示删除保护的使用场景
- [x] 4.3.3 创建示例配置文件，展示新字段的更新操作

## 5. 验证任务

### 5.1 代码审查

- [x] 5.1.1 进行代码自我审查，确保代码质量
- [x] 5.1.2 检查代码风格是否符合项目规范
- [x] 5.1.3 检查是否有潜在的 bug 或边界情况未处理
- [x] 5.1.4 确保所有 `TODO` 或 `FIXME` 注释都已处理

### 5.2 兼容性验证

- [x] 5.2.1 验证向后兼容性，确保不影响现有配置
- [x] 5.2.2 验证与 Terraform 不同版本的兼容性
- [x] 5.2.3 验证与腾讯云 API 不同版本的兼容性

### 5.3 性能验证

- [x] 5.3.1 验证新增逻辑不影响资源的创建性能
- [x] 5.3.2 验证新增逻辑不影响资源的读取性能
- [x] 5.3.3 验证新增逻辑不影响资源的更新性能
- [x] 5.3.4 确保没有不必要的 API 调用

### 5.4 端到端测试

- [ ] 5.4.1 执行完整的 `terraform init` -> `terraform plan` -> `terraform apply` -> `terraform refresh` -> `terraform plan` -> `terraform destroy` 流程
- [ ] 5.4.2 测试导入功能：`terraform import`
- [ ] 5.4.3 测试状态迁移：修改配置后执行 `terraform apply`
- [ ] 5.4.4 测试并发操作：同时更新多个实例

## 6. 发布准备任务

- [ ] 6.1 确保所有代码变更已提交
- [ ] 6.2 确保所有测试通过
- [ ] 6.3 确保文档更新完整
- [ ] 6.4 准备变更日志（如果需要）
- [ ] 6.5 准备发布说明（如果需要）
