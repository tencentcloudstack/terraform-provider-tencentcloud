## Why

为了在 Terraform Provider 的 teo_l7_acc_rule 资源中新增 TotalCount 字段，该字段用于显示规则总数。这个字段是从 DescribeL7AccRules API 的响应中获取的，可以让用户在查询资源时了解规则的总数量。

## What Changes

- 在 tencentcloud_teo_l7_acc_rule 资源的 Schema 中新增 TotalCount 字段（int 类型，Optional）
- 更新 Read 函数以从 DescribeL7AccRules API 响应中读取并设置 TotalCount 字段
- TotalCount 字段为只读字段（Computed），仅从 API 响应中获取，不需要在 Create/Update/Delete 中处理
- 更新相关的单元测试代码以测试新字段的读取逻辑

## Capabilities

### New Capabilities
- `teo-l7-acc-rule-totalcount`: 为 tencentcloud_teo_l7_acc_rule 资源添加 TotalCount 只读字段

### Modified Capabilities
- 无

## Impact

受影响的代码文件：
- tencentcloud/services/teo/resource_tc_teo_l7_acc_rule.go（资源实现文件）
- tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_test.go（单元测试文件）
- website/docs/r/teo_l7_acc_rule.md（文档文件，如需要）
