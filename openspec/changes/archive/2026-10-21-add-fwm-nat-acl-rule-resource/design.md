# Design: tencentcloud_fwm_nat_acl_rule Resource

## Architecture

Follows `tencentcloud_igtm_strategy` style (single-file CRUD resource + service layer). fwm service is brand-new in the provider, so it also needs the connectivity client binding and provider registration.

```
provider.go
    └─ import services/fwm
    └─ ResourcesMap["tencentcloud_fwm_nat_acl_rule"] = fwm.ResourceTencentCloudFwmNatAclRule()
           └─ resource_tc_fwm_nat_acl_rule.go (CRUD handlers + schema)
                  └─ service_tencentcloud_fwm.go (DescribeFwmNatAclRuleById)
                         └─ fwm SDK v20250611 (CreateNatAclRule / DescribeNatAclRules / ModifyNatAclRule / DeleteNatAclRule)
```

## File Layout

| File | Action |
|---|---|
| `tencentcloud/connectivity/client.go` | Modified — add fwm SDK import, `fwmv20250611Conn` field, `UseFwmV20250611Client()` method |
| `tencentcloud/services/fwm/resource_tc_fwm_nat_acl_rule.go` | New — CRUD + schema |
| `tencentcloud/services/fwm/resource_tc_fwm_nat_acl_rule.md` | New — resource doc + import example |
| `tencentcloud/services/fwm/resource_tc_fwm_nat_acl_rule_test.go` | New — unit tests via gomonkey mock |
| `tencentcloud/services/fwm/service_tencentcloud_fwm.go` | New — FwmService + DescribeFwmNatAclRuleById |
| `tencentcloud/provider.go` | Modified — register `tencentcloud_fwm_nat_acl_rule` |
| `tencentcloud/provider.md` | Modified — add resource doc entry |

## Connectivity Client Binding

fwm SDK is already vendored at `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611/`. The provider currently has NO fwm client binding. Mirror the existing CFW binding pattern (`client.go:2055` `UseCfwV20190904Client`):

```go
// import (client.go ~line 46 area, alphabetical)
fwmv20250611 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"

// conn field (client.go ~line 251 area, near cfwv20190904Conn)
fwmv20250611Conn            *fwmv20250611.Client

// method (append near UseCfwV20190904Client)
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

## Cloud API SDK Reference (vendored fwm/v20250611/models.go)

All four interfaces are **synchronous** (no FlowId/TaskId/JobId, no polling). The SDK client methods `CreateNatAclRule`, `DescribeNatAclRules`, `ModifyNatAclRule`, `DeleteNatAclRule` each also have a `...WithContext` variant. Confirmed from `client.go:302` (`func (c *Client) CreateNatAclRule(...) -> c.CreateNatAclRuleWithContext(context.Background(), request)`).

### CommonFilter (models.go:321) — used for Describe filters
```go
type CommonFilter struct {
    Name        *string   // filter field name, e.g. "RuleId"
    Values      []*string // filter values
    OperatorType *int64   // 1=EQUAL, 7=in, 9=fuzzy
}
```

### CreateNatAclRuleRequest (models.go:614)
```go
type CreateNatAclRuleRequest struct {
    GroupId *string          // required
    Rules   []*NatAclRule    // required; this resource passes exactly 1 element
}
```

### CreateNatAclRuleResponseParams (models.go:645)
```go
type CreateNatAclRuleResponseParams struct {
    RuleIds  []*string  // server-generated rule ids, ordered to match input Rules
    RequestId *string
}
```
Bind RuleId = `RuleIds[0]`.

### NatAclRule (models.go:3718) — input struct for Create & Modify
```go
type NatAclRule struct {
    SourceContent   *string
    SourceType      *string  // ip/net/url/domain/template/addressbook/instance/tag/location/vendor/dnsparse
    TargetContent   *string
    TargetType      *string
    Protocol        *string  // TCP/UDP/ICMP/ANY/HTTP/HTTPS/DNS/FTP...
    RuleAction      *string  // accept/drop/log (lowercase)
    OrderIndex      *int64   // >=1, per-direction numbering
    Scope           *string  // ""|ALL|region(id like ap-guangzhou)|nat instance id (cfwnat-xxx / nat-xxx)
    Direction        *int64   // 0=outbound, 1=inbound
    RuleId          *string  // required on modify
    Port            *string  // e.g. 80 / 80,443 / 80-90 ; ANY = -1/-1
    Description     *string
    ParamTemplateId *string
    BelongMemberId  *string  // required when source/target type in instance/tag/template/addressbook
}
```

### DescribeNatAclRulesRequest (models.go:1569)
```go
type DescribeNatAclRulesRequest struct {
    GroupId   *string          // required
    Direction *int64           // required: 0=outbound, 1=inbound  (NOT optional at API level)
    Filters   []*CommonFilter  // use {Name:"RuleId",Values:[ruleId],OperatorType:1} for single-rule lookup
    Offset    *int64           // default 0; negative errors
    Limit     *int64           // default 100; negative errors. Use 100 (API max) for single-rule read.
    Order     *string          // asc|desc
    By        *string
}
```

### DescribeNatAclRulesResponseParams (models.go:1620)
```go
type DescribeNatAclRulesResponseParams struct {
    Rules        []*NatAclRuleResp  // list; take Rules[0] for this single-instance resource
    TotalCount   *int64
    AllTotalCount *int64
    RequestId    *string
}
```

### NatAclRuleResp (models.go:3767) — output struct
```go
type NatAclRuleResp struct {
    RuleId            *string
    Sequence          *int64   // <-> input OrderIndex (name mismatch!)
    Direction         *int64
    SourceContent    *string
    SourceType        *string
    SourceName        *string  // display only, computed, not compared
    TargetContent     *string
    TargetType        *string
    TargetName        *string  // display only
    Port              *string
    Protocol          *string
    RuleAction        *string
    Description       *string
    Scope             *string
    ScopeDesc          *string  // display only (language-dependent), computed, NOT in state compare
    FwInsId           *string  // display/computed
    CountryName       *string
    CityName          *string
    CountryCode       *int64
    CityCode          *int64
    CountryKey        *string
    CityKey           *string
    IsRegion          *int64   // 0|1
    CloudCode         *string
    IsCloud           *int64   // 0|1
    ParamTemplateId   *string
    ParamTemplateName  *string  // display only
    Invalid            *int64   // 0=valid,1=invalid
    BelongMember       *MemberInfo // <- input BelongMemberId is string
    CreateTime         *string
    UpdateTime         *string
    DnsParseCnt        *int64
    HitCount           *int64
    LastHitTime         *string
}
```

### MemberInfo (models.go:3049) — BelongMember
```go
type MemberInfo struct {
    AppId    *string
    Uin      *string
    Nickname *string
    MemberId *string  // <-> input BelongMemberId (type mismatch!)
    NodeName *string
}
```

### ModifyNatAclRuleRequest (models.go:3204)
```go
type ModifyNatAclRuleRequest struct {
    GroupId *string      // required (must match original, not changeable)
    Rule    *NatAclRule  // required; full overwrite, include RuleId + all writable fields
}
```
Response = `{ RequestId }` only.

### DeleteNatAclRuleRequest (models.go:1158)
```go
type DeleteNatAclRuleRequest struct {
    GroupId *string     // required
    RuleIds []*string   // required; pass 1 element for this resource
}
```
Response = `{ RequestId }` only. Idempotent: if RuleId no longer exists, API returns `ResourceNotFound` → treat delete as success.

## Schema

### Required
| Field | Type | ForceNew | Description |
|---|---|---|---|
| `group_id` | String | yes | Rule group ID (Parent resource ref; all 4 APIs require it) |
| `direction` | Int | yes | 0=outbound, 1=inbound (required by Describe, ForceNew) |
| `source_content` | String | no | Source address content |
| `source_type` | String | yes | Source type: ip/net/url/domain/template/addressbook/instance/tag/location/vendor/dnsparse |
| `target_content` | String | no | Target address content |
| `target_type` | String | yes | Target type (same enum as source_type) |
| `protocol` | String | no | Protocol: TCP/UDP/ICMP/ANY/HTTP/HTTPS/DNS/FTP... |
| `rule_action` | String | no | Action: accept/drop/log (lowercase) |
| `order_index` | Int | no | Priority, >=1, per-direction numbering (mutable) |
| `scope` | String | no | Effective range: ALL/region/nat-instance-id |

### Optional
| Field | Type | ForceNew | Description |
|---|---|---|---|
| `port` | String | no | Port e.g. 80, 80,443, 80-90; ANY = -1/-1 |
| `description` | String | no | Rule description |
| `param_template_id` | String | no | Port template ID |
| `belong_member_id` | String | no | Member account ID (required when source/target type in instance/tag/template/addressbook) |

### Computed (read-only, only filled by Describe)
| Field | Type | Source (NatAclRuleResp) | Note |
|---|---|---|---|
| `rule_id` | String | RuleId | composite id part |
| `sequence` | Int | Sequence | display mirror of order_index |
| `scope_desc` | String | ScopeDesc | language-dependent display, do not compare |
| `create_time` | String | CreateTime | |
| `update_time` | String | UpdateTime | |
| `source_name` | String | SourceName | display only |
| `target_name` | String | TargetName | display only |
| `param_template_name` | String | ParamTemplateName | display only |
| `fw_ins_id` | String | FwInsId | |
| `country_name` | String | CountryName | |
| `city_name` | String | CityName | |
| `country_code` | Int | CountryCode | |
| `city_code` | Int | CityCode | |
| `country_key` | String | CountryKey | |
| `city_key` | String | CityKey | |
| `is_region` | Int | IsRegion | |
| `cloud_code` | String | CloudCode | |
| `is_cloud` | Int | IsCloud | |
| `invalid` | Int | Invalid | 0=valid,1=invalid |
| `app_id` | String | BelongMember.AppId | |
| `uin` | String | BelongMember.Uin | |
| `nickname` | String | BelongMember.Nickname | |
| `member_id` | String | BelongMember.MemberId | |
| `node_name` | String | BelongMember.NodeName | |
| `dns_parse_cnt` | Int | DnsParseCnt | |
| `hit_count` | Int | HitCount | |
| `last_hit_time` | String | LastHitTime | |

## Composite ID

`d.SetId(strings.Join([]string{groupId, ruleId}, tccommon.FILED_SP))` where `FILED_SP = "#"` (common.go:32).
On read/update/delete, split by `tccommon.FILED_SP` → `idSplit[0]=groupId, idSplit[1]=ruleId`.

Import: user passes `GroupId:RuleId` (colon). Use a custom `Importer.StateContext` that replaces `:` with `#` then calls `schema.ImportStatePassthroughContext`, OR document the `GroupId#RuleId` form. Per requirements, import supports `GroupId:RuleId` two-segment; the state file stores `#`-separated.

## Field Name Mismatch Mapping (critical for plan-stable)

| Input (Create/Modify) | Output (Describe) | Mapping rule |
|---|---|---|
| `NatAclRule.OrderIndex` (int64) | `NatAclRuleResp.Sequence` (int64) | On read: `d.Set("order_index", Sequence)`; on write: `rule.OrderIndex = d.Get("order_index")` |
| `NatAclRule.BelongMemberId` (string) | `NatAclRuleResp.BelongMember` (MemberInfo) | On read: if `BelongMember != nil && BelongMember.MemberId != nil` then `d.Set("belong_member_id", *BelongMember.MemberId)`; on write: `rule.BelongMemberId = d.Get("belong_member_id")` (only if non-empty) |

## CRUD Strategy (all synchronous, NO polling)

### Create (`resourceTencentCloudFwmNatAclRuleCreate`)
1. Build `CreateNatAclRuleRequest{ GroupId, Rules: []*NatAclRule{ singleRule } }`.
2. `singleRule` populated from schema (SourceContent, SourceType, TargetContent, TargetType, Protocol, RuleAction, OrderIndex, Scope, Direction, Port, Description, ParamTemplateId, BelongMemberId). Use `d.GetOkExists` for `direction` (int). For `order_index` use `d.GetOkExists` to avoid 0 treated unset. Required string fields use `d.GetOk`.
3. `resource.Retry(tccommon.WriteRetryTimeout, ...)` calling `UseFwmV20250611Client().CreateNatAclRuleWithContext(ctx, request)`; on error `tccommon.RetryError(e)`.
4. After success: check `response == nil || response.Response == nil || len(response.Response.RuleIds) == 0` → `NonRetryableError`. Print `logId` + `d.Id()` before check.
5. `ruleId = *response.Response.RuleIds[0]`; `d.SetId(strings.Join([]string{groupId, ruleId}, tccommon.FILED_SP))`. SetId BEFORE returning (sync, no async wait).
6. Return `resourceTencentCloudFwmNatAclRuleRead(d, meta)`.

### Read (`resourceTencentCloudFwmNatAclRuleRead`)
1. Split id → groupId, ruleId. Check `len(idSplit) != 2` error.
2. Call `service.DescribeFwmNatAclRuleById(ctx, groupId, ruleId, direction)`. **The service must pass Direction.**
3. If err is `ResourceNotFound` → `log.Printf("[CRUD] fwm_nat_acl_rule id=%s", d.Id())` then `d.SetId("")` return nil (rule externally deleted).
4. If response empty (Rules len 0) → same logging + `d.SetId("")` (NOT in retry for non-datasource; this is resource read).
5. Fill writable + computed fields from `Rules[0]`. Apply mismatch mapping (Sequence→order_index, BelongMember.MemberId→belong_member_id). Guard all nil checks before `d.Set`.
6. Direction: read from `Rules[0].Direction`.

### Update (`resourceTencentCloudFwmNatAclRuleUpdate`)
1. If `d.HasChange("source_type")`, `d.HasChange("target_type")`, or `d.HasChange("direction")` → these are ForceNew, Terraform handles rebuild (won't reach update). (Defensive: no action.)
2. Build `ModifyNatAclRuleRequest{ GroupId, Rule: &NatAclRule{...} }` with **all writable fields** (full overwrite). Rule.RuleId = ruleId from id. Rule.Direction from schema.
3. `resource.Retry(tccommon.WriteRetryTimeout, ...)` calling `ModifyNatAclRuleWithContext`. On error `RetryError`.
4. Check `response == nil || response.Response == nil` → `NonRetryableError`.
5. Return read.

### Delete (`resourceTencentCloudFwmNatAclRuleDelete`)
1. Split id → groupId, ruleId.
2. `request = DeleteNatAclRuleRequest{ GroupId, RuleIds: []*string{ &ruleId } }`.
3. `resource.Retry(tccommon.WriteRetryTimeout, ...)`: if error and `tccommon.IsExpectError(e, []string{"ResourceNotFound"})` → return nil (idempotent success). Else `RetryError(e)`.

## Service Layer (`service_tencentcloud_fwm.go`)

Mirror `service_tencentcloud_cfw.go` structure:
```go
package fwm

type FwmService struct {
    client *connectivity.TencentCloudClient
}

func (me *FwmService) DescribeFwmNatAclRuleById(ctx context.Context, groupId, ruleId string, direction int64) (rule *fwmv20250611.NatAclRuleResp, errRet error) {
    // request: GroupId, Direction, Limit=100, Offset=0
    // Filters: {Name:"RuleId", Values:[ruleId], OperatorType:1}(EQUAL)
    // ratelimit.Check + UseFwmV20250611Client().DescribeNatAclRules
    // if len(Rules) == 0 -> return nil, nil (caller decides SetId)
    // return Rules[0]
}
```

## ResourceNotFound / empty handling summary

| Operation | Condition | Action |
|---|---|---|
| Read | err contains `ResourceNotFound` | log + `d.SetId("")`, return nil |
| Read | Rules empty (len 0) | log + `d.SetId("")`, return nil |
| Read | other err | return err |
| Delete | err contains `ResourceNotFound` | return nil (idempotent) |

## Provider Registration (provider.go)

- Add import: `"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/fwm"` (alphabetical, near line 84 `services/igtm`).
- Add to `ResourcesMap`: `"tencentcloud_fwm_nat_acl_rule": fwm.ResourceTencentCloudFwmNatAclRule()` (near line 2783 `tencentcloud_igtm_strategy`).

## Test Strategy (gomonkey mock, NOT terraform suite)

Follow `resource_tc_bh_bind_device_resource_test.go` pattern:
- Test package `fwm_test`, mock `ProviderMeta`, mock `UseFwmV20250611Client`.
- Mock `CreateNatAclRuleWithContext` / `DescribeNatAclRulesWithContext` / `ModifyNatAclRuleWithContext` / `DeleteNatAclRuleWithContext`.
- Tests: TestCreate (assert request fields, return RuleIds[0]=fake-id), TestRead (return Rules[0]), TestUpdate (assert full Rule sent), TestDelete (success + ResourceNotFound idempotent).