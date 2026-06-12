## Why

TencentCloud EdgeOne (TEO) 支持共享 CNAME 功能，允许用户创建共享 CNAME 并将多个加速域名绑定到同一个 CNAME 记录上，简化 DNS 管理。当前 Terraform Provider 缺少对该资源的支持，用户无法通过 IaC 方式管理共享 CNAME 的生命周期。

## What Changes

- 新增 Terraform 资源 `tencentcloud_teo_shared_cname`，支持共享 CNAME 的完整 CRUD 生命周期管理
- 支持创建共享 CNAME（指定前缀和描述）
- 支持查询共享 CNAME 详情
- 支持修改共享 CNAME 的描述和 IP SSL 设置
- 支持删除共享 CNAME
- 资源 ID 使用 `zone_id` 和 `shared_cname` 的联合 ID

## Capabilities

### New Capabilities
- `teo-shared-cname-resource`: 提供 `tencentcloud_teo_shared_cname` 资源的完整 CRUD 生命周期管理，包括创建、读取、更新和删除共享 CNAME

### Modified Capabilities

## Impact

- 新增文件: `tencentcloud/services/teo/resource_tc_teo_shared_cname.go`
- 新增测试文件: `tencentcloud/services/teo/resource_tc_teo_shared_cname_test.go`
- 新增文档文件: `tencentcloud/services/teo/resource_tc_teo_shared_cname.md`
- 修改文件: `tencentcloud/provider.go`（注册新资源）
- 修改文件: `tencentcloud/provider.md`（添加资源文档引用）
- 依赖云 API SDK: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901`
