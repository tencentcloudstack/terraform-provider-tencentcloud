## Context

FWM (Tencent Cloud Firewall Manager / 云防火墙规则组管理) is a brand-new product surface in the provider: `tencentcloud/services/fwm/` does not exist yet, `tencentcloud/connectivity/client.go` has no FWM client, and `tencentcloud/provider.go` registers no `tencentcloud_fwm_*` entry.

### SDK state

`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm v1.3.196` (namespace `v20250611`) is already wired into the build:

- `go.mod` → `require github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm v1.3.196`
- `tools.go` → `_ "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"`
- `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611/{client.go,models.go,errors.go}` exists

Import alias to use everywhere: `fwmv20250611 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"`.

### Relevant SDK structs (verbatim, from `vendor/.../fwm/v20250611/models.go`)

`CreateRuleGroupRequest` (line 678) / `CreateRuleGroupRequestParams` (line 670):

```go
type CreateRuleGroupRequest struct {
	*tchttp.BaseRequest

	// 规则组名称，不可为空，最长50字符（按字符数计，中文按1个字符计）
	GroupName *string `json:"GroupName,omitnil,omitempty" name:"GroupName"`

	// 规则组归属产品。enterprise_sg-企业安全组，cfw_vpc_acl-VPC边界防火墙，
	// cfw_nat_acl-NAT边界防火墙，cfw_edge_acl-互联网边界防火墙。必填，不接受 all
	Product *string `json:"Product,omitnil,omitempty" name:"Product"`
}
```

`CreateRuleGroupResponse` (line 717) / `CreateRuleGroupResponseParams` (line 709):

```go
type CreateRuleGroupResponseParams struct {
	// 创建成功的规则组唯一标识，fwmrg- 前缀
	GroupId *string `json:"GroupId,omitnil,omitempty" name:"GroupId"`
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}

type CreateRuleGroupResponse struct {
	*tchttp.BaseResponse
	Response *CreateRuleGroupResponseParams `json:"Response"`
}
```

`DescribeRuleGroupsRequest` (line 2178) / `DescribeRuleGroupsRequestParams` (line 2158):

```go
type DescribeRuleGroupsRequest struct {
	*tchttp.BaseRequest

	// 产品类型
	Product *string `json:"Product,omitnil,omitempty" name:"Product"`
	// 产品类型 (sic — doc comment is copy-pasted; actual type is CommonFilter)
	Filters []*CommonFilter `json:"Filters,omitnil,omitempty" name:"Filters"`
	// 排序
	Order *string `json:"Order,omitnil,omitempty" name:"Order"`
	// 排序字段
	By *string `json:"By,omitnil,omitempty" name:"By"`
	// 起始位置
	Offset *int64 `json:"Offset,omitnil,omitempty" name:"Offset"`
	// 查询数量
	Limit *int64 `json:"Limit,omitnil,omitempty" name:"Limit"`
}
```

`DescribeRuleGroupsResponse` (line 2236) / `DescribeRuleGroupsResponseParams` (line 2225):

```go
type DescribeRuleGroupsResponseParams struct {
	// 总记录数
	TotalCount *int64 `json:"TotalCount,omitnil,omitempty" name:"TotalCount"`
	// 规则组列表
	RuleGroups []*RuleGroupResp `json:"RuleGroups,omitnil,omitempty" name:"RuleGroups"`
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}
```

`CommonFilter` (line 321) — note: the filter struct is named **`CommonFilter`**, not `Filter`:

```go
type CommonFilter struct {
	// 筛选字段名
	Name *string `json:"Name,omitnil,omitempty" name:"Name"`
	// 筛选值列表
	Values []*string `json:"Values,omitnil,omitempty" name:"Values"`
	// 操作类型：1=精确匹配，7=in，9=模糊匹配
	OperatorType *int64 `json:"OperatorType,omitnil,omitempty" name:"OperatorType"`
}
```

`RuleGroupResp` (line 4135) — the describe list element:

```go
type RuleGroupResp struct {
	AppId *string `json:"AppId,omitnil,omitempty" name:"AppId"`                 // 租户Id
	GroupId *string `json:"GroupId,omitnil,omitempty" name:"GroupId"`           // 规则组Id
	GroupName *string `json:"GroupName,omitnil,omitempty" name:"GroupName"`     // 规则组名称
	RuleCount *int64 `json:"RuleCount,omitnil,omitempty" name:"RuleCount"`      // 规则数量
	Product *string `json:"Product,omitnil,omitempty" name:"Product"`           // 产品类型（接口返回数字码）
	CreateTime *string `json:"CreateTime,omitnil,omitempty" name:"CreateTime"`  // 创建时间
	UpdateTime *string `json:"UpdateTime,omitnil,omitempty" name:"UpdateTime"`  // 更新时间
	CreateBy *string `json:"CreateBy,omitnil,omitempty" name:"CreateBy"`        // 创建人SubUin
	CreateName *string `json:"CreateName,omitnil,omitempty" name:"CreateName"`  // 创建人名称
	UpdateBy *string `json:"UpdateBy,omitnil,omitempty" name:"UpdateBy"`        // 更新人SubUin
	UpdateName *string `json:"UpdateName,omitnil,omitempty" name:"UpdateName"`  // 更新人名称
	ReceiveAccount []*ReceiveAccount `json:"ReceiveAccount,omitnil,omitempty" name:"ReceiveAccount"` // 下发关联账号
	ExecAreas []*string `json:"ExecAreas,omitnil,omitempty" name:"ExecAreas"`   // 所属阶段
	Direction *int64 `json:"Direction,omitnil,omitempty" name:"Direction"`      // 1-入站, 0-出站, -1 无方向
}
```

`ReceiveAccount` (line 4068) and `MemberInfo` (line 3049):

```go
type ReceiveAccount struct {
	Uin *string `json:"Uin,omitnil,omitempty" name:"Uin"`                       // 租户 uin
	Nickname *string `json:"Nickname,omitnil,omitempty" name:"Nickname"`        // 租户名称
	ReceiverType *int64 `json:"ReceiverType,omitnil,omitempty" name:"ReceiverType"` // 0=账号uin，1=账号组
	Members []*MemberInfo `json:"Members,omitnil,omitempty" name:"Members"`     // 仅 ReceiverType==1 时返回
}

type MemberInfo struct {
	AppId *string `json:"AppId,omitnil,omitempty" name:"AppId"`
	Uin *string `json:"Uin,omitnil,omitempty" name:"Uin"`
	Nickname *string `json:"Nickname,omitnil,omitempty" name:"Nickname"`
	MemberId *string `json:"MemberId,omitnil,omitempty" name:"MemberId"`
	NodeName *string `json:"NodeName,omitnil,omitempty" name:"NodeName"`
}
```

`ModifyRuleGroupRequest` (line 3333) / `ModifyRuleGroupRequestParams` (line 3325):

```go
type ModifyRuleGroupRequest struct {
	*tchttp.BaseRequest

	// 规则组名称 (sic — doc comment is copy-pasted; this is the group id)
	GroupId *string `json:"GroupId,omitnil,omitempty" name:"GroupId"`
	// 规则组名称
	GroupName *string `json:"GroupName,omitnil,omitempty" name:"GroupName"`
}
```

`ModifyRuleGroupResponse` (line 3369): only `RequestId`, no business payload.

`DeleteRuleGroupRequest` (line 1216) / `DeleteRuleGroupRequestParams` (line 1211):

```go
type DeleteRuleGroupRequest struct {
	*tchttp.BaseRequest

	// 规则组Id列表
	GroupIds []*string `json:"GroupIds,omitnil,omitempty" name:"GroupIds"`
}
```

`DeleteRuleGroupResponse` (line 1248): only `RequestId`, no business payload.

### API behaviour facts (from `vendor/.../fwm/v20250611/client.go` doc comments)

| API | client.go line | Sync/Async | Notes |
|---|---|---|---|
| `CreateRuleGroup` | 402 | synchronous | returns `GroupId` (`fwmrg-` prefix), initial `RuleCount = 0`; counts against rule-group quota; the three walls (VPC/NAT/Edge) are subject to prepaid write gating, enterprise security group is not |
| `DescribeRuleGroups` | 1486 | synchronous | list/query; error code `ResourceNotFound.NotFound` |
| `ModifyRuleGroup` | 2270 | synchronous | `GroupId` mandatory; only `GroupName` is mutable — **`Product` cannot be changed** |
| `DeleteRuleGroup` | 833 | synchronous | accepts `GroupIds` list; single delete = pass one element |

There is **no** `DescribeRuleGroup` (singular) API in this SDK version — verified by grepping `client.go`. Read must therefore go through `DescribeRuleGroups` with a `GroupId` filter.

### Product code mapping (from requirement, needed for Read)

`DescribeRuleGroups` returns `Product` as a numeric code, while `CreateRuleGroup` takes the string form. Read must map back:

| numeric code returned by Describe | string accepted by Create / stored in state |
|---|---|
| `1` | `enterprise_sg` |
| `2` | `cfw_nat_acl` |
| `3` | `cfw_vpc_acl` |
| `4` | `cfw_edge_acl` |

### Style baseline

Primary reference: `tencentcloud/services/igtm/resource_tc_igtm_strategy.go` (full CRUD, service-layer delegation, retry wrapping) and `tencentcloud/services/igtm/service_tencentcloud_igtm.go` (paginated describe helpers). Structure to copy:

- `func ResourceTencentCloudXxx() *schema.Resource` returning `{Create, Read, Update, Delete, Importer: schema.ImportStatePassthrough, Schema: map[string]*schema.Schema{...}}`
- Every CRUD func starts with `defer tccommon.LogElapsed("resource.tencentcloud_fwm_rule_group.create")()` and `defer tccommon.InconsistentCheck(d, meta)()`
- Var block: `logId = tccommon.GetLogId(tccommon.ContextNil)`, `ctx = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)`, `request = fwmv20250611.NewXxxRequest()`, `response = fwmv20250611.NewXxxResponse()`
- SDK call inside `resource.Retry(tccommon.WriteRetryTimeout, ...)` for mutations and `resource.Retry(tccommon.ReadRetryTimeout, ...)` for reads; on error `return tccommon.RetryError(e)`; on success `log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())`
- Nil guard inside retry: `if result == nil || result.Response == nil { return resource.NonRetryableError(fmt.Errorf("... failed, Response is nil.")) }`
- Service helper `func (me *FwmService) DescribeXxxById(ctx context.Context, id string) (ret *fwmv20250611.RuleGroupResp, errRet error)` with the `defer func(){ if errRet != nil { log.Printf("[CRITAL]...") } }()` pattern and `ratelimit.Check(request.GetAction())` before each call
- Read sets each field guarded: `if respData.Field != nil { _ = d.Set("field", respData.Field) }`

## Goals / Non-Goals

**Goals:**

- New resource `tencentcloud_fwm_rule_group` (RESOURCE_KIND_GENERAL) with full synchronous CRUD plus `terraform import`.
- Schema exposes every user-settable/observable attribute of the four APIs, with the describe list flattened to top level (no `rule_groups` wrapper layer).
- Resource ID = bare `group_id` (`fwmrg-xxx`).
- Read looks up a single group via `DescribeRuleGroups` + `Filters=[{Name:"GroupId", Values:[id], OperatorType:9}]`, paginating with `Limit=100`.
- `product` is `ForceNew` + validated against the four allowed strings; Read normalizes the returned numeric code back to the string.
- All SDK calls wrapped in `resource.Retry` with `tccommon.ReadRetryTimeout` / `tccommon.WriteRetryTimeout` and `tccommon.RetryError`.
- New connectivity client `UseFwmV20250611Client`, provider registration, `.md` doc, gomonkey unit tests.

**Non-Goals:**

- No data source for rule groups in this change.
- No rule (`CreateSecurityGroupRule` / `CreateVpcAclRule` / `CreateNatAclRule` / `CreateEdgeAclRule`) resources — this change only ships the rule *group* container.
- Do **not** expose query-control parameters as resource attributes: `filters`, `name`, `values`, `operator_type`, `order`, `by`, `offset`, `limit` (from `DescribeRuleGroups` request) and `group_ids` (from `DeleteRuleGroup` request). They are internal query/batch controls, not properties of a single rule group; exposing them would create write-only, non-idempotent state. Read/Delete hardcode them instead (see D4/D6).
- No async task polling: all four APIs are synchronous, so no `Timeouts` block and no wait loops are required beyond a read-after-create visibility retry (D7).

## Decisions

### D1: Resource naming and file layout

- Resource: `tencentcloud_fwm_rule_group`
- Package: `tencentcloud/services/fwm` (package `fwm`)
- Files: `resource_tc_fwm_rule_group.go`, `resource_tc_fwm_rule_group.md`, `resource_tc_fwm_rule_group_test.go`, `service_tencentcloud_fwm.go`
- Import alias: `fwmv20250611`

### D2: Schema — describe list flattened to top level

`DescribeRuleGroups` returns `RuleGroups []*RuleGroupResp`. Per project rules, **no** `rule_groups` wrapper layer is created; the fields of one `RuleGroupResp` are flattened into top-level attributes:

| Schema field | Type | Requirement | SDK source | Notes |
|---|---|---|---|---|
| `group_name` | string | Required | `CreateRuleGroupRequest.GroupName`, `ModifyRuleGroupRequest.GroupName`, `RuleGroupResp.GroupName` | mutable; max 50 chars |
| `product` | string | Required + ForceNew | `CreateRuleGroupRequest.Product`, `RuleGroupResp.Product` | validated enum; immutable because Modify cannot change it |
| `group_id` | string | Computed | `CreateRuleGroupResponse.GroupId`, `RuleGroupResp.GroupId` | also the resource ID |
| `app_id` | string | Computed | `RuleGroupResp.AppId` | |
| `rule_count` | int | Computed | `RuleGroupResp.RuleCount` | |
| `create_time` | string | Computed | `RuleGroupResp.CreateTime` | |
| `update_time` | string | Computed | `RuleGroupResp.UpdateTime` | |
| `create_by` | string | Computed | `RuleGroupResp.CreateBy` | |
| `create_name` | string | Computed | `RuleGroupResp.CreateName` | |
| `update_by` | string | Computed | `RuleGroupResp.UpdateBy` | |
| `update_name` | string | Computed | `RuleGroupResp.UpdateName` | |
| `receive_account` | list(Resource) | Computed | `RuleGroupResp.ReceiveAccount` | nested, see below |
| `exec_areas` | list(string) | Computed | `RuleGroupResp.ExecAreas` | |
| `direction` | int | Computed | `RuleGroupResp.Direction` | `1` inbound, `0` outbound, `-1` no direction |

`receive_account` element schema:

| Field | Type | Requirement | SDK source |
|---|---|---|---|
| `uin` | string | Computed | `ReceiveAccount.Uin` |
| `nickname` | string | Computed | `ReceiveAccount.Nickname` |
| `receiver_type` | int | Computed | `ReceiveAccount.ReceiverType` |
| `members` | list(Resource) | Computed | `ReceiveAccount.Members` (`[]*MemberInfo`) |

`members` element schema: `app_id` (string), `uin` (string), `nickname` (string), `member_id` (string), `node_name` (string) — all Computed.

`product` validation: `ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"enterprise_sg", "cfw_nat_acl", "cfw_vpc_acl", "cfw_edge_acl"})` — helper is `ValidateAllowedStringValue` in `tencentcloud/common/validators.go` (line 146), already used by e.g. `tencentcloud/services/ssm/resource_tc_ssm_ssh_key_pair_secret.go:75`.

### D3: Resource ID = `group_id`

`CreateRuleGroup` returns exactly one identifier and all other APIs take only `GroupId` (Delete takes a list containing it). Therefore the ID is the bare `fwmrg-` string, set with `d.SetId(*response.Response.GroupId)`, and `Importer: &schema.ResourceImporter{State: schema.ImportStatePassthrough}`. No `tccommon.FILED_SP` joining is needed. Read/Update/Delete use `d.Id()` directly as `GroupId` / `GroupIds[0]`.

### D4: Read = paginated `DescribeRuleGroups` filtered by `GroupId`

```go
request.Filters = []*fwmv20250611.CommonFilter{
	{
		Name:         helper.String("GroupId"),
		Values:       []*string{helper.String(groupId)},
		OperatorType: helper.Int64(9), // 9 = 模糊匹配, per API owner guidance
	},
}
```

Pagination: `Offset` starting at `0`, `Limit = 100` (no documented maximum in the SDK comment → use 100, the provider-wide convention). Loop until `len(RuleGroups) < limit`. Because `OperatorType=9` is a fuzzy match, the helper must do an **exact** final match on `GroupId` client-side and return only that element. If nothing matches → return `(nil, nil)`; the Read function then logs `[CRUD] tencentcloud_fwm_rule_group id=%s` and calls `d.SetId("")`.

Helper signature (service layer):

```go
func (me *FwmService) DescribeFwmRuleGroupById(ctx context.Context, groupId string) (ret *fwmv20250611.RuleGroupResp, errRet error)
func (me *FwmService) DescribeFwmRuleGroupsByFilter(ctx context.Context, param map[string]interface{}) (ret []*fwmv20250611.RuleGroupResp, errRet error)
```

`DescribeFwmRuleGroupsByFilter` accepts a `param` map keyed like the existing igtm helpers: `"Product"` (`*string`), `"Filters"` (`[]*fwmv20250611.CommonFilter`), `"Order"` (`*string`), `"By"` (`*string`); always overrides `Offset`/`Limit` internally.

### D5: `product` numeric→string normalization on Read

The describe response returns the product as a numeric code while the request takes a string. Read must translate so Terraform does not see a diff:

```go
func fwmRuleGroupProductCodeToName(v string) string {
	switch v {
	case "1":
		return "enterprise_sg"
	case "2":
		return "cfw_nat_acl"
	case "3":
		return "cfw_vpc_acl"
	case "4":
		return "cfw_edge_acl"
	default:
		return v // already a name
	}
}
```

Applied as `if respData.Product != nil { _ = d.Set("product", fwmRuleGroupProductCodeToName(*respData.Product)) }`. The helper is tolerant: if the API ever starts returning the string name, it passes through unchanged.

Implementation note: keep this as a small unexported func in `resource_tc_fwm_rule_group.go` so the unit test can exercise it indirectly through Read.

### D6: Create / Update / Delete shapes

Create:
1. Build `CreateRuleGroupRequest` from `group_name` + `product`.
2. `resource.Retry(tccommon.WriteRetryTimeout, ...)` → `CreateRuleGroupWithContext`.
3. Nil guard, then require `response.Response.GroupId != nil && *response.Response.GroupId != ""`; otherwise `resource.NonRetryableError` (never write an empty id).
4. `d.SetId(*response.Response.GroupId)` — synchronous API, so ID is set immediately, before any read-back.
5. Call Read.

Update: only `group_name` is mutable. Guard with `if d.HasChange("group_name")`, then `ModifyRuleGroup` with `GroupId = d.Id()` and `GroupName`. `product` is `ForceNew`, so it never reaches Update.

Delete: `DeleteRuleGroup` with `GroupIds: []*string{&groupId}`.

### D7: Read-after-create visibility retry

`CreateRuleGroup` is synchronous, but list APIs are eventually consistent in Tencent Cloud. If Create immediately called Read and the group was not yet visible, Read would `d.SetId("")` and silently drop the freshly created resource. To avoid this, Create wraps a **visibility poll** before calling Read:

```go
reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
	group, e := service.DescribeFwmRuleGroupById(ctx, groupId)
	if e != nil {
		return tccommon.RetryError(e)
	}
	if group == nil {
		return resource.RetryableError(fmt.Errorf("tencentcloud_fwm_rule_group %s not visible yet", groupId))
	}
	return nil
})
```

The same helper is *not* used in Read (Read treats not-found as "deleted" → `d.SetId("")`).

### D8: Connectivity client

Add to `tencentcloud/connectivity/client.go`, mirroring `UseIgtmV20231024Client` (line 2247):

```go
// import block
fwmv20250611 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"

// struct field
fwmv20250611Conn *fwmv20250611.Client

// new method
// UseFwmV20250611Client return FWM client for service
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

### D9: Provider registration

`tencentcloud/provider.go`: import `"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/fwm"` and add to the resources map (alphabetical slot, formatted by gofmt):

```go
"tencentcloud_fwm_rule_group": fwm.ResourceTencentCloudFwmRuleGroup(),
```

`tencentcloud/provider.md`: the file is a plain list grouped by product, each group being a product heading line followed by `Data Source` and/or `Resource` subsections (e.g. `Global Accelerator(GA2)` → `Resource` → one resource name per line). Add a new product group in the position matching the file's existing product ordering:

```
Tencent Cloud Firewall Manager(FWM)
Resource
tencentcloud_fwm_rule_group
```

### D10: Documentation `.md`

`tencentcloud/services/fwm/resource_tc_fwm_rule_group.md` follows the project format:

1. One-line description naming the product: `Provides a resource to create a FWM rule group`
2. `Example Usage` with an HCL block using `tencentcloud_fwm_rule_group` (`group_name`, `product = "enterprise_sg"`)
3. `Import` section: `terraform import tencentcloud_fwm_rule_group.example fwmrg-wsrrhuvrib`
4. **No** `Argument Reference` / `Attribute Reference` sections — generated by `make doc`

### D11: Tests — gomonkey unit tests, not acceptance tests

New resources are unit-tested with `gomonkey` (no real cloud calls). Reference pattern: `tencentcloud/services/sqlserver/resource_tc_sqlserver_db_instance_ssl_config_test.go`.

- Package `fwm_test`
- Local mock meta implementing `tccommon.ProviderMeta` (`GetAPIV3Conn() *connectivity.TencentCloudClient`), constructed as `&connectivity.TencentCloudClient{}`
- `patches.ApplyMethodReturn(metaClient, "UseFwmV20250611Client", fwmClient)` and `patches.ApplyMethodFunc(fwmClient, "<Api>WithContext", func(ctx context.Context, req *fwmv20250611.XxxRequest) (*fwmv20250611.XxxResponse, error) {...})`
- `schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{...})` to build `*schema.ResourceData`
- Cases: `Create_Success` (asserts `d.Id() == "fwmrg-xxx"`, `group_name`, normalized `product`), `Create_NilGroupId` (expects error), `Read_Success` (asserts nested `receive_account.0.members.0.*`, `exec_areas`, `direction`), `Read_ProductCodeMapping` (code `2` → `cfw_nat_acl`), `Read_NotFound` (asserts `d.Id() == ""`), `Update_GroupName` (asserts `ModifyRuleGroupRequest` carries `GroupId` + `GroupName`), `Delete_SingleGroupId` (asserts `GroupIds` has exactly one element), `Schema` (Required/ForceNew/Computed flags).
- Because the service helper wraps calls in `resource.Retry`, mocked funcs must return immediately (success path) so tests do not sleep.
- Tests must not be executed in this pipeline; they only need to compile.

### D12: Error handling and logging conventions

- All errors from SDK calls go through `tccommon.RetryError(e)` inside retry blocks and are logged with `log.Printf("[CRITAL]%s ...", logId, err)` before returning.
- Log wording uses the resource name `tencentcloud_fwm_rule_group` (never "该资源"/"当前资源").
- Read's not-found path logs `log.Printf("[CRUD] tencentcloud_fwm_rule_group id=%s", d.Id())` **before** `d.SetId("")`.

## Risks / Trade-offs

- [Risk] `DescribeRuleGroups` is a list API used as get-by-id; a fuzzy `OperatorType=9` filter could match more than one group. → Mitigation: the service helper matches `GroupId` exactly client-side and returns a single element (D4).
- [Risk] `Product` wire type is genuinely numeric while the SDK declares `*string`; if the API returned a raw JSON number, the whole response decode would fail with `cannot unmarshal number into Go struct field`. → Mitigation: the requirement states the value is a numeric *code*; the normalization helper accepts both `"1"` and `"enterprise_sg"`. If decode errors surface in practice, the fallback is an `encoding/json` `Number`-tolerant workaround in the service layer (documented but not implemented).
- [Risk] Create→Read eventual consistency could clear a freshly created ID. → Mitigation: D7 visibility poll inside `Create`.
- [Risk] Rule groups consume quota and the three-walls products have prepaid write gating (`CreateRuleGroup` doc comment). → Not solvable in code; the `.md` example uses `enterprise_sg` to avoid gating in examples, and API errors are surfaced verbatim.
- [Trade-off] `filters`/`order`/`by`/`offset`/`limit`/`group_ids` are not exposed as attributes even though they appear in the API mapping; they are list-query/batch controls, not resource properties (see Non-Goals).
- [Trade-off] All describe-derived fields are `Computed` only, so a user cannot pre-set `exec_areas` / `direction` at create time — the create API accepts neither.
