## Why

TEO（边缘安全加速平台）目前缺少对 API 资源（Security API Resource）的 Terraform 管理能力。用户需要通过 Terraform 来管理站点下的 API 资源配置，包括创建、查询、修改和删除 API 资源，以实现基础设施即代码的安全防护管理。

## What Changes

- 新增 Terraform 资源 `tencentcloud_teo_security_api_resource`，管理 TEO 站点下的 API 资源
- 支持通过 CreateSecurityAPIResource 接口创建 API 资源（每次仅1个）
- 支持通过 DescribeSecurityAPIResource 接口查询 API 资源（含分页，按 apiResourceId 匹配）
- 支持通过 ModifySecurityAPIResource 接口修改 API 资源（每次仅1个）
- 支持通过 DeleteSecurityAPIResource 接口删除 API 资源
- 使用复合 ID `zoneId#apiResourceId` 标识每个 API 资源实例
- 在 provider.go 和 provider.md 中注册新资源
- 生成对应的 .md 文档

## Capabilities

### New Capabilities
- `teo-security-api-resource`: 管理 TEO 站点下的单个 API 安全资源配置，支持完整的 CRUD 生命周期管理，使用复合 ID 标识

### Modified Capabilities

## Impact

- 新增文件: `tencentcloud/services/teo/resource_tc_teo_security_api_resource.go`
- 新增文件: `tencentcloud/services/teo/resource_tc_teo_security_api_resource_test.go`
- 新增文件: `tencentcloud/services/teo/resource_tc_teo_security_api_resource.md`
- 修改文件: `tencentcloud/provider.go` (注册新资源)
- 修改文件: `tencentcloud/provider.md` (添加资源文档条目)
- 修改文件: `tencentcloud/services/teo/service_tencentcloud_teo.go` (新增 DescribeTeoSecurityAPIResourceById 服务方法)
- 依赖云 API 包: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901`
