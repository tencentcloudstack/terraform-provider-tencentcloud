## 1. Service 层实现

- [x] 1.1 在 `tencentcloud/services/cat/service_tencentcloud_cat.go` 中新增 `DescribeCatProbeMetricTagValuesByFilter` 方法
- [x] 1.2 方法接收 `paramMap` 参数，创建 `DescribeProbeMetricTagValuesRequest` 并填充 `AnalyzeTaskType`、`Key`、`Filter`、`Filters`、`TimeRange` 字段
- [x] 1.3 调用 `me.client.UseCatClient().DescribeProbeMetricTagValues(request)` 并返回 `*cat.DescribeProbeMetricTagValuesResponseParams`

## 2. 数据源核心功能实现

- [x] 2.1 创建 `tencentcloud/services/cat/data_source_tc_cat_probe_metric_tag_values.go` 文件
- [x] 2.2 实现 `DataSourceTencentCloudCatProbeMetricTagValues` 函数，定义 Schema：
  - `analyze_task_type` (String, Optional)
  - `key` (String, Optional)
  - `filter` (String, Optional)
  - `filters` (TypeSet of String, Optional)
  - `time_range` (String, Optional)
  - `tag_value_set` (String, Computed)
  - `result_output_file` (String, Optional)
- [x] 2.3 实现 `dataSourceTencentCloudCatProbeMetricTagValuesRead` 函数，构建 paramMap 并调用 service 层方法
- [x] 2.4 在 Read 函数中使用 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 包裹 API 调用
- [x] 2.5 在 retry 块内检查响应是否为空，若为空返回 `NonRetryableError`，并打印 `[DATASOURCE] read empty, skip SetId` 日志
- [x] 2.6 使用 `tccommon.RetryError(e)` 包装 API 错误
- [x] 2.7 使用 `helper.DataResourceIdsHash([]string{tagValueSet})` 生成数据源 ID
- [x] 2.8 支持 `result_output_file` 参数将结果写入文件

## 3. Provider 注册

- [x] 3.1 在 `tencentcloud/provider.go` 中注册 `tencentcloud_cat_probe_metric_tag_values` 数据源，映射到 `cat.DataSourceTencentCloudCatProbeMetricTagValues`

## 4. 文档编写

- [x] 4.1 创建 `tencentcloud/services/cat/data_source_tc_cat_probe_metric_tag_values.md` 文件
- [x] 4.2 包含一句话描述（带云产品名称 CAT）
- [x] 4.3 包含 Example Usage 部分和使用示例

## 5. 单元测试编写

- [x] 5.1 创建 `tencentcloud/services/cat/data_source_tc_cat_probe_metric_tag_values_test.go` 文件
- [x] 5.2 使用 gomonkey mock 云 API 方法，编写业务逻辑单元测试
- [x] 5.3 测试 Read 函数的正常查询场景
- [x] 5.4 测试 Read 函数的空响应处理场景

## 6. 验证任务

- [ ] 6.1 运行 `make doc` 生成 website/docs/ 下的文档
- [ ] 6.2 运行 `gofmt` 格式化代码
