## Context

tencentcloud_teo_rule_engine 资源已经实现了基本的规则引擎功能，可以读取和配置规则引擎的基本信息。然而，DescribeRules API 返回的响应中包含了 RuleItems 参数，该参数包含了规则引擎的详细规则项配置信息（如规则类型、规则条件、规则操作等）。目前该参数没有被映射到 Terraform schema 中，导致用户无法获取完整的规则配置信息。

当前实现：
- tencentcloud_teo_rule_engine 资源通过 DescribeRules API 读取规则引擎基本信息
- 响应数据中的 RuleItems 字段被忽略，未映射到 schema

约束条件：
- 必须保持向后兼容，不能破坏现有 TF 配置和 state
- 只能新增 Optional 参数，不能修改已有资源的 schema 结构
- 需要更新相关的文档和示例

## Goals / Non-Goals

**Goals:**
- 为 tencentcloud_teo_rule_engine 资源添加 RuleItems 参数的 schema 定义
- 在资源读取操作中解析 DescribeRules API 响应中的 RuleItems 数据
- 更新相关的文档和示例，展示如何使用 RuleItems 参数

**Non-Goals:**
- 不修改现有的 schema 字段或删除已有功能
- 不涉及 RuleItems 的写操作（只读）
- 不改变现有的 API 调用逻辑

## Decisions

**1. Schema 设计**
- 将 RuleItems 定义为 List 或 Set 类型，包含多个规则项
- 每个规则项包含必要的字段（根据 API 响应结构定义）
- 设置为 Optional 字段，保持向后兼容性

**2. 数据映射**
- 在资源读取函数中添加 RuleItems 的解析逻辑
- 将 DescribeRules API 响应中的 RuleItems 字段映射到 schema
- 使用 Terraform Plugin SDK v2 提供的类型转换函数进行数据转换

**3. 字段结构**
- 根据 DescribeRules API 返回的 RuleItems 结构定义 schema
- 如果 RuleItems 包含嵌套结构，使用 List 或 Map 类型表示
- 确保字段命名符合 Terraform 命名规范（snake_case）

## Risks / Trade-offs

**风险 1：** API 响应结构变化
- **缓解措施：** 使用可选字段，如果 API 结构变化导致解析失败，不影响现有功能

**风险 2：** RuleItems 数据量大导致性能问题
- **缓解措施：** RuleItems 为 Optional 参数，用户可以选择性读取，不强制解析

**权衡：**
- 方案 A：完全映射所有 RuleItems 字段 → 提供完整信息，但增加 schema 复杂度
- 方案 B：只映射关键字段 → 简化 schema，但可能丢失有用信息
- **决定：** 采用方案 A，完全映射所有 RuleItems 字段，因为用户需要完整的规则配置信息
