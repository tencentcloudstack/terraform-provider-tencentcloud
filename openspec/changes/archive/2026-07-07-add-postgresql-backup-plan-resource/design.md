## Context

TencentCloud PostgreSQL 提供自定义备份计划（BackupPlan）能力。云API `postgres/v20170312` 已提供四个接口：
- `CreateBackupPlan`：创建备份计划，返回 `PlanId`
- `DescribeBackupPlans`：按 `DBInstanceId` 查询实例下所有备份计划，返回 `Plans []*BackupPlan`
- `ModifyBackupPlan`：修改备份计划（不含 `BackupPeriodType`）
- `DeleteBackupPlan`：删除备份计划

当前 provider 中已存在配置型资源 `tencentcloud_postgresql_backup_plan_config`，它只调用 `ModifyBackupPlan` 修改默认备份计划、Create/Delete 为空操作，无法创建/删除自定义备份计划。本变更新增一个完整的 CRUD 资源 `tencentcloud_postgresql_backup_plan`，与已有资源互不冲突。

约束：参考 `tencentcloud_igtm_strategy` 资源的代码风格；vendor 模式管理依赖（接口已在 vendor 中）；使用联合 ID；异步操作接口需轮询（经核对，CreateBackupPlan 等接口为同步返回 PlanId，无需额外轮询任务接口）。

## Goals / Non-Goals

**Goals:**
- 提供完整的 `tencentcloud_postgresql_backup_plan` CRUD 资源，管理 PostgreSQL 自定义备份计划全生命周期。
- 使用联合 ID（`db_instance_id{FILED_SP}plan_id`）支持资源导入与定位。
- 遵循项目代码风格（参考 `resource_tc_igtm_strategy.go`）：`resource.Retry` 包装、nil 检查、日志规范。
- 在 `provider.go` 注册资源，并生成文档与单元测试。

**Non-Goals:**
- 不修改已有的 `tencentcloud_postgresql_backup_plan_config` 资源。
- 不暴露 `DescribeBackupPlans` 的列表型嵌套结构（本资源是单资源 CRUD，非 datasource）。
- 不实现异步任务轮询（`CreateBackupPlan` 同步返回 `PlanId`，无 `TaskId`）。

## Decisions

### 1. 资源 ID 采用联合 ID

**决策**：使用 `db_instance_id{FILED_SP}plan_id` 作为资源 ID。

**理由**：`DescribeBackupPlans` 仅按 `DBInstanceId` 查询，返回该实例下所有备份计划，无法仅凭 `PlanId` 定位实例。Read/Delete/Update 都需要同时知道 `DBInstanceId` 与 `PlanId`。联合 ID 是项目标准模式（见 `tencentcloud_igtm_strategy` 的 `instanceId{FILED_SP}strategyId`）。

**实现**：Create 成功后 `d.SetId(strings.Join([]string{dbInstanceId, planId}, tccommon.FILED_SP))`；Read/Update/Delete 中 `strings.Split(d.Id(), tccommon.FILED_SP)` 解析为两段。

### 2. `backup_period_type` 标记为 ForceNew

**决策**：`backup_period_type`（备份周期类型，当前仅支持 `month`）设为 ForceNew。

**理由**：云API `ModifyBackupPlan` 入参不包含 `BackupPeriodType`，无法更新该字段。按项目规范，创建后不可变的字段需 ForceNew。`db_instance_id` 同样 ForceNew（实例切换本质是不同实例）。

### 3. `backup_period` 字段的读写不一致处理

**决策**：
- 写入（Create/Modify）：`backup_period` schema 为 `TypeSet(TypeString)`，写入时遍历 Set 转为 `[]*string` 传给云API。
- 读取（Read）：`BackupPlan.BackupPeriod` 在出参结构体中为 `*string`（一个 JSON 字符串，如 `["1","2"]`），需 `json.Unmarshal` 反序列化为切片后再 `d.Set("backup_period", ...)`。

**理由**：与已有 `tencentcloud_postgresql_backup_plan_config` 的 Read 处理方式保持一致（它同样对 `BackupPeriod` 做 `json.Unmarshal`）。云API出参用字符串承载列表是服务端约定。

### 4. Read 中匹配 PlanId

**决策**：`DescribeBackupPlans` 返回 `Plans` 列表，在 Read 中遍历列表匹配 `PlanId == 解析出的 plan_id` 的元素作为当前资源数据。

**理由**：`DescribeBackupPlans` 不支持按 `PlanId` 过滤，只能按 `DBInstanceId` 查询全部计划。本资源是单资源 CRUD，需从列表中精确定位到目标 `PlanId`。

### 5. 不生成 `_extension.go` 文件

**决策**：按项目规范"若非必须，则不要生成 `_extension.go` 文件"，CRUD 逻辑直接在资源文件与服务层中实现，不生成扩展文件。

### 6. 单元测试使用 gomonkey mock

**决策**：新增资源使用 gomonkey 对云API client 方法打桩，编写 `*_test.go` 单元测试，用 `go test -gcflags=all=-l` 运行，不使用 terraform 测试套件。

**理由**：项目规范要求"对于新增的 terraform 资源，使用 mock（gomonkey）的方法对云API进行 mock 处理，只进行业务代码逻辑的单元测试"。

## Risks / Trade-offs

- **[Read 列表匹配可能因服务端返回顺序变化而定位失败]** → 通过精确匹配 `PlanId` 而非索引定位来规避；未匹配到时按规范打印 `[CRUD]` 日志后 `d.SetId("")`。
- **[BackupPeriod 出参为字符串、入参为列表，类型不对称]** → Read 时统一 `json.Unmarshal` 处理，与已有 `backup_plan_config` 一致，避免状态漂移。
- **[CreateBackupPlan 是否为异步接口未明确]** → 经核对 SDK，`CreateBackupPlanResponse` 直接返回 `PlanId`，无 `TaskId`/异步任务句柄，视为同步接口，Create 后直接 SetId 并调用 Read，无需轮询任务接口。
- **[与已有 `backup_plan_config` 资源语义重叠]** → 两者定位不同：`backup_plan_config` 管理默认备份计划（无真实创建/删除），新资源管理自定义备份计划全生命周期；资源名不同，state 不冲突。在文档中说明差异即可。
