## ADDED Requirements

### Requirement: Data source registration
数据源 MUST 在 Provider 中注册为 `tencentcloud_dbdc_db_custom_zones`，使其可在 Terraform 配置中使用。

#### Scenario: Data source is accessible
- **WHEN** 用户在 Terraform 配置中使用 `data "tencentcloud_dbdc_db_custom_zones"`
- **THEN** Terraform 能够成功识别并初始化该数据源

#### Scenario: Provider registration entry
- **WHEN** 检查 `tencentcloud/provider.go` 的 DataSourcesMap
- **THEN** 存在键 `tencentcloud_dbdc_db_custom_zones` 指向 `dbdc.DataSourceTencentCloudDbdcDbCustomZones()`

#### Scenario: Provider doc entry
- **WHEN** 检查 `tencentcloud/provider.md` 的数据源列表
- **THEN** 存在 `tencentcloud_dbdc_db_custom_zones` 条目

### Requirement: Query zones without input parameters
数据源 MUST 支持不带任何输入参数的查询，返回当前地域 DBDC 支持的可用区列表。

#### Scenario: Query all zones
- **WHEN** 用户未设置任何过滤参数（`data "tencentcloud_dbdc_db_custom_zones" "example" {}`）
- **THEN** 数据源调用 `DescribeDBCustomZones` API 且不传递任何请求参数
- **THEN** 返回当前地域全部可用区数据

### Requirement: Return zone set list
数据源 MUST 返回 `zone_set` 列表，列表中每个元素对应一个可用区，并将可用区的字段平铺到元素中（不再额外嵌套一层列表 schema）。

#### Scenario: Zone set is populated
- **WHEN** `DescribeDBCustomZones` API 返回非空的 `ZoneSet`
- **THEN** 输出属性 `zone_set` 列表长度等于返回的 `ZoneSet` 元素数量
- **THEN** 列表中每个元素包含 `zone` 与 `zone_state` 字段

#### Scenario: Zone set is empty
- **WHEN** `DescribeDBCustomZones` API 返回空的 `ZoneSet`
- **THEN** `zone_set` 为空列表，数据源不报错

### Requirement: Return zone field
数据源 MUST 在 `zone_set` 的每个元素中返回 `zone` 字段，对应云 API `ZoneInfo.Zone`。

#### Scenario: Zone value is present
- **WHEN** API 返回的 `ZoneInfo.Zone` 非空
- **THEN** 元素的 `zone` 字段包含可用区标识（如 `ap-guangzhou-3`）

#### Scenario: Zone value is nil
- **WHEN** API 返回的 `ZoneInfo.Zone` 为 nil
- **THEN** 该元素的 `zone` 字段不被设置（先判断 nil 再 set）

### Requirement: Return zone state field
数据源 MUST 在 `zone_set` 的每个元素中返回 `zone_state` 字段，对应云 API `ZoneInfo.ZoneState`。

#### Scenario: Zone state SELL
- **WHEN** API 返回的 `ZoneInfo.ZoneState` 为 `SELL`
- **THEN** 元素的 `zone_state` 字段值为 `SELL`，表示正常售卖

#### Scenario: Zone state SOLD_OUT
- **WHEN** API 返回的 `ZoneInfo.ZoneState` 为 `SOLD_OUT`
- **THEN** 元素的 `zone_state` 字段值为 `SOLD_OUT`，表示售罄

#### Scenario: Zone state is nil
- **WHEN** API 返回的 `ZoneInfo.ZoneState` 为 nil
- **THEN** 该元素的 `zone_state` 字段不被设置（先判断 nil 再 set）

### Requirement: Support result output file
数据源 MUST 支持通过 `result_output_file` 参数将查询结果保存到文件。

#### Scenario: Save to file
- **WHEN** 用户设置 `result_output_file = "./zones.json"`
- **THEN** 查询结果以 JSON 格式保存到指定文件

#### Scenario: File path is optional
- **WHEN** 用户未设置 `result_output_file`
- **THEN** 不保存文件，仅返回数据到 Terraform state

### Requirement: Handle API empty response without clearing state id
数据源 MUST 在 `DescribeDBCustomZones` API 返回空（`response == nil` / `response.Response == nil` / `ZoneSet == nil`）时，不直接调用 `d.SetId("")`，而是返回 `NonRetryableError` 让外层 retry 继续尝试。

#### Scenario: API returns nil response
- **WHEN** `DescribeDBCustomZones` 返回 `result == nil` 或 `result.Response == nil` 或 `result.Response.ZoneSet == nil`
- **THEN** service 方法打印 `log.Printf("[DATASOURCE] read empty, skip SetId")`
- **THEN** service 方法返回 `resource.NonRetryableError`
- **THEN** 数据源 Read 方法不调用 `d.SetId("")`，而是将错误向上返回由 retry 处理

### Requirement: Retry on API failure
数据源 MUST 使用 `tccommon.ReadRetryTimeout` 作为超时时间对 `DescribeDBCustomZones` 调用进行 retry，API 调用失败时使用 `tccommon.RetryError()` 包装错误。

#### Scenario: API call fails and retries
- **WHEN** `DescribeDBCustomZones` API 调用返回错误
- **THEN** 错误被 `tccommon.RetryError()` 包装
- **THEN** 在 `tccommon.ReadRetryTimeout` 超时范围内重试调用

#### Scenario: Retry exhausted
- **WHEN** retry 重试耗尽仍失败
- **THEN** 数据源 Read 方法返回最终错误，不清空 state id

### Requirement: Resource id generation
数据源 MUST 在查询成功后通过 `helper.BuildToken()` 生成并设置资源 id，作为数据源的标识。

#### Scenario: Set id after successful read
- **WHEN** `DescribeDBCustomZones` 调用成功并完成 `zone_set` 设置
- **THEN** 调用 `d.SetId(helper.BuildToken())` 设置数据源 id

### Requirement: Service layer method
service 层 MUST 在 `service_tencentcloud_dbdc.go` 中新增 `DescribeDBCustomZonesByFilter` 方法封装 `DescribeDBCustomZones` API 调用。

#### Scenario: Service method signature
- **WHEN** 调用 `DbdcService.DescribeDBCustomZonesByFilter(ctx, param)`
- **THEN** 方法返回 `([]*dbdcv20201029.ZoneInfo, int64, error)`
- **THEN** 方法内部使用 `dbdcv20201029.NewDescribeDBCustomZonesRequest()` 构造请求
- **THEN** 方法内部使用 `me.client.UseDbdcV20201029Client().DescribeDBCustomZones(request)` 发起调用
