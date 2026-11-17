# Monitor Grafana Versions DataSource Specification

## Requirement: Data Source Schema Definition
Data Source `tencentcloud_monitor_grafana_versions` 必须支持以下输入参数和输出属性。

**Input Parameters:**
- `result_output_file` (String, Optional): 输出结果到文件

**Output Attributes:**
- `versions` (List): Grafana 可用版本列表，每个版本包含：
  - `alias` (String): 版本别名
  - `version` (String): 版本号

#### Scenario: Query grafana versions successfully
- **WHEN** user queries the `tencentcloud_monitor_grafana_versions` data source without any required parameters
- **THEN** the provider SHALL call `DescribeGrafanaVersions` API
- **AND** the provider SHALL return the list of available Grafana versions
- **AND** each version SHALL contain `alias` and `version` fields

```hcl
data "tencentcloud_monitor_grafana_versions" "grafana_versions" {
}

output "available_versions" {
  value = data.tencentcloud_monitor_grafana_versions.grafana_versions.versions
}
```

#### Scenario: Output to file
- **WHEN** user specifies `result_output_file` parameter
- **THEN** the version information SHALL be written to the specified file in JSON format

## Requirement: API Integration
Data Source 必须集成 Monitor API v20180724 的 `DescribeGrafanaVersions` 接口。

#### Scenario: API call execution
- **WHEN** the data source is read
- **THEN** it SHALL call `DescribeGrafanaVersions` API through the MonitorService service layer
- **AND** it SHALL handle API rate limiting and retries appropriately using `tccommon.ReadRetryTimeout`
- **AND** it SHALL wrap API errors using `tccommon.RetryError`
- **AND** it SHALL log API calls using `tccommon.LogElapsed` and debug logging

#### Scenario: Empty response handling
- **WHEN** `DescribeGrafanaVersions` returns empty result (`response == nil` or `response.Response == nil` or `len(response.Response.Versions) == 0`)
- **THEN** the provider SHALL return `NonRetryableError` instead of clearing the data source ID
- **AND** the provider SHALL log `[DATASOURCE] read empty, skip SetId` for troubleshooting

## Requirement: Error Handling
Data Source 必须正确处理各种错误情况。

#### Scenario: API transient error retry
- **WHEN** `DescribeGrafanaVersions` API call fails with retryable error
- **THEN** the provider SHALL retry the operation using `resource.Retry` with `tccommon.ReadRetryTimeout`
- **AND** if retry succeeds, the operation SHALL complete normally

#### Scenario: API permission error
- **WHEN** the API call fails due to insufficient permissions
- **THEN** the data source SHALL return the original API error message to help with troubleshooting

## Requirement: Code Structure Compliance
Data Source 实现必须遵循项目代码规范和参考实现模式。

#### Scenario: Code structure follows grafana_plugin_overviews pattern
- **WHEN** reviewing the data source implementation code
- **THEN** the code structure SHALL follow `data_source_tc_monitor_grafana_plugin_overviews.go` patterns
- **AND** function naming SHALL follow convention `dataSourceTencentCloudMonitorGrafanaVersionsRead`
- **AND** SHALL use `tccommon.LogElapsed` defer pattern for operation logging
- **AND** SHALL use `defer tccommon.InconsistentCheck(d, meta)()` for state validation
- **AND** the data source file SHALL be placed in `tencentcloud/services/tcmg/` directory

#### Scenario: Service layer method
- **WHEN** the data source needs to call the cloud API
- **THEN** a service layer method SHALL be added to `tencentcloud/services/monitor/service_tencentcloud_monitor.go`
- **AND** the method SHALL use `me.client.UseMonitorClient().DescribeGrafanaVersions(request)` to call the API
- **AND** the method SHALL follow the naming convention `DescribeMonitorGrafanaVersionsByFilter`

## Requirement: Provider Registration
Data Source 必须在 Provider 中正确注册。

#### Scenario: Register data source in provider
- **WHEN** the provider is initialized
- **THEN** `tencentcloud/provider.go` SHALL register `tencentcloud_monitor_grafana_versions` data source
- **AND** the registration SHALL use `tcmg.DataSourceTencentCloudMonitorGrafanaVersions()`

## Requirement: Unit Test
Data Source 必须提供单元测试，使用 gomonkey mock 云 API。

#### Scenario: Unit test with mock
- **WHEN** running unit tests for the data source
- **THEN** the test SHALL use gomonkey to mock the `DescribeMonitorGrafanaVersionsByFilter` service layer method
- **AND** the test SHALL verify the data source correctly maps API response fields to Terraform schema
- **AND** the test SHALL NOT use terraform acceptance test suite
