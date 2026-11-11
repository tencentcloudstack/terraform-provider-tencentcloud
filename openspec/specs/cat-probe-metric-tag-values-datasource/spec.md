# cat-probe-metric-tag-values-datasource Specification

## Purpose
TBD - created by archiving change add-cat-probe-metric-tag-values-datasource. Update Purpose after archive.
## Requirements
### Requirement: Data Source Schema Definition

The data source `tencentcloud_cat_probe_metric_tag_values` SHALL support the following input parameters and output attributes:

**Input Parameters:**
- `analyze_task_type` (String, Optional): 分析任务类型，支持 `AnalyzeTaskType_Network`（网络质量）、`AnalyzeTaskType_Browse`（页面性能）、`AnalyzeTaskType_Transport`（端口性能）、`AnalyzeTaskType_UploadDownload`（文件传输）、`AnalyzeTaskType_MediaStream`（音视频体验）
- `key` (String, Optional): 维度标签值，支持 `host`（任务域名）、`errorInfo`（状态类型）、`area`（拨测点地区）、`operator`（拨测点运营商）、`taskId`（任务ID）
- `filter` (String, Optional): 过滤条件，可以传单个过滤条件也可以拼接多个参数，支持正则匹配
- `filters` (Set of String, Optional): 过滤条件数组
- `time_range` (String, Optional): 时间范围
- `result_output_file` (String, Optional): 用于保存结果

**Output Attributes:**
- `tag_value_set` (String): 标签值序列化后的字符串

#### Scenario: Query tag values with all optional parameters

```hcl
data "tencentcloud_cat_probe_metric_tag_values" "example" {
  analyze_task_type = "AnalyzeTaskType_Network"
  key               = "host"
  filter            = "www.qq.com"
  time_range        = "1h"
}

output "tag_values" {
  value = data.tencentcloud_cat_probe_metric_tag_values.example.tag_value_set
}
```

- **WHEN** user provides all optional parameters to query probe metric tag values
- **THEN** the data source calls `DescribeProbeMetricTagValues` API with the provided parameters
- **AND** returns the `tag_value_set` attribute containing the serialized tag value string

#### Scenario: Query tag values with filters array

```hcl
data "tencentcloud_cat_probe_metric_tag_values" "example" {
  analyze_task_type = "AnalyzeTaskType_Network"
  key               = "area"
  filters = [
    "\"host\" = 'www.qq.com'",
    "time >= now()-1h",
  ]
}
```

- **WHEN** user provides the `filters` array parameter
- **THEN** the data source passes the filters as a string array to the API
- **AND** returns the matching tag values in `tag_value_set`

#### Scenario: Output results to file

- **WHEN** user specifies `result_output_file` parameter
- **THEN** the tag value set is written to the specified file

### Requirement: API Integration

The data source SHALL integrate CAT API v20180409 `DescribeProbeMetricTagValues` interface.

#### Scenario: API call execution with retry

- **WHEN** the data source is read
- **THEN** it calls `DescribeProbeMetricTagValues` via the service layer method `DescribeCatProbeMetricTagValuesByFilter`
- **AND** wraps the API call in `resource.Retry` with `tccommon.ReadRetryTimeout` timeout
- **AND** handles API errors using `tccommon.RetryError()` for retryable errors
- **AND** logs API calls using `tccommon.LogElapsed` and debug logging

#### Scenario: Handle empty API response

- **WHEN** the API returns an empty response (`response == nil` or `response.Response == nil`)
- **THEN** the data source returns `NonRetryableError` instead of clearing the state ID
- **AND** logs `[DATASOURCE] read empty, skip SetId` for troubleshooting

### Requirement: Service Layer Method

The service layer SHALL add a `DescribeCatProbeMetricTagValuesByFilter` method that wraps the `DescribeProbeMetricTagValues` API call.

#### Scenario: Service method implementation

- **WHEN** the Read function calls the service method with a parameter map
- **THEN** the method creates a `DescribeProbeMetricTagValuesRequest` and populates fields from the parameter map:
  - `AnalyzeTaskType` from `analyze_task_type`
  - `Key` from `key`
  - `Filter` from `filter`
  - `Filters` from `filters`
  - `TimeRange` from `time_range`
- **AND** calls `me.client.UseCatClient().DescribeProbeMetricTagValues(request)`
- **AND** returns the `Response` params for the caller to process

### Requirement: Provider Registration

The data source SHALL be registered in the Provider with the name `tencentcloud_cat_probe_metric_tag_values`.

#### Scenario: Register data source in provider

- **WHEN** the provider is initialized
- **THEN** `tencentcloud_cat_probe_metric_tag_values` is registered as a data source
- **AND** maps to `cat.DataSourceTencentCloudCatProbeMetricTagValues`

### Requirement: Error Handling

The data source SHALL properly handle various error scenarios.

#### Scenario: API permission error

- **WHEN** the API call fails due to insufficient permissions
- **THEN** the data source returns the original API error message to help with troubleshooting

#### Scenario: Network or service error

- **WHEN** the API call fails due to network issues or service unavailability
- **THEN** the data source retries according to the configured retry policy using `tccommon.RetryError()`
- **AND** returns an appropriate error if all retries fail

