## 1. 服务层实现

- [x] 1.1 在 `tencentcloud/services/teo/service_tencentcloud_teo.go` 中新增服务方法 `DescribeTeoInferenceServiceDeploymentLogs(ctx, param map[string]interface{}) (ret []*teov20220901.InferenceServiceDeploymentLogInfo, errRet error)`，内部使用 `Offset`（从 0 递增）与 `Limit`（固定 1000）自动分页循环调用云 API `DescribeInferenceServiceDeploymentLogs`，累积所有日志条目，直到某次返回数量小于 `Limit` 退出循环；从 `param` 中读取 `ZoneId`、`ServiceId`、`RecordId`、`StartTime`、`EndTime`、`SortBy`、`SortOrder` 赋值给请求体。
- [x] 1.2 服务方法中调用云 API 失败时记录日志并返回错误；响应为空时退出分页循环；保留 `defer` 错误日志打印（参考既有 `DescribeTeoConfigGroupVersionsByFilter` 实现）。

## 2. 数据源 schema 与 Read 实现

- [x] 2.1 新建 `tencentcloud/services/teo/data_source_tc_teo_inference_service_deployment_logs.go`，定义 `DataSourceTencentCloudTeoInferenceServiceDeploymentLogs()` 返回 schema：必填 `zone_id`、`service_id`、`record_id`；可选 `start_time`、`end_time`、`sort_by`、`sort_order`、`result_output_file`；Computed 出参 `deployment_log_info_set`（TypeList，元素为 TypeResource，含平铺字段 `log_message`、`timestamp`）。不暴露 `limit`/`offset`。
- [x] 2.2 实现 `dataSourceTencentCloudTeoInferenceServiceDeploymentLogsRead`：`defer tccommon.LogElapsed()` 与 `defer tccommon.InconsistentCheck()`；构造 `paramMap` 传入服务方法；使用 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 包装服务方法调用，失败时用 `tccommon.RetryError(e)` 包装。
- [x] 2.3 Read 的 retry 块内，若云 API 返回空列表（`len(respData) == 0`），返回 `NonRetryableError`，不直接 `d.SetId("")`；外层 retry 失败路径打印 `log.Printf("[DATASOURCE] read empty, skip SetId")`。
- [x] 2.4 Read 中将日志列表映射到 `deployment_log_info_set`，设置每个元素的 `log_message`、`timestamp` 前先判断对应云 API 字段非 nil；使用 `helper.DataResourceIdsHash(ids)`（ids 由各日志 `timestamp` 组成）设置数据源 ID；支持 `result_output_file` 写出结果。

## 3. provider 注册

- [x] 3.1 在 `tencentcloud/provider.go` 中注册数据源 `tencentcloud_teo_inference_service_deployment_logs`，映射到 `DataSourceTencentCloudTeoInferenceServiceDeploymentLogs`。
- [x] 3.2 在 `tencentcloud/provider.md` 中新增数据源 `tencentcloud_teo_inference_service_deployment_logs` 的文档索引行。

## 4. 数据源文档

- [x] 4.1 新建 `tencentcloud/services/teo/data_source_tc_teo_inference_service_deployment_logs.md`，包含以 "Use this data source to query" 开头且包含 "TEO" 的一句话描述、Example Usage 部分；不包含 `Argument Reference` 与 `Attribute Reference` 部分。

## 5. 单元测试

- [x] 5.1 新建 `tencentcloud/services/teo/data_source_tc_teo_inference_service_deployment_logs_test.go`，使用 gomonkey mock 云 API `DescribeInferenceServiceDeploymentLogs`，补充数据源 Read 业务逻辑的单元测试用例（包含正常返回日志列表、空返回的场景），不使用 terraform 测试套件。

## 6. 验证

- [x] 6.1 检查所有新增/修改的 *.go 文件可正确编译（由后续流程通过构建验证，本步骤不执行 go build）。
- [x] 6.2 确认云 API 入参（ZoneId、ServiceId、RecordId、StartTime、EndTime、SortBy、SortOrder）与出参（DeploymentLogInfoSet 含 LogMessage、Timestamp）映射与 vendor 中 `DescribeInferenceServiceDeploymentLogs` 接口定义一致。