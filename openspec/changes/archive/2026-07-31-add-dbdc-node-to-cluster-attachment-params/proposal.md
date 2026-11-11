## Why

The `tencentcloud_dbdc_node_to_db_custom_cluster_attachment` resource currently only supports `cluster_id`, `node_id`, `image_id`, and `login_settings` arguments. The underlying `AddNodesToDBCustomCluster` API also accepts `Labels`, `Taints`, `HostName`, and `HostNameType` parameters, and `DescribeDBCustomNodes` returns `NetworkMode` and `EniIP` fields, but these are not exposed. Users cannot declaratively configure node labels/taints/hostnames at bind time, nor read the node's network mode and ENI IP, limiting the resource's usefulness for DB Custom cluster node management.

## What Changes

- Add new optional `ForceNew` schema arguments to the `tencentcloud_dbdc_node_to_db_custom_cluster_attachment` resource:
  - `labels` — TypeList (MaxItems 20) of schema.Resource with `key` (Required) and `value` (Optional) fields; maps to `AddNodesToDBCustomClusterRequest.Labels`
  - `taints` — TypeList (MaxItems 5) of schema.Resource with `key` (Required), `effect` (Required), and `value` (Optional) fields; maps to `AddNodesToDBCustomClusterRequest.Taints`
  - `host_name` — Optional TypeString; maps to `AddNodesToDBCustomClusterRequest.HostName`
  - `host_name_type` — Optional TypeInt; maps to `AddNodesToDBCustomClusterRequest.HostNameType`
- Add new computed schema fields:
  - `network_mode` — Computed TypeString; read from `DescribeDBCustomNodes` response `NodeSet[].NetworkMode`
  - `eni_ip` — Computed TypeString; read from `DescribeDBCustomNodes` response `NodeSet[].EniIP`
- Populate the new input parameters in the Create function when calling `AddNodesToDBCustomCluster`
- Populate the new computed fields in the Read function by calling `DescribeDBCustomNodes` (which returns `NetworkMode` and `EniIP` on the `DBCustomNode` struct)
- Pass `login_settings` to the Delete function when calling `RemoveNodesFromDBCustomCluster` (the API accepts `LoginSettings` as an input parameter on remove)
- Update the resource documentation `.md` file with the new arguments in the example usage
- Add/extend unit tests in the `_test.go` file using gomonkey mocks

## Capabilities

### New Capabilities
- `dbdc-node-to-db-custom-cluster-attachment-params`: Adds labels, taints, host_name, host_name_type input parameters and network_mode, eni_ip computed fields to the dbdc node-to-cluster attachment resource, plus login_settings support on delete.

### Modified Capabilities
<!-- No existing spec-level requirements are being modified. The original attachment resource had no spec file. -->

## Impact

- **Affected code**:
  - `tencentcloud/services/dbdc/resource_tc_dbdc_node_to_db_custom_cluster_attachment.go` — schema additions, Create/Read/Delete logic updates
  - `tencentcloud/services/dbdc/resource_tc_dbdc_node_to_db_custom_cluster_attachment_test.go` — test additions
  - `tencentcloud/services/dbdc/resource_tc_dbdc_node_to_db_custom_cluster_attachment.md` — doc update
- **APIs used** (all already present in vendored SDK `tencentcloud-sdk-go/tencentcloud/dbdc/v20201029`):
  - `AddNodesToDBCustomCluster` — new input params: `Labels`, `Taints`, `HostName`, `HostNameType`
  - `DescribeDBCustomNodes` — new output fields: `NetworkMode`, `EniIP` (on `DBCustomNode`)
  - `RemoveNodesFromDBCustomCluster` — pass `LoginSettings` on delete
- **Dependencies**: No SDK upgrade required; all fields are already in the vendored `models.go`.
- **Backward compatibility**: All new arguments are Optional (ForceNew) or Computed — no breaking changes to existing TF configs or state.
