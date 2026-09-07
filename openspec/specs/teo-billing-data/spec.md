## ADDED Requirements

### Requirement: 数据源注册与声明

系统 SHALL 在 Terraform Provider 中注册名为 `tencentcloud_teo_billing_data` 的数据源，类型为 RESOURCE_KIND_DATASOURCE，并在 `tencentcloud/provider.go` 中通过 `dataProvider` 工厂注册，同时在 `tencentcloud/provider.md` 中声明对应文档条目。

#### Scenario: Provider 中注册数据源
- **WHEN** 用户在 Terraform 配置中声明 `data "tencentcloud_teo_billing_data" "example"`
- **THEN** Provider 能识别该数据源类型并调用其 Read 函数

#### Scenario: 文档注册
- **WHEN** 执行 `make doc` 文档生成
- **THEN** `provider.md` 中包含 `tencentcloud_teo_billing_data` 数据源条目

### Requirement: 查询计费数据

数据源 SHALL 调用云 API `DescribeBillingData`（teo v20220901）查询 EdgeOne 计费数据，将入参映射到请求字段，将响应 `Data` 列表映射到 `data` 输出参数。

#### Scenario: 带完整参数查询成功
- **WHEN** 用户配置 `start_time`、`end_time`、`zone_ids`、`metric_name`、`interval`、`filters`、`group_by` 并执行 `terraform plan`
- **THEN** 数据源将参数组装为 `DescribeBillingDataRequest` 调用云 API，并将响应中 `Data` 列表每项的 `Time`、`Value`、`ZoneId`、`Host`、`ProxyId`、`RegionId` 映射到 `data` 列表输出字段 `time`、`value`、`zone_id`、`host`、`proxy_id`、`region_id`

#### Scenario: 仅必填参数查询成功
- **WHEN** 用户仅配置 `start_time`、`end_time`、`zone_ids`、`metric_name`
- **THEN** 数据源以 Optional 参数为空调用云 API，仍能正确返回数据并设置 `data` 输出

#### Scenario: 查询返回空数据
- **WHEN** 云 API 返回 `Data` 为 null 或空列表
- **THEN** 数据源将 `data` 设置为空列表，且不报错

### Requirement: 入参 Schema 定义

数据源 SHALL 定义以下入参 schema 字段：
- `start_time`（string, Required）：起始时间。
- `end_time`（string, Required）：结束时间，查询范围 ≤ 31 天。
- `zone_ids`（list of string, Required）：站点 ID 集合，最多 100 个；`*` 表示账号级别。
- `metric_name`（string, Required）：计费指标名。
- `interval`（string, Optional）：查询时间粒度，取值 `5min` / `hour` / `day`。
- `filters`（list of object, Optional）：过滤条件，每项含 `type`（string, Required）、`value`（string, Required）。
- `group_by`（list of string, Optional）：分组聚合维度，取值 `zone-id` / `host` / `proxy-id` / `region-id`，最多两个维度。

#### Scenario: Required 字段缺失
- **WHEN** 用户配置缺少 `start_time` 或 `metric_name` 等必填字段
- **THEN** Terraform 在 plan 阶段报错提示字段缺失

#### Scenario: filters 结构
- **WHEN** 用户配置 `filters` 包含多个 `{ type = "host", value = "test.example.com" }` 项
- **THEN** 数据源将每个项映射为 `BillingDataFilter`（`Type`、`Value`）并传入请求 `Filters` 字段

### Requirement: 出参 Schema 定义

数据源 SHALL 定义以下出参 schema 字段：
- `id`（string, Computed）：数据源 ID，由查询参数组合生成。
- `data`（list of object, Computed）：计费数据点列表，每个 object 含 `time`（string）、`value`（integer）、`zone_id`（string）、`host`（string）、`proxy_id`（string）、`region_id`（string）。

#### Scenario: id 生成
- **WHEN** 数据源 Read 成功执行
- **THEN** `id` 被设置为 `metric_name#start_time#end_time` 格式（使用 `FIELD_SP` 分隔符）

#### Scenario: data 字段展开
- **WHEN** 云 API 返回 `Data` 列表包含多个 `BillingData` 项
- **THEN** `data` 输出列表每项对应一个 `BillingData`，且字段 `Time`/`Value`/`ZoneId`/`Host`/`ProxyId`/`RegionId` 分别映射到 `time`/`value`/`zone_id`/`host`/`proxy_id`/`region_id`

### Requirement: 重试与错误处理

数据源 Read 调用 `DescribeBillingData` 时 SHALL 使用 `tccommon.ReadRetryTimeout` 作为超时时间，通过 `helper.Retry()` 包装；失败时使用 `tccommon.RetryError()` 包装错误返回。retry 块内仅执行 API 调用，不执行 set 等操作。

#### Scenario: 云 API 短暂失败重试
- **WHEN** `DescribeBillingData` 调用返回临时性错误
- **THEN** 数据源在 `ReadRetryTimeout` 内重试调用，直至成功或超时

### Requirement: 空响应保护

数据源 Read 的 retry 块内 SHALL 检查云 API 返回是否为空（`response == nil` 或 `response.Response == nil`），若为空则返回 `NonRetryableError` 而非直接 `d.SetId("")`，并在 retry 失败路径保留 `log.Printf("[DATASOURCE] read empty, skip SetId")` 提示。

#### Scenario: 云 API 返回 nil 响应
- **WHEN** `DescribeBillingData` 返回 `response == nil` 或 `response.Response == nil`
- **THEN** 数据源返回 `NonRetryableError`，不清空 state 中的 id，并在日志中打印 `[DATASOURCE] read empty, skip SetId` 提示

### Requirement: 单元测试

数据源 SHALL 提供单元测试文件 `data_source_tc_teo_billing_data_test.go`，使用 gomonkey mock `DescribeBillingData` 云 API，测试业务逻辑（参数组装、响应解析、字段 set），不使用 Terraform 测试套件，通过 `go test -gcflags=all=-l` 运行。

#### Scenario: mock 云 API 返回数据并验证字段映射
- **WHEN** 运行单元测试，gomonkey mock `DescribeBillingData` 返回包含多个 `BillingData` 项的响应
- **THEN** 测试验证 Read 函数正确组装请求参数并将响应字段设置到 `data` 输出

#### Scenario: mock 云 API 返回空数据
- **WHEN** 运行单元测试，gomonkey mock `DescribeBillingData` 返回 `Data` 为空
- **THEN** 测试验证 Read 函数将 `data` 设置为空列表且不报错
