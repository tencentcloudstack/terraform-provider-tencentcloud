## 1. Service Layer

- [x] 1.1 Append `DescribeDBCustomClusterNodeById(ctx, clusterId, nodeId)` to `service_tencentcloud_dbdc.go` — wraps `DescribeDBCustomClusterNodes` (`ClusterId`, `Limit=100`, pagination) with `resource.Retry` + `ratelimit.Check`; returns matching `*dbdcv20201029.DBCustomClusterNode` (nil if not found), nil/length-safe
- [x] 1.2 Reuse existing `waitDBCustomTaskSucceeded` (already in the `dbdc` package)

## 2. Resource Implementation

- [x] 2.1 Create `resource_tc_dbdc_node_to_db_custom_cluster_attachment.go` with schema (`cluster_id`, `node_id`, `image_id`, `login_settings` — all ForceNew; computed node fields)
- [x] 2.2 Implement Create: call `AddNodesToDBCustomCluster` (`NodeIds=[node_id]`, retry), poll task until `Succeeded`, `SetId(ClusterId#NodeId)`
- [x] 2.3 Implement Read: split composite ID, call `DescribeDBCustomClusterNodeById`, populate fields with nil guards, clear ID when not found
- [x] 2.4 Implement Delete: split composite ID, call `RemoveNodesFromDBCustomCluster` (`NodeIds=[node_id]`, retry), poll task until `Succeeded`

## 3. Provider Registration

- [x] 3.1 Register `tencentcloud_dbdc_node_to_db_custom_cluster_attachment` in `provider.go` ResourcesMap

## 4. Documentation & Tests

- [x] 4.1 Create `resource_tc_dbdc_node_to_db_custom_cluster_attachment.md` (Example Usage + Import, per `resource_tc_config_compliance_pack.md` convention)
- [x] 4.2 Create `resource_tc_dbdc_node_to_db_custom_cluster_attachment_test.go` (per `resource_tc_config_compliance_pack_test.go` convention)
- [x] 4.3 Add website doc `website/docs/r/dbdc_node_to_db_custom_cluster_attachment.html.markdown` (+ erb nav entry)

## 5. Verification

- [x] 5.1 `go build ./tencentcloud/...` compiles cleanly
- [x] 5.2 `go vet ./tencentcloud/services/dbdc/...` passes
