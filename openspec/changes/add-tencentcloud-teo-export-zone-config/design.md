## Context

TEO (TencentCloud EdgeOne) 是腾讯云的边缘加速和网络安全服务，需要通过 Terraform Provider 支持完整的资源管理能力。当前 TEO 服务可能已有其他 Resource，本次变更旨在添加导出站点配置的 Resource，以支持用户通过 Terraform 管理 TEO 站点配置的导出操作。

Terraform Provider 采用 Go 语言编写，使用 Terraform Plugin SDK v2，通过 tencentcloud-sdk-go 调用腾讯云 CAPI 接口。项目遵循标准的文件组织结构和编码规范，所有 Resource 必须提供完整的 CRUD 操作、单元测试和验收测试。

## Goals / Non-Goals

**Goals:**
- 实现完整的 `tencentcloud_teo_export_zone_config` Resource，包含 Create、Read、Update、Delete 操作
- 根据 CAPI 接口定义生成准确的 Resource Schema，保持参数的 Required/Optional 属性一致
- 提供完整的单元测试覆盖和验收测试
- 遵循现有代码规范和最佳实践（错误处理、重试机制、ID 格式等）
- 支持异步操作的 Timeout 配置
- 提供 Resource 文档和使用示例

**Non-Goals:**
- 不涉及修改 TEO 服务已有的其他 Resource
- 不涉及 Provider 架构级别的变更
- 不涉及跨服务或跨模块的变更

## Decisions

**Resource Schema 设计：**
- 使用 Terraform Plugin SDK v2 的 `schema.Resource` 构建 Schema
- 根据 CAPI 接口的请求参数定义 Schema 属性，准确映射 Required/Optional/Computed 标记
- 使用 `schema.TypeString`、`schema.TypeInt` 等基础类型，复杂嵌套结构使用 `schema.TypeList` + `schema.Resource`
- 导出配置内容使用 `schema.TypeList` + `schema.Schema{Type: schema.TypeString}` 存储 JSON 配置

**CRUD 操作实现：**
- **Create**: 调用 CAPI 创建导出任务，返回任务 ID，使用 `resource.StateRefreshFunc` 轮询任务状态直到完成
- **Read**: 根据 Resource ID（可能包含 Zone ID 和导出 ID）查询导出结果，刷新状态
- **Update**: 如果 CAPI 支持更新导出参数，则实现 Update 操作；否则返回无操作
- **Delete**: 如果 CAPI 提供删除导出记录的接口，则实现删除操作；否则设置 `d.SetId("")` 逻辑删除

**错误处理和重试：**
- 使用 `tccommon.LogElapsed(ctx)` 和 `tccommon.InconsistentCheck(d, meta)` 处理日志记录和最终一致性检查
- 调用 CAPI 时使用 `helper.Retry()` 进行重试，处理网络抖动和临时性错误
- 使用 `log.Printf` 记录详细的调试信息，便于问题排查

**Resource ID 设计：**
- 复合 ID 格式：`zoneId#exportId` 或 `exportId`（根据 CAPI 返回的数据结构确定）
- 使用 `helper.HashResource` 确保 ID 的唯一性和稳定性

**异步操作支持：**
- 在 Schema 中声明 `schema.Resource{Timeouts: &schema.ResourceTimeout{...}}` 支持自定义超时配置
- 使用 `resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), refreshFunc)` 等待异步操作完成

**测试策略：**
- **单元测试**：测试 Resource 的 CRUD 逻辑，使用 mock CAPI 响应，不依赖真实 API
- **验收测试**：使用真实 TEO 环境执行完整流程，需要设置 `TF_ACC=1` 和相关环境变量

## Risks / Trade-offs

**API 复杂性：**
- [Risk] 导出站点配置的 API 可能返回复杂的嵌套 JSON 结构，映射到 Terraform Schema 可能比较困难
  - **Mitigation**: 仔细分析 CAPI 接口文档，设计合理的 Schema 结构，必要时使用 `schema.TypeMap` 或自定义类型

**异步操作耗时：**
- [Risk] 导出站点配置可能是一个耗时的异步操作，默认超时时间可能不够
  - **Mitigation**: 在 Schema 中提供合理的默认超时时间，并允许用户自定义配置

**向后兼容性：**
- [Risk] CAPI 接口可能在未来的版本中发生变化，导致 Terraform Resource 不兼容
  - **Mitigation**: 严格遵循 Provider 的向后兼容原则，新版本只添加 Optional 字段，不修改或删除已有字段

**测试环境依赖：**
- [Risk] 验收测试需要真实的 TEO 环境和有效凭证，可能在某些环境下无法执行
  - **Mitigation**: 单元测试提供完整的逻辑覆盖，验收测试作为补充验证，使用环境变量控制执行

**ID 格式设计：**
- [Risk] 复合 ID 格式可能不够灵活，未来 API 变更时需要迁移
  - **Mitigation**: 选择足够抽象的 ID 格式（如 UUID），减少对 API 具体实现的依赖
