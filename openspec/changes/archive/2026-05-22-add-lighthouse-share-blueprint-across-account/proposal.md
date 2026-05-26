## Why

用户需要通过 Terraform 管理轻量应用服务器（Lighthouse）的跨账号共享镜像操作。目前该操作只能通过控制台或云 API 手动执行，无法通过 Terraform 进行自动化管理。新增此资源可以让用户在基础设施即代码的工作流中统一管理跨账号镜像共享。

## What Changes

- 新增 Terraform 资源 `tencentcloud_lighthouse_share_blueprint_across_account`，类型为 RESOURCE_KIND_OPERATION（一次性操作资源）
- 该资源调用 `ShareBlueprintAcrossAccounts` API 实现跨账号共享镜像功能
- 资源仅实现 Create 接口，Read/Update/Delete 接口为空（一次性操作，无需记录状态）
- 资源文件命名为 `resource_tc_lighthouse_share_blueprint_across_account_operation.go`

## Capabilities

### New Capabilities
- `lighthouse-share-blueprint-across-account-operation`: 跨账号共享镜像的一次性操作资源，支持指定镜像ID和接收账号ID列表

### Modified Capabilities

## Impact

- 新增文件: `tencentcloud/services/lighthouse/resource_tc_lighthouse_share_blueprint_across_account_operation.go`
- 修改文件: `tencentcloud/provider.go`（注册新资源）
- 修改文件: `tencentcloud/provider.md`（添加资源文档索引）
- 新增文件: `tencentcloud/services/lighthouse/resource_tc_lighthouse_share_blueprint_across_account_operation_test.go`（单元测试）
- 新增文件: `tencentcloud/services/lighthouse/resource_tc_lighthouse_share_blueprint_across_account_operation.md`（资源样例文档）
- 依赖云 API: `ShareBlueprintAcrossAccounts`（lighthouse v20200324）
