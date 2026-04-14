## Why

为 `tencentcloud_teo_origin_group` 资源的创建操作（`CreateOriginGroup` API）添加 `HostHeader` 参数支持，使其能够在创建源站组时设置回源 Host Header。该参数在 Schema 中已定义，在读取和更新操作中已实现，但创建操作中缺失此功能，导致无法在资源创建时设置该参数。

## What Changes

- 在 `resource_tc_teo_origin_group.go` 的 `resourceTencentCloudTeoOriginGroupCreate` 函数中添加 `HostHeader` 参数的设置逻辑
- 将 `host_header` schema 参数映射到 `CreateOriginGroupRequest.HostHeader` 字段

## Capabilities

### New Capabilities

No new capabilities - this is a bug fix to complete missing functionality in the create operation.

### Modified Capabilities

- `teo-origin-group`: 完善 `tencentcloud_teo_origin_group` 资源的创建功能，确保 `host_header` 参数在创建操作中生效

## Impact

- 修改文件: `tencentcloud/services/teo/resource_tc_teo_origin_group.go`
- 不破坏向后兼容: 现有配置和 state 不受影响
- 不需要文档更新: `host_header` 参数已在文档中存在
