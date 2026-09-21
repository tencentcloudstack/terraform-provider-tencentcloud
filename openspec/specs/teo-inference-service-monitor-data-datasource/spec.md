# teo-inference-service-monitor-data-datasource Specification

## Purpose
TBD - created by archiving change add-teo-inference-service-monitor-data-datasource. Update Purpose after archive.
## Requirements
### Requirement: Data Source Schema Definition

The data source `tencentcloud_teo_inference_service_monitor_data` MUST support the following input parameters and output attributes, fully mapping the `DescribeInferenceServiceMonitorData` API. The schema SHALL define all listed input parameters and output attributes.

**Input Parameters:**
- `zone_id` (String, Required): 站点 ID
- `service_ids` (List of String, Required): 推理服务 ID 列表，最多 10 个
- `metric_names` (List of String, Required): 指标列表，最多 10 个
- `start_time` (String, Required): 开始时间
- `end_time` (String, Required): 结束时间（EndTime - StartTime ≤ 30 天）
- `interval` (String, Optional): 查询时间粒度，取值 min/5min/hour/day
- `result_output_file` (String, Optional): 用于保存查询结果

**Output Attributes:**
- `inference_service_monitor_records` (List): 推理服务监控记录列表，每条记录包含:
  - `service_id` (String): 推理服务 ID
  - `metric_name` (String): 指标名称
  - `inference_service_monitor_items` (List): 监控数据明细，每项包含:
    - `timestamp` (String): 监控数据对应时间点
    - `value` (Float): 具体数值

#### Scenario: Query inference service monitor data with required parameters

```hcl
data "tencentcloud_teo_inference_service_monitor_data" "basic" {
  zone_id      = "zone-xxxxx"
  service_ids  = ["service-xxxxx"]
  metric_names = ["cpu_usage_average"]
  start_time    = "2024-01-01T00:00:00Z"
  end_time      = "2024-01-01T01:00:00Z"
}

output "records" {
  value = data.tencentcloud_teo_inference_service_monitor_data.basic.inference_service_monitor_records
}
```

- **WHEN** 用户使用 zone_id、service_ids、metric_names、start_time、end_time 必填参数查询推理服务监控数据
- **THEN** 调用 `DescribeInferenceServiceMonitorData` API 并返回监控记录列表，每条记录包含 service_id、metric_name 及对应的监控数据明细

#### Scenario: Query inference service monitor data with interval

```hcl
data "tencentcloud_teo_inference_service_monitor_data" "with_interval" {
  zone_id      = "zone-xxxxx"
  service_ids  = ["service-xxxxx"]
  metric_names = ["cpu_usage_average", "gpu_usage_max"]
  start_time   = "2024-01-01T00:00:00Z"
  end_time     = "2024-01-02T00:00:00Z"
  interval      = "hour"
}

output "records" {
  value = data.tencentcloud_teo_inference_service_monitor_data.with_interval.inference_service_monitor_records
}
```

- **WHEN** 用户在必填参数基础上额外指定 interval 查询时间粒度
- **THEN** API 请求中携带 Interval 参数，返回按指定粒度聚合的监控数据明细

#### Scenario: Export query results to file

```hcl
data "tencentcloud_teo_inference_service_monitor_data" "export" {
  zone_id            = "zone-xxxxx"
  service_ids         = ["service-xxxxx"]
  metric_names        = ["cpu_usage_average"]
  start_time          = "2024-01-01T00:00:00Z"
  end_time            = "2024-01-01T01:00:00Z"
  result_output_file = "./monitor_data.json"
}
```

- **WHEN** 用户设置 result_output_file 参数
- **THEN** 将查询结果以 JSON 格式写入指定文件

### Requirement: Read Function Implementation

数据源 Read 函数 MUST 实现：
1. 从 ResourceData 读取所有输入参数，构造 `DescribeInferenceServiceMonitorDataRequest`
2. 使用 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 包装 API 调用，失败时通过 `tccommon.RetryError(e)` 包装返回
3. 在 retry 块内检查响应是否为空（`response == nil` / `response.Response == nil` / `len(response.Response.InferenceServiceMonitorRecords) == 0`），若为空返回 `NonRetryableError`，不直接 `d.SetId("")`
4. retry 块外处理成功响应：逐层 nil 检查后 set 输出字段
5. 使用 `defer tccommon.LogElapsed()` 与 `defer tccommon.InconsistentCheck()`
6. retry 失败路径保留 `log.Printf("[DATASOURCE] read empty, skip SetId")` 提示

#### Scenario: API call failure triggers retry

- **WHEN** `DescribeInferenceServiceMonitorData` API 调用返回错误
- **THEN** 错误经 `tccommon.RetryError(e)` 包装后由外层 retry 继续尝试，直到重试耗尽后向上返回错误

#### Scenario: Empty response does not clear state id

- **WHEN** API 返回空响应（response 或 Response 为 nil，或监控记录列表为空）
- **THEN** 返回 `NonRetryableError` 并打印 `[DATASOURCE] read empty, skip SetId` 日志，不清空数据源 ID，避免数据丢失

#### Scenario: Successful response maps to schema

- **WHEN** API 返回非空监控记录列表
- **THEN** 将每条 `InferenceServiceMonitorRecord` 转换为 map（含 service_id、metric_name），并将其 `InferenceServiceMonitorItems` 转换为嵌套列表（含 timestamp、value），set 到 `inference_service_monitor_records`

### Requirement: Provider Registration

MUST 在 `tencentcloud/provider.go` 的 `DataSourcesMap` 注册新数据源：

```go
"tencentcloud_teo_inference_service_monitor_data": teo.DataSourceTencentCloudTeoInferenceServiceMonitorData(),
```

#### Scenario: Data source is accessible via Terraform

- **WHEN** 用户在 Terraform 配置中使用 `data "tencentcloud_teo_inference_service_monitor_data"`
- **THEN** provider 正确识别该数据源，不出现 "provider doesn't support data source" 错误

### Requirement: Documentation

MUST 提供文档样例文件 `tencentcloud/services/teo/data_source_tc_teo_inference_service_monitor_data.md`：
1. 一句话描述，带上所属云产品名称（TEO），格式为 "Use this data source to query ..."
2. Example Usage 部分，包含基本查询与带 interval 查询示例
3. 不添加 Argument Reference 和 Attribute Reference 部分（由工具自动生成）

#### Scenario: User finds data source documentation

- **WHEN** 用户查看数据源文档
- **THEN** 能看到包含云产品名称的一句话描述、使用示例，参数说明由 make doc 自动生成

### Requirement: Unit Testing

MUST 创建测试文件 `data_source_tc_teo_inference_service_monitor_data_test.go`，使用 gomonkey mock 云 API 客户端方法，仅测试业务逻辑，不使用 Terraform 测试套件，不依赖真实云 API。

#### Scenario: Mocked unit test verifies read logic

- **WHEN** 运行业务逻辑单元测试
- **THEN** 通过 mock `DescribeInferenceServiceMonitorDataWithContext` 验证 Read 函数的参数构造、响应映射与空响应处理逻辑正确

