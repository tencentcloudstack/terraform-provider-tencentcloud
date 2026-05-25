## Why

MQTT 实例的 DescribeInstance 接口新增了 `BlockRuleLimit` 字段（封禁规则最大数量），需要在 Terraform 资源 `tencentcloud_mqtt_instance` 中暴露该只读属性，以便用户能够通过 Terraform state 获取实例的封禁规则配额信息。

## What Changes

- 在 `tencentcloud_mqtt_instance` 资源的 Schema 中新增 `block_rule_limit` 计算属性（Computed, TypeInt），映射到 DescribeInstance 响应中的 `BlockRuleLimit` 字段
- 在资源的 Read 方法中，从 DescribeInstance 响应中读取 `BlockRuleLimit` 并设置到 state

## Capabilities

### New Capabilities

- `mqtt-instance-block-rule-limit`: 为 tencentcloud_mqtt_instance 资源新增 block_rule_limit 只读计算属性，暴露 DescribeInstance 接口返回的封禁规则最大数量

### Modified Capabilities

## Impact

- `tencentcloud/services/mqtt/resource_tc_mqtt_instance.go`: 新增 schema 字段和 Read 逻辑
- `tencentcloud/services/mqtt/resource_tc_mqtt_instance.md`: 新增属性文档
