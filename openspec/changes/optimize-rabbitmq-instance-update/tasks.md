# 实现任务清单

## 1. 代码实现
- [ ] 1.1 修改 `resourceTencentCloudTdmqRabbitmqVipInstanceUpdate` 函数
  - [ ] 1.1.1 更新不可变字段列表,仅保留真正不可变的字段
  - [ ] 1.1.2 添加 `auto_renew_flag` 字段更新逻辑
  - [ ] 1.1.3 添加 `enable_public_access` 字段更新逻辑
  - [ ] 1.1.4 添加 `band_width` 字段更新逻辑
  - [ ] 1.1.5 优化 `resource_tags` 字段更新逻辑(已有)
  - [ ] 1.1.6 优化 `cluster_name` 字段更新逻辑(已有)
- [ ] 1.2 检查并移除 Schema 中可变字段的 `ForceNew` 标记
  - [ ] 1.2.1 检查 `auto_renew_flag` 是否有 `ForceNew` 标记,如有则移除
  - [ ] 1.2.2 检查 `enable_public_access` 是否有 `ForceNew` 标记,如有则移除
  - [ ] 1.2.3 检查 `band_width` 是否有 `ForceNew` 标记,如有则移除
- [ ] 1.3 优化错误信息提示
  - [ ] 1.3.1 为不可变字段返回更友好的错误提示
  - [ ] 1.3.2 为暂时不支持的字段提供重建建议

## 2. API 调用验证
- [ ] 2.1 验证 `ModifyRabbitMQVipInstance` API 支持的字段
  - [ ] 2.1.1 确认 API 支持 `AutoRenewFlag` 参数
  - [ ] 2.1.2 确认 API 支持 `EnablePublicAccess` 参数
  - [ ] 2.1.3 确认 API 支持 `Bandwidth` 参数
  - [ ] 2.1.4 确认 API 支持 `ClusterName` 参数
  - [ ] 2.1.5 确认 API 支持 `Tags` 参数
- [ ] 2.2 查阅腾讯云 TDMQ RabbitMQ API 文档
  - [ ] 2.2.1 确认每个字段的类型和限制
  - [ ] 2.2.2 确认字段更新后是否需要等待实例状态稳定

## 3. 测试实现
- [ ] 3.1 创建/更新 `resource_tc_tdmq_rabbitmq_vip_instance_test.go`
  - [ ] 3.1.1 实现自动续费标志更新测试 `TestAccTencentCloudTdmqRabbitmqVipInstance_updateAutoRenewFlag`
  - [ ] 3.1.2 实现公网访问开关更新测试 `TestAccTencentCloudTdmqRabbitmqVipInstance_updatePublicAccess`
  - [ ] 3.1.3 实现带宽更新测试 `TestAccTencentCloudTdmqRabbitmqVipInstance_updateBandWidth`
  - [ ] 3.1.4 实现标签更新测试 `TestAccTencentCloudTdmqRabbitmqVipInstance_updateTags`
  - [ ] 3.1.5 实现不可变字段更新错误测试 `TestAccTencentCloudTdmqRabbitmqVipInstance_updateImmutableFields`
- [ ] 3.2 编写测试配置模板
  - [ ] 3.2.1 创建包含自动续费标志的测试配置
  - [ ] 3.2.2 创建包含公网访问配置的测试配置
  - [ ] 3.2.3 创建包含带宽配置的测试配置
  - [ ] 3.2.4 创建包含标签的测试配置
- [ ] 3.3 运行验收测试
  - [ ] 3.3.1 运行自动续费标志更新测试并验证通过
  - [ ] 3.3.2 运行公网访问开关更新测试并验证通过
  - [ ] 3.3.3 运行带宽更新测试并验证通过
  - [ ] 3.3.4 运行标签更新测试并验证通过
  - [ ] 3.3.5 运行不可变字段更新错误测试并验证错误提示正确

## 4. 文档更新
- [ ] 4.1 更新资源文档 `resource_tc_tdmq_rabbitmq_vip_instance.md`
  - [ ] 4.1.1 更新 `auto_renew_flag` 字段说明,注明可更新
  - [ ] 4.1.2 更新 `enable_public_access` 字段说明,注明可更新
  - [ ] 4.1.3 更新 `band_width` 字段说明,注明可更新
  - [ ] 4.1.4 更新不可变字段列表,提供清晰的说明
  - [ ] 4.1.5 添加字段更新限制的说明
- [ ] 4.2 更新网站文档 `website/docs/r/tdmq_rabbitmq_vip_instance.html.markdown`
  - [ ] 4.2.1 添加字段更新说明
  - [ ] 4.2.2 添加更新示例代码
  - [ ] 4.2.3 添加不可变字段说明和重建建议
- [ ] 4.3 生成文档
  - [ ] 4.3.1 运行 `make doc` 生成文档
  - [ ] 4.3.2 验证文档生成无误

## 5. 代码质量检查
- [ ] 5.1 运行代码格式化
  - [ ] 5.1.1 运行 `make fmt` 格式化代码
  - [ ] 5.1.2 检查格式化结果
- [ ] 5.2 运行代码检查
  - [ ] 5.2.1 运行 `make lint` 确保无 lint 错误
  - [ ] 5.2.2 修复所有 lint 错误
- [ ] 5.3 代码审查
  - [ ] 5.3.1 检查错误处理逻辑
  - [ ] 5.3.2 检查日志记录
  - [ ] 5.3.3 检查注释完整性
  - [ ] 5.3.4 检查所有字段都有正确的 Description

## 6. 集成测试
- [ ] 6.1 准备测试环境
  - [ ] 6.1.1 准备测试用的 RabbitMQ 实例
  - [ ] 6.1.2 准备测试用的 VPC 和子网
- [ ] 6.2 执行完整测试流程
  - [ ] 6.2.1 创建实例
  - [ ] 6.2.2 更新可变字段
  - [ ] 6.2.3 验证更新结果
  - [ ] 6.2.4 尝试更新不可变字段,验证错误提示
  - [ ] 6.2.5 清理测试资源

## 7. 验证和发布
- [ ] 7.1 功能验证
  - [ ] 7.1.1 验证自动续费标志更新功能正常
  - [ ] 7.1.2 验证公网访问开关更新功能正常
  - [ ] 7.1.3 验证带宽更新功能正常
  - [ ] 7.1.4 验证标签更新功能正常
  - [ ] 7.1.5 验证不可变字段限制正常
- [ ] 7.2 向后兼容性验证
  - [ ] 7.2.1 验证现有 Terraform 配置无需修改即可工作
  - [ ] 7.2.2 验证升级后不会破坏现有资源状态
- [ ] 7.3 性能验证
  - [ ] 7.3.1 验证更新操作响应时间在可接受范围内
  - [ ] 7.3.2 验证更新操作不会导致实例不稳定
- [ ] 7.4 提交代码
  - [ ] 7.4.1 提交代码到特性分支
  - [ ] 7.4.2 创建 Pull Request
  - [ ] 7.4.3 关联此 proposal

## 8. 未来优化(可选)
- [ ] 8.1 评估规格变更支持
  - [ ] 8.1.1 研究是否支持 `node_spec` 更新
  - [ ] 8.1.2 研究是否支持 `node_num` 更新
  - [ ] 8.1.3 研究是否支持 `storage_size` 更新
  - [ ] 8.1.4 如支持,添加对应的 update 逻辑
- [ ] 8.2 添加状态等待机制
  - [ ] 8.2.1 对于影响实例状态的更新,添加状态等待逻辑
  - [ ] 8.2.2 确保更新完成后再返回
- [ ] 8.3 优化更新性能
  - [ ] 8.3.1 批量更新多个字段,减少 API 调用次数
  - [ ] 8.3.2 添加更新操作的并发控制
