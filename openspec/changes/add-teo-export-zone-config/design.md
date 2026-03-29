## Context

Terraform Provider for TencentCloud 需要添加 TEO (TencentCloud EdgeOne) 站点配置导出功能。当前 Provider 已经支持 TEO 服务的基础资源（如 tencentcloud_teo_zone），但缺少导出完整站点配置的数据源资源。TEO API 提供了获取站点详细配置的接口，需要在 Provider 中封装此功能，使其能够通过 Terraform 数据源使用。

## Goals / Non-Goals

**Goals:**
- 实现一个数据源资源，能够导出 TEO 站点的完整配置信息
- 支持通过 ZoneId 或 ZoneName 作为查询参数
- 提供完整的配置信息，包括基础配置、加速策略、安全规则等
- 遵循 Terraform Provider 的最佳实践和现有代码结构
- 包含完整的单元测试和验收测试

**Non-Goals:**
- 不修改站点配置（只读操作）
- 不涉及配置的创建、更新或删除
- 不改变现有 TEO 资源的行为

## Decisions

**1. API 选择**
- 决策：使用腾讯云 TEO API 的 DescribeZoneConfig 接口（或类似接口）来获取站点配置信息
- 理由：此接口提供了站点配置的完整信息，满足导出需求
- 替代方案：使用多个 API 分别获取不同配置项 → 复杂度更高，性能较差

**2. 查询参数设计**
- 决策：支持 ZoneId 和 ZoneName 两个可选参数，至少提供其中一个
- 理由：灵活性高，用户可以根据自己掌握的信息进行查询
- 替代方案：只支持 ZoneId → 简单但灵活性不足

**3. 配置信息导出结构**
- 决策：将配置信息以嵌套的 schema 结构导出，保持 API 原始结构的层次关系
- 理由：配置信息较为复杂，保持层次关系便于用户理解和使用
- 替代方案：扁平化所有字段 → 失去结构化信息，不便于维护

**4. 错误处理**
- 决策：使用 Provider 标准的错误处理模式，包括重试机制和详细的错误信息
- 理由：TEO API 可能存在最终一致性问题，重试机制可以提高稳定性

**5. 文档和测试**
- 决策：提供完整的 website 文档和示例代码，包含单元测试和验收测试
- 理由：确保用户能够正确使用此数据源，保证代码质量

## Risks / Trade-offs

**Risk: TEO API 返回的配置信息结构复杂且可能频繁变化**
- Mitigation：使用灵活的 schema 设计，对未知字段进行忽略，避免因 API 变更导致 Provider 失败

**Risk: 大型站点配置信息量大，可能影响 Terraform plan/apply 性能**
- Mitigation：这是数据源资源，只在需要时读取，不会频繁执行；如果确实存在性能问题，可以考虑添加配置选项控制导出的详细程度

**Risk: API 调用可能因权限或网络问题失败**
- Mitigation：使用标准错误处理和重试机制，提供清晰的错误信息帮助用户排查问题

## Migration Plan

这是一个新增的数据源资源，不涉及迁移：
1. 用户可以立即开始使用新的数据源资源
2. 不影响现有资源和配置
3. 如果后续需要扩展功能，可以在现有基础上迭代

## Open Questions

无
