## Context

The `tencentcloud_dbdc_db_custom_node` resource manages a single Tencent Cloud
DBDC (DB Custom) node. Its CRUD was introduced by the archived change
`2026-06-26-add-dbdc-db-custom-node`. The Read function populates fields from
the `DBCustomNode` struct returned by `DescribeDBCustomNodes` (via the service
helper `DescribeDBCustomNodeById`).

**Current state:**
- Resource file: `tencentcloud/services/dbdc/resource_tc_dbdc_db_custom_node.go`
- Service layer: `tencentcloud/services/dbdc/service_tencentcloud_dbdc.go`
  (`DescribeDBCustomNodeById` returns `*dbdcv20201029.DBCustomNode`)
- SDK: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbdc/v20201029`

**API behavior analysis (vendored SDK):**

The `DBCustomNode` struct (element type of
`DescribeDBCustomNodesResponse.NodeSet`) already contains:

| Field (SDK) | Type | Description |
|---|---|---|
| `NetworkMode` | `*string` | Network mode. `NetworkModePrivateLink` (four-layer SSH service connectivity) / `NetworkModeCrossTenantENI` (three-layer dual-NIC access) |
| `EniIP` | `*string` | The node access IP address when `NetworkModeCrossTenantENI` mode is selected |

These two fields are **output-only**. They are not present in
`CreateDBCustomNodesRequest`, `ModifyDBCustomNodeTagsRequest`, or
`RenewDBCustomNodeRequest`, so they cannot be written and are modeled as
Computed.

## Goals / Non-Goals

**Goals:**
- Add `network_mode` (TypeString, Computed) and `eni_ip` (TypeString, Computed)
  to the `tencentcloud_dbdc_db_custom_node` schema.
- Populate both fields in `resourceTencentCloudDbdcDbCustomNodeRead` from the
  `DescribeDBCustomNodeById` response, using the existing nil-guard pattern
  (`if respData.<Field> != nil { _ = d.Set(...) }`).
- Maintain full backward compatibility (Computed-only additions).

**Non-Goals:**
- Making `network_mode` or `eni_ip` user-configurable (the API does not accept
  them as input).
- Adding these fields to the `tencentcloud_dbdc_db_custom_nodes` data source
  (out of scope for this change).

## Decisions

### Decision 1: Model both fields as Computed-only

**Rationale:** `NetworkMode` and `EniIP` are returned by `DescribeDBCustomNodes`
but are not accepted by any Create/Update API. Modeling them as Computed
(Terraform `Computed: true`, no `Optional`) is the only correct representation.
No `ForceNew` is needed since they are not settable.

### Decision 2: Reuse the existing Read path and nil-guard convention

**Rationale:** The Read function already calls
`service.DescribeDBCustomNodeById(ctx, nodeId)` and sets each field with a
`if respData.X != nil` guard. The two new fields follow the identical pattern,
placed alongside the other computed string fields (e.g. after `isolated_time`
and before `system_disk`), keeping the diff minimal and reviewable.

### Decision 3: No SDK upgrade required

**Rationale:** The vendored SDK (upgraded to v1.3.149 in the prior module step)
already defines `NetworkMode` and `EniIP` on `DBCustomNode`. No further
`go.mod`/vendor changes are needed for this change.

## Risks / Trade-offs

- **[Risk] `EniIP` may be empty for non-CrossTenantENI nodes:** The API
  documentation states `EniIP` is the access IP only when
  `NetworkModeCrossTenantENI` is selected; it may be nil/empty otherwise.
  - **Mitigation:** The nil guard (`if respData.EniIP != nil`) prevents setting
    an empty value; the field simply stays unset in state when the API omits
    it, which is the correct behavior for a Computed field.
