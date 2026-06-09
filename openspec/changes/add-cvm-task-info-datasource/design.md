## Context

腾讯云 CVM 提供了 `DescribeTaskInfo` API 用于查询云服务器维修任务列表。该 API 支持多种过滤条件(实例ID、任务状态、任务类型、时间范围等)和分页查询。

当前 terraform-provider-tencentcloud 在 `tencentcloud/services/cvm/` 目录下已有多个 CVM 相关的数据源实现,遵循统一的代码组织模式。新增数据源需要遵循现有的架构模式和编码规范。

## Goals / Non-Goals

**Goals:**
- 实现 `tencentcloud_cvm_repair_tasks` 数据源,完整映射 `DescribeTaskInfo` API 的请求参数和响应字段
- 支持所有官方 API 提供的过滤条件:产品类型、任务状态、任务类型、实例ID、实例名称、时间范围、排序方式
- 遵循 Provider 现有的代码组织结构和命名规范
- 提供完整的验收测试和文档

**Non-Goals:**
- 不实现维修任务的创建、修改、删除操作(这些操作由云平台自动触发,用户只能查询和授权)
- 不实现任务授权操作(这属于资源管理,应该是 resource 而非 data source)
- 不对 API 返回的数据进行额外的聚合或统计

## Decisions

### 1. 数据源命名
**决策**: 使用 `tencentcloud_cvm_repair_tasks` 而非 `tencentcloud_cvm_task_info`

**理由**: 
- API 名称 `DescribeTaskInfo` 中的 "Task" 指的是"维修任务"(Repair Task)
- 使用 `repair_tasks` 更语义化,用户一看就知道是查询维修任务
- 与 API 响应中的 `RepairTaskInfo` 类型命名保持一致

### 2. Schema 设计
**决策**: 将所有 API 请求参数映射为数据源的可选参数,将 API 响应映射为 computed 字段

**理由**:
- 数据源的作用是查询,所有过滤条件应该是 Optional 而非 Required
- 遵循 Terraform 数据源的最佳实践:输入参数用于过滤,输出参数用于读取结果
- 参考现有数据源(如 `data_source_tc_instances.go`)的实现模式

**Schema 映射**:
- 输入参数(Optional):
  - `product` → `Product`
  - `task_status` → `TaskStatus` (TypeSet of Int)
  - `task_type_ids` → `TaskTypeIds` (TypeSet of Int)
  - `task_ids` → `TaskIds` (TypeSet of String)
  - `instance_ids` → `InstanceIds` (TypeSet of String)
  - `aliases` → `Aliases` (TypeSet of String)
  - `start_date` → `StartDate`
  - `end_date` → `EndDate`
  - `order_field` → `OrderField`
  - `order` → `Order` (Int, 0=升序, 1=降序)
  - `limit` → `Limit` (默认 20, 最大 100)
  - `offset` → `Offset` (默认 0)
  - `result_output_file` (标准字段,用于导出结果到 JSON 文件)

- 输出参数(Computed):
  - `repair_task_list` (TypeList): 维修任务列表,每个元素包含:
    - `task_id`, `instance_id`, `alias`, `task_type_id`, `task_status`
    - `create_time`, `auth_time`, `end_time`
    - `task_detail`, `task_type_name`, `task_sub_type`
    - `device_status`, `operate_status`, `auth_type`, `auth_source`
    - `zone`, `region`, `vpc_id`, `vpc_name`, `subnet_id`, `subnet_name`
    - `wan_ip`, `lan_ip`, `product`
  - `total_count` (TypeInt): 总数量

### 3. 分页处理
**决策**: 暴露 `limit` 和 `offset` 参数给用户,不实现自动分页

**理由**:
- 数据源应该提供查询灵活性,让用户自行控制分页
- 自动分页可能导致大量 API 调用,影响性能和费用
- 与现有数据源(如 `data_source_tc_instances.go`)保持一致的分页策略

### 4. 错误处理和重试
**决策**: 使用 `helper.Retry()` 包装 API 调用,处理最终一致性问题

**理由**:
- 遵循 Provider 的标准错误处理模式
- 云 API 可能出现暂时性故障,重试可以提高稳定性
- 参考 `defer tccommon.LogElapsed()` 和 `defer tccommon.InconsistentCheck()` 模式

### 5. 测试策略
**决策**: 提供基础验收测试,测试核心查询功能

**测试用例**:
- 测试不带任何过滤条件的查询(返回所有结果)
- 测试带任务状态过滤的查询
- 测试带实例ID过滤的查询
- 验证自动分页功能(当数据量超过100条时)

**备选方案**: 使用 Mock 数据而非真实 API 调用
- **不采用理由**: Terraform Provider 的验收测试标准是使用真实 API,需要设置 `TF_ACC=1` 环境变量

## Risks / Trade-offs

**[风险] API 返回大量数据导致 Terraform state 过大**
- → 缓解: 在文档中建议用户使用 `limit` 参数限制返回数量,避免一次查询过多数据
- → 缓解: 提供多种过滤条件,帮助用户精确查询所需数据

**[风险] API 字段变更导致兼容性问题**
- → 缓解: 所有字段设置为 Optional 或 Computed,即使 API 新增字段也不会破坏现有配置
- → 缓解: 使用 SDK 的结构体字段,通过 nil 检查防止空指针

**[权衡] 不实现自动分页 vs 用户体验**
- 优势: 避免不必要的 API 调用,降低成本和延迟
- 劣势: 用户需要手动处理分页逻辑
- 决策: 保持数据源的简单性,复杂的数据聚合应该在 Terraform 配置层面处理

**[权衡] 使用 TypeSet vs TypeList 存储数组参数**
- 决策: 对于 `task_ids`, `instance_ids` 等 ID 列表使用 TypeSet
- 理由: ID 列表不应该有重复值,TypeSet 自动去重;顺序无关紧要
- 决策: 对于输出的 `repair_task_list` 使用 TypeList
- 理由: 保持查询结果的顺序(根据 `order_field` 排序)

## Migration Plan

**部署步骤**:
1. 开发并通过本地测试
2. 提交 PR 到 terraform-provider-tencentcloud 仓库
3. 等待 CI/CD 验收测试通过
4. Merge 后随下一个版本发布

**回滚策略**:
- 纯新增功能,不影响现有资源,无需回滚
- 如果发现严重 bug,可以在下一个版本中修复

**文档更新**:
- 在 `website/docs/d/cvm_repair_tasks.html.markdown` 中提供详细的使用示例
- 在 changelog 中记录新增功能

## Open Questions

无待解决问题。API 文档完整,实现方案明确。
