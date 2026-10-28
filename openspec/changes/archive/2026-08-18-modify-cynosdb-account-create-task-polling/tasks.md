## 1. 修改 Create 函数：添加异步任务轮询

- [x] 1.1 在 `resourceTencentCloudCynosdbAccountCreate` 的 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 块内，调用 `CreateAccounts` 成功后检查返回值：若 `result == nil || result.Response == nil || result.Response.TaskId == nil`，返回 `resource.NonRetryableError`；否则保存 `TaskId` 到外层变量 `taskId`
- [x] 1.2 在 retry 块外、`d.SetId(...)` 之前，判断 `taskId` 是否有效（非 nil 且大于 0），若有效则使用 `tccommon.BuildStateChangeConf` + `service.taskStateRefreshFunc` 轮询任务状态至 `success`；若 `taskId` 为 nil 或 0 则跳过轮询
- [x] 1.3 在 `resource_tc_cynosdb_account.go` 中确保新增 `import` 所需的包（`strconv`、`time` 等，参考 `resource_tc_cynosdb_ssl.go`）

## 2. 修改 Read 函数：添加有限重试

- [x] 2.1 在 `resourceTencentCloudCynosdbAccountRead` 中，使用 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 包裹 `DescribeCynosdbAccountById` 调用，当 `account == nil` 时返回 `resource.RetryableError`，API 错误时返回 `tccommon.RetryError(e)`
- [x] 2.2 重试耗尽后仍未找到账号时，打印 `log.Printf("[CRUD] cynosdb_account id=%s", d.Id())` 保留现场日志，再执行 `d.SetId("")`

## 3. 补充单元测试

- [x] 3.1 在 `resource_tc_cynosdb_account_test.go` 中使用 gomonkey mock 方式补充单元测试，覆盖 Create 函数中 TaskId 有效时轮询至成功的逻辑
- [x] 3.2 补充单元测试，覆盖 Create 函数中 TaskId 为 nil/0 时跳过轮询直接 Read 的逻辑
- [x] 3.3 补充单元测试，覆盖 Read 函数中首次查不到账号时重试成功和重试耗尽的逻辑

## 4. 验证

- [x] 4.1 检查 `resource_tc_cynosdb_account.go` 中所有函数返回的 error 均已处理，无未使用变量
- [x] 4.2 确认未修改 Schema 定义，未新增/删除 Schema 字段，保持向后兼容
