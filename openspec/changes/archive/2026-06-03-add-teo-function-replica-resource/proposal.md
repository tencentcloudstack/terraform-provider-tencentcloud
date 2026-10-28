## Why

TencentCloud EdgeOne (TEO) 支持边缘函数副本功能，允许用户为边缘函数创建副本以便进行版本管理和灰度测试。当前 Terraform Provider 缺少对该资源的管理能力，用户无法通过 IaC 方式管理边缘函数副本的生命周期。

## What Changes

- 新增 Terraform 资源 `tencentcloud_teo_function_replica`，支持边缘函数副本的完整 CRUD 生命周期管理
- 资源通过 `zone_id`、`function_id`、`replica_name` 联合标识唯一副本
- 支持创建副本（CreateFunctionReplica）、查询副本（DescribeFunctionReplicas）、修改副本内容和描述（ModifyFunctionReplica）、删除副本（DeleteFunctionReplica）
- 在 provider.go 和 provider.md 中注册新资源

## Capabilities

### New Capabilities
- `teo-function-replica`: 管理 TEO 边缘函数副本资源的 CRUD 操作，包括创建、读取、更新和删除边缘函数副本

### Modified Capabilities

## Impact

- 新增文件: `tencentcloud/services/teo/resource_tc_teo_function_replica.go`
- 新增文件: `tencentcloud/services/teo/resource_tc_teo_function_replica_test.go`
- 新增文件: `tencentcloud/services/teo/resource_tc_teo_function_replica.md`
- 修改文件: `tencentcloud/provider.go`（注册新资源）
- 修改文件: `tencentcloud/provider.md`（添加资源文档引用）
- 依赖: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901`（已在 vendor 中）
