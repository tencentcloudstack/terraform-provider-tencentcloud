## 1. Resource Implementation

- [x] 1.1 Create `tencentcloud/services/tke/resource_tc_kubernetes_cluster_roll_out_sequence_tag_config.go` with constructor `ResourceTencentCloudKubernetesClusterRollOutSequenceTagConfig()`, config-type code style following `tencentcloud_waf_owasp_rule_status_config`
- [x] 1.2 Define schema matching `ModifyClusterRollOutSequenceTags` input: `cluster_id` (String, Required, ForceNew) and `tags` (List, Required) with nested `key` (String, Required) and `value` (String, Required)
- [x] 1.3 Implement Create (config-type): set `d.SetId(cluster_id)` then delegate to Update (no direct API call in Create)
- [x] 1.4 Implement Read: delegate to service-layer `TkeService.DescribeKubernetesClusterRollOutSequenceTagConfigById` (paginated Limit=100 + retry + ratelimit, client-side match on `ClusterID`), nil-safe flatten `Tags` into `tags`; `d.SetId("")` when not found or tag list empty
- [x] 1.5 Implement Update: build full `Tags`, call `ModifyClusterRollOutSequenceTagsWithContext` inside retry, then call Read
- [x] 1.6 Implement Delete: call `ModifyClusterRollOutSequenceTags` with cluster id and empty `Tags` list inside retry (removes all tags)
- [x] 1.7 Ensure all response-value access is nil-safe and every API call uses the retry mechanism

## 2. Provider Registration

- [x] 2.1 Register `tencentcloud_kubernetes_cluster_roll_out_sequence_tag_config` in `tencentcloud/provider.go` ResourcesMap
- [x] 2.2 Add the resource entry to `tencentcloud/provider.md`

## 3. Documentation

- [x] 3.1 Create `tencentcloud/services/tke/resource_tc_kubernetes_cluster_roll_out_sequence_tag_config.md` example file, format following `resource_tc_config_compliance_pack.md` (Example Usage + Import)
- [x] 3.2 Generate `website/docs/r/kubernetes_cluster_roll_out_sequence_tag_config.html.markdown` via `make doc` (do not hand-write) and confirm `website/tencentcloud.erb` link entry

## 4. Unit Test

- [x] 4.1 Create `tencentcloud/services/tke/resource_tc_kubernetes_cluster_roll_out_sequence_tag_config_test.go`, naming/format following `resource_tc_config_compliance_pack_test.go` (basic create + update + import steps)

## 5. Verification

- [x] 5.1 Run `gofmt`/`go build ./tencentcloud/...` and ensure no compile errors
- [x] 5.2 Run `go vet ./tencentcloud/services/tke/` and `read_lints` to confirm no newly introduced errors
- [x] 5.3 Run `make doc` and verify generated website doc is consistent with the `.md` example
