## 1. API 响应结构分析

- [x] 1.1 查阅 vendor 目录下 DescribeRabbitMQVipInstance API 响应结构，确认新参数（remark、enable_deletion_protection、enable_risk_warning）在响应中的字段名称和位置
- [x] 1.2 确认 API 响应中这些参数的返回类型和默认值（nil、空字符串、false 等）
- [x] 1.3 记录 API 响应结构信息，为后续实现做准备

## 2. Schema 定义更新

- [x] 2.1 在 `resourceTencentCloudTdmqRabbitmqVipInstance()` 函数的 Schema 中新增 `remark` 字段（Type: String, Optional: true, Computed: true, Description: 实例备注信息）
- [x] 2.2 在 Schema 中新增 `enable_deletion_protection` 字段（Type: Bool, Optional: true, Computed: true, Description: 是否开启删除保护）
- [x] 2.3 在 Schema 中新增 `enable_risk_warning` 字段（Type: Bool, Optional: true, Computed: true, Description: 是否开启集群风险提示）

## 3. Create 函数更新

- [x] 3.1 在 `resourceTencentCloudTdmqRabbitmqVipInstanceCreate()` 函数中，为 `remark` 参数添加到 CreateRabbitMQVipInstance API 请求的逻辑（使用 `d.GetOkExists` 检查）
- [x] 3.2 在 Create 函数中，为 `enable_deletion_protection` 参数添加到 API 请求的逻辑（使用 `d.GetOkExists` 检查）
- [x] 3.3 在 Create 函数中，为 `enable_risk_warning` 参数添加到 API 请求的逻辑（使用 `d.GetOkExists` 检查）

## 4. Update 函数实现

- [x] 4.1 在 `resourceTencentCloudTdmqRabbitmqVipInstanceUpdate()` 函数中，添加对 `remark` 参数的更新逻辑（使用 `d.HasChange` 检查，有变化时设置 `request.Remark` 并标记 `needUpdate = true`）
- [x] 4.2 在 Update 函数中，添加对 `enable_deletion_protection` 参数的更新逻辑（使用 `d.HasChange` 检查，有变化时设置 `request.EnableDeletionProtection` 并标记 `needUpdate = true`）
- [x] 4.3 在 Update 函数中，添加对 `enable_risk_warning` 参数的更新逻辑（使用 `d.HasChange` 检查，有变化时设置 `request.EnableRiskWarning` 并标记 `needUpdate = true`）
- [x] 4.4 确保只有当 `needUpdate` 为 true 时才调用 ModifyRabbitMQVipInstance API，避免不必要的 API 调用

## 5. Read 函数实现

- [x] 5.1 在 `resourceTencentCloudTdmqRabbitmqVipInstanceRead()` 函数中，添加从 API 响应中读取 `remark` 参数的逻辑，并使用 `d.Set()` 设置到 state
- [x] 5.2 在 Read 函数中，添加从 API 响应中读取 `enable_deletion_protection` 参数的逻辑，并使用 `d.Set()` 设置到 state
- [x] 5.3 在 Read 函数中，添加从 API 响应中读取 `enable_risk_warning` 参数的逻辑，并使用 `d.Set()` 设置到 state
- [x] 5.4 确保正确处理 API 响应中可能为 nil 的情况，避免 panic

## 6. 测试代码编写

- [x] 6.1 在 `resource_tc_tdmq_rabbitmq_vip_instance_test.go` 中，为 `remark` 参数添加更新测试用例（使用 mock API 模拟成功更新场景）
- [x] 6.2 在测试文件中，为 `enable_deletion_protection` 参数添加更新测试用例
- [x] 6.3 在测试文件中，为 `enable_risk_warning` 参数添加更新测试用例
- [x] 6.4 添加同时更新多个新参数的测试用例，验证所有参数都能正确更新
- [x] 6.5 添加读取包含新参数的实例的测试用例，验证 Read 函数能正确读取这些参数
- [x] 6.6 添加读取不包含新参数的实例的测试用例，验证能正确处理 nil 或默认值

## 7. 示例文档更新

- [x] 7.1 更新 `resource_tc_tdmq_rabbitmq_vip_instance.md` 示例文件，添加新参数的使用示例
- [x] 7.2 在示例中展示如何设置 `remark` 参数
- [x] 7.3 在示例中展示如何设置 `enable_deletion_protection` 参数
- [x] 7.4 在示例中展示如何设置 `enable_risk_warning` 参数

## 8. 验证与测试

- [x] 8.1 运行 `gofmt` 对修改的代码文件进行格式化检查
- [x] 8.2 运行单元测试 `go test -v ./tencentcloud/services/trabbit -run TestResourceTencentCloudTdmqRabbitmqVipInstance`，确保所有测试用例通过
- [x] 8.3 验证现有功能不受影响，运行所有 trabbit 服务的测试 `go test ./tencentcloud/services/trabbit/...`
- [x] 8.4 使用 `make doc` 命令生成 website/docs/ 下的 markdown 文档，验证文档生成正确
- [x] 8.5 手动验证 Terraform plan 和 apply 操作对新参数的支持（可选，在开发环境中测试）
