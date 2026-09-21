## 1. Service 层封装 RestoreDBInstance

- [x] 1.1 在 `tencentcloud/services/mongodb/service_tencentcloud_mongodb.go` 新增 `RestoreDBInstance(ctx, instanceId, restoreTime string, databases []*mongodb.RestoreDatabases) (flowId int64, errRet error)` 方法：构造 `mongodb.NewRestoreDBInstanceRequest()`，填充 `InstanceId`、`RestoreTime`、`Databases`
- [x] 1.2 在该方法内使用 `resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {...})` 包装 `me.client.UseMongodbClient().RestoreDBInstanceWithContext(ctx, request)` 调用，调用前 `ratelimit.Check(request.GetAction())`，失败用 `tccommon.RetryError(e)` 包装，成功记录 DEBUG 日志
- [x] 1.3 retry 成功后检查 `result == nil || result.Response == nil || result.Response.FlowId == nil`，为空则返回错误；非空则提取并返回 `*result.Response.FlowId`（int64）
- [x] 1.4 添加 defer 函数：失败时记录 `[CRITAL]%s api[%s] fail, request body [%s], reason[%s]` 日志（含 logId、action、请求体、错误）

## 2. Framework Action 资源实现

- [x] 2.1 新建 `tencentcloud/services/mongodb/action_tc_mongodb_restore_db_instance.go`，参考 `action_tc_bdrc_run_copy_pair_tasks.go` 结构：`var _ action.ActionWithConfigure = &MongodbRestoreDbInstance{}`、工厂 `NewMongodbRestoreDbInstance() action.Action`、结构体 `MongodbRestoreDbInstance` 内嵌 `fw.ActionWithConfigure`
- [x] 2.2 定义模型 `MongodbRestoreDbInstanceModel`：`InstanceId types.String`、`RestoreTime types.String`、`Databases types.List`（元素为嵌套对象类型），嵌套对象含 `Db types.String` 与 `Collections types.List`（元素为含 `OldCollection`、`NewCollection` 的嵌套对象）
- [x] 2.3 实现 `Metadata`：`resp.TypeName = "tencentcloud_mongodb_restore_db_instance"`
- [x] 2.4 实现 `Schema`：`instance_id`（Required, schema.StringAttribute）、`restore_time`（Required, schema.StringAttribute）、`databases`（Required, schema.ListNestedBlock，嵌套对象含 `db` 与 `collections` ListNestedBlock），各字段含 Description
- [x] 2.5 实现 `Invoke`：解析 `req.Config` 到模型；校验 `instance_id`/`restore_time`/`databases` 为空/null 时通过 `resp.Diagnostics.AddError` 返回；校验 `a.Client() == nil` 时返回；将 `types.List` 转 `[]*mongodb.RestoreDatabases`，逐层填充 `RestoreDatabases`（`Db`、`Collections []*RestoreCollection`）；调用 `service.RestoreDBInstance(ctx, ...)` 拿到 `flowId`，再调用 `service.DescribeAsyncRequestInfo(ctx, helper.Int64ToStr(flowId), 3*tccommon.ReadRetryTimeout)` 轮询；失败时 `resp.Diagnostics.AddError` 上报
- [x] 2.6 不实现 Read/Update/Delete、不设置 id、不回写任何输出属性（一次性操作）

## 3. Action 注册

- [x] 3.1 在 `tencentcloud/framework/registry.go` 的 import 块中新增 `mongodb` 服务包导入 `"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/mongodb"`（若不存在）
- [x] 3.2 在 `actionFactories` 切片末尾追加 `mongodb.NewMongodbRestoreDbInstance`
- [x] 3.3 不修改 `tencentcloud/provider.go`（framework action 通过 registry.go 统一注册）

## 4. 单元测试

- [x] 4.1 新建 `tencentcloud/services/mongodb/action_tc_mongodb_restore_db_instance_test.go`，参考 `action_tc_bdrc_run_copy_pair_tasks_test.go` 结构：定义 ptrString 辅助、schema/config 构造辅助（含嵌套 databases/collections 的 tftypes.Object 构造）、setClient 注入辅助
- [x] 4.2 新增 `TestMongodbRestoreDbInstance_Invoke_Success`：mock `UseMongodbClient`/`RestoreDBInstanceWithContext`（返回 FlowId）与 `DescribeAsyncRequestInfo`（返回 success），断言 request 的 InstanceId/RestoreTime/Databases 符合预期、无错误 diagnostic
- [x] 4.3 新增 `TestMongodbRestoreDbInstance_Invoke_APIError`：mock `RestoreDBInstanceWithContext` 返回错误，断言产生错误 diagnostic 且包含错误信息
- [x] 4.4 新增 `TestMongodbRestoreDbInstance_Invoke_AsyncTaskFailed`：mock `RestoreDBInstanceWithContext` 返回 FlowId 但 `DescribeAsyncRequestInfo` 报 failed，断言产生错误 diagnostic
- [x] 4.5 新增 `TestMongodbRestoreDbInstance_Invoke_MissingInput`：校验 instance_id/restore_time/databases 为空时直接返回错误 diagnostic、未调用 API
- [x] 4.6 新增 `TestMongodbRestoreDbInstance_Invoke_ClientNotConfigured`：未注入 client 时断言 "Provider not configured" 错误
- [x] 4.7 新增 `TestMongodbRestoreDbInstance_MetadataAndSchema`：断言 TypeName 与所有必填属性/block 存在且类型正确

## 5. 文档

- [x] 5.1 新建 `tencentcloud/services/mongodb/action_tc_mongodb_restore_db_instance.md`：一句话描述带上 MongoDB 云产品名称；Example Usage 使用 `action "tencentcloud_mongodb_restore_db_instance" "example" { config { instance_id = ... restore_time = ... databases { db = ... collections { old_collection = ... new_collection = ... } } } }` 语法（含嵌套块）；含 NOTE 提示 Terraform 1.14+ 支持；无 Import 段、无手动 Argument/Attribute Reference

## 6. 验证

- [x] 6.1 代码正确性检查：确认 RestoreDBInstance 入参（InstanceId、RestoreTime、Databases、Databases[].Db、Databases[].Collections、Databases[].Collections[].OldCollection、Databases[].Collections[].NewCollection）均存在于 SDK 的 `RestoreDBInstanceRequest`/`RestoreDatabases`/`RestoreCollection` 中
- [x] 6.2 确认 framework action 的 InvokeResponse 仅暴露 Diagnostics 与 SendProgress，未尝试回写 flow_id 输出属性
- [x] 6.3 确认 retry 块内仅调用 API，成功设置/轮询操作在 retry 块外
- [x] 6.4 确认所有函数返回的 error 均被检查；必不出错函数用 `_ = func()` 忽略 err

## 7. 收尾

- [ ] 7.1 执行收尾阶段 `tfpacer-finalize` skill：gofmt 格式化变更的 Go 代码、`make doc` 生成 website/docs 文档、生成 changelog 文件并统一 amend 推送