## Why

云产品 CloudAudit 的 `tencentcloud_audit_track` 资源目前不支持配置数据投递存储的压缩选项。云 API 已在 `Storage` 结构体中新增 `Compress` 字段（`1:压缩 2:不压缩`），允许用户控制投递到 cos/cls/ckafka 的审计日志是否压缩。为保持 Terraform 资源与云 API 能力同步，需要在资源中暴露该参数。

## What Changes

- 在 `tencentcloud_audit_track` 资源的 `storage` block 中新增 `compress` 子字段（`TypeInt`, `Optional`），取值为 `1`（压缩）或 `2`（不压缩）。
- 在资源 Create 方法（`CreateAuditTrack`）中，将 `compress` 映射到 `request.Storage.Compress`。
- 在资源 Read 方法（`DescribeAuditTrack`）中，将 `response.Storage.Compress` 回填到 `compress` 字段。
- 在资源 Update 方法（`ModifyAuditTrack`）中，将 `compress` 映射到 `request.Storage.Compress`。
- 资源 Delete 方法（`DeleteAuditTrack`）仅接收 `TrackId`，不涉及 `Storage` 参数，无需变更。
- 更新资源文档 `resource_tc_audit_track.md` 补充 `compress` 字段说明。

## Capabilities

### New Capabilities
<!-- 无新增能力，复用现有资源 -->

### Modified Capabilities
- `cloudaudit-track-resource`: 在 `storage` block 中新增 `compress` 可选参数，支持配置数据投递存储的压缩选项

## Impact

- **代码**:
  - `tencentcloud/services/cloudaudit/resource_tc_audit_track.go`：Schema、Create、Read、Update 方法
  - `tencentcloud/services/cloudaudit/resource_tc_audit_track_test.go`：补充单元测试用例
- **文档**: `tencentcloud/services/cloudaudit/resource_tc_audit_track.md`
- **依赖**: 已通过 vendor 引入 `Compress *uint64` 字段（`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cloudaudit/v20190319`）
- **兼容性**: 新增 `Optional` 字段，向后兼容，不影响现有 TF 配置和 state
