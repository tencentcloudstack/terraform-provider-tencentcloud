## 1. Service Layer

- [x] 1.1 Add `DescribeDBCustomNodeSecurityGroupsById(ctx, nodeId)` to `tencentcloud/services/dbdc/service_tencentcloud_dbdc.go` — wraps `DescribeDBCustomNodeSecurityGroups` (`NodeId`, `resource.Retry` + `ratelimit.Check`), returns `[]string` of `SecurityGroupId` from `Response.Groups` (nil/length-safe)
- [x] 1.2 (No other service-layer changes needed; `DescribeDBCustomNodeById` already returns `*DBCustomNode` carrying `ChargeType`, `NetworkMode`, `SystemDisk`, `DataDisks`)

## 2. Schema Changes

- [x] 2.1 Promote `charge_type` from Computed to Optional+Computed, ForceNew in `ResourceTencentCloudDbdcDbCustomNode()` (keep description; note valid values `PREPAID`/`POSTPAID`)
- [x] 2.2 Add `network_mode` (TypeString, Optional, ForceNew) — valid values `privatelink`, `cross_tenant_eni`
- [x] 2.3 Promote `system_disk` block from Computed to Optional+Computed, ForceNew (MaxItems 1; keep `disk_type`/`disk_size` sub-fields)
- [x] 2.4 Promote `data_disks` block from Computed to Optional+Computed, ForceNew (keep `disk_type`/`disk_size`/`disk_name` sub-fields)
- [x] 2.5 Add `host_name` (TypeString, Optional, ForceNew; NOT Computed — write-only)
- [x] 2.6 Add `security_group_ids` (TypeList of String, Optional; mutable — NOT ForceNew)

## 3. Create Function Changes

- [x] 3.1 Pass `charge_type` to `request.ChargeType` when set (`d.GetOk`)
- [x] 3.2 Pass `network_mode` to `request.NetworkMode` when set (`d.GetOk`)
- [x] 3.3 Build `request.SystemDisk` (`DiskType`, `DiskSize`) from the `system_disk` block when present (`helper.InterfacesHeadMap`)
- [x] 3.4 Build `request.DataDisks` (`[]*DataDisk{DiskType, DiskSize, DiskName}`) from the `data_disks` list when present
- [x] 3.5 Pass `host_name` to `request.HostName` when set (`d.GetOk`)
- [x] 3.6 Pass `security_group_ids` to `request.SecurityGroupIds` (`[]*string`) when set

## 4. Read Function Changes

- [x] 4.1 Set `charge_type` from `respData.ChargeType` (existing nil-guard) — already present, keep
- [x] 4.2 Add setting `network_mode` from `respData.NetworkMode` with nil-guard
- [x] 4.3 Keep existing `system_disk` read (from `respData.SystemDisk`) — now also reflects user input on refresh
- [x] 4.4 Keep existing `data_disks` read (from `respData.DataDisks`) — now also reflects user input on refresh
- [x] 4.5 Do NOT read `host_name` from the API (write-only); leave configured value in state untouched
- [x] 4.6 Read `security_group_ids` via the new `DescribeDBCustomNodeSecurityGroupsById` service helper; set into state with nil-guard (empty list when no groups)

## 5. Update Function Changes

- [x] 5.1 Add `d.HasChange("security_group_ids")` branch: call `ModifyDBCustomNodeSecurityGroups` with `NodeId` + full new `SecurityGroupIds` list, wrapped in `resource.Retry(WriteRetryTimeout)` + `tccommon.RetryError`, nil-guard the response
- [x] 5.2 Confirm `charge_type`, `network_mode`, `system_disk`, `data_disks`, `host_name` are ForceNew (no Update handling required — Terraform recreates)

## 6. Documentation

- [x] 6.1 Update `tencentcloud/services/dbdc/resource_tc_dbdc_db_custom_node.md` example usage to illustrate the new parameters (`charge_type`, `network_mode`, `system_disk`, `data_disks`, `host_name`, `security_group_ids`)

## 7. Tests

- [x] 7.1 Extend `tencentcloud/services/dbdc/resource_tc_dbdc_db_custom_node_test.go` acceptance tests to cover the new parameters (basic create with new args, update of `security_group_ids`, import with `ImportStateVerifyIgnore` for write-only `host_name`)

## 8. Verification

- [x] 8.1 Verify the code compiles successfully
- [x] 8.2 Verify no lint errors
