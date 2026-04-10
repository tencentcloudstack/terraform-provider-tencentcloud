## Context

TencentCloud Terraform Provider 的 `tencentcloud_teo_function` 资源目前实现了基本的 CRUD 操作，支持创建、读取、更新和删除 TEO（EdgeOne）函数。当前资源包含基本的字段如 name、remark、content 等，能够满足简单的函数管理需求。

然而，在实际的企业级应用场景中，用户需要更丰富的函数配置能力：
1. **环境变量**：函数需要在不同环境（开发、测试、生产）中使用不同的配置参数
2. **规则绑定**：需要根据请求条件动态触发函数执行
3. **区域选择**：为了合规性和性能优化，需要控制函数的部署区域

TEO SDK 提供了 `FunctionEnvironmentVariable`、`FunctionRule` 和 `FunctionRegionSelection` 等结构体来支持这些高级功能，但目前 Terraform Provider 中并未暴露这些功能。

## Goals / Non-Goals

**Goals:**
- 在 `tencentcloud_teo_function` 资源中添加 `environment_variables` 字段，支持配置函数环境变量
- 在 `tencentcloud_teo_function` 资源中添加 `rules` 字段，支持配置函数规则
- 在 `tencentcloud_teo_function` 资源中添加 `region_selection` 字段，支持配置函数部署区域
- 确保所有新增字段与现有的 CRUD 操作完整集成
- 保持向后兼容性，不影响现有用户配置

**Non-Goals:**
- 不修改 TEO SDK 本身
- 不改变现有资源的字段定义（仅新增字段）
- 不影响 TEO 服务的其他资源
- 不修改 Provider 的核心框架代码

## Decisions

### 1. Schema 字段设计

**Decision:** 使用 Terraform SDK v2 的 `TypeList` 和 `TypeMap` 来表示复杂结构

**Rationale:**
- `environment_variables` 使用 `TypeList` + `MaxItems: 50` 限制数量，每个元素包含 key、value、type 三个字段
- `rules` 使用 `TypeList` + `MaxItems: 100` 限制数量，每个元素包含 rule_id、priority、conditions、actions 字段
- `region_selection` 使用 `TypeList` 存储 region codes，不限制数量

**Alternatives considered:**
- 使用 `TypeMap`：更简洁但无法保持顺序，且难以进行复杂的嵌套结构验证
- 使用嵌套的 `TypeSet`：可以自动去重但无法保持顺序，不利于规则的优先级控制

### 2. API 调用策略

**Decision:** 将新增字段映射到 TEO SDK 的对应结构体，在 Create 和 Update 操作中传递给 API

**Rationale:**
- `environment_variables` 映射到 `FunctionEnvironmentVariable` 结构体数组
- `rules` 映射到 `FunctionRule` 结构体数组
- `region_selection` 映射到 `FunctionRegionSelection` 结构体

**API 集成点:**
- Create 操作：在 `CreateFunctionRequest` 中添加新字段
- Update 操作：在 `ModifyFunctionRequest` 中添加新字段（如果 API 支持）
- Read 操作：从 `DescribeFunctionsResponse` 中读取新字段
- Delete 操作：保持不变

**Alternatives considered:**
- 使用单独的 API 调用：需要额外的网络开销，且可能造成状态不一致
- 使用自定义的数据结构映射：增加复杂性，且难以维护

### 3. 向后兼容性处理

**Decision:** 所有新增字段都设置为 Optional，不修改现有 Required 字段

**Rationale:**
- 现有用户配置不需要修改即可继续使用
- 新字段默认值为空或 nil，不影响现有行为
- 遵循 Terraform Provider 的最佳实践

**Implementation:**
- `environment_variables`：Optional，默认为空列表
- `rules`：Optional，默认为空列表
- `region_selection`：Optional，默认为 nil（表示全局部署）

### 4. 数据验证和错误处理

**Decision:** 在 Schema 中使用 `ValidateFunc` 和在 CRUD 函数中进行数据验证

**Rationale:**
- Schema 层验证可以在早期发现配置错误
- CRUD 函数层验证可以处理 API 返回的错误

**Validation rules:**
- `environment_variables`：
  - key：最大 64 字节，只允许字母、数字和 @ . - _ 字符
  - value：最大 5000 字节
  - type：只能是 "string" 或 "json"
- `rules`：
  - rule_id：必须唯一
  - priority：必须为正整数
  - actions：至少包含一个 action
- `region_selection`：
  - regions：必须是有效的 region code

### 5. 测试策略

**Decision:** 为每个新增功能添加单元测试和验收测试

**Rationale:**
- 单元测试验证 CRUD 操作的正确性
- 验收测试验证与实际 API 的集成

**Test coverage:**
- `resource_tc_teo_function_test.go`：添加新字段的测试用例
- 测试场景：
  - Create with new fields
  - Read with new fields
  - Update new fields
  - Delete with new fields
  - Validation errors

## Risks / Trade-offs

### Risk 1: TEO API 可能不支持所有字段

**Description:** TEO SDK 的 `CreateFunction` 和 `ModifyFunction` API 可能不完全支持环境变量、规则和区域选择参数。

**Mitigation:**
- 先验证 TEO API 的文档和 SDK 结构体定义
- 如果 API 不支持某些参数，考虑使用额外的 API 调用（如 `ModifyFunctionEnvironmentVariables`）
- 在 design 阶段确认 API 支持情况，避免实施过程中发现 API 限制

### Risk 2: 复杂嵌套结构可能导致状态不一致

**Description:** 新增的字段都是复杂嵌套结构，在 Update 操作时可能出现状态不一致的问题。

**Mitigation:**
- 使用 Terraform 的 `diff.Suppress` 功能处理不重要的差异
- 在 Read 操作中完整读取所有字段，确保状态一致
- 添加详细的日志记录，便于调试状态不一致问题

### Risk 3: 性能影响

**Description:** 新增字段可能会增加 API 调用的数据量和处理时间。

**Mitigation:**
- 使用分页查询（如果 API 支持）限制返回的数据量
- 在 Read 操作中只读取必要的字段
- 使用缓存机制（如果适用）减少重复的 API 调用

### Trade-off: 配置复杂度 vs 功能完整性

**Description:** 新增功能增加了配置的复杂度，用户需要学习新的字段和概念。

**Mitigation:**
- 提供详细的文档和示例
- 使用合理的默认值
- 在文档中提供常见配置模式的示例
- 考虑提供数据源（DataSource）来简化配置

## Migration Plan

### Deployment Steps

1. **Pre-deployment**
   - 验证 TEO API 是否支持所有新字段
   - 编写完整的单元测试和验收测试
   - 更新文档和示例代码

2. **Deployment**
   - 提交代码到主分支
   - 发布新版本的 Terraform Provider
   - 通知用户新功能和变更

3. **Post-deployment**
   - 监控用户反馈和错误日志
   - 收集性能指标
   - 根据反馈进行优化

### Rollback Strategy

如果发现严重问题，可以采取以下回滚策略：
1. 发布新版本，将新字段设置为 Deprecated
2. 添加配置迁移工具，帮助用户迁移到旧的配置方式
3. 在文档中明确说明回滚步骤

## Open Questions

1. **API Support:** TEO SDK 的 `CreateFunction` 和 `ModifyFunction` API 是否完全支持环境变量、规则和区域选择参数？如果不支持，需要使用哪些替代 API？

2. **State Migration:** 对于已经创建的函数，如果用户首次添加新字段，如何确保状态的一致性？

3. **Documentation:** 需要确定文档的详细程度和示例数量，以帮助用户理解和使用新功能。

4. **Validation Level:** 数据验证应该在 Schema 层、CRUD 函数层，还是 API 层进行？如何在保证正确性的同时避免重复验证？

5. **Error Messages:** 需要确定错误消息的详细程度和本地化策略，以提供良好的用户体验。
