## Context

Tencent Cloud **Firewall Manager** (`fwm`, 云防火墙企业版 / 企业安全组) exposes an enterprise security group rule model through `fwm/v20250611`. The provider has **no** `fwm` service package today (`tencentcloud/services/fwm/` does not exist) and `tencentcloud/connectivity/client.go` has **no** `fwm` client accessor, so this change bootstraps the whole service layer.

Vendored SDK state (already present on the branch, verified):

- `vendor/modules.txt` line 1341 lists `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611` under `fwm v1.3.196`.
- `vendor/.../fwm/v20250611/{client.go,models.go,errors.go}` exist (`models.go` 4632 lines, `client.go` 2562 lines).
- No `go.mod` change is required for this change itself.

The four CRUD-ish APIs used here are **synchronous** (they return only `RequestId`, or a resource id immediately — no `TaskId` / `FlowId`), so no polling loop and no `Timeouts` block are introduced.

Vendor method inventory (line numbers are in `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611/`):

| SDK method | non-ctx | ctx variant | Returns |
|---|---|---|---|
| `CreateRuleGroup` | `client.go:402` | `client.go:411` | `GroupId *string` |
| `CreateSecurityGroupRule` | `client.go:456` | `client.go:469` | `RuleIds []*string` |
| `DescribeSecurityGroupRules` | `client.go:1596` | `client.go:1607` | `Rules []*SecGroupRuleResp`, `TotalCount`, `AllTotalCount` |
| `ModifySecurityGroupRule` | `client.go:2324` | `client.go:2335` | only `RequestId` |
| `DeleteSecurityGroupRule` | `client.go:897` | `client.go:917` | only `RequestId` |
| `DeleteRuleGroup` | `client.go:833` | `client.go:845` | only `RequestId` |

Code-style reference files (read these; they define the exact shape to imitate):

- `tencentcloud/services/igtm/resource_tc_igtm_strategy.go` (680 lines) — the mandated style reference for RESOURCE_KIND_GENERAL: `ResourceTencentCloudXxx()` → `resourceTencentCloudXxxCreate/Read/Update/Delete`, `defer tccommon.LogElapsed(...)` + `defer tccommon.InconsistentCheck(d, meta)()`, `logId = tccommon.GetLogId(tccommon.ContextNil)`, `ctx = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)`, `resource.Retry(tccommon.WriteRetryTimeout|ReadRetryTimeout, ...)` + `tccommon.RetryError(e)`, `if v, ok := d.GetOk("x"); ok { request.X = helper.String(v.(string)) }`, and Read guards `if resp.Field != nil { _ = d.Set(...) }`.
- `tencentcloud/services/cfw/resource_tc_cfw_edge_policy.go` (421 lines) — the closest **sibling product** precedent: a *single-rule* resource on top of a *batch* ACL API (`AddAclRule` with a one-element `Data` list), `uuid` returned → `ValidRuleId` used as `d.Id()`, Read via `Describe*` + filter, Delete by id.
- `tencentcloud/services/cfw/resource_tc_cfw_cluster_fw_bypass_config_test.go` — gomonkey unit-test precedent: a local `mockMetaXxx` struct implementing `tccommon.ProviderMeta` (`GetAPIV3Conn() *connectivity.TencentCloudClient`), `patches.ApplyMethodReturn(client, "UseCfwV20190904Client", cfwClient)`, then `patches.ApplyMethodFunc(cfwClient, "XxxWithContext", func(ctx, req) (resp, error) {...})`.
- `tencentcloud/services/igtm/service_tencentcloud_igtm.go` — service-layer shape: `NewXxxService(client *connectivity.TencentCloudClient)`, `type XxxService struct { client *connectivity.TencentCloudClient }`, `DescribeXxxById` with `defer func(){ if errRet != nil { log.Printf("[CRITAL]%s api[%s] fail ...") } }()`, pagination loop with `offset`/`limit` and `ratelimit.Check(request.GetAction())`.
- `tencentcloud/connectivity/client.go:2247-2257` (`UseIgtmV20231024Client`) and `:2259-2269` (`UseGa2V20250115Client`) — the exact accessor pattern to copy; cache fields live in the struct near line 268 (`igtmv20231024Conn *igtmv20231024.Client`, `ga2v20250115Conn *ga2v20250115.Client`). Import alias style: `ga2v20250115 "github.com/tencentcloud/.../ga2/v20250115"` (line 74) → use `fwmv20250611 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"`.

### Vendor struct excerpts the implementation must match

`CreateRuleGroupRequest` (models.go:678-686) / `CreateRuleGroupResponseParams` (:709-715):
```go
type CreateRuleGroupRequest struct {
	*tchttp.BaseRequest
	GroupName *string `json:"GroupName,omitnil,omitempty" name:"GroupName"` // 规则组名称，不可为空，最长50字符
	Product   *string `json:"Product,omitnil,omitempty" name:"Product"`     // enterprise_sg | cfw_vpc_acl | cfw_nat_acl | cfw_edge_acl，不接受 all
}
type CreateRuleGroupResponseParams struct {
	GroupId   *string `json:"GroupId,omitnil,omitempty" name:"GroupId"` // fwmrg- 前缀
	RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}
```

`CreateSecurityGroupRuleRequest` (models.go:813-821) / `CreateSecurityGroupRuleResponseParams` (:844-850):
```go
type CreateSecurityGroupRuleRequest struct {
	*tchttp.BaseRequest
	GroupId *string              `json:"GroupId,omitnil,omitempty" name:"GroupId"`
	Rules   []*SecurityGroupRule `json:"Rules,omitnil,omitempty" name:"Rules"`
}
type CreateSecurityGroupRuleResponseParams struct {
	RuleIds   []*string `json:"RuleIds,omitnil,omitempty" name:"RuleIds"`
	RequestId *string   `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}
```

`DescribeSecurityGroupRulesRequest` (models.go:2331-2345) / `...ResponseParams` (:2370-2382):
```go
type DescribeSecurityGroupRulesRequest struct {
	*tchttp.BaseRequest
	GroupId *string         `json:"GroupId,omitnil,omitempty" name:"GroupId"`
	Filters []*CommonFilter `json:"Filters,omitnil,omitempty" name:"Filters"`
	Offset  *int64          `json:"Offset,omitnil,omitempty" name:"Offset"`
	Limit   *int64          `json:"Limit,omitnil,omitempty" name:"Limit"`
}
type DescribeSecurityGroupRulesResponseParams struct {
	Rules         []*SecGroupRuleResp `json:"Rules,omitnil,omitempty" name:"Rules"`
	TotalCount    *int64              `json:"TotalCount,omitnil,omitempty" name:"TotalCount"`
	AllTotalCount *int64              `json:"AllTotalCount,omitnil,omitempty" name:"AllTotalCount"`
	RequestId     *string             `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}
```

`ModifySecurityGroupRuleRequest` (models.go:3394-3402) — note the **singular** `Rule`, and the response carries only `RequestId` (no echo of the rule):
```go
type ModifySecurityGroupRuleRequest struct {
	*tchttp.BaseRequest
	GroupId *string             `json:"GroupId,omitnil,omitempty" name:"GroupId"`
	Rule    *SecurityGroupRule  `json:"Rule,omitnil,omitempty" name:"Rule"` // 必须包含 RuleId
}
```

`DeleteSecurityGroupRuleRequest` (models.go:1273-1281) / `DeleteRuleGroupRequest` (:1216-1221):
```go
type DeleteSecurityGroupRuleRequest struct {
	*tchttp.BaseRequest
	GroupId *string   `json:"GroupId,omitnil,omitempty" name:"GroupId"`
	RuleIds []*string `json:"RuleIds,omitnil,omitempty" name:"RuleIds"`
}
type DeleteRuleGroupRequest struct {
	*tchttp.BaseRequest
	GroupIds []*string `json:"GroupIds,omitnil,omitempty" name:"GroupIds"`
}
```

`SecurityGroupRule` — the write model, used by both `CreateSecurityGroupRule.Rules[]` and `ModifySecurityGroupRule.Rule` (models.go:4314-4371). **All fields except `IpVersion` are documented as nullable**:
```go
type SecurityGroupRule struct {
	IpVersion        *string `json:"IpVersion,omitnil,omitempty" name:"IpVersion"`        // ip 类型
	SourceContent    *string `json:"SourceContent,omitnil,omitempty" name:"SourceContent"` // 源地址
	SourceType       *string `json:"SourceType,omitnil,omitempty" name:"SourceType"`       // 源类型（string!）
	DestContent      *string `json:"DestContent,omitnil,omitempty" name:"DestContent"`     // 目的地址
	DestType         *string `json:"DestType,omitnil,omitempty" name:"DestType"`           // 目的类型（string!）
	Protocol         *string `json:"Protocol,omitnil,omitempty" name:"Protocol"`           // 协议
	Port             *string `json:"Port,omitnil,omitempty" name:"Port"`                   // 端口
	ServiceTemplateId *string `json:"ServiceTemplateId,omitnil,omitempty" name:"ServiceTemplateId"`
	RuleAction       *string `json:"RuleAction,omitnil,omitempty" name:"RuleAction"`       // 动作（string）
	Description      *string `json:"Description,omitnil,omitempty" name:"Description"`
	OrderIndex       *int64  `json:"OrderIndex,omitnil,omitempty" name:"OrderIndex"`       // 优先级 int64
	RuleId           *string `json:"RuleId,omitnil,omitempty" name:"RuleId"`               // modify 时必填
	Scope            *string `json:"Scope,omitnil,omitempty" name:"Scope"`                 // 生效范围
	ProtocolPortType *int64  `json:"ProtocolPortType,omitnil,omitempty" name:"ProtocolPortType"`
	BelongMemberId   *string `json:"BelongMemberId,omitnil,omitempty" name:"BelongMemberId"`
}
```

`SecGroupRuleResp` — the read model (models.go:4179-4268). **Watch the type/semantic drift vs `SecurityGroupRule`**:
```go
type SecGroupRuleResp struct {
	OrderIndex        *int64          // 排序        (int64, same as request)
	RuleId            *string         // 主键 id
	IpVersion         *string
	SourceId          *string         // 源规则内容   (≈ request.SourceContent)
	SourceType        *int64          // 源规则类型   (int64! request.SourceType is string)
	TargetId          *string         // 目的规则内容 (≈ request.DestContent)
	TargetType        *int64          // 目的规则类型 (int64! request.DestType is string)
	Protocol          *string
	Port              *string
	Strategy          *int64          // 策略         (int64! request.RuleAction is string)
	Detail            *string         // 描述         (request.Description is not echoed back)
	Region            *string
	ServiceTemplateId *string
	SouInstanceName   *string
	SouPublicIp       *string
	SouPrivateIp      *string
	SouCidr           *string
	SouParameterName  *string
	InstanceName      *string
	PublicIp          *string
	PrivateIp         *string
	Cidr              *string
	ParameterName     *string
	ProtocolPortName  *string
	Id                *int64          // 规则 id，等同 RuleUuid（int64）
	DnsParseCount     *SgDnsParseCount // nullable
	Scope             *string
	IsNew             *int64          // 0/1 最近一次是否有改动
	BelongMember      *MemberInfo     // nullable
}
```

`CommonFilter` (models.go:321-333), `SgDnsParseCount` (:4391-4397), `MemberInfo` (:3049-3064):
```go
type CommonFilter struct {
	Name         *string   // 筛选字段名
	Values       []*string // 筛选值列表
	OperatorType *int64    // 1=精确匹配，7=in，9=模糊匹配
}
type SgDnsParseCount struct {
	ValidCount   *int64 // 有效下发的 IP 个数
	InvalidCount *int64 // 未下发的 IP 个数
}
type MemberInfo struct {
	AppId    *string
	Uin      *string
	Nickname *string
	MemberId *string // 个人场景为空
	NodeName *string // 所属部门
}
```

## Goals / Non-Goals

**Goals:**

- Full lifecycle (create / read / update / delete / import) of one enterprise security group rule through Terraform, backed by `fwm/v20250611`.
- Honour the API's actual shape: `CreateSecurityGroupRule` is **batch-shaped** but every other call (`ModifySecurityGroupRule.Rule`, per-`RuleId` delete, `RuleId`-filtered describe) is **single-rule shaped**.
- Bootstrap the owning rule group automatically (`CreateRuleGroup` with `Product = "enterprise_sg"`) when the user does not point at an existing `group_id`, and clean it up on destroy so no `fwmrg-*` group leaks.
- Keep the crash surface small: every SDK call wrapped in `resource.Retry`, every response nil-checked, one deterministic `<group_id>#<rule_id>` id.
- Land the service layer (`FwmService`, `UseFwmV20250611Client`) and register the resource in `provider.go` / `provider.md`.

**Non-Goals:**

- A `tencentcloud_fwm_rule_group` resource (group bootstrap/teardown is a side effect of this resource; a standalone group resource is a separate change).
- A `tencentcloud_fwm_security_group_rules` data source.
- Supporting `cfw_vpc_acl` / `cfw_nat_acl` / `cfw_edge_acl` rule groups (this resource hard-codes `Product = "enterprise_sg"`; the other products have their own `cfw` resources).
- Modelling the whole `Rules[]` batch as one resource (see D1).
- Any `Timeouts` block: all four entry points are synchronous (see D8).

## Decisions

### D1. One rule per resource instance (batch API called with a one-element slice)

`CreateSecurityGroupRule` takes `Rules []*SecurityGroupRule` and returns `RuleIds []*string`, which superficially suggests a "rule list" resource. It is rejected for these reasons:

- `DescribeSecurityGroupRules` can pinpoint exactly **one** rule (`Filters = [{Name:"RuleId", Values:[<id>], OperatorType:9}]`), which is the operation a per-resource Read needs.
- `ModifySecurityGroupRule` mutates exactly **one** rule (`Rule` must carry `RuleId`), so a batch resource would fan out into N modify calls anyway.
- `DeleteSecurityGroupRule` semantics given by the product team: one `RuleId` in `RuleIds` for a single-rule delete.
- A batch resource has no stable element identity before create (rule ids are server-generated), so Terraform-side element matching would have to be content-based — fragile and diff-prone.

Therefore: schema describes **one rule**, Create submits `Rules: []*SecurityGroupRule{rule}`, and `RuleIds[0]` is persisted as `rule_id` (and as segment 2 of `d.Id()`).

Alternative considered: a `rules` `TypeList` block whose elements carry a computed `rule_id`; rejected because list-length changes would require create/delete reconciliation with no reliable element identity, and would ship a resource that cannot reliably converge.

### D2. Resource ID = `<group_id>#<rule_id>` (2 segments, `tccommon.FILED_SP`)

Read needs `(GroupId, RuleId)`, Modify needs `(GroupId, Rule.RuleId)`, Delete needs `(GroupId, RuleIds=[ruleId])`. Persisting the composite removes any need to re-discover the group on every operation and makes `terraform import` a plain two-value string:

```
terraform import tencentcloud_fwm_security_group_rule.example fwmrg-xxxxxxxx#fwmrule-xxxxxxxx
```

Parsing: `strings.Split(d.Id(), tccommon.FILED_SP)`, require `len == 2` and both segments non-empty, otherwise `fmt.Errorf("id is broken,%s", d.Id())` (identical to `igtm_strategy`).

Alternative considered: bare `rule_id` as the id; rejected because Read/Modify/Delete all require `GroupId`, which would then have to be re-derived (impossible without listing all groups) or read from state (broken after import).

### D3. Rule-group ownership: create on demand, delete on destroy, never with `group_id`

Schema exposes:

- `group_id` — `Optional` + `Computed` + `ForceNew`: id of an existing rule group to attach the rule to.
- `group_name` — `Optional` + `ForceNew`: name used to create a **new** rule group (`CreateRuleGroup.GroupName`, max 50 chars) when `group_id` is absent.

Exactly one of the two must be supplied. Since `group_id` is `Computed`, an `ExactlyOneOf` schema constraint would fire on the unknown/computed value, so the check is done at the top of Create and returns `resource.NonRetryableError`:

```go
_, hasGroupId := d.GetOk("group_id")
_, hasGroupName := d.GetOk("group_name")
if !hasGroupId && !hasGroupName { ... "either group_id or group_name must be specified." }
if hasGroupId && hasGroupName { ... "group_id and group_name are mutually exclusive." }
```

Ownership rule (also drives Delete): `group_name != ""` ⇒ this resource created the group ⇒ Delete also calls `DeleteRuleGroup(GroupIds: [groupId])` after the rule delete. `group_id` supplied ⇒ the group pre-exists ⇒ Delete leaves it alone. After `terraform import` (where `group_name` is empty) the group is therefore never deleted — the safe default.

Violation of `CreateRuleGroup.Product`: `Product` is hard-coded to `"enterprise_sg"` (the only product compatible with `CreateSecurityGroupRule`); it is not exposed in the schema.

### D4. Retry topology

Every SDK call is wrapped in `resource.Retry`:

- Read path (`DescribeSecurityGroupRules`): `tccommon.ReadRetryTimeout`.
- Write paths (`CreateRuleGroup`, `CreateSecurityGroupRule`, `ModifySecurityGroupRule`, `DeleteSecurityGroupRule`, `DeleteRuleGroup`): `tccommon.WriteRetryTimeout`.
- Errors are classified with `tccommon.RetryError(e)`.
- Nothing but the API call happens inside the retry closure: id computation / `d.SetId` / logging-after-success stay outside, after the `reqErr != nil` branch.
- Service-layer retries are not double-wrapped: `DescribeFwmSecurityGroupRuleById` owns its own `resource.Retry` (read path), and the resource Read simply returns `tccommon.RetryError`-free errors (`if err != nil { return err }`).

### D5. Schema (complete, declaration order)

| # | Schema field | Type | Flags | SDK source |
|---|---|---|---|---|
| 1 | `group_id` | String | Optional, Computed, ForceNew | `CreateSecurityGroupRule.GroupId` / `Describe.GroupId` / `Modify.GroupId` / `Delete.GroupId` |
| 2 | `group_name` | String | Optional, ForceNew | `CreateRuleGroup.GroupName` (create-only) |
| 3 | `ip_version` | String | Optional, Computed | `SecurityGroupRule.IpVersion` |
| 4 | `source_content` | String | Optional, Computed | `SecurityGroupRule.SourceContent` |
| 5 | `source_type` | String | Optional, Computed | `SecurityGroupRule.SourceType` |
| 6 | `dest_content` | String | Optional, Computed | `SecurityGroupRule.DestContent` |
| 7 | `dest_type` | String | Optional, Computed | `SecurityGroupRule.DestType` |
| 8 | `protocol` | String | Optional, Computed | `SecurityGroupRule.Protocol` |
| 9 | `port` | String | Optional, Computed | `SecurityGroupRule.Port` |
| 10 | `service_template_id` | String | Optional, Computed | `SecurityGroupRule.ServiceTemplateId` |
| 11 | `rule_action` | String | Optional, Computed | `SecurityGroupRule.RuleAction` |
| 12 | `description` | String | Optional, Computed | `SecurityGroupRule.Description` |
| 13 | `order_index` | Int | Optional, Computed | `SecurityGroupRule.OrderIndex` |
| 14 | `scope` | String | Optional, Computed | `SecurityGroupRule.Scope` |
| 15 | `protocol_port_type` | Int | Optional, Computed | `SecurityGroupRule.ProtocolPortType` |
| 16 | `belong_member_id` | String | Optional, Computed | `SecurityGroupRule.BelongMemberId` |
| 17 | `rule_id` | String | Computed | `CreateSecurityGroupRuleResponse.RuleIds[0]`, `SecGroupRuleResp.RuleId` |
| 18 | `source_id` | String | Computed | `SecGroupRuleResp.SourceId` |
| 19 | `target_id` | String | Computed | `SecGroupRuleResp.TargetId` |
| 20 | `target_type` | Int | Computed | `SecGroupRuleResp.TargetType` (int64 in response) |
| 21 | `strategy` | Int | Computed | `SecGroupRuleResp.Strategy` |
| 22 | `detail` | String | Computed | `SecGroupRuleResp.Detail` |
| 23 | `region` | String | Computed | `SecGroupRuleResp.Region` |
| 24 | `sou_instance_name` | String | Computed | `SecGroupRuleResp.SouInstanceName` |
| 25 | `sou_public_ip` | String | Computed | `SecGroupRuleResp.SouPublicIp` |
| 26 | `sou_private_ip` | String | Computed | `SecGroupRuleResp.SouPrivateIp` |
| 27 | `sou_cidr` | String | Computed | `SecGroupRuleResp.SouCidr` |
| 28 | `sou_parameter_name` | String | Computed | `SecGroupRuleResp.SouParameterName` |
| 29 | `instance_name` | String | Computed | `SecGroupRuleResp.InstanceName` |
| 30 | `public_ip` | String | Computed | `SecGroupRuleResp.PublicIp` |
| 31 | `private_ip` | String | Computed | `SecGroupRuleResp.PrivateIp` |
| 32 | `cidr` | String | Computed | `SecGroupRuleResp.Cidr` |
| 33 | `parameter_name` | String | Computed | `SecGroupRuleResp.ParameterName` |
| 34 | `protocol_port_name` | String | Computed | `SecGroupRuleResp.ProtocolPortName` |
| 35 | `is_new` | Int | Computed | `SecGroupRuleResp.IsNew` |
| 36 | `rule_uuid` | Int | Computed | `SecGroupRuleResp.Id` |
| 37 | `dns_parse_count` | List, `MaxItems: 1` | Computed | `SecGroupRuleResp.DnsParseCount` |
| 37.1 | └ `valid_count` | Int | Computed | `SgDnsParseCount.ValidCount` |
| 37.2 | └ `invalid_count` | Int | Computed | `SgDnsParseCount.InvalidCount` |
| 38 | `belong_member` | List, `MaxItems: 1` | Computed | `SecGroupRuleResp.BelongMember` |
| 38.1 | └ `app_id` | String | Computed | `MemberInfo.AppId` |
| 38.2 | └ `uin` | String | Computed | `MemberInfo.Uin` |
| 38.3 | └ `nickname` | String | Computed | `MemberInfo.Nickname` |
| 38.4 | └ `member_id` | String | Computed | `MemberInfo.MemberId` |
| 38.5 | └ `node_name` | String | Computed | `MemberInfo.NodeName` |

Deviations from the raw field dump, all deliberate:

- `response.Response.Rules.Id` → schema **`rule_uuid`**, not `id`. `id` is Terraform's reserved attribute; declaring it inside the resource schema collides with `d.Id()` and confuses the provider/state layout. The API comment itself says `Id` is "规则id 等同 RuleUuid", so `rule_uuid` is the faithful name.
- `response.Response.Rules.SourceType` / `TargetType` are exposed as **`target_type`** (int) plus the request-side **`source_type`/`dest_type`** (string). There is no `source_type` int field: two attributes with the same name but different types are impossible, and the request-side one is the one users author.
- `AllTotalCount` / `TotalCount` are **not** exposed: they count every rule of the group, which is meaningless (and diff-noisy) on a per-rule resource.
- `rule_ids` (the Create response list) is represented by the scalar `rule_id`, since exactly one id is ever returned for a one-element `Rules` slice.

### D6. Read must not fight the API's encoding drift

`ModifySecurityGroupRule` / `CreateSecurityGroupRule` echo nothing back, and `SecGroupRuleResp` encodes several write fields differently. Read therefore **preserves** the config-authored values of these fields instead of overwriting them:

| Write field (string) | Response counterpart | Behaviour in Read |
|---|---|---|
| `source_type` | `SourceType *int64` | **not** overwritten (0/1/2… vs authored string) |
| `dest_type` | `target_type *int64` | **not** overwritten, surfaced separately as `target_type` |
| `rule_action` | `Strategy *int64` | **not** overwritten, surfaced as `strategy` |
| `description` | `Detail *string` (different field) | **not** overwritten, surfaced as `detail` |
| `source_content` | `SourceId *string` | **not** overwritten, surfaced as `source_id` |
| `dest_content` | `TargetId *string` | **not** overwritten, surfaced as `target_id` |
| `protocol_port_type` | `ProtocolPortType *int64` | **not** overwritten (request semantics may differ) |
| `belong_member_id` | `BelongMember *MemberInfo` | **not** overwritten, surfaced as `belong_member` |

Overwritten (types and semantics match 1:1): `group_id` (from id segment), `rule_id`, `ip_version`, `protocol`, `port`, `service_template_id`, `scope`, `order_index`; plus all the read-only attributes (18-38 above).

Every `d.Set` is guarded by a nil check, and nested blocks (`dns_parse_count`, `belong_member`) are built as `[]map[string]interface{}` of length 0/1.

### D7. Service layer API

`tencentcloud/services/fwm/service_tencentcloud_fwm.go`:

```go
func NewFwmService(client *connectivity.TencentCloudClient) FwmService
type FwmService struct { client *connectivity.TencentCloudClient }

// returns (nil, nil) when the rule is gone
func (me *FwmService) DescribeFwmSecurityGroupRuleById(ctx context.Context, groupId, ruleId string) (ret *fwmv20250611.SecGroupRuleResp, errRet error)
```

Implementation notes:

- Request built once: `request.GroupId = helper.String(groupId)`, `request.Filters = []*fwmv20250611.CommonFilter{{Name: helper.String("RuleId"), Values: []*string{helper.String(ruleId)}, OperatorType: helper.Int64(9)}}`, `request.Limit = helper.Int64(100)`, `request.Offset = helper.Int64(0)` (there is **no** `helper.Int64Point` in this repo — `helper.Int64(int64) *int64` is the correct pointer helper, and `helper.IntInt64(int) *int64` converts schema ints).
- `OperatorType = 9` is "模糊匹配" per `CommonFilter`'s comment; we use it verbatim because it is the documented product usage, and we then **strict-equal** `*item.RuleId == ruleId` before returning, so fuzzy matching cannot leak a wrong rule.
- `Limit` is a literal `100` (the struct's comment documents no maximum for `数量`; 100 is the provider-wide default page size). A single-rule-id filter cannot exceed one page, so no pagination loop is needed.
- `ratelimit.Check(request.GetAction())` inside the retry closure, `resource.Retry(tccommon.ReadRetryTimeout, ...)`, `tccommon.RetryError(e)` on error, `NonRetryableError` when `result == nil || result.Response == nil`, and a `[CRITAL]` log in the deferred error handler — exactly the `service_tencentcloud_igtm.go` shape.
- Not found ⇒ `ret = nil, errRet = nil` (never an error), so the resource can distinguish "deleted out of band".

### D8. No `Timeouts` block

`CreateRuleGroup` / `CreateSecurityGroupRule` / `ModifySecurityGroupRule` / `DeleteSecurityGroupRule` return resource ids or bare `RequestId` — no `TaskId`. The rule is effective immediately, so CRUD returns directly and Read is only re-invoked to refresh state (not as a polling oracle). Per the project rule "async operations must declare Timeouts", nothing is required here; adding `Create/Update/Delete` timeouts would be dead configuration.

### D9. File layout & registration

- New dir `tencentcloud/services/fwm/` with `service_tencentcloud_fwm.go`, `resource_tc_fwm_security_group_rule.go`, `resource_tc_fwm_security_group_rule.md`, `resource_tc_fwm_security_group_rule_test.go`. Nothing else (no `_extension.go`).
- `tencentcloud/connectivity/client.go`: add import `fwmv20250611 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"` — gofmt places it immediately **after** `es "…/es/v20180416"` (line 73) and **before** `ga2v20250115 "…/ga2/v20250115"` (line 74) + struct field `fwmV20250611Conn *fwmv20250611.Client` in the cache block (after `ga2v20250115Conn` at line 269, so gofmt's alignment is preserved) + accessor `UseFwmV20250611Client()` next to `UseGa2V20250115Client()` (line 2259; mirroring `NewClientProfile(300)` + `WithHttpTransport(&LogRoundTripper{})`).
- `tencentcloud/provider.go`: add `"tencentcloud_fwm_security_group_rule": fwm.ResourceTencentCloudFwmSecurityGroupRule(),` to `ResourcesMap` (gofmt will realign the map block) plus the `fwm "…/tencentcloud/services/fwm"` import.
- `tencentcloud/provider.md`: append a new section (product not yet listed anywhere):

  ```
  Firewall Manager(FWM)
  Resource
  tencentcloud_fwm_security_group_rule
  ```

- `website/docs/r/fwm_security_group_rule.html.markdown` is generated by `make doc` — never hand-edited.

### D10. Unit tests use gomonkey, not the acceptance suite

Per the project rule for *new* resources: `resource_tc_fwm_security_group_rule_test.go` (package `fwm_test`) fakes the cloud API instead of running `resource.Test`. Shape (copy `resource_tc_cfw_cluster_fw_bypass_config_test.go`):

- `type mockMetaFwm struct { client *connectivity.TencentCloudClient }` with `GetAPIV3Conn() *connectivity.TencentCloudClient` and `var _ tccommon.ProviderMeta = &mockMetaFwm{}`.
- `patches.ApplyMethodReturn(meta.client, "UseFwmV20250611Client", fwmClient)` where `fwmClient := &fwmv20250611.Client{}`.
- `patches.ApplyMethodFunc(fwmClient, "CreateRuleGroupWithContext", ...)`, `"CreateSecurityGroupRuleWithContext"`, `"DescribeSecurityGroupRulesWithContext"`, `"ModifySecurityGroupRuleWithContext"`, `"DeleteSecurityGroupRuleWithContext"`, `"DeleteRuleGroupWithContext"` to return canned responses.
- Cases: create with `group_name` (asserts `d.Id() == groupId#ruleId`, `rule_id` set, group created), create with `group_id` (asserts `CreateRuleGroup` **not** called), read normal (asserts every field incl. nested `dns_parse_count` / `belong_member`), read when describe returns empty (asserts `d.Id() == ""`), update (asserts `Rule.RuleId` populated), delete with `group_name` (asserts both deletes) and with `group_id` (asserts only the rule delete).
- No `go test` / `go build` execution in this workflow; tests must be compile-clean by construction.

## Risks / Trade-offs

- **[Risk]** `SourceType` / `DestType` (string, write) vs `SourceType` / `TargetType` (int64, read) cannot be reconciled, so Read leaves the authored values alone. A genuine out-of-band change of the *type* of a rule will not be detected by `terraform plan`.
  → **Mitigation**: the int-encoded values are still observable through the computed `target_type` / `strategy` / `detail` / `source_id` / `target_id` attributes; documented in the resource schema descriptions. Drift in `protocol` / `port` / `scope` / `ip_version` *is* detected.
- **[Risk]** `OperatorType = 9` means fuzzy match, so a differently-prefixed rule id could in theory be returned.
  → **Mitigation**: the service helper strict-equals `RuleId` and returns `(nil, nil)` otherwise; the resource then logs `[CRUD] tencentcloud_fwm_security_group_rule id=...` and clears the id.
- **[Risk]** Destroying a rule also tries to destroy the auto-created rule group; if the group was already deleted out of band, `DeleteRuleGroup` fails and the destroy aborts.
  → **Mitigation**: the group delete is wrapped in `resource.Retry(tccommon.WriteRetryTimeout, ...)` and its error is logged as `[WARN]` rather than failing the destroy when the failure is "not found" style (the outcome Terraform wants — the group is gone — is already achieved). A hard failure on other codes still surfaces.
- **[Risk]** `group_id` is `Optional + Computed + ForceNew`; a user who supplies `group_name` and later wants to move the rule to another group gets a destroy/recreate.
  → **Mitigation**: documented in the resource `.md`; moving a rule between groups is not expressible in `ModifySecurityGroupRule` anyway.
- **[Trade-off]** One rule per resource means N resources for N rules (more Terraform blocks + N API calls instead of one batch create).
  → Accepted: it is the only shape that gives each rule a stable Terraform identity, and it matches the per-`RuleId` describe/delete/modify contracts.
- **[Trade-off]** `AllTotalCount` / `TotalCount` are not surfaced, so users cannot read the group's rule count from this resource. A future data source can cover that.

## Migration Plan

Purely additive; no state migration.

1. Add `tencentcloud/connectivity/client.go` client field + `UseFwmV20250611Client()`.
2. Add `tencentcloud/services/fwm/**`.
3. Register in `tencentcloud/provider.go` and `tencentcloud/provider.md`.
4. `gofmt` (finalize stage), `make doc` (finalize stage) to generate `website/docs/r/fwm_security_group_rule.html.markdown`.
5. Users opt in by declaring `resource "tencentcloud_fwm_security_group_rule" "x" { ... }`.

Rollback: revert the new directory + the two registration spots + the connectivity accessor line. Any group/rule already created in the cloud can be removed with `terraform destroy` before reverting.

## Open Questions

None blocking. The SDK surface is fully vendored; the remaining choices (single-rule modelling, 2-segment id, group ownership, read-drift policy) are decided above and reflected in `specs/fwm-security-group-rule-resource/spec.md`.
