## Why

腾讯云 MongoDB 提供了 `RestoreDBInstance` 接口，可按指定时间点和库表清单将实例回档到新的集合中，是数据误删、数据损坏等故障场景下的核心恢复手段。当前 Terraform Provider 已具备 MongoDB 实例、SSL、透明数据加密、备份等管理能力，但缺少对回档操作的声明式编排，用户只能通过控制台或 SDK 手动触发，导致故障恢复流程无法纳入 IaC 体系、难以审计与复现。本次新增基于 Terraform Plugin Framework 的 action 资源 `tencentcloud_mongodb_restore_db_instance`，使回档操作可通过 Terraform 声明式发起，并自动轮询异步任务直至成功。

## What Changes

- 新增 Terraform Plugin Framework action 资源 `tencentcloud_mongodb_restore_db_instance`，文件 `tencentcloud/services/mongodb/action_tc_mongodb_restore_db_instance.go`，严格参考 `NewTeoConfirmOriginAclUpdate` 与 `NewBdrcRunCopyPairTasks` 的代码结构：
  - `instance_id`（Required, String）：待回档的实例 ID，格式 `cmgo-xxxxxxxx`。
  - `restore_time`（Required, String）：回档目标时间点，格式 `YYYY-MM-DD hh:mm:ss`，须处于备份保留期内。
  - `databases`（Required, List Nested Block）：回档的库表信息列表，每项含 `db`（Required, String）与 `collections`（Required, List Nested Block）。
  - `collections`（Required, List Nested Block，位于 `databases` 块内）：集合回档信息，每项含 `old_collection`（Required, String）与 `new_collection`（Required, String）。
  - 仅实现 `Invoke`（对应 Create）：调用 `RestoreDBInstanceWithContext` 触发异步回档任务，拿到 `FlowId` 后通过 `MongodbService.DescribeAsyncRequestInfo` 轮询直到任务成功（`success`）或失败（`failed`）。
- 在 `tencentcloud/services/mongodb/service_tencentcloud_mongodb.go` 新增 `RestoreDBInstance` service 层方法封装，内部使用 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 包装 `RestoreDBInstanceWithContext` 调用，失败用 `tccommon.RetryError(e)` 包装；成功后返回 `FlowId` 供 Invoke 轮询。
- 在 `tencentcloud/framework/registry.go` 的 `actionFactories` 末尾注册 `mongodb.NewMongodbRestoreDbInstance`，并在 import 块中新增 `services/mongodb` 包导入（若不存在）。
- 新增单元测试 `tencentcloud/services/mongodb/action_tc_mongodb_restore_db_instance_test.go`，使用 gomonkey mock 云 API，覆盖 Invoke 成功、API 错误、缺必填输入、client 未配置、Metadata/Schema 校验等路径。
- 新增资源文档 `tencentcloud/services/mongodb/action_tc_mongodb_restore_db_instance.md`（一句话描述带上 mongodb 云产品名称；action 语法 Example Usage；NOTE 提示 Terraform 1.14+ 支持；无 Import 段；无手动 Argument/Attribute Reference）。
- 无需修改 `tencentcloud/provider.go`（framework action 资源通过 `framework/registry.go` 统一注册，不在 SDKv2 ResourcesMap 中注册）。
- connectivity 层无需改动：mongodb 客户端访问方法 `UseMongodbClient()` 已存在。

## Capabilities

### New Capabilities
- `mongodb-restore-db-instance-action`: 提供一次性 action 资源 `tencentcloud_mongodb_restore_db_instance`，通过 `RestoreDBInstance` API 将 MongoDB 实例回档到指定时间点的指定库表集合，调用后异步轮询任务直至成功，操作完成后不持久化任何云端状态。

### Modified Capabilities
无。

## Impact

- 代码：
  - `tencentcloud/services/mongodb/action_tc_mongodb_restore_db_instance.go`（新增 — framework action 资源）
  - `tencentcloud/services/mongodb/service_tencentcloud_mongodb.go`（修改 — 新增 `RestoreDBInstance` service 方法封装，含 retry 与 FlowId 返回）
  - `tencentcloud/services/mongodb/action_tc_mongodb_restore_db_instance_test.go`（新增 — 单元测试）
  - `tencentcloud/services/mongodb/action_tc_mongodb_restore_db_instance.md`（新增 — 文档）
  - `tencentcloud/framework/registry.go`（修改 — 注册 action factory 并新增 services/mongodb import）
- 依赖：使用已 vendored 的 `tencentcloud-sdk-go` 中 `mongodb/v20190725` 包（`RestoreDBInstanceRequest` / `RestoreDBInstanceResponse` / `RestoreDatabases` / `RestoreCollection`），无需新增 vendor；复用现有 `MongodbService.DescribeAsyncRequestInfo` 轮询异步任务。
- 向后兼容：纯新增资源，不影响已有资源配置与 state。
- 文档：需通过收尾阶段 `make doc` 生成 `website/docs/` 下对应文档（禁止手改 website 目录）。