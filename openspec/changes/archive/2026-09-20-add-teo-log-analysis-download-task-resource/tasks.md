## 1. Service 层实现

- [x] 1.1 在 `tencentcloud/services/teo/service_tencentcloud_teo.go` 中新增方法 `DescribeTeoLogAnalysisDownloadTaskById(ctx, zoneId, area, taskId string) (*teov20220901.LogAnalysisDownloadTask, error)`：在 for 循环外构建 `DescribeLogAnalysisDownloadTasksRequest` 与 `Filters`（`Name: "task-id"`, `Values: [taskId]`），仅循环内变更 `Offset`/`Limit`；`Limit=100`（API 注释最大值）；分页直到累计结果数 >= `TotalCount`；严格相等 `*task.TaskId == taskId` 命中后返回；未找到返回 `(nil, nil)`；每次分页 SDK 调用包装在 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 中并防御性检查 `result == nil || result.Response == nil`
- [x] 1.2 确认 `connectivity.UseTeoV20220901Client` 已就绪，SDK 导入路径为 `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901`

## 2. 资源实现

- [x] 2.1 创建 `tencentcloud/services/teo/resource_tc_teo_log_analysis_download_task.go`（单文件），顶层布局：package + imports → `ResourceTencentCloudTeoLogAnalysisDownloadTask()` schema → `resourceTencentCloudTeoLogAnalysisDownloadTaskCreate/Read/Update/Delete`
- [x] 2.2 schema 字段声明顺序：`zone_id`（Required, ForceNew）、`area`（Required, ForceNew）、`start_time`（Required, ForceNew）、`end_time`（Required, ForceNew）、`log_type`（Optional, ForceNew）、`condition`（Optional, ForceNew）、`format`（Optional, ForceNew）、`sort`（Optional, ForceNew）、`task_id`（Computed）、`status`（Computed）、`create_time`（Computed）、`url`（Computed）、`expire_time`（Computed）
- [x] 2.3 实现 `Create`：从 schema 构建 `CreateLogAnalysisDownloadTaskRequest` 填充 ZoneId/Area/StartTime/EndTime/LogType/Condition/Format/Sort；SDK 调用包在 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 中，retry 块内检查 `result == nil || result.Response == nil` 返回 `NonRetryableError`，成功打印 `[DEBUG]`；retry 块外打印 `logId` 与 `d.Id()`，检查 `response.Response.TaskId` 非空否则返回错误；`d.SetId(zoneId + tccommon.FILED_SP + area + tccommon.FILED_SP + taskId)`；返回 Read
- [x] 2.4 实现 `Read`：拆分复合 ID `ZoneId#Area#TaskId`（校验 3 段，否则报错 ID broken）；调用 `DescribeTeoLogAnalysisDownloadTaskById`；返回 `(nil,nil)` 时打印 `log.Printf("[CRUD] teo_log_analysis_download_task id=%s", d.Id())` 后 `d.SetId("")` 返回；命中则对每个字段先 nil 检查再 `_ = d.Set(...)`
- [x] 2.5 实现 `Update`：拆分复合 ID；定义 `immutableArgs` 含 `zone_id, area, start_time, end_time, log_type, condition, format, sort`；遍历检查 `d.HasChange(v)` 命中则返回 error 提示字段不可变更需重建；无变更则返回 Read
- [x] 2.6 实现 `Delete`：no-op，仅 `defer tccommon.LogElapsed(...)` 与 `defer tccommon.InconsistentCheck(...)`，`return nil`
- [x] 2.7 在每个 CRUD 函数顶部添加 `defer tccommon.LogElapsed("resource.tencentcloud_teo_log_analysis_download_task.<op>")()` 与 `defer tccommon.InconsistentCheck(d, meta)()`
- [x] 2.8 Importer：使用 `schema.ImportStatePassthrough`

## 3. Provider 注册

- [x] 3.1 在 `tencentcloud/provider.go` 的 `ResourcesMap` 中注册 `"tencentcloud_teo_log_analysis_download_task": teo.ResourceTencentCloudTeoLogAnalysisDownloadTask()`，放置于现有 `tencentcloud_teo_*` 条目附近以保持命名空间连续

## 4. 文档示例

- [x] 4.1 创建 `tencentcloud/services/teo/resource_tc_teo_log_analysis_download_task.md`：一句话描述提及 EdgeOne（TEO）；Example Usage 展示 zone_id/area/start_time/end_time 等主要输入字段；Import 部分展示 `terraform import tencentcloud_teo_log_analysis_download_task.example zone-xxx#mainland#task-yyy`（说明使用联合 ID ZoneId#Area#TaskId）；不要添加 Argument Reference 和 Attribute Reference 部分

## 5. 单元测试

- [x] 5.1 创建 `tencentcloud/services/teo/resource_tc_teo_log_analysis_download_task_test.go`，使用 gomonkey mock 云 API（不使用 terraform 测试套件）
- [x] 5.2 新增 `TestAccTencentCloudTeoLogAnalysisDownloadTaskCreate` 用例：mock `CreateLogAnalysisDownloadTaskWithContext` 返回 TaskId，mock `DescribeLogAnalysisDownloadTasksWithContext` 返回任务详情，验证 Create 后 state 正确回填
- [x] 5.3 新增 `TestAccTencentCloudTeoLogAnalysisDownloadTaskRead` 用例：mock `DescribeLogAnalysisDownloadTasksWithContext` 返回任务，验证 Read 各字段回填
- [x] 5.4 新增 `TestAccTencentCloudTeoLogAnalysisDownloadTaskUpdateImmutable` 用例：模拟字段变更，验证 Update 返回 immutable error
- [x] 5.5 新增 `TestAccTencentCloudTeoLogAnalysisDownloadTaskDelete` 用例：验证 Delete 为 no-op 不调用云 API 且返回 nil

## 6. 验证

- [x] 6.1 运行 `go build ./tencentcloud/...` 确认编译通过（由其他流程执行，本任务仅确认代码生成正确）
- [x] 6.2 确认所有函数返回的 error 已检查；必定不出错的函数用 `_ = func()` 处理
- [x] 6.3 确认 Create/Update 调用的接口入参在云 API 创建/更新接口中存在，Read 调用的入参在查询接口中存在（参数一致性检查）