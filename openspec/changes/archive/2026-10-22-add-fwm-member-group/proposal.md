## Why

腾讯云云防火墙（FWM，fwm/v20250611）支持将集团成员划入账号组进行统一策略管理，但目前 Terraform Provider 无法管理「成员归属哪个账号组」这一关系。用户需要以基础设施即代码的方式声明每个集团成员所属的账号组，以便在成员纳管后自动完成分组编排、在变更时平滑迁移分组。由于一个成员最多归属一个组（DB 唯一索引硬约束），该关系是一个单值配置，适合用 config 型资源表达。

## What Changes

- 新增 Terraform 资源 `tencentcloud_fwm_member_group`（RESOURCE_KIND_CONFIG，文件 `resource_tc_fwm_member_group_config.go`），用于管理「某集团成员所属账号组」这一单调配置关系：
  - 资源标识为成员 `MemberId`，一个实例管一个成员；用 `app_id#member_id` 复合 id，导入用两段式（Read 只能用 AppId 定位）。
  - `group_id` 为单值字符串（空字符串表示未分组，可写可清空），业务语义来自 `AccountGroup.GroupId`。
  - Create 为空实现（`d.SetId(member_id)` 后转 Update），Delete 为空实现（成员与组本就存在，归属只能置空，不销毁实体）；仅 Read + Update 有实际逻辑。
- Read：用 `DescribeOrganMembers` + `Filters=[{Name:"AppId",Values:[appId],OperatorType:1}]` 定位成员列表，再按 `MemberId` 匹配目标成员，回读 `AccountGroup.GroupId`（null 时置空字符串，表示未分组，不视为资源不存在）。
- Update：调用 `ModifyMemberGroup`，入参 `UserMemberId` 传成员 `MemberId`（非 AppId），`GroupId` 传目标组 id 或空字符串（移出）。目标与当前一致则跳过写操作。同步事务完成，无需轮询。
- 新增 `tencentcloud/services/fwm/` 服务目录、服务层 `service_tencentcloud_fwm.go`、资源文件与测试、`.md` 文档，并在 `tencentcloud/provider.go` / `provider.md` 注册资源；在 `tencentcloud/connectivity/client.go` 补充 `UseFwmClient()` / `fwmConn` 字段（fwm 是新服务，客户端尚未接入）。
- HCL 示例需体现业务硬约束：`group_id` 非空时该成员必须已纳管（`tencentcloud_fwm_cfw_managed` 资源且 `cfw_managed=1`，必须 `depends_on` 指向 cfw_managed）；`group_id` 建议引用账号组资源 id 形成隐式依赖。

## Capabilities

### New Capabilities
- `fwm-member-group`: 管理 FWM 集团成员所属账号组的配置型资源，支持按 AppId 定位成员、读取/更新成员归属的账号组、以空 group_id 移出账号组。

### Modified Capabilities
<!-- 无需修改已有 spec -->

## Impact

- **新增代码**：`tencentcloud/services/fwm/resource_tc_fwm_member_group_config.go`、`service_tencentcloud_fwm.go`、资源 `_test.go`、资源 `.md`；`website/docs/r/tencentcloud_fwm_member_group.html.markdown` 由 `make doc` 生成。
- **修改代码**：`tencentcloud/provider.go`、`tencentcloud/provider.md`（注册资源）；`tencentcloud/connectivity/client.go`（新增 `UseFwmClient` 与 `fwmConn` 字段、import）。
- **云 API 依赖**：`DescribeOrganMembers`（读，fwm/v20250611）、`ModifyMemberGroup`（写，fwm/v20250611）。SDK 已 vendor（go.mod 声明 `tencentcloud/fwm v1.3.197`）。
- **业务联动**：与 `tencentcloud_fwm_cfw_managed` / `tencentcloud_fwm_account_group` 存在前置依赖与反向联动（关闭纳管会自动移出组），需在 description 中说明。