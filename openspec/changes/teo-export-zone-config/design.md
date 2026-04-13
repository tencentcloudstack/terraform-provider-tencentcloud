## Context

当前 Terraform Provider for TencentCloud 已支持多种腾讯云服务资源，包括 EdgeOne (TEO) 服务。TEO 服务提供站点配置管理功能，但目前缺少通过 Terraform 导出站点配置的能力。用户需要通过手动调用 API 或使用腾讯云控制台来导出站点配置，这限制了基础设施即代码 (IaC) 的自动化管理。

现有的 TEO 资源遵循标准 Terraform Provider 模式：使用 tencentcloud-sdk-go 调用云 API，实现 CRUD 操作，并通过最终一致性重试机制处理异步操作。

## Goals / Non-Goals

**Goals:**
- 实现 `tencentcloud_teo_export_zone_config` 资源的完整 CRUD 操作
- 根据 CAPI 接口定义生成准确的 Schema
- 提供完整的测试覆盖（单元测试和验收测试）
- 确保与现有 TEO 资源的一致性和可维护性

**Non-Goals:**
- 不修改现有的 TEO 资源
- 不引入新的外部依赖（使用已有的 tencentcloud-sdk-go）
- 不涉及数据迁移或 schema 变更

## Decisions

**1. 资源命名和文件组织**
- 使用标准命名: `resource_tc_teo_export_zone_config.go`
- 放置在 `tencentcloud/services/teo/` 目录
- 测试文件: `resource_tc_teo_export_zone_config_test.go`
- 文档: `website/docs/r/teo_export_zone_config.md`
- 示例: `examples/resources/teo_export_zone_config/resource.tf`

**理由**: 遵循 Terraform Provider 现有命名规范，便于代码维护和用户理解。

**2. API 调用方式**
- 使用 tencentcloud-sdk-go 调用 TEO CAPI 接口
- 复用 `service_tencentcloud_teo.go` 中的现有服务层函数
- 如需新的 API 调用，遵循服务层模式封装

**理由**: 保持代码一致性，复用现有基础架构。

**3. Schema 定义**
- 根据 CAPI 接口的请求参数生成 Schema
- Required/Optional 属性与 CAPI 接口保持一致
- 使用标准 Terraform SDK v2 的 schema 类型

**理由**: 确保与云 API 的正确映射，避免参数错误。

**4. ID 构建方式**
- 如果需要复合 ID，使用 `zoneId#exportType` 格式（分隔符: #）
- 如果是简单资源，使用单一标识符

**理由**: 遵循 Terraform Provider 的复合 ID 模式。

**5. 错误处理和重试**
- 使用 `helper.Retry()` 处理最终一致性
- 使用 `defer tccommon.LogElapsed()` 记录耗时
- 使用 `defer tccommon.InconsistentCheck()` 检查不一致状态

**理由**: 确保异步操作的可靠性和可观测性。

**6. 测试策略**
- 单元测试: 测试各个 CRUD 函数的逻辑
- 验收测试: 使用 TF_ACC=1 进行端到端测试
- 测试数据: 使用真实的 TEO 资源

**理由**: 确保代码质量和功能正确性。

## Risks / Trade-offs

**风险 1**: CAPI 接口参数定义可能不完整或更新频繁
**缓解**: 使用最新的 CAPI 接口文档，并在测试中验证参数映射的正确性

**风险 2**: 异步操作可能导致最终一致性问题
**缓解**: 实现重试机制和状态检查，确保操作完成

**权衡**: 虽然增加了代码复杂度，但通过重试机制确保了可靠性

**风险 3**: 测试环境需要真实的 TEO 资源
**缓解**: 使用测试账号的 TEO 资源，并在文档中说明测试前提条件
