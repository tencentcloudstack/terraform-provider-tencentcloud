## Why

TencentCloud PostgreSQL 提供了自定义备份计划（BackupPlan）能力，允许用户为实例创建自定义备份计划并管理其整个生命周期（创建、查询、修改、删除）。当前 provider 中仅存在一个配置型资源 `tencentcloud_postgresql_backup_plan_config`，它只能修改实例的默认备份计划、且没有真正的创建与删除语义，无法满足用户通过 Terraform 管理多个自定义备份计划的需求。新增一个完整的 CRUD 资源 `tencentcloud_postgresql_backup_plan` 可以让用户以基础设施即代码的方式创建并管理 PostgreSQL 自定义备份计划的完整生命周期。

## What Changes

- 新增 Terraform 资源 `tencentcloud_postgresql_backup_plan`（RESOURCE_KIND_GENERAL），通过云API `CreateBackupPlan`、`DescribeBackupPlans`、`ModifyBackupPlan`、`DeleteBackupPlan` 实现备份计划的创建、读取、更新、删除。
- 资源 schema 包含：`db_instance_id`（Required/ForceNew）、`plan_name`（Required）、`backup_period_type`（Required/ForceNew）、`backup_period`（Required）、`min_backup_start_time`（Optional）、`max_backup_start_time`（Optional）、`base_backup_retention_period`（Optional）、`log_backup_retention_period`（Optional）、`backup_method`（Optional）以及 `plan_id`（Computed）等字段。
- 资源 ID 采用 `db_instance_id` 与 `plan_id` 通过 `tccommon.FILED_SP` 分隔符拼接的联合 ID。
- 在 `provider.go` 中注册新资源 `tencentcloud_postgresql_backup_plan`。
- 同步创建资源文档 `resource_tc_postgresql_backup_plan.md`（用于 `make doc` 生成）。
- 同步创建单元测试文件 `resource_tc_postgresql_backup_plan_test.go`（使用 gomonkey mock 云API，不使用 terraform 测试套件）。

## Capabilities

### New Capabilities
- `postgresql-backup-plan-resource`: 管理 TencentCloud PostgreSQL 自定义备份计划（BackupPlan）的完整生命周期，包括创建备份计划、查询备份计划、修改备份计划、删除备份计划。

### Modified Capabilities
<!-- 无需修改已有能力。已有的 `tencentcloud_postgresql_backup_plan_config` 是独立的配置型资源，仅修改默认备份计划，本变更不改动它。 -->

## Impact

- **新增代码文件**:
  - `tencentcloud/services/postgresql/resource_tc_postgresql_backup_plan.go`
  - `tencentcloud/services/postgresql/resource_tc_postgresql_backup_plan.md`
  - `tencentcloud/services/postgresql/resource_tc_postgresql_backup_plan_test.go`
- **修改代码文件**:
  - `tencentcloud/provider.go`：注册 `tencentcloud_postgresql_backup_plan` 资源。
  - `tencentcloud/provider.md`：新增资源文档条目（由 `make doc` 自动生成）。
- **云API依赖**: 依赖 `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312` 中已有的 `CreateBackupPlan`、`DescribeBackupPlans`、`ModifyBackupPlan`、`DeleteBackupPlan` 接口（vendor 中已存在）。
- **向后兼容**: 新增资源，不修改已有资源 schema，完全向后兼容。
- **与已有资源关系**: 已有 `tencentcloud_postgresql_backup_plan_config`（配置型，修改默认备份计划），本资源为独立的完整 CRUD 资源，二者不冲突、互不影响。
