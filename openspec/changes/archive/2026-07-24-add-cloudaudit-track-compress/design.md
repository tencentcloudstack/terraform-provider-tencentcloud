## Context

`tencentcloud_audit_track` 资源（`tencentcloud/services/cloudaudit/resource_tc_audit_track.go`）管理 CloudAudit 跟踪集，支持将审计日志投递到 cos/cls/ckafka。当前 `storage` block 包含 `storage_type`、`storage_region`、`storage_name`、`storage_prefix`、`storage_account_id`、`storage_app_id` 等字段，但缺少云 API 已支持的 `Compress`（压缩）配置。

云 SDK（`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cloudaudit/v20190319`）的 `Storage` 结构体已新增 `Compress *uint64` 字段（`1:压缩 2:不压缩`），且该字段通过 `CreateAuditTrackRequest.Storage`、`ModifyAuditTrackRequest.Storage`、`DescribeAuditTrackResponseParams.Storage` 三个接口均可用。`DeleteAuditTrackRequest` 仅含 `TrackId`，不涉及 `Storage`。

## Goals / Non-Goals

**Goals:**
- 在 `tencentcloud_audit_track` 资源的 `storage` block 中新增 `compress` 可选字段
- 在 Create / Read / Update 三个 CRUD 方法中正确映射 `compress` 字段与云 API 的 `Storage.Compress`
- 保持向后兼容：`compress` 为 `Optional`，不影响现有配置和 state

**Non-Goals:**
- 不修改 `DeleteAuditTrack` 相关逻辑（该接口不涉及 `Storage`）
- 不修改 `storage` block 中其他已有字段
- 不改变资源的 ID 规则、import 行为
- 不引入 `_extension.go` 文件

## Decisions

### 决策 1：`compress` 字段类型与位置
- **选择**：在 `storage` block（`TypeList`, `MaxItems:1`）的子 schema 中新增 `compress`，类型为 `schema.TypeInt`，`Optional: true`。
- **理由**：云 API 的 `Compress` 为 `*uint64`，使用 `schema.TypeInt` 与资源中已有的 `status`（同为 `uint64`）保持一致，便于用户书写 HCL。
- **取值约束**：`1` 表示压缩，`2` 表示不压缩（与云 API 文档一致）。在 schema `Description` 中注明取值含义，不做 `ValidateFunc` 硬校验（与现有同类字段风格一致）。

### 决策 2：Create / Update 中 compress 的映射
- **选择**：在现有 `if dMap, ok := helper.InterfacesHeadMap(d, "storage"); ok` 分支内，新增 `if v, ok := dMap["compress"]; ok { storage.Compress = helper.IntUint64(v.(int)) }`。
- **理由**：复用现有的 storage 读取逻辑，最小化改动；`helper.IntUint64` 已用于同文件的 `status` / `track_for_all_members` 字段。
- **Update 中的处理**：`ModifyAuditTrack` 在 `d.HasChange("storage")` 分支内执行，`compress` 作为 `storage` 子字段，其变更会触发 `d.HasChange("storage")` 为 true，无需单独判断。

### 决策 3：Read 中 compress 的回填
- **选择**：在 `if track.Storage != nil` 分支的 `storageMap` 中新增 `if track.Storage.Compress != nil { storageMap["compress"] = track.Storage.Compress }`。
- **理由**：遵循"先判断 nil 再 set"的硬约束，与同分支其他字段处理方式一致。

### 决策 4：compress 默认不设置
- **选择**：当用户未配置 `compress` 时，不向云 API 传值（`Compress` 保持 nil），由云服务端使用默认行为。
- **理由**：保持向后兼容，避免对现有资源产生意外变更。

## Risks / Trade-offs

- **[Risk] 云 API 返回的 `Compress` 可能为 nil（历史数据）** → 在 Read 中先判断 `track.Storage.Compress != nil` 再 set，nil 时跳过，不影响其他字段回填。
- **[Risk] 用户配置非法取值（非 1/2）** → 云 API 会返回错误，由现有 retry / 错误处理机制透传给用户；不在 Terraform 侧做硬校验，与现有字段风格一致。
- **[Trade-off] `compress` 为 Optional 但不设置时不传值** → 行为依赖云服务端默认值，可接受，因为这是新增字段。
