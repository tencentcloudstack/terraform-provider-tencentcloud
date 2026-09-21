## Why

腾讯云 TEO（EdgeOne）提供了 `DescribeInferenceServiceMonitorData` API 用于查询推理服务的监控数据（CPU/GPU 使用率、实例数量、内存使用率等指标）。当前 Provider 缺少对应的数据源，用户无法在 Terraform 中查询推理服务的监控数据，只能手动调用 API 或控制台查询，影响基础设施即代码场景下的运维监控自动化能力。

## What Changes

- 新增 Data Source: `tencentcloud_teo_inference_service_monitor_data`，对应腾讯云 TEO 服务的 `DescribeInferenceServiceMonitorData` API
- 支持按站点 ID、推理服务 ID 列表、指标列表、时间范围查询推理服务监控数据，可选时间粒度
- 返回推理服务监控记录列表，每条记录包含服务 ID、指标名称及监控数据明细（时间点、数值）

## Capabilities

### New Capabilities
- `teo-inference-service-monitor-data-datasource`: 实现查询 TEO 推理服务监控数据的数据源，支持按站点、服务、指标、时间范围查询监控记录

### Modified Capabilities
<!-- 无现有能力需要修改 -->

## Impact

- **新增文件**:
  - `tencentcloud/services/teo/data_source_tc_teo_inference_service_monitor_data.go`（数据源实现）
  - `tencentcloud/services/teo/data_source_tc_teo_inference_service_monitor_data_test.go`（单元测试，使用 mock）
  - `tencentcloud/services/teo/data_source_tc_teo_inference_service_monitor_data.md`（文档样例）
- **修改文件**:
  - `tencentcloud/provider.go`（在 DataSourcesMap 注册新数据源）
- **依赖**: 使用 `tencentcloud-sdk-go/tencentcloud/teo/v20220901` 包中的 `DescribeInferenceServiceMonitorData` 接口（vendor 已支持）
- **兼容性**: 纯新增功能，不影响现有资源和数据源，完全向后兼容