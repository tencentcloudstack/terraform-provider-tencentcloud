## Why

CreateOriginGroup API 支持 HostHeader 参数，用于配置回源 Host Header，但当前 Terraform 资源的 Create 函数未传入此参数，导致用户无法在创建源站组时直接设置该参数。

## What Changes

- 在 `tencentcloud_teo_origin_group` 资源的 Create 函数中添加对 `host_header` 参数的处理，将用户配置的值传递给 CreateOriginGroup API 的 HostHeader 参数
- 确保与现有的 Update 和 Read 函数保持一致的参数处理逻辑

## Capabilities

### New Capabilities
- `teo-origin-group-hostheader`: 为 tencentcloud_teo_origin_group 资源的 Create 操作添加 HostHeader 参数支持，允许用户在创建源站组时配置回源 Host Header

### Modified Capabilities

## Impact

- 修改文件: `tencentcloud/services/teo/resource_tc_teo_origin_group.go`
- 影响函数: `resourceTencentCloudTeoOriginGroupCreate`
- 无 breaking changes，仅新增参数传递功能
