## Why

需要为堡垒机（BH）产品提供 Terraform 资源来管理资产与堡垒机服务的绑定关系。用户需要通过 Terraform 将设备资产绑定到指定的堡垒机服务实例，支持 K8S 集群托管场景下的维度、账号、凭证等配置。

## What Changes

- 新增 Terraform 资源 `tencentcloud_bh_bind_device_resource`，类型为 RESOURCE_KIND_GENERAL
- 该资源使用 `BindDeviceResource` 接口完成 Create/Update/Delete 操作（绑定/更新绑定/解绑）
- 该资源使用 `DescribeDevices` 接口完成 Read 操作（查询设备绑定状态）
- 支持的参数包括：device_id_set、resource_id、domain_id、manage_dimension、manage_account_id、manage_account、manage_kubeconfig、namespace、workload
- 在 provider.go 和 provider.md 中注册新资源

## Capabilities

### New Capabilities
- `bh-bind-device-resource`: 管理堡垒机资产与服务实例的绑定关系，支持 K8S 集群托管场景

### Modified Capabilities

## Impact

- 新增文件：`tencentcloud/services/bh/resource_tc_bh_bind_device_resource.go`
- 新增文件：`tencentcloud/services/bh/resource_tc_bh_bind_device_resource_test.go`
- 新增文件：`tencentcloud/services/bh/resource_tc_bh_bind_device_resource.md`
- 修改文件：`tencentcloud/provider.go`（注册资源）
- 修改文件：`tencentcloud/provider.md`（添加资源文档引用）
- 依赖：`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/bh/v20230418`（已在 vendor 中）
