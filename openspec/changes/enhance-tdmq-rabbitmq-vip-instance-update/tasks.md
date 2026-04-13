## 1. Schema 修改

- [x] 1.1 在 `resource_tc_tdmq_rabbitmq_vip_instance.go` 的 Schema 中新增 `remark` 字段（TypeString, Optional）
- [x] 1.2 在 Schema 中新增 `enable_deletion_protection` 字段（TypeBool, Optional）
- [x] 1.3 在 Schema 中新增 `enable_risk_warning` 字段（TypeBool, Optional）

## 2. Create 方法修改

- [x] 2.1 在 `resourceTencentCloudTdmqRabbitmqVipInstanceCreate` 函数中添加对 `remark` 字段的处理
- [x] 2.2 在 Create 函数中添加对 `enable_deletion_protection` 字段的处理
- [x] 2.3 在 Create 函数中添加对 `enable_risk_warning` 字段的处理（如果 Create API 支持）

## 3. Read 方法修改

- [x] 3.1 在 `resourceTencentCloudTdmqRabbitmqVipInstanceRead` 函数中添加从 API 响应读取 `remark` 字段的逻辑
- [x] 3.2 在 Read 函数中添加从 API 响应读取 `enable_deletion_protection` 字段的逻辑
- [x] 3.3 在 Read 函数中添加从 API 响应读取 `enable_risk_warning` 字段的逻辑（如果 API 响应包含）
- [x] 3.4 添加 nil 值检查，避免因 API 响应中字段为 nil 导致 panic

## 4. Update 方法修改

- [x] 4.1 从不可变字段列表中移除 `enable_deletion_protection`（如果存在）
- [x] 4.2 在 `resourceTencentCloudTdmqRabbitmqVipInstanceUpdate` 函数中添加对 `remark` 字段更新的检测和处理
- [x] 4.3 在 Update 函数中添加对 `enable_deletion_protection` 字段更新的检测和处理
- [x] 4.4 在 Update 函数中添加对 `enable_risk_warning` 字段更新的检测和处理

## 5. 测试实现

- [x] 5.1 在 `resource_tc_tdmq_rabbitmq_vip_instance_test.go` 中添加测试用例：创建实例时设置 remark 字段
- [x] 5.2 添加测试用例：创建实例时设置 enable_deletion_protection 字段
- [x] 5.3 添加测试用例：创建实例时设置 enable_risk_warning 字段
- [x] 5.4 添加测试用例：更新实例的 remark 字段
- [x] 5.5 添加测试用例：更新实例的 enable_deletion_protection 字段
- [x] 5.6 添加测试用例：更新实例的 enable_risk_warning 字段
- [x] 5.7 使用 mock 云 API 的方式实现测试用例，不调用真实的云 API

## 6. 文档更新

- [x] 6.1 在 `resource_tc_tdmq_rabbitmq_vip_instance.md` 中添加 `remark` 字段的文档说明和使用示例
- [x] 6.2 在文档中添加 `enable_deletion_protection` 字段的文档说明和使用示例
- [x] 6.3 在文档中添加 `enable_risk_warning` 字段的文档说明和使用示例
- [x] 6.4 更新文档中的不可变字段列表，移除 enable_deletion_protection
- [x] 6.5 运行 `make doc` 命令生成 `website/docs/r/tdmq_rabbitmq_vip_instance.html.markdown`

## 7. 代码验证

- [x] 7.1 运行 `go fmt ./tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance.go` 格式化代码
- [x] 7.2 运行 `go fmt ./tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance_test.go` 格式化测试代码
- [x] 7.3 运行单元测试验证修改的正确性 (跳过 - 根据禁止事项，不执行测试)

## 8. 提交和 PR

- [x] 8.1 提交代码变更，commit message 清晰描述本次修改的内容
- [x] 8.2 推送到远程仓库
- [x] 8.3 创建 Pull Request，关联本次变更提案
