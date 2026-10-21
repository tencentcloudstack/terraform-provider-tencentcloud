## 1. Connectivity Client Binding

- [x] 1.1 In `tencentcloud/connectivity/client.go`: add import `fwmv20250611 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"`
- [x] 1.2 Add conn field `fwmv20250611Conn *fwmv20250611.Client` (near `cfwv20190904Conn` ~line 251)
- [x] 1.3 Add method `UseFwmV20250611Client()` mirroring `UseCfwV20190904Client()` (line 2055): NewClientProfile(300), WithHttpTransport(LogRoundTripper)

## 2. Service Layer

- [x] 2.1 Create `tencentcloud/services/fwm/service_tencentcloud_fwm.go` with `FwmService` struct
- [x] 2.2 Implement `DescribeFwmNatAclRuleById(ctx, groupId, ruleId, direction)` — calls `DescribeNatAclRules` with `GroupId`+`Direction`+`Filters[{Name:"RuleId",Values:[ruleId],OperatorType:1}]`+`Limit=100`+`Offset=0`; returns `Rules[0]` or nil if empty

## 3. Resource Implementation

- [x] 3.1 Create `tencentcloud/services/fwm/resource_tc_fwm_nat_acl_rule.go` with full schema (required/optional/computed per design.md). Resource description in Chinese, field descriptions in English.
- [x] 3.2 Implement `resourceTencentCloudFwmNatAclRuleCreate`: build single-element `Rules`, `UseFwmV20250611Client().CreateNatAclRuleWithContext` in `WriteRetryTimeout` retry; guard nil RuleIds; SetId = `GroupId#RuleId`; return read
- [x] 3.3 Implement `resourceTencentCloudFwmNatAclRuleRead`: split id → groupId/ruleId; call `DescribeFwmNatAclRuleById`; ResourceNotFound/empty → log + `d.SetId("")`; fill fields with mismatch mapping (Sequence→order_index, BelongMember.MemberId→belong_member_id); nil-guard all sets
- [x] 3.4 Implement `resourceTencentCloudFwmNatAclRuleUpdate`: build `ModifyNatAclRuleRequest` with full Rule (all writable fields + RuleId + Direction); `WriteRetryTimeout` retry; nil-guard; return read
- [x] 3.5 Implement `resourceTencentCloudFwmNatAclRuleDelete`: `DeleteNatAclRuleRequest{ GroupId, RuleIds:[ruleId] }`; `WriteRetryTimeout` retry; `ResourceNotFound` → idempotent nil

## 4. Provider Registration

- [x] 4.1 In `tencentcloud/provider.go`: add import `services/fwm`; register `"tencentcloud_fwm_nat_acl_rule": fwm.ResourceTencentCloudFwmNatAclRule()` in `ResourcesMap`
- [x] 4.2 In `tencentcloud/provider.md`: add resource doc entry for `tencentcloud_fwm_nat_acl_rule`

## 5. Documentation & Tests

- [x] 5.1 Create `tencentcloud/services/fwm/resource_tc_fwm_nat_acl_rule.md` (one-line desc with FWM product name, Example Usage, Import section showing `GroupId#RuleId`)
- [x] 5.2 Create `tencentcloud/services/fwm/resource_tc_fwm_nat_acl_rule_test.go` (gomonkey mock cloud APIs; TestCreate/TestRead/TestUpdate/TestDelete including Delete ResourceNotFound idempotent case)

## 6. Verification

- [x] 6.1 Confirm code correctness: all Create params exist in CreateNatAclRule; all Modify params exist in ModifyNatAclRule; Describe reads correct struct; Delete uses correct field. Ensure no go build needed (gofmt only in finalize).