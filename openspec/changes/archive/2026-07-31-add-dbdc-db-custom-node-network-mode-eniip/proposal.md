## Why

The `tencentcloud_dbdc_db_custom_node` resource currently does not expose the
node's network mode (`NetworkMode`) or the ENI access IP address (`EniIP`).
These two fields are returned by the `DescribeDBCustomNodes` API (in the
`DBCustomNode` struct of `response.NodeSet`), but are neither declared in the
resource schema nor populated in the Read function. Users therefore cannot
discover which network access mode a node is using, nor the access IP when the
node uses the `NetworkModeCrossTenantENI` (three-layer dual-NIC) mode, through
Terraform state.

## What Changes

- Add a new computed parameter `network_mode` (TypeString, Computed) to the
  `tencentcloud_dbdc_db_custom_node` resource. It is refreshed from
  `DBCustomNode.NetworkMode` returned by `DescribeDBCustomNodes`. Valid values:
  `NetworkModePrivateLink` (four-layer SSH service connectivity mode) and
  `NetworkModeCrossTenantENI` (three-layer dual-NIC access mode).
- Add a new computed parameter `eni_ip` (TypeString, Computed) to the
  `tencentcloud_dbdc_db_custom_node` resource. It is refreshed from
  `DBCustomNode.EniIP` returned by `DescribeDBCustomNodes`, and is the access IP
  address of the node when the `NetworkModeCrossTenantENI` mode is selected.
- Wire both fields into the Read function with the standard nil-guard pattern
  (`if respData.NetworkMode != nil { _ = d.Set("network_mode", ...) }`).
- Both fields are read-only output parameters of `DescribeDBCustomNodes`; they
  are not accepted by `CreateDBCustomNodes`, `ModifyDBCustomNodeTags`, or
  `RenewDBCustomNode`, so they are purely computed and require no Create/Update
  wiring.

## Capabilities

### New Capabilities
- `dbdc-db-custom-node-network-mode-eniip`: Expose the `network_mode` and
  `eni_ip` computed parameters on the `tencentcloud_dbdc_db_custom_node`
  resource, refreshed from the `DescribeDBCustomNodes` API response.

### Modified Capabilities
<!-- No existing specs require modification -->

## Impact

- **Affected files:**
  - `tencentcloud/services/dbdc/resource_tc_dbdc_db_custom_node.go` — add
    `network_mode` and `eni_ip` schema fields (Computed), and populate them in
    `resourceTencentCloudDbdcDbCustomNodeRead` with nil guards.
  - `tencentcloud/services/dbdc/resource_tc_dbdc_db_custom_node_test.go` — add
    unit test coverage for the two new read fields.
- **SDK dependency:** No SDK upgrade required. The vendored
  `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbdc/v20201029`
  already defines `NetworkMode` and `EniIP` on the `DBCustomNode` struct.
- **API constraints:** `NetworkMode` and `EniIP` are only available as output
  parameters of `DescribeDBCustomNodes`; they cannot be written through any
  Create/Update API, so they are modeled as Computed only.
- **Backward compatibility:** Fully backward compatible — both parameters are
  Optional-not-set (Computed) additions; existing configurations and state are
  unaffected.
