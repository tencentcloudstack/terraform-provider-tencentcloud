## Why

`tencentcloud_clb_instance` 资源在更新 `sla_type` 时已调用云 API `ModifyLoadBalancerSla`，但当前调用未透传 `Force` 字段。当需要将共享型负载均衡升级为性能容量型实例（或调整性能容量型规格）时，部分场景需要"强制升级"才能成功，而现有 Provider 无法表达这一意图，导致升级在受保护场景下失败且无法通过 Terraform 绕过。需要新增 `force` 参数，让用户在调用 `ModifyLoadBalancerSla` 时能够控制是否强制升级。

## What Changes

- 在 `tencentcloud_clb_instance` 资源 schema 中新增可选参数 `force`（bool 类型，默认 false），对应云 API `ModifyLoadBalancerSla` 接口的 `Force` 入参（"是否强制升级，默认否"）。
- 在 `resourceTencentCloudClbInstanceUpdate` 中，当 `sla_type` 发生变更调用 `ModifyLoadBalancerSla` 时，将 `force` 参数透传到 `slaRequest.Force`。
- 在 `resourceTencentCloudClbInstanceRead` 中读取并回填 `force`（当云 API 返回该字段时）。
- 更新资源文档 `resource_tc_clb_instance.md`，补充 `force` 参数的示例与说明。
- 在单元测试 `resource_tc_clb_instance_test.go` 中补充针对 `force` 参数的测试用例。

> 注：需求中另提到 `ModifyTags` 接口的 `tags` 入参。经核查 `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/` 目录，clb 包下不存在 `ModifyTags` 接口（该接口仅存在于 tke、emr 等其他产品包中），且 `tencentcloud_clb_instance` 资源已存在 `tags` 参数（schema.TypeMap）。依据"只有云 API 支持的行为才可添加到 Provider"的约束，`tags` 不纳入本次变更，本次仅新增 `force` 一个参数。

## Capabilities

### New Capabilities
<!-- 无新增能力 -->

### Modified Capabilities
- `clb-instance-resource`: 在 `tencentcloud_clb_instance` 资源上新增 `force` 参数，使更新 `sla_type`（调用 `ModifyLoadBalancerSla`）时支持透传"是否强制升级"标志。

## Impact

- **受影响代码**:
  - `tencentcloud/services/clb/resource_tc_clb_instance.go`：schema 新增 `force` 字段；`Update` 中 `ModifyLoadBalancerSla` 调用块透传 `Force`；`Read` 回填 `force`。
  - `tencentcloud/services/clb/resource_tc_clb_instance_test.go`：补充 `force` 参数测试用例。
  - `tencentcloud/services/clb/resource_tc_clb_instance.md`：补充文档说明。
- **云 API 依赖**: `ModifyLoadBalancerSla`（`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317`），入参 `Force *bool`，已存在于 vendor 中，无需更新依赖。
- **兼容性**: `force` 为新增 Optional 字段，默认 false，向后兼容，不影响现有配置与 state。
