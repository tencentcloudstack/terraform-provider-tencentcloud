## Why

CynosDB (TDSQL-C) clusters support read-only instance groups, but currently there is no Terraform resource to enable (open) access to these groups. Users need a one-shot operation resource to programmatically open read-only instance group access for their CynosDB clusters via Terraform.

## What Changes

- Add a new `RESOURCE_KIND_OPERATION` resource `tencentcloud_cynosdb_cluster_read_only_instance_group_acces_operation` that calls the `OpenClusterReadOnlyInstanceGroupAccess` API to enable read-only instance group access.
- The resource is a one-shot operation: it only has a Create implementation (no Read/Update/Delete).
- The Create operation is asynchronous (returns a `FlowId`), so it must poll `DescribeFlow` until the operation completes.
- Input parameters: `cluster_id`, `port`, `security_group_ids`.
- Output (computed) parameter: `flow_id`.

## Capabilities

### New Capabilities
- `cynosdb-cluster-read-only-instance-group-access-operation`: One-shot operation resource to open read-only instance group access for a CynosDB cluster, wrapping the `OpenClusterReadOnlyInstanceGroupAccess` API with async flow polling.

### Modified Capabilities

(none)

## Impact

- New files:
  - `tencentcloud/services/cynosdb/resource_tc_cynosdb_cluster_read_only_instance_group_acces_operation.go`
  - `tencentcloud/services/cynosdb/resource_tc_cynosdb_cluster_read_only_instance_group_acces_operation_test.go`
  - `tencentcloud/services/cynosdb/resource_tc_cynosdb_cluster_read_only_instance_group_acces_operation.md`
- Modified files:
  - `tencentcloud/provider.go` (register the new resource)
  - `tencentcloud/provider.md` (add resource to documentation list)
- Dependencies: Uses existing `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107` SDK package (already vendored).
