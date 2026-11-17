## ADDED Requirements

### Requirement: 资源提供云防火墙共享角色配置能力

系统 SHALL 提供 Terraform 资源 `tencentcloud_fwm_cfw_share_role`，用于声明式管理某个集团成员的云防火墙共享角色（`cfw_share_role`，取值 `sharer`/`user`/`none`）字段的期望状态。该资源为配置型资源（RESOURCE_KIND_CONFIG）：不创建/不销毁实体成员，仅管理该字段的读取与更新。

#### Scenario: 资源 schema 字段定义

- **WHEN** 用户编写 `tencentcloud_fwm_cfw_share_role` 资源配置
- **THEN** schema SHALL 包含 `member_id`（required）、`app_id`（required 且 ForceNew）、`cfw_share_role`（required，sharer/user/none）、`cfw_sharer_app_id`（computed，共享者 AppId，仅 user 角色有值）
- **AND** 资源描述 SHALL 为中文并带云产品名 fwm，字段 description SHALL 为英文

#### Scenario: 资源标识与 import

- **WHEN** 资源被创建或导入
- **THEN** 资源 ID SHALL 采用 `app_id + tccommon.FILED_SP + member_id` 复合形式
- **AND** import SHALL 支持 `app_id:member_id`（经 FILED_SP 分隔）两段式，Read 需 app_id 定位成员

### Requirement: Create 为 empty 实现

资源 Create SHALL 不创建任何实体成员（成员本就存在于集团中）。Create SHALL 仅执行 `d.SetId(app_id + FILED_SP + member_id)`，随后转调 Update 将目标共享角色同步到云上。

#### Scenario: 首次 apply 设置共享角色

- **WHEN** 用户首次 apply 一个 `tencentcloud_fwm_cfw_share_role` 资源，指定 member_id、app_id 和 cfw_share_role
- **THEN** provider SHALL 设置资源 ID 为 `app_id#member_id`
- **AND** SHALL 调用 Update 流程，将 cfw_share_role 期望值同步到云上

### Requirement: Delete 为 empty 实现

资源 Delete SHALL 不销毁任何实体成员。Delete SHALL 直接返回 nil，不调用任何云 API（共享角色只能置为 none，不能"删除"成员）。

#### Scenario: 资源销毁

- **WHEN** 用户执行 `terraform destroy` 销毁 `tencentcloud_fwm_cfw_share_role` 资源
- **THEN** provider SHALL 不调用任何云 API，直接返回成功
- **AND** 云上的成员与其实际共享角色保持不变

### Requirement: Read 用 AppId 过滤定位成员

Read SHALL 通过 `DescribeOrganMembers` 接口，以 `Filters=[{Name:"AppId", Values:[appId], OperatorType:1}]`（1=精确匹配）过滤定位成员。Provider MUST NOT 传 `Filters[{Name:"MemberId"}]`（接口过滤名白名单不含 MemberId，传错会报错）。响应返回后 SHALL 从 `Members` 列表按 MemberId 匹配出目标成员，回读其 `cfw_share_role` 与 `cfw_sharer_app_id`。

#### Scenario: 成员存在时回读共享角色

- **WHEN** Read 调用 DescribeOrganMembers 且返回的 Members 中存在 MemberId 匹配的成员
- **THEN** provider SHALL 回读该成员的 `cfw_share_role` 与 `cfw_sharer_app_id`（非 nil 才 set）
- **AND** SHALL 设置 `member_id`、`app_id` 到 state

#### Scenario: 成员不存在时清空资源

- **WHEN** Read 返回 Members 为空，或列表中无 MemberId 匹配的成员
- **THEN** provider SHALL 先 `log.Printf("[CRUD] fwm_cfw_share_role id=%s", d.Id())`
- **AND** 然后 `d.SetId("")`，返回 nil

### Requirement: Update 全量设置共享角色

Update SHALL 调用 `ModifyCfwShareRole` 接口，请求仅包含 `CfwShareRole`（目标角色 sharer/user/none）与 `UserMemberId`（成员 MemberId）两个字段。Provider MUST NOT 传 `SharerMemberId`（服务端在设为 user 时自动定位集团内唯一共享者）。写操作 SHALL 同步返回（无任务 ID，不轮询）。

#### Scenario: 目标角色与当前一致时跳过写操作

- **WHEN** Update 时当前 cfw_share_role 与目标 cfw_share_role 相同
- **THEN** provider SHALL 跳过 ModifyCfwShareRole 调用，直接进入 Read 回读

#### Scenario: 直接变更到目标角色

- **WHEN** Update 时当前角色与目标角色不同，且不构成 sharer↔user 直转
- **THEN** provider SHALL 调用一次 ModifyCfwShareRole（CfwShareRole=target, UserMemberId=member_id）

### Requirement: sharer↔user 经 none 中转

sharer 与 user 之间禁止直接互转（服务端 ValidateCfwShareRoleEdit 拒绝 sharer→user、user→sharer 直接变更）。当目标角色与当前角色不同且二者都属于 {sharer, user} 时，Update SHALL 先调用一次 ModifyCfwShareRole 设 `none`，再调用一次设目标角色。中转过程中任一步失败 SHALL 整体失败并保留原 state。

#### Scenario: sharer 转 user 经 none 中转

- **WHEN** 当前角色为 sharer，目标角色为 user
- **THEN** provider SHALL 先调用 ModifyCfwShareRole（CfwShareRole=none）
- **AND** 成功后再调用 ModifyCfwShareRole（CfwShareRole=user）

#### Scenario: user 转 sharer 经 none 中转

- **WHEN** 当前角色为 user，目标角色为 sharer
- **THEN** provider SHALL 先调用 ModifyCfwShareRole（CfwShareRole=none）
- **AND** 成功后再调用 ModifyCfwShareRole（CfwShareRole=sharer）

#### Scenario: 中转中间步失败保留原 state

- **WHEN** sharer↔user 中转的第一步（设 none）成功，但第二步（设目标角色）失败
- **THEN** provider SHALL 返回错误，保留原 state
- **AND** 下次 apply 将从 none 出发继续收敛到目标角色

### Requirement: 错误透传与重试

Update 的 ModifyCfwShareRole 调用 SHALL 用 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 包装，失败用 `tccommon.RetryError(e)` 包装返回。Read 的 DescribeOrganMembers 调用 SHALL 用 `tccommon.ReadRetryTimeout` 包装。云 API 返回的业务校验错误（如 InvalidParameter、UnauthorizedOperation）SHALL 直接透传给用户，不在 provider 层前置拦截。

#### Scenario: 非管理员调用被拒

- **WHEN** 调用账号非集团管理员，调用 DescribeOrganMembers 或 ModifyCfwShareRole
- **THEN** 云 API 返回 UnauthorizedOperation，provider SHALL 将该错误透传给用户

#### Scenario: sharer 前置约束不满足

- **WHEN** 设为 sharer 但成员非管理员/未购 CFW/集团已存在 sharer/CFW 版本不符
- **THEN** 云 API 返回 InvalidParameter，provider SHALL 透传该错误

### Requirement: provider 注册与服务客户端

系统 SHALL 在 `tencentcloud/connectivity/client.go` 中注册 fwm 客户端（`UseFwmClient()` 方法，使用 `fwm/v20250611` 包）。系统 SHALL 在 `tencentcloud/provider.go` 中注册 `tencentcloud_fwm_cfw_share_role` 资源，并在 `tencentcloud/provider.md` 中新增资源条目。系统 SHALL 新增 `tencentcloud/services/fwm/` 服务目录，包含资源配置文件 `resource_tc_fwm_cfw_share_role_config.go`、服务层 `service_tencentcloud_fwm.go`、单元测试 `resource_tc_fwm_cfw_share_role_config_test.go` 与文档样例 `resource_tc_fwm_cfw_share_role_config.md`。

#### Scenario: 首次使用 fwm 客户端

- **WHEN** 资源 Read/Update 首次调用 `UseFwmClient()`
- **THEN** provider SHALL 惰性创建并缓存 fwm 客户端（仿 UseCfwClient，NewClientProfile(300)）

### Requirement: 单元测试使用 mock

单元测试 SHALL 使用 gomonkey 对云 API 进行 mock，不使用 terraform 测试套件，仅测试业务代码逻辑。测试 SHALL 覆盖：Create empty 实现、Read 成员存在/不存在、Update 直接变更、Update 经 none 中转、Update 目标与当前一致跳过写、Delete empty 实现。

#### Scenario: Update 中转逻辑的单元测试

- **WHEN** 测试用例设置当前角色为 sharer、目标角色为 user
- **THEN** mock 的 ModifyCfwShareRole 调用 SHALL 被触发两次，第一次 CfwShareRole=none，第二次 CfwShareRole=user