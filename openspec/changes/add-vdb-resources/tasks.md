## 1. Infrastructure Setup

- [x] 1.1 Create `tencentcloud/services/vdb/` directory
- [x] 1.2 Register VDB client in `tencentcloud/connectivity/client.go`: add import for `vdb/v20230616`, add `vdbV20230616Conn` field, add `UseVdbV20230616Client()` method
- [x] 1.3 Add `vdb` service import in `tencentcloud/provider.go`

## 2. Service Layer

- [x] 2.1 Create `tencentcloud/services/vdb/service_tencentcloud_vdb.go` with `VdbClientInterface` interface, `VdbService` struct (supports mock injection via `vdbClient` field), `getVdbClient()` method, and helper methods: `DescribeVdbInstanceById` (calls DescribeInstances), `DescribeVdbInstanceNodesById` (calls DescribeInstanceNodes), `DescribeDBSecurityGroupsByInstanceId` (calls DescribeDBSecurityGroups)

## 3. VDB Instance Resource — general (CRUD: CreateInstance / DescribeInstances+DescribeInstanceNodes / ScaleUpInstance+ScaleOutInstance / IsolateInstance+DestroyInstances)

- [x] 3.1 Create `tencentcloud/services/vdb/resource_tc_vdb_instance.go`: schema definition with `force_delete` bool param (default false), computed `nodes` attribute (Pod name+status from DescribeInstanceNodes), Timeouts block, full CRUD. Delete function: if force_delete=false only call IsolateInstance; if force_delete=true call IsolateInstance then DestroyInstances
- [x] 3.2 Create `tencentcloud/services/vdb/resource_tc_vdb_instance.md` with two example HCL (pay-as-you-go single + monthly subscription cluster covering all parameters) and import instructions
- [x] 3.3 Create `tencentcloud/services/vdb/resource_tc_vdb_instance_test.go` with mock-based unit tests: MockVdbClient, service layer tests (Describe/WaitForStatus/WaitForScaleUp/WaitForScaleOut/WaitForSecurityGroupsMatch), schema validation tests (required/optional/computed/ForceNew fields, sub-fields)
- [x] 3.4 Register `tencentcloud_vdb_instance` in `provider.go` ResourcesMap

## 4. VDB Instance Maintenance Window Resource — REMOVED

*This section has been removed. The `tencentcloud_vdb_instance_maintenance_window` resource is no longer part of this change. Maintenance window APIs (`DescribeInstanceMaintenanceWindow`, `ModifyInstanceMaintenanceWindow`) are not exposed.*

## 5. VDB Security Group Attachment Resource — REMOVED

*This section has been removed. The `tencentcloud_vdb_security_group_attachment` resource is no longer part of this change. Security groups are now managed inline in `tencentcloud_vdb_instance` via `security_group_ids` (Required) with `ModifyDBInstanceSecurityGroups` API for updates.*

## 6. VDB Instance Recover Operation Resource — REMOVED

*This section has been removed. The `tencentcloud_vdb_instance_recover_operation` resource is no longer part of this change. `RecoverInstance` API is not exposed as a Terraform resource.*

## 7. Provider Documentation

- [x] 7.1 Update `tencentcloud/provider.md` to add VDB section with 1 resource listed
- [x] 7.2 Create `tencentcloud/services/vdb/resource_test.go` with package declaration

## 8. Verification

- [x] 8.1 Run `go build ./...` to verify compilation
- [x] 8.2 Run `make doc` to generate website documentation in `website/docs/r/` directory

## 9. Bugfix: VDB Instance Status Values

- [x] 9.1 Fix `resource_tc_vdb_instance.go` polling status: replace `"Running"` with `"online"` (VDB API returns `online` for running state, `creating` for creating state), use `strings.EqualFold` for case-insensitive comparison
- [x] 9.2 Fix `resource_tc_vdb_instance.go` delete polling: replace `"Isolated"` with `"isolated"`, use `strings.EqualFold`
- [x] 9.3 Add DEBUG log in all polling loops to print actual `Status` value for troubleshooting

## 10. Enhancement: Abstract Polling Helpers (Status + Value-Based)

- [x] 10.1 Add `WaitForInstanceStatus(ctx, instanceId, targetStatus, timeout) error` method to `service_tencentcloud_vdb.go`. Polls `DescribeVdbInstanceById` with `resource.Retry`, uses `strings.EqualFold` for comparison, logs actual status at DEBUG level each iteration
- [x] 10.2 Add `WaitForInstanceScaleUp(ctx, instanceId, targetCpu, targetMemory, targetDiskSize, timeout) error` method. Polls `DescribeVdbInstanceById` until CPU/Memory/Disk values match target
- [x] 10.3 Add `WaitForInstanceScaleOut(ctx, instanceId, targetReplicaNum, timeout) error` method. Polls `DescribeVdbInstanceById` until ReplicaNum matches target
- [x] 10.4 Add `WaitForInstanceNotFound(ctx, instanceId, timeout) error` method. Polls `DescribeVdbInstanceById` until it returns nil (instance no longer exists), used after `DestroyInstances` with `force_delete=true`
- [x] 10.5 Refactor `resource_tc_vdb_instance.go` Create to use `WaitForInstanceStatus(ctx, instanceId, "online", d.Timeout(schema.TimeoutCreate))`
- [x] 10.6 Refactor `resource_tc_vdb_instance.go` Update ScaleUp to use `WaitForInstanceScaleUp(ctx, instanceId, targetCpu, targetMemory, targetDiskSize, d.Timeout(schema.TimeoutUpdate))`
- [x] 10.7 Refactor `resource_tc_vdb_instance.go` Update ScaleOut to use `WaitForInstanceScaleOut(ctx, instanceId, targetReplicaNum, d.Timeout(schema.TimeoutUpdate))`
- [x] 10.8 Refactor `resource_tc_vdb_instance.go` Delete isolate wait to use `WaitForInstanceStatus(ctx, instanceId, "isolated", d.Timeout(schema.TimeoutDelete))`
- [x] 10.9 Refactor `resource_tc_vdb_instance.go` Delete with `force_delete=true`: after `DestroyInstances`, call `WaitForInstanceNotFound(ctx, instanceId, d.Timeout(schema.TimeoutDelete))` to poll until instance is completely gone

## 11. Enhancement: Complete SDK Parameter Coverage for vdb_instance

- [x] 11.1 Add missing Create input parameter to schema: `params` (Optional, string, ForceNew) — instance extra parameters via JSON
- [x] 11.2 Add missing Computed fields to schema: `region`, `zone`, `product`, `shard_num`, `api_version`, `extend`, `expired_at`, `is_no_expired`, `wan_address`, `isolate_at`, `task_status`
- [x] 11.3 Add missing Network sub-fields: `networks[].preserve_duration` (Computed, int), `networks[].expire_time` (Computed, string)
- [x] 11.4 Update Create function to pass `params` to `CreateInstanceRequest.Params`
- [x] 11.5 Update Read function to set all new computed fields from `InstanceInfo` response
- [x] 11.6 Update Read function to set `networks[].preserve_duration` and `networks[].expire_time`
- [x] 11.7 Update service `DescribeVdbInstanceNodesById` to set `Limit` param to avoid truncated results
- [x] 11.8 Update `resource_tc_vdb_instance.md` example HCL and `resource_tc_vdb_instance_test.go` to reflect new fields

## 12. Enhancement: Serial Scale Operations in Update

- [x] 12.1 Refactor `resource_tc_vdb_instance.go` Update function: when both ScaleUp (cpu/memory/disk_size changed) and ScaleOut (worker_node_num changed) are needed, execute strictly serially — ScaleUp + WaitForInstanceScaleUp first, then ScaleOut + WaitForInstanceScaleOut second

## 12a. Enhancement: Inline Security Group Update via ModifyDBInstanceSecurityGroups

- [x] 12a.1 Change `security_group_ids` schema from Optional+ForceNew to Required (not ForceNew) in `resource_tc_vdb_instance.go`
- [x] 12a.2 Add security group update logic in `resource_tc_vdb_instance.go` Update function: when `security_group_ids` changes, call `ModifyDBInstanceSecurityGroups` API with full target list, then call `WaitForSecurityGroupsMatch` to poll `DescribeDBSecurityGroups` until result matches target
- [x] 12a.3 Add `WaitForSecurityGroupsMatch(ctx, instanceId, targetSgIds, timeout)` polling helper to `service_tencentcloud_vdb.go`
- [x] 12a.4 Add `ModifyDBInstanceSecurityGroupsWithContext` to `VdbClientInterface` (now 9 methods total)
- [x] 12a.5 Remove `tencentcloud_vdb_security_group_attachment` resource files, provider.go registration, provider.md entry, and website docs
- [x] 12a.6 Add unit tests for `WaitForSecurityGroupsMatch` polling helper (2 new tests)

## 13. Enhancement: VdbClientInterface for Mock-Based Unit Testing

- [x] 13.1 Define `VdbClientInterface` in `service_tencentcloud_vdb.go` listing all 9 SDK client methods used by VdbService and resource CRUD functions (including `ModifyDBInstanceSecurityGroupsWithContext`)
- [x] 13.2 Add `vdbClient VdbClientInterface` field to `VdbService` struct, add `getVdbClient()` method that returns injected mock or real SDK client
- [x] 13.3 Refactor all 3 service methods (`DescribeVdbInstanceById`, `DescribeVdbInstanceNodesById`, `DescribeDBSecurityGroupsByInstanceId`) to use `me.getVdbClient()` instead of `me.client.UseVdbV20230616Client()`
- [x] 13.4 Create `MockVdbClient` using `testify/mock` in test files, implementing `VdbClientInterface`
- [x] 13.5 Write 25 mock-based unit tests covering: service layer (DescribeInstance/Nodes/SecurityGroups), polling helpers (WaitForStatus/ScaleUp/ScaleOut/NotFound/SecurityGroupsMatch), schema validation (1 resource), sub-fields (networks/nodes)

## 14. Enhancement: Comprehensive md Examples

- [x] 14.1 Update `resource_tc_vdb_instance.md`: add second example (monthly subscription cluster) covering `security_group_ids`, `mode`, `goods_num`, `product_type`, `params`, `resource_tags`, `pay_period`, `auto_renew`

## 14a. Enhancement: ForceNew for VPC/Subnet + Immutable Fields Check in Update

- [x] 14a.1 Set `ForceNew: true` on `vpc_id` and `subnet_id` schema fields (Terraform handles recreation at plan time)
- [x] 14a.2 Add `immutableFields` array at the beginning of Update function, listing non-updatable input fields (excluding `vpc_id`/`subnet_id` which are handled by ForceNew): `pay_mode`, `instance_name`, `pay_period`, `auto_renew`, `params`, `resource_tags`, `instance_type`, `mode`, `goods_num`, `product_type`, `node_type`
- [x] 14a.3 Iterate `immutableFields` array, for each field call `d.HasChange(field)`, if changed return `fmt.Errorf("argument '%s' cannot be changed", field)`

## 15. Final Verification

- [x] 15.1 Run `go build .` to verify compilation after all enhancements
- [x] 15.2 Run `go test ./tencentcloud/services/vdb/ -v` to verify all 25 unit tests pass
- [x] 15.3 Run `make doc` to regenerate website documentation
