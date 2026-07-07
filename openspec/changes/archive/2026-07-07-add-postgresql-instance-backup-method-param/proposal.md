## Why

腾讯云 PostgreSQL 的 `ModifyBackupPlan` 云API新增了 `BackupMethod` 入参，支持配置备份方式（物理备份 `physical`、逻辑备份 `logical`、快照备份 `snapshot`）。当前 `tencentcloud_postgresql_instance` 资源的 `backup_plan` 块未暴露该参数，用户无法通过 Terraform 指定备份计划采用的备份方式，导致无法使用云API最新提供的能力。本次变更新增该参数，使 Terraform 资源与云API能力保持同步。

## What Changes

- 在 `tencentcloud_postgresql_instance` 资源的 `backup_plan` 嵌套块 schema 中新增 `backup_method` 字段（`schema.TypeString`，Optional，Computed）。
- 在资源 Create 流程中，调用 `ModifyBackupPlan` 设置默认备份计划时，将用户配置的 `backup_method` 传入 `ModifyBackupPlanRequest.BackupMethod`。
- 在资源 Update 流程中，当 `backup_plan` 发生变更时，将 `backup_method` 传入 `ModifyBackupPlanRequest.BackupMethod`。
- 在资源 Read 流程中，从 `DescribeBackupPlans` 返回的 `BackupPlan.BackupMethod` 读取并回填到 state。
- 更新资源文档 `resource_tc_postgresql_instance.md`，补充 `backup_method` 参数说明与示例。

## Capabilities

### New Capabilities
- `postgresql-instance-backup-method`: 为 `tencentcloud_postgresql_instance` 资源的 `backup_plan` 块新增 `backup_method` 参数，支持配置备份方式（physical/logical/snapshot），覆盖 schema 定义、Create/Update/Read 流程及文档。

### Modified Capabilities
<!-- 无现有 spec 需要修改 -->

## Impact

- **受影响代码**:
  - `tencentcloud/services/postgresql/resource_tc_postgresql_instance.go` - schema 定义及 CRUD 逻辑（Create、Update、Read 中 `backup_plan` 相关处理）
  - `tencentcloud/services/postgresql/resource_tc_postgresql_instance.md` - 资源示例文档
- **云API依赖**: `ModifyBackupPlan` 接口（`postgres/v20170312`）已新增 `BackupMethod` 入参；`BackupPlan` 响应结构已包含 `BackupMethod` 字段。vendor 中的 SDK 已更新到位，无需额外升级 SDK。
- **向后兼容**: `backup_method` 为 Optional 且 Computed 字段，不破坏现有配置与 state。
