## Why

当前 `tencentcloud_lite_hbase_instance` 资源将 `disk_size` 参数硬编码为不可修改（在 Update 函数的 `immutableArgs` 列表中，第 380 行），用户修改 `disk_size` 会直接返回 `argument cannot be changed` 错误，无法通过 Terraform 对轻量 HBase 实例进行磁盘扩容。这导致用户必须手动通过控制台或 API 扩容，破坏基础设施即代码的一致性原则，并引发 Terraform 状态与实际配置漂移。

腾讯云 EMR 服务已提供 `ModifySLInstance` 接口（https://cloud.tencent.com/document/api/589/110869），支持对轻量 HBase 实例的存储空间进行扩容（只允许扩容，不允许缩容）。该接口为异步接口，调用后需通过 `DescribeSLInstance` 接口（https://cloud.tencent.com/document/product/589/110871）轮询实例状态：`Status` 值为 22 表示变配中（ADJUSTING），值为 2 表示运行中（变配完成，RUNNING）。

vendor 目录中的 SDK 已包含 `ModifySLInstanceRequest` 结构体及其 `DiskSize` 参数，且 `ModifySLInstanceWithContext` 方法已存在，无需升级 SDK。

## What Changes

- 从 Update 函数的 `immutableArgs` 列表中移除 `disk_size`，使其变为可更新字段
- 在 Update 函数中新增 `disk_size` 变更处理逻辑：当检测到 `d.HasChange("disk_size")` 时，构造 `ModifySLInstanceRequest` 并设置 `InstanceId` 与目标 `DiskSize`，调用 `ModifySLInstanceWithContext`
- 调用 `ModifySLInstance` 后，复用现有的 `SLInstanceStateRefreshFunc` 异步等待机制，等待实例状态从 `22`（变配中）变为 `2`（运行中/变配完成）
- 同步更新资源 `.md` 示例文档，补充 `disk_size` 扩容的使用说明
- 补充单元测试用例覆盖 `disk_size` 变更逻辑（修改现有资源，使用 Terraform 测试套件补充测试用例）

## Capabilities

### New Capabilities
- `emr-lite-hbase-instance-resource`: 轻量 HBase 实例资源管理，新增 `disk_size` 字段的在线扩容（变配）能力

### Modified Capabilities
<!-- 无已有 spec 需要修改 -->

## Impact

### 受影响的代码
- `tencentcloud/services/emr/resource_tc_lite_hbase_instance.go` - Update 函数：从 immutableArgs 移除 `disk_size`，新增 `disk_size` 变更处理块（构造 ModifySLInstanceRequest、调用接口、异步等待状态变为 2）
- `tencentcloud/services/emr/resource_tc_lite_hbase_instance_test.go` - 补充 `disk_size` 变更场景的验收测试用例
- `tencentcloud/services/emr/resource_tc_lite_hbase_instance.md` - 同步更新示例文档，补充 disk_size 扩容说明

### API 依赖
- `ModifySLInstance`（https://cloud.tencent.com/document/api/589/110869）- 已存在于 vendor SDK 中，请求参数包含 `InstanceId`、`Zone`、`NodeNum`、`ClientToken`、`DiskSize`，本次仅使用 `InstanceId` 和 `DiskSize`
- `DescribeSLInstance`（https://cloud.tencent.com/document/product/589/110871）- 已被现有 `SLInstanceStateRefreshFunc` / `GetSLInstanceStatus` 用于状态轮询，返回值 `Status` 字段：2=运行中，22=变配中

### 向后兼容性
- 完全向后兼容：不修改 `disk_size` 的现有配置行为与之前完全一致
- 行为变更：修改 `disk_size` 从"返回错误"变为"原地扩容更新"，符合用户预期，属正向增强
- 云 API 约束：`DiskSize` 必须是 100 的倍数，且不允许比当前存储空间小（仅允许扩容），该约束由云 API 校验
