## Context

The `tencentcloud_tcaplus_cluster` resource currently supports creating shared TcaplusDB clusters with `idl_type`, `cluster_name`, `vpc_id`, `subnet_id`, `password`, and `old_password_expire_last`. The TcaplusDB `CreateCluster` API also accepts `ResourceTags` (cluster tags), `ServerList` (dedicated server machines), `ProxyList` (dedicated proxy machines), and `ClusterType` (1=shared, 2=dedicated), but the Terraform resource does not expose them.

**Current state:**
- Resource file: `tencentcloud/services/tcaplusdb/resource_tc_tcaplus_cluster.go`
- Service layer: `tencentcloud/services/tcaplusdb/service_tencentcloud_tcaplus.go` (`CreateCluster` function)
- SDK: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcaplusdb/v20190823` (vendored)

**API behavior analysis (verified against vendored SDK):**

| API | ClusterType (req) | ResourceTags (req) | ServerList (req) | ProxyList (req) | ClusterType (resp) | ServerList (resp) | ProxyList (resp) |
|-----|-------------------|--------------------|-------------------|------------------|--------------------|-------------------|------------------|
| `CreateCluster` | Yes (`*int64`) | Yes (`[]*TagInfoUnit`) | Yes (`[]*MachineInfo`) | Yes (`[]*MachineInfo`) | No | No | No |
| `DescribeClusters` | N/A | N/A | N/A | N/A | Yes (`*int64`) | Yes (`[]*ServerDetailInfo`) | Yes (`[]*ProxyDetailInfo`) |
| `ModifyClusterName` | No | No | No | No | N/A | N/A | N/A |
| `ModifyClusterPassword` | No | No | No | No | N/A | N/A | N/A |
| `DeleteCluster` | No | No | No | No | N/A | N/A | N/A |

**SDK struct field mapping (verified):**

- `CreateClusterRequest`: `ClusterType *int64`, `ResourceTags []*TagInfoUnit`, `ServerList []*MachineInfo`, `ProxyList []*MachineInfo`
- `MachineInfo`: `MachineType *string`, `MachineNum *int64` (used for both ServerList and ProxyList on Create)
- `ClusterInfo` (DescribeClusters response element): `ClusterType *int64`, `ServerList []*ServerDetailInfo`, `ProxyList []*ProxyDetailInfo`
- `ServerDetailInfo`: `ServerUid *string`, `MachineType *string`, `MemoryRate *int64`, `DiskRate *int64`, `ReadNum *int64`, `WriteNum *int64`, `Version *string`
- `ProxyDetailInfo`: `ProxyUid *string`, `MachineType *string`, `ProcessSpeed *int64`, `AverageProcessDelay *int64`, `SlowProcessSpeed *int64`, `Version *string`
- `TagInfoUnit`: `TagKey *string`, `TagValue *string`

**Key constraint:** `ClusterType`, `ResourceTags`, `ServerList`, and `ProxyList` are only accepted by `CreateCluster`. No Modify API supports them, so they are immutable after creation.

**Name collision analysis:** The leaf field names `MachineType`, `MachineNum` appear under both `ServerList` and `ProxyList`, and `Version`/`MachineType` appear under both `ServerDetailInfo` and `ProxyDetailInfo`. Therefore these cannot be flattened to top-level schema fields (they would collide). They must be modeled as nested blocks `server_list` and `proxy_list`, each containing the respective fields. This mirrors the existing `resource_tags` nested-block pattern already used by the sibling `tencentcloud_tcaplus_tablegroup` resource.

## Goals / Non-Goals

**Goals:**
- Add `cluster_type` (Optional, Computed, immutable after creation) top-level parameter mapped to `CreateClusterRequest.ClusterType` (`*int64`), refreshed from `ClusterInfo.ClusterType` on Read.
- Add `resource_tags` (Optional, TypeList, immutable after creation) nested block mapped to `CreateClusterRequest.ResourceTags` (`[]*TagInfoUnit`), with `tag_key` and `tag_value` sub-fields. (Read does not populate this block from DescribeClusters because the `DescribeClusters` response does not return tags; tags are write-only at the cluster level through this resource.)
- Add `server_list` (Optional, TypeList, immutable after creation, Computed) nested block mapped to `CreateClusterRequest.ServerList` (`[]*MachineInfo`) on Create, with `machine_type` and `machine_num` sub-fields. On Read, refreshed from `ClusterInfo.ServerList` (`[]*ServerDetailInfo`), exposing `server_uid`, `machine_type`, `memory_rate`, `disk_rate`, `read_num`, `write_num`, and `version`.
- Add `proxy_list` (Optional, TypeList, immutable after creation, Computed) nested block mapped to `CreateClusterRequest.ProxyList` (`[]*MachineInfo`) on Create, with `machine_type` and `machine_num` sub-fields. On Read, refreshed from `ClusterInfo.ProxyList` (`[]*ProxyDetailInfo`), exposing `proxy_uid`, `machine_type`, `process_speed`, `average_process_delay`, `slow_process_speed`, and `version`.
- Extend the `CreateCluster` service-layer function to accept the new parameters and pass them into the SDK request.
- Use the `immutableArgs` array pattern in the Update function (consistent with the provider convention) covering `cluster_type`, `resource_tags`, `server_list`, and `proxy_list`, returning a clear error when any of them changes.
- Maintain full backward compatibility — all new parameters are Optional; existing configurations continue to create shared clusters.

**Non-Goals:**
- Making `cluster_type`, `resource_tags`, `server_list`, or `proxy_list` updatable (no API supports it; they are immutable).
- Adding tag modify/describe APIs (the TcaplusDB cluster API surface has no `ModifyClusterTags`/`DescribeClusterTags` analog; tags are set at creation only through this resource).
- Updating the `tencentcloud_tcaplus_clusters` datasource (out of scope for this change).

## Decisions

### Decision 1: Model `server_list` and `proxy_list` as nested blocks, not flattened fields

**Rationale:** The leaf field names `MachineType`, `MachineNum` (Create) and `MachineType`, `Version` (Read) collide between `ServerList` and `ProxyList`. Flattening them to top-level schema fields would cause ambiguous mapping and violate the Terraform schema uniqueness constraint. Nested blocks `server_list` and `proxy_list` keep each list's fields scoped, matching the existing `resource_tags` block pattern in `tencentcloud_tcaplus_tablegroup`.

### Decision 2: Use `TypeList` (not `TypeSet`) for the nested blocks

**Rationale:** Each `MachineInfo`/`ServerDetailInfo`/`ProxyDetailInfo` entry is positional and order-stable; there is no natural unique key per entry. `TypeList` is simpler and consistent with the API's array semantics. The `resource_tags` block in the sibling `tablegroup` resource uses `TypeSet`; however, tags have a natural key (`tag_key`). Here we follow `TypeList` for server/proxy lists and `TypeList` for tags (to keep the implementation simple and because tags here are write-only at creation).

### Decision 3: `cluster_type` is a top-level `TypeInt` field

**Rationale:** `ClusterType` is a single scalar with no naming collision, so it maps cleanly to a top-level `cluster_type` schema field (`TypeInt`, Optional, Computed, immutable). The API uses `*int64` with values `1` (shared) and `2` (dedicated).

### Decision 4: Immutable (not ForceNew) for the new creation-only parameters

**Rationale:** The Modify APIs (`ModifyClusterName`, `ModifyClusterPassword`) do not accept `ClusterType`, `ResourceTags`, `ServerList`, or `ProxyList`. Rather than using `ForceNew: true` (which silently destroys and recreates the cluster), the Update function uses an `immutableArgs` array that returns a clear error when any of these fields changes. This gives users a better error message than silent destruction.

### Decision 5: `resource_tags` is write-only (not populated on Read)

**Rationale:** The `DescribeClusters` response (`ClusterInfo`) does not include cluster-level tags, and the TcaplusDB cluster API surface has no `DescribeClusterTags` analog (unlike table groups which have `DescribeTableGroupTags`). Therefore `resource_tags` is set on Create only and is not refreshed on Read. Terraform will keep the configured value in state. This is an acceptable trade-off because tags are immutable through this resource anyway.

### Decision 6: Extend existing `CreateCluster` service function signature

**Rationale:** The current `CreateCluster(ctx, idlType, clusterName, vpcId, subnetId, password string)` signature is extended to `CreateCluster(ctx, idlType, clusterName, vpcId, subnetId, password string, resourceTags []*tcaplusdb.TagInfoUnit, serverList []*tcaplusdb.MachineInfo, proxyList []*tcaplusdb.MachineInfo, clusterType int64)`. The new parameters are passed through to the SDK request only when non-empty/non-zero, preserving backward compatibility. The return value stays `(string, error)` (the cluster id).

## Risks / Trade-offs

- **[Risk] `resource_tags` is not refreshed on Read**: The `DescribeClusters` API does not return cluster tags.
  - **Mitigation:** `resource_tags` is immutable after creation, so the configured value remains authoritative in state. This is documented in the resource description.

- **[Risk] Changing any immutable field destroys the resource**: Using the `immutableArgs` pattern (not ForceNew).
  - **Mitigation:** Users receive a clear error explaining which field cannot be changed, rather than silent destruction.

- **[Risk] `server_list`/`proxy_list` Read population differs from Create input**: Create uses `MachineInfo` (`machine_type`, `machine_num`), while Read uses `ServerDetailInfo`/`ProxyDetailInfo` (which have more fields and lack `machine_num`).
  - **Mitigation:** The nested block schema includes all Read fields (`server_uid`, `memory_rate`, etc. for `server_list`; `proxy_uid`, `process_speed`, etc. for `proxy_list`) as Computed, plus the Create-only `machine_num`/`machine_type` as Optional. On Read, all available fields are populated with nil checks before each `Set`.

- **[Risk] Backward compatibility for existing callers of `CreateCluster`**: The service function signature changes.
  - **Mitigation:** The only caller is `resourceTencentCloudTcaplusClusterCreate`, which is updated in the same change. No external callers exist.
