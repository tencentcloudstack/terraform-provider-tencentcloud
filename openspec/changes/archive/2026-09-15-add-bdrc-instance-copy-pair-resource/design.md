## Context

腾讯云备份与容灾（BDRC）产品提供 CVM 实例级别的容灾复制对（Instance Copy Pair）管理能力。当前本 Provider 的 `tencentcloud/services/` 下尚不存在 `bdrc` 服务目录，`tencentcloud/connectivity/client.go` 中也没有 bdrc 客户端绑定。

vendored SDK（`vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/bdrc/v20260330/`）已提供以下接口：

- `CreateInstanceCopyPairWithContext`：创建 CVM 复制对。入参含 `ProtectGroupId`、`CreateTargetInstanceParameters`（`[]*CreateInstanceModel`，1~10 个）、`InstanceCopyPairName`、`ClientToken`、`RecoveryPointObjective`。出参返回 `CopyPairIds []*string`（创建的复制对 ID 列表）。该接口为异步接口。
- `DescribeCopyPairsWithContext`：查询容灾复制对。入参含 `CopyPairType`（必填，CVM 场景传 `INSTANCE`）、`CopyPairIds`、`Filters`、分页 `Offset`/`Limit`、`Order`、`OrderField`、`QueryProtectionTime`、`GetAllCopyPair`、`QueryCvmCreateParams`、`CreateFrom`。出参返回 `CopyPairSet []*CopyPair`，每个 `CopyPair` 含 `CopyPairId`、`CopyPairName`、`CopyPairState`、`SourceRegion`、`TargetRegion`、`SourceResourceId`、`TargetResourceId`、`Percent`、`RecoveryPointObjective`、`DataDirection`、`CreateTime`、`DiskCopyPairSet`、`CvmCreateParams` 等大量字段。
- `ModifyCopyPairAttributeWithContext`：修改复制对属性。入参仅 `CopyPairId`（必填）、`CopyPairType`（可选，默认 `INSTANCE`）、`CopyPairName`（可选）。
- `DeleteCopyPairsWithContext`：删除复制对。入参 `CopyPairIds`（必填）、`CopyPairType`（必填）、`DeleteTargetResource`（可选）。

`CreateInstanceModel` 结构嵌套了 `Placement`、`SystemDisk`（`DiskModel`）、`InstanceChargePrepaid`、`DataDisks`（`[]*DiskModel`）、`VirtualPrivateCloud`、`InternetAccessible`、`LoginSettings`、`EnhancedService`（含 `SecurityService`/`MonitorService`/`AutomationService`/`BasicService`，各含 `Enabled *bool`）等子结构。

参考资源为 `tencentcloud_igtm_strategy`，其单文件布局（schema + CRUD + 辅助函数）、retry 拓扑、nil 防御、日志惯例均符合本 Provider 规范。

## Goals / Non-Goals

**Goals:**
- 提供 `tencentcloud_bdrc_instance_copy_pair` 资源的完整 CRUD 生命周期管理（创建/读取/更新/删除/导入）。
- Schema 字段完整覆盖 `CreateInstanceCopyPairRequest` 的全部入参，包括 `CreateTargetInstanceParameters` 列表及其嵌套子结构（Placement、SystemDisk、InstanceChargePrepaid、DataDisks、VirtualPrivateCloud、InternetAccessible、LoginSettings、EnhancedService 等）。
- 异步创建后通过轮询 `DescribeCopyPairs` 直到复制对状态生效（非 `INIT` 状态）。
- 所有 SDK 调用包裹在 `resource.Retry` 中（写路径 `WriteRetryTimeout`，读路径 `ReadRetryTimeout`），对 nil Response 做 `NonRetryableError` 防御。
- 代码风格与 `tencentcloud_igtm_strategy` 一致：单文件资源布局、`defer tccommon.LogElapsed/InconsistentCheck`、`log.Printf` 带资源名 `bdrc_instance_copy_pair`。
- 在 `connectivity/client.go` 和 `provider.go`/`provider.md` 中完成 bdrc 客户端与资源注册。

**Non-Goals:**
- 不实现 BDRC 数据源（datasource），本变更仅新增资源。
- 不管理 DISK 或 CFS 类型的复制对，本资源仅针对 CVM（INSTANCE）类型复制对。
- 不管理容灾保护组（Protect Group）资源——属于另一个资源范畴。
- 不管理容灾演练组（Drill Group）资源。
- 不处理 failover / rollback 等灾备切换操作——这些是运维操作，不属于 IaC 生命周期管理。

## Decisions

### D1. 资源 ID 设计：单 ID = `CopyPairId`
`CreateInstanceCopyPairResponse` 返回 `CopyPairIds []*string`（1~10 个）。由于一个 Terraform 资源实例对应一个复制对，我们在 Create 中取返回列表的第一项作为资源 ID（`CopyPairId`，格式 `cvmcopypair-xxxxxxxx`）。`DescribeCopyPairs`、`ModifyCopyPairAttribute`、`DeleteCopyPairs` 均以 `CopyPairId` 为键，无需复合 ID。

**备选方案**：复合 ID `<ProtectGroupId>#<CopyPairId>`。否决——`ModifyCopyPairAttribute` 和 `DeleteCopyPairs` 都不需要 `ProtectGroupId`，使用复合 ID 只增加解析负担而无价值。

### D2. `create_target_instance_parameters` 列表 Schema
`CreateTargetInstanceParameters` 是 `[]*CreateInstanceModel`，SDK 注释允许 1~10 个。但 Terraform 资源管理的语义是"一个资源块 = 一个复制对"。设计决策：

- Schema 中 `create_target_instance_parameters` 为 `TypeList`，`MaxItems: 1`，`MinItems: 1`。这样既忠实映射了 API 的列表语义（该参数本身就是数组），又限制了只创建一个复制对，与 Terraform 1:1 资源语义一致。
- 列表元素是 `TypeList`/`TypeSet` 嵌套 Resource，其子字段完整覆盖 `CreateInstanceModel` 的全部属性。

**备选方案**：将 `CreateInstanceModel` 的所有字段平铺到顶层。否决——`CreateTargetInstanceParameters` 是 API 的一级参数（不是响应字段），保持列表结构更贴近 API 语义，且 `source_instance_id`/`image_id` 等名称若无前缀容易与响应中同名字段混淆。使用列表嵌套结构使创建参数与读取参数（computed）在 schema 层面清晰分离。

### D3. 创建参数 vs 读取属性的命名隔离
`CreateInstanceModel` 的子字段（如 `disk_type`、`disk_size`、`delete_with_instance`）与 `CopyPair` 响应和 `DiskCopyPairForCvm` 响应中的字段存在重名。为避免冲突：

- 创建参数全部嵌套在 `create_target_instance_parameters` 列表的嵌套 Resource 内（`system_disk`、`data_disks`、`virtual_private_cloud` 等各为独立的嵌套结构）。
- 读取（computed）属性来自 `CopyPair` 响应，放在顶层（如 `copy_pair_state`、`source_region`、`target_region`、`percent` 等）。
- `copy_pair_id`（顶层，computed）= `d.Id()`，与 `CreateInstanceModel` 内的 `CopyPairId`（容灾演练用的参数）无冲突——后者属于创建参数列表内部。

### D4. `CreateInstanceCopyPair` 异步轮询
`CreateInstanceCopyPair` 为异步接口。创建成功后，资源 ID 已从 `CopyPairIds[0]` 获得，但复制对可能仍处于 `INIT` 状态。设计决策：

- 在 Create 的 retry 块外（获取 ID 后），调用一个内部轮询函数：循环调用 `DescribeCopyPairs`（`CopyPairType=INSTANCE`，`CopyPairIds=[id]`），检查 `CopyPairState`。
- 轮询成功条件：`CopyPairState` 不为 `INIT`（表示异步创建已推进到 `RUNNING`/`FULL_COPYING`/`NORMAL` 等状态）。
- 轮询超时：使用 `tccommon.ReadRetryTimeout` 作为单次 Describe 重试预算，整体轮询使用 `resource.Retry` 包裹。
- 轮询失败处理：若超时仍未脱离 `INIT`，返回 `NonRetryableError` 提示用户人工介入。

### D5. Update 路径：仅 `instance_copy_pair_name` 可变
`ModifyCopyPairAttribute` 接口仅接受 `CopyPairId`、`CopyPairType`、`CopyPairName` 三个参数。除 `instance_copy_pair_name`（映射 `CopyPairName`）外，`CreateInstanceCopyPair` 的所有其他入参（`protect_group_id`、`create_target_instance_parameters`、`client_token`、`recovery_point_objective`）均不可通过更新接口修改。

设计决策：
- `instance_copy_pair_name`：Optional，可更新（不设 ForceNew）。
- `protect_group_id`：Required，ForceNew（修改需重建）。
- `create_target_instance_parameters`：Required，ForceNew（CVM 创建参数变更需重建）。
- `client_token`：Optional，ForceNew。
- `recovery_point_objective`：Optional，ForceNew。
- Update 函数中，`mutableArgs = []string{"instance_copy_pair_name"}`，若检测到其他字段变更则返回 error。

### D6. Delete 路径
`DeleteCopyPairs` 需要 `CopyPairIds`（列表）、`CopyPairType`（固定 `INSTANCE`）、`DeleteTargetResource`（可选 bool）。

设计决策：
- 资源 schema 新增一个 computed/optional 字段 `delete_target_resource`（bool），默认不传（API 默认 `true`）。用户可显式设置。
- Delete 时构造 `CopyPairIds = [d.Id()]`，`CopyPairType = "INSTANCE"`。

### D7. Read 路径中的空响应处理
按代码规范，Read 中若 `DescribeCopyPairs` 返回空（`response == nil` 或 `len(CopyPairSet) == 0`），先 `log.Printf("[CRUD] bdrc_instance_copy_pair id=%s", d.Id())` 保留现场，再 `d.SetId("")` 返回 nil。

但该逻辑在 service 层 `DescribeBdrcInstanceCopyPairById` 中实现：service 层 retry 块内若返回空则返回 `NonRetryableError`（避免因短暂波动清空 id），resource 层在 service 返回 `(nil, nil)` 时才 `d.SetId("")`。这遵循 BDRC/通用规范第 14 条的数据源空响应处理原则，虽然本资源是 RESOURCE_KIND_GENERAL 而非 DATASOURCE，但 Read 逻辑保持一致：service 层 retry 内返回 NonRetryableError，service 层在 retry 失败路径打印 `[CRUD]` 日志。

### D8. connectivity 与 provider 注册
- `connectivity/client.go`：新增 `bdrcv20260330 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/bdrc/v20260330"` 导入，`TencentCloudClient` 结构体新增 `bdrcv20260330Conn *bdrcv20260330.Client` 字段，新增 `UseBdrcV20260330Client()` 方法（参照 `UseIgtmV20231024Client`，`NewClientProfile(300)`）。
- `provider.go`：新增 `svcBdrc "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/bdrc"` 导入别名，在 `ResourcesMap` 中注册 `"tencentcloud_bdrc_instance_copy_pair": bdrc.ResourceTencentCloudBdrcInstanceCopyPair()`。
- `provider.md`：新增 `Backup and Disaster Recovery(BDRC)` 章节及 `tencentcloud_bdrc_instance_copy_pair` 条目。

### D9. 测试策略
按代码生成要求，新增 terraform 资源使用 mock（gomonkey）方式对云 API 进行 mock，只进行业务代码逻辑的单元测试，不使用 terraform 测试套件。测试文件 `resource_tc_bdrc_instance_copy_pair_test.go` 中使用 `gomonkey` patch `UseBdrcV20260330Client` 返回的 client 方法，验证 Create/Read/Update/Delete 各函数的逻辑分支。

## Risks / Trade-offs

- **[Risk]** `CreateInstanceCopyPair` 一次可创建多个复制对（1~10），但 Terraform 资源语义为 1:1。→ **Mitigation**：Schema 中 `create_target_instance_parameters` 设 `MaxItems: 1`，只取返回 `CopyPairIds` 的第一项作为资源 ID，确保 1:1 语义。
- **[Risk]** 异步创建轮询可能因 API 延迟而超时。→ **Mitigation**：使用 `resource.Retry` 包裹轮询，超时返回明确错误信息提示人工介入，不会导致 state 不一致（ID 在轮询前已设置，若轮询失败 Terraform 会报告 create failed）。
- **[Risk]** `DescribeCopyPairs` 的 `CopyPairType` 必须传 `INSTANCE`，若传错可能查不到结果。→ **Mitigation**：在 Read service 层硬编码 `CopyPairType = "INSTANCE"`。
- **[Trade-off]** 创建参数嵌套在 `create_target_instance_parameters` 列表中而非平铺顶层。这使得 HCL 稍显嵌套，但保持了与 API 结构的忠实映射，避免了创建参数与 computed 属性的命名冲突。
- **[Trade-off]** `delete_target_resource` 作为 schema 字段而非 delete 时动态参数。这使其在 state 中可见，用户可明确控制删除行为。

## Migration Plan

纯新增变更，无需迁移：
1. 合并 bdrc 客户端绑定、服务目录、资源文件和 provider 注册代码。
2. 发布后用户通过添加 `resource "tencentcloud_bdrc_instance_copy_pair" "..." {}` 开始使用。
3. 无需修改任何已有资源或 state。

回滚：纯 revert 新增文件 + `provider.go`/`connectivity/client.go` 的注册行，无 state 变更需撤销。

## Open Questions

- 无需用户输入的待决问题。SDK 已提供全部所需 API，提案级决策（单 ID、创建参数嵌套列表、异步轮询、仅名称可更新）均已确定。
