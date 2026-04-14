## Why

用户需要通过 Terraform 管理腾讯云 EdgeOne（TEO）产品的别称域名资源，但目前 terraform-provider-tencentcloud 缺少对 TEO 别称域名的支持。用户必须手动在控制台或使用云 API 管理别称域名，无法实现基础设施即代码（IaC）的自动化管理。

## What Changes

为 teo 服务新增以下功能：
- 新增 `resource_tc_teo_alias_domain` 资源，支持创建、查询、修改、删除别称域名
- 支持 zone_id、alias_name、target_name 三个必填参数
- 支持 paused 状态管理（启用/暂停别称域名）
- 对 CreateAliasDomain、ModifyAliasDomain、ModifyAliasDomainStatus、DeleteAliasDomain 等异步操作提供轮询等待机制
- 生成对应的单元测试文件，使用 mock 云 API 方式测试代码逻辑

## Capabilities

### New Capabilities
- `teo-alias-domain`: 管理 TEO 别称域名的创建、查询、修改、删除和状态管理功能，支持 zone_id、alias_name、target_name 基础参数以及 paused 状态控制

### Modified Capabilities
- 无

## Impact

**新增文件**：
- `tencentcloud/services/teo/resource_tc_teo_alias_domain.go` - 资源实现文件
- `tencentcloud/services/teo/resource_tc_teo_alias_domain_test.go` - 单元测试文件
- `website/docs/r/teo_alias_domain.html.markdown` - 资源文档

**影响代码**：
- `tencentcloud/services/teo/service_tencentcloud_teo.go` - 注册新资源到 provider

**云 API 依赖**：
- `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901` 包中的以下接口：
  - CreateAliasDomain
  - DescribeAliasDomains
  - ModifyAliasDomain
  - ModifyAliasDomainStatus
  - DeleteAliasDomain
