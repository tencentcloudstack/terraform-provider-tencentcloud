# Implementation Tasks: add-emr-cluster-v2-resource

## 1. Service Layer

- [x] 1.1 Append `DescribeEmrClusterV2ById(ctx, instanceId) (*emr.ClusterInstancesInfo, error)` to `tencentcloud/services/emr/service_tencentcloud_emr.go` — wraps `DescribeInstances` with `InstanceIds=[instanceId]`, returns first element or `nil`
- [x] 1.2 Append `DescribeEmrClusterV2Nodes(ctx, instanceId, nodeFlag string) ([]*emr.NodeHardwareInfo, error)` — paginates `DescribeClusterNodes` with `Limit=100`, iterates `Offset` until exhausted
- [x] 1.3 Both helpers use `resource.Retry(tccommon.ReadRetryTimeout, ...)` with `tccommon.RetryError` and emit `[DEBUG]` log lines matching existing EMR service style

## 2. Resource Schema

- [x] 2.1 Create `tencentcloud/services/emr/resource_tc_emr_cluster_v2.go` skeleton with package/imports, `ResourceTencentCloudEmrClusterV2()` constructor, and Importer set to `schema.ImportStatePassthrough`
- [x] 2.2 Declare top-level scalar schema fields: `product_version`, `enable_support_ha_flag`, `instance_name`, `instance_charge_type` (all Required) and `client_token`, `need_master_wan`, `enable_remote_login_flag`, `enable_kerberos_flag`, `custom_conf`, `enable_cbs_encrypt_flag`, `cos_bucket`, `load_balancer_id`, `default_meta_version`, `need_cdb_audit`, `sg_ip`, `partition_number`, `web_ui_version` (Optional)
- [x] 2.3 Declare string-list schema fields: `security_group_ids`, `disaster_recover_group_ids`
- [x] 2.4 Declare `login_settings` nested block (TypeList, MaxItems:1, Required) with `password` (Sensitive) and `public_key_id`
- [x] 2.5 Declare `scene_software_config` nested block (TypeList, MaxItems:1, Required) with `software []string` and `scene_name`
- [x] 2.6 Declare `instance_charge_prepaid` nested block (TypeList, MaxItems:1, Optional) with `period` and `renew_flag`
- [x] 2.7 Declare `script_bootstrap_action_config` nested block (TypeList, Optional) with `cos_file_uri`, `execution_moment`, `args []string`, `cos_file_name`, `remark`
- [x] 2.8 Declare `tags` nested block (TypeList, Optional) with `tag_key` and `tag_value`
- [x] 2.9 Declare `meta_db_info` nested block (TypeList, MaxItems:1, Optional) with `meta_data_jdbc_url`, `meta_data_user`, `meta_data_pass` (Sensitive), `meta_type`, `unify_meta_instance_id`
- [x] 2.10 Declare `depend_service` nested block (TypeList, Optional) with `service_name` and `instance_id`
- [x] 2.11 Declare `node_marks` nested block (TypeList, Optional) with `node_type`, `node_names []string`, `zone`
- [x] 2.12 Declare `zone_resource_configuration` nested block (TypeList, Optional) with child blocks: `virtual_private_cloud` (MaxItems:1, {vpc_id, subnet_id}), `placement` (MaxItems:1, {zone, project_id}), `zone_tag`, and `all_node_resource_spec` (MaxItems:1)
- [x] 2.13 Under `all_node_resource_spec`, declare counts (`master_count`, `core_count`, `task_count`, `common_count`) and four `*_resource_spec` nested blocks (MaxItems:1 each) mapping to `NodeResourceSpec`
- [x] 2.14 Define `NodeResourceSpec` block schema (reusable helper function) with `instance_type`, `system_disk`, `data_disk`, `local_data_disk`, `tags`; each disk entry mirrors `DiskSpecInfo` {count, disk_size, disk_type}
- [x] 2.15 Add `Timeouts` block: `Create: 60 * time.Minute`, `Read: 20 * time.Minute`, `Delete: 30 * time.Minute`
- [x] 2.16 Verify no schema field has `ForceNew: true` (grep sanity check)

## 3. Create Handler

- [x] 3.1 Implement `resourceTencentCloudEmrClusterV2Create` boilerplate: `defer LogElapsed`, `defer InconsistentCheck`, build `logId`/`ctx`, instantiate `emr.NewCreateClusterRequest()`
- [x] 3.2 Build request from schema: one `if v, ok := d.GetOk(...)` block per top-level field, in the same order as the schema declaration
- [x] 3.3 Implement nested builders for `login_settings`, `scene_software_config`, `instance_charge_prepaid`, `meta_db_info` (MaxItems:1 blocks)
- [x] 3.4 Implement nested builders for `script_bootstrap_action_config`, `tags`, `depend_service`, `node_marks` (list blocks) using append + helper conversions
- [x] 3.5 Implement the `zone_resource_configuration` builder including the full `AllNodeResourceSpec` → `NodeResourceSpec` → `DiskSpecInfo` tree (factor `buildNodeResourceSpec` helper)
- [x] 3.6 Call `UseEmrClient().CreateClusterWithContext(ctx, request)` wrapped in `resource.Retry(tccommon.WriteRetryTimeout, ...)`; return non-retryable error if response / `InstanceId` is nil
- [x] 3.7 Set resource ID via `d.SetId(*response.Response.InstanceId)`
- [x] 3.8 Poll `service.DescribeEmrClusterV2ById` inside `resource.Retry(d.Timeout(schema.TimeoutCreate) - time.Minute, ...)` until `Status == 2`; treat known terminated statuses as non-retryable errors
- [x] 3.9 Return `resourceTencentCloudEmrClusterV2Read(d, meta)`

## 4. Read Handler

- [x] 4.1 Implement `resourceTencentCloudEmrClusterV2Read` boilerplate and call `service.DescribeEmrClusterV2ById(ctx, d.Id())`
- [x] 4.2 If result is nil, log warning, `d.SetId("")`, return nil
- [x] 4.3 Populate scalar state fields from `ClusterInstancesInfo` (`instance_name`, `product_version` / `ProductId`→version mapping if needed, `instance_charge_type`, `enable_support_ha_flag` if present, tags, etc.) — only set fields actually returned by the API
- [x] 4.4 Call `service.DescribeEmrClusterV2Nodes(ctx, d.Id(), "all")` and aggregate nodes by type to populate `zone_resource_configuration[].all_node_resource_spec` counts + per-role specs (best effort)
- [x] 4.5 Do NOT overwrite `login_settings.password`, `meta_db_info.meta_data_pass`, `client_token`, `custom_conf` (preserve user config)

## 5. Update Handler (No-op)

- [x] 5.1 Implement `resourceTencentCloudEmrClusterV2Update` boilerplate
- [x] 5.2 Add code comment: `// NOTE: Modify APIs will be added in a follow-up change. Currently Update is a no-op.`
- [x] 5.3 Return `resourceTencentCloudEmrClusterV2Read(d, meta)`

## 6. Delete Handler

- [x] 6.1 Implement `resourceTencentCloudEmrClusterV2Delete` boilerplate and build `emr.NewTerminateInstanceRequest()` with `InstanceId = d.Id()`
- [x] 6.2 Call `UseEmrClient().TerminateInstanceWithContext(ctx, request)` wrapped in `resource.Retry(tccommon.WriteRetryTimeout, ...)`
- [x] 6.3 Poll `DescribeEmrClusterV2ById` inside `resource.Retry(d.Timeout(schema.TimeoutDelete) - time.Minute, ...)` until result is nil or status indicates terminated
- [x] 6.4 Return nil

## 7. Provider Registration

- [x] 7.1 Add `"tencentcloud_emr_cluster_v2": emr.ResourceTencentCloudEmrClusterV2(),` to `tencentcloud/provider.go` ResourcesMap, immediately after `"tencentcloud_emr_cluster"`

## 8. Documentation (Example MD)

- [x] 8.1 Create `tencentcloud/services/emr/resource_tc_emr_cluster_v2.md` following `resource_tc_config_compliance_pack.md` layout
- [x] 8.2 Include `Example Usage` section with a realistic HCL block covering `product_version`, `instance_charge_type=POSTPAID_BY_HOUR`, `login_settings`, `scene_software_config`, and one `zone_resource_configuration` with VPC + Placement + Master/Core resource specs
- [x] 8.3 Include `~> **NOTE:**` explaining that Update is not yet supported in this version
- [x] 8.4 Include `Import` section with `terraform import tencentcloud_emr_cluster_v2.example emr-xxxxxxxx`

## 9. Acceptance Test

- [x] 9.1 Create `tencentcloud/services/emr/resource_tc_emr_cluster_v2_test.go`
- [x] 9.2 Define `TestAccTencentCloudEmrClusterV2_basic` using `tcacctest.AccPreCheck` and `tcacctest.AccProviders`
- [x] 9.3 Add an apply step with minimal valid config and `resource.TestCheckResourceAttrSet("tencentcloud_emr_cluster_v2.example", "id")`
- [x] 9.4 Add an Import step with `ImportStateVerify: true` and `ImportStateVerifyIgnore: []string{"login_settings", "client_token", "meta_db_info"}`
- [x] 9.5 Define `testAccEmrClusterV2_basic` HCL constant near the bottom of the file (mirror the style of `resource_tc_igtm_strategy_test.go`)

## 10. Validation

- [x] 10.1 Run `gofmt -s -w tencentcloud/services/emr/resource_tc_emr_cluster_v2.go tencentcloud/services/emr/resource_tc_emr_cluster_v2_test.go tencentcloud/services/emr/service_tencentcloud_emr.go`
- [x] 10.2 Run `goimports -w` on the same files (skipped: `goimports` not installed locally; `gofmt` covers formatting and `go build` verified imports are valid)
- [x] 10.3 Run `go build ./...` from repo root
- [x] 10.4 Run `go vet ./tencentcloud/services/emr/...`
- [x] 10.5 Run `make lint` (or the closest available target) and resolve any new warnings (verified per-file via IDE linter: only HINT-level `resource.Retry`/`ImportStatePassthrough` deprecations, matching existing project style)
- [x] 10.6 Run `make doc` to regenerate `website/docs/r/emr_cluster_v2.html.markdown` from the `.md` example file (also added `tencentcloud_emr_cluster_v2` to `tencentcloud/provider.md` so gendoc picks it up)
- [x] 10.7 Run `openspec validate add-emr-cluster-v2-resource --strict` and confirm "valid"
- [x] 10.8 Smoke-compile the acceptance test without `TF_ACC`: `go test -run TestAccTencentCloudEmrClusterV2_basic -count=1 ./tencentcloud/services/emr/...` (expect skip, not compile error)

## 11. Phase 2: Dynamic Node Count via Resource-Spec List Length

- [x] 11.1 Schema: remove `master_count`, `core_count`, `task_count`, `common_count` from `all_node_resource_spec`
- [x] 11.2 Schema: remove `MaxItems: 1` from `master_resource_spec`, `core_resource_spec`, `task_resource_spec`, `common_resource_spec` so they become unbounded `TypeList`
- [x] 11.3 Create handler: derive `MasterCount` / `CoreCount` / `TaskCount` / `CommonCount` from `len(<role>_resource_spec)`; pass the first element of each list as the SDK's single-valued `*NodeResourceSpec` template
- [x] 11.4 Factor the spec-list reader into a small helper (`buildNodeResourceSpecAndCount`) that returns both `(*emr.NodeResourceSpec, *int64)` for a given `[]interface{}`
- [x] 11.5 Update `resource_tc_emr_cluster_v2.md` Example Usage to remove explicit `master_count`/`core_count` and instead show N blocks per role driving the count implicitly
- [x] 11.6 Update `resource_tc_emr_cluster_v2_test.go` HCL fixture to match the new layout
- [x] 11.7 Run `gofmt`, `go build ./...`, `go vet ./tencentcloud/services/emr/...`, smoke-compile test, and `make doc`

## 12. Phase 3: Same-Role Spec Uniformity Validation in Create

- [x] 12.1 Implement helper `validateEmrNodeResourceSpecUniformity(role string, specList []interface{}) error` that compares every element of the list to the first element field-by-field (instance_type, system_disk, data_disk, local_data_disk, tags); returns a descriptive error if any divergence is found
- [x] 12.2 In `resourceTencentCloudEmrClusterV2Create`, call the helper for each role (`master_resource_spec`, `core_resource_spec`, `task_resource_spec`, `common_resource_spec`) before building the SDK request, returning early on validation failure
- [x] 12.3 Run linter / IDE diagnostics and confirm no new errors

## 13. Naming cleanup: remove residual "token" wording

Background: An earlier iteration auto-generated stable identity strings (called "token") via `helper.BuildToken()`. That mechanism was fully removed and replaced by user-supplied `_node_index` / `_disk_index` Required string fields. The runtime code no longer generates any token, but variable names, function names, comments, and schema `Description` strings still say "token", which is misleading. This section is a pure rename / wording cleanup with **zero business-logic changes**.

- [ ] 13.1 Schema `Description` strings: rewrite the descriptions of `_node_index` and `_disk_index` so they no longer contain the word "token"; explain that they are user-supplied stable identifiers required because `*_resource_spec` and `data_disk` are `TypeList` with `DiffSuppressFunc`, and identical configuration blocks must still be matched 1:1 across plan/apply
- [ ] 13.2 In `resourceTencentCloudEmrClusterV2Read`, rename local variables in the config-side index extraction block (around the `cfgDiskEntry` struct, `configNodeTokens`, `configDisksPerNode`): `cfgDiskEntry.token` → `cfgDiskEntry.diskIndex`; `configNodeTokens` → `configNodeIndexes`; the inner `tok` / `dTok` → `nodeIdxStr` / `diskIdxStr`
- [ ] 13.3 In the same Read handler, rename the state-side mapping locals: `stateNodeByRID` value type and var stays a string map but rename the inner `tok` to `nodeIdxStr`; `stateDiskByID` inner `dTok` → `diskIdxStr`
- [ ] 13.4 In the same Read handler, rename the API-disk resolution locals: `nodeTok` → `userNodeIdx`; `diskTok` → `userDiskIdx`; update all references
- [ ] 13.5 Rename helper `emrNodeListEqualByToken` → `emrNodeListEqualByNodeIndex` and `emrDiskListEqualByToken` → `emrDiskListEqualByDiskIndex`; update all call sites in `emrNodeListOrderSuppressFunc`, `emrDataDiskListOrderSuppressFunc`, and `emrNodeContentEqual`
- [ ] 13.6 Inside the renamed `emrNodeListEqualByNodeIndex` and `emrDiskListEqualByDiskIndex`, rename `oldByTok` / `newByTok` → `oldByNodeIndex` / `newByNodeIndex` (and disk equivalents); rename the inner `tok` loop variable to `nodeIdxStr` / `diskIdxStr`
- [ ] 13.7 In `alignNodeListByNodeIndex` (Update path), rename `oldByToken` → `oldByNodeIndex`; `matchedTokens` → `matchedNodeIndexes`; inner `tok` → `nodeIdxStr`
- [ ] 13.8 In `handleNodeDataDiskChange` (Update path), rename `oldByToken` → `oldByDiskIndex`; `matchedTokens` → `matchedDiskIndexes`; inner `tok` → `diskIdxStr`; update all error message wording so it says `_disk_index` rather than "token" or "_disk_index token"
- [ ] 13.9 Replace every occurrence of the words "token", "Token", "tok" inside Go comments within `resource_tc_emr_cluster_v2.go` with `_node_index` / `_disk_index` / "user-supplied stable index" as appropriate
- [ ] 13.10 Run `gofmt -s -w` and `go build ./...` to confirm zero compile errors
- [ ] 13.11 Run `go vet ./tencentcloud/services/emr/...` and IDE linter to confirm no new warnings
- [ ] 13.12 Final sweep: `grep -nE "token|Token|\\btok\\b" tencentcloud/services/emr/resource_tc_emr_cluster_v2.go` should return zero hits (excluding the unrelated upstream SDK `client_token` request field, which is a separate Terraform schema attribute)

## 14. Final fix: TypeList + CustomizeDiff to truly suppress reorder drift

Background: The `DiffSuppressFunc` placed on the four `*_resource_spec` and on `data_disk` was proven to never be invoked by the SDK (`schemaMap.diff` calls a leaf field's own `DiffSuppressFunc`; list `.#` diffs go through a temporary `countSchema` that has no `DiffSuppressFunc`). As a result, any block reordering by the user produces a full positional diff that nothing in the schema layer can suppress. The only correct fix is a top-level `CustomizeDiff` that, before SDK diff calculation, reorders the new list to match the old list by user-supplied `_node_index` / `_disk_index`. Once orders match, every `.#` and leaf diff naturally collapses to the user's true intent.

- [x] 14.1 Add `CustomizeDiff: customizeDiffEmrClusterV2` to `resourceTencentCloudEmrClusterV2()` schema definition
- [x] 14.2 Implement `customizeDiffEmrClusterV2(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error`:
  - skip when resource is being created (`d.Id() == ""`) — first apply has no old state to align against
  - call `d.GetChange("zone_resource_configuration")`; cast both old and new to `[]interface{}`
  - for each zone in new (matched to old by `placement.0.zone`), reorder each role list (`master_resource_spec`, `core_resource_spec`, `task_resource_spec`, `common_resource_spec`) so blocks present in old appear in old's order; new-only blocks (i.e., `_node_index` not found in old) are appended in their original config order
  - for each retained node, reorder `data_disk` similarly by `_disk_index`; new-only disks appended at the end
  - call `d.SetNew("zone_resource_configuration", reordered)` to commit the reordered new value back into the diff
- [x] 14.3 Remove `DiffSuppressFunc: emrNodeListOrderSuppressFunc` from all four `*_resource_spec` schema entries (lines 363/372/381/390); they are confirmed dead code by SDK source review
- [x] 14.4 Remove `DiffSuppressFunc: emrDataDiskListOrderSuppressFunc` from the `data_disk` schema in `emrNodeSpecElem` (line ~2347)
- [x] 14.5 Delete the now-unused helper functions: `emrNodeListOrderSuppressFunc`, `emrDataDiskListOrderSuppressFunc`, `emrNodeListEqualByNodeIndex`, `emrDiskListEqualByDiskIndex`, `emrNodeContentEqual`, `emrSystemDiskEqual` (verify no other call sites first)
- [x] 14.6 Run `gofmt -s -w`, `go build ./tencentcloud/services/emr/...`, `go vet ./tencentcloud/services/emr/...` and confirm zero errors
- [ ] 14.7 Manual smoke test on `examples/dist/main.tf`: re-apply with no config changes (expect zero drift); reorder several blocks (expect zero drift); resize one disk_size (expect a single targeted diff line)

## 15. Final refactor: TypeList + _node_index/_disk_index + CustomizeDiff (overwrite Computed)

Background: After many iterations the only architecture that satisfies all of the user's drift / scale-out / scale-in requirements is:

- TypeList for the four `*_resource_spec` and for `data_disk` (so identical-shape blocks can coexist as N entries)
- Required user-supplied `_node_index` and `_disk_index` (logical identity for matching across plan/apply)
- A single `CustomizeDiff` that **does NOT reorder new**; it instead overwrites Computed fields (`emr_resource_id`, `order_no`, `disk_id`) on each `new` block by matching to the corresponding `old` block via `_node_index` / `_disk_index`. This eliminates positional Computed-field swap drift while leaving the user's HCL declaration order as the authoritative shape of the plan.
- `Read` writes API order to state (no attempt to align state to user config order), because plan-refresh Read cannot access user config (SDK ReadResource RPC does not transmit it).
- All business constraints (master/common length immutable, core/task scale-out/in, disk additions only, no shrinking, immutable system_disk) are enforced inside CustomizeDiff so the user sees errors at `plan` time rather than mid-apply.

This section replaces sections 13 and 14 conceptually; those earlier attempts left over half-baked patches and obsolete naming. We perform a full clean reconstruction below.

- [ ] 15.1 Remove the workspace-only state surgeries done during debugging (terraform.tfstate.bak.* backups can stay). The .tfstate file should be left as the user finds it; we do not modify it.
- [x] 15.2 Schema cleanup in `resourceTencentCloudEmrClusterV2`:
  - Keep the four `*_resource_spec` and `data_disk` as `TypeList` (no `MaxItems` change).
  - Remove every `DiffSuppressFunc` left over from earlier attempts on these list fields.
  - Keep the existing `_node_index` and `_disk_index` (Required, TypeString); rewrite their Description strings so they explicitly describe their identity-matching role and document that they must be unique within their parent list and stable across plan/apply.
  - Add `CustomizeDiff: customizeDiffEmrClusterV2` at the resource level (already present from a previous iteration; ensure its body matches the new specification below).
- [x] 15.3 Rewrite `customizeDiffEmrClusterV2` to implement the new contract:
  1. `d.Id() == ` → return nil (Create).
  2. Pull old/new `zone_resource_configuration` via `d.GetChange`.
  3. For each zone (matched by `placement.0.zone`):
     - For each role (`master_resource_spec`, `core_resource_spec`, `task_resource_spec`, `common_resource_spec`):
       - Validate uniqueness of `_node_index` within new list; uniqueness of `_disk_index` within each block's data_disk.
       - **Business rule** master/common: if `len(new) != len(old)` → error "master/common 列表长度不可变更（不支持扩缩容）".
       - **Business rule** data_disk: if `old` contains a `_disk_index` not present in `new` → error "data_disk 不支持减少".
       - **Business rule** disk_size: for each `(_disk_index)` matched pair, if new < old → error "data_disk.disk_size 不支持缩容".
       - **Business rule** system_disk: for each matched node, if old.system_disk fields differ from new.system_disk → error "system_disk 创建后不可更改".
       - For each new block:
         - find old by `_node_index`; if found, overwrite `emr_resource_id`, `order_no` with old's values
         - for each new.data_disk, find old.data_disk by `_disk_index`; if found, overwrite `disk_id`
   4. `d.SetNew("zone_resource_configuration", newZones)`.
  - The CustomizeDiff must NEVER reorder `new` lists. The user's HCL declaration order is the authoritative plan shape.
- [x] 15.4 Read handler simplification:
  - Remove all logic that tries to use `d.GetOk("zone_resource_configuration")` to build a 'config order map' for repositioning state — that approach was proven non-functional during plan-refresh (no config available there).
  - Remove all 'reorder apiDisks to match user config' / 'reorder specs to match user config' blocks.
  - Keep the existing logic that uses `stateNodeByRID` / `stateDiskByID` to preserve user-supplied `_node_index` / `_disk_index` strings across reads (this is what makes Read stable across multiple plan-refreshes).
  - Read writes API order to state. State physical order is decoupled from user HCL order; that is fine because CustomizeDiff bridges the two at plan time.
- [x] 15.5 Update handler simplification:
  - Keep `alignNodeListByNodeIndex` and `handleNodeDataDiskChange` as the matching engines (they already key on `_node_index` / `_disk_index`).
  - Verify no path uses `d.HasChange(...positional...)` for the four resource_spec lists; everything must go through `d.GetChange` + identity-keyed comparison.
- [x] 15.6 Documentation: rewrite the `Description` of `_node_index` and `_disk_index` so the user understands these are required stable identifiers, must be unique within their parent list, and must remain stable across plan/apply. Document that changing the index value mid-stream is treated as 'delete the old, add a new' and may trigger destructive operations on disabled-scale roles.
- [x] 15.7 Run `gofmt -s -w`, `go build ./...`, `go vet ./tencentcloud/services/emr/...`; resolve any new errors.
- [ ] 15.8 Manual smoke validation in `examples/dist`:
  - Re-apply with no changes → expected: zero drift
  - Reorder several blocks in main.tf → expected: zero drift
  - Resize one disk_size up → expected: single targeted diff
  - Add a new core block at any position → expected: single + diff
  - Remove a core block from any position → expected: single - diff
  - Try to change master count → expected: plan-time error from CustomizeDiff
  - Try to shrink a disk_size → expected: plan-time error
  - Try to change system_disk → expected: plan-time error

## 16. Phase 6: TypeSet refactor (final solution for ordering drift)

Background: Sections 14/15 attempted TypeList + CustomizeDiff to bridge config-order vs state-order. Multiple rounds of testing proved that TypeList's positional diff cannot be reliably suppressed when users insert/remove blocks in the middle of a list. The user and the agent jointly settled on the following permanent design:

- Each `*_resource_spec` (master/core/task/common) becomes `TypeSet` with a custom hash function that hashes ONLY `_node_index`.
- `data_disk` (inside each spec) becomes `TypeSet` with a custom hash function that hashes ONLY `_disk_index`.
- `_node_index` and `_disk_index` remain Required user-supplied strings; they are the identity keys.
- All TypeList-era patches (CustomizeDiff overwriting Computed fields, Read reordering specs/disks to match config order, two-pass node identity pre-resolution in Read) are deleted.
- CustomizeDiff retains only business-rule validation: master/common length immutable; core/task allows scale-out/in but not rename; data_disk allows scale-out (size up only) but never scale-in; disk_type immutable; system_disk immutable; `_node_index` unique within spec; `_disk_index` unique within node.

User contract (documented in schema Description):
1. `_node_index` MUST be set, MUST be unique within the same role list.
2. `_disk_index` MUST be set, MUST be unique within the same node's data_disk list.
3. Once written to state, `_node_index` and `_disk_index` MUST NOT be renamed; renaming is detected by CustomizeDiff and rejected.
4. State migration: existing TypeList state cannot be migrated; users must `terraform destroy && terraform apply` to recreate.

Tasks:

- [x] 16.1 Schema: change `master_resource_spec`, `core_resource_spec`, `task_resource_spec`, `common_resource_spec` from `TypeList` to `TypeSet`; remove any `MaxItems`. Attach `Set: hashEmrNodeResourceSpec` (custom hash hashing only `_node_index`).
- [x] 16.2 Schema: change `data_disk` (inside each spec) from `TypeList` to `TypeSet`. Attach `Set: hashEmrDataDisk` (custom hash hashing only `_disk_index`).
- [x] 16.3 Add helper `hashEmrNodeResourceSpec(v interface{}) int` — extracts `_node_index` from the map and returns `schema.HashString(_node_index)`.
- [x] 16.4 Add helper `hashEmrDataDisk(v interface{}) int` — extracts `_disk_index` and returns `schema.HashString(_disk_index)`.
- [x] 16.5 Schema Description: rewrite `_node_index` and `_disk_index` Description to reflect the TypeSet identity-key semantics and the four user contract rules above.
- [x] 16.6 Replace `customizeDiffEmrClusterV2` body with a TypeSet-native validator (implementation chose to compute add/remove via `_node_index` map diff rather than `(*schema.Set).Difference`, because the maps already need to be built for paired-node validation; the result is equivalent).
- [x] 16.7 Read handler simplification — deleted the "reorder apiDisks to match cfgDiskList" block and the "reorder specs to match cfgList" block. Kept the two-pass pre-resolution because it is still needed for first-Read-after-Create to populate `_node_index` from user config (TypeSet hashes empty `_node_index` to the same bucket — leaving it blank would collapse all nodes into one set entry).
- [x] 16.8 CustomizeDiff cleanup — deleted all Computed-field-overwrite logic (`emr_resource_id` / `order_no` / `disk_id`) and the `d.SetNew("zone_resource_configuration", ...)` call. The new CustomizeDiff is pure validation; it never mutates `new`.
- [x] 16.9 Update handler verification — `nodeRoleChanged`, `handleNodeDataDiskChange`, and all data_disk read sites already routed through `emrNodeSetToList` / `emrDiskSetToList` (defensive helpers that accept both `*schema.Set` and `[]interface{}`), so the schema-type change required no callsite edits. Helper comments updated to reflect that TypeSet is now the canonical schema type.
- [x] 16.10 Ran `gofmt -s -w`, `go build ./...`, `go vet ./tencentcloud/services/emr/...` — all EXIT=0, no new lint errors (45 HINT items are all pre-existing deprecation warnings in unrelated handlers).
- [ ] 16.11 Manual smoke validation in `examples/dist` (user-driven; provider side only confirms build):
  - `terraform destroy && terraform apply` (mandatory state rebuild).
  - `terraform plan` immediately after apply → expected: zero drift.
  - Reorder several blocks in main.tf → `terraform plan` → expected: zero drift.
  - Change `disk_size` 100 → 110 on one disk → expected: single targeted diff (in-place modify, no remove+add).
  - Add a new `core_resource_spec` block at any position → expected: single + diff (scale-out).
  - Remove a `core_resource_spec` block from any position → expected: single - diff (scale-in).
  - Try to remove a master block → expected: plan-time error.
  - Try to shrink `disk_size` → expected: plan-time error.
  - Try to remove a `data_disk` → expected: plan-time error.
  - Try to rename a `_node_index` → expected: plan-time error.
  - Try to rename a `_disk_index` → expected: plan-time error.
  - Try to change `system_disk.disk_size` → expected: plan-time error.


## 17. Move `software` from `scene_software_config` to per-node `soft_ware` (May 27 2026)

Background: the API's SoftInfo (deduped software list) better reflects reality when each node role declares which components/processes it runs. Refactor the user-facing schema to:

1. Remove `scene_software_config.software`. The block remains Required and now contains only `scene_name`.
2. Add a new `soft_ware` field inside each of the four `*_resource_spec` blocks. Shape:
   - `soft_ware` is `TypeSet`, `Optional`, `ForceNew=true` (cannot be modified after create).
   - Each element has:
     - `name`: `TypeString`, `Required` — component name with version, e.g. `hdfs-3.2.2`.
     - `process`: `TypeSet` of `TypeString`, `Required` — process list for this component on this role.
3. Create-time uniformity: every `*_resource_spec` block within the same role must declare an identical `soft_ware` set (full `name + process` content equality, multiset semantics).
4. Cross-role aggregation: at Create time, collect `name` from `soft_ware` of all four roles, dedupe, and pass that aggregated list as `request.SceneSoftwareConfig.Software` to `CreateCluster` — replacing the removed `scene_software_config.software` input path.
5. Read: do NOT overwrite `soft_ware` on Read (process list is not returned by API). `scene_software_config.scene_name` continues to round-trip from API.
6. Update: `soft_ware` is ForceNew, so any modification triggers destroy+recreate by the SDK; no Update branch needed.

Tasks:

- [x] 17.1 Schema: remove `software` field from `scene_software_config` Elem.
- [x] 17.2 Schema: add `soft_ware` (TypeSet, Optional, ForceNew) to `emrNodeSpecElem()` with `name` (Required string) and `process` (Required TypeSet of strings) sub-fields.
- [x] 17.3 Create: remove the `software` extraction from `scene_software_config`. Add a new helper that, given the four role specs of a single zone (or all zones), aggregates and dedupes `soft_ware[*].name` into a `[]*string` for `request.SceneSoftwareConfig.Software`.
- [x] 17.4 Create-time uniformity: extend `validateEmrNodeResourceSpecUniformity` (or its caller) to compare `soft_ware` content across blocks within the same role, rejecting non-uniform input with a clear error.
- [x] 17.5 Create: SDK gap — verified that `emr.NodeResourceSpec` has no `SoftWare` field (only `InstanceType`, `SystemDisk`, `Tags`, `DataDisk`, `LocalDataDisk`). Per-role `soft_ware` is therefore consumed only at the provider layer (uniformity validation + name aggregation into `SceneSoftwareConfig.Software`); the user-supplied `process` list is currently informational and not transmitted to the EMR API. Documented in the field Description.
- [x] 17.6 Read: drop the `if cluster.Config != nil && cluster.Config.SoftInfo != nil { item["software"] = cluster.Config.SoftInfo }` block (the field no longer exists in schema). Keep the `scene_name` round-trip.
- [ ] 17.7 Update existing `examples/dist/main.tf` and `resource_tc_emr_cluster_v2.md` / `_test.go` to use the new `soft_ware` shape.
- [x] 17.8 `gofmt` + `go build` + `go vet` validation.
- [ ] 17.9 Manual validation (user-driven): destroy old cluster, apply with new schema, verify second apply yields no diff; verify trying to modify `soft_ware` triggers ForceNew (destroy+recreate plan).

## 18. Refinements on Section 17 schema (May 27 2026, follow-up)

After 17 landed, two adjustments were made:

1. `scene_software_config.software` is restored as a Computed-only field. It is no longer accepted as user input, but it is round-tripped from `cluster.Config.SoftInfo` in Read so users can still observe the cluster-wide deployed component list in state.
2. The per-role `soft_ware` element fields are renamed: `name` → `services` (still the versioned component name like `hdfs-3.2.2`), `process` → `roles` (still the process list for that component on this role). Semantics unchanged.

Tasks:

- [x] 18.1 Schema: add back `software` as `TypeList` of strings, `Computed: true`, no longer Required.
- [x] 18.2 Schema: rename `soft_ware` element fields `name` → `services`, `process` → `roles`. Keep ForceNew on both.
- [x] 18.3 Helpers: `emrCollectSoftwareNames` now reads `services` instead of `name`. `emrSoftWareSetsEqual` now uses `services` and `roles`. Update both.
- [x] 18.4 Read: re-add the `item["software"] = cluster.Config.SoftInfo` round-trip alongside `scene_name`. Keep skipping `soft_ware` (process list isn't returned by API anyway).
- [x] 18.5 `gofmt` + `go build` + `go vet` validation.
