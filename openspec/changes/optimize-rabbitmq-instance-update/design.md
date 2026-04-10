## Context

当前 Terraform Provider for TencentCloud 的 RabbitMQ VIP Instance 资源（`resource_tc_tdmq_rabbitmq_vip_instance.go`）仅支持更新 `cluster_name` 和 `resource_tags` 两个参数。代码位于 `tencentcloud/services/trabbit/` 目录下，update 函数定义在第 450-523 行。

该资源使用腾讯云 TDMQ SDK 的 `ModifyRabbitMQVipInstance` API（版本 `v20200217`）进行实例配置修改。当前 update 函数中定义了大量不可变参数（第 460-465 行），包括 `zone_ids`、`vpc_id`、`subnet_id`、`node_spec`、`node_num`、`storage_size`、`enable_create_default_ha_mirror_queue`、`auto_renew_flag`、`time_span`、`pay_mode`、`cluster_version`、`band_width`、`enable_public_access`。

这些参数中很多在腾讯云侧实际支持修改，但被错误地标记为不可变，导致用户无法通过 Terraform 灵活调整实例配置。

## Goals / Non-Goals

**Goals:**
- 支持通过 Terraform 更新 RabbitMQ VIP Instance 的关键配置参数
- 移除不必要的不可变参数限制，提升资源管理的灵活性
- 确保代码与腾讯云 TDMQ API 的实际能力保持一致
- 保持向后兼容，不影响现有 Terraform 配置和 state

**Non-Goals:**
- 不修改资源的 schema 结构（不新增、删除或修改字段定义）
- 不改变资源的 create 和 read 逻辑
- 不新增资源或数据源

## Decisions

### 1. API 能力调研优先
在实现前，必须先调研腾讯云 TDMQ `ModifyRabbitMQVipInstance` API 实际支持的参数。通过查阅腾讯云 SDK 文档或 API 文档，确认哪些参数可以在实例创建后修改。

**理由**：避免编写无法工作的代码，确保实现符合云服务实际能力。如果某些参数确实不支持修改，则需要保持原有的不可变标记。

### 2. 逐步移除不可变参数
根据 API 能力调研结果，逐步将支持修改的参数从 `immutableArgs` 列表中移除，并为每个参数添加相应的更新逻辑。

**理由**：降低风险，便于测试和验证每个参数的更新功能。

### 3. 参数更新顺序优化
对于可能存在依赖关系的参数（如 `node_num` 可能影响 `storage_size`），需要确定合理的更新顺序，避免冲突或错误。

**理由**：某些参数的更新可能需要按特定顺序执行，否则可能导致 API 调用失败。

### 4. 异步操作处理
对于需要较长时间生效的参数变更（如节点数量扩容），需要添加等待逻辑，确保更新完成后才返回，避免 Terraform state 不一致。

**理由**：RabbitMQ 实例的某些配置变更是异步的，需要等待实例状态变为稳定后才能继续。

## Risks / Trade-offs

**风险 1**: ModifyRabbitMQVipInstance API 可能不支持所有预期参数的更新
→ **缓解措施**: 在实现前进行充分的 API 能力调研，必要时联系腾讯云技术支持确认。如果某些参数确实不支持更新，则保留其在 `immutableArgs` 列表中的位置。

**风险 2**: 某些参数的更新可能触发实例重启或服务中断
→ **缓解措施**: 在文档中明确标注哪些参数的更新可能导致服务中断，提醒用户在生产环境中谨慎使用。对于关键参数的更新，可以考虑添加确认机制。

**风险 3**: 向后兼容性问题
→ **缓解措施**: 不修改 schema 的任何字段定义，仅修改 update 函数的内部逻辑。确保现有的 Terraform 配置文件无需修改即可正常工作。保持原有的 `immutableArgs` 检查逻辑，仅移除被确认可更新的参数。

**风险 4**: 并发更新导致状态不一致
→ **缓解措施**: 使用 Terraform 的标准并发控制机制，确保同一资源的更新操作串行执行。

## Migration Plan

本变更不涉及数据迁移或 schema 变更，因此无需特殊的迁移步骤。用户只需更新到新版本的 Terraform Provider 即可获得增强的更新能力。

部署步骤：
1. 在开发环境中完成代码修改和单元测试
2. 在测试环境中进行完整的集成测试
3. 更新相关文档（website/docs/ 目录下的文档）
4. 发布新版本的 Provider

回滚策略：如果出现问题，可以回退到旧版本的 Provider，不会影响用户已创建的资源（因为仅修改了 update 函数的逻辑）。

## Open Questions

1. **ModifyRabbitMQVipInstance API 的确切能力是什么？**
   需要查阅腾讯云 TDMQ API 文档，确认哪些参数可以被修改。特别是：
   - `node_spec`、`node_num`、`storage_size` 是否支持修改？
   - `band_width`、`auto_renew_flag`、`enable_public_access` 是否支持修改？
   - 这些参数的修改是否有前置条件或限制？

2. **参数更新是否需要等待操作完成？**
   某些参数的更新（如节点数量变更）可能需要较长时间生效，需要确认是否需要在 update 函数中添加等待逻辑，以及等待的判断标准是什么。

3. **是否存在参数之间的依赖关系？**
   例如，`node_num` 的修改是否会影响 `storage_size`，或者反之亦然？需要确定参数更新的最佳顺序。
