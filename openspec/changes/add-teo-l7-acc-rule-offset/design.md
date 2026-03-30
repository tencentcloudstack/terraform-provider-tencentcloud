## Context

当前 tencentcloud_teo_l7_acc_rule 资源的 Read 操作通过 `DescribeTeoL7AccRuleById` 服务函数调用 DescribeL7AccRules API。该函数当前未使用分页参数，直接返回 API 的结果。当某个 zone 下的规则数量超过 API 单次查询的默认限制时，可能导致部分规则无法被获取，从而出现数据不一致的问题。

DescribeL7AccRules API 支持 Offset 和 Limit 参数用于分页查询。参考 teo 服务中其他类似函数的实现（如 DescribeZones、DescribeOriginGroup 等），这些函数都使用了分页循环来确保获取完整数据。

## Goals / Non-Goals

**Goals:**
- 修改 DescribeTeoL7AccRuleById 函数，添加 Offset 和 Limit 参数支持
- 实现分页循环逻辑，确保获取所有 L7 规则数据
- 保持向后兼容性，不修改资源 schema 和现有行为

**Non-Goals:**
- 不修改资源的 Schema 定义
- 不影响 Create、Update、Delete 操作
- 不新增配置参数或输出字段

## Decisions

**1. 分页策略选择**
- 决定使用 `for` 循环 + Offset/Limit 模式
- 参考 teo 服务中现有实现模式（如 service_tencentcloud_teo.go 中的其他函数）
- 初始 Offset 设为 0，Limit 设为 100（与 teo 服务中其他 List 操作保持一致）

**2. 数据收集方式**
- 在循环外初始化 `[]*teov20220901.Rule` 切片用于收集所有规则
- 每次循环将结果追加到切片中
- 当返回结果数量小于 Limit 或为空时，退出循环

**3. 代码修改范围**
- 仅修改 `DescribeTeoL7AccRuleById` 函数
- 保持函数签名不变（参数和返回值）
- resource 层的 `resourceTencentCloudTeoL7AccRuleRead` 函数无需修改，因为返回值类型兼容

**Rationale:**
- 参考现有的 teo 服务实现，保持代码风格一致
- 不改变函数签名，确保调用方无需修改
- 使用常见的分页模式，降低引入 bug 的风险

## Risks / Trade-offs

**Risk 1: API 返回数据量过大可能导致性能问题**
- Mitigation: Limit 设置为合理的值（100），避免单次请求过大

**Risk 2: 规则数量增长过快可能导致多次 API 调用**
- Mitigation: 分页逻辑是必要的，确保数据完整性比性能更重要

**Trade-off: 分页会增加 API 调用次数 vs 数据完整性**
- 选择数据完整性优先，因为 Terraform 资源的一致性要求准确的状态信息

**Trade-off: 初始 Offset 从 0 开始 vs 保留原有逻辑**
- 选择从 0 开始的完整分页逻辑，确保即使规则数量增加也能获取全部数据

## Migration Plan

1. 修改 `DescribeTeoL7AccRuleById` 函数，添加分页逻辑
2. 编译项目确保代码无语法错误
3. 运行相关测试用例验证功能正确性
4. 验证对现有配置无影响（向后兼容）

Rollback 策略：
- 如有问题，可以通过 git revert 快速回退代码
- 不涉及 schema 变更，不会影响 state 文件
