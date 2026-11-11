## 1. Preparation

- [x] 1.1 确认 SDK 已支持 `DescribeProbeTasks` 接口（vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cat/v20180409，已存在 `DescribeProbeTasksRequest`/`DescribeProbeTasksResponse`/`ProbeTask`/`KeyValuePair` 结构）
- [x] 1.2 确认 cat service 目录已存在（tencentcloud/services/cat/），参考现有 `data_source_tc_cat_instant_tasks.go` 与 `service_tencentcloud_cat.go` 代码风格

## 2. Service Layer

- [x] 2.1 在 `tencentcloud/services/cat/service_tencentcloud_cat.go` 中新增方法 `DescribeCatProbeTasksByFilter(ctx context.Context, param map[string]interface{}) (tasks []*cat.ProbeTask, total *int64, errRet error)`
- [x] 2.2 实现自动分页：固定 `pageSize=100`（云 API 最大值），循环累加 `response.Response.TaskSet`，直到本页数量小于 pageSize 停止
- [x] 2.3 从 `param` 读取并设置 request 字段：TaskIDs/TaskName/TargetAddress/TaskStatus/PayMode/OrderState/TaskType/TaskCategory/OrderBy/Ascend/TagFilters
- [x] 2.4 每次请求前 `ratelimit.Check(request.GetAction())`，请求后记录 [DEBUG] 日志，失败记录 [CRITAL] 日志并返回错误

## 3. Data Source Schema & Read

- [x] 3.1 创建 `tencentcloud/services/cat/data_source_tc_cat_probe_tasks.go`，定义 `DataSourceTencentCloudCatProbeTasks()` schema
- [x] 3.2 定义可选入参 schema：task_i_ds(List string)、task_name(string)、target_address(string)、task_status(List int)、pay_mode(int)、order_state(int)、task_type(List int)、task_category(List int)、order_by(string)、ascend(bool)、tag_filters(List block: key/value 必填)、result_output_file(string)
- [x] 3.3 定义出参 schema：task_set(List block，字段平铺 name/task_id/task_type/nodes/node_ip_type/interval/parameters/status/target_address/pay_mode/order_state/task_category/created_at/cron/cron_state/tag_info_list(List block: key/value)/sub_sync_flag)、total(int)
- [x] 3.4 实现 `dataSourceTencentCloudCatProbeTasksRead`：使用 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 包裹 service 调用，错误用 `tccommon.RetryError(e)` 包装
- [x] 3.5 retry 块内：nil 响应（tasks==nil && total==nil）返回 NonRetryableError 不清 id；空列表（len==0）打印 `[DATASOURCE] read empty, skip SetId` 后继续
- [x] 3.6 遍历 tasks 构造 task_set 列表 map，设置每个字段前判断 nil；设置 total；`d.SetId(helper.BuildToken())`；处理 result_output_file 导出

## 4. Provider Registration

- [x] 4.1 在 `tencentcloud/provider.go` 的 DataSourcesMap 中添加 `"tencentcloud_cat_probe_tasks": cat.DataSourceTencentCloudCatProbeTasks()`

## 5. Documentation

- [x] 5.1 创建 `tencentcloud/services/cat/data_source_tc_cat_probe_tasks.md`：一句话描述（带 CAT/云拨测产品名）、Example Usage（基本查询、按条件过滤查询），不手写 Argument/Attribute Reference
- [x] 5.2 后续通过 `make doc` 自动生成 `website/docs/` 与 provider.md 条目（不在本步骤手动编辑 website/）

## 6. Testing

- [x] 6.1 创建 `tencentcloud/services/cat/data_source_tc_cat_probe_tasks_test.go`，使用 gomonkey mock `DescribeProbeTasksWithContext`（不使用 terraform 测试套件）
- [x] 6.2 测试正常查询场景：mock 返回非空 ProbeTask 列表，验证 task_set 字段映射与 total 设置
- [x] 6.3 测试空列表场景：mock 返回空列表，验证不报错、正常设置 id
- [x] 6.4 测试 API 错误场景：mock 返回错误，验证错误被返回

## 7. Validation

- [x] 7.1 检查所有函数返回的 error 均被处理（无未使用变量；必定不出错的用 `_ =` 赋值）
- [x] 7.2 确认生成的代码可正确构建执行（不执行 go build/vet/lint，由其他流程验证）
- [x] 7.3 运行 `openspec validate add-cat-probe-tasks-datasource --strict` 验证提案
