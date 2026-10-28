## Why

用户需要通过 Terraform 查询腾讯云监控（Monitor）产品下 Grafana 的可用版本列表。当前 Provider 已经支持 Grafana 实例的创建和管理（`tencentcloud_monitor_grafana_instance`），以及 Grafana 版本升级资源（`tencentcloud_monitor_grafana_version_upgrade`），但缺少查询 Grafana 可用版本列表的数据源。

这会导致用户无法：
1. 通过 Terraform 查询 Grafana 支持的可用版本列表
2. 在创建 Grafana 实例或升级版本前获取可用版本信息，用于版本选择决策
3. 在 Terraform 配置中基于版本信息进行条件判断和资源配置

## What Changes

- 新增 Data Source: `tencentcloud_monitor_grafana_versions`
- 实现对 Monitor API `DescribeGrafanaVersions` 接口的调用（包名: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724`）
- 支持查询 Grafana 可用版本列表：
  - `result_output_file`: 输出结果到文件（可选）
- 返回版本信息：
  - `versions`: Grafana 可用版本列表，每个版本包含：
    - `alias`: 版本别名
    - `version`: 版本号

## Capabilities

### New Capabilities
- `monitor-grafana-versions-datasource`: 查询腾讯云 Grafana 可用版本列表的数据源能力

### Modified Capabilities
（无）

## Impact

- **新增能力**: Grafana 可用版本列表查询数据源
- **受影响的服务**: Monitor / TCMG (tencentcloud/services/tcmg)
- **新增文件**:
  - `tencentcloud/services/tcmg/data_source_tc_monitor_grafana_versions.go`
  - `tencentcloud/services/tcmg/data_source_tc_monitor_grafana_versions.md`
  - `tencentcloud/services/tcmg/data_source_tc_monitor_grafana_versions_test.go`
  - `tencentcloud/services/monitor/service_tencentcloud_monitor.go`（新增 service 层方法）
  - Provider 注册代码需要添加此 data source（`tencentcloud/provider.go`）
- **API 依赖**:
  - Monitor API v20180724: `DescribeGrafanaVersions`
  - 文档: https://cloud.tencent.com/document/api/248/73027
- **兼容性**: 无破坏性变更，纯新增功能
