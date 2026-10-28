## Context

CBS（云硬盘）产品提供了快照跨地域复制功能，允许用户将某个地域的快照复制到其他地域以实现数据容灾和迁移。当前 Terraform Provider 中已有 CBS 相关的 storage、snapshot、snapshot_policy 等资源，但缺少快照跨地域复制的资源支持。

CBS SDK 提供了以下三个相关接口：
- **CopySnapshotCrossRegions**（异步）：发起跨地域复制，返回每个目标地域的新快照 ID
- **DescribeSnapshots**：查询快照列表，可通过 SnapshotIds 过滤查询
- **DeleteSnapshots**：删除快照

## Goals / Non-Goals

**Goals:**
- 新增 `tencentcloud_cbs_copy_snapshot_cross_region` 资源（RESOURCE_KIND_ATTACHMENT），支持快照跨地域复制的创建、读取和删除
- 支持异步操作：Create 调用 CopySnapshotCrossRegions 后轮询 DescribeSnapshots 直到快照状态变为 NORMAL
- 支持 Import 操作，使用复合 ID（snapshot_id + "#" + copied_snapshot_id）
- 生成完整的单元测试和文档

**Non-Goals:**
- 不支持 Update 操作（ATTACHMENT 类型资源只有 CRD）
- 不支持跨地域复制的批量管理数据源（仅本次新增 ATTACHMENT 资源）
- 不修改现有 CBS 资源的 schema 或行为

## Decisions

### 1. 资源 ID 构造方式

**决策**：使用 `snapshot_id + FILED_SP + copied_snapshot_id` 作为复合 ID。

**理由**：
- `snapshot_id` 是源快照 ID，在 Create 请求中提供
- `copied_snapshot_id` 是复制后在目标地域生成的新快照 ID，由 CopySnapshotCrossRegions 接口返回的 SnapshotCopyResult 中获取
- 由于 destination_regions 是一个列表，一次创建可能生成多个目标快照。每个 terraform 资源实例代表一个源快照向一个目标地域的复制操作。当 destination_regions 包含多个地域时，使用第一个 copied_snapshot_id 作为 ID 的组成部分，其余 copied_snapshot_id 通过 `snapshot_copy_result_set` 存储在 state 中
- 采用 `snapshot_id#copied_snapshot_id` 格式，与项目中其他 ATTACHMENT 资源（如 cbs_snapshot_policy_attachment）的复合 ID 模式一致

**备选方案**：使用 `snapshot_id#destination_region#copied_snapshot_id` 三段 ID——但 destination_region 已经是 terraform 参数中明确指定的，不需要重复存储在 ID 中。

### 2. 异步轮询策略

**决策**：Create 调用 CopySnapshotCrossRegions 成功后，轮询 DescribeSnapshots 接口检查每个目标地域新快照的状态，直到 SnapshotState 为 NORMAL。

**理由**：
- CBS SDK 文档明确说明 CopySnapshotCrossRegions 是异步接口
- 复制过程中快照状态为 COPYING_FROM_REMOTE 或 CHECKING_COPIED，完成后变为 NORMAL
- 需要为每个 destination_region 中的 copied_snapshot_id 分别轮询状态
- 使用 `resource.Retry` + `tccommon.ReadRetryTimeout` 进行轮询

### 3. Delete 策略

**决策**：Delete 时调用 DeleteSnapshots 接口删除所有目标地域中由此次复制创建的快照，使用 snapshot_copy_result_set 中保存的 copied_snapshot_id 列表。

**理由**：
- DeleteSnapshots 接口接受 SnapshotIds 列表参数
- 从 state 中读取所有 copied_snapshot_id，构建删除请求
- delete_bind_images 参数设为可选，默认 false

### 4. Read 策略

**决策**：Read 时调用 DescribeSnapshots 接口，使用 copied_snapshot_id 查询快照状态和详情。

**理由**：
- DescribeSnapshots 支持 SnapshotIds 过滤查询
- 从复合 ID 中解析 copied_snapshot_id，构建查询请求
- 查询结果中可以获取快照的状态、名称、大小等信息

### 5. Schema 设计

**决策**：schema 中 Required 字段为 `snapshot_id` 和 `destination_regions`，Optional 字段为 `snapshot_name` 和 `delete_bind_images`，Computed 字段为 `snapshot_copy_result_set`。所有 Required/Optional 字段设置 ForceNew: true。

**理由**：
- ATTACHMENT 资源无 Update 操作，所有可写字段必须 ForceNew
- `snapshot_id` 和 `destination_regions` 是创建复制的必要参数
- `snapshot_name` 是可选的复制快照名称
- `delete_bind_images` 是 Delete 操作的可选参数
- `snapshot_copy_result_set` 是 Create 操作的返回结果，为 Computed 只读属性

## Risks / Trade-offs

- **[多地域复制状态轮询复杂度]** → 一次 Create 可能涉及多个目标地域，需要逐一轮询每个地域的新快照状态。缓解：在 Create 函数中对每个 SnapshotCopyResult 逐一检查状态，直到全部 NORMAL
- **[跨地域 API 调用限制]** → DescribeSnapshots 需要在目标地域查询，但当前 Terraform provider 的 client 连接是基于单一地域的。缓解：CopySnapshotCrossRegions 返回的 SnapshotCopyResult 中包含了目标地域的信息，但我们需要在源地域用 DescribeSnapshots 查询——实际上 CBS 的 DescribeSnapshots 可以在源地域查询所有快照（包括复制状态的），或者我们可以通过 snapshot_id 直接查询
- **[复制快照的删除风险]** → 删除操作会实际删除目标地域的快照数据。缓解：提供 delete_bind_images 可选参数让用户控制是否同时删除关联镜像
