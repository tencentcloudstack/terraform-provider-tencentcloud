## ADDED Requirements

### Requirement: 查询 EdgeOne 源站防护可用控制域列表

系统 SHALL 提供 Terraform 数据源 `tencentcloud_teo_available_origin_acl_family`，封装云 API `DescribeAvailableOriginACLFamily`（teo v20220901），按站点 ID 查询源站防护可配置控制域（OriginACLFamily）列表信息，包括版本号、生效时间、回源 IP 网段（IPv4/IPv6）与控制域名称。

#### Scenario: 按站点 ID 查询全部可用控制域

- **WHEN** 用户配置数据源仅提供 `zone_id`，未提供 `filters`
- **THEN** 系统 SHALL 调用 `DescribeAvailableOriginACLFamily`（ZoneId=zone_id），内部自动分页（limit=100）拉取全部控制域信息
- **AND** 系统 SHALL 将结果映射到 `origin_acl_family_infos` 列表，每个元素含 `version`、`active_time`、`entire_addresses`（含 `i_pv4`、`i_pv6`）、`origin_acl_family`
- **AND** 系统 SHALL 设置数据源 id（使用 token），不报错

#### Scenario: 按 OriginACLFamily 过滤查询控制域

- **WHEN** 用户配置 `filters`，其中 `name = "OriginACLFamily"`，`values = ["gaz"]`
- **THEN** 系统 SHALL 构造 `[]*teo.Filter{Name: "OriginACLFamily", Values: ["gaz"]}` 作为请求 `Filters`
- **AND** 系统 SHALL 仅返回匹配的控制域信息

#### Scenario: 站点下无可用控制域

- **WHEN** 云 API 返回 `OriginACLFamilyInfos` 为空列表（正常空结果）
- **THEN** 系统 SHALL 不视为错误，`origin_acl_family_infos` 设为空列表
- **AND** 系统 SHALL 正常完成 Read 并设置数据源 id

#### Scenario: 云 API 返回异常空响应

- **WHEN** 云 API 返回 `response == nil` 或 `response.Response == nil`
- **THEN** 服务层 SHALL 返回 `NonRetryableError`，数据源 Read 的外层 retry 继续尝试直至重试耗尽失败
- **AND** 系统 SHALL 不直接清空 state 中的 id

### Requirement: 数据源 schema 与字段映射

数据源 schema SHALL 遵循 RESOURCE_KIND_DATASOURCE 规范，输入参数与输出参数与云 API 字段一一映射。

#### Scenario: 输入参数定义

- **WHEN** 定义数据源 schema 输入参数
- **THEN** `zone_id` SHALL 为 `Required`、`TypeString`，对应 `request.ZoneId`
- **AND** `filters` SHALL 为 `Optional`、`TypeList`，对应 `request.Filters`（`teo.Filter` 类型）
- **AND** `filters.name` SHALL 为 `Required`、`TypeString`，对应 `request.Filters.Name`
- **AND** `filters.values` SHALL 为 `Required`、`TypeSet`（元素 TypeString），对应 `request.Filters.Values`
- **AND** `result_output_file` SHALL 为 `Optional`、`TypeString`，用于保存结果

#### Scenario: 输出参数定义

- **WHEN** 定义数据源 schema 输出参数
- **THEN** `origin_acl_family_infos` SHALL 为 `Computed`、`TypeList`，对应 `response.Response.OriginACLFamilyInfos`
- **AND** 列表中每个元素的 `version` SHALL 为 `Computed`、`TypeString`，对应 `OriginACLFamilyInfo.Version`
- **AND** `active_time` SHALL 为 `Computed`、`TypeString`，对应 `OriginACLFamilyInfo.ActiveTime`
- **AND** `entire_addresses` SHALL 为 `Computed`、`TypeList`（单元素），对应 `OriginACLFamilyInfo.EntireAddresses`
- **AND** `entire_addresses.i_pv4` SHALL 为 `Computed`、`TypeSet`（TypeString），对应 `Addresses.IPv4`
- **AND** `entire_addresses.i_pv6` SHALL 为 `Computed`、`TypeSet`（TypeString），对应 `Addresses.IPv6`
- **AND** `origin_acl_family` SHALL 为 `Computed`、`TypeString`，对应 `OriginACLFamilyInfo.OriginACLFamily`

#### Scenario: nil 字段处理

- **WHEN** 云 API 返回的某个字段为 nil
- **THEN** 系统 SHALL 不对该字段调用 set 方法，跳过该字段设置

### Requirement: 服务层自动分页与重试

服务层方法 `DescribeTeoAvailableOriginACLFamilyByFilter` SHALL 内部自动分页拉取全部数据，不向用户暴露 Offset/Limit 参数，且重试逻辑位于分页循环内部。

#### Scenario: 自动分页拉取全部数据

- **WHEN** 站点下控制域数量超过单页上限（100）
- **THEN** 服务层 SHALL 以 limit=100 循环请求，累加所有页的 `OriginACLFamilyInfos` 直至本页返回数小于 limit 或返回空列表
- **AND** 服务层 SHALL 返回全部累积结果与 totalCount

#### Scenario: 单页调用瞬时错误重试

- **WHEN** 单次 `DescribeAvailableOriginACLFamily` 调用因瞬时错误失败
- **THEN** 服务层 SHALL 通过 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 重试该单次调用
- **AND** 重试发生在分页循环内部，不影响已拉取页的结果