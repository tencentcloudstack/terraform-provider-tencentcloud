## Why

在 tencentcloud_teo_l7_acc_rule 数据源的读取操作中，DescribeL7AccRules API 返回了 TotalCount 参数，但当前实现未接入该参数。TotalCount 表示符合条件的规则总数，对于用户了解查询结果的完整数量非常重要，特别是在分页查询场景下。接入该参数可以提供更好的用户体验，让用户能够准确获取资源总数信息。

## What Changes

- 在 tencentcloud_teo_l7_acc_rule 数据源中接入 TotalCount 参数的响应处理
- 在数据源 schema 中新增 TotalCount 字段（Computed 类型）
- 更新相关的测试用例以验证 TotalCount 参数的正确性

## Capabilities

### New Capabilities

### Modified Capabilities

- `teo-l7-acc-rule`: 在数据源的读取操作中新增 TotalCount 响应字段的处理逻辑

## Impact

- 修改文件: `tencentcloud/services/teo/data_source_tc_teo_l7_acc_rule.go`
- 修改文件: `tencentcloud/services/teo/data_source_tc_teo_l7_acc_rule_test.go`
- 修改文件: `website/docs/r/teo.html.markdown`（或相应的数据源文档）
- 影响: tencentcloud_teo_l7_acc_rule 数据源的用户，可以在查询结果中获取到 TotalCount 字段
