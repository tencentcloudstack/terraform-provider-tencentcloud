## Context

腾讯云 BDRC（业务数据弹性容灾）产品已在 `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/bdrc/v20260330` 中提供 `RunCopyPairTasks` 接口（启动复制对），但 Terraform Provider 当前未接入任何 BDRC 资源。

当前状态：
- BDRC SDK 包已 vendored，包含 `RunCopyPairTasksRequest`（入参 `CopyPairIds []*string`、`CopyPairType *string`）与 `RunCopyPairTasksResponse`（出参 `Response.CopyPairIds []*string`）。
- `tencentcloud/connectivity/client.go` 中尚无 bdrc 客户端，需新增 `UseBdrcV20260330Client()` 方法及对应字段与 import。
- 该资源为 **RESOURCE_KIND_OPERATION**：一次性操作，操作完成后不需要记录任何状态，只需实现 Create（对应 framework action 的 `Invoke`），RUD 接口为空。
- 用户要求使用 **Terraform Plugin Framework** 的 action 资源类型实现，参考 `NewTeoConfirmOriginAclUpdate`（`tencentcloud/services/teo/action_tc_teo_confirm_origin_acl_update.go`）的代码格式，忽略 SDKv2 设计。
- `RunCopyPairTasks` 为**同步接口**（SDK 注释与入出参均无异步/任务流标识），调用后无需轮询 Read 接口等待生效。

约束：
- framework action 的 `InvokeResponse` 仅暴露 `Diagnostics` 与 `SendProgress`，**不提供**写入输出/Computed 值的机制，因此 `RunCopyPairTasks` 响应中的 `CopyPairIds` 无法作为 Terraform 输出属性回写。这与"一次性操作不记录任何状态"的需求一致；API 响应仅用于确认启动成功。
- framework action 资源通过 `tencentcloud/framework/registry.go` 的 `actionFactories` 统一注册，不在 SDKv2 `provider.go` 的 `ResourcesMap` 中注册，因此无需修改 `provider.go`。

## Goals / Non-Goals

**Goals:**
- 新增 framework action 资源 `tencentcloud_bdrc_copy_pair_tasks`，通过 `RunCopyPairTasks` API 声明式启动一组 BDRC 复制对任务。
- 接入 BDRC SDK 客户端到 connectivity 层，供后续 BDRC 资源复用。
- 通过单元测试覆盖 Invoke 成功与 API 错误两条路径（gomonkey mock 云 API）。
- 保持与参考实现 `NewTeoConfirmOriginAclUpdate` 一致的代码风格与结构。

**Non-Goals:**
- 不持久化任何云端状态（不设置 id、不实现 Read/Update/Delete、不回写响应字段）。
- 不为 `copy_pair_type` 做 schema 层枚举校验（DISK/INSTANCE/CFS），由云 API 返回错误即可，与参考实现保持一致。
- 不新增 `provider.go` 注册条目（framework action 走 registry.go）。
- 不新增 `_extension.go` 文件。
- 不轮询 Read 接口（接口为同步）。

## Decisions

### Decision 1: 使用 Terraform Plugin Framework action 类型，而非 SDKv2 operation 资源

**选择**：按用户要求，使用 framework action（`action.Action` 接口）实现，文件命名为 `action_tc_bdrc_copy_pair_tasks.go`，工厂函数 `NewBdrcCopyPairTasks()` 返回 `action.Action`，在 `framework/registry.go` 的 `actionFactories` 注册。

**备选**：使用 SDKv2 的 operation 资源模式（`resource_tc_bdrc_copy_pair_tasks_operation.go`，在 `provider.go` ResourcesMap 注册）。

**理由**：
- 用户明确要求使用 framework plugin 框架实现 action 资源类型，参考 `NewTeoConfirmOriginAclUpdate`，忽略其他 SDKv2 设计。
- framework action 是表达"一次性操作"语义的标准方式，与参考实现一致。

### Decision 2: schema 仅包含两个 Required 输入字段，不回写响应输出

**选择**：schema 定义 `copy_pair_ids`（Required, List of String）与 `copy_pair_type`（Required, String）。`Invoke` 中调用 `RunCopyPairTasksWithContext`，成功后直接返回（仅通过 Diagnostics 上报错误），不设置任何输出属性。

**备选**：在 schema 中增加 Computed 的 `copy_pair_ids` 输出属性以回写 API 响应。

**理由**：
- framework action 的 `InvokeResponse` 不提供写入输出/Computed 值的机制（仅 `Diagnostics` 与 `SendProgress`），技术上无法回写。
- 需求明确"一次性操作，操作完不需要记录任何状态"，与参考实现 `NewTeoConfirmOriginAclUpdate`（只读输入、不写输出）一致。
- API 响应 `CopyPairIds` 仅用于日志记录确认启动成功，不作为 Terraform 属性。

### Decision 3: 通过 service 层方法封装 API 调用与 retry

**选择**：在 `tencentcloud/services/bdrc/service_tencentcloud_bdrc.go` 中新增 `BdrcService` 与 `RunCopyPairTasks(ctx, copyPairIds, copyPairType)` 方法，内部使用 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 包装 `RunCopyPairTasksWithContext` 调用，失败时用 `tccommon.RetryError(e)` 包装。action 的 `Invoke` 方法构造参数并调用该 service 方法。

**备选**：在 `Invoke` 中直接调用 SDK client（不经 service 层），如参考实现那样。

**理由**：
- 参考实现 `TeoConfirmOriginAclUpdate` 通过 `service.ConfirmOriginACLUpdate(ctx, zoneId)` 封装 retry 逻辑，保持 action 文件简洁；本次沿用同样分层。
- 统一 retry/日志模式，便于后续 BDRC 资源复用 `BdrcService`。

### Decision 4: connectivity 层新增 bdrc 客户端

**选择**：在 `tencentcloud/connectivity/client.go` 中新增 import `bdrcv20260330 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/bdrc/v20260330"`、结构体字段 `bdrcv20260330Conn *bdrcv20260330.Client`，以及 `UseBdrcV20260330Client()` 方法（参考 `UseTeoV20220901Client()` 实现，超时 300s、挂载 `LogRoundTripper`）。

**理由**：
- 当前 connectivity 层无 bdrc 客户端，必须新增才能在 service 层访问 SDK。
- 沿用现有客户端方法命名与实现模式（`UseXxxV<version>Client`），保持一致性。

### Decision 5: 文档命名与格式

**选择**：文档文件 `action_tc_bdrc_copy_pair_tasks.md`，格式为一句话描述（带上云产品名称 BDRC）+ Example Usage（action 语法 `action "tencentcloud_bdrc_copy_pair_tasks" "example" { config { ... } }`）+ NOTE 提示（参考 teo action 文档）。无 Import 段、无 Argument/Attribute Reference。

**理由**：
- 与参考 framework action 文档 `action_tc_teo_confirm_origin_acl_update.md` 格式一致。
- 遵循资源文档规范：action 资源无 Import，Argument/Attribute Reference 由工具自动生成。

## Risks / Trade-offs

- **Risk**：framework action 功能依赖 Terraform 1.14+ 及 provider 正确的 mux 配置 → **Mitigation**：仓库已存在 `add-plugin-framework-muxing` 变更与 `framework/registry.go` 机制，action 注册路径已验证（teo action 已注册并通过），本次仅追加一条 factory 条目，风险可控。
- **Trade-off**：API 响应 `CopyPairIds` 不回写为 Terraform 输出，用户无法在下游引用"已启动的复制对 ID" → 可接受，需求明确为一次性操作不记录状态；用户可自行通过数据源或原 `copy_pair_ids` 输入获取。
- **Trade-off**：`copy_pair_type` 不做 schema 层枚举校验 → 与参考实现及 provider 多数字段一致，非法值由云 API 返回错误，减少与未来 API 扩展耦合。

## Migration Plan

- 纯新增资源与 connectivity 客户端，无 state 迁移需求。
- 存量配置不受影响（不修改任何已有资源 schema）。
- 回滚：删除 `action_tc_bdrc_copy_pair_tasks.go`、`service_tencentcloud_bdrc.go`、文档、测试，并移除 `framework/registry.go` 与 `connectivity/client.go` 中的 bdrc 条目即可。

## Open Questions

- 无
