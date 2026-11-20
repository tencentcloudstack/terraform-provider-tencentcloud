## Why

EdgeOne (TEO) 源站防护功能依赖"回源 IP 网段控制域"（OriginACLFamily）来下发 ACL 规则。控制域包含标准控制域（gaz/mlc/emc）和精简控制域（plat-gaz/plat-mlc/plat-emc）等，不同控制域提供的回源 IP 网段数量与使用限制不同。当前 Terraform Provider 缺少查询可用控制域列表的数据源，用户在配置 `tencentcloud_teo_origin_acl` 资源或相关多路径网关资源时，无法通过 Terraform 自动获取可选的控制域及其版本、回源 IP 网段信息，只能到控制台手工查询，难以实现自动化编排。新增 `tencentcloud_teo_available_origin_acl_family` 数据源可填补这一空白。

## What Changes

- 新增数据源 `tencentcloud_teo_available_origin_acl_family`，封装云 API `DescribeAvailableOriginACLFamily`，按站点 ID 查询源站防护可配置控制域列表信息。
- 新增文件 `tencentcloud/services/teo/data_source_tc_teo_available_origin_acl_family.go`，定义数据源 schema 与 Read 函数。
- 在 `tencentcloud/services/teo/service_tencentcloud_teo.go` 中新增服务层方法 `DescribeTeoAvailableOriginACLFamilyByFilter`，支持分页与重试。
- 新增文件 `tencentcloud/services/teo/data_source_tc_teo_available_origin_acl_family.md`，提供一句话描述与示例用法（Argument Reference / Attribute Reference 由 `make doc` 自动生成）。
- 在 `tencentcloud/provider.go` 的 `DataSourcesMap` 中注册 `tencentcloud_teo_available_origin_acl_family`。
- 在 `tencentcloud/provider.md` 的数据源列表中追加 `tencentcloud_teo_available_origin_acl_family`。
- 新增单元测试文件 `tencentcloud/services/teo/data_source_tc_teo_available_origin_acl_family_test.go`，使用 gomonkey mock 云 API 进行业务逻辑测试。

入参：
- `zone_id`（必填，string）：站点 ID，对应 `request.ZoneId`。
- `filters`（可选，list）：过滤条件，对应 `request.Filters`（类型 `teo.Filter`，含 `Name`/`Values`）。其中 `name`（必填）、`values`（必填）。

出参（平铺到顶层，无外层包裹的 list/set 字段）：
- `origin_acl_family_infos`（computed，list）：控制域信息列表，对应 `response.Response.OriginACLFamilyInfos`。
  - `version`（string）：源站防护版本号。
  - `active_time`（string）：版本生效时间（ISO 8601，UTC+8）。
  - `entire_addresses`（list）：回源 IP 网段详情，对应 `OriginACLFamilyInfo.EntireAddresses`（类型 `Addresses`）。
    - `i_pv4`（set）：IPv4 网段列表。
    - `i_pv6`（set）：IPv6 网段列表。
  - `origin_acl_family`（string）：源站防护回源 ACL 控制域。
- `result_output_file`（可选，string）：用于保存结果。

## Capabilities

### New Capabilities

- `teo-available-origin-acl-family-datasource`: 查询 EdgeOne 源站防护可配置控制域列表的数据源能力，封装 `DescribeAvailableOriginACLFamily` 接口。

### Modified Capabilities

（无）

## Impact

- **新增代码**：`tencentcloud/services/teo/data_source_tc_teo_available_origin_acl_family.go`、`data_source_tc_teo_available_origin_acl_family_test.go`、`data_source_tc_teo_available_origin_acl_family.md`。
- **修改代码**：`tencentcloud/services/teo/service_tencentcloud_teo.go`（新增服务方法）、`tencentcloud/provider.go`（注册数据源）、`tencentcloud/provider.md`（追加数据源条目）。
- **依赖**：使用 vendor 中已有的 `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901`，无新增外部依赖。
- **API**：`DescribeAvailableOriginACLFamily`（teo v20220901），支持分页（Offset/Limit，Limit 最大值 100），需在服务层内部自动分页拉取全部数据。
- **兼容性**：纯新增数据源，不影响任何现有资源与数据源，完全向后兼容。