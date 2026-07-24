# Proposal: GA2 Resource Not Found Error Handling

## Summary
为 tencentcloud/services/ga2 目录下所有资源的 Read 接口中的 describe 调用补充 ResourceNotFound 错误处理逻辑。

## Motivation
当 describe 接口返回 Code 为 "ResourceNotFound" 的 SDK 错误时（即根据 id 找不到资源的情况），当前的 Read 方法会返回错误，导致 terraform 重新创建资源而不是正确将其从状态中移除。应该在这种场景下将 id 设置为空，返回 nil，表示资源已被删除。

同时，如果在 Create 后调用的 Read 接口（即 `d.IsNewResource()` 为 `true`），则跳过该逻辑，正常返回错误。

## Changes
1. 在 `service_tencentcloud_ga2.go` 中添加公共函数 `isGa2ResourceNotFoundError`，用于判断错误是否为 ResourceNotFound SDK 错误
2. 在以下 describe 服务的 retry 块中添加 `isGa2ResourceNotFoundError` 检查，遇到 ResourceNotFound 错误时返回 `NonRetryableError`：
   - DescribeEndpointGroups
   - DescribeGlobalAccelerators
   - DescribeListeners
   - DescribeForwardingRule
   - DescribeAccelerateAreas
   - DescribeForwardingPolicy
3. 在上述所有 describe 服务的 retry 块外，如果错误是 ResourceNotFound，返回 `(nil, nil)`
4. 在以下资源的 Read 方法中，当 `respData == nil` 时添加 `d.IsNewResource()` 检查：
   - resource_tc_ga2_accelerate_area.go
   - resource_tc_ga2_endpoint_group.go
   - resource_tc_ga2_forwarding_policy.go
   - resource_tc_ga2_forwarding_rule.go
   - resource_tc_ga2_global_accelerator.go
   - resource_tc_ga2_listener.go
