# GA2 Resource Not Found Handling

## MODIFIED

### service_tencentcloud_ga2.go
- 新增 `IsGa2ResourceNotFoundError(err error) bool` 公共函数，判断错误是否为 `ResourceNotFound` SDK 错误（Code == "ResourceNotFound"）。由各资源 Read 方法调用。

### resource_tc_ga2_forwarding_policy.go
- 在 Read 方法中，`DescribeGa2ForwardingPolicyById` 返回后：
  - 先检查 `err != nil`，若为 `!d.IsNewResource() && IsGa2ResourceNotFoundError(err)` 则记录 WARN 日志、清除 ID、返回 nil；否则透传 err
  - 再检查 `respData == nil`，若 `!d.IsNewResource()` 则记录 WARN 日志并清除 ID；若 `IsNewResource()` 则返回 "not found after creation" 错误

### resource_tc_ga2_listener.go
- 同上，在 Read 方法中处理 `DescribeGa2ListenerById` 的 ResourceNotFound 错误和 nil 响应

### resource_tc_ga2_endpoint_group.go
- 同上，在 Read 方法中处理 `DescribeGa2EndpointGroupById` 的 ResourceNotFound 错误和 nil 响应

### resource_tc_ga2_accelerate_area.go
- 同上，在 Read 方法中处理 `DescribeGa2AccelerateAreaById` 的 ResourceNotFound 错误和 nil 响应

### resource_tc_ga2_forwarding_rule.go
- 同上，在 Read 方法中处理 `DescribeGa2ForwardingRuleById` 的 ResourceNotFound 错误和 nil 响应

### resource_tc_ga2_global_accelerator.go
- 在 Read 方法中仅处理 `respData == nil` 情况（`DescribeGlobalAccelerators` 接口不返回 ResourceNotFound 错误）