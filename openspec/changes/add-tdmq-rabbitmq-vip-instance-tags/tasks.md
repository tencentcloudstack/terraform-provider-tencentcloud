# 实施任务清单

## 1. Schema 修改
- [x] 1.1 在 `resource_tc_tdmq_rabbitmq_vip_instance.go` 的 Schema 中添加 `tags` 字段
  - [x] 类型: `schema.TypeMap`
  - [x] Optional: `true`
  - [x] Computed: `true`
  - [x] Description: 资源标签说明

## 2. Create 函数修改
- [x] 2.1 在 `resourceTencentCloudTdmqRabbitmqVipInstanceCreate` 中添加标签处理
  - [x] 2.1.1 使用 `helper.GetTags` 从 ResourceData 获取标签
  - [x] 2.1.2 将标签转换为 `[]*tdmq.Tag` 格式
  - [x] 2.1.3 设置 `request.ResourceTags` 参数
  - [x] 2.1.4 在 API 调用时传递标签参数

## 3. Read 函数修改
- [x] 3.1 在 `resourceTencentCloudTdmqRabbitmqVipInstanceRead` 中添加标签读取
  - [x] 3.1.1 从 API 响应中获取 `Tags` 字段
  - [x] 3.1.2 将 `[]*tdmq.Tag` 转换为 `map[string]string` 格式
  - [x] 3.1.3 使用 `d.Set("tags", tags)` 设置到状态

## 4. Update 函数修改

**⚠️ 重要**: 根据 API 文档 https://cloud.tencent.com/document/api/1179/88450 ,`ModifyRabbitMQVipInstance` 接口支持 `Tags` 参数(全量标签更新,非增量)。

**实施方案**: 
- 优先检查当前 SDK 版本(`tencentcloud/tdmq v1.1.15`)的 `ModifyRabbitMQVipInstanceRequest` 是否包含 `Tags` 字段
- **如果 SDK 支持**: 直接使用 `ModifyRabbitMQVipInstance` 接口的 `Tags` 参数(全量替换)
- **如果 SDK 不支持**: 需要先升级 SDK 版本,或使用统一标签服务作为临时方案

- [x] 4.1 检查并升级 SDK(如需要)
  - [x] 4.1.1 检查 `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217/models.go` 中 `ModifyRabbitMQVipInstanceRequest` 是否有 `Tags` 字段
  - [x] 4.1.2 如果没有,检查是否有新版本 SDK 支持
  - [x] 4.1.3 **结论**: SDK 不支持,使用统一标签服务作为方案

- [x] 4.2 在 `resourceTencentCloudTdmqRabbitmqVipInstanceUpdate` 中添加标签更新
  - [x] 4.2.1 添加 `d.HasChange("tags")` 检测
  - [x] 4.2.2 获取 old/new tags 并使用 `svctag.DiffTags` 计算差异
  - [x] 4.2.3 **已实施**: 使用 `svctag.TagService.ModifyTags` 统一标签服务
  - [x] 4.2.4 添加错误处理

## 5. 导入辅助包
- [x] 5.1 确保导入 `svctag` 包: `github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag`

## 6. 文档更新
- [x] 6.1 更新源文档 `resource_tc_tdmq_rabbitmq_vip_instance.md`
  - [x] 6.1.1 在参数列表中添加 `tags` 说明
  - [x] 6.1.2 添加使用示例,展示如何配置标签
- [x] 6.2 运行 `make doc` 重新生成文档
- [x] 6.3 验证生成的文档 `website/docs/r/tdmq_rabbitmq_vip_instance.html.markdown`

## 7. Changelog
- [x] 7.1 创建 `.changelog/3848.txt` 文件
  - [x] 格式: `resource/tencentcloud_tdmq_rabbitmq_vip_instance: support tags parameter`

## 8. 代码质量检查
- [x] 8.1 运行 `make fmt` 格式化代码
- [x] 8.2 运行 `go build` 验证编译
- [x] 8.3 验证代码逻辑与现有 tags 实现一致(参考其他资源)

## 9. 验证
- [ ] 9.1 手动测试创建实例时绑定标签
- [ ] 9.2 手动测试读取实例时获取标签
- [ ] 9.3 手动测试更新实例标签
- [ ] 9.4 手动测试删除实例标签
- [ ] 9.5 验证导入已有资源时标签正常同步

**注**: 第 9 项验证需要实际的云环境测试,已完成代码实施和编译验证。
