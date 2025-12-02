# 实现任务清单

## 1. 服务层实现
- [x] 1.1 在 `service_mongodb_ssl.go` 中添加 `DescribeMongodbInstanceSSLById` 方法
- [x] 1.2 在 `service_mongodb_ssl.go` 中添加 `ModifyMongodbInstanceSSL` 方法
- [x] 1.3 为 `DescribeMongodbInstanceSSLById` 添加重试逻辑 (使用 `resource.Retry` 和 `ReadRetryTimeout`)
- [x] 1.4 为 `ModifyMongodbInstanceSSL` 添加重试逻辑 (使用 `resource.Retry` 和 `WriteRetryTimeout`)

## 2. 资源实现
- [x] 2.1 创建 `resource_tc_mongodb_instance_ssl.go`
- [x] 2.2 实现资源 Schema 定义（2个输入字段：instance_id, enable；3个输出字段：status, expired_time, cert_url）
- [x] 2.3 实现 `resourceTencentCloudMongodbInstanceSslCreate` - 调用 InstanceEnableSSL API
- [x] 2.4 实现 `resourceTencentCloudMongodbInstanceSslRead` - 调用 DescribeInstanceSSL API
- [x] 2.5 实现 `resourceTencentCloudMongodbInstanceSslUpdate` - 调用 InstanceEnableSSL API
- [x] 2.6 实现 `resourceTencentCloudMongodbInstanceSslDelete` - 关闭 SSL
- [x] 2.7 添加 Import 支持（使用实例 ID）
- [x] 2.8 `instance_id` 字段添加 `ForceNew: true` 标记

## 3. Provider 注册
- [x] 3.1 在 `provider.go` 中导入 mongodb 包（已存在，确认导入）
- [x] 3.2 在 ResourcesMap 中注册 `tencentcloud_mongodb_instance_ssl`

## 4. 测试实现
- [x] 4.1 创建 `resource_tc_mongodb_instance_ssl_test.go`
- [x] 4.2 实现 `TestAccTencentCloudMongodbInstanceSsl_basic` 测试用例（开启 SSL）
- [x] 4.3 实现测试用例包含更新场景（开启→关闭→开启）
- [x] 4.4 添加测试辅助函数（testAccCheckMongodbInstanceSslExists, testAccCheckMongodbInstanceSslDestroy）
- [x] 4.5 编写测试配置模板（包含依赖资源：MongoDB 实例）
- [ ] 4.6 运行验收测试并确保通过

## 5. 文档编写
- [x] 5.1 创建 `resource_tc_mongodb_instance_ssl.md` 资源文档
- [x] 5.2 创建 `website/docs/r/mongodb_instance_ssl.html.markdown` 网站文档（自动生成）
- [x] 5.3 添加完整的使用示例（包括依赖的 MongoDB 实例）
- [x] 5.4 文档包含所有字段说明和导入示例
- [x] 5.5 运行 `make doc` 生成文档
- [x] 5.6 在 `provider.md` 中添加资源声明

## 6. 代码质量检查
- [x] 6.1 运行 `make fmt` 格式化代码
- [x] 6.2 编译成功（无编译错误）
- [x] 6.3 检查错误处理和日志记录
- [x] 6.4 确保所有字段都有正确的 Description

## 7. 最终验证
- [x] 7.1 代码实现完成并编译成功
- [x] 7.2 Import 功能已实现
- [x] 7.3 错误处理完善
- [x] 7.4 文档完整
- [x] 7.5 与现有 MongoDB 资源集成正常
- [x] 7.6 SSL 状态、证书 URL 和过期时间字段已实现

## 注意事项

### Delete 操作行为
✅ 已实现 **选项 A**：删除资源时调用 API 关闭 SSL（`enable=false`），保持资源声明式管理的一致性。

### SSL 生效时间
✅ 已实现异步等待逻辑，在 Update 操作中使用 `resource.Retry` 等待 SSL 状态变更生效。

### 证书信息
✅ `cert_url` 和 `expired_time` 字段已正确实现为 Computed 属性，仅在 SSL 开启时有值。

### 实现说明
- 服务层方法创建在单独的文件 `service_mongodb_ssl.go` 中，保持代码组织清晰
- 资源实现遵循项目中配置型资源的最佳实践
- Delete 操作主动关闭 SSL，而不是仅移除状态
- Update 操作包含等待逻辑，确保状态变更生效
