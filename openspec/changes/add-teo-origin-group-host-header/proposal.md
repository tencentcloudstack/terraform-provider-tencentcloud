## Why

支持在 tencentcloud_teo_origin_group 资源中配置回源 Host Header。这是腾讯云 TEO（TencentCloud EdgeOne）源站组管理的必要功能，允许用户在源站为 HTTP 类型时自定义回源请求的 Host Header。

## What Changes

- 在 tencentcloud_teo_origin_group 资源的 schema 中新增 `host_header` 字段（可选，string 类型）
- 更新 Create 函数，支持在创建源站组时传递 HostHeader 参数
- 更新 Read 函数，从 API 响应中读取并返回 HostHeader 配置
- 更新 Update 函数，支持修改 HostHeader 参数
- HostHeader 参数仅在源站类型为 HTTP 时生效（通过规则引擎验证）
- 更新资源单元测试，覆盖新增字段的测试场景
- 更新资源验收测试，验证字段的实际行为

## Capabilities

### New Capabilities
- `teo-origin-group-host-header`: 支持 TEO 源站组资源的 HostHeader 配置能力，允许用户自定义回源请求的 Host Header

### Modified Capabilities
无（仅新增可选字段，不改变现有资源的行为要求）

## Impact

- 影响文件：
  - `tencentcloud/services/teo/resource_tencentcloud_teo_origin_group.go`（资源实现）
  - `tencentcloud/services/teo/resource_tencentcloud_teo_origin_group_test.go`（单元测试）
  - `website/docs/r/teo_origin_group.html.markdown`（资源文档）
- 依赖的 API：
  - CreateOriginGroup API（支持传入 HostHeader 参数）
  - DescribeOriginGroup API（返回 HostHeader 配置）
  - ModifyOriginGroup API（支持更新 HostHeader 参数）
- 无第三方依赖变更
- 向后兼容：新增可选字段，不影响现有配置