## Context

The `tencentcloud_dbdc_db_custom_nodes` datasource wraps the `DescribeDBCustomNodes` API (vendored SDK package `dbdc/v20201029`). The `DBCustomNode` struct returned in `Response.NodeSet` already includes two fields that are not yet surfaced in the Terraform schema:

- `NetworkMode *string` — network mode enum: `NetworkModePrivateLink` (four-layer SSH service) or `NetworkModeCrossTenantENI` (three-layer dual-NIC access).
- `EniIP *string` — the node access IP address when `NetworkMode` is `NetworkModeCrossTenantENI`.

The datasource already flattens the other `DBCustomNode` fields (e.g. `host_ip`, `rack_id`, `switch_id`) into a `node_set` list of `schema.Resource`. This change simply adds the two missing fields following the exact same pattern.

## Goals / Non-Goals

**Goals:**
- Surface `NetworkMode` and `EniIP` from `DescribeDBCustomNodes` responses as computed attributes (`network_mode`, `eni_ip`) inside each `node_set` element.
- Maintain full backward compatibility (only new Computed fields are added; no existing schema field is altered or removed).

**Non-Goals:**
- No new input/filter parameters are added.
- No new API calls or SDK upgrade.
- No changes to the service layer pagination logic.

## Decisions

### Decision 1: Add fields as Computed inside `node_set` element schema
Both `NetworkMode` and `EniIP` are read-only attributes returned by the API, so they are `Computed: true` string fields nested under the existing `node_set` element resource — identical to how `host_ip`, `switch_id`, and `rack_id` are defined today.

**Alternative considered**: Exposing them at the top level of the datasource schema. Rejected — these are per-node attributes, so they belong inside `node_set` alongside the other per-node fields.

### Decision 2: Use snake_case Terraform names `network_mode` and `eni_ip`
Terraform convention is snake_case. The cloud API uses `NetworkMode` and `EniIP`; these map to `network_mode` and `eni_ip` respectively, consistent with existing field naming (e.g. `lan_ip` ← `LanIP`, `host_ip` ← `HostIp`).

### Decision 3: Nil-guard before `d.Set`
Following the existing convention in the Read function, each new field is only added to the `nodeMap` when the corresponding API pointer is non-nil. This matches the "Nil field handling in response" requirement already in the spec.

## Risks / Trade-offs

- **[Risk] Fields may be null for some nodes** → Mitigation: nil-guard each field before adding to the output map, consistent with the existing `SystemDisk`/`DataDisks`/`Tags` handling. The API doc notes `EniIP` is only populated when `NetworkMode` is `NetworkModeCrossTenantENI`.
- **[Risk] Breaking existing state** → Mitigation: Only Computed fields are added; Terraform treats new Computed attributes as backward compatible.
