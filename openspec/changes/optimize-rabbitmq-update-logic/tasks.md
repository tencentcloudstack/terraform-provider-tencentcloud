## 1. Schema 定义修改

- [x] 1.1 在 resource_tc_tdmq_rabbitmq_vip_instance.go 的 Schema 中添加 `remark` 字段
  - Type: schema.TypeString
  - Optional: true
  - Computed: true
  - Description: "Instance remark or notes."

- [x] 1.2 在 resource_tc_tdmq_rabbitmq_vip_instance.go 的 Schema 中添加 `enable_deletion_protection` 字段
  - Type: schema.TypeBool
  - Optional: true
  - Computed: true
  - Description: "Whether to enable deletion protection for the instance."

- [x] 1.3 在 resource_tc_tdmq_rabbitmq_vip_instance.go 的 Schema 中添加 `enable_risk_warning` 字段
  - Type: schema.TypeBool
  - Optional: true
  - Computed: true
  - Description: "Whether to enable cluster risk warning."

## 2. Create 方法修改

- [x] 2.1 在 resourceTencentCloudTdmqRabbitmqVipInstanceCreate 方法中添加 `remark` 参数处理
  - 注意：CreateRabbitMQVipInstance API 不支持 remark 参数，只能在 Update 时设置

- [x] 2.2 在 resourceTencentCloudTdmqRabbitmqVipInstanceCreate 方法中添加 `enable_deletion_protection` 参数处理

- [x] 2.3 在 resourceTencentCloudTdmqRabbitmqVipInstanceCreate 方法中添加 `enable_risk_warning` 参数处理
  - 注意：CreateRabbitMQVipInstance API 不支持 enable_risk_warning 参数，只能在 Update 时设置

## 3. Read 方法修改

- [x] 3.1 在 resourceTencentCloudTdmqRabbitmqVipInstanceRead 方法中添加 `remark` 参数读取

- [x] 3.2 在 resourceTencentCloudTdmqRabbitmqVipInstanceRead 方法中添加 `enable_deletion_protection` 参数读取

- [x] 3.3 在 resourceTencentCloudTdmqRabbitmqVipInstanceRead 方法中添加 `enable_risk_warning` 参数读取

## 4. Update 方法修改

- [x] 4.1 从 Update 方法的不可变参数列表 (immutableArgs) 中移除 `remark`
  - 确认 `remark` 不在 immutableArgs 列表中

- [x] 4.2 从 Update 方法的不可变参数列表 (immutableArgs) 中移除 `enable_deletion_protection`
  - 确认 `enable_deletion_protection` 不在 immutableArgs 列表中

- [x] 4.3 从 Update 方法的不可变参数列表 (immutableArgs) 中移除 `enable_risk_warning`
  - 确认 `enable_risk_warning` 不在 immutableArgs 列表中

- [x] 4.4 在 resourceTencentCloudTdmqRabbitmqVipInstanceUpdate 方法中添加 `remark` 参数更新逻辑

- [x] 4.5 在 resourceTencentCloudTdmqRabbitmqVipInstanceUpdate 方法中添加 `enable_deletion_protection` 参数更新逻辑

- [x] 4.6 在 resourceTencentCloudTdmqRabbitmqVipInstanceUpdate 方法中添加 `enable_risk_warning` 参数更新逻辑

## 5. 单元测试更新

- [ ] 5.1 在 resource_tc_tdmq_rabbitmq_vip_instance_test.go 中添加 `remark` 参数的单元测试
  - 创建测试用例 `TestAccTencentCloudTdmqRabbitmqVipInstance_remark`
  - 测试创建带 remark 的实例
  - 测试更新 remark 参数
  - 使用 mock 云 API 避免调用真实服务
  - 注意：测试任务将在收尾阶段执行，此处略过

- [ ] 5.2 在 resource_tc_tdmq_rabbitmq_vip_instance_test.go 中添加 `enable_deletion_protection` 参数的单元测试
  - 创建测试用例 `TestAccTencentCloudTdmqRabbitmqVipInstance_enableDeletionProtection`
  - 测试创建带 enable_deletion_protection 的实例
  - 测试更新 enable_deletion_protection 参数
  - 使用 mock 云 API 避免调用真实服务
  - 注意：测试任务将在收尾阶段执行，此处略过

- [ ] 5.3 在 resource_tc_tdmq_rabbitmq_vip_instance_test.go 中添加 `enable_risk_warning` 参数的单元测试
  - 创建测试用例 `TestAccTencentCloudTdmqRabbitmqVipInstance_enableRiskWarning`
  - 测试创建带 enable_risk_warning 的实例
  - 测试更新 enable_risk_warning 参数
  - 使用 mock 云 API 避免调用真实服务
  - 注意：测试任务将在收尾阶段执行，此处略过

- [ ] 5.4 添加多参数同时更新的单元测试
  - 创建测试用例 `TestAccTencentCloudTdmqRabbitmqVipInstance_updateMultipleFields`
  - 测试同时更新 remark、enable_deletion_protection、enable_risk_warning
  - 验证所有参数在一次 API 调用中更新
  - 使用 mock 云 API 避免调用真实服务
  - 注意：测试任务将在收尾阶段执行，此处略过

## 6. 文档和示例更新

- [x] 6.1 更新 resource_tc_tdmq_rabbitmq_vip_instance.md 示例文件
  - 添加包含 remark 参数的示例
  - 添加包含 enable_deletion_protection 参数的示例
  - 添加包含 enable_risk_warning 参数的示例
  - 添加多参数更新的示例

- [ ] 6.2 验证文档生成
  - 运行 `make doc` 命令生成 website/docs/r/tdmq_rabbitmq_vip_instance.html.markdown
  - 检查生成的文档是否包含新参数的描述
  - 确认文档格式正确
  - 注意：文档生成任务将在收尾阶段通过 make doc 执行，此处略过

## 7. 代码格式化和验证

- [ ] 7.1 执行 go fmt 格式化修改的代码
  - 运行 `go fmt ./tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance.go`
  - 运行 `go fmt ./tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance_test.go`
  - 确保所有代码符合 Go 标准格式
  - 注意：格式化任务将在收尾阶段通过 gofmt 执行，此处略过

- [x] 7.2 验证代码编译
  - 运行 `go build ./tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance.go`
  - 确保没有编译错误

## 8. 集成测试（可选）

- [ ] 8.1 运行单元测试验证功能
  - 运行 `go test -v ./tencentcloud/services/trabbit -run TestAccTencentCloudTdmqRabbitmqVipInstance_remark`
  - 运行 `go test -v ./tencentcloud/services/trabbit -run TestAccTencentCloudTdmqRabbitmqVipInstance_enableDeletionProtection`
  - 运行 `go test -v ./tencentcloud/services/trabbit -run TestAccTencentCloudTdmqRabbitmqVipInstance_enableRiskWarning`
  - 运行 `go test -v ./tencentcloud/services/trabbit -run TestAccTencentCloudTdmqRabbitmqVipInstance_updateMultipleFields`
  - 确保所有测试通过
  - 注意：测试任务在收尾阶段执行，此处略过

- [ ] 8.2 验证向后兼容性
  - 确保不设置新参数时，现有资源不受影响
  - 验证 state refresh 正确读取新参数
  - 注意：验证任务在收尾阶段执行，此处略过
