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

