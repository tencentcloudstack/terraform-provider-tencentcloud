## Context

tencentcloud_teo_function 资源目前已支持基本的 CRUD 操作，但在读取操作中不支持通过函数 ID 列表进行过滤。DescribeFunctions API 已经支持 FunctionIds 参数，该参数允许用户指定一个函数 ID 列表来过滤查询结果。为了提高查询效率和灵活性，需要在 Terraform Provider 中对接该功能。

当前状态：
- tencentcloud_teo_function 资源存在，Schema 已定义
- Read 函数调用 DescribeFunctions API，但未使用 FunctionIds 参数
- 资源使用复合 ID 格式：`zoneId#functionId`

约束条件：
- 必须保持向后兼容，不能破坏现有 TF 配置和 state
- 新增字段必须为 Optional，不能为 Required
- 需要更新相关的单元测试和验收测试
- 必须保持与 DescribeFunctions API 的接口一致性

## Goals / Non-Goals

**Goals:**
- 在 tencentcloud_teo_function 资源的 Schema 中新增 FunctionIds 字段（list, 可选）
- 更新 Read 函数，使用 FunctionIds 字段作为 DescribeFunctions API 的请求参数
- 确保 FunctionIds 字段在 Read 操作中正常工作，能够正确过滤函数列表
- 更新单元测试和验收测试，覆盖新增字段的功能

**Non-Goals:**
- 不修改 Create、Update、Delete 函数的现有逻辑（因为 FunctionIds 仅用于 Read 操作的过滤）
- 不改变现有的复合 ID 格式
- 不影响其他资源或数据源

## Decisions

1. **字段类型选择**：使用 `schema.TypeList` + `schema.TypeString` 来定义 FunctionIds 字段
   - 理由：FunctionIds 是一个字符串列表，与 DescribeFunctions API 的参数类型一致
   - 替代方案：使用 `schema.TypeSet` - 但由于函数 ID 列表不需要去重，List 更合适

2. **字段属性设置**：FunctionIds 设置为 Optional + Computed
   - 理由：
     - Optional：用户可以选择是否使用该字段进行过滤
     - Computed：虽然该字段用于 Read 操作的过滤，但为了避免混淆，不设置为 Computed
   - 替代方案：只设置为 Optional - 更符合过滤参数的语义

3. **Read 函数实现**：在 Read 函数中，检查 d.Get("function_ids") 是否存在，如果存在则将其作为 DescribeFunctions API 的 FunctionIds 参数
   - 理由：保持现有 Read 函数的逻辑结构，仅在需要时添加过滤条件
   - 替代方案：在 Resource 函数中统一处理过滤条件 - 会增加代码复杂度

4. **测试策略**：新增单元测试和验收测试，覆盖以下场景：
   - 不使用 FunctionIds 字段时的正常读取
   - 使用 FunctionIds 字段过滤单个函数 ID
   - 使用 FunctionIds 字段过滤多个函数 ID
   - 使用无效的函数 ID 时的错误处理

## Risks / Trade-offs

**风险 1：** API 兼容性风险
- 如果 DescribeFunctions API 的 FunctionIds 参数行为发生变化，可能导致过滤结果不符合预期
- 缓解措施：通过验收测试验证过滤行为的正确性，定期测试

**风险 2：** 性能影响
- 如果用户传入大量的函数 ID，可能会影响 API 调用性能
- 缓解措施：DescribeFunctions API 应该支持合理的参数列表长度，由云服务端控制

**权衡 1：** 字段设置为 Optional vs Required
- 选择 Optional：保持向后兼容，不影响现有配置
- 权衡：如果用户希望强制使用该字段，需要在代码中添加额外的验证逻辑

**权衡 2：** 在 Read 函数中处理 vs 单独创建 DataSource
- 选择在 Read 函数中处理：简化实现，用户可以在同一个资源中使用过滤功能
- 权衡：如果过滤需求非常复杂，可能需要单独创建 DataSource

## Migration Plan

由于 FunctionIds 是新增的 Optional 字段，不涉及迁移：
- 现有的 Terraform 配置无需修改，可以继续正常工作
- 用户可以选择性地在配置中添加 FunctionIds 字段来启用过滤功能
- 不需要 state 迁移，因为该字段不影响资源的唯一标识

**部署步骤：**
1. 更新 resource_tc_teo_function.go 文件，新增 FunctionIds 字段定义
2. 更新 Read 函数，添加 FunctionIds 参数的处理逻辑
3. 更新 resource_tc_teo_function_test.go 文件，新增测试用例
4. 运行单元测试验证功能
5. 更新 resource_tc_teo_function.md 文档
6. 提交代码并通过 Code Review

**回滚策略：**
- 如果出现问题，可以移除 FunctionIds 字段及其相关逻辑
- 由于字段为 Optional，不会影响现有的 Terraform 配置和 state

## Open Questions

无 - 本次变更的技术决策已经明确，不存在未解决的问题。
