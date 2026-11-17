## 1. Service Layer（service 层）

- [x] 1.1 在 `tencentcloud/services/cat/service_tencentcloud_cat.go` 新增 `DescribeCatInstantTasksByFilter` 方法，封装 `DescribeInstantTasks` 接口
  - 方法签名：`func (me *CatService) DescribeCatInstantTasksByFilter(ctx context.Context) (tasks []*cat.SingleInstantTask, total *uint64, errRet error)`
  - 使用固定 `pageSize`（云 API 支持的最大值）循环分页，累积 `response.Response.Tasks`，直到本页返回数量小于 `pageSize` 即停止
  - 每次请求前调用 `ratelimit.Check(request.GetAction())`
  - 每次请求后记录 `[DEBUG]` 级别的请求体与响应体日志
  - 失败时记录 `[CRITAL]` 错误日志（含请求体与错误原因）并返回错误
  - 返回累积后的全部 `tasks` 与最后一次响应的 `total`
- [x] 1.2 处理空列表场景：当 API 返回非 nil 响应但 `response.Response.Tasks` 为空列表时，返回空切片与对应的 `total`，不返回错误

## 2. Data Source 实现（schema 定义与 Read 函数）

- [x] 2.1 创建 `tencentcloud/services/cat/data_source_tc_cat_instant_tasks.go`
  - 声明 `package cat` 与必要的 import（`context`、`log`、`resource`、`schema`、`tccommon`、`cat`、`helper`）
- [x] 2.2 定义 `DataSourceTencentCloudCatInstantTasks() *schema.Resource` 的 schema：
  - 输出 `tasks`（`schema.TypeList`，Computed），每个元素 schema 平铺以下字段（均 Computed，遵循列表展开硬约束，不额外嵌套一层）：
    - `task_id`（TypeString）、`target_address`（TypeString）、`task_type`（TypeInt）、`probe_time`（TypeInt）、`status`（TypeString）、`success_rate`（TypeFloat）、`node_count`（TypeInt）、`task_category`（TypeInt）
  - 输出 `total`（TypeInt，Computed）
  - 输入 `result_output_file`（TypeString，Optional）
  - 不向用户暴露 `limit`/`offset` 参数（遵循 config.yaml 约束）
- [x] 2.3 实现 `dataSourceTencentCloudCatInstantTasksRead` Read 函数：
  - `defer tccommon.LogElapsed("data_source.tencentcloud_cat_instant_tasks.read")()`
  - `defer tccommon.InconsistentCheck(d, meta)()`
  - 使用 `tccommon.GetLogId(tccommon.ContextNil)` 与 `context.WithValue` 构造 ctx
  - 调用 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 包裹 service 方法 `DescribeCatInstantTasksByFilter(ctx)`
- [x] 2.4 retry 块内错误与空响应处理（遵循数据源硬约束）：
  - service 方法返回错误时使用 `tccommon.RetryError(e)` 包装返回
  - 云 API 返回空（`tasks == nil && total == nil`，即 response/Response 均为空）时直接返回 `NonRetryableError`，不调用 `d.SetId("")`
  - 列表为空（`len(tasks) == 0` 但响应正常）时打印 `log.Printf("[DATASOURCE] read empty, skip SetId")` 后继续，不返回错误
- [x] 2.5 将 API 响应映射到 Terraform state：
  - 遍历 `tasks`，为每个 `SingleInstantTask` 构造 map，set 各字段前先做 nil 判断（遵循 Read 方法硬约束）
  - `task_id`/`target_address`/`status` 为 `*string` → string
  - `task_type`/`probe_time`/`node_count`/`task_category` 为 `*uint64` → int
  - `success_rate` 为 `*float64` → float64
  - 使用 `helper.DataResourceIdsHash(ids)`（收集所有 `task_id`）设置 `d.SetId()`
  - set `tasks` 列表与 `total` 到 ResourceData

## 3. 结果导出

- [x] 3.1 在 Read 函数末尾处理 `result_output_file`：若用户指定，使用 `tccommon.WriteToFile(output.(string), tasksList)` 导出结果，失败时返回错误

## 4. Provider 注册

- [x] 4.1 在 `tencentcloud/provider.go` 的 `DataSourcesMap` 中添加注册代码：
  ```go
  "tencentcloud_cat_instant_tasks": cat.DataSourceTencentCloudCatInstantTasks(),
  ```
- [x] 4.2 确认注册命名遵循约定：`tencentcloud_cat_instant_tasks`

## 5. 文档

- [x] 5.1 创建 `tencentcloud/services/cat/data_source_tc_cat_instant_tasks.md`：
  - 一句话描述，带上所属云产品名称（CAT / 云拨测），格式为 "Use this data source to query ..."
  - Example Usage 部分提供基本查询示例
  - 不手写 `Argument Reference` 与 `Attribute Reference`（由工具自动生成）
  - RESOURCE_KIND_DATASOURCE 类型无需 Import 部分
- [ ] 5.2 运行 `make doc` 自动生成 `website/docs/` 文档与 `provider.md` 条目（由 tfpacer-finalize 阶段执行）

## 6. 单元测试

- [x] 6.1 创建 `tencentcloud/services/cat/data_source_tc_cat_instant_tasks_test.go`：
  - 使用 mock（gomonkey）方式对云 API 进行 mock，不使用 terraform 测试套件（遵循新增 terraform 资源测试硬约束）
  - 覆盖正常查询场景：mock `DescribeInstantTasks` 返回非空任务列表，验证字段映射与 `total` 设置、`d.Id()` 设置
  - 覆盖空列表场景：mock 返回空列表，验证不报错、正常设置 id
  - 覆盖 API 错误场景：mock 返回错误，验证错误被返回
  - 禁止通过 `go test` 命令执行单元测试，但需保证生成的代码在当前环境下可正确构建执行

## 7. 代码正确性检查与验证（收尾阶段）

- [x] 7.1 检查所有函数返回的 error 是否被检查；必定不出错的函数使用 `_ = func()` 将 err 赋值给 `_`
- [x] 7.2 检查 service 层调用云 API 的超时时间是否使用 `tccommon.ReadRetryTimeout`，错误是否使用 `tccommon.RetryError()` 包装
- [x] 7.3 验证字段映射与 vendor 中 `SingleInstantTask`/`DescribeInstantTasksResponseParams` 一致（TaskId/TargetAddress/TaskType/ProbeTime/Status/SuccessRate/NodeCount/TaskCategory/Tasks/Total）
- [ ] 7.4 在收尾阶段通过 tfpacer-finalize skill 执行 `gofmt`、`make doc`、生成 `.changelog` 文件并 amend 推送
