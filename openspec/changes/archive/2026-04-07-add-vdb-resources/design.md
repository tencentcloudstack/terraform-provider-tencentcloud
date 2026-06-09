## Context

TencentCloud Vector Database (VDB) is a managed vector database service. The SDK package `tencentcloud-sdk-go/tencentcloud/vdb/v20230616` is already vendored and exposes 14 个 API 接口. Currently, no Terraform resources exist for VDB.

Based on 5 种资源类型分类 (general / datasource / config / attachment / operation), we identified 1 resource to create.

## Goals / Non-Goals

**Goals:**
- Provide 1 Terraform resource:
  1. `tencentcloud_vdb_instance` (general) — 完整 CRUD 生命周期，security_group_ids 为 Required 参数，通过 ModifyDBInstanceSecurityGroups 实现安全组更新
- **All non-deprecated SDK parameters must be exposed** in schema (input or computed).
- **All async operations (CUD) must poll until status reaches expected value** before returning.
- Follow existing provider patterns for each resource type.
- Register the VDB client in `connectivity/client.go` and all resources in `provider.go`.
- Generate website documentation via `make doc`.

**Non-Goals:**
- Pricing APIs (`DescribePrice*`) are not in the 14 个接口列表, not exposed.
- `ModifyDBInstanceSecurityGroups` is used by `tencentcloud_vdb_instance` for inline security group updates (full-replace). It is NOT exposed as a separate resource.
- `AssociateSecurityGroups` / `DisassociateSecurityGroups` are not used — security groups are managed inline via `ModifyDBInstanceSecurityGroups`.
- `RecoverInstance` API is not exposed as a Terraform resource. Instance recovery is a manual operation best performed via console or CLI.
- `DescribeInstanceMaintenanceWindow` / `ModifyInstanceMaintenanceWindow` APIs are not exposed. Maintenance window management is not in scope for this change.
- Data sources are not in scope for this change.
- SDK 中标记为 `Deprecated` 的参数不暴露。

## Decisions

### Decision 1: Resource Classification by 5 Types

Based on the 5 种资源类型 framework:

| 资源 | 类型 | 接口要求 | 使用的云API |
|---|---|---|---|
| `tencentcloud_vdb_instance` | general | CRUD | C:`CreateInstance` R:`DescribeInstances`+`DescribeInstanceNodes`+`DescribeDBSecurityGroups` U:`ScaleUpInstance`+`ScaleOutInstance`+`ModifyDBInstanceSecurityGroups` D:`IsolateInstance`+`DestroyInstances` |

### Decision 2: Instance Read includes DescribeInstanceNodes

**Choice**: `DescribeInstanceNodes` is called as part of the Instance Read function to populate `nodes` computed attribute (Pod name and status list).

**Rationale**: Node info is a read-only attribute of the instance, naturally queried alongside instance details during Read.

### Decision 3: Instance Delete with force_delete Parameter

**Choice**: Instance schema includes a `force_delete` boolean parameter (default `false`):
- `force_delete = false`: Delete only calls `IsolateInstance`, moving the instance to recycle bin. User can later recover it.
- `force_delete = true`: Delete calls `IsolateInstance` first, then `DestroyInstances` to permanently destroy the instance.

**Rationale**: Provides flexibility for users — soft delete (recycle bin) by default for safety, hard delete on demand. This is a common pattern in TencentCloud resources (e.g., CDB MySQL).

### Decision 4: Config Resource Pattern for Maintenance Window

*REMOVED — `tencentcloud_vdb_instance_maintenance_window` is no longer part of this change. Maintenance window APIs (`DescribeInstanceMaintenanceWindow`, `ModifyInstanceMaintenanceWindow`) are not exposed.*

### Decision 5: Attachment Resource Pattern for Security Group

*Removed — `tencentcloud_vdb_security_group_attachment` is no longer part of this change. Security groups are now managed inline in `tencentcloud_vdb_instance` via `security_group_ids` (Required) with `ModifyDBInstanceSecurityGroups` API for updates.*

### Decision 6: Operation Resource Pattern for Recover

*Removed — `tencentcloud_vdb_instance_recover_operation` is no longer part of this change. `RecoverInstance` API is not exposed as a Terraform resource.*

### Decision 7: Service Client Registration

**Choice**: Register as `UseVdbV20230616Client()` in `connectivity/client.go` with 300s timeout profile.

### Decision 8: File Organization

```
tencentcloud/services/vdb/
├── service_tencentcloud_vdb.go
├── resource_tc_vdb_instance.go
├── resource_tc_vdb_instance.md
├── resource_tc_vdb_instance_test.go
└── resource_test.go
```

### Decision 9: Abstract Polling Helpers

**Choice**: Extract reusable polling helper methods into `service_tencentcloud_vdb.go`:

1. `WaitForInstanceStatus` — 基于状态轮询，用于 Create（等待 `online`）、Delete（等待 `isolated`）：
```go
func (me *VdbService) WaitForInstanceStatus(ctx context.Context, instanceId string, targetStatus string, timeout time.Duration) error
```
使用 `strings.EqualFold` 进行大小写无关比较，每次轮询打印实际状态值。

2. `WaitForInstanceScaleUp` — 基于值比较轮询，用于垂直扩容（CPU/Memory/DiskSize）：
```go
func (me *VdbService) WaitForInstanceScaleUp(ctx context.Context, instanceId string, targetCpu float64, targetMemory float64, targetDiskSize uint64, timeout time.Duration) error
```
轮询 `DescribeInstances` 直到实例的 CPU、Memory、Disk 值与目标值一致。

3. `WaitForInstanceScaleOut` — 基于值比较轮询，用于水平扩容（WorkerNodeNum）：
```go
func (me *VdbService) WaitForInstanceScaleOut(ctx context.Context, instanceId string, targetReplicaNum uint64, timeout time.Duration) error
```
轮询 `DescribeInstances` 直到实例的 ReplicaNum 与目标值一致。

4. `WaitForInstanceNotFound` — 基于不存在轮询，用于 `force_delete=true` 后确认实例完全销毁：
```go
func (me *VdbService) WaitForInstanceNotFound(ctx context.Context, instanceId string, timeout time.Duration) error
```
轮询 `DescribeVdbInstanceById` 直到返回 nil（实例不存在），确认 `DestroyInstances` 已彻底完成。

5. `WaitForSecurityGroupsMatch` — 基于值比较轮询，用于安全组更新后确认生效：
```go
func (me *VdbService) WaitForSecurityGroupsMatch(ctx context.Context, instanceId string, targetSgIds []string, timeout time.Duration) error
```
在 `ModifyDBInstanceSecurityGroups` 后轮询 `DescribeDBSecurityGroups` 直到返回的安全组列表与目标列表一致。

**Rationale**: Create/Delete 关注的是实例状态变化，使用状态轮询即可。`force_delete=true` 场景下，`DestroyInstances` 后实例可能短暂存在，需要通过 `WaitForInstanceNotFound` 确认彻底删除。而扩容操作中，实例状态可能在操作完成前就回到 `online`，但实际资源值尚未更新，因此必须通过值比较来确认扩容真正完成。安全组更新同理，`ModifyDBInstanceSecurityGroups` 返回后绑定可能尚未生效，通过 `WaitForSecurityGroupsMatch` 确认。

### Decision 10: Complete SDK Parameter Coverage

**Choice**: All non-deprecated parameters from the VDB SDK must be mapped in the resource schema.

For `tencentcloud_vdb_instance`:

**Create input** — add missing non-deprecated parameter:
- `params` (Optional, string) — instance extra parameters via JSON

**Read computed** — add missing non-deprecated InstanceInfo response fields:
- `region` (Computed, string)
- `zone` (Computed, string)
- `product` (Computed, string)
- `shard_num` (Computed, int)
- `api_version` (Computed, string)
- `extend` (Computed, string)
- `expired_at` (Computed, string)
- `is_no_expired` (Computed, bool)
- `wan_address` (Computed, string)
- `isolate_at` (Computed, string)
- `task_status` (Computed, int)

**Network sub-fields** — add missing:
- `networks[].preserve_duration` (Computed, int)
- `networks[].expire_time` (Computed, string)

**Deprecated SDK fields NOT to expose** (确认跳过):
`Project`, `NetworkType`, `TemplateId`, `Components`, `Zone` (Create input), `SlaveZones`, `IsNoExpired` (Create input), `EngineName` (Create input), `EngineVersion` (Create input), `Brief`, `Chief`, `DBA`, `HealthScore`, `Warning`.

**Hardcoded SDK fields NOT to expose** (内部固定值):
- `GoodsNum`: 固定为 `1`。Terraform Provider 每个 resource block 管理一个资源，购买数量由 Provider 内部控制，不暴露给用户。

Note: `NodeType` is deprecated in `CreateInstanceRequestParams` but still present in `InstanceInfo` response. Keep it as ForceNew input + Computed for backward compatibility.

### Decision 11: Serial Execution of Scale Operations in Update

**Choice**: When a Terraform plan requires both vertical scaling (CPU/memory/disk) and horizontal scaling (worker_node_num) simultaneously, they MUST execute serially:
1. First execute `ScaleUpInstance` (vertical)，然后调用 `WaitForInstanceScaleUp` 轮询直到 CPU/Memory/DiskSize 值与目标一致
2. Then execute `ScaleOutInstance` (horizontal)，然后调用 `WaitForInstanceScaleOut` 轮询直到 ReplicaNum 与目标一致

**Rationale**: VDB instances cannot process concurrent scale operations. The API will reject a second scale request while the first is in progress. Serial execution with value-based polling between operations ensures reliability — 基于值比较而非状态比较可以确保扩容操作真正完成。

### Decision 12: VdbClientInterface for Mock-Based Unit Testing

**Choice**: Define `VdbClientInterface` in `service_tencentcloud_vdb.go`，抽象 VDB SDK client 的所有方法为接口。`VdbService` 通过 `getVdbClient()` 方法获取 client：优先使用注入的 mock（`vdbClient` 字段），否则回退到真实的 SDK client（`client.UseVdbV20230616Client()`）。

```go
type VdbClientInterface interface {
    DescribeInstancesWithContext(ctx, request) (response, error)
    DescribeInstanceNodesWithContext(ctx, request) (response, error)
    DescribeDBSecurityGroupsWithContext(ctx, request) (response, error)
    CreateInstanceWithContext(ctx, request) (response, error)
    ScaleUpInstanceWithContext(ctx, request) (response, error)
    ScaleOutInstanceWithContext(ctx, request) (response, error)
    IsolateInstanceWithContext(ctx, request) (response, error)
    DestroyInstancesWithContext(ctx, request) (response, error)
    ModifyDBInstanceSecurityGroupsWithContext(ctx, request) (response, error)
}
```

测试文件使用 `testify/mock` 创建 `MockVdbClient`，通过 `newTestService(mockClient)` 注入 mock，不依赖 Terraform 验收测试框架，不调用真实云 API。

**Rationale**: 此项目所有现有测试均为验收测试（调用真实 API），缺少单元测试。通过接口抽象 + mock 注入，可以在不改变生产代码调用链的前提下，对 service 层和 resource schema 进行完整的单元测试。接口兼容真实 SDK client（`*vdb.Client` 已实现所有方法），对生产代码零侵入。

### Decision 13: Comprehensive md Examples Covering All Parameters

**Choice**: 每个 resource 的 `.md` 文件应提供足够的 HCL 示例，覆盖所有参数的用法：

- `resource_tc_vdb_instance.md`：提供两个示例 —— 按量付费单机版（最小参数集）和包年包月集群版（全量参数，包含 `security_group_ids`, `mode`, `product_type`, `params`, `resource_tags`）。

**Rationale**: 用户通过 `.md` 中的 HCL 示例了解参数用法，示例不充分会导致用户不知道某些参数的存在或用法。

### Decision 14: Security Groups as Required Inline Attribute with ModifyDBInstanceSecurityGroups

**Choice**: `security_group_ids` is a Required attribute (not ForceNew) on `tencentcloud_vdb_instance`. Security group updates are handled inline in the Update function:
1. When `security_group_ids` changes, call `ModifyDBInstanceSecurityGroups` API with the full target list (full-replace semantics).
2. After the API call, poll via `WaitForSecurityGroupsMatch(ctx, instanceId, targetSgIds, timeout)` which calls `DescribeDBSecurityGroups` until the returned list matches the target list.

**Rationale**: The `ModifyDBInstanceSecurityGroups` API performs a full-replace of all security groups on an instance, which maps naturally to a list attribute on the instance resource. An attachment pattern would require `AssociateSecurityGroups`/`DisassociateSecurityGroups` for individual bindings, adding complexity without benefit. Making `security_group_ids` Required (instead of Optional) ensures every VDB instance has explicit security group configuration, which is a security best practice.

### Decision 15: ForceNew for VPC/Subnet + Immutable Fields Check in Update

**Choice**: `vpc_id` 和 `subnet_id` 设置 `ForceNew: true`，由 Terraform 在 plan 阶段自动处理重建。其余不可修改的字段通过 Update 方法开头的 `immutableFields` 数组检查：

```go
// vpc_id, subnet_id 已设置 ForceNew: true，Terraform plan 阶段即提示需重建
// goods_num 已从 schema 移除，内部硬编码为 1
immutableFields := []string{
    "pay_mode", "instance_name",
    "pay_period", "auto_renew", "params", "resource_tags",
    "instance_type", "mode", "product_type", "node_type",
}
for _, field := range immutableFields {
    if d.HasChange(field) {
        return fmt.Errorf("argument `%s` cannot be changed", field)
    }
}
```

可修改的字段仅有：`cpu`、`memory`、`disk_size`、`worker_node_num`、`security_group_ids`。

**Rationale**:
- `vpc_id`、`subnet_id` 是网络基础设施字段，API Read 会返回，设置 ForceNew 让 Terraform plan 阶段即提示需要重建，用户体验最好。
- 其余不可变字段（如 `pay_mode`、`instance_name` 等）API 也会返回但无修改接口，通过 `immutableFields` 数组在 Update 入口拦截，提供清晰的错误信息。
- 两层防护互补：ForceNew 在 plan 阶段拦截，immutableFields 在 apply 阶段兜底。

## Risks / Trade-offs

- **[Risk] Async instance creation**: VDB instance creation is asynchronous, status transitions: `creating` → `online`. → **Mitigation**: Use `WaitForInstanceStatus` helper to poll until `online`.
- **[Risk] force_delete=false leaves resources in recycle bin**: Users may forget isolated instances incur costs. → **Mitigation**: Document clearly in `.md` that `force_delete=false` only isolates, not destroys.
- **[Risk] force_delete=true destroy may not be immediate**: `DestroyInstances` is asynchronous, instance may still exist briefly after API call returns. → **Mitigation**: After `DestroyInstances`, call `WaitForInstanceNotFound` to poll until instance is completely gone (`DescribeVdbInstanceById` returns nil).
- **[Risk] Scale operations may also be async**: → **Mitigation**: ScaleUp 后调用 `WaitForInstanceScaleUp` 基于值比较轮询直到 CPU/Memory/DiskSize 一致；ScaleOut 后调用 `WaitForInstanceScaleOut` 基于值比较轮询直到 ReplicaNum 一致。Serial execution prevents conflicts.
- **[Risk] DescribeInstanceNodes may return empty during scaling**: → **Mitigation**: Treat empty node list as valid state during Read.
- **[Risk] Concurrent scale operations rejected by API**: → **Mitigation**: Decision 11 enforces serial execution with polling between operations.

## Appendix: VDB Instance Status Values

Based on actual API testing, the VDB instance status values are:

| 状态值 | 含义 |
|---|---|
| `creating` | 创建中 |
| `online` | 正常运行 |
| `isolated` | 已隔离（回收站） |

**重要**：状态值为小写英文，轮询时使用 `strings.EqualFold` 进行大小写无关比较以确保兼容性。
