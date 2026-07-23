## Context

TencentCloud EdgeOne (TEO) 边缘函数支持绑定组件（当前仅支持 KV 命名空间类型）。云 API 提供了 `DescribeFunctionComponentBindings` 和 `ModifyFunctionComponentBindings` 两个接口来管理绑定关系。当前 Provider 已有类似的 TEO 函数配置资源（如 `tencentcloud_teo_function_runtime_environment`），本资源遵循相同的模式。

资源类型为 RESOURCE_KIND_CONFIG，意味着配置随函数存在而存在，不需要 Create/Delete 接口，只需 Read 和 Update。

## Goals / Non-Goals

**Goals:**
- 提供 `tencentcloud_teo_function_component_binding` 资源，支持声明式管理边缘函数的组件绑定配置
- 支持通过 `zone_id` + `function_id` 作为复合 ID 标识资源
- Update 使用 `rebind` 操作模式实现声明式全量绑定管理
- 支持 Import 功能
- 提供完整的单元测试（使用 gomonkey mock）

**Non-Goals:**
- 不支持增量绑定操作（bind/unbind），仅使用 rebind 实现声明式管理
- 不暴露 `filters` 参数给用户（仅在内部 Read 时使用）
- 不暴露分页参数（Offset/Limit）给用户，内部使用最大值 1000 一次性获取

## Decisions

### 1. 使用 rebind 操作模式

**决策**: Update 时使用 `Operation = "rebind"` 模式，将当前绑定列表全量替换为用户声明的列表。

**理由**: Terraform 是声明式的，用户声明的是期望状态。rebind 模式可以将所有绑定重置为传入的列表，完美匹配 Terraform 的声明式语义。相比 bind/unbind 需要计算差异，rebind 更简单且不易出错。

**替代方案**: 使用 bind + unbind 计算差异 → 复杂度高，容易出现状态不一致。

### 2. 复合 ID 格式

**决策**: 使用 `zone_id#function_id` 作为资源 ID（使用 `tccommon.FILED_SP` 分隔符）。

**理由**: 与现有 TEO 函数配置资源（如 `tencentcloud_teo_function_runtime_environment`）保持一致。

### 3. Read 实现使用分页获取全量数据

**决策**: Read 时设置 `Limit = 1000`（API 最大值），通过分页循环获取所有绑定数据。

**理由**: 确保能读取到所有绑定关系，避免数据丢失。

### 4. Create 方法实现

**决策**: Create 方法调用 `ModifyFunctionComponentBindings` 接口（Operation = "rebind"），然后设置 ID 并调用 Read。

**理由**: RESOURCE_KIND_CONFIG 类型资源的配置随函数存在而存在，Create 实际上是首次设置配置。

### 5. Delete 方法实现

**决策**: Delete 方法调用 `ModifyFunctionComponentBindings` 接口（Operation = "rebind"，传入空列表），清空所有绑定。

**理由**: 配置类资源删除时应恢复到默认状态（无绑定）。

## Risks / Trade-offs

- [Risk] rebind 操作会清空所有现有绑定再设置新的 → 如果 API 调用失败可能导致短暂的绑定丢失。Mitigation: 使用 retry 机制确保操作成功。
- [Risk] 分页查询可能遗漏数据 → Mitigation: 使用循环分页直到获取所有数据。
- [Trade-off] 不支持增量操作 → 简化了实现但每次 Update 都是全量替换，对于大量绑定可能有性能影响。实际场景中绑定数量有限，可接受。
