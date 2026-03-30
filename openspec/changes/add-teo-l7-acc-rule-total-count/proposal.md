## Why

为了支持用户在查询 tencentcloud_teo_l7_acc_rule 数据源时获取总记录数，需要在数据源中接入 TotalCount 参数。这有助于用户在进行分页查询时了解总数据量，以便进行更有效的数据处理和展示。

## What Changes

- 在 tencentcloud_teo_l7_acc_rule 数据源中新增 TotalCount 输出字段
- 调用 DescribeL7AccRules API 并解析返回的 TotalCount 参数
- 更新相关文档以说明 TotalCount 字段的用途

## Capabilities

### New Capabilities
- `teo-l7-acc-rule-total-count`: 为 tencentcloud_teo_l7_acc_rule 数据源添加 TotalCount 参数支持

### Modified Capabilities
- 无

## Impact

- 影响文件：tencentcloud/services/teo/data_source_teo_l7_acc_rule.go
- 影响文件：website/docs/t/teo_l7_acc_rule.html.markdown
- 无破坏性变更，仅新增字段
- 需要更新单元测试以验证 TotalCount 字段的正确返回
