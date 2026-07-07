## Context

`tencentcloud_postgresql_instance` 是管理腾讯云 PostgreSQL 实例全生命周期的 RESOURCE_KIND_GENERAL 资源。其中 `backup_plan` 是一个 `TypeList`（`MaxItems: 1`）嵌套块，用于管理实例的备份计划。

现有的 `backup_plan` CRUD 逻辑如下：
- **Create**：若用户配置了 `monthly_backup_period`，先调用 `CreateBackupPlan`（`BackupPeriodType=month`）创建月度备份计划；随后始终调用 `ModifyBackupPlan` 修改默认（week）备份计划，将 `min_backup_start_time`、`max_backup_start_time`、`base_backup_retention_period`、`backup_period` 等字段传入。
- **Read**：调用 `DescribeBackupPlans` 获取备份计划列表，按 `BackupPeriodType` 区分 `week`（默认计划）与 `month`（自定义月度计划），将各字段回填到 `backup_plan` 块。
- **Update**：当 `d.HasChange("backup_plan")` 时，调用 `ModifyBackupPlan` 更新默认计划，并对月度计划执行 create/modify/delete 分支处理。

腾讯云 `ModifyBackupPlan` 接口（`postgres/v20170312`）近期新增了 `BackupMethod` 入参（枚举值：`physical` 物理备份、`logical` 逻辑备份、`snapshot` 快照备份），同时 `BackupPlan` 响应结构也已包含 `BackupMethod` 字段。本次变更将该参数透传到 Terraform 资源。

vendor 中的 SDK 已更新到位（`ModifyBackupPlanRequest.BackupMethod`、`BackupPlan.BackupMethod` 均已存在），无需额外升级 SDK。

## Goals / Non-Goals

**Goals:**
- 在 `backup_plan` 嵌套块中新增 `backup_method` 参数，使用户可指定备份方式。
- 在 Create/Update 流程中将 `backup_method` 通过 `ModifyBackupPlanRequest.BackupMethod` 传递给云API。
- 在 Read 流程中从 `BackupPlan.BackupMethod` 回填到 state。
- 保持向后兼容：参数为 Optional + Computed，不影响现有配置。

**Non-Goals:**
- 不修改 `CreateBackupPlan` 接口的调用（该接口不支持 `BackupMethod`，仅用于创建 month 类型计划）。
- 不改变 `backup_plan` 块的其他字段语义。
- 不调整月度备份计划（month）的逻辑，`backup_method` 仅作用于通过 `ModifyBackupPlan` 修改的默认（week）备份计划。
- 不修改 provider.go 注册（资源已注册）。

## Decisions

### 决策 1：`backup_method` 作为 `backup_plan` 嵌套块字段，而非顶层字段
**选择**：将 `backup_method` 放入 `backup_plan` 块内。
**理由**：`BackupMethod` 是备份计划的属性，与 `min_backup_start_time`、`backup_period` 等同级，语义上属于备份计划配置。云API `ModifyBackupPlan` 也是在备份计划维度设置该参数。放入顶层会破坏与现有 `backup_plan` 字段的一致性。
**备选**：作为顶层字段——已否决，语义不符。

### 决策 2：`backup_method` 标记为 Optional + Computed
**选择**：`Optional: true, Computed: true`。
**理由**：用户可不指定该参数（由云API使用默认备份方式），此时 Read 时需从 API 回填实际值，故 Computed。这符合"不能修改已有资源 schema（除非只新增 Optional 字段）"的硬约束，保证向后兼容。

### 决策 3：仅在 `ModifyBackupPlan` 调用中传递 `BackupMethod`
**选择**：在 Create 和 Update 流程中调用 `ModifyBackupPlan`（修改默认 week 计划）时设置 `request.BackupMethod`；不为月度计划的 `CreateBackupPlan`/`ModifyBackupPlan` 设置该参数。
**理由**：`CreateBackupPlanRequest` 不支持 `BackupMethod` 字段；需求明确指定 `BackupMethod` 为 `ModifyBackupPlan` 接口的新增入参。现有代码中默认计划的修改均通过 `ModifyBackupPlanRequest` 完成，在此处透传即可覆盖主要场景。月度计划保持原有逻辑不变，避免引入不一致。
**备选**：同时在月度计划的 `ModifyBackupPlanRequest`（`request1`）中设置——暂不采用，避免扩大改动范围，聚焦需求中明确的单一参数。

### 决策 4：Read 时回填 `backup_method`
**选择**：在 Read 流程中，当 `backupPlan`（week 类型）的 `BackupMethod` 非 nil 时，写入 `planMap["backup_method"]`。
**理由**：遵循现有 nil 安全模式（先判断非 nil 再 set），与同块内其他字段处理方式一致。

## Risks / Trade-offs

- [风险] 用户设置了不被实例版本/规格支持的 `BackupMethod` 枚举值 → 由云API校验并返回错误，Terraform 透传错误信息；schema 层不强制枚举校验，以兼容未来云API可能新增的枚举值。
- [风险] 旧版本 state 中无 `backup_method` 字段 → 由于该字段为 Optional+Computed，Terraform 在 Read 后会自动回填，不影响现有资源。
- [权衡] 仅对默认（week）计划设置 `BackupMethod`，月度计划不设置 → 月度计划的备份方式由云API默认值决定；如后续云API在 `CreateBackupPlan` 支持 `BackupMethod`，可再扩展。
