## Context

当前 `tencentcloud_teo_origin_group` 资源的 schema 中已经定义了 `host_header` 参数，Update 和 Read 函数也都正确处理了该参数：
- Update 函数在第 449-451 行正确地将 `host_header` 传递给 ModifyOriginGroup API
- Read 函数在第 366-368 行正确地从 API 响应中读取 `host_header` 值

然而，在 Create 函数中（resourceTencentCloudTeoOriginGroupCreate），虽然处理了 name、type、records 等参数，但遗漏了 `host_header` 参数的处理。这导致用户无法在创建 origin_group 时直接指定 host_header。

CreateOriginGroup API 支持的参数包括：ZoneId、Records、Name、Type 和 HostHeader，但当前的实现只传递了前四个参数。

## Goals / Non-Goals

**Goals:**
- 在 Create 函数中添加对 `host_header` 参数的处理逻辑
- 确保用户创建资源时可以一次性配置 host_header，无需后续 Update
- 保持与现有 Update 和 Read 函数的实现一致性
- 不破坏现有配置和状态，确保向后兼容

**Non-Goals:**
- 不修改 schema 定义（schema 已经正确）
- 不改变 host_header 的语义或行为
- 不涉及其他参数的修改
- 不影响资源的其他功能

## Decisions

**在 Create 函数中添加 host_header 参数传递逻辑**

- **决策位置**：在 resourceTencentCloudTeoOriginGroupCreate 函数中，处理完 records 参数后（约第 220 行），添加 host_header 参数的处理代码
- **实现方式**：遵循现有参数处理的模式，使用 `d.GetOk("host_header")` 检查参数是否存在，如果存在则将其值传递给 request.HostHeader
- **代码位置**：添加在 `request.Records` 赋值之后、API 调用之前

**理由：**
1. 这个位置的代码模式与处理其他参数（如 name、type）完全一致
2. 位于 API 调用前确保所有参数都已正确设置
3. 与 Update 函数中的实现位置和方式保持一致

**为什么不其他方案：**
- 不考虑修改 schema（已定义且正确）
- 不考虑在其他位置添加代码（不符合现有模式）
- 不考虑重构整个 Create 函数（风险高，不必要）

## Risks / Trade-offs

**风险1：参数传递错误** → **缓解措施：** 严格按照现有参数处理模式编写代码，确保类型转换和指针赋值正确

**风险2：API 不支持该参数** → **缓解措施：** 已验证 CreateOriginGroup API 的结构定义明确包含 HostHeader 字段

**风险3：测试覆盖不足** → **缓解措施：** 添加专门的测试用例验证 host_header 参数在创建时的正确传递

**权衡：**
- 这是一个小范围、低风险的修复，只涉及添加几行代码
- 修改范围小，不会影响其他功能
- 不涉及架构变更，无需考虑复杂的迁移策略
