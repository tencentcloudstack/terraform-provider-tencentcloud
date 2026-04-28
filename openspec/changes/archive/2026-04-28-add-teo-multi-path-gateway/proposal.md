## Why

TEO (EdgeOne) 多通道安全加速网关（MultiPathGateway）是 TE0 产品的核心能力之一，支持用户创建云上网关和自有网关来管理多通道加速线路。当前 Terraform Provider 中缺少对该资源的支持，用户无法通过 Terraform 管理多通道安全加速网关的生命周期。需要新增 `tencentcloud_teo_multi_path_gateway` 资源，使用户能够通过 Terraform 完成网关的创建、查询、修改和删除操作。

## What Changes

- 新增 Terraform 资源 `tencentcloud_teo_multi_path_gateway`，类型为 RESOURCE_KIND_GENERAL
- 支持通过 `CreateMultiPathGateway` 接口创建多通道安全加速网关，包含 zone_id、gateway_type、gateway_name、gateway_port、region_id、gateway_ip 参数
- 支持通过 `DescribeMultiPathGateways` 接口查询网关详情
- 支持通过 `ModifyMultiPathGateway` 接口修改网关名称、IP 和端口
- 支持通过 `DeleteMultiPathGateway` 接口删除网关
- 在 provider.go 中注册新资源
- 在 provider.md 中添加资源文档入口
- 生成对应的 .md 文档文件

## Capabilities

### New Capabilities
- `teo-multi-path-gateway`: 新增 TEO 多通道安全加速网关资源的完整 CRUD 管理，包括创建、读取、更新、删除以及导入功能

### Modified Capabilities

## Impact

- 新增文件: `tencentcloud/services/teo/resource_tc_teo_multi_path_gateway.go`
- 新增文件: `tencentcloud/services/teo/resource_tc_teo_multi_path_gateway_test.go`
- 新增文件: `tencentcloud/services/teo/resource_tc_teo_multi_path_gateway.md`
- 修改文件: `tencentcloud/provider.go`（注册新资源）
- 修改文件: `tencentcloud/provider.md`（添加资源文档入口）
- 修改文件: `tencentcloud/services/teo/service_tencentcloud_teo.go`（添加服务层 Describe 方法）
- 依赖云 API: `teo/v20220901` 包中的 CreateMultiPathGateway、DescribeMultiPathGateways、ModifyMultiPathGateway、DeleteMultiPathGateway 接口
