# clb-instance-resource Specification

## Purpose
TBD - created by archiving change add-clb-instance-timeout-block. Update Purpose after archive.
## Requirements
### Requirement: Create Operation Timeout Handling
`resourceTencentCloudClbInstanceCreate` 中所有异步任务等待（包括创建 CLB、设置安全组、设置日志、修改 target_region_info、设置 delete_protect、关联 endpoint）MUST 使用 `d.Timeout(schema.TimeoutCreate)` 替代硬编码超时。

#### Scenario: Create 中多个异步任务均使用可配置超时
- **WHEN** 创建 CLB 实例并设置安全组、日志、target_region_info 等附加配置
- **THEN** 每个异步任务等待均使用 `d.Timeout(schema.TimeoutCreate)` 作为超时上限

### Requirement: Update Operation Timeout Handling
`resourceTencentCloudClbInstanceUpdate` 中所有异步任务等待（包括修改 SLA、修改属性、设置安全组、设置日志、修改 project、EIP 操作）MUST 使用 `d.Timeout(schema.TimeoutUpdate)` 替代硬编码超时。

#### Scenario: Update 中多个异步任务均使用可配置超时
- **WHEN** 更新 CLB 实例的多个属性触发多个异步任务
- **THEN** 每个异步任务等待均使用 `d.Timeout(schema.TimeoutUpdate)` 作为超时上限

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

