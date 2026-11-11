## Context

`tencentcloud_cynosdb_cluster` 与 `tencentcloud_cynosdb_cluster_v2` 资源中的 `instance_name` 字段当前为 `Computed`，用户无法通过 Terraform 修改实例名称，只能通过控制台单独修改，导致 state 漂移。腾讯云 TDSQL-C MySQL 提供了 `ModifyInstanceName` 接口（cynosdb v20190107），支持对已存在的实例修改名称，请求参数为 `InstanceId` 与 `InstanceName`，同步返回。

服务层方法 `ModifyInstanceName(ctx, instanceId, instanceName)` 已在 `service_tencentcloud_cynosdb.go` 中实现（随 `tencentcloud_cynosdb_readonly_instance` 资源一起合入 master），内部使用 `resource.Retry(tccommon.WriteRetryTimeout, ...)` + `tccommon.RetryError(err)`，可直接复用。

现有代码结构：
- 资源文件：`tencentcloud/services/cynosdb/resource_tc_cynosdb_cluster.go`、`resource_tc_cynosdb_cluster_v2.go`
- Schema 公共函数：`tencentcloud/services/cynosdb/extension_cynosdb.go` 的 `TencentCynosdbInstanceBaseInfo()`
- 服务层：`tencentcloud/services/cynosdb/service_tencentcloud_cynosdb.go`（已有 `ModifyInstanceName` 方法）

## Goals / Non-Goals

**Goals:**
- 使 `instance_name` 字段支持就地更新，调用 `ModifyInstanceName` API 完成修改。
- 复用已有的 `ModifyInstanceName` 服务层方法，无需新增。
- 保持向后兼容：不改变现有 state 格式与 ID 规则；不填写 `instance_name` 时行为不变。

**Non-Goals:**
- 不修改 `instance_name` 的 `Computed` 属性（保留 `Computed: true`，同时新增 `Optional: true`）。
- 不调整其它字段的 schema。
- 不引入新的 schema 字段。
- 不新增服务层方法（直接复用 `ModifyInstanceName`）。

## Decisions

### Decision 1: instance_name 改为 Optional + Computed

**选择**：在 `TencentCynosdbInstanceBaseInfo()` 中将 `instance_name` 从 `Computed: true` 改为 `Optional: true, Computed: true`。

**理由**：`ModifyInstanceName` API 支持修改实例名称。保留 `Computed` 保证 Read 仍会从云上回填，新增 `Optional` 允许用户显式配置并触发 Update。

### Decision 2: 在 Update 方法中增加 instance_name 变更分支

**选择**：在 `resourceTencentCloudCynosdbClusterUpdate` 与 `resourceTencentCloudCynosdbClusterV2Update` 中新增 `if d.HasChange("instance_name")` 分支，调用 `cynosdbService.ModifyInstanceName(ctx, instanceId, instanceName)`。

**理由**：与现有 `cluster_name`、`auto_renew_flag` 等分支保持一致的 `d.HasChange` 模式；`instanceId` 已在 update 函数开头从 `d.Get("instance_id").(string)` 获取，直接复用。

### Decision 3: 复用已有服务层方法

**选择**：直接调用 `cynosdbService.ModifyInstanceName(ctx, instanceId, instanceName)`。

**理由**：该方法随 `tencentcloud_cynosdb_readonly_instance` 资源一起合入 master，内部已使用 `WriteRetryTimeout` 重试与 `tccommon.RetryError` 错误包装，符合 provider 现有重试模式，无需重复实现。

### Decision 4: 测试使用 terraform acc 测试套件

**选择**：在 `resource_tc_cynosdb_cluster_test.go` 与 `resource_tc_cynosdb_cluster_v2_test.go` 中新增 `TestAccTencentCloudCynosdbClusterResourceUpdateInstanceName` / `TestAccTencentCloudCynosdbClusterV2ResourceUpdateInstanceName`，使用 terraform acc 测试套件，通过创建集群后修改 `instance_name` 验证 Update 路径。

**理由**：cluster 与 cluster_v2 为现有资源，按代码生成要求"对于修改现有的 terraform 资源的需求，仍然使用 terraform 的测试套件补充测试用例"。新增用例不使用 go test 执行。

## Risks / Trade-offs

- **[风险] 旧 state 中 instance_name 为 Computed，升级后变为 Optional+Computed** → 缓解：`Computed` 保留，不填写该字段时行为与原来完全一致；首次 apply 若用户未配置 `instance_name` 不会触发任何变更。
- **[风险] ModifyInstanceName 接口偶发失败** → 缓解：通过 `WriteRetryTimeout` 重试机制覆盖。
