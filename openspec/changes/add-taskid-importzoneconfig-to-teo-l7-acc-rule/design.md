## Context

TEO (Tencent Edge One) 是腾讯云的边缘安全加速服务，提供 L7 访问控制规则等安全能力。tencentcloud_teo_l7_acc_rule 资源用于管理这些访问控制规则。

当前的实现中，更新操作通过 ImportZoneConfig API 进行。这是一个异步 API，返回 TaskId 参数用于跟踪任务状态。但是，现有代码没有实现基于 TaskId 的等待逻辑，导致更新操作可能在任务未完成时就返回，造成状态不一致的问题。

在 terraform-provider-tencentcloud 中，异步操作的等待是一个常见模式，通常使用 helper.Retry() 函数轮询任务状态直到完成或超时。

## Goals / Non-Goals

**Goals:**
- 在 tencentcloud_teo_l7_acc_rule 的 Update 函数中实现基于 TaskId 的异步等待逻辑
- 通过轮询机制跟踪任务执行状态，确保操作完成后再返回
- 保持与现有 provider 代码风格一致
- 确保不破坏向后兼容性

**Non-Goals:**
- 不修改资源的 schema（不新增或修改 Terraform 配置参数）
- 不修改 Create 和 Delete 操作（仅影响 Update）
- 不改变现有的错误处理机制

## Decisions

### 1. 使用现有的等待机制
选择使用 helper.Retry() 函数而不是手动实现轮询循环。

**理由：**
- provider 中的其他资源（如 CDN、DNSPod）广泛使用此模式
- helper.Retry() 已集成了超时、重试和错误处理
- 代码风格一致，易于维护

### 2. 任务状态查询 API
通过 DescribeZoneConfigRollBackTasks 或类似的查询 API 来获取任务状态。

**理由：**
- TEO 提供了专门的任务查询接口
- 可以获取详细的任务状态信息（成功/失败/进行中）
- 比直接轮询原 API 更高效

### 3. 轮询间隔和超时配置
- 轮询间隔：5 秒
- 超时时间：使用资源的默认 Timeout 配置（如果存在）或使用默认 10 分钟

**理由：**
- 5 秒间隔平衡了 API 调用频率和响应时间
- 10 分钟是合理的异步操作超时时间
- 复用现有的 Timeout 配置机制

## Risks / Trade-offs

### 风险 1: API 调用失败导致无限等待
**风险**: 如果任务查询 API 返回错误，可能导致重试逻辑异常。

**缓解措施**:
- 在查询 API 调用失败时返回错误，停止等待
- 记录详细的错误日志以便排查
- 添加最终一致性检查机制

### 风险 2: 任务状态枚举值变化
**风险**: TEO API 的任务状态枚举值可能在未来版本中发生变化。

**缓解措施**:
- 使用字符串匹配而不是硬编码的整数常量
- 添加注释说明预期的状态值
- 如果 API 返回未知状态，记录警告并继续等待

### 权衡：额外的 API 调用开销
**权衡**: 实现等待逻辑会增加 API 调用次数，可能影响操作耗时。

**缓解措施**:
- 轮询间隔设置为 5 秒，避免过于频繁的调用
- 一旦任务完成立即停止轮询
- 大多数情况下任务会在较短时间内完成

## Migration Plan

1. **代码变更**:
   - 修改 `resource_tencentcloud_teo_l7_acc_rule.go` 的 Update 函数
   - 添加任务状态查询的辅助函数
   - 添加单元测试和集成测试

2. **测试验证**:
   - 运行现有的资源测试确保无回归
   - 添加针对异步等待逻辑的测试用例
   - 手动测试验证任务等待行为

3. **文档更新**:
   - 更新 `website/docs/r/teo_l7_acc_rule.md` 说明异步操作的行为
   - 如有必要，添加使用示例

4. **发布**: 通过正常的 PR 流程发布，确保 code review 和 CI 通过

## Open Questions

1. 具体的任务查询 API 接口是什么？
   - 需要查阅 TEO API 文档或现有代码以确定正确的查询接口
   - 可能在 `DescribeZoneConfigRollBackTasks` 或类似接口中

2. 任务状态的成功/失败值是什么？
   - 需要通过测试或文档确认
   - 预期类似 "Success", "Failed", "Processing" 等字符串值

3. 是否需要支持部分成功的情况？
   - 当前假设任务是原子操作（全成功或全失败）
   - 如果 API 支持部分成功，需要相应调整逻辑
