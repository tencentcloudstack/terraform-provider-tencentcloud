## Context

云防火墙纳管（fwm）产品首次接入 Terraform Provider。本资源为 **配置型资源（RESOURCE_KIND_CONFIG）**，管理「某个集团成员账号是否被云防火墙纳管」这一个开关的期望状态。成员实体本就存在于集团中，资源不创建/不销毁实体，只有 Read（查状态）+ Update（设开关）。

关键约束（已在业务语义中证实）：
- `DescribeOrganMembers` 的过滤字段映射表**不含 MemberId**，传 `Filters[{Name:"MemberId"}]` 会报错「filter name does not meet the requirements」。只能用 `AppId` 过滤定位成员列表，再从返回 Members 中按 MemberId 匹配。
- `ModifyCfwManaged` 是**同步接口**（无 FlowId/TaskId），立即返回，Read 可立即读到新值，**不要生成 waitReady/轮询**。
- 关闭纳管（cfw_managed 1→0）可能被 `ValidateCfwManagedEdit` 以 InvalidParameter 拒绝（共享套餐/已下发规则/非集团成员），provider 原样透传错误。
- 仅集团管理员可操作，非管理员返回 UnauthorizedOperation。
- fwm SDK 包 `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611` 已在 go.mod / vendor / tools.go 中就绪，但 connectivity 层尚未接入 `UseFwmClient()`，需新增。

参考模板：`tencentcloud/services/cfw/resource_tc_cfw_edge_firewall_switch.go`（同为配置型资源，Create→Update、Delete 空实现）。

## Goals / Non-Goals

**Goals:**
- 新增 `tencentcloud_fwm_cfw_managed` 资源，管理单个集团成员的纳管开关。
- Schema 暴露 `member_id`（required）、`app_id`（required + forceNew）、`cfw_managed`（required，int，1/0）。
- Read 用 AppId 过滤定位成员，再按 MemberId 匹配读 cfw_managed；列表为空或匹配不到 → `d.SetId("")` 移除 state。
- Update 调 ModifyCfwManaged 设开关，且当 cfw_managed 与当前值相同时跳过写操作。
- Create 为 empty 实现（仅 `d.SetId(member_id)` 后复用 Update）；Delete 为 empty 实现（从 state 移除，不调云 API）。
- 支持 import，使用 `app_id` + `member_id` 两段式复合 ID（分隔符用 `tccommon.FILED_SP` 即 `#`，与现有复合 ID 约定一致）。
- connectivity 层新增 `UseFwmClient()` 与 `fwmConn` 字段。
- provider.go / provider.md 注册资源。
- 提供 mock（gomonkey）单元测试。

**Non-Goals:**
- 不批量管理多个成员（一个资源实例只管一个成员）。
- 不暴露 Members 返回的宿主属性（nickname/uin/role 等）作为可比字段，避免 diff 漂移。
- 不实现关闭纳管被拒时的自动清理逻辑（透传错误，提示用户先清理规则）。
- 不处理 account_group 联动副作用（provider 无感知）。
- 不生成 _extension.go 文件。

## Decisions

### 决策 1：配置型资源 CRUD 骨架（参考 cfw_edge_firewall_switch）
- **Create**：empty 空实现 → `d.SetId(member_id)` → 复用 `Update(d, meta)` 调用 ModifyCfwManaged 设开关。理由：成员本就存在，Create 即「声明期望开关状态」，复用 Update 写入逻辑避免重复。
- **Read**：service 层 `DescribeFwmCfwManagedByAppId(ctx, appId, memberId)` 用 `Filters=[{Name:"AppId",Values:[appId],OperatorType:1}]` 过滤，从 Members 中按 MemberId 匹配，返回该成员对象；资源 Read 据此 set 字段。
- **Update**：比较 `d.Get("cfw_managed")` 与 Read 到的当前值，相同则跳过，不同则调 ModifyCfwManaged。
- **Delete**：empty 空实现，仅 `return nil`（从 state 移除即可，不调云 API）。
- **备选**：Delete 调 ModifyCfwManaged 设 cfw_managed=0。**否决**——需求明确「不能删除纳管，只能关开关」，且 destroy 语义不应隐式改云端配置。

### 决策 2：复合 ID 格式
- ID = `app_id + "#" + member_id`（使用 `tccommon.FILED_SP` 即 `#`，与现有复合 ID 约定一致）。
- Read/Update 从 `d.Id()` split 出 app_id 与 member_id。
- Import 支持 `terraform import tencentcloud_fwm_cfw_managed.xxx "appId#memberId"`。
- **备选**：用冒号 `:`。**否决**——现有代码库统一用 `tccommon.FILED_SP`(#)，保持一致。

### 决策 3：Read 定位策略
- 只能用 AppId 过滤（MemberId 不在映射表内），所以 schema 必须 required + forceNew 暴露 app_id。
- Filters 用 `OperatorType=1`（EQUAL，AppId→ou.appid 在映射内）。
- Limit 取云 API 注释标注的默认上限场景的合理值（默认 20，但为保证能匹配到，传较大值如 100）。无上限，传负报错。
- 列表为空或匹配不到 MemberId → `d.SetId("")`，不报错。

### 决策 4：connectivity 层接入
- 在 `tencentcloud/connectivity/client.go` 中：
  - import `fwm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"`
  - 结构体加 `fwmConn *fwm.Client` 字段
  - 新增 `UseFwmClient() *fwm.Client` 方法（参考 UseCfwClient 模式）

### 决策 5：响应字段裁剪
- 只关心 cfw_managed，其余宿主属性（uin/nickname/role/account_group 等）一律不暴露进 schema，避免随集团变化漂移导致 diff。
- schema 仅 3 个字段：member_id、app_id、cfw_managed。

## Risks / Trade-offs

- **[Read 用 AppId 过滤可能返回多条]** → Mitigation：从返回 Members 列表里再用 MemberId 精确匹配，匹配不到视为成员已离开集团 → `d.SetId("")`。
- **[关闭纳管可能被服务端拒绝]** → Mitigation：透传 InvalidParameter 错误，在错误信息中提示先清理已下发规则；不静默改 state。
- **[account_group 联动副作用]** → Mitigation：provider 无感知，若同时管理 account_group 资源会产生 diff 来源，属已知联动，文档提示。
- **[非集团管理员调用失败]** → Mitigation：透传 UnauthorizedOperation，测试账号须为集团管理员。
- **[app_id 变更需重建]** → Mitigation：app_id 标 ForceNew，变更触发重建（Create empty + Update）。

## 参考信息（实施阶段直接使用，无需回读 vendor 大文件）

### 云 API request/response struct 关键字段（vendor 摘录）

#### DescribeOrganMembersRequest（fwm/v20250611/models.go:1878）
```go
type DescribeOrganMembersRequestParams struct {
    Filters []*CommonFilter `json:"Filters,omitnil,omitempty" name:"Filters"`  // 过滤条件列表
    Limit   *int64          `json:"Limit,omitnil,omitempty" name:"Limit"`       // 分页大小，默认 20
    Offset  *int64          `json:"Offset,omitnil,omitempty" name:"Offset"`     // 分页偏移量，默认 0
    By      *string         `json:"By,omitnil,omitempty" name:"By"`            // 排序字段，如 MemberCreateTime
    Order   *string         `json:"Order,omitnil,omitempty" name:"Order"`       // 排序方式：asc/desc
}
// Request 构造：fwm.NewDescribeOrganMembersRequest()
// 方法：client.DescribeOrganMembers(request) / DescribeOrganMembersWithContext(ctx, request)
```

#### CommonFilter（fwm/v20250611/models.go:342）
```go
type CommonFilter struct {
    Name         *string  `json:"Name,omitnil,omitempty" name:"Name"`           // 筛选字段名（AppId 在映射内，MemberId 不在）
    Values       []*string `json:"Values,omitnil,omitempty" name:"Values"`      // 筛选值列表
    OperatorType *int64   `json:"OperatorType,omitnil,omitempty" name:"OperatorType"` // 1=精确匹配(EQUAL)，7=in，9=模糊匹配
}
```

#### DescribeOrganMembersResponse（fwm/v20250611/models.go:1938）
```go
type DescribeOrganMembersResponseParams struct {
    TotalCount *int64              `json:"TotalCount,omitnil,omitempty" name:"TotalCount"` // 集团成员总数
    Members    []*OrganMemberItem  `json:"Members,omitnil,omitempty" name:"Members"`       // 集团成员列表
    RequestId  *string             `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}
// Response 构造：fwm.NewDescribeOrganMembersResponse()
```

#### OrganMemberItem（fwm/v20250611/models.go:4370）— 关键字段
```go
type OrganMemberItem struct {
    MemberId      *string `json:"MemberId,omitnil,omitempty" name:"MemberId"`        // 成员 ID（JSON tag 是 MemberId）
    AppId         *string `json:"AppId,omitnil,omitempty" name:"AppId"`              // 成员账号 AppId
    Uin           *string `json:"Uin,omitnil,omitempty" name:"Uin"`
    Nickname      *string `json:"Nickname,omitnil,omitempty" name:"Nickname"`
    SubAccountCount *int64 `json:"SubAccountCount,omitnil,omitempty" name:"SubAccountCount"`
    NodeName      *string `json:"NodeName,omitnil,omitempty" name:"NodeName"`
    Role          *string `json:"Role,omitnil,omitempty" name:"Role"`
    RoleDisplay   *string `json:"RoleDisplay,omitnil,omitempty" name:"RoleDisplay"`
    AccountGroup  *AccountGroupInfo `json:"AccountGroup,omitnil,omitempty" name:"AccountGroup"`
    CfwManaged    *int64  `json:"CfwManaged,omitnil,omitempty" name:"CfwManaged"`    // 0-未纳管，1-已纳管（int 不是 bool）
    CfwShareRole  *string `json:"CfwShareRole,omitnil,omitempty" name:"CfwShareRole"`
    // ... 其余宿主字段略，本资源不暴露
}
```

#### AccountGroupInfo（fwm/v20250611/models.go:70）
```go
type AccountGroupInfo struct {
    GroupId   *string `json:"GroupId,omitnil,omitempty" name:"GroupId"`
    GroupName *string `json:"GroupName,omitnil,omitempty" name:"GroupName"`
}
```

#### ModifyCfwManagedRequest（fwm/v20250611/models.go:3355）
```go
type ModifyCfwManagedRequestParams struct {
    MemberId   *string `json:"MemberId,omitnil,omitempty" name:"MemberId"`     // 目标成员 Id（必填）
    CfwManaged *int64  `json:"CfwManaged,omitnil,omitempty" name:"CfwManaged"` // 0-关闭，1-开启（必填，int）
}
// Request 构造：fwm.NewModifyCfwManagedRequest()
// 方法：client.ModifyCfwManaged(request) / ModifyCfwManagedWithContext(ctx, request)
```

#### ModifyCfwManagedResponse（fwm/v20250611/models.go:3394）
```go
type ModifyCfwManagedResponseParams struct {
    RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`  // 同步返回，无 FlowId/TaskId
}
```

### 代码风格参照

#### 参照 1：cfw_edge_firewall_switch（配置型资源 CRUD 骨架，tencentcloud/services/cfw/resource_tc_cfw_edge_firewall_switch.go）
- Create：`d.SetId(id)` → `return resourceTencentCloudXxxUpdate(d, meta)`（复用 Update）
- Read：调 service 层 `DescribeXxxById`，nil → `d.SetId("")` + WARN 日志；字段 set 前判 nil
- Update：`resource.Retry(tccommon.WriteRetryTimeout, ...)` 包裹 API 调用，`tccommon.RetryError(e)` 包装错误
- Delete：`return nil`（empty）
- Schema：Required 字段直接定义，int 字段用 `schema.TypeInt`，bool/int 可选字段用 `d.GetOkExists`
- 包导入模式：
```go
import (
    "context"
    "fmt"
    "log"
    tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    fwm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"
    "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)
```

#### 参照 2：service 层 DescribeXxxById 模式（tencentcloud/services/cfw/service_tencentcloud_cfw.go:952）
```go
func (me *CfwService) DescribeCfwEdgeFirewallSwitchById(ctx context.Context, publicIp string) (...) {
    request := cfw.NewDescribeFwEdgeIpsRequest()
    request.Filters = []*cfw.CommonFilter{{Name: common.StringPtr("PublicIp"), Values: common.StringPtrs([]string{publicIp}), OperatorType: common.Int64Ptr(1)}}
    request.Limit = helper.Int64(100)
    defer func() { if errRet != nil { log.Printf("[CRITAL]%s api[%s] fail ...") } }()
    err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
        ratelimit.Check(request.GetAction())
        result, e := me.client.UseCfwClient().DescribeFwEdgeIps(request)
        if e != nil { return tccommon.RetryError(e) }
        ...
    })
}
```
- 注意：本资源 service 层需用 `me.client.UseFwmClient()` 而非 UseCfwClient。

#### 参照 3：复合 ID 解析（tencentcloud/services/ccn/resource_tc_ccn_route_table_input_policies.go:187）
```go
items := strings.Split(d.Id(), tccommon.FILED_SP)  // FILED_SP = "#"
if len(items) < 2 { return fmt.Errorf("id is broken,%s", d.Id()) }
appId := items[0]
memberId := items[1]
```

#### 参照 4：mock（gomonkey）单测模式（tencentcloud/services/ccn/resource_tc_ccn_instances_accept_attach_test.go）
```go
type mockMeta struct { client *connectivity.TencentCloudClient }
func (m *mockMeta) GetAPIV3Conn() *connectivity.TencentCloudClient { return m.client }
var _ tccommon.ProviderMeta = &mockMeta{}

patches := gomonkey.NewPatches()
defer patches.Reset()
client := &fwm.Client{}
patches.ApplyMethodReturn(mockMeta.client, "UseFwmClient", client)
patches.ApplyMethodFunc(client, "DescribeOrganMembers", func(req *fwm.DescribeOrganMembersRequest) (*fwm.DescribeOrganMembersResponse, error) {
    resp := fwm.NewDescribeOrganMembersResponse()
    resp.Response = &fwm.DescribeOrganMembersResponseParams{ Members: []*fwm.OrganMemberItem{{MemberId: helper.String("mid"), AppId: helper.String("aid"), CfwManaged: helper.IntInt64(1)}}, }
    return resp, nil
})
d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{...})
err := res.Create(d, meta)  // 或 Read/Update
```

### connectivity 层接入要点（tencentcloud/connectivity/client.go）
- import 区（约第 46 行 cfw 导入旁）：加 `fwm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"`
- 结构体字段区（约第 230 行 cfwConn 旁）：加 `fwmConn *fwm.Client`
- 方法区（约第 1771 行 UseCfwClient 旁）：加 `UseFwmClient()` 方法，参考 UseCfwClient 实现

### provider 注册要点
- `tencentcloud/provider.go`：在 ResourcesMap 中注册 `"tencentcloud_fwm_cfw_managed": fwm.ResourceTencentCloudFwmCfwManaged()`（因 fwm 产品首次接入，新增一行即可，按字母序放置）
- `tencentcloud/provider.md`：新增 fwm 产品章节条目
