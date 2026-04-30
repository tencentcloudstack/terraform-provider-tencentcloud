## Why

TEO (EdgeOne) 多通道安全加速网关目前缺少 Terraform 资源支持，用户无法通过 Terraform 管理多通道安全加速网关的完整生命周期（创建、读取、更新、删除）。需要新增 `tencentcloud_teo_multi_path_gateway` 资源，使用户能够以基础设施即代码的方式管理 TEO 多通道安全加速网关。

## What Changes

- 新增 Terraform 资源 `tencentcloud_teo_multi_path_gateway`，支持多通道安全加速网关的完整 CRUD 操作：
  - **Create**: 调用 `CreateMultiPathGateway` 接口创建网关，返回 `gateway_id`
  - **Read**: 调用 `DescribeMultiPathGateways` 接口查询网关详情
  - **Update**: 调用 `ModifyMultiPathGateway` 接口修改网关名称、IP、端口
  - **Delete**: 调用 `DeleteMultiPathGateway` 接口删除网关
- 在 `provider.go` 和 `provider.md` 中注册新资源
- 新增资源文档 `resource_tc_teo_multi_path_gateway.md`
- 新增单元测试 `resource_tc_teo_multi_path_gateway_test.go`

## Capabilities

### New Capabilities
- `teo-multi-path-gateway-resource`: TEO 多通道安全加速网关 Terraform 资源，支持通过 CreateMultiPathGateway/DescribeMultiPathGateways/ModifyMultiPathGateway/DeleteMultiPathGateway 四个云 API 接口管理网关的完整生命周期

### Modified Capabilities
<!-- 无需修改现有 spec -->

## Impact

- **新增文件**:
  - `tencentcloud/services/teo/resource_tc_teo_multi_path_gateway.go`
  - `tencentcloud/services/teo/resource_tc_teo_multi_path_gateway_test.go`
  - `tencentcloud/services/teo/resource_tc_teo_multi_path_gateway.md`
- **修改文件**:
  - `tencentcloud/provider.go` - 注册新资源
  - `tencentcloud/provider.md` - 添加资源文档条目
- **依赖**: 使用现有 `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901` 包中的 CRUD 接口
