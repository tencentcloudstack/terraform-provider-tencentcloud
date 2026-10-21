## 1. SDK / Dependency Preconditions

- [x] 1.1 Confirm `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm v1.3.196` is present in `go.mod`, `tools.go` (`_ ".../fwm/v20250611"`) and `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611/`; do not run `go mod tidy`/`go build`.

## 2. Connectivity Client Registration

- [x] 2.1 In `tencentcloud/connectivity/client.go`, add import alias `fwmv20250611 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"` next to the other service imports.
- [x] 2.2 Add cached field `fwmv20250611Conn *fwmv20250611.Client` to the `TencentCloudClient` struct (near `igtmv20231024Conn`).
- [x] 2.3 Add `UseFwmV20250611Client() *fwmv20250611.Client` following `UseIgtmV20231024Client` (lazy init, `me.NewClientProfile(300)`, `WithHttpTransport(&LogRoundTripper{})`).

## 3. Service Layer

- [x] 3.1 Create `tencentcloud/services/fwm/service_tencentcloud_fwm.go` (package `fwm`) with `NewFwmService(client *connectivity.TencentCloudClient) FwmService` and `type FwmService struct { client *connectivity.TencentCloudClient }`.
- [x] 3.2 Implement `DescribeFwmRuleGroupsByFilter(ctx context.Context, param map[string]interface{}) (ret []*fwmv20250611.RuleGroupResp, errRet error)`: accepts `"Product"`, `"Filters"`, `"Order"`, `"By"` from `param`; internally paginates with `Offset` from `0` and `Limit = 100`; calls `DescribeRuleGroupsWithContext` inside `resource.Retry(tccommon.ReadRetryTimeout, ...)` with `ratelimit.Check(request.GetAction())`; nil-guards `result`/`result.Response`/`result.Response.RuleGroups`; appends until page length `< 100`.
- [x] 3.3 Implement `DescribeFwmRuleGroupById(ctx context.Context, groupId string) (ret *fwmv20250611.RuleGroupResp, errRet error)`: builds `Filters: []*fwmv20250611.CommonFilter{{Name: "GroupId", Values: []*string{group}, OperatorType: 9}}` (use `helper.String` / `helper.Int64`), delegates to `DescribeFwmRuleGroupsByFilter`, and returns the element whose `GroupId` exactly equals `groupId`; returns `(nil, nil)` when nothing matches.
- [x] 3.4 Use the standard `defer func() { if errRet != nil { log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error()) } }()` logging pattern in both helpers.

## 4. Resource Schema

- [x] 4.1 Create `tencentcloud/services/fwm/resource_tc_fwm_rule_group.go` (package `fwm`) with `ResourceTencentCloudFwmRuleGroup() *schema.Resource` wiring `Create`/`Read`/`Update`/`Delete` and `Importer: &schema.ResourceImporter{State: schema.ImportStatePassthrough}`.
- [x] 4.2 Declare user attributes: `group_name` (`TypeString`, Required) and `product` (`TypeString`, Required, ForceNew, `ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"enterprise_sg", "cfw_nat_acl", "cfw_vpc_acl", "cfw_edge_acl"})`).
- [x] 4.3 Declare computed attributes: `group_id` (string), `app_id` (string), `rule_count` (int), `create_time` (string), `update_time` (string), `create_by` (string), `create_name` (string), `update_by` (string), `update_name` (string), `exec_areas` (`TypeList` of `TypeString`), `direction` (int).
- [x] 4.4 Declare computed `receive_account` as `TypeList` whose `Elem` is a `*schema.Resource` with `uin` (string), `nickname` (string), `receiver_type` (int) and a nested computed `members` `TypeList` of `*schema.Resource` with `app_id`, `uin`, `nickname`, `member_id`, `node_name` (all string). Every element field must be `Computed: true`.
- [x] 4.5 Do not add a `rule_groups` wrapper attribute, and do not add `filters`/`order`/`by`/`offset`/`limit`/`group_ids` attributes.

## 5. Resource CRUD

- [x] 5.1 `Create`: `defer tccommon.LogElapsed("resource.tencentcloud_fwm_rule_group.create")()` + `defer tccommon.InconsistentCheck(d, meta)()`; build `CreateRuleGroupRequest` from `group_name`/`product`; call `CreateRuleGroupWithContext` inside `resource.Retry(tccommon.WriteRetryTimeout, ...)`; nil-guard `result`/`result.Response`; require a non-nil non-empty `GroupId` else `resource.NonRetryableError`; then `d.SetId(*response.Response.GroupId)`; then run the visibility poll of 5.2; then call Read.
- [x] 5.2 Visibility poll inside `Create`: `resource.Retry(tccommon.ReadRetryTimeout, ...)` around `service.DescribeFwmRuleGroupById(ctx, groupId)`, returning `tccommon.RetryError(e)` on error and `resource.RetryableError(...)` while the returned struct is nil.
- [x] 5.3 `Read`: `defer tccommon.LogElapsed("resource.tencentcloud_fwm_rule_group.read")()`; use `d.Id()` directly as the group id; call `service.DescribeFwmRuleGroupById`; if nil, `log.Printf("[CRUD] tencentcloud_fwm_rule_group id=%s", d.Id())` and then `d.SetId("")` / `return nil`; otherwise `_ = d.Set(...)` for each present field, each guarded by a nil check, with `product` going through `fwmRuleGroupProductCodeToName`.
- [x] 5.4 Add unexported helper `fwmRuleGroupProductCodeToName(v string) string` mapping `"1"→"enterprise_sg"`, `"2"→"cfw_nat_acl"`, `"3"→"cfw_vpc_acl"`, `"4"→"cfw_edge_acl"`, default echo input.
- [x] 5.5 Implement flattening of `receive_account` (`[]*fwmv20250611.ReceiveAccount` → `[]map[string]interface{}`, including the nested `members` list) as unexported helper(s) used by Read.
- [x] 5.6 `Update`: `defer tccommon.LogElapsed("resource.tencentcloud_fwm_rule_group.update")()`; only act when `d.HasChange("group_name")`; build `ModifyRuleGroupRequest{GroupId: &id, GroupName: ...}`; call `ModifyRuleGroupWithContext` inside `resource.Retry(tccommon.WriteRetryTimeout, ...)`; nil-guard `result`/`result.Response`; then call Read.
- [x] 5.7 `Delete`: `defer tccommon.LogElapsed("resource.tencentcloud_fwm_rule_group.delete")()`; build `DeleteRuleGroupRequest{GroupIds: []*string{&id}}` with exactly one element from `d.Id()`; call `DeleteRuleGroupWithContext` inside `resource.Retry(tccommon.WriteRetryTimeout, ...)`; nil-guard `result`/`result.Response`; log `[CRITAL]` on failure.
- [x] 5.8 Do not create any `_extension.go` file for this resource.

## 6. Provider Registration

- [x] 6.1 Add import `"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/fwm"` to `tencentcloud/provider.go`.
- [x] 6.2 Register `"tencentcloud_fwm_rule_group": fwm.ResourceTencentCloudFwmRuleGroup(),` in the `ResourcesMap` map.
- [x] 6.3 Add a product section to `tencentcloud/provider.md` following the existing product-grouped list format: heading `Tencent Cloud Firewall Manager(FWM)`, then `Resource`, then `tencentcloud_fwm_rule_group`.

## 7. Documentation

- [x] 7.1 Create `tencentcloud/services/fwm/resource_tc_fwm_rule_group.md` with: one-line description `Provides a resource to create a FWM rule group`, an `Example Usage` HCL block using `group_name` and `product = "enterprise_sg"`, and an `Import` section showing `terraform import tencentcloud_fwm_rule_group.example fwmrg-wsrrhuvrib`.
- [x] 7.2 Do not add `Argument Reference` / `Attribute Reference` sections; do not hand-write anything under `website/` (`make doc` runs in the finalize stage only).

## 8. Unit Tests (gomonkey, no real cloud calls)

- [x] 8.1 Create `tencentcloud/services/fwm/resource_tc_fwm_rule_group_test.go` in package `fwm_test`, with a local mock meta type implementing `tccommon.ProviderMeta` (`GetAPIV3Conn() *connectivity.TencentCloudClient`, backed by `&connectivity.TencentCloudClient{}`).
- [x] 8.2 Add patch helpers using `gomonkey.ApplyMethodReturn(client, "UseFwmV20250611Client", fwmClient)` and `gomonkey.ApplyMethodFunc(fwmClient, "CreateRuleGroupWithContext"/"DescribeRuleGroupsWithContext"/"ModifyRuleGroupWithContext"/"DeleteRuleGroupWithContext", ...)`; build `*schema.ResourceData` via `schema.TestResourceDataRaw`.
- [x] 8.3 Cover: Create success (asserts `d.Id()`, `group_name`, normalized `product`), Create with nil/empty `GroupId` (expects error), Read success (asserts scalar fields plus `receive_account.0.members.0.*`, `exec_areas`, `direction`), Read product code `2` → `cfw_nat_acl`, Read not-found (asserts `d.Id() == ""`), Update `group_name` (asserts request carries `GroupId` + new `GroupName`), Delete (asserts `GroupIds` length is 1), and a schema-flags test for Required/ForceNew/Computed.
- [x] 8.4 Ensure all returned errors are checked or explicitly discarded with `_ =`, so the package compiles cleanly.

## 9. Verification (do not execute build/test commands)

- [x] 9.1 Re-read the new resource and service files and verify every `fwmv20250611` symbol used exists in `vendor/.../fwm/v20250611/models.go` (`CreateRuleGroupRequest`, `CreateRuleGroupResponse`, `DescribeRuleGroupsRequest`, `DescribeRuleGroupsResponse`, `RuleGroupResp`, `CommonFilter`, `ReceiveAccount`, `MemberInfo`, `ModifyRuleGroupRequest`, `ModifyRuleGroupResponse`, `DeleteRuleGroupRequest`, `DeleteRuleGroupResponse`).
- [x] 9.2 Verify no `DescribeRuleGroup` (singular) API is referenced anywhere, since it does not exist in this SDK version.
- [x] 9.3 Verify `resource_tc_fwm_rule_group.go`, `service_tencentcloud_fwm.go`, the test file, `connectivity/client.go` and `provider.go` are mutually consistent (function names, receiver names, package import path `tencentcloud/services/fwm`).
- [x] 9.4 Leave `gofmt`, `make doc` and `.changelog/` generation to the finalize stage.
