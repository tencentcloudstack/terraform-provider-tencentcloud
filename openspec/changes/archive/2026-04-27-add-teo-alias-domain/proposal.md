# Proposal: Add tencentcloud_teo_alias_domain Resource

## What

新增 Terraform resource `tencentcloud_teo_alias_domain`，用于管理腾讯云 TEO（边缘安全加速平台）别称域名（Alias Domain）的完整生命周期。

## Why

当前 terraform-provider-tencentcloud 的 TEO 模块尚未覆盖别称域名的管理能力。别称域名功能允许用户将多个域名（别称）映射到同一个加速域名（目标域名），并支持证书配置和启停管理，是企业版套餐的核心功能之一。用户需要通过 Terraform 实现 IaC 自动化管理别称域名，包括：

- 创建别称域名并配置目标域名及证书
- 查询别称域名状态
- 修改别称域名的目标域名或证书配置
- 修改别称域名的启用/停用状态
- 删除别称域名

## Scope

- 新增 resource 文件：`tencentcloud/services/teo/resource_tc_teo_alias_domain.go`
- 新增 service 方法：在 `service_tencentcloud_teo.go` 中追加 `DescribeTeoAliasDomainById`
- 新增 md 文档：`tencentcloud/services/teo/resource_tc_teo_alias_domain.md`
- 新增单元测试：`tencentcloud/services/teo/resource_tc_teo_alias_domain_test.go`
- 注册资源：在 `tencentcloud/provider.go` 中注册新资源

## API Mapping

| CRUD 操作 | 接口名称 | 文档链接 |
|---|---|---|
| Create | CreateAliasDomain | https://cloud.tencent.com/document/api/1552/81247 |
| Read | DescribeAliasDomains | https://cloud.tencent.com/document/api/1552/81245 |
| Update (config) | ModifyAliasDomain | https://cloud.tencent.com/document/api/1552/81244 |
| Update (status) | ModifyAliasDomainStatus | https://cloud.tencent.com/document/api/1552/81243 |
| Delete | DeleteAliasDomain | https://cloud.tencent.com/document/api/1552/81246 |

## Resource ID

资源唯一 ID 格式：`{zone_id}#{alias_name}`（使用 `tccommon.FILED_SP` 分隔）

## Constraints

- 仅企业版套餐支持
- 功能当前仍在内测阶段
