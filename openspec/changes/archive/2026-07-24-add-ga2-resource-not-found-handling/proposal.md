# Proposal: GA2 Resource Not Found Error Handling

## Summary
为 tencentcloud/services/ga2 目录下资源的 Read 接口中的 describe 调用补充 ResourceNotFound 错误处理和空响应处理逻辑。

## Motivation
当 describe 接口返回 Code 为 "ResourceNotFound" 的 SDK 错误时（即根据 id 找不到资源的情况），当前的 Read 方法会返回错误，导致 terraform 重新创建资源而不是正确将其从状态中移除。应该在这种场景下将 id 设置为空，返回 nil，表示资源已被删除。

同时，如果在 Create 后调用的 Read 接口（即 `d.IsNewResource()` 为 `true`），则跳过该逻辑，正常返回错误。

此外，当 describe 接口返回空响应（respData == nil）时，也需要类似的逻辑：非新建资源清除 ID，新建资源返回错误。

**注意**：`DescribeGlobalAccelerators` 接口不会返回 ResourceNotFound 错误，因此 `resource_tc_ga2_global_accelerator.go` 的 Read 方法仅处理 `respData == nil` 情况。

## Changes

### 1. Service Layer Changes (service_tencentcloud_ga2.go)
- 添加 `IsGa2ResourceNotFoundError(err error) bool` 公共函数，用于判断错误是否为 ResourceNotFound SDK 错误。由各资源 Read 方法调用。

### 2. Resource Read Method Changes
在所有以下资源的 Read 方法中添加 ResourceNotFound 错误处理 + nil 响应处理：
- `resource_tc_ga2_accelerate_area.go`
- `resource_tc_ga2_endpoint_group.go`
- `resource_tc_ga2_forwarding_policy.go`
- `resource_tc_ga2_forwarding_rule.go`
- `resource_tc_ga2_listener.go`

处理逻辑如下（`err` 检查必须在 `respData` 检查之前）：
```go
respData, err := service.DescribeXxxById(ctx, ...)
if err != nil {
    if !d.IsNewResource() && IsGa2ResourceNotFoundError(err) {
        log.Printf("[WARN]...")
        d.SetId("")
        return nil
    }
    return err
}

if respData == nil {
    log.Printf("[WARN]...")
    if d.IsNewResource() {
        return fmt.Errorf("... not found after creation")
    }
    d.SetId("")
    return nil
}
```

### 3. Exception: global_accelerator
`resource_tc_ga2_global_accelerator.go` 的 Read 方法不添加 `IsGa2ResourceNotFoundError` 判断，因为 `DescribeGlobalAccelerators` 接口不会返回 ResourceNotFound 错误。

## Impact

- **Affected code**:
  - `tencentcloud/services/ga2/service_tencentcloud_ga2.go` — `IsGa2ResourceNotFoundError` 公共函数
  - `tencentcloud/services/ga2/resource_tc_ga2_accelerate_area.go` — Read method
  - `tencentcloud/services/ga2/resource_tc_ga2_endpoint_group.go` — Read method
  - `tencentcloud/services/ga2/resource_tc_ga2_forwarding_policy.go` — Read method
  - `tencentcloud/services/ga2/resource_tc_ga2_forwarding_rule.go` — Read method
  - `tencentcloud/services/ga2/resource_tc_ga2_listener.go` — Read method
  - `tencentcloud/services/ga2/resource_tc_ga2_global_accelerator.go` — Read method（仅 nil 响应处理）
- **API**: Uses existing describe APIs (no API changes)
- **Dependencies**: Uses already-imported `sdkErrors` package
- **Backward compatibility**: Fully compatible — this change improves error handling without changing schema or normal behavior
