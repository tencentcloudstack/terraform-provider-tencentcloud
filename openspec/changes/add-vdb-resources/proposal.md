## Why

TencentCloud Vector Database (VDB) is a cloud product that currently lacks Terraform resource support. Users need to manage VDB instances via Terraform for infrastructure-as-code workflows. The VDB SDK (`tencentcloud-sdk-go/tencentcloud/vdb/v20230616`) is already vendored and provides 14 个云 API 接口。

## What Changes

- Add new service directory `tencentcloud/services/vdb/` with service layer, resources, tests, and documentation.

### Resource 1: `tencentcloud_vdb_instance` — 类型: general（通用资源）

管理 VDB 实例的完整生命周期。

| CRUD | 云API | 说明 |
|---|---|---|
| **C** | `CreateInstance` | 创建实例（security_group_ids 为 Required 参数） |
| **R** | `DescribeInstances` + `DescribeInstanceNodes` + `DescribeDBSecurityGroups` | 查询实例详情 + 查询实例 Pod 节点列表 + 查询安全组绑定 |
| **U** | `ScaleUpInstance` + `ScaleOutInstance` + `ModifyDBInstanceSecurityGroups` | 垂直扩容（CPU/内存/磁盘）+ 水平扩容（节点数）+ 安全组全量替换 |
| **D** | `IsolateInstance` + `DestroyInstances`（条件调用） | `force_delete=false` 时仅隔离到回收站；`force_delete=true` 时隔离后再销毁，并轮询 `WaitForInstanceNotFound` 直到实例完全消失 |

### 其他变更
- Register VDB client in `tencentcloud/connectivity/client.go`.
- Register all VDB resources in `tencentcloud/provider.go` and update `tencentcloud/provider.md`.
- Generate website documentation via `make doc`.

**不生成 resource 的 API 及原因：**
- `DescribePriceCreateInstance` / `DescribePriceRenewInstance` / `DescribePriceResizeInstance`：SDK 中存在但不在本次需求的 14 个接口列表中，且仅查询价格，无 CUD。
- `AssociateSecurityGroups` / `DisassociateSecurityGroups`：安全组绑定改为在 `tencentcloud_vdb_instance` 中通过 `ModifyDBInstanceSecurityGroups` 全量替换实现，不再使用单条绑定/解绑接口。
- `RecoverInstance`：不接入。恢复隔离实例的操作不适合作为 Terraform 资源暴露，用户可通过控制台或 CLI 手动恢复。
- `DescribeInstanceMaintenanceWindow` / `ModifyInstanceMaintenanceWindow`：不接入。维护时间窗功能不在本次需求范围内。

**所有 14 个接口分配如下（9 个使用，5 个不使用）：**

| 云API | 归属资源 | 用于 |
|---|---|---|
| `CreateInstance` | vdb_instance | C |
| `DescribeInstances` | vdb_instance | R |
| `DescribeInstanceNodes` | vdb_instance | R |
| `DescribeDBSecurityGroups` | vdb_instance | R（读取安全组绑定） |
| `ScaleUpInstance` | vdb_instance | U |
| `ScaleOutInstance` | vdb_instance | U |
| `ModifyDBInstanceSecurityGroups` | vdb_instance | U（安全组全量替换） |
| `IsolateInstance` | vdb_instance | D |
| `DestroyInstances` | vdb_instance | D (force_delete=true) |
| `DescribeInstanceMaintenanceWindow` | 不接入 | — |
| `ModifyInstanceMaintenanceWindow` | 不接入 | — |
| `AssociateSecurityGroups` | 不使用 | — |
| `DisassociateSecurityGroups` | 不使用 | — |
| `RecoverInstance` | 不接入 | — |

## Capabilities

### New Capabilities
- `vdb-instance`: General resource for VDB instance lifecycle. APIs: CreateInstance, DescribeInstances, DescribeInstanceNodes, DescribeDBSecurityGroups, ScaleUpInstance, ScaleOutInstance, ModifyDBInstanceSecurityGroups, IsolateInstance, DestroyInstances. Security groups managed inline via `security_group_ids` (Required) with full-replace updates.

### Modified Capabilities
<!-- No existing capabilities are being modified -->

## VDB Instance Status Values (实测)

VDB 云 API 返回的实例状态值为小写英文，与部分腾讯云产品（如 CDB 使用 `Running`）不同：

| 状态值 | 含义 | 出现场景 |
|---|---|---|
| `creating` | 创建中 | CreateInstance 后 |
| `online` | 正常运行 | 创建完成、扩容完成 |
| `isolated` | 已隔离 | IsolateInstance 后 |

代码中所有轮询使用 `strings.EqualFold` 进行大小写无关比较，并在每次轮询打印实际状态值用于排查。

## Design Principles

### 1. 状态轮询与值比较轮询
- Create/Delete 等操作关注实例状态变化，使用 `WaitForInstanceStatus(ctx, instanceId, targetStatus, timeout)` 轮询直到状态匹配。
- Delete with `force_delete=true` 在 `DestroyInstances` 后使用 `WaitForInstanceNotFound(ctx, instanceId, timeout)` 轮询直到实例完全消失（`DescribeVdbInstanceById` 返回 nil）。
- ScaleUp 操作使用 `WaitForInstanceScaleUp(ctx, instanceId, targetCpu, targetMemory, targetDiskSize, timeout)` 基于值比较轮询，直到实例的 CPU/Memory/DiskSize 与目标值一致。
- ScaleOut 操作使用 `WaitForInstanceScaleOut(ctx, instanceId, targetReplicaNum, timeout)` 基于值比较轮询，直到实例的 ReplicaNum 与目标值一致。
- 安全组更新使用 `WaitForSecurityGroupsMatch(ctx, instanceId, targetSgIds, timeout)` 基于值比较轮询，在 `ModifyDBInstanceSecurityGroups` 后轮询 `DescribeDBSecurityGroups` 直到返回的安全组列表与目标列表一致。
- 扩容使用值比较而非状态比较的原因：实例状态可能在操作完成前就回到 `online`，但实际资源值尚未更新。

### 2. 完整 SDK 参数覆盖
所有非 Deprecated 的 SDK 参数必须在 schema 中暴露（作为 input 或 computed）。经对比，`tencentcloud_vdb_instance` 资源缺少以下参数需要补齐：
- **Create 输入**: `params` (实例额外参数，JSON 格式)
- **Read 计算属性**: `region`, `zone`, `product`, `shard_num`, `api_version`, `extend`, `expired_at`, `is_no_expired`, `wan_address`, `isolate_at`, `task_status`
- **Networks 子字段**: `preserve_duration`, `expire_time`

其他 resource 已完整覆盖 SDK 参数。

### 3. 水平扩容与垂直扩容串行执行
当 Terraform plan 同时涉及垂直扩容（CPU/内存/磁盘）和水平扩容（节点数）时，必须严格串行执行：先 ScaleUp → `WaitForInstanceScaleUp` 等待值一致 → 再 ScaleOut → `WaitForInstanceScaleOut` 等待值一致。VDB API 不支持并发扩容操作。

### 4. VdbClientInterface 接口抽象与 Mock 单元测试
通过定义 `VdbClientInterface` 接口抽象 VDB SDK client 的所有方法，`VdbService` 支持注入 mock client。测试使用 `testify/mock` 框架创建 `MockVdbClient`，不依赖 Terraform 验收测试框架，不调用真实云 API。共 25 个单元测试覆盖 service 层、轮询逻辑和 schema 验证。

### 5. 完整的 md 示例覆盖全部参数
每个资源的 `.md` 文件提供充分的 HCL 示例，确保用户可以了解全部参数的用法。参数较多的资源提供多个示例（最小参数集 + 全量参数），参数较少的资源提供单个完整示例。

### 6. ForceNew + 不可变字段检查
`vpc_id`、`subnet_id` 设置 `ForceNew: true`，由 Terraform plan 阶段自动处理重建。其余不可变字段（`pay_mode`、`instance_name`、`pay_period`、`auto_renew`、`params`、`resource_tags`、`instance_type`、`mode`、`goods_num`、`product_type`、`node_type`）通过 Update 方法入口的 `immutableFields` 数组检查，变更时直接报错。

## Impact

- **New files**: `tencentcloud/services/vdb/` directory with ~5 files (service, 1 resource × {.go, .md, _test.go}, resource_test.go).
- **Modified files**: `tencentcloud/connectivity/client.go`, `tencentcloud/provider.go`, `tencentcloud/provider.md`.
- **Dependencies**: Uses already-vendored `tencentcloud-sdk-go/tencentcloud/vdb/v20230616`.
- **Generated files**: `website/docs/r/vdb_*.html.markdown` via `make doc`.
