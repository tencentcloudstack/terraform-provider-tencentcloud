## Why

DescribeL7AccRules API 返回的 TotalCount 参数未在 tencentcloud_teo_l7_acc_rule 数据源中暴露，导致用户无法获取七层加速规则的总量信息，限制了在配置管理和自动化场景下的统计与监控能力。

## What Changes

- 在 tencentcloud_teo_l7_acc_rule 数据源中新增 TotalCount 输出字段
- 该字段从 DescribeL7AccRules API 的响应中读取并映射到数据源输出
- 保持向后兼容性，不影响现有字段

## Capabilities

### New Capabilities
- `teo-l7-acc-rule-totalcount`: 在 tencentcloud_teo_l7_acc_rule 数据源中暴露 TotalCount 字段，返回七层加速规则的总数

### Modified Capabilities
（无）

## Impact

- 修改文件：`tencentcloud/services/teo/data_source_tc_teo_l7_acc_rule.go`
- 测试文件：`tencentcloud/services/teo/data_source_tc_teo_l7_acc_rule_test.go`
- 文档文件：`website/docs/r/teo_l7_acc_rule.html.markdown`
- 不涉及依赖变更
- 不影响现有用户配置
