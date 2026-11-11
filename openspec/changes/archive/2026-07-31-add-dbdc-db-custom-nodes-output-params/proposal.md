## Why

The `tencentcloud_dbdc_db_custom_nodes` datasource currently does not expose the `NetworkMode` and `EniIP` fields that the `DescribeDBCustomNodes` API returns. Users need visibility into the network mode of each DB Custom node (e.g., `NetworkModePrivateLink` vs `NetworkModeCrossTenantENI`) and, when the cross-tenant ENI mode is selected, the corresponding access IP address. The cloud API already returns these fields in the `DBCustomNode` struct, so the Terraform datasource should surface them.

## What Changes

- Add two new computed output fields to the `node_set` element schema of `tencentcloud_dbdc_db_custom_nodes`:
  - `network_mode` (TypeString, Computed) — network mode of the node, mapped from `DBCustomNode.NetworkMode`
  - `eni_ip` (TypeString, Computed) — access IP address when the network mode is `NetworkModeCrossTenantENI`, mapped from `DBCustomNode.EniIP`
- Update the Read function to populate these fields from the API response with nil guards (consistent with existing fields).
- Update the unit tests to cover the new fields.
- Update the `.md` example doc so the new output attributes are documented (auto-generated Argument/Attribute Reference excluded per convention).

## Capabilities

### New Capabilities
<!-- None -->

### Modified Capabilities
- `dbdc-db-custom-nodes-datasource`: Add `network_mode` and `eni_ip` computed output fields to the `node_set` schema and Read logic so users can query the network mode and ENI access IP of DB Custom nodes.

## Impact

- **Code**: `tencentcloud/services/dbdc/data_source_tc_dbdc_db_custom_nodes.go` (schema + Read), `tencentcloud/services/dbdc/data_source_tc_dbdc_db_custom_nodes_test.go` (tests), `tencentcloud/services/dbdc/data_source_tc_dbdc_db_custom_nodes.md` (doc).
- **APIs**: No new API calls; the existing `DescribeDBCustomNodes` response already contains `NetworkMode` and `EniIP` in the `DBCustomNode` struct (vendored SDK `dbdc/v20201029`).
- **Dependencies**: None — no SDK upgrade required.
- **Backward compatibility**: Fully backward compatible; only new Computed fields are added to the `node_set` element.
