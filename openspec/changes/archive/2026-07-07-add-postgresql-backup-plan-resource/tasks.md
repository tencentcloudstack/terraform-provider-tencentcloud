## 1. Service 层实现

- [x] 1.1 在 `tencentcloud/services/postgresql/service_tencentcloud_postgresql.go` 中新增 `DescribePostgresqlBackupPlanById` 方法：调用 `DescribeBackupPlans`（按 `DBInstanceId` 查询），从返回的 `Plans` 列表中匹配 `PlanId` 等于目标值的 `BackupPlan`，使用 `tccommon.ReadRetryTimeout` + `resource.Retry` 包装重试逻辑，返回匹配到的 `*postgres.BackupPlan`

## 2. 资源 CRUD 实现

- [x] 2.1 创建 `tencentcloud/services/postgresql/resource_tc_postgresql_backup_plan.go`，定义 `ResourceTencentCloudPostgresqlBackupPlan()`，包含 schema：
  - `db_instance_id`（Required/ForceNew/String）
  - `plan_name`（Required/String）
  - `backup_period_type`（Required/ForceNew/String）
  - `backup_period`（Required/TypeSet(TypeString)）
  - `min_backup_start_time`（Optional/String）
  - `max_backup_start_time`（Optional/String）
  - `base_backup_retention_period`（Optional/Int）
  - `log_backup_retention_period`（Optional/Int）
  - `backup_method`（Optional/String）
  - `plan_id`（Computed/String）
  - 支持 Importer（`schema.ImportStatePassthrough`）
- [x] 2.2 实现 Create handler：构建 `CreateBackupPlanRequest`，用 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 调用 `CreateBackupPlanWithContext`；成功后检查 `response.Response.PlanId` 为空则返回 `NonRetryableError`；用 `db_instance_id` + `tccommon.FILED_SP` + `plan_id` 拼接联合 ID 并 `d.SetId()`；最后调用 Read handler
- [x] 2.3 实现 Read handler：从 `d.Id()` 按 `tccommon.FILED_SP` 解析 `db_instance_id` 与 `plan_id`；调用 service 层 `DescribePostgresqlBackupPlanById`；未匹配到时先 `log.Printf("[CRUD] postgresql backup_plan id=%s", d.Id())` 再 `d.SetId("")`；匹配到后对每个非 nil 字段调用 `d.Set()`；`backup_period` 需 `json.Unmarshal` 反序列化字符串后 set
- [x] 2.4 实现 Update handler：构建 `ModifyBackupPlanRequest`（含 `DBInstanceId`、`PlanId` 及变更字段），用 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 调用 `ModifyBackupPlanWithContext`；调用失败用 `tccommon.RetryError()` 包装；成功后调用 Read handler
- [x] 2.5 实现 Delete handler：从 `d.Id()` 解析 `db_instance_id` 与 `plan_id`；构建 `DeleteBackupPlanRequest`，用 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 调用 `DeleteBackupPlanWithContext`；调用失败用 `tccommon.RetryError()` 包装

## 3. Provider 注册

- [x] 3.1 在 `tencentcloud/provider.go` 的 ResourcesMap 中注册 `"tencentcloud_postgresql_backup_plan": postgresql.ResourceTencentCloudPostgresqlBackupPlan()`

## 4. 文档

- [x] 4.1 创建 `tencentcloud/services/postgresql/resource_tc_postgresql_backup_plan.md`：一句话描述（带云产品名 PostgreSQL）、Example Usage（含必填参数示例）、Import 部分（说明使用联合 ID `db_instance_id{FILED_SP}plan_id`）；不添加 Argument/Attribute Reference 部分

## 5. 单元测试

- [x] 5.1 创建 `tencentcloud/services/postgresql/resource_tc_postgresql_backup_plan_test.go`：使用 gomonkey mock 云API client（`CreateBackupPlanWithContext`、`DescribeBackupPlansWithContext`、`ModifyBackupPlanWithContext`、`DeleteBackupPlanWithContext`），编写 Create/Read/Update/Delete 业务逻辑单元测试
- [x] 5.2 使用 `go test -gcflags=all=-l` 运行单元测试，确保全部通过

## 6. 验证

- [x] 6.1 检查所有函数返回的 error 已正确处理；不会出错的用 `_ = func()` 赋值给 `_`
- [x] 6.2 确认资源代码无 `_extension.go` 文件、go 文件开头无注释
- [x] 6.3 确认 Create handler 已检查返回值 `Response.PlanId` 为空的场景
- [x] 6.4 确认 Read handler 在返回空时先打印 `[CRUD]` 日志再 `d.SetId("")`
