## Why

TencentCloud Global Accelerator (GA2) 产品需要通过 Terraform 管理全球加速实例的生命周期（创建、查询、修改、删除），当前 provider 中尚未提供该资源，用户无法通过 IaC 方式管理 GA2 全球加速实例。

## What Changes

- 新增 Terraform 资源 `tencentcloud_ga2_global_accelerator`，支持全球加速实例的完整 CRUD 生命周期管理
- 支持创建时指定名称、计费模式、描述、跨境类型、跨境承诺标志和标签
- 支持修改名称、描述、跨境类型和跨境承诺标志
- 支持删除全球加速实例
- 所有写操作（Create/Modify/Delete）为异步接口，需轮询实例状态直到操作完成
- 在 provider.go 和 provider.md 中注册新资源

## Capabilities

### New Capabilities

- `ga2-global-accelerator-resource`: 全球加速实例的 CRUD 资源管理，包括创建、读取、更新、删除操作，以及异步操作的状态轮询

### Modified Capabilities

（无）

## Impact

- 新增文件: `tencentcloud/services/ga2/resource_tc_ga2_global_accelerator.go`
- 新增文件: `tencentcloud/services/ga2/resource_tc_ga2_global_accelerator_test.go`
- 新增文件: `tencentcloud/services/ga2/resource_tc_ga2_global_accelerator.md`
- 修改文件: `tencentcloud/provider.go`（注册资源）
- 修改文件: `tencentcloud/provider.md`（添加资源文档引用）
- 可能修改: `tencentcloud/services/ga2/service_tencentcloud_ga2.go`（添加服务层方法）
- 依赖: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ga2/v20250115`（已在 vendor 中）
