## Why

用户需要在 Terraform 中查询云服务器的维修任务列表,以便在基础设施即代码的场景下监控和管理实例的维护状态。当前 Provider 缺少对应的数据源,用户无法通过 Terraform 获取维修任务信息,必须手动调用 API 或使用控制台查询,影响了自动化运维效率。

## What Changes

- 新增数据源 `tencentcloud_cvm_repair_tasks`,对应腾讯云 CVM 服务的 `DescribeTaskInfo` API
- 支持通过多种条件过滤查询维修任务:实例ID、实例名称、任务状态、任务类型、产品类型、时间范围等
- 支持分页查询和自定义排序
- 返回维修任务的完整信息:任务ID、实例信息、任务状态、创建/授权/结束时间、网络信息等

## Capabilities

### New Capabilities
- `cvm-repair-tasks-datasource`: 实现 CVM 维修任务数据源,支持多条件查询和过滤维修任务列表

### Modified Capabilities
<!-- 无现有能力需要修改 -->

## Impact

- **新增文件**:
  - `tencentcloud/services/cvm/data_source_tc_cvm_repair_tasks.go` (数据源实现)
  - `tencentcloud/services/cvm/data_source_tc_cvm_repair_tasks_test.go` (验收测试)
  - `website/docs/d/cvm_repair_tasks.html.markdown` (文档)
- **修改文件**:
  - `tencentcloud/provider.go` (注册新数据源)
- **依赖**: 使用现有的 `tencentcloud-sdk-go/tencentcloud/cvm` 包,无需新增依赖
- **兼容性**: 纯新增功能,不影响现有资源和数据源,完全向后兼容
