## Context

当前 TEO 服务的 Terraform Provider 缺少即时转码模板资源的管理能力。用户需要通过手动调用 API 或使用其他方式管理 TEO 的即时转码模板，这与 Terraform 的 Infrastructure as Code 理念不符。TEO 服务提供了完整的即时转码模板 API 接口（创建、查询、删除），为在 Terraform Provider 中添加该资源提供了技术基础。

## Goals / Non-Goals

**Goals:**

- 在 Terraform Provider 中新增 `tencentcloud_teo_just_in_time_transcode_template` 资源
- 实现完整的 CRUD 操作：创建、读取、更新、删除
- 支持视频流和音频流的转码配置
- 支持复合 ID 格式（zone_id#template_id）
- 实现异步操作的超时和重试机制
- 提供完整的单元测试和集成测试覆盖

**Non-Goals:**

- 不实现批量操作（如批量创建或删除多个模板）
- 不实现模板的导入功能（Terraform 的 terraform import 功能通过 Read 接口自然支持）
- 不实现复杂的模板版本管理
- 不实现模板的复制或克隆功能

## Decisions

### 资源 ID 设计

**决策**: 使用复合 ID 格式 `zone_id#template_id`

**理由**:
- TEO 服务中模板是在站点（Zone）维度下管理的
- 单独的 template_id 在不同站点间可能重复
- 符合 Provider 中复合 ID 的现有模式（如 `instanceId#userId`）
- 便于在 Read 操作中快速定位资源

**备选方案考虑**:
- 方案A: 仅使用 template_id
  - 优点: 简单
  - 缺点: 可能造成不同站点间的模板 ID 冲突
- 方案B: 使用 JSON 编码的复合键
  - 优点: 可扩展性强
  - 缺点: 不符合现有 Provider 模式，可读性差

### 更新策略

**决策**: 使用删除+重建的更新策略（ForceNew 参数）

**理由**:
- TEO API 没有提供 UpdateJustInTimeTranscodeTemplate 接口
- 模板配置参数较多，部分参数不支持直接修改
- 删除+重建是 Terraform 中处理不可变资源的标准模式
- 确保状态一致性，避免部分更新的复杂性

**备选方案考虑**:
- 方案A: 尝试部分更新（如果 API 支持）
  - 优点: 更新更高效
  - 缺点: 当前 API 不支持，需要额外适配层
- 方案B: 仅标记为不可变
  - 优点: 清晰的用户预期
  - 缺点: 用户体验差，任何修改都需要手动删除

### 资源 Schema 设计

**决策**:
- 基础参数（zone_id, template_name, comment）作为 Required
- 开关参数（video_stream_switch, audio_stream_switch）作为 Optional，默认值为 "on"
- 模板配置参数（video_template, audio_template）作为 Optional，但需要根据开关条件验证

**理由**:
- 符合 TEO API 的要求
- 提供合理的默认值，简化用户配置
- 支持仅配置视频或仅配置音频的场景

**Schema 结构设计**:
```
zone_id (String, Required, ForceNew)
template_name (String, Required, ForceNew)
comment (String, Optional, ForceNew)
video_stream_switch (String, Optional, ForceNew, Default: "on")
audio_stream_switch (String, Optional, ForceNew, Default: "on")
video_template {
  codec (String, Optional)
  fps (Float, Optional)
  bitrate (Int, Optional)
  resolution_adaptive (String, Optional)
  width (Int, Optional)
  height (Int, Optional)
  fill_type (String, Optional)
}
audio_template {
  codec (String, Optional)
  audio_channel (Int, Optional)
}
template_id (String, Computed)
create_time (String, Computed)
update_time (String, Computed)
```

### 错误处理和重试机制

**决策**:
- 使用 `tccommon.WriteRetryTimeout` 处理写操作的暂时性错误
- 使用 `tccommon.Retry()` 处理读操作的暂时性错误
- 对于认证错误、参数错误等不重试
- 在 CRUD 操作中使用 `defer tccommon.LogElapsed()` 记录操作耗时
- 使用 `defer tccommon.InconsistentCheck()` 进行状态一致性检查

**理由**:
- 符合现有 Provider 的错误处理模式
- 避免不必要的重试导致问题复杂化
- 提供充分的日志信息便于问题排查

### 异步操作处理

**决策**: Create 和 Delete 操作后调用 Read 接口轮询直到操作生效

**理由**:
- TEO API 是异步的，操作不会立即生效
- 需要通过 Read 接口确认操作状态
- 提供合理的超时机制（使用 Terraform 的 Timeouts 功能）

**实现细节**:
- 在 Create 操作后，调用 DescribeJustInTimeTranscodeTemplates 查询模板是否存在
- 在 Delete 操作后，调用 DescribeJustInTimeTranscodeTemplates 确认模板已被删除
- 使用 `resource.Retry()` 实现轮询逻辑
- 轮询间隔: 5秒
- 最大轮询次数: 根据用户配置的 Timeouts 决定

### 单元测试策略

**决策**: 使用 mock（gomonkey）方法进行单元测试

**理由**:
- 这是新增资源，不涉及修改现有资源
- 避免依赖真实的云环境和网络
- 测试速度更快，更适合 CI/CD 流程
- 可以精确控制 API 返回值，测试各种边界情况

**测试覆盖范围**:
- 成功场景：创建、读取、删除
- 错误场景：API 返回错误、参数验证失败
- 边界场景：参数为空、参数超出范围
- 重试场景：暂时性错误重试成功

## Risks / Trade-offs

### 风险 1: API 接口变更

**风险**: TEO 服务可能会修改 API 接口或参数，导致 Provider 代码需要相应调整

**缓解措施**:
- 使用 vendor 模式管理依赖，锁定 API 版本
- 定期检查 API 更新公告
- 在升级 API 版本时进行充分测试

### 风险 2: 异步操作超时

**风险**: 在某些情况下，模板创建或删除可能需要较长时间，导致 Terraform 操作超时

**缓解措施**:
- 提供可配置的 Timeouts 参数，让用户可以根据实际情况调整
- 在文档中说明可能需要较长等待时间
- 使用合理的默认超时时间（如 30 分钟）

### 风险 3: 参数验证复杂性

**风险**: video_template 和 audio_template 参数的验证逻辑较为复杂，可能存在边界情况未考虑

**缓解措施**:
- 充分的单元测试覆盖各种参数组合
- 在 Schema 中使用 ValidateFunc 进行参数验证
- 在文档中清晰说明参数的取值范围和约束条件

### 权衡 1: 更新策略选择

**权衡**: 删除+重建的更新策略会导致操作不可逆，可能造成临时服务中断

**缓解措施**:
- 在文档中明确说明更新操作的不可逆性
- 建议用户在生产环境操作前进行充分测试
- 未来如果 API 提供了 Update 接口，可以考虑优化为原地更新

### 权衡 2: 资源 ID 复杂度

**权衡**: 复合 ID 格式可能让用户在使用时感到不便

**缓解措施**:
- 在文档中详细说明 ID 的格式和组成部分
- 提供清晰的错误提示，帮助用户理解 ID 格式错误
- 考虑在 Read 操作中提供更友好的错误信息

## Migration Plan

由于这是新增资源，不涉及现有资源的迁移。

**部署步骤**:
1. 在 Provider 中注册新资源
2. 部署新版本的 Provider
3. 用户可以开始使用新资源

**回滚策略**:
- 如果发现严重问题，可以回退到之前的 Provider 版本
- 用户需要手动删除或导入已创建的模板资源
- 在删除 Provider 代码时，需要保留 schema 注册逻辑以避免 Terraform state 错误

## Open Questions

目前没有未解决的关键问题。所有技术决策都已经明确，可以开始实施。
