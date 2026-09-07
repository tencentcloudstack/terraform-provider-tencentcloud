## Context

云拨测（CAT，Cat）服务已在 Provider 中提供以下资源/数据源：
- 资源 `tencentcloud_cat_task_set`（基于 `DescribeProbeTasks`/`CreateProbeTasks`/`DeleteProbeTask` 等接口）
- 数据源 `tencentcloud_cat_node`（拨测节点）
- 数据源 `tencentcloud_cat_probe_data`（拨测明细数据）
- 数据源 `tencentcloud_cat_metric_data`（拨测指标数据）

本次新增数据源 `tencentcloud_cat_instant_tasks`，用于查询**历史即时拨测任务**列表。对应云 API 为 `DescribeInstantTasks`（`cat/v20180409`），返回结构为 `response.Response.Tasks`（`[]*SingleInstantTask`）和 `response.Response.Total`（`*uint64`）。

`SingleInstantTask` 结构体字段如下（来自 vendor）：
- `TaskId *string` —— 任务 ID
- `TargetAddress *string` —— 任务地址
- `TaskType *uint64` —— 任务类型
- `ProbeTime *uint64` —— 测试时间
- `Status *string` —— 任务状态
- `SuccessRate *float64` —— 成功率
- `NodeCount *uint64` —— 节点数量
- `TaskCategory *uint64` —— 节点类型

请求参数仅有 `Limit` 与 `Offset`（均为 `*uint64`），无其它过滤条件。

## Goals / Non-Goals

**Goals:**
- 提供 `tencentcloud_cat_instant_tasks` 数据源，封装 `DescribeInstantTasks` 接口
- 将 `SingleInstantTask` 列表字段展开为 Terraform schema 的列表项字段（遵循"列表展开，不额外嵌套一层"的硬约束）
- 内部自动分页获取全部历史即时拨测任务（不向用户暴露 `limit`/`offset` 参数）
- 遵循现有 cat 数据源代码风格（参考 `data_source_tc_cat_node.go`、`data_source_tc_cat_metric_data.go`）
- 提供 `result_output_file` 参数用于导出结果
- 使用 `resource.Retry(tccommon.ReadRetryTimeout, ...)` + `tccommon.RetryError()` 处理临时性错误
- 在 Read 的 retry 块内对空响应返回 `NonRetryableError`（数据源硬约束），避免清空 id

**Non-Goals:**
- 不修改任何已有资源/数据源的 schema（纯新增）
- 不暴露 `limit`/`offset` 参数给用户（遵循 config.yaml 约束）
- 不实现写操作（数据源仅做查询）
- 不引入新的外部依赖

## Decisions

### Decision 1: 列表字段展开到 `tasks` 列表项的顶层（不额外嵌套）
**选择**: 将 `response.Response.Tasks` 映射为 schema 中的 `tasks`（`schema.TypeList`），每个列表元素的 schema 直接平铺 `task_id`/`target_address`/`task_type`/`probe_time`/`status`/`success_rate`/`node_count`/`task_category`，顶层另设 `total` 字段。

**理由**: 项目硬约束明确要求"资源参数 schema 中禁止创建该资源列表型数据这一层 schema"，应将列表展开，使每个字段都可被 Terraform 单独 set/read。参考 `data_source_tc_cat_node.go` 中 `node_define` 列表的处理方式。

**备选方案**: 把所有字段再包一层 `instant_task_set` —— 违反硬约束，已否决。

### Decision 2: service 层方法签名与分页策略
**选择**: 在 `service_tencentcloud_cat.go` 新增：
```go
func (me *CatService) DescribeCatInstantTasksByFilter(ctx context.Context) (tasks []*cat.SingleInstantTask, total *uint64, errRet error)
```
内部以固定 `pageSize`（取云 API 支持的最大值）循环分页，累积 `response.Response.Tasks`，直到本页返回数量小于 `pageSize` 即停止；同时返回最后一次响应的 `Total`。

**理由**:
- 该 API 无任何过滤参数，因此 service 方法无需接收过滤入参。
- 遵循 config.yaml "数据源分页：不暴露 limit/offset 给用户，内部实现自动分页"。
- 参考 `DescribeCatTaskSet` 中已有的分页循环写法（offset/pageSize 累加），保持风格一致。

**备选方案**: 在 data source 中直接调用 API 不抽 service 方法 —— 违反项目"通过 service 层调用云 API"的模式，已否决。

### Decision 3: retry 块对空响应返回 NonRetryableError
**选择**: 在 data source 的 Read 方法 retry 块内，调用 service 方法后判断 `tasks == nil && total == nil`（即 response/Response 均为空），返回 `tccommon.NonRetryableError`；不在此处 `d.SetId("")`。仅当 `len(tasks) == 0`（有响应但列表为空）时视为正常空结果，不返回错误，并打印 `log.Printf("[DATASOURCE] read empty, skip SetId")` 后继续。

**理由**: 遵循 RESOURCE_KIND_DATASOURCE 硬约束——"retry 块内必须检查云 API Read 接口是否返回了空，若返回了空，不要直接 `d.SetId("")`，而是直接返回 `NonRetryableError`"。区别"API 返回 nil 响应（异常）"与"API 返回空列表（正常）"两种情况。

### Decision 4: data source id 使用任务 id 哈希
**选择**: 收集所有 `task_id` 后使用 `helper.DataResourceIdsHash(ids)` 生成 data source 的 `d.SetId()`，与 `data_source_tc_cat_node.go` 一致。

**理由**: 复用项目通用工具，保持与同类数据源一致。

## Risks / Trade-offs

- **[任务数量极大时全量分页耗时]** → 通过固定较大 pageSize 减少 API 调用次数；并在 service 层 `ratelimit.Check()` 控制速率。历史即时拨测任务总量通常可控，风险较低。
- **[空响应误判清空 state]** → 已通过 Decision 3 在 retry 块内区分"异常空"与"正常空列表"，避免误清 id。
- **[字段类型转换]*uint64 → int / *float64 → float64** → Terraform schema 使用 `TypeInt`/`TypeFloat`，读取时 SDK 返回的是指针，需做 nil 判断后再 set，参考现有数据源的 nil 检查模式。
