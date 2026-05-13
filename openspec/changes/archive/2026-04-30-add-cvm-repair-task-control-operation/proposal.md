## Why

腾讯云 CVM 在硬件维护或故障检测时会创建"维修任务"（Repair Task），处于"待授权"状态时需要用户主动授权才能进行后续维修操作（如重启、迁移等）。目前 Terraform Provider 已支持通过 `tencentcloud_cvm_repair_tasks` 数据源查询维修任务，但缺少对维修任务进行授权操作的能力。用户需要通过 Terraform 自动化对待授权的维修任务进行授权，以实现完整的运维闭环。

## What Changes

- 新增 `tencentcloud_cvm_repair_task_control` 操作型资源（operation resource），对应云 API `RepairTaskControl`
- 支持指定产品类型（CVM/CDH/CPM2.0）、实例 ID 列表、维修任务 ID 进行授权操作
- 支持立即授权与预约授权（`OrderAuthTime`）两种模式
- 支持本地盘弃盘迁移策略（`TaskSubMethod=LossyLocal`）
- 在 Provider 中注册新资源
- 提供完整的资源文档模板与生成的 website 文档
- 提供验收测试用例

## Capabilities

### New Capabilities

- `cvm-repair-task-control-operation`: 对腾讯云 CVM 待授权状态的维修任务进行授权控制操作的 Terraform 操作型资源

### Modified Capabilities

无

## Impact

- **新增代码**:
  - `tencentcloud/services/cvm/resource_tc_cvm_repair_task_control_operation.go`
  - `tencentcloud/services/cvm/resource_tc_cvm_repair_task_control_operation_test.go`
  - `tencentcloud/services/cvm/resource_tc_cvm_repair_task_control_operation.md`
- **修改代码**:
  - `tencentcloud/provider.go`: 注册 `tencentcloud_cvm_repair_task_control_operation`
  - `tencentcloud/provider.md`: 添加资源到列表
- **依赖**: 复用已有的 `tencentcloud-sdk-go/tencentcloud/cvm/v20170312`，无新增依赖
- **API 调用**: 调用 CVM `RepairTaskControl` 接口
- **API 限制**: 默认 20 次/秒，预约时间需在当前时间 5 分钟之后、48 小时之内
- **向后兼容**: 仅新增资源，不破坏任何现有 schema 或行为
