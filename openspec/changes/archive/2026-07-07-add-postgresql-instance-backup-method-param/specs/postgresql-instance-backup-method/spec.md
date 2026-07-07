## ADDED Requirements

### Requirement: backup_plan 块新增 backup_method 字段

`tencentcloud_postgresql_instance` 资源的 `backup_plan` 嵌套块 SHALL 包含 `backup_method` 字段，用于指定备份计划采用的备份方式。

#### Scenario: Schema 定义 backup_method 字段

- **假设** 定义 `tencentcloud_postgresql_instance` 资源的 `backup_plan` 嵌套块 schema
- **当** schema 定义完成时
- **那么** `backup_plan` 块 SHALL 包含 `backup_method` 字段，具备以下属性：
  - 类型: `schema.TypeString`
  - Optional: `true`
  - Computed: `true`
  - Description 说明备份方式枚举值（`physical` 物理备份、`logical` 逻辑备份、`snapshot` 快照备份）

### Requirement: Create 流程传递 BackupMethod

资源 Create 流程在调用 `ModifyBackupPlan` 修改默认备份计划时，SHALL 将用户配置的 `backup_method` 传入 `ModifyBackupPlanRequest.BackupMethod`。

#### Scenario: 创建时设置备份方式

- **假设** 用户在 `backup_plan` 块中配置了 `backup_method = "logical"`
- **当** 执行资源 Create 流程，调用 `ModifyBackupPlan` 修改默认（week）备份计划时
- **那么** SHALL 将 `request.BackupMethod` 设置为 `"logical"` 传递给云API

#### Scenario: 创建时未配置备份方式

- **假设** 用户未在 `backup_plan` 块中配置 `backup_method`
- **当** 执行资源 Create 流程，调用 `ModifyBackupPlan` 修改默认（week）备份计划时
- **那么** SHALL 不设置 `request.BackupMethod`，由云API使用默认备份方式

### Requirement: Update 流程传递 BackupMethod

资源 Update 流程在 `backup_plan` 发生变更调用 `ModifyBackupPlan` 时，SHALL 将用户配置的 `backup_method` 传入 `ModifyBackupPlanRequest.BackupMethod`。

#### Scenario: 更新时修改备份方式

- **假设** 存在一个已创建的 postgresql instance 资源
- **当** 用户修改 `backup_plan` 块中的 `backup_method` 值并执行 Update 流程时
- **那么** 调用 `ModifyBackupPlan` 修改默认（week）备份计划时，SHALL 将新的 `backup_method` 值传入 `request.BackupMethod`

#### Scenario: 更新时未配置备份方式

- **假设** 存在一个已创建的 postgresql instance 资源
- **当** 用户修改 `backup_plan` 其他字段但未配置 `backup_method` 并执行 Update 流程时
- **那么** 调用 `ModifyBackupPlan` 时 SHALL 不设置 `request.BackupMethod`

### Requirement: Read 流程回填 backup_method

资源 Read 流程 SHALL 从 `DescribeBackupPlans` 返回的默认（week）备份计划中读取 `BackupMethod` 字段并回填到 state。

#### Scenario: 读取并回填备份方式

- **假设** 调用 `DescribeBackupPlans` 返回的默认（week）备份计划中 `BackupMethod` 字段非 nil
- **当** 执行资源 Read 流程，构造 `backup_plan` 的 `planMap` 时
- **那么** SHALL 将 `backupPlan.BackupMethod` 写入 `planMap["backup_method"]` 并通过 `d.Set("backup_plan", ...)` 回填到 state

#### Scenario: 读取时备份方式为空不回填

- **假设** 调用 `DescribeBackupPlans` 返回的默认（week）备份计划中 `BackupMethod` 字段为 nil
- **当** 执行资源 Read 流程，构造 `backup_plan` 的 `planMap` 时
- **那么** SHALL 不设置 `planMap["backup_method"]`，保持 nil 安全

### Requirement: 资源文档更新

资源文档 `resource_tc_postgresql_instance.md` SHALL 包含 `backup_method` 参数的说明与示例。

#### Scenario: 文档包含 backup_method 说明

- **假设** 存在资源文档文件 `tencentcloud/services/postgresql/resource_tc_postgresql_instance.md`
- **当** 更新文档时
- **那么** 文档 SHALL 在 `backup_plan` 块示例中包含 `backup_method` 参数的用法，并说明其枚举值含义
