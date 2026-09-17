# Tasks for `add-bdrc-disaster-recovery-protect-group-resource`

## 1. Connectivity layer (BDRC client)

- [x] 1.1 In `tencentcloud/connectivity/client.go`, add import `bdrcv20260330 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/bdrc/v20260330"` alongside the other sdk version imports (keep alphabetical-ish adjacency with other service imports).
- [x] 1.2 Add struct field `bdrcv20260330Conn *bdrcv20260330.Client` to `TencentCloudClient` (mirroring `igtmv20231024Conn`).
- [x] 1.3 Add method `func (me *TencentCloudClient) UseBdrcV20260330Client() *bdrcv20260330.Client` mirroring `UseIgtmV20231024Client`: lazy-init with `me.NewClientProfile(300)`, `bdrcv20260330.NewClient(me.Credential, me.Region, cpf)`, `.WithHttpTransport(&LogRoundTripper{})`, return cached instance.

## 2. Service layer

- [x] 2.1 Create `tencentcloud/services/bdrc/service_tencentcloud_bdrc.go`, package `bdrc`. Define `type BdrcService struct { client *tccommon.ProviderMeta }` (use the same struct pattern as `igtm/service_tencentcloud_igtm.go`: `client *tccommon.ProviderMeta` or the connectivity client — match igtm's actual field type).
- [x] 2.2 Implement `func (s *BdrcService) DescribeDisasterRecoveryProtectGroupById(ctx context.Context, protectGroupId, protectGroupType string) (*bdrcv20260330.ProtectGroup, error)`:
  - Build `DescribeDisasterRecoveryProtectGroupsRequest`; set `ProtectGroupIds = []*string{&protectGroupId}`.
  - Set `ProtectGroupType` only when `protectGroupType != ""`.
  - Set `Limit = helper.Int64(100)` (API max); loop `Offset` from 0 incrementing by `Limit` until `len(ProtectGroupSet) < Limit` or match found.
  - Wrap `UseBdrcV20260330Client().DescribeDisasterRecoveryProtectGroupsWithContext(ctx, request)` in `resource.Retry(tccommon.ReadRetryTimeout, ...)`; guard `result == nil || result.Response == nil` → `resource.NonRetryableError`.
  - Iterate `result.Response.ProtectGroupSet`; if `*item.ProtectGroupId == protectGroupId` return `item, nil`.
  - After exhausting pages return `nil, nil`.
- [x] 2.3 Add the `DescribeDisasterRecoveryProtectGroups` total-count pagination loop so all pages are scanned (not just the first page) when the ID is not in the first batch.

## 3. Resource implementation

- [x] 3.1 Create `tencentcloud/services/bdrc/resource_tc_bdrc_disaster_recovery_protect_group.go`, package `bdrc`. Imports: `context`, `fmt`, `log`, terraform plugin sdk v2 helpers (`resource`, `schema`), `bdrcv20260330` SDK, `tccommon`, `helper`. No `strings` import (single-field ID).
- [x] 3.2 Declare `func ResourceTencentCloudBdrcDisasterRecoveryProtectGroup() *schema.Resource` with `Create/Read/Update/Delete` callbacks and `Importer: &schema.ResourceImporter{State: schema.ImportStatePassthrough}`.
- [x] 3.3 Schema input fields (match design D1):
  - `site_pair_id` (TypeString, Required, ForceNew)
  - `protect_group_type` (TypeString, Required, ForceNew)
  - `recovery_point_objective` (TypeInt, Required, ForceNew)
  - `protect_group_name` (TypeString, Optional)
  - `data_direction` (TypeString, Optional, ForceNew)
- [x] 3.4 Schema computed fields (match design D2, flattened — NO `protect_group_set` key):
  - `app_id` (TypeInt, Computed), `site_pair_name` (TypeString, Computed), `source_region` / `source_zone` / `source_vpc` / `target_region` / `target_zone` / `target_vpc` (TypeString, Computed), `copy_type` / `disaster_recovery_type` / `peer_cloud_name` / `create_from` / `life_state` / `account_uin` / `sub_account_uin` / `create_time` / `modify_time` (TypeString, Computed), `bind_protected_resource_count` (TypeInt, Computed), `error_recovery_point_objective_count` (TypeInt, Computed).
  - `protected_resource_status_set` (TypeList, Computed, Elem: schema.Resource with `status` TypeString Computed and `count` TypeInt Computed).
- [x] 3.5 Implement `resourceTencentCloudBdrcDisasterRecoveryProtectGroupCreate` (design D4): defer LogElapsed/InconsistentCheck; build request from 5 inputs; `resource.Retry(WriteRetryTimeout)` wrapping `CreateDisasterRecoveryProtectGroupWithContext`; nil-safe `result == nil || result.Response == nil`; on retry failure `log.Printf("[CRITAL]%s create bdrc disaster_recovery_protect_group failed, reason:%+v", logId, reqErr)`; verify `ProtectGroupId != nil` (log logId + d.Id before) else return `fmt.Errorf("ProtectGroupId is nil.")`; `d.SetId(*response.Response.ProtectGroupId)`; return Read.
- [x] 3.6 Implement `resourceTencentCloudBdrcDisasterRecoveryProtectGroupRead` (design D5): defer; build `BdrcService`; `protectGroupId := d.Id()`, `protectGroupType := d.Get("protect_group_type").(string)`; call `DescribeDisasterRecoveryProtectGroupById`; if nil → `log.Printf("[CRUD] bdrc disaster_recovery_protect_group id=%s", d.Id())` then `d.SetId("")` then `return nil`; nil-check each field before `d.Set(...)`; flatten `ProtectedResourceStatusSet` into `protected_resource_status_set` list; return nil.
- [x] 3.7 Implement `resourceTencentCloudBdrcDisasterRecoveryProtectGroupUpdate` (design D6): defer; `immutableArgs := []string{"site_pair_id", "protect_group_type", "recovery_point_objective", "data_direction"}`; for each if `d.HasChange(v)` return error "bdrc disaster_recovery_protect_group `<v>` is immutable, cannot be updated, please recreate."; if `d.HasChange("protect_group_name")` build `ModifyProtectGroupAttributeRequest` with `ProtectGroupId = helper.String(d.Id())` and `ProtectGroupName` from `d.GetOk`, wrap in `resource.Retry(WriteRetryTimeout)` with nil-safe guards, `[CRITAL]` log on failure; return Read.
- [x] 3.8 Implement `resourceTencentCloudBdrcDisasterRecoveryProtectGroupDelete` (design D7): defer; build `DeleteDisasterRecoveryProtectGroupsRequest` with `ProtectGroups = []*string{helper.String(d.Id())}`; wrap in `resource.Retry(WriteRetryTimeout)` with nil-safe guards; `[CRITAL]` log on failure; return nil.

## 4. Provider registration

- [x] 4.1 In `tencentcloud/provider.go`, add import `"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/bdrc"`.
- [x] 4.2 In `tencentcloud/provider.go` resource registration map, append `"tencentcloud_bdrc_disaster_recovery_protect_group": bdrc.ResourceTencentCloudBdrcDisasterRecoveryProtectGroup(),`.
- [x] 4.3 In `tencentcloud/provider.md`, add a new `Business Disaster Recovery Center(BDRC)` section with a `Resource` subsection listing `tencentcloud_bdrc_disaster_recovery_protect_group` so gendoc picks it up.

## 5. Documentation

- [x] 5.1 Create `tencentcloud/services/bdrc/resource_tc_bdrc_disaster_recovery_protect_group.md` containing:
  - One-line summary mentioning BDRC: `Provides a resource to create a BDRC disaster recovery protect group.`
  - `Example Usage` HCL block with the 5 input fields populated.
  - `Import` section explaining import uses `ProtectGroupId` (single field).
  - Do NOT include `Argument Reference` / `Attribute Reference` sections (auto-generated by `make doc`).
- [x] 5.2 Run `make doc` to regenerate `website/docs/r/bdrc_disaster_recovery_protect_group.html.markdown`. Hand-editing the generated file is forbidden.

## 6. Unit test (gomonkey mock, not terraform acceptance suite)

- [x] 6.1 Create `tencentcloud/services/bdrc/resource_tc_bdrc_disaster_recovery_protect_group_test.go`, package `bdrc_test`.
- [x] 6.2 Use gomonkey to patch `UseBdrcV20260330Client` (or the connectivity client) and mock the four BDRC client methods (`CreateDisasterRecoveryProtectGroupWithContext`, `DescribeDisasterRecoveryProtectGroupsWithContext`, `ModifyProtectGroupAttributeWithContext`, `DeleteDisasterRecoveryProtectGroupsWithContext`) returning canned `ProtectGroup` responses.
- [x] 6.3 Cover the business-logic flow: Create (assert `d.SetId` with mock ProtectGroupId) → Read (assert computed fields populated) → Update rename (assert `ModifyProtectGroupAttribute` called with new name) → Delete (assert `DeleteDisasterRecoveryProtectGroups` called with the id).
- [x] 6.4 Cover the immutable-arg error path: trigger `d.HasChange` on `protect_group_type` and assert Update returns the immutable error.
- [x] 6.5 Cover the Read not-found path: mock Describe returning empty `ProtectGroupSet` and assert `d.SetId("")` is called with the `[CRUD]` log preserved.
- [x] 6.6 Ensure all generated test code is compilable in the current environment (no `go test` execution required; just constructable).

## 7. Validation

- [x] 7.1 `openspec validate add-bdrc-disaster-recovery-protect-group-resource --strict` passes.
- [x] 7.2 Verify all CRUD SDK call params exist in the corresponding cloud API request structs (Create params in `CreateDisasterRecoveryProtectGroupRequest`, Modify params in `ModifyProtectGroupAttributeRequest`, Delete params in `DeleteDisasterRecoveryProtectGroupsRequest`, Describe params in `DescribeDisasterRecoveryProtectGroupsRequest`).
- [x] 7.3 Verify schema has NO `protect_group_set` / `protect_group_list` wrapper key (list data flattened to top level).
