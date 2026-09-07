## Context

腾讯云监控（Monitor）产品提供了 Grafana 托管服务（TCMG），用户可以通过云 API 创建和管理 Grafana 实例。当前 Terraform Provider 已支持 Grafana 实例的完整生命周期管理以及版本升级资源（`tencentcloud_monitor_grafana_version_upgrade`），但缺少查询 Grafana 可用版本列表的数据源。

云 API `DescribeGrafanaVersions`（Monitor v20180724）是一个无入参的只读接口，返回当前可用的 Grafana 版本列表（`[]*GrafanaVersion`），每个版本包含 `Alias`（版本别名）和 `Version`（版本号）两个字段。

现有的同类数据源 `tencentcloud_monitor_grafana_plugin_overviews` 已经采用了相同的架构模式：数据源文件放在 `tencentcloud/services/tcmg` 目录下，service 层方法放在 `tencentcloud/services/monitor/service_tencentcloud_monitor.go` 中，本次新增数据源将遵循完全一致的模式。

## Goals / Non-Goals

**Goals:**
- 新增 `tencentcloud_monitor_grafana_versions` 数据源，支持查询 Grafana 可用版本列表
- 遵循现有 `tencentcloud_monitor_grafana_plugin_overviews` 数据源的代码架构模式
- 正确映射云 API `DescribeGrafanaVersions` 响应字段到 Terraform schema
- 在 Provider 中注册新数据源
- 提供完整的单元测试（使用 gomonkey mock 云 API）

**Non-Goals:**
- 不实现 Grafana 版本的创建、修改、删除操作（版本为平台提供，仅支持查询）
- 不支持按条件过滤版本（云 API 本身无入参，不支持过滤）
- 不修改现有 Grafana 相关资源的任何行为

## Decisions

### 决策1: 数据源文件放置在 `tencentcloud/services/tcmg` 目录
**理由**: 与现有同类数据源 `tencentcloud_monitor_grafana_plugin_overviews` 保持一致，该目录已包含所有 Grafana 相关的数据源和资源文件。service 层方法仍放在 `tencentcloud/services/monitor/service_tencentcloud_monitor.go` 中，通过 `svcmonitor.NewMonitorService` 调用。

### 决策2: Schema 设计采用列表展开模式，不额外嵌套一层
**理由**: 按照项目规范，资源列表型数据应展开到 schema 顶层。`versions` 作为 `TypeList` 顶层字段，每个元素是一个 `schema.Resource`，包含 `alias` 和 `version` 两个 computed 字段。这与 `tencentcloud_monitor_grafana_plugin_overviews` 的 `plugin_set` 结构一致。

### 决策3: 数据源无输入参数（除 `result_output_file` 外）
**理由**: 云 API `DescribeGrafanaVersions` 的 Request 为空结构体，不接受任何入参，因此数据源仅需提供 `result_output_file` 可选参数用于结果导出。

### 决策4: retry 块内检查空响应返回 NonRetryableError
**理由**: 按照项目规范要求，数据源 Read 方法的 retry 块内若云 API 返回空（`response == nil` / `response.Response == nil` / `len(response.Response.Versions) == 0`），不应直接 `d.SetId("")`，而应直接返回 `NonRetryableError`，避免因云 API 短暂波动导致本地 state 中的 id 被清空。

### 决策5: 数据源 ID 使用 `helper.DataResourceIdsHash`
**理由**: 与 `tencentcloud_monitor_grafana_plugin_overviews` 一致，使用所有版本号的 hash 作为数据源 ID，确保数据源在 Terraform 中有唯一标识。

### 决策6: 单元测试使用 gomonkey mock 云 API
**理由**: 按照项目规范，新增 terraform 资源（包括数据源）的单元测试不使用 terraform 测试套件，而是使用 mock（gomonkey）方法对云 API 进行 mock 处理，只进行业务代码逻辑的单元测试。

## Risks / Trade-offs

- **[风险] 云 API 返回空版本列表**: 若平台暂时无可用版本，retry 块将返回 NonRetryableError 导致数据源读取失败。→ 缓解：这是预期行为，平台正常情况下应有可用版本，且规范要求不直接清空 id。
- **[风险] 版本列表变化导致数据源 ID 变化**: 使用版本号 hash 作为 ID，版本列表变化时 ID 会变化。→ 缓解：数据源每次 read 都会重新计算 ID，不影响功能正确性，仅影响 plan 显示。
