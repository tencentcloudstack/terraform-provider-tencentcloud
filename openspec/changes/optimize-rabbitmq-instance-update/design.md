# Design: 优化 RabbitMQ 实例的 update 逻辑

## Context

### 当前状态

`resourceTencentCloudTdmqRabbitmqVipInstanceUpdate` 函数（位于 `tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance.go`）当前存在以下问题：

1. **过度限制的不可变参数列表**：
   - 第 460-465 行定义了 12 个不可变参数，包括 `zone_ids`、`vpc_id`、`subnet_id`、`node_spec`、`node_num`、`storage_size`、`auto_renew_flag`、`time_span`、`pay_mode`、`cluster_version`、`band_width`、`enable_public_access`
   - 这些参数中，部分（如 `auto_renew_flag`、`band_width`、`enable_public_access`）实际上可以通过腾讯云 API 更新

2. **缺少异步状态等待**：
   - update 操作调用 API 后立即返回，没有等待资源状态更新
   - 可能导致后续操作在资源仍在更新状态时执行，造成竞态条件

3. **有限的更新能力**：
   - 仅支持更新 `cluster_name` 和 `resource_tags` 两个参数
   - 无法满足用户的实际运维需求

### 约束条件

1. **向后兼容性**：必须保持向后兼容，不能破坏现有 TF 配置和 state
2. **Schema 限制**：不能修改已有资源的 schema（除非只新增 Optional 字段）
3. **API 限制**：必须遵守腾讯云 TDMQ RabbitMQ API `ModifyRabbitMQVipInstance` 的参数限制
4. **错误处理**：需要遵循项目的错误处理模式，使用 `defer tccommon.LogElapsed()` 和 `defer tccommon.InconsistentCheck()`

### 利益相关者

- **用户**：需要能够灵活调整 RabbitMQ 实例配置
- **运维团队**：需要通过 Terraform 管理实例的生命周期
- **腾讯云 API**：提供实际的更新能力，但存在参数限制

## Goals / Non-Goals

**Goals:**

1. ✅ 移除不必要的不可变参数限制，支持更新 `auto_renew_flag`、`band_width`、`enable_public_access`
2. ✅ 添加异步状态等待逻辑，确保 update 操作完成后资源处于稳定状态
3. ✅ 优化错误信息，清晰区分不可变参数和其他错误
4. ✅ 保持完全向后兼容，不破坏现有配置
5. ✅ 遵循项目的代码风格和错误处理模式

**Non-Goals:**

1. ❌ 不支持更新涉及底层架构变更的参数（如 `zone_ids`、`vpc_id`、`subnet_id`）
2. ❌ 不支持预付费实例的特殊操作（如 `time_span`、`pay_mode` 的修改）
3. ❌ 不支持集群版本升级（`cluster_version` 的修改需要特殊流程）
4. ❌ 不修改资源 schema（所有参数都已存在）
5. ❌ 不新增额外的 API 依赖（使用现有的 `ModifyRabbitMQVipInstance` API）

## Decisions

### 1. 精简不可变参数列表

**决策**：将不可变参数列表从 12 个减少到 9 个，移除 `auto_renew_flag`、`band_width`、`enable_public_access`。

**理由**：
- 腾讯云 API `ModifyRabbitMQVipInstance` 支持更新这些参数
- 这些参数的更新不涉及底层架构变更，可以安全地通过 API 调用完成
- 用户有实际需求调整这些配置（如开启/关闭公网访问、调整带宽等）

**替代方案考虑**：
- **方案 A**：保持所有参数不可变 → ❌ 限制了用户的灵活性
- **方案 B**：移除所有不可变限制 → ❌ 可能导致 API 错误和资源损坏
- **方案 C（采纳）**：基于 API 支持情况，有选择地移除限制 → ✅ 平衡了灵活性和安全性

### 2. 添加异步状态等待逻辑

**决策**：在 update 操作完成后，添加状态等待逻辑，确保实例状态从 "Updating" 变为 "Running" 或 "Success"。

**实现方式**：
```go
// 使用 resource.Retry 进行状态轮询
err = resource.Retry(tccommon.ReadRetryTimeout*10, func() *resource.RetryError {
    result, e := service.DescribeTdmqRabbitmqVipInstanceByFilter(ctx, paramMap)
    if e != nil {
        return tccommon.RetryError(e)
    }

    if result[0].Status == svctdmq.RabbitMQVipInstanceUpdating {
        return resource.RetryableError(fmt.Errorf("rabbitmq_vip_instance status is updating"))
    } else if result[0].Status == svctdmq.RabbitMQVipInstanceRunning ||
               result[0].Status == svctdmq.RabbitMQVipInstanceSuccess {
        return nil
    } else {
        return resource.NonRetryableError(fmt.Errorf("rabbitmq_vip_instance status illegal: %s", *result[0].Status))
    }
})
```

**理由**：
- 遵循项目的一致性模式，与其他资源（如 create 操作）保持一致
- 防止竞态条件，确保资源处于稳定状态后再返回
- 提供更好的用户体验，避免后续操作失败

**替代方案考虑**：
- **方案 A**：不添加状态等待 → ❌ 可能导致竞态条件
- **方案 B**：使用固定等待时间 → ❌ 不可靠，可能等待时间不足或过长
- **方案 C（采纳）**：使用轮询机制等待状态更新 → ✅ 可靠且高效

### 3. 优化不可变参数的错误信息

**决策**：为不可变参数提供更清晰的错误信息，明确告知用户哪些参数不能修改以及原因。

**实现方式**：
```go
immutableArgs := []string{
    "zone_ids",
    "vpc_id",
    "subnet_id",
    "node_spec",
    "node_num",
    "storage_size",
    "enable_create_default_ha_mirror_queue",
    "time_span",
    "pay_mode",
    "cluster_version",
}

for _, v := range immutableArgs {
    if d.HasChange(v) {
        return fmt.Errorf("argument `%s` cannot be changed after instance creation. "+
                        "Please recreate the instance if you need to modify this parameter.", v)
    }
}
```

**理由**：
- 提供更好的用户体验，帮助用户理解为什么某些参数不能修改
- 减少支持请求，用户可以快速找到解决方案
- 遵循项目的错误处理最佳实践

### 4. 保持完全向后兼容

**决策**：只移除不必要的限制，不修改任何现有逻辑或 API 调用方式。

**理由**：
- 现有配置可以继续正常工作，不会导致任何破坏性变更
- 用户可以逐步迁移到新的功能，不需要立即修改配置
- 符合 Terraform Provider 的版本控制最佳实践

## Risks / Trade-offs

### 风险 1：API 参数支持不确定性

**描述**：腾讯云 API `ModifyRabbitMQVipInstance` 可能不完全支持我们认为可更新的参数。

**缓解措施**：
- 在实现前先查看 API 文档或进行小规模测试
- 如果 API 不支持某个参数，保持该参数不可变
- 在测试环境中验证所有参数的更新能力

### 风险 2：状态等待超时

**描述**：update 操作可能需要很长时间，导致状态等待超时。

**缓解措施**：
- 使用较长的超时时间（`tccommon.ReadRetryTimeout*10`）
- 提供清晰的超时错误信息
- 允许用户通过 Terraform Timeouts 配置自定义超时时间

### 风险 3：部分参数更新失败

**描述**：如果 API 不支持某些参数，update 操作可能部分失败。

**缓解措施**：
- 遵循 Terraform 的幂等性原则，确保 update 操作可以安全地重试
- 在测试中验证所有参数的更新能力
- 提供清晰的错误信息，帮助用户理解哪些参数更新失败

### 权衡：灵活性与安全性的平衡

**描述**：移除不可变限制会增加灵活性，但也可能导致用户误操作。

**权衡分析**：
- **更灵活**：用户可以调整更多参数，满足实际需求
- **风险增加**：用户可能会错误地修改某些参数

**缓解措施**：
- 对于真正不可变的参数（如 `zone_ids`、`vpc_id`），保持严格限制
- 提供清晰的错误信息和文档
- 在测试中验证所有参数的更新能力

## Migration Plan

### 实施步骤

1. **代码修改**：
   - 修改 `resourceTencentCloudTdmqRabbitmqVipInstanceUpdate` 函数
   - 移除不必要的不可变参数限制
   - 添加异步状态等待逻辑
   - 优化错误信息

2. **测试**：
   - 更新验收测试用例，验证新增的 update 功能
   - 测试异步状态等待逻辑
   - 验证不可变参数的错误处理

3. **文档更新**：
   - 更新资源文档，明确哪些参数支持 update
   - 添加使用示例和注意事项

4. **发布**：
   - 通过 CI/CD 流水线验证
   - 合并到主分支
   - 发布新版本

### 回滚策略

如果出现问题，可以通过以下方式回滚：
1. 恢复原始的 `immutableArgs` 列表
2. 移除异步状态等待逻辑
3. 发布修复版本

由于只移除了不必要的限制，回滚不会影响现有配置的安全性。

## Open Questions

1. **Q**: 腾讯云 API `ModifyRabbitMQVipInstance` 是否支持更新所有新增的参数？
   - **A**: 需要查看 API 文档或进行测试验证

2. **Q**: 是否需要支持更细粒度的参数更新控制（如通过 feature flag）？
   - **A**: 不需要，直接基于 API 支持情况决定

3. **Q**: 是否需要添加更多的状态常量（如 "Maintenance"、"Rollback" 等）？
   - **A**: 需要根据 API 返回的状态值和实际测试结果决定
