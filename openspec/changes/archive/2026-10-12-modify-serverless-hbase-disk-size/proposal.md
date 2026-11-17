## Why

The `tencentcloud_serverless_hbase_instance` resource currently treats the `disk_size` field as immutable (any change returns an error), yet the underlying cloud API `ModifySLInstance` (EMR v20190103) already supports storage expansion via its `DiskSize` parameter ("存储空间大小，需要是100的倍数，且不允许比当前存储空间小（只允许扩容）"). Users therefore cannot expand the disk capacity of an existing Serverless HBase instance through Terraform and are forced to destroy and recreate it, causing data loss and downtime. By wiring `disk_size` changes to `ModifySLInstance` and waiting for the async scaling to finish (Status `22` 变配中 → `2` 运行中), the resource gains a critical in-place expansion capability.

## What Changes

- Modify the `Update` function of `tencentcloud_serverless_hbase_instance` (`tencentcloud/services/emr/resource_tc_serverless_hbase_instance.go`):
  - Remove `disk_size` from the `immutableArgs` block so it becomes a mutable, updatable argument.
  - Add a new change branch: when `d.HasChange("disk_size")`, construct a `ModifySLInstanceRequest` with only `InstanceId` and `DiskSize`, call `ModifySLInstanceWithContext` (with `tccommon.WriteRetryTimeout` retry), then poll the instance status until it transitions from `22` (变配中) to `2` (运行中, 变配完成) using `SLInstanceStateRefreshFunc`.
- Keep the existing `node_num` expansion logic (which also uses `ModifySLInstance`) unchanged; ensure `disk_size` changes are handled independently and do not conflict with `zone_settings.node_num` changes.
- Add/adjust unit tests in `resource_tc_serverless_hbase_instance_test.go` to cover the new `disk_size` update path (using gomonkey mocks for the cloud API client), per project convention for modified resources.
- No schema shape changes: `disk_size` remains `TypeInt, Required`, only its lifecycle behavior changes from immutable to updatable.

## Capabilities

### New Capabilities
- `serverless-hbase-instance-disk-size-update`: Enable in-place expansion of `disk_size` for the `tencentcloud_serverless_hbase_instance` resource via the async `ModifySLInstance` API, polling `DescribeSLInstance` (through the existing state refresh helper) until the instance returns to running state.

### Modified Capabilities
<!-- No existing spec covers the serverless hbase instance resource; this is a new capability. -->

## Impact

- **Code**: `tencentcloud/services/emr/resource_tc_serverless_hbase_instance.go` (Update function — move `disk_size` out of `immutableArgs`, add `ModifySLInstance` async branch + status polling).
- **Tests**: `tencentcloud/services/emr/resource_tc_serverless_hbase_instance_test.go` (add unit tests for the disk expansion update path).
- **Cloud APIs**: `ModifySLInstance` (EMR v20190103, https://cloud.tencent.com/document/api/589/110869) for the change; status queried via `DescribeSLInstance`/`DescribeSLInstanceList` flow already used by `SLInstanceStateRefreshFunc` (https://cloud.tencent.com/document/product/589/110871).
- **Backward compatibility**: Fully backward compatible. Existing configurations that do not change `disk_size` are unaffected; `disk_size` was already a `Required` field with a value, and making it mutable does not alter the schema or state shape.
- **Dependencies**: No new SDK methods required — `ModifySLInstanceWithContext`, `BuildStateChangeConf`, and `SLInstanceStateRefreshFunc` are already used by the resource.