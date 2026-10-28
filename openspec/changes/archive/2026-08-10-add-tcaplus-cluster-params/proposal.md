## Why

The TcaplusDB `CreateCluster` API supports specifying cluster tags (`ResourceTags`), dedicated server machines (`ServerList`), dedicated proxy machines (`ProxyList`), and the cluster type (`ClusterType`) for creating dedicated clusters. However, the Terraform resource `tencentcloud_tcaplus_cluster` only exposes `idl_type`, `cluster_name`, `vpc_id`, `subnet_id`, `password`, and `old_password_expire_last`, which limits users to creating shared clusters without tags. Users who need to create dedicated TcaplusDB clusters or attach tags at creation time must fall back to the console or raw API.

## What Changes

- Add `cluster_type` (Optional, Computed, immutable after creation) top-level parameter to `tencentcloud_tcaplus_cluster`, mapped to the `CreateCluster` API `ClusterType` field (`*int64`, 1=shared, 2=dedicated). It is only accepted by `CreateCluster` (no update API supports it), so changes after creation SHALL be rejected via the `immutableArgs` check in the Update function. It is read back from `DescribeClusters` response `Clusters.ClusterType`.
- Add `resource_tags` (Optional, TypeList, immutable after creation) nested block parameter to `tencentcloud_tcaplus_cluster`, mapped to the `CreateCluster` API `ResourceTags` field (`[]*TagInfoUnit`). Each element has `tag_key` (string) and `tag_value` (string). It is only accepted by `CreateCluster`, so changes after creation SHALL be rejected via the `immutableArgs` check in the Update function.
- Add `server_list` (Optional, TypeList, immutable after creation, Computed) nested block parameter to `tencentcloud_tcaplus_cluster`, mapped to the `CreateCluster` API `ServerList` field (`[]*MachineInfo`). For creation, each element exposes `machine_type` (string) and `machine_num` (int). On Read, it is refreshed from the `DescribeClusters` response `Clusters.ServerList` (`ServerDetailInfo`), which additionally exposes `server_uid`, `memory_rate`, `disk_rate`, `read_num`, `write_num`, and `version`. Changes after creation SHALL be rejected via the `immutableArgs` check in the Update function.
- Add `proxy_list` (Optional, TypeList, immutable after creation, Computed) nested block parameter to `tencentcloud_tcaplus_cluster`, mapped to the `CreateCluster` API `ProxyList` field (`[]*MachineInfo`). For creation, each element exposes `machine_type` (string) and `machine_num` (int). On Read, it is refreshed from the `DescribeClusters` response `Clusters.ProxyList` (`ProxyDetailInfo`), which additionally exposes `proxy_uid`, `machine_type`, `process_speed`, `average_process_delay`, `slow_process_speed`, and `version`. Changes after creation SHALL be rejected via the `immutableArgs` check in the Update function.
- Extend the `CreateCluster` service-layer function to accept and pass the new parameters (`resourceTags`, `serverList`, `proxyList`, `clusterType`) into the SDK request.
- Update the Read function to populate the new fields from the `DescribeClusters` response (with nil checks before each `set`).
- Add the `immutableArgs` array pattern in the Update function (consistent with the provider's existing immutable-field convention) covering `cluster_type`, `resource_tags`, `server_list`, and `proxy_list`, returning a clear error when any of these change.

## Capabilities

### New Capabilities
- `tcaplus-cluster-params`: Enable the `cluster_type`, `resource_tags`, `server_list`, and `proxy_list` parameters on the `tencentcloud_tcaplus_cluster` resource to allow users to create dedicated TcaplusDB clusters with dedicated server/proxy machines and attach resource tags at creation time.

### Modified Capabilities
<!-- No existing specs require modification -->

## Impact

- **Affected files:**
  - `tencentcloud/services/tcaplusdb/resource_tc_tcaplus_cluster.go` — add the four new schema fields, wire them through the Create flow, add `immutableArgs` check in Update, and add Read support (with nil checks)
  - `tencentcloud/services/tcaplusdb/service_tencentcloud_tcaplus.go` — extend `CreateCluster` to accept `resourceTags []*tcaplusdb.TagInfoUnit`, `serverList []*tcaplusdb.MachineInfo`, `proxyList []*tcaplusdb.MachineInfo`, and `clusterType int64`, and pass them into the SDK request
  - `tencentcloud/services/tcaplusdb/resource_tc_tcaplus_cluster_test.go` — add unit test cases (using gomonkey mocks) covering the new parameters
  - `tencentcloud/services/tcaplusdb/resource_tc_tcaplus_cluster.md` — update documentation with usage examples for the new parameters
- **SDK dependency:** No SDK update required — the vendored `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcaplusdb/v20190823` already exposes `ResourceTags`, `ServerList`, `ProxyList`, and `ClusterType` on `CreateClusterRequest`, and `ClusterType`, `ServerList` (`ServerDetailInfo`), and `ProxyList` (`ProxyDetailInfo`) on the `ClusterInfo` response of `DescribeClusters`.
- **Backward compatibility:** fully backward compatible — all new parameters are Optional; existing configurations continue to work unchanged and still create shared clusters.
- **API constraints:** `ResourceTags`, `ServerList`, `ProxyList`, and `ClusterType` are only accepted by `CreateCluster` (no Modify API supports them), so they are immutable after creation and rejected with a clear error rather than silently destroying the resource.
