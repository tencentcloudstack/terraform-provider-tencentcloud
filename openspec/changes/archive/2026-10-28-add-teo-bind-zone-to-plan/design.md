## Context

TEO（EdgeOne）站点创建后需要绑定到已有套餐（Plan）才能正式使用。腾讯云 TEO SDK（`teo/v20220901`）提供了 `BindZoneToPlan` 接口，入参为 `ZoneId` 与 `PlanId`，返回值仅包含 `RequestId`，无业务数据，且为同步接口。当前 Terraform Provider 中 `tencentcloud_teo_zone` 资源不包含"绑定套餐"动作，用户必须借助控制台或 SDK 手动完成绑定，导致声明式运维流程断裂。

当前状态：
- vendor 中 `teov20220901.BindZoneToPlanRequest` 已存在，字段：`ZoneId`、`PlanId`（均为 `*string`）
- vendor 中 `teov20220901.Client.BindZoneToPlan` 与 `BindZoneToPlanWithContext` 已存在
- 返回 `BindZoneToPlanResponse`，其 `Response` 仅含 `RequestId`，无业务数据
- framework provider（`tencentcloud/framework`）已支持 action 类型，并通过 `registry.go` 的 `actionFactories` 聚合注册，可作为实现参考
- framework provider 通过 `sharedmeta.ProviderMeta` 向 action 注入共享的 `*connectivity.TencentCloudClient`

约束：
- action 是一次性操作，操作完不记录任何状态（无 Read/Update/Delete lifecycle）
- action 文件命名格式：`action_tc_teo_bind_zone_to_plan.go`
- action schema 字段不能为 Computed/Sensitive，`zone_id`、`plan_id` 均为 Required
- 使用 `tccommon.WriteRetryTimeout` 作为超时时间添加 retry 处理
- action 注册于 `tencentcloud/framework/registry.go` 的 `actionFactories`，并在 `tencentcloud/provider.md` 的 TEO 产品 `Action` 子节中登记

## Goals / Non-Goals

**Goals:**
- 通过 Terraform framework action 声明式地完成 TEO 站点与套餐的绑定
- 遵循 framework provider 既有 action 模式（参考 `tencentcloud/framework/action.Action` 接口与 ssm framework 资源）
- 提供基于 gomonkey 的单元测试，覆盖 Metadata、Schema、Configure、Invoke 成功、Invoke API 错误、Invoke 缺参数等场景
- 在 `registry.go` / `provider.md` 中完成 action 注册

**Non-Goals:**
- 不在 action 中暴露绑定后的套餐查询/校验（API 本身不返回业务数据，且 action 不管理 state）
- 不提供 Read/Update/Delete（action 一次性执行，无 lifecycle）
- 不提供 Importer（action 不可导入）
- 不修改 `tencentcloud_teo_zone` 等现有资源

## Decisions

### Decision 1: 采用 framework action 而非 SDKv2 operation 资源

**选择**：使用 `terraform-plugin-framework` 的 `action.Action` 接口实现 `tencentcloud_teo_bind_zone_to_plan`（文件 `action_tc_teo_bind_zone_to_plan.go`）。

**备选**：使用 SDKv2 的 `schema.Resource`（RESOURCE_KIND_OPERATION）实现，文件 `resource_tc_teo_bind_zone_to_plan_operation.go`。

**理由**：
- `BindZoneToPlan` 语义上是"动作"，操作完不需要记录任何状态，framework action 的语义正是"一次性操作，不造成 state 变化"，比 SDKv2 operation 资源（需要伪造 id / no-op Read / no-op Delete）更贴合
- provider framework 栈已支持 action（`Provider.Actions` / `frameworkActions`，`Configure.ActionData` 注入 `*sharedmeta.ProviderMeta`），注册成本低
- action 不在 state 中持久化任何数据，避免 SDKv2 operation 资源中需要 `helper.BuildToken()` 伪造 id 的问题

### Decision 2: action 通过 Configure 注入共享 SDK 客户端

**选择**：实现 `action.ActionWithConfigure`，在 `Configure` 中将 `req.ProviderData.(*sharedmeta.ProviderMeta).Client` 存入 action 结构体的 `Client` 字段，`Invoke` 中使用该 client 调用 `UseTeoV20220901Client().BindZoneToPlanWithContext`。

**理由**：
- 与 framework provider 的 `Configure` 注入机制一致（`resp.ActionData = meta`）
- 与 ssm framework 资源 `Configure` 注入 `*sharedmeta.ProviderMeta` 模式一致
- 复用 SDKv2 provider 构建的共享 `*connectivity.TencentCloudClient`，无需重复解析凭证

### Decision 3: retry 包装 BindZoneToPlan 调用

**选择**：使用 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 包装 `BindZoneToPlanWithContext` 调用，失败时用 `tccommon.RetryError(e)` 包装返回。

**理由**：
- 符合项目规范："调用云API接口时，需要添加 retry 处理"
- 写操作使用 `tccommon.WriteRetryTimeout`
- retry 块内只执行接口调用；成功后不写 state（action 无 state）

### Decision 4: schema 字段全部 Required

**选择**：`zone_id`、`plan_id` 均为 `Required: true`、`actionschema.StringAttribute`、`TypeString`。

**理由**：
- action 无 Update lifecycle，所有参数在 Invoke 时一次性传入
- action schema 不支持 ForceNew（ForceNew 是 SDKv2 managed resource 的概念），因此不设置
- 两个参数在 `BindZoneToPlanRequest` 中均为必需入参

### Decision 5: BindZoneToPlan 为同步接口，无需轮询

**选择**：调用 `BindZoneToPlan` 成功后直接返回，不调用 Read 接口轮询。

**理由**：
- `BindZoneToPlan` 为同步接口，返回即代表操作完成
- action 无 state，无需轮询确认

## Risks / Trade-offs

- **Risk**：用户对同一站点重复调用该 action，会重复调用 `BindZoneToPlan`，API 可能返回 `InvalidParameter.ZoneHasBeenBound` → **Mitigation**：由云 API 返回错误，retry 机制会在非可重试错误时停止；属于预期行为，action 普遍如此。
- **Risk**：`zone_id` 或 `plan_id` 不存在时 API 返回 `InvalidParameter.ZoneNotFound` / `InvalidParameter.PlanNotFound` → **Mitigation**：由 `tccommon.RetryError` 判定为非可重试错误并返回，提示用户检查参数。
- **Trade-off**：action 不记录绑定状态，用户无法通过 `terraform refresh` 感知绑定是否被外部解绑 → 可接受，符合 framework action 定义。

## Migration Plan

- 纯新增 action，无 state 迁移需求
- 存量资源：无影响
- 文档：在 `action_tc_teo_bind_zone_to_plan.md` 中提供示例；`make doc` 生成 `website/docs/a/` 文档
- 回滚：删除新增文件及 registry.go / provider.md 中的登记项即可

## Open Questions

- 无