## Why

tencentcloud_teo_l7_acc_rule 数据源需要支持 Offset 参数以实现分页查询功能。DescribeL7AccRules API 已经支持 Offset 参数，但当前数据源未暴露此参数，导致用户无法进行灵活的分页数据获取。

## What Changes

- 在 tencentcloud_teo_l7_acc_rule 数据源的 schema 中新增 Offset 参数（Optional, Int, 用于分页偏移量）
- 在数据源查询逻辑中将 Offset 参数传递给 DescribeL7AccRules API
- 更新数据源文档说明 Offset 参数的使用方法

## Capabilities

### New Capabilities

- `teo-l7-acc-rule-datasource-pagination`: 为 tencentcloud_teo_l7_acc_rule 数据源添加分页查询能力，支持 Offset 参数以控制数据查询偏移量

### Modified Capabilities

- 无

## Impact

- 受影响文件：
  - `tencentcloud/services/teo/data_source_tc_teo_l7_acc_rule.go` (添加 Offset 参数到 schema 和查询逻辑)
  - `tencentcloud/services/teo/data_source_tc_teo_l7_acc_rule_test.go` (新增测试用例)
  - `website/docs/r/teo_l7_acc_rule.html.markdown` (更新文档)
- 无 breaking changes，新增可选参数不影响现有配置
- 依赖 TencentCloud Go SDK v1.0.831+ (DescribeL7AccRules API 已支持 Offset)
