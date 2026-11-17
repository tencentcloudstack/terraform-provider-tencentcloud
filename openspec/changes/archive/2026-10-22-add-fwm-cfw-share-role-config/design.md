## Context

云防火墙 fwm 产品当前在 Terraform Provider 中没有任何资源/数据源。本次新增一个**配置型资源（RESOURCE_KIND_CONFIG）** `tencentcloud_fwm_cfw_share_role`，用于声明式管理「某个集团成员的云防火墙共享角色」这一单个字段的期望状态。资源仅 Read + Update，Create/Delete 为 empty 实现。

云上行为特征：
- 共享角色取值 `sharer`（共享者）/ `user`（使用者）/ `none`（未设置）。
- 资源标识 = 目标成员 MemberId；Read 必须用 AppId 作为过滤条件定位成员（接口不支持按 MemberId 过滤）。
- Update 是同步全量设置，只传 `UserMemberId` + `CfwShareRole` 两个字段。
- sharer↔user 禁止直接互转，必须经 `none` 中转。
- 仅集团管理员可调用（非管理员返回 `UnauthorizedOperation`）。
- 设为 sharer 有前置约束（管理员角色、已购 CFW、集团内共享者唯一、CFW 高级/企业/旗舰版）；设为 none 也有校验（user→none 需无防火墙实例；sharer→none 需无使用者）。

vendor SDK 包：`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611`（已存在于 vendor，未提交 git）。

## Goals / Non-Goals

**Goals:**
- 提供 `tencentcloud_fwm_cfw_share_role` 资源，支持声明式读取与更新成员的 cfw_share_role 字段。
- 正确实现 sharer↔user 经 none 中转的状态机逻辑，避免单步变更被服务端拒绝。
- 在 provider 中注册 fwm 服务客户端与新资源，新增 fwm 服务目录。
- 提供基于 gomonkey mock 云 API 的单元测试。

**Non-Goals:**
- 不创建/不销毁成员实体（配置型资源语义）。
- 不批量管理多个成员（一实例管一成员）。
- 不暴露 Members 列表层 schema（不嵌套 xxx_set/xxx_list）。
- 不做异步轮询（写操作同步返回，Read 立即可读）。
- 不在 provider 层校验 sharer 前置约束（由云 API 返回 InvalidParameter，直接透传错误）。

## Decisions

### 决策 1：配置型资源语义——Create/Delete 为 empty 实现

Create 调用时不发任何创建请求，直接 `d.SetId(member_id)` 后转调 Update（让 Update 把目标角色同步到云上）。Delete 直接返回 nil（角色只能置 none 不能"删除"成员）。

**参考**：`tencentcloud/services/cfw/resource_tc_cfw_edge_firewall_switch.go` 的 Create（`d.SetId` 后转 Update）与 Delete（直接 return nil）即是 config 型资源的标准模板。

### 决策 2：资源 ID 与 import

- 资源 ID = `member_id`（目标成员 MemberId）。
- import 用 `app_id:member_id` 两段式（Read 需 app_id 定位），用 `tccommon.FILED_SP` 作分隔符。但 schema 中 `app_id` 标 ForceNew，import 后从复合 ID 解析回 app_id/member_id。
- **重要**：因为 import 需要 app_id 定位，资源 ID 采用 `app_id + FILED_SP + member_id` 复合形式，Read/Delete/Update 中从 `d.Id()` 解析这两个字段（参考 igtm_strategy 的 `idSplit := strings.Split(d.Id(), tccommon.FILED_SP)`）。

### 决策 3：Read 用 AppId 过滤定位成员

- 不可传 `Filters[{Name:"MemberId"}]`（接口过滤名白名单无 MemberId，传错会报 "filter name does not meet the requirements"）。
- 正确做法：`Filters=[{Name:"AppId", Values:[appId], OperatorType:1}]`（1=精确匹配 EQUAL），从返回 `Members` 列表按 MemberId 匹配目标成员。
- 回读字段：`cfw_share_role`（来自 `OrganMemberItem.CfwShareRole`）、`cfw_sharer_app_id`（来自 `OrganMemberItem.CfwSharerAppId`，computed，user 角色时有值）。
- 成员不存在（Members 为空或无匹配）时：先 `log.Printf("[CRUD] fwm_cfw_share_role id=%s", d.Id())`，再 `d.SetId("")`。

### 决策 4：Update 全量设置 + 中转逻辑

- `ModifyCfwShareRole` 请求仅两个字段：`CfwShareRole`（sharer/user/none）+ `UserMemberId`（成员 MemberId）。**不传 SharerMemberId**（服务端设为 user 时会自动定位集团内唯一共享者）。
- 注意 JSON 字段名是 `UserMemberId`（不是 MemberId）。
- 中转逻辑：先 Read 获取 current 角色值 current；若 current == target 则跳过写操作；若 current≠target 且 {current,target} 都属于 {sharer,user}（即 sharer↔user 直转），则先调一次 `CfwShareRole=none`，再调一次目标角色；其余情况（含与 none 互转）直接调一次目标角色。中转过程中任一步失败则整体失败，保留原 state。
- 写操作用 `resource.Retry(tccommon.WriteRetryTimeout, ...)`，失败用 `tccommon.RetryError(e)` 包装。

### 决策 5：schema 字段设计

| 字段 | 类型 | 属性 | 说明 |
|------|------|------|------|
| `member_id` | string | required | 目标成员 MemberId |
| `app_id` | string | required, ForceNew | 目标成员 AppId，用于 Read 定位 |
| `cfw_share_role` | string | required | 共享角色 sharer/user/none |
| `cfw_sharer_app_id` | string | computed | 共享者 AppId，仅 user 角色有值 |

资源描述（中文）："提供云防火墙（fwm）共享角色配置资源，用于管理集团成员的云防火墙共享角色（sharer/user/none）期望状态。"
字段描述用英文。`cfw_share_role` 等 bool/int 可选字段用 `d.GetOkExists`（本资源均为 string，用 `d.GetOk`）。

### 决策 6：fwm 客户端注册

在 `tencentcloud/connectivity/client.go` 中：
- import：`fwm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"`
- 结构体字段：`fwmConn *fwm.Client`
- 方法 `UseFwmClient()`：仿 `UseCfwClient()`，用 `NewClientProfile(300)`。

资源代码中调用：`meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseFwmClient().DescribeOrganMembersWithContext(ctx, request)`。

## Risks / Trade-offs

- **[风险] sharer↔user 中转的两步写非原子** → 中转第一步成功、第二步失败时，资源处于临时 `none` 状态，与用户声明的目标态不一致。Terraform 会因返回 error 触发下次 plan/apply 重试。Mitigation：中转失败返回 error 保留原 state；下次 apply 再次从 `none` 出发继续收敛，最终达成目标态。
- **[风险] Read 依赖 AppId 过滤而非直接按 MemberId 查** → 若同一 AppId 下存在多成员需逐个匹配。Mitigation：从 Members 列表按 MemberId 精确匹配，匹配不到判为不存在并 SetId("")。
- **[风险] 测试账号须为集团管理员** → 非管理员调用返回 UnauthorizedOperation。Mitigation：单元测试用 gomonkey mock 云 API，不依赖真实账号；测试用例覆盖权限错误路径。
- **[取舍] 不在 provider 层前置校验 sharer 的四项约束** → 由云 API 返回 InvalidParameter 透传。取舍理由：约束含运行态信息（如集团内是否已有 sharer），provider 无法完整判断，强行前置校验易与云上状态不一致。

## 云 API 关键 struct（vendor 摘录）

以下摘录自 `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611/models.go`，供实现阶段直接参考，无需回读大文件。

### CommonFilter（DescribeOrganMembers 入参 Filters 元素）
```go
// models.go:342
type CommonFilter struct {
    Name *string `json:"Name,omitnil,omitempty" name:"Name"`           // 筛选字段名
    Values []*string `json:"Values,omitnil,omitempty" name:"Values"`   // 筛选值列表
    OperatorType *int64 `json:"OperatorType,omitnil,omitempty" name:"OperatorType"` // 1=精确匹配,7=in,9=模糊匹配
}
```

### DescribeOrganMembersRequest（models.go:1895）
```go
type DescribeOrganMembersRequest struct {
    *tchttp.BaseRequest
    Filters []*CommonFilter `json:"Filters,omitnil,omitempty" name:"Filters"` // 搜索过滤条件列表
    Limit *int64 `json:"Limit,omitnil,omitempty" name:"Limit"`               // 分页大小，默认 20
    Offset *int64 `json:"Offset,omitnil,omitempty" name:"Offset"`            // 分页偏移量，默认 0
    By *string `json:"By,omitnil,omitempty" name:"By"`                       // 排序字段
    Order *string `json:"Order,omitnil,omitempty" name:"Order"`              // 排序方式 asc/desc
}
```

### DescribeOrganMembersResponse（models.go:1949）
```go
type DescribeOrganMembersResponse struct {
    *tchttp.BaseResponse
    Response *DescribeOrganMembersResponseParams `json:"Response"`
}
// models.go:1938
type DescribeOrganMembersResponseParams struct {
    TotalCount *int64 `json:"TotalCount,omitnil,omitempty" name:"TotalCount"`
    Members []*OrganMemberItem `json:"Members,omitnil,omitempty" name:"Members"` // 集团成员列表
    RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}
```

### OrganMemberItem（models.go:4370，Members 元素，回读字段来源）
```go
type OrganMemberItem struct {
    MemberId *string `json:"MemberId,omitnil,omitempty" name:"MemberId"`           // 成员 ID（资源标识）
    AppId *string `json:"AppId,omitnil,omitempty" name:"AppId"`                     // 成员账号 AppId
    Uin *string `json:"Uin,omitnil,omitempty" name:"Uin"`
    Nickname *string `json:"Nickname,omitnil,omitempty" name:"Nickname"`
    SubAccountCount *int64 `json:"SubAccountCount,omitnil,omitempty"`
    NodeName *string `json:"NodeName,omitnil,omitempty"`
    Role *string `json:"Role,omitnil,omitempty"`    // admin/delegatedAdmin/member
    RoleDisplay *string `json:"RoleDisplay,omitnil,omitempty"`
    AccountGroup *AccountGroupInfo `json:"AccountGroup,omitnil,omitempty"`
    CfwManaged *int64 `json:"CfwManaged,omitnil,omitempty"`                         // 0未纳管 1已纳管
    CfwShareRole *string `json:"CfwShareRole,omitnil,omitempty"`                   // sharer/user/none ★核心回读
    CfwShareRoleDisplay *string `json:"CfwShareRoleDisplay,omitnil,omitempty"`
    CfwSharerAppId *string `json:"CfwSharerAppId,omitnil,omitempty"`                // 共享者AppId，user时有值 ★computed
    CfwInstanceId *string `json:"CfwInstanceId,omitnil,omitempty"`
    PolicyAnalysisEnabled *int64 `json:"PolicyAnalysisEnabled,omitnil,omitempty"`
    MemberCreateTime *string `json:"MemberCreateTime,omitnil,omitempty"`
    JoinType *string `json:"JoinType,omitnil,omitempty"`
    CfwPayStatus *int64 `json:"CfwPayStatus,omitnil,omitempty"`
    CfwCapable *int64 `json:"CfwCapable,omitnil,omitempty"`
    SgManaged *int64 `json:"SgManaged,omitnil,omitempty"`
    IsCfwPostPay *int64 `json:"IsCfwPostPay,omitnil,omitempty"`
}
```

### AccountGroupInfo（models.go:70）
```go
type AccountGroupInfo struct {
    GroupId *string `json:"GroupId,omitnil,omitempty"`
    GroupName *string `json:"GroupName,omitnil,omitempty"`
}
```

### ModifyCfwShareRoleRequest（models.go:3427）
```go
type ModifyCfwShareRoleRequest struct {
    *tchttp.BaseRequest
    CfwShareRole *string `json:"CfwShareRole,omitnil,omitempty" name:"CfwShareRole"`         // sharer/user/none
    UserMemberId *string `json:"UserMemberId,omitnil,omitempty" name:"UserMemberId"`         // 使用者Id ★注意字段名是 UserMemberId
    SharerMemberId *string `json:"SharerMemberId,omitnil,omitempty" name:"SharerMemberId"`   // 共享者账号ID，user时必填（本资源不传）
}
```

### ModifyCfwShareRoleResponse（models.go:3467，仅 RequestId）
```go
type ModifyCfwShareRoleResponse struct {
    *tchttp.BaseResponse
    Response *ModifyCfwShareRoleResponseParams `json:"Response"`
}
```

### client.go 方法签名
```go
// DescribeOrganMembersWithContext：仅管理员可调用，错误码 UNAUTHORIZEDOPERATION
func (c *Client) DescribeOrganMembersWithContext(ctx context.Context, request *DescribeOrganMembersRequest) (response *DescribeOrganMembersResponse, err error)

// ModifyCfwShareRoleWithContext
func (c *Client) ModifyCfwShareRoleWithContext(ctx context.Context, request *ModifyCfwShareRoleRequest) (response *ModifyCfwShareRoleResponse, err error)
```

## 代码风格参照（config 型资源模板）

参照 `tencentcloud/services/cfw/resource_tc_cfw_edge_firewall_switch.go`（同为配置型资源）：

- **Schema**：`map[string]*schema.Schema`，required 字段直接声明。
- **Create**：`defer LogElapsed/InconsistentCheck` → 取 key → `d.SetId(...)` → `return resourceTencentCloud...Update(d, meta)`。
- **Read**：`defer ...` → 取 `logId`/`ctx`/`service` → 调 service.Describe...ById → 空则 `log.Printf("[WARN]...")` + `d.SetId("")` → 非空逐字段 `if x.X != nil { _ = d.Set(...) }`。
- **Update**：`defer ...` → 构造 request → `resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {...})` → 失败 `log.Printf("[CRITAL]...")` 并返回。
- **Delete**：`defer ...` → 直接 `return nil`。

参照 `tencentcloud/services/igtm/resource_tc_igtm_strategy.go`：
- 复合 ID 解析：`idSplit := strings.Split(d.Id(), tccommon.FILED_SP)`，校验 `len(idSplit)`。
- Importer：`Importer: &schema.ResourceImporter{State: schema.ImportStatePassthrough}`。
- WriteRetry 内调用 `WithContext(ctx, request)` 形式，失败用 `tccommon.RetryError(e)` 包装，成功打印 `log.Printf("[DEBUG]%s api[%s] success...")`。

## CRUD 流程要点

### Create
1. 取 `member_id`、`app_id`。
2. `d.SetId(app_id + FILED_SP + member_id)`。
3. `return resourceTencentCloudFwmCfwShareRoleConfigUpdate(d, meta)`。

### Read
1. 解析 `d.Id()` → app_id, member_id。
2. 构造 `DescribeOrganMembersRequest`：`Filters=[{Name:"AppId", Values:[app_id], OperatorType:1}]`，`Limit` 取接口注释最大值（用一个大值确保一次取全，如 100）。
3. retry（`tccommon.ReadRetryTimeout`）调用 `DescribeOrganMembersWithContext`。
4. 响应空校验：`response==nil || response.Response==nil || len(Members)==0` → `log.Printf("[CRUD] fwm_cfw_share_role id=%s", d.Id())` 后 `d.SetId("")` 返回 nil。
5. 遍历 Members 按 MemberId 匹配；匹配不到同上处理。
6. 匹配到则回读 `cfw_share_role`（非 nil 才 set）、`cfw_sharer_app_id`、`member_id`、`app_id`。

### Update（核心）
1. 解析 `d.Id()` → app_id, member_id。
2. 先 Read 获取 current 角色值。
3. target = `d.Get("cfw_share_role")`。
4. 若 current == target：跳过写操作。
5. 若 current≠target 且 current∈{sharer,user} 且 target∈{sharer,user}（sharer↔user 直转）：
   - 先调 `ModifyCfwShareRole`（CfwShareRole=none, UserMemberId=member_id）。
   - 再调 `ModifyCfwShareRole`（CfwShareRole=target, UserMemberId=member_id）。
6. 否则：直接调一次 `ModifyCfwShareRole`（CfwShareRole=target, UserMemberId=member_id）。
7. 每次写调用用 `resource.Retry(tccommon.WriteRetryTimeout, ...)`，失败 `tccommon.RetryError(e)` 包装。
8. 写完成后 `return resourceTencentCloudFwmCfwShareRoleConfigRead(d, meta)`。

### Delete
1. `defer ...` → `return nil`（empty 实现）。

## Open Questions

无。所有关键行为已由云 API 与用户自定义要求明确。