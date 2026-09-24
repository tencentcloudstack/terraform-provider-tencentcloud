## Why

EdgeOne (TEO) 源站防护的回源 ACL 控制域存在多种可用域（标准控制域 gaz/mlc/emc、精简控制域 plat-gaz/plat-mlc/plat-emc、定制版控制域 plat-specific-gaz/plat-specific-mlc/plat-specific-emc）。用户在配置 `tencentcloud_teo_origin_acl` 资源的 `origin_acl_family` 参数时，需要先查询所在站点支持哪些可用控制域及其对应的回源 IP 网段详情，以便选择合适的控制域。目前 provider 没有提供查询该信息的数据源，用户只能通过云 API 控制台或 SDK 手动查询，体验较差。

云 API `DescribeAvailableOriginACLFamily` 已支持查询站点下所有可用控制域及其 IP 网段信息，但 provider 尚未接入。本变更新增数据源 `tencentcloud_teo_available_origin_acl_family`，将该云 API 的入参都接入，方便用户在 Terraform 中直接查询可用控制域信息。

## What Changes

- 新增数据源 `tencentcloud_teo_available_origin_acl_family`，封装云 API `DescribeAvailableOriginACLFamily`
- 将云 API 的入参都接入数据源 schema：
  - `zone_id`（必填 string）：站点 ID
  - `filters`（选填）：过滤条件，支持按 `OriginACLFamily`（控制域）进行过滤
  - 分页参数 `Offset`/`Limit` 不暴露给用户，由内部自动分页获取所有数据（遵循 openspec config 约束）
- 接入响应字段：`TotalCount` 以及 `OriginACLFamilyInfos` 列表，列表元素平铺到顶层 schema（遵循"资源列表型数据展开"约束），每个元素包含：`version`、`active_time`、`entire_addresses`（含 `ipv4`/`ipv6`）、`origin_acl_family`
- 在 `tencentcloud/provider.go` 和 `tencentcloud/provider.md` 中注册该数据源
- 新增数据源文档 `data_source_tc_teo_available_origin_acl_family.md`（用于 `make doc` 生成）
- 新增数据源测试 `data_source_tc_teo_available_origin_acl_family_test.go`（使用 gomonkey mock 云 API）

## Capabilities

### New Capabilities
- `teo-available-origin-acl-family-datasource`: 查询 TEO 源站防护可用的回源 ACL 控制域详细信息的数据源，包含版本号、生效时间、回源 IP 网段和控制域名称，支持按控制域过滤。

### Modified Capabilities
<!-- 无需修改现有 capability 的 requirement -->

## Impact

- **新增代码文件**:
  - `tencentcloud/services/teo/data_source_tc_teo_available_origin_acl_family.go`
  - `tencentcloud/services/teo/data_source_tc_teo_available_origin_acl_family_test.go`
  - `tencentcloud/services/teo/data_source_tc_teo_available_origin_acl_family.md`
  - `tencentcloud/services/teo/service_tencentcloud_teo.go`（新增 `DescribeTeoAvailableOriginAclFamilyByFilter` 方法）
- **修改代码文件**:
  - `tencentcloud/provider.go`（注册数据源）
  - `tencentcloud/provider.md`（注册数据源文档索引）
- **云 API 依赖**:`DescribeAvailableOriginACLFamily`（已存在于 vendor 中）
- **向后兼容**:纯新增数据源，不影响任何现有资源/数据源/schema/state