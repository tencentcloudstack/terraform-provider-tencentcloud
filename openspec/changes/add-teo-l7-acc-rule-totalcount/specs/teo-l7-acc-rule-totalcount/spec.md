# teo-l7-acc-rule-totalcount Specification

## Purpose
在 tencentcloud_teo_l7_acc_rule 数据源中暴露 TotalCount 字段，使用户能够获取七层加速规则的总数信息，支持配置管理和自动化场景下的统计与监控。

## Requirements

### Requirement: TotalCount 字段输出
tencentcloud_teo_l7_acc_rule 数据源 SHALL 从 DescribeL7AccRules API 响应中读取 TotalCount 字段，并将其作为 Computed 输出属性暴露。

#### Scenario: 查询时返回 TotalCount
- **WHEN** 用户查询 tencentcloud_teo_l7_acc_rule 数据源
- **THEN** 输出包含 TotalCount 字段，值为 DescribeL7AccRules API 返回的总记录数
- **AND** TotalCount 字段为 Computed 类型

#### Scenario: 向后兼容性
- **WHEN** 现有 Terraform 配置引用 tencentcloud_teo_l7_acc_rule
- **THEN** 所有现有字段保持不变
- **AND** TotalCount 作为新增字段可用，不影响现有配置
