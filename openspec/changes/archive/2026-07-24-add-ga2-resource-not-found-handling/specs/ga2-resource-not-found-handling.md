# GA2 Resource Not Found Handling

## MODIFIED

### service_tencentcloud_ga2.go
- 新增 `isGa2ResourceNotFoundError(err error) bool` 公共函数，判断错误是否为 `ResourceNotFound` SDK 错误（Code == "ResourceNotFound"）
- 在以下 describe 服务函数的 retry 块中添加 ResourceNotFound 检查，遇到时返回 `NonRetryableError`，retry 块外判断 ResourceNotFound 时返回 `(nil, nil)`：
  - `DescribeGa2EndpointGroupById`
  - `DescribeGa2GlobalAcceleratorById`
  - `DescribeGa2ListenerById`
  - `DescribeGa2ForwardingRuleById`
  - `describeGa2AccelerateAreas`
  - `DescribeGa2ForwardingPolicyById`

### resource_tc_ga2_common.go
- 新增 `HandleGa2ResourceNotFoundError(err error, d *schema.ResourceData, resourceName, logId string) bool` 函数：当 SDK 错误 Code 为 "ResourceNotFound" 且资源不是新建时（`!d.IsNewResource()`），记录 WARN 日志，清除资源 ID，返回 true
- 新增 `HandleGa2ReadNotFound(err error, respData interface{}, d *schema.ResourceData, resourceName, logId string) (bool, error)` 统一函数，处理两种场景：
  1. SDK ResourceNotFound 错误 — 委托给 `HandleGa2ResourceNotFoundError`
  2. Nil/empty response — 非新建资源记录 WARN 并清除 ID；新建资源返回错误

### resource_tc_ga2_forwarding_policy.go
- 在 Read 方法中，当 `DescribeGa2ForwardingPolicyById` 返回错误时，额外检查是否为 SDK `ResourceNotFound` 错误（防御性内联处理）；若 `!d.IsNewResource()` 则清除 ID 返回 nil
- 在 Read 方法中，当 `respData == nil` 时（服务层已将 ResourceNotFound 转为 nil），若 `!d.IsNewResource()` 清除 ID 返回 nil，若 `IsNewResource()` 返回 error

### resource_tc_ga2_accelerate_area.go
- 在 Read 方法中，当 `respData == nil` 时（服务层已将 ResourceNotFound 转为 nil），若 `!d.IsNewResource()` 清除 ID 返回 nil，若 `IsNewResource()` 返回 error

### resource_tc_ga2_endpoint_group.go
- 在 Read 方法中，当 `respData == nil` 时（服务层已将 ResourceNotFound 转为 nil），若 `!d.IsNewResource()` 清除 ID 返回 nil，若 `IsNewResource()` 返回 error

### resource_tc_ga2_global_accelerator.go
- 在 Read 方法中，当 `respData == nil` 时（服务层已将 ResourceNotFound 转为 nil），若 `!d.IsNewResource()` 清除 ID 返回 nil，若 `IsNewResource()` 返回 error

### resource_tc_ga2_listener.go
- 在 Read 方法中，当 `respData == nil` 时（服务层已将 ResourceNotFound 转为 nil），若 `!d.IsNewResource()` 清除 ID 返回 nil，若 `IsNewResource()` 返回 error

### resource_tc_ga2_forwarding_rule.go
- 在 Read 方法中，当 `respData == nil` 时（服务层已将 ResourceNotFound 转为 nil），若 `!d.IsNewResource()` 清除 ID 返回 nil，若 `IsNewResource()` 返回 error