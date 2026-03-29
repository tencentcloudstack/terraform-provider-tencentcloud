## Why

为 tencentcloud_teo_l7_acc_rule 资源添加 TotalCount 输出参数，以便用户能够获取当前站点下 L7 访问控制规则的总数。这个信息对于资源管理、监控和配置审计非常重要，可以帮助用户了解规则的总规模。

## What Changes

- 在 tencentcloud_teo_l7_acc_rule 资源的 schema 中新增 `total_count` 输出字段（Computed）
- 在 Read 方法中，从 DescribeL7AccRules API 响应中获取 TotalCount 值并设置到该字段

## Capabilities

### New Capabilities
- `l7-acc-rule-totalcount`: 为 tencentcloud_teo_l7_acc_rule 资源添加 TotalCount 输出参数，用于返回当前站点下 L7 访问控制规则的总数

### Modified Capabilities
- None

## Impact

- **Affected Code**: `/repo/tencentcloud/services/teo/resource_tc_teo_l7_acc_rule.go`
  - 修改 `ResourceTencentCloudTeoL7AccRule()` 函数，在 Schema 中添加 `total_count` 字段
  - 修改 `resourceTencentCloudTeoL7AccRuleRead()` 函数，在读取 API 响应后设置 `total_count` 字段
- **API**: 使用现有的 DescribeL7AccRules API，该 API 响应已包含 TotalCount 字段
- **Dependencies**: 无新增依赖
- **Backward Compatibility**: 新增字段为 Computed 类型，不会破坏现有配置和状态
