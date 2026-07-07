## 1. Schema 定义

- [x] 1.1 在 `tencentcloud/services/postgresql/resource_tc_postgresql_instance.go` 的 `backup_plan` 嵌套块 schema 中新增 `backup_method` 字段（`schema.TypeString`，Optional: true，Computed: true），Description 说明枚举值 `physical`（物理备份）、`logical`（逻辑备份）、`snapshot`（快照备份）

## 2. CRUD 函数实现

- [x] 2.1 在 Create 流程的 `resourceTencentCloudPostgresqlInstanceCreate` 函数中，调用 `ModifyBackupPlan` 修改默认（week）备份计划时，从 `plan["backup_method"]` 读取用户配置，非空时设置 `request.BackupMethod`
- [x] 2.2 在 Update 流程的 `resourceTencentCloudPostgresqlInstanceUpdate` 函数中，当 `d.HasChange("backup_plan")` 且调用 `ModifyBackupPlan` 修改默认（week）备份计划时，从 `plan["backup_method"]` 读取用户配置，非空时设置 `request.BackupMethod`
- [x] 2.3 在 Read 流程的 `resourceTencentCloudPostgresqlInstanceRead` 函数中，构造默认（week）备份计划的 `planMap` 时，判断 `backupPlan.BackupMethod` 非 nil 后写入 `planMap["backup_method"]`，保持 nil 安全

## 3. 文档更新

- [x] 3.1 更新 `tencentcloud/services/postgresql/resource_tc_postgresql_instance.md`，在 `backup_plan` 块示例中补充 `backup_method` 参数用法及枚举值说明

## 4. 测试

- [x] 4.1 在 `tencentcloud/services/postgresql/resource_tc_postgresql_instance_test.go` 中补充 `backup_method` 参数的单元测试用例（使用 gomonkey mock 云API，验证 Create/Update 传入 BackupMethod、Read 回填逻辑）

## 5. 验证

- [x] 5.1 使用 `go test -gcflags=all=-l` 运行涉及的单元测试文件，确保测试通过
- [x] 5.2 运行 `openspec validate add-postgresql-instance-backup-method-param --strict` 验证提案完整性
