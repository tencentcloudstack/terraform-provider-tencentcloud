## Why

云防火墙（fwm）共享角色（sharer/user/none）目前只能在腾讯云控制台手工配置，缺少 IaC 管理手段。集团管理员在多成员云防火墙规格共享场景下需要以声明式方式管理「某个集团成员的云防火墙共享角色」这一字段的期望状态，并随 Terraform 流程化变更、可审计、可重复执行。当前 provider 尚无任何 fwm 资源，无法满足该需求。

## What Changes

- 新增配置型（RESOURCE_KIND_CONFIG）Terraform 资源 `tencentcloud_fwm_cfw_share_role`，用于管理某个集团成员的云防火墙共享角色字段（sharer/user/none）的期望状态。
- 资源标识 = 目标成员的 MemberId；schema 同时暴露 `app_id`（required + ForceNew，用于 Read 定位）和 `member_id`（required，资源 ID）、`cfw_share_role`（required，sharer/user/none）。
- Create 为 empty 实现：仅 `d.SetId(member_id)` 后转调 Update（成员本就存在，不创建实体）。
- Delete 为 empty 实现：角色只能置为 none，不能"删除"实体。
- Read：通过 `DescribeOrganMembers` 接口，以 `Filters=[{Name:"AppId", Values:[appId], OperatorType:1}]`（1=精确匹配）定位成员，再从返回 Members 列表按 MemberId 匹配出目标成员，回读其 `cfw_share_role` 与 `cfw_sharer_app_id`。绝不可传 `Filters[{Name:"MemberId"}]`（接口不支持该过滤名）。
- Update：调用 `ModifyCfwShareRole` 接口，仅传 `UserMemberId`（成员 MemberId）+ `CfwShareRole`（sharer/user/none）。实现状态机中转逻辑：sharer↔user 禁止直接互转，必须先置 none 再置目标角色。
- 新增 fwm 服务目录 `tencentcloud/services/fwm/`，包含资源配置文件 `resource_tc_fwm_cfw_share_role_config.go`、`service_tencentcloud_fwm.go`、`resource_tc_fwm_cfw_share_role_config_test.go`、`resource_tc_fwm_cfw_share_role_config.md`。
- 在 `tencentcloud/connectivity/client.go` 中注册 fwm 客户端（`UseFwmClient`），在 `tencentcloud/provider.go` 与 `tencentcloud/provider.md` 中注册新资源。

## Capabilities

### New Capabilities
- `fwm-cfw-share-role-config`: 管理集团成员云防火墙共享角色（sharer/user/none）的配置型资源，包含 Read/Update 行为以及 sharer↔user 经 none 中转的状态机约束。

### Modified Capabilities
<!-- 无 -->

## Impact

- **新增代码**:
  - `tencentcloud/services/fwm/resource_tc_fwm_cfw_share_role_config.go`（资源 CRUD）
  - `tencentcloud/services/fwm/service_tencentcloud_fwm.go`（服务层，DescribeOrganMembers 封装）
  - `tencentcloud/services/fwm/resource_tc_fwm_cfw_share_role_config_test.go`（单元测试，mock 云 API）
  - `tencentcloud/services/fwm/resource_tc_fwm_cfw_share_role_config.md`（文档样例）
- **修改代码**:
  - `tencentcloud/connectivity/client.go`：新增 fwm client 字段、`UseFwmClient()` 方法、import
  - `tencentcloud/provider.go`：注册 `tencentcloud_fwm_cfw_share_role` 资源
  - `tencentcloud/provider.md`：新增资源文档条目
- **依赖**: 使用 vendor 中已有的 `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611` 包（已存在于 vendor，未提交）。
- **云 API**: 调用 `DescribeOrganMembers`（查询集团成员）、`ModifyCfwShareRole`（设置共享角色），均为同步接口，无任务 ID、无需轮询。
- **向后兼容**: 纯新增资源，不影响现有资源与 state。