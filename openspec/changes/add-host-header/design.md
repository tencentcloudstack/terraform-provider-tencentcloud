## Context

TEO（TencentCloud EdgeOne）源站组资源 `tencentcloud_teo_origin_group` 当前不支持配置回源 Host Header。根据 CreateOriginGroup API 的最新更新，新增了 `HostHeader` 参数，允许用户指定回源时使用的 Host Header。该参数仅当源站类型为 HTTP 时生效，且规则引擎修改 Host Header 配置的优先级高于源站组的 Host Header。

当前资源实现位于 `tencentcloud/services/teo/resource_tc_teo_origin_group.go`，使用 Terraform Plugin SDK v2 构建，通过 tencentcloud-sdk-go 调用云 API。

## Goals / Non-Goals

**Goals:**
- 在 tencentcloud_teo_origin_group 资源的 Schema 中新增 `host_header` 字段（string，Optional）
- 在 Create 和 Update 函数中将该字段传递给 CreateOriginGroup API
- 在 Read 函数中从 API 响应读取并填充该字段到 state
- 添加相应的单元测试和验收测试用例，确保字段正确读写

**Non-Goals:**
- 不修改现有的 Schema 结构（仅新增字段，不删除或修改现有字段）
- 不涉及规则引擎或其他高级功能
- 不修改资源的 ID 或其他核心行为

## Decisions

**Schema 设计:**
- 字段名: `host_header` (Go 风格的 snake_case)
- 类型: `schema.TypeString`
- 可选性: `Optional` (根据 API 定义)
- 无默认值: 让用户明确指定或不指定
- 不使用 `Computed` 或 `ForceNew`，因为这是一个可更新参数

**CRUD 函数修改:**
- **Create**: 在调用 CreateOriginGroup 时，如果 `host_header` 值非空，将其添加到请求参数
- **Read**: 从 API 响应中解析 `HostHeader` 字段并设置到 state
- **Update**: 调用 CreateOriginGroup API（TEO API 设计为创建/更新同一个接口），传递包括 `host_header` 在内的所有参数
- **Delete**: 无需修改（HostHeader 参数不影响删除操作）

**测试策略:**
- 单元测试: 在 `resource_tc_teo_origin_group_test.go` 中添加测试用例，验证 Schema 字段定义
- 验收测试: 创建测试资源，设置 `host_header`，验证 API 调用和 state 更新正确

**API 兼容性:**
- 假设 tencentcloud-sdk-go 已包含 CreateOriginGroup API 的 HostHeader 参数支持
- 如果 SDK 版本过旧，需要更新 vendor 依赖

## Risks / Trade-offs

**风险 1**: SDK 版本不包含 HostHeader 参数
- **缓解措施**: 实现前检查 tencentcloud-sdk-go 版本，必要时更新 vendor 依赖

**风险 2**: API 响应字段名可能不一致
- **缓解措施**: 参考 API 文档和现有代码模式，使用正确的字段映射

**风险 3**: HostHeader 字段与规则引擎的 Host Header 配置冲突
- **缓解措施**: 在文档中说明优先级（规则引擎 > 源站组），这不是代码层面的问题

**权衡**: 简单的参数添加 vs 复杂的条件逻辑
- **决策**: 使用简单的可选参数，不在代码层面添加额外的验证逻辑（如检查 Type 是否为 HTTP），因为这是业务层面的约束，由 API 负责验证

## Migration Plan

**部署步骤:**
1. 更新资源代码（Schema 和 CRUD 函数）
2. 更新测试代码
3. 运行单元测试确保通过
4. 运行验收测试（需要真实凭证）
5. 更新资源文档（website/docs/r/teo_origin_group.html.md）
6. 提交代码变更

**回滚策略:**
- 这是一个纯增量的可选字段，不会破坏现有配置
- 如果发现问题，可以简单的删除该字段相关代码
- 已创建的资源不会受到影响（新字段为可选，不影响已有 state）

## Open Questions

无。这是一个简单的字段添加任务，技术路径清晰。