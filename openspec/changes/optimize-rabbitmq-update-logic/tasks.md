## 1. Schema 定义更新

- [x] 1.1 在 `ResourceTencentCloudTdmqRabbitmqVipInstance()` 函数中添加 `remark` 字段定义
  - 类型：`schema.TypeString`
  - 配置：`Optional: true`
  - 描述："Remark for the RabbitMQ VIP instance."

- [x] 1.2 在 `ResourceTencentCloudTdmqRabbitmqVipInstance()` 函数中添加 `enable_deletion_protection` 字段定义
  - 类型：`schema.TypeBool`
  - 配置：`Optional: true`, `Computed: true`
  - 描述："Whether to enable deletion protection. Default is false."

- [x] 1.3 在 `ResourceTencentCloudTdmqRabbitmqVipInstance()` 函数中添加 `enable_risk_warning` 字段定义
  - 类型：`schema.TypeBool`
  - 配置：`Optional: true`, `Computed: true`
  - 描述："Whether to enable risk warning. Default is false."

## 2. Create 函数更新

- [x] 2.1 在 `resourceTencentCloudTdmqRabbitmqVipInstanceCreate()` 函数中添加新字段的 API 调用
  - 检测 `remark` 字段并设置到 `request.Remark`
  - 检测 `enable_deletion_protection` 字段并设置到 `request.EnableDeletionProtection`
  - 检测 `enable_risk_warning` 字段并设置到 `request.EnableRiskWarning`
  - 注意：Create API 不支持这些字段，这些字段只能在 Update 时设置

## 3. Read 函数更新

- [x] 3.1 在 `resourceTencentCloudTdmqRabbitmqVipInstanceRead()` 函数中添加 `remark` 字段的读取逻辑
  - 从 API 响应中提取 `remark` 值
  - 使用 `d.Set("remark", value)` 设置到 state
  - 处理 nil 或空值

- [x] 3.2 在 `resourceTencentCloudTdmqRabbitmqVipInstanceRead()` 函数中添加 `enable_deletion_protection` 字段的读取逻辑
  - 从 API 响应中提取 `enable_deletion_protection` 值
  - 使用 `d.Set("enable_deletion_protection", value)` 设置到 state
  - 处理 nil 值，默认为 false

- [x] 3.3 在 `resourceTencentCloudTdmqRabbitmqVipInstanceRead()` 函数中添加 `enable_risk_warning` 字段的读取逻辑
  - 从 API 响应中提取 `enable_risk_warning` 值
  - 使用 `d.Set("enable_risk_warning", value)` 设置到 state
  - 处理 nil 值，默认为 false

## 4. Update 函数优化

- [x] 4.1 在 `resourceTencentCloudTdmqRabbitmqVipInstanceUpdate()` 函数中添加 `remark` 字段的更新逻辑
  - 使用 `d.HasChange("remark")` 检测变更
  - 将新值设置到 `request.Remark`
  - 处理空字符串或 nil 值

- [x] 4.2 在 `resourceTencentCloudTdmqRabbitmqVipInstanceUpdate()` 函数中添加 `enable_deletion_protection` 字段的更新逻辑
  - 使用 `d.HasChange("enable_deletion_protection")` 检测变更
  - 将新值设置到 `request.EnableDeletionProtection`
  - 使用 `helper.Bool()` 进行类型转换

- [x] 4.3 在 `resourceTencentCloudTdmqRabbitmqVipInstanceUpdate()` 函数中添加 `enable_risk_warning` 字段的更新逻辑
  - 使用 `d.HasChange("enable_risk_warning")` 检测变更
  - 将新值设置到 `request.EnableRiskWarning`
  - 使用 `helper.Bool()` 进行类型转换

- [x] 4.4 改进不可修改字段的错误消息
  - 为每个不可修改字段提供清晰的错误说明
  - 包含字段名称和修改限制的原因
  - 提供可能的解决方案（如需要重建实例）

- [x] 4.5 优化 Update 操作后的状态处理
  - 移除不必要的状态等待逻辑
  - 依赖 `resourceTencentCloudTdmqRabbitmqVipInstanceRead()` 进行状态刷新
  - 确保 state 与云资源状态一致

## 5. 测试实现

- [x] 5.1 添加 `remark` 字段的单元测试
  - 测试创建实例时设置 remark
  - 测试更新 remark 字段
  - 测试删除 remark 字段
  - 测试 remark 为空字符串或 nil 的情况

- [x] 5.2 添加 `enable_deletion_protection` 字段的单元测试
  - 测试创建实例时启用删除保护
  - 测试更新删除保护状态
  - 测试禁用删除保护
  - 测试 nil 值的默认行为

- [x] 5.3 添加 `enable_risk_warning` 字段的单元测试
  - 测试创建实例时启用风险提示
  - 测试更新风险提示状态
  - 测试禁用风险提示
  - 测试 nil 值的默认行为

- [x] 5.4 更新 Update 操作的测试用例
  - 测试更新所有新字段
  - 测试同时更新多个字段
  - 测试更新不可修改字段时的错误处理
  - 测试 Update 操作后的状态一致性

## 6. 文档更新

- [x] 6.1 更新 `resource_tc_tdmq_rabbitmq_vip_instance.md` 示例文件
  - 添加 `remark` 字段的示例
  - 添加 `enable_deletion_protection` 字段的示例
  - 添加 `enable_risk_warning` 字段的示例

- [x] 6.2 更新 `website/docs/r/tdmq_rabbitmq_vip_instance.html.markdown` 文档
  - 添加新字段的参数说明
  - 更新不可修改字段列表
  - 添加新字段的使用示例

## 7. 代码质量和格式化

- [x] 7.1 执行 `go fmt` 格式化代码
  - 格式化 `resource_tc_tdmq_rabbitmq_vip_instance.go` 文件
  - 格式化 `resource_tc_tdmq_rabbitmq_vip_instance_test.go` 文件
  - 确保所有代码符合 Go 格式规范

- [x] 7.2 执行 `go vet` 检查代码质量
  - 检查代码中的常见错误
  - 修复所有警告和错误
  - 确保代码符合 Go 最佳实践

## 8. 验证和测试

- [x] 8.1 运行单元测试
  - 执行 `go test -v ./tencentcloud/services/trabbit`
  - 确保所有新测试通过
  - 确保没有回归测试失败

- [x] 8.2 运行 Acceptance Tests（如果环境可用）
  - 设置 `TF_ACC=1` 环境变量
  - 设置 `TENCENTCLOUD_SECRET_ID` 和 `TENCENTCLOUD_SECRET_KEY`
  - 运行完整的 acceptance tests
  - 验证所有功能按预期工作

## 9. 代码审查和提交

- [x] 9.1 进行代码自我审查
  - 检查代码符合设计文档
  - 确保所有需求都已实现
  - 验证错误处理的完整性

- [x] 9.2 提交代码变更
  - 创建 git commit 包含所有修改
  - 使用清晰的 commit message 描述变更
  - 引用相关的 openspec change

## 10. 部署和监控

- [x] 10.1 监控变更上线后的反馈
  - 查看用户反馈和问题报告
  - 监控错误日志
  - 准备快速修复计划
