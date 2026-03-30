## ADDED Requirements

### Requirement: 数据源支持 TotalCount 输出字段
tencentcloud_teo_l7_acc_rule 数据源 SHALL 提供一个名为 total_count 的输出字段，用于返回符合条件的记录总数。

#### Scenario: 成功查询时返回 TotalCount
- **WHEN** 用户查询 tencentcloud_teo_l7_acc_rule 数据源且 API 返回成功
- **THEN** 数据源 SHALL 从 DescribeL7AccRules API 响应中读取 TotalCount 字段
- **AND** 数据源 SHALL 将 TotalCount 值填充到 total_count 字段中
- **AND** total_count 字段的值 SHALL 为非负整数

#### Scenario: API 返回空列表时 TotalCount 为 0
- **WHEN** 查询条件没有任何匹配记录
- **THEN** DescribeL7AccRules API SHALL 返回 TotalCount = 0
- **AND** 数据源的 total_count 字段 SHALL 为 0

#### Scenario: TotalCount 字段为只读
- **WHEN** 用户尝试在 Terraform 配置中设置 total_count 字段
- **THEN** 配置 SHALL 被拒绝
- **AND** 系统 SHALL 返回错误提示 total_count 为只读字段（Computed 属性）
