## Context

`tencentcloud_cynosdb_readonly_instance` 资源当前的 `instance_name` 字段为 `ForceNew`，修改实例名称会触发资源销毁重建。腾讯云 TDSQL-C MySQL 提供了 `ModifyInstanceName` 接口（cynosdb v20190107），支持对已存在的实例修改名称，请求参数为 `InstanceId` 与 `InstanceName`，无需异步轮询（同步返回）。

现有代码结构：
- 资源文件：`tencentcloud/services/cynosdb/resource_tc_cynosdb_readonly_instance.go`
- 服务层：`tencentcloud/services/cynosdb/service_tencentcloud_cynosdb.go`（已有 `ModifyClusterName` 等同类方法可参考）
- 测试：`tencentcloud/services/cynosdb/resource_tc_cynosdb_readonly_instance_test.go`（现有 Terraform acc 测试套件）

## Goals / Non-Goals

**Goals:**
- 使 `instance_name` 字段支持就地更新，调用 `ModifyInstanceName` API 完成修改。
- 在服务层封装 `ModifyInstanceName` 调用，遵循 `WriteRetryTimeout` 重试与 `tccommon.RetryError` 错误包装的现有模式。
- 保持向后兼容：不改变现有 state 格式与 ID 规则。

**Non-Goals:**
- 不修改 `instance_name` 的 `Required` 属性。
- 不调整其它字段（如 `vpc_id`/`subnet_id` 仍不支持变更）。
- 不引入新的 schema 字段。

## Decisions

### Decision 1: 移除 instance_name 的 ForceNew

**选择**：将 `instance_name` 的 `ForceNew: true` 删除，使其成为可更新字段。

**理由**：`ModifyInstanceName` API 支持修改实例名称，无需重建资源。移除 `ForceNew` 是接入该 API 的前提。

**替代方案**：保留 `ForceNew` 并新增一个独立字段——会导致语义重复且用户体验割裂，弃用。

### Decision 2: 在 Update 方法中增加 instance_name 变更分支

**选择**：在 `resourceTencentCloudCynosdbReadonlyInstanceUpdate` 中新增 `if d.HasChange("instance_name")` 分支，调用 `cynosdbService.ModifyInstanceName(ctx, instanceId, instanceName)`。

**理由**：与现有 `instance_cpu_core`/`instance_memory_size`、维护窗口等分支保持一致的 `d.HasChange` 模式。

### Decision 3: 服务层方法签名与重试策略

**选择**：在 `service_tencentcloud_cynosdb.go` 中新增：
```go
func (me *CynosdbService) ModifyInstanceName(ctx context.Context, instanceId string, instanceName string) (errRet error)
```
内部使用 `resource.Retry(tccommon.WriteRetryTimeout, ...)` + `tccommon.RetryError(err)`，与 `ModifyClusterName` 方法保持一致。

**理由**：写操作遵循 provider 现有重试模式；请求参数仅 `InstanceId` 与 `InstanceName`，无需额外封装。

### Decision 4: 单元测试使用 gomonkey mock

**选择**：在 `resource_tc_cynosdb_readonly_instance_test.go` 中新增 `TestUnitCynosdbReadonlyInstance_UpdateInstanceName`，使用 gomonkey mock `UseCynosdbClient`、`ModifyInstanceName`、`DescribeInstances`（用于 Read）。

**理由**：本变更新增了 CRUD 逻辑分支，按照代码生成要求使用 mock 进行业务逻辑单元测试，可独立运行无需云资源。

## Risks / Trade-offs

- **[风险] 旧 state 中 instance_name 为 ForceNew，升级后行为变化** → 缓解：仅是从"重建"变为"就地更新"，对用户而言是行为增强，不会破坏已有配置；首次 apply 不会触发任何变更。
- **[风险] ModifyInstanceName 接口偶发失败** → 缓解：通过 `WriteRetryTimeout` 重试机制覆盖。
