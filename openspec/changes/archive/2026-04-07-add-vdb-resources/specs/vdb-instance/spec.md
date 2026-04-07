## ADDED Requirements

### Requirement: Reusable polling helpers
The service layer SHALL provide the following polling helper methods:
1. `WaitForInstanceStatus(ctx, instanceId, targetStatus, timeout)` — 基于状态轮询，用于 Create（等待 `online`）和 Delete（等待 `isolated`）。比较使用 `strings.EqualFold`，每次轮询打印实际状态。
2. `WaitForInstanceScaleUp(ctx, instanceId, targetCpu, targetMemory, targetDiskSize, timeout)` — 基于值比较轮询，用于垂直扩容。轮询 `DescribeInstances` 直到实例的 CPU/Memory/Disk 与目标值一致。
3. `WaitForInstanceScaleOut(ctx, instanceId, targetReplicaNum, timeout)` — 基于值比较轮询，用于水平扩容。轮询 `DescribeInstances` 直到实例的 ReplicaNum 与目标值一致。
4. `WaitForInstanceNotFound(ctx, instanceId, timeout)` — 基于不存在轮询，用于 `force_delete=true` 后确认实例完全销毁。轮询 `DescribeVdbInstanceById` 直到返回 nil。
5. `WaitForSecurityGroupsMatch(ctx, instanceId, targetSgIds, timeout)` — 基于值比较轮询，用于安全组更新后确认生效。轮询 `DescribeDBSecurityGroups` 直到返回的安全组列表与目标列表一致。

#### Scenario: Poll until target status (Create/Delete)
- **WHEN** Create or Delete calls `WaitForInstanceStatus` with a target status (e.g., `"online"`, `"isolated"`)
- **THEN** the helper polls `DescribeInstances` repeatedly until the instance status matches (case-insensitive), returning nil on match or error on timeout

#### Scenario: Poll until scale up values match
- **WHEN** ScaleUp calls `WaitForInstanceScaleUp` with target CPU, Memory, DiskSize values
- **THEN** the helper polls `DescribeInstances` repeatedly until the instance's CPU, Memory, Disk all equal the target values

#### Scenario: Poll until scale out value matches
- **WHEN** ScaleOut calls `WaitForInstanceScaleOut` with target ReplicaNum value
- **THEN** the helper polls `DescribeInstances` repeatedly until the instance's ReplicaNum equals the target value

#### Scenario: Poll until instance not found (delete confirmation)
- **WHEN** Delete with `force_delete=true` calls `WaitForInstanceNotFound` after `DestroyInstances`
- **THEN** the helper polls `DescribeVdbInstanceById` repeatedly until it returns nil (instance no longer exists), returning nil on not-found or error on timeout

### Requirement: Create VDB instance with all non-deprecated SDK parameters
The system SHALL allow users to create a VDB instance exposing ALL non-deprecated `CreateInstanceRequest` parameters: vpc_id, subnet_id, pay_mode, instance_name, security_group_ids (Required), pay_period, auto_renew, params (JSON string), resource_tags, instance_type, mode, product_type, node_type, cpu, memory, disk_size, worker_node_num. The `goods_num` parameter SHALL NOT be exposed to users — it SHALL be hardcoded to `1` in the Create function, since the Provider manages one resource per block. After calling `CreateInstance` API, the system SHALL call `WaitForInstanceStatus` to wait for status `online`.

#### Scenario: Successful instance creation
- **WHEN** user applies a `tencentcloud_vdb_instance` resource with valid configuration
- **THEN** the system creates the instance via `CreateInstance` API, calls `WaitForInstanceStatus(ctx, id, "online", createTimeout)`, and stores the instance ID in state

#### Scenario: Instance creation with params
- **WHEN** user specifies `params` (JSON string) in the resource configuration
- **THEN** the system passes the params to the `CreateInstance` API request

### Requirement: Read VDB instance with all non-deprecated InstanceInfo fields
The system SHALL read ALL non-deprecated fields from `InstanceInfo` response and `DescribeInstanceNodes`, populating into state:
- Input fields reflected: instance_name, pay_mode, cpu, memory, disk_size (from Disk), instance_type, node_type, product_type, auto_renew, vpc_id/subnet_id (from Networks[0]), security_group_ids, worker_node_num (from ReplicaNum)
- Computed fields: status, region, zone, product, shard_num, api_version, extend, expired_at, is_no_expired, wan_address, isolate_at, task_status, created_at, engine_name, engine_version
- Networks list: vpc_id, subnet_id, vip, port, preserve_duration, expire_time
- Nodes list (from DescribeInstanceNodes): name, status

#### Scenario: Instance exists
- **WHEN** Terraform refreshes state for an existing `tencentcloud_vdb_instance`
- **THEN** the system calls `DescribeInstances` + `DescribeInstanceNodes` and sets ALL non-deprecated attributes

#### Scenario: Instance not found
- **WHEN** Terraform refreshes state but the instance no longer exists
- **THEN** the system removes the resource from state by calling `d.SetId("")`

### Requirement: ForceNew for VPC/Subnet and immutable fields check in Update
`vpc_id` and `subnet_id` SHALL be marked `ForceNew: true` in the schema, so Terraform triggers resource recreation at plan time when these fields change. The Update function SHALL additionally define an `immutableFields` array listing fields that cannot be modified in-place (`pay_mode`, `instance_name`, `pay_period`, `auto_renew`, `params`, `resource_tags`, `instance_type`, `mode`, `product_type`, `node_type`). Before processing any update logic, the function SHALL iterate the array and reject changes to any listed field with an error message `argument 'X' cannot be changed`.

#### Scenario: Change vpc_id or subnet_id
- **WHEN** user modifies `vpc_id` or `subnet_id`
- **THEN** Terraform plan shows the resource will be destroyed and recreated (ForceNew behavior)

#### Scenario: Attempt to change other immutable field
- **WHEN** user modifies a field in the immutableFields array (e.g., `instance_name`, `pay_mode`)
- **THEN** the system returns an error: `argument 'instance_name' cannot be changed`

#### Scenario: Change mutable field
- **WHEN** user modifies a mutable field (e.g., `cpu`, `security_group_ids`)
- **THEN** the system proceeds with the update logic normally

### Requirement: Update VDB instance via serial scaling operations with value-based polling
The system SHALL support updating a VDB instance via `ScaleUpInstance` (vertical: cpu/memory/disk_size) and `ScaleOutInstance` (horizontal: worker_node_num). When both types of changes occur simultaneously, they MUST execute serially. ScaleUp SHALL use `WaitForInstanceScaleUp` to poll until CPU/Memory/DiskSize match target values. ScaleOut SHALL use `WaitForInstanceScaleOut` to poll until ReplicaNum matches target value.

#### Scenario: Scale up only
- **WHEN** user modifies cpu, memory, or disk_size attributes (but NOT worker_node_num)
- **THEN** the system calls `ScaleUpInstance` API, then `WaitForInstanceScaleUp(targetCpu, targetMemory, targetDiskSize)` until values match

#### Scenario: Scale out only
- **WHEN** user increases worker_node_num (but NOT cpu/memory/disk_size)
- **THEN** the system calls `ScaleOutInstance` API, then `WaitForInstanceScaleOut(targetReplicaNum)` until value matches

#### Scenario: Concurrent scale up and scale out
- **WHEN** user modifies both cpu/memory/disk_size AND worker_node_num simultaneously
- **THEN** the system first calls `ScaleUpInstance` + `WaitForInstanceScaleUp(targetCpu, targetMemory, targetDiskSize)`, then calls `ScaleOutInstance` + `WaitForInstanceScaleOut(targetReplicaNum)` — strict serial order

### Requirement: Update security groups via ModifyDBInstanceSecurityGroups with polling
The system SHALL support updating `security_group_ids` (Required, not ForceNew) via the `ModifyDBInstanceSecurityGroups` API, which performs a full-replace of all security groups on the instance. After calling the API, the system SHALL call `WaitForSecurityGroupsMatch(ctx, instanceId, targetSgIds, timeout)` to poll `DescribeDBSecurityGroups` until the returned security group list matches the target list.

#### Scenario: Update security groups
- **WHEN** user modifies `security_group_ids` in the resource configuration
- **THEN** the system calls `ModifyDBInstanceSecurityGroups` with the new full list, then `WaitForSecurityGroupsMatch` to poll until the actual security groups match the target list

#### Scenario: Update security groups alongside scaling
- **WHEN** user modifies both `security_group_ids` and scaling parameters (cpu/memory/disk_size/worker_node_num)
- **THEN** the system handles security group update and scaling operations independently in the Update function (security group update does not depend on scaling completion or vice versa)

### Requirement: Delete VDB instance with force_delete control and status polling
The system SHALL support a `force_delete` boolean parameter (default `false`). When `force_delete=false`, delete calls `IsolateInstance` and waits for `isolated` status. When `force_delete=true`, delete calls `IsolateInstance`, waits for `isolated`, then calls `DestroyInstances`, then calls `WaitForInstanceNotFound` to poll until the instance is completely gone.

#### Scenario: Soft delete (force_delete=false, default)
- **WHEN** user destroys with `force_delete` unset or `false`
- **THEN** the system calls `IsolateInstance`, then `WaitForInstanceStatus("isolated")`

#### Scenario: Hard delete (force_delete=true)
- **WHEN** user destroys with `force_delete=true`
- **THEN** the system calls `IsolateInstance`, `WaitForInstanceStatus("isolated")`, then `DestroyInstances`, then `WaitForInstanceNotFound(ctx, instanceId, timeout)` to poll until the instance is completely gone

### Requirement: Import VDB instance
The system SHALL support importing existing VDB instances by instance ID.

#### Scenario: Import by instance ID
- **WHEN** user runs `terraform import tencentcloud_vdb_instance.example vdb-xxx`
- **THEN** the system reads the instance details and populates state

### Requirement: VdbClientInterface for mock-based unit testing
The service layer SHALL define a `VdbClientInterface` interface abstracting all 9 VDB SDK client methods used by service and resource functions (including `ModifyDBInstanceSecurityGroupsWithContext`). `VdbService` SHALL support injecting a mock client via `vdbClient` field, with `getVdbClient()` returning the injected mock or falling back to the real SDK client. This enables unit testing without calling real cloud APIs.

#### Scenario: Service uses injected mock
- **WHEN** `VdbService.vdbClient` is set to a mock implementation
- **THEN** all service methods (Describe*, WaitFor*) use the mock instead of the real SDK client

#### Scenario: Service falls back to real client
- **WHEN** `VdbService.vdbClient` is nil (production code path)
- **THEN** `getVdbClient()` returns `me.client.UseVdbV20230616Client()` — the real SDK client

### Requirement: Comprehensive md examples covering all parameters
The `.md` documentation for `tencentcloud_vdb_instance` SHALL provide at least two HCL examples: a minimal pay-as-you-go example and a full monthly subscription example covering all parameters (`security_group_ids`, `mode`, `goods_num`, `product_type`, `params`, `resource_tags`, `pay_period`, `auto_renew`).

#### Scenario: User reads instance md
- **WHEN** user reads the `resource_tc_vdb_instance.md` documentation
- **THEN** they can find HCL examples demonstrating every non-deprecated parameter
