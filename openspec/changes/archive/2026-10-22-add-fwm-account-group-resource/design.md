## Context

为云产品 FWM（Edge 边界防火墙账号组管理，SDK 包 `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611`）新增 Terraform 通用资源 `tencentcloud_fwm_account_group`，用于管理「账号组」整个生命周期。

业务语义要点（接口名无法体现，必须遵守）：
1. 账号组是 FWM 把集团成员账号归类、作为策略下发接收方集合的实体。资源标识 = 服务端生成的 GroupId（形如 `acg-xxxxxxxx`，8 位小写随机串），**非复合 ID**。
2. 四个接口全部同步：Create/Modify/Delete 成功即已落库，Read 立即可读，**不要生成 waitReady/async 轮询**。
3. 仅集团管理员可用：Create/Describe/Modify/Delete 全部先经 IsMemberAdmin 校验，非管理员返回 UnauthorizedOperation（业务约束非缺陷）。
4. 无单实例查询接口，Read 用列表 DescribeAccountGroups + Filters 按 GroupId 精确过滤（OperatorType=1 表示 EQUAL），取 `AccountGroups` 首元素；查不到返回空数组视为已删除（`d.SetId("")`，不报错不重试）。
5. 删除为软删（服务端置 status=1，查询侧只查 status=0，provider 无感知）。
6. DeleteAccountGroup 有两种业务拒绝（非幂等）：组内仍有成员（member_count>0）报 InvalidParameter；该组被已生效策略引用报 InvalidParameter。destroy 前须清空成员并解除策略引用，provider 原样透传错误提示用户先清理。

修复现状：FWM 服务目录在 `tencentcloud/services/fwm/` 尚不存在，需新建；`tencentcloud/connectivity/client.go` 中尚无 FWM 客户端，需新增 `UseFwmV20250611Client()`。

## Goals / Non-Goals

**Goals:**
- 新增 `tencentcloud_fwm_account_group` 资源，CRUD 完整、支持 import。
- 接入 FWM SDK（v20250611）到 connectivity 层。
- description 字段标 ForceNew（接口缺口：ModifyAccountGroup 不接受 description）。
- 准确处理「读不到即删除」语义。

**Non-Goals:**
- 不实现账号组成员管理（成员关系走独立接口 ModifyMemberGroup，另建 `tencentcloud_fwm_account_group_member` 资源）。
- 不内联 members 集合（无批量接口、逐个调用会限流、成员增减改 MemberCount 造成漂移）。
- 不处理异步任务（接口本身同步）。

## Decisions

### 决策 1：description 标 ForceNew 而非 mutable

`ModifyAccountGroupRequest` 仅含 `GroupId` 与 `GroupName` 两个字段，不接受 Description；CreateAccountGroup 接受 Description。若用 mutable 实现 description 变更，改了不生效产生永久 diff。因此 description 标 `ForceNew`：用户改描述时整组重建。

**备选**：透传用户任意输入让服务端忽略 —— 拒绝，会永久 diff。

### 决策 2：Read 用列表 + Filters 精确取首元素

服务端无单实例查询接口。DescribeAccountGroups 的字段映射表明确包含 GroupId→group_id，可用 `Filters=[{Name:"GroupId",Values:[groupId],OperatorType:1}]`（1=EQUAL）精确取单条。AccountGroups 为空 → `d.SetId("")` 移除 state 不报错。

### 决策 3：资源标识为单 GroupId，import 直接写 GroupId

GroupId 为服务端生成的 8 位随机串，全局唯一，无需复合分隔符。import 使用 ImportStatePassthrough，直接传 GroupId。Read/Delete/Update 均从 `d.Id()` 取 GroupId 作为请求参数。

### 决策 4：客户端命名 `UseFwmV20250611Client`

遵循 igtm 的 `UseIgtmV20231024Client` 命名风格，新增连接字段 `fwmv20250611Conn *fwmv20250611.Client` 及方法 `UseFwmV20250611Client()`。

### 决策 5：新增服务层 `FwmService`

参考 `CfwService`，新增 `FwmService`，实现 `DescribeFwmAccountGroupById(ctx, groupId)`（调用 DescribeAccountGroups + Filters 取首元素）与 `DeleteFwmAccountGroupById(ctx, groupId)`（调用 DeleteAccountGroup）。CRUD 函数中 Read/Delete 经过服务层，Create/Update 直接在资源文件内调用 SDK（参考 igtm_strategy 风格）。

### 决策 6：测试使用 gomonkey mock，不使用 TF 测试套件

新增资源，按规范使用 gomonkey 对 FWM Client 方法做 mock，做业务代码单元测试，参考 `resource_tc_cfw_cluster_fw_bypass_config_test.go` 的 mockMeta 模式。

## Risks / Trade-offs

- [DeleteAccountGroup 非幂等业务拒绝] → provider 原样透传 InvalidParameter，提示用户先清空成员并解除策略引用，不当成功。
- [Description 响应可能 null/空串] → Read 中做空值保护，nil 不调用 set。
- [Order 枚举无白名单，服务端直接拼 SQL] → provider 端自定 ValidateFunc 限定小写 asc/desc，不透传任意输入（本资源 schema 不暴露 Order/By，仅 Read 内部固定，故该风险仅存在于潜在的数据源场景，本资源无）。

## 参考信息（vendor 摘录，供实施阶段直接使用，避免回读大文件）

### 云 API Struct 关键字段定义（源自 vendor/.../fwm/v20250611/models.go）

**AccountGroup（列表元素、Read 回填对象）— models.go:49**
```go
type AccountGroup struct {
    GroupId      *string `json:"GroupId,omitnil,omitempty" name:"GroupId"`        // 账户组ID 如 acg-1a2b3c4d
    GroupName    *string `json:"GroupName,omitnil,omitempty" name:"GroupName"`    // 账户组名称，<=50字符
    Description  *string `json:"Description,omitnil,omitempty" name:"Description"`// 描述，注意：可能返回 null
    MemberCount  *int64  `json:"MemberCount,omitnil,omitempty" name:"MemberCount"`// 成员数量（只读）
    CreateTime   *string `json:"CreateTime,omitnil,omitempty" name:"CreateTime"`  // 创建时间（只读）
    CanDelete    *bool   `json:"CanDelete,omitnil,omitempty" name:"CanDelete"`    // 是否可删除（只读计算值）
}
```

**CommonFilter（过滤器）— models.go:342**
```go
type CommonFilter struct {
    Name         *string  `json:"Name,omitnil,omitempty" name:"Name"`             // 筛选字段名 如 GroupId
    Values       []*string `json:"Values,omitnil,omitempty" name:"Values"`         // 筛选值列表
    OperatorType *int64   `json:"OperatorType,omitnil,omitempty" name:"OperatorType"` // 1=精确匹配,7=in,9=模糊匹配
}
```

**CreateAccountGroupRequest — models.go:365**
```go
type CreateAccountGroupRequest struct {
    *tchttp.BaseRequest
    GroupName   *string `json:"GroupName,omitnil,omitempty" name:"GroupName"`     // 必填，账户组名称
    Description *string `json:"Description,omitnil,omitempty" name:"Description"` // 可选，描述
}
// CreateAccountGroupResponse.Response.GroupId *string — models.go:398（新建账户组ID）
```

**DescribeAccountGroupsRequest — models.go:1597**
```go
type DescribeAccountGroupsRequest struct {
    *tchttp.BaseRequest
    Filters []*CommonFilter `json:"Filters,omitnil,omitempty" name:"Filters"`  // 过滤器
    Limit   *int64          `json:"Limit,omitnil,omitempty" name:"Limit"`      // 分页大小，默认20
    Offset  *int64          `json:"Offset,omitnil,omitempty" name:"Offset"`    // 偏移，默认0
    Order   *string         `json:"Order,omitnil,omitempty" name:"Order"`      // 顺序 asc/desc
    By      *string         `json:"By,omitnil,omitempty" name:"By"`           // 排序字段
}
// DescribeAccountGroupsResponse.Response.AccountGroups []*AccountGroup — models.go:1645
// DescribeAccountGroupsResponse.Response.TotalCount *int64 — models.go:1642
```

**ModifyAccountGroupRequest — models.go:3302（注意：仅 GroupId+GroupName，无 Description）**
```go
type ModifyAccountGroupRequest struct {
    *tchttp.BaseRequest
    GroupId   *string `json:"GroupId,omitnil,omitempty" name:"GroupId"`   // 必填，账户组ID
    GroupName *string `json:"GroupName,omitnil,omitempty" name:"GroupName"` // 必填，新名称
}
```

**DeleteAccountGroupRequest — models.go:1179**
```go
type DeleteAccountGroupRequest struct {
    *tchttp.BaseRequest
    GroupId *string `json:"GroupId,omitnil,omitempty" name:"GroupId"` // 必填，账户组ID
}
```

### Client 方法签名（vendor/.../fwm/v20250611/client.go）

```go
func NewClient(credential common.CredentialIface, region string, clientProfile *profile.ClientProfile) (client *Client, err error)  // client.go:39
func (c *Client) CreateAccountGroupWithContext(ctx, *CreateAccountGroupRequest) (*CreateAccountGroupResponse, error)      // client.go:143
func (c *Client) DescribeAccountGroupsWithContext(ctx, *DescribeAccountGroupsRequest) (*DescribeAccountGroupsResponse, error)  // client.go:1175
func (c *Client) ModifyAccountGroupWithContext(ctx, *ModifyAccountGroupRequest) (*ModifyAccountGroupResponse, error)     // client.go:2219
func (c *Client) DeleteAccountGroupWithContext(ctx, *DeleteAccountGroupRequest) (*DeleteAccountGroupResponse, error)  // client.go:787
```

### 代码风格参照：resource_tc_igtm_strategy.go（关键结构）

- 包声明、import 分组：SDK → tccommon → helper。包别名 `fwmv20250611`。
- `ResourceTencentCloudIgtmStrategy() *schema.Resource`：Create/Read/Update/Delete + Importer（ImportStatePassthrough）。
- Schema：required/optional/computed 分组；search-by-id 字段设 `ForceNew`（本资源 group_id 为 computed 不用；description 标 ForceNew）。
- Create：`defer tccommon.LogElapsed(...)` + `defer tccommon.InconsistentCheck(d, meta)`；`resource.Retry(tccommon.WriteRetryTimeout, ...)` 包裹 SDK 调用；内部校验 `result==nil||result.Response==nil` 返回 NonRetryableError；成功后检查 id 字段非空，`d.SetId(id)`，再调 Read。
- Update：用 `mutableArgs` + `d.HasChange` 判断是否需要变更；本次 group_name 可变（mutableArgs=["group_name"]），description 变更因 ForceNew 直接由框架重建不走 Update。
- Delete：`resource.Retry(WriteRetryTimeout)` 调用 DeleteAccountGroup，错误透传。
- helper：`helper.String(s)` → `*string`；Description 用 `d.GetOk("description")` 判断是否显式设置。

### 服务层参照：cfw/service_tencentcloud_cfw.go（关键结构）

```go
type CfwService struct { client *connectivity.TencentCloudClient }
// FwmService 同构：
//   DescribeFwmAccountGroupById(ctx, groupId) (*fwmv20250611.AccountGroup, error)
//     构建 Filters=[{Name:"GroupId",Values:[groupId],OperatorType:1}]，Limit=20，调 DescribeAccountGroups，
//     AccountGroups 为空返回 nil,nil（视为删除）；首元素返回。
//   DeleteFwmAccountGroupById(ctx, groupId) error
//     调用 DeleteAccountGroup。
```

### connectivity/client.go 变更点

- import 区新增：`fwmv20250611 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"`（按字母序，es 之后、ga2 之前的 fwm 段）。
- struct `TencentCloudClient` 新增字段：`fwmv20250611Conn *fwmv20250611.Client`。
- 新增方法（参照 UseIgtmV20231024Client）：
```go
func (me *TencentCloudClient) UseFwmV20250611Client() *fwmv20250611.Client {
    if me.fwmv20250611Conn != nil {
        return me.fwmv20250611Conn
    }
    cpf := me.NewClientProfile(300)
    me.fwmv20250611Conn, _ = fwmv20250611.NewClient(me.Credential, me.Region, cpf)
    me.fwmv20250611Conn.WithHttpTransport(&LogRoundTripper{})
    return me.fwmv20250611Conn
}
```

### provider.go 变更点

- import 区新增：`"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/fwm"`（按字母序，ciam/ckafka...clb 段之后或 config 之前，ef 后置）。
- ResourcesMap 新增：`"tencentcloud_fwm_account_group": fwm.ResourceTencentCloudFwmAccountGroup(),`

### Schema 定义

| SchemaName | Type | 属性 | 云API路径 / 说明 |
|---|---|---|---|
| `group_name` | TypeString | Required | CreateAccountGroup.GroupName / ModifyAccountGroup.GroupName，<=50字符，同adminAppId下唯一 |
| `description` | TypeString | Optional, ForceNew | CreateAccountGroup.Description，**仅建时可设** |
| `group_id` | TypeString | Computed (ForceNew隐含) | CreateAccountGroup.Response.GroupId / AccountGroup.GroupId |
| `member_count` | TypeInt | Computed | AccountGroup.MemberCount（只读） |
| `can_delete` | TypeBool | Computed | AccountGroup.CanDelete（只读计算值，不可当删除判据） |
| `create_time` | TypeString | Computed | AccountGroup.CreateTime（只读） |

- testValue：group_name=test-account-group、description=test account group。
- description 中文描述、字段 description 英文描述。

### CRUD 行为摘要

- **Create**：request.GroupName/Description 取值；Retry 调 CreateAccountGroupWithContext；校验 Response.GroupId 非空；`d.SetId(groupId)`；调 Read。
- **Read**：`d.Id()` 为 GroupId；服务层 DescribeFwmAccountGroupById；nil 视为删除 → `log.Printf("[CRUD] fwm_account_group id=%s", d.Id())` 后 `d.SetId("")`；回填各字段前判 nil。
- **Update**：mutableArgs=["group_name"]；HasChange(group_name) 时构建 ModifyAccountGroupRequest（GroupId+GroupName），Retry 调 ModifyAccountGroupWithContext；再 Read。description 变更因 ForceNew 被框架重建，不走 Update。
- **Delete**：`d.Id()` 为 GroupId；服务层 DeleteFwmAccountGroupById；DeleteAccountGroup 的 InvalidParameter 业务拒绝原样透传。
- **Import**：直接写 GroupId（如 `terraform import tencentcloud_fwm_account_group.example acg-1a2b3c4d`）。