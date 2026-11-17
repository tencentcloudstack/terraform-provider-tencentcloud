## 1. Service Layer Changes

- [x] 1.1 在 service_tencentcloud_ga2.go 中添加 `IsGa2ResourceNotFoundError` 公共函数，判断错误是否为 ResourceNotFound SDK 错误

## 2. Resource Read Method Changes

- [x] 2.1 修改 resource_tc_ga2_forwarding_policy.go Read 方法，添加 ResourceNotFound 和 nil 响应处理
- [x] 2.2 修改 resource_tc_ga2_listener.go Read 方法，添加 ResourceNotFound 和 nil 响应处理
- [x] 2.3 修改 resource_tc_ga2_endpoint_group.go Read 方法，添加 ResourceNotFound 和 nil 响应处理
- [x] 2.4 修改 resource_tc_ga2_accelerate_area.go Read 方法，添加 ResourceNotFound 和 nil 响应处理
- [x] 2.5 修改 resource_tc_ga2_forwarding_rule.go Read 方法，添加 ResourceNotFound 和 nil 响应处理
- [x] 2.6 修改 resource_tc_ga2_global_accelerator.go Read 方法，仅添加 nil 响应处理（该接口不返回 ResourceNotFound）

## 3. Verification

- [x] 3.1 运行 `go build ./tencentcloud/services/ga2/` 确保编译通过