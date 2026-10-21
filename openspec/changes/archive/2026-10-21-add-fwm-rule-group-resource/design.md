## Context

云防火墙 FWM（Firewall Manager）在 SDK 中是一个独立产品包 `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611`（`APIVersion = "2025-06-11"`），已 vendor 到仓库（`vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611/`，包含 `client.go` / `models.go` / `errors.go`）。provider 目前**没有** `fwm` 服务包，也没有任何 `fwm` 相关资源/数据源，因此本次需要新建服务包 + 新建 connectivity 客户端工厂。

「规则组」（RuleGroup）是防火墙规则的容器。四个接口语义：

| 操作 | 接口 | 同步/异步 | 关键入参 | 关键出参 |
|---|---|---|---|---|
| Create | `CreateRuleGroup` | 同步 | `GroupName`、`Product` | `GroupId` |
| Read | `DescribeRuleGroups` | 同步 | `Filters`（`GroupId` 精确/模糊匹配）、`Product`、`Order`、`By`、`Offset`、`Limit` | `TotalCount`、`RuleGroups []*RuleGroupResp` |
| Update | `ModifyRuleGroup` | 同步 | `GroupId`（必选）、`GroupName` | 无（仅 `RequestId`） |
| Delete | `DeleteRuleGroup` | 同步 | `GroupIds []*string` | 无（仅 `RequestId`） |

全部接口**同步**返回，没有任何 `FlowId` / `TaskId` / `JobId`，因此**不需要**异步轮询、**不需要**声明 `Timeouts` 块。

代码风格参照资源：`tencentcloud/services/igtm/resource_tc_igtm_strategy.go`（CRUD 骨架、`defer tccommon.LogElapsed` / `defer tccommon.InconsistentCheck`、`var (...)` 声明块、`resource.Retry(tccommon.WriteRetryTimeout, ...)`、`d.Set` 前的 nil 判断）。
服务层参照：`tencentcloud/services/igtm/service_tencentcloud_igtm.go` 的 `DescribeIgtmStrategyById`（`resource.Retry(tccommon.ReadRetryTimeout, ...)` + `ratelimit.Check(request.GetAction())`）。
单元测试参照：`tencentcloud/services/bh/resource_tc_bh_bind_device_resource_test.go`（gomonkey + `mockMeta` + `schema.TestResourceDataRaw`）。

## Goals / Non-Goals

**Goals:**

- 交付 RESOURCE_KIND_GENERAL 资源 `tencentcloud_fwm_rule_group`，覆盖 Create / Read / Update / Delete 完整生命周期与 import。
- Schema 只暴露云 API **真实可写**的字段（`group_name`、`product`），其余返回字段平铺为 Computed，不创建 `rule_groups` 列表包装层。
- 处理 `Product` 的「入参字符串 / 出参数字」不对称问题，保证 read 后 state 与 config 一致，无永久 diff。
- 全部 SDK 调用包 `resource.Retry`，并对 `result == nil`、`Response == nil`、`GroupId == nil`、`GroupId == ""` 做保护。
- 新建 `tencentcloud/services/fwm/` 服务包、connectivity 客户端工厂、provider 注册、`.md` 文档、gomonkey 单元测试。

**Non-Goals:**

- 不实现 `tencentcloud_fwm_rule_group` 的 data source（本 change 只交付 resource；`DescribeRuleGroups` 的查询参数 `filters` / `order` / `by` / `offset` / `limit` 不对外暴露）。
- 不实现规则（rule）级别的资源（`CreateSecurityGroupRuleGroup` / `ModifyEdgeAclRule` 等不在本次范围）。
- 不修改任何既有资源/数据源/schema，不修改已有 service 方法。
- 不做 SDK 升级（`fwm` 包已在 vendor 中）。

## Decisions

### D1 — Schema 字段集（顶层，平铺）

| HCL 字段 | SDK 字段（来源接口） | 类型 | Required/Optional/Computed | ForceNew | 说明 |
|---|---|---|---|---|---|
| `group_name` | `CreateRuleGroup.GroupName` / `ModifyRuleGroup.GroupName` / `RuleGroupResp.GroupName` | `TypeString` | Required | No | 规则组名称，最长 50 字符 |
| `product` | `CreateRuleGroup.Product` / `RuleGroupResp.Product` | `TypeString` | Required | **Yes** | 归属产品，`ModifyRuleGroup` 不可变更 |
| `group_id` | `CreateRuleGroupResponse.GroupId` / `RuleGroupResp.GroupId` | `TypeString` | Computed | - | 规则组 ID，同时作为 `d.Id()` |
| `app_id` | `RuleGroupResp.AppId` | `TypeString` | Computed | - | 租户 AppId |
| `rule_count` | `RuleGroupResp.RuleCount` | `TypeInt` | Computed | - | 规则数量 |
| `create_time` | `RuleGroupResp.CreateTime` | `TypeString` | Computed | - | 创建时间 |
| `update_time` | `RuleGroupResp.UpdateTime` | `TypeString` | Computed | - | 更新时间 |
| `create_by` | `RuleGroupResp.CreateBy` | `TypeString` | Computed | - | 创建人 SubUin |
| `create_name` | `RuleGroupResp.CreateName` | `TypeString` | Computed | - | 创建人名称 |
| `update_by` | `RuleGroupResp.UpdateBy` | `TypeString` | Computed | - | 更新人 SubUin |
| `update_name` | `RuleGroupResp.UpdateName` | `TypeString` | Computed | - | 更新人名称 |
| `receive_account` | `RuleGroupResp.ReceiveAccount` | `TypeList` | Computed | - | 下发关联账号，嵌套结构见 D2 |
| `exec_areas` | `RuleGroupResp.ExecAreas` | `TypeList(Elem: TypeString)` | Computed | - | 所属阶段 |
| `direction` | `RuleGroupResp.Direction` | `TypeInt` | Computed | - | 1-入站，0-出站，-1 无方向 |

**`product` 的 Required 判定**：云 API 注释明确「规则组归属产品……必填，不接受 all」，且 `DescribeRuleGroups` 若不传 `Product` 无法收敛范围；用户需求也明确 `GroupName` 与 `Product` 必传 → `Required: true`。同时 `ValidateFunc: validation.StringInSlice([]string{"enterprise_sg","cfw_nat_acl","cfw_vpc_acl","cfw_edge_acl"}, false)`，在 plan 阶段就拦截非法值。

**`product` 的 ForceNew 判定**：`ModifyRuleGroupRequest` 只有 `GroupId` + `GroupName` 两个字段（models.go 3333-3341），**没有** `Product`；用户需求也说明「ModifyRuleGroup 不可变更 Product，在 CreateRuleGroup 时就固定了」→ `ForceNew: true`。

### D2 — `receive_account` 嵌套结构（唯一的多层结构）

```hcl
receive_account {
  uin           = "100000000001"
  nickname      = "xxx"
  receiver_type = 1
  members {
    app_id    = "1300000001"
    uin       = "100000000002"
    nickname  = "yyy"
    member_id = "m-xxxx"
    node_name = "部门A"
  }
}
```

- `receive_account`: `TypeList` / `Computed` / `Elem: &schema.Resource`，含 `uin`(String)、`nickname`(String)、`receiver_type`(Int)、`members`(TypeList)。
- `members`: `TypeList` / `Computed` / `Elem: &schema.Resource`，含 `app_id`(String)、`uin`(String)、`nickname`(String)、`member_id`(String)、`node_name`(String)。
- 依据：`ReceiveAccount` 是 `RuleGroupResp` 的**子结构**（一个规则组的下发账号配置），不是「资源列表层」，因此保留其嵌套结构；而 `RuleGroups` 数组本身**不允许**建 `rule_groups` 包装层，全部平铺到顶层（见 D1）。

### D3 — 不暴露 `filters` / `order` / `by` / `offset` / `limit` / `name` / `values` / `operator_type`

这 8 个参数是 `DescribeRuleGroups` 的**查询侧**入参，不属于资源属性。若暴露到 resource schema：
- `name` / `values` / `operator_type` 会与 `group_name` 语义冲突且用户无法稳定维护；
- `offset` / `limit` 属于分页实现细节（项目规范也要求不暴露分页参数给用户）；
- 用户一旦改动它们，会产生无意义的 diff。

因此 Read 内部**硬编码**过滤器：`Filters = [{Name: "GroupId", Values: [groupId], OperatorType: 9}]`（9 = 模糊匹配，`CommonFilter` 注释见 models.go 330-332；用户已确认该写法可按 `GroupId` 精准命中，因 `GroupId` 为 `fwmrg-` 前缀的唯一值）。`Product` 不参与过滤，避免「入参字符串 / 出参数字」不一致导致的查不到。

### D4 — 资源 ID = `GroupId`（单字段，无复合 ID）

- Create 成功后 `d.SetId(*response.Response.GroupId)`，并 `d.Set("group_id", ...)`。
- `GroupId` 形如 `fwmrg-wsrrhuvrib`，全局唯一（云 API 注释「创建成功的规则组唯一标识，fwmrg- 前缀」），无需与 `Product` 组合 → 不需要 `tccommon.FILED_SP` 拼接/切分。
- Importer：`schema.ImportStatePassthrough`，`terraform import tencentcloud_fwm_rule_group.example fwmrg-wsrrhuvrib`。
- Read / Update / Delete 中直接从 `d.Id()` 取 `groupId` 作为请求参数（符合项目「read/update/delete 从 d.Id() 获取 id 作为 request 一部分」的要求）。

### D5 — Create：同步，无轮询；id 空值必须报错

```go
if result == nil || result.Response == nil {
    return resource.NonRetryableError(fmt.Errorf("Create fwm rule group failed, Response is nil."))
}
log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
if result.Response.GroupId == nil || *result.Response.GroupId == "" {
    return resource.NonRetryableError(fmt.Errorf("Create fwm rule group failed, GroupId is empty."))
}
d.SetId(*result.Response.GroupId)
return resourceTencentCloudFwmRuleGroupRead(d, meta)
```

由于接口同步生效（无 `FlowId`/`TaskId`），拿到 `GroupId` 后**不做**轮询；但为保证 `terraform apply` 后 state 完整，Create 末尾统一调用一次 Read（与 `tencentcloud_igtm_strategy` 一致）。

### D6 — Read：`Product` 数字 → 字符串归一化

`RuleGroupResp.Product` 在 SDK 中声明为 `*string`（models.go 4148-4149），但云 API 实际返回的是**数字字符串**（`"1"`~`"4"`），而 Create 接收的是产品名。若不归一化，`product` 会永久 diff。

```go
var fwmRuleGroupProductToName = map[string]string{
	"1": "enterprise_sg",
	"2": "cfw_nat_acl",
	"3": "cfw_vpc_acl",
	"4": "cfw_edge_acl",
}

func fwmRuleGroupNormalizeProduct(product string) string {
	if name, ok := fwmRuleGroupProductToName[product]; ok {
		return name
	}
	// 已经是产品名（或未来新增产品）时原样返回，避免破坏 state
	return product
}
```

Read 中：`if item.Product != nil { _ = d.Set("product", fwmRuleGroupNormalizeProduct(*item.Product)) }`。

Read 的其他约定：
- `respData == nil`（资源已被删除）时，**先** `log.Printf("[CRUD] tencentcloud_fwm_rule_group id=%s", d.Id())` 保留现场，**再** `d.SetId("")`，返回 nil。
- 每个 `d.Set` 前判断对应 SDK 字段是否为 nil（遵循项目规范）。
- 所有字段均来自 `DescribeRuleGroups` 返回的第一条匹配项（按 `GroupId` 精确比对后再使用）。

### D7 — Update：`mutableArgs = []string{"group_name"}`

```go
mutableArgs := []string{"group_name"}
needChange := false
for _, v := range mutableArgs {
    if d.HasChange(v) { needChange = true; break }
}
if needChange {
    request := fwmv20250611.NewModifyRuleGroupRequest()
    request.GroupId = helper.String(d.Id())
    if v, ok := d.GetOk("group_name"); ok {
        request.GroupName = helper.String(v.(string))
    }
    // resource.Retry(tccommon.WriteRetryTimeout, ...) 调 ModifyRuleGroupWithContext
}
return resourceTencentCloudFwmRuleGroupRead(d, meta)
```

`product` 是 ForceNew，不在 `mutableArgs` 内（Terraform 会走 destroy + create，不会进入 Update）。为了防御性正确性，Update 开头对 `immutableArgs`（`[]string{"product"}`）做检查：若 `d.HasChange` 命中则返回 error，避免未来 schema 改动导致误改产品。

### D8 — Delete：`GroupIds = [groupId]`

```go
request := fwmv20250611.NewDeleteRuleGroupRequest()
request.GroupIds = []*string{helper.String(d.Id())}
// resource.Retry(tccommon.WriteRetryTimeout, ...) 调 DeleteRuleGroupWithContext
```

单删只传一个元素（用户需求）。接口同步成功即视为删除完成，不额外轮询（避免与外层 retry 形成嵌套 retry）。

### D9 — Service 层 + connectivity 客户端工厂

新增 `tencentcloud/services/fwm/service_tencentcloud_fwm.go`：

```go
func NewFwmService(client *connectivity.TencentCloudClient) FwmService
type FwmService struct { client *connectivity.TencentCloudClient }

// 返回 (nil, nil) 表示未找到，错误通过 errRet 返回
func (me *FwmService) DescribeFwmRuleGroupById(ctx context.Context, groupId string) (ret *fwmv20250611.RuleGroupResp, errRet error)
```

- 请求：`DescribeRuleGroupsRequest{ Filters: []*CommonFilter{{Name: helper.String("GroupId"), Values: []*string{helper.String(groupId)}, OperatorType: helper.Int64(9)}}, Limit: helper.Int64(100), Offset: helper.Int64(0) }`。
- **分页**：`Offset` / `Limit` 存在（models.go 2193-2197），但云 API 注释未标注最大值，取 `Limit = 100` 作为保守值；实现为循环分页（每次 `Offset += Limit`），当返回条数 `< Limit` 或累计达到 `TotalCount` 时结束，再按 `GroupId` 精确比对返回。这样即使过滤失效也不会漏。
- `ratelimit.Check(request.GetAction())` + `resource.Retry(tccommon.ReadRetryTimeout, ...)`；`result == nil || result.Response == nil` → `resource.NonRetryableError`。

`tencentcloud/connectivity/client.go` 新增（照抄 igtm 2247-2257 的写法）：

```go
// UseFwmV20250611Client return FWM client for service
func (me *TencentCloudClient) UseFwmV20250611Client() *fwmv20250611.Client {
	if me.fwmV20250611Conn != nil {
		return me.fwmV20250611Conn
	}
	cpf := me.NewClientProfile(300)
	me.fwmV20250611Conn, _ = fwmv20250611.NewClient(me.Credential, me.Region, cpf)
	me.fwmV20250611Conn.WithHttpTransport(&LogRoundTripper{})
	return me.fwmV20250611Conn
}
```

struct 字段（追加到其他 `*Conn` 字段附近，如 269 行 `ga2v20250115Conn` 之后）：`fwmV20250611Conn *fwmv20250611.Client`；import 别名：`fwmv20250611 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"`。

### D10 — Provider 注册

- `tencentcloud/provider.go`：import 块新增 `"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/fwm"`；`ResourcesMap` 新增 `"tencentcloud_fwm_rule_group": fwm.ResourceTencentCloudFwmRuleGroup(),`（与 `tencentcloud_igtm_strategy` 同层，直接新增一行，不改动其他行）。
- `tencentcloud/provider.md`：在文件末尾（`Backup and Disaster Recovery Center(BDRC)` 段之后）追加：

```
Firewall Manager(FWM)
Resource
tencentcloud_fwm_rule_group
```

### D11 — 代码风格参照（实施时无需再回读源文件）

`tencentcloud/services/igtm/resource_tc_igtm_strategy.go` 的结构骨架：

1. `package igtm` → 本次 `package fwm`；imports：`context`/`fmt`/`log`、`helper/resource`、`helper/schema`、SDK 包、`tccommon`、`internal/helper`；本资源额外需要 `helper/validation` 与 `strings`（仅 Update 的 immutableArgs 检查可用 `d.HasChange`，不必 import strings）。
2. `func ResourceTencentCloudXxx() *schema.Resource`：`Create/Read/Update/Delete` + `Importer{schema.ImportStatePassthrough}` + `Schema: map[string]*schema.Schema{...}`，字段之间空行分隔，Computed 字段集中在 `// computed` 注释下（本资源为：`group_id`、`app_id`、`rule_count`、`create_time`、`update_time`、`create_by`、`create_name`、`update_by`、`update_name`、`receive_account`、`exec_areas`、`direction`）。
3. 每个 CRUD 函数首两行：
   ```go
   defer tccommon.LogElapsed("resource.tencentcloud_fwm_rule_group.create")()
   defer tccommon.InconsistentCheck(d, meta)()
   ```
4. 变量声明块：
   ```go
   var (
       logId    = tccommon.GetLogId(tccommon.ContextNil)
       ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
       request  = fwmv20250611.NewCreateRuleGroupRequest()
       response = fwmv20250611.NewCreateRuleGroupResponse()
   )
   ```
5. 参数填充：`if v, ok := d.GetOk("group_name"); ok { request.GroupName = helper.String(v.(string)) }`。
6. 调用与日志：`resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError { result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseFwmV20250611Client().CreateRuleGroupWithContext(ctx, request); if e != nil { return tccommon.RetryError(e) } else { log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString()) } ... })`。
7. 错误日志前缀：`log.Printf("[CRITAL]%s create fwm rule group failed, reason:%+v", logId, reqErr)`（资源统一使用 `fwm_rule_group` 小写蛇形命名）。

### D12 — 单元测试（gomonkey，不使用 terraform 测试套件）

文件 `tencentcloud/services/fwm/resource_tc_fwm_rule_group_test.go`，`package fwm_test`，参照 `resource_tc_bh_bind_device_resource_test.go`：

- `mockMetaFwmRuleGroup{client *connectivity.TencentCloudClient}` 实现 `tccommon.ProviderMeta`。
- `patches.ApplyMethodReturn(meta.client, "UseFwmV20250611Client", &fwmv20250611.Client{})`。
- `patches.ApplyMethodFunc(fwmClient, "CreateRuleGroupWithContext", ...)` / `"DescribeRuleGroupsWithContext"` / `"ModifyRuleGroupWithContext"` / `"DeleteRuleGroupWithContext"`。
- 用例：
  1. `TestFwmRuleGroupResource_Create`：断言请求 `GroupName`/`Product`，返回 `GroupId=fwmrg-xxx`，断言 `d.Id() == "fwmrg-xxx"`、`group_id == "fwmrg-xxx"`。
  2. `TestFwmRuleGroupResource_ReadProductNormalize`：Describe 返回 `Product="2"`、`Direction=-1`，断言 state `product == "cfw_nat_acl"`；`receive_account` / `members` / `exec_areas` 正确落盘。
  3. `TestFwmRuleGroupResource_ReadNotFound`：Describe 返回空 `RuleGroups`，断言 `d.Id() == ""` 且无 error。
  4. `TestFwmRuleGroupResource_Update`：`group_name` 变更触发 `ModifyRuleGroupWithContext`，断言 `GroupId` 取自 `d.Id()`。
  5. `TestFwmRuleGroupResource_Delete`：断言 `GroupIds` 长度为 1 且等于 `d.Id()`。
- 禁止执行 `go test`（由后续流程统一校验），仅保证可编译。

### D13 — 文档 `resource_tc_fwm_rule_group.md`

格式遵循 `resource_tc_igtm_strategy.md`：

1. 一句话描述：`Provides a resource to create a FWM rule group.`
2. `Example Usage`（`product` 使用 `enterprise_sg`）。
3. `Import` 段：`terraform import tencentcloud_fwm_rule_group.example fwmrg-wsrrhuvrib`。
4. 不写 `Argument Reference` / `Attribute Reference`（由 `make doc` 生成）。

### D14 — 关键云 API struct 定义（vendor 摘录，实施时直接引用）

`vendor/.../fwm/v20250611/models.go`：

```go
// 678-686
type CreateRuleGroupRequest struct {
	*tchttp.BaseRequest
	GroupName *string `json:"GroupName,omitnil,omitempty" name:"GroupName"` // 不可为空，最长50字符
	Product   *string `json:"Product,omitnil,omitempty" name:"Product"`     // enterprise_sg/cfw_vpc_acl/cfw_nat_acl/cfw_edge_acl，必填
}

// 709-720
type CreateRuleGroupResponseParams struct {
	GroupId   *string `json:"GroupId,omitnil,omitempty" name:"GroupId"`     // fwmrg- 前缀
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}
type CreateRuleGroupResponse struct {
	*tchttp.BaseResponse
	Response *CreateRuleGroupResponseParams `json:"Response"`
}

// 321-333
type CommonFilter struct {
	Name         *string   `json:"Name,omitnil,omitempty" name:"Name"`
	Values       []*string `json:"Values,omitnil,omitempty" name:"Values"`
	OperatorType *int64    `json:"OperatorType,omitnil,omitempty" name:"OperatorType"` // 1=精确匹配，7=in，9=模糊匹配
}

// 2178-2198
type DescribeRuleGroupsRequest struct {
	*tchttp.BaseRequest
	Product *string          `json:"Product,omitnil,omitempty" name:"Product"`
	Filters []*CommonFilter  `json:"Filters,omitnil,omitempty" name:"Filters"`
	Order   *string          `json:"Order,omitnil,omitempty" name:"Order"`
	By      *string          `json:"By,omitnil,omitempty" name:"By"`
	Offset  *int64           `json:"Offset,omitnil,omitempty" name:"Offset"`
	Limit   *int64           `json:"Limit,omitnil,omitempty" name:"Limit"`
}

// 2225-2239
type DescribeRuleGroupsResponseParams struct {
	TotalCount *int64          `json:"TotalCount,omitnil,omitempty" name:"TotalCount"`
	RuleGroups []*RuleGroupResp `json:"RuleGroups,omitnil,omitempty" name:"RuleGroups"`
	RequestId  *string         `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}
type DescribeRuleGroupsResponse struct {
	*tchttp.BaseResponse
	Response *DescribeRuleGroupsResponseParams `json:"Response"`
}

// 3333-3341
type ModifyRuleGroupRequest struct {
	*tchttp.BaseRequest
	GroupId   *string `json:"GroupId,omitnil,omitempty" name:"GroupId"`
	GroupName *string `json:"GroupName,omitnil,omitempty" name:"GroupName"`
}

// 3364-3372（Response 只有 RequestId）
type ModifyRuleGroupResponseParams struct {
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

// 1216-1221
type DeleteRuleGroupRequest struct {
	*tchttp.BaseRequest
	GroupIds []*string `json:"GroupIds,omitnil,omitempty" name:"GroupIds"`
}

// 1243-1251（Response 只有 RequestId）
type DeleteRuleGroupResponseParams struct {
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

// 4135-4177
type RuleGroupResp struct {
	AppId          *string           `json:"AppId,omitnil,omitempty" name:"AppId"`
	GroupId        *string           `json:"GroupId,omitnil,omitempty" name:"GroupId"`
	GroupName      *string           `json:"GroupName,omitnil,omitempty" name:"GroupName"`
	RuleCount      *int64            `json:"RuleCount,omitnil,omitempty" name:"RuleCount"`
	Product        *string           `json:"Product,omitnil,omitempty" name:"Product"`
	CreateTime     *string           `json:"CreateTime,omitnil,omitempty" name:"CreateTime"`
	UpdateTime     *string           `json:"UpdateTime,omitnil,omitempty" name:"UpdateTime"`
	CreateBy       *string           `json:"CreateBy,omitnil,omitempty" name:"CreateBy"`
	CreateName     *string           `json:"CreateName,omitnil,omitempty" name:"CreateName"`
	UpdateBy       *string           `json:"UpdateBy,omitnil,omitempty" name:"UpdateBy"`
	UpdateName     *string           `json:"UpdateName,omitnil,omitempty" name:"UpdateName"`
	ReceiveAccount []*ReceiveAccount `json:"ReceiveAccount,omitnil,omitempty" name:"ReceiveAccount"`
	ExecAreas      []*string         `json:"ExecAreas,omitnil,omitempty" name:"ExecAreas"`
	Direction      *int64            `json:"Direction,omitnil,omitempty" name:"Direction"` // 1-入站，0-出站，-1 无方向
}

// 4068-4080
type ReceiveAccount struct {
	Uin          *string       `json:"Uin,omitnil,omitempty" name:"Uin"`
	Nickname     *string       `json:"Nickname,omitnil,omitempty" name:"Nickname"`
	ReceiverType *int64        `json:"ReceiverType,omitnil,omitempty" name:"ReceiverType"` // 0=账号uin，1=账号组
	Members      []*MemberInfo `json:"Members,omitnil,omitempty" name:"Members"`           // ReceiverType=1 时返回
}

// 3049-3064
type MemberInfo struct {
	AppId    *string `json:"AppId,omitnil,omitempty" name:"AppId"`
	Uin      *string `json:"Uin,omitnil,omitempty" name:"Uin"`
	Nickname *string `json:"Nickname,omitnil,omitempty" name:"Nickname"`
	MemberId *string `json:"MemberId,omitnil,omitempty" name:"MemberId"`
	NodeName *string `json:"NodeName,omitnil,omitempty" name:"NodeName"`
}
```

客户端方法（`client.go`）：`CreateRuleGroupWithContext(ctx, request) (*CreateRuleGroupResponse, error)`（411 行）、`DeleteRuleGroupWithContext`（845 行）、`DescribeRuleGroupsWithContext`（1495 行）、`ModifyRuleGroupWithContext`（2281 行）。

helper 可用函数（`tencentcloud/internal/helper/transform.go`）：`helper.String(string) *string`、`helper.Int64(int64) *int64`、`helper.Int64ToStr`、`helper.StrToInt64Point`。

## Risks / Trade-offs

- **[`Product` 出参实际为数字，SDK 声明为 string]** → 通过 D6 的归一化映射处理；未命中的值原样保留，保证向前兼容。若后续云 API 新增产品编号，仅需在 `fwmRuleGroupProductToName` 增加一条映射（不会破坏已存 state）。
- **[`DescribeRuleGroups` 的 `Limit` 未标注最大值]** → 取 100 并实现分页循环；即使单页不足也会继续翻页，避免漏查。若后端对 `Limit` 上限更小并报错，可在实施时下调（无 schema 影响）。
- **[`OperatorType=9` 为模糊匹配]** → 仅可能返回额外记录，不会漏；读到结果后按 `GroupId` 字符串精确比对，保证不会误绑定其他规则组。
- **[同步接口存在最终一致性窗口]** → Create 后立刻 Read 可能短暂查不到，因此 Create 末尾的 Read 依赖 `DescribeFwmRuleGroupById` 内部的 `resource.Retry(tccommon.ReadRetryTimeout)`；若仍查不到，Read 会 `d.SetId("")` 造成资源丢失风险 —— 为此 Read 在 `respData == nil` 时先打印 `[CRUD] tencentcloud_fwm_rule_group id=%s` 日志再清空 id，便于排障与手工 import 恢复。
- **[新增服务包与 connectivity 方法]** → 纯新增，不改动既有方法签名，对现有 80+ 服务零影响。
- **[`receive_account` 为纯 Computed 嵌套列表]** → 后端返回顺序不稳定时仅影响 state 展示，不影响 plan（Computed 字段不参与 diff）。

## Migration Plan

1. 新建 `tencentcloud/services/fwm/` 包：`service_tencentcloud_fwm.go` → `resource_tc_fwm_rule_group.go`。
2. 在 `tencentcloud/connectivity/client.go` 增加 `fwmV20250611Conn` 字段与 `UseFwmV20250611Client()`。
3. 注册 provider（`provider.go` + `provider.md`）。
4. 补 `.md` 示例与 gomonkey 单元测试。
5. 收尾阶段由 `tfpacer-finalize` 统一执行 `gofmt` / `make doc` / `.changelog` 生成。
6. 回滚策略：本次全部为新增文件与新增注册行，回滚只需删除新文件并撤销 `provider.go` / `provider.md` / `client.go` 的三处新增行；不会影响已有 state。
