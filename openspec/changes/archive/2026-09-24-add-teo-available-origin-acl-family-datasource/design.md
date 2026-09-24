## Context

EdgeOne (TEO) 源站防护通过回源 ACL 控制域（`OriginACLFamily`）来决定回源 IP 网段的地理范围。控制域分为三档：标准控制域（`gaz`/`mlc`/`emc`）、精简控制域（`plat-gaz`/`plat-mlc`/`plat-emc`）、定制版控制域（`plat-specific-gaz`/`plat-specific-mlc`/`plat-specific-emc`）。

`tencentcloud_teo_origin_acl` 资源已支持设置 `origin_acl_family` 参数，但用户在选值前无法在 Terraform 内查询某站点下可用的控制域及其回源 IP 网段详情。云 API `DescribeAvailableOriginACLFamily` 已提供该查询能力（存在于 vendor 中），其入参为：

- `ZoneId`（string，必填）：站点 ID
- `Filters`（[]*Filter，选填）：过滤条件，仅支持按 `OriginACLFamily` 过滤，`Filters.Values` 上限为 20
- `Offset`（uint64，选填）：分页偏移量，默认 0
- `Limit`（uint64，选填）：分页限制数目，默认 20，最大 100

响应 `DescribeAvailableOriginACLFamilyResponseParams` 包含：
- `TotalCount`（int64）：总数
- `OriginACLFamilyInfos`（[]*OriginACLFamilyInfo）：控制域详情列表，每个元素含 `Version`、`ActiveTime`、`EntireAddresses`（`Addresses` 类型，含 `IPv4`/`IPv6`）、`OriginACLFamily`

代码风格参考：
- 业务逻辑（schema + Read 函数）参考 `data_source_tc_teo_origin_acl.go`、`data_source_tc_igtm_instance_list.go`
- service 层分页查询参考 `service_tencentcloud_teo.go` 中 `DescribeZones` 的分页循环模式

## Goals / Non-Goals

**Goals:**
- 新增数据源 `tencentcloud_teo_available_origin_acl_family`，封装云 API `DescribeAvailableOriginACLFamily`
- 将云 API 的入参都接入数据源 schema：`zone_id`（必填）、`filters`（选填，支持按 `origin_acl_family` 过滤）
- 分页参数 `Offset`/`Limit` 不暴露给用户，在 service 层内部自动分页获取所有数据（遵循 openspec config 约束）
- 将响应列表 `OriginACLFamilyInfos` 展开，把每个元素的所有参数平铺到资源参数 schema 顶层（遵循"资源列表型数据展开"约束），即通过 `origin_acl_family_info_set` Computed 列表字段输出，列表元素包含 `version`、`active_time`、`entire_addresses`（含 `ipv4`/`ipv6`）、`origin_acl_family`
- 在 `provider.go` 和 `provider.md` 中注册该数据源
- 新增数据源文档与单元测试（使用 gomonkey mock 云 API）

**Non-Goals:**
- 不修改 `tencentcloud_teo_origin_acl` 资源或其数据源的现有 schema
- 不暴露分页参数 `Offset`/`Limit` 给用户
- 不为 `origin_acl_family` 过滤值做 ValidateFunc 校验（由云 API 拒绝非法值）
- 不生成 `_extension.go` 文件

## Decisions

### Decision 1: Filters schema 结构 — 复用通用 Filter 结构
**Choice**: `filters` 使用 `TypeList` + `schema.Resource`，子字段为 `name`（string, Required）、`values`（TypeSet of string, Required），与云 API 的 `Filter`（`Name`/`Values`）一一对应。
**Rationale**: 云 API `DescribeAvailableOriginACLFamily` 的 `Filters` 类型是 `[]*Filter`（非 AdvancedFilter，无 Fuzzy 字段）。参照 `data_source_tc_teo_zone_available_plans.go` 等同服务数据源的 filters 设计，`name` 取值 `OriginACLFamily`，`values` 取控制域名称。不使用 AdvancedFilter 以免引入 API 不支持的 `Fuzzy` 字段。

### Decision 2: 分页参数内部处理，不暴露给用户
**Choice**: 不在 schema 中声明 `offset`/`limit`，在 service 层 `DescribeTeoAvailableOriginAclFamilyByFilter` 方法内部以 `Limit=100`（API 最大值）循环分页，累加 `Offset` 直到返回数据量小于 `pageSize`。
**Rationale**: 遵循 openspec config 硬约束"数据源分页:不暴露 limit/offset 参数给用户,内部实现自动分页获取所有数据"。同时遵循"若查询接口中有分页字段，则给定值应该是云API接口注释中标注的最大值"，因此 `Limit` 取最大值 100。

### Decision 3: 响应列表展开为 Computed 列表字段
**Choice**: 使用 `origin_acl_family_info_set`（TypeList, Computed）输出控制域详情列表，每个元素是 `schema.Resource`，包含平铺字段 `version`、`active_time`、`entire_addresses`（TypeList，内含 `ipv4`/`ipv6` 两个 TypeSet）、`origin_acl_family`。同时输出 `total_count`（TypeInt, Computed）顶层字段。
**Rationale**: 遵循"资源列表型数据展开"约束——禁止创建"该资源列表型数据"这一层 schema，应将列表中每个元素的所有参数平铺。但作为数据源，返回的是一个列表，因此使用 Computed 列表字段承载展开后的元素。`result_output_file` 用于保存结果（参考同服务其他数据源）。

### Decision 4: 数据源 Read 错误处理与空响应处理
**Choice**: 在 Read 函数的 retry 块内，若云 API 返回空（`response == nil` 或 `response.Response == nil`），直接返回 `NonRetryableError`，不要 `d.SetId("")`；并在外层 retry 失败路径保留 `log.Printf("[DATASOURCE] read empty, skip SetId")` 提示。
**Rationale**: 遵循 openspec config 中 RESOURCE_KIND_DATASOURCE 资源 Read 方法的约束——避免因云 API 短暂波动导致本地 state 中的 id 被清空。

### Decision 5: 数据源 ID 生成方式
**Choice**: 使用 `helper.BuildToken()` 生成数据源 ID，不使用 `zone_id`。当列表为空时返回 `NonRetryableError` 而非设置空 ID。
**Rationale**: 参照 `data_source_tc_teo_multi_path_gateway_origin_acl.go` 的做法（复合 ID）。但因本数据源可能不返回数据时直接报错更合适，采用 `BuildToken()` 生成临时 ID，与 `data_source_tc_igtm_instance_list.go` 一致。

## Risks / Trade-offs

- **[Risk] 云 API 返回大量精简/定制版控制域 IP 网段** → 内部分页以 `Limit=100` 循环获取所有数据，可完整返回。IP 网段列表可能较长，但作为 Computed 字段不影响 plan 差异（用户不可写）。
- **[Risk] filters name 传非 `OriginACLFamily` 值** → 云 API 仅支持 `OriginACLFamily` 过滤，非法 name 由云端拒绝并返回错误，不在 provider 侧做校验，保持维护性。
- **[Trade-off] 不暴露分页参数** → 用户无法手动控制分页，但遵循 openspec config 统一约束，且自动分页对用户透明。
- **[Trade-off] 不为 origin_acl_family 加 ValidateFunc** → 可用控制域可能随产品迭代增减，依赖云 API 校验更可维护。