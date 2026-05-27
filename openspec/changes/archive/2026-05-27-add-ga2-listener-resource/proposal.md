## Why

TencentCloud Global Accelerator 2.0 (GA2) 需要通过 Terraform 管理监听器资源，以便用户可以通过 IaC 方式创建、查询、修改和删除 GA2 监听器，实现全球加速实例的监听器生命周期管理。

## What Changes

- 新增 Terraform 资源 `tencentcloud_ga2_listener`，支持完整的 CRUD 操作
- 资源通过 GA2 云 API（CreateListener、DescribeListeners、ModifyListener、DeleteListener）管理监听器
- 支持 TCP/UDP/HTTP/HTTPS 协议监听器的创建和管理
- 支持异步操作轮询（通过 DescribeTaskResult 接口）
- 在 provider.go 和 provider.md 中注册新资源

## Capabilities

### New Capabilities
- `ga2-listener-resource`: GA2 监听器 Terraform 资源的完整 CRUD 实现，包括创建、读取、更新、删除监听器，以及异步任务状态轮询

### Modified Capabilities

## Impact

- 新增文件: `tencentcloud/services/ga2/resource_tc_ga2_listener.go`
- 新增文件: `tencentcloud/services/ga2/resource_tc_ga2_listener_test.go`
- 新增文件: `tencentcloud/services/ga2/resource_tc_ga2_listener.md`
- 修改文件: `tencentcloud/provider.go`（注册新资源）
- 修改文件: `tencentcloud/provider.md`（添加资源文档引用）
- 修改文件: `tencentcloud/services/ga2/service_tencentcloud_ga2.go`（添加服务层方法）
- 依赖: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ga2/v20250115`（已在 vendor 中）
