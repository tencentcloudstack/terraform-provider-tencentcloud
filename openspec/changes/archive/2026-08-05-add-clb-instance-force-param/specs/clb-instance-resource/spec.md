## ADDED Requirements

### Requirement: CLB Instance SLA Upgrade Force Flag
`tencentcloud_clb_instance` 资源 SHALL 提供可选参数 `force`（bool，默认 false），对应云 API `ModifyLoadBalancerSla` 接口的 `Force` 入参（"是否强制升级，默认否"）。当 `sla_type` 发生变更触发 `ModifyLoadBalancerSla` 调用时，`force` 的值 MUST 被透传到 `ModifyLoadBalancerSlaRequest.Force`。`force` MUST 为非 Computed 的 Optional 字段（因 `DescribeLoadBalancers` 返回结构不含 `Force`，无法回填）。

#### Scenario: force 默认为 false 且不破坏现有行为
- **WHEN** 用户未在配置中声明 `force`，且变更 `sla_type`
- **THEN** 调用 `ModifyLoadBalancerSla` 时 `Force` 取默认值 false，行为与变更前一致

#### Scenario: force=true 在 sla_type 变更时透传
- **WHEN** 用户配置 `force = true` 并变更 `sla_type` 触发 `ModifyLoadBalancerSla` 调用
- **THEN** `ModifyLoadBalancerSlaRequest.Force` MUST 被设置为 true，并随请求发送至云 API

#### Scenario: force=true 但 sla_type 未变更时不触发额外调用
- **WHEN** 用户配置 `force = true` 但 `sla_type` 未发生变更
- **THEN** 不会因 `force` 单独触发 `ModifyLoadBalancerSla` 调用

#### Scenario: Read 不回填 force
- **WHEN** 执行 `resourceTencentCloudClbInstanceRead`
- **THEN** 不对 `force` 调用 `d.Set("force", ...)`，以用户配置值为准

#### Scenario: force 非 Computed 保证无 diff
- **WHEN** 用户配置中声明 `force` 后执行 plan/apply
- **THEN** `force` 不会因 Read 未回填而产生非预期 diff
