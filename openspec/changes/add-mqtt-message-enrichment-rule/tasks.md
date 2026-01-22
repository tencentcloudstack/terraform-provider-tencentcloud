## 1. 实现资源核心功能
- [x] 1.1 创建 `resource_tc_mqtt_message_enrichment_rule.go` 文件
- [x] 1.2 实现 Schema 定义，包含所有必要字段
- [x] 1.3 实现 `resourceTencentCloudMqttMessageEnrichmentRuleCreate` 函数
- [x] 1.4 实现 `resourceTencentCloudMqttMessageEnrichmentRuleRead` 函数
- [x] 1.5 实现 `resourceTencentCloudMqttMessageEnrichmentRuleUpdate` 函数
- [x] 1.6 实现 `resourceTencentCloudMqttMessageEnrichmentRuleDelete` 函数

## 2. 服务层集成
- [x] 2.1 在 MqttService 中添加 `CreateMqttMessageEnrichmentRule` 方法
- [x] 2.2 在 MqttService 中添加 `DescribeMqttMessageEnrichmentRuleById` 方法
- [x] 2.3 在 MqttService 中添加 `ModifyMqttMessageEnrichmentRule` 方法
- [x] 2.4 在 MqttService 中添加 `DeleteMqttMessageEnrichmentRuleById` 方法
- [x] 2.5 调用 MQTT SDK 的相应方法并处理响应

## 3. Provider 集成
- [x] 3.1 在 MQTT 服务的 provider 注册中添加新资源
- [x] 3.2 确保资源名称为 `tencentcloud_mqtt_message_enrichment_rule`

## 4. 编写测试
- [x] 4.1 创建 `resource_tc_mqtt_message_enrichment_rule_test.go` 文件
- [x] 4.2 实现基本功能验收测试 `TestAccTencentCloudMqttMessageEnrichmentRuleResource_basic`
- [x] 4.3 实现更新功能测试 `TestAccTencentCloudMqttMessageEnrichmentRuleResource_update`
- [x] 4.4 测试用例应验证 CRUD 操作
- [x] 4.5 确保测试使用真实的 MQTT 实例

## 5. 编写文档
- [x] 5.1 创建 `resource_tc_mqtt_message_enrichment_rule.md` 文档文件
- [x] 5.2 包含参数说明、属性说明和使用示例
- [x] 5.3 运行 `make doc` 生成最终文档

## 6. 代码质量保证
- [x] 6.1 运行 `make fmt` 格式化代码
- [x] 6.2 运行 `make lint` 检查代码质量
- [x] 6.3 运行 `make test` 执行单元测试
- [x] 6.4 运行 `TF_ACC=1 make testacc TEST=./tencentcloud/services/mqtt` 执行验收测试

## 7. 验证和文档
- [x] 7.1 手动测试资源功能
- [x] 7.2 验证 CRUD 操作符合预期
- [x] 7.3 确保错误处理正确（如规则不存在时）
- [x] 7.4 检查生成的文档格式正确