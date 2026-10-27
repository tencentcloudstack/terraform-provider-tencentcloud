## Context

腾讯云 MongoDB 产品的 SDK 包 `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mongodb/v20190725` 已在 vendor 中提供 `RestoreDBInstance` 接口（回档实例到指定时间点），但 Terraform Provider 当前未接入该操作。

当前状态：
- MongoDB SDK 包已 vendored，包含：
  - `RestoreDBInstanceRequest`（入参 `InstanceId *string`、`RestoreTime *string`、`Databases []*RestoreDatabases`）。
  - `RestoreDBInstanceResponse`（出参 `Response.FlowId *int64`，回档任务流程 ID）。
  - `RestoreDatabases`（`Db *string`、`Collections []*RestoreCollection`）。
  - `RestoreCollection`（`OldCollection *string`、`NewCollection *string`）。
- `tencentcloud/connectivity/client.go` 中已有 mongodb 客户端访问方法 `UseMongodbClient()`，无需新增。
- `tencentcloud/services/mongodb/service_tencentcloud_mongodb.go` 中已有 `MongodbService` 结构体与 `DescribeAsyncRequestInfo(ctx, asyncId string, timeout time.Duration) error` 方法，可通过 `FlowId` 轮询异步任务状态（`success`/`failed`）。
- 该资源为 **RESOURCE_KIND_OPERATION**：一次性操作，操作完成后不需要记录任何状态，只需实现 Create（对应 framework action 的 `Invoke`），RUD 接口为空。
- 用户要求使用 **Terraform Plugin Framework** 的 action 资源类型实现，严格参考 `NewTeoConfirmOriginAclUpdate`（`tencentcloud/services/teo/action_tc_teo_confirm_origin_acl_update.go`）与 `NewBdrcRunCopyPairTasks`（`tencentcloud/services/bdrc/action_tc_bdrc_run_copy_pair_tasks.go`）的代码结构，忽略 SDKv2 设计。
- `RestoreDBInstance` 为**异步接口**：返回 `FlowId`（int64），调用成功后需通过 `DescribeAsyncRequestInfo` 轮询直到任务成功（`success`），失败（`failed`）或超时则返回错误。

约束：
- framework action 的 `InvokeResponse` 仅暴露 `Diagnostics` 与 `SendProgress`，**不提供**写入输出/Computed 值的机制，因此 `RestoreDBInstance` 响应中的 `FlowId` 无法作为 Terraform 输出属性回写。这与"一次性操作不记录任何状态"的需求一致；`FlowId` 仅用于调用 `DescribeAsyncRequestInfo` 轮询任务完成，并在日志中确认回档任务已触发。
- framework action 资源通过 `tencentcloud/framework/registry.go` 的 `actionFactories` 统一注册，不在 SDKv2 `provider.go` 的 `ResourcesMap` 中注册，因此无需修改 `provider.go`。

## Goals / Non-Goals

**Goals:**
- 新增 framework action 资源 `tencentcloud_mongodb_restore_db_instance`，通过 `RestoreDBInstance` API 声明式发起 MongoDB 实例回档操作。
- 调用成功后通过 `DescribeAsyncRequestInfo` 轮询异步任务直到成功或失败，保证回档实际生效。
- 通过单元测试覆盖 Invoke 成功、API 错误、缺必填输入、client 未配置、Metadata/Schema 校验等路径（gomonkey mock 云 API）。
- 保持与参考实现 `NewTeoConfirmOriginAclUpdate` / `NewBdrcRunCopyPairTasks` 一致的代码风格与结构。

**Non-Goals:**
- 不持久化任何云端状态（不设置 id、不实现 Read/Update/Delete、不回写 `flow_id` 等）。
- 不为 `restore_time` 做时间格式 schema 校验，由云 API 返回错误。
- 不新增 `provider.go` 注册条目（framework action 走 registry.go）。
- 不新增 `_extension.go` 文件。
- 不修改 connectivity 层（mongodb 客户端已存在）。

## Decisions

### Decision 1: 使用 Terraform Plugin Framework action 类型，而非 SDKv2 operation 资源

**选择**：按用户要求，使用 framework action（`action.Action` 接口）实现，文件命名为 `action_tc_mongodb_restore_db_instance.go`，工厂函数 `NewMongodbRestoreDbInstance()` 返回 `action.Action`，在 `framework/registry.go` 的 `actionFactories` 注册。

**备选**：使用 SDKv2 的 operation 资源模式（`resource_tc_mongodb_restore_db_instance_operation.go`，在 `provider.go` ResourcesMap 注册）。

**理由**：
- 用户明确要求使用 framework plugin 框架实现 action 资源类型，参考 `NewTeoConfirmOriginAclUpdate` 与 `NewBdrcRunCopyPairTasks`，忽略其他 SDKv2 设计。
- framework action 是表达"一次性操作"语义的标准方式，与参考实现一致。

### Decision 2: schema 使用嵌套 block 表达 databases/collections 层级，不回写 flow_id

**选择**：schema 定义 `instance_id`（Required, String）、`restore_time`（Required, String）、`databases`（Required, List Nested Block）。`databases` 块内含 `db`（Required, String）与 `collections`（Required, List Nested Block）。`collections` 块内含 `old_collection`（Required, String）与 `new_collection`（Required, String）。`Invoke` 中调用 `RestoreDBInstanceWithContext`，拿到 `FlowId` 后调用 `DescribeAsyncRequestInfo` 轮询；成功后直接返回（仅通过 Diagnostics 上报错误），不设置任何输出属性。

**备选**：在 schema 中增加 Computed 的 `flow_id` 输出属性以回写 API 响应。

**理由**：
- framework action 的 `InvokeResponse` 不提供写入输出/Computed 值的机制（仅 `Diagnostics` 与 `SendProgress`），技术上无法回写。
- 需求明确"一次性操作，操作完不需要记录任何状态"，与参考实现 `NewTeoConfirmOriginAclUpdate`、`NewBdrcRunCopyPairTasks`（只读输入、不写输出）一致。
- API 响应 `FlowId` 仅用于轮询任务完成与日志记录确认回档已触发，不作为 Terraform 属性。
- 采用嵌套 block（`schema.ListNestedBlock`）而非扁平结构，以准确映射 `Databases[].Collections[]` 两层级关系，便于 Terraform 校验与可读 HCL。

### Decision 3: 通过 service 层方法封装 API 调用、retry 与 FlowId 返回

**选择**：在 `tencentcloud/services/mongodb/service_tencentcloud_mongodb.go` 中新增 `MongodbService.RestoreDBInstance(ctx, instanceId, restoreTime string, databases []*mongodb.RestoreDatabases) (flowId int64, errRet error)` 方法，内部使用 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 包装 `RestoreDBInstanceWithContext` 调用，失败用 `tccommon.RetryError(e)` 包装；成功后提取并返回 `FlowId`，defer 中记录失败日志。action `Invoke` 方法构造参数、调用该 service 方法、拿到 `flowId` 后调用 `DescribeAsyncRequestInfo(ctx, helper.Int64ToStr(flowId), 3*tccommon.ReadRetryTimeout)` 轮询直到任务成功。

**备选**：在 `Invoke` 中直接调用 SDK client（不经 service 层），并在 Invoke 内同时完成 retry 与轮询。

**理由**：
- 参考实现 `TeoConfirmOriginAclUpdate` 通过 `service.ConfirmOriginACLUpdate(ctx, zoneId)` 封装 retry 逻辑，保持 action 文件简洁；本次沿用同样分层。
- `RestoreDBInstance` 的 retry 仅用于 API 调用本身（触发回档任务），轮询任务完成是独立步骤，不应放在同一个 retry 块中（避免 retry 语义混淆），由 Invoke 在拿到 flowId 后单独调用 `DescribeAsyncRequestInfo`。
- 统一 retry/日志模式，便于后续 mongodb 资源复用 `MongodbService`。

### Decision 4: 异步任务轮询复用现有 DescribeAsyncRequestInfo，超时取 3 * ReadRetryTimeout

**选择**：在 `Invoke` 中调用 `MongodbService.DescribeAsyncRequestInfo(ctx, flowIdString, 3*tccommon.ReadRetryTimeout)` 轮询回档任务状态，与 `tencentcloud_mongodb_instance_transparent_data_encryption` 等现有资源一致；将 `flowId`（int64）通过 `helper.Int64ToStr` 转字符串后传入。

**备选**：轮询超时取 20 * ReadRetryTimeout（如 `tencentcloud_mongodb_instance` 创建用）。

**理由**：
- 回档操作通常较创建实例轻量，3 * ReadRetryTimeout 与现有透明数据加密等异步操作一致即可；如遇回档大库表耗时长，可后续调整。
- 复用已验证的轮询方法，避免重复实现。

### Decision 5: 不注册到 SDKv2 provider.go，仅在 framework/registry.go 注册

**选择**：仅在 `tencentcloud/framework/registry.go` 的 `actionFactories` 追加 `mongodb.NewMongodbRestoreDbInstance`，并在 import 块新增 `services/mongodb` 包导入（若不存在）；不修改 `tencentcloud/provider.go`。

**理由**：
- framework action 资源通过 `framework/registry.go` 统一注册，与 teo、bdrc action 一致；在 SDKv2 `provider.go` ResourcesMap 注册会导致重复注册与 mux 冲突。

### Decision 6: 文档命名与格式

**选择**：文档文件 `action_tc_mongodb_restore_db_instance.md`，格式为一句话描述（带上云产品名称 MongoDB）+ Example Usage（action 语法 `action "tencentcloud_mongodb_restore_db_instance" "example" { config { ... } }`）+ NOTE 提示 Terraform 1.14+ 支持。无 Import 段、无手动 Argument/Attribute Reference。

**理由**：
- 与参考 framework action 文档 `action_tc_teo_confirm_origin_acl_update.md`、`action_tc_bdrc_run_copy_pair_tasks.md` 格式一致。
- 遵循资源文档规范：action 资源无 Import，Argument/Attribute Reference 由工具自动生成。

## Risks / Trade-offs

- **Risk**：framework action 功能依赖 Terraform 1.14+ 及 provider 正确的 mux 配置 → **Mitigation**：仓库已存在 `add-plugin-framework-muxing` 变更与 `framework/registry.go` 机制，teo、bdrc action 已注册并通过，本次仅追加一条 factory 条目，风险可控。
- **Risk**：回档大库表场景下 3 * ReadRetryTimeout 可能不够 → **Mitigation**：与现有透明数据加密等异步操作一致；超时后通过 Diagnostics 上报错误，用户可重试。
- **Trade-off**：API 响应 `FlowId` 不回写为 Terraform 输出，用户无法在下游引用回档任务流程 ID → 可接受，需求明确为一次性操作不记录状态；用户可在 Apply 输出日志中查到 `FlowId`。

## Migration Plan

- 纯新增资源，无 state 迁移需求。
- 存量配置不受影响（不修改任何已有资源 schema）。
- 回滚：删除 `action_tc_mongodb_restore_db_instance.go`、测试、文档，并移除 `framework/registry.go` 中的 mongodb action 条目即可。

## Open Questions

- 无