# Change: Add MQTT Message Enrichment Rule Resource

## Why

腾讯云 MQTT 服务提供了消息属性增强规则功能,允许用户在消息发布时动态添加或修改消息属性,如设置消息过期时间、响应主题、关联数据和用户自定义属性等。这对于实现设备追踪、消息路由和消息生命周期管理等场景非常重要。

目前 Terraform Provider 尚未支持该资源类型,用户无法通过 IaC 方式管理 MQTT 消息属性增强规则,需要手动在控制台或通过 API 进行配置,不利于自动化运维和版本管理。

## What Changes

新增 Terraform 资源 `tencentcloud_mqtt_message_enrichment_rule`,支持:

- **创建**消息属性增强规则 (CreateMessageEnrichmentRule API)
- **查询**消息属性增强规则 (DescribeMessageEnrichmentRules API)
- **更新**消息属性增强规则 (ModifyMessageEnrichmentRule API)
- **删除**消息属性增强规则 (DeleteMessageEnrichmentRule API)
- **导入**现有规则到 Terraform 状态

关键功能特性:
- 资源唯一标识符组合: `InstanceId` + `Id` (规则ID)
- `Priority` 字段作为只读 computed 字段,创建时使用默认值 1
- `Condition` 和 `Actions` 字段需要 Base64 编码的 JSON 字符串
- 支持规则状态管理 (激活/不激活)

## Impact

**新增文件**:
- `tencentcloud/services/mqtt/resource_tc_mqtt_message_enrichment_rule.go` - 资源实现
- `tencentcloud/services/mqtt/resource_tc_mqtt_message_enrichment_rule_test.go` - 验收测试
- `tencentcloud/services/mqtt/resource_tc_mqtt_message_enrichment_rule.md` - 资源文档

**修改文件**:
- `tencentcloud/services/mqtt/service_tencentcloud_mqtt.go` - 添加服务层方法
- Provider 注册文件 - 注册新资源

**受影响的规范 (Specs)**:
- 新增 `mqtt-message-enrichment-rule` 能力规范

**依赖关系**:
- 依赖现有的 MQTT 服务框架和客户端初始化
- 使用 `tencentcloud-sdk-go/tencentcloud/mqtt/v20240516` SDK 包
- 无破坏性变更,完全向后兼容

**测试覆盖**:
- 基本 CRUD 操作验收测试
- 导入功能测试
- 字段验证测试

**文档更新**:
- 生成 Terraform Registry 格式的资源文档
- 包含完整的使用示例和参数说明
