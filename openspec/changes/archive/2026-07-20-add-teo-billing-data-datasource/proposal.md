## Why

Terraform Provider for TencentCloud 当前不支持查询 EdgeOne (TEO) 的计费数据。用户在以基础设施即代码方式管理 EdgeOne 站点时，无法通过 Terraform 数据源获取流量、带宽、请求数等计费指标数据，导致需要在 Terraform 之外使用额外工具来核对用量。新增 `tencentcloud_teo_billing_data` 数据源可填补这一空白，使用户能在 Terraform 中直接查询计费数据并用于配置编排或对账。

## What Changes

- 新增数据源 `tencentcloud_teo_billing_data`（RESOURCE_KIND_DATASOURCE），通过调用云 API `DescribeBillingData`（teo v20220901）读取计费数据。
- 新增资源代码文件 `tencentcloud/services/teo/data_source_tc_teo_billing_data.go`，实现数据源 Read 逻辑（仅查询，不涉及创建/更新/删除）。
- 新增单元测试文件 `tencentcloud/services/teo/data_source_tc_teo_billing_data_test.go`，使用 gomonkey mock 云 API 进行业务逻辑测试。
- 新增文档样例文件 `tencentcloud/services/teo/data_source_tc_teo_billing_data.md`。
- 在 `tencentcloud/provider.go` 和 `tencentcloud/provider.md` 中注册该数据源。

### Schema 参数说明

入参（均为查询条件，Optional）：
- `start_time` (string, Required): 起始时间。
- `end_time` (string, Required): 结束时间，查询范围需 ≤ 31 天。
- `zone_ids` (list of string, Required): 站点 ID 集合，最多 100 个；用 `*` 表示查询账号级别数据。
- `metric_name` (string, Required): 计费指标名，如 `acc_flux`、`acc_bandwidth` 等。
- `interval` (string, Optional): 查询时间粒度，取值 `5min` / `hour` / `day`。
- `filters` (list of object, Optional): 过滤条件，每项包含 `type` 与 `value`，支持 `host`、`proxy-id`、`region-id`。
- `group_by` (list of string, Optional): 分组聚合维度，取值 `zone-id`、`host`、`proxy-id`、`region-id`，最多两个维度。

出参：
- `data` (list of object, Computed): 计费数据点列表，每项包含 `time`、`value`、`zone_id`、`host`、`proxy_id`、`region_id`。
- `id` (string, Computed): 数据源 ID，由查询参数组合生成，用于 Terraform 状态标识。

## Capabilities

### New Capabilities
- `teo-billing-data`: 查询 EdgeOne 计费数据的数据源能力，封装 `DescribeBillingData` 云 API，支持按时间范围、站点、指标、时间粒度、过滤条件、分组维度查询用量数据。

### Modified Capabilities
<!-- 无 -->

## Impact

- 新增代码：`data_source_tc_teo_billing_data.go`、`data_source_tc_teo_billing_data_test.go`、`data_source_tc_teo_billing_data.md`。
- 修改文件：`tencentcloud/provider.go`（注册数据源）、`tencentcloud/provider.md`（文档注册）。
- 依赖：使用 vendor 目录下已有的 `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901` 包，无需新增第三方依赖。
- 不影响任何现有资源/数据源的 schema 与行为，向后完全兼容。
