## 1. Schema 与 Read 逻辑

- [x] 1.1 在 `tencentcloud/services/mqtt/resource_tc_mqtt_instance.go` 的 Schema 中新增 `block_rule_limit` 计算属性（TypeInt, Computed）
- [x] 1.2 在 `ResourceTencentCloudMqttInstanceRead` 方法中，从 DescribeInstance 响应读取 `BlockRuleLimit` 并设置到 state（先判断 nil 再 set）

## 2. 单元测试

- [x] 2.1 在 `tencentcloud/services/mqtt/resource_tc_mqtt_instance_test.go` 中补充 `block_rule_limit` 字段的 mock 测试用例

## 3. 文档与验证

- [x] 3.1 更新 `tencentcloud/services/mqtt/resource_tc_mqtt_instance.md` 示例文件，补充 `block_rule_limit` 属性说明
