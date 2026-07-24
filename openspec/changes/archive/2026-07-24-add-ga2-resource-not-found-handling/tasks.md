## 1. Service Layer Changes

- [x] 1.1 在 service_tencentcloud_ga2.go 中添加 `isGa2ResourceNotFoundError` 公共函数
- [x] 1.2 修改所有 describe 服务函数的 retry 块，添加 ResourceNotFound 错误处理
- [x] 1.3 在 retry 块外处理 ResourceNotFound 错误，返回 (nil, nil)

## 2. Common Helper Functions

- [x] 2.1 在 resource_tc_ga2_common.go 中添加 `HandleGa2ResourceNotFoundError` 函数
- [x] 2.2 在 resource_tc_ga2_common.go 中添加 `HandleGa2ReadNotFound` 统一函数

## 3. Resource Read Method Changes

- [x] 3.1 修改 resource_tc_ga2_forwarding_policy.go Read 方法，添加 ResourceNotFound 和 nil 响应处理
- [x] 3.2 修改 resource_tc_ga2_global_accelerator.go Read 方法，添加 ResourceNotFound 和 nil 响应处理
- [x] 3.3 修改 resource_tc_ga2_listener.go Read 方法，添加 ResourceNotFound 和 nil 响应处理
- [x] 3.4 修改 resource_tc_ga2_endpoint_group.go Read 方法，添加 ResourceNotFound 和 nil 响应处理
- [x] 3.5 修改 resource_tc_ga2_accelerate_area.go Read 方法，添加 ResourceNotFound 和 nil 响应处理
- [x] 3.6 修改 resource_tc_ga2_forwarding_rule.go Read 方法，添加 ResourceNotFound 和 nil 响应处理

## 4. Verification

- [x] 4.1 运行 `go build ./tencentcloud/services/ga2/` 确保编译通过
- [x] 4.2 创建 changelog 文件（合并至 .changelog/4338.txt）