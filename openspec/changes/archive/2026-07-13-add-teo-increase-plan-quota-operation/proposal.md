## Why

当用户的 TEO（EdgeOne）企业版套餐绑定的站点数、Web 防护自定义规则精准匹配策略规则数或速率限制精准速率限制模块规则数达到套餐允许的配额上限时，需要通过 `IncreasePlanQuota` 云 API 增购对应配额。目前 Terraform Provider 不提供该操作资源，用户只能通过控制台手动操作，无法实现 IaC 自动化管理。

## What Changes

- 新增 **operation-type**（一次性操作）资源 `tencentcloud_teo_increase_plan_quota`，调用 `IncreasePlanQuota` 云 API 增购套餐配额。
- 入参：`plan_id`（套餐ID）、`quota_type`（配额类型）、`quota_number`（配额数量），所有入参均为 Required + ForceNew。
- 出参：`deal_name`（订单号），Computed。
- 资源的唯一 ID 由 `helper.BuildToken()` 自动生成。
- Create 生命周期调用 `IncreasePlanQuota` 并返回 `deal_name`；Read 和 Delete 生命周期为空操作（no-op）。
- 注册到 Provider、生成资源文档（`.md`）、网站文档（通过 `make doc`）、单元测试。

## Capabilities

### New Capabilities
- `teo-increase-plan-quota`: 一次性操作资源，用于增购 TEO 套餐配额，调用 `IncreasePlanQuota` 云 API 并返回订单号。

### Modified Capabilities
<!-- None: brand-new resource, no existing requirements change. -->

## Impact

- New file: `tencentcloud/services/teo/resource_tc_teo_increase_plan_quota_operation.go`
- New file: `tencentcloud/services/teo/resource_tc_teo_increase_plan_quota_operation.md`
- New file: `tencentcloud/services/teo/resource_tc_teo_increase_plan_quota_operation_test.go`
- New file: `website/docs/r/teo_increase_plan_quota_operation.html.markdown`（通过 `make doc` 生成）
- Modified: `tencentcloud/provider.go`（注册资源），`tencentcloud/provider.md`，`website/tencentcloud.erb`
- SDK: 使用已有的 `teo/v20220901` API `IncreasePlanQuota`（已在 vendor 中验证存在，无需 SDK 变更）。
- 无破坏性变更；纯增量添加。