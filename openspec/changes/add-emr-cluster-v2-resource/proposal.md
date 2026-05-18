# Add tencentcloud_emr_cluster_v2 Resource

## Why

The existing `tencentcloud_emr_cluster` resource is based on the legacy `CreateInstance` API, which uses a flat parameter layout and does not support newer EMR capabilities such as `SceneSoftwareConfig`, multi-AZ `ZoneResourceConfiguration`, `NodeMark`-based node identification, shared metadata DB (`CustomMetaDBInfo`), cross-cluster `DependService`, or StarRocks storage-compute separation (`CosBucket`).

The EMR team ships a new, recommended creation API `CreateCluster` (EMR 2019-01-03, documented at 589/83953) which exposes the full set of modern parameters. Users need a first-class Terraform resource that aligns 1:1 with this new API so they can provision EMR clusters via IaC using the current product capabilities (new scene/software model, multi-AZ, TF-specific node marks, etc.).

## What Changes

- **NEW** Terraform resource `tencentcloud_emr_cluster_v2` backed by the `CreateCluster` / `DescribeInstances` / `DescribeClusterNodes` / `TerminateInstance` APIs.
- Schema fields strictly mirror the `CreateCluster` request parameters (no renames, no merges) including nested blocks for `LoginSettings`, `SceneSoftwareConfig`, `InstanceChargePrepaid`, `ScriptBootstrapActionConfig`, `Tags`, `MetaDBInfo`, `DependService`, `ZoneResourceConfiguration` (with `VirtualPrivateCloud` / `Placement` / `AllNodeResourceSpec` / `NodeResourceSpec` / `DiskSpecInfo`), and `NodeMarks`.
- Resource unique ID is the `InstanceId` returned by `CreateCluster` (e.g. `emr-f2da1cd`).
- Create waits for `DescribeInstances.Clusters[0].Status == 2` (running) before returning, using `resource.Retry` with `schema.Timeouts` (create timeout default 60 min).
- Read populates state from both `DescribeInstances` (cluster-level info) and `DescribeClusterNodes` (node-level info).
- Delete calls `TerminateInstance` and waits until the cluster is no longer returned by `DescribeInstances`.
- Update CRUD is intentionally **not implemented** in this change — all schema fields are declared without `ForceNew` (per requirements, modify APIs will be added in a follow-up change). The current Update handler is a no-op that only re-reads state.
- New files under `tencentcloud/services/emr/`:
  - `resource_tc_emr_cluster_v2.go`
  - `resource_tc_emr_cluster_v2.md`
  - `resource_tc_emr_cluster_v2_test.go`
- Service-layer additions in `service_tencentcloud_emr.go`: `DescribeEmrClusterV2ById`, `DescribeEmrClusterV2Nodes`, `EmrClusterV2StateRefreshFunc`.
- Register resource in `tencentcloud/provider.go` ResourcesMap.

## Capabilities

### New Capabilities

- `emr-cluster-v2-resource`: Defines the lifecycle contract and field mapping for the new `tencentcloud_emr_cluster_v2` Terraform resource (Create/Read/Delete using the `CreateCluster` family of APIs).

### Modified Capabilities

_None._ The existing `tencentcloud_emr_cluster` resource is untouched; `tencentcloud_emr_cluster_v2` is an additive, independent resource.

## Impact

- **Code**:
  - New: `tencentcloud/services/emr/resource_tc_emr_cluster_v2.go`, `.md`, `_test.go`
  - Modified: `tencentcloud/services/emr/service_tencentcloud_emr.go` (append helpers only)
  - Modified: `tencentcloud/provider.go` (register `tencentcloud_emr_cluster_v2`)
  - Generated: `website/docs/r/emr_cluster_v2.html.markdown` (via `make doc`)
- **APIs**: Consumes EMR v20190103 SDK (already vendored). No new vendor dependencies.
- **Dependencies**: None.
- **Backward compatibility**: Purely additive. Existing `tencentcloud_emr_cluster` state and configurations are unaffected.
- **Breaking changes**: None.
- **Follow-up**: A subsequent change will add the `Update` lifecycle (scale-out/scale-in, ModifyResourcesTags, etc.) and, if needed, a companion data source.
