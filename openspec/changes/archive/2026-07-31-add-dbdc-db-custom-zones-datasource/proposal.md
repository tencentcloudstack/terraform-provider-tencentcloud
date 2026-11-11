## Why

用户需要在 Terraform 中查询腾讯云 DBDC（数据库专用集群）指定地域支持的可用区列表及其售卖状态，以便在编排 DBDC 资源（如 `tencentcloud_dbdc_db_custom_cluster`、`tencentcloud_dbdc_db_custom_node`）时选择可用区、判断可用区是否可正常售卖（SELL/SOLD_OUT）。目前 Provider 已有多个 DBDC 数据源（clusters、nodes、cluster_nodes、images），但缺少可用区查询能力，用户无法在同一 Terraform 体系中获知可用区售卖状态，可能导致选用已售罄可用区而创建失败。

## What Changes

- 新增数据源 `tencentcloud_dbdc_db_custom_zones`，调用 DBDC `DescribeDBCustomZones` API 查询指定地域的可用区列表
- 返回 `zone_set` 列表，每个元素包含 `zone`（可用区）和 `zone_state`（可用区状态：SELL 正常售卖 / SOLD_OUT 售罄）
- 支持通过 `result_output_file` 参数将查询结果保存到文件
- 在 `tencentcloud/provider.go` 与 `tencentcloud/provider.md` 中注册该数据源

## Capabilities

### New Capabilities
- `dbdc-db-custom-zones-datasource`: 查询 DBDC 指定地域可用区列表及售卖状态的只读数据源

### Modified Capabilities
<!-- 无现有功能需要修改 -->

## Impact

**新增文件:**
- `tencentcloud/services/dbdc/data_source_tc_dbdc_db_custom_zones.go` - 数据源实现
- `tencentcloud/services/dbdc/data_source_tc_dbdc_db_custom_zones_test.go` - 单元测试（使用 gomonkey mock 云 API）
- `tencentcloud/services/dbdc/data_source_tc_dbdc_db_custom_zones.md` - 文档模板

**修改文件:**
- `tencentcloud/services/dbdc/service_tencentcloud_dbdc.go` - 新增 `DescribeDBCustomZonesByFilter` service 方法
- `tencentcloud/provider.go` - 在 DataSourcesMap 中注册 `tencentcloud_dbdc_db_custom_zones`
- `tencentcloud/provider.md` - 在数据源列表中添加 `tencentcloud_dbdc_db_custom_zones`

**依赖:**
- 已有的 tencentcloud-sdk-go dbdc v20201029 包（`DescribeDBCustomZones` 接口已存在于 vendor 中）
- 无需新增外部依赖

**API 映射:**
- `DescribeDBCustomZones`（无入参）→ `response.Response.ZoneSet`（`[]*ZoneInfo`）映射为 Terraform `zone_set`
