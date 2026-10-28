## Why

TEO（EdgeOne）站点在创建后需要绑定到已有套餐（Plan）才能正式生效。当前 Terraform Provider 没有暴露"为站点绑定套餐"这一操作，用户必须通过控制台或 SDK 手动调用 `BindZoneToPlan`，导致声明式运维流程断裂。TEO SDK（`teo/v20220901`）已提供 `BindZoneToPlan` 接口，可通过一个一次性 action 将其纳入 Terraform 管理。

## What Changes

- 新增 Terraform 框架（terraform-plugin-framework）action `tencentcloud_teo_bind_zone_to_plan`，用于调用 `BindZoneToPlan` 接口将未绑定套餐的站点绑定到已有套餐。
- action schema 包含两个 Required、TypeString 参数：`zone_id`（站点 ID）和 `plan_id`（套餐 ID）。
- Invoke 阶段调用 `BindZoneToPlan` 接口（同步接口，返回值仅含 RequestId，无业务数据），使用 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 包装 `BindZoneToPlanWithContext` 调用，失败用 `tccommon.RetryError(e)` 包装返回；成功时不写入任何 state（action 天然不管理状态）。
- 在 `tencentcloud/framework/registry.go` 的 `actionFactories` 中注册该 action，并在 `tencentcloud/provider.md` 中 TEO 产品下新增 `Action` 子节登记该 action。
- 提供对应的 `.md` 文档（`action_tc_teo_bind_zone_to_plan.md`）与基于 gomonkey 的单元测试（覆盖 Metadata、Schema、Configure、Invoke 成功、Invoke API 错误、Invoke 缺参数等场景）。

## Capabilities

### New Capabilities
- `teo-bind-zone-to-plan-action`: 新增 framework action `tencentcloud_teo_bind_zone_to_plan`，调用 TEO `BindZoneToPlan` API 完成站点与套餐的绑定。

### Modified Capabilities
<!-- 无现有 capability 需要修改 -->

## Impact

- 代码：
  - `tencentcloud/services/teo/action_tc_teo_bind_zone_to_plan.go`（新增 framework action 实现）
  - `tencentcloud/services/teo/action_tc_teo_bind_zone_to_plan_test.go`（新增 gomonkey 单元测试）
  - `tencentcloud/services/teo/action_tc_teo_bind_zone_to_plan.md`（新增 action 文档）
  - `tencentcloud/framework/registry.go`（注册 action 工厂）
  - `tencentcloud/provider.go`（不注册 SDKv2 资源；action 由 framework 注册）
  - `tencentcloud/provider.md`（在 TEO 下新增 Action 子节登记 action）
- 依赖：使用已 vendored 的 `tencentcloud-sdk-go` 中 `teov20220901.BindZoneToPlanRequest`，无需变更 vendor。
- 向后兼容：原 SDKv2 一次性 operation 资源尚未发布，改为 framework action 不影响现有资源与 state。
- 文档：需要同步生成 `website/docs/a/` 文档（由 `make doc` 流程读取 `.md` 文件）。