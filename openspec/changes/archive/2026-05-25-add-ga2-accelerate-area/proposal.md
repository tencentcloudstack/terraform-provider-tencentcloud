## Why

TencentCloud Global Accelerator (GA2) 产品需要通过 Terraform 管理加速地域资源，当前 provider 中缺少对 GA2 加速地域的支持，用户无法通过 IaC 方式创建、查询、修改和删除加速地域配置。

## What Changes

- 新增 Terraform 资源 `tencentcloud_ga2_accelerate_area`，支持完整的 CRUD 生命周期管理
- 资源通过 GA2 云 API（CreateAccelerateAreas / DescribeAccelerateAreas / ModifyAccelerateAreas / DeleteAccelerateAreas）实现
- Create/Modify/Delete 接口为异步接口，调用后需轮询 DescribeAccelerateAreas 直到生效
- 在 provider.go 和 provider.md 中注册新资源

## Capabilities

### New Capabilities
- `ga2-accelerate-area-crud`: 提供 GA2 加速地域资源的完整 CRUD 操作，包括创建加速地域、查询加速地域、修改加速地域带宽/ISP类型等配置、删除加速地域，以及异步操作的轮询等待机制

### Modified Capabilities

## Impact

- 新增文件: `tencentcloud/services/ga2/resource_tc_ga2_accelerate_area.go`
- 新增文件: `tencentcloud/services/ga2/resource_tc_ga2_accelerate_area_test.go`
- 新增文件: `tencentcloud/services/ga2/service_tencentcloud_ga2.go`
- 新增文档: `tencentcloud/services/ga2/resource_tc_ga2_accelerate_area.md`
- 修改文件: `tencentcloud/provider.go`（注册资源）
- 修改文件: `tencentcloud/provider.md`（添加资源文档引用）
- 依赖: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ga2/v20250115`（已在 vendor 中）
