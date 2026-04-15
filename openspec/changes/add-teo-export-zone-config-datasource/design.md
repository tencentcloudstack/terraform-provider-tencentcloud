## Context

TEO (Tencent Edge One) 是腾讯云的边缘安全加速平台服务。目前 Terraform Provider 已经支持了 TEO 的部分资源管理功能，但缺少配置导出的能力。用户需要能够导出站点配置进行备份、迁移或版本控制。

当前的 TEO 服务实现位于 `tencentcloud/services/teo/` 目录，遵循标准的 Terraform Provider 模式。已有的 TEO 数据源可以作为参考模式。

## Goals / Non-Goals

**Goals:**
- 实现一个 Terraform 数据源，能够调用 `ExportZoneConfig` API 导出站点配置
- 支持按站点 ID 导出配置
- 支持可选的配置类型过滤
- 提供完整的测试用例和文档

**Non-Goals:**
- 不支持配置导入功能
- 不支持配置差异对比
- 不涉及配置的解析和验证

## Decisions

1. **使用标准 Terraform 数据源模式**
   - 采用与其他 TEO 数据源相同的架构模式
   - 文件命名遵循 `data_source_tc_teo_export_zone_config.go` 规范
   - 使用 Terraform Plugin SDK v2 的 `datasource.DataSource` 接口

2. **Schema 设计**
   - `zone_id`：必填参数，字符串类型，指定要导出的站点 ID
   - `types`：可选参数，列表类型，支持配置类型过滤（如 `["L7AccelerationConfig"]`）
   - `content`：输出参数，字符串类型，导出的配置内容（JSON 格式）

3. **错误处理和重试机制**
   - 使用 `helper.Retry()` 处理最终一致性问题
   - 采用标准的错误处理模式：`defer tccommon.LogElapsed()` 和 `defer tccommon.InconsistentCheck()`
   - 对于 API 错误，返回清晰的错误信息

4. **测试策略**
   - 使用 mock 方式测试云 API 调用（使用 gomonkey）
   - 不使用真实的 Terraform 测试套件，避免需要真实的云资源
   - 覆盖正常流程和错误场景

5. **文档生成**
   - 通过 `make doc` 命令自动生成文档到 `website/docs/r/` 目录
   - 文档包含参数说明和使用示例

## Risks / Trade-offs

1. **导出配置可能很大**
   - [Risk] 如果站点配置很大，导出的 JSON 内容可能超出 Terraform 的限制
   - [Mitigation] 在文档中提醒用户使用 `types` 参数过滤需要的配置类型，减少导出内容

2. **API 兼容性**
   - [Risk] ExportZoneConfig API 可能会有变更，导致参数或返回值变化
   - [Mitigation] 严格依赖 vendor 模式管理的 SDK 版本，确保 API 调用的稳定性

3. **异步操作处理**
   - [Risk] ExportZoneConfig 是同步 API，不需要轮询
   - [Mitigation] 无需特殊处理，直接调用并返回结果
