## Why

CBS（云硬盘）目前不支持通过 Terraform 进行快照跨地域复制操作。用户需要手动在多个地域之间复制快照以实现数据容灾和迁移，无法通过基础设施即代码的方式统一管理跨地域快照的生命周期。新增 `tencentcloud_cbs_copy_snapshot_cross_region` 资源可以让用户在 Terraform 中声明式地管理快照跨地域复制，实现多地域数据同步的自动化运维。

## What Changes

- 新增 Terraform 资源 `tencentcloud_cbs_copy_snapshot_cross_region`（RESOURCE_KIND_ATTACHMENT 类型），支持快照跨地域复制的创建、查询和删除
- Create：调用 `CopySnapshotCrossRegions` 接口发起跨地域复制（异步接口），完成后轮询 `DescribeSnapshots` 直到快照状态变为 NORMAL
- Read：调用 `DescribeSnapshots` 接口查询复制的快照状态和详情
- Delete：调用 `DeleteSnapshots` 接口删除目标地域的复制快照
- 在 `provider.go` 和 `provider.md` 中注册新资源
- 生成对应的资源文档 `.md` 文件

## Capabilities

### New Capabilities
- `cbs-copy-snapshot-cross-region`: 快照跨地域复制资源，管理 CBS 快照从一个地域复制到一个或多个目标地域的绑定关系，包括创建（绑定）、读取、删除（解绑）操作

### Modified Capabilities

## Impact

- 新增文件：`tencentcloud/services/cbs/resource_tc_cbs_copy_snapshot_cross_region_attachment.go`
- 新增测试文件：`tencentcloud/services/cbs/resource_tc_cbs_copy_snapshot_cross_region_attachment_test.go`
- 新增文档：`tencentcloud/services/cbs/resource_tc_cbs_copy_snapshot_cross_region_attachment.md`
- 修改文件：`tencentcloud/provider.go`（注册新资源）、`tencentcloud/provider.md`（添加资源文档条目）
- 依赖的云 API 接口：`CopySnapshotCrossRegions`、`DescribeSnapshots`、`DeleteSnapshots`（均来自 cbs v20170312 SDK）
- CopySnapshotCrossRegions 为异步接口，需在 Create 中实现轮询逻辑
