## Why

TEO（EdgeOne）站点在创建后需要绑定到已有套餐（Plan）才能正式生效。当前 Terraform Provider 没有暴露"为站点绑定套餐"这一操作，用户必须通过控制台或 SDK 手动调用 `BindZoneToPlan`，导致声明式运维流程断裂。TEO SDK（`teo/v20220901`）已提供 `BindZoneToPlan` 接口，可通过一个一次性 operation 资源将其纳入 Terraform 管理。

## What Changes

- 新增 Terraform operation 资源 `tencentcloud_teo_bind_zone_to_plan`（RESOURCE_KIND_OPERATION），用于调用 `BindZoneToPlan` 接口将未绑定套餐的站点绑定到已有套餐。
- 资源 schema 包含两个 Required、ForceNew、TypeString 参数：`zone_id`（站点 ID）和 `plan_id`（套餐 ID）。
- Create 阶段调用 `BindZoneToPlan` 接口（同步接口，返回值仅含 RequestId，无业务数据），调用成功后使用 `helper.BuildToken()` 设置资源 id；Read/Delete 为 no-op；无 Update、无 Importer。
- 在 `provider.go` 与 `provider.md` 中注册该资源。
- 提供对应的 `.md` 文档与基于 gomonkey 的单元测试。

## Capabilities

### New Capabilities
- `teo-bind-zone-to-plan-operation`: 新增一次性 operation 资源 `tencentcloud_teo_bind_zone_to_plan`，调用 TEO `BindZoneToPlan` API 完成站点与套餐的绑定。

### Modified Capabilities
<!-- 无现有 capability 需要修改 -->

## Impact

- 代码：
  - `tencentcloud/services/teo/resource_tc_teo_bind_zone_to_plan_operation.go`（新增 operation 资源）
  - `tencentcloud/services/teo/resource_tc_teo_bind_zone_to_plan_operation_test.go`（新增 gomonkey 单元测试）
  - `tencentcloud/services/teo/resource_tc_teo_bind_zone_to_plan_operation.md`（新增资源文档）
  - `tencentcloud/provider.go`（注册资源）
  - `tencentcloud/provider.md`（登记资源名）
- 依赖：使用已 vendored 的 `tencentcloud-sdk-go` 中 `teov20220901.BindZoneToPlanRequest`，无需变更 vendor。
- 向后兼容：纯新增资源，不影响现有资源与 state。
- 文档：需要同步生成 `website/docs/` 文档（由 `make doc` 流程读取 `.md` 文件）。
