## Why

在 Terraform Provider 中接入 tencentcloud_teo_rule_engine 资源的 RuleItems 参数，以满足用户在 TEO (EdgeOne) 规则引擎中配置复杂规则项的需求。该参数允许用户定义更细粒度的规则条件、动作和优先级，增强规则引擎的配置灵活性。

## What Changes

- **新增**: 在 `tencentcloud_teo_rule_engine` 资源中添加 `rule_items` 参数
- **新增**: 支持从 DescribeRules API 读取 RuleItems 数据结构
- **新增**: 实现规则项的 CRUD 操作（创建、读取、更新、删除）
- **新增**: 添加相应的数据源字段和测试用例

## Capabilities

### New Capabilities
- `teo-rule-items`: 支持在 TEO 规则引擎资源中配置和管理 RuleItems 参数，包括规则项的条件、动作、优先级等属性

### Modified Capabilities
- (无现有能力变更)

## Impact

- **代码**: 需要修改 `resource_tc_teo_rule_engine.go` 和相关的服务层代码
- **测试**: 需要添加新的测试用例覆盖 RuleItems 参数
- **文档**: 需要更新 `resource_tc_teo_rule_engine.md` 文档
- **API**: 需要调用 DescribeRules API 获取 RuleItems 数据
