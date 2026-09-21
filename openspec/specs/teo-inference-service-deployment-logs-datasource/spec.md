# teo-inference-service-deployment-logs-datasource Specification

## Purpose
TBD - created by archiving change add-teo-inference-service-deployment-logs-datasource. Update Purpose after archive.
## Requirements
### Requirement: 数据源 schema 定义

数据源 `tencentcloud_teo_inference_service_deployment_logs` SHALL 定义以下 schema 字段：

必填入参：
- `zone_id`（TypeString，Required）：站点 ID。
- `service_id`（TypeString，Required）：推理服务 ID。
- `record_id`（TypeString，Required）：部署记录 ID。

可选入参：
- `start_time`（TypeString，Optional）：需检索日志的开始时间。
- `end_time`（TypeString，Optional）：需检索日志的结束时间。
- `sort_by`（TypeString，Optional）：排序字段。
- `sort_order`（TypeString，Optional）：排序方式。

出参：
- `deployment_log_info_set`（TypeList，Computed）：部署日志列表，每个元素为 TypeResource，包含平铺字段 `log_message`（TypeString，Computed，日志消息内容）与 `timestamp`（TypeString，Computed，日志产生时间）。
- `result_output_file`（TypeString，Optional）：用于保存结果。

数据源 SHALL NOT 向用户暴露 `limit`/`offset` 分页参数。

#### Scenario: schema 包含所有入参与出参字段
- **WHEN** 检查 `DataSourceTencentCloudTeoInferenceServiceDeploymentLogs()` 返回的 schema
- **THEN** schema 中存在 `zone_id`、`service_id`、`record_id`（Required）、`start_time`、`end_time`、`sort_by`、`sort_order`（Optional）、`deployment_log_info_set`（Computed）、`result_output_file`（Optional）字段
- **AND** schema 中不存在 `limit` 与 `offset` 字段

### Requirement: 查询推理服务部署日志

数据源 Read 方法 SHALL 调用云 API `DescribeInferenceServiceDeploymentLogs`（teo v20220901）查询推理服务指定部署记录的日志。

Read 方法 SHALL：
1. 从 schema 读取 `zone_id`、`service_id`、`record_id`（必填）及可选的 `start_time`、`end_time`、`sort_by`、`sort_order`，构造请求参数。
2. 通过服务层方法调用云 API，服务层方法 SHALL 在内部使用 `Offset`/`Limit` 自动分页（`Limit` 取 1000）拉取所有日志，直到某次返回数量小于单页大小。
3. 在 Read 方法中使用 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 包装调用，失败时使用 `tccommon.RetryError(e)` 包装错误。
4. 将返回的日志列表映射到 `deployment_log_info_set`，每个元素的 `log_message` 与 `timestamp` 在对应云 API 字段非 nil 时设置。

#### Scenario: 正常查询返回日志列表
- **WHEN** 用户配置 `zone_id`、`service_id`、`record_id` 并执行 `terraform refresh`
- **THEN** 数据源调用 `DescribeInferenceServiceDeploymentLogs` 并通过内部自动分页获取所有日志
- **AND** 将日志列表设置到 `deployment_log_info_set`，每条日志的 `log_message` 与 `timestamp` 被正确填充
- **AND** 数据源 ID 被设置为基于返回内容计算的 hash

#### Scenario: 云 API 返回空列表
- **WHEN** 调用云 API 返回的 `DeploymentLogInfoSet` 为空（长度为 0）
- **THEN** Read 方法在 retry 块内返回 `NonRetryableError`，SHALL NOT 直接调用 `d.SetId("")`
- **AND** 在外层 retry 失败路径打印 `log.Printf("[DATASOURCE] read empty, skip SetId")` 提示

#### Scenario: 云 API 调用失败并重试
- **WHEN** 调用云 API 返回错误
- **THEN** retry 块使用 `tccommon.RetryError(e)` 包装错误以触发重试
- **AND** 在 `tccommon.ReadRetryTimeout` 超时前持续重试

### Requirement: 设置字段前判断非 nil

数据源 Read 方法在将云 API 返回字段设置到 schema 前，SHALL 判断对应字段是否为 nil；若为 nil 则跳过该字段的 set 操作。

#### Scenario: 日志条目字段为 nil 时不 set
- **WHEN** 云 API 返回的某条日志的 `LogMessage` 或 `Timestamp` 为 nil
- **THEN** Read 方法跳过对该 nil 字段的 set 操作，不将其设置到 `deployment_log_info_set` 对应元素中

### Requirement: 自动分页获取所有日志

服务层方法 `DescribeTeoInferenceServiceDeploymentLogs` SHALL 在内部循环调用云 API，使用 `Offset`（从 0 递增）与 `Limit`（固定为 1000，即云 API 注释标注的最大值）自动分页，累积所有页的日志条目，直到某次返回数量小于 `Limit` 时退出循环。

#### Scenario: 多页日志自动合并
- **WHEN** 符合条件的日志总数超过 1000 条
- **THEN** 服务层方法多次调用云 API（Offset 递增 1000）
- **AND** 所有页的日志条目被合并到同一返回列表中

### Requirement: provider 注册数据源

`tencentcloud/provider.go` SHALL 注册数据源 `tencentcloud_teo_inference_service_deployment_logs`，`tencentcloud/provider.md` SHALL 新增该数据源的文档索引。

#### Scenario: provider 中可使用该数据源
- **WHEN** 用户在 Terraform 配置中声明 `data "tencentcloud_teo_inference_service_deployment_logs" "example"`
- **THEN** provider 正确识别该数据源类型并可执行查询

### Requirement: 数据源文档

SHALL 生成数据源文档 `data_source_tc_teo_inference_service_deployment_logs.md`，包含：
- 一句话描述（包含所属云产品名称 TEO），格式为 "Use this data source to query ..."。
- Example Usage 部分。
- 不包含 `Argument Reference` 与 `Attribute Reference` 部分（由工具自动生成）。

#### Scenario: 文档内容符合规范
- **WHEN** 检查 `data_source_tc_teo_inference_service_deployment_logs.md`
- **THEN** 文件包含以 "Use this data source to query" 开头的一句话描述，包含 "TEO"
- **AND** 包含 Example Usage 部分
- **AND** 不包含 `Argument Reference` 与 `Attribute Reference` 部分

