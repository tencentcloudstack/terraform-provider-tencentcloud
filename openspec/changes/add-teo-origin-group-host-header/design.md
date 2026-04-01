## Context

tencentcloud_teo_origin_group 资源用于管理腾讯云 TEO（TencentCloud EdgeOne）的源站组。当前实现支持配置源站类型（HTTP/HTTPS）、源站地址等基础信息。根据最新的 API 变更，CreateOriginGroup、DescribeOriginGroup 和 ModifyOriginGroup 接口已支持 HostHeader 参数，用于在回源时自定义 HTTP 请求的 Host Header。

HostHeader 参数仅在源站类型为 HTTP 时生效，且规则引擎修改 Host Header 配置的优先级高于源站组的 Host Header。这一功能对于 CDN/边缘计算场景非常重要，用户可以通过配置不同的 Host Header 来实现多域名回源、灰度发布等高级策略。

## Goals / Non-Goals

**Goals:**
- 在 tencentcloud_teo_origin_group 资源中添加 `host_header` 字段，支持用户配置回源 Host Header
- 确保字段与 API 接口定义一致（可选，string 类型）
- 更新所有 CRUD 操作函数，正确处理 HostHeader 参数的读写更新
- 添加完整的单元测试和验收测试，确保功能正确性
- 保持向后兼容，不影响现有用户配置

**Non-Goals:**
- 不涉及其他 TEO 资源的修改
- 不修改源站组的其他字段或行为
- 不实现规则引擎的 Host Header 配置（规则引擎优先级更高）
- 不添加额外的验证逻辑（由 API 服务端处理）

## Decisions

**1. Schema 字段命名选择**
- **决策**: 使用 `host_header`（snake_case）作为 Terraform schema 字段名
- **理由**: 遵循 Terraform Provider 的命名约定，使用 snake_case 保持一致性。API 参数名称为 `HostHeader`（PascalCase），在代码中进行映射转换。

**2. 字段属性定义**
- **决策**: `host_header` 定义为 `TypeString` 和 `Optional`
- **理由**: 根据 API 定义，HostHeader 是可选参数，且类型为 string。设置为 Optional 可以确保向后兼容，现有配置无需修改。

**3. 字段位置**
- **决策**: 将 `host_header` 字段添加在源站类型相关字段之后
- **理由**: 保持相关字段在一起，提高代码可读性。参考现有 schema 结构，将新增字段放在类型相关的 `origin_type` 字段附近。

**4. 测试策略**
- **决策**: 单元测试覆盖字段的读写逻辑，验收测试验证与 API 的集成
- **理由**: 单元测试确保 CRUD 函数正确处理字段，验收测试验证与实际 API 的交互。由于是可选字段，需要测试有值和无值两种场景。

**5. API 请求处理**
- **决策**: 仅在源站类型为 HTTP 时传递 HostHeader 参数
- **理由**: 根据 API 文档，HostHeader 仅在 Type = HTTP 时生效。虽然在代码中不强制验证（由服务端处理），但在测试中需要验证这一行为。

## Risks / Trade-offs

**1. API 行为变更风险**
- **风险**: 如果 API 服务端的 HostHeader 逻辑发生变化，可能导致 Provider 行为不一致
- **缓解**: 通过验收测试验证实际 API 行为，定期检查 API 文档更新

**2. 字段冲突风险**
- **风险**: 未来可能添加与 HostHeader 冲突的其他字段（如规则引擎配置）
- **缓解**: 保持字段命名清晰明确，在文档中说明与规则引擎的优先级关系

**3. 向后兼容性**
- **风险**: 如果未来 API 改变 HostHeader 的 Required 属性，可能需要调整 schema
- **缓解**: 当前实现使用 Optional，即使 API 变为 Required 也不影响现有配置，只需更新文档说明

**4. 测试覆盖不完整**
- **风险**: 测试未覆盖所有边界情况（如特殊字符、超长字符串等）
- **缓解**: 依赖 API 服务端的验证逻辑，Provider 层主要负责正确传递参数