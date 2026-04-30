## Why

Terraform 用户在使用 TEO（边缘安全加速平台）多通道安全加速网关功能时，需要查询可用地域列表以便配置网关资源。当前 provider 中缺少对应的数据源，用户无法通过 Terraform 查询多通道安全加速网关可用地域信息。

## What Changes

- 新增数据源 `tencentcloud_teo_multi_path_gateway_region`，用于查询多通道安全加速网关可用地域列表
- 数据源调用云 API `DescribeMultiPathGatewayRegions`，入参为 `zone_id`（站点 ID），出参为 `gateway_regions`（可用地域列表）
- 在 `provider.go` 和 `provider.md` 中注册新数据源
- 生成对应的 `.md` 文档文件

## Capabilities

### New Capabilities
- `teo-multi-path-gateway-region-datasource`: 新增 TEO 多通道安全加速网关可用地域数据源，支持按站点 ID 查询可用地域列表（包含地域 ID、中文名称、英文名称）

### Modified Capabilities

## Impact

- 新增文件: `tencentcloud/services/teo/data_source_tc_teo_multi_path_gateway_region.go`
- 新增文件: `tencentcloud/services/teo/data_source_tc_teo_multi_path_gateway_region_test.go`
- 新增文件: `tencentcloud/services/teo/data_source_tc_teo_multi_path_gateway_region.md`
- 修改文件: `tencentcloud/provider.go`（注册新数据源）
- 修改文件: `tencentcloud/provider.md`（注册新数据源）
- 依赖云 API: `teo/v20220901.DescribeMultiPathGatewayRegions`
