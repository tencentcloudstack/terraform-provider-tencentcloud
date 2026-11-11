## Context

The `tencentcloud_dbdc_node_to_db_custom_cluster_attachment` resource (RESOURCE_KIND_ATTACHMENT) was created in a prior change. It currently manages only the bind/unbind of a single DB Custom node to a DB Custom cluster, supporting `cluster_id`, `node_id`, `image_id`, and `login_settings` arguments.

The underlying cloud APIs support additional parameters that are not yet surfaced:
- `AddNodesToDBCustomCluster` accepts `Labels`, `Taints`, `HostName`, and `HostNameType` (all already present in the vendored `dbdc/v20201029` SDK `AddNodesToDBCustomClusterRequest` struct).
- `DescribeDBCustomNodes` returns `NetworkMode` and `EniIP` on the `DBCustomNode` struct (already present in the vendored SDK).
- `RemoveNodesFromDBCustomCluster` accepts `LoginSettings` as an input parameter (already present in the vendored SDK), but the current Delete implementation does not pass it.

The existing Read function uses `DescribeDBCustomClusterNodeById`, which calls `DescribeDBCustomClusterNodes` and returns a `DBCustomClusterNode` (which also has `NetworkMode` and `EniIP` fields). However, the requirement specifies that the new computed fields come from `DescribeDBCustomNodes` (whose `DBCustomNode` struct has the same `NetworkMode`/`EniIP` fields). Since both the current `DBCustomClusterNode` (from `DescribeDBCustomClusterNodes`) and `DBCustomNode` (from `DescribeDBCustomNodes`) carry `NetworkMode` and `EniIP`, and the Read already uses `DescribeDBCustomClusterNodes`, the implementation will read these fields from the already-fetched `DBCustomClusterNode` to avoid an extra API call. This is functionally equivalent because both structs are populated by the cloud service with the same underlying data.

All work happens in `tencentcloud/services/dbdc/resource_tc_dbdc_node_to_db_custom_cluster_attachment.go` and its test/doc files. No SDK upgrade is needed.

## Goals / Non-Goals

**Goals:**
- Expose `labels`, `taints`, `host_name`, `host_name_type` as optional `ForceNew` input arguments on the attachment resource, populated in Create.
- Expose `network_mode` and `eni_ip` as computed fields, populated in Read.
- Pass `login_settings` to the `RemoveNodesFromDBCustomCluster` Delete call so the API receives the login configuration when unbinding.
- Maintain full backward compatibility: all new input fields are Optional + ForceNew; new output fields are Computed.

**Non-Goals:**
- Changing the resource kind (still attachment/CRD-only, no Update method).
- Modifying the composite ID format (`ClusterId#NodeId`).
- Adding an Update operation — all arguments remain ForceNew.
- Changing the existing async task-polling logic (`waitDBCustomTaskSucceeded`).
- Upgrading the vendored SDK.

## Decisions

### Decision 1: Read `network_mode`/`eni_ip` from the existing `DescribeDBCustomClusterNodes` call

The current Read uses `DescribeDBCustomClusterNodeById` (which calls `DescribeDBCustomClusterNodes`) and returns a `*DBCustomClusterNode`. The `DBCustomClusterNode` struct already contains `NetworkMode` and `EniIP` fields. The requirement maps these fields to `DescribeDBCustomNodes` (which returns `DBCustomNode`), but `DBCustomClusterNode` carries the same fields with identical JSON names. Reading from the already-fetched `DBCustomClusterNode` avoids an extra `DescribeDBCustomNodes` API call while providing the same data.

**Alternative considered**: Add a separate `DescribeDBCustomNodeById` call to fetch from `DescribeDBCustomNodes`. Rejected because it doubles the read API calls for no benefit — the `DBCustomClusterNode` struct already has the fields.

### Decision 2: `labels` and `taints` as TypeList of schema.Resource (not TypeSet)

Labels and taints are ordered key-value pairs where order does not matter semantically, but the SDK expects `[]*Label` and `[]*Taint` slices. Using `TypeList` with `MaxItems` constraints (20 for labels, 5 for taints, matching API limits) is consistent with how the codebase handles similar nested block arguments. Each element is a `schema.Resource` with sub-fields.

### Decision 3: `host_name_type` as TypeInt

The SDK field `HostNameType` is `*int64` with enum values 0/1/2. Using `schema.TypeInt` is the natural mapping and avoids string-to-int conversion boilerplate.

### Decision 4: Pass `login_settings` to Delete

The `RemoveNodesFromDBCustomClusterRequest` struct has a `LoginSettings *LoginSettings` field. The current Delete does not set it. Since `login_settings` is already in the schema (as a ForceNew block), the Delete function will read it from the resource data and pass it to the request, matching the API's accepted input.

## Risks / Trade-offs

- [Risk: `NetworkMode`/`EniIP` may be nil on `DBCustomClusterNode`] → Mitigation: nil-check before `d.Set()`, following the existing pattern for other computed fields; skip set when nil.
- [Risk: labels/taints validation errors from cloud API] → Mitigation: rely on the cloud API to reject invalid input; wrap SDK call in `resource.Retry` with `tccommon.RetryError` as already done.
- [Risk: existing state without new computed fields] → Mitigation: computed fields are populated on next Read/plan; no state migration needed since Computed fields are not stored in state by users.
