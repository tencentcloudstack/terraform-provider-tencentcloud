## Context

Terraform Provider for TencentCloud 通过 `tencentcloud-sdk-go` 调用云 API 管理腾讯云资源。EdgeOne (TEO) 服务已存在多个数据源（如 `tencentcloud_teo_zones`、`tencentcloud_teo_plans` 等），均位于 `tencentcloud/services/teo/` 目录下，遵循 `data_source_tc_teo_<name>.go` 命名规范。

当前 TEO 数据源覆盖了站点、套餐、配置组、安全模板等查询场景，但尚未覆盖计费数据查询。云 API `DescribeBillingData`（teo v20220901）提供了按时间范围、站点、指标、时间粒度、过滤条件、分组维度查询计费用的能力，返回数据点列表（`BillingData`），每个数据点包含时间戳、数值、站点 ID、域名、四层代理实例 ID、计费大区 ID 等字段。

本变更新增 `tencentcloud_teo_billing_data` 数据源（RESOURCE_KIND_DATASOURCE），封装该 API 的查询能力。该数据源为只读查询，不涉及任何资源创建/更新/删除操作。

### 云 API 结构（来自 vendor 校验）

- 请求 `DescribeBillingDataRequest`：
  - `StartTime *string`（起始时间）
  - `EndTime *string`（结束时间，范围 ≤ 31 天）
  - `ZoneIds []*string`（站点 ID 集合，必填，最多 100 个；`*` 表示账号级别）
  - `MetricName *string`（指标名）
  - `Interval *string`（时间粒度：`5min` / `hour` / `day`）
  - `Filters []*BillingDataFilter`（过滤条件，每项含 `Type`、`Value`）
  - `GroupBy []*string`（分组维度：`zone-id` / `host` / `proxy-id` / `region-id`，最多两个）
- 响应 `DescribeBillingDataResponse`：
  - `Response.Data []*BillingData`（数据点列表，可能返回 null）
  - `BillingData` 结构：`Time *string`、`Value *uint64`、`ZoneId *string`、`Host *string`、`ProxyId *string`、`RegionId *string`

该接口为同步接口（无异步轮询需求），无分页字段。

## Goals / Non-Goals

**Goals:**
- 新增 `tencentcloud_teo_billing_data` 数据源，支持通过 `DescribeBillingData` API 查询 TEO 计费数据。
- 完整映射入参（StartTime/EndTime/ZoneIds/MetricName/Interval/Filters/GroupBy）与出参（Data 列表展开字段）。
- 在 `provider.go` 和 `provider.md` 中注册数据源。
- 提供基于 gomonkey mock 的单元测试（不使用 Terraform 测试套件）。
- 生成 `.md` 文档样例。

**Non-Goals:**
- 不实现计费数据的写入、修改或删除（该 API 仅为查询接口）。
- 不对接其他计费相关 API。
- 不变更任何现有 TEO 数据源/资源的 schema 或行为。
- 不处理异步轮询（本接口为同步接口）。
- 不暴露分页参数（该 API 无分页字段）。

## Decisions

### 决策 1：数据源 ID 生成方式

数据源为纯查询接口，无服务端持久化 ID。采用查询关键参数组合生成确定性 ID，格式为 `metric_name#start_time#end_time`，使用 `tccommon.FIELD_SP`（`#`）作为分隔符。

**理由**：与项目中其他查询型数据源的惯例一致；该组合在单次 `terraform plan/apply` 范围内稳定，满足 Terraform state 标识需求。

**备选方案**：使用 `RequestId`。否决，因为每次查询 RequestId 不同，会导致 state 频繁变更。

### 决策 2：出参 data 列表字段展开

按照代码生成规范要求，`response.Response.Data` 是一个列表，将列表中元素的字段（`time`、`value`、`zone_id`、`host`、`proxy_id`、`region_id`）平铺到 `data` 这个 list 类型的 schema 中，每个元素为一个 object，object 内部包含上述字段。不在 schema 顶层再嵌套一层"列表型数据"包裹结构。

**理由**：遵循项目规范，使每个字段都可被 Terraform 单独 set/read，且与云 API 返回结构直接对应。

### 决策 3：空响应处理

在 Read 方法的 retry 块内，若云 API 返回空（`response == nil` / `response.Response == nil`），不直接 `d.SetId("")`，而是返回 `NonRetryableError`，让外层 retry 继续尝试，并在 retry 失败路径上保留 `log.Printf("[DATASOURCE] read empty, skip SetId")` 提示。

**理由**：遵循 DATASOURCE 资源代码生成规范第 14 条，避免因云 API 短暂波动导致本地 state 中的 id 被清空，造成数据丢失。

### 决策 4：filters 参数结构

`filters` 为 list of object，每个 object 含 `type`（string）和 `value`（string），与云 API `BillingDataFilter` 结构一一对应。设为 Optional。

### 决策 5：retry 与超时

调用 `DescribeBillingData` 时使用 `tccommon.ReadRetryTimeout` 作为超时时间，通过 `helper.Retry()` 包装。失败时使用 `tccommon.RetryError()` 包装错误返回。retry 块内仅执行 API 调用，不执行 set 等操作。

### 决策 6：单元测试方式

使用 gomonkey 对 `DescribeBillingData` 进行 mock，仅测试业务逻辑（参数组装、响应解析、字段 set），不使用 Terraform 测试套件，不依赖真实云环境。使用 `go test -gcflags=all=-l` 运行。

## Risks / Trade-offs

- **[指标名取值范围广]** → `metric_name` 为自由字符串，云 API 会在服务端校验取值。Terraform schema 不做枚举限制，避免与云 API 取值变化不同步；非法值由云 API 返回错误，经 `RetryError` 透传给用户。
- **[ZoneIds 用 `*` 表示账号级别]** → 保留原样透传，文档中说明该特殊用法。
- **[Data 可能为 null]** → Read 中对 `Data == nil` 进行判空处理，set 空列表而非报错，保证查询无数据时 Terraform 不报错。
- **[查询时间范围限制 31 天]** → 由云 API 服务端校验，schema 层不强加校验，错误经 retry 包装透传。
