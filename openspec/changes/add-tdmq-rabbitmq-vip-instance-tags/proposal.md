# Change: 为 tencentcloud_tdmq_rabbitmq_vip_instance 资源添加 tags 支持

## Why

当前 `tencentcloud_tdmq_rabbitmq_vip_instance` 资源不支持标签(tags)参数,而腾讯云 TDMQ RabbitMQ VIP 实例的 API 已经支持标签功能:
- `CreateRabbitMQVipInstance` API 支持 `ResourceTags` 参数用于创建时绑定标签
- `DescribeRabbitMQVipInstance` API 返回的实例信息中包含 `Tags` 字段
- 标签是腾讯云资源管理的重要功能,用于资源分类、成本分摊和权限管理

用户无法通过 Terraform 为 RabbitMQ VIP 实例配置标签,需要手动在控制台操作,影响了 IaC 的完整性和自动化程度。

## What Changes

- 在 `tencentcloud_tdmq_rabbitmq_vip_instance` 资源的 Schema 中添加 `tags` 字段
- 在 Create 操作中支持通过 `CreateRabbitMQVipInstance.ResourceTags` 参数创建时绑定标签
- 在 Read 操作中读取实例的 `Tags` 字段并同步到 Terraform 状态
- 在 Update 操作中支持标签的修改和删除,使用 `ModifyRabbitMQVipInstance.Tags` 参数(全量替换)
  - 注意: API 文档明确支持此参数,但需要检查当前 SDK 版本是否包含该字段
- 更新资源文档,添加 `tags` 参数说明和使用示例

## Impact

- **影响的资源**: `tencentcloud_tdmq_rabbitmq_vip_instance`
- **影响的文件**:
  - `tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance.go` - 资源实现
  - `tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance.md` - 源文档
  - `website/docs/r/tdmq_rabbitmq_vip_instance.html.markdown` - 生成的文档
  - `go.mod` / `go.sum` - 可能需要升级 SDK 版本
- **向后兼容性**: ✅ 完全兼容
  - 新增可选参数,不影响现有配置
  - 现有实例不受影响
- **用户影响**: 正面
  - 提供完整的标签管理能力
  - 统一资源管理体验
  - 提升自动化程度

### SDK 版本说明

- **当前版本**: `tencentcloud/tdmq v1.1.15`
- **API 支持情况**:
  - `CreateRabbitMQVipInstance.ResourceTags` ✅ 已验证(SDK 包含)
  - `ModifyRabbitMQVipInstance.Tags` ✅ API 文档支持,但当前 SDK 版本**不包含**此字段
  - `DescribeRabbitMQVipInstance` 返回 `Tags` ✅ 已验证(SDK 包含)
- **实施要求**: 需要在实施时验证最新 SDK 版本,如果 `ModifyRabbitMQVipInstanceRequest` 仍不包含 `Tags` 字段,可能需要升级 SDK 或使用统一标签服务作为临时方案
