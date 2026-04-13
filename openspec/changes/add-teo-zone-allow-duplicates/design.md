## Context

`TencentCloud Provider` 已经支持 `tencentcloud_teo_zone` 资源，该资源对应 TEO（Tencent EdgeOne）站点。目前该资源的 Schema 支持创建站点所需的核心字段，但 CreateZone API 新增了 `allow_duplicates` 参数，用于控制是否允许在站点中创建重复的规则配置。为了保持与云 API 功能的对齐，需要在 Provider 资源中添加对该参数的支持。

当前 `tencentcloud_teo_zone` 资源位于 `tencentcloud/services/teo/resource_tc_teo_zone.go`，包含基本的 CRUD 操作，通过 Terraform Plugin SDK v2 实现资源管理。

## Goals / Non-Goals

**Goals:**

- 在 `tencentcloud_teo_zone` 资源的 Schema 中新增 `allow_duplicates` 字段（Optional 类型，布尔值）
- 在 Create 函数中调用 CreateZone API 时传入 `allow_duplicates` 参数
- 在 Read 函数中从 DescribeZone API 响应中读取并映射 `allow_duplicates` 字段
- 在 Update 函数中支持通过 ModifyZone API 更新 `allow_duplicates` 参数（如果 API 支持）
- 更新单元测试，覆盖新字段的各种场景（true/false/未设置）
- 更新验收测试，验证 `allow_duplicates` 参数的正确性

**Non-Goals:**

- 不修改 `tencentcloud_teo_zone` 资源的其他现有字段
- 不改变资源的基本 CRUD 流程和架构
- 不涉及其他 TEO 资源的修改

## Decisions

**1. 字段类型选择：布尔类型**

- 决策：`allow_duplicates` 字段定义为 `schema.TypeBool`
- 理由：根据 CreateZone API 文档，`allow_duplicates` 参数为布尔值类型，用于控制是否允许重复配置
- 备选方案：无（API 定义明确）

**2. 字段属性：Optional**

- 决策：`allow_duplicates` 字段设置为 Optional 而非 Required
- 理由：
  - 保持向后兼容性，现有配置不需要修改
  - API 允许不传递该参数，使用默认值
  - 符合 Terraform Provider 的最佳实践
- 备选方案：设置为 Required（会破坏现有配置）

**3. Update 策略：仅当字段发生变化时调用 API**

- 决策：在 Update 函数中，仅当 `allow_duplicates` 值发生变化时才调用 ModifyZone API
- 理由：
  - 减少不必要的 API 调用
  - 避免潜在的 API 限流问题
  - 提高性能
- 备选方案：每次 Update 都调用 API（增加不必要的开销）

**4. Read 策略：从 API 响应中读取实际值**

- 决策：在 Read 函数中，从 DescribeZone API 响应中读取 `allow_duplicates` 字段值并设置到 state
- 理由：
  - 确保状态与云服务实际配置一致
  - 支持状态刷新和重建
  - 符合 Terraform Provider 的最佳实践
- 备选方案：仅读取用户配置值（无法反映实际状态）

## Risks / Trade-offs

**Risk 1: API 行为不确定**

- 风险：如果 ModifyZone API 不支持修改 `allow_duplicates` 参数，Update 操作会失败
- 缓解：在实施前通过测试验证 API 是否支持该参数的修改；如果不支持，则在文档中明确说明该字段仅在创建时可设置

**Risk 2: 默认值不一致**

- 风险：用户未设置 `allow_duplicates` 时，API 的默认值可能与预期不符
- 缓解：在 Read 函数中明确读取 API 返回的实际值，不假设任何默认值；在文档中说明 API 的默认行为

**Risk 3: 向后兼容性**

- 风险：添加新字段可能影响现有用户的 state 读取
- 缓解：新字段设置为 Optional，state 升级时会自动使用默认值，不会导致配置错误

**Trade-off: 测试覆盖范围**

- 权衡：是否需要为所有 `allow_duplicates` 值组合（true/false/未设置）编写完整的验收测试
- 决策：编写关键路径的验收测试（true 和 false），单元测试覆盖所有场景，以平衡测试成本和覆盖率
