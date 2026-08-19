## Context

`tencentcloud_cynosdb_account` 资源的 Create 函数调用 `CreateAccounts` 接口后，该接口返回 `TaskId` 表示这是一个异步任务。当前代码没有处理 `TaskId`，而是在调用成功后直接 `d.SetId(...)` 并执行 Read。但由于异步任务尚未完成，Read 时 `DescribeAccounts` 接口可能查不到刚创建的账号，导致 `d.SetId("")` 清空 ID，Terraform 误认为创建失败。

cynosdb 服务中已存在两种成熟的任务轮询模式：
1. `resource.Retry` + `DescribeTasks` 检查 Status 是否为 `success`（见 `resource_tc_cynosdb_cls_delivery.go`）
2. `tccommon.BuildStateChangeConf` + `service.taskStateRefreshFunc`（见 `resource_tc_cynosdb_ssl.go`、`resource_tc_cynosdb_cluster_transparent_encrypt.go`）

`CynosdbService.taskStateRefreshFunc` 方法已封装了通过 `DescribeTasks` 接口按 `TaskId` 过滤查询任务状态的逻辑，返回任务 `Status` 字符串。

云 API 验证：
- `CreateAccountsResponse` 返回 `TaskId (*int64)` 字段
- `DescribeTasksRequest` 支持通过 `Filters` 中的 `TaskId` 字段精确过滤
- `DescribeTasksResponse` 返回 `TaskList []*BizTaskInfo`，其中 `BizTaskInfo.Status (*string)` 为任务状态，`success` 表示完成

## Goals / Non-Goals

**Goals:**
- 在 Create 流程中，`CreateAccounts` 返回 `TaskId` 后，轮询 `DescribeTasks` 直到任务 Status 为 `success`，再执行 Read
- 在 Read 流程中，当账号查不到时进行有限重试，避免因最终一致性延迟导致 ID 被误清空
- 保持向后兼容，不修改 Schema 定义

**Non-Goals:**
- 不修改 Update/Delete 流程（它们使用同步接口，不涉及异步任务）
- 不新增 Schema 字段或 Timeouts 块
- 不修改 `CynosdbService.taskStateRefreshFunc` 方法

## Decisions

### 决策 1：Create 中使用 `BuildStateChangeConf` + `taskStateRefreshFunc` 轮询任务

选择 `tccommon.BuildStateChangeConf` + `service.taskStateRefreshFunc` 模式，而非 `resource.Retry` 模式。

**理由：**
- `taskStateRefreshFunc` 已存在于 `CynosdbService` 中，专门封装了 `DescribeTasks` 按 `TaskId` 过滤查询的逻辑，复用最为自然
- `BuildStateChangeConf` 是 Terraform Plugin SDK 推荐的异步任务等待模式，支持设定 pending/target 状态、超时时间、轮询间隔
- 与 `resource_tc_cynosdb_ssl.go`、`resource_tc_cynosdb_cluster_transparent_encrypt.go` 保持一致

**备选方案：** 使用 `resource.Retry` + 内联 `DescribeTasks` 调用（如 `cls_delivery.go`）。此方案需要内联编写 `DescribeTasks` 查询逻辑，与已有 `taskStateRefreshFunc` 重复，不优先。

### 决策 2：Create 中 `CreateAccounts` 返回值检查

在 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 块内，调用 `CreateAccounts` 成功后检查返回值：
- 若 `result == nil || result.Response == nil || result.Response.TaskId == nil`，返回 `NonRetryableError`，避免后续逻辑写入空 TaskId
- 若检查通过，保存 TaskId 到外层变量

### 决策 3：TaskId 为 0 时的处理

`CreateAccounts` 可能不总是返回有效的 `TaskId`（某些场景可能直接同步完成）。若 `TaskId` 为 nil 或为 0，跳过任务轮询，直接执行 Read，保持向后兼容。

### 决策 4：Read 中有限重试逻辑

在 `resourceTencentCloudCynosdbAccountRead` 中，使用 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 包裹 `DescribeCynosdbAccountById` 调用：
- 若账号查不到（返回 `account == nil`），返回 `RetryableError` 继续重试
- 若 API 返回错误，使用 `tccommon.RetryError(e)` 处理
- 重试耗尽后仍为 nil，打印日志并 `d.SetId("")`

**备选方案：** 不修改 Read，仅在 Create 中轮询。但 Read 可能被 Terraform 在其他时机调用（如 Import、refresh），添加重试可增强鲁棒性。

### 决策 5：任务失败状态处理

在 `BuildStateChangeConf` 中，将 `failStates` 设置为常见的失败状态（如 `failed`、`FAIL`），当任务进入失败状态时立即返回错误，避免无效等待直到超时。

## Risks / Trade-offs

- **[风险] TaskId 可能为 nil 或 0** → 在调用 `BuildStateChangeConf` 前判断 `TaskId` 是否有效，若无效则跳过轮询直接 Read
- **[风险] 任务轮询增加 Create 耗时** → 设置合理的超时时间 `10 * tccommon.ReadRetryTimeout`，与同类资源（ssl、transparent_encrypt）保持一致
- **[风险] Read 重试可能掩盖真正的资源删除** → 仅在 Create 流程中调用 Read 时重试有意义；但为保持代码统一，Read 函数统一使用重试逻辑，重试耗尽后仍会 `d.SetId("")`
- **[风险] `taskStateRefreshFunc` 返回 nil 时 WaitForState 行为** → 当 `DescribeTasks` 返回空列表时，`taskStateRefreshFunc` 返回 `(nil, "", nil)`，`BuildStateChangeConf` 会将其视为 pending 状态继续轮询，符合预期
