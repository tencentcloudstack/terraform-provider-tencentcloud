# postgresql-backup-plan-resource Specification

## Purpose
TBD - created by archiving change add-postgresql-backup-plan-resource. Update Purpose after archive.
## Requirements
### Requirement: PostgreSQL backup plan resource lifecycle

系统 SHALL 提供一个名为 `tencentcloud_postgresql_backup_plan` 的 Terraform 资源，用于管理 TencentCloud PostgreSQL 自定义备份计划（BackupPlan）的完整生命周期，包括创建、读取、更新、删除。该资源 MUST 调用云API `CreateBackupPlan`、`DescribeBackupPlans`、`ModifyBackupPlan`、`DeleteBackupPlan`（位于 `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312` 包）。

资源 ID MUST 使用 `db_instance_id` 与 `plan_id` 通过 `tccommon.FILED_SP` 分隔符拼接的联合 ID，格式为 `db_instance_id{FILED_SP}plan_id`。

#### Scenario: 创建备份计划

- **WHEN** 用户在 Terraform 配置中声明 `tencentcloud_postgresql_backup_plan` 资源，并提供 `db_instance_id`、`plan_name`、`backup_period_type`、`backup_period` 等必填参数
- **THEN** 系统 MUST 调用 `CreateBackupPlan` 接口创建备份计划，使用 `tccommon.WriteRetryTimeout` 作为超时时间并通过 `resource.Retry` 包装重试逻辑
- **AND** 调用成功后 MUST 检查返回值 `response.Response.PlanId` 是否为空，若为空 MUST 返回 `NonRetryableError`
- **AND** 系统 MUST 使用 `db_instance_id` 与返回的 `plan_id` 拼接联合 ID 并调用 `d.SetId()`，随后调用 Read handler 同步状态

#### Scenario: 读取备份计划

- **WHEN** 系统执行 Read 操作
- **THEN** 系统 MUST 从 `d.Id()` 中解析出 `db_instance_id` 与 `plan_id`（以 `tccommon.FILED_SP` 分隔）
- **AND** 系统 MUST 调用 `DescribeBackupPlans` 接口（入参为 `DBInstanceId`），使用 `tccommon.ReadRetryTimeout` 作为超时时间并通过 `resource.Retry` 包装重试逻辑
- **AND** 系统 MUST 从返回的 `Plans` 列表中匹配 `PlanId` 等于解析出的 `plan_id` 的备份计划
- **AND** 若云API返回为空或未匹配到对应 `plan_id` 的备份计划，系统 MUST 先打印 `log.Printf("[CRUD] postgresql backup_plan id=%s", d.Id())` 保留现场，再执行 `d.SetId("")`
- **AND** 匹配成功后，系统 MUST 在调用 `d.Set()` 设置字段前判断 Response 中对应字段是否为 nil，若为 nil 则不调用 `d.Set()`

#### Scenario: 更新备份计划

- **WHEN** 用户修改 `tencentcloud_postgresql_backup_plan` 资源的可变参数（如 `plan_name`、`backup_period`、`min_backup_start_time`、`max_backup_start_time`、`base_backup_retention_period`、`log_backup_retention_period`、`backup_method`）
- **THEN** 系统 MUST 调用 `ModifyBackupPlan` 接口更新备份计划，使用 `tccommon.WriteRetryTimeout` 作为超时时间并通过 `resource.Retry` 包装重试逻辑
- **AND** `ModifyBackupPlan` 的请求中 MUST 包含 `DBInstanceId` 与 `PlanId`（从 `d.Id()` 解析）
- **AND** 调用成功后 MUST 调用 Read handler 同步状态

#### Scenario: 删除备份计划

- **WHEN** 用户从 Terraform 配置中移除 `tencentcloud_postgresql_backup_plan` 资源
- **THEN** 系统 MUST 从 `d.Id()` 中解析出 `db_instance_id` 与 `plan_id`
- **AND** 系统 MUST 调用 `DeleteBackupPlan` 接口（入参为 `DBInstanceId` 与 `PlanId`），使用 `tccommon.WriteRetryTimeout` 作为超时时间并通过 `resource.Retry` 包装重试逻辑
- **AND** 调用失败时 MUST 使用 `tccommon.RetryError()` 将错误进行包装并返回

### Requirement: PostgreSQL backup plan resource schema

`tencentcloud_postgresql_backup_plan` 资源的 schema MUST 定义以下字段：

| 字段 | 类型 | 约束 | 说明 |
|---|---|---|---|
| `db_instance_id` | TypeString | Required, ForceNew | 实例 ID |
| `plan_name` | TypeString | Required | 备份计划名称 |
| `backup_period_type` | TypeString | Required, ForceNew | 备份周期类型，当前仅支持 month |
| `backup_period` | TypeSet(TypeString) | Required | 备份日期 |
| `min_backup_start_time` | TypeString | Optional | 备份最早开始时间 |
| `max_backup_start_time` | TypeString | Optional | 备份最晚开始时间 |
| `base_backup_retention_period` | TypeInt | Optional | 数据备份保留时长（天） |
| `log_backup_retention_period` | TypeInt | Optional | 日志备份保留时长（天） |
| `backup_method` | TypeString | Optional | 备份方式（physical/logical/snapshot） |
| `plan_id` | TypeString | Computed | 备份计划 ID |

由于 `backup_period_type` 只有创建时可选（云API `ModifyBackupPlan` 入参不包含该字段），该字段 MUST 标记为 ForceNew。

资源 MUST 支持 Import（`schema.ImportStatePassthrough`），导入时需使用联合 ID `db_instance_id{FILED_SP}plan_id`。

#### Scenario: schema 字段定义

- **WHEN** 检查 `tencentcloud_postgresql_backup_plan` 资源的 schema 定义
- **THEN** schema MUST 包含上述所有字段，且各字段的类型与约束符合上表
- **AND** `db_instance_id` 与 `backup_period_type` MUST 标记为 ForceNew

#### Scenario: backup_period 字段的读取处理

- **WHEN** Read handler 处理 `DescribeBackupPlans` 返回的 `BackupPlan.BackupPeriod`（类型为 `*string`）
- **THEN** 系统 MUST 将该字符串 JSON 反序列化为切片后设置到 `backup_period` 字段（schema 为 TypeSet），以保持与写入时的列表结构一致

### Requirement: Provider registration

系统 MUST 在 `tencentcloud/provider.go` 中注册 `tencentcloud_postgresql_backup_plan` 资源，映射到 `postgresql.ResourceTencentCloudPostgresqlBackupPlan()`。

#### Scenario: 资源注册

- **WHEN** 检查 `provider.go` 的 ResourcesMap
- **THEN** MUST 存在键 `"tencentcloud_postgresql_backup_plan"`，其值为 `postgresql.ResourceTencentCloudPostgresqlBackupPlan()`

### Requirement: Documentation and tests

系统 MUST 同步生成资源文档与单元测试：

- 资源文档 `resource_tc_postgresql_backup_plan.md` MUST 遵循 gendoc/README.md 格式，包含一句话描述（带云产品名称 PostgreSQL）、Example Usage、Import 部分（说明使用联合 ID）。
- 单元测试文件 `resource_tc_postgresql_backup_plan_test.go` MUST 使用 gomonkey mock 云API进行业务逻辑单元测试，不使用 terraform 测试套件，并使用 `go test -gcflags=all=-l` 跑通。

#### Scenario: 文档生成

- **WHEN** 执行 `make doc`
- **THEN** 系统 MUST 在 `website/docs/r/` 下生成 `postgresql_backup_plan.html.markdown` 文档

#### Scenario: 单元测试

- **WHEN** 执行 `go test -gcflags=all=-l` 运行 `resource_tc_postgresql_backup_plan_test.go`
- **THEN** 所有测试用例 MUST 通过，覆盖 Create、Read、Update、Delete 业务逻辑

