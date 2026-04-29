## Why

腾讯云 EdgeOne (TEO) 多通道安全加速网关（Multi-Path Gateway）目前没有对应的 Terraform 资源，用户无法通过 Terraform 管理多通道安全加速网关的完整生命周期。需要新增 `tencentcloud_teo_multi_path_gateway` 资源，支持网关的创建、查询、修改和删除操作，实现基础设施即代码管理。

## What Changes

- 新增 Terraform 资源 `tencentcloud_teo_multi_path_gateway`，支持多通道安全加速网关的 CRUD 操作
  - Create: 调用 `CreateMultiPathGateway` 接口创建网关
  - Read: 调用 `DescribeMultiPathGateways` 接口查询网关详情
  - Update: 调用 `ModifyMultiPathGateway` 接口修改网关信息
  - Delete: 调用 `DeleteMultiPathGateway` 接口删除网关
- 在 `provider.go` 和 `provider.md` 中注册新资源
- 新增资源文档 `.md` 文件
- 新增单元测试文件，使用 gomonkey mock 方式测试业务逻辑

## Capabilities

### New Capabilities
- `teo-multi-path-gateway-resource`: 管理腾讯云 EdgeOne 多通道安全加速网关资源的完整生命周期，包括创建、读取、更新和删除

### Modified Capabilities
<!-- 无现有能力需要修改 -->

## Impact

- 新增文件: `tencentcloud/services/teo/resource_tc_teo_multi_path_gateway.go`
- 新增文件: `tencentcloud/services/teo/resource_tc_teo_multi_path_gateway_test.go`
- 新增文件: `tencentcloud/services/teo/resource_tc_teo_multi_path_gateway.md`
- 修改文件: `tencentcloud/provider.go` (注册新资源)
- 修改文件: `tencentcloud/provider.md` (添加资源文档链接)
- 依赖 SDK: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901`
