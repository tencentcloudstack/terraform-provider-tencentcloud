## Context

CAT（云拨测）服务当前已在 Provider 中支持拨测任务的创建与管理（`tencentcloud_cat_task_set` 资源，对应 `DescribeProbeTasks` / `DeleteProbeTask` 等接口）。但缺少一个面向列表查询的 Data Source，用户无法在 Terraform 中按条件读取已有拨测任务信息。

现有 cat service 层已封装了 `DescribeCatTaskSet`（按单个 taskId 查询）方法，但它仅返回单条任务，无法满足按多条件过滤、分页获取完整列表的需求。本次新增 Data Source `tencentcloud_cat_probe_tasks` 复用同一 `DescribeProbeTasks` 云 API，但需要一个新的 service 层方法支持所有可选过滤参数并将全部页数据累积返回。

云 API `DescribeProbeTasks` 入参包含 `Offset`/`Limit`（最大值 100），且支持 `TaskIDs`/`TaskName`/`TargetAddress`/`TaskStatus`/`PayMode`/`OrderState`/`TaskType`/`TaskCategory`/`OrderBy`/`Ascend`/`TagFilters` 等过滤参数；出参为 `TaskSet`（`[]*ProbeTask`）与 `Total`。其中 `TaskStatus`/`TaskType`/`TaskCategory` 在入参中为 `[]*int64`（数组），而在出参 `ProbeTask` 中为 `*int64`（单值）。`TagFilters` 入参为 `[]*KeyValuePair`（Key/Value）。

## Goals / Non-Goals

**Goals:**
- 新增 Data Source `tencentcloud_cat_probe_tasks`，支持所有 `DescribeProbeTasks` 入参作为可选过滤条件
- service 层新增 `DescribeCatProbeTasksByFilter` 方法，内部自动分页（pageSize 取云 API 最大值 100），累积全部任务返回
- 字段映射严格遵循接口与参数映射关系，列表项字段平铺在 `task_set` 列表元素内（嵌套的 `tag_info_list` 作为子列表）
- Read 方法使用 retry 包裹并对空响应做安全处理（nil 响应返回 NonRetryableError，空列表正常返回）
- 提供基于 gomonkey 的 mock 单元测试

**Non-Goals:**
- 不修改现有 `tencentcloud_cat_task_set` 资源及其 service 方法
- 不向用户暴露 `limit`/`offset` 分页参数
- 不新增异步接口轮询逻辑（`DescribeProbeTasks` 为同步查询接口）

## Decisions

### 1. Service 层方法签名与分页策略

新增方法：
```go
func (me *CatService) DescribeCatProbeTasksByFilter(ctx context.Context, param map[string]interface{}) (tasks []*cat.ProbeTask, total *int64, errRet error)
```

- 采用固定 `pageSize = 100`（云 API 注释标注的最大值）循环分页，累积 `response.Response.TaskSet`，直到本页返回数量小于 `pageSize` 即停止。
- 不向用户暴露 `limit`/`offset`，分页在 service 层内部完成，与现有 `DescribeCatInstantTasksByFilter` / `DescribeCatTaskSet` 保持一致风格。
- 参数通过 `param map[string]interface{}` 传入，在方法内按 key 赋值到 request 字段，便于扩展。

**理由**: 复用现有 cat service 的分页模式，保持代码风格一致；map 传参与 igtm 等现有数据源一致。

### 2. 入参类型处理

| Terraform Schema 字段 | Schema 类型 | 云 API 字段 | 云 API 类型 | 转换说明 |
|---|---|---|---|---|
| `task_i_ds` | TypeList(string) | `TaskIDs` | `[]*string` | 列表转指针切片 |
| `task_name` | TypeString | `TaskName` | `*string` | 直接转换 |
| `target_address` | TypeString | `TargetAddress` | `*string` | 直接转换 |
| `task_status` | TypeList(int) | `TaskStatus` | `[]*int64` | 列表转指针切片 |
| `pay_mode` | TypeInt | `PayMode` | `*int64` | 直接转换 |
| `order_state` | TypeInt | `OrderState` | `*int64` | 直接转换 |
| `task_type` | TypeList(int) | `TaskType` | `[]*int64` | 列表转指针切片 |
| `task_category` | TypeList(int) | `TaskCategory` | `[]*int64` | 列表转指针切片 |
| `order_by` | TypeString | `OrderBy` | `*string` | 直接转换 |
| `ascend` | TypeBool | `Ascend` | `*bool` | 直接转换 |
| `tag_filters` | TypeList(block) | `TagFilters` | `[]*KeyValuePair` | 块含 key/value 转 KeyValuePair 切片 |

`task_status`/`task_type`/`task_category` 入参为切片类型（支持多值过滤），因此 schema 使用 TypeList(int)。`tag_filters` 为含 `key`/`value`（均必填）的 block 列表，映射到 `KeyValuePair`。

### 3. 出参字段映射（task_set 列表元素）

`task_set` 为 TypeList，每个元素为 Resource，字段平铺（不额外嵌套"列表型数据"一层），嵌套的 `tag_info_list` 作为子 TypeList：

| Terraform 字段 | Schema 类型 | 云 API 字段(ProbeTask) |
|---|---|---|
| `name` | TypeString | `Name` |
| `task_id` | TypeString | `TaskId` |
| `task_type` | TypeInt | `TaskType` |
| `nodes` | TypeList(string) | `Nodes` |
| `node_ip_type` | TypeInt | `NodeIpType` |
| `interval` | TypeInt | `Interval` |
| `parameters` | TypeString | `Parameters` |
| `status` | TypeInt | `Status` |
| `target_address` | TypeString | `TargetAddress` |
| `pay_mode` | TypeInt | `PayMode` |
| `order_state` | TypeInt | `OrderState` |
| `task_category` | TypeInt | `TaskCategory` |
| `created_at` | TypeString | `CreatedAt` |
| `cron` | TypeString | `Cron` |
| `cron_state` | TypeInt | `CronState` |
| `tag_info_list` | TypeList(block: key/value) | `TagInfoList` ([]*KeyValuePair) |
| `sub_sync_flag` | TypeInt | `SubSyncFlag` |

`total` 为顶层 TypeInt（映射 `response.Response.Total`）。

### 4. Read 方法 retry 与空响应处理

遵循数据源规范：
- retry 块内若 service 返回错误，用 `tccommon.RetryError(e)` 包装
- 若云 API 返回空（`response == nil` / `response.Response == nil`，表现为 `tasks == nil && total == nil`），直接返回 `NonRetryableError`，不调用 `d.SetId("")`，保留现场
- 若返回正常但列表为空（`len(tasks) == 0`），打印 `log.Printf("[DATASOURCE] read empty, skip SetId")` 后继续，不报错
- `d.SetId(helper.BuildToken())` 设置 id（与 igtm 数据源一致）

### 5. 测试策略

使用 gomonkey mock `DescribeProbeTasksWithContext`，不使用 terraform 测试套件，覆盖：
- 正常查询（mock 返回非空任务列表，验证字段映射与 total 设置）
- 空列表场景（mock 返回空列表，验证不报错、正常设置 id）
- API 错误场景（mock 返回错误，验证错误被返回）

## Risks / Trade-offs

- **[入参切片字段语义差异]** → `TaskStatus`/`TaskType`/`TaskCategory` 入参为切片（多值过滤）但出参为单值，schema 中分别用 TypeList(int) 入参与 TypeInt 出参，由于 Terraform 数据源 Read 不回写入参（入参由用户配置），二者互不干扰，无冲突。
- **[TagFilters Key/Value 必填]** → 映射要求 `tag_filters` 块内 `key`/`value` 必填，与云 API `KeyValuePair` 结构一致；若用户不传 `tag_filters` 则不设置该 request 字段。
- **[分页累积性能]** → 当任务数量极大时全量累积可能较慢，但符合数据源"不暴露分页给用户"的约定，且 `Limit` 最大值 100 已尽量减少请求次数。
