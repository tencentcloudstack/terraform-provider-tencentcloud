# fwm-edge-acl-rule Specification

## Purpose
TBD - created by archiving change add-fwm-edge-acl-rule-ruleids-param. Update Purpose after archive.
## Requirements
### Requirement: RuleIds 出参记录创建返回的规则 ID 列表
The `tencentcloud_fwm_edge_acl_rule` resource SHALL 暴露一个名为 `RuleIds` 的出参（Computed）字段，其类型为字符串列表，用于承接云 API `CreateEdgeAclRule` 返回的 `response.RuleIds`。该字段 SHALL NOT 可由用户配置（非 Optional / 非 Required），且 SHALL NOT 触发 ForceNew。

#### Scenario: 创建规则成功后回填 RuleIds
- **WHEN** 调用 `CreateEdgeAclRule` 成功且 `response.Response.RuleIds` 非空
- **THEN** 资源 SHALL 通过 `d.Set("RuleIds", RuleIds)` 将返回的规则 ID 列表写入 Terraform state

#### Scenario: 创建返回空 RuleIds 时报错
- **WHEN** 调用 `CreateEdgeAclRule` 成功但 `response` / `response.Response` / `RuleIds` 为空
- **THEN** 资源 SHALL 返回 `NonRetryableError`，SHALL NOT 写入空 id 到 state

#### Scenario: Read 刷新 RuleIds
- **WHEN** 执行资源 Read 并通过 `DescribeEdgeAclRules` 查询到规则列表
- **THEN** 资源 SHALL 将查询到的规则 ID 列表回填到 `RuleIds` 字段，保持 state 与云端一致

#### Scenario: RuleIds 不可由用户配置
- **WHEN** 用户在 Terraform 配置中显式为 `RuleIds` 赋值
- **THEN** Terraform SHALL 拒绝该配置（字段为 Computed，不接受用户输入）