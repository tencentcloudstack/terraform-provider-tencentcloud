## Why

tencentcloud_teo_rule_engine 资源当前缺少对 RuleItems 参数的支持，该参数在 DescribeRules API 中返回，包含了规则引擎的详细规则项配置信息。缺少该参数导致用户无法通过 Terraform 获取完整的规则引擎配置，影响了对规则的查看和验证需求。

## What Changes

- 为 tencentcloud_teo_rule_engine 资源添加 RuleItems 参数的读取支持
- 从 DescribeRules API 的响应中解析并映射 RuleItems 数据到 Terraform schema
- 新增 RuleItems 参数的相关文档和示例

## Capabilities

### New Capabilities
- `teo-rule-items-read`: 支持从 DescribeRules API 读取并返回规则引擎的 RuleItems 配置信息

### Modified Capabilities
- 无

## Impact

- 受影响的资源：tencentcloud_teo_rule_engine
- 涉及文件：
  - `tencentcloud/services/teo/resource_tencentcloud_teo_rule_engine.go`（如果是资源文件）
  - `tencentcloud/services/teo/data_source_tencentcloud_teo_rule_engine.go`（如果是数据源文件）
- API 调用：DescribeRules API 的响应解析逻辑需要更新
- 兼容性：非破坏性变更，仅新增 Optional 参数
