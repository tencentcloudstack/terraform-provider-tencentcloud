## Why

腾讯云容灾中心（BDRC）已在云 API（`bdrc/v20260330`）发布了"容灾保护组"的完整 CRUD 能力（`CreateDisasterRecoveryProtectGroup`、`DescribeDisasterRecoveryProtectGroups`、`ModifyProtectGroupAttribute`、`DeleteDisasterRecoveryProtectGroups`），用于管理云盘 / 实例 / 文件存储的容灾保护组。Terraform Provider 尚未暴露该能力，用户只能在控制台或调用云 API 手动管理保护组，无法实现基础设施代码化（IaC）。

本 change 新增通用型（general）资源 `tencentcloud_bdrc_disaster_recovery_protect_group`，让用户在 HCL 中声明式地创建、查询、改名、删除容灾保护组，由 Provider 调用 BDRC 云 API 完成整个生命周期管理。

## What Changes

- 新增通用型资源 `tencentcloud_bdrc_disaster_recovery_protect_group`，完整实现 Create / Read / Update / Delete 回调。
- 新增 Schema 字段，严格按 `CreateDisasterRecoveryProtectGroup` 接口入参映射：`site_pair_id`（必填，ForceNew）、`protect_group_type`（必填，ForceNew）、`recovery_point_objective`（必填，ForceNew）、`protect_group_name`（可选，可变更）、`data_direction`（可选，ForceNew）。
- 新增只读 Computed 字段，来源于 `DescribeDisasterRecoveryProtectGroups` 返回的 `ProtectGroup` 结构（如 `app_id`、`site_pair_name`、`source_region`、`source_zone`、`source_vpc`、`target_region`、`target_zone`、`target_vpc`、`copy_type`、`disaster_recovery_type`、`peer_cloud_name`、`create_from`、`life_state`、`account_uin`、`sub_account_uin`、`create_time`、`modify_time`、`bind_protected_resource_count`、`error_recovery_point_objective_count`、`protected_resource_status_set`）。列表型数据展开平铺到顶层 schema，不引入 `protect_group_set` 这一层嵌套。
- 资源 ID 设为 `ProtectGroupId`（单字段，无需复合分隔符）。支持 import（`schema.ImportStatePassthrough`）。
- `site_pair_id`、`protect_group_type`、`recovery_point_objective`、`data_direction` 标记为 `ForceNew`，因为 `ModifyProtectGroupAttribute` 仅支持修改 `protect_group_name`。Update 回调中将这些字段加入 `immutableArgs`，若发生变更则返回 error。
- 新增 BDRC 客户端方法 `UseBdrcV20260330Client()` 及 `bdrcv20260330Conn` 字段到 `tencentcloud/connectivity/client.go`（bdrc 为全新云产品）。
- 新增服务层 `tencentcloud/services/bdrc/service_tencentcloud_bdrc.go`，封装 `DescribeDisasterRecoveryProtectGroupById` 查询方法（按 `ProtectGroupId` + `ProtectGroupType` 调用 `DescribeDisasterRecoveryProtectGroups`，内部自动分页 `Limit=100`）。
- 所有 SDK 调用用 `resource.Retry(...)` 包装（Read 用 `tccommon.ReadRetryTimeout`，Write/Delete 用 `tccommon.WriteRetryTimeout`），错误经 `tccommon.RetryError(e)` 包装。
- 新增 `.md` 资源文档 + 单元测试 `_test.go`（使用 gomonkey mock 云 API，不使用 terraform 验收测试套件）。
- 在 `tencentcloud/provider.go` 注册新资源，并在 `tencentcloud/provider.md` 新增 BDRC 段。

## Capabilities

### New Capabilities

- `bdrc-disaster-recovery-protect-group-resource`: 新增 `tencentcloud_bdrc_disaster_recovery_protect_group` 资源的 schema、CRUD（Create / Read / Update / Delete）行为、ID 约定、字段约束（ForceNew / immutableArgs）、客户端方法、服务层、文档与测试规范，作为容灾保护组能力的 IaC 入口。

### Modified Capabilities

<!-- 不修改任何已有 capability 的 requirement，仅在 provider.go / provider.md / connectivity/client.go 增加注册与客户端方法。 -->

## Impact

- 新文件：
  - `tencentcloud/services/bdrc/resource_tc_bdrc_disaster_recovery_protect_group.go`
  - `tencentcloud/services/bdrc/resource_tc_bdrc_disaster_recovery_protect_group.md`
  - `tencentcloud/services/bdrc/resource_tc_bdrc_disaster_recovery_protect_group_test.go`
  - `tencentcloud/services/bdrc/service_tencentcloud_bdrc.go`
  - `website/docs/r/bdrc_disaster_recovery_protect_group.html.markdown`（由 `make doc` 生成）
- 既有文件改动：
  - `tencentcloud/connectivity/client.go`：新增 `bdrcv20260330` 包导入、`bdrcv20260330Conn` 字段、`UseBdrcV20260330Client()` 方法。
  - `tencentcloud/provider.go`：新增 `bdrc` 服务包导入 + 资源注册行 `"tencentcloud_bdrc_disaster_recovery_protect_group": bdrc.ResourceTencentCloudBdrcDisasterRecoveryProtectGroup()`。
  - `tencentcloud/provider.md`：新增 Business Disaster Recovery Center(BDRC) 段并追加资源名，使 gendoc 索引可识别。
- 依赖：现有 SDK `tencentcloud-sdk-go/tencentcloud/bdrc/v20260330` 已包含全部四个接口及对应 model，无需升级 SDK。
- 不修改任何既有资源/数据源 schema，对现有用户零影响。
