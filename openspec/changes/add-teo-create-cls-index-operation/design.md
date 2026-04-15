## Context

TEO (TencentCloud EdgeOne) 服务提供了实时日志投递功能，允许将边缘访问日志投递到腾讯云 CLS (Cloud Log Service)。在创建实时日志投递任务后，需要为该任务创建 CLS 索引以便进行日志检索和分析。

当前 TEO 服务已支持创建实时日志投递任务（`resource_tc_teo_realtime_log_delivery`），但缺少为这些任务创建 CLS 索引的 Terraform 资源。用户需要手动通过云 API 或控制台创建 CLS 索引，这在自动化场景中不够方便。

## Goals / Non-Goals

**Goals:**
- 为 TEO 服务提供创建 CLS 索引的 Terraform operation 资源
- 支持通过站点 ID 和实时日志投递任务 ID 创建 CLS 索引
- 遵循 Terraform provider 的代码规范和模式
- 提供完整的单元测试和文档

**Non-Goals:**
- 不提供读取、更新、删除 CLS 索引的功能（这是 operation 资源的特性）
- 不管理 CLS 索引的生命周期（索引创建后的管理由 CLS 服务负责）
- 不支持异步轮询（CreateCLSIndex 是同步接口，调用后立即返回）

## Decisions

### 资源类型选择为 operation

**决策：** 将资源类型定义为 operation（一次性操作），而不是常规 resource。

**理由：**
- CreateCLSIndex 接口是幂等操作，可以多次调用
- 操作执行完成后不需要维护状态，不需要 Read、Update、Delete 接口
- 参考现有的 operation 资源实现（如 `resource_tc_teo_l7_acc_rule_priority_operation`）
- 符合 TEO 服务的操作类资源模式

**替代方案考虑：**
- 常规 resource：需要实现 CRUD 接口，但云 API 不提供对应的 Read、Update、Delete 接口，会增加复杂度

### Schema 设计

**决策：** Schema 包含两个 Required 参数：`zone_id` 和 `task_id`。

**理由：**
- CreateCLSIndexRequest 接口只需要这两个参数
- 两个参数都是操作必需的
- 所有参数设置 ForceNew: true（operation 资源的典型模式）

**Schema 结构：**
```go
Schema: map[string]*schema.Schema{
    "zone_id": {
        Type:        schema.TypeString,
        Required:    true,
        ForceNew:    true,
        Description: "站点 ID。",
    },
    "task_id": {
        Type:        schema.TypeString,
        Required:    true,
        ForceNew:    true,
        Description: "实时日志投递任务 ID。",
    },
}
```

### CRUD 接口实现

**决策：**
- Create：调用 `CreateCLSIndex` 接口，使用 `tccommon.WriteRetryTimeout` 处理重试
- Read：空实现（返回 nil）
- Delete：空实现（返回 nil）
- Update：不提供（operation 资源不支持）

**理由：**
- CreateCLSIndex 是同步接口，调用后立即返回
- operation 资源不需要维护状态
- Read 和 Delete 空实现可以避免 Terraform state 不一致的问题
- 所有参数都是 ForceNew，不需要 Update 接口

### 资源 ID 设置

**决策：** 将资源 ID 设置为 `zone_id`。

**理由：**
- operation 资源的 ID 主要用于区分不同的资源实例
- zone_id 具有唯一性，可以标识资源
- 参考现有 operation 资源实现

### 错误处理

**决策：**
- 使用 `tccommon.LogElapsed` 记录操作耗时
- 使用 `tccommon.InconsistentCheck` 进行一致性检查
- 使用 `tccommon.RetryError` 处理可重试错误
- API 调用失败时返回原始错误信息

**理由：**
- 遵循 Terraform provider 的错误处理模式
- 提供足够的日志信息用于问题排查
- 处理云 API 的限流和临时性错误

## Risks / Trade-offs

### 风险：API 参数变更

**风险：** CreateCLSIndex 接口参数可能在未来版本中发生变更。

**缓解措施：**
- 依赖 vendor 目录中的云 SDK，定期更新依赖
- 在 CI/CD 中添加参数校验测试
- 关注 TEO 服务的 API 变更公告

### 风险：幂等性

**风险：** CreateCLSIndex 接口的幂等性未在文档中明确说明，可能存在重复创建的问题。

**缓解措施：**
- 使用 ForceNew: true 确保参数变更时重新创建
- 在测试中验证多次调用同一参数的行为
- 在文档中说明操作的幂等性

### 权衡：简化 vs. 完整性

**权衡：** Read 和 Delete 接口为空实现，简化了代码但限制了资源管理的完整性。

**说明：**
- 这是 operation 资源的设计选择，符合需求
- 如果将来需要管理索引生命周期，可以添加新的 resource 资源
- 保持当前实现的简洁性

## Migration Plan

**部署步骤：**
1. 在 `tencentcloud/services/teo/` 目录下创建资源文件
2. 在 `service_tencentcloud_teo.go` 中注册新资源
3. 运行单元测试确保代码正确性
4. 运行集成测试验证与云 API 的集成
5. 更新文档并提供使用示例

**回滚策略：**
- 如果发现问题，可以从 `service_tencentcloud_teo.go` 中移除资源注册
- 删除新增的资源文件
- 不会影响现有功能和用户配置

## Open Questions

无（当前设计已明确实现方案）。
