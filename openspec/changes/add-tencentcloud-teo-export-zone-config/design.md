# Design for tencentcloud_teo_export_zone_config

## Context

Terraform Provider for TencentCloud 是一个使用 Terraform Plugin SDK v2 开发的 Provider，用于管理腾讯云资源。目前 TEO（TencentCloud EdgeOne）服务中缺少导出站点配置的功能。

TEO 提供了 `ExportZoneConfig` API，用于导出站点的配置信息，包括七层加速配置等。这个 API 支持通过 ZoneId 和可选的 Types 参数来指定要导出的配置类型。

当前状态：
- TEO 服务已有多个 resources 和 datasources
- 缺少导出站点配置的 resource
- 用户需要能够将站点配置导出为 JSON 格式，便于备份和迁移

## Goals / Non-Goals

**Goals:**
- 添加 `tencentcloud_teo_export_zone_config` resource，支持导出站点配置
- 实现完整的 CRUD 操作（Create、Read、Update、Delete）
- 支持通过 ZoneId（必填）和 Types（可选）参数导出配置
- 返回导出的配置内容（JSON 格式）
- 提供单元测试和验收测试
- 添加资源文档和使用样例

**Non-Goals:**
- 不支持配置导入功能（仅导出）
- 不支持配置修改功能（只读 resource）
- 不支持增量导出或版本控制

## Decisions

### Resource Schema 设计

**Decision**: 使用 `tencentcloud_teo_export_zone_config` 作为资源名称，参数映射如下：
- `zone_id` (string, Required): 映射 CAPI 的 `ZoneId`
- `types` (list of string, Optional): 映射 CAPI 的 `Types`
- `content` (string, Computed): 映射 CAPI 响应的 `Content`

**Rationale**:
- 遵循 Terraform Provider 的命名规范
- 参数命名使用 snake_case，符合 Go 和 Terraform 习惯
- `content` 设置为 Computed，因为它是 API 返回的结果

**Alternative considered**:
- 使用 `export_zone_config` 作为参数名 → 不采用，避免与 resource 名称混淆
- 将 `content` 设为 Optional → 不采用，因为 content 是只读的，不应该由用户输入

### CRUD 操作实现

**Decision**: 实现 Create、Read、Update、Delete 四个操作函数，但实际逻辑如下：
- **Create**: 调用 `ExportZoneConfig` API 导出配置
- **Read**: 重新调用 `ExportZoneConfig` API 获取最新配置
- **Update**: 重新调用 `ExportZoneConfig` API 获取配置（如果 Types 参数变化）
- **Delete**: 从 state 中删除资源（不从云平台删除，因为这是导出操作）

**Rationale**:
- 导出操作不创建或修改云平台资源，只是读取配置
- Update 时如果 Types 参数变化，需要重新导出
- Delete 操作仅清理本地 state，符合 Terraform 的语义

**Alternative considered**:
- 只实现 Read 操作（类似 datasource）→ 不采用，用户希望在 state 中管理导出配置的快照
- Delete 时调用删除 API → 不采用，导出操作本身不需要删除

### 测试策略

**Decision**: 
- 单元测试：mock API 调用，测试 schema 定义和错误处理
- 验收测试：使用真实的 TEO 站点进行端到端测试（需要 TF_ACC=1）

**Rationale**:
- 单元测试确保代码逻辑正确
- 验收测试确保与真实 API 的集成正常

**Alternative considered**:
- 仅使用验收测试 → 不采用，运行速度慢且需要真实的云资源

### 错误处理

**Decision**:
- 使用 `helper.Retry()` 实现最终一致性重试
- 使用 `defer tccommon.LogElapsed()` 记录操作耗时
- 使用 `defer tccommon.InconsistentCheck()` 检查状态一致性
- API 错误直接返回，由 Terraform 处理

**Rationale**:
- 遵循 Provider 的通用错误处理模式
- 提供详细的日志和错误信息

## Risks / Trade-offs

**Risk 1**: API 响应的 Content 字段可能很大，导致 Terraform state 文件过大
- **Mitigation**: 建议用户使用 Types 参数指定需要导出的配置类型，避免导出所有配置

**Risk 2**: 导出的配置内容是 JSON 字符串，用户需要自行解析
- **Mitigation**: 在文档中提供示例，说明如何解析和使用导出的配置

**Trade-off**: 导出操作不会创建或修改云平台资源，但在 Terraform 中作为 resource 管理
- **Rationale**: 用户可以在 state 中跟踪导出操作的快照，便于管理和审计

**Trade-off**: Delete 操作仅清理 state，不调用云 API
- **Rationale**: 导出操作本身不需要删除，符合 Terraform 的语义

## Migration Plan

由于这是一个新 resource，不需要迁移计划。用户可以直接开始使用新的 resource。

## Open Questions

无。所有技术决策已经明确。
