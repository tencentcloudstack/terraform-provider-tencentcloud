# Spec: L7 Access Rule TotalCount

## ADDED Requirements

### Requirement: 数据源应支持 TotalCount 字段

数据源 `tencentcloud_teo_l7_acc_rule` 应在 schema 中定义一个 `total_count` 字段（类型为 integer），该字段从 DescribeL7AccRules API 响应中读取 `TotalCount` 参数的值。

#### Scenario: 查询七层访问规则时返回 TotalCount

- **WHEN** 用户通过数据源 `tencentcloud_teo_l7_acc_rule` 查询七层访问规则
- **THEN** 数据源的响应中应包含 `total_count` 字段，且值为 API 返回的 `TotalCount` 参数值

#### Scenario: API 响应中不包含 TotalCount

- **WHEN** DescribeL7AccRules API 响应中不包含 `TotalCount` 参数或该参数为 null
- **THEN** 数据源的 `total_count` 字段应返回 0

### Requirement: TotalCount 字段应为只读（Computed）

`total_count` 字段应定义为 computed 类型，用户无法在配置中指定该字段的值。

#### Scenario: 尝试在配置中设置 total_count 值

- **WHEN** 用户在数据源配置中设置 `total_count` 字段
- **THEN** Terraform 应忽略该值并使用 API 返回的实际值

### Requirement: 保持向后兼容

添加 `total_count` 字段不应影响现有数据源的其他字段或行为。

#### Scenario: 查询不包含 total_count 的现有数据源配置

- **WHEN** 用户的现有数据源配置不包含 `total_count` 字段
- **THEN** 数据源应正常工作，仅返回原有的字段，不报错
