## Context

当前 `tencentcloud_teo_l7_acc_rule` 资源已经使用 `DescribeL7AccRules` API 来读取规则数据。虽然 API 支持 `Filters` 参数用于筛选查询结果，但当前实现只在内部使用了固定过滤器（`rule-id`），并未暴露给用户配置。

Teo 服务的实现遵循 Terraform Provider SDK v2 的标准模式，使用 `DescribeTeoL7AccRuleById` 服务层函数封装 API 调用。当前的过滤器实现（第 1629-1634 行）硬编码为只支持 `rule-id` 过滤。

## Goals / Non-Goals

**Goals:**
- 在 `tencentcloud_teo_l7_acc_rule` 资源的 schema 中添加 `filters` 参数
- 支持用户通过 Terraform 配置自定义过滤器条件
- 保持向后兼容性，不影响现有使用方式
- 遵循 Terraform Provider SDK v2 的标准过滤参数模式

**Non-Goals:**
- 修改 `DescribeL7AccRules` API 的行为
- 支持非 `Filters` 参数的查询方式（如分页、排序等）
- 修改现有资源的状态管理逻辑

## Decisions

### 1. Schema 设计：复用 Terraform 标准过滤器模式

**选择：** 使用 Terraform Plugin SDK v2 的标准过滤器 schema 模式，定义 `filters` 为 `TypeSet` 类型，包含 `name` 和 `values` 字段。

**理由：**
- 这与 Terraform 生态系统中其他 Provider 的实现保持一致（如 AWS、Azure Provider）
- 用户熟悉这种模式，学习成本低
- SDK 提供了内置的验证和类型转换支持

**替代方案考虑：**
- *自定义过滤器结构*：虽然更灵活，但会增加学习成本，不符合 Terraform 最佳实践

### 2. 过滤器传递方式：通过服务层函数参数

**选择：** 修改 `DescribeTeoL7AccRuleById` 函数签名，添加 `filters` 参数，将用户配置的过滤器映射到 API 请求。

**理由：**
- 保持现有的服务层封装模式
- 在资源层和 API 层之间提供清晰的转换逻辑
- 便于单元测试和错误处理

**替代方案考虑：**
- *直接在资源层构建 API 请求*：虽然简单，但会破坏封装性，不利于代码维护

### 3. 兼容性处理：可选参数 + 默认行为

**选择：** 将 `filters` 设置为 `Optional` 参数，当用户未提供时，保持现有查询行为（只按 zone_id 查询）。

**理由：**
- 确保向后兼容性，现有配置无需修改
- 用户可以根据需要选择是否使用过滤器
- 避免不必要的 API 调用开销

**替代方案考虑：**
- *强制要求过滤器*：会破坏现有配置，不符合最小变更原则

### 4. 过滤器实现位置：资源层逻辑

**选择：** 在 `resourceTencentCloudTeoL7AccRuleRead` 函数中，从 schema 读取 `filters` 参数，转换为 API 所需的 `Filter` 结构。

**理由：**
- 资源层是 Terraform 配置和 Provider 逻辑之间的桥梁
- 便于进行参数验证和错误处理
- 保持服务层函数的职责单一（只负责 API 调用）

**替代方案考虑：**
- *在服务层处理过滤逻辑*：虽然可行，但会增加服务层复杂度

## Risks / Trade-offs

### 风险 1：不兼容的过滤器名称或值
**风险：** 用户可能使用 API 不支持的过滤器名称或值，导致 API 调用失败。

**缓解措施：**
- 在 Terraform 文档中明确列出支持的过滤器类型和值
- 在资源层添加参数验证逻辑（可选）
- 将 API 错误信息清晰地返回给用户

### 风险 2：过滤器组合导致的查询结果为空
**风险：** 多个过滤器组合可能导致查询结果为空，用户无法理解原因。

**缓解措施：**
- 在文档中说明过滤器的 AND 组合逻辑
- 在 README 或示例中提供使用示例

### 权衡：性能 vs 灵活性
**权衡：** 支持多个过滤器会增加 API 调用的复杂度，但提供了更强大的查询能力。

**决策：** 选择灵活性，允许用户配置多个过滤器，因为这是用户的核心需求。

## Migration Plan

### 部署步骤
1. 修改 `resource_tc_teo_l7_acc_rule.go`，在 schema 中添加 `filters` 参数
2. 修改 `service_tencentcloud_teo.go` 中的 `DescribeTeoL7AccRuleById` 函数，支持接收自定义过滤器
3. 在 `resourceTencentCloudTeoL7AccRuleRead` 函数中，从 schema 读取 `filters` 并传递给服务层
4. 添加单元测试和集成测试
5. 更新文档和示例

### 回滚策略
- 使用版本控制，如果出现问题可以快速回滚到之前的版本
- 保持向后兼容性，即使新功能有问题，现有配置仍可正常工作

## Open Questions

1. **过滤器值验证范围：** 是否需要在代码层面验证过滤器名称和值的有效性？
   - *建议：* 初期不进行验证，依赖 API 返回错误信息；后续根据用户反馈决定是否添加验证

2. **过滤器组合逻辑：** 多个过滤器之间是 AND 还是 OR 关系？
   - *建议：* 使用 AND 关系（即结果必须满足所有过滤器条件），这是 Terraform 过滤器的常见做法
