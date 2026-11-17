# cat-instant-tasks-datasource Specification

## Purpose
TBD - created by syncing change add-cat-instant-tasks-datasource. Update Purpose after archive.
## Requirements
### Requirement: Data Source Schema Definition

`tencentcloud_cat_instant_tasks` 数据源 SHALL 支持以下输出属性，所有列表项字段平铺在 `tasks` 列表元素内（不额外嵌套一层），顶层提供 `total` 与 `result_output_file`：

**Output Attributes:**
- `tasks` (List): 历史即时拨测任务列表，每个元素包含：
  - `task_id` (String): 任务 ID（映射 `SingleInstantTask.TaskId`）
  - `target_address` (String): 任务地址（映射 `SingleInstantTask.TargetAddress`）
  - `task_type` (Int): 任务类型（映射 `SingleInstantTask.TaskType`）
  - `probe_time` (Int): 测试时间（映射 `SingleInstantTask.ProbeTime`）
  - `status` (String): 任务状态（映射 `SingleInstantTask.Status`）
  - `success_rate` (Float): 成功率（映射 `SingleInstantTask.SuccessRate`）
  - `node_count` (Int): 节点数量（映射 `SingleInstantTask.NodeCount`）
  - `task_category` (Int): 节点类型（映射 `SingleInstantTask.TaskCategory`）
- `total` (Int): 任务总数（映射 `response.Response.Total`）
- `result_output_file` (String, Optional): 将查询结果导出到指定文件

数据源 SHALL NOT 向用户暴露 `limit`/`offset` 分页参数；分页 SHALL 在 service 层内部自动完成，获取全部历史即时拨测任务。

#### Scenario: Query all instant tasks

```hcl
data "tencentcloud_cat_instant_tasks" "example" {
}

output "tasks" {
  value = data.tencentcloud_cat_instant_tasks.example.tasks
}

output "total" {
  value = data.tencentcloud_cat_instant_tasks.example.total
}
```

- **WHEN** 用户不带任何参数查询 `tencentcloud_cat_instant_tasks`
- **THEN** 返回当前账号下全部历史即时拨测任务
- **AND** `tasks` 列表中每个元素包含 `task_id`/`target_address`/`task_type`/`probe_time`/`status`/`success_rate`/`node_count`/`task_category` 字段
- **AND** `total` 反映任务总数

#### Scenario: Export query results to file

```hcl
data "tencentcloud_cat_instant_tasks" "export" {
  result_output_file = "/tmp/cat_instant_tasks.json"
}
```

- **WHEN** 用户指定 `result_output_file`
- **THEN** 查询结果以 JSON 格式写入指定文件
- **AND** 文件内容包含完整的任务列表信息

### Requirement: Service Layer Implementation

`service_tencentcloud_cat.go` SHALL 新增方法用于封装 `DescribeInstantTasks` 接口，内部实现自动分页：

```go
func (me *CatService) DescribeCatInstantTasksByFilter(ctx context.Context) (tasks []*cat.SingleInstantTask, total *uint64, errRet error)
```

实现约束：
1. 使用固定 `pageSize`（云 API 支持的最大值）循环分页，累积 `response.Response.Tasks`，直到本页返回数量小于 `pageSize` 即停止
2. 每次请求前调用 `ratelimit.Check(request.GetAction())`
3. 每次请求后记录 `[DEBUG]` 级别的请求体与响应体日志
4. 失败时记录 `[CRITAL]` 级别错误日志并返回错误
5. 返回累积后的全部 `tasks` 与最后一次响应的 `total`

#### Scenario: Service method handles empty list

- **WHEN** API 返回非 nil 响应但 `response.Response.Tasks` 为空列表
- **THEN** service 方法返回空（长度为 0）的 `tasks` 切片与对应的 `total`
- **AND** 不返回错误

#### Scenario: Service method handles API errors

- **WHEN** API 调用返回错误
- **THEN** service 方法记录 `[CRITAL]` 错误日志（含请求体与错误原因）
- **AND** 返回该错误，由上层 retry 机制决定是否重试

#### Scenario: Service method paginates all results

- **WHEN** 历史即时拨测任务数量超过单页 `pageSize`
- **THEN** service 方法循环调用 `DescribeInstantTasks`，逐步累加 `offset`
- **AND** 返回所有页累积合并后的完整 `tasks` 列表

### Requirement: Read Retry and Empty-Response Handling

Data Source 的 Read 方法 SHALL 使用 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 包裹 service 方法调用，并通过 `tccommon.RetryError(e)` 包装错误。在 retry 块内 SHALL 满足：

1. 若 service 方法返回错误，使用 `tccommon.RetryError(e)` 包装返回，由外层 retry 决定重试
2. 若云 API 返回空（`response == nil` / `response.Response == nil` 等，表现为 `tasks == nil && total == nil`），SHALL 直接返回 `tccommon.NonRetryableError`，不调用 `d.SetId("")`，保留现场便于排障
3. 若 API 返回正常但列表为空（`len(tasks) == 0`），SHALL 视为正常空结果，打印 `log.Printf("[DATASOURCE] read empty, skip SetId")` 后继续后续逻辑，不返回错误

#### Scenario: Retry on transient API failure

- **WHEN** `DescribeInstantTasks` 出现临时性错误
- **THEN** 通过 `resource.Retry` 自动重试，直到成功或达到 `tccommon.ReadRetryTimeout` 超时
- **AND** 重试耗尽后以"重试耗尽"形式失败，便于人工介入

#### Scenario: Reject clearing id on nil response

- **WHEN** API 返回的响应为 nil（异常情况）
- **THEN** Read 方法返回 `NonRetryableError`，不执行 `d.SetId("")`
- **AND** 避免因 API 短暂波动导致本地 state 中的 id 被清空

### Requirement: Provider Registration

`tencentcloud/provider.go` 的 `DataSourcesMap` SHALL 注册新数据源：

```go
"tencentcloud_cat_instant_tasks": cat.DataSourceTencentCloudCatInstantTasks(),
```

#### Scenario: Data source is accessible via Terraform

- **WHEN** 用户在 Terraform 配置中使用 `data "tencentcloud_cat_instant_tasks"`
- **THEN** Provider 能识别该数据源，不报 "provider doesn't support data source" 错误

### Requirement: Documentation

数据源 SHALL 提供文档文件 `data_source_tc_cat_instant_tasks.md`，内容包括：
1. 一句话描述，须带上所属云产品名称（CAT / 云拨测）
2. 使用示例（基本查询）
3. 不手写 `Argument Reference` 与 `Attribute Reference`（由工具自动生成）

`provider.md` 中 CAT 部分的 Data Source 条目 SHALL 通过 `make doc` 自动生成，不手动编辑 `website/` 目录。

#### Scenario: User finds data source in documentation

- **WHEN** 用户查看生成的文档
- **THEN** 能看到 `tencentcloud_cat_instant_tasks` 的一句话描述与示例
- **AND** 文档由 `make doc` 自动生成，无需手动维护 `website/docs/`

### Requirement: Testing

数据源 SHALL 提供单元测试文件 `data_source_tc_cat_instant_tasks_test.go`，使用 mock（gomonkey）方式对云 API 进行 mock，只做业务代码逻辑的单元测试（不使用 terraform 测试套件）。测试 SHALL 覆盖：
1. 正常查询场景（mock 返回非空任务列表，验证字段映射与 `total` 设置）
2. 空列表场景（mock 返回空列表，验证不报错、正常设置 id）
3. API 错误场景（mock 返回错误，验证错误被返回）

#### Scenario: Unit test covers normal query

- **WHEN** mock `DescribeInstantTasks` 返回包含若干 `SingleInstantTask` 的响应
- **THEN** 单元测试验证 `tasks` 列表字段被正确映射、`total` 被正确设置、`d.Id()` 被设置为任务 id 哈希

#### Scenario: Unit test covers empty list

- **WHEN** mock `DescribeInstantTasks` 返回空任务列表
- **THEN** 单元测试验证 Read 不报错并正常完成
