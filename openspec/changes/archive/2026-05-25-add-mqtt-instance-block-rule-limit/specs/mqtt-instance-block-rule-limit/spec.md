## ADDED Requirements

### Requirement: Block Rule Limit Computed Attribute

资源 `tencentcloud_mqtt_instance` 必须在 Schema 中新增 `block_rule_limit` 计算属性（Computed, TypeInt），映射到 DescribeInstance 响应中的 `BlockRuleLimit` 字段，表示实例的封禁规则最大数量。

#### Scenario: Read block_rule_limit from DescribeInstance response
- **WHEN** 执行资源的 Read 操作
- **THEN** 调用 DescribeInstance 接口获取实例详情
- **AND** 如果响应中的 `BlockRuleLimit` 不为 nil，则将其值设置到 Terraform state 的 `block_rule_limit` 属性中
- **AND** 如果响应中的 `BlockRuleLimit` 为 nil，则不设置该属性

#### Scenario: block_rule_limit is read-only
- **WHEN** 用户尝试在 Terraform 配置中显式设置 `block_rule_limit` 的值
- **THEN** Terraform 返回错误，因为该属性为 Computed（只读），不可由用户指定

#### Scenario: block_rule_limit not in Create or Update
- **WHEN** 创建或更新 MQTT 实例
- **THEN** `block_rule_limit` 不作为请求参数传递给 CreateInstance 或 ModifyInstance 接口
