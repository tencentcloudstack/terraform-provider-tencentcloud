## Why

The `tencentcloud_dbdc_db_custom_node` resource currently exposes only a subset of the `CreateDBCustomNodes` API parameters. The API supports specifying the charge type, network mode, system/data disk configuration, host name, and security groups at node creation time, but none of these are configurable through Terraform. Users who need cloud-disk node types (e.g. `DB.SA5`), custom host names, or bound security groups must fall back to the console or raw API calls, breaking the declarative workflow.

## What Changes

- Add `charge_type` as an **Optional + Computed, ForceNew** parameter (promote the existing Computed field) mapped to `CreateDBCustomNodesRequest.ChargeType`. Valid values: `PREPAID`, `POSTPAID`. Immutable after creation because no Modify API exists.
- Add `network_mode` as an **Optional, ForceNew** parameter mapped to `CreateDBCustomNodesRequest.NetworkMode`. Valid values: `privatelink`, `cross_tenant_eni`. Immutable after creation.
- Promote the existing `system_disk` block from **Computed-only** to **Optional + Computed, ForceNew** so users can pass `disk_type`/`disk_size` (mapped to `CreateDBCustomNodesRequest.SystemDisk`). Immutable after creation.
- Promote the existing `data_disks` block from **Computed-only** to **Optional + Computed, ForceNew** so users can pass `disk_type`/`disk_size`/`disk_name` (mapped to `CreateDBCustomNodesRequest.DataDisks`). Immutable after creation.
- Add `host_name` as an **Optional, ForceNew** parameter mapped to `CreateDBCustomNodesRequest.HostName`. Immutable after creation; not returned by `DescribeDBCustomNodes`, so it is write-only (imported resources will not repopulate it).
- Add `security_group_ids` as an **Optional** (mutable, not ForceNew) parameter mapped to `CreateDBCustomNodesRequest.SecurityGroupIds` on Create and to the `ModifyDBCustomNodeSecurityGroups` API on Update. Read refreshes it via the separate `DescribeDBCustomNodeSecurityGroups` API because `DBCustomNode` does not return security groups.
- Wire the new arguments through the Create flow, add Read support (from `DescribeDBCustomNodes` where available, plus a new service helper for security groups), and add Update handling for `security_group_ids`.
- Update `resource_tc_dbdc_db_custom_node.md` example and add test coverage.

## Capabilities

### New Capabilities
- `dbdc-db-custom-node-create-params`: New optional creation parameters on the `tencentcloud_dbdc_db_custom_node` resource — charge type, network mode, system disk, data disks, host name, and security group ids — including mutable update support for security group ids.

### Modified Capabilities
<!-- No existing specs require modification. The dbdc-db-custom-node resource has no prior spec; this change creates the capability spec. -->

## Impact

- **Affected files:**
  - `tencentcloud/services/dbdc/resource_tc_dbdc_db_custom_node.go` — extend schema (promote `charge_type`, `system_disk`, `data_disks` to Optional+Computed; add `network_mode`, `host_name`, `security_group_ids`), wire Create, extend Read, extend Update for `security_group_ids`.
  - `tencentcloud/services/dbdc/service_tencentcloud_dbdc.go` — add `DescribeDBCustomNodeSecurityGroupsById(ctx, nodeId)` helper wrapping `DescribeDBCustomNodeSecurityGroups`.
  - `tencentcloud/services/dbdc/resource_tc_dbdc_db_custom_node.md` — update example usage.
  - `tencentcloud/services/dbdc/resource_tc_dbdc_db_custom_node_test.go` — extend acceptance tests for new parameters.
- **APIs used:** `CreateDBCustomNodes` (already used, new request fields), `DescribeDBCustomNodes` (already used, reads `ChargeType`/`NetworkMode`/`SystemDisk`/`DataDisks`), `DescribeDBCustomNodeSecurityGroups` (new, for `security_group_ids` Read), `ModifyDBCustomNodeSecurityGroups` (new, for `security_group_ids` Update).
- **SDK dependency:** No SDK upgrade required. All referenced request/response structs and client methods already exist in the vendored `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbdc/v20201029` package.
- **Backward compatibility:** Fully backward compatible. Every changed field is promoted from Computed to Optional+Computed or added as Optional, so existing configurations and state continue to work unchanged.
- **API constraints:**
  - `ChargeType`, `NetworkMode`, `SystemDisk`, `DataDisks`, `HostName` are accepted only by `CreateDBCustomNodes`; no Modify API exists, so they are `ForceNew` (rebuild on change).
  - `HostName` is not returned by `DescribeDBCustomNodes`, so it cannot be refreshed on Read/import (write-only).
  - `SecurityGroupIds` is accepted by `CreateDBCustomNodes` and updatable via `ModifyDBCustomNodeSecurityGroups`; it is not part of the `DBCustomNode` response, so Read uses `DescribeDBCustomNodeSecurityGroups`.
