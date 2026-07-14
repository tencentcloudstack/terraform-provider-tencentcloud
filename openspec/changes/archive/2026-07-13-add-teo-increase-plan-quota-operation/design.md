## Context

TEO 云 API `IncreasePlanQuota`（包 `teo/v20220901`）已在 vendor 中可用。该接口用于增购套餐配额，入参为 `PlanId`、`QuotaType`、`QuotaNumber`，出参为 `DealName`（订单号）。该接口是同步接口，无需异步轮询。

本资源为 **operation-type**（一次性操作），遵循已有 TEO operation 资源（如 `tencentcloud_teo_confirm_origin_acl_update_operation`）的实现模式。

## Goals / Non-Goals

**Goals:**
- 通过 `IncreasePlanQuota` API 触发配额增购，带重试。
- 生成资源 ID 使用 `helper.BuildToken()`。
- 返回 `deal_name` 作为 Computed 属性，提供订单号给用户。
- 输出资源 `.md` 样例文档、网站文档、单元测试和 Provider 注册。

**Non-Goals:**
- 无 Update 操作（operation 类型资源不支持 Update）。
- 无 Delete 操作（Read / Delete 均为空操作）。
- 无需异步轮询（接口为同步接口，非异步）。

## Decisions

### 资源名称和文件布局
- 资源名: `tencentcloud_teo_increase_plan_quota`
- 文件路径: `tencentcloud/services/teo/resource_tc_teo_increase_plan_quota_operation.go`
- 文档: `tencentcloud/services/teo/resource_tc_teo_increase_plan_quota_operation.md`
- 网站文档: `website/docs/r/teo_increase_plan_quota_operation.html.markdown`
- 构造函数: `ResourceTencentCloudTeoIncreasePlanQuotaOperation()`

### Schema 设计
- `plan_id` (TypeString, Required, ForceNew) — 套餐 ID
- `quota_type` (TypeString, Required, ForceNew) — 配额类型（site / precise_access_control_rule / rate_limiting_rule）
- `quota_number` (TypeInt, Required, ForceNew) — 配额数量（上限 100）
- `deal_name` (TypeString, Computed) — 订单号，来自 `IncreasePlanQuotaResponse.DealName`

### CRUD 映射
- **Create**: 调用 `IncreasePlanQuota` 在 `resource.Retry(WriteRetryTimeout, ...)` 内，成功后 `d.SetId(helper.BuildToken())`，nil-safe 设置 `deal_name`，然后调用 Read。
- **Read**: 空操作，返回 nil。
- **Delete**: 空操作，返回 nil。

### 客户端方法
- 使用 `meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().IncreasePlanQuota(request)` 调用 API。
- 带 `ctx` 的变体：`IncreasePlanQuotaWithContext(ctx, request)`。

### 参考实现
- 同一 teo 产品下的 `resource_tc_teo_confirm_origin_acl_update_operation.go`（带参数和 Computed 返回值的简单 operation）。

## Risks / Trade-offs

- [配额增购是不可逆操作] → Mitigation: 文档中明确说明，该资源执行后不可通过 Terraform destroy 回滚；Delete 仅移除 state。
- [接口可能因账户余额不足等原因失败] → Mitigation: 使用 retry 机制，失败时通过 `tccommon.RetryError` 包装错误返回给用户。

## Migration Plan

纯增量添加；无需迁移。新资源注册到现有 TEO 资源旁边。回滚 = 移除新增文件和 Provider 注册即可。