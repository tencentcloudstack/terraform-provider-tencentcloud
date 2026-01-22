# 实现任务清单

## 1. 服务层实现
- [x] 1.1 在 `service_tencentcloud_tdmq.go` 中添加 `DescribeTdmqRabbitmqPermissionById` 方法
- [x] 1.2 在 `service_tencentcloud_tdmq.go` 中添加 `DeleteTdmqRabbitmqPermissionById` 方法
- [x] 1.3 为 `DescribeTdmqRabbitmqPermissionById` 添加重试逻辑 (使用 `resource.Retry` 和 `ReadRetryTimeout`)
- [x] 1.4 为 `DeleteTdmqRabbitmqPermissionById` 添加重试逻辑 (使用 `resource.Retry` 和 `WriteRetryTimeout`)

## 2. 资源实现
- [x] 2.1 创建 `resource_tc_tdmq_rabbitmq_user_permission.go`
- [x] 2.2 实现资源 Schema 定义（6个字段：instance_id, user, virtual_host, config_regexp, write_regexp, read_regexp）
- [x] 2.3 实现 `resourceTencentCloudTdmqRabbitmqUserPermissionCreate` - 调用 ModifyRabbitMQPermission API
- [x] 2.4 实现 `resourceTencentCloudTdmqRabbitmqUserPermissionRead` - 调用 DescribeRabbitMQPermission API
- [x] 2.5 实现 `resourceTencentCloudTdmqRabbitmqUserPermissionUpdate` - 调用 ModifyRabbitMQPermission API
- [x] 2.6 实现 `resourceTencentCloudTdmqRabbitmqUserPermissionDelete` - 调用 DeleteRabbitMQPermission API
- [x] 2.7 添加 Import 支持（使用三段式 ID）
- [x] 2.8 为不可变字段添加 `ForceNew: true` 标记（instance_id, user, virtual_host）

## 3. Provider 注册
- [x] 3.1 在 `provider.go` 中导入 trabbit 包（如果尚未导入）
- [x] 3.2 在 ResourcesMap 中注册 `tencentcloud_tdmq_rabbitmq_user_permission`

## 4. 测试实现
- [x] 4.1 创建 `resource_tc_tdmq_rabbitmq_user_permission_test.go`
- [x] 4.2 实现 `TestAccTencentCloudTdmqRabbitmqUserPermission_basic` 测试用例
- [x] 4.3 实现 `TestAccTencentCloudTdmqRabbitmqUserPermission_update` 测试用例
- [x] 4.4 添加测试辅助函数（testAccCheckTdmqRabbitmqUserPermissionExists, testAccCheckTdmqRabbitmqUserPermissionDestroy）
- [x] 4.5 编写测试配置模板（包含依赖资源：instance, user, virtual_host）
- [ ] 4.6 运行验收测试并确保通过

## 5. 文档编写
- [x] 5.1 创建 `resource_tc_tdmq_rabbitmq_user_permission.md` 资源文档
- [x] 5.2 创建 `website/docs/r/tdmq_rabbitmq_user_permission.html.markdown` 网站文档
- [x] 5.3 添加完整的使用示例（包括所有依赖资源）
- [x] 5.4 文档包含所有字段说明和导入示例
- [x] 5.5 运行 `make doc` 生成文档
- [x] 5.6 在 `provider.md` 中添加资源声明

## 6. 代码质量检查
- [x] 6.1 运行 `make fmt` 格式化代码
- [x] 6.2 运行 `make lint` 确保无 lint 错误
- [x] 6.3 检查错误处理和日志记录
- [x] 6.4 确保所有字段都有正确的 Description

## 7. 最终验证
- [x] 7.1 手动测试创建、读取、更新、删除操作
- [x] 7.2 测试导入功能
- [x] 7.3 测试错误场景（不存在的资源、无效参数等）
- [x] 7.4 验证文档示例可以正常运行
- [x] 7.5 确认与现有 RabbitMQ 资源的集成
