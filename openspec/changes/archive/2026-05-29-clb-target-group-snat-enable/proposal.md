## Why

`tencentcloud_clb_target_group` 资源当前未暴露 `SnatEnable` 参数，但腾讯云 CLB SDK 的 `CreateTargetGroup`、`DescribeTargetGroupList`（响应 `TargetGroupInfo` 中包含 `SnatEnable`）以及 `ModifyTargetGroupAttribute` 三个 API 均已原生支持该字段。用户在使用 Terraform 管理 CLB 目标组时，无法通过 IaC 控制源 IP 替换（SNAT）行为，必须借助控制台或外部脚本兜底，影响自动化闭环。本变更将该能力打通到 Provider，使目标组的 SNAT 配置可被 Terraform 完整地创建、读取与变更。

## What Changes

- 在 `tencentcloud_clb_target_group` 资源 schema 中新增 `snat_enable` 字段（Optional + Computed，类型 `Bool`）。
- `Create` 流程中：当用户显式设置 `snat_enable` 时，将值传入 `CreateTargetGroupRequest.SnatEnable`。
- `Read` 流程中：从 `DescribeTargetGroupList` 返回的 `TargetGroupInfo.SnatEnable` 回写到 state。
- `Update` 流程中：当 `snat_enable` 发生变化时，通过 `ModifyTargetGroupAttribute` 接口同步到云端。
- 扩展 service 层 `ClbService.CreateTargetGroup` 入参以承载 `snatEnable *bool`（非破坏性扩展，nil 表示不传）。
- 同步更新资源示例 `resource_tc_clb_target_group.md`，并通过 `make doc` 重新生成 `website/docs/` 下的文档。
- 补充验收测试覆盖 SNAT 启用/关闭以及更新场景。

不涉及破坏性变更：仅向已有资源新增 Optional 字段，旧 TF 配置与 state 完全兼容。

## Capabilities

### New Capabilities
<!-- 无 -->

### Modified Capabilities
- `clb-target-group-query`: 在目标组的查询 / 资源 CRUD 行为中纳入 `SnatEnable` 字段的读写要求（Create/Read/Update 三处均须支持）。

## Impact

- 代码：
  - `tencentcloud/services/clb/resource_tc_clb_target_group.go`（schema、Create/Read/Update）
  - `tencentcloud/services/clb/service_tencentcloud_clb.go`（`CreateTargetGroup` 方法签名扩展）
  - `tencentcloud/services/clb/resource_tc_clb_target_group_test.go`（新增/扩展验收测试）
  - `tencentcloud/services/clb/resource_tc_clb_target_group.md`（示例 HCL 增加 `snat_enable`）
  - `website/docs/r/clb_target_group.html.markdown`（由 `make doc` 自动生成，无需手写）
- API：使用既有 SDK 字段，无需升级 SDK 版本。
- 依赖：vendor 中的 `tencentcloud-sdk-go/tencentcloud/clb/v20180317` 已包含 `SnatEnable`，无需更新。
- 兼容性：纯增量 Optional 字段，不影响存量用户。
