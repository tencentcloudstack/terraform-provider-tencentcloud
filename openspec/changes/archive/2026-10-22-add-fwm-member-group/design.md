## Context

腾讯云云防火墙（FWM，SDK 包 `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611`）已 vendor（go.mod 声明 `tencentcloud/fwm v1.3.197`），但 Provider 尚未接入任何 fwm 资源，`tencentcloud/connectivity/client.go` 中没有 `UseFwmClient` / `fwmConn`。本变更新增第一个 fwm 资源 `tencentcloud_fwm_member_group`，类型为 RESOURCE_KIND_CONFIG：管理「某集团成员归属哪个账号组」这一单调配置关系。

业务背景与硬约束（来自需求，实施时必须遵守）：
1. config 型，非 CRUD：成员与组本就存在，归属关系只能置空不能销毁实体。Create=空实现（`d.SetId(member_id)` 后转 Update），Delete=空实现（直接 return nil），只有 Read + Update 有逻辑。
2. 一个成员最多属一个组（DB 唯一索引 `uk_admin_appid_appid(admin_appid,appid)` 硬约束），`group_id` 是**单值字符串不是列表**。服务端 `AccountGroup` 为单对象（`*AccountGroupInfo`），无组时为 `null`。若平台 schema 把 AccountGroup 标成 list 是平台标注错误，按单对象处理。
3. Read 只能用 AppId 定位，**不能**传 MemberId 过滤：`DescribeOrganMembers` 的过滤映射只支持 AppId/SubAccountCount/NodeName/Role/CfwManaged/SgManaged/CfwShareRole/CfwInstanceId/PolicyAnalysisEnabled/AdminAppId/Keyword，不含 MemberId。传 MemberId 过滤器会报「filter name does not meet the requirements」。做法：schema 除 `member_id` 外必须暴露 `app_id`（required + forceNew），用 `Filters=[{Name:"AppId",Values:[appId],OperatorType:1}]`（1=EQUAL）过滤，再从 Members 列表按 `MemberId` 匹配目标成员。
4. 回读：`Members[].AccountGroup` 为单对象 `{GroupId, GroupName}`，无组时 `null`。Read 取 `AccountGroup.GroupId` 回填 `group_id`；`null` 时置空字符串（表示未分组），**不要**当成资源不存在。
5. Update 字段映射：**入参扁平、出参嵌套**。请求只有 `UserMemberId`（成员 MemberId，JSON 键 `UserMemberId`，与 Describe 的 `MemberId` 同概念异名）+ `GroupId`（顶层扁平，对应 state 的 `group_id`←`AccountGroup.GroupId`）。`group_id` 为空字符串时传 `GroupId=""`，服务端自动移出（软删关系 + member_count-1），无需额外 Delete。
6. `UserMemberId` 传的是 **MemberId 不是 AppId**：服务端 `GetMemberByMemberId` 按 `member_id = ?` 精确匹配，传 AppId 返回 `ResourceNotFound`。MemberId 形如 `mem-xxxx`，app_id 是数字字符串，两者不同。
7. 加入组前置硬约束：成员必须已纳管。`GroupId` 非空时服务端校验 `targetMember.CfwManaged == 1`，否则 `InvalidParameter`「成员未开启云防火墙纳管，不能加入账号组」。`group_id` 非空时必须存在该成员的 `tencentcloud_fwm_cfw_managed` 资源且 `cfw_managed=1`，**HCL 示例必须加 `depends_on` 指向 cfw_managed**。移出（GroupId=""）不做此校验。
8. 目标组必须存在且未删：`GroupId` 非空时按 group_id+admin_appid+status=0 查组，查不到返回 `ResourceNotFound`。建议 `group_id` 引用 `tencentcloud_fwm_account_group` 的 ID 形成隐式依赖。
9. 权限前置：仅集团管理员。`IsMemberAdmin` 在业务逻辑之前，非管理员/普通成员返回 `UnauthorizedOperation`。测试账号必须是集团管理员。
10. 同步性：写操作单事务完成（软删旧关系→member_count-1→插入恢复新关系→member_count+1→写变动记录→通知策略状态），成功即返回（无任务 ID），Read 立即可读，**不要** waitReady/轮询。
11. 与 cfw_managed 的反向联动（须写进 resource description）：关闭纳管（cfw_managed 1→0）时服务端会在同一事务内自动把成员移出组并 member_count-1。同时管理两个资源时，应先把 group_id 置空再关纳管，否则产生预期外 state 漂移。
12. 失败语义：成员不存在→`ResourceNotFound`；成员不属于该集团→`InvalidParameter`「成员不属于该集团」；未纳管→`InvalidParameter`（见7）；组不存在→`ResourceNotFound`。原样透传，不静默改 state。
13. 响应裁剪：Members 返回整个成员对象，本资源只关心 `AccountGroup.GroupId`。其余字段（nickname/uin/role/account_group_name/cfw_managed/sg_managed/cfw_share_role/cfw_instance_id/cfw_pay_status 等）标 computed 或不暴露，不进可比字段（会随集团变化漂移）。
14. schema：`member_id`（required）、`app_id`（required + forceNew）、`group_id`（optional，空=未分组，可写可清空）。description 中文、字段 description 英文；testValue：app_id=<成员AppId>、member_id=<成员MemberId>、group_id=<已存在组ID>（该成员需已纳管）。

## Goals / Non-Goals

**Goals:**
- 提供 `tencentcloud_fwm_member_group` 资源，支持读取/更新成员所属账号组。
- 接入 fwm 客户端（`UseFwmClient` / `fwmConn`），填齐 Provider 对新服务 fwm 的支持基座。
- 正确处理「未分组=空字符串而非资源消失」「单值而非列表」「只能 AppId 定位」等非标准业务语义。

**Non-Goals:**
- 不管理成员纳管（cfw_managed）本身，纳管由独立资源 `tencentcloud_fwm_cfw_managed` 处理。
- 不管理账号组的创建/删除，账号组由独立资源 `tencentcloud_fwm_account_group` 处理。
- 不暴露 Members 中的派生/漂移字段作为可比参数（避免 state 漂移）。

## Decisions

### D1：config 型 CRUD 骨架（参考 `cfw_edge_firewall_switch`）

参考 `tencentcloud/services/cfw/resource_tc_cfw_edge_firewall_switch.go`（同为 config 型）：
- Create：`d.SetId(<id>)` 后直接 `return resourceXxxUpdate(d, meta)`，不做实际创建。
- Delete：空实现 `return nil`。
- Read：调用 Describe 接口查现值并回填；查不到实体（成员不存在）才 `d.SetId("")`，但 `AccountGroup==null`（未分组）**不**清空，仅置 `group_id=""`。
- Importer：`State: schema.ImportStatePassthrough`，复合 id `app_id#member_id`（用 `tccommon.FILED_SP`=`#` 分隔）。

### D2：复合 ID 与 Read 定位

- 资源 id = `app_id` + `tccommon.FILED_SP`(`#`) + `member_id`，前者用于 `DescribeOrganMembers` 过滤定位，后者用于在 Members 列表内匹配目标成员。
- Read 解析 `d.Id()` 得到 app_id 与 member_id，构造 `Filters=[{Name:"AppId",Values:[appId],OperatorType:1}]` 查询；遍历 `Members` 用 `MemberId` 匹配，匹配不到视为成员不存在（`d.SetId("")`）。
- 替代方案：直接传 MemberId 过滤——已否决，过滤映射不支持，会报错。

### D3：group_id 单值语义

- schema `group_id` 用 `schema.TypeString`，**不**用 list。Read 时 `AccountGroup != nil` 取 `GroupI else 置 `""`。Update 时 `group_id` 直接传 `ModifyMemberGroupRequest.GroupId`，空字符串表示移出。

### D4：Update 跳过逻辑

- Read 后比较目标 `group_id` 与现状：若 `d.HasChange("group_id")` 才调用 `ModifyMemberGroup`，一致则跳过写操作，减少无谓调用与 member_count 抖动。

### D5：fwm 客户端接入

- 在 `tencentcloud/connectivity/client.go`：import `fwm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"`；结构体加 `fwmConn *fwm.Client`；新增 `UseFwmClient()`（仿 `UseCfwClient`，`NewClientProfile(300)`）。该步由 tfpacer 代码生成流程处理，本变更的 tasks 中列出但实际接入在实施阶段。

### D6：不使用 _extension.go、不写注释头、统一资源名

- 资源文件 `resource_tc_fwm_member_group_config.go`，开头不加注释；日志/错误描述统一用资源名 `fwm_member_group`（小写蛇形）。

## 云 API 结构摘录（实施阶段无需回读 vendor 大文件）

### DescribeOrganMembers（读）— fwm/v20250611/models.go

请求 `DescribeOrganMembersRequest`（models.go:1895）：
```go
type DescribeOrganMembersRequest struct {
    Filters []*CommonFilter `json:"Filters,omitnil,omitempty" name:"Filters"` // 过滤条件列表
    Limit   *int64          `json:"Limit,omitnil,omitempty" name:"Limit"`      // 分页大小，默认 20
    Offset  *int64          `json:"Offset,omitnil,omitempty" name:"Offset"`    // 分页偏移，默认 0
    By      *string         `json:"By,omitnil,omitempty" name:"By"`           // 排序字段，如 MemberCreateTime
    Order   *string         `json:"Order,omitnil,omitempty" name:"Order"`      // asc/desc
}
```
响应 `DescribeOrganMembersResponse`（models.go:1949）→ `Response *DescribeOrganMembersResponseParams`：
```go
type DescribeOrganMembersResponseParams struct {
    TotalCount *int64             `json:"TotalCount"` // 集团成员总数
    Members    []*OrganMemberItem `json:"Members"`    // 集团成员列表
    RequestId  *string            `json:"RequestId"`
}
```
`CommonFilter`（models.go:342）：
```go
type CommonFilter struct {
    Name        *string  `json:"Name"`        // 筛选字段名（AppId/SubAccountCount/NodeName/Role/CfwManaged/SgManaged/CfwShareRole/CfwInstanceId/PolicyAnalysisEnabled/AdminAppId/Keyword）
    Values      []*string `json:"Values"`     // 筛选值列表
    OperatorType *int64  `json:"OperatorType"` // 1=精确匹配,7=in,9=模糊匹配
}
```
`OrganMemberItem`（models.go:4370）关键字段：
```go
type OrganMemberItem struct {
    MemberId  *string           `json:"MemberId"`   // 成员 ID（如 mem-xxxx）
    AppId     *string           `json:"AppId"`      // 成员账号 AppId（数字字符串）
    Uin       *string           `json:"Uin"`
    Nickname  *string           `json:"Nickname"`
    AccountGroup *AccountGroupInfo `json:"AccountGroup"` // 所属账户组（单对象，无组时 null）
    CfwManaged   *int64         `json:"CfwManaged"`   // 0-未纳管,1-已纳管
    // ... 其余字段（SubAccountCount/NodeName/Role/RoleDisplay/CfwShareRole/CfwShareRoleDisplay/
    //     CfwSharerAppId/CfwInstanceId/PolicyAnalysisEnabled/MemberCreateTime/JoinType/
    //     CfwPayStatus/CfwCapable/SgManaged/IsCfwPostPay）标 computed 或不暴露
}
```
`AccountGroupInfo`（models.go:70）：
```go
type AccountGroupInfo struct {
    GroupId   *string `json:"GroupId"`   // 账户组ID
    GroupName *string `json:"GroupName"` // 账户组名称
}
```
客户端方法（client.go:1328）：`func (c *Client) DescribeOrganMembers(request) (*DescribeOrganMembersResponse, error)`，另有 `DescribeOrganMembersWithContext`。

### ModifyMemberGroup（写/更新）— fwm/v20250611/models.go

请求 `ModifyMemberGroupRequest`（models.go:3621）：
```go
type ModifyMemberGroupRequest struct {
    UserMemberId *string `json:"UserMemberId,omitnil,omitempty" name:"UserMemberId"` // 目标成员 AppId（注释文字如此，但业务上实为成员 MemberId，见 D/需求6）
    GroupId      *string `json:"GroupId,omitnil,omitempty" name:"GroupId"`           // 账户组ID，空字符串表示移出
}
```
响应 `ModifyMemberGroupResponse`（models.go:3657）→ `Response *ModifyMemberGroupResponseParams`：
```go
type ModifyMemberGroupResponseParams struct {
    RequestId *string `json:"RequestId"` // 仅 RequestId，无任务 ID（同步）
}
```
客户端方法（client.go:2472）：`func (c *Client) ModifyMemberGroup(request) (*ModifyMemberGroupResponse, error)`，另有 `ModifyMemberGroupWithContext`。

> 注：`UserMemberId` 字段在 SDK 注释里写「目标成员 AppId」，但服务端 `GetMemberByMemberId` 按 `member_id = ?` 精确匹配，业务上必须传成员 **MemberId**（`mem-xxxx`），传 AppId 会返回 `ResourceNotFound`。实施时传 `member_id`。

## 代码风格参照文件结构

### 参照一：config 型骨架 `tencentcloud/services/cfw/resource_tc_cfw_edge_firewall_switch.go`
- `ResourceTencentCloudCfwEdgeFirewallSwitch()`：Schema + Create/Read/Update/Delete + Importer。
- Create：`d.SetId(publicIp)` 后 `return resourceTencentCloudCfwEdgeFirewallSwitchUpdate(d, meta)`。
- Read：调用 service 层 `DescribeCfwEdgeFirewallSwitchById`，nil 时 `log.Printf("[WARN]...")` + `d.SetId("")`；非 nil 逐字段判空再 `d.Set`。
- Update：`resource.Retry(tccommon.WriteRetryTimeout, ...)`，retry 内调 API，`e != nil` 用 `tccommon.RetryError(e)`。
- Delete：空实现 `return nil`。

### 参照二：资源代码风格 `tencentcloud/services/igtm/resource_tc_igtm_strategy.go`
- 包名 = 服务名（`fwm`）。
- import 顺序：标准库→terraform SDK→tencentcloud SDK（`fwmv20250611 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"`）→`tccommon`→`helper`。
- `ResourceTencentCloudFwmMemberGroup()` 返回 `*schema.Resource`。
- 延迟函数：`defer tccommon.LogElapsed("resource.tencentcloud_fwm_member_group.create")()`、`defer tccommon.InconsistentCheck(d, meta)()`。

### 注册：`tencentcloud/provider.go`
- 资源 map 项：`"tencentcloud_fwm_member_group": fwm.ResourceTencentCloudFwmMemberGroup(),`（参考 `tencentcloud_igtm_strategy` 行，provider.go:2784 附近）。
- 注意：fwm 服务目录在实施后才存在，注册前需先建服务目录与文件。

## Risks / Trade-offs

- **[AppId 过滤返回多个成员需二次匹配]** → Read 用 `MemberId` 在 Members 列表内精确匹配；匹配不到视为不存在。AppId 一般唯一，但用 AppId 是因为过滤映射不支持 MemberId。
- **[未分组 vs 资源消失歧义]** → `AccountGroup == null` 时 `group_id=""` 且**不** `d.SetId("")`；仅当整条成员记录查不到（Members 内无匹配 MemberId）才清空 id。需在 Read 内严格区分两个分支并各打日志。
- **[并发纳管与加组顺序]** → HCL 强制 `depends_on = [tencentcloud_fwm_cfw_managed.xxx]`；description 说明反向联动避免 state 漂移。
- **[SDK 字段 UserMemberId 注释误导]** → design/spec 明确传 MemberId，实施者勿被注释误导传 AppId。
- **[同步写无需轮询]** → Update 成功即 Read，不引入 waitReady，避免无谓重试与超时误判。
- **[fwm 客户端首次接入]** → client.go 改动属于基座性变更，需保证不影响其他服务；仿 `UseCfwClient` 写法。