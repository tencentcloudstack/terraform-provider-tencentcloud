## Context

`tencentcloud_mqtt_instance` 是管理 MQTT 实例的 Terraform 通用资源（RESOURCE_KIND_GENERAL），支持 CRUD 操作。当前资源的 Schema 已包含 DescribeInstance 返回的多个属性，但缺少新增的 `BlockRuleLimit`（封禁规则最大数量）字段。

该字段仅在 DescribeInstance 响应中返回，不在 CreateInstance 或 ModifyInstance 请求参数中，因此是一个只读的计算属性。

SDK 中 `DescribeInstanceResponseParams` 新增字段定义：
```go
BlockRuleLimit *int64 `json:"BlockRuleLimit,omitnil,omitempty" name:"BlockRuleLimit"`
```

## Goals / Non-Goals

**Goals:**
- 在 `tencentcloud_mqtt_instance` 资源中新增 `block_rule_limit` 计算属性
- 在 Read 方法中从 DescribeInstance 响应读取并设置该属性值
- 保持向后兼容，不破坏现有 TF 配置和 state

**Non-Goals:**
- 不将 `block_rule_limit` 设置为可写参数（Create/Update 不支持该字段）
- 不修改其他 MQTT 相关资源

## Decisions

1. **Schema 定义**: 新增 `block_rule_limit` 为 `TypeInt`、`Computed: true` 的计算属性。该字段只在 DescribeInstance 响应中返回，不在 Create/Modify 请求中，因此只能作为 Computed 字段。

2. **Read 方法**: 在 `ResourceTencentCloudMqttInstanceRead` 中，当 `respData.BlockRuleLimit` 不为 nil 时，调用 `d.Set("block_rule_limit", respData.BlockRuleLimit)` 设置值。遵循现有代码模式，先判断 nil 再 set。

3. **单元测试**: 在现有的 `resource_tc_mqtt_instance_test.go` 中补充 `block_rule_limit` 的 mock 测试用例。

## Risks / Trade-offs

- [风险] 新增 Computed 字段会在下次 terraform plan/apply 时出现在 state 中 → 无实际风险，Computed 字段不影响现有配置
