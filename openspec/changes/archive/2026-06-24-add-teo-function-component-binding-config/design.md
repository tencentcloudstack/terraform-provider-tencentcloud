## Context

TencentCloud EdgeOne (TEO) 边缘函数支持绑定组件（当前仅支持 KV 命名空间类型），通过 `DescribeFunctionComponentBindings` 和 `ModifyFunctionComponentBindings` 两个云 API 进行管理。该资源为 RESOURCE_KIND_CONFIG 类型，即配置型资源，只需要 Read 和 Update 方法，不需要 Create 和 Delete。

当前 TEO 服务已有多个资源实现（如 `tencentcloud_teo_acceleration_domain`、`tencentcloud_teo_certificate_config` 等），可复用现有的服务层模式。

## Goals / Non-Goals

**Goals:**
- 实现 `tencentcloud_teo_function_component_binding` 资源，支持声明式管理边缘函数的组件绑定配置
- 使用 `zone_id` + `function_id` 作为联合 ID（使用 `tccommon.FILED_SP` 分隔）
- Read 方法通过 `DescribeFunctionComponentBindings` 接口获取绑定列表（支持分页，Limit 最大值 1000）
- Update 方法通过 `ModifyFunctionComponentBindings` 接口使用 `rebind` 操作模式实现声明式更新
- 提供完整的单元测试（使用 gomonkey mock 云 API）

**Non-Goals:**
- 不实现 Create/Delete 方法（CONFIG 类型资源无需）
- 不支持 import（CONFIG 类型资源通常不需要 import）
- 不实现过滤条件（filters）参数，Read 时获取全量绑定列表

## Decisions

1. **Update 操作使用 `rebind` 模式**
   - 理由：`rebind` 模式会清空所有现有绑定并设置为传入的绑定列表，符合 Terraform 声明式管理的语义（用户声明期望状态，系统自动对齐）
   - 替代方案：使用 `bind`/`unbind` 组合操作需要计算差异，增加复杂度且容易出错

2. **资源 ID 使用 zone_id + function_id 联合 ID**
   - 理由：CONFIG 类型资源的配置依附于 zone 和 function，两者共同唯一标识一个配置
   - 使用 `tccommon.FILED_SP` 作为分隔符，与项目其他资源保持一致

3. **Read 方法使用分页获取全量数据**
   - 理由：`DescribeFunctionComponentBindings` 接口支持分页，Limit 最大值为 1000，需要循环获取所有绑定
   - 在 Read 中设置 Limit 为 1000（接口最大值），通过 Offset 分页获取全部数据

4. **zone_id 和 function_id 设置为 Required + ForceNew**
   - 理由：作为资源标识的一部分，变更时需要重建资源

5. **function_component_bindings 字段设计为 List 类型**
   - 理由：绑定列表是有序的，每个绑定包含 type、variable_name 和 kv_namespace_parameters 子结构

## Risks / Trade-offs

- [Risk] `rebind` 操作会清空所有现有绑定再重新设置 → 如果 Terraform 配置中遗漏了某些绑定，会导致这些绑定被删除。通过文档明确说明该行为来缓解。
- [Risk] 分页查询可能因网络问题导致数据不完整 → 使用 retry 机制包装 API 调用，确保可靠性。
