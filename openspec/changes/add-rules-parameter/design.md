## Context

当前 tencentcloud_teo_l7_acc_rule 资源缺少对 Rules 参数的支持。根据 DescribeL7AccRules API 的定义，Rules 是一个复杂的嵌套结构，包含规则的详细信息：
- Status: 规则状态（enable/disable）
- RuleId: 规则唯一标识
- RuleName: 规则名称
- Description: 规则注释列表
- Branches: 子规则分支，包含 Condition 和 Actions
- RulePriority: 规则优先级（仅作为出参）

该变更需要在现有的资源中新增这个复杂的嵌套字段，同时确保与 CAPI 接口的数据结构完全一致，并保持向后兼容。

## Goals / Non-Goals

**Goals:**
- 在 tencentcloud_teo_l7_acc_rule 资源 Schema 中新增 Rules 参数（可选，列表类型）
- 正确映射 Rules 的所有嵌套字段（包括多层次的 Actions 和 SubRules）
- 更新 Read 函数以从 DescribeL7AccRules API 响应中读取 Rules 数据
- 更新 Create/Update 函数以处理 Rules 字段的转换和提交
- 保持向后兼容，不破坏现有用户配置

**Non-Goals:**
- 不修改现有资源的其他字段
- 不引入新的外部依赖
- 不改变现有 API 的调用方式

## Decisions

### 数据结构设计

Rules 参数采用嵌套的 Schema 结构，每个 Rule 包含以下字段：
- `status`: 规则状态（Computed, 可选）
- `rule_id`: 规则 ID（Computed, 可选）
- `rule_name`: 规则名称（Computed, 可选）
- `description`: 规则注释列表（Computed, 可选）
- `branches`: 子规则分支列表（Computed, 可选）

每个 Branch 包含：
- `condition`: 匹配条件（Computed, 可选）
- `actions`: 操作列表（Computed, 可选）
- `sub_rules`: 子规则列表（Computed, 可选）

Actions 是一个复杂的多态结构，根据 Name 字段的不同，会有不同的 Parameters：
- Cache, CacheKey, CachePrefresh 等不同类型的参数

### Schema 字段类型选择

所有 Rules 相关字段标记为 Computed，因为：
1. 这些字段由服务端返回和管理
2. RuleId 是系统生成的唯一标识
3. Status、RulePriority 等是服务端维护的状态信息

如果未来需要支持用户创建和修改规则，可以逐步将部分字段改为 Optional。

### API 集成

在 Read 函数中：
1. 调用 DescribeL7AccRules API 获取 Rules 列表
2. 将 API 响应映射到 Terraform Schema 结构
3. 使用 `d.Set()` 设置 Rules 数据

Create/Update 函数暂不处理 Rules 的创建和修改，因为这些字段当前仅作为 Computed 字段存在。

### 错误处理

- API 调用失败时，返回标准错误信息
- Rules 字段为 nil 或空时，跳过设置（避免 nil pointer 错误）
- 使用 `d.Set()` 时的错误需要正确处理和返回

## Risks / Trade-offs

**[Risk] Rules 结构复杂，嵌套层级深**
→ **Mitigation**: 使用递归或清晰的辅助函数来处理 Actions 的多态结构，确保代码可维护性

**[Risk] Actions 的 Parameters 类型众多（30+ 种），实现工作量大**
→ **Mitigation**: 优先实现最常用的 Actions 类型（如 Cache, CacheKey, AccessURLRedirect），其他类型可以逐步补充，或者在文档中说明当前支持的范围

**[Risk] API 返回的数据结构可能与 Terraform Schema 不完全匹配**
→ **Mitigation**: 在代码中添加详细的类型转换逻辑，确保数据类型正确转换，并在 Read 函数中添加日志以便调试

**[Risk] Computed 字段的值在 Create 后可能与预期不符**
→ **Mitigation**: 在文档中明确说明 Rules 字段的含义和预期行为，建议用户在 Create 后通过 Read 操作查看完整的 Rules 配置

**[Trade-off] 当前将所有字段设为 Computed 可能限制用户灵活性**
→ **Justification**: 基于当前需求，Rules 主要用于读取服务端配置。未来如果需要支持用户修改规则，可以逐步将相应字段改为 Optional，这样可以在保持兼容性的同时扩展功能
