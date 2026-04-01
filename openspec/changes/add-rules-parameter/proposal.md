## Why

tencentcloud_teo_l7_acc_rule 资源当前缺少 Rules 参数支持，导致无法通过 Terraform 完整管理七层加速规则。根据 DescribeL7AccRules API 的响应结构，Rules 是一个复杂的嵌套结构，包含规则状态、规则ID、规则名称、描述以及分支（Branches）等信息。接入此参数可以让用户能够完整地管理和配置七层加速规则的各项属性。

## What Changes

- 在 tencentcloud_teo_l7_acc_rule 资源的 Schema 中新增 Rules 参数（可选，列表类型）
- 更新 Read 函数，从 DescribeL7AccRules API 响应中读取并填充 Rules 数据
- 更新 Create/Update/Delete 函数，处理 Rules 字段的转换和提交逻辑
- 确保 Rules 的所有嵌套字段（Status、RuleId、RuleName、Description、Branches 等）与 CAPI 接口定义一致

## Capabilities

### New Capabilities

- `teo-l7-acc-rules`: 支持七层加速规则（L7 Acceleration Rules）的完整管理，包括规则的状态、ID、名称、描述以及规则分支等复杂嵌套结构

### Modified Capabilities

(无 spec 级别的行为变更)

## Impact

- **受影响的代码**: tencentcloud/services/teo/resource_teo_l7_acc_rule.go
- **受影响的函数**: Create, Read, Update, Delete CRUD 操作
- **测试影响**: 需要更新 resource_teo_l7_acc_rule_test.go 单元测试
- **API 调用**: 需要调用 DescribeL7AccRules API 获取 Rules 数据，以及相应的 Create/Update API
- **兼容性**: 新增 Optional 字段，不破坏现有配置和 state，保持向后兼容
