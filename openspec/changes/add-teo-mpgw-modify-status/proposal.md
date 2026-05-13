## Why

多路径网关（Multi-Path Gateway）在实际运维中经常需要在线启用/停用（`online`/`offline`），当前 `tencentcloud_teo_multi_path_gateway` 资源只暴露 `status` 为 Computed 字段，用户无法通过 Terraform 配置修改网关状态，必须借助控制台或 SDK 手动操作，导致声明式运维流程断裂。腾讯云 TEO SDK 已提供 `ModifyMultiPathGatewayStatus` 接口，可以通过扩展现有资源在 Update 阶段调用该接口实现状态切换。

## What Changes

- 将 `tencentcloud_teo_multi_path_gateway` 资源的 `status` 字段改为 `Optional + Computed`（向后兼容：未配置时仍由 Read 回填，不会破坏现有 TF 配置与 state）。
- 扩展 Update 操作：当 `status` 发生变化时，调用 `ModifyMultiPathGatewayStatus` API，传入 `ZoneId`、`GatewayId`、`GatewayStatus`（`online` / `offline`）。
- Update 操作在调用 `ModifyMultiPathGatewayStatus` 后，使用 `DescribeTeoMultiPathGatewayById` 轮询等待，直到网关 `Status` 达到目标值或中间态（如 `creating`）结束。
- 完善已有的 `resource_tc_teo_multi_path_gateway_test.go` 单测，新增覆盖 status 变更分支的用例。
- 更新资源文档 `resource_tc_teo_multi_path_gateway.md`，在示例中展示 `status` 字段的可选使用方式。

非破坏性：未设置 `status` 的既有配置仍然保持原有 Computed 行为，无 state drift。

## Capabilities

### New Capabilities
<!-- 无新增 capability，本次仅修改现有 capability -->

### Modified Capabilities
- `teo-multi-path-gateway-resource`: 将 `status` 字段从 Computed 变更为 Optional+Computed；Update 操作新增对 `ModifyMultiPathGatewayStatus` 接口的调用分支以支持启停网关；单元测试新增对应场景。

## Impact

- 代码：
  - `tencentcloud/services/teo/resource_tc_teo_multi_path_gateway.go`（schema 和 Update 逻辑）
  - `tencentcloud/services/teo/resource_tc_teo_multi_path_gateway_test.go`（新增 status 变更用例）
  - `tencentcloud/services/teo/resource_tc_teo_multi_path_gateway.md`（示例补充 status）
- 依赖：使用已 vendored 的 `tencentcloud-sdk-go` 中 `teov20220901.ModifyMultiPathGatewayStatusRequest`，无需变更 vendor。
- 向后兼容：新增 Optional 字段，不影响已有 state 与 TF 配置。
- 文档：需要同步更新 website docs（由自动生成流程读取 `.md` 文件）。
