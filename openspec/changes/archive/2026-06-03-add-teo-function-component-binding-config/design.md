## Context

TEO (TencentCloud EdgeOne) 边缘函数支持绑定组件（当前支持 KV 命名空间），通过 `DescribeFunctionComponentBindings` 和 `ModifyFunctionComponentBindings` 两个云 API 管理绑定关系。本资源为 RESOURCE_KIND_CONFIG 类型，配置随函数存在而存在，仅需 Read 和 Update 操作。

现有 TEO 资源（如 `tencentcloud_teo_prefetch_origin_limit`）已建立了 CONFIG 类型资源的实现模式：Create 方法调用 Modify 接口设置初始配置，Read 方法查询配置，Update 方法修改配置，Delete 方法重置/清空配置。

## Goals / Non-Goals

**Goals:**
- 提供 `tencentcloud_teo_function_component_binding` 资源，允许用户通过 Terraform 管理边缘函数的组件绑定配置
- 支持 Import 功能，使用 `zone_id#function_id` 联合 ID
- Read 方法自动分页获取所有绑定
- Update 方法使用 `rebind` 操作模式实现声明式配置管理

**Non-Goals:**
- 不支持单个绑定的增量管理（bind/unbind），统一使用 rebind 实现声明式全量管理
- 不暴露 filters 参数给用户（Read 内部获取全量数据）
- 不支持异步轮询（ModifyFunctionComponentBindings 为同步接口）

## Decisions

### 1. 使用 rebind 操作模式实现声明式管理

**选择**: Update 方法固定使用 `rebind` 操作类型
**理由**: Terraform 资源管理是声明式的，用户在 HCL 中声明期望的绑定列表，provider 应确保实际状态与声明一致。`rebind` 操作会清空所有现有绑定并设置为传入的绑定列表，完美匹配 Terraform 的声明式语义。
**替代方案**: 使用 `bind`/`unbind` 进行增量管理，但这需要 diff 计算且容易出现状态不一致。

### 2. 资源 ID 使用 zone_id 和 function_id 联合 ID

**选择**: `d.SetId(strings.Join([]string{zoneId, functionId}, tccommon.FILED_SP))`
**理由**: 函数组件绑定配置由 zone_id 和 function_id 唯一确定，使用 `tccommon.FILED_SP` 作为分隔符符合项目约定。

### 3. Create 方法调用 ModifyFunctionComponentBindings 设置初始配置

**选择**: Create 方法使用 `rebind` 操作设置初始绑定列表
**理由**: CONFIG 类型资源没有独立的 Create 接口，配置随资源存在。Create 方法实际上是设置初始配置值。

### 4. Delete 方法使用 rebind 传入空列表清空绑定

**选择**: Delete 方法调用 `ModifyFunctionComponentBindings` 并传入空的 `FunctionComponentBindings` 列表，`Operation` 为 `rebind`
**理由**: 根据 API 文档，"当 Operation 为 rebind 且传入空列表时，表示清空所有绑定"。这实现了 CONFIG 资源的 Delete 语义（恢复默认状态）。

### 5. operation 字段不暴露给用户

**选择**: `operation` 不作为 schema 字段暴露，内部固定使用 `rebind`
**理由**: 声明式管理下，操作类型由 provider 内部决定，用户只需声明期望的绑定列表。

## Risks / Trade-offs

- [Risk] rebind 操作会清空所有现有绑定再重新设置 → 如果用户在 Terraform 外手动添加了绑定，执行 apply 时会被清除。这是 Terraform 声明式管理的预期行为。
- [Risk] 分页查询可能返回大量数据 → 使用最大 Limit=1000 减少 API 调用次数，循环分页获取全部数据。
