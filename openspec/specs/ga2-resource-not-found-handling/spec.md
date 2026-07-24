# GA2 Resource Not Found Handling

## MODIFIED

### service_tencentcloud_ga2.go
- 新增 `IsGa2ResourceNotFoundError(err error) bool` 公共函数，判断错误是否为 `ResourceNotFound` SDK 错误

### resource_tc_ga2_forwarding_policy.go
- 在 Read 方法中添加 ResourceNotFound 错误处理和 nil 响应处理：先检查 err，再检查 respData，`!d.IsNewResource()` 时清除 ID

### resource_tc_ga2_listener.go
- 在 Read 方法中添加 ResourceNotFound 错误处理和 nil 响应处理

### resource_tc_ga2_endpoint_group.go
- 在 Read 方法中添加 ResourceNotFound 错误处理和 nil 响应处理

### resource_tc_ga2_accelerate_area.go
- 在 Read 方法中添加 ResourceNotFound 错误处理和 nil 响应处理

### resource_tc_ga2_forwarding_rule.go
- 在 Read 方法中添加 ResourceNotFound 错误处理和 nil 响应处理

### resource_tc_ga2_global_accelerator.go
- 在 Read 方法中仅处理 `respData == nil`（该接口不返回 ResourceNotFound 错误）