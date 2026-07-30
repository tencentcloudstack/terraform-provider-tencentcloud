## 1. Service Layer

- [x] 1.1 Append `DescribeDBCustomClusterById(ctx, clusterId)` to `service_tencentcloud_dbdc.go` — wraps `DescribeDBCustomClusterDetail` with `resource.Retry` + `ratelimit.Check`, returns `*dbdcv20201029.DescribeDBCustomClusterDetailResponseParams`, nil-safe
- [x] 1.2 Append `DescribeDBCustomTaskStatusById(ctx, taskId)` — wraps `DescribeDBCustomTaskStatus` with retry, returns task `Status` string, nil-safe

## 2. Resource Implementation

- [x] 2.1 Create `resource_tc_dbdc_db_custom_cluster.go` with full schema (arguments aligned to `CreateDBCustomClusterRequest`; `container_network`, `api_server_network` nested blocks; `tags` map; computed read-only fields)
- [x] 2.2 Implement Create: build request, call `CreateDBCustomCluster` (retry), poll `DescribeDBCustomTaskStatus` until `Succeeded`, `SetId(ClusterId)`
- [x] 2.3 Implement Read: call `DescribeDBCustomClusterById`, populate all fields with nil guards, clear ID when not found
- [x] 2.4 Implement Update: on `tags` change compute `AddTags`/`DeleteTagKeys`, call `ModifyDBCustomClusterTags` (retry)
- [x] 2.5 Implement Delete: call `DestroyDBCustomCluster` (retry), poll `DescribeDBCustomTaskStatus` until `Succeeded`
- [x] 2.6 Add shared `waitDBCustomTaskSucceeded` task-polling helper

## 3. Provider Registration

- [x] 3.1 Register `tencentcloud_dbdc_db_custom_cluster` in `provider.go` ResourcesMap

## 4. Documentation & Tests

- [x] 4.1 Create `resource_tc_dbdc_db_custom_cluster.md` (Example Usage + Import, per `resource_tc_config_compliance_pack.md` convention)
- [x] 4.2 Create `resource_tc_dbdc_db_custom_cluster_test.go` (per `resource_tc_config_compliance_pack_test.go` convention)
- [x] 4.3 Add website doc `website/docs/r/dbdc_db_custom_cluster.html.markdown` (+ erb nav entry)

## 5. Verification

- [x] 5.1 `go build ./tencentcloud/...` compiles cleanly
- [x] 5.2 `go vet ./tencentcloud/services/dbdc/...` passes
