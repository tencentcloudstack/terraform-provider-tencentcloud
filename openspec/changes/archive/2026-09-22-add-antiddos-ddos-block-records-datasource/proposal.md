## Why

Terraform Provider for TencentCloud 当前缺少查询 AntiDDoS（DDoS 防护）封堵解封记录的数据源。用户无法在 Terraform 中查询某个时间段内资源的 DDoS 封堵状态与解封配额信息，无法在自动化编排中引用这些数据（如根据封堵状态触发解封流程、审计封堵记录、生成配额用量报表）。新增 `tencentcloud_antiddos_ddos_block_records` 数据源可填补这一空白，对应云 API `DescribeDDoSBlockRecords`（antiddos v20250903 版本）。

## What Changes

- 新增数据源 `tencentcloud_antiddos_ddos_block_records`（RESOURCE_KIND_DATASOURCE），封装 AntiDDoS `DescribeDDoSBlockRecords` 接口，用于按时间范围和过滤条件查询封堵解封记录列表及解封配额信息。
- 新增文件 `tencentcloud/services/antiddos/data_source_tc_antiddos_ddos_block_records.go`：数据源定义与 Read 函数。
- 修改文件 `tencentcloud/services/antiddos/service_tencentcloud_antiddos.go`：新增 service 层方法 `DescribeAntiddosDDoSBlockRecordsByFilter`，内部实现自动分页并（关键约束）将 retry 放在分页循环内部。
- 修改文件 `tencentcloud/connectivity/client.go`：新增 `UseAntiddosV20250903Client()` 访问器方法及对应 `antiddosV20250903Conn` 字段，因为现有 `UseAntiddosClient()` 仅返回 v20200309 的 client，而 `DescribeDDoSBlockRecords` 接口位于 v20250903 包中。
- 修改文件 `tencentcloud/provider.go`：在 DataSourcesMap 中注册 `tencentcloud_antiddos_ddos_block_records`。
- 新增文档文件 `tencentcloud/services/antiddos/data_source_tc_antiddos_ddos_block_records.md`（用于 `make doc` 生成 website 文档）。
- 新增单元测试文件 `tencentcloud/services/antiddos/data_source_tc_antiddos_ddos_block_records_test.go`，使用 gomonkey mock 云 API 进行业务逻辑测试。

## Capabilities

### New Capabilities
- `antiddos-ddos-block-records-datasource`: 通过 AntiDDoS `DescribeDDoSBlockRecords` 接口查询 DDoS 封堵解封记录列表及解封配额信息的数据源能力，支持按时间范围和过滤条件查询，内部自动分页。

### Modified Capabilities

无。本次为纯新增数据源，不修改现有资源的 spec 行为。

## Impact

- **新增代码**：antiddos 服务目录下新增数据源文件、service 方法、文档与测试文件。
- **修改代码**：
  - `tencentcloud/connectivity/client.go`：新增 v20250903 antiddos client 访问器（需新增 import 与字段、方法）。
  - `tencentcloud/provider.go`：DataSourcesMap 注册一行。
  - `tencentcloud/services/antiddos/service_tencentcloud_antiddos.go`：追加 service 方法。
- **SDK 依赖**：使用 vendor 中已有的 `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/antiddos/v20250903` 包（已 vendor，无需新增依赖，但需在 connectivity 层接入）。
- **向后兼容**：纯新增数据源，不影响现有资源与 state，完全向后兼容。
- **文档**：新增数据源文档，由 `make doc` 生成到 `website/docs/d/`。