## Context

当前 `tencentcloud_teo_function` 资源在 Read 操作中内部使用了 `DescribeFunctions` API，并通过 `FunctionIds` 参数传入单个函数 ID 进行查询。但是，这个参数并未作为资源属性暴露给 Terraform 用户。在某些特定场景下（如批量同步、数据迁移等），用户可能需要一次性查询多个函数的信息以提高效率。

现有的实现中：
- `resourceTencentCloudTeoFunctionRead` 函数通过 `DescribeTeoFunctionById` 服务层方法读取单个函数
- `DescribeTeoFunctionById` 内部调用 `DescribeFunctions` API 并传入 `FunctionIds` 参数
- Resource ID 格式为 `{zone_id}#{function_id}`，使用复合 ID

根据 Terraform Provider 的最佳实践，单个 resource 实例应该代表一个具体的云资源实例。因此，新增的 `function_ids` 参数主要用于优化读取操作的灵活性，而不是改变 resource 的基本语义。

## Goals / Non-Goals

**Goals:**

1. 在 `tencentcloud_teo_function` 资源 Schema 中新增 `function_ids` 参数字段（类型为 `TypeList`，元素类型 `TypeString`），属性为 `Optional`
2. 更新 Read 函数逻辑，支持使用 `function_ids` 参数进行多函数查询
3. 保持向后兼容，不破坏现有的单个函数读取功能
4. 确保 `function_ids` 参数与 CAPI 的 `DescribeFunctions` API 的 `FunctionIds` 参数定义一致
5. 更新单元测试和验收测试，覆盖新的查询场景

**Non-Goals:**

1. 不修改 Create、Update、Delete 操作的逻辑（这些操作保持单个函数的基本行为）
2. 不改变 Resource ID 的格式和语义
3. 不引入新的 API 依赖或外部库
4. 不改变其他 TEO 相关资源的实现

## Decisions

### 1. Schema 字段设计

**决策**: 在 `tencentcloud_teo_function` 资源中新增 `function_ids` 字段，类型为 `schema.TypeList`，元素类型为 `schema.TypeString`，属性为 `Optional` 和 `Computed`

**理由**:
- 与 CAPI 的 `FunctionIds` 参数类型一致（`[]*string`）
- `Optional` 属性确保向后兼容，不影响现有用户
- `Computed` 属性表示该字段可以由系统返回数据，但用户也可以主动设置

**考虑的替代方案**:
- 使用 `TypeSet` 而非 `TypeList`: `Set` 会自动去重，但在函数 ID 场景下，用户可能需要保持顺序或允许重复（虽然不太常见），因此选择 `List`
- 不添加 `Computed` 属性: 但考虑到某些场景下 API 可能返回多个函数信息，设置为 `Computed` 更灵活

### 2. Read 函数逻辑调整

**决策**: 在 `resourceTencentCloudTeoFunctionRead` 函数中，优先检查 `function_ids` 参数。如果设置了 `function_ids`，使用该参数调用 `DescribeFunctions` API；否则，继续使用现有的单个 `function_id` 查询逻辑

**理由**:
- 保持向后兼容，不破坏现有用户的读取行为
- 新的查询模式通过可选参数触发，用户可以自主选择
- 减少对现有代码的影响范围

**考虑的替代方案**:
- 强制所有查询都使用 `function_ids`（单元素数组）: 这会破坏向后兼容性，不符合要求
- 创建新的服务层方法专门处理多函数查询: 这会增加代码复杂度，当前在现有方法中扩展即可

### 3. 服务层方法扩展

**决策**: 扩展 `DescribeTeoFunctionById` 服务层方法，或创建新的方法 `DescribeTeoFunctionsByIds` 以支持多函数查询

**理由**:
- 保持服务层的职责单一性
- 如果新方法与现有方法差异较大，创建新方法更清晰
- 如果只是参数层面的差异，可以考虑扩展现有方法

**实现细节**:
- 新方法接收 `[]string` 类型的 `functionIds` 参数
- 内部调用 `DescribeFunctions` API 时传入 `FunctionIds` 参数
- 返回 `[]*Function` 类型的结果

### 4. 向后兼容性保证

**决策**: 确保 `function_ids` 参数为 `Optional`，且默认行为与现有实现一致

**理由**:
- 必须满足 Terraform Provider 的硬性约束：不能破坏现有 TF 配置和 state
- 用户升级后无需修改现有配置即可继续使用

**实现细节**:
- 在 Schema 中将 `function_ids` 设置为 `Optional`
- 在 Read 函数中，首先检查 `function_ids` 是否设置，未设置则使用现有逻辑
- 不修改 Resource ID 的生成和解析逻辑

### 5. 测试策略

**决策**: 扩展现有单元测试和验收测试，添加针对 `function_ids` 参数的测试用例

**理由**:
- 确保新功能的正确性
- 验证向后兼容性
- 覆盖边界场景（空列表、单元素列表、多元素列表等）

**测试场景**:
- 单个函数查询（未设置 `function_ids`）: 验证现有功能不受影响
- 单个函数查询（通过 `function_ids`）: 验证新参数的基本功能
- 多个函数查询（通过 `function_ids`）: 验证批量查询能力
- 边界场景（空列表、错误 ID 等）: 验证错误处理

## Risks / Trade-offs

### 风险 1: Resource 语义混淆

**描述**: 在单个 resource 实例中添加 `function_ids` 参数可能导致用户混淆，不清楚何时应该使用该参数

**缓解措施**:
- 在文档中明确说明 `function_ids` 参数的用途和使用场景
- 在 Schema 的 Description 中提供清晰的说明
- 在示例中展示正确的使用方式

### 风险 2: 性能影响

**描述**: 批量查询可能返回大量数据，影响 Terraform 的性能

**缓解措施**:
- 限制 `function_ids` 列表的最大长度（与 CAPI 的限制一致）
- 在文档中建议合理的使用范围
- 在服务层添加必要的参数校验

### 风险 3: 测试覆盖不足

**描述**: 新功能的测试可能不够全面，导致边界场景出现问题

**缓解措施**:
- 添加全面的单元测试和验收测试
- 使用真实 API 进行验收测试（需要 TENCENTCLOUD_SECRET_ID/KEY 环境变量）
- 考虑添加性能测试

### 权衡 1: 灵活性 vs 复杂性

**描述**: 添加 `function_ids` 参数增加了使用灵活性，但也增加了代码复杂度

**权衡**: 选择灵活性，因为用户在特定场景下确实有批量查询的需求，而增加的复杂度在可控范围内

### 权衡 2: 新方法 vs 扩展现有方法

**描述**: 创建新的服务层方法增加了代码量，但保持了代码清晰度；扩展现有方法减少了代码量，但可能违反单一职责原则

**权衡**: 根据实际代码结构选择，如果现有方法简单且易于扩展，则扩展现有方法；否则创建新方法

## Migration Plan

### 部署步骤

1. **代码修改阶段**:
   - 修改 `resource_tc_teo_function.go` 中的 Schema 定义
   - 更新 `resourceTencentCloudTeoFunctionRead` 函数逻辑
   - 扩展或创建服务层方法以支持多函数查询
   - 更新文档 `resource_tc_teo_function.md`

2. **测试阶段**:
   - 运行单元测试（`go test ./tencentcloud/services/teo/...`）
   - 运行验收测试（`TF_ACC=1 go test ./tencentcloud/services/teo/...`）
   - 验证所有测试通过

3. **文档更新阶段**:
   - 更新资源文档，说明 `function_ids` 参数的用途和使用方法
   - 提供使用示例

4. **代码审查阶段**:
   - 提交 Pull Request 进行代码审查
   - 根据反馈进行必要的修改

5. **合并发布阶段**:
   - 合并代码到主分支
   - 发布新版本的 Terraform Provider

### 回滚策略

如果出现问题，可以通过以下方式回滚：
1. 移除 `function_ids` 参数（由于该参数是 Optional 的，移除后不会影响现有配置）
2. 恢复服务层方法的原始实现
3. 发布修复版本

## Open Questions

1. **`function_ids` 参数的使用限制**: 是否需要在 Schema 中添加 `MaxItems` 约束，以与 CAPI 的限制保持一致？
   - 建议: 查询 CAPI 文档，确定 `FunctionIds` 参数的最大长度限制，然后在 Schema 中添加相应的约束

2. **多函数查询的返回值处理**: 如果 `function_ids` 包含多个 ID，Read 函数应该如何处理返回的多个函数信息？
   - 建议: 由于 Terraform resource 实例应该对应单个资源，这种情况可能更适合使用 data source。如果用户确实需要这种能力，可能需要考虑创建新的 data source 而非修改 resource

3. **服务层方法的命名**: 应该扩展 `DescribeTeoFunctionById` 方法，还是创建新的 `DescribeTeoFunctionsByIds` 方法？
   - 建议: 如果现有方法的逻辑简单且易于扩展，则扩展现有方法；否则创建新方法以保持代码清晰度
