## Why

需要为 TEO 服务添加创建 CLS 索引的 Terraform operation 资源，以便用户可以通过 Terraform 管理实时日志投递任务的 CLS 索引创建操作。

## What Changes

为云产品 teo 新增 Terraform operation 资源：

- 新增资源 `resource_tc_teo_create_cls_index_operation`
  - 资源类型为 operation（一次性操作）
  - 仅实现 Create 接口，Read、Update、Delete 接口为空
  - 调用云 API 接口 `CreateCLSIndex`
  - 支持参数：`zone_id`（站点 ID）、`task_id`（实时日志投递任务 ID）

## Capabilities

### New Capabilities
- `teo-create-cls-index-operation`: 为 TEO 服务提供创建 CLS 索引的 operation 资源能力

### Modified Capabilities
- （无）

## Impact

- 新增资源文件：`tencentcloud/services/teo/resource_tc_teo_create_cls_index_operation.go`
- 新增测试文件：`tencentcloud/services/teo/resource_tc_teo_create_cls_index_operation_test.go`
- 新增文档文件：`tencentcloud/services/teo/resource_tc_teo_create_cls_index_operation.md`
- 在 `tencentcloud/services/teo/service_tencentcloud_teo.go` 中注册新资源
- 依赖云 API：`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901` 的 `CreateCLSIndex` 接口
