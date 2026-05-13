## Context

腾讯云 CVM 通过 `RepairTaskControl` 接口对"待授权"维修任务进行授权操作。该 API 是一次性的副作用操作，并非用于管理持久化资源——它没有对应的 Read/Update/Delete API，本质上是触发一次云侧动作。

Terraform Provider 内已有大量类似的"操作型资源"（operation resource）设计模式，例如：
- `tencentcloud_dcdb_flush_binlog_operation`
- `tencentcloud_mysql_restart_db_instances_operation`
- `tencentcloud_dts_sync_job_resume_operation`

这些资源命名后缀为 `_operation`，采用"Create 时触发动作、Read/Delete 为空操作"的模式。本次新增遵循同一模式。

约束:
- 必须使用 `tencentcloud-sdk-go/tencentcloud/cvm/v20170312` 中的 `RepairTaskControlRequest`
- 必须保持向后兼容
- 必须有完整文档与测试
- 必须支持 Terraform 重试机制处理瞬时错误

## Goals / Non-Goals

**Goals:**
- 提供 `tencentcloud_cvm_repair_task_control_operation` 资源以触发维修任务授权
- 完整支持 API 文档的所有参数：`Product`、`InstanceIds`、`TaskId`、`Operate`、`OrderAuthTime`、`TaskSubMethod`
- 与现有 `tencentcloud_cvm_repair_tasks` 数据源形成"查询 + 授权"完整闭环
- 遵循 Provider 既有的 operation 资源模式，保持代码一致性

**Non-Goals:**
- 不实现维修任务的取消、查询状态变化等其他控制操作（API 当前仅支持 `AuthorizeRepair`）
- 不在资源中维护任务的实时状态（Read 为空操作，state 不持久化授权进度）
- 不实现 Update（操作为一次性，更新参数等价于重新创建）
- 不引入新的 SDK 依赖

## Decisions

### Decision 1: 采用 Operation Resource 模式而非 Resource

**选择**: 使用以 `_operation` 结尾的操作型资源命名与模式

**理由**:
- `RepairTaskControl` 没有对应的 Describe/Modify/Delete API，不符合传统 CRUD 资源的契约
- Provider 已有数十个同类 `_operation` 资源，模式成熟一致
- 用户可以通过 `taint` 重新触发，符合 Terraform 用户的预期

**替代方案**:
- 实现为常规 Resource: 不可行——缺少 Read 接口无法保证一致性
- 实现为 Provisioner: 不符合 Provider 项目结构，可发现性差

### Decision 2: Schema 设计

字段映射如下：

| Terraform 字段 | 类型 | Required | ForceNew | 对应 API 字段 |
|---|---|---|---|---|
| `product` | String | 是 | 是 | `Product` |
| `instance_ids` | List(String) | 是 | 是 | `InstanceIds` |
| `task_id` | String | 是 | 是 | `TaskId` |
| `operate` | String | 是 | 是 | `Operate` |
| `order_auth_time` | String | 否 | 是 | `OrderAuthTime` |
| `task_sub_method` | String | 否 | 是 | `TaskSubMethod` |

**关键决策**:
- 所有字段都设置 `ForceNew: true`：因为操作是一次性的，参数变更应当触发重新执行
- `operate` 即使当前只支持 `AuthorizeRepair`，也作为参数暴露以便将来扩展
- 资源 ID 使用 `task_id` 作为主键（已完成授权的维修任务 ID 由 API 返回）

**替代方案**:
- 将 `operate` 硬编码为 `AuthorizeRepair`：放弃灵活性，将来 API 扩展需破坏性变更
- ID 用 UUID：失去与原始任务的关联，不利于调试

### Decision 3: 错误处理与重试

**选择**: 使用 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 包装 API 调用，使用 `tccommon.RetryError(e)` 区分可重试与不可重试错误

**理由**: 与项目所有 operation 资源保持一致，处理网络瞬时故障

### Decision 4: Read 与 Delete 实现

**选择**: Read 与 Delete 均为空操作，仅记录日志

**理由**:
- API 没有反向操作（无法"取消授权"）
- 无 Describe 接口供 Read 校验状态
- 这是项目内 operation 资源的统一模式

### Decision 5: 文件命名

- 实现文件: `resource_tc_cvm_repair_task_control_operation.go`
- 测试文件: `resource_tc_cvm_repair_task_control_operation_test.go`
- 文档模板: `resource_tc_cvm_repair_task_control_operation.md`
- Terraform 资源名: `tencentcloud_cvm_repair_task_control_operation`

遵循 `tencentcloud_<service>_<api-name-snake>_operation` 的命名约定。

## Risks / Trade-offs

- **[风险] 操作不可逆**: 授权后无法撤销 → **缓解**: 在文档中明确说明，并明确警告 `task_sub_method=LossyLocal` 会清空本地盘数据
- **[风险] 预约时间校验复杂**: API 要求 5 分钟之后、48 小时之内 → **缓解**: 由云端校验，错误信息透传给用户；不在 Provider 端做时间窗口校验（避免时区漂移问题）
- **[风险] State 不反映真实云侧状态**: Read 为空，云端状态变化 Terraform 不感知 → **缓解**: 这是 operation 资源固有特性，文档中明确说明，建议配合 `tencentcloud_cvm_repair_tasks` 数据源查询实际状态
- **[风险] taint 重新触发可能失败**: 如果任务已不在"待授权"状态，重新执行会报错 → **缓解**: 错误信息由云端透传，符合用户预期
- **[权衡] ForceNew 全字段**: 任何参数变化都重建资源，提升清晰度但失去 in-place 更新能力 → 接受，因为操作本质上无 update 语义

## Migration Plan

无迁移工作。这是一个全新的资源，不影响任何现有 TF 配置或 state。

部署步骤：
1. 合入代码后用户更新 Provider 版本即可使用
2. 用户在 `.tf` 文件中添加 `tencentcloud_cvm_repair_task_control_operation` 资源块
3. 执行 `terraform apply` 触发授权

回滚: 移除资源块并 `terraform state rm` 即可，云端已执行的授权动作不会回滚（也无法回滚，符合预期）。

## Open Questions

- 是否需要在 Provider 端对 `product` 字段做枚举校验？  
  **决策**: 暂不校验，保持向前兼容（API 未来可能新增产品类型）。云端会校验非法值。
- 是否需要在 Provider 端对 `operate` 字段做枚举校验？  
  **决策**: 暂不校验，理由同上。
