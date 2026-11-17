## ADDED Requirements

### Requirement: Storage compress field

`tencentcloud_audit_track` 资源的 `storage` block SHALL 提供可选字段 `compress`（`TypeInt`, `Optional`），用于控制数据投递存储的压缩行为，取值为 `1`（压缩）或 `2`（不压缩）。当用户未配置该字段时，资源 SHALL NOT 向云 API 传值，由云服务端使用默认行为。

#### Scenario: Create with compress specified

- **WHEN** 用户在 `storage` block 中配置 `compress = 1`
- **THEN** 资源 Create 方法调用 `CreateAuditTrack` 时，`request.Storage.Compress` SHALL 被设置为 `1`
- **AND** 创建完成后 Read 方法 SHALL 将 `compress` 回填到 state

#### Scenario: Create without compress

- **WHEN** 用户在 `storage` block 中未配置 `compress`
- **THEN** 资源 Create 方法调用 `CreateAuditTrack` 时，`request.Storage.Compress` SHALL NOT 被设置（保持 nil）
- **AND** 创建完成后若云 API 返回 `Compress` 非 nil，Read 方法 SHALL 将其回填；若为 nil，Read 方法 SHALL 跳过回填

#### Scenario: Update compress

- **WHEN** 用户修改 `storage.compress` 的值
- **THEN** `d.HasChange("storage")` SHALL 返回 true
- **AND** 资源 Update 方法调用 `ModifyAuditTrack` 时，`request.Storage.Compress` SHALL 被设置为新值
- **AND** 更新完成后 Read 方法 SHALL 将新值回填到 state

#### Scenario: Read compress from cloud API

- **WHEN** 资源执行 Read 并调用 `DescribeAuditTrack`
- **AND** 云 API 返回 `response.Storage.Compress` 非 nil
- **THEN** Read 方法 SHALL 将该值回填到 `storage.compress` 字段

#### Scenario: Read with nil compress

- **WHEN** 资源执行 Read 并调用 `DescribeAuditTrack`
- **AND** 云 API 返回 `response.Storage` 为 nil 或 `Compress` 为 nil
- **THEN** Read 方法 SHALL 跳过 `compress` 字段的 set，不影响 `storage` block 中其他字段的回填

#### Scenario: Delete unaffected

- **WHEN** 资源执行 Delete 并调用 `DeleteAuditTrack`
- **THEN** 该方法仅使用 `TrackId`，SHALL NOT 涉及 `Storage.Compress` 参数
