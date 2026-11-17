## 1. Schema Definition

- [x] 1.1 Add `cluster_type` top-level parameter (`TypeInt`, Optional, Computed) to the `tencentcloud_tcaplus_cluster` resource schema in `tencentcloud/services/tcaplusdb/resource_tc_tcaplus_cluster.go`
- [x] 1.2 Add `resource_tags` nested block (`TypeList`, Optional, Elem with `tag_key`/`tag_value` string sub-fields) to the resource schema
- [x] 1.3 Add `server_list` nested block (`TypeList`, Optional, Computed, Elem with `machine_type` string, `machine_num` int for Create plus `server_uid`/`memory_rate`/`disk_rate`/`read_num`/`write_num`/`version` Computed sub-fields for Read) to the resource schema
- [x] 1.4 Add `proxy_list` nested block (`TypeList`, Optional, Computed, Elem with `machine_type` string, `machine_num` int for Create plus `proxy_uid`/`process_speed`/`average_process_delay`/`slow_process_speed`/`version` Computed sub-fields for Read) to the resource schema

## 2. Service Layer

- [x] 2.1 Extend the `CreateCluster` function signature in `tencentcloud/services/tcaplusdb/service_tencentcloud_tcaplus.go` to accept `resourceTags []*tcaplusdb.TagInfoUnit`, `serverList []*tcaplusdb.MachineInfo`, `proxyList []*tcaplusdb.MachineInfo`, and `clusterType int64`
- [x] 2.2 Pass the new parameters into the `CreateClusterRequest` only when non-empty/non-zero, preserving backward compatibility for existing callers

## 3. CRUD Implementation

- [x] 3.1 Update `resourceTencentCloudTcaplusClusterCreate` to read the new schema fields from `*schema.ResourceData`, build the SDK parameter objects (`TagInfoUnit`, `MachineInfo`), and pass them to the extended `CreateCluster` service function
- [x] 3.2 Update `resourceTencentCloudTcaplusClusterRead` to populate `cluster_type`, `server_list`, and `proxy_list` from the `DescribeClusters` response (`ClusterInfo`) with nil checks before each `set`; note `resource_tags` is write-only and not refreshed on Read
- [x] 3.3 Update `resourceTencentCloudTcaplusClusterUpdate` to add `cluster_type`, `resource_tags`, `server_list`, and `proxy_list` to the `immutableArgs` array, returning a clear error when any of them changes (they are not accepted by any Modify API)

## 4. Tests

- [x] 4.1 Add unit test cases (using gomonkey mocks for the cloud API) in `tencentcloud/services/tcaplusdb/resource_tc_tcaplus_cluster_test.go` covering creation with the new parameters (dedicated cluster with `cluster_type`, `resource_tags`, `server_list`, `proxy_list`)
- [x] 4.2 Add unit test cases covering the Read path for the new nested blocks (`server_list`/`proxy_list` field population from `DescribeClusters` response with nil checks)
- [x] 4.3 Add unit test cases covering the `immutableArgs` rejection in the Update function when `cluster_type`, `resource_tags`, `server_list`, or `proxy_list` change

## 5. Documentation

- [x] 5.1 Update `tencentcloud/services/tcaplusdb/resource_tc_tcaplus_cluster.md` with example usage demonstrating the new `cluster_type`, `resource_tags`, `server_list`, and `proxy_list` parameters (using `jsonencode()` for any JSON-string field values if applicable)
