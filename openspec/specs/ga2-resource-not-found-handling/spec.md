# GA2 Resource Not Found Handling

## MODIFIED

### service_tencentcloud_ga2.go
- 新增 `isGa2ResourceNotFoundError(err error) bool` 公共函数，判断错误是否为 `ResourceNotFound` SDK 错误
- 在 `DescribeGa2EndpointGroupById` 的 retry 块中添加 ResourceNotFound 处理
- 在 `DescribeGa2GlobalAcceleratorById` 的 retry 块中添加 ResourceNotFound 处理
- 在 `DescribeGa2ListenerById` 的 retry 块中添加 ResourceNotFound 处理
- 在 `DescribeGa2ForwardingRuleById` 的 retry 块中添加 ResourceNotFound 处理
- 在 `describeGa2AccelerateAreas` 的 retry 块中添加 ResourceNotFound 处理
- 在 `DescribeGa2ForwardingPolicyById` 的 retry 块中添加 ResourceNotFound 处理

### resource_tc_ga2_accelerate_area.go
- 在 Read 方法中，当 `respData == nil` 时添加 `d.IsNewResource()` 检查，若 IsNewResource 则返回 error

### resource_tc_ga2_endpoint_group.go
- 在 Read 方法中，当 `respData == nil` 时添加 `d.IsNewResource()` 检查，若 IsNewResource 则返回 error

### resource_tc_ga2_forwarding_policy.go
- 在 Read 方法中，当 `respData == nil` 时添加 `d.IsNewResource()` 检查，若 IsNewResource 则返回 error

### resource_tc_ga2_forwarding_rule.go
- 在 Read 方法中，当 `respData == nil` 时添加 `d.IsNewResource()` 检查，若 IsNewResource 则返回 error

### resource_tc_ga2_global_accelerator.go
- 在 Read 方法中，当 `respData == nil` 时添加 `d.IsNewResource()` 检查，若 IsNewResource 则返回 error

### resource_tc_ga2_listener.go
- 在 Read 方法中，当 `respData == nil` 时添加 `d.IsNewResource()` 检查，若 IsNewResource 则返回 error