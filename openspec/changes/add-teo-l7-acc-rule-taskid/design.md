## Context

tencentcloud_teo_l7_acc_rule 资源目前通过 ImportZoneConfig API 进行更新操作。业务需求要求在更新操作时能够传入 TaskId 参数，用于跟踪异步任务和管理批量操作。

## Goals / Non-Goals

**Goals:**
- 在 tencentcloud_teo_l7_acc_rule 资源的 schema 中添加 TaskId 可选参数
- 在 update 函数中传递 TaskId 参数到 ImportZoneConfig API
- 保持向后兼容性，不影响现有资源配置

**Non-Goals:**
- 修改 create 或 read 操作的行为
- 改变 TaskId 的语义或使用方式
- 修改其他资源的 TaskId 参数

## Decisions

1. **Schema Design**
   - TaskId 参数定义为 Optional 和 Computed 类型
   - 使用 TypeString 类型存储 TaskId 值
   - 不设置 ForceNew，允许在不重新创建资源的情况下更新 TaskId

2. **API 集成**
   - 在 resourceTencentCloudTeoL7AccRuleUpdate 函数中获取 TaskId 值
   - 构建请求参数时，如果 TaskId 不为空，则将其添加到 ImportZoneConfig API 请求中

3. **向后兼容性**
   - TaskId 是可选参数，不影响现有配置
   - 未指定 TaskId 时，API 调用行为保持不变

## Risks / Trade-offs

- [API 兼容性] → ImportZoneConfig API 需要支持 TaskId 参数，已确认 API 支持该参数
- [状态管理] → TaskId 仅用于 API 调用，不需要在 state 中持久化，避免增加状态复杂度
