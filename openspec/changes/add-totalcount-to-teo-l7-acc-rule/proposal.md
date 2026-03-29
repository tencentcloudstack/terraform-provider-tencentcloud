## Why

DescribeL7AccRules API 响应中包含 TotalCount 字段，用于返回规则总数。当前数据源未接入此参数，导致用户无法直接获取规则总数，需要通过遍历规则列表手动计算。

## What Changes

- 在 tencentcloud_teo_l7_acc_rule 数据源中新增 TotalCount 字段
- 从 DescribeL7AccRules API 响应中读取 TotalCount 参数并映射到 schema
- 确保 TotalCount 字段在数据源配置中可访问

## Capabilities

### New Capabilities
- `l7-acc-rule-totalcount`: 在 TEOS 七层访问规则数据源中支持 TotalCount 字段，用于返回规则总数

### Modified Capabilities

## Impact

- 修改 `tencentcloud/services/teo/data_source_tc_teo_l7_acc_rule.go` 文件
- 更新数据源 schema，新增 TotalCount 字段定义
- 更新 Read 函数，从 API 响应中提取并设置 TotalCount 值
- 更新相关文档和示例
