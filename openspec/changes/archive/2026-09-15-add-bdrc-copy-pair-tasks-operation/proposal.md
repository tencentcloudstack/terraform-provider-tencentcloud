## Why

腾讯云 BDRC（业务数据弹性容灾）提供了复制对（Copy Pair）的启动能力，运维人员需要在容灾演练或故障恢复流程中通过 `RunCopyPairTasks` 接口一次性启动一组复制对任务。当前 Terraform Provider 尚未覆盖 BDRC 产品的任何资源，用户无法以声明式方式触发该操作，必须借助控制台或 SDK 手动调用，导致自动化容灾编排流程断裂。本次新增基于 Terraform Plugin Framework 的 action 资源 `tencentcloud_bdrc_copy_pair_tasks`，将该一次性操作纳入 Terraform 编排体系。

## What Changes

- 新增 BDRC 产品的 SDK 客户端接入：在 `tencentcloud/connectivity/client.go` 中新增 `bdrcv20260330` 客户端字段及 `UseBdrcV20260330Client()` 访问方法（vendor 中已存在 `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/bdrc/v20260330`）。
- 新增 Terraform Plugin Framework action 资源 `tencentcloud_bdrc_copy_pair_tasks`，文件 `tencentcloud/services/bdrc/action_tc_bdrc_copy_pair_tasks.go`：
  - `copy_pair_ids`（Required, List of String）：复制对 ID 列表。
  - `copy_pair_type`（Required, String）：要启动复制对的类型，取值 DISK / INSTANCE / CFS。
  - 仅实现 `Invoke`（对应 Create）：调用 `RunCopyPairTasksWithContext`，操作完成后不记录任何云端状态。
- 在 `tencentcloud/framework/registry.go` 的 `actionFactories` 中注册 `bdrc.NewBdrcCopyPairTasks`。
- 新增单元测试 `tencentcloud/services/bdrc/action_tc_bdrc_copy_pair_tasks_test.go`，使用 gomonkey mock 云 API，仅做业务逻辑测试。
- 新增资源文档 `tencentcloud/services/bdrc/action_tc_bdrc_copy_pair_tasks.md`（一句话描述 + Example Usage；无 Import 段；action 资源无 Argument/Attribute Reference）。
- 无需修改 `tencentcloud/provider.go`（framework action 资源通过 `framework/registry.go` 统一注册，不在 SDKv2 ResourcesMap 中注册）。

## Capabilities

### New Capabilities

- `bdrc-copy-pair-tasks-action`: 提供一次性 action 资源 `tencentcloud_bdrc_copy_pair_tasks`，通过 `RunCopyPairTasks` API 启动一组 BDRC 复制对任务，操作完成后不持久化任何云端状态。

### Modified Capabilities

无。

## Impact

- 代码：
  - `tencentcloud/connectivity/client.go`（新增 bdrc 客户端 import、字段及 `UseBdrcV20260330Client()` 方法）
  - `tencentcloud/services/bdrc/action_tc_bdrc_copy_pair_tasks.go`（新增 — framework action 资源）
  - `tencentcloud/services/bdrc/service_tencentcloud_bdrc.go`（新增 — 服务层 RunCopyPairTasks 方法封装，含 retry）
  - `tencentcloud/services/bdrc/action_tc_bdrc_copy_pair_tasks_test.go`（新增 — 单元测试）
  - `tencentcloud/services/bdrc/action_tc_bdrc_copy_pair_tasks.md`（新增 — 文档）
  - `tencentcloud/framework/registry.go`（修改 — 注册 action factory）
- 依赖：使用已 vendored 的 `tencentcloud-sdk-go` 中 `bdrc/v20260330` 包（`RunCopyPairTasksRequest` / `RunCopyPairTasksResponse`），无需新增 vendor。
- 向后兼容：纯新增资源，不影响已有资源配置与 state。
- 文档：需通过收尾阶段 `make doc` 生成 `website/docs/` 下对应文档（禁止手改 website 目录）。
