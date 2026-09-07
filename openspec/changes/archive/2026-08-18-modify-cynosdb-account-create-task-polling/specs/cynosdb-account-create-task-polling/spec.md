## ADDED Requirements

### Requirement: Create 流程轮询异步任务完成状态

`tencentcloud_cynosdb_account` 资源的 Create 函数在调用 `CreateAccounts` 接口成功后，SHALL 检查返回的 `TaskId`。若 `TaskId` 非空且大于 0，SHALL 通过 `DescribeTasks` 接口轮询任务状态，直到 `Status` 为 `success` 后再执行 Read 操作。

#### Scenario: CreateAccounts 返回有效 TaskId 并轮询至成功

- **WHEN** `CreateAccounts` 返回非空且大于 0 的 `TaskId`
- **THEN** Create 函数 SHALL 使用 `DescribeTasks` 接口按 `TaskId` 轮询任务状态，直到 `Status` 为 `success`，再执行 Read 操作

#### Scenario: CreateAccounts 返回空 TaskId

- **WHEN** `CreateAccounts` 返回的 `TaskId` 为 nil 或为 0
- **THEN** Create 函数 SHALL 跳过任务轮询，直接设置 ID 并执行 Read 操作

#### Scenario: CreateAccounts 返回 nil Response

- **WHEN** `CreateAccounts` 返回 `result == nil` 或 `result.Response == nil` 或 `result.Response.TaskId == nil`
- **THEN** Create 函数 SHALL 返回 `NonRetryableError`，不设置资源 ID

#### Scenario: 异步任务进入失败状态

- **WHEN** `DescribeTasks` 轮询发现任务 `Status` 为失败状态（如 `failed`、`FAIL`）
- **THEN** Create 函数 SHALL 立即返回错误，不继续等待

#### Scenario: 异步任务轮询超时

- **WHEN** 任务在超时时间内未达到 `success` 状态
- **THEN** Create 函数 SHALL 返回超时错误

### Requirement: Read 流程有限重试

`tencentcloud_cynosdb_account` 资源的 Read 函数在调用 `DescribeCynosdbAccountById` 查不到账号时，SHALL 在 `tccommon.ReadRetryTimeout` 时间内进行有限重试，仅在重试耗尽后仍未找到账号时才清空资源 ID。

#### Scenario: Read 首次查不到账号后重试成功

- **WHEN** `DescribeCynosdbAccountById` 首次返回 `account == nil`
- **THEN** Read 函数 SHALL 在重试周期内继续查询，若后续查询到账号则正常设置字段

#### Scenario: Read 重试耗尽仍未找到账号

- **WHEN** `DescribeCynosdbAccountById` 在 `ReadRetryTimeout` 时间内持续返回 `account == nil`
- **THEN** Read 函数 SHALL 打印日志并执行 `d.SetId("")`

#### Scenario: Read 调用 API 返回错误

- **WHEN** `DescribeCynosdbAccountById` 返回 API 错误
- **THEN** Read 函数 SHALL 使用 `tccommon.RetryError(e)` 处理错误，可重试错误继续重试，不可重试错误直接返回

### Requirement: 向后兼容

变更 SHALL 保持现有 Schema 定义不变，不新增、删除或修改任何 Schema 字段。现有 Terraform 配置和 state 结构 SHALL 不受影响。

#### Scenario: 现有配置无需修改

- **WHEN** 用户使用现有的 `tencentcloud_cynosdb_account` 资源配置执行 `terraform apply`
- **THEN** 资源 SHALL 正常创建，无需用户修改任何配置

#### Scenario: 现有 state 兼容

- **WHEN** Terraform 加载已有的 `tencentcloud_cynosdb_account` state 执行 refresh
- **THEN** Read 函数 SHALL 正常工作，不会因变更导致 state 不兼容
