## Why

TEO (EdgeOne) 用户需要通过 Terraform 管理别称域名（Alias Domain）资源的完整生命周期。当前 Provider 缺少对 TEO 别称域名的支持，用户只能通过控制台或 API 手动创建、修改和删除别称域名，无法实现基础设施即代码的自动化管理。

## What Changes

- 新增 `tencentcloud_teo_alias_domain` 资源（RESOURCE_KIND_GENERAL），支持别称域名的完整 CRUD 操作
  - Create：调用 `CreateAliasDomain` API 创建别称域名
  - Read：调用 `DescribeAliasDomains` API 查询别称域名详情
  - Update：调用 `ModifyAliasDomain` API 修改别称域名配置
  - Delete：调用 `DeleteAliasDomain` API 删除别称域名
- 在 `provider.go` 中注册新资源
- 在 `service_tencentcloud_teo.go` 中添加服务层查询方法
- 添加单元测试和资源文档

## Capabilities

### New Capabilities
- `teo-alias-domain-resource`: 管理 TEO 别称域名资源的完整生命周期，包括创建、读取、更新、删除和导入操作

### Modified Capabilities

（无已有规格需要修改）

## Impact

- **新增文件**: `tencentcloud/services/teo/resource_tc_teo_alias_domain.go`（资源实现）、`resource_tc_teo_alias_domain_test.go`（测试）、`resource_tc_teo_alias_domain.md`（文档）
- **修改文件**: `tencentcloud/provider.go`（注册资源）、`tencentcloud/services/teo/service_tencentcloud_teo.go`（服务层方法）
- **API 依赖**: `teo/v20220901` 包中的 CreateAliasDomain、DescribeAliasDomains、ModifyAliasDomain、DeleteAliasDomain 接口
- **复合 ID**: `zone_id#alias_name`，使用 `tccommon.FILED_SP` 分隔符
- **向后兼容**: 纯新增功能，不影响现有资源
