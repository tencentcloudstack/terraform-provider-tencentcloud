# Proposal: GA2 Resource Not Found Error Handling

## Summary
为 tencentcloud/services/ga2 目录下所有资源的 Read 接口中的 describe 调用补充 ResourceNotFound 错误处理和空响应处理逻辑。

## Motivation
当 describe 接口返回 Code 为 "ResourceNotFound" 的 SDK 错误时（即根据 id 找不到资源的情况），当前的 Read 方法会返回错误，导致 terraform 重新创建资源而不是正确将其从状态中移除。应该在这种场景下将 id 设置为空，返回 nil，表示资源已被删除。

同时，如果在 Create 后调用的 Read 接口（即 `d.IsNewResource()` 为 `true`），则跳过该逻辑，正常返回错误。

此外，当 describe 接口返回空响应（respData == nil）时，也需要类似的逻辑：非新建资源清除 ID，新建资源返回错误。

## Changes

### 1. Service Layer Changes (service_tencentcloud_ga2.go)
- 添加 `isGa2ResourceNotFoundError` 公共函数，用于判断错误是否为 ResourceNotFound SDK 错误
- 在所有 describe 服务函数的 retry 块中添加 ResourceNotFound 检查，遇到 ResourceNotFound 错误时返回 `NonRetryableError`：
  - DescribeEndpointGroups
  - DescribeGlobalAccelerators
  - DescribeListeners
  - DescribeForwardingRule
  - DescribeAccelerateAreas
  - DescribeForwardingPolicy
- 在上述所有 describe 服务函数的 retry 块外，如果错误是 ResourceNotFound，返回 `(nil, nil)`

### 2. Common Helper Functions (resource_tc_ga2_common.go)
- 添加 `HandleGa2ResourceNotFoundError` 函数：当 SDK 错误 Code 为 "ResourceNotFound" 且资源不是新建时，记录 WARN 日志，清除资源 ID，返回 true 表示已处理
- 添加 `HandleGa2ReadNotFound` 统一函数，处理两种场景：
  1. SDK ResourceNotFound 错误 — 委托给 `HandleGa2ResourceNotFoundError`
  2. Nil/empty response — 非新建资源记录 WARN 并清除 ID；新建资源返回错误

### 3. Resource Read Method Changes
在所有以下资源的 Read 方法中添加 ResourceNotFound 错误处理 + nil 响应处理，使用 `d.IsNewResource()` 守卫：
- `resource_tc_ga2_accelerate_area.go`
- `resource_tc_ga2_endpoint_group.go`
- `resource_tc_ga2_forwarding_policy.go`
- `resource_tc_ga2_forwarding_rule.go`
- `resource_tc_ga2_global_accelerator.go`
- `resource_tc_ga2_listener.go`

## Impact

- **Affected code**:
  - `tencentcloud/services/ga2/service_tencentcloud_ga2.go` — isGa2ResourceNotFoundError + retry block changes
  - `tencentcloud/services/ga2/resource_tc_ga2_common.go` — HandleGa2ResourceNotFoundError, HandleGa2ReadNotFound
  - `tencentcloud/services/ga2/resource_tc_ga2_accelerate_area.go` — Read method
  - `tencentcloud/services/ga2/resource_tc_ga2_endpoint_group.go` — Read method
  - `tencentcloud/services/ga2/resource_tc_ga2_forwarding_policy.go` — Read method
  - `tencentcloud/services/ga2/resource_tc_ga2_forwarding_rule.go` — Read method
  - `tencentcloud/services/ga2/resource_tc_ga2_global_accelerator.go` — Read method
  - `tencentcloud/services/ga2/resource_tc_ga2_listener.go` — Read method
- **API**: Uses existing describe APIs (no API changes)
- **Dependencies**: Uses already-imported `sdkErrors` package
- **Backward compatibility**: Fully compatible — this change improves error handling without changing schema or normal behavior
