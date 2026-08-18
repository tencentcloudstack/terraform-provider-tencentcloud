## Why

CDB (Cloud Database MySQL) currently has no Terraform resource to clone (clone/rollback) an existing instance into a new one. The cloud API provides `CreateCloneInstance` (async) to create a clone instance from a source instance, optionally rolling back to a specified time or backup set. Users need to manage the full lifecycle of cloned CDB instances through Terraform to ensure consistent, reproducible database provisioning and to avoid manual console operations that drift from desired state.

## What Changes

- Add a new Terraform RESOURCE_KIND_GENERAL resource `tencentcloud_cdb_clone_instance` to manage the full CRUD lifecycle of cloned CDB instances
- Create resource file: `tencentcloud/services/cdb/resource_tc_cdb_clone_instance.go`
- Create test file: `tencentcloud/services/cdb/resource_tc_cdb_clone_instance_test.go`
- Create documentation file: `tencentcloud/services/cdb/resource_tc_cdb_clone_instance.md`
- Register the resource in `tencentcloud/provider.go` and `tencentcloud/provider.md`
- The resource will support CRUD operations:
  - **Create**: Call `CreateCloneInstance` API (async) to clone a CDB instance; poll `DescribeAsyncRequestInfo` until the async task succeeds, then extract the cloned instance ID from `DescribeDBInstances`
  - **Read**: Call `DescribeDBInstances` API to query the cloned instance details and populate schema fields
  - **Update**: Call `UpgradeDBInstance` API (async) to adjust instance configuration (memory/volume/cpu/etc.); poll `DescribeAsyncRequestInfo` until the upgrade task succeeds
  - **Delete**: Call `OfflineIsolatedInstances` API to offline (destroy) the isolated cloned instance

## Capabilities

### New Capabilities
- `cdb-clone-instance`: Manages the full lifecycle of a cloned CDB instance — creation from a source instance with optional rollback time/backup, configuration upgrades, querying instance details, and deletion via offline

### Modified Capabilities
<!-- No existing capabilities are modified -->

## Impact

- New resource registration in `tencentcloud/provider.go` and `tencentcloud/provider.md`
- New resource implementation in `tencentcloud/services/cdb/resource_tc_cdb_clone_instance.go`
- New test file in `tencentcloud/services/cdb/resource_tc_cdb_clone_instance_test.go`
- New documentation in `tencentcloud/services/cdb/resource_tc_cdb_clone_instance.md` and `website/docs/r/` (generated via `make doc` in finalize phase)
- Cloud API dependencies from `cdb/v20170320` SDK package: `CreateCloneInstance` (async), `DescribeDBInstances`, `UpgradeDBInstance` (async), `OfflineIsolatedInstances`, `DescribeAsyncRequestInfo` (for async task polling)
