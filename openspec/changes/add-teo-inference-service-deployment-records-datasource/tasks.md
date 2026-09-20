## 1. Service Layer

- [x] 1.1 Add `DescribeTeoInferenceServiceDeploymentRecordsByFilter` method to `tencentcloud/services/teo/service_tencentcloud_teo.go` that calls `DescribeInferenceServiceDeploymentRecords` API with paramMap containing `ZoneId`、`ServiceId`、`SortBy`、`SortOrder`, implements internal pagination (Limit fixed to API max value 100, Offset incremented until all records fetched), merges all `RecordSet` results, and returns `[]*teov20220901.InferenceServiceDeploymentRecord`

## 2. Data Source Schema and Read Function

- [x] 2.1 Create `tencentcloud/services/teo/data_source_tc_teo_inference_service_deployment_records.go` with `DataSourceTencentCloudTeoInferenceServiceDeploymentRecords()` function defining the schema: `zone_id` (Required)、`service_id` (Required)、`sort_by` (Optional)、`sort_order` (Optional)、`result_output_file` (Optional), and computed `record_set` (TypeList) with flattened nested fields (record_id、operation、status、duration、create_time、active_status、inference_service_config 及其嵌套的 containers/resource_config/affinity_config 等全部字段)
- [x] 2.2 Implement `dataSourceTencentCloudTeoInferenceServiceDeploymentRecordsRead()` function with: retry logic using `tccommon.ReadRetryTimeout`; empty response handling returning `NonRetryableError` (不调用 `d.SetId("")`) 并保留 `[DATASOURCE] read empty, skip SetId` 日志; nil-safe nested response mapping at each level; float64 (gpu_num/cpu_num) to int conversion; composite ID using `zone_id + FILED_SP + service_id`; result output file support

## 3. Provider Registration

- [x] 3.1 Add data source entry `"tencentcloud_teo_inference_service_deployment_records": teo.DataSourceTencentCloudTeoInferenceServiceDeploymentRecords()` to `tencentcloud/provider.go` dataSources map
- [x] 3.2 Add data source entry to `tencentcloud/provider.md` data source list

## 4. Documentation

- [x] 4.1 Create `tencentcloud/services/teo/data_source_tc_teo_inference_service_deployment_records.md` with one-line description mentioning TEO product name, Example Usage section (using jsonencode() for json string fields if applicable), no Argument/Attribute Reference sections, no Import section (DATASOURCE)

## 5. Unit Tests

- [x] 5.1 Create `tencentcloud/services/teo/data_source_tc_teo_inference_service_deployment_records_test.go` with unit tests using gomonkey mock approach to mock the cloud API / service method and verify the Read function correctly maps nested response fields into the schema (不使用 Terraform 测试套件)

## 6. Verification

- [x] 6.1 Verify all new and modified files compile correctly and follow the project patterns (代码正确性检查: 确认入参 ZoneId/ServiceId/SortBy/SortOrder 与云 API 入参一致, 出参字段与云 API 出参路径一致)
- [x] 6.2 Run `gofmt` formatting and `make doc` for documentation generation (由 tfpacer-finalize 阶段统一执行)