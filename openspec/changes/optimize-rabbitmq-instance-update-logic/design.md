## Context

当前 `tencentcloud_tdmq_rabbitmq_vip_instance` 资源的 update 函数（`resourceTencentCloudTdmqRabbitmqVipInstanceUpdate`）将 11 个参数标记为不可修改参数（immutableArgs），包括：

- `zone_ids` (可用区)
- `vpc_id` (VPC ID)
- `subnet_id` (子网 ID)
- `node_spec` (节点规格)
- `node_num` (节点数量)
- `storage_size` (存储大小)
- `enable_create_default_ha_mirror_queue` (镜像队列)
- `auto_renew_flag` (自动续费标识)
- `time_span` (购买时长)
- `pay_mode` (付费模式)
- `cluster_version` (集群版本)
- `band_width` (公网带宽)
- `enable_public_access` (公网访问开关)

当前只支持修改 `cluster_name` 和 `resource_tags` 两个参数。

然而，腾讯云的 `ModifyRabbitMQVipInstance` API 实际上支持更多参数的修改，包括：
- `AutoRenewFlag` - 自动续费标识
- `Bandwidth` - 公网带宽
- `EnablePublicAccess` - 公网访问开关

这种过度限制导致用户无法通过 Terraform 灵活地管理实例配置，必须手动通过控制台或 API 进行修改。

## Goals / Non-Goals

**Goals:**
- 扩展 update 函数支持更多参数的在线修改
- 实现对 `auto_renew_flag`、`band_width`、`enable_public_access` 参数的更新支持
- 添加异步操作等待机制，确保长时间更新操作的可靠性
- 保持完全向后兼容，不破坏现有用户配置

**Non-Goals:**
- 不支持 `zone_ids`、`vpc_id`、`subnet_id` 等核心网络参数的修改（这些参数需要重建实例）
- 不支持 `node_spec`、`node_num`、`storage_size` 等规格参数的修改（这些参数可能需要扩容操作，不在本次优化范围内）
- 不修改资源的 schema 定义（仅增强 update 逻辑）

## Decisions

### 1. 参数修改范围选择

**决策**: 从 `immutableArgs` 列表中移除 `auto_renew_flag`、`band_width`、`enable_public_access` 三个参数

**理由**:
- 腾讯云 `ModifyRabbitMQVipInstance` API 明确支持这三个参数的修改
- 这三个参数是常见的运维需求，用户需要频繁调整
- 修改这三个参数不需要重建实例，可以在线完成
- 这些参数的修改操作相对简单，风险较低

**考虑的替代方案**:
- 方案 A: 将所有 API 支持的参数都移出 `immutableArgs`
  - 优点: 最大化灵活性
  - 缺点: 风险较高，需要更多测试和验证，可能触发实例重建或长时间扩容操作

- 方案 B: 只移除这三个低风险参数
  - 优点: 风险可控，实施成本低，易于测试
  - 缺点: 功能扩展有限
  - **最终选择**: 方案 B，平衡了风险和收益

### 2. 异步操作等待机制

**决策**: 对于 `band_width` 和 `enable_public_access` 两个参数的修改，添加异步等待机制

**理由**:
- 公网带宽修改和公网访问开关可能需要较长时间生效
- 在操作完成前立即返回可能导致后续操作失败或状态不一致
- 用户需要知道操作何时真正完成

**实现方式**:
- 使用现有的 `resource.Retry()` 机制进行等待
- 等待超时使用 `tccommon.WriteRetryTimeout` 或更长的 `tccommon.WriteRetryTimeout * 2`
- 通过 `DescribeRabbitMQVipInstances` API 轮询检查操作是否完成
- 判断依据: 检查 `band_width` 和 `enable_public_access` 的值是否更新为目标值

### 3. 错误处理和日志记录

**决策**: 完善错误处理和日志记录

**理由**:
- 提供更好的调试信息，帮助用户排查问题
- 符合项目现有的错误处理模式
- 便于追踪更新操作的执行情况

**实现方式**:
- 在修改操作完成后，添加日志记录修改的参数和结果
- 对于 API 返回的错误，提供更详细的错误信息
- 使用现有的 `log.Printf` 模式进行日志输出

## Risks / Trade-offs

### 1. API 兼容性风险

**风险**: 腾讯云 API 可能在未来版本中修改或删除对这些参数的支持

**缓解措施**:
- 在代码中添加 API 错误处理，当 API 返回不支持参数的错误时，给出明确的错误提示
- 通过集成测试验证 API 行为，确保参数修改功能正常工作
- 定期检查腾讯云 API 文档更新

### 2. 状态一致性风险

**风险**: 在异步等待期间，状态可能出现短暂的不一致

**缓解措施**:
- 使用现有的 `defer tccommon.InconsistentCheck(d, meta)()` 机制检测状态不一致
- 在 update 函数结束时，调用 `resourceTencentCloudTdmqRabbitmqVipInstanceRead` 刷新状态
- 确保在等待超时后仍然刷新状态，避免状态不一致

### 3. 测试覆盖风险

**风险**: 新增的 update 逻辑可能缺少足够的测试覆盖

**缓解措施**:
- 在 `resource_tc_tdmq_rabbitmq_vip_instance_test.go` 中添加对应的测试用例
- 测试用例覆盖:
  - 修改 `auto_renew_flag` 的场景
  - 修改 `band_width` 的场景
  - 修改 `enable_public_access` 的场景
  - 同时修改多个参数的场景
  - 修改不可变参数应返回错误

### 4. 性能影响

**风险**: 添加异步等待机制可能增加 update 操作的执行时间

**缓解措施**:
- 只对必要的参数添加等待机制（`band_width` 和 `enable_public_access`）
- 对于 `auto_renew_flag` 的修改，不添加等待，因为该修改通常立即生效
- 合理设置等待超时时间，避免过长的等待时间

## Migration Plan

### 部署步骤

1. 修改 `resourceTencentCloudTdmqRabbitmqVipInstanceUpdate` 函数
   - 从 `immutableArgs` 列表中移除 `auto_renew_flag`、`band_width`、`enable_public_access`
   - 添加这三个参数的修改逻辑
   - 为 `band_width` 和 `enable_public_access` 添加异步等待机制

2. 添加测试用例
   - 在 `resource_tc_tdmq_rabbitmq_vip_instance_test.go` 中添加 update 相关的测试用例
   - 确保测试覆盖所有新增的修改逻辑

3. 更新文档（如果需要）
   - 检查 `website/docs/r/tdmq_rabbitmq_vip_instance.md` 是否需要更新
   - 在文档中明确说明哪些参数可以在线修改

### 回滚策略

如果新功能出现问题，可以快速回滚：

1. 将 `auto_renew_flag`、`band_width`、`enable_public_access` 重新添加到 `immutableArgs` 列表
2. 删除或注释掉新增的参数修改逻辑和异步等待代码
3. 通过 git revert 回滚代码变更

由于本次变更只增加了可修改参数的范围，不删除现有功能，因此回滚不会影响现有用户的使用体验。

## Open Questions

1. **等待超时时间**: 对于 `band_width` 和 `enable_public_access` 的修改，应该使用多长的等待超时时间？建议使用 `tccommon.WriteRetryTimeout * 2`，但需要根据实际测试结果调整。

2. **状态检查频率**: 在异步等待期间，轮询检查状态的间隔时间是多少？建议使用默认的 retry 间隔，或者根据操作类型调整。

3. **错误提示信息**: 当修改不可变参数时，当前的错误提示信息是否足够清晰？是否需要改进以告知用户哪些参数可以修改？

4. **文档更新**: 是否需要在文档中明确说明所有可修改的参数列表？建议在资源文档中添加一个"可修改参数"章节。
