## Context

腾讯云 Redis 实例支持配置密码复杂度策略（密码最小总长度、字母/数字/特殊字符最小数量、是否启用），目前仅能通过控制台或云 API 管理。Provider 中 Redis 服务位于 `tencentcloud/services/crs/` 包，已有 `tencentcloud_redis_connection_config`、`tencentcloud_redis_backup_config` 等配置型资源作为参考。

本资源为 RESOURCE_KIND_CONFIG 类型：配置随实例存在而存在，只需 Read + Update（RU）接口，无独立 Create/Delete。

## Goals / Non-Goals

**Goals:**
- 提供声明式管理 Redis 实例密码复杂度策略的 Terraform 资源
- 支持 `enabled`、`min_letter_count`、`min_digit_count`、`min_special_count`、`min_length` 字段的读取与更新
- 支持 Import（使用 instance_id 导入）
- 遵循项目现有配置型资源模式（参考 `tencentcloud_redis_connection_config`）

**Non-Goals:**
- 不管理 Redis 实例本身的生命周期（由 `tencentcloud_redis_instance` 负责）
- 不实现 Delete 真实操作（无对应云 API，配置随实例存在）
- 不暴露 `password_policy` 嵌套对象作为中间 schema 层（按规则将字段平铺到顶层）

## Decisions

### 1. 字段平铺而非嵌套 password_policy 对象

**决策**: 将 `PasswordPolicy` 的子字段（`enabled`、`min_letter_count`、`min_digit_count`、`min_special_count`、`min_length`）平铺到资源 schema 顶层，不创建 `password_policy` 嵌套层。

**理由**: 云 API 出参中 `PasswordPolicy` 是单对象（非列表），用户需求映射也将其展开为顶层字段。按代码生成规则，避免多余的嵌套层使每个字段可直接 set/read。

**替代方案**: 使用 `password_policy` TypeList(MaxItems=1) 嵌套对象 —— 增加不必要的复杂度，且与需求映射不一致。

### 2. 资源 ID 使用 instance_id

**决策**: `d.SetId(instanceId)`，资源 ID 即实例 ID。

**理由**: 这是单例配置（每个实例只有一份密码策略），符合 RESOURCE_KIND_CONFIG 模式。与 `tencentcloud_redis_connection_config` 一致。

### 3. Create = SetId + Update 模式

**决策**: Create 函数设置 `d.SetId(instanceId)` 后直接调用 Update 函数。

**理由**: 配置型资源无独立创建接口，首次 apply 即为修改配置。遵循 `tencentcloud_redis_connection_config` 等现有资源的模式。

### 4. 同步接口，无需异步轮询

**决策**: `ModifyInstancePasswordPolicy` 返回值仅含 `RequestId`，不含 `TaskId`，为同步接口，调用后直接调用 Read 刷新状态即可，无需 `DescribeTaskInfo` 轮询。

**理由**: 经 vendor 中云 API 模型确认，`ModifyInstancePasswordPolicyResponse` 无 TaskId 字段。

### 5. Delete 为 no-op

**决策**: Delete 函数不做任何操作，仅返回 nil。

**理由**: 云 API 无"重置/删除密码策略"接口，密码策略配置随实例存在而存在。terraform destroy 仅从 state 中移除。

### 6. 可选字段在 Update 中使用 GetOkExists

**决策**: `enabled` 为必填字段使用 `GetOk`；`min_letter_count`、`min_digit_count`、`min_special_count`、`min_length` 为可选字段使用 `GetOkExists`，允许设置为 0。

**理由**: 这些计数字段取值范围为 [1,16]（min_length 为 [8,64]），但在 Terraform 中可选字段需支持用户不设置或设置具体值，使用 `GetOkExists` 区分"未设置"与"设置为 0"。

## Risks / Trade-offs

- **[实例不存在导致 Read 清空 state]** → 在 Read 中，若 `DescribeInstancePasswordPolicy` 返回实例不存在错误或空响应，先打印 `log.Printf("[CRUD] ... id=%s", d.Id())` 保留现场，再 `d.SetId("")`。
- **[ModifyInstancePasswordPolicy 可能因实例状态不可操作而失败]** → 使用 `tccommon.RetryError` 包装错误，让外层 retry 重试；不可重试错误（如实例不存在）返回 `NonRetryableError`。
- **[可选字段传 0 与未传值的区分]** → 使用 `GetOkExists` 处理可选整数字段，避免 0 值被误判为未设置。
