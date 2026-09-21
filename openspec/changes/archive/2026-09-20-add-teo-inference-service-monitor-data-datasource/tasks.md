## 1. 数据源 Schema 与 Read 函数实现

- [x] 1.1 在 `tencentcloud/services/teo/` 目录下创建 `data_source_tc_teo_inference_service_monitor_data.go` 文件
- [x] 1.2 定义 `DataSourceTencentCloudTeoInferenceServiceMonitorData()` 函数，返回 `*schema.Resource`
- [x] 1.3 定义输入参数 schema：`zone_id`(String,Required)、`service_ids`(List of String,Required)、`metric_names`(List of String,Required)、`start_time`(String,Required)、`end_time`(String,Required)、`interval`(String,Optional)、`result_output_file`(String,Optional)
- [x] 1.4 定义输出参数 schema：`inference_service_monitor_records`(TypeList,Computed)，每条记录含 `service_id`(String)、`metric_name`(String)、`inference_service_monitor_items`(TypeList)，每项含 `timestamp`(String)、`value`(Float)
- [x] 1.5 为所有字段添加 Description 说明
- [x] 1.6 实现 `dataSourceTencentCloudTeoInferenceServiceMonitorDataRead()`，添加 `defer tccommon.LogElapsed()` 与 `defer tccommon.InconsistentCheck()`
- [x] 1.7 使用 `ratelimit.Check(request.GetAction())` 进行限流
- [x] 1.8 使用 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 包装 `DescribeInferenceServiceMonitorDataWithContext` 调用，失败时用 `tccommon.RetryError(e)` 包装
- [x] 1.9 在 retry 块内检查空响应（response/Response 为 nil 或记录列表为空），返回 `NonRetryableError` 并打印 `[DATASOURCE] read empty, skip SetId`，不直接 `d.SetId("")`
- [x] 1.10 retry 块外逐层 nil 检查后 set 输出字段（records → items）
- [x] 1.11 使用 `helper.BuildToken()` 生成数据源 ID
- [x] 1.12 处理 `result_output_file` 导出逻辑

## 2. Provider 注册

- [x] 2.1 在 `tencentcloud/provider.go` 的 `DataSourcesMap` 中添加 `"tencentcloud_teo_inference_service_monitor_data": teo.DataSourceTencentCloudTeoInferenceServiceMonitorData()`

## 3. 文档样例

- [x] 3.1 在 `tencentcloud/services/teo/` 目录下创建 `data_source_tc_teo_inference_service_monitor_data.md` 文件
- [x] 3.2 添加一句话描述（含云产品名称 TEO，格式为 "Use this data source to query ..."）
- [x] 3.3 添加 Example Usage（基本查询、带 interval 查询、导出到文件）
- [x] 3.4 不添加 Argument Reference 和 Attribute Reference 部分（由 make doc 自动生成）

## 4. 单元测试

- [x] 4.1 在 `tencentcloud/services/teo/` 目录下创建 `data_source_tc_teo_inference_service_monitor_data_test.go` 文件
- [x] 4.2 使用 gomonkey mock `DescribeInferenceServiceMonitorDataWithContext` 方法，编写业务逻辑单元测试（不使用 Terraform 测试套件）
- [x] 4.3 测试正常响应的字段映射
- [x] 4.4 测试空响应处理逻辑
- [x] 4.5 测试 API 调用失败的重试逻辑

## 5. 收尾与文档生成

- [x] 5.1 通过收尾阶段执行 gofmt 格式化
- [x] 5.2 通过收尾阶段执行 make doc 生成 website 文档
- [x] 5.3 通过收尾阶段创建 .changelog 文件