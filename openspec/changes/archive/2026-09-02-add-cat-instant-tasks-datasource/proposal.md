## Why

用户需要通过 Terraform 查询云拨测（CAT，Cat）的历史即时拨测任务列表。当前 Provider 已支持 cat 资源/数据源（`tencentcloud_cat_task_set`、`tencentcloud_cat_node`、`tencentcloud_cat_probe_data`、`tencentcloud_cat_metric_data`），但缺少针对"历史即时拨测任务"的查询数据源。

这会导致用户无法：
1. 通过 Terraform 查询已有的即时拨测任务信息（任务 ID、目标地址、任务类型、拨测时间、状态、成功率、节点数量、节点类型等）
2. 在 Terraform 配置中引用即时拨测任务的结果做下游编排或告警关联
3. 查看即时拨测任务的总量（Total）

## What Changes

- 新增 Data Source: `tencentcloud_cat_instant_tasks`
- 实现对 CAT API `DescribeInstantTasks` 接口的调用，获取历史即时拨测任务列表
- 支持内部自动分页获取所有数据（不向用户暴露 limit/offset 参数，遵循 openspec config.yaml 约束）
- 返回即时拨测任务列表及总量信息，字段映射如下：
  - `response.Response.Tasks` → `tasks`（列表，展开每个任务的字段到顶层 schema 元素中）
  - `response.Response.Tasks.TaskId` → `task_id`
  - `response.Response.Tasks.TargetAddress` → `target_address`
  - `response.Response.Tasks.TaskType` → `task_type`
  - `response.Response.Tasks.ProbeTime` → `probe_time`
  - `response.Response.Tasks.Status` → `status`
  - `response.Response.Tasks.SuccessRate` → `success_rate`
  - `response.Response.Tasks.NodeCount` → `node_count`
  - `response.Response.Tasks.TaskCategory` → `task_category`
  - `response.Response.Total` → `total`

## Capabilities

### New Capabilities
- `cat-instant-tasks-datasource`: 查询云拨测（CAT）历史即时拨测任务列表的 Data Source，封装 `DescribeInstantTasks` 接口，支持自动分页、结果导出与字段展开

### Modified Capabilities
<!-- 无现有 spec 需要修改 -->

## Impact

- **新增能力**: CAT 历史即时拨测任务列表查询
- **受影响的服务**: CAT (tencentcloud/services/cat)
- **新增文件**:
  - `tencentcloud/services/cat/data_source_tc_cat_instant_tasks.go`
  - `tencentcloud/services/cat/data_source_tc_cat_instant_tasks.md`
  - `tencentcloud/services/cat/data_source_tc_cat_instant_tasks_test.go`
- **修改文件**:
  - `tencentcloud/services/cat/service_tencentcloud_cat.go`（新增 `DescribeCatInstantTasksByFilter` service 方法，内部实现分页）
  - `tencentcloud/provider.go`（在 `DataSourcesMap` 注册 `tencentcloud_cat_instant_tasks`）
  - `tencentcloud/provider.md`（在 CAT 部分 Data Source 节添加条目，由 `make doc` 自动生成）
- **API 依赖**:
  - CAT API v20180409: `DescribeInstantTasks`（同步接口，无需轮询 Read）
  - SDK 包: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cat/v20180409`（已在 vendor 中）
- **兼容性**: 无破坏性变更，纯新增功能
