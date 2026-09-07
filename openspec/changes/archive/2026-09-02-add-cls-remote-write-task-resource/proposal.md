## Why

CLS (Cloud Log Service) 提供 RemoteWrite 任务功能，允许用户将日志通过 Prometheus RemoteWrite 协议投递到外部目标服务。当前 Terraform Provider 尚未支持该资源的编排，用户只能通过控制台或 API 手动创建和管理 RemoteWrite 任务，无法实现基础设施即代码（IaC）管理。新增 `tencentcloud_cls_remote_write_task` 资源后，用户可通过 Terraform 声明式地管理 RemoteWrite 任务的完整生命周期（创建、读取、更新、删除），消除配置漂移并支持版本化管理。

## What Changes

- 新增 Terraform 资源 `tencentcloud_cls_remote_write_task`（RESOURCE_KIND_GENERAL），覆盖完整 CRUD 生命周期
  - **Create**: 调用 `CreateRemoteWriteTask` 创建任务，返回 `task_id`
  - **Read**: 调用 `DescribeRemoteWriteTasks` 按 `taskId` 过滤查询任务详情
  - **Update**: 调用 `ModifyRemoteWriteTask` 修改任务配置（名称、网络类型、鉴权信息、目标地址等）
  - **Delete**: 调用 `DeleteRemoteWriteTask` 删除任务
- 资源使用 `task_id` + `topic_id` 组合作为复合 ID（`task_id#topic_id`），因为 Delete 和 Modify 接口均要求同时传入这两个字段
- Schema 包含顶层平铺参数：`topic_id`、`name`、`target`、`remote_write_url`、`auth_type`、`net_type`、`vpc_id`、`virtual_gateway_type`、`enable` 以及嵌套的 `auth_info` 块（含 `username`、`password`、`token`）
- Read 接口额外返回只读字段：`status`、`create_time`、`update_time`、`logset_id`
- 在 `tencentcloud/provider.go` 和 `tencentcloud/provider.md` 中注册新资源
- 生成对应的资源文档 `resource_tc_cls_remote_write_task.md`

## Capabilities

### New Capabilities
- `cls-remote-write-task-resource`: 管理 CLS RemoteWrite 任务的完整生命周期，包括创建、查询、修改、删除 RemoteWrite 投递任务，支持配置目标地址、鉴权类型、网络类型等参数

### Modified Capabilities
<!-- 无需修改已有 capability -->

## Impact

- **新增文件**:
  - `tencentcloud/services/cls/resource_tc_cls_remote_write_task.go` — 资源实现（schema + CRUD）
  - `tencentcloud/services/cls/resource_tc_cls_remote_write_task.md` — 资源文档示例
  - `tencentcloud/services/cls/resource_tc_cls_remote_write_task_test.go` — 单元测试（使用 gomonkey mock 云 API）
- **修改文件**:
  - `tencentcloud/provider.go` — 注册 `tencentcloud_cls_remote_write_task` 资源
  - `tencentcloud/provider.md` — 新增资源文档条目
- **依赖的云 API**（已在 vendor 中可用）:
  - `CreateRemoteWriteTask`（cls/v20201016）
  - `DescribeRemoteWriteTasks`（cls/v20201016）
  - `ModifyRemoteWriteTask`（cls/v20201016）
  - `DeleteRemoteWriteTask`（cls/v20201016）
- **向后兼容**: 纯新增资源，不影响现有资源和配置
