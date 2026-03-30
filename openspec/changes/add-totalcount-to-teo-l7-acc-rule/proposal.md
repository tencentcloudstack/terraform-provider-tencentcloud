## Why

DescribeL7AccRules API 返回 TotalCount 字段（规则总数），但当前的 tencentcloud_teo_l7_acc_rule 资源没有将这个参数暴露给用户。用户无法直接通过 Terraform 获取站点 L7 加速规则的总数，这对于需要了解规则规模、监控规则数量变化或进行容量规划的用户来说是一个重要的缺失信息。

## What Changes

- 在 tencentcloud_teo_l7_acc_rule 资源的 schema 中新增 `total_count` 字段（Computed 类型）
- 在 Read 函数中从 DescribeL7AccRules API 响应的 TotalCount 字段读取值并设置到状态中
- 在资源文档中添加 `total_count` 字段的说明

## Capabilities

### New Capabilities

- `teo-l7-acc-rule-totalcount`: 在 tencentcloud_teo_l7_acc_rule 资源中添加 TotalCount 参数支持，使用户能够查询 L7 加速规则的总数

### Modified Capabilities

(无 - 这是新增参数，不修改现有行为)

## Impact

- **修改文件**:
  - `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule.go` - 在 schema 中添加 `total_count` 字段，在 Read 函数中设置该字段值
  - `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule.md` - 在文档中添加 `total_count` 字段说明
  - `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_test.go` - 添加 TotalCount 参数的测试用例

- **受影响的 API**: DescribeL7AccRules（仅读取 API 响应的 TotalCount 字段）

- **依赖**: 无新增外部依赖，使用现有的 tencentcloud-sdk-go v20220901 SDK

- **向后兼容性**: 完全兼容，`total_count` 为 Computed 字段，不影响现有配置和 state
