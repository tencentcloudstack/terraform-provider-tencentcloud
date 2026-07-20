## Context

TEO（EdgeOne）站点创建后需要绑定到已有套餐（Plan）才能正式使用。腾讯云 TEO SDK（`teo/v20220901`）提供了 `BindZoneToPlan` 接口，入参为 `ZoneId` 与 `PlanId`，返回值仅包含 `RequestId`，无业务数据，且为同步接口。当前 Terraform Provider 中 `tencentcloud_teo_zone` 资源不包含"绑定套餐"动作，用户必须借助控制台或 SDK 手动完成绑定，导致声明式运维流程断裂。

当前状态：
- vendor 中 `teov20220901.BindZoneToPlanRequest` 已存在，字段：`ZoneId`、`PlanId`（均为 `*string`）
- vendor 中 `teov20220901.Client.BindZoneToPlan` 与 `BindZoneToPlanWithContext` 已存在
- 返回 `BindZoneToPlanResponse`，其 `Response` 仅含 `RequestId`，无业务数据
- provider 已有同类一次性 operation 资源（如 `tencentcloud_teo_confirm_origin_acl_update_operation`），可作为代码风格参考

约束：
- RESOURCE_KIND_OPERATION：一次性操作，操作完不需要记录任何状态
- Create 调用 C 接口（`BindZoneToPlan`），Read/Update/Delete 为空
- 资源文件命名格式：`resource_tc_teo_bind_zone_to_plan_operation.go`
- 不需要 `_extension.go` 文件
- 使用 `tccommon.WriteRetryTimeout` 作为超时时间添加 retry 处理

## Goals / Non-Goals

**Goals:**
- 通过 Terraform 声明式地完成 TEO 站点与套餐的绑定
- 遵循 provider 既有的一次性 operation 资源模式（参考 `tencentcloud_teo_confirm_origin_acl_update_operation`）
- 提供基于 gomonkey 的单元测试，覆盖成功、API 错误、no-op Read/Delete、schema 校验等场景
- 在 `provider.go` / `provider.md` 中完成资源注册

**Non-Goals:**
- 不在资源中暴露绑定后的套餐查询/校验（API 本身不返回业务数据，且 operation 资源 Read 为 no-op）
- 不提供 Update 操作（一次性操作，参数变更只能重建）
- 不提供 Importer（一次性 operation 资源不可导入）
- 不修改 `tencentcloud_teo_zone` 等现有资源

## Decisions

### Decision 1: 采用一次性 operation 资源而非在 zone 资源上扩展

**选择**：新建独立资源 `tencentcloud_teo_bind_zone_to_plan`（文件 `resource_tc_teo_bind_zone_to_plan_operation.go`）。

**备选**：在 `tencentcloud_teo_zone` 资源上新增 `plan_id` 字段并在 Create 中调用 `BindZoneToPlan`。

**理由**：
- `BindZoneToPlan` 语义上是"动作"而非"属性"，且是一次性操作，操作完不需要记录任何状态，符合 RESOURCE_KIND_OPERATION 的定义
- 独立资源不会污染 `tencentcloud_teo_zone` 已有的 Create/Read/Update/Delete 流程，保持向后兼容
- 与 provider 中既有 operation 资源（`tencentcloud_teo_confirm_origin_acl_update_operation`、`tencentcloud_teo_identify_zone_operation` 等）保持一致风格

### Decision 2: id 使用 helper.BuildToken()

**选择**：Create 调用 `BindZoneToPlan` 成功后，使用 `helper.BuildToken()` 生成随机 token 作为资源 id。

**备选**：使用 `zone_id` 与 `plan_id` 组合作为复合 id。

**理由**：
- `BindZoneToPlan` 返回值仅含 `RequestId`，无业务 id 可用
- operation 资源不记录状态、不可 import、Read 为 no-op，id 仅用于 Terraform state 标识，无需具备业务语义
- 与同类 operation 资源（`tencentcloud_teo_confirm_origin_acl_update_operation`）一致

### Decision 3: retry 包装 BindZoneToPlan 调用

**选择**：使用 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 包装 `BindZoneToPlanWithContext` 调用，失败时用 `tccommon.RetryError(e)` 包装返回。

**理由**：
- 符合项目规范："调用云API接口时，需要以 tccommon.ReadRetryTimeout 作为超时时间，添加 retry 处理"
- 写操作使用 `tccommon.WriteRetryTimeout`
- retry 块内只执行接口调用，`d.SetId()` 放在 retry 块外、错误处理后

### Decision 4: schema 字段全部 Required + ForceNew

**选择**：`zone_id`、`plan_id` 均为 `Required: true, ForceNew: true, TypeString`。

**理由**：
- 一次性 operation 资源无 Update，所有参数变更必须重建（ForceNew）
- 两个参数在 `BindZoneToPlanRequest` 中均为必需入参

### Decision 5: BindZoneToPlan 为同步接口，无需轮询

**选择**：调用 `BindZoneToPlan` 成功后直接 `d.SetId()` 并返回，不调用 Read 接口轮询。

**理由**：
- `BindZoneToPlan` 为同步接口（需求描述未标注"异步接口"），返回即代表操作完成
- operation 资源 Read 本身为 no-op，无业务状态可轮询

## Risks / Trade-offs

- **Risk**：用户对同一站点重复 apply 该资源，会重复调用 `BindZoneToPlan`，API 可能返回 `InvalidParameter.ZoneHasBeenBound` → **Mitigation**：由云 API 返回错误，retry 机制会在非可重试错误时停止；用户可通过 `terraform taint` 或重建资源处理。属于预期行为，operation 资源普遍如此。
- **Risk**：`zone_id` 或 `plan_id` 不存在时 API 返回 `InvalidParameter.ZoneNotFound` / `InvalidParameter.PlanNotFound` → **Mitigation**：由 `tccommon.RetryError` 判定为非可重试错误并返回，提示用户检查参数。
- **Trade-off**：operation 资源不记录绑定后的实际状态，用户无法通过 `terraform refresh` 感知绑定是否被外部解绑 → 可接受，符合 RESOURCE_KIND_OPERATION 定义。

## Migration Plan

- 纯新增资源，无 state 迁移需求
- 存量资源：无影响
- 文档：在 `resource_tc_teo_bind_zone_to_plan_operation.md` 中提供示例；`make doc` 生成 `website/docs/` 文档
- 回滚：删除新增文件及 provider.go 注册项即可

## Open Questions

- 无
